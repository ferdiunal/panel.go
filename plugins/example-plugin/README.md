# Example Plugin

Panel.go için örnek plugin implementasyonu. Plugin sisteminin tüm özelliklerini gösterir.

## Özellikler

- ✅ Custom resource (Example Resource)
- ✅ Custom field (Example Field)
- ✅ Custom middleware (Request logging)
- ✅ Custom route (`/api/example-plugin/hello`)
- ✅ Database migration (Example table)
- ✅ Frontend integration

## Kurulum

### Backend

1. Plugin'i import edin:

```go
import _ "github.com/ferdiunal/panel.go/plugins/example-plugin"
```

2. Panel'i başlatın:

```go
func main() {
    config := panel.Config{
        Database: panel.DatabaseConfig{Instance: db},
        Server:   panel.ServerConfig{Host: "localhost", Port: "8080"},
    }

    p := panel.New(config)
    p.Start()
}
```

> Not: `panel.New(config)` plugin'i registry'den otomatik boot eder.

### Frontend

1. Plugin'i kaydedin:

```typescript
// main.tsx
import { ExamplePlugin } from './plugins/example-plugin/frontend';
import { pluginRegistry } from '@/plugins/PluginRegistry';

pluginRegistry.register(ExamplePlugin);
await pluginRegistry.initialize();
```

## Kullanım

### Example Resource

Plugin yüklendikten sonra, yan menüde "Plugin Examples" grubunda "Example Resource" görünecektir.

**Özellikler:**
- CRUD işlemleri (Create, Read, Update, Delete)
- Arama ve filtreleme
- Sıralama
- Sayfalama

**Alanlar:**
- Name (Text)
- Description (Textarea)
- Active (Switch)
- Created At (DateTime)
- Updated At (DateTime)

### Example Field

Custom field component örneği. Backend'de şu şekilde kullanılır:

```go
fields.Custom("Example Field", "example_field").
    Type("example-field").
    Label("Example Field").
    Placeholder("Enter value").
    HelpText("This is an example custom field")
```

Frontend'de otomatik olarak `ExampleField` component'i render edilir.

### Custom Route

Plugin, custom bir API endpoint sağlar:

```bash
curl http://localhost:8080/api/example-plugin/hello
```

**Response:**
```json
{
  "message": "Hello from ExamplePlugin!",
  "version": "1.0.0"
}
```

### Custom Middleware

Plugin, tüm istekleri loglayan bir middleware ekler:

```
ExamplePlugin Middleware: GET /api/resource/example
ExamplePlugin Middleware: POST /api/resource/example
```

## Yapı

```
plugins/example-plugin/
├── plugin.yaml              # Plugin metadata
├── plugin.go                # Backend implementation
├── README.md                # Bu dosya
└── frontend/
    ├── index.ts             # Frontend plugin entry
    └── fields/
        └── ExampleField.tsx # Custom field component
```

## Plugin Metadata

```yaml
name: example-plugin
version: 1.0.0
author: Panel.go Team
description: Example plugin demonstrating plugin system capabilities
enabled: true
entry: plugin.go
```

## Backend Implementation

### Plugin Struct

```go
type ExamplePlugin struct {
    plugin.BasePlugin
}
```

### Metadata

```go
func (p *ExamplePlugin) Name() string        { return "example-plugin" }
func (p *ExamplePlugin) Version() string     { return "1.0.0" }
func (p *ExamplePlugin) Author() string      { return "Panel.go Team" }
func (p *ExamplePlugin) Description() string { return "Example plugin..." }
```

### Lifecycle

```go
func (p *ExamplePlugin) Register(panel interface{}) error {
    // Plugin kaydı
    return nil
}

func (p *ExamplePlugin) Boot(panel interface{}) error {
    // Plugin başlatma
    return nil
}
```

### Capabilities

```go
// Resources
func (p *ExamplePlugin) Resources() []resource.Resource {
    return []resource.Resource{NewExampleResource()}
}

// Middleware
func (p *ExamplePlugin) Middleware() []fiber.Handler {
    return []fiber.Handler{/* ... */}
}

// Routes
func (p *ExamplePlugin) Routes(router fiber.Router) {
    router.Get("/api/example-plugin/hello", /* ... */)
}

// Migrations
func (p *ExamplePlugin) Migrations() []plugin.Migration {
    return []plugin.Migration{&CreateExampleTable{}}
}
```

## Frontend Implementation

### Plugin Definition

```typescript
export const ExamplePlugin: Plugin = {
  name: 'example-plugin',
  version: '1.0.0',
  description: 'Example plugin demonstrating plugin system',
  author: 'Panel.go Team',

  fields: [
    {
      type: 'example-field',
      component: ExampleField,
    },
  ],

  init: async () => {
    console.log('ExamplePlugin initialized');
  },

  cleanup: async () => {
    console.log('ExamplePlugin cleaned up');
  },
};
```

### Custom Field Component

```typescript
export function ExampleField({
  field,
  value,
  onChange,
  error,
  disabled,
}: ExampleFieldProps) {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{field.label || field.key}</CardTitle>
        {field.helpText && <CardDescription>{field.helpText}</CardDescription>}
      </CardHeader>
      <CardContent>
        <Input
          value={value || ''}
          onChange={(e) => onChange(e.target.value)}
          placeholder={field.placeholder}
          disabled={disabled}
        />
      </CardContent>
    </Card>
  );
}
```

## Database Schema

Plugin, aşağıdaki tabloyu oluşturur:

```sql
CREATE TABLE example_models (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Test

### Backend Test

```bash
cd plugins/example-plugin
go test -v
```

### Frontend Test

```bash
cd plugins/example-plugin/frontend
npm test
```

## Geliştirme

### Yeni Özellik Ekleme

1. Backend'de yeni metod ekleyin:

```go
func (p *ExamplePlugin) NewFeature() {
    // Implementation
}
```

2. Frontend'de yeni component ekleyin:

```typescript
export function NewComponent() {
    // Implementation
}
```

3. Plugin'i yeniden build edin:

```bash
go build
```

### Debug

Backend debug için:

```go
func (p *ExamplePlugin) Boot(panel interface{}) error {
    fmt.Println("ExamplePlugin: Boot called")
    fmt.Printf("Resources: %d\n", len(p.Resources()))
    return nil
}
```

Frontend debug için:

```typescript
init: async () => {
    console.log('ExamplePlugin initialized');
    console.log('Fields:', ExamplePlugin.fields);
},
```

## Sorun Giderme

### Plugin Yüklenmiyor

1. `init()` fonksiyonunun çağrıldığından emin olun
2. Import path'in doğru olduğunu kontrol edin
3. Plugin paketinin uygulamada blank import edildiğini doğrulayın

### Resource Görünmüyor

1. `Resources()` metodunun nil döndürmediğini kontrol edin
2. Resource'un `Visible()` metodunun true döndürdüğünü kontrol edin

### Custom Field Çalışmıyor

1. Field type'ın benzersiz olduğunu kontrol edin (`example-field`)
2. Frontend plugin'in kaydedildiğini kontrol edin
3. Component'in doğru props'ları aldığını kontrol edin

## Lisans

MIT License

## Destek

Sorularınız için:
- GitHub Issues: https://github.com/ferdiunal/panel.go/issues
- Documentation: https://panel-go.dev/docs/plugins
