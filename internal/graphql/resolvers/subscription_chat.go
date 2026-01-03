package resolvers

import (
	"context"
	"fmt"

	"github.com/barzurustami/bozor/internal/graphql/model"
	"github.com/barzurustami/bozor/internal/middleware"
	"github.com/google/uuid"
)

func resolveChatMessageAdded(ctx context.Context, r *Resolver, chatID string) (<-chan *model.ChatMessage, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	parsedID, err := uuid.Parse(chatID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id")
	}

	domainCh, err := r.ChatService.SubscribeMessages(ctx, parsedID, userID)
	if err != nil {
		return nil, err
	}

	out := make(chan *model.ChatMessage, 1)
	go func() {
		defer close(out)
		for msg := range domainCh {
			out <- toModelChatMessage(msg)
		}
	}()

	return out, nil
}

func resolveChatMessageRead(ctx context.Context, r *Resolver, chatID string) (<-chan *model.ChatMessage, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	parsedID, err := uuid.Parse(chatID)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id")
	}

	domainCh, err := r.ChatService.SubscribeReads(ctx, parsedID, userID)
	if err != nil {
		return nil, err
	}

	out := make(chan *model.ChatMessage, 1)
	go func() {
		defer close(out)
		for msg := range domainCh {
			out <- toModelChatMessage(msg)
		}
	}()

	return out, nil
}
