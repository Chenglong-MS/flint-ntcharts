import { resolvePivotAggregate } from "flint-chart/core";
import type { ChannelSemantics } from "flint-chart/core";

// Shared reading of a banded (category) axis for the bar family: which
// categories it shows, in what order, under what names, and how several rows
// on one category combine. Terminal bars are always bands, as Flint's Chart.js
// bars are, so this follows that backend: ordinal order, then a sort-by-measure
// request, else first appearance; dates in time order.

type Row = Record<string, unknown>;

/** A category's key: its value as written, so a label and a row agree. */
export const categoryKey = (value: unknown): string => String(value);

const isDateString = (s: string) => /^\d{4}-\d{2}(-\d{2})?([T ]\d{2}:\d{2}(:\d{2}(\.\d+)?)?(Z|[+-]\d{2}:?\d{2})?)?$/.test(s);

/**
 * The categories along a banded axis, in display order. A semantic order
 * (months, weekdays) comes first; dates run in time order; a sort-by-measure
 * request (`sortBy`, `sortOrder`) orders by each category's combined measure;
 * otherwise the order the rows give.
 */
export function categoryOrder(
  table: Row[],
  catCS: ChannelSemantics,
  encoding: { sortBy?: string; sortOrder?: "ascending" | "descending" } | undefined,
  measure: (rows: Row[]) => number,
): string[] {
  const field = catCS.field;
  const seen: string[] = [];
  const rowsOf = new Map<string, Row[]>();
  for (const row of table) {
    if (row[field] == null) continue;
    const key = categoryKey(row[field]);
    if (!rowsOf.has(key)) { rowsOf.set(key, []); seen.push(key); }
    rowsOf.get(key)!.push(row);
  }
  let order = seen;
  const preferred = catCS.ordinalSortOrder;
  if (preferred?.length) {
    const rank = new Map(preferred.map((v, i) => [categoryKey(v), i]));
    order = [...seen].sort((a, b) => (rank.get(a) ?? Infinity) - (rank.get(b) ?? Infinity));
  } else if (seen.length > 1 && seen.every(isDateString)) {
    order = [...seen].sort((a, b) => Date.parse(a) - Date.parse(b));
  }
  if (encoding?.sortBy && encoding.sortOrder) {
    const dir = encoding.sortOrder === "ascending" ? 1 : -1;
    const value = new Map(order.map((k) => [k, measure(rowsOf.get(k) ?? [])]));
    order = [...order].sort((a, b) => dir * ((value.get(a) ?? 0) - (value.get(b) ?? 0)));
  }
  return order;
}

/**
 * How several rows on one category combine into the one bar the terminal
 * draws: the encoding's own `aggregate`, else the measure's semantic default
 * (sum for additive measures, average for intensive ones and for `Number`,
 * whose additivity is unknown). Flint 1.1's rule for exporters that pivot raw
 * rows (`resolvePivotAggregate`).
 */
export function combineOp(valCS: ChannelSemantics, aggregate: string | undefined): "sum" | "average" | "count" {
  return resolvePivotAggregate(aggregate, valCS.aggregationDefault);
}

export function combine(values: number[], op: "sum" | "average" | "count"): number {
  if (op === "count") return values.length;
  const total = values.reduce((a, b) => a + b, 0);
  return op === "average" ? total / values.length : total;
}

/**
 * The name a date category is drawn under: as short as its values allow and
 * still tell them apart — the year where every date is a year's first day, the
 * month where every date is a month's first day, else the day. Times of day
 * keep their minutes. A value that is not a date is drawn as written.
 */
export function categoryLabels(keys: string[]): string[] {
  if (!keys.length || !keys.every(isDateString)) return keys;
  const parts = keys.map((k) => ({ k, date: k.slice(0, 10), time: k.length > 10 ? k.slice(11, 16) : "" }));
  if (parts.some((p) => p.time && p.time !== "00:00")) return keys.map((k) => k.slice(0, 16).replace("T", " "));
  const unique = (labels: string[]) => new Set(labels).size === labels.length;
  const day = (d: string) => (d.length === 7 ? "01" : d.slice(8, 10));
  const firstOfMonth = parts.every((p) => day(p.date) === "01");
  const years = parts.map((p) => p.date.slice(0, 4));
  if (firstOfMonth && parts.every((p) => p.date.slice(5, 7) === "01") && unique(years)) return years;
  const months = parts.map((p) => p.date.slice(0, 7));
  if (firstOfMonth && unique(months)) return months;
  return parts.map((p) => p.date);
}

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

/** A shorter way to name dated bars, and whether it leaves out what the axis needs to be read. */
interface LabelForm { labels: string[]; lossy: boolean }

// Shorter names for date categories, most telling first. A mark for the
// larger period (the year over months, the month over days) sits on the first
// bar and wherever that period starts; the other bars may go unnamed. A day of
// the month only names a bar when the bars run day by day.
function dateForms(keys: string[], labels: string[]): LabelForm[] {
  type D = { y: string; m: string; day: string; t: string };
  const d: D[] = keys.map((k) => ({ y: k.slice(0, 4), m: k.slice(5, 7), day: k.length >= 10 ? k.slice(8, 10) : "01", t: k.length > 10 ? k.slice(11, 16) : "" }));
  const oneYear = new Set(d.map((p) => p.y)).size === 1;
  const startsYear = (i: number) => i === 0 || d[i].y !== d[i - 1].y;
  const startsMonth = (i: number) => startsYear(i) || d[i].m !== d[i - 1].m;
  const startsDay = (i: number) => i === 0 || keys[i].slice(0, 10) !== keys[i - 1].slice(0, 10);
  const time = keys.map((k) => Date.parse(k.slice(0, 10)));
  const daily = time.every((t, i) => i === 0 || t - time[i - 1] <= 3 * 86_400_000);
  const month = (p: D) => MONTHS[Number(p.m) - 1] ?? p.m;
  const form = (lossy: boolean, f: (p: D, i: number) => string): LabelForm => ({ labels: d.map(f), lossy });
  const yearMarks = [
    form(true, (p, i) => (startsYear(i) ? `'${p.y.slice(2)}` : "")),
    form(true, (p, i) => (startsYear(i) ? p.y.slice(2) : "")),
  ];
  switch (labels[0]?.length) {
    case 4: // years
      return [form(false, (p) => `'${p.y.slice(2)}`), form(false, (p) => p.y.slice(2))];
    case 7: // months
      return [
        form(false, (p, i) => (startsYear(i) ? p.y : month(p))),
        form(false, (p, i) => (startsYear(i) ? `'${p.y.slice(2)}` : month(p))),
        form(!oneYear, (p) => p.m),
        ...(oneYear ? [form(false, (p) => month(p)[0])] : yearMarks),
      ];
    case 10: // days
      return [
        form(!oneYear, (p) => `${p.m}-${p.day}`),
        ...(daily ? [form(!oneYear, (p, i) => (startsMonth(i) ? month(p) : p.day))] : []),
        ...(oneYear ? [form(true, (p, i) => (startsMonth(i) ? month(p) : ""))] : yearMarks),
        ...(daily ? [form(true, (p) => p.day)] : []),
      ];
    default: // times of day
      return [
        form(!oneYear, (p) => `${p.m}-${p.day} ${p.t}`),
        form(true, (p, i) => (startsDay(i) ? `${p.m}-${p.day}` : p.t)),
        form(true, (p) => p.t),
        form(true, (p) => p.t.slice(0, 2)),
      ];
  }
}

/**
 * Names for vertical bars `room` columns wide. The renderer cuts a longer
 * name to the bar's width, which leaves a name's prefix (fine for words) but
 * turns `2021-03` into `20` and `13` into `1`. So a date steps down to the
 * first shorter form that fits, and a number that does not fit is left out
 * rather than cut. `lossy` says the axis now needs its range stated.
 */
export function fitBandLabels(keys: string[], labels: string[], room: number): LabelForm {
  const fits = (ls: string[]) => ls.every((l) => l.length <= room);
  if (fits(labels)) return { labels, lossy: false };
  if (keys.every(isDateString)) {
    return dateForms(keys, labels).find((f) => fits(f.labels)) ?? { labels: keys.map(() => ""), lossy: true };
  }
  if (keys.every((k) => /^\d{1,2}:\d{2}(:\d{2})?$/.test(k))) {
    const onTheHour = keys.every((k) => /:00(:00)?$/.test(k));
    const hours = { labels: keys.map((k) => k.split(":")[0].padStart(2, "0")), lossy: !onTheHour };
    return [{ labels: keys.map((k) => k.slice(0, 5)), lossy: false }, hours].find((f) => fits(f.labels)) ?? { labels: keys.map(() => ""), lossy: true };
  }
  if (keys.every((k) => k.trim() !== "" && Number.isFinite(Number(k)))) {
    return { labels: labels.map((l) => (l.length <= room ? l : "")), lossy: true };
  }
  return { labels, lossy: false };
}

/** The width ntcharts gives each of `n` vertical bars on a `width`-column canvas (one-column gaps). */
export const barRoom = (width: number, n: number) => Math.max(1, Math.floor((width - (n - 1)) / Math.max(1, n)));
