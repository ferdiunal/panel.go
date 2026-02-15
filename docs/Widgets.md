# Widget'lar (Cards)

Widget'lar (Cards), resource ve page ekranlarında metrik/özet verileri göstermek için kullanılır.

## Hızlı Akış

1. Kart tipini seç (`value`, `trend`, `partition`, `progress`, `table`).
2. Resource veya Page içinde `Cards()` metodunda kartları döndür.
3. Kart verisini `Resolve(...)` içinde üret.
4. Frontend, `component` adına göre doğru chart/UI bileşenini otomatik render eder.

## Desteklenen Kart Tipleri

1. `value-metric`: Tekil sayı/metin gösterimi.
2. `trend-metric`: `Area Chart - Axes` (shadcn/ui örneği).
3. `partition-metric`: `Pie Chart - Interactive` (shadcn/ui örneği).
4. `progress-metric`: `Line Chart - Interactive` (shadcn/ui örneği).
5. `table-metric`: Tablo metrik görünümü.

## Kart Oluşturma Örnekleri

### Value

```go
func (r *UserResource) Cards() []widget.Card {
	return []widget.Card{
		widget.NewCountWidget("Toplam Kullanıcı", &User{}),
	}
}
```

### Trend (Area Chart - Axes)

```go
func (r *UserResource) Cards() []widget.Card {
	return []widget.Card{
		widget.NewTrendWidget("Kayıt Trendi", &User{}, "created_at"),
	}
}
```

### Partition (Pie Chart - Interactive)

```go
func (r *OrderResource) Cards() []widget.Card {
	return []widget.Card{
		metric.NewPartition("Sipariş Durumları").
			Query(func(db *gorm.DB) (map[string]int64, error) {
				return metric.GroupByColumn(db, &Order{}, "status")
			}).
			SetColors(map[string]string{
				"pending":   "var(--chart-1)",
				"completed": "var(--chart-2)",
				"cancelled": "var(--chart-3)",
			}),
	}
}
```

### Progress (Line Chart - Interactive)

```go
func (r *OrderResource) Cards() []widget.Card {
	return []widget.Card{
		metric.NewProgress("Aylık Hedef", 1000).
			SetSeriesKey("desktop", "siparis").
			SetSeriesLabel("desktop", "Sipariş").
			SetSeriesKey("mobile", "hedef").
			SetSeriesLabel("mobile", "Hedef").
			SetActiveSeries("siparis").
			Current(func(db *gorm.DB) (int64, error) {
				return metric.CountWhere(db, &Order{}, "created_at >= ?", startOfMonth())
			}).
			History(func(db *gorm.DB) ([]map[string]interface{}, error) {
				// date + SetSeriesKey ile tanımlanan data key alanları beklenir
				return []map[string]interface{}{
					{"date": "2026-02-01", "siparis": 120, "hedef": 1000},
					{"date": "2026-02-02", "siparis": 180, "hedef": 1000},
				}, nil
			}),
	}
}
```

> `History(...)` verilmezse backend, `current/target` değerlerinden 30 günlük fallback `chartData` üretir.

## Frontend Bileşen Eşleşmesi

- `trend-metric` -> `/web/src/components/widgets/trend-metric.tsx`
- `partition-metric` -> `/web/src/components/metrics/PartitionMetric.tsx`
- `progress-metric` -> `/web/src/components/metrics/ProgressMetric.tsx`
- Router: `/web/src/components/widget-renderer.tsx`

## Veri Sözleşmesi (Data Contract)

Chart kartları için backend, legacy alanları koruyup ek olarak `chartData` üretir.
Detaylı alan yapıları ve JSON örnekleri için:

- [Charts Data Contract](Charts-Data-Contract)

## Migration Notları

- Eski kartlar (sadece `data.value`, `data.current`, `data.target`) çalışmaya devam eder.
- Yeni interaktif chart deneyimi için backend tarafında `chartData` üretmeniz önerilir.
- `trend-metric` için `widget.NewTrendWidget(...)` kullanımıyla `chartData` otomatik normalize edilir.

## Sorun Giderme

- Kart görünmüyor: `component` adı frontend switch ile eşleşmeli.
- Pie chart renkleri yanlış: `SetColors(...)` anahtarlarıyla kategori isimlerini eşleştirin.
- Line chart düz çizgi: `History(...)` çıktısında `date` ve `SetSeriesKey(...)` ile tanımlı seri key alanlarını doğrulayın.
- Area chart boş: trend sorgusunun tarih alanı (`created_at` vb.) doğru ve dolu olmalı.
