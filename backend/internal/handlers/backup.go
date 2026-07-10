package handlers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type BackupHandler struct {
	db *gorm.DB
}

func NewBackupHandler(db *gorm.DB) *BackupHandler {
	return &BackupHandler{db: db}
}

func (h *BackupHandler) ExportAll(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var accounts []models.Account
	h.db.Where("user_id = ?", userID).Find(&accounts)

	var categories []models.Category
	h.db.Where("user_id = ?", userID).Find(&categories)

	var transactions []models.Transaction
	h.db.Where("user_id = ?", userID).Find(&transactions)

	var budgets []models.Budget
	h.db.Where("user_id = ?", userID).Find(&budgets)

	var tags []models.Tag
	h.db.Where("user_id = ?", userID).Find(&tags)

	var reminders []models.Reminder
	h.db.Where("user_id = ?", userID).Find(&reminders)

	var goals []models.FinancialGoal
	h.db.Where("user_id = ?", userID).Find(&goals)

	var insurances []models.Insurance
	h.db.Where("user_id = ?", userID).Find(&insurances)

	var assetChanges []models.AssetChange
	h.db.Where("user_id = ?", userID).Find(&assetChanges)

	response.Success(c, gin.H{
		"version":       "1.0",
		"exported_at":   time.Now().Format("2006-01-02 15:04:05"),
		"accounts":      accounts,
		"categories":    categories,
		"transactions":  transactions,
		"budgets":      budgets,
		"tags":          tags,
		"reminders":     reminders,
		"goals":         goals,
		"insurances":    insurances,
		"asset_changes": assetChanges,
	})
}

func (h *BackupHandler) ImportAll(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Accounts      []models.Account      `json:"accounts"`
		Categories    []models.Category      `json:"categories"`
		Transactions  []models.Transaction   `json:"transactions"`
		Budgets       []models.Budget        `json:"budgets"`
		Tags           []models.Tag           `json:"tags"`
		Reminders     []models.Reminder      `json:"reminders"`
		Goals         []models.FinancialGoal `json:"goals"`
		Insurances    []models.Insurance     `json:"insurances"`
		AssetChanges  []models.AssetChange   `json:"asset_changes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	imported := map[string]int{}

	if req.Accounts != nil {
		for i := range req.Accounts {
			req.Accounts[i].UserID = userID
			req.Accounts[i].ID = uuid.Nil
		}
		h.db.Create(&req.Accounts)
		imported["accounts"] = len(req.Accounts)
	}

	if req.Categories != nil {
		for i := range req.Categories {
			req.Categories[i].UserID = userID
			req.Categories[i].ID = uuid.Nil
		}
		h.db.Create(&req.Categories)
		imported["categories"] = len(req.Categories)
	}

	if req.Transactions != nil {
		for i := range req.Transactions {
			req.Transactions[i].UserID = userID
			req.Transactions[i].ID = uuid.Nil
		}
		h.db.Create(&req.Transactions)
		imported["transactions"] = len(req.Transactions)
	}

	if req.Budgets != nil {
		for i := range req.Budgets {
			req.Budgets[i].UserID = userID
			req.Budgets[i].ID = uuid.Nil
		}
		h.db.Create(&req.Budgets)
		imported["budgets"] = len(req.Budgets)
	}

	if req.Tags != nil {
		for i := range req.Tags {
			req.Tags[i].UserID = userID
			req.Tags[i].ID = uuid.Nil
		}
		h.db.Create(&req.Tags)
		imported["tags"] = len(req.Tags)
	}

	if req.Reminders != nil {
		for i := range req.Reminders {
			req.Reminders[i].UserID = userID
			req.Reminders[i].ID = uuid.Nil
		}
		h.db.Create(&req.Reminders)
		imported["reminders"] = len(req.Reminders)
	}

	if req.Goals != nil {
		for i := range req.Goals {
			req.Goals[i].UserID = userID
			req.Goals[i].ID = uuid.Nil
		}
		h.db.Create(&req.Goals)
		imported["goals"] = len(req.Goals)
	}

	if req.Insurances != nil {
		for i := range req.Insurances {
			req.Insurances[i].UserID = userID
			req.Insurances[i].ID = uuid.Nil
		}
		h.db.Create(&req.Insurances)
		imported["insurances"] = len(req.Insurances)
	}

	if req.AssetChanges != nil {
		for i := range req.AssetChanges {
			req.AssetChanges[i].UserID = userID
			req.AssetChanges[i].ID = uuid.Nil
		}
		h.db.Create(&req.AssetChanges)
		imported["asset_changes"] = len(req.AssetChanges)
	}

	backup := models.Backup{
		UserID:     userID,
		BackupType: "manual",
		FileName:   "import_" + time.Now().Format("20060102150405"),
	}
	h.db.Create(&backup)

	response.Success(c, gin.H{
		"message":  "Import completed",
		"imported": imported,
	})
}

func (h *BackupHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var backups []models.Backup
	h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&backups)

	response.Success(c, backups)
}
