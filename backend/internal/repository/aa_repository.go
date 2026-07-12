package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// AARepository defines the interface for AA group data access
type AARepository interface {
	CreateGroup(ctx context.Context, group *models.AAGroup) error
	FindGroupByID(ctx context.Context, id uuid.UUID) (*models.AAGroup, error)
	FindGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]models.AAGroup, error)
	AddMember(ctx context.Context, member *models.AAMember) error
	AddExpense(ctx context.Context, groupID uuid.UUID, amount float64) error
	CreateSettlement(ctx context.Context, settlement *models.AASettlement) error
	GetSettlements(ctx context.Context, groupID uuid.UUID) ([]models.AASettlement, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) error
}

// aaRepository implements AARepository
type aaRepository struct {
	db *gorm.DB
}

// NewAARepository creates a new AA repository
func NewAARepository(db *gorm.DB) AARepository {
	return &aaRepository{db: db}
}

func (r *aaRepository) CreateGroup(ctx context.Context, group *models.AAGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *aaRepository) FindGroupByID(ctx context.Context, id uuid.UUID) (*models.AAGroup, error) {
	var group models.AAGroup
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Preload("Settlements").
		First(&group, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *aaRepository) FindGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]models.AAGroup, error) {
	var groups []models.AAGroup
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *aaRepository) AddMember(ctx context.Context, member *models.AAMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *aaRepository) AddExpense(ctx context.Context, groupID uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&models.AAGroup{}).
		Where("id = ?", groupID).
		UpdateColumn("total_amount", gorm.Expr("total_amount + ?", amount)).Error
}

func (r *aaRepository) CreateSettlement(ctx context.Context, settlement *models.AASettlement) error {
	return r.db.WithContext(ctx).Create(settlement).Error
}

func (r *aaRepository) GetSettlements(ctx context.Context, groupID uuid.UUID) ([]models.AASettlement, error) {
	var settlements []models.AASettlement
	if err := r.db.WithContext(ctx).
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Find(&settlements).Error; err != nil {
		return nil, err
	}
	return settlements, nil
}

func (r *aaRepository) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.AAGroup{}, "id = ?", id).Error
}