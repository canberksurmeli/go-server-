package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Sent      bool      `json:"sent"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateMessageRequest struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type UpdateMessageRequest struct {
	Content string `json:"content,omitempty"`
	Author  string `json:"author,omitempty"`
}