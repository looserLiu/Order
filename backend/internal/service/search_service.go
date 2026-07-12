package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// SearchService handles search business logic
type SearchService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
}

// NewSearchService creates a new search service
func NewSearchService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) *SearchService {
	return &SearchService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

// SearchResult represents search results
type SearchResult struct {
	Transactions []models.Transaction `json:"transactions"`
	Accounts     []models.Account     `json:"accounts"`
	Categories   []models.Category    `json:"categories"`
}

// Search searches across transactions, accounts, and categories
func (s *SearchService) Search(ctx context.Context, userID uuid.UUID, keyword, searchType string, limit int) (*SearchResult, error) {
	result := &SearchResult{}

	if limit == 0 {
		limit = 20
	}

	// Search transactions
	if searchType == "all" || searchType == "transactions" {
		// This would use a search method in repository
		// For now, return empty
	}

	// Search accounts
	if searchType == "all" || searchType == "accounts" {
		// This would use a search method in repository
	}

	// Search categories
	if searchType == "all" || searchType == "categories" {
		// This would use a search method in repository
	}

	return result, nil
}

// SearchTransactions searches transactions by keyword
func (s *SearchService) SearchTransactions(ctx context.Context, userID uuid.UUID, keyword string, limit int) ([]models.Transaction, error) {
	return []models.Transaction{}, nil
}

// SearchAccounts searches accounts by name
func (s *SearchService) SearchAccounts(ctx context.Context, userID uuid.UUID, keyword string, limit int) ([]models.Account, error) {
	return []models.Account{}, nil
}

// SearchCategories searches categories by name
func (s *SearchService) SearchCategories(ctx context.Context, userID uuid.UUID, keyword string, limit int) ([]models.Category, error) {
	return []models.Category{}, nil
}