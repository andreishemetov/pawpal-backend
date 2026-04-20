package repo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/andreishemetov/pawpal/internal/data"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestTx(t *testing.T) (*sql.Tx, *sql.DB) {
	dsn := "postgres://pawpal:pawpal_pass@localhost:5432/pawpal_dev?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	return tx, db
}

func insertTestUser(t *testing.T, tx *sql.Tx, email string) int {
	t.Helper()

	var id int
	err := tx.QueryRow(
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`,
		email, "test-hash", "user",
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	return id
}

func TestPetPostgresRepo_AddAndGet(t *testing.T) {
	tx, db := setupTestTx(t)
	defer tx.Rollback()
	defer db.Close()

	repo := NewPetPostgresRepo(tx)

	ctx := context.Background()
	userID := insertTestUser(t, tx, "charlie@example.com")

	pet := data.Pet{
		Name:   "Charlie",
		Type:   "Dog",
		Age:    3,
		UserID: userID,
	}

	created, err := repo.Add(ctx, pet)
	if err != nil {
		t.Fatal(err)
	}

	if created.ID == 0 {
		t.Fatal("expected ID to be set")
	}

	found, err := repo.GetByID(ctx, userID, "user", created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if found.Name != "Charlie" {
		t.Fatalf("expected Charlie, got %s", found.Name)
	}
}

func TestPetPostgresRepo_Delete(t *testing.T) {
	tx, db := setupTestTx(t)
	defer tx.Rollback()
	defer db.Close()

	repo := NewPetPostgresRepo(tx)
	ctx := context.Background()
	userID := insertTestUser(t, tx, "milo@example.com")

	pet, _ := repo.Add(ctx, data.Pet{Name: "Milo", UserID: userID})

	deleted, err := repo.DeleteByID(ctx, userID, "user", pet.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Fatal("expected deleted = true")
	}

	_, err = repo.GetByID(ctx, userID, "user", pet.ID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
