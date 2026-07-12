package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// GoalService handles financial goal business logic
type GoalService struct {
	repo repository.GoalRepository
}

// NewGoalService creates a new goal service
func NewGoalService(repo repository.GoalRepository) *GoalService {
	return &GoalService{repo: repo}
}

// CreateRequest represents goal creation data
type GoalCreateRequest struct {
	Name         string     `json:"name" binding:"required"`
	TargetAmount float64    `json:"target_amount" binding:"required"`
	Deadline     *time.Time `json:"deadline"`
	Category     string     `json:"category"`
	Note         string     `json:"note"`
}

// Create creates a new financial goal
func (s *GoalService) Create(ctx context.Context, userID uuid.UUID, req *GoalCreateRequest) (*models.FinancialGoal, error) {
	goal := &models.FinancialGoal{
		UserID:       userID,
		Name:         req.Name,
		TargetAmount: req.TargetAmount,
		CurrentAmount: 0,
		Deadline:     req.Deadline,
		Category:     req.Category,
		Status:       "in_progress",
		Note:         req.Note,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, goal); err != nil {
		return nil, err
	}

	return goal, nil
}

// List retrieves all goals for a user
func (s *GoalService) List(ctx context.Context, userID uuid.UUID) ([]models.FinancialGoal, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Update updates a goal
func (s *GoalService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.FinancialGoal, error) {
	goal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if goal.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if name, ok := updates["name"].(string); ok {
		goal.Name = name
	}
	if targetAmount, ok := updates["target_amount"].(float64); ok {
		goal.TargetAmount = targetAmount
	}
	if deadline, ok := updates["deadline"].(*time.Time); ok {
		goal.Deadline = deadline
	}
	if status, ok := updates["status"].(string); ok {
		goal.Status = status
	}
	if note, ok := updates["note"].(string); ok {
		goal.Note = note
	}

	goal.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, goal); err != nil {
		return nil, err
	}

	return goal, nil
}

// AddAmount adds amount to a goal
func (s *GoalService) AddAmount(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount float64) (*models.FinancialGoal, error) {
	goal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if goal.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if err := s.repo.AddAmount(ctx, id, amount); err != nil {
		return nil, err
	}

	goal.CurrentAmount += amount
	return goal, nil
}

// Delete removes a goal
func (s *GoalService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	goal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if goal.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}