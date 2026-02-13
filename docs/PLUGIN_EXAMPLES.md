# Plugin Örnekleri

Panel.go için gerçek dünya plugin örnekleri. Bu döküman, yaygın kullanım senaryolarını ve tam plugin implementasyonlarını içerir.

## İçindekiler

- [CSV/XLSX Importer Plugin](#csvxlsx-importer-plugin)
- [Custom Field Plugin](#custom-field-plugin)
- [API Integration Plugin](#api-integration-plugin)
- [Minimal Plugin](#minimal-plugin)
- [Comprehensive Example Plugin (--with-example)](#comprehensive-example-plugin---with-example)

## CSV/XLSX Importer Plugin

Dosya import işlemleri için tam özellikli plugin örneği.

### Özellikler

- CSV ve XLSX dosya desteği
- Dosya validasyonu
- Progress tracking
- Error handling
- Custom field component

### Backend (plugin.go)

```go
package importer

import (
    "fmt"
    "github.com/ferdiunal/panel.go/pkg/plugin"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

type Plugin struct {
    plugin.BasePlugin
    db *gorm.DB
}

func init() {
    plugin.Register(&Plugin{})
}

func (p *Plugin) Name() string    { return "importer" }
func (p *Plugin) Version() string { return "1.0.0" }

func (p *Plugin) Register(panel interface{}) error {
    panelApp := panel.(*panel.Panel)
    p.db = panelApp.Db
    return nil
}

func (p *Plugin) Boot(panel interface{}) error {
    return p.db.AutoMigrate(&Import{})
}

func (p *Plugin) Resources() []resource.Resource {
    return []resource.Resource{
        NewImportResource(p.db),
    }
}

func (p *Plugin) Routes(router fiber.Router) {
    router.Post("/api/import/upload", p.handleUpload)
    router.Get("/api/import/:id/status", p.handleStatus)
}

func (p *Plugin) handleUpload(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "File required"})
    }

    // Validate file type
    if !isValidFileType(file.Filename) {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid file type"})
    }

    // Process file
    importRecord := &Import{
        Filename: file.Filename,
        Status:   "processing",
    }

    if err := p.db.Create(importRecord).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Database error"})
    }

    // Process async
    go p.processFile(file, importRecord.ID)

    return c.JSON(fiber.Map{
        "id":      importRecord.ID,
        "status":  "processing",
        "message": "File uploaded successfully",
    })
}

func (p *Plugin) processFile(file *multipart.FileHeader, importID uint) {
    // File processing logic
}

func isValidFileType(filename string) bool {
    ext := filepath.Ext(filename)
    return ext == ".csv" || ext == ".xlsx"
}
```

### Model (import.go)

```go
package importer

import "time"

type Import struct {
    ID        uint      `gorm:"primarykey"`
    Filename  string    `gorm:"not null"`
    Status    string    `gorm:"not null"` // processing, completed, failed
    Progress  int       `gorm:"default:0"`
    Total     int       `gorm:"default:0"`
    Errors    string    `gorm:"type:text"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Resource (import_resource.go)

```go
package importer

import (
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "gorm.io/gorm"
)

type ImportResource struct {
    resource.BaseResource
    db *gorm.DB
}

func NewImportResource(db *gorm.DB) *ImportResource {
    return &ImportResource{db: db}
}

func (r *ImportResource) Slug() string  { return "imports" }
func (r *ImportResource) Title() string { return "Imports" }
func (r *ImportResource) Model() interface{} { return &Import{} }

func (r *ImportResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Custom("File", "file").
            Type("import-field").
            Label("Import File").
            HelpText("Upload CSV or XLSX file").
            Required().
            OnlyOnForms(),
        fields.Text("Filename", "filename").
            Label("File Name").
            HideOnForms(),
        fields.Badge("Status", "status").
            Label("Status").
            Options(map[string]string{
                "processing": "Processing",
                "completed":  "Completed",
                "failed":     "Failed",
            }).
            Colors(map[string]string{
                "processing": "blue",
                "completed":  "green",
                "failed":     "red",
            }).
            HideOnForms(),
        fields.Number("Progress", "progress").
            Label("Progress (%)").
            HideOnForms(),
        fields.Datetime("Created At", "created_at").
            Label("Created").
            HideOnForms(),
    }
}
```

### Frontend (index.ts)

```typescript
import { Plugin } from '@/plugins/types';
import { ImportField } from './fields/ImportField';

export const ImporterPlugin: Plugin = {
  name: 'importer',
  version: '1.0.0',
  description: 'CSV/XLSX import plugin',
  author: 'Panel.go Team',

  fields: [
    {
      type: 'import-field',
      component: ImportField,
    },
  ],
};

export default ImporterPlugin;
```

### Frontend Field (ImportField.tsx)

```typescript
import React, { useState } from 'react';
import { FieldLayout } from '@/components/fields/FieldLayout';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';

interface ImportFieldProps {
  field: any;
  name: string;
  value: any;
  onChange: (value: any) => void;
  error?: string;
}

export const ImportField: React.FC<ImportFieldProps> = ({
  field,
  name,
  value,
  onChange,
  error,
}) => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    setProgress(0);

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await fetch('/api/import/upload', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) throw new Error('Upload failed');

      const data = await response.json();
      onChange(data.id);
      setProgress(100);
    } catch (err) {
      console.error('Upload error:', err);
    } finally {
      setUploading(false);
    }
  };

  return (
    <FieldLayout
      name={name}
      label={field.label}
      error={error}
      required={field.required}
      helpText={field.helpText}
    >
      <div className="space-y-4">
        <input
          type="file"
          accept=".csv,.xlsx"
          onChange={handleFileChange}
          disabled={uploading}
          className="block w-full text-sm text-gray-500
            file:mr-4 file:py-2 file:px-4
            file:rounded-md file:border-0
            file:text-sm file:font-semibold
            file:bg-primary file:text-primary-foreground
            hover:file:bg-primary/90
            disabled:opacity-50"
        />

        {uploading && (
          <div className="space-y-2">
            <Progress value={progress} />
            <p className="text-sm text-muted-foreground">
              Uploading... {progress}%
            </p>
          </div>
        )}
      </div>
    </FieldLayout>
  );
};
```

### Metadata (plugin.yaml)

```yaml
name: importer
version: 1.0.0
author: Panel.go Team
description: CSV/XLSX import plugin for Panel.go
repository: https://github.com/ferdiunal/panel.go-importer
license: MIT
```

### Kullanım

```bash
# Plugin oluştur
panel plugin create importer

# Dosyaları kopyala
# (Yukarıdaki kod örneklerini ilgili dosyalara kopyalayın)

# Plugin'i import et
# main.go
import _ "your-module/plugins/importer"

# Build ve başlat
panel plugin build
go run main.go
```

## Custom Field Plugin

Özel field component'i için minimal plugin örneği.

### Backend (plugin.go)

```go
package colorpicker

import "github.com/ferdiunal/panel.go/pkg/plugin"

type Plugin struct {
    plugin.BasePlugin
}

func init() {
    plugin.Register(&Plugin{})
}

func (p *Plugin) Name() string    { return "colorpicker" }
func (p *Plugin) Version() string { return "1.0.0" }

func (p *Plugin) Register(panel interface{}) error { return nil }
func (p *Plugin) Boot(panel interface{}) error    { return nil }
```

### Frontend (index.ts)

```typescript
import { Plugin } from '@/plugins/types';
import { ColorPickerField } from './fields/ColorPickerField';

export const ColorPickerPlugin: Plugin = {
  name: 'colorpicker',
  version: '1.0.0',
  description: 'Color picker field',
  author: 'Your Name',

  fields: [
    {
      type: 'color-picker',
      component: ColorPickerField,
    },
  ],
};

export default ColorPickerPlugin;
```

### Frontend Field (ColorPickerField.tsx)

```typescript
import React from 'react';
import { FieldLayout } from '@/components/fields/FieldLayout';

interface ColorPickerFieldProps {
  field: any;
  name: string;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export const ColorPickerField: React.FC<ColorPickerFieldProps> = ({
  field,
  name,
  value,
  onChange,
  error,
}) => {
  return (
    <FieldLayout
      name={name}
      label={field.label}
      error={error}
      required={field.required}
      helpText={field.helpText}
    >
      <div className="flex items-center gap-4">
        <input
          type="color"
          value={value || '#000000'}
          onChange={(e) => onChange(e.target.value)}
          className="h-10 w-20 rounded border cursor-pointer"
        />
        <input
          type="text"
          value={value || ''}
          onChange={(e) => onChange(e.target.value)}
          placeholder="#000000"
          className="flex-1 px-3 py-2 border rounded"
        />
      </div>
    </FieldLayout>
  );
};
```

### Backend Kullanımı

```go
fields.Custom("Brand Color", "brand_color").
    Type("color-picker").
    Label("Brand Color").
    Default("#3B82F6")
```

## API Integration Plugin

Dış API entegrasyonu için plugin örneği.

### Backend (plugin.go)

```go
package analytics

import (
    "github.com/ferdiunal/panel.go/pkg/plugin"
    "github.com/gofiber/fiber/v2"
)

type Plugin struct {
    plugin.BasePlugin
    apiKey string
}

func init() {
    plugin.Register(&Plugin{})
}

func (p *Plugin) Name() string    { return "analytics" }
func (p *Plugin) Version() string { return "1.0.0" }

func (p *Plugin) Register(panel interface{}) error {
    p.apiKey = os.Getenv("ANALYTICS_API_KEY")
    return nil
}

func (p *Plugin) Routes(router fiber.Router) {
    router.Get("/api/analytics/stats", p.handleStats)
    router.Get("/api/analytics/visitors", p.handleVisitors)
}

func (p *Plugin) handleStats(c *fiber.Ctx) error {
    // Fetch from external API
    stats, err := p.fetchStats()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(stats)
}

func (p *Plugin) fetchStats() (map[string]interface{}, error) {
    // External API call
    return map[string]interface{}{
        "visitors": 1234,
        "pageviews": 5678,
    }, nil
}
```

## Minimal Plugin

En basit plugin örneği.

### Backend (plugin.go)

```go
package minimal

import "github.com/ferdiunal/panel.go/pkg/plugin"

type Plugin struct {
    plugin.BasePlugin
}

func init() {
    plugin.Register(&Plugin{})
}

func (p *Plugin) Name() string    { return "minimal" }
func (p *Plugin) Version() string { return "1.0.0" }

func (p *Plugin) Register(panel interface{}) error { return nil }
func (p *Plugin) Boot(panel interface{}) error    { return nil }
```

### Metadata (plugin.yaml)

```yaml
name: minimal
version: 1.0.0
author: Your Name
description: Minimal plugin example
```

## Best Practices

### 1. Error Handling

```go
func (p *Plugin) Boot(panel interface{}) error {
    if err := p.initialize(); err != nil {
        return fmt.Errorf("plugin boot failed: %w", err)
    }
    return nil
}
```

### 2. Type Assertion

```go
func (p *Plugin) Register(panel interface{}) error {
    panelApp, ok := panel.(*panel.Panel)
    if !ok {
        return fmt.Errorf("invalid panel type")
    }
    p.db = panelApp.Db
    return nil
}
```

### 3. Resource Naming

```go
// ✅ İyi
func (r *ImportResource) Slug() string { return "imports" }

// ❌ Kötü
func (r *ImportResource) Slug() string { return "Import" }
```

### 4. Field Type Naming

```typescript
// ✅ İyi
fields: [
  { type: 'import-field', component: ImportField },
]

// ❌ Kötü
fields: [
  { type: 'field', component: ImportField },
]
```

## İleri Okuma

- [Plugin Sistemi](./PLUGIN_SYSTEM.md) - Genel bakış
- [CLI Komutları](./PLUGIN_CLI.md) - CLI referansı
- [Troubleshooting](./PLUGIN_TROUBLESHOOTING.md) - Sorun giderme

## Comprehensive Example Plugin (--with-example)

`--with-example` flag'i ile oluşturulan tam özellikli örnek plugin. Tüm GORM relationship türlerini ve best practice'leri içerir.

### Özellikler

- **9 Entity**: Organization, BillingInfo, Address, Product, Category, Shipment, ShipmentRow, Comment, Tag
- **7 Resource**: organization, billing_info, address, product, category, shipment, shipment_row
- **Tüm Relationship Türleri**: BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo, MorphToMany
- **Entity.go Pattern**: Tüm GORM struct'lar tek dosyada (circular dependency önlenir)
- **Registry Pattern**: Resource'lar arası erişim için merkezi kayıt
- **Production-Ready**: Compile-ready, tam dökümante edilmiş kod

### Kullanım

```bash
# Comprehensive örnek oluştur
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

### Oluşturulan Yapı

```
plugins/example/
├── entity/
│   └── entity.go              # 9 entity (tüm GORM struct'lar)
├── resources/
│   ├── registry.go            # Registry pattern
│   ├── organization/main.go   # Organization resource
│   ├── billing_info/main.go   # BillingInfo resource (HasOne örneği)
│   ├── address/main.go        # Address resource
│   ├── product/main.go        # Product resource
│   ├── category/main.go       # Category resource (BelongsToMany örneği)
│   ├── shipment/main.go       # Shipment resource
│   └── shipment_row/main.go   # ShipmentRow resource
├── plugin.go                  # Plugin entry point (tüm resource'ları import/register)
└── plugin.yaml                # Plugin metadata
```

### Entity'ler ve Relationship Türleri

#### 1. Organization (Ana Entity)
```go
type Organization struct {
    ID          uint64
    Name        string
    Email       string
    Phone       string
    Addresses   []Address    // HasMany
    BillingInfo *BillingInfo // HasOne
    Products    []Product    // HasMany
    Shipments   []Shipment   // HasMany
}
```

**Relationship'ler:**
- `HasMany`: Addresses, Products, Shipments
- `HasOne`: BillingInfo

#### 2. BillingInfo (HasOne Örneği)
```go
type BillingInfo struct {
    ID             uint64
    OrganizationID uint64        // Foreign key (unique)
    Organization   *Organization // BelongsTo
    TaxNumber      string
    TaxOffice      string
}
```

**Relationship'ler:**
- `BelongsTo`: Organization

**GORM Tag:**
```go
`gorm:"not null;unique;column:organization_id;bigint;index"`
```

#### 3. Address (BelongsTo Örneği)
```go
type Address struct {
    ID             uint64
    OrganizationID uint64
    Organization   *Organization // BelongsTo
    Name           string
    Address        string
    City           string
    Country        string
}
```

**Relationship'ler:**
- `BelongsTo`: Organization

#### 4. Product (BelongsToMany Örneği)
```go
type Product struct {
    ID             uint64
    OrganizationID uint64
    Organization   *Organization // BelongsTo
    Name           string
    Categories     []Category    // BelongsToMany
    ShipmentRows   []ShipmentRow // HasMany
}
```

**Relationship'ler:**
- `BelongsTo`: Organization
- `BelongsToMany`: Categories (pivot table: product_categories)
- `HasMany`: ShipmentRows

**GORM Tag:**
```go
`gorm:"many2many:product_categories;"`
```

#### 5. Category (BelongsToMany Örneği)
```go
type Category struct {
    ID       uint64
    Name     string
    Products []Product // BelongsToMany
}
```

**Relationship'ler:**
- `BelongsToMany`: Products (pivot table: product_categories)

#### 6. Shipment (HasMany Örneği)
```go
type Shipment struct {
    ID             uint64
    OrganizationID uint64
    Organization   *Organization // BelongsTo
    Name           string
    ShipmentRows   []ShipmentRow // HasMany
}
```

**Relationship'ler:**
- `BelongsTo`: Organization
- `HasMany`: ShipmentRows

#### 7. ShipmentRow (Multiple BelongsTo)
```go
type ShipmentRow struct {
    ID         uint64
    ShipmentID uint64
    Shipment   *Shipment // BelongsTo
    ProductID  uint64
    Product    *Product  // BelongsTo
    Quantity   int
}
```

**Relationship'ler:**
- `BelongsTo`: Shipment
- `BelongsTo`: Product

#### 8. Comment (MorphTo Örneği)
```go
type Comment struct {
    ID              uint64
    CommentableID   uint64 // Polymorphic foreign key
    CommentableType string // Polymorphic type
    Content         string
}
```

**Relationship'ler:**
- `MorphTo`: Commentable (Product/Shipment)

**GORM Tag:**
```go
`gorm:"not null;column:commentable_id;bigint;index:idx_commentable"`
`gorm:"not null;column:commentable_type;varchar(255);index:idx_commentable"`
```

#### 9. Tag (MorphToMany Örneği)
```go
type Tag struct {
    ID   uint64
    Name string
}
```

**Relationship'ler:**
- `MorphToMany`: Taggable (Product/Shipment)

**GORM Tag:**
```go
`gorm:"many2many:taggables;polymorphic:Taggable;"`
```

### Resource'lar

Her resource şu yapıyı içerir:

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
        fields.Email("Email", "email").Label("E-posta").Required(),
        fields.Tel("Phone", "phone").Label("Telefon").Required(),
        
        // Registry pattern ile circular dependency çözüldü
        fields.HasOne("BillingInfo", "billing_info", resources.GetOrPanic("billing-info")).
            ForeignKey("organization_id").OwnerKey("id").WithEagerLoad(),
        
        fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses")).
            ForeignKey("organization_id").OwnerKey("id").WithEagerLoad(),
        
        fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
        fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
    }
}

func (r *OrganizationResource) RecordTitle(record interface{}) string {
    if org, ok := record.(*entity.Organization); ok {
        return org.Name
    }
    return ""
}

func (r *OrganizationResource) With() []string {
    return []string{"BillingInfo", "Addresses", "Products", "Shipments"}
}
```

### Registry Pattern

Registry pattern ile resource'lar birbirini import etmeden erişebilir:

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

// Register bir resource factory'sini kayıt eder
func Register(slug string, factory func() resource.Resource) {
    mu.Lock()
    defer mu.Unlock()
    registry[slug] = factory
    resource.Register(slug, factory())
}

// Get kayıtlı bir resource'u slug'ına göre alır
func Get(slug string) resource.Resource {
    mu.RLock()
    defer mu.RUnlock()
    if factory, ok := registry[slug]; ok {
        return factory()
    }
    return nil
}

// GetOrPanic kayıtlı bir resource'u alır, bulamazsa panic yapar
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

### Önemli Pattern'ler

#### 1. Entity.go Pattern
Tüm GORM struct'lar tek bir `entity/entity.go` dosyasında tutulur. Bu, circular dependency'yi önler:

```go
// ❌ Yanlış: Her entity ayrı dosyada
// organization.go
type Organization struct {
    Addresses []Address // Import cycle!
}

// address.go
type Address struct {
    Organization *Organization // Import cycle!
}

// ✅ Doğru: Tüm entity'ler tek dosyada
// entity/entity.go
type Organization struct {
    Addresses []Address // No import cycle
}

type Address struct {
    Organization *Organization // No import cycle
}
```

#### 2. Registry Pattern
Resource'lar registry pattern kullanarak birbirini import etmeden erişir:

```go
// ❌ Yanlış: Direct import
import "example/resources/address"

fields.HasMany("Addresses", "addresses", address.NewAddressResource())

// ✅ Doğru: Registry pattern
fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses"))
```

#### 3. Package Naming
Plugin adı tire (-) içeriyorsa, Go package adı underscore (_) kullanır:

```bash
# Plugin adı: my-plugin
# Package adı: my_plugin

panel plugin create my-plugin --with-example
# → package my_plugin
```

#### 4. Relationship Field'ları

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

### Best Practices

1. **Entity.go Pattern Kullan**: Tüm entity'leri tek dosyada tut
2. **Registry Pattern Kullan**: Resource'lar arası erişim için registry kullan
3. **Eager Loading**: N+1 query problemini önlemek için `WithEagerLoad()` kullan
4. **Policy Tanımla**: Her resource için policy tanımla (ViewAny, View, Create, Update, Delete)
5. **RecordTitle Implement Et**: Her resource için `RecordTitle()` method'unu implement et
6. **With() Method'u**: Eager load edilecek relationship'leri `With()` method'unda belirt
7. **Türkçe Label'lar**: Field'lara Türkçe label'lar ekle
8. **GORM Tag'leri**: Index, constraint ve cascade rule'ları ekle

### Notlar

- Comment ve Tag entity'leri mevcut ama resource'ları oluşturulmadı (MorphTo/MorphToMany field'ları henüz panel.go'da implement edilmemiş)
- Entity'lerde tüm relationship türleri tanımlı (MorphTo ve MorphToMany dahil)
- Production-ready: Compile-ready, tam dökümante edilmiş kod
- Cargo.go örneğinden referans alınarak oluşturuldu

