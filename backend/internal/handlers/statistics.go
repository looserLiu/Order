package handlers

import (
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type StatisticsHandler struct {
	db *gorm.DB
}

func NewStatisticsHandler(db *gorm.DB) *StatisticsHandler {
	return &StatisticsHandler{db: db}
}

func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalAccounts int64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Count(&totalAccounts)

	var totalTransactions int64
	h.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&totalTransactions)

	var totalBudgets int64
	h.db.Model(&models.Budget{}).Where("user_id = ?", userID).Count(&totalBudgets)

	// This month
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	endOfMonth := time.Now()

	var thisMonthIncome, thisMonthExpense float64
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ? AND bill_date <= ?", userID, "income", startOfMonth, endOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthIncome)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ? AND bill_date <= ?", userID, "expense", startOfMonth, endOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthExpense)

	// Last month
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfMonth.AddDate(0, 0, -1)

	var lastMonthIncome, lastMonthExpense float64
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ? AND bill_date <= ?", userID, "income", startOfLastMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&lastMonthIncome)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ? AND bill_date <= ?", userID, "expense", startOfLastMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&lastMonthExpense)

	// Total balance
	var totalBalance float64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance)

	// Unread notifications
	var unreadNotifications int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadNotifications)

	// Income change percentage
	incomeChange := 0.0
	if lastMonthIncome > 0 {
		incomeChange = (thisMonthIncome - lastMonthIncome) / lastMonthIncome * 100
	}

	// Expense change percentage
	expenseChange := 0.0
	if lastMonthExpense > 0 {
		expenseChange = (thisMonthExpense - lastMonthExpense) / lastMonthExpense * 100
	}

	response.Success(c, gin.H{
		"total_accounts":        totalAccounts,
		"total_transactions":     totalTransactions,
		"total_budgets":         totalBudgets,
		"total_balance":         totalBalance,
		"this_month_income":     thisMonthIncome,
		"this_month_expense":    thisMonthExpense,
		"last_month_income":     lastMonthIncome,
		"last_month_expense":    lastMonthExpense,
		"income_change":         incomeChange,
		"expense_change":       expenseChange,
		"unread_notifications": unreadNotifications,
	})
}

// UploadHandler handles file uploads
