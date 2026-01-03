package postgres

import (
	"context"
	"time"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{pool: pool}
}

func (r *ChatRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	const query = `
		SELECT id, request_id, creator_id, initiator_id, created_at, last_message_at
		FROM chats
		WHERE id = $1
	`

	chat := domain.Chat{}
	var lastMessageAt *time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&chat.ID,
		&chat.RequestID,
		&chat.CreatorID,
		&chat.InitiatorID,
		&chat.CreatedAt,
		&lastMessageAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	chat.LastMessageAt = lastMessageAt
	return &chat, nil
}

func (r *ChatRepository) GetByRequestAndInitiator(ctx context.Context, requestID, initiatorID uuid.UUID) (*domain.Chat, error) {
	const query = `
		SELECT id, request_id, creator_id, initiator_id, created_at, last_message_at
		FROM chats
		WHERE request_id = $1 AND initiator_id = $2
	`

	chat := domain.Chat{}
	var lastMessageAt *time.Time

	err := r.pool.QueryRow(ctx, query, requestID, initiatorID).Scan(
		&chat.ID,
		&chat.RequestID,
		&chat.CreatorID,
		&chat.InitiatorID,
		&chat.CreatedAt,
		&lastMessageAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	chat.LastMessageAt = lastMessageAt
	return &chat, nil
}

func (r *ChatRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	const query = `
		SELECT id, request_id, creator_id, initiator_id, created_at, last_message_at
		FROM chats
		WHERE creator_id = $1 OR initiator_id = $1
		ORDER BY COALESCE(last_message_at, created_at) DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []domain.Chat
	for rows.Next() {
		chat := domain.Chat{}
		var lastMessageAt *time.Time

		if err := rows.Scan(
			&chat.ID,
			&chat.RequestID,
			&chat.CreatorID,
			&chat.InitiatorID,
			&chat.CreatedAt,
			&lastMessageAt,
		); err != nil {
			return nil, err
		}
		chat.LastMessageAt = lastMessageAt
		chats = append(chats, chat)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return chats, nil
}

func (r *ChatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	const query = `
		INSERT INTO chats (id, request_id, creator_id, initiator_id, created_at, last_message_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		chat.ID,
		chat.RequestID,
		chat.CreatorID,
		chat.InitiatorID,
		chat.CreatedAt,
		chat.LastMessageAt,
	)
	return err
}

func (r *ChatRepository) UpdateLastMessageAt(ctx context.Context, chatID uuid.UUID, at time.Time) error {
	const query = `
		UPDATE chats
		SET last_message_at = $2
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, chatID, at)
	return err
}
