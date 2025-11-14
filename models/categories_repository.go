package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category

	if err := r.db.Order("id").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories from DB: %w", err)
	}

	return categories, nil
}

func (r *CategoriesRepository) GetByCode(code string) (*Category, error) {
	var category Category

	switch err := r.db.Where("code = ?", code).Take(&category).Error; {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, fmt.Errorf("category with code %s not found: %w", code, ErrNotFound)
	case err != nil:
		return nil, fmt.Errorf("failed to get category with code %s from DB: %w", code, err)
	}

	return &category, nil
}

// CreateCategory inserts a category on the DB. The model gets updated with the ID of the inserted category
// Ideally we would generate a UUID in the service instead of relying on the DB,
// but the provided code uses uint for IDs so we'll stick to this aproach
func (r *CategoriesRepository) CreateCategory(category *Category) error {
	err := r.db.Create(&category).Error
	if err != nil {
		return fmt.Errorf("failed to insert category in the DB: %w", err)
	}
	return nil
}
