package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/model"
	"personal-blog/backend/internal/service"
)

const maximumHomepageConfigBytes = 1 << 20

type HomepageAdminHandler struct {
	service *service.HomepageAdminService
}

func NewHomepageAdminHandler(homepageService *service.HomepageAdminService) *HomepageAdminHandler {
	return &HomepageAdminHandler{service: homepageService}
}

func (handler *HomepageAdminHandler) GetDraft(context *gin.Context) {
	config, err := handler.service.GetDraft(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeAdminConfig(context, config)
}

func (handler *HomepageAdminHandler) SaveDraft(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	config, ok := decodeHomepageConfigRequest(context)
	if !ok {
		return
	}

	saved, err := handler.service.SaveDraft(context.Request.Context(), config, administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeAdminConfig(context, saved)
}

func (handler *HomepageAdminHandler) GetPublished(context *gin.Context) {
	config, err := handler.service.GetPublished(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeAdminConfig(context, config)
}

func (handler *HomepageAdminHandler) Publish(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	result, err := handler.service.Publish(context.Request.Context(), administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.Header("Cache-Control", "no-store")
	writeSuccess(context, http.StatusOK, result)
}

func (handler *HomepageAdminHandler) ResetDraft(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	config, err := handler.service.ResetDraft(context.Request.Context(), administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeAdminConfig(context, config)
}

func (handler *HomepageAdminHandler) Preview(context *gin.Context) {
	preview, err := handler.service.PreviewDraft(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.Header("Cache-Control", "no-store")
	writeSuccess(context, http.StatusOK, preview)
}

func (handler *HomepageAdminHandler) writeError(context *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidHomepageConfig):
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "首页配置不正确")
	case errors.Is(err, service.ErrHomepageVersionNotFound):
		writeError(context, http.StatusNotFound, "HOMEPAGE_CONFIG_NOT_FOUND", "首页配置不存在")
	default:
		slog.Error("administrator homepage config request failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
	}
}

func decodeHomepageConfigRequest(context *gin.Context) (model.HomepageConfig, bool) {
	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, maximumHomepageConfigBytes)
	data, err := io.ReadAll(context.Request.Body)
	if err != nil {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "首页配置不正确")
		return model.HomepageConfig{}, false
	}
	config, err := model.DecodeHomepageConfig(data)
	if err != nil {
		writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "首页配置不正确")
		return model.HomepageConfig{}, false
	}
	return config, true
}

func administratorID(context *gin.Context) (uint64, bool) {
	administratorID, exists := middleware.AdministratorID(context)
	if !exists {
		writeError(context, http.StatusUnauthorized, "UNAUTHORIZED", "未登录或 Token 无效")
		return 0, false
	}
	return administratorID, true
}

func writeAdminConfig(context *gin.Context, config service.AdminHomepageConfig) {
	context.Header("Cache-Control", "no-store")
	writeSuccess(context, http.StatusOK, config)
}
