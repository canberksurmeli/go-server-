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
	fmt.Println("Setting up scheduled jobs...")
	err := sched.AddJob("fetch-messages", 10*time.Second, FetchMessagesJob)
	if err != nil {
		fmt.Printf("Failed to add fetch messages job: %v\n", err)
	}
	sched.Start()
}

func FetchMessagesJob(ctx context.Context) error {
	db := database.Get()

	repo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(repo)

	err := messageService.ProcessUnsentMessages(ctx, 2)
	if err != nil {
		fmt.Printf("Error processing unsent messages: %v\n", err)
		return err
	}

	return nil
}
