package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var user map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if user["name"] != "Aarav Shah" {
		t.Errorf("Expected user name 'Aarav Shah', got %v", user["name"])
	}
}

func TestGetUserNotFound(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/users/999", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code for non-existent user: got %v want %v", status, http.StatusNotFound)
	}
}

func TestGetUserBookings(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/users/1/bookings", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var bookings []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &bookings); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(bookings) < 1 {
		t.Errorf("Expected at least 1 booking, got %d", len(bookings))
	}
}