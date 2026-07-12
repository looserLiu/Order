package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// InsuranceService handles insurance business logic
type InsuranceService struct {
	repo repository.InsuranceRepository
}

// NewInsuranceService creates a new insurance service
func NewInsuranceService(repo repository.InsuranceRepository) *InsuranceService {
	return &InsuranceService{repo: repo}
}

// CreateRequest represents insurance creation data
type InsuranceCreateRequest struct {
	Name         string     `json:"name" binding:"required"`
	InsuranceType string    `json:"insurance_type"`
	Company      string     `json:"company"`
	Premium      float64    `json:"premium"`
	PaymentType  string     `json:"payment_type"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Coverage     float64    `json:"coverage"`
	Beneficiary  string     `json:"beneficiary"`
	Note         string     `json:"note"`
}

// Create creates a new insurance
func (s *InsuranceService) Create(ctx context.Context, userID uuid.UUID, req *InsuranceCreateRequest) (*models.Insurance, error) {
	insurance := &models.Insurance{
		UserID:        userID,
		Name:          req.Name,
		InsuranceType: req.InsuranceType,
		Company:       req.Company,
		Premium:       req.Premium,
		PaymentType:   req.PaymentType,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Coverage:      req.Coverage,
		Beneficiary:   req.Beneficiary,
		Note:          req.Note,
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, insurance); err != nil {
		return nil, err
	}

	return insurance, nil
}

// List retrieves all insurances for a user
func (s *InsuranceService) List(ctx context.Context, userID uuid.UUID) ([]models.Insurance, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// Update updates an insurance
func (s *InsuranceService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.Insurance, error) {
	insurance, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if insurance.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	if name, ok := updates["name"].(string); ok {
		insurance.Name = name
	}
	if company, ok := updates["company"].(string); ok {
		insurance.Company = company
	}
	if premium, ok := updates["premium"].(float64); ok {
		insurance.Premium = premium
	}
	if status, ok := updates["status"].(string); ok {
		insurance.Status = status
	}
	if note, ok := updates["note"].(string); ok {
		insurance.Note = note
	}

	insurance.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, insurance); err != nil {
		return nil, err
	}

	return insurance, nil
}

// Delete removes an insurance
func (s *InsuranceService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	insurance, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if insurance.UserID != userID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}

// GetSummary gets insurance summary
func (s *InsuranceService) GetSummary(ctx context.Context, userID uuid.UUID) (map[string]float64, error) {
	return s.repo.GetSummary(ctx, userID)
}