# Plugin Troubleshooting

Panel.go plugin sistemi için yaygın sorunlar ve çözümleri.

## İçindekiler

- [Build Hataları](#build-hataları)
- [Plugin Yükleme Sorunları](#plugin-yükleme-sorunları)
- [Frontend Entegrasyon](#frontend-entegrasyon)
- [Git ve Workspace](#git-ve-workspace)
- [Performance](#performance)

## Build Hataları

### web-ui Clone Edilemiyor

**Hata:**
```
Error: git clone hatası: repository not found
```

**Neden:**
- Git kurulu değil
- Network bağlantısı yok
- Repository erişim sorunu

**Çözüm:**
```bash
# Git kurulu mu kontrol et
git --version

# Manuel clone
git clone https://github.com/ferdiunal/panel.web web-ui

# Sonra build
panel plugin build
```

### pnpm/npm Bulunamıyor

**Hata:**
```
Error: command not found: pnpm
```

**Neden:**
- Node.js kurulu değil
- pnpm kurulu değil

**Çözüm:**
```bash
# Node.js kur (18+)
# macOS
brew install node

# pnpm kur
npm install -g pnpm

# Veya npm kullan
cd web-ui
npm install
npm run build
```

### Build Timeout

**Hata:**
```
Error: build timeout after 10 minutes
```

**Neden:**
- Yavaş internet bağlantısı
- Çok fazla dependency
- Sistem kaynakları yetersiz

**Çözüm:**
```bash
# Dependencies'i manuel yükle
cd web-ui
pnpm install --no-frozen-lockfile

# Build al
pnpm build

# Output'u kopyala
cp -r dist ../assets/ui
```

### TypeScript Hataları

**Hata:**
```
Error: Type 'string' is not assignable to type 'number'
```

**Neden:**
- Plugin field'ında type mismatch
- Props interface yanlış

**Çözüm:**
```typescript
// ❌ Yanlış
interface MyFieldProps {
  value: number;
}

// ✅ Doğru
interface MyFieldProps {
  value: any; // veya string | number
}
```

## Plugin Yükleme Sorunları

### Plugin Bulunamıyor

**Hata:**
```
Error: plugin bulunamadı: importer
```

**Neden:**
- Plugin dizini yanlış
- Plugin dosyaları eksik

**Çözüm:**
```bash
# Plugin listesini kontrol et
panel plugin list

# Plugin dizinini kontrol et
ls -la plugins/

# Plugin path'i belirt
panel plugin list --path ./plugins
```

### Plugin Kaydedilmedi

**Hata:**
```
Warning: Plugin 'importer' not registered
```

**Neden:**
- `init()` fonksiyonu çağrılmadı
- Import path yanlış
- `plugin.Register()` çağrılmadı

**Çözüm:**
```go
// ✅ Doğru
package importer

import "github.com/ferdiunal/panel.go/pkg/plugin"

func init() {
    plugin.Register(&Plugin{})
}

// main.go
import _ "your-module/plugins/importer"
```

### Plugin Boot Hatası

**Hata:**
```
Error: plugin 'importer' boot failed: database error
```

**Neden:**
- Database bağlantısı yok
- Migration hatası
- Type assertion hatası

**Çözüm:**
```go
func (p *Plugin) Boot(panel interface{}) error {
    // Type assertion
    panelApp, ok := panel.(*panel.Panel)
    if !ok {
        return fmt.Errorf("invalid panel type")
    }

    // Database kontrol
    if panelApp.Db == nil {
        return fmt.Errorf("database not initialized")
    }

    // Migration
    if err := panelApp.Db.AutoMigrate(&Import{}); err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }

    return nil
}
```

## Frontend Entegrasyon

### Custom Field Görünmüyor

**Hata:**
Field backend'de tanımlı ama frontend'de render edilmiyor.

**Neden:**
- Plugin kaydedilmedi
- Field type yanlış
- Component export edilmedi

**Çözüm:**
```typescript
// 1. Plugin'i kaydet
// web-ui/src/plugins/index.ts
import { ImporterPlugin } from '../../plugins/importer/frontend';
pluginRegistry.register(ImporterPlugin);

// 2. Field type'ı kontrol et
// Backend
fields.Custom("File", "file").Type("import-field")

// Frontend
export const ImporterPlugin: Plugin = {
  fields: [
    { type: 'import-field', component: ImportField },
  ],
};
```

### Field Props Hatası

**Hata:**
```
Error: Cannot read property 'label' of undefined
```

**Neden:**
- Props interface yanlış
- Field data eksik

**Çözüm:**
```typescript
// ✅ Doğru
interface ImportFieldProps {
  field: any;
  name: string;
  value: any;
  onChange: (value: any) => void;
  error?: string; // Optional
}

export const ImportField: React.FC<ImportFieldProps> = ({
  field,
  name,
  value,
  onChange,
  error,
}) => {
  // field.label optional olabilir
  const label = field?.label || name;

  return (
    <FieldLayout name={name} label={label} error={error}>
      {/* ... */}
    </FieldLayout>
  );
};
```

### Hot Reload Çalışmıyor

**Hata:**
Kod değişiklikleri yansımıyor.

**Neden:**
- Watch mode çalışmıyor
- Browser cache
- Symlink sorunu

**Çözüm:**
```bash
# 1. Watch mode'u yeniden başlat
cd web-ui
pnpm dev

# 2. Browser cache'i temizle
# Chrome: Cmd+Shift+R (macOS) / Ctrl+Shift+R (Windows)

# 3. Symlink'i kontrol et
ls -la web-ui/plugins/
# importer -> ../../plugins/importer/frontend olmalı

# 4. Symlink'i yeniden oluştur
rm web-ui/plugins/importer
ln -s ../../plugins/importer/frontend web-ui/plugins/importer
```

## Git ve Workspace

### Workspace Config Bozuk

**Hata:**
```
Error: workspace config parse edilemedi
```

**Neden:**
- YAML syntax hatası
- Geçersiz path

**Çözüm:**
```yaml
# web-ui/pnpm-workspace.yaml
packages:
  - "plugins/*"
  - "../plugins/*/frontend"

# Syntax kontrol
cat web-ui/pnpm-workspace.yaml | pnpm exec yaml-validator
```

### Symlink Oluşturulamıyor

**Hata:**
```
Error: symlink oluşturulamadı: permission denied
```

**Neden:**
- Dosya sistemi izinleri
- Windows'ta symlink desteği yok

**Çözüm:**
```bash
# macOS/Linux
ln -s ../../plugins/importer/frontend web-ui/plugins/importer

# Windows (Admin gerekli)
mklink /D web-ui\plugins\importer ..\..\plugins\importer\frontend

# Veya junction kullan (Admin gerektirmez)
mklink /J web-ui\plugins\importer ..\..\plugins\importer\frontend
```

### Git Submodule Çakışması

**Hata:**
```
Error: 'web-ui' already exists in the index
```

**Neden:**
- web-ui git submodule olarak eklenmiş
- .gitmodules dosyası var

**Çözüm:**
```bash
# Submodule'ü kaldır
git submodule deinit -f web-ui
git rm -f web-ui
rm -rf .git/modules/web-ui

# Normal dizin olarak ekle
git clone https://github.com/ferdiunal/panel.web web-ui
echo "web-ui/" >> .gitignore
```

## Performance

### Build Çok Yavaş

**Sorun:**
Build 5+ dakika sürüyor.

**Neden:**
- Çok fazla plugin
- Büyük dependencies
- Disk I/O yavaş

**Çözüm:**
```bash
# 1. Cache'i temizle
cd web-ui
rm -rf node_modules .pnpm-store
pnpm install

# 2. Production build yerine dev build
panel plugin build --dev

# 3. Watch mode kullan (development)
cd web-ui
pnpm dev
```

### Bundle Size Çok Büyük

**Sorun:**
assets/ui/ dizini 10+ MB.

**Neden:**
- Gereksiz dependencies
- Source maps dahil
- Minification yok

**Çözüm:**
```bash
# 1. Production build al
panel plugin build

# 2. Bundle analyzer kullan
cd web-ui
pnpm add -D rollup-plugin-visualizer
pnpm build

# 3. Gereksiz dependencies'i kaldır
pnpm remove unused-package
```

### Memory Leak

**Sorun:**
Panel uzun süre çalıştıktan sonra yavaşlıyor.

**Neden:**
- Plugin cleanup yapılmıyor
- Event listener'lar temizlenmiyor

**Çözüm:**
```typescript
// Plugin cleanup
export const MyPlugin: Plugin = {
  name: 'my-plugin',

  init: async () => {
    // Setup
    window.addEventListener('resize', handleResize);
  },

  cleanup: async () => {
    // Cleanup
    window.removeEventListener('resize', handleResize);
  },
};
```

## Debugging

### Debug Mode

```bash
# Backend debug
go run -race main.go

# Frontend debug
cd web-ui
pnpm dev --debug
```

### Log Seviyesi

```go
// Backend logging
import "log"

func (p *Plugin) Boot(panel interface{}) error {
    log.Printf("[%s] Booting plugin...", p.Name())
    // ...
    log.Printf("[%s] Plugin booted successfully", p.Name())
    return nil
}
```

```typescript
// Frontend logging
export const MyPlugin: Plugin = {
  init: async () => {
    console.log('[MyPlugin] Initializing...');
    // ...
    console.log('[MyPlugin] Initialized successfully');
  },
};
```

### Breakpoint Debugging

```bash
# Backend (Delve)
dlv debug main.go

# Frontend (Chrome DevTools)
# 1. cd web-ui && pnpm dev
# 2. Chrome'da F12
# 3. Sources tab -> Breakpoint ekle
```

## Yaygın Hatalar ve Çözümleri

### "Cannot find module"

```bash
# Node modules'ı yeniden yükle
cd web-ui
rm -rf node_modules
pnpm install
```

### "Port already in use"

```bash
# Port'u kullanan process'i bul
lsof -i :5177

# Process'i öldür
kill -9 <PID>

# Veya farklı port kullan
cd web-ui
pnpm dev --port 5178
```

### "ENOSPC: System limit for number of file watchers reached"

```bash
# Linux
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# macOS (genellikle sorun olmaz)
```

### "Permission denied"

```bash
# Dosya izinlerini düzelt
chmod -R 755 plugins/
chmod -R 755 web-ui/

# Ownership'i düzelt
sudo chown -R $USER:$USER plugins/
sudo chown -R $USER:$USER web-ui/
```

## Destek Alma

### 1. Log Toplama

```bash
# Backend logs
go run main.go 2>&1 | tee panel.log

# Frontend logs
cd web-ui
pnpm dev 2>&1 | tee build.log
```

### 2. Sistem Bilgisi

```bash
# Go version
go version

# Node version
node --version

# pnpm version
pnpm --version

# OS info
uname -a
```

### 3. Plugin Bilgisi

```bash
# Plugin listesi
panel plugin list --json > plugins.json

# Plugin metadata
cat plugins/importer/plugin.yaml
```

### 4. Issue Oluşturma

GitHub Issues: https://github.com/ferdiunal/panel.go/issues

**Issue Template:**
```markdown
## Sorun
[Sorunun kısa açıklaması]

## Adımlar
1. [Adım 1]
2. [Adım 2]
3. [Adım 3]

## Beklenen Davranış
[Ne olmasını bekliyordunuz?]

## Gerçek Davranış
[Ne oldu?]

## Ortam
- OS: [macOS 14.0]
- Go: [1.25.0]
- Node: [18.0.0]
- Panel.go: [1.0.0]

## Loglar
```
[Log çıktısı]
```

## Ek Bilgi
[Ekran görüntüsü, kod snippet, vb.]
```

## İleri Okuma

- [Plugin Sistemi](./PLUGIN_SYSTEM.md) - Genel bakış
- [CLI Komutları](./PLUGIN_CLI.md) - CLI referansı
- [Plugin Geliştirme](./PLUGIN_DEVELOPMENT.md) - Backend ve frontend
- [Örnekler](./PLUGIN_EXAMPLES.md) - Gerçek dünya örnekleri
