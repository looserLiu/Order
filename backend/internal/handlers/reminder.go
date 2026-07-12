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

type ReminderHandler struct {
	reminderService *service.ReminderService
}

func NewReminderHandler(reminderService *service.ReminderService) *ReminderHandler {
	return &ReminderHandler{reminderService: reminderService}
}

func (h *ReminderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	reminders, err := h.reminderService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, reminders)
}

func (h *ReminderHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	remindTime, _ := time.Parse("2006-01-02T15:04", req.RemindTime)

	createReq := &service.ReminderCreateRequest{
		Title:      req.Title,
		Content:    req.Content,
		RemindTime: remindTime,
		RepeatType: req.RepeatType,
		CategoryID: req.CategoryID,
	}

	reminder, err := h.reminderService.Create(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, reminder)
}

func (h *ReminderHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	updates["title"] = req.Title
	updates["content"] = req.Content
	remindTime, _ := time.Parse("2006-01-02T15:04", req.RemindTime)
	updates["remind_time"] = remindTime
	updates["repeat_type"] = req.RepeatType
	updates["is_active"] = req.IsActive

	reminder, err := h.reminderService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, reminder)
}

func (h *ReminderHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.reminderService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Reminder deleted", nil)
}
