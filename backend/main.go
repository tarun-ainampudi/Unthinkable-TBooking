
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

	// Load environment variables from .env file
	err := godotenv.Load()
    if err != nil {
        log.Fatal("[Error] Error loading .env file")
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

	// Register Routes
	r.Get("/api/events", h.GetEvents)
	r.Get("/api/events/{id}/seats", h.GetEventSeats)
	r.Get("/api/pricing", h.GetCategoryMeta)
	r.Get("/api/users/{id}", h.GetUser)
	r.Get("/api/users/{id}/bookings", h.GetUserBookings)
	r.Post("/api/bookings", h.CreateBooking)

	// Start Server
	log.Println("[Info] Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}