package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

type ChatSession struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatMessage struct {
	ID        int         `json:"id" db:"id"`
	SessionID uuid.UUID   `json:"session_id" db:"session_id"`
	Role      MessageRole `json:"role" db:"role"`
	Content   string      `json:"content" db:"content"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}
