package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type BudgetHandler struct {
	budgetService *service.BudgetService
}

func NewBudgetHandler(budgetService *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{budgetService: budgetService}
}

func (h *BudgetHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	budgets, err := h.budgetService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

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

	createReq := &service.BudgetCreateRequest{
		CategoryID:     req.CategoryID,
		Amount:         req.Amount,
		Period:         req.Period,
		StartDate:      startDate,
		EndDate:        endDate,
		AlertThreshold: req.AlertThreshold,
	}

	budget, err := h.budgetService.Create(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, budget)
}

func (h *BudgetHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	budget, err := h.budgetService.Get(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	response.Success(c, budget)
}

func (h *BudgetHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	updates["amount"] = req.Amount
	updates["category_id"] = req.CategoryID
	updates["period"] = req.Period
	updates["alert_threshold"] = req.AlertThreshold

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	updates["start_date"] = startDate

	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		updates["end_date"] = &t
	}

	budget, err := h.budgetService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, budget)
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.budgetService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Budget deleted", nil)
}

func (h *BudgetHandler) GetProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	budget, err := h.budgetService.Get(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	startDate := budget.StartDate
	endDate := time.Now()
	if budget.EndDate != nil {
		endDate = *budget.EndDate
	}

	// Calculate spent - simplified
	spent := 0.0

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
