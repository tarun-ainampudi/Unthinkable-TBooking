package tests

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"unthinkable-tbooking-backend/database"
	"unthinkable-tbooking-backend/handlers"
)

var testDB *pgxpool.Pool

// TestMain sets up the test environment before any tests run
func TestMain(m *testing.M) {
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://tarun:postgres@localhost:5432/ticketdb?sslmode=disable")
		log.Println("[Info] DATABASE_URL not set for tests; using local fallback connection string.")
	}

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

func loginTestUser(router *chi.Mux) string {
	body := `{"email":"aarav.shah@example.com","password":"ticket123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		panic("failed to create authenticated test session: " + rr.Body.String())
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		panic("failed to parse login response: " + err.Error())
	}
	return payload["sessionToken"].(string)
}

// setupTestRouter creates a chi router with all routes mapped to our handlers
func setupTestRouter() *chi.Mux {
	r := chi.NewRouter()
	h := handlers.NewHandler(testDB)

	r.Get("/api/events", h.GetEvents)
	r.Get("/api/events/{id}/seats", h.GetEventSeats)
	r.Get("/api/pricing", h.GetCategoryMeta)
	r.Get("/api/session", h.GetSession)
	r.Get("/api/users/{id}", h.GetUser)
	r.Get("/api/users/{id}/bookings", h.GetUserBookings)
	r.Post("/api/login", h.Login)
	r.Post("/api/logout", h.Logout)
	r.Post("/api/seats/hold", h.HoldSeats)
	r.Delete("/api/seats/hold", h.ReleaseHeldSeats)
	r.Post("/api/bookings", h.CreateBooking)

	return r
}