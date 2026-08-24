package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"unthinkable-tbooking-backend/models"
)

func (h *Handler) GetEventSeats(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	rows, err := h.DB.Query(ctx, `SELECT seat_label, row_name, col_num, category, status, price FROM seats WHERE event_id = $1 ORDER BY row_name, col_num`, eventID)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var s models.Seat
		rows.Scan(&s.ID, &s.Row, &s.Col, &s.Category, &s.Status, &s.Price)
		seats = append(seats, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seats)
}

func (h *Handler) HoldSeats(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthorizedUser(w, r)
	if !ok {
		return
	}

	var req struct {
		UserID     int      `json:"userId"`
		EventID    int      `json:"eventId"`
		SeatLabels []string `json:"seatLabels"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	if req.UserID != 0 && req.UserID != userID {
		http.Error(w, `{"error": "Unauthorized user session"}`, http.StatusForbidden)
		return
	}
	if len(req.SeatLabels) == 0 {
		http.Error(w, `{"error": "No seats selected"}`, http.StatusBadRequest)
		return
	}
	req.UserID = userID

	ctx := r.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		http.Error(w, `{"error": "Failed to start transaction"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var totalAmount float64
	for _, seatLabel := range req.SeatLabels {
		var status string
		var price float64
		err = tx.QueryRow(ctx, `SELECT status, price FROM seats WHERE event_id = $1 AND seat_label = $2`, req.EventID, seatLabel).Scan(&status, &price)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf(`{"error": "Seat %s does not exist"}`, seatLabel), http.StatusConflict)
				return
			}
			http.Error(w, `{"error": "Database error while checking seat availability"}`, http.StatusInternalServerError)
			return
		}
		if status == "booked" || status == "held" {
			http.Error(w, fmt.Sprintf(`{"error": "Seat %s is already held or booked"}`, seatLabel), http.StatusConflict)
			return
		}
		_, err = tx.Exec(ctx, `UPDATE seats SET status = 'held', held_until = NOW() + INTERVAL '8 minutes' WHERE event_id = $1 AND seat_label = $2 AND status IN ('available', 'held')`, req.EventID, seatLabel)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to hold seat %s"}`, seatLabel), http.StatusInternalServerError)
			return
		}
		totalAmount += price
	}

	if err = tx.Commit(ctx); err != nil {
		http.Error(w, `{"error": "Failed to reserve seats"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"eventId":      req.EventID,
		"seatLabels":   req.SeatLabels,
		"total":        totalAmount,
		"expiresInSec": 480,
	})
}

func (h *Handler) ReleaseHeldSeats(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthorizedUser(w, r)
	if !ok {
		return
	}

	var req struct {
		UserID     int      `json:"userId"`
		EventID    int      `json:"eventId"`
		SeatLabels []string `json:"seatLabels"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	if req.UserID != 0 && req.UserID != userID {
		http.Error(w, `{"error": "Unauthorized user session"}`, http.StatusForbidden)
		return
	}
	if len(req.SeatLabels) == 0 {
		http.Error(w, `{"error": "No seats selected"}`, http.StatusBadRequest)
		return
	}
	req.UserID = userID

	ctx := r.Context()
	for _, seatLabel := range req.SeatLabels {
		var status string
		err := h.DB.QueryRow(ctx, `SELECT status FROM seats WHERE event_id = $1 AND seat_label = $2`, req.EventID, seatLabel).Scan(&status)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			http.Error(w, `{"error": "Database error while releasing seat hold"}`, http.StatusInternalServerError)
			return
		}
		if status == "booked" {
			http.Error(w, fmt.Sprintf(`{"error": "Seat %s is already booked and cannot be released"}`, seatLabel), http.StatusConflict)
			return
		}
		_, err = h.DB.Exec(ctx, `UPDATE seats SET status = 'available', held_until = NULL WHERE event_id = $1 AND seat_label = $2 AND status = 'held'`, req.EventID, seatLabel)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to release seat %s"}`, seatLabel), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"eventId":    req.EventID,
		"seatLabels": req.SeatLabels,
	})
}
