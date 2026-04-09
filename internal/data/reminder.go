package data

import "time"

type Reminder struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PetID     int       `json:"pet_id"`
	Message   string    `json:"message"`
	RemindAt  time.Time `json:"remind_at"`
	Processed bool      `json:"processed"`
	Channel   string    `json:"channel"`
	// SentAt is a nullable DB field: nil maps to NULL (not sent yet).
	// omitempty keeps sent_at out of JSON when it is nil.
	SentAt       *time.Time `json:"sent_at,omitempty"`
	ErrorMessage *string    `json:"error_message,omitempty"`
}

/*
В БД у тебя для sent_at обычно хранится:

конкретный timestamp, или
NULL.
А в Go это удобно отразить так:

*time.Time:
nil ↔ NULL в БД
&timeValue ↔ заполненное значение времени
Почему не просто time.Time:

у time.Time всегда есть значение (хотя бы 0001-01-01...), и сложно отличить “не задано” от “задано нулевым временем”;
с указателем явно видно, что поле опциональное;
с omitempty в JSON поле не отправляется вообще, если nil.
Итого: указатель нужен для корректной работы с nullable полями (NULL) и более чистой сериализации в JSON.
*/