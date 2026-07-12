package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type ReportHandler struct {
	reportService *service.ReportService
	db            *gorm.DB
}

func NewReportHandler(reportService *service.ReportService, db *gorm.DB) *ReportHandler {
	return &ReportHandler{reportService: reportService, db: db}
}

func (h *ReportHandler) Summary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	sd, _ := time.Parse("2006-01-02", startDate)
	ed, _ := time.Parse("2006-01-02", endDate)

	summary, err := h.reportService.GetSummary(c.Request.Context(), userID, sd, ed)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, summary)
}

func (h *ReportHandler) Trend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	sd, _ := time.Parse("2006-01-02", startDate)
	ed, _ := time.Parse("2006-01-02", endDate)

	trends, err := h.reportService.GetTrend(c.Request.Context(), userID, sd, ed)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, trends)
}

func (h *ReportHandler) ByCategory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	txType := c.DefaultQuery("type", "expense")

	sd, _ := time.Parse("2006-01-02", startDate)
	ed, _ := time.Parse("2006-01-02", endDate)

	data, err := h.reportService.GetByCategory(c.Request.Context(), userID, txType, sd, ed)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *ReportHandler) ByAccount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	sd, _ := time.Parse("2006-01-02", startDate)
	ed, _ := time.Parse("2006-01-02", endDate)

	data, err := h.reportService.GetByAccount(c.Request.Context(), userID, sd, ed)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *ReportHandler) ByMerchant(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	sd, _ := time.Parse("2006-01-02", startDate)
	ed, _ := time.Parse("2006-01-02", endDate)

	data, err := h.reportService.GetByMerchant(c.Request.Context(), userID, sd, ed)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *ReportHandler) MonthlyCompare(c *gin.Context) {
	userID := middleware.GetUserID(c)
	months := parseInt(c.DefaultQuery("months", "6"))

	data, err := h.reportService.GetMonthlyCompare(c.Request.Context(), userID, months)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}

func (h *ReportHandler) Export(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -365).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	sd, _ := time.Parse("2006-01-02", startDate)
	ed, _ := time.Parse("2006-01-02", endDate)

	transactions, err := h.reportService.ExportTransactions(c.Request.Context(), userID, sd, ed)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, transactions)
}
