# Charts Data Contract

Bu doküman, shadcn/ui chart bileşenleri ile backend payload formatı arasındaki sözleşmeyi tanımlar.

## Özet

- `trend-metric` -> Area Chart - Axes
- `partition-metric` -> Pie Chart - Interactive
- `progress-metric` -> Line Chart - Interactive

Backend, mümkün olduğunda hem eski alanları (`data`, `current`, `target`) hem yeni alanları (`chartData`, `chartColors`) döndürür.

## 1. Trend Metric (`trend-metric`)

### Beklenen yapı

```json
{
  "component": "trend-metric",
  "title": "Kayıt Trendi",
  "width": "1/3",
  "data": {
    "data": [
      { "date": "2026-02-01", "value": 12 },
      { "date": "2026-02-02", "value": 18 }
    ],
    "chartData": [
      { "month": "2026-02-01", "date": "2026-02-01", "desktop": 12, "mobile": 0 },
      { "month": "2026-02-02", "date": "2026-02-02", "desktop": 18, "mobile": 0 }
    ]
  }
}
```

### Notlar

- `widget.NewTrendWidget(...)` ile `chartData` otomatik normalize edilir.
- `desktop` serisi zorunlu, `mobile` opsiyonel (yoksa `0`).

## 2. Partition Metric (`partition-metric`)

### Beklenen yapı

```json
{
  "component": "partition-metric",
  "title": "Sipariş Durumları",
  "width": "1/3",
  "data": {
    "data": {
      "pending": 15,
      "completed": 28,
      "cancelled": 4
    },
    "chartData": [
      { "month": "pending", "label": "pending", "desktop": 15, "fill": "var(--color-pending)" },
      { "month": "completed", "label": "completed", "desktop": 28, "fill": "var(--color-completed)" },
      { "month": "cancelled", "label": "cancelled", "desktop": 4, "fill": "var(--color-cancelled)" }
    ],
    "chartColors": {
      "pending": "var(--chart-1)",
      "completed": "var(--chart-2)",
      "cancelled": "var(--chart-3)"
    }
  }
}
```

### Notlar

- Frontend pie kartında aktif dilim seçimi `month` anahtarı üzerinden yapılır.
- `label`, select içinde kullanıcıya gösterilen metindir.

## 3. Progress Metric (`progress-metric`)

### Beklenen yapı

```json
{
  "component": "progress-metric",
  "title": "Aylık Hedef",
  "width": "1/3",
  "data": {
    "current": 320,
    "target": 1000,
    "percentage": 32,
    "activeSeries": "siparis",
    "series": {
      "desktop": { "key": "siparis", "label": "Sipariş", "color": "var(--chart-1)", "enabled": true },
      "mobile": { "key": "hedef", "label": "Hedef", "color": "var(--chart-2)", "enabled": true }
    },
    "chartData": [
      { "date": "2026-02-01", "siparis": 120, "hedef": 1000 },
      { "date": "2026-02-02", "siparis": 180, "hedef": 1000 },
      { "date": "2026-02-03", "siparis": 220, "hedef": 1000 }
    ]
  }
}
```

### Notlar

- `series.desktop.key` ve `series.mobile.key` değerleri line chart `dataKey` alanlarını belirler.
- `activeSeries` değeri alias (`desktop/mobile`) veya data key (`siparis/hedef`) olabilir.
- `History(...)` verilmezse backend, `current/target` değerlerinden 30 günlük fallback `chartData` üretir.

## Önerilen Backend Kullanımı

```go
metric.NewProgress("Aylık Hedef", 1000).
  SetSeriesKey("desktop", "siparis").
  SetSeriesLabel("desktop", "Sipariş").
  SetSeriesKey("mobile", "hedef").
  SetSeriesLabel("mobile", "Hedef").
  SetActiveSeries("siparis").
  Current(func(db *gorm.DB) (int64, error) {
    return metric.CountWhere(db, &Order{}, "created_at >= ?", startOfMonth)
  }).
  History(func(db *gorm.DB) ([]map[string]interface{}, error) {
    return []map[string]interface{}{
      {"date": "2026-02-01", "siparis": 120, "hedef": 1000},
      {"date": "2026-02-02", "siparis": 180, "hedef": 1000},
    }, nil
  })
```

## Frontend Kaynakları

- `/web/src/components/widgets/trend-metric.tsx`
- `/web/src/components/metrics/PartitionMetric.tsx`
- `/web/src/components/metrics/ProgressMetric.tsx`
- `/web/src/components/widget-renderer.tsx`
