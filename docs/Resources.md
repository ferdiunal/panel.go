# Resource Rehberi (Legacy Teknik Akış)

Resource, Panel.go'da bir veri modelini panelde yönetilebilir hale getiren ana birimdir.

Bu doküman, özellikle legacy/düşük seviye akışta `model + field resolver + policy + repository + register` zincirini doğru kurmak için hazırlanmıştır.

## Resource Ne Tanımlar?

Bir resource tipik olarak şu parçaları içerir:
- Model
- Field'lar
- Policy
- Repository
- Menü bilgileri (`slug`, `title`, `icon`, `group`)

## Hızlı Akış (Önerilen Sıra)

1. Model'i oluştur
2. Field resolver'ı tanımla
3. Policy'yi ekle
4. Resource struct'ını kur
5. `resource.Register(...)` ile global registry'ye kaydet
6. Panel başlangıcında resource'u yükle

## Minimum Çalışan Resource

```go
package organization

import (
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/data"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "gorm.io/gorm"
)

type Organization struct {
    ID   uint   `gorm:"primaryKey"`
    Name string
}

func init() {
    resource.Register("organizations", NewOrganizationResource())
}

type OrganizationResource struct {
    resource.OptimizedBase
}

type OrganizationFieldResolver struct{}

type OrganizationPolicy struct{}

func NewOrganizationResource() *OrganizationResource {
    r := &OrganizationResource{}

    r.SetModel(&Organization{})
    r.SetSlug("organizations")
    r.SetTitle("Organizations")
    r.SetIcon("building")
    r.SetGroup("Operations")
    r.SetVisible(true)
    r.SetNavigationOrder(10)

    // Relationship alanlarında insan-okunur etiket için kritik
    r.SetRecordTitleKey("name")

    r.SetFieldResolver(&OrganizationFieldResolver{})
    r.SetPolicy(&OrganizationPolicy{})

    return r
}

func (r *OrganizationResource) Repository(db *gorm.DB) data.DataProvider {
    return data.NewGormDataProvider(db, &Organization{})
}

func (r *OrganizationFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    return []fields.Element{
        fields.ID("ID"),
        fields.Text("Name", "name").Required().OnList().OnDetail().OnForm(),
    }
}

func (p *OrganizationPolicy) ViewAny(ctx *context.Context) bool { return true }
func (p *OrganizationPolicy) View(ctx *context.Context, model any) bool { return true }
func (p *OrganizationPolicy) Create(ctx *context.Context) bool { return true }
func (p *OrganizationPolicy) Update(ctx *context.Context, model any) bool { return true }
func (p *OrganizationPolicy) Delete(ctx *context.Context, model any) bool { return true }
```

## Resource Register Akışı

İlişkilerin ve `AutoOptions` gibi özelliklerin stabil çalışması için register akışını doğru kurmak kritiktir.

### 1) Global Registry API (`pkg/resource`)

```go
resource.Register(slug string, res resource.Resource)
resource.Get(slug string) resource.Resource
resource.List() map[string]resource.Resource
resource.Clear()
```

Notlar:
- `Register`, doğrudan resource instance alır.
- Çekirdek API'de `resource.GetOrPanic` yoktur.
- `Clear` test amaçlıdır.

### 2) `init()` ile otomatik kayıt

En yaygın desen:

```go
func init() {
    resource.Register("products", NewProductResource())
}
```

Böylece package import edildiğinde resource global registry'ye girer.

### 3) Panel başlangıcında resource yükleme

Panel başlatılırken:
- Config içinden gelen resource'lar register edilir.
- Global registry'deki resource'lar (`resource.List()`) otomatik panel'e eklenir.

Bu nedenle register edilmiş resource'lar, panel tarafında ayrıca tek tek eklenmeden de çalışabilir.

## Resource Konfigürasyonu

```go
r.SetSlug("products")
r.SetTitle("Ürünler")
r.SetIcon("shopping-bag")
r.SetGroup("Satış")
r.SetNavigationOrder(10)
r.SetVisible(true)
```

## Dialog Tipi ve Genişliği

Resource bazında create/edit modal davranışını yönetebilirsiniz:

```go
r.SetDialogType(resource.DialogTypeSheet)   // sheet | drawer | modal
r.SetDialogSize(resource.DialogSize4XL)     // sm | md | lg | xl | 2xl | 3xl | 4xl | 5xl | full
```

Notlar:
- `SetDialogSize` değeri frontend'e `meta.dialog_size` olarak iletilir.
- `dialog_size` gönderilmezse frontend varsayılanı `md` kullanılır.
- Aynı değer hem dialog hem sheet genişliği için uygulanır.

## Index Table Davranışı

`IndexView` üzerinde satır tıklama aksiyonu ve drag-drop reorder artık resource tarafından yönetilebilir.

### 1) Satır Tıklama Aksiyonu (`row_click_action`)

Varsayılan davranış `edit` modalını açmaktır. İsterseniz satır tıklamasını `detail` modalına çevirebilirsiniz.

```go
// Varsayılan (opsiyonel, yazmasanız da edit)
r.SetIndexRowClickAction(resource.IndexRowClickActionEdit)

// Satıra tıklanınca detay modalı aç
r.SetIndexRowClickAction(resource.IndexRowClickActionDetail)
```

Frontend tarafına `GET /api/resource/:resource` yanıtında şu meta alanı gelir:

```json
{
  "meta": {
    "row_click_action": "edit"
  }
}
```

Desteklenen değerler:
- `edit`
- `detail`

### 2) Index Pagination Tipi (`pagination.type`)

Resource bazında index sayfasında hangi pagination UI'nin kullanılacağını belirleyebilirsiniz.

```go
// Klasik sayfa numaraları (varsayılan)
r.SetIndexPaginationType(resource.IndexPaginationTypeLinks)

// Sadece İleri / Geri
r.SetIndexPaginationType(resource.IndexPaginationTypeSimple)

// Daha fazla yükle
r.SetIndexPaginationType(resource.IndexPaginationTypeLoadMore)
```

Frontend tarafına `GET /api/resource/:resource` yanıtında şu meta alanı gelir:

```json
{
  "meta": {
    "pagination": {
      "type": "links"
    }
  }
}
```

Desteklenen değerler:
- `links`: Klasik sayılı pagination
- `simple`: İleri / geri butonları
- `load_more`: Daha fazla yükle davranışı (append)

Notlar:
- Varsayılan değer `links` olarak normalize edilir.
- Frontend, `type` değerine göre otomatik uygun pagination component'ini render eder.

### 3) Satır Drag-Drop Reorder

Tablo satırlarının sürükle-bırak ile yeniden sıralanması için resource bazında bir order kolonu tanımlanır.

```go
// Kısa kullanım
r.EnableIndexReorder("order_column")

// Alternatif
r.SetIndexReorder(true, "order_column")

// Kapatmak için
r.DisableIndexReorder()
```

Frontend tarafına `meta.reorder` bilgisi döner:

```json
{
  "meta": {
    "reorder": {
      "enabled": true,
      "column": "order_column"
    }
  }
}
```

Notlar:
- `enabled=true` olsa bile `column` boşsa reorder kapalı kabul edilir.
- Reorder güncellemesi, listede görünen sıraya göre `1..n` şeklinde yazılır.

### 4) Reorder API Endpoint'i

Satır sıralaması değiştiğinde frontend aşağıdaki endpoint'i çağırır:

`POST /api/resource/:resource/reorder`

İstek gövdesi:

```json
{
  "ids": [7, 3, 9]
}
```

Davranış:
- `ids` sırası yeni tablo sırasıdır.
- Resource'da tanımlı reorder kolonu (`meta.reorder.column`) transaction içinde güncellenir.
- Policy kontrolünde `Update` yetkisi gerekir.

Örnek başarılı yanıt:

```json
{
  "success": true,
  "column": "order_column",
  "ids": ["7", "3", "9"]
}
```

Olası hata durumları:
- `400 Bad Request`: reorder aktif değil / geçersiz body / `ids` boş
- `403 Forbidden`: update yetkisi yok

### 5) Edit ↔ Detail Modal Geçişi

Kaynak index sayfasında modal akışı şu şekilde çalışır:

- Edit modal başlığında `Detaya Don` butonu bulunur.
- Detail modal başlığında `Duzenle` butonu bulunur.
- Bu geçişler route senkronizasyonu ile çalışır (`/show` ↔ `/edit`).

Bu davranış frontend'de varsayılan olarak aktiftir; ekstra resource ayarı gerektirmez.

## Varsayılan Sıralama

```go
r.SetSortable([]resource.Sortable{
    {Column: "created_at", Direction: "desc"},
})
```

## Record Title (İlişkiler İçin Kritik)

İlişkilerde görünen etiketleri kullanıcı dostu hale getirir.

```go
r.SetRecordTitleKey("name")

// veya
r.SetRecordTitleFunc(func(record any) string {
    // custom format
    return "..."
})
```

## İlişki Alanlarıyla Kullanım

```go
fields.BelongsTo("Category", "category_id", "categories").
    DisplayUsing("name").
    AutoOptions("name")

fields.HasMany("Prices", "prices", "prices").
    ForeignKey("product_id")
```

Detaylı ilişki API'si için:
- [İlişkiler](Relationships)
- [İlişkiler API Referansı](Relationships-API)

## Local Registry (Opsiyonel, Uygulama Seviyesi)

Büyük projelerde circular dependency yönetimi için local registry kullanabilirsiniz:

```go
// app/resources/registry.go
func Register(slug string, factory func() resource.Resource) {
    registry[slug] = factory
    resource.Register(slug, factory())
}

func GetOrPanic(slug string) resource.Resource {
    r := Get(slug)
    if r == nil {
        panic("resource not found: " + slug)
    }
    return r
}
```

Bu desen uygulamanıza aittir; Panel.go'nun çekirdek API'si değildir.

## Sık Hata Kontrolü

- Slug çakışması: Aynı slug son register edilenle ezilir.
- Resource görünmüyor: Package import edilmemiş olabilir, `init()` çalışmamıştır.
- İlişki dropdown boş: `AutoOptions` var ama ilişkili resource register edilmemiş olabilir.
- Relationship veri gelmiyor: FK/OwnerKey/Pivot ayarları modelle uyuşmuyor olabilir.

## Sonraki Adım

- Field detayları için: [Alanlar (Fields)](Fields)
- İlişki akışı için: [İlişkiler (Relationships)](Relationships)
- Tam başlangıç akışı için: [Başlarken](Getting-Started)
