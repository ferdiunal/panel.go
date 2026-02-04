# Resource'lar Rehberi

Resource, veritabanınızdaki bir tabloyu yönetim panelinde göstermek ve yönetmek için kullanılan bir yapıdır.

## Resource Nedir?

Resource, şunları tanımlar:
- **Model** - Veritabanı tablosu
- **Alanlar** - Gösterilecek sütunlar
- **Politika** - Yetkilendirme kuralları
- **Repository** - Veritabanı işlemleri

## Basit Resource Oluşturma

```go
package resources

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"gorm.io/gorm"
)

// Veritabanı modeli
type Category struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Slug  string
	Count int
}

// Resource tanımı
type CategoryResource struct {
	resource.OptimizedBase
}

// Yeni resource oluştur
func NewCategoryResource() *CategoryResource {
	r := &CategoryResource{}

	// Temel bilgiler
	r.SetModel(&Category{})
	r.SetSlug("categories")
	r.SetTitle("Kategoriler")
	r.SetIcon("folder")
	r.SetGroup("İçerik")

	// Alanları tanımla
	r.SetFieldResolver(&CategoryFieldResolver{})

	// Politikayı tanımla
	r.SetPolicy(&CategoryPolicy{})

	return r
}

// Veritabanı işlemleri
func (r *CategoryResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &Category{})
}

// Alanları tanımla
type CategoryFieldResolver struct{}

func (r *CategoryFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "name",
			Name:  "Kategori Adı",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required().Searchable(),

		(&fields.Schema{
			Key:   "slug",
			Name:  "URL Slug",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnDetail().OnForm().Required(),

		(&fields.Schema{
			Key:   "count",
			Name:  "Yazı Sayısı",
			View:  "number",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),
	}
}

// Politika tanımla
type CategoryPolicy struct{}

func (p *CategoryPolicy) ViewAny(ctx *context.Context) bool {
	return true
}

func (p *CategoryPolicy) View(ctx *context.Context, model any) bool {
	return true
}

func (p *CategoryPolicy) Create(ctx *context.Context) bool {
	return isAdmin(ctx)
}

func (p *CategoryPolicy) Update(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}

func (p *CategoryPolicy) Delete(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}

func (p *CategoryPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

func (p *CategoryPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}

func isAdmin(ctx *context.Context) bool {
	user := ctx.User()
	return user != nil && user.Role == "admin"
}
```

## Resource Özellikleri

### Temel Bilgiler

```go
// Slug - URL'de kullanılır
r.SetSlug("products")

// Başlık - Panelde gösterilir
r.SetTitle("Ürünler")

// İkon - Menüde gösterilir
r.SetIcon("shopping-bag")

// Grup - Menü grubunda gösterilir
r.SetGroup("Satış")

// Navigasyon sırası
r.SetNavigationOrder(10)

// Görünürlük
r.SetVisible(true)
```

### Sıralama

Varsayılan sıralama ayarlarını belirleyin:

```go
func (r *ProductResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{
			Column:    "created_at",
			Direction: "desc", // En yeni önce
		},
	}
}
```

### İlişkiler

Diğer resource'larla ilişkili alanlar:

```go
(&fields.Schema{
	Key:   "category_id",
	Name:  "Kategori",
	View:  "belongs-to",
	Props: make(map[string]interface{}),
}).OnList().OnDetail().OnForm()
```

## Örnek: Tam E-ticaret Resource'ları

### Ürün Resource'u

```go
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	CategoryID  string
	Stock       int
	IsActive    bool
	CreatedAt   time.Time
}

type ProductResource struct {
	resource.OptimizedBase
}

func NewProductResource() *ProductResource {
	r := &ProductResource{}

	r.SetModel(&Product{})
	r.SetSlug("products")
	r.SetTitle("Ürünler")
	r.SetIcon("shopping-bag")
	r.SetGroup("Satış")

	r.SetFieldResolver(&ProductFieldResolver{})
	r.SetPolicy(&ProductPolicy{})

	return r
}

func (r *ProductResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &Product{})
}

func (r *ProductResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{Column: "created_at", Direction: "desc"},
	}
}

type ProductFieldResolver struct{}

func (r *ProductFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "name",
			Name:  "Ürün Adı",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required().Searchable().Sortable(),

		(&fields.Schema{
			Key:   "description",
			Name:  "Açıklama",
			View:  "textarea",
			Props: make(map[string]interface{}),
		}).OnForm(),

		(&fields.Schema{
			Key:   "price",
			Name:  "Fiyat",
			View:  "number",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required().Min(0).Sortable(),

		(&fields.Schema{
			Key:   "category_id",
			Name:  "Kategori",
			View:  "belongs-to",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required(),

		(&fields.Schema{
			Key:   "stock",
			Name:  "Stok",
			View:  "number",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required().Min(0),

		(&fields.Schema{
			Key:   "is_active",
			Name:  "Aktif mi?",
			View:  "switch",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "created_at",
			Name:  "Oluşturulma Tarihi",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),
	}
}

type ProductPolicy struct{}

func (p *ProductPolicy) ViewAny(ctx *context.Context) bool {
	return true
}

func (p *ProductPolicy) View(ctx *context.Context, model any) bool {
	return true
}

func (p *ProductPolicy) Create(ctx *context.Context) bool {
	return isAdmin(ctx)
}

func (p *ProductPolicy) Update(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}

func (p *ProductPolicy) Delete(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}

func (p *ProductPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

func (p *ProductPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}
```

### Sipariş Resource'u

```go
type Order struct {
	ID        string
	OrderNo   string
	UserID    string
	Total     float64
	Status    string // pending, processing, shipped, delivered
	CreatedAt time.Time
}

type OrderResource struct {
	resource.OptimizedBase
}

func NewOrderResource() *OrderResource {
	r := &OrderResource{}

	r.SetModel(&Order{})
	r.SetSlug("orders")
	r.SetTitle("Siparişler")
	r.SetIcon("package")
	r.SetGroup("Satış")

	r.SetFieldResolver(&OrderFieldResolver{})
	r.SetPolicy(&OrderPolicy{})

	return r
}

func (r *OrderResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &Order{})
}

type OrderFieldResolver struct{}

func (r *OrderFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "order_no",
			Name:  "Sipariş No",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().Searchable(),

		(&fields.Schema{
			Key:   "user_id",
			Name:  "Müşteri",
			View:  "belongs-to",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail(),

		(&fields.Schema{
			Key:   "total",
			Name:  "Toplam",
			View:  "number",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().ReadOnly(),

		(&fields.Schema{
			Key:   "status",
			Name:  "Durum",
			View:  "select",
			Props: map[string]interface{}{
				"options": map[string]string{
					"pending":    "Beklemede",
					"processing": "İşleniyor",
					"shipped":    "Gönderildi",
					"delivered":  "Teslim Edildi",
				},
			},
		}).OnList().OnDetail().OnForm().Filterable(),

		(&fields.Schema{
			Key:   "created_at",
			Name:  "Oluşturulma Tarihi",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),
	}
}

type OrderPolicy struct{}

func (p *OrderPolicy) ViewAny(ctx *context.Context) bool {
	return isAdmin(ctx) || isStaff(ctx)
}

func (p *OrderPolicy) View(ctx *context.Context, model any) bool {
	return isAdmin(ctx) || isStaff(ctx)
}

func (p *OrderPolicy) Create(ctx *context.Context) bool {
	return false // Siparişler sistem tarafından oluşturulur
}

func (p *OrderPolicy) Update(ctx *context.Context, model any) bool {
	return isAdmin(ctx) || isStaff(ctx)
}

func (p *OrderPolicy) Delete(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}

func (p *OrderPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

func (p *OrderPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}

func isStaff(ctx *context.Context) bool {
	user := ctx.User()
	return user != nil && (user.Role == "admin" || user.Role == "staff")
}
```

## Panel'e Resource'ları Kaydetme

```go
func main() {
	db, _ := gorm.Open(sqlite.Open("shop.db"), &gorm.Config{})
	db.AutoMigrate(&Product{}, &Category{}, &Order{})

	cfg := panel.Config{
		Database: panel.DatabaseConfig{
			Instance: db,
		},
		Server: panel.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "development",
		Resources: []resource.Resource{
			NewCategoryResource(),
			NewProductResource(),
			NewOrderResource(),
		},
	}

	app := panel.New(cfg)
	app.Start()
}
```

## İpuçları

1. **Slug**: URL-friendly olmalı (lowercase, tire ile ayrılmış)
2. **İkon**: Lucide icon set'ten seçin
3. **Grup**: İlişkili resource'ları aynı grupta toplayın
4. **Sıralama**: Sık kullanılan alanları sıralanabilir yapın
5. **Performans**: Büyük tablolar için pagination kullanın

## Sonraki Adımlar

- [Alanlar Rehberi](./Fields.md) - Alan tanımı
- [Politika Rehberi](./Authorization.md) - Yetkilendirme
- [İlişkiler Rehberi](./Relationships.md) - Tablo ilişkileri
