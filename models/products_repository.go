package models

import (
	"fmt"

	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(params *CatalogParams) ([]Product, int64, error) {
	var products []Product

	query := r.db.Model(&Product{})

	if params.CategoryCode != nil {
		query = query.
			Joins("JOIN product_categories ON products.category_id = product_categories.id").
			Where("product_categories.code = ?", *params.CategoryCode)
	}

	if params.PriceLt != nil {
		query = query.Where("price < ?", *params.PriceLt)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get count from DB: %w", err)
	}

	query = query.
		Order("id"). // to ensure deterministic sorting for pagination
		Preload("Variants").
		Preload("Category").
		Limit(params.Limit).
		Offset(params.Offset)

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get products from DB: %w", err)
	}

	return products, count, nil
}
