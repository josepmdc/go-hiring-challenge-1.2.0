package categories_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAll(t *testing.T) {
	t.Run("given a request to get all categories, all categories should be returned", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllCategories(func() ([]models.Category, error) {
			return []models.Category{
				{ID: 432, Code: "CLOTHING", Name: "Category 1"},
				{ID: 234, Code: "SHOES", Name: "Category 2"},
			}, nil
		})

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		res := httptest.NewRecorder()

		categories.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `
			{
				"categories": [
					{ "code": "CLOTHING", "name": "Category 1" },
					{ "code": "SHOES", "name": "Category 2" }
				]
			}`,
			res.Body.String(),
		)
	})

	t.Run("given a request to get all categories and the service returns an error, 500 should be returned", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockGetAllCategories(func() ([]models.Category, error) {
			return nil, fmt.Errorf("failed to get categories")
		})

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		res := httptest.NewRecorder()

		categories.NewHandler(svc).HandleGet(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, `{ "error": "failed to get categories" }`, res.Body.String())
	})
}

func TestCreate(t *testing.T) {
	t.Run("given a valid request to create a category, it should return 201 and the created category", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockCreateCategory(func(req categories.CreateCategoryReq) (*models.Category, error) {
			return &models.Category{ID: 432, Code: "CLOTHING", Name: "Category 1"}, nil
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/categories",
			strings.NewReader(`{ "code": "CLOTHING", "name": "Category 1" }`),
		)
		res := httptest.NewRecorder()

		categories.NewHandler(svc).HandlePost(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.JSONEq(t, `{ "code": "CLOTHING", "name": "Category 1" }`, res.Body.String())
	})

	t.Run("given an invalid JSON body, it should return 400", func(t *testing.T) {
		svc := NewMockService(t)

		req := httptest.NewRequest(
			http.MethodPost,
			"/categories",
			strings.NewReader(`Hello World`),
		)
		res := httptest.NewRecorder()

		categories.NewHandler(svc).HandlePost(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{ "error": "invalid JSON body" }`, res.Body.String())
	})

	t.Run("given a valid JSON that doesn not match the schema, it should return 400", func(t *testing.T) {
		tests := []struct {
			body  string
			error string
		}{
			{body: `{"test": "A"}`, error: "code is required"},
			{body: `{"code": "A"}`, error: "name is required"},
		}

		for _, tt := range tests {
			svc := NewMockService(t)

			req := httptest.NewRequest(
				http.MethodPost,
				"/categories",
				strings.NewReader(tt.body),
			)
			res := httptest.NewRecorder()

			categories.NewHandler(svc).HandlePost(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.JSONEq(t, `{ "error": "`+tt.error+`" }`, res.Body.String())
		}
	})

	t.Run("given a duplicate code, it should return 400", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockCreateCategory(func(req categories.CreateCategoryReq) (*models.Category, error) {
			return nil, categories.ErrDuplicateCode
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/categories",
			strings.NewReader(`{ "code": "CLOTHING", "name": "Category 1" }`),
		)
		res := httptest.NewRecorder()

		categories.NewHandler(svc).HandlePost(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{ "error": "failed to create category: category code must be unique" }`, res.Body.String())
	})

	t.Run("given an unkown errror, it should return 500", func(t *testing.T) {
		svc := NewMockService(t)

		svc.MockCreateCategory(func(req categories.CreateCategoryReq) (*models.Category, error) {
			return nil, errors.New("something bad happened")
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/categories",
			strings.NewReader(`{ "code": "CLOTHING", "name": "Category 1" }`),
		)
		res := httptest.NewRecorder()

		categories.NewHandler(svc).HandlePost(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, `{ "error": "failed to create category: something bad happened" }`, res.Body.String())
	})
}

type MockService struct {
	t                    *testing.T
	getAllCategoriesMock func() ([]models.Category, error)
	createCategoryMock   func(categories.CreateCategoryReq) (*models.Category, error)
}

func NewMockService(t *testing.T) *MockService {
	return &MockService{t: t}
}

func (svc *MockService) MockGetAllCategories(mock func() ([]models.Category, error)) {
	svc.getAllCategoriesMock = mock
}

func (svc *MockService) MockCreateCategory(mock func(categories.CreateCategoryReq) (*models.Category, error)) {
	svc.createCategoryMock = mock
}

func (svc *MockService) GetAllCategories() ([]models.Category, error) {
	require.NotNil(svc.t, svc.getAllCategoriesMock, "unexpected call to GetAllCategories")
	return svc.getAllCategoriesMock()
}

func (svc *MockService) CreateCategory(req categories.CreateCategoryReq) (*models.Category, error) {
	require.NotNil(svc.t, svc.createCategoryMock, "unexpected call to CreateCategory")
	return svc.createCategoryMock(req)
}
