package service

import (
	"context"
	"fmt"
	"message-provider-go/internal/models"
	"message-provider-go/internal/repository"
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

// ProcessUnsentMessages fetches unsent messages, sends them, and marks them as sent
func (s *MessageService) ProcessUnsentMessages(ctx context.Context, limit int) error {
	// Get unsent messages
	messages, err := s.repo.GetUnsentMessages(ctx, limit)
	if err != nil {
		return fmt.Errorf("failed to get unsent messages: %w", err)
	}

	if len(messages) == 0 {
		fmt.Println("No unsent messages to process")
		return nil
	}

	tx, err := s.repo.GetDB().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var sentMessageIDs []int
	for _, msg := range messages {
		// Simulate sending the message (replace this with actual sending logic)
		err := s.sendMessage(msg)
		if err != nil {
			fmt.Printf("Failed to send message ID %d: %v\n", msg.ID, err)
			continue
		}

		fmt.Printf("Successfully sent message ID %d: Content='%s', Author='%s'\n",
			msg.ID, msg.Content, msg.Author)
		sentMessageIDs = append(sentMessageIDs, msg.ID)
	}

	// Mark successfully sent messages as sent in database
	if len(sentMessageIDs) > 0 {
		err = s.repo.MarkMultipleAsSent(ctx, tx, sentMessageIDs)
		if err != nil {
			return fmt.Errorf("failed to mark messages as sent: %w", err)
		}

		// Commit transaction
		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		fmt.Printf("Successfully processed %d messages\n", len(sentMessageIDs))
	}

	return nil
}

// sendMessage simulates sending a message to an external service
// Replace this with your actual message sending logic (HTTP request, queue, etc.)
func (s *MessageService) sendMessage(msg *models.Message) error {
	// TODO: Implement actual message sending logic here
	// For example:
	// - Send HTTP request to external API
	// - Publish to message queue
	// - Send email/SMS
	// etc.

	// For now, just simulate success
	fmt.Printf("📤 Sending message: ID=%d, Content='%s', Author='%s'\n",
		msg.ID, msg.Content, msg.Author)

	// Simulate some processing
	// time.Sleep(100 * time.Millisecond)

	return nil
}

func (s *MessageService) GetSentMessages(ctx context.Context) ([]*models.Message, error) {
	messages, err := s.repo.GetSentMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent messages: %w", err)
	}

	fmt.Printf("Retrieved %d sent messages\n", len(messages))
	return messages, nil
}
