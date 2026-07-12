package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// BackupService handles backup business logic
type BackupService struct {
	repo repository.BackupRepository
}

// NewBackupService creates a new backup service
func NewBackupService(repo repository.BackupRepository) *BackupService {
	return &BackupService{repo: repo}
}

// ExportData represents exported backup data
type ExportData struct {
	Version     string                `json:"version"`
	ExportedAt  string                `json:"exported_at"`
	Accounts    []models.Account      `json:"accounts"`
	Categories  []models.Category     `json:"categories"`
	Transactions []models.Transaction `json:"transactions"`
	Budgets     []models.Budget       `json:"budgets"`
	Tags        []models.Tag          `json:"tags"`
	Reminders   []models.Reminder     `json:"reminders"`
	Goals       []models.FinancialGoal `json:"goals"`
	Insurances  []models.Insurance    `json:"insurances"`
	AssetChanges []models.AssetChange  `json:"asset_changes"`
}

// ImportData represents imported backup data
type ImportData struct {
	Accounts     []models.Account      `json:"accounts"`
	Categories   []models.Category     `json:"categories"`
	Transactions []models.Transaction  `json:"transactions"`
	Budgets      []models.Budget       `json:"budgets"`
	Tags         []models.Tag          `json:"tags"`
	Reminders    []models.Reminder     `json:"reminders"`
	Goals        []models.FinancialGoal `json:"goals"`
	Insurances   []models.Insurance    `json:"insurances"`
	AssetChanges []models.AssetChange   `json:"asset_changes"`
}

// ExportAll exports all user data
func (s *BackupService) ExportAll(ctx context.Context, userID uuid.UUID) (*ExportData, error) {
	// This would typically use other repositories to fetch data
	// For now, return empty data
	return &ExportData{
		Version:    "1.0",
		ExportedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// ImportAll imports user data
func (s *BackupService) ImportAll(ctx context.Context, userID uuid.UUID, data *ImportData) (map[string]int, error) {
	imported := make(map[string]int)
	return imported, nil
}

// List retrieves all backups for a user
func (s *BackupService) List(ctx context.Context, userID uuid.UUID) ([]models.Backup, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Create creates a backup record
func (s *BackupService) Create(ctx context.Context, userID uuid.UUID, backupType, fileName string, fileSize int64) (*models.Backup, error) {
	backup := &models.Backup{
		UserID:     userID,
		FileName:   fileName,
		FileSize:   fileSize,
		BackupType: backupType,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, backup); err != nil {
		return nil, err
	}

	return backup, nil
}

// GetByID retrieves a backup by ID
func (s *BackupService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Backup, error) {
	backup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if backup.UserID != userID {
		return nil, errors.New("unauthorized access")
	}
	return backup, nil
}