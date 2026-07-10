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

type GoalHandler struct {
	db *gorm.DB
}

func NewGoalHandler(db *gorm.DB) *GoalHandler {
	return &GoalHandler{db: db}
}

func (h *GoalHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var goals []models.FinancialGoal
	h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals)

	response.Success(c, goals)
}

func (h *GoalHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var deadline *time.Time
	if req.Deadline != "" {
		t, _ := time.Parse("2006-01-02", req.Deadline)
		deadline = &t
	}

	goal := models.FinancialGoal{
		UserID:        userID,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
		Deadline:     deadline,
		Category:      req.Category,
		Note:          req.Note,
	}

	h.db.Create(&goal)
	response.Success(c, goal)
}

func (h *GoalHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var goal models.FinancialGoal
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Goal not found")
		return
	}

	if req.Deadline != "" {
		t, _ := time.Parse("2006-01-02", req.Deadline)
		goal.Deadline = &t
	}

	goal.Name = req.Name
	goal.TargetAmount = req.TargetAmount
	goal.CurrentAmount = req.CurrentAmount
	goal.Category = req.Category
	goal.Note = req.Note

	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Status = "completed"
	}

	h.db.Save(&goal)
	response.Success(c, goal)
}

func (h *GoalHandler) AddAmount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	type AddAmountRequest struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	var req AddAmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var goal models.FinancialGoal
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Goal not found")
		return
	}

	goal.CurrentAmount += req.Amount
	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Status = "completed"
	}

	h.db.Save(&goal)
	response.Success(c, goal)
}

func (h *GoalHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.FinancialGoal{})
	response.SuccessWithMessage(c, "Goal deleted", nil)
}
