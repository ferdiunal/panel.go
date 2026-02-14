# Plugin Development Guide

Panel.go için plugin geliştirme rehberi. Bu dokümanda plugin sisteminin nasıl kullanılacağı ve custom plugin'lerin nasıl oluşturulacağı anlatılmaktadır.

## İçindekiler

1. [Plugin Sistemi Genel Bakış](#plugin-sistemi-genel-bakış)
2. [Backend Plugin Geliştirme](#backend-plugin-geliştirme)
3. [Frontend Plugin Geliştirme](#frontend-plugin-geliştirme)
4. [Plugin Yükleme Stratejileri](#plugin-yükleme-stratejileri)
5. [Best Practices](#best-practices)
6. [Örnek Plugin](#örnek-plugin)

## Plugin Sistemi Genel Bakış

Panel.go plugin sistemi, uygulamaya yeni özellikler eklemenizi sağlar:

- **Resources**: Yeni CRUD resource'ları
- **Pages**: Custom sayfalar
- **Fields**: Custom field component'leri
- **Middleware**: HTTP middleware'ler
- **Routes**: Custom API endpoint'leri
- **Migrations**: Veritabanı migration'ları

### Mimari

Plugin sistemi iki katmandan oluşur:

1. **Backend (Go)**: Resource, middleware, route, migration yönetimi
2. **Frontend (TypeScript/React)**: Custom field, widget, page component'leri

## Backend Plugin Geliştirme

### 1. Plugin Yapısı

```
plugins/
└── my-plugin/
    ├── plugin.yaml          # Plugin metadata
    ├── plugin.go            # Plugin implementation
    ├── resources/           # Custom resources
    ├── migrations/          # Database migrations
    └── frontend/            # Frontend components
        ├── index.ts
        └── fields/
```

### 2. Plugin Interface

Her plugin `plugin.Plugin` interface'ini implement etmelidir:

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/plugin"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/gofiber/fiber/v2"
)

type MyPlugin struct {
    plugin.BasePlugin
}

// Metadata
func (p *MyPlugin) Name() string        { return "my-plugin" }
func (p *MyPlugin) Version() string     { return "1.0.0" }
func (p *MyPlugin) Author() string      { return "Your Name" }
func (p *MyPlugin) Description() string { return "Plugin description" }

// Lifecycle
func (p *MyPlugin) Register(panel interface{}) error {
    // Plugin kaydı sırasında çağrılır
    return nil
}

func (p *MyPlugin) Boot(panel interface{}) error {
    // Panel başlatıldığında çağrılır
    return nil
}

// Capabilities
func (p *MyPlugin) Resources() []resource.Resource {
    return []resource.Resource{
        // Custom resources
    }
}

func (p *MyPlugin) Middleware() []fiber.Handler {
    return []fiber.Handler{
        // Custom middleware
    }
}

func (p *MyPlugin) Routes(router fiber.Router) {
    // Custom routes
}

func (p *MyPlugin) Migrations() []plugin.Migration {
    return []plugin.Migration{
        // Database migrations
    }
}
```

### 3. Plugin Kaydı

Plugin'i global registry'ye kaydetmek için `init()` fonksiyonu kullanın:

```go
func init() {
    plugin.Register(&MyPlugin{})
}
```

### 4. Custom Resource Ekleme

```go
func (p *MyPlugin) Resources() []resource.Resource {
    return []resource.Resource{
        resource.New(
            &MyModel{},
            "my-resource",
            "My Resource",
        ).Icon("cube").
            Group("My Group").
            Fields(func(r *resource.Resource) []fields.Field {
                return []fields.Field{
                    fields.ID(),
                    fields.Text("Name", "name").Required(),
                    fields.Textarea("Description", "description"),
                }
            }),
    }
}
```

### 5. Custom Middleware Ekleme

```go
func (p *MyPlugin) Middleware() []fiber.Handler {
    return []fiber.Handler{
        func(c *fiber.Ctx) error {
            // Middleware logic
            return c.Next()
        },
    }
}
```

### 6. Custom Route Ekleme

```go
func (p *MyPlugin) Routes(router fiber.Router) {
    router.Get("/api/my-plugin/hello", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Hello from MyPlugin!",
        })
    })
}
```

### 7. Migration Ekleme

```go
type CreateMyTable struct{}

func (m *CreateMyTable) Name() string {
    return "create_my_table"
}

func (m *CreateMyTable) Up(db interface{}) error {
    gormDB := db.(*gorm.DB)
    return gormDB.AutoMigrate(&MyModel{})
}

func (m *CreateMyTable) Down(db interface{}) error {
    gormDB := db.(*gorm.DB)
    return gormDB.Migrator().DropTable(&MyModel{})
}

func (p *MyPlugin) Migrations() []plugin.Migration {
    return []plugin.Migration{
        &CreateMyTable{},
    }
}
```

## Frontend Plugin Geliştirme

### 1. Plugin Yapısı

```typescript
// plugins/my-plugin/frontend/index.ts
import { Plugin } from '@/plugins/types';
import { MyCustomField } from './fields/MyCustomField';

export const MyPlugin: Plugin = {
  name: 'my-plugin',
  version: '1.0.0',
  description: 'My custom plugin',
  author: 'Your Name',

  // Custom fields
  fields: [
    {
      type: 'my-custom-field',
      component: MyCustomField,
    },
  ],

  // Initialization
  init: async () => {
    console.log('MyPlugin initialized');
  },

  // Cleanup
  cleanup: async () => {
    console.log('MyPlugin cleaned up');
  },
};
```

### 2. Custom Field Component

```typescript
// plugins/my-plugin/frontend/fields/MyCustomField.tsx
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface MyCustomFieldProps {
  field: {
    key: string;
    label?: string;
    placeholder?: string;
  };
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function MyCustomField({
  field,
  value,
  onChange,
  error,
}: MyCustomFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={field.key}>{field.label || field.key}</Label>
      <Input
        id={field.key}
        value={value || ''}
        onChange={(e) => onChange(e.target.value)}
        placeholder={field.placeholder}
        className={error ? 'border-destructive' : ''}
      />
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  );
}
```

### 3. Plugin Kaydı

```typescript
// main.tsx veya plugin loader
import { pluginRegistry } from '@/plugins/PluginRegistry';
import { MyPlugin } from './plugins/my-plugin/frontend';

pluginRegistry.register(MyPlugin);
await pluginRegistry.initialize();
```

## Frontend Plugin UI Integration

### Plugin Field Registration

Plugin field'larını ana UI ile entegre etmek için şu adımları izleyin:

#### 1. Plugin Oluştur

```typescript
// plugins/my-plugin/frontend/index.ts
import { Plugin } from '@/plugins/types';
import { MyCustomField } from './fields/MyCustomField';

export const MyPlugin: Plugin = {
  name: 'my-plugin',
  version: '1.0.0',
  description: 'My custom plugin',
  author: 'Your Name',

  fields: [
    {
      type: 'my-custom-field',
      component: MyCustomField,
    },
  ],

  init: async () => {
    console.log('MyPlugin initialized');
  },
};
```

#### 2. Plugin'i Register Et

```typescript
// web/src/plugins/index.ts
import { pluginRegistry } from './PluginRegistry';
import { MyPlugin } from '../plugins/my-plugin/frontend';

// Register plugin
pluginRegistry.register(MyPlugin);
```

#### 3. Backend'de Kullan

```go
fields.Custom("My Field", "my_field").
    Type("my-custom-field").
    Label("My Custom Field").
    Placeholder("Enter value")
```

### Field Component Development

Plugin field component'leri şu props'ları almalı:

```typescript
interface FieldComponentProps {
  field: FieldDefinition;  // Field definition (key, label, placeholder, etc.)
  value: any;              // Current value
  onChange: (value: any) => void;  // Value change handler
  onBlur?: () => void;     // Blur handler (optional)
  error?: string;          // Error message (optional)
  disabled?: boolean;      // Disabled state (optional)
  required?: boolean;      // Required state (optional)
}
```

**Örnek Field Component:**

```typescript
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export function MyCustomField({ field, value, onChange, error, disabled }: FieldComponentProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={field.key}>
        {field.label || field.key}
        {field.required && <span className="text-destructive ml-1">*</span>}
      </Label>
      <Input
        id={field.key}
        value={value || ''}
        onChange={(e) => onChange(e.target.value)}
        placeholder={field.placeholder}
        disabled={disabled}
        className={error ? 'border-destructive' : ''}
      />
      {error && <p className="text-xs text-destructive">{error}</p>}
      {field.helpText && <p className="text-xs text-muted-foreground">{field.helpText}</p>}
    </div>
  );
}
```

### View-Specific Variants

Plugin field'ları için 3 farklı variant oluşturabilirsiniz:

- **`my-field-form`**: Form view (düzenlenebilir)
- **`my-field-index`**: Index/table view (read-only, compact)
- **`my-field-detail`**: Detail view (read-only, full display)

**Önemli:** Sadece base type register ederseniz (`my-field`), sistem otomatik olarak form variant'ı (`my-field-form`) oluşturur.

**Örnek - Tüm Variant'lar:**

```typescript
export const MyPlugin: Plugin = {
  name: 'my-plugin',
  version: '1.0.0',
  fields: [
    {
      type: 'my-field-form',
      component: MyFieldForm,
    },
    {
      type: 'my-field-index',
      component: MyFieldIndex,
    },
    {
      type: 'my-field-detail',
      component: MyFieldDetail,
    },
  ],
};
```

### Plugin Initialization Flow

Plugin'ler şu sırayla initialize edilir:

1. **App Başlatılır** (`main.tsx`)
2. **Plugin'ler Import Edilir** (`plugins/index.ts`)
3. **Plugin Registry'ye Kaydedilir** (`pluginRegistry.register()`)
4. **Plugin'ler Initialize Edilir** (`initializePlugins()`)
   - Plugin `init()` hook'ları çağrılır
   - Plugin field'ları toplanır
   - Field'lar `fieldRegistry`'ye kaydedilir
5. **App Render Edilir**
6. **Field'lar Kullanıma Hazır**

### Conflict Resolution

Plugin field'ları core field'larla çakışırsa:

- **Core field'lar önceliklidir**: Plugin field'ı skip edilir
- **Console warning gösterilir**: Conflict log'lanır
- **App çalışmaya devam eder**: Graceful degradation

**Örnek Log:**

```
[Plugin System] Field type 'text' conflicts with core field. Skipping.
[Plugin System] Registered 5 plugin fields, skipped 1 conflicts.
```

## Plugin Yükleme Stratejileri

### 1. Manuel Import (Önerilen)

Compile-time'da plugin'i import edin:

```go
// main.go
import (
    _ "github.com/user/my-plugin"
)

func main() {
    config := panel.Config{...}
    p := panel.New(config)
    p.Start()
}
```

> Not: `panel.New(config)` plugin registry'deki plugin'leri otomatik register+boot eder.

**Avantajlar:**
- Type-safe
- Compile-time kontrol
- Performanslı

### 2. Auto-discovery (Sınırlı)

`AutoDiscover`, şu an sadece `plugin.yaml` descriptor keşfini etkinleştirir.
Runtime'da plugin paketini otomatik import edip boot etmez; plugin yine compile-time import edilmelidir.

```go
config := panel.Config{
    Plugins: panel.PluginConfig{
        AutoDiscover: true,
        Path:         "./plugins",
    },
}
```

**Avantajlar:**
- Plugin klasörlerini descriptor bazlı tarama
- Geçersiz descriptor'ları erken görme

**Dezavantajlar:**
- Runtime plugin package yükleme yok
- Üretimde yine manuel import gerekir

## Best Practices

### 1. Plugin Adlandırma

- Kebab-case kullanın: `my-plugin`
- Benzersiz isimler seçin
- Açıklayıcı isimler kullanın

### 2. Versioning

- Semantic versioning kullanın: `1.0.0`
- Breaking change'lerde major version artırın
- Yeni özellikler için minor version artırın
- Bug fix'ler için patch version artırın

### 3. Error Handling

```go
func (p *MyPlugin) Boot(panel interface{}) error {
    if err := p.initialize(); err != nil {
        return fmt.Errorf("plugin boot failed: %w", err)
    }
    return nil
}
```

### 4. Resource Naming

- Singular form kullanın: `product` (not `products`)
- Lowercase kullanın
- Kebab-case kullanın: `product-category`

### 5. Field Type Naming

- Unique type name kullanın
- Kebab-case kullanın: `my-custom-field`
- Prefix ekleyin: `myplugin-field`

### 6. Testing

```go
func TestMyPlugin(t *testing.T) {
    plugin := &MyPlugin{}

    // Test metadata
    assert.Equal(t, "my-plugin", plugin.Name())
    assert.Equal(t, "1.0.0", plugin.Version())

    // Test resources
    resources := plugin.Resources()
    assert.NotEmpty(t, resources)
}
```

### 7. Documentation

Her plugin için README.md oluşturun:

```markdown
# My Plugin

Plugin açıklaması

## Installation

\`\`\`bash
go get github.com/user/my-plugin
\`\`\`

## Usage

\`\`\`go
import _ "github.com/user/my-plugin"
\`\`\`

## Features

- Feature 1
- Feature 2

## Configuration

...
```

## Örnek Plugin

Tam bir örnek için `plugins/example-plugin` klasörüne bakın:

```bash
plugins/example-plugin/
├── plugin.yaml              # Metadata
├── plugin.go                # Backend implementation
├── README.md                # Documentation
└── frontend/
    ├── index.ts             # Frontend entry
    └── fields/
        └── ExampleField.tsx # Custom field
```

### Çalıştırma

1. Plugin'i import edin:

```go
import _ "github.com/ferdiunal/panel.go/plugins/example-plugin"
```

2. Panel'i başlatın:

```go
func main() {
    config := panel.Config{...}
    p := panel.New(config)
    p.Start()
}
```

3. Frontend plugin'i kaydedin:

```typescript
import { ExamplePlugin } from './plugins/example-plugin/frontend';
import { pluginRegistry } from '@/plugins/PluginRegistry';

pluginRegistry.register(ExamplePlugin);
```

## Troubleshooting

### Plugin Yüklenmiyor

1. `init()` fonksiyonunun çağrıldığından emin olun
2. Import path'in doğru olduğunu kontrol edin
3. Plugin'in `plugin.Register()` ile kaydedildiğini kontrol edin

### Resource Görünmüyor

1. `Resources()` metodunun nil döndürmediğini kontrol edin
2. Resource'un `Visible()` metodunun true döndürdüğünü kontrol edin
3. Plugin paketinin `main.go` içinde blank import ile eklendiğini kontrol edin

### Custom Field Çalışmıyor

1. Field type'ın benzersiz olduğunu kontrol edin
2. Frontend plugin'in kaydedildiğini kontrol edin
3. Component'in doğru props'ları aldığını kontrol edin

## Kaynaklar

- [Example Plugin](../plugins/example-plugin/)
- [Plugin API Reference](./API.md)
- [Field Development Guide](./FIELD_DEVELOPMENT.md)

## Destek

Sorularınız için:
- GitHub Issues: https://github.com/ferdiunal/panel.go/issues
- Documentation: https://panel-go.dev/docs
