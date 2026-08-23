/* ============================== MOCK DATA ============================== */
/* Replace all of this with real API calls to your Go backend. */

export const EVENTS = [
  {
    id: 1,
    title: "Midnight Frequencies",
    category: "Concert",
    tag: "Electronic",
    venue: "Marquee Arena, Amaravati",
    date: "2026-09-12",
    time: "8:00 PM",
    priceFrom: 899,
    rating: 4.8,
    accent: "gold",
    blurb: "A night-long set from three touring electronic acts, built for a room that likes its bass loud.",
  },
  {
    id: 2,
    title: "The Glass Menagerie",
    category: "Theatre",
    tag: "Drama",
    venue: "Riverside Playhouse",
    date: "2026-09-18",
    time: "7:00 PM",
    priceFrom: 499,
    rating: 4.6,
    accent: "teal",
    blurb: "Tennessee Williams' family drama, staged in the round by the Riverside repertory company.",
  },
  {
    id: 3,
    title: "Laugh Riot Live",
    category: "Comedy",
    tag: "Stand-up",
    venue: "The Attic Club",
    date: "2026-08-29",
    time: "9:00 PM",
    priceFrom: 349,
    rating: 4.9,
    accent: "red",
    blurb: "Four comics, one mic, zero warning about what happens if you sit in the front row.",
  },
  {
    id: 4,
    title: "Nocturne: A Jazz Evening",
    category: "Concert",
    tag: "Jazz",
    venue: "Blue Room Hall",
    date: "2026-09-05",
    time: "8:30 PM",
    priceFrom: 699,
    rating: 4.7,
    accent: "teal",
    blurb: "A quartet residency closing out its run with a set of standards and one new commission.",
  },
  {
    id: 5,
    title: "Selene — Dance Reimagined",
    category: "Dance",
    tag: "Contemporary",
    venue: "Marquee Arena, Amaravati",
    date: "2026-10-02",
    time: "6:30 PM",
    priceFrom: 799,
    rating: 4.5,
    accent: "gold",
    blurb: "A contemporary dance piece built around the phases of the moon, performed with a live score.",
  },
  {
    id: 6,
    title: "Open Mic Sundays",
    category: "Comedy",
    tag: "Variety",
    venue: "The Attic Club",
    date: "2026-08-31",
    time: "6:00 PM",
    priceFrom: 149,
    rating: 4.2,
    accent: "red",
    blurb: "Ten new acts, unfiltered. The weekly night where the club's future headliners get found.",
  },
];

export const CATEGORY_META = {
  VIP: { label: "VIP", price: 2499, color: "var(--gold)" },
  Premium: { label: "Premium", price: 1499, color: "var(--teal)" },
  Standard: { label: "Standard", price: 899, color: "var(--paper-dim)" },
};

export const BOOKED_SEATS = new Set([
  "A3", "A4", "B7", "B8", "C1", "C12", "D5", "D6", "E9",
  "F2", "F3", "G10", "G11", "H4", "H5", "H6", "A9", "B2",
]);
export const HELD_SEATS = new Set(["C5", "C6", "D9", "F8"]);

export const ROWS = ["A", "B", "C", "D", "E", "F", "G", "H"];
export const COLS = 12;

export function rowCategory(row) {
  if (row === "A" || row === "B") return "VIP";
  if (row === "C" || row === "D" || row === "E") return "Premium";
  return "Standard";
}

export function buildSeats() {
  const seats = [];
  ROWS.forEach((row) => {
    for (let c = 1; c <= COLS; c++) {
      const id = `${row}${c}`;
      let status = "available";
      if (BOOKED_SEATS.has(id)) status = "booked";
      else if (HELD_SEATS.has(id)) status = "held";
      seats.push({ id, row, col: c, category: rowCategory(row), status });
    }
  });
  return seats;
}

export const PAST_BOOKINGS = [
  {
    id: "bk_9021",
    eventTitle: "Winterlight Sessions",
    venue: "Marquee Arena, Amaravati",
    date: "2026-06-14",
    time: "8:00 PM",
    seatLabels: ["D5", "D6"],
    total: 2998,
    status: "completed",
    code: "WL-9021-XQ",
  },
  {
    id: "bk_8877",
    eventTitle: "Comedy Store: Uncensored",
    venue: "The Attic Club",
    date: "2026-05-02",
    time: "9:00 PM",
    seatLabels: ["F2"],
    total: 349,
    status: "cancelled",
    code: "CS-8877-KT",
  },
  {
    id: "bk_9143",
    eventTitle: "Raga Nights",
    venue: "Blue Room Hall",
    date: "2026-07-20",
    time: "7:30 PM",
    seatLabels: ["B3", "B4"],
    total: 4998,
    status: "completed",
    code: "RN-9143-PL",
  },
];

export const USER = {
  name: "Aarav Shah",
  email: "aarav.shah@example.com",
  phone: "+91 98765 43210",
  memberSince: "March 2025",
};
