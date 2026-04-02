package service

import (
	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) GetByID(id uint) (*model.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Create(req model.CategoryRequest) (*model.Category, error) {
	category := &model.Category{Name: req.Name}
	if err := s.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(id uint, req model.CategoryRequest) (*model.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	category.Name = req.Name
	if err := s.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}
