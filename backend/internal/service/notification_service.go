package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// NotificationService handles notification business logic
type NotificationService struct {
	repo repository.NotificationRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// CreateRequest represents notification creation data
type NotificationCreateRequest struct {
	Type      string      `json:"type" binding:"required"`
	Title     string      `json:"title" binding:"required"`
	Content   string      `json:"content"`
	RelatedID *uuid.UUID  `json:"related_id"`
}

// Create creates a new notification
func (s *NotificationService) Create(ctx context.Context, userID uuid.UUID, req *NotificationCreateRequest) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:    userID,
		Type:      req.Type,
		Title:     req.Title,
		Content:   req.Content,
		RelatedID: req.RelatedID,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// List retrieves all notifications for a user
func (s *NotificationService) List(ctx context.Context, userID uuid.UUID) ([]models.Notification, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	notification, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if notification.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.MarkAsRead(ctx, id)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// Delete removes a notification
func (s *NotificationService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	notification, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if notification.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}