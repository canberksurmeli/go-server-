package models

import "time"

// Message represents a message entity
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Sent      bool      `json:"sent"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMessageRequest represents a request to create a new message
type CreateMessageRequest struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

// UpdateMessageRequest represents a request to update a message
type UpdateMessageRequest struct {
	Content string `json:"content,omitempty"`
	Author  string `json:"author,omitempty"`
}