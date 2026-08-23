package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"unthinkable-tbooking-backend/models"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := h.DB.Query(ctx, `SELECT id, title, category, tag, venue, event_date, event_time, price_from, rating, accent, blurb FROM events ORDER BY id`)
	if err != nil {
		log.Println("[Info] GetEvents: %v\n",err)
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		rows.Scan(&e.ID, &e.Title, &e.Category, &e.Tag, &e.Venue, &e.Date, &e.Time, &e.PriceFrom, &e.Rating, &e.Accent, &e.Blurb)
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
		SELECT b.id, e.title, e.venue, e.event_date, e.event_time, b.seat_labels, b.total_amount, b.status, b.booking_code
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
		rows.Scan(&b.ID, &b.EventTitle, &b.Venue, &b.Date, &b.Time, &seatLabelsJSON, &b.Total, &b.Status, &b.Code)
		json.Unmarshal(seatLabelsJSON, &b.SeatLabels)
		bookings = append(bookings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     int      `json:"userId"`
		EventID    int      `json:"eventId"`
		SeatLabels []string `json:"seatLabels"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if len(req.SeatLabels) == 0 {
		http.Error(w, `{"error": "No seats selected"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		http.Error(w, `{"error": "Failed to start transaction"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx) // Safe to call even after Commit

	var totalAmount float64
	for _, seatLabel := range req.SeatLabels {
		var price float64
		// Check if seat exists and is available
		err = tx.QueryRow(ctx, `SELECT price FROM seats WHERE event_id = $1 AND seat_label = $2 AND status = 'available'`, req.EventID, seatLabel).Scan(&price)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Seat %s is not available or does not exist"}`, seatLabel), http.StatusConflict)
			return
		}
		totalAmount += price

		// Update seat status
		_, err = tx.Exec(ctx, `UPDATE seats SET status = 'booked' WHERE event_id = $1 AND seat_label = $2`, req.EventID, seatLabel)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to update seat %s"}`, seatLabel), http.StatusInternalServerError)
			return
		}
	}

	bookingID := fmt.Sprintf("bk_%d", time.Now().Unix())
	bookingCode := fmt.Sprintf("BK-%d-%s", req.EventID, strconv.FormatInt(time.Now().UnixNano()%10000, 10))
	seatLabelsJSON, _ := json.Marshal(req.SeatLabels)

	_, err = tx.Exec(ctx, `
		INSERT INTO bookings (id, user_id, event_id, seat_labels, total_amount, status, booking_code, created_at)
		VALUES ($1, $2, $3, $4, $5, 'completed', $6, NOW())`,
		bookingID, req.UserID, req.EventID, seatLabelsJSON, totalAmount, bookingCode)
	if err != nil {
		http.Error(w, `{"error": "Failed to create booking record"}`, http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(ctx); err != nil {
		http.Error(w, `{"error": "Failed to commit transaction"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"bookingId":   bookingID,
		"bookingCode": bookingCode,
		"total":       totalAmount,
	})
}

func (h *Handler) GetCategoryMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.CategoryPricing)
}