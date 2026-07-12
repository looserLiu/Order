package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetMe(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Currency != "" {
		updates["currency"] = req.Currency
	}
	if req.Timezone != "" {
		updates["timezone"] = req.Timezone
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if err := h.userService.UpdateMe(c.Request.Context(), userID, updates); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, _ := h.userService.GetMe(c.Request.Context(), userID)
	response.Success(c, user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	type ChangePasswordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Password changed successfully", nil)
}
