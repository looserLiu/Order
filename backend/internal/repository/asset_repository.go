package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// AssetRepository defines the interface for asset data access
type AssetRepository interface {
	Create(ctx context.Context, asset *models.AssetChange) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.AssetChange, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.AssetChange, error)
	Update(ctx context.Context, asset *models.AssetChange) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSummary(ctx context.Context, userID uuid.UUID) (map[string]float64, error)
}

// assetRepository implements AssetRepository
type assetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) Create(ctx context.Context, asset *models.AssetChange) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *assetRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.AssetChange, error) {
	var asset models.AssetChange
	if err := r.db.WithContext(ctx).First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.AssetChange, error) {
	var assets []models.AssetChange
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *assetRepository) Update(ctx context.Context, asset *models.AssetChange) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *assetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.AssetChange{}, "id = ?", id).Error
}

func (r *assetRepository) GetSummary(ctx context.Context, userID uuid.UUID) (map[string]float64, error) {
	var results []struct {
		AssetType string
		Total     float64
	}

	if err := r.db.WithContext(ctx).Model(&models.AssetChange{}).
		Select("asset_type, SUM(amount) as total").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group("asset_type").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	summary := make(map[string]float64)
	for _, r := range results {
		summary[r.AssetType] = r.Total
	}
	return summary, nil
}