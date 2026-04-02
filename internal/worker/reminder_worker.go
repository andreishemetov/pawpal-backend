package worker

import (
	"context"
	"log"
	"time"

	"github.com/andreishemetov/pawpal/internal/repo"
)

type ReminderWorker struct {
	repo *repo.ReminderPostgresRepo
}

func NewReminderWorker(r *repo.ReminderPostgresRepo) *ReminderWorker {
	return &ReminderWorker{repo: r}
}

func (w *ReminderWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped")
			return

		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *ReminderWorker) process(ctx context.Context) {
	reminders, err := w.repo.GetDue(ctx)
	if err != nil {
		log.Println("worker error:", err)
		return
	}

	for _, r := range reminders {
		log.Printf("🔔 Reminder: user=%d pet=%d message=%s\n",
			r.UserID, r.PetID, r.Message)

		_ = w.repo.MarkProcessed(ctx, r.ID)
	}
}