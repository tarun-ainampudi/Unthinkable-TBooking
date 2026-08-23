package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserBookingsIncludeQrCode(t *testing.T) {
	router := setupTestRouter()

	bookingReq := httptest.NewRequest(http.MethodPost, "/api/bookings", strings.NewReader(`{"userId":1,"eventId":1,"seatLabels":["A5","A6"]}`))
	bookingReq.Header.Set("Content-Type", "application/json")
	bookingRes := httptest.NewRecorder()
	router.ServeHTTP(bookingRes, bookingReq)

	if status := bookingRes.Code; status != http.StatusCreated {
		t.Fatalf("booking request returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	bookingsReq := httptest.NewRequest(http.MethodGet, "/api/users/1/bookings", nil)
	bookingsRes := httptest.NewRecorder()
	router.ServeHTTP(bookingsRes, bookingsReq)

	if status := bookingsRes.Code; status != http.StatusOK {
		t.Fatalf("bookings request returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var bookings []map[string]interface{}
	if err := json.Unmarshal(bookingsRes.Body.Bytes(), &bookings); err != nil {
		t.Fatalf("Failed to parse bookings response: %v", err)
	}

	if len(bookings) == 0 {
		t.Fatal("Expected at least one booking in the response")
	}

	if _, exists := bookings[0]["qrCodeDataUrl"]; !exists {
		t.Fatal("Expected the user bookings payload to include a QR code image URL")
	}
}
