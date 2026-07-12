package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// GoalRepository defines the interface for financial goal data access
type GoalRepository interface {
	Create(ctx context.Context, goal *models.FinancialGoal) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.FinancialGoal, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.FinancialGoal, error)
	Update(ctx context.Context, goal *models.FinancialGoal) error
	AddAmount(ctx context.Context, id uuid.UUID, amount float64) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// goalRepository implements GoalRepository
type goalRepository struct {
	db *gorm.DB
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) Create(ctx context.Context, goal *models.FinancialGoal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

func (r *goalRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FinancialGoal, error) {
	var goal models.FinancialGoal
	if err := r.db.WithContext(ctx).First(&goal, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *goalRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.FinancialGoal, error) {
	var goals []models.FinancialGoal
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *goalRepository) Update(ctx context.Context, goal *models.FinancialGoal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

func (r *goalRepository) AddAmount(ctx context.Context, id uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&models.FinancialGoal{}).
		Where("id = ?", id).
		UpdateColumn("current_amount", gorm.Expr("current_amount + ?", amount)).Error
}

func (r *goalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.FinancialGoal{}, "id = ?", id).Error
}