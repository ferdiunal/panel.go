# OpenAPI/Swagger Entegrasyonu

Panel.go, external API resource endpoint'leri için otomatik OpenAPI 3.0.3 spesifikasyonu oluşturur ve Swagger UI arayüzü sağlar.

## İçindekiler

- [Özellikler](#özellikler)
- [Endpoint'ler](#endpointler)
- [Kapsam ve Sınırlar](#kapsam-ve-sınırlar)
- [Kullanım](#kullanım)
- [Swagger UI](#swagger-ui)
- [OpenAPI Spec](#openapi-spec)
- [Özelleştirme](#özelleştirme)

## Özellikler

- ✅ **Otomatik Spec Oluşturma**: Resource'lar ve field'lar otomatik olarak OpenAPI schema'ya çevrilir
- ✅ **CRUD Endpoint'leri**: List, Get, Create, Update, Delete endpoint'leri otomatik oluşturulur
- ✅ **Action Endpoint'leri**: Resource action'ları OpenAPI spec'e dahil edilir
- ✅ **Filter Parametreleri**: Resource filter'ları query parametreleri olarak eklenir
- ✅ **Swagger UI**: Modern, interaktif API dokümantasyonu
- ✅ **Custom Mapping**: Field type'ları için özel mapping desteği

## Endpoint'ler

Panel.go başlatıldığında aşağıdaki endpoint'ler otomatik olarak oluşturulur:

### Dokümantasyon Arayüzleri

| Endpoint | Açıklama |
|----------|----------|
| `GET /api/docs` | Swagger UI arayüzü |
| `GET /api/openapi.json` | OpenAPI 3.0.3 spec (JSON) |

### Resource Endpoint'leri (Otomatik)

Her resource için aşağıdaki endpoint'ler otomatik oluşturulur:

| Endpoint | Method | Açıklama |
|----------|--------|----------|
| `/api/{slug}` | GET | Resource listesi (paginated) |
| `/api/{slug}` | POST | Yeni kayıt oluştur |
| `/api/{slug}/{id}` | GET | Tek kayıt getir |
| `/api/{slug}/{id}` | PUT | Kayıt güncelle |
| `/api/{slug}/{id}` | DELETE | Kayıt sil |
| `/api/{slug}/actions/{action}` | POST | Action çalıştır |

## Kapsam ve Sınırlar

- OpenAPI/Swagger çıktısı yalnızca external API (`/api/:resource`) yüzeyini dokümante eder.
- Internal panel endpoint'leri (`/api/internal/*`) ve internal REST endpoint'leri (`/api/internal/rest/*`) `/api/openapi.json` içinde yer almaz.
- Swagger operasyon tag'leri resource `Group()` değerine göre değil, resource bazlı (title/slug tabanlı) üretilir.
- `/api/docs` ve `/api/openapi.json` endpoint'leri publictir; session gerektirmez.

## Kullanım

### Swagger UI'ya Erişim

Tarayıcınızda aşağıdaki URL'yi açın:

```
http://localhost:8080/api/docs
```

Swagger UI, tüm API endpoint'lerini interaktif olarak gösterir. Her endpoint'i test edebilir, request/response örneklerini görebilirsiniz.

### OpenAPI Spec'i İndirme

OpenAPI spec'i JSON formatında indirmek için:

```bash
curl http://localhost:8080/api/openapi.json > openapi.json
```

## Swagger UI

Swagger UI, API'nizi interaktif olarak test etmenizi sağlar:

### Özellikler

- **Try it out**: Endpoint'leri doğrudan tarayıcıdan test edin
- **Authentication**: API key ile kimlik doğrulama
- **Request/Response**: Örnek request ve response'ları görün
- **Schema**: Veri modellerini inceleyin
- **Filter**: Endpoint'leri arayın ve filtreleyin

### Kullanım

1. Swagger UI'yı açın: `http://localhost:8080/api/docs`
2. Bir endpoint seçin (örn: `GET /api/users`)
3. "Try it out" butonuna tıklayın
4. Parametreleri doldurun
5. "Execute" butonuna tıklayın
6. Response'u görün

### Authentication

Swagger UI'da authentication kullanmak için:

1. Sağ üstteki "Authorize" butonuna tıklayın
2. API key ile çalışacaksanız `apiKeyAuth` bölümüne key değerini girin
3. "Authorize" butonuna tıklayın
4. Artık endpoint'leri API key ile test edebilirsiniz

## OpenAPI Spec

OpenAPI spec'i programatik olarak kullanabilirsiniz:

### JSON Formatında

```bash
curl http://localhost:8080/api/openapi.json
```

### Postman'e Import

1. OpenAPI spec'i indirin
2. Postman'i açın
3. File > Import
4. openapi.json dosyasını seçin
5. Tüm endpoint'ler Postman'e import edilir

### Code Generation

OpenAPI spec'i kullanarak client kodu oluşturabilirsiniz:

```bash
# TypeScript client
npx @openapitools/openapi-generator-cli generate \
  -i http://localhost:8080/api/openapi.json \
  -g typescript-axios \
  -o ./client

# Go client
openapi-generator-cli generate \
  -i http://localhost:8080/api/openapi.json \
  -g go \
  -o ./client

# Python client
openapi-generator-cli generate \
  -i http://localhost:8080/api/openapi.json \
  -g python \
  -o ./client
```

## Özelleştirme

### OpenAPI Config

Panel.go başlatırken OpenAPI config'i özelleştirebilirsiniz:

```go
config := panel.Config{
    Name: "My Admin Panel",
    // ... diğer config'ler
}

p := panel.New(config)
```

OpenAPI spec'te `config.Name` değeri title olarak kullanılır.

### Custom Field Mapping

Field type'larını OpenAPI schema'ya nasıl map edileceğini özelleştirebilirsiniz.

Detaylı bilgi için [API-CUSTOM-MAPPING.md](./API-CUSTOM-MAPPING.md) dosyasına bakın.

### Spec'i Yenileme

Custom mapping'ler ekledikten sonra spec'i yenilemek için:

```go
panel.RefreshOpenAPISpec()
```

## Örnek Kullanım

### Basit Kullanım

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
)

func main() {
    config := panel.Config{
        Name: "My Admin Panel",
        Database: panel.DatabaseConfig{
            Instance: db,
        },
        Server: panel.ServerConfig{
            Host: "localhost",
            Port: "8080",
        },
    }

    p := panel.New(config)

    // OpenAPI endpoint'leri otomatik olarak oluşturulur
    // http://localhost:8080/api/docs - Swagger UI
    // http://localhost:8080/api/openapi.json - OpenAPI Spec

    p.Start()
}
```

### Custom Mapping ile Kullanım

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/openapi"
)

func main() {
    config := panel.Config{
        Name: "My Admin Panel",
        // ... config
    }

    p := panel.New(config)

    // Custom field type mapping
    p.OpenAPI().MapFieldType(fields.TYPE_RICHTEXT, func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "string",
            Format: "html",
            Description: "Rich text content (HTML)",
        }
    })

    // Spec'i yenile
    p.RefreshOpenAPISpec()

    p.Start()
}
```

## Best Practices

### 1. Authentication

Swagger UI'da authentication kullanırken:

- API key kullanın
- Key'i güvenli bir şekilde saklayın
- Production'da HTTPS kullanın

### 2. Versioning

API versioning için:

- URL'de version kullanın: `/api/v1/resources/users`
- OpenAPI spec'te version belirtin

### 3. Documentation

- Field'lara açıklayıcı label'lar ekleyin
- Validation rules ekleyin
- Help text kullanın

### 4. Security

- Production'da `/api/docs` endpoint'ini kapatmayı düşünün
- Veya authentication gerektirin
- Rate limiting uygulayın

## Sorun Giderme

### Spec Oluşturulmuyor

Eğer OpenAPI spec oluşturulmuyorsa:

1. Resource'ların doğru register edildiğinden emin olun
2. Field'ların doğru tanımlandığından emin olun
3. Log'ları kontrol edin

### Swagger UI Açılmıyor

Eğer Swagger UI açılmıyorsa:

1. URL'nin doğru olduğundan emin olun: `/api/docs`
2. Panel.go'nun başlatıldığından emin olun
3. Port'un doğru olduğundan emin olun

### Custom Mapping Çalışmıyor

Eğer custom mapping çalışmıyorsa:

1. `RefreshOpenAPISpec()` metodunu çağırdığınızdan emin olun
2. Mapping fonksiyonunun doğru olduğundan emin olun
3. Field type'ın doğru olduğundan emin olun

## İleri Seviye

### OpenAPI Extensions

OpenAPI spec'e custom extension'lar ekleyebilirsiniz:

```go
// TODO: Extension desteği eklenecek
```

### Custom Response Codes

Custom response code'ları ekleyebilirsiniz:

```go
// TODO: Custom response code desteği eklenecek
```

### Custom Security Schemes

Custom security scheme'leri ekleyebilirsiniz:

```go
// TODO: Custom security scheme desteği eklenecek
```

## Kaynaklar

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Generator](https://openapi-generator.tech/)

## Sonraki Adımlar

- [Custom Mapping Dokümantasyonu](./API-CUSTOM-MAPPING.md)
- [Field Dokümantasyonu](./Fields.md)
- [Resource Dokümantasyonu](./Resources.md)
