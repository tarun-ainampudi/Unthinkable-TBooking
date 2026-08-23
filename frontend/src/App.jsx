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

const HOLD_SECONDS = 8 * 60;
const SESSION_KEY = "ticket_session_token";

export default function App() {
  const [page, setPage] = useState("events");
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [selectedSeats, setSelectedSeats] = useState([]);
  const [seats, setSeats] = useState([]);
  const [bookings, setBookings] = useState([]);
  const [lastBooking, setLastBooking] = useState(null);
  const [user, setUser] = useState(null);
  const [sessionToken, setSessionToken] = useState(() => localStorage.getItem(SESSION_KEY) || "");
  const [bookingError, setBookingError] = useState("");
  const [holdDeadline, setHoldDeadline] = useState(null);

  const requestWithSession = (path, options = {}) => {
    const headers = {
      ...(options.headers || {}),
    };

    if (sessionToken) {
      headers.Authorization = `Bearer ${sessionToken}`;
    }

    return apiFetch(path, {
      ...options,
      headers,
    });
  };

  const loadBookings = async () => {
    if (!user || !sessionToken) {
      setBookings([]);
      return;
    }

    try {
      const items = await requestWithSession(`/api/users/${user.id}/bookings`);
      setBookings(items || []);
    } catch (error) {
      console.error("Failed to load bookings", error);
      setBookings([]);
    }
  };

  useEffect(() => {
    const loadSession = async () => {
      if (!sessionToken) {
        setUser(null);
        setBookings([]);
        return;
      }

      try {
        const session = await requestWithSession("/api/session");
        setUser(session?.user || null);
        if (!session?.user) {
          localStorage.removeItem(SESSION_KEY);
          setSessionToken("");
        }
      } catch (error) {
        localStorage.removeItem(SESSION_KEY);
        setSessionToken("");
        setUser(null);
      }
    };

    loadSession();
  }, [sessionToken]);

  useEffect(() => {
    if (user && sessionToken) {
      loadBookings();
    }
  }, [user, sessionToken]);

  const handleLogin = async ({ email, password }) => {
    try {
      const payload = await apiFetch("/api/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });

      const nextToken = payload.sessionToken;
      localStorage.setItem(SESSION_KEY, nextToken);
      setSessionToken(nextToken);
      setUser(payload.user || null);
      setPage("profile");
      return payload;
    } catch (error) {
      throw error;
    }
  };

  const handleLogout = async () => {
    try {
      if (sessionToken) {
        await requestWithSession("/api/logout", { method: "POST" });
      }
    } catch (error) {
      console.error("Logout failed", error);
    } finally {
      localStorage.removeItem(SESSION_KEY);
      setSessionToken("");
      setUser(null);
      setBookings([]);
      setSelectedSeats([]);
      setHoldDeadline(null);
      setPage("profile");
    }
  };

  const requireLogin = (nextPage = "profile") => {
    if (!sessionToken || !user) {
      setPage("profile");
      return false;
    }
    if (nextPage) setPage(nextPage);
    return true;
  };

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
    if (!requireLogin()) return;
    setSelectedEvent(event);
    setSelectedSeats([]);
    setBookingError("");
    setSeats([]);
    await loadSeatsForEvent(event);
    setPage("seats");
  }

  async function goToPayment() {
    if (!requireLogin()) return;
    setBookingError("");
    if (!selectedEvent || selectedSeats.length === 0) return;

    try {
      const payload = {
        userId: user.id,
        eventId: Number(selectedEvent.id),
        seatLabels: selectedSeats.slice().sort(),
      };
      await requestWithSession("/api/seats/hold", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      setHoldDeadline(Date.now() + HOLD_SECONDS * 1000);
      setPage("payment");
    } catch (error) {
      setBookingError(error.message || "Seats could not be held. Please try again.");
      setPage("seats");
    }
  }

  async function releaseHeldSeats() {
    if (!selectedEvent || selectedSeats.length === 0) return;

    try {
      await requestWithSession("/api/seats/hold", {
        method: "DELETE",
        body: JSON.stringify({
          userId: user?.id || 1,
          eventId: Number(selectedEvent.id),
          seatLabels: selectedSeats.slice().sort(),
        }),
      });
    } catch (error) {
      console.error("Failed to release held seats", error);
    } finally {
      setSelectedSeats([]);
      setHoldDeadline(null);
      await loadSeatsForEvent(selectedEvent);
      setPage("seats");
    }
  }

  async function confirmBooking() {
    const seatLabels = selectedSeats.slice().sort();
    if (!selectedEvent || seatLabels.length === 0) return;

    try {
      setBookingError("");
      const response = await requestWithSession("/api/bookings", {
        method: "POST",
        body: JSON.stringify({
          userId: user.id,
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
        email: response.email || user?.email || "aarav.shah@example.com",
        qrCodeDataUrl: response.qrCodeDataUrl || null,
      };

      setBookings((prev) => [booking, ...prev]);
      setLastBooking(booking);
      setSelectedSeats([]);
      setHoldDeadline(null);
      setPage("confirmation");
    } catch (error) {
      setBookingError(error.message || "Booking failed. Please try again.");
      setPage("payment");
    }
  }

  return (
    <div className="app-root">
      <GlobalStyles />
      <TopNav page={page} setPage={setPage} user={user} onEventsClick={() => setPage("events")} onBookingsClick={() => (user ? setPage("bookings") : setPage("profile"))} />
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
            holdDeadline={holdDeadline}
            user={user}
            onBack={releaseHeldSeats}
            onPay={confirmBooking}
            bookingError={bookingError}
          />
        )}

        {page === "confirmation" && lastBooking && (
          <ConfirmationPage booking={lastBooking} onDone={() => setPage(user ? "bookings" : "profile")} />
        )}

        {page === "profile" && <ProfilePage user={user} onLogin={handleLogin} onLogout={handleLogout} />}

        {page === "bookings" && <BookingsPage bookings={bookings} user={user} onRequireLogin={() => setPage("profile")} />}
      </main>
    </div>
  );
}
