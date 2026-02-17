# Doğrulama (Validation)

Bu doküman, Panel.go'nun create/update isteklerinde çalışan server-side doğrulama akışını açıklar.

Kapsam:
- `go-playground/validator/v10` tabanlı rule doğrulama
- Field bazlı (`errors[field]`) hata response sözleşmesi
- Özelleştirilebilir mesajlar
- i18n entegrasyonu
- Frontend'de hatanın ilgili field altında gösterimi

## Nerede çalışır?

Validasyon aşağıdaki endpoint akışlarında çalışır:
- `POST /api/resource/:resource` (create)
- `PUT/PATCH /api/resource/:resource/:id` (update)

İlgili dosyalar:
- `pkg/handler/request_validation.go`
- `pkg/handler/resource_store_controller.go`
- `pkg/handler/resource_update_controller.go`
- `web/src/components/forms/UniversalResourceForm.tsx`
- `web/src/pages/resource/index.tsx`

## Desteklenen kurallar

Backend tarafında desteklenen kurallar:
- `required`
- `email`
- `url`
- `min`
- `max`
- `minLength` / `min_length` / `minlength`
- `maxLength` / `max_length` / `maxlength`
- `pattern`
- `unique`
- `exists`

Notlar:
- `pattern` kuralı için özel `panel_regex` validator kullanılır.
- `unique` ve `exists` kuralları GORM üzerinden veritabanı sorgusu ile doğrulanır.
- `read_only` veya `disabled` field'lar validasyondan atlanır.
- Görünür olmayan field'lar (`IsVisible == false`) validasyondan atlanır.

## Field tanımı örneği

### Klasik API (geriye uyumlu)

```go
fieldDefs := []fields.Element{
	fields.Text("Ad Soyad", "full_name").
		Required().
		AddValidationRule(fields.MinLength(3)).
		AddValidationRule(fields.MaxLength(80)).
		WithProps("validation_messages", map[string]interface{}{
			"required":  "validation.required",
			"minLength": "validation.minLength",
		}),

	fields.Email("E-posta", "email").
		Required().
		AddValidationRule(fields.EmailRule()).
		AddValidationRule(fields.Unique("users", "email")),

	fields.Text("Profil URL", "profile_url").
		AddValidationRule(fields.URL()),
}
```

### Rules / CreationRules / UpdateRules API

`Rules()`, `CreationRules()` ve `UpdateRules()` metotları ile create/update bağlamına duyarlı
doğrulama kuralları tanımlanabilir.

- `Rules(...)` — Her iki akışta (create + update) geçerli olan temel kurallar.
- `CreationRules(...)` — Sadece create akışında ek olarak uygulanan kurallar.
- `UpdateRules(...)` — Sadece update akışında ek olarak uygulanan kurallar.

Create akışında çalışan kurallar: `Rules() + CreationRules()` birleşimi.
Update akışında çalışan kurallar: `Rules() + UpdateRules()` birleşimi.

```go
fieldDefs := []fields.Element{
	// Temel kullanım: variadic Rules ile toplu kural ekleme
	fields.Text("Ad Soyad", "full_name").
		Rules(fields.Required(), fields.MinLength(3), fields.MaxLength(80)),

	// Create/Update ayrımı
	fields.Email("E-posta", "email").
		Rules(fields.EmailRule()).
		CreationRules(fields.Required(), fields.Unique("users", "email")).
		UpdateRules(fields.Unique("users", "email")),

	// Sadece create'de zorunlu
	fields.Password("Şifre", "password").
		CreationRules(fields.Required(), fields.MinLength(8)).
		UpdateRules(fields.MinLength(8)),

	// Klasik API ile birlikte kullanılabilir
	fields.Text("Profil URL", "profile_url").
		Rules(fields.URL()),
}
```

## Mesaj öncelik sırası

Bir kural hata verdiğinde mesaj şu öncelikle seçilir:
1. `props.validation_messages[rule]` veya `props.messages[rule]`
2. Rule içindeki custom message
3. i18n key: `validation.<rule>`
4. Rule fallback message
5. `validation.invalid` (son fallback)

Mesaj anahtarı bir i18n key ise (ör. `validation.required`) otomatik lokalize edilir.

## i18n anahtarları

Varsayılan doğrulama anahtarları:
- `validation.required`
- `validation.email`
- `validation.url`
- `validation.min`
- `validation.max`
- `validation.minLength`
- `validation.maxLength`
- `validation.pattern`
- `validation.unique`
- `validation.exists`
- `validation.invalid`

Dil dosyaları:
- `locales/tr.yaml`
- `locales/en.yaml`
- `pkg/panel/locales/tr.yaml`
- `pkg/panel/locales/en.yaml`

## HTTP 422 response sözleşmesi

Doğrulama hatasında response:

```json
{
  "error": "Validation error",
  "code": "VALIDATION_ERROR",
  "errors": {
    "email": ["Please enter a valid email address"],
    "full_name": ["full_name is required"]
  },
  "details": {
    "email": ["Please enter a valid email address"],
    "full_name": ["full_name is required"]
  }
}
```

Not:
- `errors` ve `details` aynı field bazlı payload'ı taşır.
- Frontend ilk mesajı alıp ilgili field'a bağlar.

## Frontend field-level hata gösterimi

`UniversalResourceForm` server response içindeki `errors/details` map'ini parse eder ve:
- `form.setError(field, { type: "server", message })` uygular
- Böylece hata doğrudan ilgili input altında görünür

`resource/index` sayfasında:
- `422` için generic toast bastırılır (inline field hataları önceliklidir)
- `422` dışı hatalarda standart hata toast akışı devam eder

## Create vs Update davranışı

`CreationRules()` ve `UpdateRules()` kullanıldığında:
- Create akışında: `Rules() + CreationRules()` birleşimi uygulanır
- Update akışında: `Rules() + UpdateRules()` birleşimi uygulanır
- `CreationRules`/`UpdateRules` tanımlanmadığında yalnızca `Rules()` (base rules) uygulanır

`required` kuralı için:
- Create: field yoksa/boşsa hata döner
- Update: field request içinde hiç yoksa `required` atlanır (partial update dostu)

Diğer kurallar:
- Field request'te yoksa atlanır
- Field boşsa (empty string/nil/boş dizi) atlanır

## `unique` ve `exists` notları

`unique(table, column)`:
- Tek değer için doğrulama yapar
- Update akışında aynı tablodaki mevcut kaydı primary key'e göre hariç tutar
- Çoklu değer (`[]`) için `unique` kontrolü uygulanmaz

`exists(table, column)`:
- Tek değer ve çoklu değer (`[]`) destekler
- Çoklu değerde tüm benzersiz değerlerin tabloda bulunması gerekir

## Öneriler

- Kritik alanlarda `required + rule` kombinasyonunu birlikte tanımlayın.
- Mesaj override için i18n key kullanın (`validation.*`) ve metinleri locale dosyalarında yönetin.
- `unique/exists` için table/column isimlerini doğru ve güvenli formatta verin.
- Validation davranışını resource testleri ile doğrulayın (`422`, `errors`, `code`).
