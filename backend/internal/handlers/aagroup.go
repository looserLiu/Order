package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AAGroupHandler struct {
	aaService *service.AAService
}

func NewAAGroupHandler(aaService *service.AAService) *AAGroupHandler {
	return &AAGroupHandler{aaService: aaService}
}

func (h *AAGroupHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	groups, err := h.aaService.ListGroups(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, groups)
}

func (h *AAGroupHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAAGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	createReq := &service.CreateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
		Members:       req.Members,
	}

	group, err := h.aaService.CreateGroup(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, group)
}

func (h *AAGroupHandler) AddExpense(c *gin.Context) {
	userID := middleware.GetUserID(c)
	groupID, _ := uuid.Parse(c.Param("id"))

	var req service.AddExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.aaService.AddExpense(c.Request.Context(), groupID, userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Expense added", nil)
}

func (h *AAGroupHandler) GetSettlements(c *gin.Context) {
	userID := middleware.GetUserID(c)
	groupID, _ := uuid.Parse(c.Param("id"))

	settlements, err := h.aaService.GetSettlements(c.Request.Context(), groupID, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, settlements)
}

func (h *AAGroupHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.aaService.DeleteGroup(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Group deleted", nil)
}
