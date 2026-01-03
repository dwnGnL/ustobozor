package domain

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID            uuid.UUID
	RequestID     uuid.UUID
	CreatorID     uuid.UUID
	InitiatorID   uuid.UUID
	CreatedAt     time.Time
	LastMessageAt *time.Time
}
