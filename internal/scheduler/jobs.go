package scheduler

import (
	"context"
	"fmt"
	"message-provider-go/internal/database"
	"message-provider-go/internal/repository"
	"message-provider-go/internal/service"
	"time"
)

func SetupScheduledJobs(sched *Scheduler) {
	// Add jobs to the scheduler
	fmt.Println("Setting up scheduled jobs...")
	err := sched.AddJob("fetch-messages", 10*time.Second, FetchMessagesJob)
	if err != nil {
		fmt.Printf("Failed to add fetch messages job: %v\n", err)
	}
	sched.Start()
}

// FetchMessagesJob fetches 2 unsent messages, sends them, and updates their status in the database
func FetchMessagesJob(ctx context.Context) error {
	// Get database connection
	db := database.Get()

	// Initialize repository and service
	repo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(repo)

	// Process 2 unsent messages
	err := messageService.ProcessUnsentMessages(ctx, 2)
	if err != nil {
		fmt.Printf("Error processing unsent messages: %v\n", err)
		return err
	}

	return nil
}
