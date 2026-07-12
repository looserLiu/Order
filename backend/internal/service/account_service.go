package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// AccountService handles account business logic
type AccountService struct {
	repo repository.AccountRepository
}

// NewAccountService creates a new account service
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

// CreateRequest represents account creation data
type AccountCreateRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Currency  string  `json:"currency"`
	Icon      string  `json:"icon"`
	Color     string  `json:"color"`
	IsDefault bool    `json:"is_default"`
}

// Create creates a new account
func (s *AccountService) Create(ctx context.Context, userID uuid.UUID, req *AccountCreateRequest) (*models.Account, error) {
	account := &models.Account{
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Currency:  req.Currency,
		Icon:      req.Icon,
		Color:     req.Color,
		IsDefault: req.IsDefault,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// List retrieves all accounts for a user
func (s *AccountService) List(ctx context.Context, userID uuid.UUID) ([]models.Account, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Get retrieves a single account
func (s *AccountService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Account, error) {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account.UserID != userID {
		return nil, errors.New("unauthorized access")
	}
	return account, nil
}

// Update updates an account
func (s *AccountService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.Account, error) {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if name, ok := updates["name"].(string); ok {
		account.Name = name
	}
	if accountType, ok := updates["type"].(string); ok {
		account.Type = accountType
	}
	if currency, ok := updates["currency"].(string); ok {
		account.Currency = currency
	}
	if icon, ok := updates["icon"].(string); ok {
		account.Icon = icon
	}
	if color, ok := updates["color"].(string); ok {
		account.Color = color
	}

	account.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// Delete removes an account
func (s *AccountService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if account.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}

// GetTotalBalance gets total balance for a user
func (s *AccountService) GetTotalBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	return s.repo.GetTotalBalance(ctx, userID)
}