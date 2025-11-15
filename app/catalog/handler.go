package catalog

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/lib/fn"
	"github.com/mytheresa/go-hiring-challenge/models"
)

const (
	defaultOffset = 0
	defaultLimit  = 10
	maxLimit      = 100
)

type Response struct {
	Products   []Product `json:"products"`
	TotalCount int64     `json:"totalCount"`
}

type Product struct {
	Code     string    `json:"code"`
	Price    string    `json:"price"` // string instead of float to not lose precision
	Category *Category `json:"category"`
}

type ProductDetailed struct {
	Code     string    `json:"code"`
	Price    string    `json:"price"` // string instead of float to not lose precision
	Variants []Variant `json:"variants"`
	Category *Category `json:"category"`
}

type Variant struct {
	Name  string          `json:"name"`
	SKU   string          `json:"sku"`
	Price decimal.Decimal `json:"price"`
}

type Category struct {
	// ID not included since it's for internal use only
	Code string `json:"code"`
	Name string `json:"name"`
}

type CatalogService interface {
	GetAllProducts(params *models.CatalogParams) ([]models.Product, int64, error)
	GetProductByCode(code string) (*models.Product, error)
}

type CatalogHandler struct {
	svc CatalogService
}

func NewHandler(svc CatalogService) *CatalogHandler {
	return &CatalogHandler{
		svc: svc,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	params, err := getQueryParams(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	products, count, err := h.svc.GetAllProducts(params)
	if err != nil {
		// we shouldn't really dump internal error messages in the response
		// like this, since they can expose critical information of the system.
		// We should have some middleware that cleans error messages,
		// but for the simplicity of the task we'll leave it like this
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	response := Response{
		Products:   fn.Map(products, toApiProduct),
		TotalCount: count,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "code must be set")
		return
	}

	product, err := h.svc.GetProductByCode(code)
	switch {
	case errors.Is(err, models.ErrNotFound):
		api.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("could not find product with code %s", code))
		return
	case err != nil:
		api.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get product by code: %s", err))
		return
	}
	api.OKResponse(w, toApiProductDetailed(product))
}

func toApiProduct(p models.Product) Product {
	res := Product{
		Code:  p.Code,
		Price: p.Price.String(),
	}

	if p.Category != nil {
		res.Category = &Category{
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	}

	return res
}

func toApiProductDetailed(p *models.Product) ProductDetailed {
	res := ProductDetailed{
		Code:     p.Code,
		Price:    p.Price.String(),
		Variants: fn.Map(p.Variants, toApiVariant),
	}

	if p.Category != nil {
		res.Category = &Category{
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	}

	return res
}

func toApiVariant(variant models.Variant) Variant {
	return Variant{
		Name:  variant.Name,
		SKU:   variant.SKU,
		Price: variant.Price,
	}
}

func getQueryParams(r *http.Request) (*models.CatalogParams, error) {
	paginationParams, err := getPaginationParams(r)
	if err != nil {
		return nil, err
	}

	priceLt, err := parsePriceLt(r)
	if err != nil {
		return nil, err
	}

	var categoryCode *string
	if code := r.URL.Query().Get("categoryCode"); code != "" {
		categoryCode = &code
	}

	return &models.CatalogParams{
		PriceLt:          priceLt,
		CategoryCode:     categoryCode,
		PaginationParams: paginationParams,
	}, nil
}

func getPaginationParams(r *http.Request) (*models.PaginationParams, error) {
	offset, limit := defaultOffset, defaultLimit

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid offset param '%s': must be a number", offsetStr)
		}
		if parsedOffset < 0 {
			return nil, fmt.Errorf("invalid offset param '%d': must be positive", parsedOffset)
		}
		offset = parsedOffset
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid limit param '%s': must be a number", limitStr)
		}
		if parsedLimit <= 0 {
			return nil, fmt.Errorf("invalid limit param '%d': must be greater than 0", parsedLimit)
		}
		if parsedLimit > maxLimit {
			return nil, fmt.Errorf("invalid limit param '%d': must not exceed %d", parsedLimit, maxLimit)
		}
		limit = parsedLimit
	}

	return &models.PaginationParams{Offset: offset, Limit: limit}, nil
}

func parsePriceLt(r *http.Request) (*decimal.Decimal, error) {
	valueStr := r.URL.Query().Get("priceLt")
	if valueStr == "" {
		return nil, nil
	}

	value, err := decimal.NewFromString(valueStr)
	if err != nil {
		return nil, fmt.Errorf("invalid priceLt param '%s': must be a number", valueStr)
	}

	if value.LessThanOrEqual(decimal.NewFromInt(0)) {
		return nil, fmt.Errorf("invalid priceLt param '%s': must be greater than 0", value.String())
	}

	return &value, nil
}
