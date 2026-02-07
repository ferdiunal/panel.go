# Action'lar Rehberi

Action'lar, resource'lar üzerinde toplu işlemler gerçekleştirmek için kullanılan güçlü bir mekanizmadır. Seçili kayıtlar üzerinde onaylama, silme, dışa aktarma gibi işlemleri tek seferde yapabilirsiniz.

## Action Nedir?

Action, şunları sağlar:
- **Toplu İşlemler** - Birden fazla kayıt üzerinde aynı anda işlem yapma
- **Özel Alanlar** - Action'a özel form alanları tanımlama
- **Yetkilendirme** - Action'ın çalışıp çalışmayacağını kontrol etme
- **Async Execution** - Paralel model yükleme ile yüksek performans

## Basit Action Oluşturma

```go
package actions

import (
	"github.com/ferdiunal/panel.go/pkg/action"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

type ApproveAction struct {
	action.Base
}

func NewApproveAction() *ApproveAction {
	a := &ApproveAction{}

	// Temel bilgiler
	a.SetName("Onayla")
	a.SetSlug("approve")
	a.SetIcon("check-circle")

	// Onay mesajı
	a.SetConfirmText("Seçili kayıtları onaylamak istediğinize emin misiniz?")
	a.SetConfirmButtonText("Onayla")
	a.SetCancelButtonText("İptal")

	return a
}

// Execute - Action'ın çalıştırılacağı fonksiyon
func (a *ApproveAction) Execute(c *context.Context, models []interface{}) error {
	db := c.Locals("db").(*gorm.DB)

	for _, model := range models {
		if post, ok := model.(*Post); ok {
			post.Status = "approved"
			post.ApprovedAt = time.Now()

			if err := db.Save(post).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
```

## Action'a Alan Ekleme

Action'lar, kullanıcıdan ek bilgi almak için form alanları içerebilir:

```go
func NewPublishAction() *PublishAction {
	a := &PublishAction{}

	a.SetName("Yayınla")
	a.SetSlug("publish")

	// Action'a özel alanlar
	a.SetFields([]fields.Element{
		fields.DateTime("published_at").
			Label("Yayın Tarihi").
			Default(time.Now()).
			Required(),

		fields.Select("category").
			Label("Kategori").
			Options(map[string]string{
				"news":    "Haberler",
				"blog":    "Blog",
				"article": "Makale",
			}).
			Required(),

		fields.Textarea("note").
			Label("Not").
			Placeholder("Yayın notu ekleyin..."),
	})

	return a
}

func (a *PublishAction) Execute(c *context.Context, models []interface{}) error {
	// Action alanlarına erişim
	fields := c.Locals("action_fields").(map[string]interface{})

	publishedAt := fields["published_at"].(time.Time)
	category := fields["category"].(string)
	note := fields["note"].(string)

	db := c.Locals("db").(*gorm.DB)

	for _, model := range models {
		if post, ok := model.(*Post); ok {
			post.Status = "published"
			post.PublishedAt = publishedAt
			post.Category = category
			post.Note = note

			db.Save(post)
		}
	}

	return nil
}
```

## Action Yetkilendirme

`CanRun` metodu ile action'ın çalışıp çalışmayacağını kontrol edebilirsiniz:

```go
func (a *DeleteAction) CanRun(ctx *action.ActionContext) bool {
	// Sadece admin kullanıcılar silebilir
	user := ctx.User.(*User)
	if user.Role != "admin" {
		return false
	}

	// Yayınlanmış kayıtlar silinemez
	for _, model := range ctx.Models {
		if post, ok := model.(*Post); ok {
			if post.Status == "published" {
				return false
			}
		}
	}

	return true
}
```

## Action Görünürlük Ayarları

Action'ların nerede görüneceğini kontrol edebilirsiniz:

```go
func NewExportAction() *ExportAction {
	a := &ExportAction{}

	a.SetName("Dışa Aktar")
	a.SetSlug("export")

	// Sadece liste sayfasında göster
	a.SetOnlyOnIndex(true)

	// Detay sayfasında gösterme
	a.SetOnlyOnDetail(false)

	// Inline olarak göster (dropdown yerine buton)
	a.SetShowInline(true)

	return a
}
```

## Destructive Action

Tehlikeli işlemler için `destructive` flag'i kullanın:

```go
func NewDeleteAction() *DeleteAction {
	a := &DeleteAction{}

	a.SetName("Sil")
	a.SetSlug("delete")
	a.SetIcon("trash-2")

	// Tehlikeli işlem olarak işaretle (kırmızı renk)
	a.SetDestructive(true)

	a.SetConfirmText("Bu kayıtları kalıcı olarak silmek istediğinize emin misiniz?")
	a.SetConfirmButtonText("Evet, Sil")

	return a
}
```

## Resource'a Action Ekleme

Action'ları resource'a eklemek için `GetActions` metodunu kullanın:

```go
type PostResource struct {
	resource.OptimizedBase
}

func (r *PostResource) GetActions() []interface{} {
	return []interface{}{
		NewApproveAction(),
		NewPublishAction(),
		NewExportAction(),
		NewDeleteAction(),
	}
}
```

## Async Pattern: Paralel Model Loading

**v1.2.0+** Panel.go, action execution sırasında modelleri paralel olarak yükler. Bu, yüzlerce kayıt üzerinde işlem yaparken önemli performans artışı sağlar.

### Nasıl Çalışır?

```go
// Eski Yöntem (Sequential) - 100 model için ~5-10 saniye
for _, id := range body.IDs {
    model := reflect.New(modelType).Interface()
    db.First(model, "id = ?", id)
    models = append(models, model)
}

// Yeni Yöntem (Parallel) - 100 model için <1 saniye
// Fan-out: Her model için goroutine başlat
for _, id := range body.IDs {
    go func(id string) {
        defer wg.Done()
        model := reflect.New(modelType).Interface()
        err := db.First(model, "id = ?", id).Error
        results <- modelResult{model: model, err: err, id: id}
    }(id)
}

// Fan-in: Sonuçları topla
for result := range results {
    if result.err == nil {
        models = append(models, result.model)
    }
}
```

### Performans Karşılaştırması

| Model Sayısı | Sequential | Parallel | Hız Artışı |
|--------------|-----------|----------|------------|
| 10           | ~500ms    | ~50ms    | 10x        |
| 50           | ~2.5s     | ~250ms   | 10x        |
| 100          | ~5s       | ~500ms   | 10x        |
| 500          | ~25s      | ~2.5s    | 10x        |

### Teknik Detaylar

**Goroutine Pool:**
- Her model için 1 goroutine
- Buffered channel (size = model count)
- sync.WaitGroup ile completion tracking

**Error Handling:**
- First error wins stratejisi
- Diğer goroutine'ler devam eder
- Partial success desteklenir

**Dosya:** `pkg/handler/action_handler.go:104-160`

## Action Context

Action execution sırasında kullanılabilecek context bilgileri:

```go
type ActionContext struct {
	Models   []interface{}          // Seçili modeller
	Fields   map[string]interface{} // Action form alanları
	User     interface{}            // Mevcut kullanıcı
	Resource string                 // Resource slug
	DB       *gorm.DB              // Database instance
	Ctx      *fasthttp.RequestCtx  // HTTP context
}
```

## Gelişmiş Örnekler

### Toplu E-posta Gönderme

```go
type SendEmailAction struct {
	action.Base
}

func NewSendEmailAction() *SendEmailAction {
	a := &SendEmailAction{}

	a.SetName("E-posta Gönder")
	a.SetSlug("send-email")
	a.SetIcon("mail")

	a.SetFields([]fields.Element{
		fields.Text("subject").
			Label("Konu").
			Required(),

		fields.Textarea("message").
			Label("Mesaj").
			Required(),

		fields.Switch("send_immediately").
			Label("Hemen Gönder").
			Default(true),
	})

	return a
}

func (a *SendEmailAction) Execute(c *context.Context, models []interface{}) error {
	fields := c.Locals("action_fields").(map[string]interface{})

	subject := fields["subject"].(string)
	message := fields["message"].(string)
	sendImmediately := fields["send_immediately"].(bool)

	// Paralel e-posta gönderimi
	var wg sync.WaitGroup
	for _, model := range models {
		if user, ok := model.(*User); ok {
			wg.Add(1)
			go func(u *User) {
				defer wg.Done()

				if sendImmediately {
					sendEmail(u.Email, subject, message)
				} else {
					queueEmail(u.Email, subject, message)
				}
			}(user)
		}
	}

	wg.Wait()
	return nil
}
```

### CSV Dışa Aktarma

```go
type ExportCSVAction struct {
	action.Base
}

func (a *ExportCSVAction) Execute(c *context.Context, models []interface{}) error {
	// CSV buffer oluştur
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header yaz
	writer.Write([]string{"ID", "Name", "Email", "Created At"})

	// Verileri yaz
	for _, model := range models {
		if user, ok := model.(*User); ok {
			writer.Write([]string{
				user.ID,
				user.Name,
				user.Email,
				user.CreatedAt.Format("2006-01-02"),
			})
		}
	}

	writer.Flush()

	// Dosya olarak indir
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=export.csv")
	return c.Send(buf.Bytes())
}
```

### Durum Değiştirme (State Machine)

```go
type ChangeStatusAction struct {
	action.Base
}

func NewChangeStatusAction() *ChangeStatusAction {
	a := &ChangeStatusAction{}

	a.SetName("Durum Değiştir")
	a.SetSlug("change-status")

	a.SetFields([]fields.Element{
		fields.Select("status").
			Label("Yeni Durum").
			Options(map[string]string{
				"draft":     "Taslak",
				"review":    "İncelemede",
				"approved":  "Onaylandı",
				"published": "Yayınlandı",
				"archived":  "Arşivlendi",
			}).
			Required(),

		fields.Textarea("reason").
			Label("Değişiklik Nedeni").
			Placeholder("Durum değişikliği nedenini açıklayın..."),
	})

	return a
}

func (a *ChangeStatusAction) Execute(c *context.Context, models []interface{}) error {
	fields := c.Locals("action_fields").(map[string]interface{})
	db := c.Locals("db").(*gorm.DB)

	newStatus := fields["status"].(string)
	reason := fields["reason"].(string)

	for _, model := range models {
		if post, ok := model.(*Post); ok {
			oldStatus := post.Status

			// Durum geçişini kontrol et
			if !isValidTransition(oldStatus, newStatus) {
				return fmt.Errorf("geçersiz durum geçişi: %s -> %s", oldStatus, newStatus)
			}

			// Durumu güncelle
			post.Status = newStatus

			// Audit log oluştur
			db.Create(&StatusLog{
				PostID:    post.ID,
				OldStatus: oldStatus,
				NewStatus: newStatus,
				Reason:    reason,
				UserID:    c.Locals("user_id").(string),
			})

			db.Save(post)
		}
	}

	return nil
}

func isValidTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"draft":     {"review", "archived"},
		"review":    {"draft", "approved", "archived"},
		"approved":  {"published", "review"},
		"published": {"archived"},
		"archived":  {"draft"},
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}
```

## API Endpoints

### Action Listesi

```http
GET /api/resource/:resource/actions
```

**Response:**
```json
{
  "actions": [
    {
      "name": "Onayla",
      "slug": "approve",
      "icon": "check-circle",
      "confirmText": "Seçili kayıtları onaylamak istediğinize emin misiniz?",
      "confirmButtonText": "Onayla",
      "cancelButtonText": "İptal",
      "destructive": false,
      "onlyOnIndex": false,
      "onlyOnDetail": false,
      "showInline": false,
      "fields": []
    }
  ]
}
```

### Action Çalıştırma

```http
POST /api/resource/:resource/actions/:action
Content-Type: application/json

{
  "ids": ["1", "2", "3"],
  "fields": {
    "published_at": "2024-01-15T10:00:00Z",
    "category": "news"
  }
}
```

**Response:**
```json
{
  "message": "Action executed successfully on 3 item(s)",
  "count": 3
}
```

## Best Practices

### 1. Performans

```go
// ✅ İyi: Paralel işlem
var wg sync.WaitGroup
for _, model := range models {
	wg.Add(1)
	go func(m interface{}) {
		defer wg.Done()
		processModel(m)
	}(model)
}
wg.Wait()

// ❌ Kötü: Sequential işlem
for _, model := range models {
	processModel(model)
}
```

### 2. Error Handling

```go
// ✅ İyi: Detaylı hata mesajı
if err := db.Save(post).Error; err != nil {
	return fmt.Errorf("post kaydedilemedi (ID: %s): %w", post.ID, err)
}

// ❌ Kötü: Generic hata
if err := db.Save(post).Error; err != nil {
	return err
}
```

### 3. Transaction Kullanımı

```go
// ✅ İyi: Transaction ile güvenli işlem
func (a *TransferAction) Execute(c *context.Context, models []interface{}) error {
	db := c.Locals("db").(*gorm.DB)

	return db.Transaction(func(tx *gorm.DB) error {
		for _, model := range models {
			if err := processModel(tx, model); err != nil {
				return err // Rollback
			}
		}
		return nil // Commit
	})
}
```

### 4. Validation

```go
// ✅ İyi: Erken validation
func (a *PublishAction) Execute(c *context.Context, models []interface{}) error {
	// Önce tüm modelleri validate et
	for _, model := range models {
		if post, ok := model.(*Post); ok {
			if post.Title == "" {
				return fmt.Errorf("post başlığı boş olamaz (ID: %s)", post.ID)
			}
		}
	}

	// Sonra işle
	for _, model := range models {
		// ...
	}

	return nil
}
```

## Troubleshooting

### Action Görünmüyor

1. `GetActions()` metodunun doğru implement edildiğinden emin olun
2. Action'ın `OnlyOnIndex` veya `OnlyOnDetail` ayarlarını kontrol edin
3. `CanRun()` metodunun `true` döndüğünden emin olun

### Performans Sorunları

1. Paralel model loading kullanıldığından emin olun (v1.2.0+)
2. Database query'lerini optimize edin (N+1 problem)
3. Büyük işlemler için background job kullanın

### Memory Leak

1. Goroutine'lerin düzgün kapatıldığından emin olun
2. Channel'ları close edin
3. Context timeout kullanın

## İleri Seviye Konular

- **Background Jobs**: Uzun süren action'lar için job queue kullanımı
- **Event Broadcasting**: Action sonrası event yayınlama
- **Audit Logging**: Action'ların loglanması
- **Rate Limiting**: Action çalıştırma sınırlaması
- **Webhook Integration**: Action sonrası webhook tetikleme

Bu konular için [Advanced-Usage.md](Advanced-Usage.md) dökümanına bakın.
