package handlers

import (
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalAccounts int64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Count(&totalAccounts)

	var totalTransactions int64
	h.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&totalTransactions)

	var totalBudgets int64
	h.db.Model(&models.Budget{}).Where("user_id = ?", userID).Count(&totalBudgets)

	var thisMonthIncome, thisMonthExpense float64
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ?", userID, "income", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthIncome)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ?", userID, "expense", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthExpense)

	var unreadNotifications int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadNotifications)

	response.Success(c, gin.H{
		"total_accounts":      totalAccounts,
		"total_transactions":  totalTransactions,
		"total_budgets":       totalBudgets,
		"this_month_income":   thisMonthIncome,
		"this_month_expense":  thisMonthExpense,
		"unread_notifications": unreadNotifications,
	})
}
