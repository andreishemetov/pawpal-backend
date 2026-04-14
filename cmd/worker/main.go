package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/andreishemetov/pawpal/internal/worker"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	// dsn := "postgres://pawpal:pawpal_pass@localhost:5432/pawpal_dev?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Initialize repositories
	reminderRepo := repo.NewReminderPostgresRepo(db)
	userRepo := repo.NewUserPostgresRepo(db)

	dispatcher := notify.NewDispatcher(map[string]notify.Notifier{
		"log":   notify.NewLogNotifier(),
		"email": notify.NewEmailNotifier("noreply@pawpal.local"),
	})

	// worker
	reminderWorker := worker.NewReminderWorker(reminderRepo, userRepo, dispatcher)

	log.Println("worker started")
	reminderWorker.Start(context.Background())

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
}
