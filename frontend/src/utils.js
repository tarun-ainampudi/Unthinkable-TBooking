export function formatDate(d) {
  const date = new Date(d + "T00:00:00");
  return date.toLocaleDateString("en-IN", { day: "numeric", month: "short", year: "numeric" });
}

export function formatINR(n) {
  return "₹" + n.toLocaleString("en-IN");
}

export function pad(n) {
  return n.toString().padStart(2, "0");
}
