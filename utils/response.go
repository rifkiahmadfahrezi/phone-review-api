package utils

import "github.com/gin-gonic/gin"

func ResponseJSON(message string, status_code int, data any) map[string]any {
	if data == nil {
		return gin.H{
			"message":     message,
			"status_code": status_code,
		}
	}
	if message == "" {
		return gin.H{
			"data":        data,
			"status_code": status_code,
		}
	}
	return gin.H{
		"data":        data,
		"message":     message,
		"status_code": status_code,
	}
}
