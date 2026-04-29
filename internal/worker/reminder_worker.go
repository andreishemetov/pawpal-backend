package worker

import (
	"context"
	"time"

	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/rs/zerolog"
)

type ReminderWorker struct {
	repoReminders *repo.ReminderPostgresRepo
	userRepo      *repo.UserPostgresRepo
	notifier      notify.Notifier
	pollInterval  time.Duration
	logger        zerolog.Logger
}

/*
Потому что notify.Notifier — это интерфейс, а не конкретная структура.
В Go для интерфейсов обычно хранят значение как есть (notifier notify.Notifier),
а не *notify.Notifier, потому что:
интерфейс уже сам по себе “ссылка-обертка” (внутри: тип + значение);
*/

func NewReminderWorker(repoReminders *repo.ReminderPostgresRepo, userRepo *repo.UserPostgresRepo, notifier notify.Notifier, pollInterval time.Duration, logger zerolog.Logger) *ReminderWorker {
	if pollInterval < time.Second {
		pollInterval = 5 * time.Second
	}
	return &ReminderWorker{repoReminders: repoReminders, userRepo: userRepo, notifier: notifier, pollInterval: pollInterval, logger: logger}
}

func (w *ReminderWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg("worker stopped")
			return

		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *ReminderWorker) process(ctx context.Context) {
	reminders, err := w.repoReminders.GetDue(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg("failed to fetch due reminders")
		return
	}

	for _, r := range reminders {
		if err := ctx.Err(); err != nil {
			return
		}
		w.logger.Info().
			Int("user_id", r.UserID).
			Int("pet_id", r.PetID).
			Str("channel", r.Channel).
			Msg("processing reminder")

		user, err := w.userRepo.GetByID(ctx, r.UserID)
		if err != nil {
			w.logger.Error().
				Err(err).
				Int("reminder_id", r.ID).
				Int("user_id", r.UserID).
				Msg("failed to load user for reminder")
			_ = w.repoReminders.MarkFailed(ctx, r.ID, err.Error())
			continue
		}

		msg := notify.Message{
			To:      user.Email,
			Title:   "PawPal reminder",
			Body:    r.Message,
			UserID:  r.UserID,
			PetID:   r.PetID,
			Channel: r.Channel,
		}

		if err := w.notifier.Send(ctx, msg); err != nil {
			w.logger.Error().
				Err(err).
				Int("reminder_id", r.ID).
				Msg("failed to send reminder")
			_ = w.repoReminders.MarkFailed(ctx, r.ID, err.Error())
			continue
		}

		if err := w.repoReminders.MarkSent(ctx, r.ID); err != nil {
			w.logger.Error().
				Err(err).
				Int("reminder_id", r.ID).
				Msg("failed to mark reminder as sent")
		}
	}
}
