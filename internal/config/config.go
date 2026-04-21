package config

import (
	"fmt"
	"os"
)

// Config holds runtime settings loaded from the environment.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	RedisAddr   string
	HTTPAddr    string
}

// LoadForAPI reads and validates env vars required by the HTTP API.
func LoadForAPI() (*Config, error) {
	c := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		HTTPAddr:    os.Getenv("HTTP_ADDR"),
	}
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if c.RedisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR is required")
	}
	if c.HTTPAddr == "" {
		c.HTTPAddr = ":8080"
	}
	return c, nil
}

// LoadForWorker reads and validates env vars required by background workers.
func LoadForWorker() (*Config, error) {
	c := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	return c, nil
}
