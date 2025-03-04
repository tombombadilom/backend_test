package config

import (
	"flag"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port int

	// Database configuration
	DBPath string

	// Logging configuration
	LogLevel string

	// API configuration
	APIKeyExpireDays int
	RateLimitPerMin  int
}

// New creates a new configuration with values from environment variables or flags
func New() *Config {
	cfg := &Config{
		Port:             8080,
		DBPath:           "./liveops.db",
		LogLevel:         "info",
		APIKeyExpireDays: 30,
		RateLimitPerMin:  60,
	}

	// Override with environment variables if present
	if port, err := strconv.Atoi(os.Getenv("LIVEOPS_PORT")); err == nil && port > 0 {
		cfg.Port = port
	}

	if dbPath := os.Getenv("LIVEOPS_DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}

	if logLevel := os.Getenv("LIVEOPS_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	if days, err := strconv.Atoi(os.Getenv("LIVEOPS_API_KEY_EXPIRE_DAYS")); err == nil && days > 0 {
		cfg.APIKeyExpireDays = days
	}

	if rate, err := strconv.Atoi(os.Getenv("LIVEOPS_RATE_LIMIT")); err == nil && rate > 0 {
		cfg.RateLimitPerMin = rate
	}

	return cfg
}

// ParseFlags parses command line flags and updates the configuration
func (c *Config) ParseFlags() {
	flag.IntVar(&c.Port, "port", c.Port, "Server port")
	flag.StringVar(&c.DBPath, "db", c.DBPath, "SQLite database path")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Log level (debug, info, warn, error)")
	flag.IntVar(&c.APIKeyExpireDays, "api-key-expire", c.APIKeyExpireDays, "API key expiration in days")
	flag.IntVar(&c.RateLimitPerMin, "rate-limit", c.RateLimitPerMin, "Rate limit per minute")

	flag.Parse()
}
