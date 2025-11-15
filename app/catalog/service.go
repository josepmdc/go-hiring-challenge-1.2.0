package catalog

import (
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogRepo interface {
	GetAllProducts(params *models.CatalogParams) ([]models.Product, int64, error)
	GetProductByCode(code string) (*models.Product, error)
}

type Service struct {
	repo CatalogRepo
}

func NewService(repo CatalogRepo) *Service {
	return &Service{
		repo: repo,
	}
}

// The service returns the model. Optionally we could return a DTO to decouple it
// even more, but for this simple usecase I think returning the model is good enough
// Same thing with the repo, it could return a DTO to decouple  but for such a simple
// usecase, using the model is good enough
func (svc *Service) GetAllProducts(params *models.CatalogParams) ([]models.Product, int64, error) {
	res, count, err := svc.repo.GetAllProducts(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all products: %w", err)
	}

	for _, p := range res {
		p.NormalizeVariantPrices()
	}

	return res, count, nil
}

func (svc *Service) GetProductByCode(code string) (*models.Product, error) {
	product, err := svc.repo.GetProductByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by code from repo: %w", err)
	}

	product.NormalizeVariantPrices()

	return product, nil
}
