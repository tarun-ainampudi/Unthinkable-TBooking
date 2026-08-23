import React from "react";
import { MapPin, Calendar, Clock, Armchair, TicketCheck, CheckCircle2, TicketX, Download } from "lucide-react";
import { formatDate, formatINR } from "../utils";
import QRPlaceholder from "./QRPlaceholder";

const STATUS_META = {
  upcoming: { label: "Upcoming", icon: TicketCheck, cls: "status-upcoming" },
  completed: { label: "Completed", icon: CheckCircle2, cls: "status-completed" },
  cancelled: { label: "Cancelled", icon: TicketX, cls: "status-cancelled" },
};

export default function BookingRow({ booking }) {
  const statusMeta = STATUS_META[booking.status];

  return (
    <div className="booking-row ticket-stub">
      <div className="stub-notch stub-notch-left" />
      <div className="stub-notch stub-notch-right" />
      <div className="booking-row-main">
        <div>
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
        </div>
        <div className="booking-row-side">
          <span className={"status-badge " + statusMeta.cls}>
            <statusMeta.icon size={13} strokeWidth={2.5} /> {statusMeta.label}
          </span>
          <span className="booking-row-total">{formatINR(booking.total)}</span>
          <span className="booking-row-code">{booking.code}</span>
        </div>
      </div>
      {booking.status === "upcoming" && (
        <>
          <div className="stub-divider stub-divider-dashed" />
          <div className="booking-row-footer">
            <QRPlaceholder seed={booking.code} />
            <button className="btn btn-secondary">
              <Download size={14} strokeWidth={2.25} /> Download ticket
            </button>
          </div>
        </>
      )}
    </div>
  );
}
