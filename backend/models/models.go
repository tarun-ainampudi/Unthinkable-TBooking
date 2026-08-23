package models

type Event struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Category  string  `json:"category"`
	Tag       string  `json:"tag"`
	Venue     string  `json:"venue"`
	Date      string  `json:"date"`
	Time      string  `json:"time"`
	PriceFrom float64 `json:"priceFrom"`
	Rating    float64 `json:"rating"`
	Accent    string  `json:"accent"`
	Blurb     string  `json:"blurb"`
}

type Seat struct {
	ID       string  `json:"id"`
	Row      string  `json:"row"`
	Col      int     `json:"col"`
	Category string  `json:"category"`
	Status   string  `json:"status"`
	Price    float64 `json:"price"`
}

type Booking struct {
	ID             string   `json:"id"`
	EventTitle     string   `json:"eventTitle"`
	Venue          string   `json:"venue"`
	Date           string   `json:"date"`
	Time           string   `json:"time"`
	SeatLabels     []string `json:"seatLabels"`
	Total          float64  `json:"total"`
	Status         string   `json:"status"`
	Code           string   `json:"code"`
	QRCodeDataUrl  string   `json:"qrCodeDataUrl,omitempty"`
}

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	MemberSince string `json:"memberSince"`
}

type CategoryMeta struct {
	Label string  `json:"label"`
	Price float64 `json:"price"`
	Color string  `json:"color"`
}

// CategoryPricing holds the global pricing rules
var CategoryPricing = map[string]CategoryMeta{
	"VIP":      {Label: "VIP", Price: 2499, Color: "var(--gold)"},
	"Premium":  {Label: "Premium", Price: 1499, Color: "var(--teal)"},
	"Standard": {Label: "Standard", Price: 899, Color: "var(--paper-dim)"},
}