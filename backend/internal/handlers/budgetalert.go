package handlers

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/response"
)

type BudgetAlertHandler struct {
	db *gorm.DB
}

func NewBudgetAlertHandler(db *gorm.DB) *BudgetAlertHandler {
	return &BudgetAlertHandler{db: db}
}

func (h *BudgetAlertHandler) GetAlerts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		response.Error(c, 401, "Unauthorized")
		return
	}

	uid := userID

	// Get current month date range
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// Get all budgets
	var budgets []models.Budget
	if err := h.db.Where("user_id = ?", uid).Find(&budgets).Error; err != nil {
		response.Error(c, 500, "Failed to get budgets")
		return
	}

	alerts := make([]BudgetAlert, 0)

	for _, budget := range budgets {
		// Calculate spent amount
		var spent float64
		query := h.db.Model(&models.Transaction{}).
			Where("user_id = ? AND type = ? AND bill_date BETWEEN ? AND ?",
				uid, "expense", startOfMonth, endOfMonth)

		if budget.CategoryID != nil {
			query = query.Where("category_id = ?", budget.CategoryID)
		}

		if err := query.Select("COALESCE(SUM(amount * exchange_rate), 0)").Scan(&spent).Error; err != nil {
			continue
		}

		percentage := 0.0
		if budget.Amount > 0 {
			percentage = (spent / budget.Amount) * 100
		}

		// Check if alert needed
		alertType := ""
		if percentage >= 100 {
			alertType = "exceeded"
		} else if percentage >= budget.AlertThreshold*100 {
			alertType = "warning"
		}

		if alertType != "" {
			// Get category name
			categoryName := "所有类别"
			if budget.CategoryID != nil {
				var cat models.Category
				if err := h.db.First(&cat, budget.CategoryID).Error; err == nil {
					categoryName = cat.Name
				}
			}

			alerts = append(alerts, BudgetAlert{
				BudgetID:     budget.ID.String(),
				CategoryName: categoryName,
				BudgetAmount: budget.Amount,
				SpentAmount:  spent,
				Remaining:    budget.Amount - spent,
				AlertType:    alertType,
				Percentage:   percentage,
			})
		}
	}

	response.Success(c, alerts)
}
