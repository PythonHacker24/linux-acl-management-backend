package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	/*
		load config file
		if there is an error in loading the config file, then it will exit with code 1
	*/
	config.LoadConfig("./config.yaml")

	/* true for production, false for development mode */
	utils.InitLogger(false)

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

	/* setting up http mux and routes */
	mux := http.NewServeMux()

	/* routes declared in /api/routes */
	routes.RegisterRoutes(mux)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s",
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

	return err
}
