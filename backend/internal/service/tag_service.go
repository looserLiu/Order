package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// TagService handles tag business logic
type TagService struct {
	repo repository.TagRepository
}

// NewTagService creates a new tag service
func NewTagService(repo repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// CreateRequest represents tag creation data
type TagCreateRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

// Create creates a new tag
func (s *TagService) Create(ctx context.Context, userID uuid.UUID, req *TagCreateRequest) (*models.Tag, error) {
	tag := &models.Tag{
		UserID:    userID,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// List retrieves all tags for a user
func (s *TagService) List(ctx context.Context, userID uuid.UUID) ([]models.Tag, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Update updates a tag
func (s *TagService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.Tag, error) {
	tag, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tag.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if name, ok := updates["name"].(string); ok {
		tag.Name = name
	}
	if color, ok := updates["color"].(string); ok {
		tag.Color = color
	}

	if err := s.repo.Update(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete removes a tag
func (s *TagService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	tag, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if tag.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}