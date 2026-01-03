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

type RequestService struct {
	requests repository.RequestRepository
}

func NewRequestService(requests repository.RequestRepository) *RequestService {
	return &RequestService{requests: requests}
}

func (s *RequestService) Create(ctx context.Context, customerID uuid.UUID, title, description, address string) (*domain.JobRequest, error) {
	request := &domain.JobRequest{
		ID:          uuid.New(),
		CustomerID:  customerID,
		Title:       title,
		Description: description,
		Address:     address,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.requests.Create(ctx, request); err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info("job request created", zap.String("request_id", request.ID.String()))
	return request, nil
}
