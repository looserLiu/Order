package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// BudgetRepository defines the interface for budget data access
type BudgetRepository interface {
	Create(ctx context.Context, budget *models.Budget) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Budget, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Budget, error)
	FindByCategoryAndPeriod(ctx context.Context, userID, categoryID uuid.UUID, period string, startDate time.Time) (*models.Budget, error)
	Update(ctx context.Context, budget *models.Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetProgress(ctx context.Context, budgetID uuid.UUID) (float64, float64, error)
}

// budgetRepository implements BudgetRepository
type budgetRepository struct {
	db *gorm.DB
}

// NewBudgetRepository creates a new budget repository
func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(ctx context.Context, budget *models.Budget) error {
	return r.db.WithContext(ctx).Create(budget).Error
}

func (r *budgetRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Budget, error) {
	var budget models.Budget
	if err := r.db.WithContext(ctx).
		Preload("Category").
		First(&budget, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("start_date DESC").
		Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *budgetRepository) FindByCategoryAndPeriod(ctx context.Context, userID, categoryID uuid.UUID, period string, startDate time.Time) (*models.Budget, error) {
	var budget models.Budget
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ? AND period = ? AND start_date = ? AND deleted_at IS NULL",
			userID, categoryID, period, startDate).
		First(&budget).Error; err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) Update(ctx context.Context, budget *models.Budget) error {
	return r.db.WithContext(ctx).Save(budget).Error
}

func (r *budgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Budget{}, "id = ?", id).Error
}

func (r *budgetRepository) GetProgress(ctx context.Context, budgetID uuid.UUID) (float64, float64, error) {
	var spent, budgetAmount float64
	
	// Get budget amount
	if err := r.db.WithContext(ctx).Model(&models.Budget{}).
		Where("id = ?", budgetID).
		Select("amount").
		Scan(&budgetAmount).Error; err != nil {
		return 0, 0, err
	}

	// Get spent amount
	if err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("category_id = (SELECT category_id FROM budgets WHERE id = ?)", budgetID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&spent).Error; err != nil {
		return 0, 0, err
	}

	return spent, budgetAmount, nil
}