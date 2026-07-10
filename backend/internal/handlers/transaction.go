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

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	categoryID := c.Query("category_id")
	accountID := c.Query("account_id")
	txType := c.Query("type")
	tagID := c.Query("tag_id")
	merchant := c.Query("merchant")

	var transactions []models.Transaction
	query := h.db.Where("user_id = ?", userID)

	if startDate != "" {
		query = query.Where("bill_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("bill_date <= ?", endDate)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if merchant != "" {
		query = query.Where("merchant ILIKE ?", "%"+merchant+"%")
	}
	if tagID != "" {
		query = query.Where("EXISTS (SELECT 1 FROM transaction_tags WHERE transaction_id = transactions.id AND tag_id = ?)", tagID)
	}

	var total int64
	query.Count(&total)

	query.Order("bill_date DESC, created_at DESC").Offset((parseInt(page) - 1) * parseInt(pageSize)).Limit(parseInt(pageSize)).Preload("Category").Preload("Account").Preload("Tags").Find(&transactions)

	response.Paginate(c, transactions, total, parseInt(page), parseInt(pageSize))
}

func (h *TransactionHandler) ListByDate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	date := c.Query("date")

	var transactions []models.Transaction
	h.db.Where("user_id = ? AND bill_date = ?", userID, date).
		Preload("Category").Preload("Account").
		Order("created_at DESC").
		Find(&transactions)

	response.Success(c, transactions)
}

func (h *TransactionHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	billDate, _ := time.Parse("2006-01-02", req.BillDate)

	transaction := models.Transaction{
		UserID:          userID,
		AccountID:       *req.AccountID,
		TargetAccountID: req.TargetAccountID,
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Merchant:        req.Merchant,
		Note:            req.Note,
		Tags:            req.Tags,
		BillDate:        billDate,
		IsRecurring:     req.IsRecurring,
		RecurringRule:   req.RecurringRule,
	}

	h.db.Create(&transaction)

	h.updateAccountBalance(req.AccountID, req.Type, req.Amount)

	if req.Type == "transfer" && req.TargetAccountID != nil {
		h.updateAccountBalance(req.TargetAccountID, "income", req.Amount)
	}

	response.Success(c, transaction)
}

func (h *TransactionHandler) updateAccountBalance(accountID *uuid.UUID, txType string, amount float64) {
	var account models.Account
	if err := h.db.First(&account, accountID).Error; err == nil {
		if txType == "expense" {
			account.Balance -= amount
		} else if txType == "income" {
			account.Balance += amount
		}
		h.db.Save(&account)
	}
}

func (h *TransactionHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Category").Preload("Account").Preload("Tags").First(&transaction).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	response.Success(c, transaction)
}

func (h *TransactionHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	billDate, _ := time.Parse("2006-01-02", req.BillDate)
	transaction.AccountID = *req.AccountID
	transaction.TargetAccountID = req.TargetAccountID
	transaction.CategoryID = req.CategoryID
	transaction.Type = req.Type
	transaction.Amount = req.Amount
	transaction.Currency = req.Currency
	transaction.Merchant = req.Merchant
	transaction.Note = req.Note
	transaction.Tags = req.Tags
	transaction.BillDate = billDate

	h.db.Save(&transaction)
	response.Success(c, transaction)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	h.db.Delete(&transaction)
	response.SuccessWithMessage(c, "Transaction deleted", nil)
}

func (h *TransactionHandler) BatchDelete(c *gin.Context) {
	userID := middleware.GetUserID(c)

	type BatchDeleteRequest struct {
		IDs []uuid.UUID `json:"ids" binding:"required"`
	}

	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	h.db.Where("id IN ? AND user_id = ?", req.IDs, userID).Delete(&models.Transaction{})
	response.SuccessWithMessage(c, "Transactions deleted", nil)
}
