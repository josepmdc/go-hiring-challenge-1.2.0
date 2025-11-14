package categories_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/models"
)

func TestService_GetAllCategories(t *testing.T) {
	t.Run("given a request to get all categories, all categories should be returned", func(t *testing.T) {
		repo := NewMockCategoriesRepo(t)

		expected := []models.Category{
			{ID: 123, Code: "CAT1", Name: "Category 1"},
			{ID: 234, Code: "CAT2", Name: "Category 2"},
			{ID: 345, Code: "CAT3", Name: "Category 3"},
		}

		repo.MockGetAllCategories(func() ([]models.Category, error) {
			return expected, nil
		})

		actual, err := categories.NewService(repo).GetAllCategories()
		require.NoError(t, err)
		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("given a the repo returns an error, an error should be returned", func(t *testing.T) {
		repo := NewMockCategoriesRepo(t)

		expectedErr := errors.New("something happened in the DB")
		repo.MockGetAllCategories(func() ([]models.Category, error) {
			return nil, expectedErr
		})

		actual, err := categories.NewService(repo).GetAllCategories()
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, actual)
	})
}

func TestService_CreateCategory(t *testing.T) {
	req := categories.CreateCategoryReq{
		Code: "NEWCAT",
		Name: "New Category",
	}

	t.Run("given a valid request, it should create the category", func(t *testing.T) {
		mockRepo := NewMockCategoriesRepo(t)

		svc := categories.NewService(mockRepo)

		mockRepo.MockGetByCode(func(code string) (*models.Category, error) {
			assert.Equal(t, req.Code, code)
			return nil, models.ErrNotFound
		})

		mockRepo.MockCreateCategory(func(cat *models.Category) error {
			assert.Equal(t, req.Code, cat.Code)
			assert.Equal(t, req.Name, cat.Name)
			cat.ID = 123 // ID comes from the DB
			return nil
		})

		expected := &models.Category{ID: 123, Code: req.Code, Name: req.Name}

		actual, err := svc.CreateCategory(req)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("given a duplicate code, it should return ErrDuplicateCode", func(t *testing.T) {
		mockRepo := NewMockCategoriesRepo(t)

		svc := categories.NewService(mockRepo)

		existing := &models.Category{Code: req.Code, Name: "Exists"}

		mockRepo.MockGetByCode(func(code string) (*models.Category, error) {
			assert.Equal(t, req.Code, code)
			return existing, nil
		})

		actual, err := svc.CreateCategory(req)
		assert.ErrorIs(t, err, categories.ErrDuplicateCode)
		assert.Nil(t, actual)
	})

	t.Run("given a get by code returns an unexpected error, it should return the error", func(t *testing.T) {
		mockRepo := NewMockCategoriesRepo(t)

		svc := categories.NewService(mockRepo)

		expectedErr := errors.New("something happened in the DB")

		mockRepo.MockGetByCode(func(code string) (*models.Category, error) {
			assert.Equal(t, req.Code, code)
			return nil, expectedErr
		})

		actual, err := svc.CreateCategory(req)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, actual)
	})

	t.Run("given create category returns an error, it should return that error", func(t *testing.T) {
		mockRepo := NewMockCategoriesRepo(t)

		svc := categories.NewService(mockRepo)

		mockRepo.MockGetByCode(func(code string) (*models.Category, error) {
			assert.Equal(t, req.Code, code)
			return nil, models.ErrNotFound
		})

		expectedErr := errors.New("something happened in the DB")

		mockRepo.MockCreateCategory(func(cat *models.Category) error {
			assert.Equal(t, req.Code, cat.Code)
			assert.Equal(t, req.Name, cat.Name)
			return expectedErr
		})

		actual, err := svc.CreateCategory(req)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, actual)
	})
}

type MockCategoriesRepo struct {
	t *testing.T

	mockGetAllCategories func() ([]models.Category, error)
	mockGetByCode        func(code string) (*models.Category, error)
	mockCreateCategory   func(category *models.Category) error
}

func NewMockCategoriesRepo(t *testing.T) *MockCategoriesRepo {
	return &MockCategoriesRepo{t: t}
}

func (m *MockCategoriesRepo) MockGetAllCategories(mock func() ([]models.Category, error)) {
	m.mockGetAllCategories = mock
}

func (m *MockCategoriesRepo) GetAllCategories() ([]models.Category, error) {
	assert.NotNil(m.t, m.mockGetAllCategories, "unexpected call to GetAllCategories")
	return m.mockGetAllCategories()
}

func (m *MockCategoriesRepo) MockGetByCode(mock func(code string) (*models.Category, error)) {
	m.mockGetByCode = mock
}

func (m *MockCategoriesRepo) GetByCode(code string) (*models.Category, error) {
	assert.NotNil(m.t, m.mockGetByCode, "unexpected call to GetByCode")
	return m.mockGetByCode(code)
}

func (m *MockCategoriesRepo) MockCreateCategory(mock func(cat *models.Category) error) {
	m.mockCreateCategory = mock
}

func (m *MockCategoriesRepo) CreateCategory(cat *models.Category) error {
	assert.NotNil(m.t, m.mockCreateCategory, "unexpected call to CreateCategory")
	return m.mockCreateCategory(cat)
}
