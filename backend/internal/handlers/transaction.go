package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	limit := parseInt(pageSize)
	offset := (parseInt(page) - 1) * limit

	transactions, err := h.transactionService.List(c.Request.Context(), userID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var total int64 = int64(len(transactions))
	response.Paginate(c, transactions, total, parseInt(page), limit)
}

func (h *TransactionHandler) ListByDate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	date := c.Query("date")

	if date != "" {
		billDate, _ := time.Parse("2006-01-02", date)
		txs, _ := h.transactionService.ListByDate(c.Request.Context(), userID, billDate, billDate)
		response.Success(c, txs)
		return
	}

	response.Success(c, []models.Transaction{})
}

func (h *TransactionHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	billDate, _ := time.Parse("2006-01-02", req.BillDate)

	createReq := &service.CreateRequest{
		AccountID:       *req.AccountID,
		TargetAccountID: req.TargetAccountID,
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Merchant:        req.Merchant,
		Note:            req.Note,
		BillDate:        billDate,
		IsRecurring:     req.IsRecurring,
		RecurringRule:   req.RecurringRule,
	}

	transaction, err := h.transactionService.Create(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, transaction)
}

func (h *TransactionHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	transaction, err := h.transactionService.Get(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	response.Success(c, transaction)
}

func (h *TransactionHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.AccountID != nil {
		updates["account_id"] = *req.AccountID
	}
	if req.CategoryID != nil {
		updates["category_id"] = req.CategoryID
	}
	updates["amount"] = req.Amount
	updates["merchant"] = req.Merchant
	updates["note"] = req.Note
	billDate, _ := time.Parse("2006-01-02", req.BillDate)
	updates["bill_date"] = billDate

	transaction, err := h.transactionService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, transaction)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.transactionService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

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

	if err := h.transactionService.BatchDelete(c.Request.Context(), req.IDs, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Transactions deleted", nil)
}
