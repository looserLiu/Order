package test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountService_Create_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockAccountRepository)

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Account")).Return(nil)

	svc := service.NewAccountService(mockRepo)

	req := &service.AccountCreateRequest{
		Name:     "Test Account",
		Type:     "bank",
		Currency: "CNY",
		Icon:     "bank-icon",
		Color:    "#123456",
	}

	account, err := svc.Create(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, "Test Account", account.Name)
	assert.Equal(t, "bank", account.Type)
	assert.Equal(t, 0.0, account.Balance)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_List_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockAccountRepository)

	accounts := []models.Account{
		{ID: uuid.New(), UserID: userID, Name: "Account 1", Type: "bank", Balance: 1000.0},
		{ID: uuid.New(), UserID: userID, Name: "Account 2", Type: "cash", Balance: 500.0},
	}

	mockRepo.On("FindByUserID", ctx, userID).Return(accounts, nil)

	svc := service.NewAccountService(mockRepo)

	result, err := svc.List(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Account 1", result[0].Name)
	assert.Equal(t, "Account 2", result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_Get_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	mockRepo := new(MockAccountRepository)

	account := &models.Account{
		ID:     accountID,
		UserID: userID,
		Name:   "Test Account",
		Type:   "bank",
	}

	mockRepo.On("FindByID", ctx, accountID).Return(account, nil)

	svc := service.NewAccountService(mockRepo)

	result, err := svc.Get(ctx, accountID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, accountID, result.ID)
	assert.Equal(t, userID, result.UserID)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_Get_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	accountID := uuid.New()

	mockRepo := new(MockAccountRepository)

	account := &models.Account{
		ID:     accountID,
		UserID: otherUserID,
		Name:   "Other User Account",
		Type:   "bank",
	}

	mockRepo.On("FindByID", ctx, accountID).Return(account, nil)

	svc := service.NewAccountService(mockRepo)

	result, err := svc.Get(ctx, accountID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "unauthorized access", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAccountService_Get_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	mockRepo := new(MockAccountRepository)

	mockRepo.On("FindByID", ctx, accountID).Return(nil, errors.New("not found"))

	svc := service.NewAccountService(mockRepo)

	result, err := svc.Get(ctx, accountID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_Update_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	mockRepo := new(MockAccountRepository)

	account := &models.Account{
		ID:       accountID,
		UserID:   userID,
		Name:     "Old Name",
		Type:     "bank",
		Balance:  1000.0,
		Currency: "CNY",
	}

	mockRepo.On("FindByID", ctx, accountID).Return(account, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Account")).Return(nil)

	svc := service.NewAccountService(mockRepo)

	updates := map[string]interface{}{
		"name":     "New Name",
		"currency": "USD",
	}

	result, err := svc.Update(ctx, accountID, userID, updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "USD", result.Currency)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_Update_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	accountID := uuid.New()

	mockRepo := new(MockAccountRepository)

	account := &models.Account{
		ID:     accountID,
		UserID: otherUserID,
		Name:   "Other User Account",
	}

	mockRepo.On("FindByID", ctx, accountID).Return(account, nil)

	svc := service.NewAccountService(mockRepo)

	updates := map[string]interface{}{
		"name": "New Name",
	}

	result, err := svc.Update(ctx, accountID, userID, updates)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "unauthorized access", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAccountService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	mockRepo := new(MockAccountRepository)

	account := &models.Account{
		ID:     accountID,
		UserID: userID,
		Name:   "Test Account",
	}

	mockRepo.On("FindByID", ctx, accountID).Return(account, nil)
	mockRepo.On("Delete", ctx, accountID).Return(nil)

	svc := service.NewAccountService(mockRepo)

	err := svc.Delete(ctx, accountID, userID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_Delete_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	accountID := uuid.New()

	mockRepo := new(MockAccountRepository)

	account := &models.Account{
		ID:     accountID,
		UserID: otherUserID,
		Name:   "Other User Account",
	}

	mockRepo.On("FindByID", ctx, accountID).Return(account, nil)

	svc := service.NewAccountService(mockRepo)

	err := svc.Delete(ctx, accountID, userID)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized access", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestAccountService_GetTotalBalance_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockAccountRepository)

	mockRepo.On("GetTotalBalance", ctx, userID).Return(1500.0, nil)

	svc := service.NewAccountService(mockRepo)

	balance, err := svc.GetTotalBalance(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 1500.0, balance)

	mockRepo.AssertExpectations(t)
}