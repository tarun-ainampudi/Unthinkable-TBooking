# Unthinkable-TBooking — Backend

Go REST API powering the event ticketing app. It serves events and seat maps,
manages seat holds with expiry, creates bookings with QR ticket codes, and
handles session-based authentication. On first run it bootstraps the entire
PostgreSQL database (create → schema → seed) automatically.

## Tech Stack

- **Go 1.21**
- [chi](https://github.com/go-chi/chi) v5 — HTTP router + middleware
- [pgx](https://github.com/jackc/pgx) v5 — PostgreSQL driver/connection pool
- [godotenv](https://github.com/joho/godotenv) — loads `.env` on startup
- [go-qrcode](https://github.com/skip2/go-qrcode) — generates booking QR codes

## Getting Started

```bash
# 1. Configure the database connection (optional — has a local fallback)
echo 'DATABASE_URL="postgres://username:password@localhost:5432/dbname?sslmode=disable"' > .env

# 2. Install dependencies and run
go mod tidy
go run .
```

The API starts on **http://localhost:8080**.

### Startup behavior

1. Loads `.env` if present; otherwise uses environment variables/defaults.
2. Connects to PostgreSQL; if the target database doesn't exist it connects to
   the `postgres` maintenance database and runs `CREATE DATABASE`.
3. Applies the schema (`CREATE TABLE IF NOT EXISTS` … idempotent).
4. Seeds mock data: a demo user, six events, an 8×12 seat map per event
   (with some seats pre-booked/held), and one past booking.

## Project Structure

Each folder/file groups code by responsibility:

```
backend/
├── main.go                  App entry point: env, DB init, routes, CORS, server
├── .env                     DATABASE_URL (not committed in real deployments)
├── go.mod / go.sum          Go module definition and checksums
│
├── handlers/                HTTP layer — one file per domain
│   ├── handler.go           Handler struct (DB pool dependency) + constructor
│   ├── auth.go              Session auth middleware, login/logout/get-session
│   ├── events.go            List events, pricing category metadata
│   ├── seats.go             Seat map for an event, hold & release seats
│   ├── bookings.go          Create booking (QR + email sim), user bookings list
│   ├── users.go             Fetch a single user profile
│   └── utils.go             QR generation, email simulation, session tokens
│
├── database/                Persistence layer
│   ├── database.go          InitDB: connect, create DB if missing
│   ├── schema.go            EnsureSchema: tables, constraints, indexes
│   └── seed.go              SeedDB: demo users/events/seats/bookings
│
├── models/                  Shared data structures
│   └── models.go            Event, Seat, Booking, User, CategoryPricing
│
└── tests/                   Integration tests (hit a real local PostgreSQL)
    ├── setup_test.go        Test bootstrap/router wiring
    ├── event_test.go        Events + pricing endpoints
    ├── seats_test.go        Seat map endpoint
    ├── bookings_test.go     Bookings endpoints incl. QR data
    └── user_test.go         User profile endpoints
```

## API Reference

All endpoints are prefixed with `/api`. Authenticated endpoints expect an
`Authorization: Bearer <sessionToken>` header.

| Method   | Path                      | Auth | Description                                        |
|----------|---------------------------|------|----------------------------------------------------|
| `GET`    | `/api/events`             | No   | All events                                         |
| `GET`    | `/api/events/{id}/seats`  | No   | Seat map for an event                              |
| `GET`    | `/api/pricing`            | No   | Seat category pricing metadata                     |
| `POST`   | `/api/login`              | No   | Login → `{ sessionToken, expiresAt, user }`        |
| `GET`    | `/api/session`            | Yes  | Current session's user + refreshed expiry          |
| `POST`   | `/api/logout`             | Yes  | Revoke the current session                         |
| `GET`    | `/api/users/{id}`         | No   | User profile                                       |
| `GET`    | `/api/users/{id}/bookings`| No   | Booking history for a user                         |
| `POST`   | `/api/seats/hold`         | Yes  | Hold seats for 8 minutes (`{eventId, seatLabels}`) |
| `DELETE` | `/api/seats/hold`         | Yes  | Release held seats                                 |
| `POST`   | `/api/bookings`           | Yes  | Confirm booking → booking code + QR data URL       |

### Domain rules

- **Seat holds** expire after **8 minutes**; expired holds become available again.
- **Sessions** expire after **15 minutes** and are refreshed on every
  authenticated request.
- **Bookings** are created transactionally: every selected seat is re-validated
  before being marked `booked`; a QR code is generated and a ticket email is
  simulated.

## Testing

Tests require a reachable local PostgreSQL server.

```bash
go test ./tests -v
```

## Configuration

| Variable       | Default                          | Purpose                        |
|----------------|----------------------------------|--------------------------------|
| `DATABASE_URL` | `postgres://username:password@localhost:5432/dbname?sslmode=disable` | PostgreSQL connection string |
