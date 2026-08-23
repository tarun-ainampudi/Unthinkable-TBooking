import React from "react";

export default function SeatButton({ seat, isSelected, onClick }) {
  let cls = "seat";
  if (seat.status === "booked") cls += " seat-booked";
  else if (seat.status === "held") cls += " seat-held";
  else if (isSelected) cls += " seat-selected";
  else cls += ` seat-${seat.category.toLowerCase()}`;

  const disabled = seat.status !== "available";

  return (
    <button
      className={cls}
      onClick={onClick}
      disabled={disabled}
      title={`${seat.id} · ${seat.category}${disabled ? ` · ${seat.status}` : ""}`}
    >
      {seat.col}
    </button>
  );
}
