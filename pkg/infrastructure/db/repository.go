package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain_errors "github.com/danik-tro/weather-subscriber/pkg/domain"
	domain "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) EnsureSchema() error {
	return r.db.AutoMigrate(&SubscriptionModel{})
}

func (r *GormRepository) Save(ctx context.Context, s *domain.Subscription) error {
	model := ToModel(s)

	tx := r.db.WithContext(ctx)

	result := tx.Save(model)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *GormRepository) IsSubscribed(ctx context.Context, email, city string) (bool, error) {
	var count int64

	tx := r.db.WithContext(ctx)

	result := tx.Model(&SubscriptionModel{}).
		Where("email = ? AND city = ? AND confirmed = ?", email, city, true).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (r *GormRepository) IsComfirmationTokenExists(ctx context.Context, token string) (bool, error) {
	var count int64

	tx := r.db.WithContext(ctx)

	result := tx.Model(&SubscriptionModel{}).
		Where("confirmation_token = ?", token).
		Count(&count)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return count > 0, nil
}

func (r *GormRepository) IsUnsubscribeTokenExists(ctx context.Context, token string) (bool, error) {
	var count int64

	tx := r.db.WithContext(ctx)

	result := tx.Model(&SubscriptionModel{}).
		Where("unsubscribe_token = ?", token).
		Count(&count)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return count > 0, nil
}

func (r *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx)

	result := tx.Delete(&SubscriptionModel{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain_errors.ErrSubscriptionNotFound
	}

	return nil
}

func (r *GormRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var model SubscriptionModel

	tx := r.db.WithContext(ctx)

	result := tx.First(&model, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain_errors.ErrSubscriptionNotFound
		}
		return nil, result.Error
	}

	return ToDomain(&model), nil
}

func (r *GormRepository) Confirm(ctx context.Context, confirmationToken string) error {
	now := time.Now()

	tx := r.db.WithContext(ctx)

	result := tx.Model(&SubscriptionModel{}).
		Where("confirmation_token = ?", confirmationToken).
		Updates(map[string]interface{}{
			"confirmed":    true,
			"confirmed_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain_errors.ErrSubscriptionNotFound
	}

	return nil
}

func (r *GormRepository) GetConfirmedSubscriptions(ctx context.Context, frequency domain.Frequency) ([]domain.Subscriber, error) {
	var models []SubscriptionModel

	tx := r.db.WithContext(ctx)

	result := tx.Where("confirmed = ? AND frequency = ?", true, frequency).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	subscriptions := make([]domain.Subscriber, len(models))
	for i, model := range models {
		subscriptions[i] = domain.Subscriber{
			Email:            model.Email,
			City:             model.City,
			Frequency:        domain.Frequency(model.Frequency),
			UnsubscribeToken: model.UnsubscribeToken,
		}
	}

	return subscriptions, nil
}

func (r *GormRepository) FindByConfirmationToken(ctx context.Context, token string) (*domain.Subscription, error) {
	var model SubscriptionModel

	tx := r.db.WithContext(ctx)

	result := tx.Model(&SubscriptionModel{}).
		Where("confirmation_token = ?", token).
		First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain_errors.ErrSubscriptionNotFound
		}
		return nil, result.Error
	}

	return ToDomain(&model), nil
}

func (r *GormRepository) FindByUnsubscribeToken(ctx context.Context, token string) (*domain.Subscription, error) {
	var model SubscriptionModel

	tx := r.db.WithContext(ctx)

	result := tx.Model(&SubscriptionModel{}).
		Where("unsubscribe_token = ?", token).
		First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain_errors.ErrSubscriptionNotFound
		}
		return nil, result.Error
	}

	return ToDomain(&model), nil
}
