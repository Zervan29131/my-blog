package config

import (
	"strings"
	"testing"
)

func TestLoadReportsMissingRequiredVariables(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("ADMIN_USERNAME", "")
	t.Setenv("ADMIN_PASSWORD", "")
	t.Setenv("JWT_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected missing environment variables to be rejected")
	}
	for _, key := range []string{"DATABASE_URL", "ADMIN_USERNAME", "ADMIN_PASSWORD", "JWT_SECRET"} {
		if !strings.Contains(err.Error(), key) {
			t.Fatalf("expected error to mention %s, got %v", key, err)
		}
	}
}

func TestLoadRejectsShortJWTSecret(t *testing.T) {
	t.Setenv("DATABASE_URL", "host=localhost dbname=blog")
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD", "safe-password")
	t.Setenv("JWT_SECRET", "too-short")

	if _, err := Load(); err == nil {
		t.Fatal("expected short JWT secret to be rejected")
	}
}

func TestLoadRejectsExampleCredentialsInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATABASE_URL", "host=postgres password=change_me dbname=blog")
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD", "change_me")
	t.Setenv("JWT_SECRET", exampleJWTSecret)

	if _, err := Load(); err == nil {
		t.Fatal("expected example production credentials to be rejected")
	}
}

func TestLoadParsesCORSAllowedOrigins(t *testing.T) {
	t.Setenv("DATABASE_URL", "host=localhost dbname=blog")
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD", "safe-password")
	t.Setenv("JWT_SECRET", "a-safe-jwt-secret-with-at-least-32-characters")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://blog.example.com, http://localhost:5173,https://blog.example.com")

	configuration, err := Load()
	if err != nil {
		t.Fatalf("load configuration: %v", err)
	}
	if len(configuration.CORSAllowedOrigins) != 2 ||
		configuration.CORSAllowedOrigins[0] != "https://blog.example.com" ||
		configuration.CORSAllowedOrigins[1] != "http://localhost:5173" {
		t.Fatalf("unexpected CORS origins: %v", configuration.CORSAllowedOrigins)
	}
}

func TestLoadRejectsInvalidCORSOrigin(t *testing.T) {
	t.Setenv("DATABASE_URL", "host=localhost dbname=blog")
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD", "safe-password")
	t.Setenv("JWT_SECRET", "a-safe-jwt-secret-with-at-least-32-characters")
	t.Setenv("CORS_ALLOWED_ORIGINS", "*")

	if _, err := Load(); err == nil {
		t.Fatal("expected wildcard CORS origin to be rejected")
	}
}
