# Plugin CLI Komutları Referansı

Panel.go plugin CLI komutlarının detaylı referans dökümanı.

## İçindekiler

- [plugin create](#plugin-create)
- [plugin add](#plugin-add)
- [plugin remove](#plugin-remove)
- [plugin list](#plugin-list)
- [plugin build](#plugin-build)

## plugin create

Yeni plugin oluşturur. Backend ve frontend dosyalarını scaffold eder, workspace config'i günceller ve build alır.

### Kullanım

```bash
panel plugin create <plugin-name> [flags]
```

### Parametreler

- `<plugin-name>`: Plugin adı (required, kebab-case önerilir)

### Flags

| Flag | Tip | Default | Açıklama |
|------|-----|---------|----------|
| `--path` | string | `./plugins` | Plugin dizini |
| `--no-frontend` | bool | `false` | Frontend scaffold etme |
| `--no-build` | bool | `false` | Otomatik build yapma |
| `--with-example` | bool | `false` | Tüm relationship türlerini içeren comprehensive örnek oluştur |

### Örnekler

**Temel Kullanım:**
```bash
panel plugin create importer
```

**Frontend Olmadan:**
```bash
panel plugin create analytics --no-frontend
```

**Custom Path:**
```bash
panel plugin create exporter --path ./custom-plugins
```

**Comprehensive Example (Tüm Relationship Türleri):**
```bash
panel plugin create example --with-example
```

Bu komut, tüm GORM relationship türlerini (BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo, MorphToMany) içeren tam özellikli bir örnek plugin oluşturur. 9 entity ve 7 resource ile production-ready bir yapı sağlar.

**Build Olmadan:**
```bash
panel plugin create logger --no-build
```

### Oluşturulan Dosyalar

**Backend:**
```
plugins/importer/
├── plugin.go          # Plugin implementation
└── plugin.yaml        # Plugin metadata
```

**Frontend (--no-frontend değilse):**
```
plugins/importer/
└── frontend/
    ├── index.ts       # Plugin export
    ├── package.json   # Plugin package
    ├── tsconfig.json  # TypeScript config
    └── fields/        # Custom field'lar
        └── .gitkeep
```

### İşlem Adımları

1. Plugin dizini oluşturulur: `plugins/<name>/`
2. Backend dosyaları oluşturulur: `plugin.go`, `plugin.yaml`
3. Frontend dosyaları oluşturulur (eğer `--no-frontend` değilse)
4. web-ui clone edilir (ilk kez)
5. Workspace config güncellenir: `web-ui/pnpm-workspace.yaml`
6. Plugin symlink oluşturulur: `web-ui/plugins/<name>`
7. Build alınır (eğer `--no-build` değilse): `assets/ui/`

### Çıktı

```
🚀 Plugin oluşturuluyor: importer

✓ Plugin dizini oluşturuldu: plugins/importer/
✓ Backend dosyaları oluşturuldu: plugin.go, plugin.yaml
✓ Frontend dosyaları oluşturuldu: index.ts, package.json, tsconfig.json
✓ web-ui clone edildi: web-ui
✓ Workspace config güncellendi: web-ui/pnpm-workspace.yaml
✓ Workspace reference oluşturuldu: web-ui/plugins/importer
✓ UI build alınıyor...
✓ Build tamamlandı: assets/ui/

✅ Plugin 'importer' başarıyla oluşturuldu!

Sonraki adımlar:
  1. Backend implement et: plugins/importer/plugin.go
  2. Frontend field'ları ekle: plugins/importer/frontend/fields/
  3. Plugin'i import et: import _ "your-module/plugins/importer"
  4. Rebuild: panel plugin build
```

## plugin add

Git repository'den plugin ekler. Repository'yi clone eder, validate eder ve workspace'e entegre eder.

### Kullanım

```bash
panel plugin add <git-url> [flags]
```

### Parametreler

- `<git-url>`: Git repository URL'si (required)

### Flags

| Flag | Tip | Default | Açıklama |
|------|-----|---------|----------|
| `--path` | string | `./plugins` | Plugin dizini |
| `--branch` | string | `main` | Git branch |
| `--no-build` | bool | `false` | Otomatik build yapma |

### Örnekler

**GitHub'dan Ekle:**
```bash
panel plugin add https://github.com/user/analytics-plugin
```

**Belirli Branch:**
```bash
panel plugin add https://github.com/user/exporter --branch develop
```

**Custom Path:**
```bash
panel plugin add https://github.com/user/logger --path ./custom-plugins
```

### İşlem Adımları

1. Git URL parse edilir, plugin adı çıkarılır
2. Plugin clone edilir: `git clone <url> plugins/<name>`
3. Plugin validate edilir: `plugin.yaml`, `plugin.go` kontrol edilir
4. web-ui clone edilir (ilk kez)
5. Frontend varsa workspace config güncellenir
6. Plugin symlink oluşturulur (frontend varsa)
7. Build alınır (eğer `--no-build` değilse)

### Validation

Plugin geçerli olması için:
- `plugin.yaml` dosyası olmalı
- `plugin.go` dosyası olmalı
- Metadata geçerli olmalı

### Çıktı

```
📦 Plugin ekleniyor: https://github.com/user/analytics-plugin

✓ Plugin clone ediliyor: https://github.com/user/analytics-plugin
✓ Plugin clone edildi: plugins/analytics-plugin
✓ Plugin validate edildi
✓ web-ui clone edildi: web-ui
✓ Workspace config güncellendi
✓ Workspace reference oluşturuldu: web-ui/plugins/analytics-plugin
✓ UI build alınıyor...
✓ Build tamamlandı: assets/ui/

✅ Plugin 'analytics-plugin' başarıyla eklendi!

Sonraki adımlar:
  1. Plugin'i import et: import _ "your-module/plugins/analytics-plugin"
  2. Rebuild: panel plugin build
```

## plugin remove

Plugin'i siler. Workspace reference'ı kaldırır, plugin dosyalarını siler ve build alır.

### Kullanım

```bash
panel plugin remove <plugin-name> [flags]
```

### Parametreler

- `<plugin-name>`: Plugin adı (required)

### Flags

| Flag | Tip | Default | Açıklama |
|------|-----|---------|----------|
| `--path` | string | `./plugins` | Plugin dizini |
| `--keep-files` | bool | `false` | Plugin dosyalarını silme |
| `--no-build` | bool | `false` | Otomatik build yapma |

### Örnekler

**Temel Kullanım:**
```bash
panel plugin remove importer
```

**Dosyaları Koru:**
```bash
panel plugin remove analytics --keep-files
```

**Build Olmadan:**
```bash
panel plugin remove exporter --no-build
```

### İşlem Adımları

1. Plugin varlığı kontrol edilir
2. Workspace reference silinir: `web-ui/plugins/<name>`
3. Workspace config güncellenir
4. Plugin dosyaları silinir (eğer `--keep-files` değilse)
5. Build alınır (eğer `--no-build` değilse)

### Çıktı

```
🗑️  Plugin siliniyor: importer

✓ Workspace reference silindi: web-ui/plugins/importer
✓ Workspace config güncellendi
✓ Plugin dosyaları silindi: plugins/importer
✓ UI build alınıyor...
✓ Build tamamlandı: assets/ui/

✅ Plugin 'importer' başarıyla silindi!
```

## plugin list

Yüklü plugin'leri listeler. Plugin metadata'sını okur ve tablo formatında gösterir.

### Kullanım

```bash
panel plugin list [flags]
```

### Flags

| Flag | Tip | Default | Açıklama |
|------|-----|---------|----------|
| `--path` | string | `./plugins` | Plugin dizini |
| `--json` | bool | `false` | JSON output |

### Örnekler

**Tablo Format:**
```bash
panel plugin list
```

**JSON Format:**
```bash
panel plugin list --json
```

**Custom Path:**
```bash
panel plugin list --path ./custom-plugins
```

### Çıktı (Tablo)

```
Yüklü Plugin'ler:

NAME              VERSION    AUTHOR         FRONTEND    STATUS
importer          1.0.0      Panel.go Team  Yes         Active
analytics-plugin  1.2.0      John Doe       Yes         Active
logger            1.0.0      Panel.go Team  No          Active

Toplam: 3 plugin
```

### Çıktı (JSON)

```json
[
  {
    "name": "importer",
    "version": "1.0.0",
    "author": "Panel.go Team",
    "description": "CSV/XLSX import plugin",
    "has_frontend": true,
    "valid": true,
    "path": "plugins/importer"
  },
  {
    "name": "analytics-plugin",
    "version": "1.2.0",
    "author": "John Doe",
    "description": "Analytics dashboard plugin",
    "has_frontend": true,
    "valid": true,
    "path": "plugins/analytics-plugin"
  }
]
```

## plugin build

UI build alır. web-ui'yi build eder ve output'u `assets/ui/`'ye kopyalar.

### Kullanım

```bash
panel plugin build [flags]
```

### Flags

| Flag | Tip | Default | Açıklama |
|------|-----|---------|----------|
| `--dev` | bool | `false` | Development build (no minification) |
| `--watch` | bool | `false` | Watch mode (continuous build) |

### Örnekler

**Production Build:**
```bash
panel plugin build
```

**Development Build:**
```bash
panel plugin build --dev
```

**Watch Mode:**
```bash
panel plugin build --watch
```

### İşlem Adımları

1. web-ui varlığı kontrol edilir (yoksa clone edilir)
2. Package manager detect edilir (pnpm > npm)
3. Dependencies yüklenir: `pnpm install`
4. Build alınır:
   - Production: `pnpm build`
   - Development: `pnpm build --mode development`
   - Watch: `pnpm dev`
5. Output kopyalanır: `web-ui/dist/` → `assets/ui/`

### Çıktı (Production)

```
🔨 UI build alınıyor...

✓ Package manager: pnpm
✓ Dependencies yükleniyor...
✓ Dependencies yüklendi
✓ Build alınıyor (build)...
✓ Build tamamlandı
✓ Output kopyalanıyor: web-ui/dist -> assets/ui
✓ Output kopyalandı: assets/ui

✅ Build başarıyla tamamlandı!

Build output: assets/ui/
```

### Çıktı (Watch Mode)

```
🔨 UI build alınıyor...

✓ Package manager: pnpm
✓ Dependencies yükleniyor...
✓ Dependencies yüklendi
✓ Watch mode başlatılıyor...
  (Ctrl+C ile durdurun)

VITE v5.0.0  ready in 1234 ms

  ➜  Local:   http://localhost:5177/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```

## Yaygın Kullanım Senaryoları

### Yeni Plugin Geliştirme

```bash
# 1. Plugin oluştur
panel plugin create my-plugin

# 2. Backend implement et
# Edit: plugins/my-plugin/plugin.go

# 3. Frontend field ekle
# Edit: plugins/my-plugin/frontend/fields/MyField.tsx

# 4. Plugin'i import et
# Edit: main.go
# import _ "your-module/plugins/my-plugin"

# 5. Build ve test
panel plugin build
go run main.go
```

### Mevcut Plugin Ekleme

```bash
# 1. Plugin ekle
panel plugin add https://github.com/user/analytics-plugin

# 2. Plugin'i import et
# Edit: main.go
# import _ "your-module/plugins/analytics-plugin"

# 3. Başlat
go run main.go
```

### Development Workflow

```bash
# Terminal 1: Watch mode
cd web-ui
pnpm dev

# Terminal 2: Panel başlat
go run main.go

# Panel otomatik olarak Vite dev server'a proxy eder
# http://localhost:8787 -> http://localhost:5177
```

### Production Build

```bash
# 1. Build al
panel plugin build

# 2. Binary oluştur
go build -o panel-app

# 3. Başlat
./panel-app
```

## Troubleshooting

### web-ui Clone Edilemiyor

**Hata:**
```
Error: git clone hatası: repository not found
```

**Çözüm:**
```bash
# Manuel clone
git clone https://github.com/ferdiunal/panel.web web-ui

# Sonra build
panel plugin build
```

### Build Hatası

**Hata:**
```
Error: build hatası: command not found: pnpm
```

**Çözüm:**
```bash
# pnpm kur
npm install -g pnpm

# Veya npm kullan
cd web-ui
npm install
npm run build
```

### Plugin Bulunamıyor

**Hata:**
```
Error: plugin bulunamadı: importer
```

**Çözüm:**
```bash
# Plugin listesini kontrol et
panel plugin list

# Plugin path'i kontrol et
ls -la plugins/
```

## İleri Okuma

- [Plugin Sistemi](./PLUGIN_SYSTEM.md) - Genel bakış
- [Plugin Geliştirme](./PLUGIN_DEVELOPMENT.md) - Backend ve frontend
- [Örnekler](./PLUGIN_EXAMPLES.md) - Gerçek dünya örnekleri
