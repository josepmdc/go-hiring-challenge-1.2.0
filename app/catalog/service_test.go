package catalog_test

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/lib/fn"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetAllProducts(t *testing.T) {
	t.Run("given all params, they should be sent to the repo", func(t *testing.T) {
		repo := NewMockCatalogRepo(t)

		expectedParams := models.CatalogParams{
			CategoryCode: fn.Ptr("CODE"),
			PriceLt:      fn.Ptr(decimal.NewFromFloat(3.2)),
			PaginationParams: &models.PaginationParams{
				Offset: 25,
				Limit:  5,
			},
		}

		expectedProducts := []models.Product{{
			ID:         123,
			Code:       "CODE",
			Price:      decimal.NewFromFloat(3.2),
			Variants:   []models.Variant{},
			CategoryID: 321,
			Category: &models.Category{
				ID:   321,
				Code: "CAT1",
				Name: "Category 1",
			},
		}}

		expectedCount := int64(300)

		repo.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.Equal(t, expectedParams.Limit, params.Limit)
			assert.Equal(t, expectedParams.Offset, params.Offset)
			assert.Equal(t, expectedParams.CategoryCode, params.CategoryCode)
			assert.Equal(t, expectedParams.PriceLt, params.PriceLt)
			return expectedProducts, expectedCount, nil
		})

		actualProducts, actualCount, err := catalog.NewService(repo).GetAllProducts(&expectedParams)
		require.NoError(t, err)
		assert.ElementsMatch(t, expectedProducts, actualProducts)
		assert.Equal(t, expectedCount, actualCount)
	})

	t.Run("given repo returns an error, an error should be returned", func(t *testing.T) {
		repo := NewMockCatalogRepo(t)

		expectedParams := models.CatalogParams{
			CategoryCode: fn.Ptr("CODE"),
			PriceLt:      fn.Ptr(decimal.NewFromFloat(3.2)),
			PaginationParams: &models.PaginationParams{
				Offset: 25,
				Limit:  5,
			},
		}

		expectedError := errors.New("something happened")

		repo.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.Equal(t, expectedParams.Limit, params.Limit)
			assert.Equal(t, expectedParams.Offset, params.Offset)
			assert.Equal(t, expectedParams.CategoryCode, params.CategoryCode)
			assert.Equal(t, expectedParams.PriceLt, params.PriceLt)
			return nil, 0, expectedError
		})

		actualProducts, actualCount, err := catalog.NewService(repo).GetAllProducts(&expectedParams)
		assert.ErrorIs(t, err, expectedError)
		assert.Nil(t, actualProducts)
		assert.Equal(t, int64(0), actualCount)
	})
}

func TestService_GetProductByCode(t *testing.T) {
	t.Run("given variants with zero price, they inherit the product price", func(t *testing.T) {
		mockRepo := NewMockCatalogRepo(t)

		productPrice := decimal.NewFromFloat(43.21)
		variant1 := models.Variant{Price: decimal.Zero}
		variant2 := models.Variant{Price: decimal.NewFromFloat(50.43)}

		mockRepo.MockGetProductByCode(func(code string) (*models.Product, error) {
			return &models.Product{
				Code:     "P1",
				Price:    productPrice,
				Variants: []models.Variant{variant1, variant2},
			}, nil
		})

		svc := catalog.NewService(mockRepo)
		result, err := svc.GetProductByCode("P1")

		assert.NoError(t, err)
		assert.Equal(t, productPrice, result.Variants[0].Price)
		assert.Equal(t, decimal.NewFromFloat(50.43), result.Variants[1].Price)
	})

	t.Run("given the repo returns an error, the error is returned", func(t *testing.T) {
		mockRepo := NewMockCatalogRepo(t)

		expectedError := errors.New("oops")

		mockRepo.MockGetProductByCode(func(code string) (*models.Product, error) {
			return nil, expectedError
		})

		result, err := catalog.NewService(mockRepo).GetProductByCode("P1")
		assert.ErrorIs(t, err, expectedError)
		assert.Nil(t, result)
	})
}

type MockCatalogRepo struct {
	t *testing.T

	mockGetAllProducts   func(params *models.CatalogParams) ([]models.Product, int64, error)
	mockGetProductByCode func(code string) (*models.Product, error)
}

func NewMockCatalogRepo(t *testing.T) *MockCatalogRepo {
	return &MockCatalogRepo{t: t}
}

func (m *MockCatalogRepo) MockGetAllProducts(fn func(params *models.CatalogParams) ([]models.Product, int64, error)) {
	m.mockGetAllProducts = fn
}

func (m *MockCatalogRepo) GetAllProducts(params *models.CatalogParams) ([]models.Product, int64, error) {
	assert.NotNil(m.t, m.mockGetAllProducts, "unexpected call to GetAllProducts")
	return m.mockGetAllProducts(params)
}

func (m *MockCatalogRepo) MockGetProductByCode(fn func(code string) (*models.Product, error)) {
	m.mockGetProductByCode = fn
}

func (m *MockCatalogRepo) GetProductByCode(code string) (*models.Product, error) {
	assert.NotNil(m.t, m.mockGetProductByCode, "unexpected call to GetProductByCode")
	return m.mockGetProductByCode(code)
}
