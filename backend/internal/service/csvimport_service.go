package service

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// CSVImportService handles CSV import business logic
type CSVImportService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
}

// NewCSVImportService creates a new CSV import service
func NewCSVImportService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) *CSVImportService {
	return &CSVImportService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

// CSVTransaction represents a CSV transaction row
type CSVTransaction struct {
	Date     string `json:"date"`
	Type     string `json:"type"`
	Amount   string `json:"amount"`
	Category string `json:"category"`
	Account  string `json:"account"`
	Merchant string `json:"merchant"`
	Note     string `json:"note"`
}

// ImportResult represents import result
type ImportResult struct {
	Imported int              `json:"imported"`
	Failed   int              `json:"failed"`
	Errors   []CSVImportError `json:"errors,omitempty"`
}

// CSVImportError represents an import error
type CSVImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// ImportCSV imports transactions from CSV data
func (s *CSVImportService) ImportCSV(ctx context.Context, userID uuid.UUID, transactions []CSVTransaction) (*ImportResult, error) {
	result := &ImportResult{
		Imported: 0,
		Failed:   0,
		Errors:   make([]CSVImportError, 0),
	}

	for i, item := range transactions {
		billDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			billDate, err = time.Parse("2006/01/02", item.Date)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, CSVImportError{
					Row:   i + 1,
					Error: "Invalid date format",
				})
				continue
			}
		}

		amount, err := strconv.ParseFloat(item.Amount, 64)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, CSVImportError{
				Row:   i + 1,
				Error: "Invalid amount",
			})
			continue
		}

		txType := "expense"
		if item.Type == "收入" || item.Type == "income" {
			txType = "income"
		} else if item.Type == "转账" || item.Type == "transfer" {
			txType = "transfer"
		}

		var accountID uuid.UUID
		if item.Account != "" {
			// Try to find existing account
			accounts, _ := s.accountRepo.FindByUserID(ctx, userID)
			found := false
			for _, acc := range accounts {
				if acc.Name == item.Account {
					accountID = acc.ID
					found = true
					break
				}
			}
			if !found {
				// Create new account
				account := &models.Account{
					UserID: userID,
					Name:   item.Account,
					Type:   "cash",
				}
				_ = s.accountRepo.Create(ctx, account)
				accountID = account.ID
			}
		}

		var categoryID *uuid.UUID
		if item.Category != "" {
			categories, _ := s.categoryRepo.FindByUserID(ctx, userID)
			for _, cat := range categories {
				if cat.Name == item.Category {
					categoryID = &cat.ID
					break
				}
			}
		}

		tx := &models.Transaction{
			UserID:     userID,
			AccountID:  accountID,
			CategoryID: categoryID,
			Type:       txType,
			Amount:     amount,
			Merchant:   item.Merchant,
			Note:       item.Note,
			BillDate:   billDate,
		}

		if err := s.transactionRepo.Create(ctx, nil, tx); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, CSVImportError{
				Row:   i + 1,
				Error: err.Error(),
			})
			continue
		}

		result.Imported++
	}

	return result, nil
}

// ImportTransactions imports transactions from JSON data
func (s *CSVImportService) ImportTransactions(ctx context.Context, userID uuid.UUID, transactions []ImportTransactionRequest) (*ImportResult, error) {
	result := &ImportResult{
		Imported: 0,
		Failed:   0,
		Errors:   make([]CSVImportError, 0),
	}

	for i, item := range transactions {
		billDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, CSVImportError{
				Row:   i + 1,
				Error: "Invalid date format",
			})
			continue
		}

		tx := &models.Transaction{
			UserID:     userID,
			AccountID:  *item.AccountID,
			CategoryID: item.CategoryID,
			Type:       item.Type,
			Amount:     item.Amount,
			Merchant:   item.Merchant,
			Note:       item.Note,
			BillDate:   billDate,
		}

		if err := s.transactionRepo.Create(ctx, nil, tx); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, CSVImportError{
				Row:   i + 1,
				Error: err.Error(),
			})
			continue
		}

		result.Imported++
	}

	return result, nil
}

// ImportTransactionRequest represents a single import transaction
type ImportTransactionRequest struct {
	Date       string     `json:"date" binding:"required"`
	Type       string     `json:"type" binding:"required"`
	Amount     float64    `json:"amount" binding:"required"`
	CategoryID *uuid.UUID `json:"category_id"`
	AccountID  *uuid.UUID `json:"account_id" binding:"required"`
	Merchant   string     `json:"merchant"`
	Note       string     `json:"note"`
}