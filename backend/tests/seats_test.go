package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEventSeats(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/events/1/seats", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var seats []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &seats); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// 8 rows * 12 cols = 96 seats
	if len(seats) != 96 {
		t.Errorf("Expected 96 seats, got %d", len(seats))
	}
}