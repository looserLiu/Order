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

type BudgetHandler struct {
	db *gorm.DB
}

func NewBudgetHandler(db *gorm.DB) *BudgetHandler {
	return &BudgetHandler{db: db}
}

func (h *BudgetHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var budgets []models.Budget
	h.db.Where("user_id = ?", userID).Preload("Category").Order("start_date DESC").Find(&budgets)

	response.Success(c, budgets)
}

func (h *BudgetHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateBudgetRequest
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

	budget := models.Budget{
		UserID:          userID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Period:          req.Period,
		StartDate:       startDate,
		EndDate:         endDate,
		AlertThreshold:  req.AlertThreshold,
	}

	h.db.Create(&budget)
	response.Success(c, budget)
}

func (h *BudgetHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var budget models.Budget
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Category").First(&budget).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	response.Success(c, budget)
}

func (h *BudgetHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var budget models.Budget
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&budget).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		budget.EndDate = &t
	}

	budget.CategoryID = req.CategoryID
	budget.Amount = req.Amount
	budget.Period = req.Period
	budget.StartDate = startDate
	budget.AlertThreshold = req.AlertThreshold

	h.db.Save(&budget)
	response.Success(c, budget)
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Budget{})
	response.SuccessWithMessage(c, "Budget deleted", nil)
}

func (h *BudgetHandler) GetProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var budget models.Budget
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Category").First(&budget).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	startDate := budget.StartDate
	endDate := time.Now()
	if budget.EndDate != nil {
		endDate = *budget.EndDate
	}

	var spent float64
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date BETWEEN ? AND ?",
			userID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&spent)

	if budget.CategoryID != nil {
		h.db.Model(&models.Transaction{}).
			Where("user_id = ? AND type = ? AND category_id = ? AND bill_date BETWEEN ? AND ?",
				userID, "expense", budget.CategoryID, startDate, endDate).
			Select("COALESCE(SUM(amount), 0)").Scan(&spent)
	}

	progress := 0.0
	if budget.Amount > 0 {
		progress = spent / budget.Amount * 100
	}

	alert := progress >= budget.AlertThreshold*100

	response.Success(c, gin.H{
		"budget":    budget,
		"spent":     spent,
		"remaining": budget.Amount - spent,
		"progress":  progress,
		"alert":     alert,
	})
}
