package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/ledger/backend/internal/config"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AuthHandler struct {
	userService *service.UserService
	cfg         *config.Config
}

func NewAuthHandler(userService *service.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{userService: userService, cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Email, req.Phone, req.Password, req.Nickname)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.generateToken(*user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// For login, we need to find user by email
	user, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.generateToken(*user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetMe(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	token, err := h.generateToken(*user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) generateToken(user models.User) (string, error) {
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}
