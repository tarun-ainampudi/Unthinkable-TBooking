import React from "react";
import { ArrowLeft, MapPin, Calendar, Clock, ChevronRight, Armchair } from "lucide-react";
import { formatDate, formatINR } from "../utils";
import SeatButton from "../components/SeatButton";

const ROWS = ["A", "B", "C", "D", "E", "F", "G", "H"];
const CATEGORY_META = {
  VIP: { label: "VIP", price: 2499, color: "var(--gold)" },
  Premium: { label: "Premium", price: 1499, color: "var(--teal)" },
  Standard: { label: "Standard", price: 899, color: "var(--paper-dim)" },
};

const HOLD_SECONDS = 8 * 60;
const MAX_SEATS = 8;

export default function SeatSelectionPage({ event, seats, selectedSeats, setSelectedSeats, onBack, onContinue }) {
  function toggleSeat(seat) {
    if (seat.status !== "available") return;
    setSelectedSeats((prev) => {
      const isSelected = prev.includes(seat.id);
      if (isSelected) return prev.filter((id) => id !== seat.id);
      if (prev.length >= MAX_SEATS) return prev;
      return [...prev, seat.id];
    });
  }

  const total = selectedSeats.reduce((sum, id) => {
    const seat = seats.find((s) => s.id === id);
    if (!seat || !CATEGORY_META[seat.category]) return sum;
    return sum + CATEGORY_META[seat.category].price;
  }, 0);

  const grouped = ROWS.map((row) => seats.filter((s) => s.row === row)).filter((rowSeats) => rowSeats.length > 0);

  return (
    <div className="page page-seats">
      <button className="back-link" onClick={onBack}>
        <ArrowLeft size={15} strokeWidth={2.25} /> Back to events
      </button>

      <div className="seats-header">
        <div>
          <h1>{event.title}</h1>
          <p className="seats-header-meta">
            <MapPin size={13} strokeWidth={2.25} /> {event.venue} &nbsp;·&nbsp;
            <Calendar size={13} strokeWidth={2.25} /> {formatDate(event.date)} &nbsp;·&nbsp;
            <Clock size={13} strokeWidth={2.25} /> {event.time}
          </p>
        </div>
      </div>

      <div className="seat-layout">
        <div className="seat-map-panel">
          <div className="stage">SCREEN / STAGE</div>
          {grouped.length === 0 ? (
            <div className="empty-state">
              <p>Seat map is still loading or unavailable for this event.</p>
            </div>
          ) : (
            <div className="seat-grid">
              {grouped.map((rowSeats) => (
                <div className="seat-row" key={rowSeats[0].row}>
                  <span className="row-label">{rowSeats[0].row}</span>
                  <div className="seat-row-seats">
                    {rowSeats.map((seat, idx) => (
                      <React.Fragment key={seat.id}>
                        {idx === 7 && <span className="aisle-gap" />}
                        <SeatButton
                          seat={seat}
                          isSelected={selectedSeats.includes(seat.id)}
                          onClick={() => toggleSeat(seat)}
                        />
                      </React.Fragment>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}

          <div className="seat-legend">
            {Object.entries(CATEGORY_META).map(([key, meta]) => (
              <span className="legend-item" key={key}>
                <i className="legend-swatch" style={{ background: meta.color }} />
                {meta.label} · {formatINR(meta.price)}
              </span>
            ))}
            <span className="legend-item">
              <i className="legend-swatch swatch-booked" /> Booked
            </span>
            <span className="legend-item">
              <i className="legend-swatch swatch-held" /> On hold
            </span>
          </div>
        </div>

        <aside className="seat-summary-panel">
          <h3>Your selection</h3>
          {selectedSeats.length === 0 ? (
            <p className="summary-empty">
              <Armchair size={18} strokeWidth={1.75} /> Tap available seats to add them here.
            </p>
          ) : (
            <ul className="summary-seat-list">
              {selectedSeats.slice().sort().map((id) => {
                const seat = seats.find((s) => s.id === id);
                return (
                  <li key={id}>
                    <span>
                      <span className="summary-seat-id">{id}</span>
                      <span className="summary-seat-cat">{seat.category}</span>
                    </span>
                    <span>{formatINR(CATEGORY_META[seat.category].price)}</span>
                  </li>
                );
              })}
            </ul>
          )}
          <div className="stub-divider" />
          <div className="summary-total-row">
            <span>Total</span>
            <span>{formatINR(total)}</span>
          </div>
          <button
            className="btn btn-primary btn-block"
            disabled={selectedSeats.length === 0 || grouped.length === 0}
            onClick={onContinue}
          >
            Continue to payment <ChevronRight size={16} strokeWidth={2.5} />
          </button>
          <p className="summary-note">Max {MAX_SEATS} seats per booking. Hold releases automatically on expiry.</p>
        </aside>
      </div>
    </div>
  );
}
