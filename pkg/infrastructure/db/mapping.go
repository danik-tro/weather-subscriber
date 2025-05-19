package db

import (
	domain "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
)

func ToModel(s *domain.Subscription) *SubscriptionModel {
	return &SubscriptionModel{
		ID:                s.ID,
		Email:             s.Email,
		City:              s.City,
		Frequency:         Frequency(s.Frequency),
		ConfirmationToken: s.ConfirmationToken,
		UnsubscribeToken:  s.UnsubscribeToken,
		Confirmed:         s.Confirmed,
		CreatedAt:         s.CreatedAt,
		ConfirmedAt:       s.ConfirmedAt,
		LastSentAt:        s.LastSentAt,
	}
}

func ToDomain(m *SubscriptionModel) *domain.Subscription {
	return &domain.Subscription{
		ID:                m.ID,
		Email:             m.Email,
		City:              m.City,
		Frequency:         domain.Frequency(m.Frequency),
		ConfirmationToken: m.ConfirmationToken,
		UnsubscribeToken:  m.UnsubscribeToken,
		Confirmed:         m.Confirmed,
		CreatedAt:         m.CreatedAt,
		ConfirmedAt:       m.ConfirmedAt,
		LastSentAt:        m.LastSentAt,
	}
}

func ToDomainList(models []*SubscriptionModel) []*domain.Subscription {
	subscriptions := make([]*domain.Subscription, len(models))
	for i, model := range models {
		subscriptions[i] = ToDomain(model)
	}
	return subscriptions
}
