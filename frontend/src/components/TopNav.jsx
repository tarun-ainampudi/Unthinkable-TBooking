import React from "react";
import { Ticket, History, User } from "lucide-react";

export default function TopNav({ page, setPage, user, onEventsClick, onBookingsClick }) {
  const profileLabel = user ? "Profile" : "Login";

  const links = [
    { id: "events", label: "Events", icon: Ticket, onClick: onEventsClick || (() => setPage("events")) },
    ...(user ? [{ id: "bookings", label: "My Bookings", icon: History, onClick: onBookingsClick || (() => setPage("bookings")) }] : []),
    { id: "profile", label: profileLabel, icon: User, onClick: () => setPage("profile") },
  ];

  return (
    <header className="topnav">
      <div className="topnav-inner">
        <button className="brand" onClick={() => (user ? setPage("events") : setPage("profile"))}>
          <span className="brand-mark">
            <Ticket size={18} strokeWidth={2.25} />
          </span>
          <span className="brand-word">
            Un<span className="brand-word-accent">thinkable</span>
          </span>
        </button>
        <nav className="topnav-links">
          {links.map((l) => (
            <button
              key={l.id}
              className={"topnav-link" + (page === l.id ? " is-active" : "")}
              onClick={l.onClick}
            >
              <l.icon size={15} strokeWidth={2.25} />
              <span>{l.label}</span>
            </button>
          ))}
        </nav>
      </div>
    </header>
  );
}
