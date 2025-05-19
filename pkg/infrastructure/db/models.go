package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Frequency string

const (
	FrequencyHourly Frequency = "HOURLY"
	FrequencyDaily  Frequency = "DAILY"
)

type SubscriptionModel struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key"`
	Email             string    `gorm:"uniqueIndex:idx_email_city"`
	City              string    `gorm:"uniqueIndex:idx_email_city"`
	Frequency         Frequency `gorm:"type:varchar(10);default:'DAILY'"`
	ConfirmationToken string    `gorm:"uniqueIndex;type:varchar(100)"`
	UnsubscribeToken  string    `gorm:"uniqueIndex;type:varchar(100)"`
	Confirmed         bool      `gorm:"default:false"`
	CreatedAt         time.Time
	ConfirmedAt       *time.Time
	LastSentAt        *time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}
