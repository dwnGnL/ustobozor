package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Text      string
	PhotoPath string
	CreatedAt time.Time
	ReadAt    *time.Time
}
