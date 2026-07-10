package handlers

import (
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/response"
)

type CashFlowHandler struct {
	db *gorm.DB
}

func NewCashFlowHandler(db *gorm.DB) *CashFlowHandler {
	return &CashFlowHandler{db: db}
}

func (h *CashFlowHandler) GetProjection(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, 401, "Unauthorized")
		return
	}

	uid, _ := uuid.Parse(userID)

	// Get days parameter (default 30)
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	// Get current total balance
	var totalBalance float64
	if err := h.db.Model(&models.Account{}).Where("user_id = ?", uid).Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance).Error; err != nil {
		response.Error(c, 500, "Failed to get balance")
		return
	}

	// Get recurring transactions
	var recurringTxs []models.Transaction
	if err := h.db.Where("user_id = ? AND is_recurring = ?", uid, true).Find(&recurringTxs).Error; err != nil {
		response.Error(c, 500, "Failed to get recurring transactions")
		return
	}

	// Calculate projections
	projections := make([]CashFlowProjection, days)
	currentBalance := totalBalance

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		dayIncome := 0.0
		dayExpense := 0.0
		recurringCount := 0

		for _, tx := range recurringTxs {
			if tx.BillDate.Day() == date.Day() || shouldRecurToday(&tx.BillDate, date) {
				if tx.Type == "income" {
					dayIncome += tx.Amount
				} else if tx.Type == "expense" {
					dayExpense += tx.Amount
				}
				recurringCount++
			}
		}

		currentBalance = currentBalance + dayIncome - dayExpense
		projections[i] = CashFlowProjection{
			Date:          dateStr,
			ProjectedBal:  currentBalance,
			Income:        dayIncome,
			Expense:       dayExpense,
			RecurringTx:   recurringCount,
		}
	}

	response.Success(c, gin.H{
		"current_balance": totalBalance,
		"projections":     projections,
	})
}
