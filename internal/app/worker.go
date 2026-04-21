package app

import (
	"context"
	"fmt"
	"log"

	"github.com/andreishemetov/pawpal/internal/config"
	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/andreishemetov/pawpal/internal/worker"
)

// RunReminderWorker runs the reminder delivery loop until ctx is cancelled (e.g. SIGTERM on Render).
func RunReminderWorker(ctx context.Context, cfg *config.Config) error {
	db, err := OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	reminderRepo := repo.NewReminderPostgresRepo(db)
	userRepo := repo.NewUserPostgresRepo(db)

	dispatcher := notify.NewDispatcher(map[string]notify.Notifier{
		"log":   notify.NewLogNotifier(),
		"email": notify.NewEmailNotifier(cfg.NotifyEmailFrom),
	})

	reminderWorker := worker.NewReminderWorker(reminderRepo, userRepo, dispatcher, cfg.WorkerPollInterval)

	log.Println("reminder worker started")
	reminderWorker.Start(ctx)
	log.Println("reminder worker stopped")
	return nil
}
