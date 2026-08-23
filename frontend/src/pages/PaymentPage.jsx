import React, { useState } from "react";
import { ArrowLeft, CreditCard, Smartphone, ShieldCheck, MapPin, Calendar, Clock } from "lucide-react";
import { formatDate, formatINR } from "../utils";

const CATEGORY_META = {
  VIP: { label: "VIP", price: 2499, color: "var(--gold)" },
  Premium: { label: "Premium", price: 1499, color: "var(--teal)" },
  Standard: { label: "Standard", price: 899, color: "var(--paper-dim)" },
};

export default function PaymentPage({ event, seats, selectedSeats, onBack, onPay, bookingError }) {
  const [method, setMethod] = useState("card");
  const [processing, setProcessing] = useState(false);
  const [card, setCard] = useState({ number: "", name: "", expiry: "", cvv: "" });
  const [upi, setUpi] = useState("");

  const seatDetails = selectedSeats
    .slice()
    .sort()
    .map((id) => seats.find((s) => s.id === id));

  const subtotal = seatDetails.reduce((sum, s) => sum + CATEGORY_META[s.category].price, 0);
  const fee = Math.round(subtotal * 0.03);
  const total = subtotal + fee;

  const canPay =
    method === "card"
      ? card.number.replace(/\s/g, "").length >= 12 && card.name.trim() && card.expiry.length >= 4 && card.cvv.length >= 3
      : upi.includes("@") && upi.length > 4;

  function handlePay(e) {
    e.preventDefault();
    if (!canPay || processing) return;
    setProcessing(true);
    setTimeout(() => {
      setProcessing(false);
      onPay();
    }, 1100);
  }

  return (
    <div className="page page-payment">
      <button className="back-link" onClick={onBack}>
        <ArrowLeft size={15} strokeWidth={2.25} /> Back to seat selection
      </button>

      <h1 className="page-title">Payment</h1>

      <div className="payment-layout">
        <form className="payment-form" onSubmit={handlePay}>
          <div className="method-switch">
            <button
              type="button"
              className={"method-tab" + (method === "card" ? " is-active" : "")}
              onClick={() => setMethod("card")}
            >
              <CreditCard size={16} strokeWidth={2.25} /> Card
            </button>
            <button
              type="button"
              className={"method-tab" + (method === "upi" ? " is-active" : "")}
              onClick={() => setMethod("upi")}
            >
              <Smartphone size={16} strokeWidth={2.25} /> UPI
            </button>
          </div>

          {method === "card" ? (
            <div className="field-group">
              <label className="field">
                <span>Card number</span>
                <input
                  inputMode="numeric"
                  placeholder="1234 5678 9012 3456"
                  value={card.number}
                  maxLength={19}
                  onChange={(e) =>
                    setCard({
                      ...card,
                      number: e.target.value.replace(/[^\d]/g, "").replace(/(.{4})/g, "$1 ").trim(),
                    })
                  }
                />
              </label>
              <label className="field">
                <span>Name on card</span>
                <input
                  placeholder="Aarav Shah"
                  value={card.name}
                  onChange={(e) => setCard({ ...card, name: e.target.value })}
                />
              </label>
              <div className="field-row">
                <label className="field">
                  <span>Expiry</span>
                  <input
                    placeholder="MM/YY"
                    maxLength={5}
                    value={card.expiry}
                    onChange={(e) => {
                      let v = e.target.value.replace(/[^\d]/g, "");
                      if (v.length > 2) v = v.slice(0, 2) + "/" + v.slice(2, 4);
                      setCard({ ...card, expiry: v });
                    }}
                  />
                </label>
                <label className="field">
                  <span>CVV</span>
                  <input
                    inputMode="numeric"
                    placeholder="123"
                    maxLength={3}
                    value={card.cvv}
                    onChange={(e) => setCard({ ...card, cvv: e.target.value.replace(/[^\d]/g, "") })}
                  />
                </label>
              </div>
            </div>
          ) : (
            <div className="field-group">
              <label className="field">
                <span>UPI ID</span>
                <input
                  placeholder="yourname@upi"
                  value={upi}
                  onChange={(e) => setUpi(e.target.value)}
                />
              </label>
              <p className="field-hint">You'll get a payment request on your UPI app to approve.</p>
            </div>
          )}

          <div className="secure-note">
            <ShieldCheck size={15} strokeWidth={2.25} /> This is a prototype form — no real payment is processed.
          </div>

          {bookingError && <p className="field-hint" style={{ color: "#b42318" }}>{bookingError}</p>}

          <button className="btn btn-primary btn-block" type="submit" disabled={!canPay || processing}>
            {processing ? "Processing…" : `Pay ${formatINR(total)}`}
          </button>
        </form>

        <aside className="order-summary ticket-stub">
          <h3>{event.title}</h3>
          <p className="order-summary-meta">
            <MapPin size={13} strokeWidth={2.25} /> {event.venue}
          </p>
          <p className="order-summary-meta">
            <Calendar size={13} strokeWidth={2.25} /> {formatDate(event.date)} &nbsp;·&nbsp;
            <Clock size={13} strokeWidth={2.25} /> {event.time}
          </p>
          <div className="stub-divider stub-divider-dashed" />
          <ul className="summary-seat-list">
            {seatDetails.map((s) => (
              <li key={s.id}>
                <span>
                  <span className="summary-seat-id">{s.id}</span>
                  <span className="summary-seat-cat">{s.category}</span>
                </span>
                <span>{formatINR(CATEGORY_META[s.category].price)}</span>
              </li>
            ))}
          </ul>
          <div className="stub-divider stub-divider-dashed" />
          <div className="summary-line">
            <span>Subtotal</span>
            <span>{formatINR(subtotal)}</span>
          </div>
          <div className="summary-line">
            <span>Convenience fee</span>
            <span>{formatINR(fee)}</span>
          </div>
          <div className="summary-total-row">
            <span>Total</span>
            <span>{formatINR(total)}</span>
          </div>
        </aside>
      </div>
    </div>
  );
}
