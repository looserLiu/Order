package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// AssetService handles asset business logic
type AssetService struct {
	repo repository.AssetRepository
}

// NewAssetService creates a new asset service
func NewAssetService(repo repository.AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

// CreateRequest represents asset creation data
type AssetCreateRequest struct {
	AssetType   string     `json:"asset_type" binding:"required"`
	RelatedUser string     `json:"related_user"`
	Name        string     `json:"name" binding:"required"`
	Amount      float64    `json:"amount" binding:"required"`
	InterestRate float64   `json:"interest_rate"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Note        string     `json:"note"`
}

// Create creates a new asset
func (s *AssetService) Create(ctx context.Context, userID uuid.UUID, req *AssetCreateRequest) (*models.AssetChange, error) {
	asset := &models.AssetChange{
		UserID:       userID,
		AssetType:    req.AssetType,
		RelatedUser:  req.RelatedUser,
		Name:         req.Name,
		Amount:       req.Amount,
		InterestRate:  req.InterestRate,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Status:       "active",
		Note:         req.Note,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// List retrieves all assets for a user
func (s *AssetService) List(ctx context.Context, userID uuid.UUID) ([]models.AssetChange, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Get retrieves a single asset
func (s *AssetService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.AssetChange, error) {
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset.UserID != userID {
		return nil, errors.New("unauthorized access")
	}
	return asset, nil
}

// Update updates an asset
func (s *AssetService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.AssetChange, error) {
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if name, ok := updates["name"].(string); ok {
		asset.Name = name
	}
	if amount, ok := updates["amount"].(float64); ok {
		asset.Amount = amount
	}
	if status, ok := updates["status"].(string); ok {
		asset.Status = status
	}
	if note, ok := updates["note"].(string); ok {
		asset.Note = note
	}

	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// Delete removes an asset
func (s *AssetService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if asset.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}

// GetSummary gets asset summary
func (s *AssetService) GetSummary(ctx context.Context, userID uuid.UUID) (map[string]float64, error) {
	return s.repo.GetSummary(ctx, userID)
}