package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// FamilyRepository defines the interface for family data access
type FamilyRepository interface {
	Create(ctx context.Context, family *models.Family) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Family, error)
	FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]models.Family, error)
	AddMember(ctx context.Context, member *models.FamilyMember) error
	RemoveMember(ctx context.Context, familyID, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// familyRepository implements FamilyRepository
type familyRepository struct {
	db *gorm.DB
}

// NewFamilyRepository creates a new family repository
func NewFamilyRepository(db *gorm.DB) FamilyRepository {
	return &familyRepository{db: db}
}

func (r *familyRepository) Create(ctx context.Context, family *models.Family) error {
	return r.db.WithContext(ctx).Create(family).Error
}

func (r *familyRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Family, error) {
	var family models.Family
	if err := r.db.WithContext(ctx).
		Preload("Members", "deleted_at IS NULL").
		Preload("Members.User").
		First(&family, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &family, nil
}

func (r *familyRepository) FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]models.Family, error) {
	var families []models.Family
	if err := r.db.WithContext(ctx).
		Preload("Members", "deleted_at IS NULL").
		Where("owner_id = ? AND deleted_at IS NULL", ownerID).
		Order("created_at DESC").
		Find(&families).Error; err != nil {
		return nil, err
	}
	return families, nil
}

func (r *familyRepository) AddMember(ctx context.Context, member *models.FamilyMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *familyRepository) RemoveMember(ctx context.Context, familyID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.FamilyMember{},
		"family_id = ? AND user_id = ?", familyID, userID).Error
}

func (r *familyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Family{}, "id = ?", id).Error
}