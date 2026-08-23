export default function GlobalStyles() {
  return (
    <style>{`
      @import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&family=Inter:wght@400;500;600;700;800&family=IBM+Plex+Mono:wght@400;500;600&display=swap');

      :root {
        --bg: #14121a;
        --bg-alt: #1b1826;
        --paper: #f3eee3;
        --paper-dim: #cfc6b3;
        --ink: #1e1b24;
        --muted-ink: #948c7a;
        --muted-on-dark: #a39cb3;
        --gold: #e3b23c;
        --gold-dim: #a9873a;
        --teal: #4fa393;
        --teal-dim: #3f7d73;
        --red: #d1584f;
        --red-dim: #c1443c;
        --line: rgba(243,238,227,0.12);
        --line-strong: rgba(243,238,227,0.22);
        --font-display: 'Bebas Neue', sans-serif;
        --font-body: 'Inter', sans-serif;
        --font-mono: 'IBM Plex Mono', monospace;
        --radius: 14px;
      }

      * { box-sizing: border-box; }

      .app-root {
        background: var(--bg);
        background-image:
          radial-gradient(circle at 15% 0%, rgba(227,178,60,0.07), transparent 45%),
          radial-gradient(circle at 85% 20%, rgba(79,163,147,0.08), transparent 40%);
        color: var(--paper);
        font-family: var(--font-body);
        min-height: 100vh;
        width: 100%;
      }

      button { font-family: inherit; cursor: pointer; }
      input { font-family: inherit; }
      ul { list-style: none; margin: 0; padding: 0; }

      /* ---------- Nav ---------- */
      .topnav {
        position: sticky;
        top: 0;
        z-index: 20;
        background: rgba(20,18,26,0.85);
        backdrop-filter: blur(10px);
        border-bottom: 1px solid var(--line);
      }
      .topnav-inner {
        max-width: 1180px;
        margin: 0 auto;
        padding: 14px 24px;
        display: flex;
        align-items: center;
        justify-content: space-between;
      }
      .brand {
        display: flex;
        align-items: center;
        gap: 8px;
        background: none;
        border: none;
        color: var(--paper);
        padding: 0;
      }
      .brand-mark {
        width: 30px;
        height: 30px;
        border-radius: 8px;
        background: var(--gold);
        color: var(--ink);
        display: flex;
        align-items: center;
        justify-content: center;
        transform: rotate(-8deg);
      }
      .brand-word {
        font-family: var(--font-display);
        font-size: 22px;
        letter-spacing: 0.06em;
      }
      .brand-word-accent { color: var(--gold); }

      .topnav-links { display: flex; gap: 4px; }
      .topnav-link {
        display: flex;
        align-items: center;
        gap: 6px;
        background: none;
        border: none;
        color: var(--muted-on-dark);
        font-size: 13.5px;
        font-weight: 600;
        padding: 8px 14px;
        border-radius: 100px;
        transition: color 0.15s, background 0.15s;
      }
      .topnav-link:hover { color: var(--paper); background: rgba(255,255,255,0.05); }
      .topnav-link.is-active { color: var(--ink); background: var(--paper); }

      /* ---------- Layout ---------- */
      .app-main { max-width: 1180px; margin: 0 auto; padding: 32px 24px 80px; }
      .page-title {
        font-family: var(--font-display);
        font-size: 34px;
        letter-spacing: 0.02em;
        margin: 4px 0 22px;
      }
      .back-link {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        background: none;
        border: none;
        color: var(--muted-on-dark);
        font-size: 13px;
        font-weight: 600;
        margin-bottom: 18px;
        padding: 0;
      }
      .back-link:hover { color: var(--gold); }

      /* ---------- Events hero ---------- */
      .events-hero { padding: 8px 0 28px; border-bottom: 1px solid var(--line); margin-bottom: 24px; }
      .eyebrow {
        font-family: var(--font-mono);
        text-transform: uppercase;
        font-size: 11.5px;
        letter-spacing: 0.14em;
        color: var(--gold);
        margin: 0 0 10px;
      }
      .hero-title {
        font-family: var(--font-display);
        font-size: 52px;
        line-height: 1.02;
        letter-spacing: 0.01em;
        margin: 0 0 12px;
        max-width: 640px;
      }
      .hero-title-accent { color: var(--teal); }
      .hero-sub { color: var(--muted-on-dark); font-size: 15px; max-width: 480px; line-height: 1.55; }

      .events-controls {
        display: flex;
        flex-wrap: wrap;
        gap: 14px;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 26px;
      }
      .search-box {
        display: flex;
        align-items: center;
        gap: 8px;
        background: var(--bg-alt);
        border: 1px solid var(--line);
        border-radius: 100px;
        padding: 10px 16px;
        color: var(--muted-on-dark);
        min-width: 260px;
      }
      .search-box input {
        background: none;
        border: none;
        outline: none;
        color: var(--paper);
        font-size: 13.5px;
        width: 100%;
      }
      .search-box input::placeholder { color: var(--muted-on-dark); }

      .filter-chips { display: flex; gap: 8px; flex-wrap: wrap; }
      .chip {
        background: var(--bg-alt);
        border: 1px solid var(--line);
        color: var(--muted-on-dark);
        font-size: 12.5px;
        font-weight: 600;
        padding: 8px 14px;
        border-radius: 100px;
        transition: all 0.15s;
      }
      .chip:hover { border-color: var(--line-strong); color: var(--paper); }
      .chip.is-active { background: var(--gold); border-color: var(--gold); color: var(--ink); }

      /* ---------- Event cards ---------- */
      .events-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 18px;
      }
      .event-card {
        text-align: left;
        background: var(--paper);
        color: var(--ink);
        border: none;
        border-radius: var(--radius);
        overflow: hidden;
        display: flex;
        flex-direction: column;
        transition: transform 0.18s, box-shadow 0.18s;
        padding: 0;
      }
      .event-card:hover { transform: translateY(-3px); box-shadow: 0 14px 30px rgba(0,0,0,0.35); }
      .event-card-header {
        height: 100px;
        padding: 14px 16px;
      }
      .event-card-art {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
      }
      .accent-gold .event-card-header { background: linear-gradient(135deg, #e3b23c, #a9873a); }
      .accent-teal .event-card-header { background: linear-gradient(135deg, #4fa393, #2c5951); }
      .accent-red .event-card-header { background: linear-gradient(135deg, #d1584f, #8a332c); }
      .event-card-category {
        font-family: var(--font-mono);
        font-size: 11px;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        background: rgba(20,18,26,0.28);
        color: #fff;
        padding: 4px 9px;
        border-radius: 100px;
      }
      .event-card-tag {
        font-family: var(--font-display);
        font-size: 22px;
        color: rgba(255,255,255,0.85);
        letter-spacing: 0.03em;
      }
      .event-card-title {
        color: white;
        font-style: bold;
        font-size: 20px;
      }
      .stub-divider {
        height: 0;
        border-top: 1.5px dashed rgba(30,27,36,0.18);
        margin: 0 16px;
      }
      .event-card-body { padding: 14px 16px 18px; display: flex; flex-direction: column; gap: 10px; flex: 1; }
      .event-card-body h3 { font-size: 17px; margin: 0; letter-spacing: -0.01em; }
      .event-card-blurb { font-size: 12.5px; color: var(--muted-ink); line-height: 1.5; margin: 0; }
      .event-card-meta { display: flex; flex-direction: column; gap: 4px; font-size: 12px; color: var(--muted-ink); }
      .event-card-meta span { display: flex; align-items: center; gap: 6px; }
      .event-card-footer {
        margin-top: auto;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding-top: 10px;
      }
      .from-label { display: block; font-size: 10.5px; color: var(--muted-ink); text-transform: uppercase; letter-spacing: 0.08em; }
      .from-price { font-family: var(--font-display); font-size: 20px; }
      .event-card-cta {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12.5px;
        font-weight: 700;
        color: var(--ink);
      }

      .empty-state {
        grid-column: 1 / -1;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 10px;
        color: var(--muted-on-dark);
        padding: 60px 0;
        text-align: center;
      }

      /* ---------- Buttons / fields ---------- */
      .btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        font-size: 13.5px;
        font-weight: 700;
        border-radius: 10px;
        padding: 12px 18px;
        border: 1px solid transparent;
        transition: transform 0.1s, opacity 0.15s, background 0.15s;
      }
      .btn:active { transform: scale(0.98); }
      .btn:disabled { opacity: 0.4; cursor: not-allowed; }
      .btn-primary { background: var(--gold); color: var(--ink); }
      .btn-primary:hover:not(:disabled) { background: #edc25a; }
      .btn-secondary { background: transparent; border-color: var(--line-strong); color: var(--paper); }
      .btn-secondary:hover { border-color: var(--gold); color: var(--gold); }
      .btn-block { width: 100%; }

      .field-group { display: flex; flex-direction: column; gap: 14px; }
      .field-row { display: flex; gap: 12px; }
      .field-row .field { flex: 1; }
      .field { display: flex; flex-direction: column; gap: 6px; font-size: 12.5px; color: var(--muted-on-dark); font-weight: 600; }
      .field input {
        background: var(--bg-alt);
        border: 1px solid var(--line);
        border-radius: 9px;
        padding: 11px 13px;
        color: var(--paper);
        font-size: 13.5px;
        outline: none;
      }
      .field input:focus { border-color: var(--gold); }
      .field-hint { font-size: 11.5px; color: var(--muted-on-dark); margin: -4px 0 0; }

      /* ---------- Seat selection ---------- */
      .seats-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 20px;
        margin-bottom: 22px;
        flex-wrap: wrap;
      }
      .seats-header h1 { font-family: var(--font-display); font-size: 30px; margin: 0 0 6px; }
      .seats-header-meta {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12.5px;
        color: var(--muted-on-dark);
        flex-wrap: wrap;
      }
      .hold-timer {
        background: var(--bg-alt);
        border: 1px solid var(--gold-dim);
        border-radius: 10px;
        padding: 8px 16px;
        text-align: center;
      }
      .hold-timer.is-expired { border-color: var(--red); }
      .hold-timer-label { display: block; font-size: 10px; text-transform: uppercase; letter-spacing: 0.08em; color: var(--muted-on-dark); }
      .hold-timer-clock { font-family: var(--font-mono); font-size: 19px; color: var(--gold); }
      .is-expired .hold-timer-clock { color: var(--red); }

      .seat-layout { display: grid; grid-template-columns: 1fr 300px; gap: 20px; align-items: start; }
      .seat-map-panel {
        background: var(--bg-alt);
        border: 1px solid var(--line);
        border-radius: var(--radius);
        padding: 26px;
      }
      .stage {
        text-align: center;
        font-family: var(--font-mono);
        font-size: 11px;
        letter-spacing: 0.25em;
        color: var(--muted-on-dark);
        background: linear-gradient(180deg, rgba(227,178,60,0.14), transparent);
        border-bottom: 2px solid var(--gold-dim);
        padding: 10px 0 14px;
        margin-bottom: 22px;
        border-radius: 100px 100px 0 0;
      }
      .seat-grid { display: flex; flex-direction: column; gap: 7px; overflow-x: auto; }
      .seat-row { display: flex; align-items: center; gap: 10px; }
      .row-label { width: 16px; font-family: var(--font-mono); font-size: 11px; color: var(--muted-on-dark); flex-shrink: 0; }
      .seat-row-seats { display: flex; gap: 5px; }
      .aisle-gap { width: 14px; flex-shrink: 0; }

      .seat {
        width: 26px;
        height: 24px;
        border-radius: 6px 6px 3px 3px;
        border: 1px solid var(--line-strong);
        background: transparent;
        color: var(--muted-on-dark);
        font-family: var(--font-mono);
        font-size: 9.5px;
        flex-shrink: 0;
        transition: transform 0.1s, background 0.15s;
      }
      .seat:not(:disabled):hover { transform: translateY(-2px); border-color: var(--paper); color: var(--paper); }
      .seat-vip { border-color: rgba(227,178,60,0.45); color: rgba(227,178,60,0.85); }
      .seat-premium { border-color: rgba(79,163,147,0.45); color: rgba(79,163,147,0.85); }
      .seat-standard { border-color: var(--line-strong); }
      .seat-selected { background: var(--gold); border-color: var(--gold); color: var(--ink); font-weight: 700; }
      .seat-booked { background: rgba(209,88,79,0.12); border-color: rgba(209,88,79,0.3); color: rgba(209,88,79,0.6); cursor: not-allowed; }
      .seat-held { background: repeating-linear-gradient(45deg, rgba(148,140,122,0.15), rgba(148,140,122,0.15) 3px, transparent 3px, transparent 6px); border-color: var(--line-strong); color: var(--muted-on-dark); cursor: not-allowed; }

      .seat-legend { display: flex; flex-wrap: wrap; gap: 16px; margin-top: 22px; padding-top: 18px; border-top: 1px solid var(--line); }
      .legend-item { display: flex; align-items: center; gap: 7px; font-size: 12px; color: var(--muted-on-dark); }
      .legend-swatch { width: 12px; height: 12px; border-radius: 4px; display: inline-block; }
      .swatch-booked { background: rgba(209,88,79,0.35); }
      .swatch-held { background: repeating-linear-gradient(45deg, rgba(148,140,122,0.4), rgba(148,140,122,0.4) 2px, transparent 2px, transparent 4px); }

      .seat-summary-panel {
        background: var(--paper);
        color: var(--ink);
        border-radius: var(--radius);
        padding: 22px;
        position: sticky;
        top: 84px;
      }
      .seat-summary-panel h3 { font-family: var(--font-display); font-size: 19px; margin: 0 0 14px; letter-spacing: 0.02em; }
      .summary-empty { display: flex; flex-direction: column; align-items: center; gap: 8px; color: var(--muted-ink); font-size: 12.5px; text-align: center; padding: 22px 0; }
      .summary-seat-list { display: flex; flex-direction: column; gap: 10px; margin-bottom: 4px; }
      .summary-seat-list li { display: flex; justify-content: space-between; font-size: 13px; }
      .summary-seat-id { font-family: var(--font-mono); font-weight: 700; margin-right: 8px; }
      .summary-seat-cat { font-size: 11px; color: var(--muted-ink); }
      .summary-total-row { display: flex; justify-content: space-between; font-family: var(--font-display); font-size: 21px; margin: 14px 0 16px; letter-spacing: 0.02em; }
      .summary-line { display: flex; justify-content: space-between; font-size: 12.5px; color: var(--muted-ink); margin-bottom: 8px; }
      .summary-note { font-size: 11px; color: var(--muted-ink); text-align: center; margin: 12px 0 0; line-height: 1.5; }

      /* ---------- Payment ---------- */
      .payment-layout { display: grid; grid-template-columns: 1fr 320px; gap: 22px; align-items: start; }
      .payment-form {
        background: var(--bg-alt);
        border: 1px solid var(--line);
        border-radius: var(--radius);
        padding: 24px;
        display: flex;
        flex-direction: column;
        gap: 18px;
      }
      .method-switch { display: flex; gap: 8px; background: rgba(255,255,255,0.03); padding: 4px; border-radius: 10px; }
      .method-tab {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 7px;
        background: none;
        border: none;
        color: var(--muted-on-dark);
        font-size: 13px;
        font-weight: 700;
        padding: 9px 0;
        border-radius: 8px;
      }
      .method-tab.is-active { background: var(--gold); color: var(--ink); }
      .secure-note { display: flex; align-items: center; gap: 7px; font-size: 11.5px; color: var(--muted-on-dark); background: rgba(255,255,255,0.03); padding: 10px 12px; border-radius: 8px; }

      /* ---------- Ticket stub shared component ---------- */
      .ticket-stub {
        background: var(--paper);
        color: var(--ink);
        border-radius: var(--radius);
        padding: 24px;
        position: relative;
      }
      .stub-notch { position: absolute; width: 20px; height: 20px; background: var(--bg); border-radius: 50%; top: 50%; transform: translateY(-50%); }
      .order-summary .stub-notch, .confirmation-card .stub-notch, .booking-row .stub-notch { display: none; }
      .stub-divider-dashed { border-top: 1.5px dashed rgba(30,27,36,0.2); margin: 16px 0; }
      .order-summary h3 { font-family: var(--font-display); font-size: 21px; margin: 0 0 8px; letter-spacing: 0.02em; }
      .order-summary-meta { display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--muted-ink); margin: 0 0 4px; }

      /* ---------- Confirmation ---------- */
      .page-confirmation { display: flex; justify-content: center; padding-top: 20px; }
      .confirmation-card { max-width: 400px; width: 100%; text-align: center; padding: 34px 28px; }
      .confirmation-icon { width: 56px; height: 56px; border-radius: 50%; background: var(--teal); color: #fff; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
      .confirmation-card h1 { font-family: var(--font-display); font-size: 26px; margin: 0 0 6px; letter-spacing: 0.02em; }
      .confirmation-sub { font-size: 12.5px; color: var(--muted-ink); margin: 0 0 4px; }
      .confirmation-card h3 { font-family: var(--font-display); font-size: 20px; margin: 4px 0 8px; }

      .qr-wrap { display: flex; justify-content: center; margin: 6px 0 10px; }
      .qr-grid { display: grid; grid-template-columns: repeat(7, 8px); grid-template-rows: repeat(7, 8px); gap: 2px; background: #fff; padding: 10px; border-radius: 8px; }
      .qr-cell { width: 8px; height: 8px; background: transparent; }
      .qr-cell.on { background: var(--ink); }
      .booking-code { font-family: var(--font-mono); letter-spacing: 0.08em; font-size: 13px; color: var(--muted-ink); margin: 0 0 20px; }

      /* ---------- Profile ---------- */
      .page-profile {
        display: flex;
        flex-direction: column;
        align-items: center;
      }
      .page-profile.profile-view {
        align-items: stretch;
      }
      .profile-layout { display: grid; grid-template-columns: 320px 1fr; gap: 20px; align-items: start; }
      .login-layout {
        width: min(100%, 420px);
        grid-template-columns: 1fr;
        justify-items: center;
      }
      .login-card {
        width: 100%;
      }
      .profile-card {
        background: var(--paper);
        color: var(--ink);
        border-radius: var(--radius);
        padding: 28px;
        text-align: center;
      }
      .profile-avatar {
        width: 68px;
        height: 68px;
        border-radius: 50%;
        background: var(--gold);
        color: var(--ink);
        font-family: var(--font-display);
        font-size: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 14px;
      }
      .profile-card h2 { margin: 0 0 4px; font-size: 20px; }
      .profile-member-since { font-size: 12px; color: var(--muted-ink); margin: 0 0 18px; }
      .profile-detail-list { display: flex; flex-direction: column; gap: 10px; margin-bottom: 20px; text-align: left; }
      .profile-detail { display: flex; align-items: center; gap: 9px; font-size: 13px; background: rgba(30,27,36,0.05); padding: 10px 12px; border-radius: 8px; }
      .profile-stats { display: flex; flex-direction: column; gap: 14px; }
      .stat-card { background: var(--bg-alt); border: 1px solid var(--line); border-radius: 12px; padding: 18px 20px; display: flex; align-items: baseline; justify-content: space-between; }
      .stat-number { font-family: var(--font-display); font-size: 30px; color: var(--gold); }
      .stat-label { font-size: 12.5px; color: var(--muted-on-dark); }
      .logout-btn { margin-top: 6px; }

      /* ---------- Bookings ---------- */
      .bookings-list { display: flex; flex-direction: column; gap: 16px; margin-top: 20px; }
      .booking-row { padding: 22px 26px; }
      .booking-row-main { display: flex; justify-content: space-between; gap: 20px; flex-wrap: wrap; }
      .booking-row-main h3 { font-family: var(--font-display); font-size: 19px; margin: 0 0 8px; letter-spacing: 0.02em; }
      .booking-row-side { display: flex; flex-direction: column; align-items: flex-end; gap: 6px; }
      .booking-row-total { font-family: var(--font-display); font-size: 19px; }
      .booking-row-code { font-family: var(--font-mono); font-size: 11px; color: var(--muted-ink); }
      .status-badge { display: flex; align-items: center; gap: 5px; font-size: 11px; font-weight: 700; padding: 4px 10px; border-radius: 100px; text-transform: uppercase; letter-spacing: 0.04em; }
      .status-upcoming { background: rgba(79,163,147,0.15); color: var(--teal-dim); }
      .status-completed { background: rgba(30,27,36,0.08); color: var(--muted-ink); }
      .status-cancelled { background: rgba(209,88,79,0.15); color: var(--red-dim); }
      .booking-row-footer { display: flex; align-items: center; justify-content: space-between; gap: 16px; }

      /* ---------- Responsive ---------- */
      @media (max-width: 860px) {
        .hero-title { font-size: 38px; }
        .seat-layout, .payment-layout, .profile-layout { grid-template-columns: 1fr; }
        .seat-summary-panel { position: static; }
        .topnav-link span { display: none; }
      }
      @media (max-width: 520px) {
        .app-main { padding: 22px 16px 60px; }
        .topnav-inner { padding: 12px 16px; }
        .topnav-link { padding: 8px 10px; }
        .events-grid { grid-template-columns: 1fr; }
        .seat { width: 22px; height: 21px; }
      }

      @media (prefers-reduced-motion: reduce) {
        * { transition: none !important; }
      }
    `}</style>
  );
}
