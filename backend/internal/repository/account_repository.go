package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// AccountRepository defines the interface for account data access
type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalBalance(ctx context.Context, userID uuid.UUID) (float64, error)
}

// accountRepository implements AccountRepository
type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("sort_order ASC").
		Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) Update(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Account{}, "id = ?", id).Error
}

func (r *accountRepository) GetTotalBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total float64
	if err := r.db.WithContext(ctx).Model(&models.Account{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}