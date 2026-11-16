import type { ChartTemplateDef } from "flint-chart/core";
import { detectBandedAxisForceDiscrete, makeCartesianPivot, makeSortAction } from "flint-chart/core";
import { paletteForScheme } from "../colormap.js";
import { barRoom, categoryKey, categoryLabels, categoryOrder, combine, combineOp, fitBandLabels } from "../bands.js";
import type { NtInstantiateContext } from "../assemble.js";

type Row = Record<string, unknown>;

export const ntBarChartDef: ChartTemplateDef = {
  chart: "Bar Chart",
  template: { mark: "bar" },
  channels: ["x", "y", "color", "group"],
  markCognitiveChannel: "length",
  declareLayoutMode: (cs, table) => {
    // Terminal bars are always bands, as Flint's Chart.js bars are: a date or
    // number axis is laid out as categories, so flint budgets and truncates it
    // (with its usual overflow warning) rather than handing the renderer more
    // bars than it has cells for.
    const result = detectBandedAxisForceDiscrete(cs, table, { preferAxis: "x" });
    return {
      axisFlags: result ? { [result.axis]: { banded: true } } : { x: { banded: true } },
      resolvedTypes: result?.resolvedTypes,
      // The ntcharts bar model spends two cells per category — the bar and
      // the gap after it — and draws nothing at all once the bar width
      // reaches zero. A 2-cell step makes flint budget to what it can draw.
      paramOverrides: { minStep: 2 },
    };
  },
  instantiate: (_spec: any, rawCtx: any) => {
    const ctx = rawCtx as NtInstantiateContext;
    const { emit, channelSemantics: cs, table, encodings } = ctx;
    // The measured axis is the quantitative one; with both quantitative, y.
    const horizontal = cs.x?.type === "quantitative" && cs.y != null && cs.y.type !== "quantitative";
    const catAxis = horizontal ? "y" : "x";
    const valAxis = horizontal ? "x" : "y";
    const catCS = cs[catAxis];
    const valCS = cs[valAxis];
    if (!catCS || !valCS || valCS.type !== "quantitative") {
      const shown = (["x", "y"] as const).map((a) => cs[a] ? `${a} = ${cs[a].field} (${cs[a].type})` : `${a} = none`).join(", ");
      throw new Error(`A terminal ${ctx.chartType} needs one quantitative axis to measure; got ${shown}. ` +
        `Use a Heatmap for two categorical axes.`);
    }

    // A facet becomes one more dimension the single terminal plot cannot show:
    // its rows are combined into each bar, and the chart says so.
    const facets = [cs.column, cs.row].filter(Boolean).map((f) => f!.field);
    const groupCS = cs.color?.type === "quantitative" ? undefined : (cs.color ?? cs.group);
    const op = combineOp(valCS, encodings[valAxis]?.aggregate);
    const values = (rows: Row[]) => rows.map((r) => Number(r[valCS.field])).filter(Number.isFinite);
    const measure = (rows: Row[]) => (rows.length ? combine(values(rows), op) : 0);

    emit.type = "bar";
    if (horizontal) emit.options = { ...emit.options, orientation: "horizontal" };
    const keys = categoryOrder(table, catCS, encodings[catAxis], measure);
    let labels = categoryLabels(keys);
    if (!horizontal && emit.width) {
      const fitted = fitBandLabels(keys, labels, barRoom(emit.width, keys.length));
      if (fitted.lossy) {
        ctx.warn({ severity: "info", code: "labels-shortened",
          message: `bars are too narrow to name every ${catCS.field}: the axis runs ${labels[0]} to ${labels[labels.length - 1]}`,
          channel: catAxis, field: catCS.field } as any);
      }
      labels = fitted.labels;
    }
    emit.x_axis = { ...emit.x_axis, type: "category", labels, title: catCS.field };
    emit.y_axis = { ...emit.y_axis, title: valCS.field };
    const valueFormat = horizontal ? ctx.ntFormatX : ctx.ntFormatY;
    if (valueFormat) emit.y_axis.format = valueFormat;

    // One series per colour (or group) value, each aligned to the category
    // order; several rows on one bar combine by the measure's rule.
    const seriesKeys: string[] = [];
    const cells = new Map<string, Map<string, Row[]>>();
    const preferred = groupCS?.ordinalSortOrder?.map(String);
    if (preferred) for (const g of preferred) { seriesKeys.push(g); cells.set(g, new Map()); }
    for (const row of table) {
      if (row[catCS.field] == null) continue;
      const g = groupCS ? String(row[groupCS.field]) : valCS.field;
      if (!cells.has(g)) { seriesKeys.push(g); cells.set(g, new Map()); }
      const k = categoryKey(row[catCS.field]);
      const bucket = cells.get(g)!;
      (bucket.get(k) ?? bucket.set(k, []).get(k)!).push(row);
    }
    const palette = paletteForScheme(groupCS?.colorScheme?.scheme);
    let combined = false;
    const series = seriesKeys.map((name, i) => {
      const bucket = cells.get(name)!;
      for (const rows of bucket.values()) if (rows.length > 1) combined = true;
      return {
        name,
        values: keys.map((k) => ({ y: measure(bucket.get(k) ?? []) })),
        ...(groupCS ? { color: palette[i % palette.length] } : {}),
      };
    }).filter((s) => cells.get(s.name)!.size > 0);
    emit.data.series = series;
    if (series.length > 1) emit.options = { ...emit.options, stacked: true, show_legend: true };

    // The bar model always draws from zero; pinning zero as the minimum would
    // only hide bars that go below it.
    const lowest = Math.min(0, ...series.flatMap((s) => s.values.map((p) => p.y)));
    if (valCS.zero?.zero && lowest >= 0) emit.y_axis.min = 0;

    if (facets.length) {
      ctx.warn({ severity: "info", code: "facets-merged",
        message: `the terminal draws one plot: each bar ${op === "average" ? "averages" : "combines"} its rows across ${facets.join(" and ")} (${op})`,
        channel: "column", field: facets[0] } as any);
    } else if (combined && !encodings[valAxis]?.aggregate) {
      ctx.warn({ severity: "info", code: "rows-combined",
        message: `several rows share a ${catCS.field}${groupCS ? ` and ${groupCS.field}` : ""}: each bar shows their ${op}`,
        channel: valAxis, field: valCS.field } as any);
    }
  },
  encodingActions: [makeSortAction()],
  pivot: makeCartesianPivot({
    transpose: [["x", "y"]],
    permute: [["x", "y", "color"]],
    shift: ["color", "group"],
  }),
};
