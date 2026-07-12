package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// BudgetService handles budget business logic
type BudgetService struct {
	repo repository.BudgetRepository
}

// NewBudgetService creates a new budget service
func NewBudgetService(repo repository.BudgetRepository) *BudgetService {
	return &BudgetService{repo: repo}
}

// CreateRequest represents budget creation data
type BudgetCreateRequest struct {
	CategoryID     *uuid.UUID `json:"category_id"`
	Amount         float64    `json:"amount" binding:"required"`
	Period         string     `json:"period" binding:"required"`
	StartDate      time.Time  `json:"start_date" binding:"required"`
	EndDate        *time.Time `json:"end_date"`
	AlertThreshold float64    `json:"alert_threshold"`
}

// Create creates a new budget
func (s *BudgetService) Create(ctx context.Context, userID uuid.UUID, req *BudgetCreateRequest) (*models.Budget, error) {
	budget := &models.Budget{
		UserID:         userID,
		CategoryID:     req.CategoryID,
		Amount:         req.Amount,
		Period:         req.Period,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		AlertThreshold:  req.AlertThreshold,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(ctx, budget); err != nil {
		return nil, err
	}

	return budget, nil
}

// List retrieves all budgets for a user
func (s *BudgetService) List(ctx context.Context, userID uuid.UUID) ([]models.Budget, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Get retrieves a single budget
func (s *BudgetService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Budget, error) {
	budget, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if budget.UserID != userID {
		return nil, errors.New("unauthorized access")
	}
	return budget, nil
}

// Update updates a budget
func (s *BudgetService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.Budget, error) {
	budget, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if budget.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if amount, ok := updates["amount"].(float64); ok {
		budget.Amount = amount
	}
	if endDate, ok := updates["end_date"].(*time.Time); ok {
		budget.EndDate = endDate
	}
	if alertThreshold, ok := updates["alert_threshold"].(float64); ok {
		budget.AlertThreshold = alertThreshold
	}

	budget.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, budget); err != nil {
		return nil, err
	}

	return budget, nil
}

// Delete removes a budget
func (s *BudgetService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	budget, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if budget.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}

// GetProgress gets budget progress
func (s *BudgetService) GetProgress(ctx context.Context, id uuid.UUID, userID uuid.UUID) (float64, float64, error) {
	budget, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return 0, 0, err
	}
	if budget.UserID != userID {
		return 0, 0, errors.New("unauthorized access")
	}
	return s.repo.GetProgress(ctx, id)
}