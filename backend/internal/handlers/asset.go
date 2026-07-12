package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AssetHandler struct {
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (h *AssetHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	assetType := c.Query("type")

	assets, err := h.assetService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Filter by type if provided
	var filtered []models.AssetChange
	for _, a := range assets {
		if assetType == "" || a.AssetType == assetType {
			filtered = append(filtered, a)
		}
	}

	response.Success(c, filtered)
}

func (h *AssetHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var startDate, endDate *time.Time
	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		startDate = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		endDate = &t
	}

	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	createReq := &service.AssetCreateRequest{
		AssetType:    req.AssetType,
		RelatedUser:  req.RelatedUser,
		Name:         req.Name,
		Amount:       req.Amount,
		InterestRate: req.InterestRate,
		StartDate:    startDate,
		EndDate:      endDate,
		Note:         req.Note,
	}

	asset, err := h.assetService.Create(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, asset)
}

func (h *AssetHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	asset, err := h.assetService.Get(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Asset not found")
		return
	}

	response.Success(c, asset)
}

func (h *AssetHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	updates["name"] = req.Name
	updates["amount"] = req.Amount
	updates["status"] = req.Status
	updates["note"] = req.Note

	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		updates["start_date"] = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		updates["end_date"] = &t
	}

	asset, err := h.assetService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, asset)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.assetService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Asset deleted", nil)
}

func (h *AssetHandler) GetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)

	summary, err := h.assetService.GetSummary(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, summary)
}
