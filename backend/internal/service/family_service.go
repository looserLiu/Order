package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
)

// FamilyService handles family business logic
type FamilyService struct {
	repo repository.FamilyRepository
}

// NewFamilyService creates a new family service
func NewFamilyService(repo repository.FamilyRepository) *FamilyService {
	return &FamilyService{repo: repo}
}

// CreateRequest represents family creation data
type FamilyCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

// Create creates a new family
func (s *FamilyService) Create(ctx context.Context, ownerID uuid.UUID, req *FamilyCreateRequest) (*models.Family, error) {
	family := &models.Family{
		Name:     req.Name,
		OwnerID:  ownerID,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, family); err != nil {
		return nil, err
	}

	return family, nil
}

// List retrieves all families for a user
func (s *FamilyService) List(ctx context.Context, ownerID uuid.UUID) ([]models.Family, error) {
	return s.repo.FindByOwnerID(ctx, ownerID)
}

// Get retrieves a single family
func (s *FamilyService) Get(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*models.Family, error) {
	family, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if family.OwnerID != ownerID {
		return nil, errors.New("unauthorized access")
	}
	return family, nil
}

// AddMember adds a member to a family
func (s *FamilyService) AddMember(ctx context.Context, familyID, userID uuid.UUID) error {
	member := &models.FamilyMember{
		FamilyID:  familyID,
		UserID:    userID,
		Role:      "member",
		JoinedAt:  time.Now(),
	}
	return s.repo.AddMember(ctx, member)
}

// RemoveMember removes a member from a family
func (s *FamilyService) RemoveMember(ctx context.Context, familyID, userID uuid.UUID) error {
	return s.repo.RemoveMember(ctx, familyID, userID)
}

// Delete removes a family
func (s *FamilyService) Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	family, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if family.OwnerID != ownerID {
		return errors.New("unauthorized access")
	}
	return s.repo.Delete(ctx, id)
}