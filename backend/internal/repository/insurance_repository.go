package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// InsuranceRepository defines the interface for insurance data access
type InsuranceRepository interface {
	Create(ctx context.Context, insurance *models.Insurance) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Insurance, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Insurance, error)
	Update(ctx context.Context, insurance *models.Insurance) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSummary(ctx context.Context, userID uuid.UUID) (map[string]float64, error)
}

// insuranceRepository implements InsuranceRepository
type insuranceRepository struct {
	db *gorm.DB
}

// NewInsuranceRepository creates a new insurance repository
func NewInsuranceRepository(db *gorm.DB) InsuranceRepository {
	return &insuranceRepository{db: db}
}

func (r *insuranceRepository) Create(ctx context.Context, insurance *models.Insurance) error {
	return r.db.WithContext(ctx).Create(insurance).Error
}

func (r *insuranceRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Insurance, error) {
	var insurance models.Insurance
	if err := r.db.WithContext(ctx).First(&insurance, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &insurance, nil
}

func (r *insuranceRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Insurance, error) {
	var insurances []models.Insurance
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&insurances).Error; err != nil {
		return nil, err
	}
	return insurances, nil
}

func (r *insuranceRepository) Update(ctx context.Context, insurance *models.Insurance) error {
	return r.db.WithContext(ctx).Save(insurance).Error
}

func (r *insuranceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Insurance{}, "id = ?", id).Error
}

func (r *insuranceRepository) GetSummary(ctx context.Context, userID uuid.UUID) (map[string]float64, error) {
	var results []struct {
		InsuranceType string
		Total         float64
	}

	if err := r.db.WithContext(ctx).Model(&models.Insurance{}).
		Select("insurance_type, SUM(premium) as total").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group("insurance_type").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	summary := make(map[string]float64)
	for _, r := range results {
		summary[r.InsuranceType] = r.Total
	}
	return summary, nil
}