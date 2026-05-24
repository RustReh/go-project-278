package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	SentryDsn    string
	DataBaseUrl  string
	BaseURL      string
	CORSOrigin   string
}

func normalizeBaseURL(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return ""
	}
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return base
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func NewAppConfig() (*AppConfig, error) {
	host := envOrDefault("SERVER_HOST", "localhost")
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = envOrDefault("SERVER_PORT", "8080")
	}
	if os.Getenv("PORT") != "" && os.Getenv("SERVER_HOST") == "" {
		host = "0.0.0.0"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT %q: %w", portStr, err)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("SERVER_PORT must be between 1 and 65535, got %d", port)
	}
	readTimeoutStr := envOrDefault("SERVER_READ_TIMEOUT", "10")
	readSec, err := strconv.Atoi(readTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT %q: %w", readTimeoutStr, err)
	}
	readTimeout := time.Duration(readSec) * time.Second
	writeTimeoutStr := envOrDefault("SERVER_WRITE_TIMEOUT", "10")
	writeSec, err := strconv.Atoi(writeTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT %q: %w", writeTimeoutStr, err)
	}
	writeTimeout := time.Duration(writeSec) * time.Second

	address := fmt.Sprintf("%s:%d", host, port)
	sentryDsn := envOrDefault("SENTRY_DSN", "")
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	baseURL := normalizeBaseURL(os.Getenv("BASE_URL"))
	if baseURL == "" {
		return nil, fmt.Errorf("BASE_URL is required")
	}
	corsOrigin := envOrDefault("CORS_ORIGIN", "http://localhost:5173")

	return &AppConfig{
		Addr:         address,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		SentryDsn:    sentryDsn,
		DataBaseUrl:  databaseURL,
		BaseURL:      baseURL,
		CORSOrigin:   corsOrigin,
	}, nil
}
