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

type ReminderHandler struct {
	db *gorm.DB
}

func NewReminderHandler(db *gorm.DB) *ReminderHandler {
	return &ReminderHandler{db: db}
}

func (h *ReminderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var reminders []models.Reminder
	h.db.Where("user_id = ?", userID).Order("remind_time ASC").Find(&reminders)

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

	reminder := models.Reminder{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		RemindTime: remindTime,
		RepeatType: req.RepeatType,
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
	}

	h.db.Create(&reminder)
	response.Success(c, reminder)
}

func (h *ReminderHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var reminder models.Reminder
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&reminder).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Reminder not found")
		return
	}

	remindTime, _ := time.Parse("2006-01-02T15:04", req.RemindTime)
	reminder.Title = req.Title
	reminder.Content = req.Content
	reminder.RemindTime = remindTime
	reminder.RepeatType = req.RepeatType
	reminder.CategoryID = req.CategoryID
	reminder.IsActive = req.IsActive

	h.db.Save(&reminder)
	response.Success(c, reminder)
}

func (h *ReminderHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Reminder{})
	response.SuccessWithMessage(c, "Reminder deleted", nil)
}
