import React from "react";
import { CheckCircle2, MapPin, Calendar, Clock, Armchair } from "lucide-react";
import { formatDate } from "../utils";
import QRPlaceholder from "../components/QRPlaceholder";

export default function ConfirmationPage({ booking, onDone }) {
  return (
    <div className="page page-confirmation">
      <div className="confirmation-card ticket-stub">
        <div className="confirmation-icon">
          <CheckCircle2 size={30} strokeWidth={2} />
        </div>
        <h1>Booking confirmed</h1>
        <p className="confirmation-sub">
          A copy of this ticket and QR code has been emailed to {booking.email || "your inbox"}.
        </p>
        <div className="stub-divider stub-divider-dashed" />
        <h3>{booking.eventTitle}</h3>
        <p className="order-summary-meta">
          <MapPin size={13} strokeWidth={2.25} /> {booking.venue}
        </p>
        <p className="order-summary-meta">
          <Calendar size={13} strokeWidth={2.25} /> {formatDate(booking.date)} &nbsp;·&nbsp;
          <Clock size={13} strokeWidth={2.25} /> {booking.time}
        </p>
        <p className="order-summary-meta">
          <Armchair size={13} strokeWidth={2.25} /> Seats {booking.seatLabels.join(", ")}
        </p>
        <div className="stub-divider stub-divider-dashed" />
        {booking.qrCodeDataUrl ? (
          <img src={booking.qrCodeDataUrl} alt="Booking QR code" style={{ width: 180, height: 180, margin: "8px auto 12px", display: "block", background: "#fff", padding: 8, borderRadius: 12 }} />
        ) : (
          <QRPlaceholder seed={booking.code} />
        )}
        <p className="booking-code">{booking.code}</p>
        <button className="btn btn-primary btn-block" onClick={onDone}>
          View my bookings
        </button>
      </div>
    </div>
  );
}
