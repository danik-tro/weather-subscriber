package domain

import (
	"context"
)

type ConfirmSubscriptionUseCase interface {
	Confirm(ctx context.Context, token string) error
}
