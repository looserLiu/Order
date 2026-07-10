package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type SearchHandler struct {
	db *gorm.DB
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

func (h *SearchHandler) Search(c *gin.Context) {
	userID := middleware.GetUserID(c)
	keyword := c.Query("q")
	searchType := c.DefaultQuery("type", "all")

	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "Keyword is required")
		return
	}

	limit := 20

	var transactions []models.Transaction
	if searchType == "all" || searchType == "transactions" {
		h.db.Where("user_id = ? AND (note ILIKE ? OR merchant ILIKE ?)", userID, "%"+keyword+"%", "%"+keyword+"%").
			Order("bill_date DESC").Limit(limit).Preload("Category").Preload("Account").Find(&transactions)
	}

	var accounts []models.Account
	if searchType == "all" || searchType == "accounts" {
		h.db.Where("user_id = ? AND name ILIKE ?", userID, "%"+keyword+"%").Limit(limit).Find(&accounts)
	}

	var categories []models.Category
	if searchType == "all" || searchType == "categories" {
		h.db.Where("user_id = ? AND name ILIKE ?", userID, "%"+keyword+"%").Limit(limit).Find(&categories)
	}

	response.Success(c, gin.H{
		"transactions": transactions,
		"accounts":     accounts,
		"categories":   categories,
	})
}
