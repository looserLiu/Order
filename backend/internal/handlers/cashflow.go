package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/response"
)

type CashFlowHandler struct {
	cashFlowService *service.CashFlowService
}

func NewCashFlowHandler(cashFlowService *service.CashFlowService) *CashFlowHandler {
	return &CashFlowHandler{cashFlowService: cashFlowService}
}

func (h *CashFlowHandler) GetProjection(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		response.Error(c, 401, "Unauthorized")
		return
	}

	uid := userID

	// Get days parameter (default 30)
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	result, err := h.cashFlowService.GetProjection(c.Request.Context(), uid, days)
	if err != nil {
		response.Error(c, 500, "Failed to get projection")
		return
	}

	response.Success(c, gin.H{
		"current_balance": result.CurrentBalance,
		"projections":     result.Projections,
	})
}

// shouldRecurToday checks if a transaction should recur on the given date
func shouldRecurToday(lastDate *time.Time, today time.Time) bool {
	if lastDate == nil {
		return false
	}
	return today.After(*lastDate) && today.Day() == lastDate.Day()
}
