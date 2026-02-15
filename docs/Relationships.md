# İlişkiler (Relationships) - Legacy Teknik Akış

Bu doküman, Panel.go içinde ilişki alanlarını son kullanıcıya anlaşılır şekilde kurmak için güncel ve pratik rehberdir.

Bu içerik, `resource.Register(...)`, `SetRecordTitleKey(...)` ve relationship field builder'larının birlikte kullanıldığı düşük seviye/legacy akışı hedefler.

## Hızlı Uygulama Akışı

1. İlişkili tüm resource'lar için benzersiz slug belirle.
2. Her resource'u `init()` içinde `resource.Register(slug, ...)` ile kaydet.
3. Relationship field'da doğru `relatedResource` (slug veya resource instance) ver.
4. Gerekli FK/OwnerKey/Pivot ayarlarını açıkça tanımla.
5. Dropdown kalitesi için `SetRecordTitleKey(...)` veya `SetRecordTitleFunc(...)` ayarla.

## Hızlı Özet

Panel.go'da ilişki kurarken iki şey kritiktir:

1. Doğru field constructor'ını kullanmak
2. İlişkili resource'un doğru şekilde register edilmiş olması

Kullanılan constructor'lar:

```go
fields.BelongsTo(...)
fields.HasOne(...)
fields.HasMany(...)
fields.BelongsToMany(...)
fields.NewMorphTo(...)
fields.NewMorphToMany(...)
```

## `relatedResource` Parametresi (Slug vs Resource Instance)

`BelongsTo`, `HasOne`, `HasMany`, `BelongsToMany` için 3. parametre `relatedResource`:

- `string` slug (örn: `"users"`)
- `Slug() string` dönen resource instance

Örnek:

```go
// Slug ile
fields.BelongsTo("Organization", "organization_id", "organizations")

// Resource instance ile
fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations"))
```

Not:
- `resources.GetOrPanic(...)` çekirdek kütüphane fonksiyonu değildir.
- Bu, uygulamanızda yazdığınız yardımcı bir local registry fonksiyonu olabilir.

## İlişki Tipleri

### BelongsTo

Bir kayıt tek bir üst kayda bağlıdır.

```go
fields.BelongsTo("Organization", "organization_id", "organizations").
    DisplayUsing("name").
    WithSearchableColumns("name", "email").
    AutoOptions("name").
    WithEagerLoad().
    Required()
```

Önemli metodlar:
- `DisplayUsing(key string)`
- `WithSearchableColumns(columns ...string)`
- `AutoOptions(displayField string)`
- `WithEagerLoad()` / `WithLazyLoad()`
- `Query(func(interface{}) interface{})`

### HasOne

Bir kayıt tek bir alt kayda sahiptir.

```go
fields.HasOne("BillingInfo", "billing_info", "billing-infos").
    ForeignKey("organization_id").
    OwnerKey("id").
    AutoOptions("name").
    WithEagerLoad()
```

Önemli metodlar:
- `ForeignKey(key string)`
- `OwnerKey(key string)`
- `AutoOptions(displayField string)`
- `WithEagerLoad()` / `WithLazyLoad()`

### HasMany

Bir kayıt birden fazla alt kayda sahiptir.

```go
fields.HasMany("Addresses", "addresses", "addresses").
    ForeignKey("organization_id").
    OwnerKey("id").
    AutoOptions("name").
    WithEagerLoad()
```

Önemli metodlar:
- `ForeignKey(key string)`
- `OwnerKey(key string)`
- `AutoOptions(displayField string)`
- `WithFullData()` (gerekirse daha geniş payload)
- `WithEagerLoad()` / `WithLazyLoad()`

### BelongsToMany

Çoktan çoğa ilişki.

```go
fields.BelongsToMany("Categories", "categories", "categories").
    PivotTable("product_categories").
    ForeignKey("product_id").
    RelatedKey("category_id").
    AutoOptions("name").
    WithEagerLoad()
```

Önemli metodlar:
- `PivotTable(table string)`
- `ForeignKey(key string)`
- `RelatedKey(key string)`
- `AutoOptions(displayField string)`

### MorphTo

Polimorfik tekil ilişki.

```go
fields.NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "post":  "posts",
        "video": "videos",
    }).
    Displays(map[string]string{
        "post":  "title",
        "video": "name",
    }).
    WithEagerLoad()
```

Önemli metodlar:
- `Types(map[string]string)`
- `Displays(map[string]string)`
- `GetResourceForType(morphType string)`

### MorphToMany

Polimorfik çoktan çoğa ilişki.

```go
fields.NewMorphToMany("Taggable", "taggable").
    Types(map[string]string{
        "post":  "posts",
        "video": "videos",
    }).
    Displays(map[string]string{
        "post":  "title",
        "video": "name",
    }).
    PivotTable("taggables").
    ForeignKey("tag_id").
    RelatedKey("taggable_id").
    MorphType("taggable_type")
```

## Query Callback Doğru Kullanımı

`Query` callback imzası `func(interface{}) interface{}` şeklindedir.

```go
fields.BelongsTo("Organization", "organization_id", "organizations").
    Query(func(q interface{}) interface{} {
        db, ok := q.(*gorm.DB)
        if !ok || db == nil {
            return q
        }
        return db.Where("is_active = ?", true)
    })
```

`func(q *Query) *Query` şeklindeki örnekler bu kod tabanında güncel değildir.

## AutoOptions Nasıl Çalışır?

`AutoOptions("name")` kullanıldığında backend ilişkili tablodan `id` ve `name` alanlarını çekip seçenek üretir.

İyi çalışması için:
- Relationship field'da `related_resource` doğru olmalı
- İlgili resource registry'de bulunmalı
- HasOne için `ForeignKey(...)` doğru ayarlanmalı

## Record Title İlişkilerde Görüntü Kalitesi

İlişkilerde dropdown/etiket kalitesi için resource tarafında record title ayarlayın:

```go
r.SetRecordTitleKey("name")

// veya
r.SetRecordTitleFunc(func(record any) string {
    // custom title
    return "..."
})
```

Bu ayar özellikle `BelongsTo` ve listelerde kullanıcı dostu görüntü için kritiktir.

## Resource Register ve İlişkilerin Çalışması

### 1) Global Registry (çekirdek)

`pkg/resource` API:

```go
resource.Register(slug string, res resource.Resource)
resource.Get(slug string) resource.Resource
resource.List() map[string]resource.Resource
resource.Clear()
```

Örnek:

```go
func init() {
    resource.Register("organizations", NewOrganizationResource())
}
```

### 2) Local Registry (opsiyonel, uygulama seviyesi)

Circular dependency azaltmak için uygulamanızda factory registry kurabilirsiniz:

```go
// resources/registry.go
func Register(slug string, factory func() resource.Resource) {
    registry[slug] = factory
    resource.Register(slug, factory()) // global registry'ye de aktar
}

func GetOrPanic(slug string) resource.Resource {
    r := Get(slug)
    if r == nil {
        panic("resource not found: " + slug)
    }
    return r
}
```

Not:
- `GetOrPanic` örneği uygulama helper'ıdır.
- Çekirdek `pkg/resource` içinde yoktur.

## Circular Dependency Pratiği

Eğer `AResource` ve `BResource` birbirini referans ediyorsa:

1. Her resource kendi package'inde `init()` ile register olsun
2. Relationship'te doğrudan karşı package constructor import'u yerine slug veya local registry kullanın

Örnek:

```go
fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses"))
fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations"))
```

## Sık Hata ve Çözüm

### İlişki boş geliyor

Kontrol edin:
- Slug doğru mu (`organizations` vs `organization`)
- Resource register edilmiş mi
- Model ilişkisi ve FK doğru mu
- `With()` preload listesi gerekiyor mu

### AutoOptions boş geliyor

Kontrol edin:
- `AutoOptions("...")` çağrıldı mı
- `displayField` tabloda var mı
- İlgili resource registry'de bulunuyor mu
- HasOne için `ForeignKey(...)` doğru mu

### "resource not found" hatası

Kontrol edin:
- Resource package'i import edildi mi
- `init()` içinde register var mı
- Slug birebir aynı mı

## Önerilen Checklist

- Her resource için benzersiz slug
- `init()` içinde register
- Relationship alanlarında doğru slug/resource
- Gerekli yerlerde `ForeignKey` / `OwnerKey` / `PivotTable`
- `SetRecordTitleKey` veya `SetRecordTitleFunc`
- Büyük listelerde `Query(...)` ile kısıtlama

## Sonraki Adım

- Resource yapısı için: [Kaynaklar (Resource)](Resources)
- Field kararları için: [Alanlar (Fields)](Fields)
- Uçtan uca başlangıç akışı için: [Başlarken](Getting-Started)
