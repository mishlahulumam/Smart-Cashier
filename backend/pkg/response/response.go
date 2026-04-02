package response

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{Success: true, Message: message, Data: data})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{Success: false, Message: message})
}

func Paginated(c *gin.Context, status int, message string, data interface{}, total int64, page, limit int) {
	c.JSON(status, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}
