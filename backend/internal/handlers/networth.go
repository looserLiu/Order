package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type NetWorthHandler struct {
	db *gorm.DB
}

func NewNetWorthHandler(db *gorm.DB) *NetWorthHandler {
	return &NetWorthHandler{db: db}
}

func (h *NetWorthHandler) GetNetWorth(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalAssets float64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(balance), 0)").Scan(&totalAssets)

	var totalDebtOwed float64
	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owed", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalDebtOwed)

	var totalDebtOwing float64
	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owing", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalDebtOwing)

	var totalInvestment float64
	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "investment", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalInvestment)

	netWorth := totalAssets + totalDebtOwed - totalDebtOwing

	response.Success(c, gin.H{
		"total_assets":      totalAssets,
		"total_debt_owed":   totalDebtOwed,
		"total_debt_owing":  totalDebtOwing,
		"total_investment":  totalInvestment,
		"net_worth":         netWorth,
	})
}
