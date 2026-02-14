# Resource Template

Bu şablon, yeni resource oluştururken ilişki ve register akışını doğru kurmanız için hazırlanmıştır.

## 1) Standart Template (Global Registry)

```go
package products

import (
    "myapp/entity"

    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/data"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "gorm.io/gorm"
)

func init() {
    resource.Register("products", NewProductResource())
}

type ProductResource struct {
    resource.OptimizedBase
}

type ProductFieldResolver struct{}

func NewProductResource() *ProductResource {
    r := &ProductResource{}

    r.SetModel(&entity.Product{})
    r.SetSlug("products")
    r.SetTitle("Products")
    r.SetIcon("package")
    r.SetGroup("Catalog")
    r.SetVisible(true)

    // İlişki etiketleri için okunabilir başlık
    r.SetRecordTitleKey("name")

    r.SetFieldResolver(&ProductFieldResolver{})
    r.SetPolicy(&ProductPolicy{})
    return r
}

func (r *ProductResource) Repository(db *gorm.DB) data.DataProvider {
    return data.NewGormDataProvider(db, &entity.Product{})
}

func (r *ProductFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    return []fields.Element{
        fields.ID("ID"),
        fields.Text("Name", "name").Required(),

        // Slug ile ilişki
        fields.BelongsTo("Organization", "organization_id", "organizations").
            DisplayUsing("name").
            AutoOptions("name").
            WithEagerLoad(),

        fields.HasMany("Prices", "prices", "prices").
            ForeignKey("product_id").
            WithEagerLoad(),
    }
}

type ProductPolicy struct{}

func (p *ProductPolicy) ViewAny(ctx *context.Context) bool { return true }
func (p *ProductPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *ProductPolicy) Create(ctx *context.Context) bool { return true }
func (p *ProductPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *ProductPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }
func (p *ProductPolicy) Restore(ctx *context.Context, model interface{}) bool { return false }
func (p *ProductPolicy) ForceDelete(ctx *context.Context, model interface{}) bool { return false }
```

## 2) Circular Dependency İçin Local Registry (Opsiyonel)

Büyük projelerde package'ler birbirini import etmeye başladığında local registry yardımcı olur.

```go
// myapp/resources/registry.go
package resources

import (
    "sync"

    "github.com/ferdiunal/panel.go/pkg/resource"
)

var (
    mu       sync.RWMutex
    registry = make(map[string]func() resource.Resource)
)

func Register(slug string, factory func() resource.Resource) {
    mu.Lock()
    defer mu.Unlock()

    registry[slug] = factory

    // Panel.go global registry
    resource.Register(slug, factory())
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

İlişki tarafında kullanım:

```go
fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations"))
fields.HasMany("Prices", "prices", resources.GetOrPanic("prices"))
```

Not:
- `GetOrPanic` çekirdek `pkg/resource` fonksiyonu değildir.
- Uygulamanızda yazdığınız helper'dır.

## 3) Register Checklist

- `init()` içinde register var mı?
- Slug benzersiz mi?
- Resource package'i uygulama başlarken import ediliyor mu?
- İlişki verdiğiniz slug/resource gerçekten register edilmiş mi?

## 4) Sık Hata

### "resource not found"

- Yanlış slug
- Resource package import edilmemiş
- `init()` içinde register eksik

### İlişki alanı boş geliyor

- `BelongsTo`/`HasMany` üçüncü parametre yanlış
- `ForeignKey`/`OwnerKey` model ile uyuşmuyor
- `AutoOptions` display field yanlış

## 5) Referans

- `docs/Resources.md`
- `docs/Relationships.md`
- `docs/Relationships-API.md`
