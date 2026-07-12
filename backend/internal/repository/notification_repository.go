package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// notificationRepository implements NotificationRepository
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.WithContext(ctx).First(&notification, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Notification{}, "id = ?", id).Error
}