package middleware

import (
	"net/http"
	"strings"

	jwtpkg "smart-cashier-backend/pkg/jwt"
	"smart-cashier-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, http.StatusUnauthorized, "Token tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Format token tidak valid")
			c.Abort()
			return
		}

		claims, err := jwtpkg.ValidateToken(parts[1], secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			response.Error(c, http.StatusForbidden, "Hanya admin yang dapat mengakses")
			c.Abort()
			return
		}
		c.Next()
	}
}
