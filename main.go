package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"message-provider-go/internal/config"
	"message-provider-go/internal/database"
	"message-provider-go/internal/scheduler"
	"message-provider-go/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Get()
	fmt.Println("Configuration loaded:", *cfg)

	// Initialize database connection
	err := database.Init()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
	}

	fmt.Println("Database connection established")

	sched := scheduler.New()
	scheduler.SetupScheduledJobs(sched)
	fmt.Println("Scheduler started")

	srv := server.New(sched)
	go srv.Start()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop scheduler
	sched.Stop()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exited")
}
