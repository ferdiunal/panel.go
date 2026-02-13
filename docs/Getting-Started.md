# Panel.go - Başlangıç Rehberi

Panel.go, Go uygulamalarınız için hızlı ve güçlü bir yönetim paneli oluşturmanıza yardımcı olan bir framework'tür.

## Kurulum

### 1. Projenize Ekleyin

```bash
go get github.com/ferdiunal/panel.go
```

### 2. Temel Kurulum

```go
package main

import (
	"log"
	"github.com/ferdiunal/panel.go/pkg/panel"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Veritabanı bağlantısı
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		panic("Veritabanı bağlantısı başarısız")
	}

	// Panel yapılandırması
	cfg := panel.Config{
		Database: panel.DatabaseConfig{
			Instance: db,
		},
		Server: panel.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "development",
	}

	// Panel oluştur ve başlat
	app := panel.New(cfg)
	log.Println("Panel http://localhost:8080 adresinde çalışıyor")
	app.Start()
}
```

## Temel Kavramlar

### Resource Nedir?

Resource, veritabanınızdaki bir tabloyu yönetmek için kullanılan bir yapıdır. Örneğin, "Kullanıcılar" tablosunu yönetmek için bir User Resource oluşturursunuz.

Her resource şunları içerir:
- **Alanlar** (Fields) - Tablonun sütunları
- **Politika** (Policy) - Kimin ne yapabileceği
- **Repository** - Veritabanı işlemleri

### Basit Bir Resource Oluşturma

Diyelim ki bir Blog uygulaması yapıyorsunuz ve "Yazılar" tablosunu yönetmek istiyorsunuz.

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

// Post modeli
type Post struct {
	ID    string `gorm:"primaryKey"`
	Title string
	Body  string
	Status string // draft, published
}

// PostResource, yazıları yönetmek için resource
type PostResource struct {
	resource.OptimizedBase
}

// NewPostResource, yeni bir Post resource'u oluşturur
func NewPostResource() *PostResource {
	r := &PostResource{}

	r.SetModel(&Post{})
	r.SetSlug("posts")
	r.SetTitle("Yazılar")
	r.SetIcon("file-text")
	r.SetGroup("İçerik")

	// Alanları tanımla
	r.SetFieldResolver(&PostFieldResolver{})

	// Politikayı tanımla
	r.SetPolicy(&PostPolicy{})

	return r
}

// Repository, veritabanı işlemleri
func (r *PostResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &Post{})
}

// PostFieldResolver, yazı alanlarını tanımlar
type PostFieldResolver struct{}

func (r *PostFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "title",
			Name:  "Başlık",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required(),

		(&fields.Schema{
			Key:   "body",
			Name:  "İçerik",
			View:  "textarea",
			Props: make(map[string]interface{}),
		}).OnDetail().OnForm().Required(),

		(&fields.Schema{
			Key:   "status",
			Name:  "Durum",
			View:  "select",
			Props: map[string]interface{}{
				"options": map[string]string{
					"draft":     "Taslak",
					"published": "Yayınlandı",
				},
			},
		}).OnList().OnDetail().OnForm(),
	}
}

// PostPolicy, yazılar için yetkilendirme
type PostPolicy struct{}

func (p *PostPolicy) ViewAny(ctx *context.Context) bool {
	return true // Herkes yazıları görebilir
}

func (p *PostPolicy) View(ctx *context.Context, model any) bool {
	return true
}

func (p *PostPolicy) Create(ctx *context.Context) bool {
	return true // Herkes yazı oluşturabilir
}

func (p *PostPolicy) Update(ctx *context.Context, model any) bool {
	return true
}

func (p *PostPolicy) Delete(ctx *context.Context, model any) bool {
	return true
}

func (p *PostPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

func (p *PostPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}
```

### Resource'u Panel'e Kaydetme

```go
func main() {
	db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

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
			resources.NewPostResource(),
		},
	}

	app := panel.New(cfg)
	app.Start()
}
```

## Alanlar (Fields)

Alanlar, veritabanı sütunlarını yönetim panelinde nasıl göstereceğinizi tanımlar.

### Desteklenen Alan Türleri

```go
// Metin alanı
(&fields.Schema{
	Key:   "name",
	Name:  "Ad",
	View:  "text",
	Props: make(map[string]interface{}),
}).OnList().OnDetail().OnForm()

// Metin alanı (çok satırlı)
(&fields.Schema{
	Key:   "description",
	Name:  "Açıklama",
	View:  "textarea",
	Props: make(map[string]interface{}),
}).OnForm()

// Seçim alanı
(&fields.Schema{
	Key:   "status",
	Name:  "Durum",
	View:  "select",
	Props: map[string]interface{}{
		"options": map[string]string{
			"active":   "Aktif",
			"inactive": "İnaktif",
		},
	},
}).OnList().OnDetail().OnForm()

// Tarih alanı
(&fields.Schema{
	Key:   "created_at",
	Name:  "Oluşturulma Tarihi",
	View:  "datetime",
	Props: make(map[string]interface{}),
}).ReadOnly().OnList().OnDetail()

// Evet/Hayır alanı
(&fields.Schema{
	Key:   "is_active",
	Name:  "Aktif mi?",
	View:  "switch",
	Props: make(map[string]interface{}),
}).OnList().OnDetail().OnForm()
```

### Alan Seçenekleri

```go
// Zorunlu alan
field.Required()

// Salt okunur alan
field.ReadOnly()

// Varsayılan değer
field.Default("Varsayılan Değer")

// Yardım metni
field.HelpText("Bu alan hakkında bilgi")

// Yer tutucu metni
field.Placeholder("Buraya yazın...")

// Sadece listede göster
field.OnlyOnList()

// Sadece detayda göster
field.OnlyOnDetail()

// Sadece formda göster
field.OnlyOnForm()

// Listede gizle
field.HideOnList()

// Detayda gizle
field.HideOnDetail()

// Oluşturmada gizle
field.HideOnCreate()

// Güncellemede gizle
field.HideOnUpdate()
```

## Politika (Policy)

Politika, kullanıcıların hangi işlemleri yapabileceğini kontrol eder.

```go
type PostPolicy struct{}

// Tüm yazıları görme izni
func (p *PostPolicy) ViewAny(ctx *context.Context) bool {
	// Örnek: Sadece admin'ler görebilir
	return isAdmin(ctx)
}

// Belirli bir yazıyı görme izni
func (p *PostPolicy) View(ctx *context.Context, model any) bool {
	post := model.(*Post)
	// Örnek: Yazar veya admin görebilir
	return isAuthor(ctx, post) || isAdmin(ctx)
}

// Yazı oluşturma izni
func (p *PostPolicy) Create(ctx *context.Context) bool {
	return isLoggedIn(ctx)
}

// Yazı güncelleme izni
func (p *PostPolicy) Update(ctx *context.Context, model any) bool {
	post := model.(*Post)
	return isAuthor(ctx, post) || isAdmin(ctx)
}

// Yazı silme izni
func (p *PostPolicy) Delete(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}

// Yazı geri yükleme izni
func (p *PostPolicy) Restore(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}

// Yazı kalıcı silme izni
func (p *PostPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return isAdmin(ctx)
}
```

## Gelişmiş Özellikler

### Sıralama

Resource'ları varsayılan olarak nasıl sıralanacağını belirleyin:

```go
func (r *PostResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{
			Column:    "created_at",
			Direction: "desc", // En yeni yazılar önce
		},
	}
}
```

### Arama

Alanları aranabilir yapın:

```go
(&fields.Schema{
	Key:   "title",
	Name:  "Başlık",
	View:  "text",
	Props: make(map[string]interface{}),
}).Searchable()
```

### Sıralama

Alanları sıralanabilir yapın:

```go
(&fields.Schema{
	Key:   "created_at",
	Name:  "Oluşturulma Tarihi",
	View:  "datetime",
	Props: make(map[string]interface{}),
}).Sortable()
```

### Filtreleme

Alanları filtrelenebilir yapın:

```go
(&fields.Schema{
	Key:   "status",
	Name:  "Durum",
	View:  "select",
	Props: map[string]interface{}{
		"options": map[string]string{
			"draft":     "Taslak",
			"published": "Yayınlandı",
		},
	},
}).Filterable()
```

## Örnek: Tam Blog Uygulaması

```go
package main

import (
	"log"
	"github.com/ferdiunal/panel.go/pkg/panel"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Veritabanı
	db, _ := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	db.AutoMigrate(&Post{}, &Category{})

	// Panel yapılandırması
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
			NewPostResource(),
			NewCategoryResource(),
		},
	}

	// Panel başlat
	app := panel.New(cfg)
	log.Println("Blog paneli http://localhost:8080 adresinde çalışıyor")
	app.Start()
}
```

## Sonraki Adımlar

- [Alanlar Rehberi](./Fields.md) - Tüm alan türleri hakkında detaylı bilgi
- [Politika Rehberi](./Authorization.md) - Yetkilendirme sistemi
- [İlişkiler Rehberi](./Relationships.md) - Tablo ilişkileri
- [API Referansı](./API-Reference.md) - Tüm API metodları
