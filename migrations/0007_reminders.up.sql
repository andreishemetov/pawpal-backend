CREATE TABLE reminders (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id),
  pet_id INT NOT NULL REFERENCES pets(id),
  message TEXT NOT NULL,
  remind_at TIMESTAMPTZ NOT NULL,
  processed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_reminders_due ON reminders(remind_at, processed);