# Stubline — Event Ticketing Frontend (prototype)

A React + Vite prototype covering: event selection, seat selection with a
live hold timer, payment, booking confirmation, user profile, and booking
history. All data is mocked in `src/data.js` — swap it for real API calls
to your Go backend when ready.

## Run it

```bash
npm install
npm run dev
```

Then open the printed local URL (usually http://localhost:5173).

## Project structure

```
src/
  data.js                    Mock events, seats, bookings, user (replace with API calls)
  utils.js                   formatDate / formatINR / pad helpers
  App.jsx                    Top-level state + page routing
  main.jsx                   React entry point
  styles/
    GlobalStyles.jsx         Full CSS design system (ticket-stub theme)
  components/
    TopNav.jsx                Header nav (Events / My Bookings / Profile)
    EventCard.jsx              Card used on the event selection grid
    SeatButton.jsx             Single seat in the seat map
    QRPlaceholder.jsx          Decorative QR grid (NOT a real scannable QR)
    BookingRow.jsx              Row used on the bookings page
  pages/
    EventSelectionPage.jsx     Browse + search + filter events
    SeatSelectionPage.jsx      Seat map, hold countdown, running total
    PaymentPage.jsx            Card/UPI form + order summary
    ConfirmationPage.jsx       Post-payment confirmation + QR
    ProfilePage.jsx            User details, editable
    BookingsPage.jsx           Booking history with status filters
```

## Wiring to a real backend

- **Events / seats**: replace `EVENTS` and `buildSeats()` in `data.js`
  with fetches to your API; seat status should come from the DB, not be
  generated client-side.
- **Hold timer**: `SeatSelectionPage.jsx` currently counts down locally.
  Sync it to the server-issued hold expiry (e.g. return `expires_at` from
  your hold endpoint and compute remaining time from that, ideally
  re-validated with a periodic poll or websocket).
- **Payment**: `PaymentPage.jsx` never sends data anywhere — plug in your
  real payment provider and call your booking-confirmation endpoint on
  success instead of the local `onPay()` callback.
- **QR codes**: `QRPlaceholder.jsx` is a decorative stand-in. Generate the
  real QR server-side (e.g. `go-qrcode`) and return it as an image URL or
  base64 payload to render instead.
- **Auth / roles**: there's no auth here — add a login flow and gate
  organiser/admin views separately once your role-based auth is in place.

## Notes

- No router library is used; page switching is local `useState` in
  `App.jsx`. Swap in React Router if you want real URLs / back-button
  support.
- Styling avoids Tailwind and uses hand-written CSS (see
  `GlobalStyles.jsx`) so it renders correctly without a build-time
  Tailwind compiler.
