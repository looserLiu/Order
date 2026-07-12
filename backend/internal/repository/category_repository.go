package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Category, error)
	GetTree(ctx context.Context, userID uuid.UUID) ([]models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// categoryRepository implements CategoryRepository
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("sort_order ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) GetTree(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.WithContext(ctx).Where("user_id = ? AND parent_id IS NULL AND deleted_at IS NULL", userID).
		Preload("Children", "deleted_at IS NULL").
		Order("sort_order ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}