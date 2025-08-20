package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/PythonHacker24/linux-acl-management-backend/api/routes"
	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/grpcpool"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/postgresql"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/redis"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/scheduler"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/scheduler/fcfs"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/transprocessor"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/utils"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)
	}
}

func exec() error {

	/* exec() wraps run() protecting it with user interrupts  */

	err := godotenv.Load()
	if err != nil {
		fmt.Print("No .env file found, continuing with system environment variables\n")
	}

	/* setting up cobra for cli interactions */
	var (
		configPath string
		rootCmd    = &cobra.Command{
			Use:   "laclm <command> <subcommand>",
			Short: "Backend server for linux acl management",
			Example: heredoc.Doc(`
				$ laclm --config /path/to/config.yaml
			`),
			Run: func(cmd *cobra.Command, args []string) {
				if configPath != "" {
					fmt.Printf("Using config file: %s\n\n", configPath)
				} else {
					fmt.Printf("No config file provided.\n\n")
				}
			},
		}
	)

	/* adding --config argument */
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")

	/* Execute the command */
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("arguments error: %s", err.Error())
		os.Exit(1)
	}

	/*
		load config file
		if there is an error in loading the config file, then it will exit with code 1
	*/
	if err := config.LoadConfig(configPath); err != nil {
		fmt.Printf("Configuration Error in %s: %s",
			configPath,
			err.Error(),
		)
		/* since the configuration is invalid, don't proceed */
		os.Exit(1)
	}

	/*
		true for production, false for development mode
		logger is only for http server and core components (after this step)
		using logger for cli issues doesn't make sense
	*/
	utils.InitLogger(!config.BackendConfig.AppInfo.DebugMode)

	/* zap.L() can be used all over the code for global level logging */
	zap.L().Info("Logger Initiated ...")

	/* calculate max procs accurately (runtime.GOMAXPROCS(0)) */
	if _, err := maxprocs.Set(); err != nil {
		zap.L().Error("automaxprocs: failed to set GOMAXPROCS",
			zap.Error(err),
		)
	}

	/* preparing graceful shutdown for CTRL+C and docker */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		cancel()
	}()

	return run(ctx)
}

func run(ctx context.Context) error {
	var (
		err error
		wg  sync.WaitGroup
	)

	/* create a context and waitgroup for the logging goroutine */
	logCtx, logCancel := context.WithCancel(context.Background())
	var logWg sync.WaitGroup

	defer func() {
		logCancel()
		logWg.Wait()
	}()

	/* create a error channel */
	errChLog := make(chan error, 1)

	/* create the client pool for daemons (via gRPC) */
	/* unsecure for now */

	/* attempting to keep connections alive all the time even with no activity */
	var kacp = keepalive.ClientParameters{
		/* send pings every 10 seconds if there is no activity */
		Time: 10 * time.Second,

		/* wait 2 second for ping ack before considering the connection dead */
		Timeout: 2 * time.Second,

		/* send pings even without active streams */
		PermitWithoutStream: true,
	}

	pool := grpcpool.NewClientPool(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
	)

	/*
		initializing scheduler
		scheduler uses context to quit - part of waitgroup
		propagates error through error channel
	*/
	errChShed := make(chan error, 1)

	logRedisClient, err := redis.NewRedisClient(
		config.BackendConfig.Database.TransactionLogRedis.Address,
		config.BackendConfig.Database.TransactionLogRedis.Password,
		config.BackendConfig.Database.TransactionLogRedis.DB,
	)
	if err != nil {
		zap.L().Fatal("Failed to connect to Redis", zap.Error(err))
	}

	pqDB := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.BackendConfig.Database.ArchivalPQ.User,
		config.BackendConfig.Database.ArchivalPQ.Password,
		config.BackendConfig.Database.ArchivalPQ.Host,
		config.BackendConfig.Database.ArchivalPQ.Port,
		config.BackendConfig.Database.ArchivalPQ.DBName,
		config.BackendConfig.Database.ArchivalPQ.SSLMode,
	)

	poolPQ, err := pgxpool.New(context.Background(), pqDB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	archivalPQ := postgresql.New(poolPQ)

	/* create a session manager */
	sessionManager := session.NewManager(logRedisClient, archivalPQ, errChLog)

	/* create a permissions processor */
	permProcessor := transprocessor.NewPermProcessor(pool, errChLog)

	/* start logging goroutine - should be last to exit */
	logWg.Add(1)
	go func(ctx context.Context) {
		defer logWg.Done()
		zap.L().Info("log error handler started")
		for {
			select {
			case err, ok := <-errChLog:
				if !ok {
					zap.L().Info("log error channel closed")
					return
				}
				if err != nil {
					zap.L().Error("log error occurred",
						zap.Error(err),
						zap.Time("timestamp", time.Now()),
					)
				}
			case <-ctx.Done():
				zap.L().Info("log error handler shutting down")
				return
			}
		}
	}(logCtx)

	/* currently FCFS scheduler */
	transSched := fcfs.NewFCFSScheduler(sessionManager, permProcessor)

	/* initialize the scheduler */
	scheduler.InitScheduler(ctx, transSched, &wg, errChShed)

	/* setting up http mux and routes */
	mux := http.NewServeMux()

	/* routes declared in /api/routes.go */
	routes.RegisterRoutes(mux, sessionManager)

	/* create a http server */
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			config.BackendConfig.Server.Host,
			config.BackendConfig.Server.Port,
		),
		Handler: mux,
	}

	/* starting http server as a goroutine */
	go func() {
		zap.L().Info("HTTP REST API server starting",
			zap.String("Host", config.BackendConfig.Server.Host),
			zap.Int("Port", config.BackendConfig.Server.Port),
		)
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			zap.L().Error("ListenAndServe error",
				zap.Error(err),
			)
		}
	}()

	/*
		whatever written here will be protected by graceful shutdowns
		all the functions called must be async here and ready for graceful shutdowns
	*/

	/*
		scheduler is a core feature of the application
		when an error occurs in the scheduler, the system needs to be shutdown
		since nothing can work without the scheduler
	*/
	select {
	case <-ctx.Done():
		zap.L().Info("Shutdown process initiated")
	case err = <-errChShed:

		/* context done can be called here (optional for now) */

		zap.L().Error("Fatal Error from scheduler",
			zap.Error(err),
		)
		return err
	}

	/*
		after this, exit signal is triggered
		following code must be executed to shutdown graceful shutdown
		call all the kill switches with context
	*/

	/* graceful shutdown of http server - 5 seconds for allowing completion current API requests */
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	/* initiate http server shutdown */
	if err = server.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("HTTP server shutdown error",
			zap.Error(err),
		)
	}

	zap.L().Info("HTTP server stopped")

	/* after the http server is stopped, rest of the components can be shutdown */

	usernames := sessionManager.GetAllUsernames()
	for _, username := range usernames {
		if err := sessionManager.ExpireSession(username); err != nil {
			zap.L().Warn("Failed to expire session during shutdown",
				zap.String("username: ", username),
			)
		}
		zap.L().Info("Session forced expired for: ",
			zap.String("username", username),
		)
	}

	wg.Wait()

	/* close connections with daemon */
	pool.CloseAll(errChLog)

	/* flush Redis data before closing */
	if err := logRedisClient.FlushAll(context.Background()); err != nil {
		zap.L().Error("Failed to flush Redis data during shutdown",
			zap.Error(err),
		)
	}

	/* close archival database connection */
	poolPQ.Close()

	zap.L().Info("All background processes closed gracefully")

	/* close the logging error channel and cancel logging context */
	close(errChLog)

	return err
}
