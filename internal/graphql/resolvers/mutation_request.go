package resolvers

import (
	"context"
	"fmt"

	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/middleware"
)

func resolveCreateRequest(ctx context.Context, r *Resolver, input model.CreateRequestInput) (*model.JobRequest, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	address := ""
	if input.Address != nil {
		address = *input.Address
	}

	request, err := r.RequestService.Create(ctx, userID, input.Title, input.Description, address)
	if err != nil {
		return nil, err
	}

	return toModelRequest(request, nil), nil
}
