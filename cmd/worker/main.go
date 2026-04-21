package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/andreishemetov/pawpal/internal/app"
	"github.com/andreishemetov/pawpal/internal/config"
)

func main() {
	cfg, err := config.LoadForWorker()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.RunReminderWorker(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
