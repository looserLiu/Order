package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AAGroupHandler struct {
	db *gorm.DB
}

func NewAAGroupHandler(db *gorm.DB) *AAGroupHandler {
	return &AAGroupHandler{db: db}
}

func (h *AAGroupHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var groups []models.AAGroup
	h.db.Where("user_id = ?", userID).Preload("Members").Order("created_at DESC").Find(&groups)

	response.Success(c, groups)
}

func (h *AAGroupHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAAGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	group := models.AAGroup{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	h.db.Create(&group)

	for _, m := range req.Members {
		member := models.AAMember{
			GroupID: group.ID,
			Name:    m.Name,
		}
		h.db.Create(&member)
	}

	h.db.Preload("Members").First(&group)
	response.Success(c, group)
}

func (h *AAGroupHandler) AddExpense(c *gin.Context) {
	userID := middleware.GetUserID(c)
	groupID := c.Param("id")

	// Verify user owns the group
	var group models.AAGroup
	if err := h.db.Where("id = ? AND user_id = ?", groupID, userID).First(&group).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Group not found")
		return
	}

	type AddExpenseRequest struct {
		MemberID string  `json:"member_id" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
		Note     string  `json:"note"`
	}

	var req AddExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	memberID, _ := uuid.Parse(req.MemberID)
	var member models.AAMember
	if err := h.db.Where("id = ? AND group_id = ?", memberID, groupID).First(&member).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Member not found")
		return
	}

	member.Paid += req.Amount
	h.db.Save(&member)

	group.TotalAmount += req.Amount
	h.db.Save(&group)

	h.calculateOwe(groupID)

	response.SuccessWithMessage(c, "Expense added", nil)
}

func (h *AAGroupHandler) calculateOwe(groupID string) {
	var members []models.AAMember
	h.db.Where("group_id = ?", groupID).Find(&members)

	var totalPaid float64
	for _, m := range members {
		totalPaid += m.Paid
	}

	perPerson := totalPaid / float64(len(members))
	for i := range members {
		members[i].Owe = perPerson - members[i].Paid
		if members[i].Owe < 0 {
			members[i].Owe = 0
		}
		h.db.Save(&members[i])
	}
}

func (h *AAGroupHandler) GetSettlements(c *gin.Context) {
	groupID := c.Param("id")

	type Settlement struct {
		FromID uuid.UUID `json:"from_id"`
		ToID   uuid.UUID `json:"to_id"`
		Amount float64   `json:"amount"`
	}

	var members []models.AAMember
	h.db.Where("group_id = ?", groupID).Find(&members)

	var settlements []Settlement

	for i := range members {
		for j := range members {
			if i == j {
				continue
			}
			if members[i].Owe > 0 && members[j].Paid > members[j].Owe {
				amount := min(members[i].Owe, members[j].Paid-members[j].Owe)
				if amount > 0 {
					settlements = append(settlements, Settlement{
						FromID: members[i].ID,
						ToID:   members[j].ID,
						Amount: amount,
					})
					members[i].Owe -= amount
					members[j].Owe += amount
				}
			}
		}
	}

	response.Success(c, settlements)
}

func (h *AAGroupHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("group_id = ?", id).Delete(&models.AAMember{})
	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.AAGroup{})
	response.SuccessWithMessage(c, "Group deleted", nil)
}
