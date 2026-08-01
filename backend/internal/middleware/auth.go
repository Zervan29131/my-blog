package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/service"
)

const administratorIDKey = "administrator_id"

func Authenticate(authService *service.AuthService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authorization := context.GetHeader("Authorization")
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			abortUnauthorized(context)
			return
		}

		claims, err := authService.ParseToken(parts[1])
		if err != nil {
			abortUnauthorized(context)
			return
		}

		administratorID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil || administratorID == 0 {
			abortUnauthorized(context)
			return
		}

		context.Set(administratorIDKey, administratorID)
		context.Next()
	}
}

func AdministratorID(context *gin.Context) (uint64, bool) {
	value, exists := context.Get(administratorIDKey)
	if !exists {
		return 0, false
	}
	id, ok := value.(uint64)
	return id, ok
}

func abortUnauthorized(context *gin.Context) {
	context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "未登录或 Token 无效",
			"details": gin.H{},
		},
	})
}
