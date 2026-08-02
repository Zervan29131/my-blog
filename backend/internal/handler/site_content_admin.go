package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/service"
)

const maximumSiteContentRequestBytes = 64 << 10

type SiteContentAdminHandler struct {
	navigation *service.NavigationAdminService
	social     *service.SocialLinkAdminService
	featured   *service.FeaturedArticleAdminService
}

type navigationRequest struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	LinkType     string `json:"link_type"`
	OpenInNewTab *bool  `json:"open_in_new_tab"`
	IsVisible    *bool  `json:"is_visible"`
	SortOrder    *int   `json:"sort_order"`
}

type socialLinkRequest struct {
	Platform    string `json:"platform"`
	DisplayName string `json:"display_name"`
	URL         string `json:"url"`
	IsVisible   *bool  `json:"is_visible"`
	SortOrder   *int   `json:"sort_order"`
}

type featuredArticleRequest struct {
	ArticleID uint64 `json:"article_id"`
	SortOrder *int   `json:"sort_order"`
	IsVisible *bool  `json:"is_visible"`
}

type featuredArticleVisibilityRequest struct {
	IsVisible *bool `json:"is_visible"`
}

type resourceOrderRequest struct {
	Items []resourceOrderItem `json:"items"`
}

type resourceOrderItem struct {
	ID        uint64 `json:"id"`
	SortOrder *int   `json:"sort_order"`
}

type featuredOrderRequest struct {
	Items []featuredOrderItem `json:"items"`
}

type featuredOrderItem struct {
	ArticleID uint64 `json:"article_id"`
	SortOrder *int   `json:"sort_order"`
}

func NewSiteContentAdminHandler(
	navigation *service.NavigationAdminService,
	social *service.SocialLinkAdminService,
	featured *service.FeaturedArticleAdminService,
) *SiteContentAdminHandler {
	return &SiteContentAdminHandler{navigation: navigation, social: social, featured: featured}
}

func (handler *SiteContentAdminHandler) ListNavigation(context *gin.Context) {
	items, err := handler.navigation.List(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, items)
}

func (handler *SiteContentAdminHandler) CreateNavigation(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	var request navigationRequest
	if !decodeSiteContentRequest(context, &request) || !request.valid() {
		writeSiteContentValidationError(context)
		return
	}
	item, err := handler.navigation.Create(context.Request.Context(), request.input(), administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusCreated, item)
}

func (handler *SiteContentAdminHandler) UpdateNavigation(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	id, ok := parseSiteContentID(context, "id")
	if !ok {
		return
	}
	var request navigationRequest
	if !decodeSiteContentRequest(context, &request) || !request.valid() {
		writeSiteContentValidationError(context)
		return
	}
	item, err := handler.navigation.Update(context.Request.Context(), id, request.input(), administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, item)
}

func (handler *SiteContentAdminHandler) DeleteNavigation(context *gin.Context) {
	handler.deleteResource(context, "id", handler.navigation.Delete)
}

func (handler *SiteContentAdminHandler) ReorderNavigation(context *gin.Context) {
	handler.reorderResources(context, handler.navigation.Reorder)
}

func (handler *SiteContentAdminHandler) ListSocialLinks(context *gin.Context) {
	links, err := handler.social.List(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, links)
}

func (handler *SiteContentAdminHandler) CreateSocialLink(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	var request socialLinkRequest
	if !decodeSiteContentRequest(context, &request) || !request.valid() {
		writeSiteContentValidationError(context)
		return
	}
	link, err := handler.social.Create(context.Request.Context(), request.input(), administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusCreated, link)
}

func (handler *SiteContentAdminHandler) UpdateSocialLink(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	id, ok := parseSiteContentID(context, "id")
	if !ok {
		return
	}
	var request socialLinkRequest
	if !decodeSiteContentRequest(context, &request) || !request.valid() {
		writeSiteContentValidationError(context)
		return
	}
	link, err := handler.social.Update(context.Request.Context(), id, request.input(), administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, link)
}

func (handler *SiteContentAdminHandler) DeleteSocialLink(context *gin.Context) {
	handler.deleteResource(context, "id", handler.social.Delete)
}

func (handler *SiteContentAdminHandler) ReorderSocialLinks(context *gin.Context) {
	handler.reorderResources(context, handler.social.Reorder)
}

func (handler *SiteContentAdminHandler) ListFeaturedArticles(context *gin.Context) {
	items, err := handler.featured.List(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, items)
}

func (handler *SiteContentAdminHandler) CreateFeaturedArticle(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	var request featuredArticleRequest
	if !decodeSiteContentRequest(context, &request) || request.ArticleID == 0 ||
		request.SortOrder == nil || request.IsVisible == nil {
		writeSiteContentValidationError(context)
		return
	}
	item, err := handler.featured.Create(context.Request.Context(), service.FeaturedArticleInput{
		ArticleID: request.ArticleID, SortOrder: *request.SortOrder, IsVisible: *request.IsVisible,
	}, administratorID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusCreated, item)
}

func (handler *SiteContentAdminHandler) DeleteFeaturedArticle(context *gin.Context) {
	handler.deleteResource(context, "articleId", handler.featured.Delete)
}

func (handler *SiteContentAdminHandler) UpdateFeaturedArticle(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	articleID, ok := parseSiteContentID(context, "articleId")
	if !ok {
		return
	}
	var request featuredArticleVisibilityRequest
	if !decodeSiteContentRequest(context, &request) || request.IsVisible == nil {
		writeSiteContentValidationError(context)
		return
	}
	item, err := handler.featured.UpdateVisibility(
		context.Request.Context(), articleID, *request.IsVisible, administratorID,
	)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, item)
}

func (handler *SiteContentAdminHandler) ReorderFeaturedArticles(context *gin.Context) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	var request featuredOrderRequest
	if !decodeSiteContentRequest(context, &request) {
		writeSiteContentValidationError(context)
		return
	}
	items := make([]service.ResourceOrder, 0, len(request.Items))
	for _, item := range request.Items {
		if item.ArticleID == 0 || item.SortOrder == nil {
			writeSiteContentValidationError(context)
			return
		}
		items = append(items, service.ResourceOrder{ID: item.ArticleID, SortOrder: *item.SortOrder})
	}
	if err := handler.featured.Reorder(context.Request.Context(), items, administratorID); err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, gin.H{"updated": len(items)})
}

func (handler *SiteContentAdminHandler) deleteResource(
	context *gin.Context,
	parameter string,
	deleteFunction func(context.Context, uint64, uint64) error,
) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	id, ok := parseSiteContentID(context, parameter)
	if !ok {
		return
	}
	if err := deleteFunction(context.Request.Context(), id, administratorID); err != nil {
		handler.writeError(context, err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (handler *SiteContentAdminHandler) reorderResources(
	context *gin.Context,
	reorderFunction func(context.Context, []service.ResourceOrder, uint64) error,
) {
	administratorID, ok := administratorID(context)
	if !ok {
		return
	}
	var request resourceOrderRequest
	if !decodeSiteContentRequest(context, &request) {
		writeSiteContentValidationError(context)
		return
	}
	items := make([]service.ResourceOrder, 0, len(request.Items))
	for _, item := range request.Items {
		if item.ID == 0 || item.SortOrder == nil {
			writeSiteContentValidationError(context)
			return
		}
		items = append(items, service.ResourceOrder{ID: item.ID, SortOrder: *item.SortOrder})
	}
	if err := reorderFunction(context.Request.Context(), items, administratorID); err != nil {
		handler.writeError(context, err)
		return
	}
	writeSuccess(context, http.StatusOK, gin.H{"updated": len(items)})
}

func (handler *SiteContentAdminHandler) writeError(context *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidSiteContent):
		writeSiteContentValidationError(context)
	case errors.Is(err, service.ErrSiteContentNotFound):
		writeError(context, http.StatusNotFound, "RESOURCE_NOT_FOUND", "资源不存在")
	case errors.Is(err, service.ErrSiteContentLimit):
		writeError(context, http.StatusConflict, "RESOURCE_LIMIT_EXCEEDED", "资源数量已达到上限")
	case errors.Is(err, service.ErrSiteContentConflict):
		writeError(context, http.StatusConflict, "RESOURCE_CONFLICT", "资源已经存在")
	case errors.Is(err, service.ErrArticleNotPublished):
		writeError(context, http.StatusBadRequest, "ARTICLE_NOT_PUBLISHED", "只能推荐已发布文章")
	default:
		slog.Error("administrator site content request failed", "error", err)
		writeError(context, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
	}
}

func (request navigationRequest) valid() bool {
	return request.OpenInNewTab != nil && request.IsVisible != nil && request.SortOrder != nil
}

func (request navigationRequest) input() service.NavigationInput {
	return service.NavigationInput{
		Name: request.Name, URL: request.URL, LinkType: request.LinkType,
		OpenInNewTab: *request.OpenInNewTab, IsVisible: *request.IsVisible, SortOrder: *request.SortOrder,
	}
}

func (request socialLinkRequest) valid() bool {
	return request.IsVisible != nil && request.SortOrder != nil
}

func (request socialLinkRequest) input() service.SocialLinkInput {
	return service.SocialLinkInput{
		Platform: request.Platform, DisplayName: request.DisplayName, URL: request.URL,
		IsVisible: *request.IsVisible, SortOrder: *request.SortOrder,
	}
}

func decodeSiteContentRequest(context *gin.Context, target any) bool {
	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, maximumSiteContentRequestBytes)
	data, err := io.ReadAll(context.Request.Body)
	if err != nil {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return false
	}
	var trailing any
	return decoder.Decode(&trailing) == io.EOF
}

func parseSiteContentID(context *gin.Context, parameter string) (uint64, bool) {
	id, err := strconv.ParseUint(context.Param(parameter), 10, 64)
	if err != nil || id == 0 {
		writeSiteContentValidationError(context)
		return 0, false
	}
	return id, true
}

func writeSiteContentValidationError(context *gin.Context) {
	writeError(context, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确")
}
