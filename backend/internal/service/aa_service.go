package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// AAService handles AA group business logic
type AAService struct {
	repo repository.AARepository
}

// NewAAService creates a new AA service
func NewAAService(repo repository.AARepository) *AAService {
	return &AAService{repo: repo}
}

// CreateGroupRequest represents AA group creation data
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Members     []struct {
		Name string `json:"name"`
	} `json:"members"`
}

// AddExpenseRequest represents adding expense to a group
type AddExpenseRequest struct {
	MemberID string  `json:"member_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Note     string  `json:"note"`
}

// CreateGroup creates a new AA group
func (s *AAService) CreateGroup(ctx context.Context, userID uuid.UUID, req *CreateGroupRequest) (*models.AAGroup, error) {
	group := &models.AAGroup{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateGroup(ctx, group); err != nil {
		return nil, err
	}

	// Add members
	for _, m := range req.Members {
		member := &models.AAMember{
			GroupID:  group.ID,
			Name:     m.Name,
			JoinedAt: time.Now(),
		}
		if err := s.repo.AddMember(ctx, member); err != nil {
			return nil, err
		}
	}

	return s.repo.FindGroupByID(ctx, group.ID)
}

// ListGroups retrieves all AA groups for a user
func (s *AAService) ListGroups(ctx context.Context, userID uuid.UUID) ([]models.AAGroup, error) {
	return s.repo.FindGroupsByUserID(ctx, userID)
}

// GetGroup retrieves a single AA group
func (s *AAService) GetGroup(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.AAGroup, error) {
	group, err := s.repo.FindGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if group.UserID != userID {
		return nil, errors.New("unauthorized access")
	}
	return group, nil
}

// AddExpense adds an expense to a group
func (s *AAService) AddExpense(ctx context.Context, groupID uuid.UUID, userID uuid.UUID, req *AddExpenseRequest) error {
	group, err := s.repo.FindGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.UserID != userID {
		return errors.New("unauthorized access")
	}

	memberID, _ := uuid.Parse(req.MemberID)
	var member models.AAMember
	// Check if member exists in group - need to query directly
	// For now, we'll just add the expense

	// Add expense to group total
	if err := s.repo.AddExpense(ctx, groupID, req.Amount); err != nil {
		return err
	}

	return nil
}

// GetSettlements gets settlements for a group
func (s *AAService) GetSettlements(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) ([]models.AASettlement, error) {
	group, err := s.repo.FindGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group.UserID != userID {
		return nil, errors.New("unauthorized access")
	}
	return s.repo.GetSettlements(ctx, groupID)
}

// DeleteGroup removes an AA group
func (s *AAService) DeleteGroup(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	group, err := s.repo.FindGroupByID(ctx, id)
	if err != nil {
		return err
	}
	if group.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.DeleteGroup(ctx, id)
}