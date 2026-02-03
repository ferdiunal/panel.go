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
- **Embedded Frontend**: Frontend dosyaları binary içine gömülerek tek bir çalıştırılabilir dosya olarak dağıtılabilir.

## 📦 Kurulum

```bash
go get github.com/ferdiunal/panel.go
```

## ⚡ Hızlı Başlangıç

Sadece 4 adımda çalışır hale getirin.

### 1. Veritabanı Modeli (GORM)

```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    FullName  string    `json:"full_name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 2. Resource Tanımı

Modelinizi ve UI alanlarını (Fields) bağlayan yapıyı kurun.

```go
import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

type UserResource struct{
    resource.Base
}

// Resource Tanımlayıcı
func GetUserResource() resource.Resource {
    return &UserResource{
        Base: resource.Base{
            DataModel: &User{},
            Label:     "Users",
            FieldsVal: []fields.Element{
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
            },
        },
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
    "github.com/ferdiunal/panel.go/pkg/panel"
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
        Environment: "production", // "development" (embedded assetleri atlar) veya "production"
        Storage: panel.StorageConfig{
            Path: "./storage/public", // Disk üzerindeki yol
            URL:  "/storage",         // URL öneki
        },
    }

    // 3. Panel Oluştur
    app := panel.New(cfg)
    
    // Resource Kaydet
    app.RegisterResource(GetUserResource())

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

### Custom Repository Kullanımı

Varsayılan olarak her resource `GormDataProvider` kullanır. Ancak karmaşık sorgulara, farklı veri kaynaklarına veya özel iş mantığına ihtiyacınız varsa kendi repository'nizi kullanabilirsiniz.

1. `data.DataProvider` interface'ini implemente eden bir struct oluşturun.
2. Resource struct'ınızda `Repository` metodunu override ederek bu provider'ı dönün.

```go
// 1. Custom Repository Oluşturma
type MyCustomRepo struct {
    // ... gerekli alanlar
}

// data.DataProvider interface metodlarını implemente edin...

// 2. Resource İçinde Tanımlama
func (r *UserResource) Repository(db *gorm.DB) data.DataProvider {
    return &MyCustomRepo{}
}
```

## 📝 Lisans

MIT License.
