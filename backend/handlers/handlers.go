package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skip2/go-qrcode"
	"unthinkable-tbooking-backend/models"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

func generateTicketQRCode(bookingCode string) string {
	png, err := qrcode.Encode(bookingCode, qrcode.Medium, 256)
	if err != nil {
		log.Printf("[Warn] QR generation failed for %s: %v", bookingCode, err)
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}

func simulateBookingEmail(recipient, bookingCode, qrDataURL string) error {
	log.Printf("[Info] Sending ticket email to %s for booking %s (QR length=%d)", recipient, bookingCode, len(qrDataURL))
	return nil
}

func randomSessionToken() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("session-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

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

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := h.DB.Query(ctx, `SELECT id, title, category, tag, venue, event_date, event_time, price_from, rating, accent, blurb FROM events ORDER BY id`)
	if err != nil {
		log.Printf("[Info] GetEvents: %v\n", err)
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Category, &e.Tag, &e.Venue, &e.Date, &e.Time, &e.PriceFrom, &e.Rating, &e.Accent, &e.Blurb); err != nil {
			log.Println("[Info] GetEvents scan error:", err)
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

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

func (h *Handler) GetCategoryMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.CategoryPricing)
}