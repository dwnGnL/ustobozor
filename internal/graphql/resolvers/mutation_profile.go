package resolvers

import (
	"context"
	"fmt"

	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/middleware"
)

func resolveUpsertProfile(ctx context.Context, r *Resolver, input model.ProfileInput) (*model.Profile, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	about := ""
	if input.About != nil {
		about = *input.About
	}

	city := ""
	if input.City != nil {
		city = *input.City
	}

	profile, err := r.ProfileService.Upsert(ctx, userID, input.FullName, about, city, input.Skills)
	if err != nil {
		return nil, err
	}

	return toModelProfile(profile), nil
}
