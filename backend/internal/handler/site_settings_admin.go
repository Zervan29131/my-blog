package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/service"
)

type SiteSettingsAdminHandler struct {
	service *service.SiteSettingsAdminService
}

type siteSettingRequest struct {
	SiteName             string  `json:"site_name"`
	SiteShortName        *string `json:"site_short_name"`
	SiteDescription      string  `json:"site_description"`
	TitleSuffix          *string `json:"title_suffix"`
	LogoURL              *string `json:"logo_url"`
	FaviconURL           *string `json:"favicon_url"`
	DefaultShareImageURL *string `json:"default_share_image_url"`
	CopyrightName        string  `json:"copyright_name"`
	StartYear            *int    `json:"start_year"`
	AdditionalText       *string `json:"additional_text"`
	FilingNumber         *string `json:"filing_number"`
	FilingURL            *string `json:"filing_url"`
	ShowTechnology       *bool   `json:"show_technology"`
	TechnologyText       *string `json:"technology_text"`
}

func NewSiteSettingsAdminHandler(settingsService *service.SiteSettingsAdminService) *SiteSettingsAdminHandler {
	return &SiteSettingsAdminHandler{service: settingsService}
}

func (handler *SiteSettingsAdminHandler) Get(context *gin.Context) {
	setting, err := handler.service.Get(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.Header("Cache-Control", "no-store")
	writeSuccess(context, http.StatusOK, setting)
}

func (handler *SiteSettingsAdminHandler) Update(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	var request siteSettingRequest
	if !decodeSiteContentRequest(context, &request) || request.ShowTechnology == nil {
		writeSiteContentValidationError(context)
		return
	}
	setting, err := handler.service.Update(context.Request.Context(), request.input(), administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.Header("Cache-Control", "no-store")
	writeSuccess(context, http.StatusOK, setting)
}

func (handler *SiteSettingsAdminHandler) writeError(context *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidSiteSetting):
		writeSiteContentValidationError(context)
	case errors.Is(err, service.ErrSiteSettingNotFound):
		writeError(context, http.StatusNotFound, "SITE_SETTING_NOT_FOUND", "站点设置不存在")
	default:
		slog.Error("administrator site setting request failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
	}
}

func (request siteSettingRequest) input() service.SiteSettingInput {
	return service.SiteSettingInput{
		SiteName: request.SiteName, SiteShortName: request.SiteShortName,
		SiteDescription: request.SiteDescription, TitleSuffix: request.TitleSuffix,
		LogoURL: request.LogoURL, FaviconURL: request.FaviconURL,
		DefaultShareImageURL: request.DefaultShareImageURL,
		CopyrightName:        request.CopyrightName, StartYear: request.StartYear,
		AdditionalText: request.AdditionalText, FilingNumber: request.FilingNumber,
		FilingURL: request.FilingURL, ShowTechnology: *request.ShowTechnology,
		TechnologyText: request.TechnologyText,
	}
}
