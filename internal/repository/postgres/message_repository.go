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

type MessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{pool: pool}
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error) {
	const query = `
		SELECT id, chat_id, sender_id, text, photo_path, created_at, read_at
		FROM chat_messages
		WHERE id = $1
	`

	msg := domain.ChatMessage{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&msg.ID,
		&msg.ChatID,
		&msg.SenderID,
		&msg.Text,
		&msg.PhotoPath,
		&msg.CreatedAt,
		&msg.ReadAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &msg, nil
}

func (r *MessageRepository) ListByChat(ctx context.Context, chatID uuid.UUID, limit, offset int32) ([]domain.ChatMessage, error) {
	const query = `
		SELECT id, chat_id, sender_id, text, photo_path, created_at, read_at
		FROM chat_messages
		WHERE chat_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		msg := domain.ChatMessage{}
		if err := rows.Scan(
			&msg.ID,
			&msg.ChatID,
			&msg.SenderID,
			&msg.Text,
			&msg.PhotoPath,
			&msg.CreatedAt,
			&msg.ReadAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return messages, nil
}

func (r *MessageRepository) Create(ctx context.Context, message *domain.ChatMessage) error {
	const query = `
		INSERT INTO chat_messages (id, chat_id, sender_id, text, photo_path, created_at, read_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		message.ID,
		message.ChatID,
		message.SenderID,
		message.Text,
		message.PhotoPath,
		message.CreatedAt,
		message.ReadAt,
	)
	return err
}

func (r *MessageRepository) MarkReadByID(ctx context.Context, messageID uuid.UUID, at time.Time) (*domain.ChatMessage, error) {
	const query = `
		UPDATE chat_messages
		SET read_at = $2
		WHERE id = $1 AND read_at IS NULL
		RETURNING id, chat_id, sender_id, text, photo_path, created_at, read_at
	`

	msg := domain.ChatMessage{}
	err := r.pool.QueryRow(ctx, query, messageID, at).Scan(
		&msg.ID,
		&msg.ChatID,
		&msg.SenderID,
		&msg.Text,
		&msg.PhotoPath,
		&msg.CreatedAt,
		&msg.ReadAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &msg, nil
}

func (r *MessageRepository) MarkReadByChat(ctx context.Context, chatID, readerID uuid.UUID, at time.Time) ([]domain.ChatMessage, error) {
	const query = `
		UPDATE chat_messages
		SET read_at = $3
		WHERE chat_id = $1
			AND sender_id <> $2
			AND read_at IS NULL
		RETURNING id, chat_id, sender_id, text, photo_path, created_at, read_at
	`

	rows, err := r.pool.Query(ctx, query, chatID, readerID, at)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		msg := domain.ChatMessage{}
		if err := rows.Scan(
			&msg.ID,
			&msg.ChatID,
			&msg.SenderID,
			&msg.Text,
			&msg.PhotoPath,
			&msg.CreatedAt,
			&msg.ReadAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return messages, nil
}
