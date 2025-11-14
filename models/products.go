package models

import (
	"github.com/shopspring/decimal"

	"github.com/mytheresa/go-hiring-challenge/lib/fn"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Variants   []Variant       `gorm:"foreignKey:ProductID"`
	CategoryID uint
	Category   *Category
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) NormalizeVariantPrices() {
	p.Variants = fn.Map(p.Variants, func(v Variant) Variant {
		if v.Price.IsZero() {
			v.Price = p.Price
		}
		return v
	})
}

type CatalogParams struct {
	CategoryCode *string
	PriceLt      *decimal.Decimal
	*PaginationParams
}

type PaginationParams struct {
	Offset int
	Limit  int
}
