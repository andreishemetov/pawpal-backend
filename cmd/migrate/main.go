package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	dir := os.Getenv("MIGRATIONS_PATH")
	if dir == "" {
		dir = "migrations"
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("migrations path: %v", err)
	}

	// file:// URL with absolute path (Unix-friendly for Render/Docker).
	srcURL := fmt.Sprintf("file://%s", filepath.ToSlash(abs))

	m, err := migrate.New(srcURL, dsn)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migrations: schema already up to date")
			return
		}
		log.Fatalf("migrate up: %v", err)
	}
	log.Println("migrations: applied successfully")
}
