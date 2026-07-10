package handlers

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ledger/backend/internal/models"
)

// Shared request types
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Currency  string `json:"currency"`
	Timezone  string `json:"timezone"`
	Phone     string `json:"phone"`
}

type CreateAccountRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Icon      string  `json:"icon"`
	Color     string  `json:"color"`
	IsDefault bool    `json:"is_default"`
	SortOrder int     `json:"sort_order"`
}

type CreateCategoryRequest struct {
	Name      string     `json:"name" binding:"required"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Icon      string     `json:"icon"`
	Color     string     `json:"color"`
	Type      string     `json:"type" binding:"required"`
	SortOrder int        `json:"sort_order"`
}

type CreateTransactionRequest struct {
	AccountID       *uuid.UUID  `json:"account_id" binding:"required"`
	TargetAccountID *uuid.UUID  `json:"target_account_id"`
	CategoryID      *uuid.UUID  `json:"category_id"`
	Type            string      `json:"type" binding:"required"`
	Amount          float64     `json:"amount" binding:"required"`
	Currency        string      `json:"currency"`
	Merchant        string      `json:"merchant"`
	Note            string      `json:"note"`
	Tags            []uuid.UUID `json:"tags"`
	BillDate        string      `json:"bill_date" binding:"required"`
	IsRecurring     bool        `json:"is_recurring"`
	RecurringRule   string      `json:"recurring_rule"`
}

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

type CreateBudgetRequest struct {
	CategoryID     *uuid.UUID `json:"category_id"`
	Amount         float64    `json:"amount" binding:"required"`
	Period         string     `json:"period" binding:"required"`
	StartDate      string     `json:"start_date" binding:"required"`
	EndDate        string     `json:"end_date"`
	AlertThreshold float64    `json:"alert_threshold"`
}

type CreateFamilyRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateAssetRequest struct {
	AssetType    string  `json:"asset_type" binding:"required"`
	RelatedUser  string  `json:"related_user"`
	Name         string  `json:"name" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	InterestRate float64 `json:"interest_rate"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Status       string  `json:"status"`
	Note         string  `json:"note"`
}

type CreateReminderRequest struct {
	Title      string     `json:"title" binding:"required"`
	Content    string     `json:"content"`
	RemindTime string     `json:"remind_time" binding:"required"`
	RepeatType string     `json:"repeat_type"`
	CategoryID *uuid.UUID `json:"category_id"`
	IsActive   bool       `json:"is_active"`
}

type ImportTransactionRequest struct {
	Date     string  `json:"date" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Category string  `json:"category"`
	Account  string  `json:"account"`
	Merchant string  `json:"merchant"`
	Note     string  `json:"note"`
}

type CreateAAGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Members     []struct {
		Name string `json:"name"`
	} `json:"members"`
}

type CreateGoalRequest struct {
	Name          string  `json:"name" binding:"required"`
	TargetAmount  float64 `json:"target_amount" binding:"required"`
	CurrentAmount float64 `json:"current_amount"`
	Deadline      string  `json:"deadline"`
	Category      string  `json:"category"`
	Note          string  `json:"note"`
}

type CreateInsuranceRequest struct {
	Name          string  `json:"name" binding:"required"`
	InsuranceType string  `json:"insurance_type"`
	Company       string  `json:"company"`
	Premium       float64 `json:"premium"`
	PaymentType   string  `json:"payment_type"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Coverage      float64 `json:"coverage"`
	Beneficiary   string  `json:"beneficiary"`
	Note          string  `json:"note"`
}

type CSVTransaction struct {
	Date     string `json:"date"`
	Type     string `json:"type"`
	Amount   string `json:"amount"`
	Category string `json:"category"`
	Account  string `json:"account"`
	Merchant string `json:"merchant"`
	Note     string `json:"note"`
}

type CashFlowProjection struct {
	Date         string  `json:"date"`
	ProjectedBal float64 `json:"projected_balance"`
	Income       float64 `json:"income"`
	Expense      float64 `json:"expense"`
	RecurringTx  int     `json:"recurring_transactions"`
}

type BudgetAlert struct {
	BudgetID     string  `json:"budget_id"`
	CategoryName string  `json:"category_name"`
	BudgetAmount float64 `json:"budget_amount"`
	SpentAmount  float64 `json:"spent_amount"`
	Remaining    float64 `json:"remaining"`
	AlertType    string  `json:"alert_type"`
	Percentage   float64 `json:"percentage"`
}

type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type Currency struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Rate   float64 `json:"rate"`
}

// Shared helper functions
func buildCategoryTree(categories []models.Category) []map[string]interface{} {
	var result []map[string]interface{}
	categoryMap := make(map[uuid.UUID]*models.Category)

	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	for i := range categories {
		cat := &categories[i]
		node := map[string]interface{}{
			"id":         cat.ID,
			"name":       cat.Name,
			"icon":       cat.Icon,
			"color":      cat.Color,
			"type":       cat.Type,
			"sort_order": cat.SortOrder,
			"children":   []map[string]interface{}{},
		}

		if cat.ParentID != nil {
			if parent, ok := categoryMap[*cat.ParentID]; ok {
				parent.Children = append(parent.Children, *cat)
				continue
			}
		}
		result = append(result, node)
	}

	return result
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func seedCategories(db *gorm.DB, userID uuid.UUID) {
	expenseCategories := []struct {
		Name  string
		Icon  string
		Color string
	}{
		{"餐饮", "utensils", "#FF6B6B"},
		{"交通", "car", "#4ECDC4"},
		{"购物", "shopping-bag", "#45B7D1"},
		{"娱乐", "gamepad-2", "#96CEB4"},
		{"居住", "home", "#FFEAA7"},
		{"医疗", "heart-pulse", "#DDA0DD"},
		{"教育", "graduation-cap", "#98D8C8"},
		{"通讯", "phone", "#F7DC6F"},
		{"其他支出", "more-horizontal", "#BDC3C7"},
	}

	incomeCategories := []struct {
		Name  string
		Icon  string
		Color string
	}{
		{"工资", "briefcase", "#27AE60"},
		{"奖金", "gift", "#2ECC71"},
		{"投资收益", "trending-up", "#3498DB"},
		{"其他收入", "plus-circle", "#9B59B6"},
	}

	for i, cat := range expenseCategories {
		db.Create(&models.Category{
			UserID:    userID,
			Name:      cat.Name,
			Icon:      cat.Icon,
			Color:     cat.Color,
			Type:      "expense",
			SortOrder: i,
			IsSystem:  false,
		})
	}

	for i, cat := range incomeCategories {
		db.Create(&models.Category{
			UserID:    userID,
			Name:      cat.Name,
			Icon:      cat.Icon,
			Color:     cat.Color,
			Type:      "income",
			SortOrder: i,
			IsSystem:  false,
		})
	}
}

func seedDefaultAccount(db *gorm.DB, userID uuid.UUID) {
	db.Create(&models.Account{
		UserID:    userID,
		Name:      "现金",
		Type:      "cash",
		Balance:   0,
		IsDefault: true,
	})
}

func shouldRecurToday(lastDate *time.Time, today time.Time) bool {
	if lastDate == nil {
		return false
	}
	return today.After(*lastDate) && today.Day() == lastDate.Day()
}

var defaultExchangeRates = map[string]float64{
	"CNY": 1.0,
	"USD": 7.2,
	"EUR": 7.8,
	"JPY": 0.048,
	"GBP": 9.1,
	"HKD": 0.92,
	"KRW": 0.0054,
	"AUD": 4.7,
	"CAD": 5.3,
	"SGD": 5.35,
}
