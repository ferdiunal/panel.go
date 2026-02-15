# Panel.go - Detaylı Teknik Başlangıç (Legacy Akış)

Bu rehber, Panel.go ile düşük seviye/legacy stile yakın bir başlangıç akışını adım adım kurar. Hedef, `model + resource + field resolver + policy + repository` zincirini tek bir çalışır senaryoda netleştirmektir.

## 1) Önkoşullar

- Go kurulumu (önerilen: güncel stabil sürüm)
- GORM ve seçtiğiniz veritabanı sürücüsü
- `panel` CLI (opsiyonel ama önerilen)
- Boş bir proje dizini

Veritabanı seçimi için kısa öneri:
- Geliştirme: `sqlite`
- Production: `postgres`
- Mevcut ekosistemle uyum gerekiyorsa: `mysql`

Not:
- Bu rehberde terminoloji tekilleştirilmiştir: `Panel.go framework`, `panel CLI`, `Resource`, `Policy`.

## 2) Kurulum

### 2.1 SDK ekleme

```bash
go get github.com/ferdiunal/panel.go
```

### 2.2 CLI kurulumu (opsiyonel ama önerilen)

```bash
go install github.com/ferdiunal/panel.go/cmd/panel@latest
```

## 3) `panel init` ile proje bootstrap

Yeni bir proje dizininde:

```bash
panel init
```

Veya veritabanını doğrudan seçerek:

```bash
panel init -d sqlite
panel init -d postgres
panel init -d mysql
```

Bu adım sonrası tipik olarak aşağıdaki dosyalar oluşur:
- `main.go`
- `go.mod`
- `.env`
- `.panel/stubs/*`

CLI detayları için: [CLI - Init Komutu](CLI_INIT)

## 4) Legacy Teknik Akış: Model + Resource + Resolver + Policy + Repository

Bu bölümde örnek bir `posts` kaynağı oluşturuyoruz.

### 4.1 Model (`internal/domain/post/entity.go`)

```go
package post

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"` // draft, published
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

### 4.2 Field Resolver (`internal/resource/post/field_resolver.go`)

```go
package post

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

type PostFieldResolver struct{}

func (r *PostFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
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
					"published": "Yayında",
				},
			},
		}).OnList().OnDetail().OnForm(),
	}
}
```

### 4.3 Policy (`internal/resource/post/policy.go`)

```go
package post

import appContext "github.com/ferdiunal/panel.go/pkg/context"

type PostPolicy struct{}

func (p *PostPolicy) ViewAny(ctx *appContext.Context) bool {
	return true
}

func (p *PostPolicy) View(ctx *appContext.Context, model any) bool {
	return true
}

func (p *PostPolicy) Create(ctx *appContext.Context) bool {
	return true
}

func (p *PostPolicy) Update(ctx *appContext.Context, model any) bool {
	return true
}

func (p *PostPolicy) Delete(ctx *appContext.Context, model any) bool {
	return true
}
```

### 4.4 Resource + Repository (`internal/resource/post/resource.go`)

```go
package post

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/resource"
	domainPost "my-panel-app/internal/domain/post"
	"gorm.io/gorm"
)

func init() {
	resource.Register("posts", NewPostResource())
}

type PostResource struct {
	resource.OptimizedBase
}

func NewPostResource() *PostResource {
	r := &PostResource{}

	r.SetModel(&domainPost.Post{})
	r.SetSlug("posts")
	r.SetTitle("Yazılar")
	r.SetIcon("file-text")
	r.SetGroup("İçerik")
	r.SetVisible(true)
	r.SetNavigationOrder(10)
	r.SetRecordTitleKey("title")

	r.SetFieldResolver(&PostFieldResolver{})
	r.SetPolicy(&PostPolicy{})

	return r
}

func (r *PostResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &domainPost.Post{})
}
```

Not:
- `resource.Register("posts", ...)` ilişkilerde `relatedResource` çözümlemeleri ve otomatik option yükleme için güvenli bir başlangıç desenidir.
- `SetRecordTitleKey("title")` ilişki dropdown'larında okunabilir etiket için kritiktir.

## 5) Resource Kaydı ve Panel Başlatma

`main.go` içinde paneli bağlayın ve resource ekleyin:

```go
package main

import (
	"log"

	"github.com/ferdiunal/panel.go/pkg/panel"
	domainPost "my-panel-app/internal/domain/post"
	postResource "my-panel-app/internal/resource/post"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("veritabanı bağlantı hatası: %v", err)
	}

	// Model migration
	_ = db.AutoMigrate(
		&domainPost.Post{},
	)

	cfg := panel.Config{
		Database: panel.DatabaseConfig{Instance: db},
		Server: panel.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "development",
	}

	app := panel.New(cfg)
	app.RegisterResource(postResource.NewPostResource())
	app.Start()
}
```

## 6) Otomatik Endpoint'leri Doğrulama

Resource slug'ınız `posts` ise aşağıdaki endpoint'ler otomatik gelir:

- `GET /api/resource/posts`
- `POST /api/resource/posts`
- `GET /api/resource/posts/:id`
- `PUT /api/resource/posts/:id`
- `DELETE /api/resource/posts/:id`

Hızlı doğrulama:

```bash
curl http://localhost:8080/api/resource/posts
```

## 7) Sık Hatalar ve Çözüm Tablosu

| Problem | Olası Neden | Çözüm |
|---|---|---|
| Resource görünmüyor | Package import edilmedi / register çalışmadı | `main.go` içinde ilgili package importunu ve `RegisterResource` çağrısını kontrol et |
| Dropdown ilişki boş | İlişkili resource register edilmedi | İlgili resource için `resource.Register(slug, ...)` ve doğru `relatedResource` kullan |
| 403 yetki hatası | Policy izinleri `false` dönüyor | Policy metotlarını ve permission tanımlarını kontrol et |
| Veri yazılmıyor | Repository/model map uyumsuzluğu | `SetModel`, GORM tag'leri ve field key'lerini eşleştir |
| Endpoint 404 | Yanlış slug | `SetSlug("posts")` ile URL slug'ının aynı olduğundan emin ol |

## 8) Sonraki Adımlar

- [Kaynaklar (Resource)](Resources)
- [Alanlar (Fields)](Fields)
- [İlişkiler (Relationships)](Relationships)
- [Yetkilendirme](Authorization)

## Hızlı Geçiş

Kısa akışa dönmek istersen:
- Repo ana giriş: [README](../README.md)
- Dokümantasyon merkezi: [Home](Home)
