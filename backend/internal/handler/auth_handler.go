package handler

import (
	"net/http"

	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/service"
	"smart-cashier-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Username dan password harus diisi")
		return
	}

	result, err := h.service.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", result)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid")
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User berhasil dibuat", user)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.service.GetProfile(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "OK", user)
}
