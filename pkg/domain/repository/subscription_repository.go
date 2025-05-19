package domain

import (
	"context"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	FindByConfirmationToken(ctx context.Context, token string) (*domain.Subscription, error)
	FindByUnsubscribeToken(ctx context.Context, token string) (*domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, subscription *domain.Subscription) error
	IsSubscribed(ctx context.Context, email, city string) (bool, error)
	IsComfirmationTokenExists(ctx context.Context, token string) (bool, error)
	IsUnsubscribeTokenExists(ctx context.Context, token string) (bool, error)
	GetConfirmedSubscriptions(ctx context.Context, frequency domain.Frequency) ([]domain.Subscriber, error)
	Confirm(ctx context.Context, token string) error
}
