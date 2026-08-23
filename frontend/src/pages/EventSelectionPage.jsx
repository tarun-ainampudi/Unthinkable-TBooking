import React, { useEffect, useState } from "react";
import { Search, TicketX } from "lucide-react";
import EventCard from "../components/EventCard";
import { apiFetch } from "../api";

export default function EventSelectionPage({ onSelect }) {
  const [events, setEvents] = useState([]);
  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState("All");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadEvents = async () => {
      try {
        const data = await apiFetch("/api/events");
        setEvents(data || []);
      } catch (error) {
        console.error("Failed to load events", error);
        setEvents([]);
      } finally {
        setLoading(false);
      }
    };

    loadEvents();
  }, []);

  const categories = ["All", ...Array.from(new Set(events.map((e) => e.category)))];

  const filtered = events.filter((e) => {
    const matchesQuery =
      e.title.toLowerCase().includes(query.toLowerCase()) ||
      e.venue.toLowerCase().includes(query.toLowerCase());
    const matchesFilter = filter === "All" || e.category === filter;
    return matchesQuery && matchesFilter;
  });

  return (
    <div className="page page-events">
      <section className="events-hero">
        <p className="eyebrow">Amaravati &amp; nearby · this month</p>
        <h1 className="hero-title">
          Find your <span className="hero-title-accent">next night out.</span>
        </h1>
        <p className="hero-sub">
          Live music, theatre, comedy and dance — pick a show, pick a seat, walk in with a QR code.
        </p>
      </section>

      <section className="events-controls">
        <div className="search-box">
          <Search size={16} strokeWidth={2.25} />
          <input
            placeholder="Search by show or venue"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <div className="filter-chips">
          {categories.map((c) => (
            <button
              key={c}
              className={"chip" + (filter === c ? " is-active" : "")}
              onClick={() => setFilter(c)}
            >
              {c}
            </button>
          ))}
        </div>
      </section>

      <section className="events-grid">
        {loading ? (
          <div className="empty-state">
            <p>Loading events…</p>
          </div>
        ) : (
          filtered.map((event) => (
            <EventCard key={event.id} event={event} onSelect={() => onSelect(event)} />
          ))
        )}
        {!loading && filtered.length === 0 && (
          <div className="empty-state">
            <TicketX size={28} strokeWidth={1.75} />
            <p>No shows match that search. Try another title or venue.</p>
          </div>
        )}
      </section>
    </div>
  );
}
