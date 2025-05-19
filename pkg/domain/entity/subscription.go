package domain

import (
	"time"

	"github.com/google/uuid"
)

type Frequency string

const (
	FrequencyHourly Frequency = "HOURLY"
	FrequencyDaily  Frequency = "DAILY"
)

type Subscription struct {
	ID                uuid.UUID
	Email             string
	City              string
	Frequency         Frequency
	ConfirmationToken string
	UnsubscribeToken  string
	Confirmed         bool
	CreatedAt         time.Time
	ConfirmedAt       *time.Time
	LastSentAt        *time.Time
}

type Subscriber struct {
	Email            string
	City             string
	Frequency        Frequency
	UnsubscribeToken string
}

func NewSubscription(email, city string, freq Frequency) (*Subscription, error) {
	now := time.Now().UTC()

	confirmationToken, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	unsubscribeToken, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Subscription{
		ID:                uuid.New(),
		Email:             email,
		City:              city,
		Frequency:         freq,
		ConfirmationToken: confirmationToken.String(),
		UnsubscribeToken:  unsubscribeToken.String(),
		Confirmed:         false,
		CreatedAt:         now,
		ConfirmedAt:       nil,
		LastSentAt:        nil,
	}, nil
}
