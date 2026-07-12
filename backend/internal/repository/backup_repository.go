package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"gorm.io/gorm"
)

// BackupRepository defines the interface for backup data access
type BackupRepository interface {
	Create(ctx context.Context, backup *models.Backup) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Backup, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Backup, error)
}

// backupRepository implements BackupRepository
type backupRepository struct {
	db *gorm.DB
}

// NewBackupRepository creates a new backup repository
func NewBackupRepository(db *gorm.DB) BackupRepository {
	return &backupRepository{db: db}
}

func (r *backupRepository) Create(ctx context.Context, backup *models.Backup) error {
	return r.db.WithContext(ctx).Create(backup).Error
}

func (r *backupRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Backup, error) {
	var backup models.Backup
	if err := r.db.WithContext(ctx).First(&backup, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &backup, nil
}

func (r *backupRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Backup, error) {
	var backups []models.Backup
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&backups).Error; err != nil {
		return nil, err
	}
	return backups, nil
}