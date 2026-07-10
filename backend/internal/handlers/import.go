package handlers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type ImportHandler struct {
	db *gorm.DB
}

func NewImportHandler(db *gorm.DB) *ImportHandler {
	return &ImportHandler{db: db}
}

func (h *ImportHandler) ImportTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req []ImportTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var imported int
	var failed int

	for _, item := range req {
		billDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			failed++
			continue
		}

		var accountID uuid.UUID
		if item.Account != "" {
			var account models.Account
			if err := h.db.Where("user_id = ? AND name = ?", userID, item.Account).First(&account).Error; err == nil {
				accountID = account.ID
			} else {
				account := models.Account{
					UserID: userID,
					Name:   item.Account,
					Type:   "cash",
				}
				h.db.Create(&account)
				accountID = account.ID
			}
		}

		var categoryID *uuid.UUID
		if item.Category != "" {
			var category models.Category
			if err := h.db.Where("user_id = ? AND name = ?", userID, item.Category).First(&category).Error; err == nil {
				categoryID = &category.ID
			}
		}

		tx := models.Transaction{
			UserID:     userID,
			AccountID:  accountID,
			CategoryID: categoryID,
			Type:       item.Type,
			Amount:     item.Amount,
			Merchant:   item.Merchant,
			Note:       item.Note,
			BillDate:   billDate,
		}

		h.db.Create(&tx)
		imported++
	}

	response.Success(c, gin.H{
		"imported": imported,
		"failed":   failed,
	})
}
