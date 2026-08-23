import React, { useState } from "react";
import { TicketX } from "lucide-react";
import BookingRow from "../components/BookingRow";

export default function BookingsPage({ bookings = [] }) {
  const [tab, setTab] = useState("all");

  const filtered = bookings.filter((b) => {
    if (tab === "all") return true;
    return (b.status || "upcoming").toLowerCase() === tab;
  });

  return (
    <div className="page page-bookings">
      <h1 className="page-title">My bookings</h1>

      <div className="filter-chips">
        {["all", "upcoming", "completed", "cancelled"].map((t) => (
          <button key={t} className={"chip" + (tab === t ? " is-active" : "")} onClick={() => setTab(t)}>
            {t[0].toUpperCase() + t.slice(1)}
          </button>
        ))}
      </div>

      <div className="bookings-list">
        {filtered.map((b) => (
          <BookingRow key={b.id} booking={b} />
        ))}
        {filtered.length === 0 && (
          <div className="empty-state">
            <TicketX size={28} strokeWidth={1.75} />
            <p>Nothing here yet.</p>
          </div>
        )}
      </div>
    </div>
  );
}
