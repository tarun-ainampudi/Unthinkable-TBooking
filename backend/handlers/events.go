package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"unthinkable-tbooking-backend/models"
)

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

func (h *Handler) GetCategoryMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.CategoryPricing)
}
