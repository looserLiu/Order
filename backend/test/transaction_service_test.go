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

// MockTransactionRepository implements repository.TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx interface{}, transaction *models.Transaction) error {
	args := m.Called(ctx, tx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.Transaction, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, tx interface{}, transaction *models.Transaction) error {
	args := m.Called(ctx, tx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, tx interface{}, id uuid.UUID) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockTransactionRepository) BatchDelete(ctx context.Context, tx interface{}, ids []uuid.UUID) error {
	args := m.Called(ctx, tx, ids)
	return args.Error(0)
}

// MockAccountRepository implements repository.AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *models.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Account, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(ctx context.Context, account *models.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAccountRepository) GetTotalBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

// MockCategoryRepository implements repository.CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTransactionService_Create_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()
	categoryID := uuid.New()

	mockTxRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	mockCategoryRepo := new(MockCategoryRepository)

	account := &models.Account{
		ID:     accountID,
		UserID: userID,
		Name:   "Test Account",
		Type:   "bank",
	}

	category := &models.Category{
		ID:     categoryID,
		UserID: userID,
		Name:   "Test Category",
		Type:   "expense",
	}

	mockAccountRepo.On("FindByID", ctx, accountID).Return(account, nil)
	mockCategoryRepo.On("FindByID", ctx, categoryID).Return(category, nil)
	mockTxRepo.On("Create", ctx, nil, mock.AnythingOfType("*models.Transaction")).Return(nil)

	svc := service.NewTransactionService(mockTxRepo, mockAccountRepo, mockCategoryRepo)

	req := &service.CreateRequest{
		AccountID:  accountID,
		CategoryID: &categoryID,
		Type:       "expense",
		Amount:     100.0,
		Currency:   "CNY",
		BillDate:   time.Now(),
	}

	transaction, err := svc.Create(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, userID, transaction.UserID)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, 100.0, transaction.Amount)

	mockTxRepo.AssertExpectations(t)
	mockAccountRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestTransactionService_Create_AccountNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	mockTxRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	mockCategoryRepo := new(MockCategoryRepository)

	mockAccountRepo.On("FindByID", ctx, accountID).Return(nil, errors.New("not found"))

	svc := service.NewTransactionService(mockTxRepo, mockAccountRepo, mockCategoryRepo)

	req := &service.CreateRequest{
		AccountID: accountID,
		Type:      "expense",
		Amount:    100.0,
		BillDate:  time.Now(),
	}

	transaction, err := svc.Create(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "account not found", err.Error())

	mockAccountRepo.AssertExpectations(t)
}

func TestTransactionService_Create_UnauthorizedAccount(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	accountID := uuid.New()

	mockTxRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	mockCategoryRepo := new(MockCategoryRepository)

	account := &models.Account{
		ID:     accountID,
		UserID: otherUserID,
		Name:   "Other User Account",
		Type:   "bank",
	}

	mockAccountRepo.On("FindByID", ctx, accountID).Return(account, nil)

	svc := service.NewTransactionService(mockTxRepo, mockAccountRepo, mockCategoryRepo)

	req := &service.CreateRequest{
		AccountID: accountID,
		Type:      "expense",
		Amount:    100.0,
		BillDate:  time.Now(),
	}

	transaction, err := svc.Create(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "unauthorized account access", err.Error())

	mockAccountRepo.AssertExpectations(t)
}

func TestTransactionService_Get_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mockTxRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	mockCategoryRepo := new(MockCategoryRepository)

	transaction := &models.Transaction{
		ID:      transactionID,
		UserID:  userID,
		AccountID: uuid.New(),
		Type:    "expense",
		Amount:  100.0,
	}

	mockTxRepo.On("FindByID", ctx, transactionID).Return(transaction, nil)

	svc := service.NewTransactionService(mockTxRepo, mockAccountRepo, mockCategoryRepo)

	result, err := svc.Get(ctx, transactionID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transactionID, result.ID)

	mockTxRepo.AssertExpectations(t)
}

func TestTransactionService_Get_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()

	mockTxRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	mockCategoryRepo := new(MockCategoryRepository)

	transaction := &models.Transaction{
		ID:      transactionID,
		UserID:  otherUserID,
		AccountID: uuid.New(),
		Type:    "expense",
		Amount:  100.0,
	}

	mockTxRepo.On("FindByID", ctx, transactionID).Return(transaction, nil)

	svc := service.NewTransactionService(mockTxRepo, mockAccountRepo, mockCategoryRepo)

	result, err := svc.Get(ctx, transactionID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "unauthorized access", err.Error())

	mockTxRepo.AssertExpectations(t)
}

func TestTransactionService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mockTxRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	mockCategoryRepo := new(MockCategoryRepository)

	transaction := &models.Transaction{
		ID:      transactionID,
		UserID:  userID,
		AccountID: uuid.New(),
		Type:    "expense",
		Amount:  100.0,
	}

	mockTxRepo.On("FindByID", ctx, transactionID).Return(transaction, nil)
	mockTxRepo.On("Delete", ctx, nil, transactionID).Return(nil)

	svc := service.NewTransactionService(mockTxRepo, mockAccountRepo, mockCategoryRepo)

	err := svc.Delete(ctx, transactionID, userID)

	assert.NoError(t, err)

	mockTxRepo.AssertExpectations(t)
}

func TestTransactionService_Delete_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()

	mockTxRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	mockCategoryRepo := new(MockCategoryRepository)

	transaction := &models.Transaction{
		ID:      transactionID,
		UserID:  otherUserID,
		AccountID: uuid.New(),
		Type:    "expense",
		Amount:  100.0,
	}

	mockTxRepo.On("FindByID", ctx, transactionID).Return(transaction, nil)

	svc := service.NewTransactionService(mockTxRepo, mockAccountRepo, mockCategoryRepo)

	err := svc.Delete(ctx, transactionID, userID)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized access", err.Error())

	mockTxRepo.AssertExpectations(t)
}