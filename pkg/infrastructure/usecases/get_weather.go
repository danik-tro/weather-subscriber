package usecases

import (
	"context"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	domain_usecases "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	value_object "github.com/danik-tro/weather-subscriber/pkg/domain/value_object"
	weather "github.com/danik-tro/weather-subscriber/pkg/external/weather"
)

type GetWeatherUseCase struct {
	weatherService weather.WeatherService
}

func (uc *GetWeatherUseCase) GetWeather(ctx context.Context, city string) (value_object.Weather, error) {
	weather, err := uc.weatherService.GetWeather(ctx, city)
	if err != nil {
		return value_object.Weather{}, err
	}

	if weather == nil {
		return value_object.Weather{}, domain.ErrCityNotFound
	}

	return *weather, nil
}

func NewGetWeatherUseCase(weatherService weather.WeatherService) domain_usecases.GetWeatherUseCase {
	return &GetWeatherUseCase{
		weatherService: weatherService,
	}
}
