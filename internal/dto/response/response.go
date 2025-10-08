package response

import "github.com/gin-gonic/gin"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONResponse(c *gin.Context, statusCode int, status, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	JSONResponse(c, statusCode, "success", message, data)
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	JSONResponse(c, statusCode, "error", message, nil)
}
