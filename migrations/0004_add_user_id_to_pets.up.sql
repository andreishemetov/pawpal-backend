ALTER TABLE pets
ADD COLUMN user_id INT REFERENCES users(id);

CREATE INDEX idx_pets_user_id ON pets(user_id);