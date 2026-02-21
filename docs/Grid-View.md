# Grid View

Bu doküman, index table yanında grid görünümünü resource bazında nasıl yöneteceğinizi anlatır.

## Hızlı Başlangıç

1. Resource üzerinde grid'i aç/kapat:

```go
r.SetGridEnabled(true)  // varsayılan true
// r.SetGridEnabled(false) // bu resource için grid kapat
```

2. Kart başlığında kullanılacak alanı ayarla:

```go
r.SetRecordTitleKey("name")
```

3. Field görünürlüğünü grid için düzenle:

```go
fields.Image("Görsel", "image").OnList().OnDetail()
fields.Text("Özet", "summary").HideOnList().ShowOnGrid()
fields.Text("Rozet", "badge").ShowOnlyGrid()
fields.Text("İç Not", "internal_note").HideOnGrid()
```

## Resource Bazlı Grid Kontrolü

Grid artık resource seviyesinde bir özelliktir:

- `SetGridEnabled(true)`: table + grid toggle aktif
- `SetGridEnabled(false)`: grid kapalı, index sadece table
- `IsGridEnabled()`: resource için grid durumunu döner

Not:
- Varsayılan davranış backward-compatible şekilde `true` (grid açık) kabul edilir.
- Frontend toggle görünürlüğü backend'den gelen `grid_enabled` bilgisine göre belirlenir.

## Field Visibility Kuralları

- `HideOnGrid`: grid kart/listing görünümünde gizler (header/card body)
- `ShowOnGrid`: grid'de görünürlüğü zorlar (özellikle `HideOnList` kombinasyonunda)
- `ShowOnlyGrid`: index kapsamlarında (table + grid) görünür, create/update/detail'da gizli

Öncelik:
- `HideOnGrid` grid'de baskındır
- `HideOnList` grid'de de gizler, `ShowOnGrid` ile override edilir
- `OnlyOnDetail/OnlyOnForm/OnlyOnCreate/OnlyOnUpdate` kısıtları grid'de korunur

Önemli:
- `HideOnGrid`, kart/listing görünümünü etkiler.
- Row payload içindeki alan verisi korunur (aksiyonlar, modal akışları ve computed işlemler için).

## Kart Render Sırası

Grid kartı otomatik olarak şu sırayı kullanır:

1. Varsa ilk görünür `image-field`
2. `record_title_key` alanından kart başlığı
3. Kalan görünür alanlar (field sırası)
4. Boş değerler `—` olarak gösterilir

## Computed Kart İçeriği (Display + Stack)

Kart içi computed içerik için yeni bir API yoktur. Mevcut akış kullanılır:

```go
fields.Text("Özet", "summary").
	HideOnList().
	ShowOnGrid().
	Display(func(value interface{}, item interface{}) core.Element {
		p, ok := item.(*Product)
		if !ok || p == nil {
			return fields.Stack([]core.Element{})
		}

		return fields.Stack([]core.Element{
			fields.Badge(fmt.Sprintf("Stok: %d", p.Stock)).WithProps("variant", "secondary"),
			fields.Badge(fmt.Sprintf("Fiyat: %.2f", p.Price)).WithProps("variant", "outline"),
		})
	})
```

## Uçtan Uca Örnek Resource

```go
type ProductResource struct {
	resource.OptimizedBase
}

func NewProductResource() *ProductResource {
	r := &ProductResource{}

	r.SetSlug("products")
	r.SetTitle("Products")
	r.SetRecordTitleKey("name")
	r.SetGridEnabled(true)

	r.SetFieldResolver(&ProductFieldResolver{})
	r.SetPolicy(&ProductPolicy{})

	return r
}

type ProductFieldResolver struct{}

func (f *ProductFieldResolver) Fields() []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Image("Görsel", "image").OnList().OnDetail(),
		fields.Text("Ad", "name").OnList().OnDetail().OnForm(),
		fields.Text("Durum Rozeti", "status_badge").ShowOnlyGrid(),
		fields.Text("Özet", "summary").
			HideOnList().
			ShowOnGrid().
			Display(func(value interface{}, item interface{}) core.Element {
				p, ok := item.(*Product)
				if !ok || p == nil {
					return fields.Stack([]core.Element{})
				}

				return fields.Stack([]core.Element{
					fields.Badge(fmt.Sprintf("SKU: %s", p.SKU)).WithProps("variant", "secondary"),
					fields.Badge(fmt.Sprintf("Stok: %d", p.Stock)).WithProps("variant", "outline"),
				})
			}),
	}
}
```

## API Kullanımı

Resource index:

```http
GET /api/resource/products?products[view]=grid
```

Lens:

```http
GET /api/resource/products/lens/active-products?view=grid
```

Relationship:

```http
GET /api/resource/tags?tags[view]=grid&viaResource=products&viaResourceId=42&viaRelationship=tags
```

## Response Alanları

- Resource index: `meta.grid_enabled`, `meta.record_title_key`
- Lens: `grid_enabled`, `record_title_key`

Örnek:

```json
{
  "meta": {
    "grid_enabled": true,
    "record_title_key": "name"
  }
}
```

## Troubleshooting

1. `HideOnGrid` neden görünmüyor?

- `HideOnGrid` sadece `grid` görünümünde çalışır.
- Table varsayılan olduğu için URL'de `view` yoksa alan görünmeye devam eder.
- Kontrol için:

```http
GET /api/resource/products?products[view]=grid
```

2. `props.span` neden etkisiz?

- `span` sadece grid kart body düzeninde kullanılır.
- Table görünümünde `span` uygulanmaz.
- `span` örneği:

```go
fields.Number("Fiyat", "price").WithProps("span", 6)
fields.Number("KDV", "vat_rate").WithProps("span", 6)
```

3. UI'da grid toggle görünmüyor?

- Frontend asset'leri güncel olmayabilir.
- `web` build sonrası `pkg/panel/ui` kopyalanmalı ve backend yeniden başlatılmalıdır.

```bash
cd /Users/ferdiunal/Web/panel.go/web
bun run build
cd /Users/ferdiunal/Web/panel.go
cp -R web/dist/. pkg/panel/ui/
# sonra panel uygulamasını yeniden başlat
```

4. Hızlı doğrulama checklist

- URL/query içinde `view=grid` var mı?
- Response `meta.headers` içinde `hide_on_grid` alanları düşüyor mu?
- Row payload'da `hide_on_grid` alanlarının kalması beklenir (tasarım gereği).
- `meta.grid_enabled` true mu?
- Tarayıcı hard refresh yapıldı mı?
