package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type CSVImportHandler struct {
	csvImportService *service.CSVImportService
}

func NewCSVImportHandler(csvImportService *service.CSVImportService) *CSVImportHandler {
	return &CSVImportHandler{csvImportService: csvImportService}
}

func (h *CSVImportHandler) ImportCSV(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Transactions []service.CSVTransaction `json:"transactions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.csvImportService.ImportCSV(c.Request.Context(), userID, req.Transactions)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// Statistics Handler for dashboard

type StatisticsHandler struct {
	db *gorm.DB
}

func NewStatisticsHandler(db *gorm.DB) *StatisticsHandler {
	return &StatisticsHandler{db: db}
}

func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalIncome, totalExpense float64
	h.db.Model(&models.Transaction{}).Where("user_id = ? AND type = ?", userID, "income").Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	h.db.Model(&models.Transaction{}).Where("user_id = ? AND type = ?", userID, "expense").Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	var accountCount int64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Count(&accountCount)

	var categoryCount int64
	h.db.Model(&models.Category{}).Where("user_id = ?", userID).Count(&categoryCount)

	var budgetCount int64
	h.db.Model(&models.Budget{}).Where("user_id = ?", userID).Count(&budgetCount)

	response.Success(c, gin.H{
		"total_income":   totalIncome,
		"total_expense":  totalExpense,
		"account_count":  accountCount,
		"category_count": categoryCount,
		"budget_count":   budgetCount,
	})
}
