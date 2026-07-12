package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	accounts, err := h.accountService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, accounts)
}

func (h *AccountHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.AccountCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.accountService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, account)
}

func (h *AccountHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	account, err := h.accountService.Get(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Account not found")
		return
	}

	response.Success(c, account)
}

func (h *AccountHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var req service.AccountCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Currency != "" {
		updates["currency"] = req.Currency
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}

	account, err := h.accountService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, account)
}

func (h *AccountHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.accountService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Account deleted", nil)
}

func (h *AccountHandler) GetTotalBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)

	balance, err := h.accountService.GetTotalBalance(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"total_balance": balance})
}
