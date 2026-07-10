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

type InsuranceHandler struct {
	db *gorm.DB
}

func NewInsuranceHandler(db *gorm.DB) *InsuranceHandler {
	return &InsuranceHandler{db: db}
}

func (h *InsuranceHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var insurances []models.Insurance
	h.db.Where("user_id = ?", userID).Order("next_payment_date ASC").Find(&insurances)

	response.Success(c, insurances)
}

func (h *InsuranceHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateInsuranceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	var endDate *time.Time
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		endDate = &t
	}

	var nextPaymentDate *time.Time
	if req.PaymentType == "yearly" {
		nextPaymentDate = &time.Time{}
		*nextPaymentDate = startDate.AddDate(1, 0, 0)
	} else if req.PaymentType == "quarterly" {
		nextPaymentDate = &time.Time{}
		*nextPaymentDate = startDate.AddDate(0, 3, 0)
	} else if req.PaymentType == "monthly" {
		nextPaymentDate = &time.Time{}
		*nextPaymentDate = startDate.AddDate(0, 1, 0)
	}

	insurance := models.Insurance{
		UserID:           userID,
		Name:             req.Name,
		InsuranceType:    req.InsuranceType,
		Company:          req.Company,
		Premium:          req.Premium,
		PaymentType:      req.PaymentType,
		StartDate:        startDate,
		EndDate:          endDate,
		Coverage:         req.Coverage,
		Beneficiary:      req.Beneficiary,
		Note:             req.Note,
		NextPaymentDate:  nextPaymentDate,
	}

	h.db.Create(&insurance)
	response.Success(c, insurance)
}

func (h *InsuranceHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateInsuranceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var insurance models.Insurance
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&insurance).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Insurance not found")
		return
	}

	insurance.Name = req.Name
	insurance.InsuranceType = req.InsuranceType
	insurance.Company = req.Company
	insurance.Premium = req.Premium
	insurance.PaymentType = req.PaymentType
	insurance.Coverage = req.Coverage
	insurance.Beneficiary = req.Beneficiary
	insurance.Note = req.Note

	h.db.Save(&insurance)
	response.Success(c, insurance)
}

func (h *InsuranceHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Insurance{})
	response.SuccessWithMessage(c, "Insurance deleted", nil)
}

func (h *InsuranceHandler) GetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalPremium float64
	var activeCount int64
	var expiringCount int64

	h.db.Model(&models.Insurance{}).Where("user_id = ? AND status = ?", userID, "active").
		Select("COALESCE(SUM(premium), 0)").Scan(&totalPremium)
	h.db.Model(&models.Insurance{}).Where("user_id = ? AND status = ?", userID, "active").Count(&activeCount)

	thirtyDaysLater := time.Now().AddDate(0, 0, 30)
	h.db.Model(&models.Insurance{}).
		Where("user_id = ? AND status = ? AND next_payment_date <= ?", userID, "active", thirtyDaysLater).
		Count(&expiringCount)

	var totalCoverage float64
	h.db.Model(&models.Insurance{}).Where("user_id = ? AND status = ?", userID, "active").
		Select("COALESCE(SUM(coverage), 0)").Scan(&totalCoverage)

	response.Success(c, gin.H{
		"total_premium":   totalPremium,
		"active_count":    activeCount,
		"expiring_count":  expiringCount,
		"total_coverage":  totalCoverage,
	})
}
