# Ayarlar (Settings)

Ayarlar modülü, uygulamanız için anahtar-değer (key-value) tabanlı yapılandırmaları veritabanında saklamanızı ve yönetmenizi sağlar.

## Veri Modeli

Ayarlar, `settings` tablosunda saklanır. Değerler JSONB formatında tutulduğu için esnek bir yapı sunar.

```go
type Setting struct {
    Key       string                 `gorm:"primaryKey;type:varchar(255)" json:"key"`
    Value     map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"value"`
    CreatedAt time.Time              `json:"created_at"`
    UpdatedAt time.Time              `json:"updated_at"`
}
```

## Kullanım

### Ayarlar Sayfası

Ayarların yönetimi için önceden tanımlanmış bir `Settings` sayfası bulunur. Bu sayfa `pkg/page/settings.go` dosyasında tanımlıdır.

Ayarlar sayfasını kullanmak için `Config` içinde etkinleştirmeniz ve alanlarını tanımlamanız yeterlidir:

```go
// main.go
func main() {
    cfg := panel.Config{
        // Settings sayfasını yapılandırma
        SettingsPage: &page.Settings{
            Elements: []fields.Element{
                 // Standart alanlar
                 fields.Text("Site Name", "site_name").Required(),
                 fields.Switch("Registration Open", "register"),
                 
                 // Özel alanlar (Dinamik)
                 fields.Text("Support Email", "support_email"),
                 fields.Image("Logo", "site_logo"),
            },
        },
        // ...
    }
    
    app := panel.New(cfg)
    app.Start()
}
```

### Ayarlara Erişim

Panel instance'ı üzerinden ayarlara her yerden erişebilirsiniz:

```go
// Tüm ayarları yükler (cache mechanism eklenebilir)
settings := app.LoadSettings()

// Belirli bir değere erişim
if enabled, ok := settings.Values["register"].(bool); ok && enabled {
    // Kayıt işlemleri...
}

// veya Config üzerinden (Startup sırasında yüklenenler)
siteName := app.Config.SettingsValues.Values["site_name"]
```

### Feature Flags

Bazı ayarlar (`register`, `forgot_password`) Panel tarafından otomatik olarak tanınır ve `feature flag` olarak davranır:

- `register: true` -> Kayıt olma endpointlerini açar.
- `forgot_password: true` -> Şifre sıfırlama endpointlerini açar.
