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

type FamilyTransactionHandler struct {
	db *gorm.DB
}

func NewFamilyTransactionHandler(db *gorm.DB) *FamilyTransactionHandler {
	return &FamilyTransactionHandler{db: db}
}

func (h *FamilyTransactionHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	familyID := c.Param("id")

	type CreateFamilyTransactionRequest struct {
		AccountID   *uuid.UUID `json:"account_id" binding:"required"`
		CategoryID  *uuid.UUID `json:"category_id"`
		Type        string     `json:"type" binding:"required"`
		Amount      float64    `json:"amount" binding:"required"`
		Note        string     `json:"note"`
		BillDate    string     `json:"bill_date" binding:"required"`
	}

	var req CreateFamilyTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fID, _ := uuid.Parse(familyID)
	billDate, _ := time.Parse("2006-01-02", req.BillDate)

	tx := models.FamilyTransaction{
		FamilyID:   fID,
		UserID:     userID,
		AccountID:  *req.AccountID,
		CategoryID: req.CategoryID,
		Type:       req.Type,
		Amount:     req.Amount,
		Note:       req.Note,
		BillDate:   billDate,
	}

	h.db.Create(&tx)
	response.Success(c, tx)
}
