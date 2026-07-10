package handlers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AssetHandler struct {
	db *gorm.DB
}

func NewAssetHandler(db *gorm.DB) *AssetHandler {
	return &AssetHandler{db: db}
}

func (h *AssetHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	assetType := c.Query("type")

	var assets []models.AssetChange
	query := h.db.Where("user_id = ?", userID)
	if assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}
	query.Order("created_at DESC").Find(&assets)

	response.Success(c, assets)
}

func (h *AssetHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var startDate, endDate *time.Time
	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		startDate = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		endDate = &t
	}

	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	asset := models.AssetChange{
		UserID:       userID,
		AssetType:    req.AssetType,
		RelatedUser:  req.RelatedUser,
		Name:         req.Name,
		Amount:       req.Amount,
		InterestRate: req.InterestRate,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		Note:         req.Note,
	}

	h.db.Create(&asset)
	response.Success(c, asset)
}

func (h *AssetHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var asset models.AssetChange
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&asset).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Asset not found")
		return
	}

	response.Success(c, asset)
}

func (h *AssetHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var asset models.AssetChange
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&asset).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Asset not found")
		return
	}

	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		asset.StartDate = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		asset.EndDate = &t
	}

	asset.AssetType = req.AssetType
	asset.RelatedUser = req.RelatedUser
	asset.Name = req.Name
	asset.Amount = req.Amount
	asset.InterestRate = req.InterestRate
	asset.Status = req.Status
	asset.Note = req.Note

	h.db.Save(&asset)
	response.Success(c, asset)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.AssetChange{})
	response.SuccessWithMessage(c, "Asset deleted", nil)
}

func (h *AssetHandler) GetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var debtOwed, debtOwing, investment float64

	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owed", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&debtOwed)

	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owing", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&debtOwing)

	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "investment", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&investment)

	response.Success(c, gin.H{
		"debt_owed":   debtOwed,
		"debt_owing":  debtOwing,
		"investment":  investment,
		"net_worth":   debtOwed - debtOwing + investment,
	})
}
