package repository

import (
	"smart-cashier-backend/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindAll(search string) ([]model.Customer, error) {
	var customers []model.Customer
	q := r.db.Model(&model.Customer{})
	if search != "" {
		q = q.Where("name LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	err := q.Order("name asc").Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) FindByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.First(&customer, id).Error
	return &customer, err
}

func (r *CustomerRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) Update(customer *model.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Customer{}, id).Error
}

func (r *CustomerRepository) AddPoints(id uint, points int) error {
	return r.db.Model(&model.Customer{}).Where("id = ?", id).
		Update("points", gorm.Expr("points + ?", points)).Error
}
