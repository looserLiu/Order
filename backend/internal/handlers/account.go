package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AccountHandler struct {
	db *gorm.DB
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{db: db}
}

func (h *AccountHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var accounts []models.Account
	h.db.Where("user_id = ?", userID).Order("sort_order ASC").Find(&accounts)

	response.Success(c, accounts)
}

func (h *AccountHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account := models.Account{
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Balance:   req.Balance,
		Currency:  req.Currency,
		Icon:      req.Icon,
		Color:     req.Color,
		IsDefault: req.IsDefault,
		SortOrder: req.SortOrder,
	}

	h.db.Create(&account)
	response.Success(c, account)
}

func (h *AccountHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var account models.Account
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Account not found")
		return
	}

	response.Success(c, account)
}

func (h *AccountHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var account models.Account
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Account not found")
		return
	}

	account.Name = req.Name
	account.Type = req.Type
	account.Balance = req.Balance
	account.Currency = req.Currency
	account.Icon = req.Icon
	account.Color = req.Color
	account.IsDefault = req.IsDefault
	account.SortOrder = req.SortOrder

	h.db.Save(&account)
	response.Success(c, account)
}

func (h *AccountHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Account{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	response.SuccessWithMessage(c, "Account deleted", nil)
}

func (h *AccountHandler) GetTotalBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var total float64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Select("COALESCE(SUM(balance), 0)").Scan(&total)

	response.Success(c, gin.H{"total_balance": total})
}
