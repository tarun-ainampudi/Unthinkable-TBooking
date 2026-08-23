# Unthinkable — Event Ticketing Frontend

A React + Vite prototype covering: event selection, seat selection with a
live hold timer, payment, booking confirmation, user profile, and booking
history. All data is mocked in `src/data.js` — swap it for real API calls
to your Go backend when ready.

## Run it

```bash
npm install
npm run dev
```

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