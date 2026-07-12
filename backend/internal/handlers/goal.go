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

type GoalHandler struct {
	goalService *service.GoalService
}

func NewGoalHandler(goalService *service.GoalService) *GoalHandler {
	return &GoalHandler{goalService: goalService}
}

func (h *GoalHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	goals, err := h.goalService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

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

	createReq := &service.GoalCreateRequest{
		Name:         req.Name,
		TargetAmount: req.TargetAmount,
		Deadline:     deadline,
		Category:     req.Category,
		Note:         req.Note,
	}

	goal, err := h.goalService.Create(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, goal)
}

func (h *GoalHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	updates["name"] = req.Name
	updates["target_amount"] = req.TargetAmount
	updates["category"] = req.Category
	updates["note"] = req.Note

	if req.Deadline != "" {
		t, _ := time.Parse("2006-01-02", req.Deadline)
		updates["deadline"] = &t
	}

	goal, err := h.goalService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, goal)
}

func (h *GoalHandler) AddAmount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	type AddAmountRequest struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	var req AddAmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	goal, err := h.goalService.AddAmount(c.Request.Context(), id, userID, req.Amount)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, goal)
}

func (h *GoalHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.goalService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Goal deleted", nil)
}
