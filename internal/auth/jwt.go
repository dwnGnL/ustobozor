package auth

import (
	"errors"
	"time"

	"github.com/barzurustami/bozor/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (s *JWTService) Generate(userID uuid.UUID) (domain.TokenPair, error) {
	now := time.Now().UTC()
	accessExp := now.Add(s.accessTTL)
	refreshExp := now.Add(s.refreshTTL)

	accessToken, err := s.signToken(userID, accessExp, "access", s.accessSecret)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshToken, err := s.signToken(userID, refreshExp, "refresh", s.refreshSecret)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}, nil
}

func (s *JWTService) ParseAccess(token string) (uuid.UUID, error) {
	return s.parseToken(token, "access", s.accessSecret)
}

func (s *JWTService) ParseRefresh(token string) (uuid.UUID, error) {
	return s.parseToken(token, "refresh", s.refreshSecret)
}

func (s *JWTService) signToken(userID uuid.UUID, exp time.Time, tokenType string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
		"typ": tokenType,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func (s *JWTService) parseToken(token, tokenType string, secret []byte) (uuid.UUID, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	})
	if err != nil || !parsed.Valid {
		return uuid.UUID{}, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, ErrInvalidToken
	}

	if claims["typ"] != tokenType {
		return uuid.UUID{}, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.UUID{}, ErrInvalidToken
	}

	return uuid.Parse(sub)
}
