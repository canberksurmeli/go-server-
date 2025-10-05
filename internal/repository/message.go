package repository

import (
	"context"
	"message-provider-go/internal/database"
	"message-provider-go/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type MessageRepository struct {
	db *database.DB
}

func NewMessageRepository(db *database.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) GetUnsentMessages(ctx context.Context, limit int) ([]*models.Message, error) {
	query := `
        SELECT id, content, author, sent, sent_at, created_at, updated_at 
        FROM messages 
        WHERE sent = false 
        ORDER BY created_at ASC 
        LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID,
			&message.Content,
			&message.Author,
			&message.Sent,
			&message.SentAt,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, rows.Err()
}

func (r *MessageRepository) MarkAsSent(ctx context.Context, tx pgx.Tx, id int) error {
	query := `
        UPDATE messages 
        SET sent = true, 
            sent_at = $1, 
            updated_at = $2 
        WHERE id = $3`

	now := time.Now()
	_, err := tx.Exec(ctx, query, now, now, id)
	return err
}

func (r *MessageRepository) MarkMultipleAsSent(ctx context.Context, tx pgx.Tx, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE messages 
        SET sent = true, 
            sent_at = $1, 
            updated_at = $2 
        WHERE id = ANY($3)`

	now := time.Now()
	_, err := tx.Exec(ctx, query, now, now, ids)
	return err
}

func (r *MessageRepository) GetTwoMessagesWithTx(ctx context.Context, tx pgx.Tx) ([]*models.Message, error) {
	query := `
        SELECT id, content, author, sent, sent_at, created_at, updated_at 
        FROM messages 
        ORDER BY created_at DESC 
        LIMIT 2`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID,
			&message.Content,
			&message.Author,
			&message.Sent,
			&message.SentAt,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, rows.Err()
}

func (r *MessageRepository) GetDB() *database.DB {
	return r.db
}

func (r *MessageRepository) GetSentMessages(ctx context.Context) ([]*models.Message, error) {
	query := `
        SELECT id, content, author, sent, sent_at, created_at, updated_at
        FROM messages
        WHERE sent = true
        ORDER BY sent_at DESC
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := &models.Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.Author,
			&msg.Sent,
			&msg.SentAt,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
