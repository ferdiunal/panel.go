# RTL (Right-to-Left) Desteği Kullanım Kılavuzu

Bu dokümantasyon, Panel.go'da RTL (Right-to-Left) desteğinin nasıl kullanılacağını açıklar.

## İçindekiler

1. [Genel Bakış](#genel-bakış)
2. [Backend RTL Desteği](#backend-rtl-desteği)
3. [Frontend RTL Desteği](#frontend-rtl-desteği)
4. [Kullanım Örnekleri](#kullanım-örnekleri)

---

## Genel Bakış

Panel.go, Arapça, İbranice, Farsça ve Urduca gibi sağdan sola yazılan diller için tam RTL desteği sağlar.

### Desteklenen RTL Dilleri

| Dil | Kod | Açıklama |
|-----|-----|----------|
| Arapça | `ar` | Arabic |
| İbranice | `he` | Hebrew |
| Farsça | `fa` | Persian |
| Urduca | `ur` | Urdu |

---

## Backend RTL Desteği

### RTL Helper Fonksiyonları

**Dosya:** `pkg/rtl/rtl.go`

```go
import "github.com/ferdiunal/panel.go/pkg/rtl"

// Dil RTL mi kontrol et
if rtl.IsRTL(language.Arabic) {
    // RTL layout kullan
}

// Dil kodu RTL mi kontrol et
if rtl.IsRTLString("ar") {
    // RTL layout kullan
}

// Text direction al
dir := rtl.GetDirection(language.Arabic)
// dir = "rtl"

// Dil kodundan direction al
dir := rtl.GetDirectionString("ar")
// dir = "rtl"

// Context'ten direction al
dir := rtl.GetDirectionFromContext(c)
```

### RTL Middleware

RTL middleware, otomatik olarak `X-Text-Direction` header'ını ekler:

```go
import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/ferdiunal/panel.go/pkg/rtl"
)

func main() {
    config := panel.Config{
        // ... diğer yapılandırmalar
    }

    p := panel.New(config)

    // RTL middleware'ini ekle
    p.Fiber.Use(rtl.Middleware())

    p.Start()
}
```

### i18n ile Entegrasyon

RTL desteği, i18n helper fonksiyonları ile entegre çalışır:

```go
import (
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/rtl"
)

func MyHandler(c *fiber.Ctx) error {
    // Mevcut dili al
    lang := i18n.GetLocale(c)

    // Direction'ı al
    dir := rtl.GetDirectionString(lang)

    // Çeviri yap
    message := i18n.Trans(c, "welcome")

    return c.JSON(fiber.Map{
        "message":   message,
        "direction": dir,
        "lang":      lang,
    })
}
```

---

## Frontend RTL Desteği

### React Hook

Frontend'de RTL desteği için bir React hook oluşturun:

**Dosya:** `web/src/hooks/use-rtl.ts`

```typescript
import { useEffect } from 'react'

const RTL_LANGUAGES = ['ar', 'he', 'fa', 'ur']

export function useRTL(lang: string) {
  const isRTL = RTL_LANGUAGES.includes(lang)
  const direction = isRTL ? 'rtl' : 'ltr'

  useEffect(() => {
    // HTML dir attribute'unu güncelle
    document.documentElement.dir = direction
    document.documentElement.lang = lang

    // Body'ye RTL class ekle
    if (isRTL) {
      document.body.classList.add('rtl')
    } else {
      document.body.classList.remove('rtl')
    }
  }, [lang, direction, isRTL])

  return { isRTL, direction }
}
```

### Tailwind CSS RTL Desteği

Tailwind CSS'de RTL desteği için `rtl:` prefix'ini kullanın:

```tsx
// LTR'de margin-left, RTL'de margin-right
<div className="ml-4 rtl:mr-4 rtl:ml-0">
  Content
</div>

// LTR'de text-left, RTL'de text-right
<div className="text-left rtl:text-right">
  Text
</div>

// LTR'de float-left, RTL'de float-right
<div className="float-left rtl:float-right">
  Sidebar
</div>
```

### Logical Properties

Modern CSS logical properties kullanın:

```css
/* Fiziksel properties yerine */
.element {
  margin-left: 1rem;  /* ❌ Kötü */
  margin-right: 1rem; /* ❌ Kötü */
}

/* Logical properties kullanın */
.element {
  margin-inline-start: 1rem;  /* ✅ İyi */
  margin-inline-end: 1rem;    /* ✅ İyi */
}
```

Tailwind CSS logical properties:

```tsx
// margin-inline-start (LTR'de left, RTL'de right)
<div className="ms-4">Content</div>

// margin-inline-end (LTR'de right, RTL'de left)
<div className="me-4">Content</div>

// padding-inline-start
<div className="ps-4">Content</div>

// padding-inline-end
<div className="pe-4">Content</div>
```

---

## Kullanım Örnekleri

### Örnek 1: RTL Dil Desteği ile Resource

```go
package resource

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/ferdiunal/panel.go/pkg/rtl"
    "github.com/gofiber/fiber/v2"
)

type ProductResource struct {
    resource.OptimizedBase
}

func NewProductResource() *ProductResource {
    r := &ProductResource{}

    r.SetModel(&Product{})
    r.SetSlug("products")

    // i18n destekli başlık
    r.SetTitleFunc(func(c *fiber.Ctx) string {
        return i18n.Trans(c, "resources.products.title")
    })

    r.SetIcon("package")

    // i18n destekli grup
    r.SetGroupFunc(func(c *fiber.Ctx) string {
        return i18n.Trans(c, "resources.groups.content")
    })

    return r
}

func (r *ProductResource) Fields(c *fiber.Ctx) []fields.Field {
    // RTL direction'ı al
    dir := rtl.GetDirectionFromContext(c)

    return []fields.Field{
        fields.Text("name").
            Label(i18n.Trans(c, "product.name")).
            Placeholder(i18n.Trans(c, "product.name_placeholder")).
            // RTL için text direction ekle
            Props(map[string]interface{}{
                "dir": dir,
            }),

        fields.Textarea("description").
            Label(i18n.Trans(c, "product.description")).
            Props(map[string]interface{}{
                "dir": dir,
            }),
    }
}
```

### Örnek 2: Frontend RTL Hook Kullanımı

```tsx
import { useRTL } from '@/hooks/use-rtl'
import { useLanguage } from '@/hooks/use-language'

export function App() {
  const { lang } = useLanguage()
  const { isRTL, direction } = useRTL(lang)

  return (
    <div className="app" dir={direction}>
      {/* LTR'de left, RTL'de right */}
      <aside className="sidebar ms-0 me-auto rtl:ms-auto rtl:me-0">
        Sidebar
      </aside>

      {/* Ana içerik */}
      <main className="content">
        <h1 className="text-start rtl:text-end">
          {isRTL ? 'مرحبا' : 'Welcome'}
        </h1>
      </main>
    </div>
  )
}
```

### Örnek 3: API Response ile RTL Bilgisi

```go
func (p *Panel) handleInit(c *context.Context) error {
    // Mevcut dili al
    lang := i18n.GetLocale(c.Ctx)

    // RTL mi kontrol et
    isRTL := rtl.IsRTLString(lang)
    dir := rtl.GetDirectionString(lang)

    return c.JSON(fiber.Map{
        "features": fiber.Map{
            "register":        p.Config.Features.Register,
            "forgot_password": p.Config.Features.ForgotPassword,
        },
        "i18n": fiber.Map{
            "lang":      lang,
            "direction": dir,
            "isRTL":     isRTL,
        },
        "version":  "1.0.0",
        "settings": p.Config.SettingsValues.Values,
    })
}
```

### Örnek 4: Arapça Dil Desteği

**Dil Dosyası (`locales/ar/messages.yaml`):**

```yaml
# Genel Mesajlar
welcome:
  other: "مرحبا"

welcomeWithName:
  other: "مرحبا، {{.Name}}"

# Resource Grupları
resources:
  groups:
    system:
      other: "النظام"
    content:
      other: "المحتوى"

  # User Resource
  users:
    title:
      other: "المستخدمون"
    fields:
      name:
        other: "الاسم"
      email:
        other: "البريد الإلكتروني"
      password:
        other: "كلمة المرور"
```

**Kullanım:**

```bash
# Arapça ile API çağrısı
curl http://localhost:8080/api/resource/users?lang=ar

# Response
{
  "data": [...],
  "meta": {
    "lang": "ar",
    "direction": "rtl",
    "isRTL": true
  }
}
```

---

## Best Practices

### 1. Logical Properties Kullanın

```css
/* ❌ Kötü - Fiziksel properties */
.element {
  margin-left: 1rem;
  padding-right: 2rem;
  border-left: 1px solid;
}

/* ✅ İyi - Logical properties */
.element {
  margin-inline-start: 1rem;
  padding-inline-end: 2rem;
  border-inline-start: 1px solid;
}
```

### 2. Tailwind RTL Utilities

```tsx
/* ❌ Kötü - Sabit direction */
<div className="ml-4 text-left">Content</div>

/* ✅ İyi - RTL destekli */
<div className="ms-4 text-start rtl:text-end">Content</div>
```

### 3. Icon Flipping

Bazı iconlar RTL'de ters çevrilmelidir:

```tsx
/* Yön gösteren iconlar RTL'de ters çevrilmeli */
<ChevronRight className="rtl:rotate-180" />
<ArrowLeft className="rtl:rotate-180" />

/* Nötr iconlar ters çevrilmemeli */
<User className="" />
<Settings className="" />
```

### 4. Text Alignment

```tsx
/* ❌ Kötü */
<p className="text-left">Text</p>

/* ✅ İyi */
<p className="text-start">Text</p>
```

### 5. Flexbox Direction

```tsx
/* ❌ Kötü */
<div className="flex flex-row">
  <div>Left</div>
  <div>Right</div>
</div>

/* ✅ İyi - Otomatik RTL desteği */
<div className="flex">
  <div>Start</div>
  <div>End</div>
</div>
```

---

## Sorun Giderme

### Problem: RTL düzgün çalışmıyor

**Çözüm:**
1. HTML `dir` attribute'unun doğru ayarlandığını kontrol edin
2. Tailwind CSS'de `rtl:` prefix'inin çalıştığını kontrol edin
3. Logical properties kullandığınızdan emin olun

### Problem: Iconlar ters çevrilmiyor

**Çözüm:**
1. Yön gösteren iconlara `rtl:rotate-180` class'ı ekleyin
2. Nötr iconları olduğu gibi bırakın

### Problem: Text alignment yanlış

**Çözüm:**
1. `text-left/right` yerine `text-start/end` kullanın
2. `rtl:text-end` gibi RTL-specific class'lar ekleyin

---

## Kaynaklar

- [Tailwind CSS RTL Support](https://tailwindcss.com/docs/hover-focus-and-other-states#rtl-support)
- [CSS Logical Properties](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Logical_Properties)
- [shadcn/ui RTL Documentation](https://ui.shadcn.com/docs/rtl)
