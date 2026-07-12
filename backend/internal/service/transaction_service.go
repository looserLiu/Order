package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
	"gorm.io/gorm"
)

// TransactionService handles transaction business logic
type TransactionService struct {
	repo         repository.TransactionRepository
	accountRepo  repository.AccountRepository
	categoryRepo repository.CategoryRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	repo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) *TransactionService {
	return &TransactionService{
		repo:         repo,
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateRequest represents transaction creation data
type CreateRequest struct {
	AccountID       uuid.UUID  `json:"account_id"`
	TargetAccountID *uuid.UUID `json:"target_account_id"`
	CategoryID      *uuid.UUID `json:"category_id"`
	Type            string     `json:"type"`
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	ExchangeRate    float64    `json:"exchange_rate"`
	Tags            []uuid.UUID `json:"tags"`
	Merchant        string     `json:"merchant"`
	Note            string     `json:"note"`
	BillDate        time.Time  `json:"bill_date"`
	IsRecurring     bool       `json:"is_recurring"`
	RecurringRule   string     `json:"recurring_rule"`
}

// Create creates a new transaction with atomic account balance update
func (s *TransactionService) Create(ctx context.Context, userID uuid.UUID, req *CreateRequest) (*models.Transaction, error) {
	// Validate account ownership
	account, err := s.accountRepo.FindByID(ctx, req.AccountID)
	if err != nil {
		return nil, errors.New("account not found")
	}
	if account.UserID != userID {
		return nil, errors.New("unauthorized account access")
	}

	// Validate target account if provided
	if req.TargetAccountID != nil {
		targetAccount, err := s.accountRepo.FindByID(ctx, *req.TargetAccountID)
		if err != nil {
			return nil, errors.New("target account not found")
		}
		if targetAccount.UserID != userID {
			return nil, errors.New("unauthorized target account access")
		}
	}

	// Validate category if provided
	if req.CategoryID != nil {
		category, err := s.categoryRepo.FindByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		if category.UserID != userID {
			return nil, errors.New("unauthorized category access")
		}
	}

	transaction := &models.Transaction{
		UserID:          userID,
		AccountID:       req.AccountID,
		TargetAccountID: req.TargetAccountID,
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          req.Amount,
		Currency:        req.Currency,
		ExchangeRate:    req.ExchangeRate,
		Merchant:        req.Merchant,
		Note:            req.Note,
		BillDate:        req.BillDate,
		IsRecurring:     req.IsRecurring,
		RecurringRule:   req.RecurringRule,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Use transaction for atomic operation
	err = s.repo.Create(ctx, nil, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// List retrieves transactions for a user
func (s *TransactionService) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	return s.repo.FindByUserID(ctx, userID, limit, offset)
}

// ListByDate retrieves transactions by date range
func (s *TransactionService) ListByDate(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.Transaction, error) {
	return s.repo.FindByDateRange(ctx, userID, startDate, endDate)
}

// Get retrieves a single transaction
func (s *TransactionService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Transaction, error) {
	transaction, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userID {
		return nil, errors.New("unauthorized access")
	}
	return transaction, nil
}

// Update updates a transaction
func (s *TransactionService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.Transaction, error) {
	transaction, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// Apply updates
	if accountID, ok := updates["account_id"].(uuid.UUID); ok {
		transaction.AccountID = accountID
	}
	if categoryID, ok := updates["category_id"].(*uuid.UUID); ok {
		transaction.CategoryID = categoryID
	}
	if amount, ok := updates["amount"].(float64); ok {
		transaction.Amount = amount
	}
	if merchant, ok := updates["merchant"].(string); ok {
		transaction.Merchant = merchant
	}
	if note, ok := updates["note"].(string); ok {
		transaction.Note = note
	}
	if billDate, ok := updates["bill_date"].(time.Time); ok {
		transaction.BillDate = billDate
	}

	transaction.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, nil, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// Delete removes a transaction
func (s *TransactionService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	transaction, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if transaction.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, nil, id)
}

// BatchDelete removes multiple transactions
func (s *TransactionService) BatchDelete(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error {
	return s.repo.BatchDelete(ctx, nil, ids)
}