# Panel.go 🚀

**Panel.go**, Go (Golang) projelerinizde hızlı, tip güvenli ve yönetilebilir admin panelleri oluşturmanız için tasarlanmış modern bir SDK'dır. 

Go'nun performansına ve tip güvenliğine uygun olarak tasarlanan bu yapı, veritabanı modellerinizi dakikalar içinde tam fonksiyonel bir REST API'ye ve yönetim arayüzüne dönüştürür.

## ✨ Özellikler

- **Resource Abstraction**: Model ve UI mantığını tek bir yapıda toplayın.
- **Fluent Field API**: Zincirleme metodlarla (`Text("Ad").Sortable().Required()`) kolayca alan tanımlayın.
- **Otomatik CRUD**: Oluşturduğunuz her resource için Create, Read, Update, Delete ve Show endpointleri hazır gelir.
- **Smart Data Provider**: GORM entegrasyonu ile sayfalama, sıralama ve filtreleme otomatik halledilir.
- **Central App Config**: Tek bir `Panel` instance'ı ile tüm servisi yönetin.
- **Genişletilebilir Mimari**: Kendi özel servislerinizi ve rotalarınızı kolayca entegre edin.
- **Custom Data Providers**: Veri erişim katmanını tamamen özelleştirebilme (Custom Repository) yeteneği.

## 📦 Kurulum

```bash
go get panel.go
```

## ⚡ Hızlı Başlangıç

Sadece 4 adımda çalışır hale getirin.

### 1. Veritabanı Modeli (GORM)

```go
type User struct {
    ID        uint   `json:"id" gorm:"primaryKey"`
    FullName  string `json:"full_name"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 2. Resource Tanımı

Modelinizi ve UI alanlarını (Fields) bağlayan yapıyı kurun.

```go
import (
    "panel.go/internal/fields"
    "panel.go/internal/resource"
)

type UserResource struct{}

// Hangi model ile çalışacağını belirtin
func (u *UserResource) Model() interface{} {
    return &User{}
}

// Hangi alanların görüneceğini belirtin
func (u *UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID().Sortable(),

        fields.Text("Ad Soyad", "full_name").
            Sortable().
            Placeholder("Tam ad...").
            Required(),

        fields.Email("E-Posta", "email").
            Sortable().
            Required(),

        fields.Select("Rol", "role").
            Options(map[string]string{
                "admin": "Yönetici",
                "user":  "Kullanıcı",
            }),
            
        fields.DateTime("Kayıt Tarihi", "created_at").
            OnList().
            ReadOnly(),
    }
}
```

### 3. Uygulamayı Başlatma

`main.go` dosyanızda paneli yapılandırın ve resource'ları kaydedin.

```go
package main

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "panel.go/internal/panel"
)

func main() {
    // 1. Veritabanı Bağlantısı
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    db.AutoMigrate(&User{})

    // 2. Panel Ayarları
    cfg := panel.Config{
        Server: panel.ServerConfig{
            Host: "localhost",
            Port: "8080",
        },
        Database: panel.DatabaseConfig{
            Instance: db,
        },
    }

    // 3. Panel Oluştur ve Resource Kaydet
    app := panel.New(cfg)
    
    // "/api/resource/users" rotasını otomatik oluşturur
    app.Register("users", &UserResource{}) 

    // 4. Sunucuyu Başlat
    app.Start()
}
```

## 🔌 API Endpoints

Resource kaydedildikten sonra (örneğin `"users"` slug'ı ile), aşağıdaki endpointler otomatik olarak aktif olur:

| Metot | Endpoint | Açıklama |
|-------|----------|----------|
| `GET` | `/api/resource/users` | Listeleme (Sayfalama, Sıralama, Arama destekli) |
| `POST` | `/api/resource/users` | Yeni kayıt oluşturma |
| `GET` | `/api/resource/users/:id` | Tekil kayıt detayını görüntüleme |
| `PUT` | `/api/resource/users/:id` | Kayıt güncelleme |
| `DELETE` | `/api/resource/users/:id` | Kayıt silme |

## 🛠 Gelişmiş Kullanım


### Mevcut Uygulamayı Genişletme (Custom Services)

Panel.go, sadece admin paneli için değil, uygulamanızın tamamı için bir çatı görevi görebilir. `app.Fiber` nesnesine erişerek kendi özel route'larınızı ve servislerinizi ekleyebilirsiniz.

```go
func main() {
    // ... app kurulumu ...
    app := panel.New(cfg)

    // 1. Resource Kaydı
    app.Register("users", &UserResource{})

    // 2. Özel Servis/Route Ekleme
    // Fiber app instance'ına direkt erişiminiz vardır.
    
    // Basit bir endpoint
    app.Fiber.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })

    // Group kullanımı
    v1 := app.Fiber.Group("/api/v1")
    v1.Post("/login", authHandler.Login)
    v1.Post("/register", authHandler.Register)

    // 3. Sunucuyu Başlat
    app.Start()
}
```

### Custom Repository Kullanımı

Varsayılan olarak her resource `GormDataProvider` kullanır. Ancak karmaşık sorgulara, farklı veri kaynaklarına veya özel iş mantığına ihtiyacınız varsa kendi repository'nizi kullanabilirsiniz.

1. `data.DataProvider` interface'ini implemente eden bir struct oluşturun.
2. Resource struct'ınızda `Repository` metodunu override ederek bu provider'ı dönün.

```go
// 1. Custom Repository Oluşturma
type MyCustomRepo struct {
    // ... gerekli alanlar
}

// data.DataProvider interface metodlarını implemente edin
func (r *MyCustomRepo) Index(ctx context.Context, req data.QueryRequest) (*data.QueryResponse, error) {
    // Özel listeleme mantığı
    return &data.QueryResponse{Items: []interface{}{}, Total: 0}, nil
}
func (r *MyCustomRepo) Show(ctx context.Context, id string) (interface{}, error) { return nil, nil }
func (r *MyCustomRepo) Create(ctx context.Context, data map[string]interface{}) (interface{}, error) { return nil, nil }
func (r *MyCustomRepo) Update(ctx context.Context, id string, data map[string]interface{}) (interface{}, error) { return nil, nil }
func (r *MyCustomRepo) Delete(ctx context.Context, id string) error { return nil }
func (r *MyCustomRepo) SetSearchColumns(cols []string) {}
func (r *MyCustomRepo) SetWith(rels []string) {}

// 2. Resource İçinde Tanımlama
func (u *UserResource) Repository(db *gorm.DB) data.DataProvider {
    return &MyCustomRepo{}
    // Veya varsayılan GORM provider'ı özelleştirerek dönebilirsiniz:
    // provider := data.NewGormDataProvider(db, u.Model())
    // return provider
}
```

## 📝 Lisans

MIT License.
