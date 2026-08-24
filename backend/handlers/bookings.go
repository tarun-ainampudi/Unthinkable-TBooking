package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"unthinkable-tbooking-backend/models"
)

func (h *Handler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	ctx := r.Context()

	rows, err := h.DB.Query(ctx, `
		SELECT b.id, e.title, e.venue, e.event_date, e.event_time, b.seat_labels, b.total_amount, b.status, b.booking_code, COALESCE(b.qr_code_data_url, '')
		FROM bookings b
		JOIN events e ON b.event_id = e.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC`, userID)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		var seatLabelsJSON []byte
		var qrCodeDataURL string
		rows.Scan(&b.ID, &b.EventTitle, &b.Venue, &b.Date, &b.Time, &seatLabelsJSON, &b.Total, &b.Status, &b.Code, &qrCodeDataURL)
		json.Unmarshal(seatLabelsJSON, &b.SeatLabels)
		b.QRCodeDataUrl = qrCodeDataURL
		bookings = append(bookings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
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
		var price float64
		var status string
		err = tx.QueryRow(ctx, `SELECT status, price FROM seats WHERE event_id = $1 AND seat_label = $2`, req.EventID, seatLabel).Scan(&status, &price)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf(`{"error": "Seat %s does not exist"}`, seatLabel), http.StatusConflict)
				return
			}
			http.Error(w, `{"error": "Database error while checking seat status"}`, http.StatusInternalServerError)
			return
		}
		if status != "available" && status != "held" {
			http.Error(w, fmt.Sprintf(`{"error": "Seat %s is no longer available"}`, seatLabel), http.StatusConflict)
			return
		}
		totalAmount += price

		_, err = tx.Exec(ctx, `UPDATE seats SET status = 'booked', held_until = NULL WHERE event_id = $1 AND seat_label = $2`, req.EventID, seatLabel)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to update seat %s"}`, seatLabel), http.StatusInternalServerError)
			return
		}
	}

	bookingID := fmt.Sprintf("bk_%d", time.Now().Unix())
	bookingCode := fmt.Sprintf("BK-%d-%s", req.EventID, strconv.FormatInt(time.Now().UnixNano()%10000, 10))
	seatLabelsJSON, _ := json.Marshal(req.SeatLabels)
	qrCodeDataURL := generateTicketQRCode(bookingCode)

	_, err = tx.Exec(ctx, `
		INSERT INTO bookings (id, user_id, event_id, seat_labels, total_amount, status, booking_code, qr_code_data_url, created_at)
		VALUES ($1, $2, $3, $4, $5, 'completed', $6, $7, NOW())`,
		bookingID, req.UserID, req.EventID, seatLabelsJSON, totalAmount, bookingCode, qrCodeDataURL)
	if err != nil {
		http.Error(w, `{"error": "Failed to create booking record"}`, http.StatusInternalServerError)
		return
	}

	var recipientEmail string
	err = h.DB.QueryRow(ctx, `SELECT email FROM users WHERE id = $1`, req.UserID).Scan(&recipientEmail)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[Warn] Unable to fetch email for user %d: %v", req.UserID, err)
	}
	if recipientEmail == "" {
		recipientEmail = "guest@example.com"
	}

	if err := simulateBookingEmail(recipientEmail, bookingCode, qrCodeDataURL); err != nil {
		log.Printf("[Warn] Ticket email simulation failed: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		http.Error(w, `{"error": "Failed to commit transaction"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"bookingId":       bookingID,
		"bookingCode":     bookingCode,
		"total":           totalAmount,
		"qrCodeDataUrl":   qrCodeDataURL,
		"email":           recipientEmail,
	})
}
