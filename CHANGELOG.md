# Changelog

Tüm önemli değişiklikler bu dosyada dökümante edilir.

## [Unreleased]

### 🧩 Grid View İnce Ayarları (HideOnGrid Card-Only + Card Header Padding)

`HideOnGrid` davranışı kart/listing görünümüne odaklanacak şekilde netleştirildi ve grid kart başlığının görselden ayrışması için ek spacing eklendi.

#### Backend

- `HideOnGrid` artık grid kart/listing alanlarını (`headers`) filtreler; row payload (`data`) içindeki alanları düşürmez.
- Grid görünümünde header'da gizlenen ama index context'inde görünür olan alanlar, payload tarafında korunur.
- Uygulanan dosyalar:
  - `pkg/handler/grid_visibility.go` (yeni)
  - `pkg/handler/resource_index_controller.go`
  - `pkg/handler/lens_controller.go`
  - `pkg/handler/resource_index_controller_test.go`
  - `pkg/handler/lens_controller_test.go`

#### Frontend

- Grid kartta görsel sonrası başlık satırına üst padding eklendi (başlık/görsel arasındaki sıkışma giderildi).
- Uygulanan dosya:
  - `web/src/components/views/ResourceGridView.tsx`

#### Dokümantasyon

- `HideOnGrid` semantiği “kart/listing filtreleme” olarak netleştirildi.
- `HideOnGrid` alanlarının payload'da korunmasının beklenen davranış olduğu dokümante edildi.
- Güncellenen dosyalar:
  - `docs/Grid-View.md`
  - `docs/Fields.md`
  - `docs/API-Reference.md`

#### Doğrulama

- ✅ `go test ./pkg/handler -run 'TestHandleResourceIndex_GridViewVisibility|TestHandleResourceIndex_GridViewVisibilityWithoutContextMiddleware|TestHandleLens_GridViewVisibility|TestHandleLens_GridViewVisibilityWithoutContextMiddleware'`
- ✅ `bun test src/lib/resource-grid-card.test.ts` (`web/`)
- ✅ `bun run build` (`web/`)

### 🔐 Internal/External API Servis Ayrımı ve External Plain JSON

`/api` external API yüzeyi olarak konumlandırıldı; internal endpoint'ler `/api/internal/*` altında toplandı.

#### Backend

- `pkg/panel/rest_api_service.go`:
  - Internal REST servisi bağımsız base path + header + key ile çalışır.
  - Feature gate: `Features.RestAPI`.
  - Varsayılan base path: `/api/internal/rest`.
  - Endpoint seti: `GET index`, `GET detail`, `PUT/PATCH update`, `DELETE`.
- `pkg/panel/external_api_service.go`:
  - External API servisi bağımsız base path + header + key ile çalışır.
  - Feature gate: `Features.ExternalAPI`.
  - Varsayılan base path: `/api`.
  - Endpoint seti: `GET index`, `GET show`, `POST store`, `PUT/PATCH update`, `DELETE`.
  - Response flatten akışı eklendi: field resolver payload yerine doğrudan `name => value` JSON döner.
  - `HideOnApi` context'ine sahip alanlar external yanıtta otomatik gizlenir.
  - External servis auth'u external key yanında panel API key ve managed API key doğrulamasını da kabul eder.
- `pkg/panel/app.go`:
  - Eski `/api/*` panel endpoint'leri `/api/internal/*` altına taşındı.
  - Auth/System endpoint'leri: `/api/internal/auth/*`, `/api/internal/init`, `/api/internal/navigation`.
  - OpenAPI spec sadece external endpoint yüzeyini üretir.
- `pkg/core/types.go`, `pkg/core/element.go`, `pkg/fields/base.go`, `pkg/fields/enum.go`:
  - Yeni context: `HIDE_ON_API`.
  - Yeni fluent API: `HideOnApi()`.

#### Config / Env

- `cmd/panel/stubs/main.stub`, `cmd/panel/stubs/main-postgres.stub`, `cmd/panel/stubs/main-mysql.stub`, `cmd/panel/stubs/env.stub`:
  - Feature env anahtarları: `FEATURE_REST_API`, `FEATURE_EXTERNAL_API`
  - Internal REST env anahtarları: `INTERNAL_REST_API_BASE_PATH`, `INTERNAL_REST_API_HEADER`, `INTERNAL_REST_API_KEY`
  - External API env anahtarları: `EXTERNAL_API_BASE_PATH`, `EXTERNAL_API_HEADER`, `EXTERNAL_API_KEY`
  - Varsayılanlar güncellendi: `INTERNAL_REST_API_BASE_PATH=/api/internal/rest`, `EXTERNAL_API_BASE_PATH=/api`

#### Test / Doğrulama

- `pkg/panel/rest_api_service_test.go`: internal servis auth + detail/update/delete + validation senaryoları
- `pkg/panel/external_api_service_test.go`: flatten response, `HideOnApi`, panel API key fallback, validation uyumu
- `pkg/core/types_test.go`, `pkg/core/element_test.go`, `pkg/fields/fields_test.go`: yeni context/fluent API kapsamı
- ✅ `go test ./...`

### ✅ Create/Update Bağlamına Duyarlı Validation API (Rules / CreationRules / UpdateRules)

Field bazlı doğrulama sistemine `Rules()`, `CreationRules()`, `UpdateRules()` fluent API'si eklendi. Artık create ve update akışlarına özel doğrulama kuralları tanımlanabilir.

#### Backend

- `pkg/fields/base.go`:
  - Schema struct'a `CreationValidationRules` ve `UpdateValidationRules` alanları eklendi.
  - Yeni fluent API metotları: `Rules(...)`, `CreationRules(...)`, `UpdateRules(...)` (variadic).
  - Yeni getter metotları: `GetCreationValidationRules()`, `GetUpdateValidationRules()` (base + context rules merge).
- `pkg/core/element.go`:
  - Element interface'e `Rules()`, `CreationRules()`, `UpdateRules()`, `GetCreationValidationRules()`, `GetUpdateValidationRules()` metotları eklendi.
- `pkg/fields/validation.go`:
  - `MergeValidationRules(base, extra)` helper fonksiyonu eklendi.
- `pkg/handler/request_validation.go`:
  - `collectFieldValidationRules` fonksiyonuna `visibilityCtx` parametresi eklendi.
  - `ContextCreate` → `GetCreationValidationRules()`, `ContextUpdate` → `GetUpdateValidationRules()` ile bağlama duyarlı kural seçimi yapılıyor.

#### Kullanım Örneği

```go
// Create/Update ayrımı
fields.Email("E-posta", "email").
    Rules(fields.EmailRule()).
    CreationRules(fields.Required(), fields.Unique("users", "email")).
    UpdateRules(fields.Unique("users", "email"))

// Sadece create'de zorunlu
fields.Password("Şifre", "password").
    CreationRules(fields.Required(), fields.MinLength(8)).
    UpdateRules(fields.MinLength(8))
```

#### Geriye Uyumluluk

- Mevcut `AddValidationRule()`, `Required()`, `Email()`, `Min()`, `Max()` vb. metotlar aynen çalışmaya devam eder.
- `CreationRules`/`UpdateRules` tanımlanmadığında yalnızca base rules uygulanır.

#### Dokümantasyon

- `docs/Validation.md` güncellendi: `Rules / CreationRules / UpdateRules API` bölümü ve örnekler eklendi.

#### Doğrulama

- ✅ `go build ./...`

### ✅ Field-Level Validator Entegrasyonu (go-playground/validator)

Backend ve frontend form akışına field-level hata döndüren server-side validasyon katmanı eklendi.

#### Backend

- `pkg/handler/request_validation.go` eklendi:
  - `go-playground/validator/v10` ile rule tabanlı doğrulama çalıştırılıyor.
  - Desteklenen kurallar: `required`, `email`, `url`, `min`, `max`, `minLength`, `maxLength`, `pattern`, `unique`, `exists`.
  - `pattern` için özel `panel_regex` validator eklendi.
  - `unique` / `exists` kuralları GORM üzerinden DB kontrolü yapıyor.
  - Mesaj önceliği: `props.validation_messages` override > rule message > i18n key (`validation.*`) > fallback.
  - Hata response formatı standardize edildi:
    - HTTP `422 Unprocessable Entity`
    - `code: "VALIDATION_ERROR"`
    - `errors` ve `details` içinde field bazlı mesajlar.
- Store/Update endpoint’lerine validasyon adımı eklendi:
  - `pkg/handler/resource_store_controller.go`
  - `pkg/handler/resource_update_controller.go`

#### Frontend

- `web/src/components/forms/UniversalResourceForm.tsx`:
  - Backend `422` response’undaki `errors/details` parse edilip `react-hook-form` ile `form.setError(field, ...)` uygulanıyor.
  - Böylece hata mesajı direkt ilgili field altında gösteriliyor.
- `web/src/pages/resource/index.tsx`:
  - Validation (`422`) durumunda generic toast bastırıldı, inline field hataları önceliklendirildi.

#### i18n

- Validation mesaj anahtarları eklendi:
  - `locales/tr.yaml`
  - `locales/en.yaml`
  - `pkg/panel/locales/tr.yaml`
  - `pkg/panel/locales/en.yaml`
- Yeni anahtarlar: `validation.required`, `validation.email`, `validation.url`, `validation.min`, `validation.max`, `validation.minLength`, `validation.maxLength`, `validation.pattern`, `validation.unique`, `validation.exists`, `validation.invalid`.

#### Dokümantasyon

- Yeni doküman eklendi: `docs/Validation.md`
- Doküman; backend validator akışı, mesaj önceliği, i18n, `422` response sözleşmesi ve frontend field-level hata gösterimini kapsar.

#### Test / Doğrulama

- Handler testlerine validation senaryoları eklendi:
  - `pkg/handler/resource_store_controller_test.go`
  - `pkg/handler/resource_update_controller_test.go`
- ✅ `go test ./pkg/handler/...`
- ✅ `bun run build` (`web/`)

### 🧩 Edit Form Select Initial Value Düzeltmesi (Dependent Fields / target_type)

Edit formda backend `target_type` değeri gelse bile select alanının placeholder göstermesi sorunu giderildi.

#### Frontend

- `web/src/components/fields/form/SelectField.tsx`:
  - Select value normalize akışı güçlendirildi (`string`, `object`, JSON-string payload desteği).
  - RHF değeri boş geldiğinde `field.data` fallback'i ile seçili değer korunuyor.
  - Fallback değer form state'e senkronize edilerek dependency resolver ile tutarlılık sağlandı.
- `web/src/pages/resource/index.tsx`:
  - Edit initial data üretiminde select alanları normalize edilerek initialize ediliyor.
  - `target_type` eksik/boş payload senaryosunda `product_id` / `category_id` / `static_url` üzerinden güvenli infer eklendi.

#### Sonuç

- Edit modal açılışında `Hedef Tipi` alanı artık kayıtlı değeri seçili gösterir.
- `depends_on("target_type")` ile kontrol edilen alanlar ilk render'da doğru görünür/aktif olur.

#### Doğrulama

- ✅ `bun run build` (`web/`)

### 📊 Chart Widget Modernizasyonu (shadcn/ui + Dinamik Series)

Dashboard chart kartları shadcn/ui örneklerine taşındı ve backend/frontend veri sözleşmesi genişletildi.

#### Frontend

- `trend-metric`, `partition-metric` ve `progress-metric` bileşenleri shadcn/ui chart bileşenleri ile hizalandı.
- `progress-metric` için seri yönetimi dinamik hale getirildi:
  - `series` artık map yapısında (`desktop/mobile` zorunlu değil).
  - `seriesOrder` ile sıra kontrolü desteklendi.
  - `activeSeries` alias veya data key ile çözümleniyor.
- `ProgressMetric` ve `TrendMetric` kartlarında hardcoded alt başlık kaldırıldı; `subtitle`/`description` payload'dan okunuyor.
- Tarih/sayı formatları `Intl.DateTimeFormat` ve `Intl.NumberFormat` ile tarayıcı locale'ına göre render ediliyor.
- `web/src/main.tsx` içinde `html[lang]` ve `dir` değerleri güvenli şekilde set edilerek i18n formatlaması garanti altına alındı.

#### Backend

- `pkg/metric/metric.go` içinde `ProgressMetric` seri modeli generic hale getirildi.
- `SetSeriesLabel`, `SetSeriesColor`, `SetSeriesEnabled`, `SetSeriesKey`, `SetActiveSeries` metodları dinamik seri key'leriyle çalışacak şekilde güncellendi.
- `Resolve()` çıktısına `series`, `activeSeries`, `seriesOrder`, `subtitle` alanları eklendi.
- Line chart için history normalize/fallback üretimi dinamik seri sayısına göre çalışacak şekilde güncellendi.

#### Dokümantasyon

- `docs/Charts-Data-Contract.md` güncellendi (dinamik `series`, `seriesOrder`, `activeSeries`).
- `docs/Widgets.md` güncellendi (yeni progress kullanım örnekleri ve troubleshooting notları).

#### Doğrulama

- ✅ `go test ./pkg/widget ./pkg/metric ./pkg/handler`
- ✅ `bun run build` (`web/`)

### 🛡️ Dependency Resolver CSRF 403 Düzeltmesi

Dependency resolver endpoint'ine giden isteklerde CSRF header eksikliği nedeniyle oluşan `403` hatası giderildi.

#### Frontend

- `web/src/hooks/useFormDependencies.ts` içinde dependency çözümleme çağrısı `fetch` yerine axios tabanlı `resourceService.resolveDependencies(...)` üzerinden çalışacak şekilde güncellendi.
- Böylece `/api/resource/:resource/fields/resolve-dependencies` çağrılarında session + CSRF akışı diğer API çağrılarıyla aynı hale getirildi.
- `target_type` gibi dependency tetikleyen alan değişimlerinde görülen 403 sorunu çözüldü.

#### Doğrulama

- ✅ `bun run build` (`web/`)

### 🎨 Dashboard Kart Grid Width Desteği (Frontend)

Dashboard ve resource/lens kart grid yerleşimlerinde `card.width` değerinin gerçekten uygulanması sağlandı.

#### Frontend

- Ortak helper eklendi: `web/src/lib/card-grid.ts`
  - Yeni fonksiyon: `getCardGridSpan(width?: string): string`
  - Desteklenen width mapping:
    - `full` → `col-span-1 md:col-span-2 lg:col-span-6 xl:col-span-12`
    - `3/4` → `col-span-1 md:col-span-2 lg:col-span-5 xl:col-span-9`
    - `2/3` → `col-span-1 md:col-span-2 lg:col-span-4 xl:col-span-8`
    - `1/2` → `col-span-1 md:col-span-1 lg:col-span-3 xl:col-span-6`
    - `1/4` → `col-span-1 md:col-span-1 lg:col-span-2 xl:col-span-3`
    - varsayılan (`1/3`) → `col-span-1 md:col-span-1 lg:col-span-2 xl:col-span-4`
- Aşağıdaki ekranlarda hardcoded kart span kaldırıldı ve helper kullanıldı:
  - `web/src/pages/common/page-viewer.tsx`
  - `web/src/pages/resource/index.tsx`
  - `web/src/components/views/LensView.tsx`
- Üç ekranda da kart grid container sınıfı `grid-cols-1 md:grid-cols-2 lg:grid-cols-6 xl:grid-cols-12` olacak şekilde standardize edildi.

#### Doğrulama

- ✅ `bun run build` (`web/`)

### ✨ Resource Index Pagination Tipleri (Links / Simple / Load More)

Resource bazında index sayfası pagination davranışı yönetilebilir hale getirildi.

#### Backend

- Yeni pagination tipi enum'u eklendi:
  - `resource.IndexPaginationTypeLinks` (varsayılan)
  - `resource.IndexPaginationTypeSimple`
  - `resource.IndexPaginationTypeLoadMore`
- `Base` ve `OptimizedBase` için yeni metodlar:
  - `SetIndexPaginationType(...)`
  - `GetIndexPaginationType()`
- Handler seviyesinde pagination tipi resolve edilip varsayılanı `links` olacak şekilde normalize edildi.
- `GET /api/resource/:resource` index yanıtına `meta.pagination.type` alanı eklendi.

Örnek API meta:

```json
{
  "meta": {
    "pagination": {
      "type": "links"
    }
  }
}
```

#### Frontend

- `web/src/components/views/Pagination.tsx` üç modu destekleyecek şekilde genişletildi:
  - `links`: klasik sayılı pagination
  - `simple`: sadece ileri/geri
  - `load_more`: daha fazla yükle
- Resource index sayfası (`web/src/pages/resource/index.tsx`) artık `meta.pagination.type` değerine göre doğru pagination modunu render ediyor.
- `load_more` modunda sayfalar birleştirilerek (append) listede gösteriliyor.
- İlgili type tanımı güncellendi: `web/src/types.ts`
- Pagination testleri güncellendi: `web/src/components/views/Pagination.test.tsx`

#### Dokümantasyon

- `docs/Resources.md` dosyasında **Index Pagination Tipi (`pagination.type`)** bölümü eklendi.
- Desteklenen değerler, kullanım örnekleri ve meta çıktısı dökümante edildi.

#### Doğrulama

- ✅ `go test ./pkg/handler ./pkg/resource`
- ✅ `bun run test src/components/views/Pagination.test.tsx` (`web/`)
- ✅ `bun run build` (`web/`)

#### 🔧 Varsayılan Per Page Güncellemesi

- Resource index için varsayılan `per_page` değeri `15` yerine `10` olarak güncellendi.
- Backend query parser varsayılanı güncellendi: `pkg/query/parser.go`
- Frontend URL param varsayılanı güncellendi: `web/src/lib/resource-params.ts`
- Sonuç: İlk yüklemede pagination select varsayılan olarak `10` gösterir.

### ⚡ Full-Repo Concurrency, Sync, Channel Refactor (Güvenli Kademeli)

Repo genelinde request-path concurrency standardı, cancellation zinciri ve goroutine lifecycle yönetimi güçlendirildi. Değişiklikler kademeli rollout için feature flag yaklaşımı ile eklendi.

#### 🧩 Concurrency Config Genişletmesi

`pkg/panel/config.go` içindeki `ConcurrencyConfig` genişletildi:

- `EnableDataPipelineV2`
- `DataWorkers`
- `EnableMiddlewareV2`
- `EnableOpenAPIV2`
- `OpenAPIWorkers`

Mevcut handler alanları (`EnablePipelineV2`, `FailFast`, `MaxWorkers`, `CardWorkers`, `FieldWorkers`) korunarak backward-compatible şekilde genişletildi.

#### 🗃️ Data Katmanı (GORM Provider)

`pkg/data/gorm_provider.go` içinde relationship lazy-load akışı bounded worker-pool ve cancellation-aware hale getirildi:

- Yeni additive yapı: `RelationshipConcurrencyConfig`
- Yeni additive metod: `SetRelationshipConcurrencyConfig(...)`
- Lazy relationship load işlemleri v2 açıkken bounded pipeline ile çalışır
- Fail-fast davranışı flag üzerinden yönetilir
- V2 kapalıyken legacy davranış korunur

#### 🛡️ Middleware Concurrency/Lifecycle

`pkg/middleware/api_key.go`:

- API key doğrulama için lock-free immutable snapshot modu eklendi
- Yeni additive metod: `SetAtomicSnapshotEnabled(bool)`
- Runtime config güncellemeleri snapshot atomik state üzerinden request-path'e taşınır

`pkg/middleware/security.go`:

- `AccountLockout` için stop edilebilir lifecycle eklendi
- Yeni additive metod: `(*AccountLockout).Close()`
- Cleanup goroutine artık kontrollü şekilde sonlandırılabiliyor

#### 🧭 Panel State Concurrency (Startup-Only Register)

`pkg/panel/app.go` + `pkg/panel/resource_scope.go`:

- Resource/Page registry erişimleri immutable snapshot modeli ile request-path'e taşındı
- Startup sonrası registration freeze davranışı eklendi
- Freeze sonrası `Register` / `RegisterPage` çağrıları no-op + warning log
- `Panel.Start()` başlangıcında freeze uygulanır, `BootPlugins()` sonunda da freeze finalize edilir
- `Panel.Close()` ile background lifecycle cleanup (lockout close) eklendi

`pkg/panel/page_routes.go` ve navigation path'lerinde doğrudan mutable map yerine snapshot okumaları kullanıldı.

#### 🧱 Core Field Clone Altyapısı

`pkg/core/clone.go` eklendi:

- Yeni additive interface: `ElementCloner` (`Clone() Element`)
- `CloneElement` helper (cloner varsa onu kullanır, yoksa güvenli reflection fallback)

`pkg/core/context.go`:

- `GetOrCloneField(...)` içindeki TODO kaldırıldı
- Gerçek clone + cache akışı aktif hale getirildi

`pkg/handler/field_handler.go`:

- Field izolasyon clone helper'ı `core.CloneElement(...)` ile standardize edildi

#### 📘 OpenAPI Concurrency ve Cache Güvenliği

`pkg/openapi/spec.go`:

- Spec generation için singleflight eklendi (tek üretim)
- Cache get/set immutable clone mantığına taşındı
- Paralel dynamic build opsiyonu config ile bağlandı

`pkg/openapi/dynamic_spec.go`:

- Bounded parallel path/schema üretimi için parallel generator metodları eklendi
- V2 açık değilse mevcut serial üretim davranışı korunur

#### 🧪 Testler ve Stabilizasyon

Eklenen/güncellenen testler:

- `pkg/core/clone_test.go`
- `pkg/middleware/api_key_test.go`
- `pkg/middleware/security_test.go`
- `pkg/openapi/spec_cache_test.go`
- `pkg/panel/panel_test.go`

Panel integration timeout stabilizasyonu için:

- `pkg/panel/test_http_helper_test.go` eklendi
- Panel testlerinde merkezi `testFiberRequest(...)` helper'ı ile timeout standardı artırıldı

Doğrulama:

- ✅ `go test ./pkg/core ./pkg/middleware ./pkg/openapi ./pkg/data ./pkg/handler ./pkg/panel`
- ✅ `go test -race ./pkg/handler ./pkg/data ./pkg/middleware ./pkg/panel ./pkg/internal/concurrency`
- ⚠️ `go test -race ./...` tam repo koşusunda refactor dışı mevcut build sorunu (`pkg/metric/metric.go` unused import) nedeniyle kırılmaya devam ediyor

### ✨ Yeni Özellikler (Frontend & Backend)

#### 🚀 Detail View İyileştirmeleri (Nova Benzeri)

Detail (Detay) sayfasındaki ilişki yönetimi ve kullanıcı deneyimi önemli ölçüde geliştirildi.

**Frontend:**
- **Tablo Görünümü:** `HasMany`, `BelongsToMany` ve `MorphToMany` ilişkileri artık detay modalında **Tablo** (`RelationshipTable`) olarak listeleniyor.
- **İç İçe Modallar (Nested Modals):** Bir kaydın detayından, ilişkili başka bir kaydın detayına tıklandığında yeni bir modal açılıyor. Önceki modal kapanmıyor, geri gelindiğinde kaldığı yerden devam ediyor.
- **Dinamik Genişlik:** İlişki tablosu içeren detay modalları otomatik olarak daha geniş (`sm:max-w-5xl`) açılıyor.
- **Search & Pagination:** İlişki tabloları içinde arama yapabilir ve sayfalar arasında gezinebilirsiniz.
- **Deep Linking:** URL üzerinden (`?detail_id=...`) doğrudan detay modalını açma desteği eklendi.

**Backend:**
- **Query Parser Güncellemesi:** `pkg/query/parser.go` güncellendi. Artık `viaResource`, `viaResourceId` ve `viaRelationship` parametreleri destekleniyor. Bu sayede ilişkili kayıtlar (örneğin bir şirkete ait adresler) doğru şekilde filtreleniyor.

#### 📱 Form İyileştirmeleri

- **Tel Field (Phone Input):** `Tel` tipindeki alanlar için gelişmiş `PhoneInput` (ülke bayraklı, formatlı) bileşeni entegre edildi.
- **Akıllı Component Seçimi:** Backend `text-field` view'ı gönderse bile, eğer alanın tipi `tel` ise frontend otomatik olarak `TelInput` bileşenini kullanıyor.

#### Resource Title Pattern (Nova Uyumlu)

Panel.go'ya title pattern'i eklendi. Her resource için kayıt başlığı (record title) özelliği artık kullanılabilir. Bu, ilişki fieldlarında kayıtların okunabilir şekilde gösterilmesini sağlar.

**Özellikler:**
- `SetRecordTitleKey(key string)` - Kayıt başlığı için kullanılacak field adını ayarlar
- `GetRecordTitleKey() string` - Kayıt başlığı için kullanılacak field adını döndürür
- `SetRecordTitleFunc(fn func(any) string)` - Özel başlık fonksiyonu ayarlar
- `RecordTitle(record any) string` - Kaydın okunabilir başlığını döndürür

**Kullanım Örneği:**

```go
// UserResource'da "name" field'ını başlık olarak ayarla
func NewUserResource() *UserResource {
    r := &UserResource{}
    r.SetModel(&User{})
    r.SetSlug("users")
    r.SetRecordTitleKey("name") // ← Yeni özellik
    return r
}

// Özel başlık fonksiyonu ile
r.SetRecordTitleFunc(func(record any) string {
    user := record.(*User)
    return user.FirstName + " " + user.LastName
})
```

**İlişki Fieldları:**

Tüm ilişki fieldları artık minimal format döndürür: `{"id": ..., "title": ...}`

- **BelongsTo**: `{"id": 5, "title": "John Doe"}`
- **HasMany**: `[{"id": 1, "title": "First Post"}, {"id": 2, "title": "Second Post"}]`
- **HasOne**: `{"id": 1, "title": "User Profile"}`
- **BelongsToMany**: `[{"id": 1, "title": "Admin"}, {"id": 2, "title": "Editor"}]`

**Etkilenen Dosyalar:**
- `pkg/resource/resource.go` - Interface'e yeni metodlar eklendi
- `pkg/resource/optimized.go` - OptimizedBase implementation
- `pkg/resource/base.go` - Base implementation
- `pkg/fields/belongs_to.go` - Extract metodu eklendi
- `pkg/fields/has_many.go` - Extract metodu güncellendi
- `pkg/fields/has_one.go` - Extract metodu güncellendi
- `pkg/fields/belongs_to_many.go` - Extract metodu eklendi
- `pkg/resource/user/resource.go` - SetRecordTitleKey("name") eklendi
- `pkg/resource/account/resource.go` - SetRecordTitleKey("name") eklendi
- `pkg/resource/session/resource.go` - SetRecordTitleKey("id") eklendi
- `pkg/resource/verification/resource.go` - SetRecordTitleKey("id") eklendi

**Testler:**
- `pkg/resource/record_title_test.go` - RecordTitle için kapsamlı testler eklendi
- Tüm testler başarıyla çalışıyor ✅

### 🔧 Düzeltmeler

#### Base Resource Bug Fix

`Base.SetDialogType` ve `Base.SetOpenAPIEnabled` metodları pointer receiver'a çevrildi. Bu metodlar value receiver kullanıyordu ve değişiklikler kayboluyordu.

**Önceki (Hatalı):**
```go
func (b Base) SetDialogType(dialogType DialogType) Resource {
    b.DialogType = dialogType // Değişiklik kaybolur (kopya üzerinde)
    return b
}
```

**Sonrası (Düzeltilmiş):**
```go
func (b *Base) SetDialogType(dialogType DialogType) Resource {
    b.DialogType = dialogType // Değişiklik kalıcı
    return b
}
```

### ⚠️ Breaking Changes

1. **İlişki Field Serialize Formatı**: HasMany, HasOne, BelongsToMany fieldları artık `{"id": ..., "title": ...}` formatında döndürüyor (önceden tam kayıt veya sadece ID döndürüyordu)

2. **Base Resource Metodları**: `SetDialogType` ve `SetOpenAPIEnabled` metodları pointer receiver'a çevrildi

### 📝 Önemli Notlar

- **Eager Loading Zorunlu**: İlişki fieldlarında eager loading yapılmalı, aksi halde title null olur
- **DisplayUsing Korundu**: Mevcut DisplayUsing() callback'leri çalışmaya devam ediyor
- **Type Assertion**: RelatedResource interface{} tipinde olduğu için type assertion kullanıldı
- **MorphTo**: TypeMappings map[string]string olduğu için (resource slug'ları tutuyor) title pattern uygulanmadı

### 🧪 Test Durumu

- ✅ Resource testleri: Tüm testler başarılı
- ✅ RecordTitle testleri: Yeni testler eklendi ve başarılı
- ⚠️ Fields testleri: Mevcut test dosyalarında constructor fonksiyon adları ile ilgili sorunlar var (implementasyondan bağımsız)

### 📚 Dökümantasyon

- CHANGELOG.md oluşturuldu
- RecordTitle için kapsamlı testler ve örnekler eklendi
- Tüm metodlar Türkçe dokümantasyon ile açıklandı

---

## Önceki Sürümler

Önceki sürüm notları için git commit geçmişine bakınız.
