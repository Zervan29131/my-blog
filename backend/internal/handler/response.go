package handler

import "github.com/gin-gonic/gin"

func writeSuccess(context *gin.Context, status int, data any) {
	context.JSON(status, gin.H{
		"data":    data,
		"message": "success",
	})
}

func writeError(context *gin.Context, status int, code, message string) {
	context.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": gin.H{},
		},
	})
}
