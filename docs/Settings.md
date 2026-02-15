# Ayarlar (Settings) - Legacy Teknik Akış

Ayarlar modülü, uygulamanız için anahtar-değer (key-value) tabanlı yapılandırmaları veritabanında saklamanızı ve yönetmenizi sağlar.

## Bu Doküman Ne Zaman Okunmalı?

Önerilen sıra:
1. [Başlarken](Getting-Started)
2. [Sayfalar (Pages)](Pages)
3. [Yetkilendirme](Authorization)
4. Bu doküman (`Settings`)

## Hızlı Settings Akışı

1. Settings alanlarını `page.Settings{ Elements: ... }` ile tanımla.
2. `Config.Pages` altında settings page'i kaydet.
3. Ayarları `app.LoadSettings()` veya `app.Config.SettingsValues` üzerinden tüket.
4. Feature flag alanlarını (`register`, `forgot_password`) iş mantığına bağla.

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
        // Settings sayfasını Pages dizisi ile yapılandırma
        Pages: []page.Page{
            &page.Settings{
                Elements: []fields.Element{
                     // Standart alanlar
                     fields.Text("Site Name", "site_name").Required(),
                     fields.Switch("Registration Open", "register"),
                     
                     // Özel alanlar (Dinamik)
                     fields.Text("Support Email", "support_email"),
                     fields.Image("Logo", "site_logo"),
                },
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

## Sık Hata Kontrolü (Settings)

- Ayar değeri okunmuyor: key adı ile field key'inin birebir aynı olduğunu kontrol edin.
- JSON parse/type hatası: `map[string]interface{}` cast noktalarında tip kontrolü ekleyin.
- Feature flag etkisiz: ilgili key'in gerçekten `settings` verisinde güncellendiğini doğrulayın.
- Sayfa görünmüyor: `Pages` listesinde `page.Settings{...}` kaydının yapıldığını kontrol edin.

## Sonraki Adım

- Settings sayfası detayları için: [Settings Page](SETTINGS_PAGE)
- Kullanıcı profil ayarları için: [Account Page](ACCOUNT_PAGE)
- Bildirim tercihleri gibi dinamik alanlar için: [Bildirimler (Notifications)](Notifications)
