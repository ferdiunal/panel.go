# Panel.go Plugin Sistemi

Panel.go plugin sistemi, admin panel'inize özel özellikler eklemek için güçlü ve esnek bir yapı sunar. Bu döküman, plugin sisteminin nasıl çalıştığını ve nasıl kullanılacağını açıklar.

## İçindekiler

- [Genel Bakış](#genel-bakış)
- [Hızlı Başlangıç](#hızlı-başlangıç)
- [Mimari](#mimari)
- [CLI Komutları](#cli-komutları)
- [İleri Okuma](#ileri-okuma)

## Genel Bakış

### Plugin Sistemi Nedir?

Panel.go plugin sistemi, admin panel'inize yeni özellikler eklemek için modüler bir yaklaşım sunar. Plugin'ler:

- **Backend**: Go ile yazılır, Panel.go'nun tüm özelliklerine erişebilir
- **Frontend**: TypeScript/React ile yazılır, custom field'lar ve UI component'leri ekleyebilir
- **Type-Safe**: Compile-time type checking ile güvenli geliştirme
- **Hot Reload**: Development mode'da anlık değişiklik görme

### Neden Plugin Sistemi?

**Geleneksel Yaklaşım (Laravel Nova):**
- Runtime dynamic loading (karmaşık, type-safety kaybı)
- Network latency (plugin'ler runtime'da yüklenir)
- Güvenlik riskleri (runtime code injection)

**Panel.go Yaklaşımı:**
- Compile-time integration (basit, type-safe)
- Zero runtime overhead (plugin'ler build'e dahil)
- Project-specific UI build (her proje kendi UI'ını build eder)

### Temel Kavramlar

**1. Plugin Dizini**
```
plugins/
├── importer/              # Plugin adı
│   ├── plugin.go          # Backend plugin
│   ├── plugin.yaml        # Plugin metadata
│   └── frontend/          # Frontend plugin (opsiyonel)
│       ├── index.ts       # Plugin export
│       ├── package.json   # Plugin package
│       └── fields/        # Custom field'lar
```

**2. Workspace Pattern**
Plugin'ler npm/pnpm workspace olarak entegre edilir:
```
web-ui/
├── plugins/
│   └── importer -> ../../plugins/importer/frontend
├── pnpm-workspace.yaml
└── package.json
```

**3. Build Output**
UI build output `assets/ui/` dizinine kopyalanır:
```
assets/
└── ui/
    ├── index.html
    └── assets/
        ├── index-[hash].js
        └── index-[hash].css
```

## Hızlı Başlangıç

### Ön Gereksinimler

- Go 1.25+
- Node.js 18+ (pnpm önerilir)
- Git

### 1. Panel CLI Kurulumu

```bash
cd /path/to/panel.go
go install ./cmd/panel
```

### 2. İlk Plugin'inizi Oluşturun

```bash
cd examples/cargo.go
panel plugin create importer
```

Bu komut:
- ✓ `plugins/importer/` dizini oluşturur
- ✓ Backend dosyaları oluşturur (plugin.go, plugin.yaml)
- ✓ Frontend dosyaları oluşturur (index.ts, package.json, tsconfig.json)
- ✓ web-ui'yi clone eder (ilk kez)
- ✓ Workspace config'i günceller
- ✓ UI build alır → `assets/ui/`

### 3. Plugin'i Import Edin

```go
// main.go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    _ "your-module/plugins/importer" // Plugin'i import et
)

func main() {
    // Panel başlat
    p := panel.New(config)
    p.Start()
}
```

### 4. Uygulamayı Başlatın

```bash
go run main.go
```

Panel başladığında `assets/ui/` dizininden UI serve edilir ve plugin'iniz aktif olur.

## Mimari

### Asset Serving Priority

Panel.go, UI dosyalarını şu öncelik sırasıyla serve eder:

**1. Priority: assets/ui/ (Project-Specific Build)**
```bash
panel plugin build  # Build alınır
# Output: assets/ui/
```
- Plugin build output'u
- Production-ready
- Her proje kendi UI'ını build eder

**2. Priority: pkg/panel/ui/ (SDK Development)**
```bash
# Development mode
cd pkg/panel/ui
npm run dev
```
- SDK geliştiricileri için
- Hot reload
- Sadece development mode'da

**3. Priority: Embedded Assets (Backward Compatibility)**
```go
//go:embed ui/*
var assetsFS embed.FS
```
- Fallback
- SDK kullanıcıları için
- Binary'ye gömülü

### Plugin Lifecycle

`panel.New(config)` sırasında plugin registry okunur ve aşağıdaki sıra uygulanır:

```
1. Register (init)
   ↓
2. Boot (app startup)
   ↓
3. Resources/Pages/Middleware eklenir
   ↓
4. Routes kaydedilir
   ↓
5. Migrations çalıştırılır
```

**Backend Plugin Interface:**
```go
type Plugin interface {
    Name() string
    Version() string
    Author() string
    Description() string

    Register(panel interface{}) error
    Boot(panel interface{}) error

    Resources() []resource.Resource
    Pages() []plugin.Page
    Middleware() []fiber.Handler
    Routes(router fiber.Router)
    Migrations() []plugin.Migration
}
```

**Frontend Plugin Interface:**
```typescript
interface Plugin {
  name: string;
  version: string;
  description: string;
  author: string;
  fields?: FieldDefinition[];
}
```

## CLI Komutları

### plugin create

Yeni plugin oluşturur.

```bash
panel plugin create <plugin-name> [flags]
```

**Flags:**
- `--path <path>`: Plugin dizini (default: `./plugins`)
- `--no-frontend`: Frontend scaffold etme
- `--no-build`: Otomatik build yapma

**Örnek:**
```bash
panel plugin create importer
panel plugin create analytics --no-frontend
panel plugin create exporter --path ./custom-plugins
```

### plugin add

Git repository'den plugin ekler.

```bash
panel plugin add <git-url> [flags]
```

**Flags:**
- `--path <path>`: Plugin dizini (default: `./plugins`)
- `--branch <branch>`: Git branch (default: `main`)
- `--no-build`: Otomatik build yapma

**Örnek:**
```bash
panel plugin add https://github.com/user/analytics-plugin
panel plugin add https://github.com/user/exporter --branch develop
```

### plugin remove

Plugin'i siler.

```bash
panel plugin remove <plugin-name> [flags]
```

**Flags:**
- `--path <path>`: Plugin dizini (default: `./plugins`)
- `--keep-files`: Plugin dosyalarını silme
- `--no-build`: Otomatik build yapma

**Örnek:**
```bash
panel plugin remove importer
panel plugin remove analytics --keep-files
```

### plugin list

Yüklü plugin'leri listeler.

```bash
panel plugin list [flags]
```

**Flags:**
- `--path <path>`: Plugin dizini (default: `./plugins`)
- `--json`: JSON output

**Örnek:**
```bash
panel plugin list
panel plugin list --json
```

**Output:**
```
Yüklü Plugin'ler:

NAME              VERSION    AUTHOR         FRONTEND    STATUS
importer          1.0.0      Panel.go Team  Yes         Active
analytics-plugin  1.2.0      John Doe       Yes         Active

Toplam: 2 plugin
```

### plugin build

UI build alır.

```bash
panel plugin build [flags]
```

**Flags:**
- `--dev`: Development build (no minification)
- `--watch`: Watch mode (continuous build)

**Örnek:**
```bash
panel plugin build              # Production build
panel plugin build --dev        # Development build
panel plugin build --watch      # Watch mode
```

## İleri Okuma

- [CLI Komutları Referansı](./PLUGIN_CLI.md) - Tüm CLI komutlarının detaylı açıklaması
- [Plugin Geliştirme Rehberi](./PLUGIN_DEVELOPMENT.md) - Backend ve frontend plugin geliştirme
- [Örnekler](./PLUGIN_EXAMPLES.md) - Gerçek dünya örnekleri
- [Troubleshooting](./PLUGIN_TROUBLESHOOTING.md) - Yaygın sorunlar ve çözümleri

## Destek

- GitHub Issues: https://github.com/ferdiunal/panel.go/issues
- Dökümanlar: https://github.com/ferdiunal/panel.go/docs
