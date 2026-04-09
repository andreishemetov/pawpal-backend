package notify

import (
	"context"
	"log"
)

type LogNotifier struct{}

func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

func (n *LogNotifier) Send(ctx context.Context, msg Message) error {
	log.Printf("[notify] channel=%s to=%s user_id=%d pet_id=%d title=%q body=%q",
		msg.Channel, msg.To, msg.UserID, msg.PetID, msg.Title, msg.Body,
	)
	return nil
}