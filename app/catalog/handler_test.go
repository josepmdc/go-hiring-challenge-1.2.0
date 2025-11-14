package catalog_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/models"
)

func TestGetAll(t *testing.T) {
	testProduct := models.Product{
		ID:    123,
		Code:  "ASDF",
		Price: decimal.NewFromFloat(42.256),
		// add a variant to make sure it doesn't show up in the JSON response
		Variants: []models.Variant{{ID: 321}},
		Category: &models.Category{
			ID:   432,
			Code: "CLOTHING",
			Name: "Category 1",
		},
	}

	totalCount := int64(123)

	// we'll assert this JSON string so we catch breaking changes in the API
	expectedProductsResponse := `
	{
		"products": [
			{
				"code":  "ASDF",
				"price": "42.256",
				"category": { "code": "CLOTHING", "name": "Category 1" }
			}
		],
		"totalCount": 123
	}`

	t.Run("given no pagination params, the default values should be used", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.Equal(t, 0, params.Offset)
			assert.Equal(t, 10, params.Limit)
			return []models.Product{testProduct}, totalCount, nil
		})

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		res := httptest.NewRecorder()

		catalog.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedProductsResponse, res.Body.String())
	})

	t.Run("given pagination params, they should be used to paginate the results", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.Equal(t, 45, params.Offset)
			assert.Equal(t, 5, params.Limit)
			return []models.Product{testProduct}, totalCount, nil
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/catalog?limit=5&offset=45", nil)

		catalog.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedProductsResponse, res.Body.String())
	})

	t.Run("given an invalid offset parameter, an error should be returned", func(t *testing.T) {
		tests := []struct {
			name         string
			offset       string
			errorMessage string
		}{
			{
				name:         "not a number",
				offset:       "NaN",
				errorMessage: "invalid offset param 'NaN': must be a number",
			},
			{
				name:         "negative offset",
				offset:       "-5",
				errorMessage: "invalid offset param '-5': must be positive",
			},
		}

		for _, tt := range tests {
			svc := NewMockService(t)

			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/catalog?offset="+tt.offset, nil)

			catalog.NewHandler(svc).HandleGet(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.JSONEq(t, `{ "error": "`+tt.errorMessage+`"}`, res.Body.String())
		}
	})

	t.Run("given an invalid limit parameter, and error should be returned", func(t *testing.T) {
		tests := []struct {
			name         string
			limit        string
			errorMessage string
		}{
			{
				name:         "not a number",
				limit:        "NaN",
				errorMessage: "invalid limit param 'NaN': must be a number",
			},
			{
				name:         "negative limit",
				limit:        "-5",
				errorMessage: "invalid limit param '-5': must be greater than 0",
			},
			{
				name:         "zero limit",
				limit:        "0",
				errorMessage: "invalid limit param '0': must be greater than 0",
			},
			{
				name:         "limit greater than max",
				limit:        "101",
				errorMessage: "invalid limit param '101': must not exceed 100",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc := NewMockService(t)

				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/catalog?limit="+tt.limit, nil)

				catalog.NewHandler(svc).HandleGet(res, req)
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.JSONEq(t, `{ "error": "`+tt.errorMessage+`" }`, res.Body.String())
			})
		}
	})

	t.Run("given a valid priceLt parameter, it should send it to the service", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.NotNil(t, params.PriceLt)
			assert.Equal(t, decimal.NewFromFloat(65.432), *params.PriceLt)
			return []models.Product{}, 0, nil
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/catalog?priceLt=65.432", nil)

		catalog.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{ "products": [], "totalCount": 0 }`, res.Body.String())
	})

	t.Run("given an invalid priceLt parameter, and error should be returned", func(t *testing.T) {
		tests := []struct {
			name         string
			priceLt      string
			errorMessage string
		}{
			{
				name:         "not a number",
				priceLt:      "NaN",
				errorMessage: "invalid priceLt param 'NaN': must be a number",
			},
			{
				name:         "negative priceLt",
				priceLt:      "-5",
				errorMessage: "invalid priceLt param '-5': must be greater than 0",
			},
			{
				name:         "zero priceLt",
				priceLt:      "0",
				errorMessage: "invalid priceLt param '0': must be greater than 0",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc := NewMockService(t)

				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/catalog?priceLt="+tt.priceLt, nil)

				catalog.NewHandler(svc).HandleGet(res, req)
				assert.Equal(t, http.StatusBadRequest, res.Code)
				assert.JSONEq(t, `{ "error": "`+tt.errorMessage+`" }`, res.Body.String())
			})
		}
	})

	t.Run("given a valid category code parameter, it should send it to the service", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.NotNil(t, params.CategoryCode)
			assert.Equal(t, "CLOTHING", *params.CategoryCode)
			return []models.Product{}, 0, nil
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/catalog?categoryCode=CLOTHING", nil)

		catalog.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{ "products": [], "totalCount": 0 }`, res.Body.String())
	})

	t.Run("given both priceLt and categoryCode filters, they are passed correctly", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.NotNil(t, params.PriceLt)
			assert.Equal(t, decimal.NewFromInt(100), *params.PriceLt)
			assert.NotNil(t, params.CategoryCode)
			assert.Equal(t, "CLOTHING", *params.CategoryCode)
			return []models.Product{}, 0, nil
		})

		req := httptest.NewRequest(http.MethodGet, "/catalog?priceLt=100&categoryCode=CLOTHING", nil)
		res := httptest.NewRecorder()

		catalog.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{ "products": [], "totalCount": 0 }`, res.Body.String())
	})

	t.Run("given the service returns an error, 500 is returned", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			return nil, 0, fmt.Errorf("something went wrong")
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)

		catalog.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, `{ "error": "something went wrong" }`, res.Body.String())
	})

	t.Run("given a product with no category, the category field in the response should be null", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllProducts(func(params *models.CatalogParams) ([]models.Product, int64, error) {
			assert.Equal(t, 0, params.Offset)
			assert.Equal(t, 10, params.Limit)
			return []models.Product{{ID: 123, Code: "ASDF", Price: decimal.NewFromFloat(42.256)}}, totalCount, nil
		})

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		res := httptest.NewRecorder()

		catalog.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `
			{
				"products": [
					{
						"code":  "ASDF",
						"price": "42.256",
						"category": null
					}
				],
				"totalCount": 123
			}`,
			res.Body.String(),
		)
	})
}

type MockService struct {
	t *testing.T

	getAllProductsMock func(*models.CatalogParams) ([]models.Product, int64, error)
}

func NewMockService(t *testing.T) *MockService {
	return &MockService{t: t}
}

func (svc *MockService) MockGetAllProducts(mock func(params *models.CatalogParams) ([]models.Product, int64, error)) {
	svc.getAllProductsMock = mock
}

func (svc *MockService) GetAllProducts(params *models.CatalogParams) ([]models.Product, int64, error) {
	assert.NotNil(svc.t, svc.getAllProductsMock, "unexpected call to GetAllProducts")
	return svc.getAllProductsMock(params)
}
