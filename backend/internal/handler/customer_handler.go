package handler

import (
	"net/http"
	"strconv"

	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/service"
	"smart-cashier-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	customers, err := h.service.GetAll(search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data pelanggan")
		return
	}
	response.Success(c, http.StatusOK, "OK", customers)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	customer, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Pelanggan tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "OK", customer)
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req model.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data pelanggan tidak valid")
		return
	}
	customer, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal membuat data pelanggan")
		return
	}
	response.Success(c, http.StatusCreated, "Pelanggan berhasil ditambahkan", customer)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req model.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data pelanggan tidak valid")
		return
	}
	customer, err := h.service.Update(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal mengupdate data pelanggan")
		return
	}
	response.Success(c, http.StatusOK, "Data pelanggan berhasil diupdate", customer)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal menghapus data pelanggan")
		return
	}
	response.Success(c, http.StatusOK, "Data pelanggan berhasil dihapus", nil)
}
