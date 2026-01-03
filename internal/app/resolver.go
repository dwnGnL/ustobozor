package app

import "github.com/barzurustami/bozor/internal/graphql/resolvers"

func NewResolver(services *Services, repos *Repositories) *resolvers.Resolver {
	return &resolvers.Resolver{
		AuthService:    services.Auth,
		ProfileService: services.Profile,
		RequestService: services.Request,
		PhotoService:   services.Photo,
		ChatService:    services.Chat,
		UserRepo:       repos.Users,
		ProfileRepo:    repos.Profiles,
	}
}
