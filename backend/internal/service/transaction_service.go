package service

import (
	"errors"
	"fmt"
	"time"

	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/repository"
)

type TransactionService struct {
	txRepo       *repository.TransactionRepository
	productRepo  *repository.ProductRepository
	customerRepo *repository.CustomerRepository
}

func NewTransactionService(
	txRepo *repository.TransactionRepository,
	productRepo *repository.ProductRepository,
	customerRepo *repository.CustomerRepository,
) *TransactionService {
	return &TransactionService{
		txRepo:       txRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

func (s *TransactionService) Create(userID uint, req model.CreateTransactionRequest) (*model.Transaction, error) {
	var items []model.TransactionItem
	var subtotal float64

	for _, item := range req.Items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("produk dengan ID %d tidak ditemukan", item.ProductID)
		}
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("stok %s tidak mencukupi (tersisa: %d)", product.Name, product.Stock)
		}

		itemSubtotal := (product.Price * float64(item.Quantity)) - item.Discount
		items = append(items, model.TransactionItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Discount:    item.Discount,
			Subtotal:    itemSubtotal,
		})
		subtotal += itemSubtotal
	}

	total := subtotal - req.Discount
	if total < 0 {
		total = 0
	}

	if req.PaymentAmount < total {
		return nil, errors.New("pembayaran kurang dari total")
	}

	invoice := fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1e6)

	tx := &model.Transaction{
		InvoiceNumber: invoice,
		UserID:        userID,
		CustomerID:    req.CustomerID,
		Items:         items,
		Subtotal:      subtotal,
		Discount:      req.Discount,
		Total:         total,
		PaymentAmount: req.PaymentAmount,
		ChangeAmount:  req.PaymentAmount - total,
	}

	if err := s.txRepo.Create(tx); err != nil {
		return nil, errors.New("gagal menyimpan transaksi")
	}

	for _, item := range req.Items {
		s.productRepo.UpdateStock(item.ProductID, -item.Quantity)
	}

	if req.CustomerID != nil {
		points := int(total / 10000)
		if points > 0 {
			s.customerRepo.AddPoints(*req.CustomerID, points)
		}
	}

	return s.txRepo.FindByID(tx.ID)
}

func (s *TransactionService) GetAll(startDate, endDate time.Time, page, limit int) ([]model.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.txRepo.FindAll(startDate, endDate, page, limit)
}

func (s *TransactionService) GetByID(id uint) (*model.Transaction, error) {
	return s.txRepo.FindByID(id)
}
