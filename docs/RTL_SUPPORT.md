# RTL ve Dark Tema Desteği

Panel.go, backend-driven RTL (Right-to-Left) ve dark tema desteği sağlar. Kullanıcı sadece backend config yapar, frontend otomatik olarak handle eder.

## İçindekiler

- [Genel Bakış](#genel-bakış)
- [Backend Yapılandırma](#backend-yapılandırma)
- [HTML Injection Sistemi](#html-injection-sistemi)
- [Frontend Entegrasyon](#frontend-entegrasyon)
- [Kullanım Örnekleri](#kullanım-örnekleri)

## Genel Bakış

### Backend-Driven Yaklaşım

Panel.go, RTL ve dark tema desteğini **backend-driven** olarak sağlar:

1. ✅ **Kullanıcı sadece Go config yapar** - Frontend kod yazmaya gerek yok
2. ✅ **Otomatik HTML injection** - Backend runtime'da HTML'i modify eder
3. ✅ **Initial render'da aktif** - JavaScript gerektirmez, SEO friendly
4. ✅ **API desteği** - `/api/init` endpoint'i RTL ve theme bilgisi döndürür

### Desteklenen RTL Dilleri

| Dil | Kod | Direction |
|-----|-----|-----------|
| Arapça | `ar` | rtl |
| İbranice | `he` | rtl |
| Farsça | `fa` | rtl |
| Urduca | `ur` | rtl |

### Desteklenen Temalar

- `light` - Açık tema (varsayılan)
- `dark` - Koyu tema

## Backend Yapılandırma

### 1. i18n Config (RTL için)

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "golang.org/x/text/language"
)

func main() {
    config := panel.Config{
        I18n: panel.I18nConfig{
            Enabled:          true,
            RootPath:         "./locales",
            AcceptLanguages:  []language.Tag{
                language.Turkish,
                language.English,
                language.Arabic,  // RTL otomatik aktif
                language.Hebrew,  // RTL otomatik aktif
            },
            DefaultLanguage:  language.Turkish,
            FormatBundleFile: "yaml",
        },
    }

    p := panel.New(config)
    p.Start()
}
```

### 2. Tema Config

Tema bilgisi cookie veya query parameter'dan alınır:

```bash
# Cookie ile
curl -H "Cookie: theme=dark" http://localhost:8080

# Query parameter ile
curl http://localhost:8080?theme=dark
```

## HTML Injection Sistemi

### Nasıl Çalışır?

1. **panel.web'de placeholder'lar** (`index.html`):
   ```html
   <html lang="{{PANEL_LANG}}" dir="{{PANEL_DIR}}" data-theme="{{PANEL_THEME}}">
     <head>
       <title>{{PANEL_TITLE}}</title>
     </head>
   </html>
   ```

2. **Backend runtime'da replace** (`pkg/panel/html.go`):
   ```go
   html = strings.ReplaceAll(html, "{{PANEL_LANG}}", "ar")
   html = strings.ReplaceAll(html, "{{PANEL_DIR}}", "rtl")
   html = strings.ReplaceAll(html, "{{PANEL_THEME}}", "dark")
   html = strings.ReplaceAll(html, "{{PANEL_TITLE}}", "Panel.go")
   ```

3. **Client'a serve edilen HTML**:
   ```html
   <html lang="ar" dir="rtl" data-theme="dark">
     <head>
       <title>Panel.go</title>
     </head>
   </html>
   ```

### Placeholder'lar

| Placeholder | Açıklama | Örnek Değer |
|-------------|----------|-------------|
| `{{PANEL_LANG}}` | Dil kodu | `ar`, `en`, `tr` |
| `{{PANEL_DIR}}` | Text direction | `rtl`, `ltr` |
| `{{PANEL_THEME}}` | Tema | `dark`, `light` |
| `{{PANEL_TITLE}}` | Site başlığı | `Panel.go` |

### HTML Injection Fonksiyonları

```go
// pkg/panel/html.go

// Injection data oluştur
data := panel.GetHTMLInjectionData(c, config)

// HTML'i inject et
html := panel.InjectHTML(htmlString, data)

// HTML serve et
panel.ServeHTML(c, "assets/ui/index.html", config)
```

## Frontend Entegrasyon

### CSS ile RTL Desteği

Frontend'de `[dir="rtl"]` selector ile RTL stilleri yönetilir:

```css
/* Tailwind CSS */
[dir="rtl"] .rtl\:mr-4 {
  margin-right: 1rem;
}

[dir="rtl"] .rtl\:ml-0 {
  margin-left: 0;
}

/* Custom CSS */
[dir="rtl"] {
  direction: rtl;
}

[dir="rtl"] .sidebar {
  left: auto;
  right: 0;
}
```

### CSS ile Dark Tema Desteği

```css
/* Tailwind CSS */
[data-theme="dark"] {
  color-scheme: dark;
}

[data-theme="dark"] .dark\:bg-gray-900 {
  background-color: #111827;
}

/* Custom CSS */
[data-theme="dark"] body {
  background: #1a1a1a;
  color: #ffffff;
}
```

### JavaScript ile Dinamik Değişiklik (Opsiyonel)

Frontend JavaScript ile tema değiştirme:

```typescript
// Tema değiştir
function setTheme(theme: 'light' | 'dark') {
  document.documentElement.setAttribute('data-theme', theme);
  document.cookie = `theme=${theme}; path=/; max-age=31536000`;
}

// RTL değiştir (dil değişikliği ile)
function setLanguage(lang: string) {
  window.location.href = `/?lang=${lang}`;
}
```

## Kullanım Örnekleri

### Örnek 1: Arapça RTL Desteği

**Backend Config:**
```go
config := panel.Config{
    I18n: panel.I18nConfig{
        Enabled:         true,
        DefaultLanguage: language.Arabic, // RTL otomatik
    },
}
```

**Sonuç:**
```html
<html lang="ar" dir="rtl" data-theme="light">
```

**API Response:**
```json
{
  "i18n": {
    "lang": "ar",
    "direction": "rtl"
  },
  "theme": "light"
}
```

### Örnek 2: Dark Tema

**Request:**
```bash
curl -H "Cookie: theme=dark" http://localhost:8080
```

**Sonuç:**
```html
<html lang="en" dir="ltr" data-theme="dark">
```

### Örnek 3: Arapça + Dark Tema

**Request:**
```bash
curl -H "Cookie: theme=dark" http://localhost:8080?lang=ar
```

**Sonuç:**
```html
<html lang="ar" dir="rtl" data-theme="dark">
```

**API Response:**
```json
{
  "i18n": {
    "lang": "ar",
    "direction": "rtl"
  },
  "theme": "dark"
}
```

### Örnek 4: Resource ile RTL

```go
package resource

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

type ProductResource struct {
    resource.OptimizedBase
}

func NewProductResource() *ProductResource {
    r := &ProductResource{}
    r.SetModel(&Product{})
    r.SetSlug("products")

    // i18n destekli başlık (RTL otomatik)
    r.SetTitleFunc(func(c *fiber.Ctx) string {
        return i18n.Trans(c, "resources.products.title")
    })

    return r
}

func (r *ProductResource) Fields(c *fiber.Ctx) []fields.Field {
    return []fields.Field{
        fields.Text("name").
            Label(i18n.Trans(c, "product.name")).
            Placeholder(i18n.Trans(c, "product.name_placeholder")),

        fields.Textarea("description").
            Label(i18n.Trans(c, "product.description")),
    }
}
```

## API Endpoint

### GET /api/init

Init endpoint'i RTL ve theme bilgisini döndürür:

**Request:**
```bash
curl http://localhost:8080/api/init?lang=ar
```

**Response:**
```json
{
  "features": {
    "register": true,
    "forgot_password": false
  },
  "oauth": {
    "google": false
  },
  "i18n": {
    "lang": "ar",
    "direction": "rtl"
  },
  "theme": "light",
  "version": "1.0.0",
  "settings": {
    "site_name": "Panel.go"
  }
}
```

## Best Practices

### 1. Dil Dosyaları

Arapça dil dosyası örneği (`locales/ar/messages.yaml`):

```yaml
# Genel
welcome:
  other: "مرحبا"

# Resources
resources:
  products:
    title:
      other: "المنتجات"

# Fields
product:
  name:
    other: "الاسم"
  description:
    other: "الوصف"
```

### 2. CSS Logical Properties

Fiziksel properties yerine logical properties kullanın:

```css
/* ❌ Kötü - Fiziksel */
.element {
  margin-left: 1rem;
  padding-right: 2rem;
}

/* ✅ İyi - Logical */
.element {
  margin-inline-start: 1rem;
  padding-inline-end: 2rem;
}
```

### 3. Tailwind RTL Utilities

```tsx
/* ❌ Kötü */
<div className="ml-4 text-left">Content</div>

/* ✅ İyi */
<div className="ms-4 text-start">Content</div>
```

### 4. Icon Flipping

```tsx
/* Yön gösteren iconlar RTL'de ters çevrilmeli */
<ChevronRight className="rtl:rotate-180" />

/* Nötr iconlar ters çevrilmemeli */
<User className="" />
```

## Avantajlar

### Backend-Driven Yaklaşım

✅ **Kullanıcı Dostu**: Sadece Go config, frontend kod yok
✅ **SEO Friendly**: HTML'de doğru `lang`, `dir` attribute'ları
✅ **Initial Render**: JavaScript gerektirmez, ilk yüklemede aktif
✅ **Performanslı**: Sadece string replace (regex yok)
✅ **Type-Safe**: Go type system ile güvenli

### Geleneksel Yaklaşım (❌)

❌ **Frontend Kod Gerekli**: React hook, useEffect, vb.
❌ **JavaScript Bağımlı**: Initial render'da RTL yok
❌ **SEO Sorunları**: HTML attribute'ları eksik
❌ **Karmaşık**: Kullanıcı frontend kod yazmalı

## Troubleshooting

### Problem: RTL çalışmıyor

**Çözüm:**
1. HTML'de `dir` attribute'unu kontrol edin:
   ```bash
   curl http://localhost:8080 | grep 'dir="rtl"'
   ```

2. i18n config'i kontrol edin:
   ```go
   config.I18n.Enabled = true
   ```

3. Dil dosyalarını kontrol edin:
   ```bash
   ls locales/ar/
   ```

### Problem: Dark tema çalışmıyor

**Çözüm:**
1. Cookie'yi kontrol edin:
   ```bash
   curl -v http://localhost:8080 | grep theme
   ```

2. HTML'de `data-theme` attribute'unu kontrol edin:
   ```bash
   curl http://localhost:8080 | grep data-theme
   ```

### Problem: Placeholder'lar replace edilmiyor

**Çözüm:**
1. panel.web'de placeholder'ların doğru olduğunu kontrol edin
2. Backend'de `ServeHTML` fonksiyonunun çağrıldığını kontrol edin
3. Build alın: `panel plugin build`

## Kaynaklar

- [i18n Helpers](./I18N_HELPERS.md) - Çoklu dil desteği
- [Plugin Sistemi](./PLUGIN_SYSTEM.md) - Plugin geliştirme
- [Tailwind RTL](https://tailwindcss.com/docs/hover-focus-and-other-states#rtl-support)
- [CSS Logical Properties](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Logical_Properties)
