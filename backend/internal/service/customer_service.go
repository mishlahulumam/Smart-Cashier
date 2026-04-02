package service

import (
	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/repository"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) GetAll(search string) ([]model.Customer, error) {
	return s.repo.FindAll(search)
}

func (s *CustomerService) GetByID(id uint) (*model.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *CustomerService) Create(req model.CustomerRequest) (*model.Customer, error) {
	customer := &model.Customer{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Update(id uint, req model.CustomerRequest) (*model.Customer, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	customer.Name = req.Name
	customer.Phone = req.Phone
	customer.Email = req.Email
	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Delete(id uint) error {
	return s.repo.Delete(id)
}
