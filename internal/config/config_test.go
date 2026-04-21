package config

import (
	"testing"
)

func TestLoad_HTTPAddrPrecedence(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x/y")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("REDIS_ADDR", "localhost:6379")

	t.Run("HTTP_ADDR wins over PORT", func(t *testing.T) {
		t.Setenv("HTTP_ADDR", ":3000")
		t.Setenv("PORT", "9999")
		cfg, err := LoadForAPI()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.HTTPAddr != ":3000" {
			t.Fatalf("HTTPAddr = %q, want :3000", cfg.HTTPAddr)
		}
	})

	t.Run("PORT when HTTP_ADDR empty", func(t *testing.T) {
		t.Setenv("HTTP_ADDR", "")
		t.Setenv("PORT", "10000")
		cfg, err := LoadForAPI()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.HTTPAddr != ":10000" {
			t.Fatalf("HTTPAddr = %q, want :10000", cfg.HTTPAddr)
		}
	})

	t.Run("default when both empty", func(t *testing.T) {
		t.Setenv("HTTP_ADDR", "")
		t.Setenv("PORT", "")
		cfg, err := LoadForAPI()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.HTTPAddr != ":8080" {
			t.Fatalf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
		}
	})
}

func TestLoadForAPI_MissingRedis(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x/y")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("REDIS_ADDR", "")
	t.Cleanup(func() {
		t.Setenv("REDIS_ADDR", "")
	})
	_, err := LoadForAPI()
	if err == nil {
		t.Fatal("expected error when REDIS_ADDR missing")
	}
}

func TestLoadForWorker_OnlyDatabase(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x/y")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("REDIS_ADDR", "")
	cfg, err := LoadForWorker()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("DatabaseURL empty")
	}
}
