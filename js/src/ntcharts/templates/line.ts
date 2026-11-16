import { isCategorical, finishCategorical, categorySeries } from "../category.js";
import type { ChartTemplateDef } from "flint-chart/core";
import { makeCartesianPivot, makeSortAction } from "flint-chart/core";
import { splitSeries } from "../series.js";
import { toMs } from "../temporal.js";
import type { NtInstantiateContext } from "../assemble.js";
import { pinFittedYDomain } from "../domain.js";
import { dateSlots } from "../cadence.js";

// Columns the Y axis labels and line take from the plot's width.
const Y_LABEL_ROOM = 9;

export const ntLineChartDef: ChartTemplateDef = {
  chart: "Line Chart",
  template: { mark: "line" },
  channels: ["x", "y", "color", "group"],
  markCognitiveChannel: "position",
  instantiate: (_spec: any, rawCtx: any) => {
    let ctx = rawCtx as NtInstantiateContext;
    // Regular dates with room for every label are drawn on category slots, so
    // each label sits under its point (see dateSlots).
    const x0 = ctx.channelSemantics.x;
    const slots = x0?.type === "temporal"
      ? dateSlots(ctx.table.map((r) => r[x0.field]), x0.temporalFormat, ctx.emit.width - Y_LABEL_ROOM)
      : undefined;
    if (slots) {
      const order = [...slots.keys()].sort((a, b) => toMs(a) - toMs(b));
      ctx = { ...ctx, channelSemantics: { ...ctx.channelSemantics, x: { ...x0!, type: "ordinal", ordinalSortOrder: order } } };
    }
    const { emit, channelSemantics: cs, table } = ctx;
    if (cs.y?.type !== "quantitative") {
      throw new Error(`A terminal ${ctx.chartType} measures along y; got y = ${cs.y ? `${cs.y.field} (${cs.y.type})` : "none"}. ` +
        `Put the measure on y and the ${cs.y?.field ?? "category"} on x.`);
    }
    const temporal = cs.x?.type === "temporal";
    emit.type = temporal ? "timeseries" : "line";
    emit.x_axis = { ...emit.x_axis, type: temporal ? "time" : "value", title: cs.x?.field };
    emit.y_axis = { ...emit.y_axis, title: cs.y?.field };
    if (ctx.ntFormatX) emit.x_axis.format = ctx.ntFormatX;
    if (ctx.ntFormatY) emit.y_axis.format = ctx.ntFormatY;
    if (cs.y?.zero?.zero) emit.y_axis.min = 0;
    else pinFittedYDomain(emit, cs, table);

    if (isCategorical(ctx)) {
      finishCategorical(ctx, categorySeries(ctx), true);
      if (slots) {
        emit.x_axis!.labels = emit.x_axis!.labels!.map((k) => slots.get(k) ?? k);
        for (const s of emit.data.series) for (const p of s.values ?? []) p.x = slots.get(String(p.x)) ?? p.x;
      }
      return;
    }

    const xField = cs.x!.field;
    emit.data.series = splitSeries(table, cs, (row) => ({
      x: temporal ? toMs(row[xField]) : Number(row[xField]),
      y: Number(row[cs.y!.field]),
    }));

    // ntcharts' timeserieslinechart/linechart assume chronological push
    // order (pure append, segments drawn by insertion index), but our
    // emitted series preserve raw row order from the source data. Sort each
    // series ascending by x so the drawn segments are correct regardless of
    // input ordering. Also drop any points whose x or y failed to parse
    // (e.g. an unparsable temporal cell yields NaN from toMs) rather than
    // silently plotting a fake epoch-0/NaN point; surface a single warning
    // if anything was dropped.
    let droppedCount = 0;
    for (const series of emit.data.series) {
      const values = series.values ?? [];
      const kept = values.filter((p) => {
        const ok = !Number.isNaN(p.x as number) && !Number.isNaN(p.y as number);
        if (!ok) droppedCount++;
        return ok;
      });
      kept.sort((a, b) => (a.x as number) - (b.x as number));
      series.values = kept;
    }
    if (droppedCount > 0) {
      ctx.warn({
        severity: "warning",
        code: "invalid-temporal-x",
        message: `dropped ${droppedCount} point(s) with unparsable X values`,
      });
    }

    if (emit.data.series.length > 1) emit.options = { ...emit.options, show_legend: true };
  },
  encodingActions: [makeSortAction()],
  pivot: makeCartesianPivot({ permute: [["x", "y", "color"]], shift: ["color", "group"] }),
};
