package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type FamilyHandler struct {
	familyService *service.FamilyService
	db            *gorm.DB
}

func NewFamilyHandler(familyService *service.FamilyService, db *gorm.DB) *FamilyHandler {
	return &FamilyHandler{familyService: familyService, db: db}
}

func (h *FamilyHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var families []models.Family
	h.db.Joins("JOIN family_members ON family_members.family_id = families.id").
		Where("family_members.user_id = ?", userID).
		Preload("Members").Preload("Owner").
		Find(&families)

	response.Success(c, families)
}

func (h *FamilyHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	createReq := &service.FamilyCreateRequest{
		Name: req.Name,
	}

	family, err := h.familyService.Create(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Add owner as member
	member := models.FamilyMember{
		FamilyID: family.ID,
		UserID:   userID,
		Role:     "owner",
	}
	h.db.Create(&member)

	response.Success(c, family)
}

func (h *FamilyHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var family models.Family
	if err := h.db.Where("id = ?", id).
		Joins("JOIN family_members ON family_members.family_id = families.id").
		Where("family_members.user_id = ?", userID).
		Preload("Members").Preload("Owner").First(&family).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Family not found")
		return
	}

	response.Success(c, family)
}

func (h *FamilyHandler) AddMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	type AddMemberRequest struct {
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role"`
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var family models.Family
	if err := h.db.Where("id = ? AND owner_id = ?", id, userID).First(&family).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Family not found or not authorized")
		return
	}

	var newUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&newUser).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	role := "member"
	if req.Role != "" {
		role = req.Role
	}

	member := models.FamilyMember{
		FamilyID: family.ID,
		UserID:   newUser.ID,
		Role:     role,
	}
	h.db.Create(&member)

	response.SuccessWithMessage(c, "Member added", nil)
}

func (h *FamilyHandler) RemoveMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))
	memberID := c.Param("member_id")

	var family models.Family
	if err := h.db.Where("id = ? AND owner_id = ?", id, userID).First(&family).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Family not found or not authorized")
		return
	}

	h.db.Where("id = ? AND family_id = ?", memberID, id).Delete(&models.FamilyMember{})
	response.SuccessWithMessage(c, "Member removed", nil)
}

func (h *FamilyHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var family models.Family
	if err := h.db.Where("id = ? AND owner_id = ?", id, userID).First(&family).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Family not found or not authorized")
		return
	}

	h.db.Where("family_id = ?", id).Delete(&models.FamilyMember{})
	h.db.Delete(&family)
	response.SuccessWithMessage(c, "Family deleted", nil)
}

func (h *FamilyHandler) GetTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var family models.Family
	if err := h.db.Where("id = ?", id).
		Joins("JOIN family_members ON family_members.family_id = families.id").
		Where("family_members.user_id = ?", userID).
		First(&family).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Family not found or not authorized")
		return
	}

	var transactions []models.FamilyTransaction
	h.db.Where("family_id = ?", id).Preload("Category").Preload("Account").Preload("User").Order("bill_date DESC").Find(&transactions)

	response.Success(c, transactions)
}
