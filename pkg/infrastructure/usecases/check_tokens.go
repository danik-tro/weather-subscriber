package usecases

import (
	"context"

	domain_repository "github.com/danik-tro/weather-subscriber/pkg/domain/repository"
	domain_usecases "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
)

type CheckTokens struct {
	repository domain_repository.SubscriptionRepository
}

func (c *CheckTokens) CheckConfirmationToken(ctx context.Context, token string) (bool, error) {
	return c.repository.IsComfirmationTokenExists(ctx, token)
}

func (c *CheckTokens) CheckUnsubscribeToken(ctx context.Context, token string) (bool, error) {
	return c.repository.IsUnsubscribeTokenExists(ctx, token)
}

func NewCheckTokens(repository domain_repository.SubscriptionRepository) domain_usecases.CheckTokensUseCase {
	return &CheckTokens{repository: repository}
}
