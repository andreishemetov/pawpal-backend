package notify

import "context"

type Message struct {
	To      string
	Title   string
	Body    string
	UserID  int
	PetID   int
	Channel string
}

type Notifier interface {
	Send(ctx context.Context, msg Message) error
}