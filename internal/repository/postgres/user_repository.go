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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	const query = `
		SELECT id, phone, created_at
		FROM users
		WHERE phone = $1
	`

	var (
		id        uuid.UUID
		createdAt time.Time
	)

	err := r.pool.QueryRow(ctx, query, phone).Scan(&id, &phone, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Phone:     phone,
		CreatedAt: createdAt,
	}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `
		SELECT id, phone, created_at
		FROM users
		WHERE id = $1
	`

	var (
		phone     string
		createdAt time.Time
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(&id, &phone, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Phone:     phone,
		CreatedAt: createdAt,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (id, phone, created_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Phone, user.CreatedAt)
	return err
}
