package usecases

import (
	"context"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	domain_repository "github.com/danik-tro/weather-subscriber/pkg/domain/repository"
	domain_usecases "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
)

type ConfirmSubscription struct {
	repo domain_repository.SubscriptionRepository
}

func (uc *ConfirmSubscription) Confirm(ctx context.Context, token string) error {
	subscription, err := uc.repo.IsComfirmationTokenExists(ctx, token)

	if err != nil {
		return err
	}

	if !subscription {
		return domain.ErrSubscriptionNotFound
	}

	if err := uc.repo.Confirm(ctx, token); err != nil {
		return err
	}

	return nil
}

func NewConfirmSubscription(repo domain_repository.SubscriptionRepository) domain_usecases.ConfirmSubscriptionUseCase {
	return &ConfirmSubscription{
		repo: repo,
	}
}
