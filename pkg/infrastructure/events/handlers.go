package events

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	config "github.com/danik-tro/weather-subscriber/pkg"
	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	entity "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
	domain_repository "github.com/danik-tro/weather-subscriber/pkg/domain/repository"
	value_objects "github.com/danik-tro/weather-subscriber/pkg/domain/value_object"
	"github.com/danik-tro/weather-subscriber/pkg/external/weather"
)

type Handler struct {
	EmailService   domain.EmailService
	WeatherService *weather.WeatherService
	Repository     domain_repository.SubscriptionRepository
	Config         config.Config
}

func (h *Handler) UserSubscribed() domain.EventHandler {
	return func(ctx context.Context, event domain.Event) error {
		subscription, ok := event.Payload.(*entity.Subscription)

		if !ok {
			return fmt.Errorf("invalid event type: %T", event)
		}

		confirmationLink := fmt.Sprintf("%s/confirm/%s", h.Config.BaseURL, subscription.ConfirmationToken)

		confirmationTmpl, err := template.ParseFiles("templates/confirmation.html")

		if err != nil {
			return fmt.Errorf("failed to parse confirmation template: %w", err)
		}

		confirmData := struct {
			ConfirmURL string
		}{
			ConfirmURL: confirmationLink,
		}

		var bodyBuffer bytes.Buffer
		if err := confirmationTmpl.Execute(&bodyBuffer, confirmData); err != nil {
			return fmt.Errorf("failed to execute confirmation template: %w", err)
		}

		message := bodyBuffer.String()

		h.EmailService.SendMessage(ctx, subscription.Email, "Confirm your email", message)

		return nil
	}
}

func (h *Handler) WeatherEvent() domain.EventHandler {
	return func(ctx context.Context, event domain.Event) error {
		weather, ok := event.Payload.(value_objects.WeatherEvent)

		if !ok {
			return fmt.Errorf("invalid event type: %T", event)
		}

		unsubscribeLink := fmt.Sprintf("%s/unsubscribe/%s", h.Config.BaseURL, weather.UnsubscribeToken)

		weatherTmpl, err := template.ParseFiles("templates/weather_update.html")

		if err != nil {
			return fmt.Errorf("failed to parse confirmation template: %w", err)
		}

		weatherData := struct {
			UnsubscribeURL string
			City           string
			Temperature    float64
			Humidity       float64
			Description    string
		}{
			UnsubscribeURL: unsubscribeLink,
			City:           weather.City,
			Temperature:    weather.Temperature,
			Humidity:       weather.Humidity,
			Description:    weather.Description,
		}

		var bodyBuffer bytes.Buffer
		if err := weatherTmpl.Execute(&bodyBuffer, weatherData); err != nil {
			return fmt.Errorf("failed to execute weather update template: %w", err)
		}

		message := bodyBuffer.String()

		h.EmailService.SendMessage(ctx, weather.Email, "Weather update", message)

		return nil
	}
}

func (h *Handler) FetchAndUpdateWeatherSubscribers(frequency entity.Frequency) domain.EventHandler {
	return func(ctx context.Context, event domain.Event) error {
		subscriptions, err := h.Repository.GetConfirmedSubscriptions(ctx, frequency)
		if err != nil {
			return fmt.Errorf("failed to get daily subscriptions: %w", err)
		}

		for _, subscription := range subscriptions {
			weather, err := h.WeatherService.GetWeather(ctx, subscription.City)
			if err != nil {
				fmt.Printf("Failed to get weather for city %s: %v\n", subscription.City, err)
				continue
			}
			emailData := struct {
				City           string
				Temperature    float64
				Humidity       float64
				Description    string
				UnsubscribeURL string
			}{
				City:           subscription.City,
				Temperature:    weather.Temperature,
				Humidity:       weather.Humidity,
				Description:    weather.Description,
				UnsubscribeURL: fmt.Sprintf("%s/unsubscribe/%s", h.Config.BaseURL, subscription.UnsubscribeToken),
			}

			weatherTmpl, err := template.ParseFiles("templates/weather_update.html")
			if err != nil {
				fmt.Printf("Failed to parse weather template: %v\n", err)
				continue
			}

			var bodyBuffer bytes.Buffer
			if err := weatherTmpl.Execute(&bodyBuffer, emailData); err != nil {
				fmt.Printf("Failed to execute weather template: %v\n", err)
				continue
			}

			if err := h.EmailService.SendMessage(
				ctx,
				subscription.Email,
				fmt.Sprintf("Daily Weather Update for %s", subscription.City),
				bodyBuffer.String(),
			); err != nil {
				fmt.Printf("Failed to send weather update email: %v\n", err)
				continue
			}

		}

		return nil
	}
}
