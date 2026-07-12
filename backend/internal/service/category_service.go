package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// CategoryService handles category business logic
type CategoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateRequest represents category creation data
type CategoryCreateRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name" binding:"required"`
	Icon     string     `json:"icon"`
	Color    string     `json:"color"`
	Type     string     `json:"type" binding:"required"`
}

// Create creates a new category
func (s *CategoryService) Create(ctx context.Context, userID uuid.UUID, req *CategoryCreateRequest) (*models.Category, error) {
	category := &models.Category{
		UserID:    userID,
		ParentID:  req.ParentID,
		Name:      req.Name,
		Icon:      req.Icon,
		Color:     req.Color,
		Type:      req.Type,
		IsSystem:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// List retrieves all categories for a user
func (s *CategoryService) List(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// GetTree retrieves category tree for a user
func (s *CategoryService) GetTree(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	return s.repo.GetTree(ctx, userID)
}

// Update updates a category
func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if name, ok := updates["name"].(string); ok {
		category.Name = name
	}
	if icon, ok := updates["icon"].(string); ok {
		category.Icon = icon
	}
	if color, ok := updates["color"].(string); ok {
		category.Color = color
	}

	category.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete removes a category
func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}