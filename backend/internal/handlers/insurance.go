package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type InsuranceHandler struct {
	insuranceService *service.InsuranceService
}

func NewInsuranceHandler(insuranceService *service.InsuranceService) *InsuranceHandler {
	return &InsuranceHandler{insuranceService: insuranceService}
}

func (h *InsuranceHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	insurances, err := h.insuranceService.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, insurances)
}

func (h *InsuranceHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateInsuranceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	var endDate *time.Time
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		endDate = &t
	}

	createReq := &service.InsuranceCreateRequest{
		Name:         req.Name,
		InsuranceType: req.InsuranceType,
		Company:      req.Company,
		Premium:      req.Premium,
		PaymentType:  req.PaymentType,
		StartDate:    startDate,
		EndDate:      endDate,
		Coverage:     req.Coverage,
		Beneficiary:  req.Beneficiary,
		Note:         req.Note,
	}

	insurance, err := h.insuranceService.Create(c.Request.Context(), userID, createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, insurance)
}

func (h *InsuranceHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	var req CreateInsuranceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	updates["name"] = req.Name
	updates["insurance_type"] = req.InsuranceType
	updates["company"] = req.Company
	updates["premium"] = req.Premium
	updates["payment_type"] = req.PaymentType
	updates["coverage"] = req.Coverage
	updates["beneficiary"] = req.Beneficiary
	updates["note"] = req.Note

	insurance, err := h.insuranceService.Update(c.Request.Context(), id, userID, updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, insurance)
}

func (h *InsuranceHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := uuid.Parse(c.Param("id"))

	if err := h.insuranceService.Delete(c.Request.Context(), id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Insurance deleted", nil)
}

func (h *InsuranceHandler) GetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)

	summary, err := h.insuranceService.GetSummary(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, summary)
}
