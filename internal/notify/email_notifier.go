package notify

import (
	"context"
	"fmt"
)

type EmailNotifier struct {
	from string
}

func NewEmailNotifier(from string) *EmailNotifier {
	return &EmailNotifier{from: from}
}

func (n *EmailNotifier) Send(ctx context.Context, msg Message) error {
	// Placeholder for real SMTP/provider integration
	fmt.Printf("[email] from=%s to=%s title=%q body=%q\n",
		n.from, msg.To, msg.Title, msg.Body,
	)
	return nil
}