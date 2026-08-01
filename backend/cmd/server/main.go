package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"personal-blog/backend/internal/config"
	"personal-blog/backend/internal/database"
	"personal-blog/backend/internal/handler"
	"personal-blog/backend/internal/middleware"
	"personal-blog/backend/internal/router"
	"personal-blog/backend/internal/service"
)

func main() {
	configuration, err := config.Load()
	if err != nil {
		slog.Error("configuration validation failed", "error", err)
		os.Exit(1)
	}

	connectContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databaseConnection, err := database.Open(connectContext, configuration.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed")
		os.Exit(1)
	}
	slog.Info("database connected")

	sqlDatabase, err := database.SQL(databaseConnection)
	if err != nil {
		slog.Error("database connection setup failed")
		os.Exit(1)
	}
	defer sqlDatabase.Close()

	if err := database.Migrate(databaseConnection); err != nil {
		slog.Error("database migration failed", "error", err)
		os.Exit(1)
	}
	if err := database.InitializeHomepageCMS(databaseConnection); err != nil {
		slog.Error("homepage CMS initialization failed", "error", err)
		os.Exit(1)
	}

	authService := service.NewAuthService(
		databaseConnection,
		configuration.JWTSecret,
		configuration.JWTExpiresHours,
	)
	created, err := authService.InitializeAdministrator(
		context.Background(),
		configuration.AdminUsername,
		configuration.AdminPassword,
	)
	if err != nil {
		slog.Error("administrator initialization failed", "error", err)
		os.Exit(1)
	}
	if created {
		slog.Info("administrator initialized", "username", configuration.AdminUsername)
	}

	authHandler := handler.NewAuthHandler(authService)
	articleService := service.NewArticleService(databaseConnection)
	articleHandler := handler.NewArticleHandler(articleService)
	commentService := service.NewCommentService(databaseConnection)
	commentHandler := handler.NewCommentHandler(commentService)
	dashboardHandler := handler.NewDashboardHandler(service.NewDashboardService(databaseConnection))
	publicConfigHandler := handler.NewPublicConfigHandler(service.NewPublicConfigService(databaseConnection))
	homepageAdminHandler := handler.NewHomepageAdminHandler(service.NewHomepageAdminService(databaseConnection))
	siteContentAdminHandler := handler.NewSiteContentAdminHandler(
		service.NewNavigationAdminService(databaseConnection),
		service.NewSocialLinkAdminService(databaseConnection),
		service.NewFeaturedArticleAdminService(databaseConnection),
	)
	siteSettingsAdminHandler := handler.NewSiteSettingsAdminHandler(
		service.NewSiteSettingsAdminService(databaseConnection),
	)
	commentRateLimiter := middleware.NewCommentRateLimiter(5, time.Minute)
	engine := router.New(router.Dependencies{
		AuthHandler:       authHandler,
		ArticleHandler:    articleHandler,
		CommentHandler:    commentHandler,
		DashboardHandler:  dashboardHandler,
		PublicConfig:      publicConfigHandler,
		HomepageAdmin:     homepageAdminHandler,
		SiteContentAdmin:  siteContentAdminHandler,
		SiteSettingsAdmin: siteSettingsAdminHandler,
		AuthMiddleware:    middleware.Authenticate(authService),
		CommentRateLimit:  commentRateLimiter.Handler(),
		CORSMiddleware:    middleware.CORS(configuration.CORSAllowedOrigins),
	})
	address := ":" + configuration.ServerPort
	slog.Info("server starting", "address", address)

	if err := engine.Run(address); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
