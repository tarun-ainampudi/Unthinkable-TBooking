package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureSchema(ctx context.Context, db *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			phone TEXT,
			member_since TEXT,
			password TEXT NOT NULL DEFAULT 'ticket123'
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS password TEXT NOT NULL DEFAULT 'ticket123'`,
		`CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			category TEXT NOT NULL,
			tag TEXT,
			venue TEXT NOT NULL,
			event_date TEXT NOT NULL,
			event_time TEXT NOT NULL,
			price_from NUMERIC(10,2) NOT NULL DEFAULT 0,
			rating NUMERIC(3,1) NOT NULL DEFAULT 0,
			accent TEXT,
			blurb TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS seats (
			event_id INTEGER NOT NULL,
			seat_label TEXT NOT NULL,
			row_name TEXT NOT NULL,
			col_num INTEGER NOT NULL,
			category TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'available',
			price NUMERIC(10,2) NOT NULL DEFAULT 0,
			held_until TIMESTAMPTZ,
			PRIMARY KEY (event_id, seat_label),
			CONSTRAINT fk_seats_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			event_id INTEGER NOT NULL,
			seat_labels JSONB NOT NULL,
			total_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'completed',
			booking_code TEXT NOT NULL,
			qr_code_data_url TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_bookings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_bookings_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMPTZ NOT NULL,
			revoked BOOLEAN NOT NULL DEFAULT FALSE,
			CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`ALTER TABLE bookings ADD COLUMN IF NOT EXISTS qr_code_data_url TEXT`,
		`CREATE INDEX IF NOT EXISTS idx_events_category ON events(category)`,
		`CREATE INDEX IF NOT EXISTS idx_seats_event_status ON seats(event_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_user_event ON bookings(user_id, event_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id)`,
		`ALTER TABLE seats ADD COLUMN IF NOT EXISTS held_until TIMESTAMPTZ`,
	}

	for _, query := range queries {
		if _, err := db.Exec(ctx, query); err != nil {
			return err
		}
	}

	return nil
}
