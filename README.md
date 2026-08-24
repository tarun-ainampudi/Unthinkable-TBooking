# Unthinkable-TBooking

A full-stack event ticketing application. Users can browse events (concerts,
theatre, comedy, dance), pick seats on a live seat map with a hold timer, pay
with card/UPI, and receive a booking confirmation with a QR ticket. Bookings
and profile details are persisted in PostgreSQL.

## Tech Stack

| Layer    | Technology                                                        |
|----------|-------------------------------------------------------------------|
| Frontend | React 18, Vite, lucide-react                                      |
| Backend  | Go 1.21, chi (router), pgx v5 (PostgreSQL driver), go-qrcode      |
| Database | PostgreSQL (auto-created and seeded on first run)                 |

## Project Structure

```
Unthinkable-TBooking/
├── frontend/     React SPA (event browse → seats → payment → confirmation)
├── backend/      Go REST API + database bootstrap
└── data/         Local PostgreSQL data directory (created by initdb)
```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL installed locally (default: `localhost:5432`)

### 0. Initialize and start PostgreSQL

The repository keeps its database cluster in the local `data/` folder. Using
`pg_ctl` (adjust the path to where PostgreSQL is installed, e.g. `C:\pgsql\bin`
on Windows or `/usr/lib/postgresql/<version>/bin` on Linux):

```bash
# One-time: initialize a new database cluster in ./data
pg_ctl initdb -D data

# Start the server (logs to postgres.log)
pg_ctl start -D data -l postgres.log
```

To check the status and stop it:

```bash
pg_ctl status -D data     # check if the server is running
pg_ctl stop -D data       # stop the server
```

### 1. Start the backend

```bash
cd backend
go mod tidy        # first time only
go run .
```

The server listens on **http://localhost:8080**. On startup it will:

1. Read `DATABASE_URL` from `.env` (falls back to a local default).
2. Create the database if it does not exist.
3. Create all tables, indexes, and seed mock events/seats/bookings.

### 2. Start the frontend

```bash
cd frontend
npm install        # first time only
npm run dev
```

Open the URL Vite prints (default **http://localhost:5173**).

## Demo Login

| Email                     | Password    |
|---------------------------|-------------|
| `aarav.shah@example.com`  | `ticket123` |

## Configuration

| Location            | Variable       | Purpose                                    |
|---------------------|----------------|--------------------------------------------|
| `backend/.env`      | `DATABASE_URL` | PostgreSQL connection string               |
| `frontend/.env`     | `VITE_API_URL` | Backend base URL (defaults to port 8080)   |

## Documentation

- [`backend/README.md`](backend/README.md) — API reference, architecture, tests
- [`frontend/README.md`](frontend/README.md) — UI flow, components, configuration
