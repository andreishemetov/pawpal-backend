package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime settings read from the environment.
// Call LoadForAPI or LoadForWorker to obtain a validated instance for each binary.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	RedisAddr   string

	// HTTPAddr is the address passed to http.ListenAndServe (e.g. ":8080").
	// Resolved from HTTP_ADDR, or from PORT (Render), or default ":8080".
	HTTPAddr string

	// Worker / notifications
	WorkerPollInterval time.Duration
	NotifyEmailFrom    string

	// API tuning (cache & rate limit)
	PetCacheTTL     time.Duration
	RateLimitMax    int
	RateLimitWindow time.Duration
	HTTPTimeout     time.Duration
}

// LoadForAPI loads configuration from the environment and validates fields required by the HTTP API.
func LoadForAPI() (*Config, error) {
	c, err := load()
	if err != nil {
		return nil, err
	}
	if err := c.validateAPI(); err != nil {
		return nil, err
	}
	return c, nil
}

// LoadForWorker loads configuration from the environment and validates fields required by background workers.
func LoadForWorker() (*Config, error) {
	c, err := load()
	if err != nil {
		return nil, err
	}
	if err := c.validateWorker(); err != nil {
		return nil, err
	}
	return c, nil
}

func load() (*Config, error) {
	c := &Config{
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSecret:   strings.TrimSpace(os.Getenv("JWT_SECRET")),
		RedisAddr:   strings.TrimSpace(os.Getenv("REDIS_ADDR")),
		NotifyEmailFrom: strings.TrimSpace(os.Getenv("NOTIFY_EMAIL_FROM")),
	}

	httpAddr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	port := strings.TrimSpace(os.Getenv("PORT"))
	switch {
	case httpAddr != "":
		c.HTTPAddr = httpAddr
	case port != "":
		c.HTTPAddr = ":" + port
	default:
		c.HTTPAddr = ":8080"
	}

	var err error
	if c.WorkerPollInterval, err = parseDurationEnv("REMINDER_WORKER_POLL_INTERVAL", 5*time.Second); err != nil {
		return nil, err
	}
	if c.WorkerPollInterval < time.Second {
		return nil, fmt.Errorf("REMINDER_WORKER_POLL_INTERVAL must be at least 1s")
	}

	if c.NotifyEmailFrom == "" {
		c.NotifyEmailFrom = "noreply@pawpal.local"
	}

	if c.PetCacheTTL, err = parseDurationEnv("PET_CACHE_TTL", 30*time.Second); err != nil {
		return nil, err
	}
	if c.PetCacheTTL < time.Second {
		return nil, fmt.Errorf("PET_CACHE_TTL must be at least 1s")
	}

	if c.RateLimitMax, err = parsePositiveIntEnv("RATE_LIMIT_MAX", 10); err != nil {
		return nil, err
	}

	if c.RateLimitWindow, err = parseDurationEnv("RATE_LIMIT_WINDOW", time.Minute); err != nil {
		return nil, err
	}
	if c.RateLimitWindow < time.Second {
		return nil, fmt.Errorf("RATE_LIMIT_WINDOW must be at least 1s")
	}

	if c.HTTPTimeout, err = parseDurationEnv("HTTP_TIMEOUT", 15*time.Second); err != nil {
		return nil, err
	}
	if c.HTTPTimeout < time.Second {
		return nil, fmt.Errorf("HTTP_TIMEOUT must be at least 1s")
	}

	return c, nil
}

func (c *Config) validateAPI() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.RedisAddr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}
	return nil
}

func (c *Config) validateWorker() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func parseDurationEnv(key string, defaultVal time.Duration) (time.Duration, error) {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return defaultVal, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid duration %q: %w", key, s, err)
	}
	return d, nil
}

func parsePositiveIntEnv(key string, defaultVal int) (int, error) {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return n, nil
}
