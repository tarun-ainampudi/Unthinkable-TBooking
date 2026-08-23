package tests

import (
	"context"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"unthinkable-tbooking-backend/database"
	"unthinkable-tbooking-backend/handlers"
)

var testDB *pgxpool.Pool

// TestMain sets up the test environment before any tests run
func TestMain(m *testing.M) {
	// 1. Connect to a dedicated test database
	testDB = database.InitDB()

	ctx := context.Background()

	// 2. Clean the database before tests to ensure a consistent state
	_, err := testDB.Exec(ctx, `TRUNCATE TABLE bookings, seats, events, users RESTART IDENTITY CASCADE`)
	if err != nil {
		// If tables don't exist yet, ignore the error (they will be created by SeedDB or migrations)
	}

	// 3. Seed the database with mock data
	database.SeedDB(testDB)

	// 4. Run all tests
	code := m.Run()

	// 5. Cleanup
	testDB.Close()
	os.Exit(code)
}

// setupTestRouter creates a chi router with all routes mapped to our handlers
func setupTestRouter() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewHandler(testDB)

	r.Get("/api/events", h.GetEvents)
	r.Get("/api/events/{id}/seats", h.GetEventSeats)
	r.Get("/api/pricing", h.GetCategoryMeta)
	r.Get("/api/users/{id}", h.GetUser)
	r.Get("/api/users/{id}/bookings", h.GetUserBookings)
	r.Post("/api/bookings", h.CreateBooking)

	return r
}