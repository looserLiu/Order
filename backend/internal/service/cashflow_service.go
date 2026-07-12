package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// CashFlowService handles cash flow projection business logic
type CashFlowService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

// NewCashFlowService creates a new cash flow service
func NewCashFlowService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
) *CashFlowService {
	return &CashFlowService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// CashFlowProjection represents a single day's cash flow projection
type CashFlowProjection struct {
	Date          string  `json:"date"`
	ProjectedBal  float64 `json:"projected_balance"`
	Income        float64 `json:"income"`
	Expense       float64 `json:"expense"`
	RecurringTx   int     `json:"recurring_transactions"`
}

// GetProjection gets cash flow projection for a number of days
func (s *CashFlowService) GetProjection(ctx context.Context, userID uuid.UUID, days int) (*CashFlowProjectionResult, error) {
	if days == 0 {
		days = 30
	}

	// Get current total balance
		var totalBalance float64
		accounts, _ := s.accountRepo.FindByUserID(ctx, userID)
		for _, acc := range accounts {
			totalBalance += acc.Balance
		}
	
		// Get recurring transactions
		recurringTxs, _ := s.transactionRepo.FindByUserID(ctx, userID, 1000, 0)
		var recurringList []models.Transaction
		for _, tx := range recurringTxs {
			if tx.IsRecurring {
				recurringList = append(recurringList, tx)
			}
		}

	// Calculate projections
	projections := make([]CashFlowProjection, days)
	currentBalance := totalBalance

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		dayIncome := 0.0
		dayExpense := 0.0
		recurringCount := 0

		for _, tx := range recurringList {
			if tx.BillDate.Day() == date.Day() || shouldRecurToday(&tx.BillDate, date) {
				if tx.Type == "income" {
					dayIncome += tx.Amount
				} else if tx.Type == "expense" {
					dayExpense += tx.Amount
				}
				recurringCount++
			}
		}

		currentBalance = currentBalance + dayIncome - dayExpense
		projections[i] = CashFlowProjection{
			Date:         dateStr,
			ProjectedBal: currentBalance,
			Income:       dayIncome,
			Expense:      dayExpense,
			RecurringTx:  recurringCount,
		}
	}

	return &CashFlowProjectionResult{
		CurrentBalance: totalBalance,
		Projections:    projections,
	}, nil
}

// CashFlowProjectionResult represents the full projection result
type CashFlowProjectionResult struct {
	CurrentBalance float64              `json:"current_balance"`
	Projections    []CashFlowProjection `json:"projections"`
}

// shouldRecurToday checks if a transaction should recur on the given date
func shouldRecurToday(lastDate *time.Time, today time.Time) bool {
	if lastDate == nil {
		return false
	}
	return today.After(*lastDate) && today.Day() == lastDate.Day()
}