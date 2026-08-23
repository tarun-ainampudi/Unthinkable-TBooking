import React, { useEffect, useState } from "react";
import GlobalStyles from "./styles/GlobalStyles";
import TopNav from "./components/TopNav";
import EventSelectionPage from "./pages/EventSelectionPage";
import SeatSelectionPage from "./pages/SeatSelectionPage";
import PaymentPage from "./pages/PaymentPage";
import ConfirmationPage from "./pages/ConfirmationPage";
import ProfilePage from "./pages/ProfilePage";
import BookingsPage from "./pages/BookingsPage";
import { apiFetch } from "./api";

export default function App() {
  const [page, setPage] = useState("events");
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [selectedSeats, setSelectedSeats] = useState([]);
  const [seats, setSeats] = useState([]);
  const [bookings, setBookings] = useState([]);
  const [lastBooking, setLastBooking] = useState(null);
  const [user, setUser] = useState(null);
  const [bookingError, setBookingError] = useState("");

  useEffect(() => {
    const loadUser = async () => {
      try {
        const profile = await apiFetch("/api/users/1");
        setUser(profile);
      } catch (error) {
        console.error("Failed to load user profile", error);
      }
    };

    const loadBookings = async () => {
      try {
        const items = await apiFetch("/api/users/1/bookings");
        setBookings(items || []);
      } catch (error) {
        console.error("Failed to load bookings", error);
      }
    };

    loadUser();
    loadBookings();
  }, []);

  async function loadSeatsForEvent(event) {
    try {
      const nextSeats = await apiFetch(`/api/events/${event.id}/seats`);
      const normalized = (nextSeats || []).map((seat) => ({
        id: seat.id || seat.seat_label || `${seat.row}${seat.col}`,
        row: seat.row || seat.row_name || seat.id?.slice(0, 1),
        col: typeof seat.col === "number" ? seat.col : Number(seat.col_num || seat.col || 1),
        category: seat.category || "Standard",
        status: seat.status || "available",
        price: Number(seat.price || 0),
      }));
      setSeats(normalized);
    } catch (error) {
      console.error("Failed to load seat map", error);
      setSeats([]);
    }
  }

  async function goToSeats(event) {
    setSelectedEvent(event);
    setSelectedSeats([]);
    setBookingError("");
    setSeats([]);
    await loadSeatsForEvent(event);
    setPage("seats");
  }

  function goToPayment() {
    setBookingError("");
    setPage("payment");
  }

  async function confirmBooking() {
    const seatLabels = selectedSeats.slice().sort();
    if (!selectedEvent || seatLabels.length === 0) return;

    try {
      setBookingError("");
      const response = await apiFetch("/api/bookings", {
        method: "POST",
        body: JSON.stringify({
          userId: 1,
          eventId: Number(selectedEvent.id),
          seatLabels,
        }),
      });

      const booking = {
        id: response.bookingId || `bk_${Date.now()}`,
        eventTitle: selectedEvent.title,
        venue: selectedEvent.venue,
        date: selectedEvent.date,
        time: selectedEvent.time,
        seatLabels,
        total: Number(response.total || 0),
        status: "upcoming",
        code: response.bookingCode || "BK-BOOKED",
      };

      setBookings((prev) => [booking, ...prev]);
      setLastBooking(booking);
      setSelectedSeats([]);
      setPage("confirmation");
    } catch (error) {
      setBookingError(error.message || "Booking failed. Please try again.");
      setPage("payment");
    }
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
            bookingError={bookingError}
          />
        )}

        {page === "confirmation" && lastBooking && (
          <ConfirmationPage booking={lastBooking} onDone={() => setPage("bookings")} />
        )}

        {page === "profile" && <ProfilePage user={user} />}

        {page === "bookings" && <BookingsPage bookings={bookings} />}
      </main>
    </div>
  );
}
