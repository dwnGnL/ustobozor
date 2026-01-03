package resolvers

import (
	"context"

	"github.com/barzurustami/bozor/internal/graphql/generated"
	"github.com/barzurustami/bozor/internal/graphql/model"
)

func (r *mutationResolver) RequestSMSCode(ctx context.Context, phone string) (bool, error) {
	return resolveRequestSMSCode(ctx, r.Resolver, phone)
}

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthPayload, error) {
	return resolveRegister(ctx, r.Resolver, input)
}

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	return resolveLogin(ctx, r.Resolver, input)
}

func (r *mutationResolver) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthPayload, error) {
	return resolveRefreshToken(ctx, r.Resolver, refreshToken)
}

func (r *mutationResolver) CreateRequest(ctx context.Context, input model.CreateRequestInput) (*model.JobRequest, error) {
	return resolveCreateRequest(ctx, r.Resolver, input)
}

func (r *mutationResolver) CreateChat(ctx context.Context, requestID string) (*model.Chat, error) {
	return resolveCreateChat(ctx, r.Resolver, requestID)
}

func (r *mutationResolver) SendMessage(ctx context.Context, input model.SendMessageInput) (*model.ChatMessage, error) {
	return resolveSendMessage(ctx, r.Resolver, input)
}

func (r *mutationResolver) MarkChatRead(ctx context.Context, chatID string) ([]*model.ChatMessage, error) {
	return resolveMarkChatRead(ctx, r.Resolver, chatID)
}

func (r *mutationResolver) MarkMessageRead(ctx context.Context, messageID string) (*model.ChatMessage, error) {
	return resolveMarkMessageRead(ctx, r.Resolver, messageID)
}

func (r *mutationResolver) UploadPhotos(ctx context.Context, input model.UploadPhotosInput) ([]*model.Photo, error) {
	return resolveUploadPhotos(ctx, r.Resolver, input)
}

func (r *mutationResolver) UpsertProfile(ctx context.Context, input model.ProfileInput) (*model.Profile, error) {
	return resolveUpsertProfile(ctx, r.Resolver, input)
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return resolveMe(ctx, r.Resolver)
}

func (r *queryResolver) Chats(ctx context.Context) ([]*model.Chat, error) {
	return resolveChats(ctx, r.Resolver)
}

func (r *queryResolver) ChatMessages(ctx context.Context, chatID string, limit *int, offset *int) ([]*model.ChatMessage, error) {
	return resolveChatMessages(ctx, r.Resolver, chatID, limit, offset)
}

func (r *subscriptionResolver) ChatMessageAdded(ctx context.Context, chatID string) (<-chan *model.ChatMessage, error) {
	return resolveChatMessageAdded(ctx, r.Resolver, chatID)
}

func (r *subscriptionResolver) ChatMessageRead(ctx context.Context, chatID string) (<-chan *model.ChatMessage, error) {
	return resolveChatMessageRead(ctx, r.Resolver, chatID)
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

type subscriptionResolver struct{ *Resolver }
