import React, { useState } from "react";
import GlobalStyles from "./styles/GlobalStyles";
import TopNav from "./components/TopNav";
import EventSelectionPage from "./pages/EventSelectionPage";
import SeatSelectionPage from "./pages/SeatSelectionPage";
import PaymentPage from "./pages/PaymentPage";
import ConfirmationPage from "./pages/ConfirmationPage";
import ProfilePage from "./pages/ProfilePage";
import BookingsPage from "./pages/BookingsPage";
import { buildSeats, CATEGORY_META, PAST_BOOKINGS } from "./data";

export default function App() {
  const [page, setPage] = useState("events");
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [selectedSeats, setSelectedSeats] = useState([]);
  const [seats, setSeats] = useState(buildSeats);
  const [bookings, setBookings] = useState(PAST_BOOKINGS);
  const [lastBooking, setLastBooking] = useState(null);

  function goToSeats(event) {
    setSelectedEvent(event);
    setSelectedSeats([]);
    setSeats(buildSeats());
    setPage("seats");
  }

  function goToPayment() {
    setPage("payment");
  }

  function confirmBooking() {
    const seatLabels = selectedSeats.slice().sort();
    const total = selectedSeats.reduce((sum, id) => {
      const seat = seats.find((s) => s.id === id);
      return sum + CATEGORY_META[seat.category].price;
    }, 0);
    const code = `${selectedEvent.title.slice(0, 2).toUpperCase()}-${Math.floor(1000 + Math.random() * 9000)}-${
      "ABCDEFGHJKLMNP"[Math.floor(Math.random() * 14)]
    }${"ABCDEFGHJKLMNP"[Math.floor(Math.random() * 14)]}`;
    const booking = {
      id: `bk_${Date.now()}`,
      eventTitle: selectedEvent.title,
      venue: selectedEvent.venue,
      date: selectedEvent.date,
      time: selectedEvent.time,
      seatLabels,
      total,
      status: "upcoming",
      code,
    };
    setBookings((prev) => [booking, ...prev]);
    setLastBooking(booking);
    setPage("confirmation");
  }

  return (
    <div className="app-root">
      <GlobalStyles />
      <TopNav page={page} setPage={setPage} />
      <main className="app-main">
        {page === "events" && <EventSelectionPage onSelect={goToSeats} />}

        {page === "seats" && selectedEvent && (
          <SeatSelectionPage
            event={selectedEvent}
            seats={seats}
            selectedSeats={selectedSeats}
            setSelectedSeats={setSelectedSeats}
            onBack={() => setPage("events")}
            onContinue={goToPayment}
          />
        )}

        {page === "payment" && selectedEvent && (
          <PaymentPage
            event={selectedEvent}
            seats={seats}
            selectedSeats={selectedSeats}
            onBack={() => setPage("seats")}
            onPay={confirmBooking}
          />
        )}

        {page === "confirmation" && lastBooking && (
          <ConfirmationPage booking={lastBooking} onDone={() => setPage("bookings")} />
        )}

        {page === "profile" && <ProfilePage />}

        {page === "bookings" && <BookingsPage bookings={bookings} />}
      </main>
    </div>
  );
}
