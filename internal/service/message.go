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

func (s *MessageService) ProcessUnsentMessages(ctx context.Context, limit int) error {
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

	var sentMessageIds []int
	for _, msg := range messages {
		err := s.sendMessage(msg)
		if err != nil {
			fmt.Printf("Failed to send message ID %d: %v\n", msg.ID, err)
			continue
		}

		fmt.Printf("Successfully sent message ID %d: Content='%s', Author='%s'\n",
			msg.ID, msg.Content, msg.Author)
		sentMessageIds = append(sentMessageIds, msg.ID)
	}

	if len(sentMessageIds) > 0 {
		err = s.repo.MarkMultipleAsSent(ctx, tx, sentMessageIds)
		if err != nil {
			return fmt.Errorf("failed to mark messages as sent: %w", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		fmt.Printf("Successfully processed %d messages\n", len(sentMessageIds))
	}

	return nil
}

func (s *MessageService) sendMessage(msg *models.Message) error {
	fmt.Printf("Sending message: ID=%d, Content='%s', Author='%s'\n",
		msg.ID, msg.Content, msg.Author)

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
