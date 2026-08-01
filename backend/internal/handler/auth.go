package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

type loginRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=72"`
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Login(context *gin.Context) {
	var request loginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}
	request.Username = strings.TrimSpace(request.Username)
	if request.Username == "" {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
		return
	}

	result, err := handler.authService.Login(context.Request.Context(), request.Username, request.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		slog.Warn("administrator login failed", "username", request.Username)
		writeError(context, http.StatusUnauthorized, "INVALID_CREDENTIALS", "用户名或密码错误")
		return
	}
	if err != nil {
		slog.Error("administrator login failed due to internal error", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}

	writeSuccess(context, http.StatusOK, gin.H{
		"token":      result.Token,
		"expires_in": result.ExpiresIn,
	})
}

func (handler *AuthHandler) Me(context *gin.Context) {
	administratorID, exists := middleware.AdministratorID(context)
	if !exists {
		writeError(context, http.StatusUnauthorized, "UNAUTHORIZED", "未登录或 Token 无效")
		return
	}

	administrator, err := handler.authService.CurrentAdministrator(context.Request.Context(), administratorID)
	if errors.Is(err, service.ErrAdministratorGone) {
		writeError(context, http.StatusUnauthorized, "UNAUTHORIZED", "未登录或 Token 无效")
		return
	}
	if err != nil {
		slog.Error("get current administrator failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}

	writeSuccess(context, http.StatusOK, gin.H{
		"id":         administrator.ID,
		"username":   administrator.Username,
		"created_at": administrator.CreatedAt,
		"updated_at": administrator.UpdatedAt,
	})
}
