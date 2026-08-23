import React from "react";
import { MapPin, Calendar, Clock, ChevronRight } from "lucide-react";
import { formatDate, formatINR } from "../utils";

export default function EventCard({ event, onSelect }) {
  return (
    <button className={`event-card accent-${event.accent}`} onClick={onSelect}>
      <div className="event-card-header">
        <div className="event-card-art">
          <span className="event-card-category">{event.category}</span>
          <span className="event-card-tag">{event.tag}</span>
        </div>
        <h3 className="event-card-title">{event.title}</h3>
      </div>
      <div className="stub-divider" />
      <div className="event-card-body">
        <p className="event-card-blurb">{event.blurb}</p>
        <div className="event-card-meta">
          <span><MapPin size={13} strokeWidth={2.25} /> {event.venue}</span>
          <span><Calendar size={13} strokeWidth={2.25} /> {formatDate(event.date)}</span>
          <span><Clock size={13} strokeWidth={2.25} /> {event.time}</span>
        </div>
        <div className="event-card-footer">
          <div>
            <span className="from-label">From</span>
            <span className="from-price">{formatINR(event.priceFrom)}</span>
          </div>
          <span className="event-card-cta">
            Select seats <ChevronRight size={15} strokeWidth={2.5} />
          </span>
        </div>
      </div>
    </button>
  );
}
