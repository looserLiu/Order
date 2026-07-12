package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	unreadOnly := c.Query("unread") == "true"

	notifications, err := h.notificationService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var unreadCount int64
	// Get unread count - would need to add method to service
	_ = unreadOnly
	_ = unreadCount

	response.Success(c, gin.H{
		"list":         notifications,
		"unread_count": 0,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.notificationService.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "All notifications marked as read", nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.notificationService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Notification deleted", nil)
}
