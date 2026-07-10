package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type NotificationHandler struct {
	db *gorm.DB
}

func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{db: db}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	 unreadOnly := c.Query("unread")

	var notifications []models.Notification
	query := h.db.Where("user_id = ?", userID)
	
	if unreadOnly == "true" {
		query = query.Where("is_read = ?", false)
	}
	
	query.Order("created_at DESC").Limit(50).Find(&notifications)

	var unreadCount int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	response.Success(c, gin.H{
		"list":       notifications,
		"unread_count": unreadCount,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true)
	response.SuccessWithMessage(c, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	h.db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true)
	response.SuccessWithMessage(c, "All notifications marked as read", nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})
	response.SuccessWithMessage(c, "Notification deleted", nil)
}
