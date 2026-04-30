package events

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	client *redis.Client
}

func NewProducer(client *redis.Client) *Producer {
	return &Producer{client: client}
}

func (p *Producer) PublishReminder(ctx context.Context, reminderID int, userID int, petID int, channel string) error {
	_, err := p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: ReminderStream,
		Values: map[string]any{
			"reminder_id": strconv.Itoa(reminderID),
			"user_id":     strconv.Itoa(userID),
			"pet_id":      strconv.Itoa(petID),
			"channel":     channel,
		},
	}).Result()

	return err
}