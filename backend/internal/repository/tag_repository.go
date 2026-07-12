package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// TagRepository defines the interface for tag data access
type TagRepository interface {
	Create(ctx context.Context, tag *models.Tag) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Tag, error)
	Update(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// tagRepository implements TagRepository
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.WithContext(ctx).First(&tag, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Tag{}, "id = ?", id).Error
}