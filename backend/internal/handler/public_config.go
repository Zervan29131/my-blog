package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/service"
)

type PublicConfigHandler struct {
	service *service.PublicConfigService
}

func NewPublicConfigHandler(configService *service.PublicConfigService) *PublicConfigHandler {
	return &PublicConfigHandler{service: configService}
}

func (handler *PublicConfigHandler) SiteConfig(context *gin.Context) {
	config, err := handler.service.GetSiteConfig(context.Request.Context())
	if err != nil {
		slog.Error("get public site config failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}
	writeSuccess(context, http.StatusOK, config)
}

func (handler *PublicConfigHandler) Homepage(context *gin.Context) {
	config, err := handler.service.GetHomepage(context.Request.Context())
	if err != nil {
		slog.Error("get public homepage config failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}
	writeSuccess(context, http.StatusOK, config)
}
