package worker

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/andreishemetov/pawpal/internal/events"
	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/repo"
)

type ReminderStreamWorker struct {
	client    *redis.Client
	reminders *repo.ReminderPostgresRepo
	users     *repo.UserPostgresRepo
	notifier  notify.Notifier
	consumer  string
}

func NewReminderStreamWorker(
	client *redis.Client,
	reminders *repo.ReminderPostgresRepo,
	users *repo.UserPostgresRepo,
	notifier notify.Notifier,
	consumer string,
) *ReminderStreamWorker {
	return &ReminderStreamWorker{
		client:    client,
		reminders: reminders,
		users:     users,
		notifier:  notifier,
		consumer:  consumer,
	}
}

func (w *ReminderStreamWorker) EnsureGroup(ctx context.Context) error {
	err := w.client.XGroupCreateMkStream(ctx, events.ReminderStream, events.ReminderGroup, "$").Err()
	if err != nil && !isBusyGroupErr(err) {
		return err
	}
	return nil
}

func isBusyGroupErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "BUSYGROUP")
}

func (w *ReminderStreamWorker) Start(ctx context.Context) error {
	if err := w.EnsureGroup(ctx); err != nil {
		return err
	}

	log.Println("stream worker started")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		streams, err := w.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    events.ReminderGroup,
			Consumer: w.consumer,
			Streams:  []string{events.ReminderStream, ">"},
			Count:    10,
			Block:    0,
		}).Result()

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			log.Println("xreadgroup error:", err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				if err := w.handleMessage(ctx, msg); err != nil {
					log.Println("handle message error:", err)
					continue
				}

				if err := w.client.XAck(ctx, events.ReminderStream, events.ReminderGroup, msg.ID).Err(); err != nil {
					log.Println("xack error:", err)
				}
			}
		}
	}
}

func (w *ReminderStreamWorker) handleMessage(ctx context.Context, msg redis.XMessage) error {
	reminderID, err := strconv.Atoi(asString(msg.Values["reminder_id"]))
	if err != nil {
		return err
	}

	rem, err := w.reminders.GetByID(ctx, reminderID)
	if err != nil {
		return err
	}

	user, err := w.users.GetByID(ctx, rem.UserID)
	if err != nil {
		return err
	}

	nmsg := notify.Message{
		To:      user.Email,
		Title:   "PawPal reminder",
		Body:    rem.Message,
		UserID:  rem.UserID,
		PetID:   rem.PetID,
		Channel: rem.Channel,
	}

	if err := w.notifier.Send(ctx, nmsg); err != nil {
		_ = w.reminders.MarkFailed(ctx, rem.ID, err.Error())
		return err
	}

	return w.reminders.MarkSent(ctx, rem.ID)
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		return ""
	}
}