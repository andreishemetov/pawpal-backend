package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreishemetov/pawpal/internal/config"
	"github.com/andreishemetov/pawpal/internal/logx"
	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/redisx"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/andreishemetov/pawpal/internal/worker"
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

	redisClient := redisx.NewClient(cfg.RedisAddr)

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

	logger.Info().Msg("stream worker started")
	err = streamWorker.Start(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("stream worker: %w", err)
	}
	logger.Info().Msg("stream worker stopped")
	return nil
}
