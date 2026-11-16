import { toMs } from "./temporal.js";

const MONTHS = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
const DAYS = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
const pad = (n: number, w = 2) => String(n).padStart(w, "0");

/**
 * A UTC time formatted by a d3 time-format string (the common specifiers
 * flint emits); an unknown specifier yields undefined.
 */
export function formatTime(ms: number, d3: string): string | undefined {
  const t = new Date(ms);
  const h = t.getUTCHours();
  const parts: Record<string, string> = {
    Y: String(t.getUTCFullYear()), y: pad(t.getUTCFullYear() % 100), m: pad(t.getUTCMonth() + 1),
    b: MONTHS[t.getUTCMonth()].slice(0, 3), B: MONTHS[t.getUTCMonth()], d: pad(t.getUTCDate()), e: String(t.getUTCDate()).padStart(2, " "),
    a: DAYS[t.getUTCDay()].slice(0, 3), A: DAYS[t.getUTCDay()], H: pad(h), I: pad(h % 12 || 12), M: pad(t.getUTCMinutes()),
    S: pad(t.getUTCSeconds()), p: h < 12 ? "AM" : "PM", q: String(Math.floor(t.getUTCMonth() / 3) + 1), "%": "%",
  };
  let unknown = false;
  const out = d3.replace(/%[-_0]?([a-zA-Z%])/g, (_, c: string) => parts[c] ?? ((unknown = true), ""));
  return unknown ? undefined : out.trim();
}

/**
 * The dates of a series that runs at a regular calendar cadence without gaps —
 * year starts a year apart, month starts one or three months apart, or a fixed
 * step (days, weeks, hours) — sorted; undefined otherwise.
 */
export function regularDates(values: unknown[]): number[] | undefined {
  const ms = [...new Set(values.map(toMs).filter(Number.isFinite))].sort((a, b) => a - b);
  if (ms.length < 2) return undefined;
  const d = ms.map((v) => new Date(v));
  const monthIndex = d.map((t) => t.getUTCFullYear() * 12 + t.getUTCMonth());
  const monthStarts = d.every((t) => t.getUTCDate() === 1 && t.getUTCHours() === 0 && t.getUTCMinutes() === 0);
  if (monthStarts) {
    const step = monthIndex[1] - monthIndex[0];
    if (step > 0 && monthIndex.every((m, i) => i === 0 || m - monthIndex[i - 1] === step)) return ms;
  }
  const step = ms[1] - ms[0];
  return ms.every((v, i) => i === 0 || v - ms[i - 1] === step) ? ms : undefined;
}

/**
 * Category labels for a regular date series when each fits under its point
 * in `width` columns: one label per date in flint's format. ntcharts places
 * time-axis labels on a fixed 4-column step, so a "Feb" can sit up to three
 * columns past Feb 1; category slots put each label under its point.
 */
export function dateSlots(values: unknown[], d3: string | undefined, width: number): Map<string, string> | undefined {
  if (!values.every((v) => v == null || Number.isFinite(toMs(v)))) return undefined;
  const dates = regularDates(values.filter((v) => v != null));
  if (!dates || !d3) return undefined;
  const labels = dates.map((ms) => formatTime(ms, d3));
  if (labels.some((l) => l == null) || new Set(labels).size !== labels.length) return undefined;
  const widest = Math.max(...labels.map((l) => l!.length));
  if (dates.length * (widest + 1) > width) return undefined;
  const byKey = new Map<string, string>();
  for (const v of values) {
    const ms = toMs(v);
    const i = dates.indexOf(ms);
    if (i >= 0) byKey.set(String(v), labels[i]!);
  }
  return byKey;
}
