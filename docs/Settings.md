# Ayarlar (Settings)

Ayarlar modülü, uygulamanız için anahtar-değer (key-value) tabanlı yapılandırmaları veritabanında saklamanızı ve yönetmenizi sağlar. Bu özellik, `internal/domain/setting` altında tanımlanan `Setting` modeli ve Panel'in sayfa sistemi üzerine kuruludur.

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

### Ayar Tanımlama ve Seeding

Kaynaklarınızda varsayılan ayarları veya gruplandırılmış yapılandırmaları tanımlamak için `resource.Base` içerisindeki `SettingsSeed` yapısını kullanabilirsiniz.

```go
type MyResource struct {
    resource.Base
}

func NewMyResource() *MyResource {
    r := &MyResource{}
    r.Seed = resource.SettingsSeed{
        Key: "my_resource_config",
        Value: map[string]interface{}{
            "feature_enabled": true,
            "max_limit": 100,
        },
    }
    return r
}
```

### Ayarlar Sayfası

Ayarların yönetimi için önceden tanımlanmış bir `Settings` sayfası bulunur. Bu sayfa `internal/page/settings.go` dosyasında tanımlıdır ve uygulamanıza `main.go` üzerinden kaydedilir.

```go
// main.go
settingsPage := &page.Settings{}
p.RegisterPage(settingsPage)
```

Sayfa kaydedildikten sonra `/api/pages/settings` endpoint'i üzerinden erişilebilir hale gelir.
