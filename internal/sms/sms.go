package sms

import "context"

type Sender interface {
	Send(ctx context.Context, phone, message string) error
}
