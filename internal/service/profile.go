package service

import (
	"context"
	"time"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/logger"
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProfileService struct {
	profiles repository.ProfileRepository
}

func NewProfileService(profiles repository.ProfileRepository) *ProfileService {
	return &ProfileService{profiles: profiles}
}

func (s *ProfileService) Upsert(ctx context.Context, userID uuid.UUID, fullName, about, city string, skills []string) (*domain.Profile, error) {
	if skills == nil {
		skills = []string{}
	}

	profile := &domain.Profile{
		ID:        uuid.New(),
		UserID:    userID,
		FullName:  fullName,
		About:     about,
		City:      city,
		Skills:    skills,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.profiles.Upsert(ctx, profile); err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info("profile upserted", zap.String("user_id", userID.String()))
	return profile, nil
}

func (s *ProfileService) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	return s.profiles.GetByUserID(ctx, userID)
}
