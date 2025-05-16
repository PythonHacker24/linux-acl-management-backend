package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/api/routes"
	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/utils"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)
	}
}

func exec() error {

	/* exec() wraps run() protecting it with user interrupts  */

	/* setting up cobra for cli interactions */
	var(
		configPath string
		rootCmd = &cobra.Command{
			Use:   "laclm <command> <subcommand>",
			Short: "Backend server for linux acl management",
			Example: heredoc.Doc(`
				$ laclm
				$ laclm --config /path/to/config.yaml
			`),
			Run: func(cmd *cobra.Command, args []string) {
				if configPath != "" {
					fmt.Printf("Using config file: %s\n\n", configPath)
				} else {
					fmt.Println("No config file provided.\n\n")
				}
			},
		}
	)

	/* adding --config arguement */
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")

	/* Execute the command */
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("arguements error: %s", err.Error())
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
		load environment variables
		if there is an error or environment variables are not set, then it will exit with code 1
	*/
	config.LoadEnv()

	fmt.Println("loaded config")
	
	/* 
		true for production, false for development mode 
		logger is only for http server and core components (after this step)
		using logger for cli issues doesn't make sense
	*/
	utils.InitLogger(!config.BackendConfig.AppInfo.DebugMode)

	/* zap.L() can be used all over the code for global level logging */
	zap.L().Info("Logger Initiated ...")

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
	var err error

	/* complete backend system must initiate before http server starts */

	/* setting up http mux and routes */
	mux := http.NewServeMux()

	/* routes declared in /api/routes */
	routes.RegisterRoutes(mux)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			config.BackendConfig.Server.Host,
			config.BackendConfig.Server.Port,
		),
		Handler: mux,
	}

	/* starting http server as a goroutine */
	go func() {
		zap.L().Info("HTTP REST API server starting on :8080")
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

	<-ctx.Done()

	/*
		after this, exit signal is triggered
		following code must be executed to shutdown graceful shutdown
		call all the kill switches with context
	*/

	/* graceful shutdown of http server */
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

	return err
}
