package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)
	}
}

func exec() error {

	/*
		exec() wraps run() protecting it with user interrupts  
	*/

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

	return nil
}
