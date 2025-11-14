package categories

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/lib/fn"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Categories struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	// not adding ID since exercise says it's only for internal use
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesService interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(req CreateCategoryReq) (*models.Category, error)
}

type CategoriesHandler struct {
	svc CategoriesService
}

func NewHandler(svc CategoriesService) *CategoriesHandler {
	return &CategoriesHandler{
		svc: svc,
	}
}

func (h *CategoriesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, Categories{
		Categories: fn.Map(res, toApiCategory),
	})
}

type CreateCategoryBody struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (h *CategoriesHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var body CreateCategoryBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if body.Code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "code is required")
		return
	}

	if body.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "name is required")
		return
	}

	newCategory, err := h.svc.CreateCategory(CreateCategoryReq{
		Code: body.Code,
		Name: body.Name,
	})
	switch {
	case errors.Is(err, ErrDuplicateCode):
		api.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to create category: %s", err))
		return
	case err != nil:
		api.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create category: %s", err))
		return
	}

	api.CreatedResponse(w, Category{
		Code: newCategory.Code,
		Name: newCategory.Name,
	})
}

func toApiCategory(category models.Category) Category {
	return Category{
		Code: category.Code,
		Name: category.Name,
	}
}
