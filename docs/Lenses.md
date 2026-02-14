# Lensler (Lenses)

Lens, bir resource'un "özel liste görünümü"dür. Index'in aynısı değildir; kendi query'si, field'ları ve kartları olabilir.

Tipik kullanım:
- "Onay Bekleyen Siparişler"
- "Son 7 Gün Aktif Kullanıcılar"
- "Yüksek Riskli İşlemler"

## Lens Arayüzü

`pkg/resource/lens.go` içindeki lens kontratı:

- `Name() string`
- `Slug() string`
- `Query(db *gorm.DB) *gorm.DB`
- `Fields() []fields.Element`
- `GetFields(ctx *context.Context) []fields.Element`
- `GetCards(ctx *context.Context) []widget.Card`

Not:
- `GetFields()` ile dinamik field üretilebilir.
- Lens field'ı boş dönerse resource field'ları fallback olarak kullanılır.

## Örnek Lens

```go
type ActiveUsersLens struct{}

func (l *ActiveUsersLens) Name() string { return "Aktif Kullanıcılar" }
func (l *ActiveUsersLens) Slug() string { return "active-users" }

func (l *ActiveUsersLens) Query(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", "active").Order("created_at DESC")
}

func (l *ActiveUsersLens) Fields() []fields.Element {
	return []fields.Element{
		fields.ID("ID"),
		fields.Text("Ad", "name"),
		fields.Text("E-posta", "email"),
	}
}

func (l *ActiveUsersLens) GetFields(ctx *context.Context) []fields.Element {
	return l.Fields()
}

func (l *ActiveUsersLens) GetCards(ctx *context.Context) []widget.Card {
	return []widget.Card{}
}
```

## Resource'a Ekleme

```go
func (r *UserResource) Lenses() []resource.Lens {
	return []resource.Lens{
		&ActiveUsersLens{},
	}
}
```

## Endpoint'ler

### Lens listesi

`GET /api/resource/:resource/lenses`

Örnek yanıt:

```json
{
  "data": [
    { "name": "Aktif Kullanıcılar", "slug": "active-users" }
  ]
}
```

### Lens verisi

`GET /api/resource/:resource/lens/:lens`

Desteklenen query parametreleri:
- `page`
- `per_page`
- `search`
- `sort_by` + `sort_order`

Uyumluluk için:
- `sort_column` + `sort_direction` da desteklenir.

Örnek yanıt:

```json
{
  "name": "Aktif Kullanıcılar",
  "resources": [],
  "prevPageUrl": null,
  "nextPageUrl": null,
  "perPage": 25,
  "softDeletes": false,
  "hasId": true,
  "headers": [],
  "data": []
}
```

Not:
- `data`, eski istemciler için geriye dönük uyumluluk alanıdır.
- Yeni istemci için ana liste `resources` alanıdır.

### Lens kartları

`GET /api/resource/:resource/lens/:lens/cards`

Örnek yanıt:

```json
{
  "data": []
}
```

## Yetki

Lens endpoint'leri `ViewAny` policy kontrolünden geçer. Policy false dönerse `403 Unauthorized` alınır.
