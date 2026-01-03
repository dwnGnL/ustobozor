package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleCustomer Role = "CUSTOMER"
	RoleWorker   Role = "WORKER"
)

type User struct {
	ID        uuid.UUID
	Phone     string
	Roles     []Role
	CreatedAt time.Time
}
