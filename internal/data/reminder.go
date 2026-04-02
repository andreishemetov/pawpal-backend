package data

import "time"

type Reminder struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PetID     int       `json:"pet_id"`
	Message   string    `json:"message"`
	RemindAt  time.Time `json:"remind_at"`
	Processed bool      `json:"processed"`
}