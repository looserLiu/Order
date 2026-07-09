package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ledger/backend/internal/config"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

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

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
	}

	if err := h.db.Create(&user).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Email already exists")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	seedCategories(h.db, user.ID)
	seedDefaultAccount(h.db, user.ID)

	response.Success(c, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) generateToken(user models.User) (string, error) {
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, user)
}

type UpdateUserRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Currency  string `json:"currency"`
	Timezone  string `json:"timezone"`
	Phone     string `json:"phone"`
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Currency != "" {
		user.Currency = req.Currency
	}
	if req.Timezone != "" {
		user.Timezone = req.Timezone
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	h.db.Save(&user)
	response.Success(c, user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	type ChangePasswordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid old password")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hash)
	h.db.Save(&user)

	response.SuccessWithMessage(c, "Password changed successfully", nil)
}

type AccountHandler struct {
	db *gorm.DB
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{db: db}
}

func (h *AccountHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var accounts []models.Account
	h.db.Where("user_id = ?", userID).Order("sort_order ASC").Find(&accounts)

	response.Success(c, accounts)
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

func (h *AccountHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account := models.Account{
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Balance:   req.Balance,
		Currency:  req.Currency,
		Icon:      req.Icon,
		Color:     req.Color,
		IsDefault: req.IsDefault,
		SortOrder: req.SortOrder,
	}

	h.db.Create(&account)
	response.Success(c, account)
}

func (h *AccountHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var account models.Account
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Account not found")
		return
	}

	response.Success(c, account)
}

func (h *AccountHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var account models.Account
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Account not found")
		return
	}

	account.Name = req.Name
	account.Type = req.Type
	account.Balance = req.Balance
	account.Currency = req.Currency
	account.Icon = req.Icon
	account.Color = req.Color
	account.IsDefault = req.IsDefault
	account.SortOrder = req.SortOrder

	h.db.Save(&account)
	response.Success(c, account)
}

func (h *AccountHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Account{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	response.SuccessWithMessage(c, "Account deleted", nil)
}

func (h *AccountHandler) GetTotalBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var total float64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Select("COALESCE(SUM(balance), 0)").Scan(&total)

	response.Success(c, gin.H{"total_balance": total})
}

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var categories []models.Category
	h.db.Where("user_id = ? OR is_system = ?", userID, true).Order("sort_order ASC").Find(&categories)

	response.Success(c, categories)
}

func (h *CategoryHandler) GetTree(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var categories []models.Category
	h.db.Where("user_id = ? OR is_system = ?", userID, true).Order("sort_order ASC").Find(&categories)

	tree := buildCategoryTree(categories)
	response.Success(c, tree)
}

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

type CreateCategoryRequest struct {
	Name      string     `json:"name" binding:"required"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Icon      string     `json:"icon"`
	Color     string     `json:"color"`
	Type      string     `json:"type" binding:"required"`
	SortOrder int        `json:"sort_order"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category := models.Category{
		UserID:    userID,
		ParentID:  req.ParentID,
		Name:      req.Name,
		Icon:      req.Icon,
		Color:     req.Color,
		Type:      req.Type,
		SortOrder: req.SortOrder,
	}

	h.db.Create(&category)
	response.Success(c, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var category models.Category
	if err := h.db.Where("id = ? AND user_id = ? AND is_system = ?", id, userID, false).First(&category).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	category.Name = req.Name
	category.ParentID = req.ParentID
	category.Icon = req.Icon
	category.Color = req.Color
	category.Type = req.Type

	h.db.Save(&category)
	response.Success(c, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	if err := h.db.Where("id = ? AND user_id = ? AND is_system = ?", id, userID, false).Delete(&models.Category{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	response.SuccessWithMessage(c, "Category deleted", nil)
}

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	categoryID := c.Query("category_id")
	accountID := c.Query("account_id")
	txType := c.Query("type")
	tagID := c.Query("tag_id")
	merchant := c.Query("merchant")

	var transactions []models.Transaction
	query := h.db.Where("user_id = ?", userID)

	if startDate != "" {
		query = query.Where("bill_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("bill_date <= ?", endDate)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if merchant != "" {
		query = query.Where("merchant ILIKE ?", "%"+merchant+"%")
	}

	var total int64
	query.Count(&total)

	query.Order("bill_date DESC, created_at DESC").Offset((parseInt(page) - 1) * parseInt(pageSize)).Limit(parseInt(pageSize)).Preload("Category").Preload("Account").Preload("Tags").Find(&transactions)

	response.Paginate(c, transactions, total, parseInt(page), parseInt(pageSize))
}

func (h *TransactionHandler) ListByDate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	date := c.Query("date")

	var transactions []models.Transaction
	h.db.Where("user_id = ? AND bill_date = ?", userID, date).
		Preload("Category").Preload("Account").
		Order("created_at DESC").
		Find(&transactions)

	response.Success(c, transactions)
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

func (h *TransactionHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	billDate, _ := time.Parse("2006-01-02", req.BillDate)

	transaction := models.Transaction{
		UserID:          userID,
		AccountID:       *req.AccountID,
		TargetAccountID: req.TargetAccountID,
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Merchant:        req.Merchant,
		Note:            req.Note,
		Tags:            req.Tags,
		BillDate:        billDate,
		IsRecurring:     req.IsRecurring,
		RecurringRule:   req.RecurringRule,
	}

	h.db.Create(&transaction)

	h.updateAccountBalance(req.AccountID, req.Type, req.Amount)

	if req.Type == "transfer" && req.TargetAccountID != nil {
		h.updateAccountBalance(req.TargetAccountID, "income", req.Amount)
	}

	response.Success(c, transaction)
}

func (h *TransactionHandler) updateAccountBalance(accountID *uuid.UUID, txType string, amount float64) {
	var account models.Account
	if err := h.db.First(&account, accountID).Error; err == nil {
		if txType == "expense" {
			account.Balance -= amount
		} else if txType == "income" {
			account.Balance += amount
		}
		h.db.Save(&account)
	}
}

func (h *TransactionHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Category").Preload("Account").Preload("Tags").First(&transaction).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	response.Success(c, transaction)
}

func (h *TransactionHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	billDate, _ := time.Parse("2006-01-02", req.BillDate)
	transaction.AccountID = *req.AccountID
	transaction.TargetAccountID = req.TargetAccountID
	transaction.CategoryID = req.CategoryID
	transaction.Type = req.Type
	transaction.Amount = req.Amount
	transaction.Currency = req.Currency
	transaction.Merchant = req.Merchant
	transaction.Note = req.Note
	transaction.Tags = req.Tags
	transaction.BillDate = billDate

	h.db.Save(&transaction)
	response.Success(c, transaction)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	h.db.Delete(&transaction)
	response.SuccessWithMessage(c, "Transaction deleted", nil)
}

func (h *TransactionHandler) BatchDelete(c *gin.Context) {
	userID := middleware.GetUserID(c)

	type BatchDeleteRequest struct {
		IDs []uuid.UUID `json:"ids" binding:"required"`
	}

	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	h.db.Where("id IN ? AND user_id = ?", req.IDs, userID).Delete(&models.Transaction{})
	response.SuccessWithMessage(c, "Transactions deleted", nil)
}

type TagHandler struct {
	db *gorm.DB
}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{db: db}
}

func (h *TagHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var tags []models.Tag
	h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tags)

	response.Success(c, tags)
}

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

func (h *TagHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tag := models.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	h.db.Create(&tag)
	response.Success(c, tag)
}

func (h *TagHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var tag models.Tag
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Tag not found")
		return
	}

	tag.Name = req.Name
	tag.Color = req.Color
	h.db.Save(&tag)

	response.Success(c, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Tag{})
	response.SuccessWithMessage(c, "Tag deleted", nil)
}

type BudgetHandler struct {
	db *gorm.DB
}

func NewBudgetHandler(db *gorm.DB) *BudgetHandler {
	return &BudgetHandler{db: db}
}

func (h *BudgetHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var budgets []models.Budget
	h.db.Where("user_id = ?", userID).Preload("Category").Order("start_date DESC").Find(&budgets)

	response.Success(c, budgets)
}

type CreateBudgetRequest struct {
	CategoryID     *uuid.UUID `json:"category_id"`
	Amount         float64    `json:"amount" binding:"required"`
	Period         string     `json:"period" binding:"required"`
	StartDate      string     `json:"start_date" binding:"required"`
	EndDate        string     `json:"end_date"`
	AlertThreshold float64    `json:"alert_threshold"`
}

func (h *BudgetHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateBudgetRequest
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

	budget := models.Budget{
		UserID:          userID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Period:          req.Period,
		StartDate:       startDate,
		EndDate:         endDate,
		AlertThreshold:  req.AlertThreshold,
	}

	h.db.Create(&budget)
	response.Success(c, budget)
}

func (h *BudgetHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var budget models.Budget
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Category").First(&budget).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	response.Success(c, budget)
}

func (h *BudgetHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var budget models.Budget
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&budget).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		budget.EndDate = &t
	}

	budget.CategoryID = req.CategoryID
	budget.Amount = req.Amount
	budget.Period = req.Period
	budget.StartDate = startDate
	budget.AlertThreshold = req.AlertThreshold

	h.db.Save(&budget)
	response.Success(c, budget)
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Budget{})
	response.SuccessWithMessage(c, "Budget deleted", nil)
}

func (h *BudgetHandler) GetProgress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var budget models.Budget
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Category").First(&budget).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Budget not found")
		return
	}

	startDate := budget.StartDate
	endDate := time.Now()
	if budget.EndDate != nil {
		endDate = *budget.EndDate
	}

	var spent float64
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date BETWEEN ? AND ?",
			userID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&spent)

	if budget.CategoryID != nil {
		h.db.Model(&models.Transaction{}).
			Where("user_id = ? AND type = ? AND category_id = ? AND bill_date BETWEEN ? AND ?",
				userID, "expense", budget.CategoryID, startDate, endDate).
			Select("COALESCE(SUM(amount), 0)").Scan(&spent)
	}

	progress := 0.0
	if budget.Amount > 0 {
		progress = spent / budget.Amount * 100
	}

	alert := progress >= budget.AlertThreshold*100

	response.Success(c, gin.H{
		"budget":    budget,
		"spent":     spent,
		"remaining": budget.Amount - spent,
		"progress":  progress,
		"alert":     alert,
	})
}

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

type FamilyHandler struct {
	db *gorm.DB
}

func NewFamilyHandler(db *gorm.DB) *FamilyHandler {
	return &FamilyHandler{db: db}
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

type CreateFamilyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *FamilyHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	family := models.Family{
		Name:    req.Name,
		OwnerID: userID,
	}

	h.db.Create(&family)

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
	id := c.Param("id")

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
	id := c.Param("id")

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
	id := c.Param("id")
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
	id := c.Param("id")

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

type FamilyTransactionHandler struct {
	db *gorm.DB
}

func NewFamilyTransactionHandler(db *gorm.DB) *FamilyTransactionHandler {
	return &FamilyTransactionHandler{db: db}
}

func (h *FamilyTransactionHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	familyID := c.Param("id")

	type CreateFamilyTransactionRequest struct {
		AccountID   *uuid.UUID `json:"account_id" binding:"required"`
		CategoryID  *uuid.UUID `json:"category_id"`
		Type        string     `json:"type" binding:"required"`
		Amount      float64    `json:"amount" binding:"required"`
		Note        string     `json:"note"`
		BillDate    string     `json:"bill_date" binding:"required"`
	}

	var req CreateFamilyTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fID, _ := uuid.Parse(familyID)
	billDate, _ := time.Parse("2006-01-02", req.BillDate)

	tx := models.FamilyTransaction{
		FamilyID:   fID,
		UserID:     userID,
		AccountID:  *req.AccountID,
		CategoryID: req.CategoryID,
		Type:       req.Type,
		Amount:     req.Amount,
		Note:       req.Note,
		BillDate:   billDate,
	}

	h.db.Create(&tx)
	response.Success(c, tx)
}

type AssetHandler struct {
	db *gorm.DB
}

func NewAssetHandler(db *gorm.DB) *AssetHandler {
	return &AssetHandler{db: db}
}

func (h *AssetHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	assetType := c.Query("type")

	var assets []models.AssetChange
	query := h.db.Where("user_id = ?", userID)
	if assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}
	query.Order("created_at DESC").Find(&assets)

	response.Success(c, assets)
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

	asset := models.AssetChange{
		UserID:       userID,
		AssetType:    req.AssetType,
		RelatedUser:  req.RelatedUser,
		Name:         req.Name,
		Amount:       req.Amount,
		InterestRate: req.InterestRate,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		Note:         req.Note,
	}

	h.db.Create(&asset)
	response.Success(c, asset)
}

func (h *AssetHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var asset models.AssetChange
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&asset).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Asset not found")
		return
	}

	response.Success(c, asset)
}

func (h *AssetHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var asset models.AssetChange
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&asset).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Asset not found")
		return
	}

	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		asset.StartDate = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		asset.EndDate = &t
	}

	asset.AssetType = req.AssetType
	asset.RelatedUser = req.RelatedUser
	asset.Name = req.Name
	asset.Amount = req.Amount
	asset.InterestRate = req.InterestRate
	asset.Status = req.Status
	asset.Note = req.Note

	h.db.Save(&asset)
	response.Success(c, asset)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.AssetChange{})
	response.SuccessWithMessage(c, "Asset deleted", nil)
}

func (h *AssetHandler) GetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var debtOwed, debtOwing, investment float64

	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owed", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&debtOwed)

	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owing", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&debtOwing)

	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "investment", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&investment)

	response.Success(c, gin.H{
		"debt_owed":   debtOwed,
		"debt_owing":  debtOwing,
		"investment":  investment,
		"net_worth":   debtOwed - debtOwing + investment,
	})
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

type ReminderHandler struct {
	db *gorm.DB
}

func NewReminderHandler(db *gorm.DB) *ReminderHandler {
	return &ReminderHandler{db: db}
}

func (h *ReminderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var reminders []models.Reminder
	h.db.Where("user_id = ?", userID).Order("remind_time ASC").Find(&reminders)

	response.Success(c, reminders)
}

type CreateReminderRequest struct {
	Title      string     `json:"title" binding:"required"`
	Content    string     `json:"content"`
	RemindTime string     `json:"remind_time" binding:"required"`
	RepeatType string     `json:"repeat_type"`
	CategoryID *uuid.UUID `json:"category_id"`
	IsActive   bool       `json:"is_active"`
}

func (h *ReminderHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	remindTime, _ := time.Parse("2006-01-02T15:04", req.RemindTime)

	reminder := models.Reminder{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		RemindTime: remindTime,
		RepeatType: req.RepeatType,
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
	}

	h.db.Create(&reminder)
	response.Success(c, reminder)
}

func (h *ReminderHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var reminder models.Reminder
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&reminder).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Reminder not found")
		return
	}

	remindTime, _ := time.Parse("2006-01-02T15:04", req.RemindTime)
	reminder.Title = req.Title
	reminder.Content = req.Content
	reminder.RemindTime = remindTime
	reminder.RepeatType = req.RepeatType
	reminder.CategoryID = req.CategoryID
	reminder.IsActive = req.IsActive

	h.db.Save(&reminder)
	response.Success(c, reminder)
}

func (h *ReminderHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Reminder{})
	response.SuccessWithMessage(c, "Reminder deleted", nil)
}

type NotificationHandler struct {
	db *gorm.DB
}

func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{db: db}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	 unreadOnly := c.Query("unread")

	var notifications []models.Notification
	query := h.db.Where("user_id = ?", userID)
	
	if unreadOnly == "true" {
		query = query.Where("is_read = ?", false)
	}
	
	query.Order("created_at DESC").Limit(50).Find(&notifications)

	var unreadCount int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	response.Success(c, gin.H{
		"list":       notifications,
		"unread_count": unreadCount,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true)
	response.SuccessWithMessage(c, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	h.db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true)
	response.SuccessWithMessage(c, "All notifications marked as read", nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})
	response.SuccessWithMessage(c, "Notification deleted", nil)
}

type SearchHandler struct {
	db *gorm.DB
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

func (h *SearchHandler) Search(c *gin.Context) {
	userID := middleware.GetUserID(c)
	keyword := c.Query("q")
	searchType := c.DefaultQuery("type", "all")

	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "Keyword is required")
		return
	}

	limit := 20

	var transactions []models.Transaction
	if searchType == "all" || searchType == "transactions" {
		h.db.Where("user_id = ? AND (note ILIKE ? OR merchant ILIKE ?)", userID, "%"+keyword+"%", "%"+keyword+"%").
			Order("bill_date DESC").Limit(limit).Preload("Category").Preload("Account").Find(&transactions)
	}

	var accounts []models.Account
	if searchType == "all" || searchType == "accounts" {
		h.db.Where("user_id = ? AND name ILIKE ?", userID, "%"+keyword+"%").Limit(limit).Find(&accounts)
	}

	var categories []models.Category
	if searchType == "all" || searchType == "categories" {
		h.db.Where("user_id = ? AND name ILIKE ?", userID, "%"+keyword+"%").Limit(limit).Find(&categories)
	}

	response.Success(c, gin.H{
		"transactions": transactions,
		"accounts":     accounts,
		"categories":   categories,
	})
}

type ImportHandler struct {
	db *gorm.DB
}

func NewImportHandler(db *gorm.DB) *ImportHandler {
	return &ImportHandler{db: db}
}

type ImportTransactionRequest struct {
	Date       string  `json:"date" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Category   string  `json:"category"`
	Account    string  `json:"account"`
	Merchant   string  `json:"merchant"`
	Note       string  `json:"note"`
}

func (h *ImportHandler) ImportTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req []ImportTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var imported int
	var failed int

	for _, item := range req {
		billDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			failed++
			continue
		}

		var accountID uuid.UUID
		if item.Account != "" {
			var account models.Account
			if err := h.db.Where("user_id = ? AND name = ?", userID, item.Account).First(&account).Error; err == nil {
				accountID = account.ID
			} else {
				account := models.Account{
					UserID: userID,
					Name:   item.Account,
					Type:   "cash",
				}
				h.db.Create(&account)
				accountID = account.ID
			}
		}

		var categoryID *uuid.UUID
		if item.Category != "" {
			var category models.Category
			if err := h.db.Where("user_id = ? AND name = ?", userID, item.Category).First(&category).Error; err == nil {
				categoryID = &category.ID
			}
		}

		tx := models.Transaction{
			UserID:     userID,
			AccountID:  accountID,
			CategoryID: categoryID,
			Type:       item.Type,
			Amount:     item.Amount,
			Merchant:   item.Merchant,
			Note:       item.Note,
			BillDate:   billDate,
		}

		h.db.Create(&tx)
		imported++
	}

	response.Success(c, gin.H{
		"imported": imported,
		"failed":   failed,
	})
}

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalAccounts int64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Count(&totalAccounts)

	var totalTransactions int64
	h.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&totalTransactions)

	var totalBudgets int64
	h.db.Model(&models.Budget{}).Where("user_id = ?", userID).Count(&totalBudgets)

	var thisMonthIncome, thisMonthExpense float64
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ?", userID, "income", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthIncome)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ?", userID, "expense", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthExpense)

	var unreadNotifications int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadNotifications)

	response.Success(c, gin.H{
		"total_accounts":      totalAccounts,
		"total_transactions":  totalTransactions,
		"total_budgets":       totalBudgets,
		"this_month_income":   thisMonthIncome,
		"this_month_expense":  thisMonthExpense,
		"unread_notifications": unreadNotifications,
	})
}

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

type CreateAAGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Members     []struct {
		Name string `json:"name"`
	} `json:"members"`
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

	var group models.AAGroup
	h.db.First(&group, groupID)
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

type GoalHandler struct {
	db *gorm.DB
}

func NewGoalHandler(db *gorm.DB) *GoalHandler {
	return &GoalHandler{db: db}
}

func (h *GoalHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var goals []models.FinancialGoal
	h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals)

	response.Success(c, goals)
}

type CreateGoalRequest struct {
	Name         string  `json:"name" binding:"required"`
	TargetAmount float64 `json:"target_amount" binding:"required"`
	CurrentAmount float64 `json:"current_amount"`
	Deadline    string  `json:"deadline"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
}

func (h *GoalHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var deadline *time.Time
	if req.Deadline != "" {
		t, _ := time.Parse("2006-01-02", req.Deadline)
		deadline = &t
	}

	goal := models.FinancialGoal{
		UserID:        userID,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
		Deadline:     deadline,
		Category:      req.Category,
		Note:          req.Note,
	}

	h.db.Create(&goal)
	response.Success(c, goal)
}

func (h *GoalHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var goal models.FinancialGoal
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Goal not found")
		return
	}

	if req.Deadline != "" {
		t, _ := time.Parse("2006-01-02", req.Deadline)
		goal.Deadline = &t
	}

	goal.Name = req.Name
	goal.TargetAmount = req.TargetAmount
	goal.CurrentAmount = req.CurrentAmount
	goal.Category = req.Category
	goal.Note = req.Note

	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Status = "completed"
	}

	h.db.Save(&goal)
	response.Success(c, goal)
}

func (h *GoalHandler) AddAmount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	type AddAmountRequest struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	var req AddAmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var goal models.FinancialGoal
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Goal not found")
		return
	}

	goal.CurrentAmount += req.Amount
	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Status = "completed"
	}

	h.db.Save(&goal)
	response.Success(c, goal)
}

func (h *GoalHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.FinancialGoal{})
	response.SuccessWithMessage(c, "Goal deleted", nil)
}

type InsuranceHandler struct {
	db *gorm.DB
}

func NewInsuranceHandler(db *gorm.DB) *InsuranceHandler {
	return &InsuranceHandler{db: db}
}

func (h *InsuranceHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var insurances []models.Insurance
	h.db.Where("user_id = ?", userID).Order("next_payment_date ASC").Find(&insurances)

	response.Success(c, insurances)
}

type CreateInsuranceRequest struct {
	Name             string  `json:"name" binding:"required"`
	InsuranceType   string  `json:"insurance_type"`
	Company         string  `json:"company"`
	Premium         float64 `json:"premium"`
	PaymentType     string  `json:"payment_type"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	Coverage        float64 `json:"coverage"`
	Beneficiary     string  `json:"beneficiary"`
	Note            string  `json:"note"`
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

	var nextPaymentDate *time.Time
	if req.PaymentType == "yearly" {
		nextPaymentDate = &time.Time{}
		*nextPaymentDate = startDate.AddDate(1, 0, 0)
	} else if req.PaymentType == "quarterly" {
		nextPaymentDate = &time.Time{}
		*nextPaymentDate = startDate.AddDate(0, 3, 0)
	} else if req.PaymentType == "monthly" {
		nextPaymentDate = &time.Time{}
		*nextPaymentDate = startDate.AddDate(0, 1, 0)
	}

	insurance := models.Insurance{
		UserID:           userID,
		Name:             req.Name,
		InsuranceType:    req.InsuranceType,
		Company:          req.Company,
		Premium:          req.Premium,
		PaymentType:      req.PaymentType,
		StartDate:        startDate,
		EndDate:          endDate,
		Coverage:         req.Coverage,
		Beneficiary:      req.Beneficiary,
		Note:             req.Note,
		NextPaymentDate:  nextPaymentDate,
	}

	h.db.Create(&insurance)
	response.Success(c, insurance)
}

func (h *InsuranceHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req CreateInsuranceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var insurance models.Insurance
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&insurance).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Insurance not found")
		return
	}

	insurance.Name = req.Name
	insurance.InsuranceType = req.InsuranceType
	insurance.Company = req.Company
	insurance.Premium = req.Premium
	insurance.PaymentType = req.PaymentType
	insurance.Coverage = req.Coverage
	insurance.Beneficiary = req.Beneficiary
	insurance.Note = req.Note

	h.db.Save(&insurance)
	response.Success(c, insurance)
}

func (h *InsuranceHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Insurance{})
	response.SuccessWithMessage(c, "Insurance deleted", nil)
}

func (h *InsuranceHandler) GetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalPremium float64
	var activeCount int64
	var expiringCount int64

	h.db.Model(&models.Insurance{}).Where("user_id = ? AND status = ?", userID, "active").
		Select("COALESCE(SUM(premium), 0)").Scan(&totalPremium)
	h.db.Model(&models.Insurance{}).Where("user_id = ? AND status = ?", userID, "active").Count(&activeCount)

	thirtyDaysLater := time.Now().AddDate(0, 0, 30)
	h.db.Model(&models.Insurance{}).
		Where("user_id = ? AND status = ? AND next_payment_date <= ?", userID, "active", thirtyDaysLater).
		Count(&expiringCount)

	var totalCoverage float64
	h.db.Model(&models.Insurance{}).Where("user_id = ? AND status = ?", userID, "active").
		Select("COALESCE(SUM(coverage), 0)").Scan(&totalCoverage)

	response.Success(c, gin.H{
		"total_premium":   totalPremium,
		"active_count":    activeCount,
		"expiring_count":  expiringCount,
		"total_coverage":  totalCoverage,
	})
}

type NetWorthHandler struct {
	db *gorm.DB
}

func NewNetWorthHandler(db *gorm.DB) *NetWorthHandler {
	return &NetWorthHandler{db: db}
}

func (h *NetWorthHandler) GetNetWorth(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalAssets float64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(balance), 0)").Scan(&totalAssets)

	var totalDebtOwed float64
	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owed", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalDebtOwed)

	var totalDebtOwing float64
	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "debt_owing", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalDebtOwing)

	var totalInvestment float64
	h.db.Model(&models.AssetChange{}).
		Where("user_id = ? AND asset_type = ? AND status = ?", userID, "investment", "active").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalInvestment)

	netWorth := totalAssets + totalDebtOwed - totalDebtOwing

	response.Success(c, gin.H{
		"total_assets":      totalAssets,
		"total_debt_owed":   totalDebtOwed,
		"total_debt_owing":  totalDebtOwing,
		"total_investment":  totalInvestment,
		"net_worth":         netWorth,
	})
}

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

type CurrencyHandler struct {
	db *gorm.DB
}

func NewCurrencyHandler(db *gorm.DB) *CurrencyHandler {
	return &CurrencyHandler{db: db}
}

func (h *CurrencyHandler) GetRates(c *gin.Context) {
	response.Success(c, defaultExchangeRates)
}

func (h *CurrencyHandler) Convert(c *gin.Context) {
	var req struct {
		From   string  `json:"from" binding:"required"`
		To     string  `json:"to" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fromRate, ok := defaultExchangeRates[req.From]
	if !ok {
		response.Error(c, http.StatusBadRequest, "Unsupported currency: "+req.From)
		return
	}

	toRate, ok := defaultExchangeRates[req.To]
	if !ok {
		response.Error(c, http.StatusBadRequest, "Unsupported currency: "+req.To)
		return
	}

	converted := (req.Amount / fromRate) * toRate

	response.Success(c, gin.H{
		"from":       req.From,
		"to":         req.To,
		"amount":     req.Amount,
		"converted": converted,
		"rate":       toRate / fromRate,
	})
}

type Currency struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Symbol string `json:"symbol"`
	Rate  float64 `json:"rate"`
}

func (h *CurrencyHandler) ListCurrencies(c *gin.Context) {
	currencies := []Currency{
		{Code: "CNY", Name: "人民币", Symbol: "¥", Rate: defaultExchangeRates["CNY"]},
		{Code: "USD", Name: "美元", Symbol: "$", Rate: defaultExchangeRates["USD"]},
		{Code: "EUR", Name: "欧元", Symbol: "€", Rate: defaultExchangeRates["EUR"]},
		{Code: "JPY", Name: "日元", Symbol: "¥", Rate: defaultExchangeRates["JPY"]},
		{Code: "GBP", Name: "英镑", Symbol: "£", Rate: defaultExchangeRates["GBP"]},
		{Code: "HKD", Name: "港币", Symbol: "HK$", Rate: defaultExchangeRates["HKD"]},
		{Code: "KRW", Name: "韩元", Symbol: "₩", Rate: defaultExchangeRates["KRW"]},
		{Code: "AUD", Name: "澳元", Symbol: "A$", Rate: defaultExchangeRates["AUD"]},
		{Code: "CAD", Name: "加元", Symbol: "C$", Rate: defaultExchangeRates["CAD"]},
		{Code: "SGD", Name: "新加坡元", Symbol: "S$", Rate: defaultExchangeRates["SGD"]},
	}

	response.Success(c, currencies)
}

// CSV Import Handler
type CSVImportHandler struct {
	db *gorm.DB
}

func NewCSVImportHandler(db *gorm.DB) *CSVImportHandler {
	return &CSVImportHandler{db: db}
}

type CSVTransaction struct {
	Date       string `json:"date"`
	Type       string `json:"type"`
	Amount     string `json:"amount"`
	Category   string `json:"category"`
	Account    string `json:"account"`
	Merchant   string `json:"merchant"`
	Note       string `json:"note"`
}

func (h *CSVImportHandler) ImportCSV(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Transactions []CSVTransaction `json:"transactions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var imported, failed int

	for _, item := range req.Transactions {
		billDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			// Try alternative format
			billDate, err = time.Parse("2006/01/02", item.Date)
			if err != nil {
				failed++
				continue
			}
		}

		amount, err := strconv.ParseFloat(item.Amount, 64)
		if err != nil {
			failed++
			continue
		}

		txType := "expense"
		if item.Type == "收入" || item.Type == "income" {
			txType = "income"
		} else if item.Type == "转账" || item.Type == "transfer" {
			txType = "transfer"
		}

		var accountID uuid.UUID
		if item.Account != "" {
			var account models.Account
			if err := h.db.Where("user_id = ? AND name = ?", userID, item.Account).First(&account).Error; err == nil {
				accountID = account.ID
			} else {
				account := models.Account{
					UserID: userID,
					Name:   item.Account,
					Type:   "cash",
				}
				h.db.Create(&account)
				accountID = account.ID
			}
		}

		var categoryID *uuid.UUID
		if item.Category != "" {
			var category models.Category
			if err := h.db.Where("user_id = ? AND name = ?", userID, item.Category).First(&category).Error; err == nil {
				categoryID = &category.ID
			}
		}

		tx := models.Transaction{
			UserID:     userID,
			AccountID:  accountID,
			CategoryID: categoryID,
			Type:       txType,
			Amount:     amount,
			Merchant:   item.Merchant,
			Note:       item.Note,
			BillDate:   billDate,
		}

		h.db.Create(&tx)
		imported++
	}

	response.Success(c, gin.H{
		"imported": imported,
		"failed":   failed,
	})
}

// Statistics Handler for dashboard
type StatisticsHandler struct {
	db *gorm.DB
}

func NewStatisticsHandler(db *gorm.DB) *StatisticsHandler {
	return &StatisticsHandler{db: db}
}

func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var totalAccounts int64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Count(&totalAccounts)

	var totalTransactions int64
	h.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&totalTransactions)

	var totalBudgets int64
	h.db.Model(&models.Budget{}).Where("user_id = ?", userID).Count(&totalBudgets)

	// This month
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	endOfMonth := time.Now()

	var thisMonthIncome, thisMonthExpense float64
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ?", userID, "income", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthIncome)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ?", userID, "expense", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthExpense)

	// Last month
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfMonth.AddDate(0, 0, -1)

	var lastMonthIncome, lastMonthExpense float64
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ? AND bill_date <= ?", userID, "income", startOfLastMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&lastMonthIncome)
	h.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND bill_date >= ? AND bill_date <= ?", userID, "expense", startOfLastMonth, endOfLastMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&lastMonthExpense)

	// Total balance
	var totalBalance float64
	h.db.Model(&models.Account{}).Where("user_id = ?", userID).Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance)

	// Unread notifications
	var unreadNotifications int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadNotifications)

	// Income change percentage
	incomeChange := 0.0
	if lastMonthIncome > 0 {
		incomeChange = (thisMonthIncome - lastMonthIncome) / lastMonthIncome * 100
	}

	// Expense change percentage
	expenseChange := 0.0
	if lastMonthExpense > 0 {
		expenseChange = (thisMonthExpense - lastMonthExpense) / lastMonthExpense * 100
	}

	response.Success(c, gin.H{
		"total_accounts":        totalAccounts,
		"total_transactions":     totalTransactions,
		"total_budgets":         totalBudgets,
		"total_balance":         totalBalance,
		"this_month_income":     thisMonthIncome,
		"this_month_expense":    thisMonthExpense,
		"last_month_income":     lastMonthIncome,
		"last_month_expense":    lastMonthExpense,
		"income_change":         incomeChange,
		"expense_change":       expenseChange,
		"unread_notifications": unreadNotifications,
	})
}

// UploadHandler handles file uploads
type UploadHandler struct {
	db *gorm.DB
}

func NewUploadHandler(db *gorm.DB) *UploadHandler {
	return &UploadHandler{db: db}
}

type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func (h *UploadHandler) Upload(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, 401, "Unauthorized")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "No file uploaded")
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		response.Error(c, 400, "Invalid file type. Only images are allowed")
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10<<20 {
		response.Error(c, 400, "File too large. Max 10MB")
		return
	}

	// Create upload directory
	userDir := filepath.Join("uploads", userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		response.Error(c, 500, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(userDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		response.Error(c, 500, "Failed to save file")
		return
	}

	// Return URL (in production, this would be a CDN URL)
	url := fmt.Sprintf("/api/v1/uploads/%s", filename)
	response.Success(c, UploadResponse{
		URL:      url,
		Filename: file.Filename,
		Size:     file.Size,
	})
}

// CashFlowHandler handles cash flow projections
type CashFlowHandler struct {
	db *gorm.DB
}

func NewCashFlowHandler(db *gorm.DB) *CashFlowHandler {
	return &CashFlowHandler{db: db}
}

type CashFlowProjection struct {
	Date          string  `json:"date"`
	ProjectedBal  float64 `json:"projected_balance"`
	Income        float64 `json:"income"`
	Expense       float64 `json:"expense"`
	RecurringTx   int     `json:"recurring_transactions"`
}

func (h *CashFlowHandler) GetProjection(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, 401, "Unauthorized")
		return
	}

	uid, _ := uuid.Parse(userID)

	// Get days parameter (default 30)
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	// Get current total balance
	var totalBalance float64
	if err := h.db.Model(&models.Account{}).Where("user_id = ?", uid).Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance).Error; err != nil {
		response.Error(c, 500, "Failed to get balance")
		return
	}

	// Get recurring transactions
	var recurringTxs []models.Transaction
	if err := h.db.Where("user_id = ? AND is_recurring = ?", uid, true).Find(&recurringTxs).Error; err != nil {
		response.Error(c, 500, "Failed to get recurring transactions")
		return
	}

	// Calculate projections
	projections := make([]CashFlowProjection, days)
	currentBalance := totalBalance

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		dayIncome := 0.0
		dayExpense := 0.0
		recurringCount := 0

		for _, tx := range recurringTxs {
			if tx.BillDate.Day() == date.Day() || shouldRecurToday(&tx.BillDate, date) {
				if tx.Type == "income" {
					dayIncome += tx.Amount
				} else if tx.Type == "expense" {
					dayExpense += tx.Amount
				}
				recurringCount++
			}
		}

		currentBalance = currentBalance + dayIncome - dayExpense
		projections[i] = CashFlowProjection{
			Date:          dateStr,
			ProjectedBal:  currentBalance,
			Income:        dayIncome,
			Expense:       dayExpense,
			RecurringTx:   recurringCount,
		}
	}

	response.Success(c, gin.H{
		"current_balance": totalBalance,
		"projections":     projections,
	})
}

func shouldRecurToday(lastDate *time.Time, today time.Time) bool {
	if lastDate == nil {
		return false
	}
	// Simple check - if the day matches and it's after the last date
	return today.After(*lastDate) && today.Day() == lastDate.Day()
}

// BudgetAlertHandler handles budget alerts
type BudgetAlertHandler struct {
	db *gorm.DB
}

func NewBudgetAlertHandler(db *gorm.DB) *BudgetAlertHandler {
	return &BudgetAlertHandler{db: db}
}

type BudgetAlert struct {
	BudgetID       string  `json:"budget_id"`
	CategoryName   string  `json:"category_name"`
	BudgetAmount   float64 `json:"budget_amount"`
	SpentAmount    float64 `json:"spent_amount"`
	Remaining     float64 `json:"remaining"`
	AlertType     string  `json:"alert_type"` // warning, exceeded
	Percentage    float64 `json:"percentage"`
}

func (h *BudgetAlertHandler) GetAlerts(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, 401, "Unauthorized")
		return
	}

	uid, _ := uuid.Parse(userID)

	// Get current month date range
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// Get all budgets
	var budgets []models.Budget
	if err := h.db.Where("user_id = ?", uid).Find(&budgets).Error; err != nil {
		response.Error(c, 500, "Failed to get budgets")
		return
	}

	alerts := make([]BudgetAlert, 0)

	for _, budget := range budgets {
		// Calculate spent amount
		var spent float64
		query := h.db.Model(&models.Transaction{}).
			Where("user_id = ? AND type = ? AND bill_date BETWEEN ? AND ?",
				uid, "expense", startOfMonth, endOfMonth)

		if budget.CategoryID != nil {
			query = query.Where("category_id = ?", budget.CategoryID)
		}

		if err := query.Select("COALESCE(SUM(amount * exchange_rate), 0)").Scan(&spent).Error; err != nil {
			continue
		}

		percentage := 0.0
		if budget.Amount > 0 {
			percentage = (spent / budget.Amount) * 100
		}

		// Check if alert needed
		alertType := ""
		if percentage >= 100 {
			alertType = "exceeded"
		} else if percentage >= budget.AlertThreshold*100 {
			alertType = "warning"
		}

		if alertType != "" {
			// Get category name
			categoryName := "所有类别"
			if budget.CategoryID != nil {
				var cat models.Category
				if err := h.db.First(&cat, budget.CategoryID).Error; err == nil {
					categoryName = cat.Name
				}
			}

			alerts = append(alerts, BudgetAlert{
				BudgetID:     budget.ID.String(),
				CategoryName: categoryName,
				BudgetAmount: budget.Amount,
				SpentAmount:  spent,
				Remaining:    budget.Amount - spent,
				AlertType:    alertType,
				Percentage:   percentage,
			})
		}
	}

	response.Success(c, alerts)
}
