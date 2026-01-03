package domain

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FullName  string
	About     string
	City      string
	Skills    []string
	UpdatedAt time.Time
}
