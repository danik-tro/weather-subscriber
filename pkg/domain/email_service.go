package domain

import "context"

type EmailService interface {
	SendMessage(ctx context.Context, recipient string, subject string, body string) error
}
