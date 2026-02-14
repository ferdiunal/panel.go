# i18n Helper Fonksiyonları Kullanım Kılavuzu

Bu dokümantasyon, Panel.go'da Laravel'deki `__()` helper'ına benzer şekilde çalışan i18n helper fonksiyonlarının nasıl kullanılacağını açıklar.

## İçindekiler

1. [Genel Bakış](#genel-bakış)
2. [Middleware Yapılandırması](#middleware-yapılandırması)
3. [Helper Fonksiyonları](#helper-fonksiyonları)
4. [Fields'larda Kullanım](#fieldslarda-kullanım)
5. [Kullanım Örnekleri](#kullanım-örnekleri)

---

## Genel Bakış

Panel.go, Laravel'deki `__()` helper'ına benzer şekilde çalışan i18n helper fonksiyonları sağlar. Bu fonksiyonlar, fields, labels, placeholders ve diğer UI elementlerinde çoklu dil desteği sağlamak için kullanılır.

### Temel Özellikler

- ✅ Laravel'deki `__()` helper'ına benzer API
- ✅ Template değişkenleri desteği
- ✅ Çoğul form desteği (TransChoice)
- ✅ Fallback değer desteği
- ✅ Çeviri varlık kontrolü
- ✅ Mevcut dil bilgisi

---

## Middleware Yapılandırması

Panel.go'da i18n desteği, Fiber i18n middleware'i kullanılarak sağlanır. Middleware, dil seçimini otomatik olarak yönetir ve çeviri fonksiyonlarını kullanıma hazır hale getirir.

### Yapılandırma

```go
import "golang.org/x/text/language"

config := panel.Config{
    // ... diğer yapılandırmalar
    I18n: panel.I18nConfig{
        Enabled:          true,
        RootPath:         "./locales",
        AcceptLanguages:  []language.Tag{language.Turkish, language.English},
        DefaultLanguage:  language.Turkish,
        FormatBundleFile: "yaml",
    },
}
```

### Parametreler

| Parametre | Tip | Varsayılan | Açıklama |
|-----------|-----|------------|----------|
| `Enabled` | bool | false | i18n'i etkinleştirir |
| `RootPath` | string | "./locales" | Dil dosyalarının bulunduğu dizin |
| `AcceptLanguages` | []language.Tag | [tr, en] | Desteklenen diller listesi |
| `DefaultLanguage` | language.Tag | Turkish | Varsayılan dil (fallback) |
| `FormatBundleFile` | string | "yaml" | Dil dosyası formatı (yaml, json, toml) |

### Dil Seçimi

Dil, şu sırayla belirlenir:

1. **Query Parametresi**: `?lang=tr`
2. **Accept-Language Header**: `Accept-Language: tr-TR,tr;q=0.9,en;q=0.8`
3. **DefaultLanguage**: Fallback dil

### Dil Değiştirme

**Query parametresi ile:**
```bash
curl http://localhost:8080/api/resource/users?lang=en
```

**Header ile:**
```bash
curl -H "Accept-Language: en-US,en;q=0.9" http://localhost:8080/api/resource/users
```

---

## Helper Fonksiyonları

### 1. Trans() - Basit Çeviri

Laravel'deki `__()` helper'ına benzer şekilde çalışır.

```go
import "github.com/ferdiunal/panel.go/pkg/i18n"

// Basit kullanım
message := i18n.Trans(c, "welcome")
// Çıktı: "Hoş geldiniz"

// Template değişkenleri ile
message := i18n.Trans(c, "welcomeWithName", map[string]interface{}{
    "Name": "Ahmet",
})
// Çıktı: "Hoş geldiniz, Ahmet"
```

### 2. TransChoice() - Çoğul Çeviri

Laravel'deki `trans_choice()` helper'ına benzer şekilde çalışır.

```go
// Tekil
message := i18n.TransChoice(c, "items", 1)
// Çıktı: "1 öğe"

// Çoğul
message := i18n.TransChoice(c, "items", 5)
// Çıktı: "5 öğe"

// Template değişkenleri ile
message := i18n.TransChoice(c, "itemsWithName", 3, map[string]interface{}{
    "Name": "Ürün",
})
// Çıktı: "3 Ürün"
```

### 3. GetLocale() - Mevcut Dil

Laravel'deki `app()->getLocale()` metoduna benzer.

```go
lang := i18n.GetLocale(c)
// Çıktı: "tr" veya "en"
```

### 4. HasTranslation() - Çeviri Kontrolü

Laravel'deki `Lang::has()` metoduna benzer.

```go
if i18n.HasTranslation(c, "welcome") {
    message := i18n.Trans(c, "welcome")
}
```

### 5. TransWithFallback() - Fallback ile Çeviri

Çeviri yoksa varsayılan değer döndürür.

```go
message := i18n.TransWithFallback(c, "unknown.key", "Varsayılan Mesaj")
// Çıktı: "Varsayılan Mesaj" (çeviri yoksa)
```

---

## Fields'larda Kullanım

### Resource Title ve Group i18n

Resource'ların başlık ve grup isimlerini i18n ile yönetmek için `SetTitleFunc` ve `SetGroupFunc` kullanın:

```go
package resource

import (
	"github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserResource struct {
	resource.OptimizedBase
}

func NewUserResource() *UserResource {
	r := &UserResource{}
	r.SetModel(&User{})
	r.SetSlug("users")

	// i18n desteği için SetTitleFunc ve SetGroupFunc kullanın
	r.SetTitleFunc(func(c *fiber.Ctx) string {
		return i18n.Trans(c, "resources.users.title")
	})
	r.SetGroupFunc(func(c *fiber.Ctx) string {
		return i18n.Trans(c, "resources.groups.user_management")
	})

	// Alternatif olarak statik değerler için:
	// r.SetTitle("Users")
	// r.SetGroup("User Management")

	r.SetIcon("users")
	r.SetNavigationOrder(1)
	r.SetVisible(true)
	return r
}

func (r *UserResource) Repository(db *gorm.DB) data.DataProvider {
	return NewUserRepository(db)
}
```

**Dil Dosyası:**
```yaml
# locales/tr.yaml
resources:
  groups:
    user_management:
      other: "Kullanıcı Yönetimi"
  users:
    title:
      other: "Kullanıcılar"
```

### FieldResolver'da Güvenli i18n Kullanımı

FieldResolver'da `ctx.Ctx` nil olabilir (resource initialization sırasında). Güvenli i18n kullanımı için `trans()` helper metodu kullanın:

```go
package resource

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/i18n"
)

type UserFieldResolver struct{}

// trans, güvenli i18n çevirisi yapar. ctx veya ctx.Ctx nil ise fallback değeri döner.
func (r *UserFieldResolver) trans(ctx *context.Context, key string, fallback string) string {
	if ctx == nil || ctx.Ctx == nil {
		return fallback
	}
	return i18n.Trans(ctx.Ctx, key)
}

func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID().ReadOnly().OnlyOnDetail(),

		fields.Text("name").
			Label(r.trans(ctx, "fields.name", "Name")).
			Placeholder(r.trans(ctx, "fields.name_placeholder", "Enter name")).
			Required().
			Searchable(),

		fields.Email("email").
			Label(r.trans(ctx, "fields.email", "Email")).
			Placeholder(r.trans(ctx, "fields.email_placeholder", "Enter email")).
			Required().
			Searchable(),

		fields.Select("role").
			Label(r.trans(ctx, "fields.role", "Role")).
			Placeholder(r.trans(ctx, "fields.role_placeholder", "Select role")).
			Options(map[string]string{
				"admin":  r.trans(ctx, "roles.admin", "Admin"),
				"editor": r.trans(ctx, "roles.editor", "Editor"),
				"viewer": r.trans(ctx, "roles.viewer", "Viewer"),
			}),
	}
}
```

**Önemli Notlar:**
- ✅ `ctx.Ctx` nil kontrolü yapın (resource initialization sırasında nil olabilir)
- ✅ Fallback değerleri sağlayın (i18n dosyası yoksa veya context nil ise)
- ✅ Select Options'larında da i18n kullanın
- ✅ trans() helper metodunu struct'a ekleyin (kod tekrarını önler)

---

## Fields'larda Kullanım (Devamı)

### Resource Tanımında

```go
package resource

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/gofiber/fiber/v2"
)

type UserResource struct {
    resource.BaseResource
}

func (r *UserResource) Fields(c *fiber.Ctx) []fields.Field {
    return []fields.Field{
        // Basit çeviri
        fields.Text("name").
            Label(i18n.Trans(c, "fields.name")).
            Placeholder(i18n.Trans(c, "fields.name_placeholder")),

        // Template değişkenleri ile
        fields.Email("email").
            Label(i18n.Trans(c, "fields.email")).
            Help(i18n.Trans(c, "fields.email_help", map[string]interface{}{
                "Domain": "example.com",
            })),

        // Fallback ile
        fields.Text("phone").
            Label(i18n.TransWithFallback(c, "fields.phone", "Phone Number")),

        // Select seçenekleri
        fields.Select("role").
            Label(i18n.Trans(c, "fields.role")).
            Options(map[string]string{
                "admin":  i18n.Trans(c, "roles.admin"),
                "editor": i18n.Trans(c, "roles.editor"),
                "viewer": i18n.Trans(c, "roles.viewer"),
            }),
    }
}
```

### Dil Dosyası Yapısı

**locales/tr.yaml:**
```yaml
# Fields
fields:
  name:
    other: "Ad"
  name_placeholder:
    other: "Adınızı girin"
  email:
    other: "E-posta"
  email_help:
    other: "{{.Domain}} uzantılı e-posta adresi kullanın"
  phone:
    other: "Telefon"
  role:
    other: "Rol"

# Roles
roles:
  admin:
    other: "Yönetici"
  editor:
    other: "Editör"
  viewer:
    other: "Görüntüleyici"
```

**locales/en.yaml:**
```yaml
# Fields
fields:
  name:
    other: "Name"
  name_placeholder:
    other: "Enter your name"
  email:
    other: "Email"
  email_help:
    other: "Use email address with {{.Domain}} extension"
  phone:
    other: "Phone"
  role:
    other: "Role"

# Roles
roles:
  admin:
    other: "Administrator"
  editor:
    other: "Editor"
  viewer:
    other: "Viewer"
```

---

## Kullanım Örnekleri

### Örnek 1: Basit Resource

```go
package resource

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/gofiber/fiber/v2"
)

type ProductResource struct {
    resource.BaseResource
}

func (r *ProductResource) Fields(c *fiber.Ctx) []fields.Field {
    return []fields.Field{
        fields.Text("name").
            Label(i18n.Trans(c, "product.name")).
            Placeholder(i18n.Trans(c, "product.name_placeholder")).
            Required(),

        fields.Textarea("description").
            Label(i18n.Trans(c, "product.description")).
            Help(i18n.Trans(c, "product.description_help")),

        fields.Number("price").
            Label(i18n.Trans(c, "product.price")).
            Min(0).
            Step(0.01),

        fields.Number("stock").
            Label(i18n.Trans(c, "product.stock")).
            Help(i18n.TransChoice(c, "product.stock_help", 10)),

        fields.Select("category").
            Label(i18n.Trans(c, "product.category")).
            Options(r.getCategoryOptions(c)),
    }
}

func (r *ProductResource) getCategoryOptions(c *fiber.Ctx) map[string]string {
    return map[string]string{
        "electronics": i18n.Trans(c, "categories.electronics"),
        "clothing":    i18n.Trans(c, "categories.clothing"),
        "food":        i18n.Trans(c, "categories.food"),
        "books":       i18n.Trans(c, "categories.books"),
    }
}
```

### Örnek 2: Validation Mesajları

```go
package resource

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/gofiber/fiber/v2"
)

type UserResource struct {
    resource.BaseResource
}

func (r *UserResource) Fields(c *fiber.Ctx) []fields.Field {
    return []fields.Field{
        fields.Text("username").
            Label(i18n.Trans(c, "user.username")).
            Required().
            MinLength(3).
            MaxLength(20).
            Rules([]string{
                i18n.Trans(c, "validation.required"),
                i18n.Trans(c, "validation.min_length", map[string]interface{}{
                    "Min": 3,
                }),
                i18n.Trans(c, "validation.max_length", map[string]interface{}{
                    "Max": 20,
                }),
            }),

        fields.Email("email").
            Label(i18n.Trans(c, "user.email")).
            Required().
            Rules([]string{
                i18n.Trans(c, "validation.required"),
                i18n.Trans(c, "validation.email"),
            }),

        fields.Password("password").
            Label(i18n.Trans(c, "user.password")).
            Required().
            MinLength(8).
            Rules([]string{
                i18n.Trans(c, "validation.required"),
                i18n.Trans(c, "validation.min_length", map[string]interface{}{
                    "Min": 8,
                }),
            }),
    }
}
```

### Örnek 3: Dinamik Seçenekler

```go
package resource

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/gofiber/fiber/v2"
)

type OrderResource struct {
    resource.BaseResource
}

func (r *OrderResource) Fields(c *fiber.Ctx) []fields.Field {
    return []fields.Field{
        fields.Select("status").
            Label(i18n.Trans(c, "order.status")).
            Options(r.getStatusOptions(c)).
            Help(i18n.Trans(c, "order.status_help")),

        fields.Select("payment_method").
            Label(i18n.Trans(c, "order.payment_method")).
            Options(r.getPaymentMethodOptions(c)),

        fields.Number("total").
            Label(i18n.Trans(c, "order.total")).
            Readonly().
            Help(i18n.Trans(c, "order.total_help", map[string]interface{}{
                "Currency": "TL",
            })),
    }
}

func (r *OrderResource) getStatusOptions(c *fiber.Ctx) map[string]string {
    return map[string]string{
        "pending":    i18n.Trans(c, "order.status.pending"),
        "processing": i18n.Trans(c, "order.status.processing"),
        "shipped":    i18n.Trans(c, "order.status.shipped"),
        "delivered":  i18n.Trans(c, "order.status.delivered"),
        "cancelled":  i18n.Trans(c, "order.status.cancelled"),
    }
}

func (r *OrderResource) getPaymentMethodOptions(c *fiber.Ctx) map[string]string {
    return map[string]string{
        "credit_card": i18n.Trans(c, "payment.credit_card"),
        "debit_card":  i18n.Trans(c, "payment.debit_card"),
        "paypal":      i18n.Trans(c, "payment.paypal"),
        "bank":        i18n.Trans(c, "payment.bank_transfer"),
    }
}
```

### Örnek 4: Action'larda Kullanım

```go
package action

import (
    "github.com/ferdiunal/panel.go/pkg/action"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/gofiber/fiber/v2"
)

type PublishAction struct {
    action.BaseAction
}

func (a *PublishAction) Name(c *fiber.Ctx) string {
    return i18n.Trans(c, "actions.publish")
}

func (a *PublishAction) ConfirmText(c *fiber.Ctx) string {
    return i18n.Trans(c, "actions.publish_confirm")
}

func (a *PublishAction) ConfirmButtonText(c *fiber.Ctx) string {
    return i18n.Trans(c, "actions.publish_button")
}

func (a *PublishAction) CancelButtonText(c *fiber.Ctx) string {
    return i18n.Trans(c, "actions.cancel")
}

func (a *PublishAction) SuccessMessage(c *fiber.Ctx, count int) string {
    return i18n.TransChoice(c, "actions.publish_success", count, map[string]interface{}{
        "Count": count,
    })
}
```

### Örnek 5: Page'lerde Kullanım

```go
package page

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/i18n"
    "github.com/ferdiunal/panel.go/pkg/page"
    "github.com/gofiber/fiber/v2"
)

type SettingsPage struct {
    page.BasePage
}

func (p *SettingsPage) Title(c *fiber.Ctx) string {
    return i18n.Trans(c, "pages.settings")
}

func (p *SettingsPage) Description(c *fiber.Ctx) string {
    return i18n.Trans(c, "pages.settings_description")
}

func (p *SettingsPage) Elements(c *fiber.Ctx) []fields.Element {
    return []fields.Element{
        fields.Text("site_name").
            Label(i18n.Trans(c, "settings.site_name")).
            Placeholder(i18n.Trans(c, "settings.site_name_placeholder")).
            Default("Panel.go"),

        fields.Switch("register_enable").
            Label(i18n.Trans(c, "settings.register_enable")).
            Help(i18n.Trans(c, "settings.register_enable_help")).
            Default(true),

        fields.Switch("forgot_password_enable").
            Label(i18n.Trans(c, "settings.forgot_password_enable")).
            Help(i18n.Trans(c, "settings.forgot_password_enable_help")).
            Default(false),
    }
}
```

---

## Best Practices

### 1. Çeviri Anahtarlarını Organize Edin

```yaml
# İyi ✅
fields:
  user:
    name:
      other: "Ad"
    email:
      other: "E-posta"

# Kötü ❌
user_name:
  other: "Ad"
user_email:
  other: "E-posta"
```

### 2. Template Değişkenlerini Kullanın

```go
// İyi ✅
i18n.Trans(c, "welcome_message", map[string]interface{}{
    "Name": user.Name,
    "Date": time.Now().Format("2006-01-02"),
})

// Kötü ❌
"Hoş geldiniz " + user.Name + ", tarih: " + time.Now().Format("2006-01-02")
```

### 3. Fallback Değerleri Kullanın

```go
// İyi ✅
i18n.TransWithFallback(c, "custom.field", "Custom Field")

// Kötü ❌
i18n.Trans(c, "custom.field") // Çeviri yoksa hata
```

### 4. Çoğul Formları Kullanın

```go
// İyi ✅
i18n.TransChoice(c, "items_count", count)

// Kötü ❌
if count == 1 {
    return "1 öğe"
} else {
    return fmt.Sprintf("%d öğe", count)
}
```

### 5. Dil Kontrolü Yapın

```go
// İyi ✅
if i18n.HasTranslation(c, "custom.message") {
    message := i18n.Trans(c, "custom.message")
} else {
    message := "Default message"
}

// Kötü ❌
message := i18n.Trans(c, "custom.message") // Panic olabilir
```

---

## Sorun Giderme

### Problem: Çeviriler gösterilmiyor

**Çözüm:**
1. Dil dosyalarının doğru dizinde olduğunu kontrol edin (`locales/tr.yaml`)
2. YAML formatının doğru olduğunu kontrol edin
3. i18n middleware'inin etkin olduğunu kontrol edin
4. Çeviri anahtarının doğru olduğunu kontrol edin

### Problem: Template değişkenleri çalışmıyor

**Çözüm:**
1. Template değişkenlerinin `{{.VariableName}}` formatında olduğunu kontrol edin
2. `map[string]interface{}` kullandığınızdan emin olun
3. Değişken isimlerinin büyük harfle başladığından emin olun

### Problem: Çoğul formlar çalışmıyor

**Çözüm:**
1. `TransChoice()` fonksiyonunu kullandığınızdan emin olun
2. Dil dosyasında `one`, `other` formlarını tanımladığınızdan emin olun
3. `PluralCount` parametresinin doğru olduğunu kontrol edin

---

## Kaynaklar

- [Circuit Breaker & i18n Dokümantasyonu](CIRCUIT_BREAKER_I18N.md)
- [go-i18n Kütüphanesi](https://github.com/nicksnyder/go-i18n)
- [Fiber i18n Middleware](https://docs.gofiber.io/contrib/fiberi18n/)
