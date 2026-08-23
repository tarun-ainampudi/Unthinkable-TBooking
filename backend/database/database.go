package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"unthinkable-tbooking-backend/models"
)

func InitDB() *pgxpool.Pool {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("[Error] DATABASE_URL environment variable is missing")
	}

	ctx := context.Background()

	// First, verify that the target database exists and is reachable.
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("[Error] Invalid DATABASE_URL: %v", err)
	}

	targetDB := cfg.ConnConfig.Database

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err == nil {
		// pgxpool.New doesn't guarantee an actual connection,
		// so explicitly ping the database.
		err = db.Ping(ctx)
	}

	if err == nil {
		log.Println("[Info] Connected to PostgreSQL successfully!")
		return db
	}

	db.Close()

	// If the database doesn't exist, connect to the default
	// PostgreSQL "postgres" database.
	maintenanceCfg := *cfg
	maintenanceCfg.ConnConfig.Database = "postgres"

	maintenanceDB, err := pgxpool.NewWithConfig(ctx, &maintenanceCfg)
	if err != nil {
		log.Fatalf("[Error] Unable to connect to PostgreSQL server: %v", err)
	}

	if err := maintenanceDB.Ping(ctx); err != nil {
		maintenanceDB.Close()
		log.Fatalf("[Error] PostgreSQL server is not reachable: %v", err)
	}

	// Create the target database.
	_, err = maintenanceDB.Exec(
		ctx,
		fmt.Sprintf(`CREATE DATABASE "%s"`, targetDB),
	)

	if err != nil {
		maintenanceDB.Close()
		log.Fatalf("[Error] Unable to create database %q: %v", targetDB, err)
	}

	log.Printf("[Info] Database %q created successfully!", targetDB)

	maintenanceDB.Close()

	// Now connect to the newly created database.
	db, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("[Error] Unable to connect to newly created database: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		log.Fatalf("[Error] Unable to ping newly created database: %v", err)
	}

	log.Println("[Info] Connected to PostgreSQL successfully!")

	return db
}

func SeedDB(db *pgxpool.Pool) {
	ctx := context.Background()

	// 1. Seed User
	_, _ = db.Exec(ctx, `
		INSERT INTO users (id, name, email, phone, member_since) 
		VALUES (1, 'Aarav Shah', 'aarav.shah@example.com', '+91 98765 43210', 'March 2025')
		ON CONFLICT (id) DO NOTHING`)

	// 2. Seed Events
	events := []models.Event{
		{ID: 1, Title: "Midnight Frequencies", Category: "Concert", Tag: "Electronic", Venue: "Marquee Arena, Amaravati", Date: "2026-09-12", Time: "8:00 PM", PriceFrom: 899, Rating: 4.8, Accent: "gold", Blurb: "A night-long set from three touring electronic acts, built for a room that likes its bass loud."},
		{ID: 2, Title: "The Glass Menagerie", Category: "Theatre", Tag: "Drama", Venue: "Riverside Playhouse", Date: "2026-09-18", Time: "7:00 PM", PriceFrom: 499, Rating: 4.6, Accent: "teal", Blurb: "Tennessee Williams' family drama, staged in the round by the Riverside repertory company."},
		{ID: 3, Title: "Laugh Riot Live", Category: "Comedy", Tag: "Stand-up", Venue: "The Attic Club", Date: "2026-08-29", Time: "9:00 PM", PriceFrom: 349, Rating: 4.9, Accent: "red", Blurb: "Four comics, one mic, zero warning about what happens if you sit in the front row."},
		{ID: 4, Title: "Nocturne: A Jazz Evening", Category: "Concert", Tag: "Jazz", Venue: "Blue Room Hall", Date: "2026-09-05", Time: "8:30 PM", PriceFrom: 699, Rating: 4.7, Accent: "teal", Blurb: "A quartet residency closing out its run with a set of standards and one new commission."},
		{ID: 5, Title: "Selene — Dance Reimagined", Category: "Dance", Tag: "Contemporary", Venue: "Marquee Arena, Amaravati", Date: "2026-10-02", Time: "6:30 PM", PriceFrom: 799, Rating: 4.5, Accent: "gold", Blurb: "A contemporary dance piece built around the phases of the moon, performed with a live score."},
		{ID: 6, Title: "Open Mic Sundays", Category: "Comedy", Tag: "Variety", Venue: "The Attic Club", Date: "2026-08-31", Time: "6:00 PM", PriceFrom: 149, Rating: 4.2, Accent: "red", Blurb: "Ten new acts, unfiltered. The weekly night where the club's future headliners get found."},
	}

	for _, ev := range events {
		_, _ = db.Exec(ctx, `
			INSERT INTO events (id, title, category, tag, venue, event_date, event_time, price_from, rating, accent, blurb)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (id) DO NOTHING`,
			ev.ID, ev.Title, ev.Category, ev.Tag, ev.Venue, ev.Date, ev.Time, ev.PriceFrom, ev.Rating, ev.Accent, ev.Blurb)
	}

	// 3. Seed Seats
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	cols := 12
	bookedSeats := map[string]bool{"A3": true, "A4": true, "B7": true, "B8": true, "C1": true, "C12": true, "D5": true, "D6": true, "E9": true, "F2": true, "F3": true, "G10": true, "G11": true, "H4": true, "H5": true, "H6": true, "A9": true, "B2": true}
	heldSeats := map[string]bool{"C5": true, "C6": true, "D9": true, "F8": true}

	for _, ev := range events {
		for _, row := range rows {
			for c := 1; c <= cols; c++ {
				seatLabel := fmt.Sprintf("%s%d", row, c)
				category := "Standard"
				price := 899.0
				if row == "A" || row == "B" {
					category = "VIP"
					price = 2499.0
				} else if row == "C" || row == "D" || row == "E" {
					category = "Premium"
					price = 1499.0
				}

				status := "available"
				if ev.ID == 1 { // Apply mock booked/held seats to Event 1 for demonstration
					if bookedSeats[seatLabel] {
						status = "booked"
					} else if heldSeats[seatLabel] {
						status = "held"
					}
				}

				_, _ = db.Exec(ctx, `
					INSERT INTO seats (event_id, seat_label, row_name, col_num, category, status, price)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					ON CONFLICT (event_id, seat_label) DO NOTHING`,
					ev.ID, seatLabel, row, c, category, status, price)
			}
		}
	}

	// 4. Seed Past Booking
	_, _ = db.Exec(ctx, `
		INSERT INTO bookings (id, user_id, event_id, seat_labels, total_amount, status, booking_code, created_at)
		VALUES ('bk_9021', 1, 1, '["D5", "D6"]', 2998, 'completed', 'WL-9021-XQ', '2026-06-14 20:00:00')
		ON CONFLICT (id) DO NOTHING`)

	log.Println("[Info] Database seeded with mock data successfully!")
}