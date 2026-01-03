package resolvers

import (
	"time"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/graphql/model"
)

func toModelUser(user *domain.User, profile *domain.Profile) *model.User {
	if user == nil {
		return nil
	}

	roles := make([]model.Role, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, model.Role(role))
	}

	return &model.User{
		ID:      user.ID.String(),
		Phone:   user.Phone,
		Roles:   roles,
		Profile: toModelProfile(profile),
	}
}

func toModelProfile(profile *domain.Profile) *model.Profile {
	if profile == nil {
		return nil
	}

	return &model.Profile{
		ID:       profile.ID.String(),
		FullName: profile.FullName,
		About:    stringPtr(profile.About),
		City:     stringPtr(profile.City),
		Skills:   profile.Skills,
	}
}

func toModelRequest(req *domain.JobRequest, photos []domain.Photo) *model.JobRequest {
	if req == nil {
		return nil
	}

	photoModels := make([]*model.Photo, 0, len(photos))
	for _, photo := range photos {
		photoModels = append(photoModels, toModelPhoto(photo))
	}

	return &model.JobRequest{
		ID:          req.ID.String(),
		Title:       req.Title,
		Description: req.Description,
		Address:     stringPtr(req.Address),
		CreatedAt:   model.Time(req.CreatedAt),
		Photos:      photoModels,
	}
}

func toModelPhoto(photo domain.Photo) *model.Photo {
	return &model.Photo{
		ID:        photo.ID.String(),
		Path:      photo.Path,
		CreatedAt: model.Time(photo.CreatedAt),
	}
}

func toModelTokenPair(tokens domain.TokenPair) *model.TokenPair {
	return &model.TokenPair{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  model.Time(tokens.AccessExpiresAt),
		RefreshExpiresAt: model.Time(tokens.RefreshExpiresAt),
	}
}

func toModelChat(chat domain.Chat) *model.Chat {
	lastMessageAt := timePtr(chat.LastMessageAt)
	return &model.Chat{
		ID:            chat.ID.String(),
		RequestID:     chat.RequestID.String(),
		CreatorID:     chat.CreatorID.String(),
		InitiatorID:   chat.InitiatorID.String(),
		CreatedAt:     model.Time(chat.CreatedAt),
		LastMessageAt: lastMessageAt,
	}
}

func toModelChatMessage(message domain.ChatMessage) *model.ChatMessage {
	return &model.ChatMessage{
		ID:        message.ID.String(),
		ChatID:    message.ChatID.String(),
		SenderID:  message.SenderID.String(),
		Text:      stringPtr(message.Text),
		Photo:     stringPtr(message.PhotoPath),
		CreatedAt: model.Time(message.CreatedAt),
		ReadAt:    timePtr(message.ReadAt),
	}
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func timePtr(value *time.Time) *model.Time {
	if value == nil {
		return nil
	}
	converted := model.Time(*value)
	return &converted
}
