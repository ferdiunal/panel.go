# Cargo.go Migration Guide - Panel.go Yeni Yapıya Geçiş

Bu dokümantasyon, cargo.go projesindeki mevcut resource ve entity yapısını panel.go'daki yeni yapıya göre nasıl güncelleyeceğinizi açıklar.

## Mevcut Durum Analizi

### Entity Yapısı
- ✅ Tek dosyada tüm entity'ler (`entity/entity.go`)
- ✅ GORM tag'leri doğru kullanılmış
- ⚠️ İlişki field isimleri karışık: `OrganizationId *Organization` (düzeltilmeli)
- ✅ Detaylı Türkçe dokümantasyon

### Resource Yapısı
- ✅ Her resource için ayrı klasör
- ✅ FieldResolver pattern kullanılıyor
- ⚠️ String-based slug'lar: `"organizations"` (resource-based'e geçilmeli)
- ✅ Policy pattern doğru kullanılmış
- ⚠️ `With()` metodu kullanılıyor (gerekli mi kontrol edilmeli)

## Yeni Yapıya Geçiş Adımları

### 1. Resource-Based Relationships (ÖNCELİKLİ)

**Eski Yaklaşım:**
```go
fields.BelongsTo("Organization", "organization_id", "organizations")
fields.HasMany("Addresses", "addresses", "addresses")
```

**Yeni Yaklaşım:**
```go
fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource())
fields.HasMany("Addresses", "addresses", address.NewAddressResource())
```

**Avantajlar:**
- ✅ Tip güvenliği (compile-time hata kontrolü)
- ✅ Refactoring desteği (resource adı değişirse otomatik güncellenir)
- ✅ IDE autocomplete
- ✅ Tablo adı otomatik alınır

**Migration:**
```go
// resources/address/main.go
import (
    "cargo.go/resources/organization"
)

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),
        // Eski: fields.BelongsTo("Organization", "organization_id", "organizations"),
        // Yeni:
        fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()),
        // ...
    }
}
```

### 2. Entity Field Naming Düzeltmesi

**Eski Yaklaşım:**
```go
type Address struct {
    OrganizationID uint64        `gorm:"not null;column:organization_id"`
    OrganizationId *Organization `gorm:"foreignKey:OrganizationID"` // ❌ Karışık
}
```

**Yeni Yaklaşım:**
```go
type Address struct {
    OrganizationID uint64        `gorm:"not null;column:organization_id"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID"` // ✅ Temiz
}
```

**Migration:**
1. Entity field isimlerini düzelt: `OrganizationId` → `Organization`
2. GORM tag'lerini kontrol et (değişmemeli)
3. Tüm referansları güncelle

### 3. GORM Relationship Loading

Panel.go'da artık GORM'un native API'leri kullanılıyor:

**Lazy Loading:**
```go
// GORM Association API kullanır
db.Model(item).Association("Organization").Find(&organization)
```

**Eager Loading:**
```go
// GORM Query Builder kullanır
db.Table("organizations").Where("id IN ?", ids).Find(&organizations)
```

**With() Metodu:**
```go
// Mevcut kullanım
func (r *AddressResource) With() []string {
    return []string{"Organization"}
}
```

⚠️ **Kontrol Edilmeli:** Panel.go'daki yeni yapıda `With()` metodu hala gerekli mi? GORM Association API kullanıldığı için bu metod gereksiz olabilir.

### 4. Eager Loading Stratejisi

**Best Practice:**
```go
// Liste sayfaları için eager loading kullan
fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()).
    WithEagerLoad()

// Detay sayfaları için lazy loading yeterli
fields.HasMany("Addresses", "addresses", address.NewAddressResource()).
    WithLazyLoad()
```

## Öncelikli Değişiklikler

### 1. Organization Resource

```go
// resources/organization/main.go
func (r *OrganizationResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),
        fields.Text("Name", "name").Label("Şirket Adı").Required(),
        fields.Email("Email", "email").Label("Şirket E-posta").Required(),
        fields.Tel("Phone", "phone").Label("Şirket Telefon").Required(),

        // Eski: fields.HasMany("Addresses", "addresses", "addresses")
        // Yeni:
        fields.HasMany("Addresses", "addresses", address.NewAddressResource()).
            ForeignKey("organization_id").
            OwnerKey("id").
            WithEagerLoad(),

        fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
        fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
    }
}
```

### 2. Address Resource

```go
// resources/address/main.go
import (
    "cargo.go/resources/organization"
)

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),

        // Eski: fields.BelongsTo("Organization", "organization_id", "organizations"),
        // Yeni:
        fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()).
            DisplayUsing("name").
            WithSearchableColumns("name", "email").
            WithEagerLoad(),

        fields.Text("Name", "name"),
        // ...
    }
}
```

### 3. Entity Düzeltmeleri

```go
// entity/entity.go

// Eski:
type Address struct {
    OrganizationID uint64        `gorm:"not null;column:organization_id"`
    OrganizationId *Organization `gorm:"foreignKey:OrganizationID"` // ❌
}

// Yeni:
type Address struct {
    OrganizationID uint64        `gorm:"not null;column:organization_id"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID"` // ✅
}
```

## Migration Checklist

### Adım 1: Entity Düzeltmeleri
- [ ] `OrganizationId` → `Organization` (tüm entity'lerde)
- [ ] `PriceListId` → `PriceList`
- [ ] `CargoCompanyId` → `CargoCompany`
- [ ] `ShipmentId` → `Shipment`
- [ ] `ProductId` → `Product`

### Adım 2: Resource Import'ları
- [ ] Her resource için gerekli import'ları ekle
- [ ] Circular dependency kontrolü yap

### Adım 3: BelongsTo Relationships
- [ ] Address → Organization
- [ ] Commission → Organization
- [ ] Commission → PriceList
- [ ] Price → PriceList
- [ ] PriceList → CargoCompany
- [ ] Shipment → Organization
- [ ] Shipment → SenderAddress
- [ ] Shipment → ReceiverAddress
- [ ] ShipmentRow → Shipment
- [ ] ShipmentRow → Product
- [ ] Product → Organization

### Adım 4: HasMany Relationships
- [ ] Organization → Addresses
- [ ] Organization → Commissions
- [ ] Organization → Shipments
- [ ] Organization → Products
- [ ] CargoCompany → PriceLists
- [ ] PriceList → Prices
- [ ] PriceList → Commissions
- [ ] Shipment → ShipmentRows
- [ ] Product → ShipmentRows

### Adım 5: Test ve Doğrulama
- [ ] Compile hataları kontrol et
- [ ] Resource'ları test et
- [ ] İlişkilerin doğru yüklendiğini kontrol et
- [ ] N+1 query problemi var mı kontrol et

## Örnek Migration: Address Resource

**Önce:**
```go
package address

import (
    "cargo.go/entity"
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/core"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),
        fields.BelongsTo("Organization", "organization_id", "organizations"),
        fields.Text("Name", "name"),
        // ...
    }
}
```

**Sonra:**
```go
package address

import (
    "cargo.go/entity"
    "cargo.go/resources/organization" // Yeni import
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/core"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),
        fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()).
            DisplayUsing("name").
            WithSearchableColumns("name", "email", "phone").
            WithEagerLoad(),
        fields.Text("Name", "name"),
        // ...
    }
}
```

## Best Practices

### 1. Eager Loading
```go
// Liste sayfalarında N+1 query problemini önlemek için
fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()).
    WithEagerLoad()
```

### 2. Searchable Columns
```go
// Arama performansı için indexed sütunları kullan
fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()).
    WithSearchableColumns("name", "email") // Indexed sütunlar
```

### 3. Display Field
```go
// Kullanıcı dostu field seç
fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()).
    DisplayUsing("name") // Anlamlı field
```

### 4. AutoOptions
```go
// Manuel options tanımlamak yerine AutoOptions kullan
fields.BelongsTo("Organization", "organization_id", organization.NewOrganizationResource()).
    AutoOptions("name")
```

## Circular Dependency ve Cycle Import Çözümleri

### Problem: Import Cycle

Resource-based relationships kullanırken en yaygın sorun **import cycle** (circular dependency) hatasıdır.

**Örnek Senaryo:**
```go
// resources/address/main.go
import "cargo.go/resources/organization" // address → organization

// resources/organization/main.go
import "cargo.go/resources/address" // organization → address

// HATA: import cycle not allowed
```

Bu durum şu ilişkilerde ortaya çıkar:
- Organization → HasMany → Addresses
- Address → BelongsTo → Organization

### Çözüm 1: Resource Registry Pattern (ÖNERİLEN)

Merkezi bir registry kullanarak resource'ları kaydetmek ve cycle'ı kırmak.

**1. Resource Registry Oluştur:**

```go
// resources/registry.go
package resources

import (
    "sync"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

var (
    registry = make(map[string]func() resource.Resource)
    mu       sync.RWMutex
)

// Register bir resource factory'sini kaydet
func Register(slug string, factory func() resource.Resource) {
    mu.Lock()
    defer mu.Unlock()
    registry[slug] = factory
}

// Get kayıtlı bir resource'u al
func Get(slug string) resource.Resource {
    mu.RLock()
    defer mu.RUnlock()
    if factory, ok := registry[slug]; ok {
        return factory()
    }
    return nil
}

// GetOrPanic kayıtlı bir resource'u al, yoksa panic
func GetOrPanic(slug string) resource.Resource {
    r := Get(slug)
    if r == nil {
        panic("resource not found: " + slug)
    }
    return r
}
```

**2. Resource'ları Registry'ye Kaydet:**

```go
// resources/organization/main.go
package organization

import (
    "cargo.go/entity"
    "cargo.go/resources" // Sadece registry import et
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/core"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
    // Resource'u registry'ye kaydet
    resources.Register("organizations", func() resource.Resource {
        return NewOrganizationResource()
    })
}

type OrganizationResource struct {
    resource.OptimizedBase
}

func NewOrganizationResource() *OrganizationResource {
    r := &OrganizationResource{}
    r.SetSlug("organizations")
    r.SetTitle("Organizations")
    r.SetModel(&entity.Organization{})
    r.SetFieldResolver(&OrganizationResolveFields{})
    return r
}

type OrganizationResolveFields struct{}

func (r *OrganizationResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),
        fields.Text("Name", "name").Required(),

        // Registry'den address resource'unu al
        fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses")).
            ForeignKey("organization_id").
            WithEagerLoad(),
    }
}
```

**3. Address Resource'unu Kaydet:**

```go
// resources/address/main.go
package address

import (
    "cargo.go/entity"
    "cargo.go/resources" // Sadece registry import et
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/core"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
    // Resource'u registry'ye kaydet
    resources.Register("addresses", func() resource.Resource {
        return NewAddressResource()
    })
}

type AddressResource struct {
    resource.OptimizedBase
}

func NewAddressResource() *AddressResource {
    r := &AddressResource{}
    r.SetSlug("addresses")
    r.SetTitle("Addresses")
    r.SetModel(&entity.Address{})
    r.SetFieldResolver(&AddressResolveFields{})
    return r
}

type AddressResolveFields struct{}

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),

        // Registry'den organization resource'unu al
        fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
            DisplayUsing("name").
            WithEagerLoad(),

        fields.Text("Name", "name"),
    }
}
```

**4. Resource'ları Yükle:**

```go
// main.go veya init.go
package main

import (
    _ "cargo.go/resources/address"      // init() çalışır
    _ "cargo.go/resources/organization" // init() çalışır
    // Diğer resource'lar...
)

func main() {
    // Resource'lar otomatik olarak registry'ye kaydedildi
    // ...
}
```

**Avantajlar:**
- ✅ Circular dependency yok
- ✅ Resource'lar birbirini import etmez
- ✅ Merkezi yönetim
- ✅ Lazy loading desteği
- ✅ Test edilebilir

### Çözüm 2: Lazy Resource Loading

Resource instance'ları lazy olarak yüklemek için function kullanmak.

```go
// resources/address/main.go
package address

import (
    "cargo.go/entity"
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/core"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

// Lazy loading için function
var getOrganizationResource = func() resource.Resource {
    // Import cycle'ı kırmak için lazy loading
    // Bu function main.go'da set edilecek
    return nil
}

func SetOrganizationResourceGetter(getter func() resource.Resource) {
    getOrganizationResource = getter
}

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),

        // Lazy loading ile resource al
        fields.BelongsTo("Organization", "organization_id", getOrganizationResource()).
            DisplayUsing("name").
            WithEagerLoad(),
    }
}
```

```go
// main.go
package main

import (
    "cargo.go/resources/address"
    "cargo.go/resources/organization"
)

func init() {
    // Lazy loading function'ları set et
    address.SetOrganizationResourceGetter(func() resource.Resource {
        return organization.NewOrganizationResource()
    })
}
```

**Dezavantajlar:**
- ⚠️ Manuel setup gerekli
- ⚠️ Daha fazla boilerplate kod

### Çözüm 3: String-Based Slugs (Geçici Çözüm)

Resource-based yerine string-based slug kullanmak.

```go
// resources/address/main.go
func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),

        // String-based slug (cycle yok ama tip güvenliği yok)
        fields.BelongsTo("Organization", "organization_id", "organizations").
            DisplayUsing("name").
            WithEagerLoad(),
    }
}
```

**Dezavantajlar:**
- ❌ Tip güvenliği yok
- ❌ Refactoring desteği yok
- ❌ IDE autocomplete yok

### Entity Field Naming ve GORM İlişkileri

**Problem:** Entity field isimleri karışık ve GORM Association API ile uyumsuz.

```go
// ❌ YANLIŞ
type Address struct {
    OrganizationID uint64        `gorm:"not null;column:organization_id"`
    OrganizationId *Organization `gorm:"foreignKey:OrganizationID"` // Karışık!
}
```

**Neden Sorun:**
1. Field adı `OrganizationId` ama tip `*Organization` (karışık)
2. GORM Association API `Association("Organization")` bekler ama field adı `OrganizationId`
3. Panel.go'nun yeni yapısı GORM Association API kullanır

**Çözüm:**

```go
// ✅ DOĞRU
type Address struct {
    OrganizationID uint64        `gorm:"not null;column:organization_id"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID"` // Temiz!
}
```

**GORM Association API Kullanımı:**

```go
// Lazy loading
var address entity.Address
db.Model(&address).Association("Organization").Find(&address.Organization)

// Eager loading (Preload)
var addresses []entity.Address
db.Preload("Organization").Find(&addresses)
```

### Entity Yapısı Best Practices

**1. Tek Dosyada Tüm Entity'ler:**

```go
// entity/entity.go
package entity

// Tüm entity'ler bu dosyada
type Organization struct { ... }
type Address struct { ... }
type Product struct { ... }
// ...
```

**Avantajlar:**
- ✅ Circular dependency riski yok
- ✅ Entity'ler birbirini import etmez
- ✅ GORM ilişkileri doğrudan kullanılabilir

**2. GORM İlişki Tag'leri:**

```go
type Organization struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement"`
    Addresses []Address `gorm:"foreignKey:OrganizationID;references:ID"`
}

type Address struct {
    ID             uint64        `gorm:"primaryKey;autoIncrement"`
    OrganizationID uint64        `gorm:"not null"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID"`
}
```

**3. Field Naming Convention:**

```go
// ✅ DOĞRU: Field adı = İlişki adı
Organization   *Organization  // Association("Organization")
Product        *Product       // Association("Product")
Addresses      []Address      // Association("Addresses")

// ❌ YANLIŞ: Field adı ≠ İlişki adı
OrganizationId *Organization  // Association("Organization") ama field "OrganizationId"
ProductId      *Product       // Karışık!
```

### Resource Package Yapısı

**Önerilen Yapı:**

```
resources/
├── registry.go              # Resource registry
├── address/
│   └── main.go             # AddressResource
├── organization/
│   └── main.go             # OrganizationResource
├── product/
│   └── main.go             # ProductResource
└── ...
```

**registry.go:**
```go
package resources

import (
    "sync"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

var (
    registry = make(map[string]func() resource.Resource)
    mu       sync.RWMutex
)

func Register(slug string, factory func() resource.Resource) {
    mu.Lock()
    defer mu.Unlock()
    registry[slug] = factory
}

func Get(slug string) resource.Resource {
    mu.RLock()
    defer mu.RUnlock()
    if factory, ok := registry[slug]; ok {
        return factory()
    }
    return nil
}

func GetOrPanic(slug string) resource.Resource {
    r := Get(slug)
    if r == nil {
        panic("resource not found: " + slug)
    }
    return r
}
```

### Migration Stratejisi

**Adım 1: Registry Oluştur**
```bash
touch resources/registry.go
```

**Adım 2: Entity Field'larını Düzelt**
```go
// entity/entity.go
// OrganizationId → Organization
// ProductId → Product
// vb.
```

**Adım 3: Resource'ları Registry'ye Kaydet**
```go
// Her resource/*/main.go dosyasına init() ekle
func init() {
    resources.Register("organizations", func() resource.Resource {
        return NewOrganizationResource()
    })
}
```

**Adım 4: ResolveFields'ı Güncelle**
```go
// Registry'den resource al
fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations"))
```

**Adım 5: main.go'da Resource'ları Yükle**
```go
import (
    _ "cargo.go/resources/address"
    _ "cargo.go/resources/organization"
    // ...
)
```

## Sorun Giderme

### Compile Hatası: "import cycle not allowed"
**Problem:** Resource'lar birbirini import ediyor
**Çözüm:** Resource Registry Pattern kullan (yukarıda detaylı açıklandı)

### Compile Hatası: "OrganizationId field bulunamıyor"
**Problem:** Entity field isimleri düzeltilmemiş
**Çözüm:** Entity field isimlerini düzelt (`OrganizationId` → `Organization`)

### Runtime Hatası: "resource not found: organizations"
**Problem:** Resource registry'ye kaydedilmemiş
**Çözüm:** Resource'un `init()` fonksiyonunda `resources.Register()` çağrıldığından emin ol

### N+1 Query Problemi
**Problem:** Her kayıt için ayrı query çalışıyor
**Çözüm:** `WithEagerLoad()` kullan

### GORM Association Hatası
**Problem:** `Association("Organization")` çalışmıyor
**Çözüm:** Entity field adının ilişki adıyla aynı olduğundan emin ol

## Referanslar

- [Panel.go Relationships](../panel.go/docs/Relationships.md)
- [Panel.go Relationship Loading Implementation](../panel.go/docs/RELATIONSHIP_LOADING_IMPLEMENTATION.md)
- [GORM Associations](https://gorm.io/docs/associations.html)

## Değişiklik Geçmişi

### 2026-02-08: İlk Versiyon
- Mevcut yapı analizi
- Resource-based relationships migration
- Entity field naming düzeltmeleri
- Best practices ve örnekler
