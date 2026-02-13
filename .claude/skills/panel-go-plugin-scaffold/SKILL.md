---
name: panel-go-plugin-scaffold
description: Create and scaffold Panel.go plugins with comprehensive examples. Use when creating plugins, especially with --with-example flag for full relationship examples.
allowed-tools: Read, Write, Edit, Bash, Glob, Grep
---

# Panel.go Plugin Scaffold Expert

Expert in creating and scaffolding Panel.go plugins with comprehensive examples, including all GORM relationship types and best practices.

## Expertise

You are a Panel.go plugin scaffold expert. You understand:

- **Plugin CLI**: `panel plugin create` command and all flags
- **Scaffold System**: Backend and frontend file generation
- **Entity.go Pattern**: All GORM struct'lar tek dosyada (circular dependency önlenir)
- **Registry Pattern**: Resource'lar arası erişim için merkezi kayıt
- **Relationship Types**: BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo, MorphToMany
- **Resource Structure**: OptimizedBase, Policy, FieldResolver, RecordTitle
- **Comprehensive Examples**: 9 entity, 7 resource ile production-ready örnekler

## Patterns

### Plugin Creation

**Basit Plugin:**
```bash
panel plugin create my-plugin
```

**Frontend Olmadan:**
```bash
panel plugin create my-plugin --no-frontend
```

**Comprehensive Example (Tüm Relationship Türleri):**
```bash
panel plugin create example --with-example
```

### Plugin Yapısı

**Basit Plugin:**
```
plugins/my-plugin/
├── plugin.go              # Plugin entry point
├── plugin.yaml            # Plugin metadata
└── frontend/              # Frontend files (optional)
    ├── index.ts
    ├── package.json
    └── fields/
```

**Comprehensive Example Plugin (--with-example):**
```
plugins/example/
├── entity/
│   └── entity.go          # 9 entity (tüm GORM struct'lar)
├── resources/
│   ├── registry.go        # Registry pattern
│   ├── organization/main.go
│   ├── billing_info/main.go  # HasOne örneği
│   ├── address/main.go
│   ├── product/main.go
│   ├── category/main.go      # BelongsToMany örneği
│   ├── shipment/main.go
│   └── shipment_row/main.go
├── plugin.go              # Tüm resource'ları import/register
└── plugin.yaml
```

### Entity.go Pattern

Tüm GORM struct'lar tek bir `entity/entity.go` dosyasında tutulur:

```go
package entity

import "time"

// Organization - Ana entity
type Organization struct {
    ID          uint64
    Name        string
    Addresses   []Address    // HasMany
    BillingInfo *BillingInfo // HasOne
    Products    []Product    // HasMany
}

// BillingInfo - HasOne örneği
type BillingInfo struct {
    ID             uint64
    OrganizationID uint64        `gorm:"not null;unique;index"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID"`
}

// Address - BelongsTo örneği
type Address struct {
    ID             uint64
    OrganizationID uint64
    Organization   *Organization `gorm:"foreignKey:OrganizationID"`
}

// Product - BelongsToMany örneği
type Product struct {
    ID         uint64
    Categories []Category `gorm:"many2many:product_categories;"`
}

// Category - BelongsToMany örneği
type Category struct {
    ID       uint64
    Products []Product `gorm:"many2many:product_categories;"`
}

// Comment - MorphTo örneği
type Comment struct {
    ID              uint64
    CommentableID   uint64 `gorm:"index:idx_commentable"`
    CommentableType string `gorm:"index:idx_commentable"`
}

// Tag - MorphToMany örneği
type Tag struct {
    ID   uint64
    Name string
}
```

### Registry Pattern

Registry pattern ile resource'lar birbirini import etmeden erişir:

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

func Register(slug string, factory func() resource.Resource) {
    mu.Lock()
    defer mu.Unlock()
    registry[slug] = factory
    resource.Register(slug, factory())
}

func GetOrPanic(slug string) resource.Resource {
    r := Get(slug)
    if r == nil {
        panic("resource not found: " + slug)
    }
    return r
}
```

**Kullanım:**
```go
// Resource'lar birbirini import etmeden erişebilir
fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses"))
```

### Resource Structure

```go
package organization

import (
    "example/entity"
    "example/resources"
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

// init fonksiyonu ile resource'u registry'ye kaydet
func init() {
    resources.Register("organizations", func() resource.Resource {
        return NewOrganizationResource()
    })
}

type OrganizationResource struct {
    resource.OptimizedBase
}

type OrganizationPolicy struct{}

func (p *OrganizationPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *OrganizationPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *OrganizationPolicy) Create(ctx *context.Context) bool    { return true }
func (p *OrganizationPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *OrganizationPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type OrganizationResolveFields struct{}

func NewOrganizationResource() *OrganizationResource {
    r := &OrganizationResource{}
    r.SetSlug("organizations")
    r.SetTitle("Organizations")
    r.SetIcon("building")
    r.SetGroup("Management")
    r.SetModel(&entity.Organization{})
    r.SetFieldResolver(&OrganizationResolveFields{})
    r.SetVisible(true)
    r.SetPolicy(&OrganizationPolicy{})
    return r
}

func (r *OrganizationResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
    return []fields.Element{
        fields.ID("ID", "id"),
        fields.Text("Name", "name").Label("Organizasyon Adı").Required(),

        // Registry pattern ile circular dependency çözüldü
        fields.HasOne("BillingInfo", "billing_info", resources.GetOrPanic("billing-info")).
            ForeignKey("organization_id").OwnerKey("id").WithEagerLoad(),

        fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses")).
            ForeignKey("organization_id").OwnerKey("id").WithEagerLoad(),
    }
}

func (r *OrganizationResource) RecordTitle(record interface{}) string {
    if org, ok := record.(*entity.Organization); ok {
        return org.Name
    }
    return ""
}

func (r *OrganizationResource) With() []string {
    return []string{"BillingInfo", "Addresses"}
}
```

### Relationship Field'ları

**BelongsTo:**
```go
fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
    DisplayUsing("name").
    WithSearchableColumns("name", "email").
    WithEagerLoad()
```

**HasMany:**
```go
fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses")).
    ForeignKey("organization_id").
    OwnerKey("id").
    WithEagerLoad()
```

**HasOne:**
```go
fields.HasOne("BillingInfo", "billing_info", resources.GetOrPanic("billing-info")).
    ForeignKey("organization_id").
    OwnerKey("id").
    WithEagerLoad()
```

**BelongsToMany:**
```go
fields.BelongsToMany("Categories", "categories", resources.GetOrPanic("categories")).
    PivotTable("product_categories").
    WithEagerLoad()
```

## Anti-Patterns

❌ **Don't use separate files for each entity** - Use entity.go pattern
❌ **Don't import resources directly** - Use registry pattern
❌ **Don't skip eager loading** - Use WithEagerLoad() for N+1 prevention
❌ **Don't forget GORM tags** - Add indexes, constraints, cascade rules
❌ **Don't use tire (-) in package names** - Go converts to underscore (_)
❌ **Don't skip Policy** - Always implement all policy methods
❌ **Don't forget RecordTitle** - Implement for better UX
❌ **Don't skip With() method** - Define eager load relationships

## Decisions

### Entity.go Pattern
- **All entities in one file** - Prevents circular dependency
- **Package name: entity** - Simple and clear
- **Location: entity/entity.go** - Standard location

### Registry Pattern
- **Centralized registration** - All resources register via init()
- **Factory pattern** - Resources created on demand
- **GetOrPanic** - Fail fast for missing resources

### Package Naming
- **Plugin name with tire**: my-plugin
- **Go package name**: my_plugin (tire → underscore)
- **Automatic conversion** - Scaffold system handles this

### Relationship Types
- **BelongsTo**: Foreign key in current model
- **HasMany**: Foreign key in related model
- **HasOne**: Foreign key in related model (unique)
- **BelongsToMany**: Pivot table required
- **MorphTo**: Polymorphic foreign key
- **MorphToMany**: Polymorphic pivot table

## Sharp Edges

⚠️ **Package Naming**: Plugin adı tire (-) içeriyorsa, Go package adı underscore (_) kullanır
⚠️ **Circular Dependency**: Entity'leri ayrı dosyalarda tutmayın, entity.go pattern kullanın
⚠️ **Registry Timing**: init() fonksiyonları import sırasına göre çalışır
⚠️ **Eager Loading**: WithEagerLoad() kullanmazsanız N+1 query problemi olur
⚠️ **GORM Tags**: Index ve constraint'leri unutmayın
⚠️ **Policy Methods**: Tüm policy method'larını implement edin
⚠️ **Local Module**: go.mod'a replace directive ekleyin (local development için)

## Usage

### Creating a Simple Plugin

```bash
# Basit plugin oluştur
panel plugin create my-plugin

# Frontend olmadan
panel plugin create my-plugin --no-frontend

# Build olmadan
panel plugin create my-plugin --no-build

# Custom path
panel plugin create my-plugin --path ./custom-plugins
```

### Creating a Comprehensive Example

```bash
# Tüm relationship türlerini içeren örnek oluştur
panel plugin create example --with-example --no-frontend --no-build

# Plugin dizinine git
cd plugins/example

# Go module başlat
go mod init example

# Local panel.go modülünü kullan
echo 'replace github.com/ferdiunal/panel.go => ../../' >> go.mod

# Dependencies yükle
go mod tidy

# Compile et
go build .
```

### Comprehensive Example İçeriği

**9 Entity:**
1. Organization (HasMany, HasOne)
2. BillingInfo (BelongsTo, HasOne örneği)
3. Address (BelongsTo)
4. Product (BelongsTo, BelongsToMany, HasMany)
5. Category (BelongsToMany örneği)
6. Shipment (BelongsTo, HasMany)
7. ShipmentRow (Multiple BelongsTo)
8. Comment (MorphTo örneği)
9. Tag (MorphToMany örneği)

**7 Resource:**
1. organization
2. billing_info (HasOne field örneği)
3. address
4. product
5. category (BelongsToMany field örneği)
6. shipment
7. shipment_row

**Not:** Comment ve Tag entity'leri mevcut ama resource'ları oluşturulmadı (MorphTo/MorphToMany field'ları henüz panel.go'da implement edilmemiş)

## Best Practices

1. **Entity.go Pattern Kullan**: Tüm entity'leri tek dosyada tut
2. **Registry Pattern Kullan**: Resource'lar arası erişim için registry kullan
3. **Eager Loading**: N+1 query problemini önlemek için `WithEagerLoad()` kullan
4. **Policy Tanımla**: Her resource için tüm policy method'larını implement et
5. **RecordTitle Implement Et**: Her resource için `RecordTitle()` method'unu implement et
6. **With() Method'u**: Eager load edilecek relationship'leri `With()` method'unda belirt
7. **Türkçe Label'lar**: Field'lara Türkçe label'lar ekle
8. **GORM Tag'leri**: Index, constraint ve cascade rule'ları ekle
9. **Local Development**: go.mod'a replace directive ekle
10. **Compile Test**: Plugin oluşturduktan sonra compile et

## Documentation

- **CLI Reference**: docs/PLUGIN_CLI.md
- **Examples**: docs/PLUGIN_EXAMPLES.md (Comprehensive Example Plugin bölümü)
- **Development Guide**: docs/PLUGIN_DEVELOPMENT.md
- **System Architecture**: docs/PLUGIN_SYSTEM.md
- **Troubleshooting**: docs/PLUGIN_TROUBLESHOOTING.md
