package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"unthinkable-tbooking-backend/models"
)

func (h *Handler) requireAuthorizedUser(w http.ResponseWriter, r *http.Request) (int, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": "Authorization required"}`, http.StatusUnauthorized)
		return 0, false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		http.Error(w, `{"error": "Session token missing"}`, http.StatusUnauthorized)
		return 0, false
	}

	var userID int
	var expiresAt time.Time
	err := h.DB.QueryRow(r.Context(), `SELECT user_id, expires_at FROM sessions WHERE id = $1 AND revoked = false AND expires_at > NOW()`, token).Scan(&userID, &expiresAt)
	if err != nil {
		http.Error(w, `{"error": "Session expired or invalid"}`, http.StatusUnauthorized)
		return 0, false
	}

	_, err = h.DB.Exec(r.Context(), `UPDATE sessions SET expires_at = NOW() + INTERVAL '15 minutes' WHERE id = $1`, token)
	if err != nil {
		log.Printf("[Warn] Failed to refresh session expiry for token %s: %v", token, err)
	}

	return userID, true
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthorizedUser(w, r)
	if !ok {
		return
	}

	var u models.User
	err := h.DB.QueryRow(r.Context(), `SELECT id, name, email, phone, member_since FROM users WHERE id = $1`, userID).
		Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.MemberSince)
	if err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":      u,
		"expiresAt": time.Now().Add(15 * time.Minute).Format(time.RFC3339),
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid login payload"}`, http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	var u models.User
	var storedPassword string
	err := h.DB.QueryRow(r.Context(), `SELECT id, name, email, phone, member_since, password FROM users WHERE email = $1`, req.Email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.MemberSince, &storedPassword)
	if err != nil || storedPassword != req.Password {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	token := randomSessionToken()
	_, err = h.DB.Exec(r.Context(), `INSERT INTO sessions (id, user_id, expires_at, revoked) VALUES ($1, $2, NOW() + INTERVAL '15 minutes', false) ON CONFLICT (id) DO UPDATE SET user_id = EXCLUDED.user_id, created_at = NOW(), expires_at = NOW() + INTERVAL '15 minutes', revoked = false`, token, u.ID)
	if err != nil {
		http.Error(w, `{"error": "Failed to create session"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"sessionToken": token,
		"expiresAt":    time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		"user":         u,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthorizedUser(w, r)
	if !ok {
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	_, err := h.DB.Exec(r.Context(), `UPDATE sessions SET revoked = true WHERE id = $1 AND user_id = $2`, token, userID)
	if err != nil {
		http.Error(w, `{"error": "Failed to log out"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
