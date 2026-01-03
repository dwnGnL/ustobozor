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
		SELECT id, phone, roles, created_at
		FROM users
		WHERE phone = $1
	`

	var (
		id        uuid.UUID
		roles     []string
		createdAt time.Time
	)

	err := r.pool.QueryRow(ctx, query, phone).Scan(&id, &phone, &roles, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Phone:     phone,
		Roles:     toRoles(roles),
		CreatedAt: createdAt,
	}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `
		SELECT id, phone, roles, created_at
		FROM users
		WHERE id = $1
	`

	var (
		phone     string
		roles     []string
		createdAt time.Time
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(&id, &phone, &roles, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:        id,
		Phone:     phone,
		Roles:     toRoles(roles),
		CreatedAt: createdAt,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (id, phone, roles, created_at)
		VALUES ($1, $2, $3, $4)
	`
	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, string(role))
	}

	_, err := r.pool.Exec(ctx, query, user.ID, user.Phone, roles, user.CreatedAt)
	return err
}

func toRoles(values []string) []domain.Role {
	roles := make([]domain.Role, 0, len(values))
	for _, role := range values {
		roles = append(roles, domain.Role(role))
	}
	return roles
}
