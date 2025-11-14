package categories

import (
	"errors"
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/models"
)

var ErrDuplicateCode = errors.New("category code must be unique")

type CategoriesRepo interface {
	GetAllCategories() ([]models.Category, error)
	GetByCode(code string) (*models.Category, error)
	CreateCategory(category *models.Category) error
}

type Service struct {
	repo CategoriesRepo
}

func NewService(repo CategoriesRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (svc *Service) GetAllCategories() ([]models.Category, error) {
	res, err := svc.repo.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}

	return res, nil
}

type CreateCategoryReq struct {
	Code string
	Name string
}

func (svc *Service) CreateCategory(req CreateCategoryReq) (*models.Category, error) {
	switch _, err := svc.repo.GetByCode(req.Code); {
	case errors.Is(err, models.ErrNotFound):
		break // this is what we want, so continue
	case err != nil:
		return nil, fmt.Errorf("failed to get category by code: %w", err)
	default:
		return nil, ErrDuplicateCode
	}

	category := &models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := svc.repo.CreateCategory(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}
