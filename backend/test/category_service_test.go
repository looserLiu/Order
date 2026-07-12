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

func (m *MockCategoryRepository) GetTree(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
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

func TestCategoryService_Create_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockCategoryRepository)

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Category")).Return(nil)

	svc := service.NewCategoryService(mockRepo)

	req := &service.CategoryCreateRequest{
		Name:  "Food",
		Type:  "expense",
		Color: "#FF0000",
	}

	category, err := svc.Create(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, userID, category.UserID)
	assert.Equal(t, "Food", category.Name)
	assert.Equal(t, "expense", category.Type)
	assert.False(t, category.IsSystem)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_List_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockCategoryRepository)

	categories := []models.Category{
		{ID: uuid.New(), UserID: userID, Name: "Food", Type: "expense"},
		{ID: uuid.New(), UserID: userID, Name: "Salary", Type: "income"},
	}

	mockRepo.On("FindByUserID", ctx, userID).Return(categories, nil)

	svc := service.NewCategoryService(mockRepo)

	result, err := svc.List(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Food", result[0].Name)
	assert.Equal(t, "Salary", result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetTree_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := new(MockCategoryRepository)

	categories := []models.Category{
		{ID: uuid.New(), UserID: userID, Name: "Food", Type: "expense"},
	}

	mockRepo.On("GetTree", ctx, userID).Return(categories, nil)

	svc := service.NewCategoryService(mockRepo)

	result, err := svc.GetTree(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Food", result[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mockRepo := new(MockCategoryRepository)

	category := &models.Category{
		ID:     categoryID,
		UserID:   userID,
		Name:     "Old Name",
		Type:     "expense",
	}

	mockRepo.On("FindByID", ctx, categoryID).Return(category, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Category")).Return(nil)

	svc := service.NewCategoryService(mockRepo)

	updates := map[string]interface{}{
		"name": "New Name",
	}

	result, err := svc.Update(ctx, categoryID, userID, updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()

	mockRepo := new(MockCategoryRepository)

	category := &models.Category{
		ID:     categoryID,
		UserID:   otherUserID,
		Name:     "Other User Category",
	}

	mockRepo.On("FindByID", ctx, categoryID).Return(category, nil)

	svc := service.NewCategoryService(mockRepo)

	updates := map[string]interface{}{
		"name": "New Name",
	}

	result, err := svc.Update(ctx, categoryID, userID, updates)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "unauthorized access", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mockRepo := new(MockCategoryRepository)

	category := &models.Category{
		ID:     categoryID,
		UserID:   userID,
		Name:     "Test Category",
	}

	mockRepo.On("FindByID", ctx, categoryID).Return(category, nil)
	mockRepo.On("Delete", ctx, categoryID).Return(nil)

	svc := service.NewCategoryService(mockRepo)

	err := svc.Delete(ctx, categoryID, userID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()

	mockRepo := new(MockCategoryRepository)

	category := &models.Category{
		ID:     categoryID,
		UserID:   otherUserID,
		Name:     "Other User Category",
	}

	mockRepo.On("FindByID", ctx, categoryID).Return(category, nil)

	svc := service.NewCategoryService(mockRepo)

	err := svc.Delete(ctx, categoryID, userID)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized access", err.Error())

	mockRepo.AssertExpectations(t)
}