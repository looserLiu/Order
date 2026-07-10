package handlers

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) Summary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var income, expense float64
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date BETWEEN ? AND ?", userID, "income", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&income)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date BETWEEN ? AND ?", userID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&expense)

	response.Success(c, gin.H{
		"income":  income,
		"expense": expense,
		"balance": income - expense,
	})
}

func (h *ReportHandler) Trend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type TrendData struct {
		Date     string  `json:"date"`
		Income   float64 `json:"income"`
		Expense  float64 `json:"expense"`
		Count    int     `json:"count"`
	}

	var trends []TrendData
	h.db.Raw(`
		SELECT 
			bill_date as date,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense,
			COUNT(*) as count
		FROM transactions
		WHERE user_id = ? AND bill_date BETWEEN ? AND ? AND deleted_at IS NULL
		GROUP BY bill_date
		ORDER BY bill_date
	`, userID, startDate, endDate).Scan(&trends)

	response.Success(c, trends)
}

func (h *ReportHandler) ByCategory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	txType := c.DefaultQuery("type", "expense")

	type CategoryData struct {
		CategoryID    *uuid.UUID `json:"category_id"`
		CategoryName string     `json:"category_name"`
		CategoryIcon string     `json:"category_icon"`
		CategoryColor string    `json:"category_color"`
		Total        float64    `json:"total"`
		Percentage   float64    `json:"percentage"`
		Count        int        `json:"count"`
	}

	var data []CategoryData
	h.db.Raw(`
		SELECT 
			c.id as category_id,
			COALESCE(c.name, 'Uncategorized') as category_name,
			COALESCE(c.icon, '') as category_icon,
			COALESCE(c.color, '#999999') as category_color,
			COALESCE(SUM(t.amount), 0) as total,
			COUNT(t.id) as count
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? AND t.type = ? AND t.bill_date BETWEEN ? AND ? AND t.deleted_at IS NULL
		GROUP BY c.id, c.name, c.icon, c.color
		ORDER BY total DESC
	`, userID, txType, startDate, endDate).Scan(&data)

	var total float64
	for _, d := range data {
		total += d.Total
	}
	for i := range data {
		if total > 0 {
			data[i].Percentage = data[i].Total / total * 100
		}
	}

	response.Success(c, data)
}

func (h *ReportHandler) ByAccount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type AccountData struct {
		AccountID   uuid.UUID `json:"account_id"`
		AccountName string    `json:"account_name"`
		Total       float64   `json:"total"`
		Count       int       `json:"count"`
	}

	var data []AccountData
	h.db.Raw(`
		SELECT 
			a.id as account_id,
			a.name as account_name,
			COALESCE(SUM(t.amount), 0) as total,
			COUNT(t.id) as count
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		WHERE t.user_id = ? AND t.bill_date BETWEEN ? AND ? AND t.deleted_at IS NULL
		GROUP BY a.id, a.name
		ORDER BY total DESC
	`, userID, startDate, endDate).Scan(&data)

	response.Success(c, data)
}

func (h *ReportHandler) ByMerchant(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type MerchantData struct {
		Merchant string  `json:"merchant"`
		Total    float64 `json:"total"`
		Count    int     `json:"count"`
	}

	var data []MerchantData
	h.db.Raw(`
		SELECT 
			COALESCE(merchant, 'Unknown') as merchant,
			COALESCE(SUM(amount), 0) as total,
			COUNT(*) as count
		FROM transactions
		WHERE user_id = ? AND type = 'expense' AND bill_date BETWEEN ? AND ? AND deleted_at IS NULL
		GROUP BY merchant
		ORDER BY total DESC
		LIMIT 20
	`, userID, startDate, endDate).Scan(&data)

	response.Success(c, data)
}

func (h *ReportHandler) MonthlyCompare(c *gin.Context) {
	userID := middleware.GetUserID(c)
	months := c.DefaultQuery("months", "6")

	type MonthlyData struct {
		Month   string  `json:"month"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
	}

	var data []MonthlyData
	h.db.Raw(`
		SELECT 
			TO_CHAR(bill_date, 'YYYY-MM') as month,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense
		FROM transactions
		WHERE user_id = ? AND bill_date >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 month' * ?
		GROUP BY TO_CHAR(bill_date, 'YYYY-MM')
		ORDER BY month
	`, userID, months).Scan(&data)

	response.Success(c, data)
}

func (h *ReportHandler) Export(c *gin.Context) {
	userID := middleware.GetUserID(c)
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -365).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var transactions []models.Transaction
	h.db.Where("user_id = ? AND bill_date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("bill_date DESC").
		Preload("Category").Preload("Account").
		Find(&transactions)

	response.Success(c, transactions)
}
