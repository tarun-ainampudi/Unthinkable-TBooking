# Unthinkable-TBooking — Frontend

React single-page app for the event ticketing platform. Users browse events,
select seats on an interactive seat map with a live hold countdown, pay via
card/UPI, and get a confirmation with a QR ticket. Profile and booking history
are fetched from the Go backend API.

## Tech Stack

- **React 18** + **Vite 5** (dev server & build tool)
- **lucide-react** — icon set
- Plain fetch wrapper (`src/api.js`) — no HTTP client library
- CSS-in-JS design system (`src/styles/GlobalStyles.jsx`, ticket-stub theme)

## Getting Started

```bash
npm install    # first time only
npm run dev    # start dev server (default: http://localhost:5173)
```

Other scripts:

```bash
npm run build     # production build to dist/
npm run preview   # serve the production build locally
```

### Connecting to the backend

The app expects the Go backend running on **http://localhost:8080**. To point
it elsewhere, set `VITE_API_URL` in a `.env` file in this folder:

```
VITE_API_URL=http://localhost:9090
```

## Application Flow

1. **Events** — browse, search, and filter events by category.
2. **Seat Selection** — pick seats; held seats are reserved for 8 minutes with
   a visible countdown and running total.
3. **Payment** — card or UPI form with order summary.
4. **Confirmation** — booking code + QR ticket.
5. **Bookings / Profile** — history with status filters and editable profile.

Login is available from the top nav (demo user: `aarav.shah@example.com` /
`ticket123`). Holding seats and booking require being logged in.

## Project Structure

```
src/
├── main.jsx                    React entry point
├── App.jsx                     Top-level state + page routing
├── api.js                      Fetch wrapper (API_BASE, JSON handling, errors)
├── utils.js                    formatDate / formatINR / pad helpers
│
├── styles/
│   └── GlobalStyles.jsx        Full CSS design system (ticket-stub theme)
│
├── components/
│   ├── TopNav.jsx              Header nav (Events / My Bookings / Profile)
│   ├── EventCard.jsx           Card used on the event selection grid
│   ├── SeatButton.jsx          Single seat in the seat map
│   ├── BookingRow.jsx          Row used on the bookings page
│   └── QRPlaceholder.jsx       Decorative QR grid (NOT a real scannable QR)
│
└── pages/
    ├── EventSelectionPage.jsx  Browse + search + filter events
    ├── SeatSelectionPage.jsx   Seat map, hold countdown, running total
    ├── PaymentPage.jsx         Card/UPI form + order summary
    ├── ConfirmationPage.jsx    Post-payment confirmation + QR
    ├── BookingsPage.jsx        Booking history with status filters
    └── ProfilePage.jsx         User details, editable
```
