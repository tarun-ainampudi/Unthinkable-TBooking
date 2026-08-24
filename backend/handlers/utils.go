package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/skip2/go-qrcode"
)

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
