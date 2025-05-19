package domain

import (
	"context"
)

type UnsubscribeUseCase interface {
	Unsubscribe(ctx context.Context, token string) error
}
