# Panel.go

Panel.go, Go + GORM ile admin paneli ve CRUD API'yi hızlıca ayağa kaldırmak için geliştirilmiş bir framework'tür.

Bu repo içinde hem SDK hem de `panel` CLI bulunur.

## Kimin için?

- Go backend geliştiricileri
- GORM kullanan ekipler
- Admin panelini sıfırdan yazmak yerine hızlıca üretmek isteyen projeler

## Neler sunar?

- Resource tabanlı yapı (model, field, policy, repository)
- Otomatik CRUD endpoint'leri
- Hazır admin UI (Go binary içine gömülü)
- Relationship field'ları (`BelongsTo`, `HasMany`, `BelongsToMany`, `MorphTo`)
- Policy ve rol/izin yönetimi
- Lens, Action, Page ve Widget desteği
- OpenAPI/Swagger üretimi
- Plugin sistemi

## Hızlı Başlangıç

### 1) SDK'yı projene ekle

```bash
go get github.com/ferdiunal/panel.go
```

### 2) CLI kur (önerilen)

```bash
go install github.com/ferdiunal/panel.go/cmd/panel@latest
```

### 3) Proje iskeletini üret

```bash
panel init
```

Bu komut:
- başlangıç dosyalarını oluşturur
- veritabanı seçimine göre örnek konfigürasyon yazar
- `.panel/stubs/` ve `.claude/skills/` dosyalarını yayınlar

Detay: [`docs/CLI_INIT.md`](docs/CLI_INIT.md)

### 4) Uygulamayı çalıştır

`main.go` içinde paneli başlatıp resource'larını kaydet:

```go
cfg := panel.Config{
    Server: panel.ServerConfig{Host: "localhost", Port: "8080"},
    Database: panel.DatabaseConfig{Instance: db},
    Environment: "development",
}

app := panel.New(cfg)
app.RegisterResource(GetUserResource())
app.Start()
```

İlk resource örneği için: [`docs/Getting-Started.md`](docs/Getting-Started.md)

## Otomatik Açılan API Yapısı

Bir resource register edildiğinde bu endpoint'ler otomatik gelir:

- `GET /api/resource/{slug}`
- `POST /api/resource/{slug}`
- `GET /api/resource/{slug}/:id`
- `PUT /api/resource/{slug}/:id`
- `DELETE /api/resource/{slug}/:id`

## Dokümantasyon Rotası (Son Kullanıcı)

### 1. Kurulum ve temel kullanım
- [`docs/Getting-Started.md`](docs/Getting-Started.md)
- [`docs/Resources.md`](docs/Resources.md)
- [`docs/Fields.md`](docs/Fields.md)
- [`docs/Validation.md`](docs/Validation.md)
- [`docs/Relationships.md`](docs/Relationships.md)

### 2. Güvenlik ve erişim
- [`docs/Authentication.md`](docs/Authentication.md)
- [`docs/Authorization.md`](docs/Authorization.md)

### 3. Arayüzü zenginleştirme
- [`docs/Actions.md`](docs/Actions.md)
- [`docs/Lenses.md`](docs/Lenses.md)
- [`docs/Widgets.md`](docs/Widgets.md)
- [`docs/Pages.md`](docs/Pages.md)
- [`docs/Settings.md`](docs/Settings.md)
- [`docs/Notifications.md`](docs/Notifications.md)

### 4. API ve entegrasyon
- [`docs/API-Reference.md`](docs/API-Reference.md)
- [`docs/API-OPENAPI.md`](docs/API-OPENAPI.md)
- [`docs/API-CUSTOM-MAPPING.md`](docs/API-CUSTOM-MAPPING.md)

### 5. İleri seviye
- [`docs/Advanced-Usage.md`](docs/Advanced-Usage.md)
- [`docs/Optimization-Guide.md`](docs/Optimization-Guide.md)
- [`docs/PLUGIN_SYSTEM.md`](docs/PLUGIN_SYSTEM.md)
- [`docs/PLUGIN_DEVELOPMENT.md`](docs/PLUGIN_DEVELOPMENT.md)

Tam menü: [`docs/_Sidebar.md`](docs/_Sidebar.md)

## CLI Kısa Komutlar

```bash
panel make:resource Product
panel make:model Product
panel make:page Dashboard
panel make:lens ActiveProducts --resource product
panel make:action Publish --resource post
```

Plugin komutları için: [`docs/PLUGIN_CLI.md`](docs/PLUGIN_CLI.md)

## Notlar

- UI dosyaları repo içinde gömülü gelir; normal kullanımda ayrıca frontend build zorunlu değildir.
- Frontend tarafında değişiklik yaparsan `make build-ui` ile UI varlıklarını yeniden üretmelisin.

## Lisans

[MIT](LICENSE)
