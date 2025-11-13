package catalog

import (
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogRepo interface {
	GetAllProducts(*models.CatalogParams) ([]models.Product, int64, error)
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
// Same thing with the repo, it could return a DTO to decouple the repo from the
// model, but for such a simple usecase, using the model is good enough
func (svc *Service) GetAllProducts(params *models.CatalogParams) ([]models.Product, int64, error) {
	res, count, err := svc.repo.GetAllProducts(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all products: %w", err)
	}

	return res, count, nil
}
