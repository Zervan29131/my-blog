package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(context *gin.Context) {
		origin := context.GetHeader("Origin")
		if origin == "" {
			context.Next()
			return
		}

		context.Header("Vary", "Origin")
		if _, exists := allowed[origin]; !exists {
			if context.Request.Method == http.MethodOptions {
				context.AbortWithStatus(http.StatusForbidden)
				return
			}
			context.Next()
			return
		}

		context.Header("Access-Control-Allow-Origin", origin)
		context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		context.Header("Access-Control-Max-Age", "600")
		if context.Request.Method == http.MethodOptions {
			requestedMethod := context.GetHeader("Access-Control-Request-Method")
			if requestedMethod == "" || !allowedCORSMethod(requestedMethod) {
				context.AbortWithStatus(http.StatusForbidden)
				return
			}
			context.AbortWithStatus(http.StatusNoContent)
			return
		}
		context.Next()
	}
}

func allowedCORSMethod(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}
