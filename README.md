# Panel.go 🚀

**Panel.go**, Go (Golang) projelerinizde hızlı, tip güvenli ve yönetilebilir admin panelleri oluşturmanız için tasarlanmış modern bir SDK'dır.

Go'nun performansına ve tip güvenliğine uygun olarak tasarlanan bu yapı, veritabanı modellerinizi dakikalar içinde tam fonksiyonel bir REST API'ye ve yönetim arayüzüne dönüştürür.

## ✨ Özellikler

- **Resource Abstraction**: Model ve UI mantığını tek bir yapıda toplayın.
- **Fluent Field API**: Zincirleme metodlarla (`Text("Ad").Sortable().Required()`) kolayca alan tanımlayın.
- **Relationship Fields**: BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo ilişkilerini destekler.
- **Otomatik CRUD**: Oluşturduğunuz her resource için Create, Read, Update, Delete ve Show endpointleri hazır gelir.
- **Smart Data Provider**: GORM entegrasyonu ile sayfalama, sıralama ve filtreleme otomatik halledilir.
- **Central App Config**: Tek bir `Panel` instance'ı ile tüm servisi yönetin.
- **Genişletilebilir Mimari**: Kendi özel servislerinizi ve rotalarınızı kolayca entegre edin.
- **Embedded Frontend**: Frontend dosyaları binary içine gömülerek tek bir çalıştırılabilir dosya olarak dağıtılabilir.
- **Kapsamlı Dokümantasyon**: Türkçe yazılmış, 70+ örnek içeren detaylı rehberler.

## � Dokümantasyon

Panel.go için kapsamlı, Türkçe yazılmış dokümantasyon mevcuttur. Tüm rehberlere `docs/` klasöründen erişebilirsiniz.

### Başlarken
- **[Başlarken](docs/Getting-Started.md)** - Kurulum ve ilk resource oluşturma
- **[Kaynaklar (Resources)](docs/Resources.md)** - Resource tanımı ve yapılandırması
- **[Alanlar (Fields)](docs/Fields.md)** - 10+ alan türü ve seçenekleri

### Temel Kavramlar
- **[İlişkiler (Relationships)](docs/Relationships.md)** - BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo
- **[Yetkilendirme (Authorization)](docs/Authorization.md)** - Policy yazma ve rol tabanlı erişim kontrolü

### İleri Seviye
- **[Gelişmiş Kullanım (Advanced Usage)](docs/Advanced-Usage.md)** - Özel alanlar, middleware, hooks, optimizasyon
- **[API Referansı (API Reference)](docs/API-Reference.md)** - Tüm metodlar ve parametreler
- **[Lensler (Lenses)](docs/Lenses.md)** - Özel raporlar ve görünümler
- **[Sayfalar (Pages)](docs/Pages.md)** - Özel gösterge panelleri
- **[Ayarlar (Settings)](docs/Settings.md)** - Uygulama ayarları
- **[Widgets](docs/Widgets.md)** - Gösterge paneli widget'ları

### Diğer
- **[Kimlik Doğrulama (Authentication)](docs/Authentication.md)** - Kullanıcı kimlik doğrulaması

**Toplam:** 2000+ satır, 70+ gerçek dünya örneği

## 📊 Proje Durumu

```
✅ 453 Test (tümü geçiyor)
✅ 0 Derleme Hatası
✅ 0 Lint Hatası
✅ Kapsamlı Türkçe Dokümantasyon
✅ Üretim Hazır
```

| Metrik | Değer |
|--------|-------|
| Test Coverage | 453/453 (%100) |
| Compilation | ✅ 0 errors |
| Linting | ✅ 0 errors |
| Documentation | ✅ 14 files |
| Examples | ✅ 70+ |
| Status | ✅ Production Ready |

## 📦 Kurulum

```bash
go get github.com/ferdiunal/panel.go
```

### UI Dosyaları

Panel.go, frontend dosyalarını Go binary'sine gömer (embed). Projeyi klonladığınızda UI dosyaları zaten `pkg/panel/ui/` klasöründe hazır olarak gelir, bu yüzden Node.js veya Bun kurmanıza gerek yoktur.

#### Frontend'i Yeniden Build Etme (Opsiyonel)

Eğer frontend kodunda değişiklik yaparsanız, UI'ı yeniden build etmek için:

```bash
# Önce web bağımlılıklarını yükleyin (sadece ilk seferde)
cd web && bun install

# UI'ı build edin ve pkg/panel/ui'a kopyalayın
make build-ui
```

Bu komut:
1. `web/` klasöründeki React uygulamasını build eder
2. Build edilen dosyaları `pkg/panel/ui/` klasörüne kopyalar
3. Bir sonraki Go build'de bu dosyalar otomatik olarak binary'e gömülür

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
        Permissions: panel.PermissionConfig{
            Path: "permissions.toml", // İzin dosyası yolu
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

## 🛡 İzin Sistemi (RBAC)

Panel.go, rol tabanlı erişim kontrolü (RBAC) için yerleşik bir yapı sunar. İzinler bir `TOML` dosyasında tanımlanır ve her kullanıcı rolüne göre yönetilir.

### 1. İzin Dosyası (permissions.toml)

Proje kök dizininde (veya config'de belirttiğiniz yolda) bir TOML dosyası oluşturun:

```toml
# Sistemde kullanılacak roller
system_roles = ["admin", "editor", "user"]

[resources]
  # 'users' kaynağı için izinler
  [resources.users]
  label = "Kullanıcı Yönetimi"
  # Bu kaynağa ait aksiyonlar (backend policy'de kontrol edilir)
  actions = ["view_any", "view", "create", "update", "delete", "block"]

  [resources.posts]
  label = "İçerik Yönetimi"
  actions = ["view_any", "create", "update"]
```

### 2. Policy Entegrasyonu

Otomatik oluşturulan policy dosyalarınızda (`pkg/policy/`) `HasPermission` metodunu kullanarak yetki kontrolü yapabilirsiniz:

```go
func (p UserPolicy) View(ctx *appContext.Context, model interface{}) bool {
    // Kullanıcının "users" kaynağında "view" yetkisi var mı?
    // Format: {resource_identifier}.{action}
    return ctx.HasPermission("users.view")
}

func (p UserPolicy) Create(ctx *appContext.Context) bool {
    return ctx.HasPermission("users.create")
}
```

> **Not:** `admin` rolüne sahip kullanıcılar varsayılan olarak tüm yetkilere sahiptir (`HasPermission` her zaman `true` döner).

### 3. Kullanıcıya Rol Atama

Kullanıcı modelinizde `Role` alanı, `system_roles` içinde tanımlanan değerlerden biri olmalıdır.

```go
user := User{
    FullName: "Ahmet Yılmaz",
    Role:     "editor",
}
```

## 📝 Lisans

MIT License.
