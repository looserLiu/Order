package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBudgetRepository implements repository.BudgetRepository
type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) Create(ctx context.Context, budget *models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Budget, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Budget), args.Error(1)
}

func (m *MockBudgetRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Budget, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Budget), args.Error(1)
}

func (m *MockBudgetRepository) Update(ctx context.Context, budget *models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBudgetRepository) GetProgress(ctx context.Context, id uuid.UUID) (float64, float64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func TestBudgetService_Create_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockBudgetRepository)

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Budget")).Return(nil)

	svc := service.NewBudgetService(mockRepo)

	req := &service.BudgetCreateRequest{
		Amount:    5000.0,
		Period:    "monthly",
		StartDate: time.Now(),
	}

	budget, err := svc.Create(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, budget)
	assert.Equal(t, userID, budget.UserID)
	assert.Equal(t, 5000.0, budget.Amount)
	assert.Equal(t, "monthly", budget.Period)

	mockRepo.AssertExpectations(t)
}

func TestBudgetService_List_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockBudgetRepository)

	budgets := []models.Budget{
		{ID: uuid.New(), UserID: userID, Amount: 5000.0, Period: "monthly"},
		{ID: uuid.New(), UserID: userID, Amount: 10000.0, Period: "yearly"},
	}

	mockRepo.On("FindByUserID", ctx, userID).Return(budgets, nil)

	svc := service.NewBudgetService(mockRepo)

	result, err := svc.List(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 5000.0, result[0].Amount)
	assert.Equal(t, 10000.0, result[1].Amount)

	mockRepo.AssertExpectations(t)
}

func TestBudgetService_Get_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	mockRepo := new(MockBudgetRepository)

	budget := &models.Budget{
		ID:     budgetID,
		UserID: userID,
		Amount: 5000.0,
		Period: "monthly",
	}

	mockRepo.On("FindByID", ctx, budgetID).Return(budget, nil)

	svc := service.NewBudgetService(mockRepo)

	result, err := svc.Get(ctx, budgetID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, budgetID, result.ID)
	assert.Equal(t, userID, result.UserID)

	mockRepo.AssertExpectations(t)
}

func TestBudgetService_Get_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	budgetID := uuid.New()

	mockRepo := new(MockBudgetRepository)

	budget := &models.Budget{
		ID:     budgetID,
		UserID: otherUserID,
		Amount: 5000.0,
	}

	mockRepo.On("FindByID", ctx, budgetID).Return(budget, nil)

	svc := service.NewBudgetService(mockRepo)

	result, err := svc.Get(ctx, budgetID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "unauthorized access", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestBudgetService_Update_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	mockRepo := new(MockBudgetRepository)

	budget := &models.Budget{
		ID:     budgetID,
		UserID: userID,
		Amount: 5000.0,
		Period: "monthly",
	}

	mockRepo.On("FindByID", ctx, budgetID).Return(budget, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Budget")).Return(nil)

	svc := service.NewBudgetService(mockRepo)

	updates := map[string]interface{}{
		"amount": 6000.0,
	}

	result, err := svc.Update(ctx, budgetID, userID, updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 6000.0, result.Amount)

	mockRepo.AssertExpectations(t)
}

func TestBudgetService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	mockRepo := new(MockBudgetRepository)

	budget := &models.Budget{
		ID:     budgetID,
		UserID: userID,
		Amount: 5000.0,
	}

	mockRepo.On("FindByID", ctx, budgetID).Return(budget, nil)
	mockRepo.On("Delete", ctx, budgetID).Return(nil)

	svc := service.NewBudgetService(mockRepo)

	err := svc.Delete(ctx, budgetID, userID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestBudgetService_GetProgress_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	mockRepo := new(MockBudgetRepository)

	budget := &models.Budget{
		ID:     budgetID,
		UserID: userID,
		Amount: 5000.0,
	}

	mockRepo.On("FindByID", ctx, budgetID).Return(budget, nil)
	mockRepo.On("GetProgress", ctx, budgetID).Return(3000.0, 5000.0, nil)

	svc := service.NewBudgetService(mockRepo)

	spent, total, err := svc.GetProgress(ctx, budgetID, userID)

	assert.NoError(t, err)
	assert.Equal(t, 3000.0, spent)
	assert.Equal(t, 5000.0, total)

	mockRepo.AssertExpectations(t)
}