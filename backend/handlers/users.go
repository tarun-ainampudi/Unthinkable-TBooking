package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"unthinkable-tbooking-backend/models"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	ctx := r.Context()

	var u models.User
	err := h.DB.QueryRow(ctx, `SELECT id, name, email, phone, member_since FROM users WHERE id = $1`, userID).
		Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.MemberSince)
	if err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
