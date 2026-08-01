package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"personal-blog/backend/internal/handler"
)

type successResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type Dependencies struct {
	AuthHandler       *handler.AuthHandler
	ArticleHandler    *handler.ArticleHandler
	CommentHandler    *handler.CommentHandler
	DashboardHandler  *handler.DashboardHandler
	PublicConfig      *handler.PublicConfigHandler
	HomepageAdmin     *handler.HomepageAdminHandler
	SiteContentAdmin  *handler.SiteContentAdminHandler
	SiteSettingsAdmin *handler.SiteSettingsAdminHandler
	AuthMiddleware    gin.HandlerFunc
	CommentRateLimit  gin.HandlerFunc
	CORSMiddleware    gin.HandlerFunc
}

func New(dependencies Dependencies) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	if dependencies.CORSMiddleware != nil {
		engine.Use(dependencies.CORSMiddleware)
	}

	api := engine.Group("/api/v1")
	api.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, successResponse{
			Data:    gin.H{"status": "ok"},
			Message: "success",
		})
	})
	if dependencies.PublicConfig != nil {
		api.GET("/site/config", dependencies.PublicConfig.SiteConfig)
		api.GET("/homepage", dependencies.PublicConfig.Homepage)
	}
	if dependencies.ArticleHandler != nil {
		api.GET("/articles", dependencies.ArticleHandler.ListPublished)
		api.GET("/articles/:slug", dependencies.ArticleHandler.GetPublished)
	}
	if dependencies.CommentHandler != nil {
		api.GET("/articles/:slug/comments", dependencies.CommentHandler.ListApproved)
		if dependencies.CommentRateLimit != nil {
			api.POST(
				"/articles/:slug/comments",
				dependencies.CommentRateLimit,
				dependencies.CommentHandler.Submit,
			)
		} else {
			api.POST("/articles/:slug/comments", dependencies.CommentHandler.Submit)
		}
	}

	if dependencies.AuthHandler != nil && dependencies.AuthMiddleware != nil {
		admin := api.Group("/admin")
		admin.POST("/auth/login", dependencies.AuthHandler.Login)

		protectedAdmin := admin.Group("")
		protectedAdmin.Use(dependencies.AuthMiddleware)
		protectedAdmin.GET("/auth/me", dependencies.AuthHandler.Me)
		if dependencies.DashboardHandler != nil {
			protectedAdmin.GET("/dashboard", dependencies.DashboardHandler.Stats)
		}
		if dependencies.HomepageAdmin != nil {
			protectedAdmin.GET("/homepage/draft", dependencies.HomepageAdmin.GetDraft)
			protectedAdmin.PUT("/homepage/draft", dependencies.HomepageAdmin.SaveDraft)
			protectedAdmin.GET("/homepage/published", dependencies.HomepageAdmin.GetPublished)
			protectedAdmin.POST("/homepage/publish", dependencies.HomepageAdmin.Publish)
			protectedAdmin.POST("/homepage/reset-draft", dependencies.HomepageAdmin.ResetDraft)
			protectedAdmin.GET("/homepage/preview", dependencies.HomepageAdmin.Preview)
		}
		if dependencies.SiteContentAdmin != nil {
			protectedAdmin.GET("/navigation", dependencies.SiteContentAdmin.ListNavigation)
			protectedAdmin.POST("/navigation", dependencies.SiteContentAdmin.CreateNavigation)
			protectedAdmin.PUT("/navigation/order", dependencies.SiteContentAdmin.ReorderNavigation)
			protectedAdmin.PUT("/navigation/:id", dependencies.SiteContentAdmin.UpdateNavigation)
			protectedAdmin.DELETE("/navigation/:id", dependencies.SiteContentAdmin.DeleteNavigation)
			protectedAdmin.GET("/social-links", dependencies.SiteContentAdmin.ListSocialLinks)
			protectedAdmin.POST("/social-links", dependencies.SiteContentAdmin.CreateSocialLink)
			protectedAdmin.PUT("/social-links/order", dependencies.SiteContentAdmin.ReorderSocialLinks)
			protectedAdmin.PUT("/social-links/:id", dependencies.SiteContentAdmin.UpdateSocialLink)
			protectedAdmin.DELETE("/social-links/:id", dependencies.SiteContentAdmin.DeleteSocialLink)
			protectedAdmin.GET("/featured-articles", dependencies.SiteContentAdmin.ListFeaturedArticles)
			protectedAdmin.POST("/featured-articles", dependencies.SiteContentAdmin.CreateFeaturedArticle)
			protectedAdmin.PUT("/featured-articles/order", dependencies.SiteContentAdmin.ReorderFeaturedArticles)
			protectedAdmin.DELETE("/featured-articles/:articleId", dependencies.SiteContentAdmin.DeleteFeaturedArticle)
		}
		if dependencies.SiteSettingsAdmin != nil {
			protectedAdmin.GET("/site/settings", dependencies.SiteSettingsAdmin.Get)
			protectedAdmin.PUT("/site/settings", dependencies.SiteSettingsAdmin.Update)
		}
		if dependencies.ArticleHandler != nil {
			protectedAdmin.GET("/articles", dependencies.ArticleHandler.ListAdmin)
			protectedAdmin.GET("/articles/:id", dependencies.ArticleHandler.GetAdmin)
			protectedAdmin.POST("/articles", dependencies.ArticleHandler.Create)
			protectedAdmin.PUT("/articles/:id", dependencies.ArticleHandler.Update)
			protectedAdmin.DELETE("/articles/:id", dependencies.ArticleHandler.Delete)
		}
		if dependencies.CommentHandler != nil {
			protectedAdmin.GET("/comments", dependencies.CommentHandler.ListAdmin)
			protectedAdmin.PUT("/comments/:id/status", dependencies.CommentHandler.UpdateStatus)
			protectedAdmin.DELETE("/comments/:id", dependencies.CommentHandler.Delete)
		}
	}

	return engine
}
