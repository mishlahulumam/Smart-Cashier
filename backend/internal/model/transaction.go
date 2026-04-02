package model

import "time"

type Transaction struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	InvoiceNumber string            `json:"invoice_number" gorm:"size:50;uniqueIndex;not null"`
	UserID        uint              `json:"user_id"`
	User          User              `json:"user" gorm:"foreignKey:UserID"`
	CustomerID    *uint             `json:"customer_id"`
	Customer      *Customer         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Items         []TransactionItem `json:"items" gorm:"foreignKey:TransactionID"`
	Subtotal      float64           `json:"subtotal" gorm:"not null"`
	Discount      float64           `json:"discount" gorm:"default:0"`
	Total         float64           `json:"total" gorm:"not null"`
	PaymentAmount float64           `json:"payment_amount" gorm:"not null"`
	ChangeAmount  float64           `json:"change_amount" gorm:"not null"`
	CreatedAt     time.Time         `json:"created_at"`
}

type TransactionItem struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	TransactionID uint    `json:"transaction_id"`
	ProductID     uint    `json:"product_id"`
	ProductName   string  `json:"product_name" gorm:"size:200;not null"`
	Price         float64 `json:"price" gorm:"not null"`
	Quantity      int     `json:"quantity" gorm:"not null"`
	Discount      float64 `json:"discount" gorm:"default:0"`
	Subtotal      float64 `json:"subtotal" gorm:"not null"`
}

type CreateTransactionRequest struct {
	CustomerID    *uint                        `json:"customer_id"`
	Items         []CreateTransactionItemInput `json:"items" binding:"required,min=1"`
	Discount      float64                      `json:"discount" binding:"gte=0"`
	PaymentAmount float64                      `json:"payment_amount" binding:"required,gt=0"`
}

type CreateTransactionItemInput struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gte=1"`
	Discount  float64 `json:"discount" binding:"gte=0"`
}
