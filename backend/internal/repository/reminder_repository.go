package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// ReminderRepository defines the interface for reminder data access
type ReminderRepository interface {
	Create(ctx context.Context, reminder *models.Reminder) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Reminder, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Reminder, error)
	Update(ctx context.Context, reminder *models.Reminder) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// reminderRepository implements ReminderRepository
type reminderRepository struct {
	db *gorm.DB
}

// NewReminderRepository creates a new reminder repository
func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

func (r *reminderRepository) Create(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

func (r *reminderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Reminder, error) {
	var reminder models.Reminder
	if err := r.db.WithContext(ctx).
		Preload("Category").
		First(&reminder, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Reminder, error) {
	var reminders []models.Reminder
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("remind_time ASC").
		Find(&reminders).Error; err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *reminderRepository) Update(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Save(reminder).Error
}

func (r *reminderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Reminder{}, "id = ?", id).Error
}