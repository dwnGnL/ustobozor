package resolvers

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/middleware"
	"github.com/google/uuid"
)

func resolveUploadPhotos(ctx context.Context, r *Resolver, input model.UploadPhotosInput) ([]*model.Photo, error) {
	if _, ok := middleware.UserIDFromContext(ctx); !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	requestID, err := uuid.Parse(input.RequestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request id")
	}

	uploads := make([]graphql.Upload, 0, len(input.Files))
	for _, file := range input.Files {
		if file == nil {
			continue
		}
		uploads = append(uploads, *file)
	}

	photos, err := r.PhotoService.Upload(ctx, requestID, uploads)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Photo, 0, len(photos))
	for _, photo := range photos {
		result = append(result, toModelPhoto(photo))
	}

	return result, nil
}
