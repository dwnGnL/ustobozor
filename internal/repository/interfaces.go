package repository

import (
	"context"
	"time"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type ProfileRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error)
	Upsert(ctx context.Context, profile *domain.Profile) error
}

type RequestRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.JobRequest, error)
	Create(ctx context.Context, req *domain.JobRequest) error
}

type PhotoRepository interface {
	CreateMany(ctx context.Context, photos []domain.Photo) error
}

type ChatRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
	GetByRequestAndInitiator(ctx context.Context, requestID, initiatorID uuid.UUID) (*domain.Chat, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error)
	Create(ctx context.Context, chat *domain.Chat) error
	UpdateLastMessageAt(ctx context.Context, chatID uuid.UUID, at time.Time) error
}

type MessageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error)
	ListByChat(ctx context.Context, chatID uuid.UUID, limit, offset int32) ([]domain.ChatMessage, error)
	Create(ctx context.Context, message *domain.ChatMessage) error
	MarkReadByID(ctx context.Context, messageID uuid.UUID, at time.Time) (*domain.ChatMessage, error)
	MarkReadByChat(ctx context.Context, chatID, readerID uuid.UUID, at time.Time) ([]domain.ChatMessage, error)
}
