package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// ReportService handles report business logic
type ReportService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
}

// NewReportService creates a new report service
func NewReportService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) *ReportService {
	return &ReportService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

// SummaryData represents summary report data
type SummaryData struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

// TrendData represents trend report data
type TrendData struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Count   int     `json:"count"`
}

// CategoryData represents category report data
type CategoryReportData struct {
	CategoryID    *uuid.UUID `json:"category_id"`
	CategoryName  string     `json:"category_name"`
	CategoryIcon  string     `json:"category_icon"`
	CategoryColor string     `json:"category_color"`
	Total         float64    `json:"total"`
	Percentage    float64    `json:"percentage"`
	Count         int        `json:"count"`
}

// AccountData represents account report data
type AccountReportData struct {
	AccountID   uuid.UUID `json:"account_id"`
	AccountName string    `json:"account_name"`
	Total       float64   `json:"total"`
	Count       int       `json:"count"`
}

// MerchantData represents merchant report data
type MerchantReportData struct {
	Merchant string  `json:"merchant"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

// MonthlyData represents monthly comparison data
type MonthlyReportData struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

// GetSummary gets transaction summary for a date range
func (s *ReportService) GetSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*SummaryData, error) {
	var income, expense float64

	// Get all transactions in date range
	transactions, err := s.transactionRepo.FindByDateRange(ctx, userID, startDate, endDate)
	if err == nil {
		for _, tx := range transactions {
			if tx.Type == "income" {
				income += tx.Amount
			} else if tx.Type == "expense" {
				expense += tx.Amount
			}
		}
	}

	return &SummaryData{
		Income:  income,
		Expense: expense,
		Balance: income - expense,
	}, nil
}

// GetTrend gets transaction trend data
func (s *ReportService) GetTrend(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]TrendData, error) {
	// This would typically use raw SQL for efficiency
	// For now, return empty slice
	return []TrendData{}, nil
}

// GetByCategory gets spending by category
func (s *ReportService) GetByCategory(ctx context.Context, userID uuid.UUID, txType string, startDate, endDate time.Time) ([]CategoryReportData, error) {
	return []CategoryReportData{}, nil
}

// GetByAccount gets spending by account
func (s *ReportService) GetByAccount(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]AccountReportData, error) {
	return []AccountReportData{}, nil
}

// GetByMerchant gets spending by merchant
func (s *ReportService) GetByMerchant(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]MerchantReportData, error) {
	return []MerchantReportData{}, nil
}

// GetMonthlyCompare gets monthly comparison data
func (s *ReportService) GetMonthlyCompare(ctx context.Context, userID uuid.UUID, months int) ([]MonthlyReportData, error) {
	return []MonthlyReportData{}, nil
}

// ExportTransactions exports transactions for a date range
func (s *ReportService) ExportTransactions(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.Transaction, error) {
	return s.transactionRepo.FindByDateRange(ctx, userID, startDate, endDate)
}