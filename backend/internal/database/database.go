package database

import (
	"log"

	"smart-cashier-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Customer{},
		&model.Transaction{},
		&model.TransactionItem{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	seedAdmin(db)
	return db
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := model.User{
		Username: "admin",
		Password: string(hash),
		Name:     "Administrator",
		Role:     "admin",
	}
	db.Create(&admin)

	hash2, _ := bcrypt.GenerateFromPassword([]byte("kasir123"), bcrypt.DefaultCost)
	cashier := model.User{
		Username: "kasir",
		Password: string(hash2),
		Name:     "Kasir 1",
		Role:     "cashier",
	}
	db.Create(&cashier)

	categories := []model.Category{
		{Name: "Makanan"},
		{Name: "Minuman"},
		{Name: "Snack"},
		{Name: "Lainnya"},
	}
	db.Create(&categories)
}
