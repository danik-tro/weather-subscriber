package domain

import (
	"context"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
)

type SubscribeWeatherUseCase interface {
	Subscribe(ctx context.Context, email, city string, freq domain.Frequency) error
}
