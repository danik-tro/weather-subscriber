package weather

import (
	"context"
	"time"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain/value_object"
)

const (
	cacheTTL = 20 * time.Minute
)

type WeatherClient interface {
	GetCurrentWeather(ctx context.Context, city string) (*WeatherData, error)
}

type WeatherCache interface {
	GetWeather(ctx context.Context, city string) (*domain.Weather, error)

	SetWeather(ctx context.Context, city string, weather *domain.Weather, ttl time.Duration) error
}

type WeatherService struct {
	weatherClient WeatherClient
	cache         WeatherCache
}

func NewWeatherService(weatherClient WeatherClient, cache WeatherCache) *WeatherService {
	return &WeatherService{weatherClient: weatherClient, cache: cache}
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (*domain.Weather, error) {
	weather, err := s.cache.GetWeather(ctx, city)
	if err != nil {
		return nil, err
	}

	if weather != nil {
		return weather, nil
	}

	weatherData, err := s.weatherClient.GetCurrentWeather(ctx, city)

	if err != nil {
		return nil, err
	}

	weather = &domain.Weather{
		Temperature: weatherData.Current.TempC,
		Humidity:    weatherData.Current.Humidity,
		Description: weatherData.Current.Condition.Text,
	}

	if err := s.cache.SetWeather(ctx, city, weather, cacheTTL); err != nil {
		return nil, err
	}

	return weather, nil
}
