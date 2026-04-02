package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/middleware"
	"github.com/andreishemetov/pawpal/internal/repo"
)

type ReminderHandler struct {
	repo *repo.ReminderPostgresRepo
}

func NewReminderHandler(r *repo.ReminderPostgresRepo) *ReminderHandler {
	return &ReminderHandler{repo: r}
}

func (h *ReminderHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var rem data.Reminder
	if err := json.NewDecoder(r.Body).Decode(&rem); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	rem.UserID = userID

	if rem.PetID <= 0 {
		http.Error(w, "pet_id is required and must be positive", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(rem.Message) == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}
	if rem.RemindAt.IsZero() {
		http.Error(w, "remind_at is required (use RFC3339, e.g. 2026-03-30T15:00:00Z)", http.StatusBadRequest)
		return
	}

	created, err := h.repo.Create(r.Context(), rem)
	if err != nil {
		log.Printf("reminder Create: %v", err)
		status, msg := reminderCreateHTTPError(err)
		http.Error(w, msg, status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func reminderCreateHTTPError(err error) (status int, msg string) {
	s := err.Error()
	switch {
	case strings.Contains(s, "foreign key constraint"):
		return http.StatusBadRequest, "invalid user_id or pet_id (must reference existing user and pet)"
	case strings.Contains(s, "violates not-null constraint"), strings.Contains(s, "NOT NULL"):
		return http.StatusBadRequest, "missing required field for database"
	default:
		return http.StatusInternalServerError, "failed to create reminder"
	}
}
