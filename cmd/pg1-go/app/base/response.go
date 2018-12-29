package base

import "github.com/gin-gonic/gin"

// StandardJSON generates standard JSON response
func StandardJSON(msg string, data interface{}) map[string]interface{} {
	return gin.H{
		"status":  "OK",
		"message": msg,
		"data":    data,
	}
}

// ErrorJSON generates error JSON response
func ErrorJSON(msg string, data interface{}) map[string]interface{} {
	return gin.H{
		"status":  "error",
		"message": msg,
		"data":    data,
	}
}
