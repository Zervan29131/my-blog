package config

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	defaultServerPort      = "8080"
	defaultJWTExpiresHours = 24
	examplePassword        = "change_me"
	exampleJWTSecret       = "replace_with_a_random_string_at_least_32_characters"
)

type Config struct {
	AppEnv             string
	ServerPort         string
	DatabaseURL        string
	AdminUsername      string
	AdminPassword      string
	JWTSecret          string
	JWTExpiresHours    int
	CORSAllowedOrigins []string
}

func Load() (Config, error) {
	config := Config{
		AppEnv:          valueOrDefault("APP_ENV", "development"),
		ServerPort:      valueOrDefault("SERVER_PORT", defaultServerPort),
		DatabaseURL:     strings.TrimSpace(os.Getenv("DATABASE_URL")),
		AdminUsername:   strings.TrimSpace(os.Getenv("ADMIN_USERNAME")),
		AdminPassword:   os.Getenv("ADMIN_PASSWORD"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExpiresHours: defaultJWTExpiresHours,
	}
	allowedOrigins, err := parseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if err != nil {
		return Config{}, err
	}
	config.CORSAllowedOrigins = allowedOrigins

	if value := strings.TrimSpace(os.Getenv("JWT_EXPIRES_HOURS")); value != "" {
		hours, err := strconv.Atoi(value)
		if err != nil || hours <= 0 {
			return Config{}, fmt.Errorf("JWT_EXPIRES_HOURS must be a positive integer")
		}
		config.JWTExpiresHours = hours
	}

	if err := config.validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func parseAllowedOrigins(value string) ([]string, error) {
	origins := make([]string, 0)
	seen := make(map[string]struct{})
	for _, rawOrigin := range strings.Split(value, ",") {
		origin := strings.TrimSpace(rawOrigin)
		if origin == "" {
			continue
		}
		parsed, err := url.Parse(origin)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
			parsed.Host == "" || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS contains invalid origin %q", origin)
		}
		if _, exists := seen[origin]; exists {
			continue
		}
		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}
	return origins, nil
}

func (config Config) validate() error {
	missing := make([]string, 0, 4)
	for key, value := range map[string]string{
		"DATABASE_URL":   config.DatabaseURL,
		"ADMIN_USERNAME": config.AdminUsername,
		"ADMIN_PASSWORD": config.AdminPassword,
		"JWT_SECRET":     config.JWTSecret,
	} {
		if value == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	if utf8.RuneCountInString(config.AdminUsername) > 50 {
		return fmt.Errorf("ADMIN_USERNAME must not exceed 50 characters")
	}
	if len(config.AdminPassword) > 72 {
		return fmt.Errorf("ADMIN_PASSWORD must not exceed 72 bytes")
	}
	if len(config.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	if strings.EqualFold(config.AppEnv, "production") {
		if config.AdminPassword == examplePassword ||
			config.JWTSecret == exampleJWTSecret ||
			strings.Contains(config.DatabaseURL, "password="+examplePassword) {
			return fmt.Errorf("example credentials are not allowed in production")
		}
	}

	return nil
}

func valueOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
