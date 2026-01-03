package domain

import (
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID        uuid.UUID
	RequestID uuid.UUID
	Path      string
	CreatedAt time.Time
}
