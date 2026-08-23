import React, { useMemo } from "react";

/**
 * Decorative QR-style grid derived from a seed string (e.g. booking code).
 * This is NOT a real scannable QR code — it's a visual stand-in for the
 * frontend prototype. Real QR generation should happen server-side
 * (e.g. with a Go library like `go-qrcode`) and be rendered as an image
 * or embedded in the confirmation email.
 */
export default function QRPlaceholder({ seed }) {
  const cells = useMemo(() => {
    let h = 0;
    for (let i = 0; i < seed.length; i++) h = (h * 31 + seed.charCodeAt(i)) >>> 0;
    const grid = [];
    for (let i = 0; i < 49; i++) {
      h = (h * 1103515245 + 12345) >>> 0;
      grid.push((h >> 16) % 3 === 0);
    }
    return grid;
  }, [seed]);

  return (
    <div className="qr-wrap">
      <div className="qr-grid">
        {cells.map((on, i) => (
          <span key={i} className={on ? "qr-cell on" : "qr-cell"} />
        ))}
      </div>
    </div>
  );
}
