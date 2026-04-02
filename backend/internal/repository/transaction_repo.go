package repository

import (
	"time"

	"smart-cashier-backend/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TransactionRepository) FindAll(startDate, endDate time.Time, page, limit int) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	q := r.db.Model(&model.Transaction{})
	if !startDate.IsZero() {
		q = q.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		q = q.Where("created_at <= ?", endDate)
	}

	q.Count(&total)

	offset := (page - 1) * limit
	err := q.Preload("User").Preload("Customer").Preload("Items").
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *TransactionRepository) FindByID(id uint) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Preload("User").Preload("Customer").Preload("Items").
		First(&tx, id).Error
	return &tx, err
}

func (r *TransactionRepository) GetDailySummary(date time.Time) (float64, int64, int64, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var result struct {
		Revenue      float64
		Transactions int64
		Items        int64
	}

	err := r.db.Model(&model.Transaction{}).
		Select("COALESCE(SUM(total), 0) as revenue, COUNT(*) as transactions").
		Where("created_at >= ? AND created_at < ?", start, end).
		Scan(&result).Error
	if err != nil {
		return 0, 0, 0, err
	}

	r.db.Model(&model.TransactionItem{}).
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.created_at >= ? AND transactions.created_at < ?", start, end).
		Select("COALESCE(SUM(transaction_items.quantity), 0) as items").
		Scan(&result)

	return result.Revenue, result.Transactions, result.Items, nil
}

func (r *TransactionRepository) GetMonthlySummary(year, month int) ([]map[string]interface{}, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)

	var results []map[string]interface{}
	err := r.db.Model(&model.Transaction{}).
		Select("DATE(created_at) as date, COALESCE(SUM(total), 0) as revenue, COUNT(*) as transactions").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at)").
		Order("date asc").
		Find(&results).Error

	return results, err
}

func (r *TransactionRepository) GetTopProducts(startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	q := r.db.Model(&model.TransactionItem{}).
		Select("transaction_items.product_id, transaction_items.product_name, SUM(transaction_items.quantity) as quantity, SUM(transaction_items.subtotal) as revenue").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id")

	if !startDate.IsZero() {
		q = q.Where("transactions.created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		q = q.Where("transactions.created_at <= ?", endDate)
	}

	err := q.Group("transaction_items.product_id, transaction_items.product_name").
		Order("quantity desc").
		Limit(limit).
		Find(&results).Error

	return results, err
}

func (r *TransactionRepository) GetSummary(startDate, endDate time.Time) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"total_revenue":      0.0,
		"total_transactions": int64(0),
		"total_items":        int64(0),
		"avg_transaction":    0.0,
	}

	var summary struct {
		Revenue      float64
		Transactions int64
	}

	q := r.db.Model(&model.Transaction{}).
		Select("COALESCE(SUM(total), 0) as revenue, COUNT(*) as transactions")

	if !startDate.IsZero() {
		q = q.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		q = q.Where("created_at <= ?", endDate)
	}

	if err := q.Scan(&summary).Error; err != nil {
		return result, err
	}

	var itemCount int64
	iq := r.db.Model(&model.TransactionItem{}).
		Select("COALESCE(SUM(quantity), 0)").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id")
	if !startDate.IsZero() {
		iq = iq.Where("transactions.created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		iq = iq.Where("transactions.created_at <= ?", endDate)
	}
	iq.Scan(&itemCount)

	avg := 0.0
	if summary.Transactions > 0 {
		avg = summary.Revenue / float64(summary.Transactions)
	}

	result["total_revenue"] = summary.Revenue
	result["total_transactions"] = summary.Transactions
	result["total_items"] = itemCount
	result["avg_transaction"] = avg

	return result, nil
}

func (r *TransactionRepository) GetRecentTransactions(limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("User").Preload("Customer").Preload("Items").
		Order("created_at desc").Limit(limit).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) CountToday() (int64, error) {
	var count int64
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	err := r.db.Model(&model.Transaction{}).
		Where("created_at >= ?", start).Count(&count).Error
	return count, err
}

func (r *TransactionRepository) TodayRevenue() (float64, error) {
	var revenue float64
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	err := r.db.Model(&model.Transaction{}).
		Select("COALESCE(SUM(total), 0)").
		Where("created_at >= ?", start).
		Scan(&revenue).Error
	return revenue, err
}
