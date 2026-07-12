package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// ReminderService handles reminder business logic
type ReminderService struct {
	repo repository.ReminderRepository
}

// NewReminderService creates a new reminder service
func NewReminderService(repo repository.ReminderRepository) *ReminderService {
	return &ReminderService{repo: repo}
}

// CreateRequest represents reminder creation data
type ReminderCreateRequest struct {
	Title      string     `json:"title" binding:"required"`
	Content    string     `json:"content"`
	RemindTime time.Time  `json:"remind_time" binding:"required"`
	RepeatType string     `json:"repeat_type"`
	CategoryID *uuid.UUID `json:"category_id"`
}

// Create creates a new reminder
func (s *ReminderService) Create(ctx context.Context, userID uuid.UUID, req *ReminderCreateRequest) (*models.Reminder, error) {
	reminder := &models.Reminder{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		RemindTime: req.RemindTime,
		RepeatType: req.RepeatType,
		CategoryID: req.CategoryID,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, reminder); err != nil {
		return nil, err
	}

	return reminder, nil
}

// List retrieves all reminders for a user
func (s *ReminderService) List(ctx context.Context, userID uuid.UUID) ([]models.Reminder, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Update updates a reminder
func (s *ReminderService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.Reminder, error) {
	reminder, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reminder.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if title, ok := updates["title"].(string); ok {
		reminder.Title = title
	}
	if content, ok := updates["content"].(string); ok {
		reminder.Content = content
	}
	if remindTime, ok := updates["remind_time"].(time.Time); ok {
		reminder.RemindTime = remindTime
	}
	if repeatType, ok := updates["repeat_type"].(string); ok {
		reminder.RepeatType = repeatType
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		reminder.IsActive = isActive
	}

	reminder.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, reminder); err != nil {
		return nil, err
	}

	return reminder, nil
}

// Delete removes a reminder
func (s *ReminderService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	reminder, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if reminder.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}