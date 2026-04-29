package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/andreishemetov/pawpal/internal/app"
	"github.com/andreishemetov/pawpal/internal/config"
	"github.com/andreishemetov/pawpal/internal/logx"
)

func main() {
	logger := logx.New()

	cfg, err := config.LoadForWorker()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	// Root context: OS signals cancel ctx → ReminderWorker.Start exits its select on ctx.Done().
	// DATABASE_URL is validated in config.LoadForWorker and again in app.OpenPostgres before Ping.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.RunReminderWorker(ctx, cfg); err != nil {
		logger.Fatal().Err(err).Msg("worker exited with error")
	}
}
