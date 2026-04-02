package handler

import (
	"net/http"
	"strconv"
	"time"

	"smart-cashier-backend/internal/service"
	"smart-cashier-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

func (h *ReportHandler) Dashboard(c *gin.Context) {
	data, err := h.service.GetDashboard()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data dashboard")
		return
	}
	response.Success(c, http.StatusOK, "OK", data)
}

func (h *ReportHandler) DailySummary(c *gin.Context) {
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Format tanggal tidak valid (YYYY-MM-DD)")
		return
	}

	data, err := h.service.GetDailySummary(date)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil laporan harian")
		return
	}
	response.Success(c, http.StatusOK, "OK", data)
}

func (h *ReportHandler) MonthlySummary(c *gin.Context) {
	now := time.Now()
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))

	data, err := h.service.GetMonthlySummary(year, month)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil laporan bulanan")
		return
	}
	response.Success(c, http.StatusOK, "OK", data)
}

func (h *ReportHandler) TopProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var startDate, endDate time.Time
	if sd := c.Query("start_date"); sd != "" {
		startDate, _ = time.Parse("2006-01-02", sd)
	}
	if ed := c.Query("end_date"); ed != "" {
		endDate, _ = time.Parse("2006-01-02", ed)
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	data, err := h.service.GetTopProducts(startDate, endDate, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data produk terlaris")
		return
	}
	response.Success(c, http.StatusOK, "OK", data)
}

func (h *ReportHandler) Summary(c *gin.Context) {
	var startDate, endDate time.Time
	if sd := c.Query("start_date"); sd != "" {
		startDate, _ = time.Parse("2006-01-02", sd)
	}
	if ed := c.Query("end_date"); ed != "" {
		endDate, _ = time.Parse("2006-01-02", ed)
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	data, err := h.service.GetSummary(startDate, endDate)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil ringkasan laporan")
		return
	}
	response.Success(c, http.StatusOK, "OK", data)
}
