package service

import (
	"time"

	"smart-cashier-backend/internal/repository"
)

type ReportService struct {
	txRepo      *repository.TransactionRepository
	productRepo *repository.ProductRepository
}

func NewReportService(txRepo *repository.TransactionRepository, productRepo *repository.ProductRepository) *ReportService {
	return &ReportService{txRepo: txRepo, productRepo: productRepo}
}

func (s *ReportService) GetDailySummary(date time.Time) (map[string]interface{}, error) {
	revenue, transactions, items, err := s.txRepo.GetDailySummary(date)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"date":         date.Format("2006-01-02"),
		"revenue":      revenue,
		"transactions": transactions,
		"items":        items,
	}, nil
}

func (s *ReportService) GetMonthlySummary(year, month int) ([]map[string]interface{}, error) {
	return s.txRepo.GetMonthlySummary(year, month)
}

func (s *ReportService) GetTopProducts(startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.txRepo.GetTopProducts(startDate, endDate, limit)
}

func (s *ReportService) GetSummary(startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.txRepo.GetSummary(startDate, endDate)
}

func (s *ReportService) GetDashboard() (map[string]interface{}, error) {
	todayRevenue, _ := s.txRepo.TodayRevenue()
	todayTx, _ := s.txRepo.CountToday()
	totalProducts, _ := s.productRepo.Count()
	lowStock, _ := s.productRepo.CountLowStock()
	recentTx, _ := s.txRepo.GetRecentTransactions(5)

	now := time.Now()
	topProducts, _ := s.txRepo.GetTopProducts(
		time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
		now, 5,
	)

	return map[string]interface{}{
		"today_revenue":       todayRevenue,
		"today_transactions":  todayTx,
		"total_products":      totalProducts,
		"low_stock_count":     lowStock,
		"recent_transactions": recentTx,
		"top_products":        topProducts,
	}, nil
}
