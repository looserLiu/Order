package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error)
	FindByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.Transaction, error)
	Update(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	BatchDelete(ctx context.Context, tx *gorm.DB, ids []uuid.UUID) error
	UpdateAccountBalance(ctx context.Context, tx *gorm.DB, accountID uuid.UUID, amount float64) error
}

// transactionRepository implements TransactionRepository
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	return tx.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("bill_date DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ? AND bill_date BETWEEN ? AND ? AND deleted_at IS NULL", userID, startDate, endDate).
		Order("bill_date ASC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *gorm.DB, transaction *models.Transaction) error {
	return tx.WithContext(ctx).Save(transaction).Error
}

func (r *transactionRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	return tx.WithContext(ctx).Delete(&models.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) BatchDelete(ctx context.Context, tx *gorm.DB, ids []uuid.UUID) error {
	return tx.WithContext(ctx).Delete(&models.Transaction{}, "id IN (?)", ids).Error
}

func (r *transactionRepository) UpdateAccountBalance(ctx context.Context, tx *gorm.DB, accountID uuid.UUID, amount float64) error {
	return tx.WithContext(ctx).Model(&models.Account{}).
		Where("id = ?", accountID).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}