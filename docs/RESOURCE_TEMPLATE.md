# Resource Template - Circular Dependency'den Kaçınma

Bu template, yeni resource oluştururken circular dependency probleminden kaçınmak için kullanılır.

## Resource Registry Pattern

Panel.go, circular dependency problemini çözmek için **Resource Registry Pattern** kullanır. Bu pattern sayesinde resource'lar birbirini direkt import etmez, bunun yerine merkezi bir registry'den alır.

## Temel Resource Yapısı

```go
package myresource

import (
    "yourproject/entity"
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/core"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

// init fonksiyonu ile resource'u registry'ye kaydet
func init() {
    resource.Register("myresources", func() resource.Resource {
        return NewMyResource()
    })
}

// MyResource - Resource struct'ı
type MyResource struct {
    resource.OptimizedBase
}

// MyResolveFields - Field resolver struct'ı
type MyResolveFields struct{}

// NewMyResource - Resource constructor
func NewMyResource() *MyResource {
    r := &MyResource{}
    r.SetSlug("myresources")
    r.SetTitle("My Resources")
    r.SetIcon("myresources")
    r.SetGroup("My Group")
    r.SetModel(&entity.MyEntity{})
    r.SetFieldResolver(&MyResolveFields{})
    r.SetVisible(true)
    return r
}

// ResolveFields - Field'ları döndür
func (r *MyResolveFields) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID("ID", "id"),

        // BelongsTo ilişkisi - Registry'den al
        fields.BelongsTo("Organization", "organization_id", resource.GetOrPanic("organizations")).
            DisplayUsing("name").
            WithSearchableColumns("name", "email").
            WithEagerLoad(),

        // HasMany ilişkisi - Registry'den al
        fields.HasMany("Items", "items", resource.GetOrPanic("items")).
            ForeignKey("myresource_id").
            WithEagerLoad(),

        fields.Text("Name", "name").Required(),
        fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
        fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
    }
}
```

## Policy Ekleme (Opsiyonel)

```go
// MyPolicy - Authorization policy
type MyPolicy struct{}

func (p *MyPolicy) ViewAny(ctx *context.Context) bool {
    return true
}

func (p *MyPolicy) View(ctx *context.Context, model interface{}) bool {
    return true
}

func (p *MyPolicy) Create(ctx *context.Context) bool {
    return true
}

func (p *MyPolicy) Update(ctx *context.Context, model interface{}) bool {
    return true
}

func (p *MyPolicy) Delete(ctx *context.Context, model interface{}) bool {
    return true
}

// Constructor'da policy'yi set et
func NewMyResource() *MyResource {
    r := &MyResource{}
    // ...
    r.SetPolicy(&MyPolicy{})
    return r
}
```

## Entity Tanımı

```go
// entity/entity.go
package entity

import "time"

// MyEntity - Entity struct'ı
//
// BelongsTo: Organization (bir kayıt bir organizasyona aittir)
// HasMany: Items (bir kaydın birden fazla item'ı olabilir)
type MyEntity struct {
    ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint"`
    OrganizationID uint64        `gorm:"not null;column:organization_id;bigint;index"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID"` // ✅ Temiz naming
    Name           string        `gorm:"not null;column:name;varchar(255)"`
    Items          []Item        `gorm:"foreignKey:MyEntityID;references:ID"`
    CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz"`
    UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz"`
}
```

## main.go'da Resource'ları Yükle

```go
package main

import (
    // Resource'ları import et (init() fonksiyonları çalışır)
    _ "yourproject/resources/address"
    _ "yourproject/resources/organization"
    _ "yourproject/resources/myresource"
    // ...
)

func main() {
    // Resource'lar otomatik olarak registry'ye kaydedildi
    // ...
}
```

## Önemli Noktalar

### ✅ YAPILMASI GEREKENLER

1. **init() Fonksiyonu Kullan**
   ```go
   func init() {
       resource.Register("myresources", func() resource.Resource {
           return NewMyResource()
       })
   }
   ```

2. **Registry'den Resource Al**
   ```go
   fields.BelongsTo("Organization", "organization_id", resource.GetOrPanic("organizations"))
   ```

3. **Entity Field Naming**
   ```go
   Organization *Organization // ✅ Temiz
   ```

4. **main.go'da Import Et**
   ```go
   _ "yourproject/resources/myresource"
   ```

### ❌ YAPILMAMASI GEREKENLER

1. **Direkt Resource Import Etme**
   ```go
   import "yourproject/resources/organization" // ❌ Circular dependency riski
   ```

2. **Karışık Entity Field Naming**
   ```go
   OrganizationId *Organization // ❌ Karışık
   ```

3. **String-Based Slug Kullanma**
   ```go
   fields.BelongsTo("Organization", "organization_id", "organizations") // ❌ Tip güvenliği yok
   ```

## Registry API

### Register
```go
resource.Register(slug string, factory func() resource.Resource)
```
Resource'u registry'ye kaydet. init() fonksiyonunda kullanılır.

### Get
```go
resource.Get(slug string) resource.Resource
```
Resource'u al. Bulamazsa nil döner.

### GetOrPanic
```go
resource.GetOrPanic(slug string) resource.Resource
```
Resource'u al. Bulamazsa panic yapar. ResolveFields içinde kullanılır.

### List
```go
resource.List() []string
```
Tüm kayıtlı resource slug'larını listele. Debug amaçlı.

### Clear
```go
resource.Clear()
```
Tüm registry'yi temizle. Test amaçlı.

## Test Örneği

```go
package myresource_test

import (
    "testing"
    "github.com/ferdiunal/panel.go/pkg/resource"
    _ "yourproject/resources/myresource"
)

func TestResourceRegistration(t *testing.T) {
    // Resource kayıtlı mı kontrol et
    r := resource.Get("myresources")
    if r == nil {
        t.Fatal("myresources resource not registered")
    }

    // Slug doğru mu kontrol et
    if r.Slug() != "myresources" {
        t.Errorf("expected slug 'myresources', got '%s'", r.Slug())
    }
}

func TestResourceRegistry(t *testing.T) {
    // Tüm resource'ları listele
    slugs := resource.List()
    if len(slugs) == 0 {
        t.Fatal("no resources registered")
    }

    // Her resource'u test et
    for _, slug := range slugs {
        r := resource.Get(slug)
        if r == nil {
            t.Errorf("resource %s not found", slug)
        }
    }
}
```

## Sorun Giderme

### "resource not found: myresources" Hatası

**Neden:** Resource registry'ye kaydedilmemiş

**Çözüm:**
1. init() fonksiyonunun olduğundan emin ol
2. main.go'da resource package'ının import edildiğinden emin ol
3. resource.Register() çağrısının doğru slug ile yapıldığından emin ol

### "import cycle not allowed" Hatası

**Neden:** Resource'lar birbirini direkt import ediyor

**Çözüm:**
1. Direkt import yerine resource.GetOrPanic() kullan
2. Registry pattern'i doğru uygulandığından emin ol

### "OrganizationId field not found" Hatası

**Neden:** Entity field naming yanlış

**Çözüm:**
1. Entity field isimlerini düzelt: `OrganizationId` → `Organization`
2. GORM tag'lerini kontrol et

## Referanslar

- [Panel.go Relationships](./Relationships.md)
- [Panel.go Relationship Loading Implementation](./RELATIONSHIP_LOADING_IMPLEMENTATION.md)
- [GORM Associations](https://gorm.io/docs/associations.html)
