package model

import "time"

type Product struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:200;not null"`
	SKU        string    `json:"sku" gorm:"size:50;uniqueIndex"`
	Price      float64   `json:"price" gorm:"not null"`
	Stock      int       `json:"stock" gorm:"not null;default:0"`
	MinStock   int       `json:"min_stock" gorm:"not null;default:5"`
	Image      string    `json:"image" gorm:"size:500"`
	CategoryID uint      `json:"category_id"`
	Category   Category  `json:"category" gorm:"foreignKey:CategoryID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductRequest struct {
	Name       string  `json:"name" binding:"required,max=200"`
	SKU        string  `json:"sku" binding:"max=50"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Stock      int     `json:"stock" binding:"gte=0"`
	MinStock   int     `json:"min_stock" binding:"gte=0"`
	Image      string  `json:"image"`
	CategoryID uint    `json:"category_id" binding:"required"`
}
