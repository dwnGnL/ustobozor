package postgres

import (
	"context"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	const query = `
		SELECT id, user_id, full_name, about, city, skills, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	profile := domain.Profile{}
	var skills []string

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.About,
		&profile.City,
		&skills,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	profile.Skills = skills
	return &profile, nil
}

func (r *ProfileRepository) Upsert(ctx context.Context, profile *domain.Profile) error {
	const query = `
		INSERT INTO profiles (id, user_id, full_name, about, city, skills, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id)
		DO UPDATE SET
			full_name = EXCLUDED.full_name,
			about = EXCLUDED.about,
			city = EXCLUDED.city,
			skills = EXCLUDED.skills,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(ctx, query,
		profile.ID,
		profile.UserID,
		profile.FullName,
		profile.About,
		profile.City,
		profile.Skills,
		profile.UpdatedAt,
	)
	return err
}
