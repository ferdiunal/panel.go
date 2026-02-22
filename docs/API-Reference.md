# API Referansı (Legacy Teknik Akış)

Panel SDK'nın tam API referansı. Tüm metodlar, parametreler ve dönüş değerleri burada belirtilmiştir.

Bu doküman, detaylı ve düşük seviye API bakışı içindir. Uçtan uca başlangıç için önce [Başlarken](Getting-Started) ve [Kaynaklar (Resource)](Resources) dokümanlarını takip edin.

## Bu Doküman Ne Zaman Kullanılmalı?

- Bir metodun imzasını veya dönüş tipini doğrularken
- Field/relationship builder zincirinde hangi seçeneklerin olduğunu ararken
- Policy ve query davranışlarını referans düzeyinde kontrol ederken

## Hızlı Referans Akışı

1. Resource davranışı için: `Resource Interface`
2. Alan davranışı için: `Field Interface`
3. İlişki kurulumları için: `Relationship Fields`
4. Yetki kontrolleri için: `Policy Interface`
5. Sorgu detayları için: `Query Builder`

## İçindekiler

- [Resource Interface](#resource-interface)
- [Field Interface](#field-interface)
- [Relationship Fields](#relationship-fields)
- [Policy Interface](#policy-interface)
- [Query Builder](#query-builder)
- [Index ve Lens View Parametreleri](#index-ve-lens-view-parametreleri)
- [Context](#context)
- [Error Handling](#error-handling)

---

## Resource Interface

### Temel Metodlar

#### `Model() interface{}`

Kaynağın temsil ettiği veri modelini döndürür.

```go
func (r *UserResource) Model() interface{} {
    return &User{}
}
```

**Dönüş Değeri:** Veri modeli (struct pointer)

---

#### `Slug() string`

Kaynağın URL slug'ını döndürür.

```go
func (r *UserResource) Slug() string {
    return "users"
}
```

**Dönüş Değeri:** Slug string'i (örn: "users", "products")

---

#### `Title() string`

Kaynağın başlığını döndürür.

```go
func (r *UserResource) Title() string {
    return "Kullanıcılar"
}
```

**Dönüş Değeri:** Başlık string'i

---

#### `Icon() string`

Kaynağın ikonunu döndürür.

```go
func (r *UserResource) Icon() string {
    return "user"
}
```

**Dönüş Değeri:** İkon adı (örn: "user", "product", "cog")

---

#### `Group() string`

Kaynağın ait olduğu grubu döndürür.

```go
func (r *UserResource) Group() string {
    return "Yönetim"
}
```

**Dönüş Değeri:** Grup adı

---

#### `Fields() []fields.Element`

Kaynağın alanlarını döndürür.

```go
func (r *UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("İsim", "name"),
        fields.Email("E-posta", "email"),
    }
}
```

**Dönüş Değeri:** Alan listesi

---

#### `SortableBy() []string`

Sıralanabilir alanları döndürür.

```go
func (r *UserResource) SortableBy() []string {
    return []string{"name", "email", "created_at"}
}
```

**Dönüş Değeri:** Sıralanabilir alan adları

---

#### `PerPage() int`

Sayfa başına kayıt sayısını döndürür.

```go
func (r *UserResource) PerPage() int {
    return 15
}
```

**Dönüş Değeri:** Sayfa başına kayıt sayısı (varsayılan: 15)

---

### Yaşam Döngüsü Metodları

#### `BeforeCreate(ctx context.Context, data map[string]interface{}) error`

Kayıt oluşturulmadan önce çalışır.

```go
func (r *UserResource) BeforeCreate(ctx context.Context, data map[string]interface{}) error {
    // Şifreyi hash'le
    if password, ok := data["password"].(string); ok {
        data["password"] = hashPassword(password)
    }
    return nil
}
```

**Parametreler:**
- `ctx`: Context
- `data`: Oluşturulacak veri

**Dönüş Değeri:** Hata (nil başarılı)

---

#### `AfterCreate(ctx context.Context, model interface{}) error`

Kayıt oluşturulduktan sonra çalışır.

```go
func (r *UserResource) AfterCreate(ctx context.Context, model interface{}) error {
    user := model.(*User)
    sendWelcomeEmail(user.Email)
    return nil
}
```

**Parametreler:**
- `ctx`: Context
- `model`: Oluşturulan model

**Dönüş Değeri:** Hata (nil başarılı)

---

#### `BeforeUpdate(ctx context.Context, id interface{}, data map[string]interface{}) error`

Kayıt güncellenecekten önce çalışır.

```go
func (r *UserResource) BeforeUpdate(ctx context.Context, id interface{}, data map[string]interface{}) error {
    // Doğrulama
    return nil
}
```

**Parametreler:**
- `ctx`: Context
- `id`: Kayıt ID'si
- `data`: Güncellenecek veri

**Dönüş Değeri:** Hata (nil başarılı)

---

#### `AfterUpdate(ctx context.Context, model interface{}) error`

Kayıt güncellendikten sonra çalışır.

```go
func (r *UserResource) AfterUpdate(ctx context.Context, model interface{}) error {
    user := model.(*User)
    logAction("user_updated", user.ID)
    return nil
}
```

**Parametreler:**
- `ctx`: Context
- `model`: Güncellenen model

**Dönüş Değeri:** Hata (nil başarılı)

---

#### `BeforeDelete(ctx context.Context, id interface{}) error`

Kayıt silinecekten önce çalışır.

```go
func (r *UserResource) BeforeDelete(ctx context.Context, id interface{}) error {
    // Silme öncesi kontrol
    return nil
}
```

**Parametreler:**
- `ctx`: Context
- `id`: Kayıt ID'si

**Dönüş Değeri:** Hata (nil başarılı)

---

#### `AfterDelete(ctx context.Context, id interface{}) error`

Kayıt silindikten sonra çalışır.

```go
func (r *UserResource) AfterDelete(ctx context.Context, id interface{}) error {
    logAction("user_deleted", id)
    return nil
}
```

**Parametreler:**
- `ctx`: Context
- `id`: Silinen kayıt ID'si

**Dönüş Değeri:** Hata (nil başarılı)

---

## Field Interface

### Temel Metodlar

#### `Resolve(data interface{}) interface{}`

Veri modelinden alan değerini çözer.

```go
func (f *TextField) Resolve(data interface{}) interface{} {
    user := data.(*User)
    return user.Name
}
```

**Parametreler:**
- `data`: Veri modeli

**Dönüş Değeri:** Alan değeri

---

#### `Serialize(data interface{}) map[string]interface{}`

Alanı JSON'a dönüştürür.

```go
func (f *TextField) Serialize(data interface{}) map[string]interface{} {
    return map[string]interface{}{
        "type":  "text",
        "value": f.Resolve(data),
    }
}
```

**Parametreler:**
- `data`: Veri modeli

**Dönüş Değeri:** JSON map'i

---

#### `Validate(value interface{}) error`

Alan değerini doğrular.

```go
func (f *TextField) Validate(value interface{}) error {
    str, ok := value.(string)
    if !ok {
        return fmt.Errorf("string olmalı")
    }
    if len(str) == 0 {
        return fmt.Errorf("boş olamaz")
    }
    return nil
}
```

**Parametreler:**
- `value`: Doğrulanacak değer

**Dönüş Değeri:** Hata (nil başarılı)

---

### Ortak UI Props (`WithProps`)

Field UI davranışı için `WithProps(key, value)` ile frontend'e ek metadata geçebilirsiniz.

Addon destekli anahtarlar:
- `startAddon`: Alanın başında render edilir (Form/Detail/Index)
- `endAddon`: Alanın sonunda render edilir (Form/Detail/Index)

Uyumluluk alias'ları:
- Baş addon: `start_component`, `prefix`, `prepend`
- Son addon: `end_component`, `suffix`, `append`

```go
fields.Text("Website", "website").
    OnForm().
    WithProps("startAddon", "https://")

fields.Number("Fiyat", "price").
    OnForm().
    WithProps("startAddon", "₺").
    WithProps("endAddon", "/ay")
```

Detaylı field davranışları için: [Fields](Fields)

---

### Görünürlük Metodları

#### `OnlyOnIndex() *TextField`

Alanı sadece liste sayfasında göster.

```go
fields.Text("İsim", "name").OnlyOnIndex()
```

**Dönüş Değeri:** Alan (fluent API)

---

#### `OnlyOnDetail() *TextField`

Alanı sadece detay sayfasında göster.

```go
fields.Text("İsim", "name").OnlyOnDetail()
```

**Dönüş Değeri:** Alan (fluent API)

---

#### `OnlyOnForms() *TextField`

Alanı sadece form'da göster.

```go
fields.Text("İsim", "name").OnlyOnForms()
```

**Dönüş Değeri:** Alan (fluent API)

---

#### `HiddenIf(callback func(interface{}) bool) *TextField`

Koşula göre alanı gizle.

```go
fields.Text("Şifre", "password").
    HiddenIf(func(data interface{}) bool {
        return true // Gizle
    })
```

**Parametreler:**
- `callback`: Gizleme koşulu

**Dönüş Değeri:** Alan (fluent API)

---

### Arama ve Sıralama Metodları

#### `Searchable() *TextField`

Alanı aranabilir yap.

```go
fields.Text("İsim", "name").Searchable()
```

**Dönüş Değeri:** Alan (fluent API)

---

#### `Sortable() *TextField`

Alanı sıralanabilir yap.

```go
fields.Text("İsim", "name").Sortable()
```

**Dönüş Değeri:** Alan (fluent API)

---

#### `WithSearchableColumns(columns ...string) *BelongsTo`

İlişkili alanda aranabilir sütunları belirle.

```go
fields.BelongsTo("Yazar", "user_id", "users").
    WithSearchableColumns("name", "email")
```

**Parametreler:**
- `columns`: Sütun adları

**Dönüş Değeri:** Alan (fluent API)

---

### Layout Metodları

#### `Span(columns int) Element`

Form ve detail görünümünde alanın 12 kolon grid içindeki genişliğini ayarlar.

```go
fields.Text("Ad", "first_name").Span(6)
fields.Text("Soyad", "last_name").Span(6)
```

**Parametreler:**
- `columns`: Kolon genişliği (`1..12`)

**Notlar:**
- Varsayılan span değeri `12` (tam genişlik)
- Geçersiz değerler clamp edilir (`<1 => 1`, `>12 => 12`)
- Sadece form/detail görünümünde etkilidir

**Dönüş Değeri:** Alan (fluent API)

---

#### `ShowNumberControls(show bool) Element`

Number field formundaki artı/eksi butonlarının görünürlüğünü ayarlar.

```go
fields.Number("Fiyat", "price").ShowNumberControls(false)
```

**Parametreler:**
- `show`: `true` ise kontroller görünür, `false` ise gizli

**Dönüş Değeri:** Alan (fluent API)

---

#### `HideNumberControls() Element`

Number field formundaki artı/eksi butonlarını gizler (`ShowNumberControls(false)` kısayolu).

```go
fields.Number("Fiyat", "price").HideNumberControls()
```

**Dönüş Değeri:** Alan (fluent API)

---

### Doğrulama Metodları

#### `Required() *TextField`

Alanı zorunlu yap.

```go
fields.Text("İsim", "name").Required()
```

**Dönüş Değeri:** Alan (fluent API)

---

#### `Rules(rules ...interface{}) *TextField`

Doğrulama kuralları ekle.

```go
fields.Text("İsim", "name").
    Rules(
        validate.MinLength(3),
        validate.MaxLength(100),
    )
```

**Parametreler:**
- `rules`: Doğrulama kuralları

**Dönüş Değeri:** Alan (fluent API)

---

## Relationship Fields

### BelongsTo

```go
fields.BelongsTo(label, foreignKey, table string) *BelongsTo
```

**Metodlar:**
- `DisplayUsing(key string) *BelongsTo` - Gösterilecek alanı belirle
- `WithSearchableColumns(columns ...string) *BelongsTo` - Aranabilir alanları belirle
- `Query(callback func(*Query) *Query) *BelongsTo` - Query'yi özelleştir
- `WithEagerLoad() *BelongsTo` - Eager loading kullan
- `WithLazyLoad() *BelongsTo` - Lazy loading kullan
- `Required() *BelongsTo` - Zorunlu yap

**Örnek:**
```go
fields.BelongsTo("Yazar", "user_id", "users").
    DisplayUsing("name").
    WithSearchableColumns("name", "email").
    WithEagerLoad()
```

---

### HasMany

```go
fields.HasMany(label, relationName, table string) *HasMany
```

**Metodlar:**
- `ForeignKey(key string) *HasMany` - Foreign key belirle
- `OwnerKey(key string) *HasMany` - Owner key belirle
- `Query(callback func(*Query) *Query) *HasMany` - Query'yi özelleştir
- `WithEagerLoad() *HasMany` - Eager loading kullan
- `WithLazyLoad() *HasMany` - Lazy loading kullan

**Örnek:**
```go
fields.HasMany("Yazılar", "posts", "posts").
    ForeignKey("author_id").
    WithEagerLoad()
```

---

### HasOne

```go
fields.HasOne(label, relationName, table string) *HasOne
```

**Metodlar:**
- `ForeignKey(key string) *HasOne` - Foreign key belirle
- `OwnerKey(key string) *HasOne` - Owner key belirle
- `Query(callback func(*Query) *Query) *HasOne` - Query'yi özelleştir

**Örnek:**
```go
fields.HasOne("Profil", "profile", "profiles").
    ForeignKey("user_id")
```

---

### BelongsToMany

```go
fields.BelongsToMany(label, pivotTable, table string) *BelongsToMany
```

**Metodlar:**
- `PivotTable(name string) *BelongsToMany` - Pivot table adını belirle
- `ForeignKey(key string) *BelongsToMany` - Foreign key belirle
- `RelatedKey(key string) *BelongsToMany` - Related key belirle
- `Query(callback func(*Query) *Query) *BelongsToMany` - Query'yi özelleştir

**Örnek:**
```go
fields.BelongsToMany("Roller", "role_user", "roles").
    PivotTable("role_user").
    ForeignKey("user_id").
    RelatedKey("role_id")
```

---

### MorphTo

```go
fields.MorphTo(label, morphKey string) *MorphTo
```

**Metodlar:**
- `Types(types map[string]string) *MorphTo` - Type mapping belirle

**Örnek:**
```go
fields.MorphTo("Yorumlanabilir", "commentable").
    Types(map[string]string{
        "post":  "posts",
        "video": "videos",
    })
```

---

## Policy Interface

### Temel Metodlar

#### `ViewAny(ctx context.Context) bool`

Herhangi bir kaynağı görüntüleme izni.

```go
func (p *UserPolicy) ViewAny(ctx context.Context) bool {
    user := ctx.Value("user").(*User)
    return user.Role == "admin"
}
```

**Parametreler:**
- `ctx`: Context

**Dönüş Değeri:** İzin (true/false)

---

#### `View(ctx context.Context, model interface{}) bool`

Belirli bir kaynağı görüntüleme izni.

```go
func (p *UserPolicy) View(ctx context.Context, model interface{}) bool {
    user := ctx.Value("user").(*User)
    target := model.(*User)
    return user.ID == target.ID || user.Role == "admin"
}
```

**Parametreler:**
- `ctx`: Context
- `model`: Kaynak modeli

**Dönüş Değeri:** İzin (true/false)

---

#### `Create(ctx context.Context) bool`

Kaynak oluşturma izni.

```go
func (p *UserPolicy) Create(ctx context.Context) bool {
    user := ctx.Value("user").(*User)
    return user.Role == "admin"
}
```

**Parametreler:**
- `ctx`: Context

**Dönüş Değeri:** İzin (true/false)

---

#### `Update(ctx context.Context, model interface{}) bool`

Kaynak güncelleme izni.

```go
func (p *UserPolicy) Update(ctx context.Context, model interface{}) bool {
    user := ctx.Value("user").(*User)
    target := model.(*User)
    return user.ID == target.ID || user.Role == "admin"
}
```

**Parametreler:**
- `ctx`: Context
- `model`: Kaynak modeli

**Dönüş Değeri:** İzin (true/false)

---

#### `Delete(ctx context.Context, model interface{}) bool`

Kaynak silme izni.

```go
func (p *UserPolicy) Delete(ctx context.Context, model interface{}) bool {
    user := ctx.Value("user").(*User)
    return user.Role == "admin"
}
```

**Parametreler:**
- `ctx`: Context
- `model`: Kaynak modeli

**Dönüş Değeri:** İzin (true/false)

---

#### `Restore(ctx context.Context, model interface{}) bool`

Kaynak geri yükleme izni (soft delete).

```go
func (p *UserPolicy) Restore(ctx context.Context, model interface{}) bool {
    user := ctx.Value("user").(*User)
    return user.Role == "admin"
}
```

**Parametreler:**
- `ctx`: Context
- `model`: Kaynak modeli

**Dönüş Değeri:** İzin (true/false)

---

#### `ForceDelete(ctx context.Context, model interface{}) bool`

Kaynak kalıcı silme izni.

```go
func (p *UserPolicy) ForceDelete(ctx context.Context, model interface{}) bool {
    user := ctx.Value("user").(*User)
    return user.Role == "admin"
}
```

**Parametreler:**
- `ctx`: Context
- `model`: Kaynak modeli

**Dönüş Değeri:** İzin (true/false)

---

## Query Builder

### Metodlar

#### `Where(column, operator, value string) *Query`

WHERE clause ekle.

```go
query.Where("status", "=", "active")
```

**Parametreler:**
- `column`: Sütun adı
- `operator`: Operatör (=, !=, >, <, >=, <=, LIKE)
- `value`: Değer

**Dönüş Değeri:** Query (fluent API)

---

#### `WhereIn(column string, values []interface{}) *Query`

WHERE IN clause ekle.

```go
query.WhereIn("status", []interface{}{"active", "pending"})
```

**Parametreler:**
- `column`: Sütun adı
- `values`: Değer listesi

**Dönüş Değeri:** Query (fluent API)

---

#### `OrderBy(column, direction string) *Query`

ORDER BY ekle.

```go
query.OrderBy("created_at", "DESC")
```

**Parametreler:**
- `column`: Sütun adı
- `direction`: Yön (ASC, DESC)

**Dönüş Değeri:** Query (fluent API)

---

#### `Limit(count int) *Query`

LIMIT ekle.

```go
query.Limit(10)
```

**Parametreler:**
- `count`: Limit sayısı

**Dönüş Değeri:** Query (fluent API)

---

#### `Offset(count int) *Query`

OFFSET ekle.

```go
query.Offset(20)
```

**Parametreler:**
- `count`: Offset sayısı

**Dönüş Değeri:** Query (fluent API)

---

#### `Select(columns ...string) *Query`

SELECT sütunları belirle.

```go
query.Select("id", "name", "email")
```

**Parametreler:**
- `columns`: Sütun adları

**Dönüş Değeri:** Query (fluent API)

---

## Index ve Lens View Parametreleri

Resource index, relationship index ve lens endpoint'lerinde `table/grid` görünümü
`view` query parametresi ile yönetilir.

### Varsayılan davranış

- `view` gönderilmezse varsayılan görünüm `table` olur.

### Resource index (nested query format)

```http
GET /api/resource/users?users[view]=grid
```

Alternatif:
```http
GET /api/resource/users?view=grid
```

### Relationship index (nested query + via parametreleri)

```http
GET /api/resource/tags?tags[view]=grid&viaResource=products&viaResourceId=42&viaRelationship=tags
```

### Lens endpoint

```http
GET /api/resource/users/lens/active-users?view=grid
```

### Desteklenen değerler

- `table`
- `grid`

Geçersiz bir değer gelirse backend otomatik olarak `table` kabul eder.

### Response'a eklenen başlık alanı

Grid kart başlığı için backend, resource tarafındaki `GetRecordTitleKey()` değerini döner:

- Resource index response: `meta.grid_enabled`
- Resource index response: `meta.record_title_key`
- Lens response: `grid_enabled` (top-level)
- Lens response: `record_title_key` (top-level)

Grid visibility notu:
- `HideOnGrid` alanları `headers`/kart listing tarafında filtreler.
- Aynı alanların row payload (`data`) içinde kalması beklenen davranıştır.

Örnek (resource index):

```json
{
  "data": [...],
  "meta": {
    "grid_enabled": true,
    "record_title_key": "name",
    "headers": [...]
  }
}
```

Örnek (lens):

```json
{
  "name": "Active Users",
  "grid_enabled": true,
  "record_title_key": "full_name",
  "headers": [...],
  "data": [...]
}
```

### Resource tarafında title key ayarı

```go
r.SetRecordTitleKey("name")
r.SetGridEnabled(true) // optional, varsayılan true
```

---

## Context

### Standart Context Anahtarları

#### `"user"`

Mevcut kullanıcı bilgisi.

```go
user := ctx.Value("user").(*User)
```

---

#### `"tenant_id"`

Kiracı ID'si (multi-tenancy).

```go
tenantID := ctx.Value("tenant_id").(string)
```

---

#### `"request"`

HTTP isteği.

```go
request := ctx.Value("request").(*http.Request)
```

---

#### `"response"`

HTTP yanıtı.

```go
response := ctx.Value("response").(http.ResponseWriter)
```

---

## Error Handling

### Hata Türleri

#### `ValidationError`

Doğrulama hatası.

```go
type ValidationError struct {
    Field   string
    Message string
}
```

---

#### `AuthorizationError`

Yetkilendirme hatası.

```go
type AuthorizationError struct {
    Message string
}
```

---

#### `NotFoundError`

Kayıt bulunamadı hatası.

```go
type NotFoundError struct {
    Message string
}
```

---

### Hata Döndürme

```go
func (r *UserResource) BeforeCreate(ctx context.Context, data map[string]interface{}) error {
    if data["email"] == "" {
        return fmt.Errorf("e-posta gerekli")
    }
    return nil
}
```

---

## Sık Hata Kontrolü (API Referansı Kullanırken)

- Yalnızca API referansına bakıp entegrasyon yapmak: Önce [Başlarken](Getting-Started) akışını baz alın.
- Field görünürlük metodlarını karıştırmak: `OnlyOn...` ve `HideOn...` kombinasyonlarını birlikte kontrol edin.
- Relationship'te yalnızca field builder'a bakmak: `resource.Register(...)` ve slug doğruluğunu ayrıca kontrol edin.
- Policy imzalarını eksik implemente etmek: Resource tarafındaki policy beklentisiyle birebir eşleştiğinden emin olun.

---

## Ayrıca Bkz.

- [Alanlar (Fields)](Fields) - Field system genel bakış
- [Kaynaklar (Resource)](Resources) - Resource tanımı
- [İlişkiler (Relationships)](Relationships) - Relationship fields
- [Yetkilendirme](Authorization) - Policy tanımı
- [Gelişmiş Kullanım](Advanced-Usage) - İleri seviye özellikler
