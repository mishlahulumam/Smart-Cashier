package service

import (
	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(categoryID uint, search string, lowStock bool) ([]model.Product, error) {
	return s.repo.FindAll(categoryID, search, lowStock)
}

func (s *ProductService) GetByID(id uint) (*model.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Create(req model.ProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		Stock:      req.Stock,
		MinStock:   req.MinStock,
		Image:      req.Image,
		CategoryID: req.CategoryID,
	}
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return s.repo.FindByID(product.ID)
}

func (s *ProductService) Update(id uint, req model.ProductRequest) (*model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	product.Name = req.Name
	product.SKU = req.SKU
	product.Price = req.Price
	product.Stock = req.Stock
	product.MinStock = req.MinStock
	product.Image = req.Image
	product.CategoryID = req.CategoryID

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return s.repo.FindByID(product.ID)
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}
