package editor

// example is a built-in demo spec, loadable with ctrl+N / alt+N where N is
// its 1-based position in examples.
type example struct {
	name string // shown on the status line
	src  string // full flint ChartAssemblyInput JSON
}

// examples is ordered: index 0 loads with key 1, and seeds the editor at
// startup. Sources are compiled in (adapted from testdata/fixtures-terminal,
// not embedded — the demo must not drift when test fixtures change) and must
// compile with at most info notes; TestExamplesCompile enforces that.
var examples = []example{
	{name: "Bar Chart", src: `{
  "data": { "values": [
    {"month": "Jan", "revenue": 120},
    {"month": "Feb", "revenue": 180},
    {"month": "Mar", "revenue": 95},
    {"month": "Apr", "revenue": 210},
    {"month": "May", "revenue": 165}
  ]},
  "semantic_types": { "revenue": {"semanticType": "Price", "unit": "USD"} },
  "chart_spec": {
    "chartType": "Bar Chart",
    "encodings": { "x": {"field": "month"}, "y": {"field": "revenue"} }
  }
}`},
	{name: "Timeseries Line", src: `{
  "data": { "values": [
    {"date": "2026-01-01", "price": 104.2},
    {"date": "2026-02-01", "price": 108.9},
    {"date": "2026-03-01", "price": 101.4},
    {"date": "2026-04-01", "price": 115.7},
    {"date": "2026-05-01", "price": 119.3},
    {"date": "2026-06-01", "price": 112.8},
    {"date": "2026-07-01", "price": 121.5}
  ]},
  "semantic_types": { "price": {"semanticType": "Price", "unit": "USD"} },
  "chart_spec": {
    "chartType": "Line Chart",
    "chartProperties": { "includeZero_y": false },
    "encodings": { "x": {"field": "date"}, "y": {"field": "price"} }
  }
}`},
	{name: "Candlestick", src: `{
  "data": { "values": [
    {"date": "2026-01-05", "open": 100, "high": 108, "low": 97,  "close": 105},
    {"date": "2026-01-06", "open": 105, "high": 112, "low": 103, "close": 110},
    {"date": "2026-01-07", "open": 110, "high": 111, "low": 98,  "close": 99},
    {"date": "2026-01-08", "open": 99,  "high": 106, "low": 96,  "close": 104},
    {"date": "2026-01-09", "open": 104, "high": 115, "low": 102, "close": 114},
    {"date": "2026-01-12", "open": 114, "high": 118, "low": 109, "close": 111},
    {"date": "2026-01-13", "open": 111, "high": 121, "low": 110, "close": 120}
  ]},
  "chart_spec": {
    "chartType": "Candlestick Chart",
    "chartProperties": { "candleStyle": "block" },
    "encodings": {
      "x": {"field": "date"},
      "open": {"field": "open"}, "high": {"field": "high"},
      "low": {"field": "low"}, "close": {"field": "close"}
    }
  }
}`},
	{name: "Log Scale", src: `{
  "data": { "values": [
    {"month": "2026-01-01", "users": 120},
    {"month": "2026-02-01", "users": 310},
    {"month": "2026-03-01", "users": 900},
    {"month": "2026-04-01", "users": 2400},
    {"month": "2026-05-01", "users": 7100},
    {"month": "2026-06-01", "users": 19500},
    {"month": "2026-07-01", "users": 61000},
    {"month": "2026-08-01", "users": 170000}
  ]},
  "chart_spec": {
    "chartType": "Line Chart",
    "chartProperties": { "logScale_y": true },
    "encodings": { "x": {"field": "month"}, "y": {"field": "users"} }
  }
}`},
	{name: "Heatmap", src: `{
  "data": { "values": [
    {"day": "Mon", "when": "Morning",   "visits": 62},
    {"day": "Mon", "when": "Afternoon", "visits": 91},
    {"day": "Mon", "when": "Evening",   "visits": 35},
    {"day": "Tue", "when": "Morning",   "visits": 58},
    {"day": "Tue", "when": "Afternoon", "visits": 97},
    {"day": "Tue", "when": "Evening",   "visits": 41},
    {"day": "Wed", "when": "Morning",   "visits": 66},
    {"day": "Wed", "when": "Afternoon", "visits": 84},
    {"day": "Wed", "when": "Evening",   "visits": 52},
    {"day": "Thu", "when": "Morning",   "visits": 49},
    {"day": "Thu", "when": "Afternoon", "visits": 88},
    {"day": "Thu", "when": "Evening",   "visits": 60},
    {"day": "Fri", "when": "Morning",   "visits": 44},
    {"day": "Fri", "when": "Afternoon", "visits": 73},
    {"day": "Fri", "when": "Evening",   "visits": 79}
  ]},
  "chart_spec": {
    "chartType": "Heatmap",
    "encodings": { "x": {"field": "when"}, "y": {"field": "day"}, "color": {"field": "visits"} }
  }
}`},
	{name: "Histogram", src: `{
  "data": { "values": [
    {"latency_ms": 113}, {"latency_ms": 134}, {"latency_ms": 114}, {"latency_ms": 111}, {"latency_ms": 94},
    {"latency_ms": 114}, {"latency_ms": 151}, {"latency_ms": 132}, {"latency_ms": 149}, {"latency_ms": 127},
    {"latency_ms": 131}, {"latency_ms": 125}, {"latency_ms": 73}, {"latency_ms": 144}, {"latency_ms": 134},
    {"latency_ms": 134}, {"latency_ms": 73}, {"latency_ms": 71}, {"latency_ms": 95}, {"latency_ms": 107},
    {"latency_ms": 129}, {"latency_ms": 119}, {"latency_ms": 135}, {"latency_ms": 102}, {"latency_ms": 129},
    {"latency_ms": 131}, {"latency_ms": 101}, {"latency_ms": 168}, {"latency_ms": 136}, {"latency_ms": 154},
    {"latency_ms": 103}, {"latency_ms": 99}, {"latency_ms": 110}, {"latency_ms": 117}, {"latency_ms": 138},
    {"latency_ms": 127}, {"latency_ms": 107}, {"latency_ms": 93}, {"latency_ms": 105}, {"latency_ms": 154},
    {"latency_ms": 240}, {"latency_ms": 266}, {"latency_ms": 271}, {"latency_ms": 223}, {"latency_ms": 261},
    {"latency_ms": 293}, {"latency_ms": 210}, {"latency_ms": 252}, {"latency_ms": 257}, {"latency_ms": 240}
  ]},
  "chart_spec": {
    "chartType": "Histogram",
    "chartProperties": { "binCount": 12 },
    "encodings": { "x": {"field": "latency_ms"} }
  }
}`},
	{name: "Calendar Heatmap", src: `{
  "data": { "values": [
    {"day": "2026-03-02", "commits": 5}, {"day": "2026-03-03", "commits": 3}, {"day": "2026-03-04", "commits": 0}, {"day": "2026-03-05", "commits": 6}, {"day": "2026-03-06", "commits": 6}, {"day": "2026-03-07", "commits": 2}, {"day": "2026-03-08", "commits": 3},
    {"day": "2026-03-09", "commits": 5}, {"day": "2026-03-10", "commits": 4}, {"day": "2026-03-11", "commits": 0}, {"day": "2026-03-12", "commits": 5}, {"day": "2026-03-13", "commits": 2}, {"day": "2026-03-14", "commits": 0}, {"day": "2026-03-15", "commits": 0},
    {"day": "2026-03-16", "commits": 1}, {"day": "2026-03-17", "commits": 2}, {"day": "2026-03-18", "commits": 7}, {"day": "2026-03-19", "commits": 0}, {"day": "2026-03-20", "commits": 0}, {"day": "2026-03-21", "commits": 1}, {"day": "2026-03-22", "commits": 3},
    {"day": "2026-03-23", "commits": 0}, {"day": "2026-03-24", "commits": 0}, {"day": "2026-03-25", "commits": 0}, {"day": "2026-03-26", "commits": 0}, {"day": "2026-03-27", "commits": 0}, {"day": "2026-03-28", "commits": 0}, {"day": "2026-03-29", "commits": 0},
    {"day": "2026-03-30", "commits": 7}, {"day": "2026-03-31", "commits": 4}, {"day": "2026-04-01", "commits": 4}, {"day": "2026-04-02", "commits": 5}, {"day": "2026-04-03", "commits": 8}, {"day": "2026-04-04", "commits": 1}, {"day": "2026-04-05", "commits": 1},
    {"day": "2026-04-06", "commits": 5}, {"day": "2026-04-07", "commits": 0}, {"day": "2026-04-08", "commits": 7}, {"day": "2026-04-09", "commits": 6}, {"day": "2026-04-10", "commits": 5}, {"day": "2026-04-11", "commits": 0}, {"day": "2026-04-12", "commits": 0},
    {"day": "2026-04-13", "commits": 6}, {"day": "2026-04-14", "commits": 0}, {"day": "2026-04-15", "commits": 3}, {"day": "2026-04-16", "commits": 7}, {"day": "2026-04-17", "commits": 0}, {"day": "2026-04-18", "commits": 3}, {"day": "2026-04-19", "commits": 1},
    {"day": "2026-04-20", "commits": 3}, {"day": "2026-04-21", "commits": 4}, {"day": "2026-04-22", "commits": 5}, {"day": "2026-04-23", "commits": 4}, {"day": "2026-04-24", "commits": 7}, {"day": "2026-04-25", "commits": 0}, {"day": "2026-04-26", "commits": 0}
  ]},
  "chart_spec": {
    "chartType": "Calendar Heatmap",
    "encodings": { "x": {"field": "day"}, "color": {"field": "commits"} }
  }
}`},
	{name: "Pie Chart (raster)", src: `{
  "renderer": "raster",
  "data": { "values": [
    {"device": "Laptop", "units": 1190},
    {"device": "Phone",  "units": 2110},
    {"device": "Tablet", "units": 780},
    {"device": "Watch",  "units": 420}
  ]},
  "chart_spec": {
    "chartType": "Pie Chart",
    "encodings": {
      "color": {"field": "device"},
      "size": {"field": "units"}
    }
  }
}`},
	{name: "Combo Chart", src: `{
  "data": {
    "values": [
      {
        "month": "Jan",
        "revenue": 120,
        "margin": 0.2
      },
      {
        "month": "Feb",
        "revenue": 180,
        "margin": 0.25
      },
      {
        "month": "Mar",
        "revenue": 150,
        "margin": 0.22
      },
      {
        "month": "Apr",
        "revenue": 220,
        "margin": 0.3
      },
      {
        "month": "May",
        "revenue": 190,
        "margin": 0.26
      },
      {
        "month": "Jun",
        "revenue": 260,
        "margin": 0.34
      }
    ]
  },
  "semantic_types": {
    "revenue": {
      "semanticType": "Amount",
      "unit": "USD"
    },
    "margin": {
      "semanticType": "Percentage",
      "intrinsicDomain": [
        0,
        1
      ]
    }
  },
  "chart_spec": {
    "chartType": "Combo Chart",
    "title": "Revenue and margin",
    "encodings": {
      "x": {
        "field": "month"
      },
      "y": {
        "field": "revenue"
      }
    },
    "chartProperties": {
      "lineField": "margin"
    }
  }
}`},
	{name: "Grouped Bars", src: `{
  "data": {
    "values": [
      {
        "device": "Laptop",
        "region": "USA",
        "units": 410
      },
      {
        "device": "Laptop",
        "region": "China",
        "units": 530
      },
      {
        "device": "Laptop",
        "region": "Japan",
        "units": 250
      },
      {
        "device": "Phone",
        "region": "USA",
        "units": 720
      },
      {
        "device": "Phone",
        "region": "China",
        "units": 910
      },
      {
        "device": "Phone",
        "region": "Japan",
        "units": 480
      },
      {
        "device": "Tablet",
        "region": "USA",
        "units": 190
      },
      {
        "device": "Tablet",
        "region": "China",
        "units": 260
      },
      {
        "device": "Tablet",
        "region": "Japan",
        "units": 330
      },
      {
        "device": "Watch",
        "region": "USA",
        "units": 140
      },
      {
        "device": "Watch",
        "region": "China",
        "units": 90
      },
      {
        "device": "Watch",
        "region": "Japan",
        "units": 120
      }
    ]
  },
  "semantic_types": {
    "units": {
      "semanticType": "Count"
    }
  },
  "chart_spec": {
    "chartType": "Grouped Bar Chart",
    "encodings": {
      "x": {
        "field": "device"
      },
      "y": {
        "field": "units"
      },
      "group": {
        "field": "region"
      }
    }
  }
}`},
	{name: "Pyramid", src: `{
  "data": {
    "values": [
      {
        "age": "0-14",
        "sex": "Male",
        "millions": 31
      },
      {
        "age": "0-14",
        "sex": "Female",
        "millions": 30
      },
      {
        "age": "15-29",
        "sex": "Male",
        "millions": 33
      },
      {
        "age": "15-29",
        "sex": "Female",
        "millions": 32
      },
      {
        "age": "30-44",
        "sex": "Male",
        "millions": 32
      },
      {
        "age": "30-44",
        "sex": "Female",
        "millions": 32
      },
      {
        "age": "45-59",
        "sex": "Male",
        "millions": 30
      },
      {
        "age": "45-59",
        "sex": "Female",
        "millions": 31
      },
      {
        "age": "60-74",
        "sex": "Male",
        "millions": 22
      },
      {
        "age": "60-74",
        "sex": "Female",
        "millions": 25
      },
      {
        "age": "75+",
        "sex": "Male",
        "millions": 9
      },
      {
        "age": "75+",
        "sex": "Female",
        "millions": 13
      }
    ]
  },
  "semantic_types": {
    "millions": {
      "semanticType": "Quantity"
    }
  },
  "chart_spec": {
    "chartType": "Pyramid Chart",
    "encodings": {
      "y": {
        "field": "age"
      },
      "x": {
        "field": "millions"
      },
      "color": {
        "field": "sex"
      }
    }
  }
}`},
	{name: "Density Plot", src: `{
  "data": {
    "values": [
      {
        "team": "Web",
        "latency_ms": 210
      },
      {
        "team": "Web",
        "latency_ms": 240
      },
      {
        "team": "Web",
        "latency_ms": 211
      },
      {
        "team": "Web",
        "latency_ms": 207
      },
      {
        "team": "Web",
        "latency_ms": 183
      },
      {
        "team": "Web",
        "latency_ms": 211
      },
      {
        "team": "Web",
        "latency_ms": 264
      },
      {
        "team": "Web",
        "latency_ms": 237
      },
      {
        "team": "Web",
        "latency_ms": 261
      },
      {
        "team": "Web",
        "latency_ms": 230
      },
      {
        "team": "Web",
        "latency_ms": 236
      },
      {
        "team": "Web",
        "latency_ms": 227
      },
      {
        "team": "Web",
        "latency_ms": 153
      },
      {
        "team": "Web",
        "latency_ms": 254
      },
      {
        "team": "Web",
        "latency_ms": 240
      },
      {
        "team": "Web",
        "latency_ms": 240
      },
      {
        "team": "Web",
        "latency_ms": 152
      },
      {
        "team": "Web",
        "latency_ms": 150
      },
      {
        "team": "Web",
        "latency_ms": 184
      },
      {
        "team": "Web",
        "latency_ms": 201
      },
      {
        "team": "Web",
        "latency_ms": 232
      },
      {
        "team": "Web",
        "latency_ms": 218
      },
      {
        "team": "Web",
        "latency_ms": 241
      },
      {
        "team": "Web",
        "latency_ms": 194
      },
      {
        "team": "Web",
        "latency_ms": 232
      },
      {
        "team": "Web",
        "latency_ms": 236
      },
      {
        "team": "Web",
        "latency_ms": 194
      },
      {
        "team": "Web",
        "latency_ms": 289
      },
      {
        "team": "Web",
        "latency_ms": 242
      },
      {
        "team": "Web",
        "latency_ms": 268
      },
      {
        "team": "Web",
        "latency_ms": 195
      },
      {
        "team": "Web",
        "latency_ms": 190
      },
      {
        "team": "Web",
        "latency_ms": 206
      },
      {
        "team": "Web",
        "latency_ms": 216
      },
      {
        "team": "Web",
        "latency_ms": 245
      },
      {
        "team": "Web",
        "latency_ms": 230
      },
      {
        "team": "Web",
        "latency_ms": 202
      },
      {
        "team": "Web",
        "latency_ms": 182
      },
      {
        "team": "Web",
        "latency_ms": 199
      },
      {
        "team": "Web",
        "latency_ms": 269
      },
      {
        "team": "Web",
        "latency_ms": 188
      },
      {
        "team": "Web",
        "latency_ms": 230
      },
      {
        "team": "Web",
        "latency_ms": 237
      },
      {
        "team": "Web",
        "latency_ms": 160
      },
      {
        "team": "Web",
        "latency_ms": 222
      },
      {
        "team": "Web",
        "latency_ms": 272
      },
      {
        "team": "Web",
        "latency_ms": 139
      },
      {
        "team": "Web",
        "latency_ms": 207
      },
      {
        "team": "Web",
        "latency_ms": 216
      },
      {
        "team": "Web",
        "latency_ms": 187
      },
      {
        "team": "Web",
        "latency_ms": 240
      },
      {
        "team": "Web",
        "latency_ms": 218
      },
      {
        "team": "Web",
        "latency_ms": 161
      },
      {
        "team": "Web",
        "latency_ms": 253
      },
      {
        "team": "Web",
        "latency_ms": 247
      },
      {
        "team": "Web",
        "latency_ms": 258
      },
      {
        "team": "Web",
        "latency_ms": 278
      },
      {
        "team": "Web",
        "latency_ms": 234
      },
      {
        "team": "Web",
        "latency_ms": 225
      },
      {
        "team": "Web",
        "latency_ms": 168
      },
      {
        "team": "Mobile",
        "latency_ms": 347
      },
      {
        "team": "Mobile",
        "latency_ms": 273
      },
      {
        "team": "Mobile",
        "latency_ms": 283
      },
      {
        "team": "Mobile",
        "latency_ms": 234
      },
      {
        "team": "Mobile",
        "latency_ms": 252
      },
      {
        "team": "Mobile",
        "latency_ms": 278
      },
      {
        "team": "Mobile",
        "latency_ms": 387
      },
      {
        "team": "Mobile",
        "latency_ms": 188
      },
      {
        "team": "Mobile",
        "latency_ms": 223
      },
      {
        "team": "Mobile",
        "latency_ms": 324
      },
      {
        "team": "Mobile",
        "latency_ms": 397
      },
      {
        "team": "Mobile",
        "latency_ms": 345
      },
      {
        "team": "Mobile",
        "latency_ms": 196
      },
      {
        "team": "Mobile",
        "latency_ms": 159
      },
      {
        "team": "Mobile",
        "latency_ms": 331
      },
      {
        "team": "Mobile",
        "latency_ms": 266
      },
      {
        "team": "Mobile",
        "latency_ms": 243
      },
      {
        "team": "Mobile",
        "latency_ms": 369
      },
      {
        "team": "Mobile",
        "latency_ms": 376
      },
      {
        "team": "Mobile",
        "latency_ms": 319
      },
      {
        "team": "Mobile",
        "latency_ms": 325
      },
      {
        "team": "Mobile",
        "latency_ms": 336
      },
      {
        "team": "Mobile",
        "latency_ms": 406
      },
      {
        "team": "Mobile",
        "latency_ms": 347
      },
      {
        "team": "Mobile",
        "latency_ms": 341
      },
      {
        "team": "Mobile",
        "latency_ms": 343
      },
      {
        "team": "Mobile",
        "latency_ms": 216
      },
      {
        "team": "Mobile",
        "latency_ms": 387
      },
      {
        "team": "Mobile",
        "latency_ms": 367
      },
      {
        "team": "Mobile",
        "latency_ms": 342
      },
      {
        "team": "Mobile",
        "latency_ms": 192
      },
      {
        "team": "Mobile",
        "latency_ms": 272
      },
      {
        "team": "Mobile",
        "latency_ms": 361
      },
      {
        "team": "Mobile",
        "latency_ms": 201
      },
      {
        "team": "Mobile",
        "latency_ms": 299
      },
      {
        "team": "Mobile",
        "latency_ms": 371
      },
      {
        "team": "Mobile",
        "latency_ms": 231
      },
      {
        "team": "Mobile",
        "latency_ms": 407
      },
      {
        "team": "Mobile",
        "latency_ms": 343
      },
      {
        "team": "Mobile",
        "latency_ms": 301
      },
      {
        "team": "Mobile",
        "latency_ms": 329
      },
      {
        "team": "Mobile",
        "latency_ms": 349
      },
      {
        "team": "Mobile",
        "latency_ms": 317
      },
      {
        "team": "Mobile",
        "latency_ms": 379
      },
      {
        "team": "Mobile",
        "latency_ms": 270
      },
      {
        "team": "Mobile",
        "latency_ms": 285
      },
      {
        "team": "Mobile",
        "latency_ms": 373
      },
      {
        "team": "Mobile",
        "latency_ms": 312
      },
      {
        "team": "Mobile",
        "latency_ms": 257
      },
      {
        "team": "Mobile",
        "latency_ms": 367
      },
      {
        "team": "Mobile",
        "latency_ms": 398
      },
      {
        "team": "Mobile",
        "latency_ms": 283
      },
      {
        "team": "Mobile",
        "latency_ms": 227
      },
      {
        "team": "Mobile",
        "latency_ms": 302
      },
      {
        "team": "Mobile",
        "latency_ms": 301
      },
      {
        "team": "Mobile",
        "latency_ms": 292
      },
      {
        "team": "Mobile",
        "latency_ms": 394
      },
      {
        "team": "Mobile",
        "latency_ms": 248
      },
      {
        "team": "Mobile",
        "latency_ms": 386
      },
      {
        "team": "Mobile",
        "latency_ms": 234
      }
    ]
  },
  "semantic_types": {
    "latency_ms": {
      "semanticType": "Duration",
      "unit": "ms"
    }
  },
  "chart_spec": {
    "chartType": "Density Plot",
    "encodings": {
      "x": {
        "field": "latency_ms"
      },
      "color": {
        "field": "team"
      }
    }
  }
}`},
	{name: "Range Area", src: `{
  "data": {
    "values": [
      {
        "month": "Jan",
        "low": 2,
        "high": 8
      },
      {
        "month": "Feb",
        "low": 3,
        "high": 10
      },
      {
        "month": "Mar",
        "low": 5,
        "high": 13
      },
      {
        "month": "Apr",
        "low": 7,
        "high": 16
      },
      {
        "month": "May",
        "low": 10,
        "high": 20
      },
      {
        "month": "Jun",
        "low": 13,
        "high": 23
      },
      {
        "month": "Jul",
        "low": 15,
        "high": 26
      },
      {
        "month": "Aug",
        "low": 15,
        "high": 26
      },
      {
        "month": "Sep",
        "low": 12,
        "high": 22
      },
      {
        "month": "Oct",
        "low": 9,
        "high": 16
      },
      {
        "month": "Nov",
        "low": 5,
        "high": 11
      },
      {
        "month": "Dec",
        "low": 3,
        "high": 8
      }
    ]
  },
  "semantic_types": {
    "month": {
      "semanticType": "Month"
    },
    "low": {
      "semanticType": "Temperature",
      "unit": "°C"
    },
    "high": {
      "semanticType": "Temperature",
      "unit": "°C"
    }
  },
  "chart_spec": {
    "chartType": "Range Area Chart",
    "encodings": {
      "x": {
        "field": "month"
      },
      "y": {
        "field": "low"
      },
      "y2": {
        "field": "high"
      }
    }
  }
}`},
	{name: "Regression", src: `{
  "data": {
    "values": [
      {
        "hours": 1,
        "score": 36
      },
      {
        "hours": 2,
        "score": 50
      },
      {
        "hours": 2,
        "score": 53
      },
      {
        "hours": 3,
        "score": 57
      },
      {
        "hours": 3,
        "score": 54
      },
      {
        "hours": 4,
        "score": 58
      },
      {
        "hours": 4,
        "score": 58
      },
      {
        "hours": 5,
        "score": 66
      },
      {
        "hours": 5,
        "score": 61
      },
      {
        "hours": 6,
        "score": 70
      },
      {
        "hours": 6,
        "score": 71
      },
      {
        "hours": 7,
        "score": 74
      },
      {
        "hours": 7,
        "score": 78
      },
      {
        "hours": 8,
        "score": 82
      },
      {
        "hours": 8,
        "score": 91
      },
      {
        "hours": 9,
        "score": 86
      },
      {
        "hours": 9,
        "score": 82
      },
      {
        "hours": 10,
        "score": 88
      }
    ]
  },
  "semantic_types": {
    "hours": {
      "semanticType": "Duration",
      "unit": "hours"
    },
    "score": {
      "semanticType": "Score"
    }
  },
  "chart_spec": {
    "chartType": "Regression",
    "encodings": {
      "x": {
        "field": "hours"
      },
      "y": {
        "field": "score"
      }
    }
  }
}`},
}
