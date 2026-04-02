package model

import "time"

type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Email     string    `json:"email" gorm:"size:100"`
	Points    int       `json:"points" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerRequest struct {
	Name  string `json:"name" binding:"required,max=100"`
	Phone string `json:"phone" binding:"max=20"`
	Email string `json:"email" binding:"max=100"`
}
