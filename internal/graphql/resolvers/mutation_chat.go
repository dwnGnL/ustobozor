package resolvers

import (
	"context"
	"fmt"

	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/middleware"
	"github.com/google/uuid"
)

func resolveCreateChat(ctx context.Context, r *Resolver, requestID string) (*model.Chat, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	parsedID, err := uuid.Parse(requestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request id")
	}

	chat, err := r.ChatService.CreateChat(ctx, parsedID, userID)
	if err != nil {
		return nil, err
	}

	return toModelChat(*chat), nil
}

func resolveSendMessage(ctx context.Context, r *Resolver, input model.SendMessageInput) (*model.ChatMessage, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	chatID, err := uuid.Parse(input.ChatID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id")
	}

	text := ""
	if input.Text != nil {
		text = *input.Text
	}

	message, err := r.ChatService.SendMessage(ctx, chatID, userID, text, input.File)
	if err != nil {
		return nil, err
	}

	return toModelChatMessage(*message), nil
}

func resolveMarkChatRead(ctx context.Context, r *Resolver, chatID string) ([]*model.ChatMessage, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	parsedID, err := uuid.Parse(chatID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id")
	}

	messages, err := r.ChatService.MarkChatRead(ctx, parsedID, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.ChatMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, toModelChatMessage(message))
	}

	return result, nil
}

func resolveMarkMessageRead(ctx context.Context, r *Resolver, messageID string) (*model.ChatMessage, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	parsedID, err := uuid.Parse(messageID)
	if err != nil {
		return nil, fmt.Errorf("invalid message id")
	}

	message, err := r.ChatService.MarkMessageRead(ctx, parsedID, userID)
	if err != nil {
		return nil, err
	}

	return toModelChatMessage(*message), nil
}
