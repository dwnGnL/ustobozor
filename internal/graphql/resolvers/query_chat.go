package resolvers

import (
	"context"
	"fmt"

	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/middleware"
	"github.com/google/uuid"
)

func resolveChats(ctx context.Context, r *Resolver) ([]*model.Chat, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	chats, err := r.ChatService.ListChats(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Chat, 0, len(chats))
	for _, chat := range chats {
		result = append(result, toModelChat(chat))
	}

	return result, nil
}

func resolveChatMessages(ctx context.Context, r *Resolver, chatID string, limit, offset *int) ([]*model.ChatMessage, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	parsedID, err := uuid.Parse(chatID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id")
	}

	limitVal := 50
	offsetVal := 0
	if limit != nil {
		limitVal = *limit
	}
	if offset != nil {
		offsetVal = *offset
	}
	if limitVal <= 0 {
		limitVal = 50
	}
	if limitVal > 200 {
		limitVal = 200
	}
	if offsetVal < 0 {
		offsetVal = 0
	}

	messages, err := r.ChatService.ListMessages(ctx, parsedID, userID, int32(limitVal), int32(offsetVal))
	if err != nil {
		return nil, err
	}

	result := make([]*model.ChatMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, toModelChatMessage(message))
	}

	return result, nil
}
