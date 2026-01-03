package resolvers

import (
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/barzurustami/bozor/internal/service"
)

type Resolver struct {
	AuthService    *service.AuthService
	ProfileService *service.ProfileService
	RequestService *service.RequestService
	PhotoService   *service.PhotoService
	ChatService    *service.ChatService
	UserRepo       repository.UserRepository
	ProfileRepo    repository.ProfileRepository
}
