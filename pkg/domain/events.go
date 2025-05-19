package domain

import (
	"context"
)

type EventType string

const (
	UserSubscribed EventType = "user_subscribed"
	WeatherEvent   EventType = "weather_event"
)

type Event struct {
	Type    EventType
	Payload interface{}
	Context context.Context
}

type EventHandler func(ctx context.Context, event Event) error

type EventPublisher interface {
	Register(eventType EventType, handler EventHandler)

	Trigger(event Event) []error

	TriggerAsync(event Event)
}
