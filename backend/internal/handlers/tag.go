package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type TagHandler struct {
	db *gorm.DB
}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{db: db}
}

func (h *TagHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var tags []models.Tag
	h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tags)

	response.Success(c, tags)
}

func (h *TagHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tag := models.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	h.db.Create(&tag)
	response.Success(c, tag)
}

func (h *TagHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var tag models.Tag
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Tag not found")
		return
	}

	tag.Name = req.Name
	tag.Color = req.Color
	h.db.Save(&tag)

	response.Success(c, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Tag{})
	response.SuccessWithMessage(c, "Tag deleted", nil)
}
