package worker

import (
	"context"
	"log"
	"time"

	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/repo"
)

type ReminderWorker struct {
	repoReminders *repo.ReminderPostgresRepo
	userRepo      *repo.UserPostgresRepo
	notifier      notify.Notifier
}

/*
Потому что notify.Notifier — это интерфейс, а не конкретная структура.
В Go для интерфейсов обычно хранят значение как есть (notifier notify.Notifier),
а не *notify.Notifier, потому что:
интерфейс уже сам по себе “ссылка-обертка” (внутри: тип + значение);
*/

func NewReminderWorker(repoReminders *repo.ReminderPostgresRepo, userRepo *repo.UserPostgresRepo, notifier notify.Notifier) *ReminderWorker {
	return &ReminderWorker{repoReminders: repoReminders, userRepo: userRepo, notifier: notifier}
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
	reminders, err := w.repoReminders.GetDue(ctx)
	if err != nil {
		log.Println("worker error:", err)
		return
	}

	for _, r := range reminders {
		user, err := w.userRepo.GetByID(ctx, r.UserID)
		if err != nil {
			log.Println("worker user lookup error:", err)
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

		if err := w.notifier.Send(ctx, msg); 
		err != nil {
			log.Println("notify error:", err)
			_ = w.repoReminders.MarkFailed(ctx, r.ID, err.Error())
			continue
		}

		if err := w.repoReminders.MarkSent(ctx, r.ID); err != nil {
			log.Println("mark sent error:", err)
		}
	}
}
