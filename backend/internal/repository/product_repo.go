package repository

import (
	"smart-cashier-backend/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll(categoryID uint, search string, lowStock bool) ([]model.Product, error) {
	var products []model.Product
	q := r.db.Preload("Category")

	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	if search != "" {
		q = q.Where("name LIKE ? OR sku LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if lowStock {
		q = q.Where("stock <= min_stock")
	}

	err := q.Order("name asc").Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *ProductRepository) CountLowStock() (int64, error) {
	var count int64
	err := r.db.Model(&model.Product{}).Where("stock <= min_stock").Count(&count).Error
	return count, err
}

func (r *ProductRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Product{}).Count(&count).Error
	return count, err
}
