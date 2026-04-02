package handler

import (
	"net/http"
	"strconv"

	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/service"
	"smart-cashier-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	search := c.Query("search")
	lowStock := c.Query("low_stock") == "true"

	products, err := h.service.GetAll(uint(categoryID), search, lowStock)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data produk")
		return
	}
	response.Success(c, http.StatusOK, "OK", products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	product, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "OK", product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req model.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data produk tidak valid")
		return
	}
	product, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal membuat produk")
		return
	}
	response.Success(c, http.StatusCreated, "Produk berhasil dibuat", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req model.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data produk tidak valid")
		return
	}
	product, err := h.service.Update(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal mengupdate produk")
		return
	}
	response.Success(c, http.StatusOK, "Produk berhasil diupdate", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal menghapus produk")
		return
	}
	response.Success(c, http.StatusOK, "Produk berhasil dihapus", nil)
}
