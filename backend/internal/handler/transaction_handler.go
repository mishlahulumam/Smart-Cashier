package handler

import (
	"net/http"
	"strconv"
	"time"

	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/service"
	"smart-cashier-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data transaksi tidak valid")
		return
	}

	userID := c.GetUint("user_id")
	tx, err := h.service.Create(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Transaksi berhasil", tx)
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var startDate, endDate time.Time
	if sd := c.Query("start_date"); sd != "" {
		startDate, _ = time.Parse("2006-01-02", sd)
	}
	if ed := c.Query("end_date"); ed != "" {
		endDate, _ = time.Parse("2006-01-02", ed)
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	transactions, total, err := h.service.GetAll(startDate, endDate, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data transaksi")
		return
	}

	response.Paginated(c, http.StatusOK, "OK", transactions, total, page, limit)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	tx, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Transaksi tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "OK", tx)
}
