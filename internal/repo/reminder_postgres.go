package repo

import (
	"context"
	"database/sql"

	"github.com/andreishemetov/pawpal/internal/data"
)

type ReminderPostgresRepo struct {
	db *sql.DB
}

func NewReminderPostgresRepo(db *sql.DB) *ReminderPostgresRepo {
	return &ReminderPostgresRepo{db: db}
}

func (r *ReminderPostgresRepo) Create(ctx context.Context, rem data.Reminder) (data.Reminder, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO reminders (user_id, pet_id, message, remind_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, rem.UserID, rem.PetID, rem.Message, rem.RemindAt).
		Scan(&rem.ID)

	return rem, err
}

func (r *ReminderPostgresRepo) GetDue(ctx context.Context) ([]data.Reminder, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, pet_id, message, remind_at
		FROM reminders
		WHERE processed = false
		  AND remind_at <= now()
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []data.Reminder

	for rows.Next() {
		var rem data.Reminder
		if err := rows.Scan(&rem.ID, &rem.UserID, &rem.PetID, &rem.Message, &rem.RemindAt); err != nil {
			return nil, err
		}
		reminders = append(reminders, rem)
	}

	return reminders, rows.Err()
}

func (r *ReminderPostgresRepo) MarkProcessed(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE reminders
		SET processed = true
		WHERE id = $1
	`, id)

	return err
}