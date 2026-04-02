package handler

import (
	"net/http"
	"strconv"

	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/service"
	"smart-cashier-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data kategori")
		return
	}
	response.Success(c, http.StatusOK, "OK", categories)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	category, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Kategori tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "OK", category)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req model.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Nama kategori harus diisi")
		return
	}
	category, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal membuat kategori")
		return
	}
	response.Success(c, http.StatusCreated, "Kategori berhasil dibuat", category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req model.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Nama kategori harus diisi")
		return
	}
	category, err := h.service.Update(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal mengupdate kategori")
		return
	}
	response.Success(c, http.StatusOK, "Kategori berhasil diupdate", category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal menghapus kategori")
		return
	}
	response.Success(c, http.StatusOK, "Kategori berhasil dihapus", nil)
}
