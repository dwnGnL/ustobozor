package app

import (
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/barzurustami/bozor/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Users    repository.UserRepository
	Profiles repository.ProfileRepository
	Requests repository.RequestRepository
	Photos   repository.PhotoRepository
	Chats    repository.ChatRepository
	Messages repository.MessageRepository
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Users:    postgres.NewUserRepository(pool),
		Profiles: postgres.NewProfileRepository(pool),
		Requests: postgres.NewRequestRepository(pool),
		Photos:   postgres.NewPhotoRepository(pool),
		Chats:    postgres.NewChatRepository(pool),
		Messages: postgres.NewMessageRepository(pool),
	}
}
