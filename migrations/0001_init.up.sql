-- docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0008_reminder_delivery_fields.up.sql

CREATE TABLE pets (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT,
  age INT,
  visits INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);