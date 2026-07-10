package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var categories []models.Category
	h.db.Where("user_id = ? OR is_system = ?", userID, true).Order("sort_order ASC").Find(&categories)

	response.Success(c, categories)
}

func (h *CategoryHandler) GetTree(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var categories []models.Category
	h.db.Where("user_id = ? OR is_system = ?", userID, true).Order("sort_order ASC").Find(&categories)

	tree := buildCategoryTree(categories)
	response.Success(c, tree)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category := models.Category{
		UserID:    userID,
		ParentID:  req.ParentID,
		Name:      req.Name,
		Icon:      req.Icon,
		Color:     req.Color,
		Type:      req.Type,
		SortOrder: req.SortOrder,
	}

	h.db.Create(&category)
	response.Success(c, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var category models.Category
	if err := h.db.Where("id = ? AND user_id = ? AND is_system = ?", id, userID, false).First(&category).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	category.Name = req.Name
	category.ParentID = req.ParentID
	category.Icon = req.Icon
	category.Color = req.Color
	category.Type = req.Type

	h.db.Save(&category)
	response.Success(c, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	if err := h.db.Where("id = ? AND user_id = ? AND is_system = ?", id, userID, false).Delete(&models.Category{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	response.SuccessWithMessage(c, "Category deleted", nil)
}
