package sms

import (
	"context"

	"github.com/barzurustami/bozor/internal/logger"
	"go.uber.org/zap"
)

type MockSender struct{}

func NewMockSender() *MockSender {
	return &MockSender{}
}

func (s *MockSender) Send(ctx context.Context, phone, message string) error {
	log := logger.FromContext(ctx)
	log.Info("sms mock send", zap.String("phone", phone), zap.String("message", message))
	return nil
}
