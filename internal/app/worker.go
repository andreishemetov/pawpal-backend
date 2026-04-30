package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/andreishemetov/pawpal/internal/config"
	"github.com/andreishemetov/pawpal/internal/logx"
	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/redisx"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/andreishemetov/pawpal/internal/worker"
	// "github.com/andreishemetov/pawpal/internal/worker"
)

// RunReminderWorker runs the reminder delivery loop until ctx is cancelled (e.g. SIGTERM on Render).
func RunReminderWorker(ctx context.Context, cfg *config.Config) error {

	logger := logx.New()
	logger.Info().Msg("starting worker")

	db, err := OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is required")
	}

	redisClient := redisx.NewClient(redisAddr)

	reminderRepo := repo.NewReminderPostgresRepo(db)
	userRepo := repo.NewUserPostgresRepo(db)

	dispatcher := notify.NewDispatcher(map[string]notify.Notifier{
		"log":   notify.NewLogNotifier(),
		"email": notify.NewEmailNotifier(cfg.NotifyEmailFrom),
	})

	// reminderWorker := worker.NewReminderWorker(reminderRepo, userRepo, dispatcher, cfg.WorkerPollInterval, logger)

	// logger.Info().Msg("reminder worker started")
	// reminderWorker.Start(ctx)
	// logger.Info().Msg("reminder worker stopped")

	streamWorker := worker.NewReminderStreamWorker(
		redisClient,
		reminderRepo,
		userRepo,
		dispatcher,
		"worker-1",
	)

	if err := streamWorker.Start(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}

	logger.Info().Msg("stream worker started")

	return nil
}
