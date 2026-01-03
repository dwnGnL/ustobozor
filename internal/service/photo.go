package service

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/logger"
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/barzurustami/bozor/internal/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PhotoService struct {
	storage  *storage.LocalStorage
	photos   repository.PhotoRepository
	requests repository.RequestRepository
}

func NewPhotoService(storage *storage.LocalStorage, photos repository.PhotoRepository, requests repository.RequestRepository) *PhotoService {
	return &PhotoService{storage: storage, photos: photos, requests: requests}
}

func (s *PhotoService) Upload(ctx context.Context, requestID uuid.UUID, uploads []graphql.Upload) ([]domain.Photo, error) {
	if _, err := s.requests.GetByID(ctx, requestID); err != nil {
		return nil, err
	}

	stored := make([]domain.Photo, 0, len(uploads))
	for _, upload := range uploads {
		path, err := s.storage.Save(ctx, upload)
		if err != nil {
			return nil, err
		}

		stored = append(stored, domain.Photo{
			ID:        uuid.New(),
			RequestID: requestID,
			Path:      path,
			CreatedAt: time.Now().UTC(),
		})
	}

	if err := s.photos.CreateMany(ctx, stored); err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info("photos uploaded", zap.String("request_id", requestID.String()), zap.Int("count", len(stored)))
	return stored, nil
}
