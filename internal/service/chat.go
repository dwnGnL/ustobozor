package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/logger"
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/barzurustami/bozor/internal/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrChatForbidden = errors.New("chat access forbidden")
	ErrChatSelf      = errors.New("cannot start chat with yourself")
	ErrEmptyMessage  = errors.New("message text or photo required")
	ErrReadOwn       = errors.New("cannot mark own message as read")
)

type ChatService struct {
	chats    repository.ChatRepository
	messages repository.MessageRepository
	requests repository.RequestRepository
	storage  *storage.LocalStorage

	mu              sync.RWMutex
	messageSubs     map[uuid.UUID]map[chan domain.ChatMessage]struct{}
	messageReadSubs map[uuid.UUID]map[chan domain.ChatMessage]struct{}
}

func NewChatService(
	chats repository.ChatRepository,
	messages repository.MessageRepository,
	requests repository.RequestRepository,
	storage *storage.LocalStorage,
) *ChatService {
	return &ChatService{
		chats:           chats,
		messages:        messages,
		requests:        requests,
		storage:         storage,
		messageSubs:     make(map[uuid.UUID]map[chan domain.ChatMessage]struct{}),
		messageReadSubs: make(map[uuid.UUID]map[chan domain.ChatMessage]struct{}),
	}
}

func (s *ChatService) CreateChat(ctx context.Context, requestID, initiatorID uuid.UUID) (*domain.Chat, error) {
	request, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	creatorID := request.CustomerID
	if initiatorID == creatorID {
		return nil, ErrChatSelf
	}

	chat, err := s.chats.GetByRequestAndInitiator(ctx, requestID, initiatorID)
	if err == nil {
		return chat, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	now := time.Now().UTC()
	chat = &domain.Chat{
		ID:          uuid.New(),
		RequestID:   requestID,
		CreatorID:   creatorID,
		InitiatorID: initiatorID,
		CreatedAt:   now,
	}

	if err := s.chats.Create(ctx, chat); err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info("chat created", zap.String("chat_id", chat.ID.String()))
	return chat, nil
}

func (s *ChatService) ListChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	return s.chats.ListByUser(ctx, userID)
}

func (s *ChatService) ListMessages(ctx context.Context, chatID, userID uuid.UUID, limit, offset int32) ([]domain.ChatMessage, error) {
	if _, err := s.ensureParticipant(ctx, chatID, userID); err != nil {
		return nil, err
	}
	return s.messages.ListByChat(ctx, chatID, limit, offset)
}

func (s *ChatService) SendMessage(ctx context.Context, chatID, senderID uuid.UUID, text string, file *graphql.Upload) (*domain.ChatMessage, error) {
	if _, err := s.ensureParticipant(ctx, chatID, senderID); err != nil {
		return nil, err
	}

	cleanText := strings.TrimSpace(text)
	if cleanText == "" && file == nil {
		return nil, ErrEmptyMessage
	}

	photoPath := ""
	if file != nil {
		path, err := s.storage.Save(ctx, *file)
		if err != nil {
			return nil, err
		}
		photoPath = path
	}

	now := time.Now().UTC()
	message := &domain.ChatMessage{
		ID:        uuid.New(),
		ChatID:    chatID,
		SenderID:  senderID,
		Text:      cleanText,
		PhotoPath: photoPath,
		CreatedAt: now,
		ReadAt:    nil,
	}

	if err := s.messages.Create(ctx, message); err != nil {
		return nil, err
	}

	if err := s.chats.UpdateLastMessageAt(ctx, chatID, now); err != nil {
		return nil, err
	}

	logger.FromContext(ctx).Info(
		"chat message sent",
		zap.String("chat_id", chatID.String()),
		zap.String("sender_id", senderID.String()),
	)

	s.publishMessage(*message)
	return message, nil
}

func (s *ChatService) MarkChatRead(ctx context.Context, chatID, userID uuid.UUID) ([]domain.ChatMessage, error) {
	if _, err := s.ensureParticipant(ctx, chatID, userID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	messages, err := s.messages.MarkReadByChat(ctx, chatID, userID, now)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		s.publishRead(message)
	}

	return messages, nil
}

func (s *ChatService) MarkMessageRead(ctx context.Context, messageID, userID uuid.UUID) (*domain.ChatMessage, error) {
	message, err := s.messages.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if _, err := s.ensureParticipant(ctx, message.ChatID, userID); err != nil {
		return nil, err
	}
	if message.SenderID == userID {
		return nil, ErrReadOwn
	}
	if message.ReadAt != nil {
		return message, nil
	}

	updated, err := s.messages.MarkReadByID(ctx, messageID, time.Now().UTC())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return s.messages.GetByID(ctx, messageID)
		}
		return nil, err
	}

	s.publishRead(*updated)
	return updated, nil
}

func (s *ChatService) SubscribeMessages(ctx context.Context, chatID, userID uuid.UUID) (<-chan domain.ChatMessage, error) {
	if _, err := s.ensureParticipant(ctx, chatID, userID); err != nil {
		return nil, err
	}

	ch := make(chan domain.ChatMessage, 1)

	s.mu.Lock()
	if s.messageSubs[chatID] == nil {
		s.messageSubs[chatID] = make(map[chan domain.ChatMessage]struct{})
	}
	s.messageSubs[chatID][ch] = struct{}{}
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.messageSubs[chatID], ch)
		if len(s.messageSubs[chatID]) == 0 {
			delete(s.messageSubs, chatID)
		}
		s.mu.Unlock()
		close(ch)
	}()

	return ch, nil
}

func (s *ChatService) SubscribeReads(ctx context.Context, chatID, userID uuid.UUID) (<-chan domain.ChatMessage, error) {
	if _, err := s.ensureParticipant(ctx, chatID, userID); err != nil {
		return nil, err
	}

	ch := make(chan domain.ChatMessage, 1)

	s.mu.Lock()
	if s.messageReadSubs[chatID] == nil {
		s.messageReadSubs[chatID] = make(map[chan domain.ChatMessage]struct{})
	}
	s.messageReadSubs[chatID][ch] = struct{}{}
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.messageReadSubs[chatID], ch)
		if len(s.messageReadSubs[chatID]) == 0 {
			delete(s.messageReadSubs, chatID)
		}
		s.mu.Unlock()
		close(ch)
	}()

	return ch, nil
}

func (s *ChatService) publishMessage(message domain.ChatMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for ch := range s.messageSubs[message.ChatID] {
		select {
		case ch <- message:
		default:
		}
	}
}

func (s *ChatService) publishRead(message domain.ChatMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for ch := range s.messageReadSubs[message.ChatID] {
		select {
		case ch <- message:
		default:
		}
	}
}

func (s *ChatService) ensureParticipant(ctx context.Context, chatID, userID uuid.UUID) (*domain.Chat, error) {
	chat, err := s.chats.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if userID != chat.CreatorID && userID != chat.InitiatorID {
		return nil, ErrChatForbidden
	}
	return chat, nil
}
