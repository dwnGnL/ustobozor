package resolvers

import (
	"context"
	"errors"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/middleware"
	"github.com/barzurustami/bozor/internal/repository"
)

func resolveMe(ctx context.Context, r *Resolver) (*model.User, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, nil
	}

	user, err := r.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var profile *domain.Profile
	if r.ProfileRepo != nil {
		profile, err = r.ProfileRepo.GetByUserID(ctx, userID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				profile = nil
			} else {
				return nil, err
			}
		}
	}

	return toModelUser(user, profile), nil
}
