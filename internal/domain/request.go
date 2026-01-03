package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobRequest struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Title       string
	Description string
	Address     string
	CreatedAt   time.Time
}
