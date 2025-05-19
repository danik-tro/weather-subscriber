package domain

import (
	"context"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain/value_object"
)

type GetWeatherUseCase interface {
	GetWeather(ctx context.Context, city string) (domain.Weather, error)
}
