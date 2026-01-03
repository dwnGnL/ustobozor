package postgres

import (
	"context"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestRepository struct {
	pool *pgxpool.Pool
}

func NewRequestRepository(pool *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{pool: pool}
}

func (r *RequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.JobRequest, error) {
	const query = `
		SELECT id, customer_id, title, description, address, created_at
		FROM job_requests
		WHERE id = $1
	`

	req := domain.JobRequest{}
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&req.ID,
		&req.CustomerID,
		&req.Title,
		&req.Description,
		&req.Address,
		&req.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &req, nil
}

func (r *RequestRepository) Create(ctx context.Context, req *domain.JobRequest) error {
	const query = `
		INSERT INTO job_requests (id, customer_id, title, description, address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		req.ID,
		req.CustomerID,
		req.Title,
		req.Description,
		req.Address,
		req.CreatedAt,
	)
	return err
}
