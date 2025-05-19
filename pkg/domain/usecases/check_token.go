package domain

import (
	"context"
)

type CheckTokensUseCase interface {
	CheckConfirmationToken(ctx context.Context, token string) (bool, error)
	CheckUnsubscribeToken(ctx context.Context, token string) (bool, error)
}
