
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5/middleware"
	"unthinkable-tbooking-backend/database"
	"unthinkable-tbooking-backend/handlers"
)

func main() {

	// Load environment variables from .env file when present.
	if err := godotenv.Load(); err != nil {
		log.Println("[Info] No .env file found; using environment variables or defaults.")
	}

	// Initialize Database
	db := database.InitDB()
	defer db.Close()

	// Seed Database with mock data
	database.SeedDB(db)

	// Initialize Handlers with DB dependency
	h := handlers.NewHandler(db)

	// Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Register Routes
	r.Get("/api/events", h.GetEvents)
	r.Get("/api/events/{id}/seats", h.GetEventSeats)
	r.Get("/api/pricing", h.GetCategoryMeta)
	r.Get("/api/users/{id}", h.GetUser)
	r.Get("/api/users/{id}/bookings", h.GetUserBookings)
	r.Get("/api/session", h.GetSession)
	r.Post("/api/login", h.Login)
	r.Post("/api/logout", h.Logout)
	r.Post("/api/seats/hold", h.HoldSeats)
	r.Delete("/api/seats/hold", h.ReleaseHeldSeats)
	r.Post("/api/bookings", h.CreateBooking)

	// Start Server
	log.Println("[Info] Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}