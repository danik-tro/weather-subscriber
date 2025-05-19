package usecases

import (
	"context"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	domain_entity "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
	domain_repository "github.com/danik-tro/weather-subscriber/pkg/domain/repository"
	domain_usecases "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	weather "github.com/danik-tro/weather-subscriber/pkg/external/weather"
)

type SubscribeWeatherUseCase struct {
	repo           domain_repository.SubscriptionRepository
	publisher      domain.EventPublisher
	weatherService weather.WeatherService
}

func (uc *SubscribeWeatherUseCase) Subscribe(ctx context.Context, email, city string, freq domain_entity.Frequency) error {
	weather, err := uc.weatherService.GetWeather(ctx, city)
	if err != nil {
		return err
	}

	if weather == nil {
		return domain.ErrCityNotFound
	}

	isSubscribed, err := uc.repo.IsSubscribed(ctx, email, city)
	if err != nil {
		return err
	}

	if isSubscribed {
		return domain.ErrSubscriptionAlreadyExists
	}

	sub, err := domain_entity.NewSubscription(email, city, freq)
	if err != nil {
		return err
	}

	if err := uc.repo.Save(ctx, sub); err != nil {
		return err
	}

	uc.publisher.TriggerAsync(
		domain.Event{
			Type:    domain.UserSubscribed,
			Payload: sub,
			Context: ctx,
		},
	)

	return nil
}

func NewSubscribeWeatherUseCase(repo domain_repository.SubscriptionRepository, publisher domain.EventPublisher, weatherService weather.WeatherService) domain_usecases.SubscribeWeatherUseCase {
	return &SubscribeWeatherUseCase{repo: repo, publisher: publisher, weatherService: weatherService}
}
