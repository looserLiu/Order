package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo repository.UserRepository
	secret   string
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, secret string) *AuthService {
	return &AuthService{userRepo: userRepo, secret: secret}
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	User         *models.User `json:"user"`
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	var user *models.User
	var err error

	if req.Email != "" {
		user, err = s.userRepo.FindByEmail(ctx, req.Email)
	} else if req.Phone != "" {
		user, err = s.userRepo.FindByPhone(ctx, req.Phone)
	} else {
		return nil, errors.New("email or phone required")
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Register creates a new user
func (s *AuthService) Register(ctx context.Context, email, phone, password, nickname string) (*LoginResponse, error) {
	// Check if email exists
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Check if phone exists
	_, err = s.userRepo.FindByPhone(ctx, phone)
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

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken generates a new token from refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	
	_, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	}, jwt.WithClaims(claims))
	
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", errors.New("invalid token subject")
	}

	// Verify user exists
	_, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	return s.generateToken(userID)
}

func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}