package models

type ProductCategory struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

func (c *ProductCategory) TableName() string {
	return "product_categories"
}
