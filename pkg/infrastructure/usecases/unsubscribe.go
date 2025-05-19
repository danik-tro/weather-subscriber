package usecases

import (
	"context"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	domain_repository "github.com/danik-tro/weather-subscriber/pkg/domain/repository"
	domain_usecases "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
)

type Unsubscribe struct {
	repo domain_repository.SubscriptionRepository
}

func (uc *Unsubscribe) Unsubscribe(ctx context.Context, token string) error {
	subscription, err := uc.repo.FindByUnsubscribeToken(ctx, token)

	if err != nil {
		return err
	}

	if subscription == nil {
		return domain.ErrSubscriptionNotFound
	}

	if err := uc.repo.Delete(ctx, subscription.ID); err != nil {
		return err
	}

	return nil
}

func NewUnsubscribeUseCase(repo domain_repository.SubscriptionRepository) domain_usecases.UnsubscribeUseCase {
	return &Unsubscribe{
		repo: repo,
	}
}
