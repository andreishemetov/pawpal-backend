package notify

import (
	"context"
	"fmt"
)

type NotifiersDispatcher struct {
	channels map[string]Notifier
}

func NewDispatcher(channels map[string]Notifier) *NotifiersDispatcher {
	return &NotifiersDispatcher{channels: channels}
}

func (d *NotifiersDispatcher) Send(ctx context.Context, msg Message) error {
	n, ok := d.channels[msg.Channel]
	if !ok {
		return fmt.Errorf("unknown channel: %s", msg.Channel)
	}
	return n.Send(ctx, msg)
}