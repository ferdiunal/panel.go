# Gelişmiş Kullanım (Advanced Usage)

Panel SDK'nın varsayılan davranışlarını ihtiyaçlarınıza göre özelleştirebilirsiniz.

## Custom Repository (Özel Veri Katmanı)

Varsayılan olarak her resource `GormDataProvider` kullanır. Ancak karmaşık sorgulara, farklı veri kaynaklarına veya özel iş mantığına ihtiyacınız varsa kendi repository'nizi kullanabilirsiniz.

### 1. Repository Oluşturma

`data.DataProvider` interface'ini implemente eden bir struct oluşturun.

```go
package repository

import (
    "context"
    "panel.go/internal/data"
)

type MyCustomRepo struct {
    // DB bağlantısı veya diğer servisler
}

func (r *MyCustomRepo) Index(ctx context.Context, req data.QueryRequest) (*data.QueryResponse, error) {
    // Özel listeleme mantığı (örn. ElasticSearch sorgusu, harici API çağrısı)
    return &data.QueryResponse{
        Items: []interface{}{}, 
        Total: 0,
    }, nil
}

// Diğer metodlar: Show, Create, Update, Delete...
func (r *MyCustomRepo) Show(ctx context.Context, id string) (interface{}, error) { return nil, nil }
func (r *MyCustomRepo) Create(ctx context.Context, data map[string]interface{}) (interface{}, error) { return nil, nil }
func (r *MyCustomRepo) Update(ctx context.Context, id string, data map[string]interface{}) (interface{}, error) { return nil, nil }
func (r *MyCustomRepo) Delete(ctx context.Context, id string) error { return nil }

// Opsiyonel: Arama ve Eager Loading yapılandırmalarını almak için
func (r *MyCustomRepo) SetSearchColumns(cols []string) {}
func (r *MyCustomRepo) SetWith(rels []string) {}
```

### 2. Resource'a Tanımlama

Resource struct'ınızda `Repository` metodunu override edin:

```go
func (u *UserResource) Repository(db *gorm.DB) data.DataProvider {
    // Kendi repository örneğinizi dönün
    return &repository.MyCustomRepo{}
    
    // Veya varsayılan GORM provider'ı wrap edebilirsiniz
    // return mywrapper.New(data.NewGormDataProvider(db, u.Model()))
}
```

## Custom Services (Özel Servisler)

`panel.New()` ile oluşturduğunuz uygulama instance'ı üzerinden altındaki Fiber uygulamasına erişebilirsiniz:

```go
app := panel.New(cfg)

app.Fiber.Get("/custom-health", func(c *fiber.Ctx) error {
    return c.SendString("OK")
})
```
