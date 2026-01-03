package app

import (
	"github.com/barzurustami/bozor/internal/auth"
	"github.com/barzurustami/bozor/internal/config"
	"github.com/barzurustami/bozor/internal/service"
	"github.com/barzurustami/bozor/internal/sms"
	"github.com/barzurustami/bozor/internal/storage"
	"go.uber.org/zap"
)

type Services struct {
	Auth    *service.AuthService
	Profile *service.ProfileService
	Request *service.RequestService
	Photo   *service.PhotoService
	Chat    *service.ChatService
	JWT     *auth.JWTService
}

func NewServices(cfg *config.Config, repos *Repositories, log *zap.Logger) *Services {
	jwtSvc := auth.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	smsSender := sms.NewMockSender()
	if cfg.SMS.Provider != "mock" {
		log.Warn("sms provider not configured, falling back to mock", zap.String("provider", cfg.SMS.Provider))
	}

	storageSvc := storage.NewLocalStorage(cfg.Upload.Dir, cfg.Upload.MaxSizeBytes)

	return &Services{
		JWT:     jwtSvc,
		Auth:    service.NewAuthService(repos.Users, smsSender, jwtSvc),
		Profile: service.NewProfileService(repos.Profiles),
		Request: service.NewRequestService(repos.Requests),
		Photo:   service.NewPhotoService(storageSvc, repos.Photos, repos.Requests),
		Chat:    service.NewChatService(repos.Chats, repos.Messages, repos.Requests, storageSvc),
	}
}
