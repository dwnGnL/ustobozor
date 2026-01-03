package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/barzurustami/bozor/internal/auth"
	"github.com/barzurustami/bozor/internal/domain"
	"github.com/barzurustami/bozor/internal/logger"
	"github.com/barzurustami/bozor/internal/repository"
	"github.com/barzurustami/bozor/internal/sms"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrInvalidCode  = errors.New("invalid code")
	ErrCodeExpired  = errors.New("code expired")
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type AuthService struct {
	users repository.UserRepository
	sms   sms.Sender
	jwt   *auth.JWTService

	mu    sync.Mutex
	codes map[string]codeEntry
}

type codeEntry struct {
	code      string
	expiresAt time.Time
}

func NewAuthService(users repository.UserRepository, sender sms.Sender, jwtSvc *auth.JWTService) *AuthService {
	return &AuthService{
		users: users,
		sms:   sender,
		jwt:   jwtSvc,
		codes: make(map[string]codeEntry),
	}
}

func (s *AuthService) RequestCode(ctx context.Context, phone string) error {
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	expires := time.Now().Add(5 * time.Minute)

	s.mu.Lock()
	s.codes[phone] = codeEntry{code: code, expiresAt: expires}
	s.mu.Unlock()

	logger.FromContext(ctx).Info("auth code generated", zap.String("phone", phone))
	return s.sms.Send(ctx, phone, fmt.Sprintf("Your verification code: %s", code))
}

func (s *AuthService) Register(ctx context.Context, phone, code string) (*domain.User, domain.TokenPair, error) {
	if err := s.verifyCode(phone, code); err != nil {
		return nil, domain.TokenPair{}, err
	}

	_, err := s.users.GetByPhone(ctx, phone)
	if err == nil {
		return nil, domain.TokenPair{}, ErrUserExists
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, domain.TokenPair{}, err
	}

	user := &domain.User{
		ID:        uuid.New(),
		Phone:     phone,
		Roles:     []domain.Role{domain.RoleCustomer, domain.RoleWorker},
		CreatedAt: time.Now().UTC(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, domain.TokenPair{}, err
	}

	tokens, err := s.jwt.Generate(user.ID)
	if err != nil {
		return nil, domain.TokenPair{}, err
	}

	logger.FromContext(ctx).Info("user registered", zap.String("phone", phone), zap.String("user_id", user.ID.String()))
	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, phone, code string) (*domain.User, domain.TokenPair, error) {
	if err := s.verifyCode(phone, code); err != nil {
		return nil, domain.TokenPair{}, err
	}

	user, err := s.users.GetByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.TokenPair{}, ErrUserNotFound
		}
		return nil, domain.TokenPair{}, err
	}

	tokens, err := s.jwt.Generate(user.ID)
	if err != nil {
		return nil, domain.TokenPair{}, err
	}

	logger.FromContext(ctx).Info("user login", zap.String("phone", phone), zap.String("user_id", user.ID.String()))
	return user, tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*domain.User, domain.TokenPair, error) {
	userID, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, domain.TokenPair{}, err
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.TokenPair{}, err
	}

	tokens, err := s.jwt.Generate(userID)
	if err != nil {
		return nil, domain.TokenPair{}, err
	}

	return user, tokens, nil
}

func (s *AuthService) verifyCode(phone, code string) error {
	s.mu.Lock()
	entry, ok := s.codes[phone]
	s.mu.Unlock()

	if !ok {
		return ErrInvalidCode
	}
	if time.Now().After(entry.expiresAt) {
		return ErrCodeExpired
	}
	if entry.code != code {
		return ErrInvalidCode
	}

	s.mu.Lock()
	delete(s.codes, phone)
	s.mu.Unlock()

	return nil
}
