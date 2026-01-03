package postgres

import (
	"context"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotoRepository struct {
	pool *pgxpool.Pool
}

func NewPhotoRepository(pool *pgxpool.Pool) *PhotoRepository {
	return &PhotoRepository{pool: pool}
}

func (r *PhotoRepository) CreateMany(ctx context.Context, photos []domain.Photo) error {
	if len(photos) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	const query = `
		INSERT INTO photos (id, request_id, path, created_at)
		VALUES ($1, $2, $3, $4)
	`

	for _, photo := range photos {
		batch.Queue(query, photo.ID, photo.RequestID, photo.Path, photo.CreatedAt)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range photos {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}
