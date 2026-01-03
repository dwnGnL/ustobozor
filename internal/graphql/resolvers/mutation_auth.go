package resolvers

import (
	"context"
	"errors"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/repository"
)

func resolveRequestSMSCode(ctx context.Context, r *Resolver, phone string) (bool, error) {
	if err := r.AuthService.RequestCode(ctx, phone); err != nil {
		return false, err
	}
	return true, nil
}

func resolveRegister(ctx context.Context, r *Resolver, input model.RegisterInput) (*model.AuthPayload, error) {
	user, tokens, err := r.AuthService.Register(ctx, input.Phone, input.Code)
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		User:   toModelUser(user, nil),
		Tokens: toModelTokenPair(tokens),
	}, nil
}

func resolveLogin(ctx context.Context, r *Resolver, input model.LoginInput) (*model.AuthPayload, error) {
	user, tokens, err := r.AuthService.Login(ctx, input.Phone, input.Code)
	if err != nil {
		return nil, err
	}

	var profile *domain.Profile
	if r.ProfileRepo != nil {
		profileDomain, err := r.ProfileRepo.GetByUserID(ctx, user.ID)
		if err == nil {
			profile = profileDomain
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	return &model.AuthPayload{
		User:   toModelUser(user, profile),
		Tokens: toModelTokenPair(tokens),
	}, nil
}

func resolveRefreshToken(ctx context.Context, r *Resolver, refreshToken string) (*model.AuthPayload, error) {
	user, tokens, err := r.AuthService.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	var profile *domain.Profile
	if r.ProfileRepo != nil {
		profileDomain, err := r.ProfileRepo.GetByUserID(ctx, user.ID)
		if err == nil {
			profile = profileDomain
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	return &model.AuthPayload{
		User:   toModelUser(user, profile),
		Tokens: toModelTokenPair(tokens),
	}, nil
}
