package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user business logic
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetMe retrieves the current user
func (s *UserService) GetMe(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(ctx, userID)
}

// UpdateMe updates user profile
func (s *UserService) UpdateMe(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if nickname, ok := updates["nickname"].(string); ok {
		user.Nickname = nickname
	}
	if avatarURL, ok := updates["avatar_url"].(string); ok {
		user.AvatarURL = avatarURL
	}
	if currency, ok := updates["currency"].(string); ok {
		user.Currency = currency
	}
	if timezone, ok := updates["timezone"].(string); ok {
		user.Timezone = timezone
	}

	return s.repo.Update(ctx, user)
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// Login authenticates a user
func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// Register creates a new user
func (s *UserService) Register(ctx context.Context, email, phone, password, nickname string) (*models.User, error) {
	// Check if email exists
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Check if phone exists
	_, err = s.repo.FindByPhone(ctx, phone)
	if err == nil {
		return nil, errors.New("phone already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Nickname:     nickname,
		Currency:     "CNY",
		Timezone:     "Asia/Shanghai",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}