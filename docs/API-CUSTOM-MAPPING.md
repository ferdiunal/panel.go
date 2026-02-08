# OpenAPI Custom Mapping

Panel.go, field type'larını OpenAPI schema'ya otomatik olarak map eder. Ancak bazı durumlarda bu mapping'i özelleştirmek isteyebilirsiniz. Bu dokümantasyon, custom mapping'in nasıl yapılacağını açıklar.

## İçindekiler

- [Giriş](#giriş)
- [Mapping Seviyeleri](#mapping-seviyeleri)
- [Global Field Type Mapping](#global-field-type-mapping)
- [Specific Field Mapping](#specific-field-mapping)
- [Resource-Level Mapping](#resource-level-mapping)
- [Örnekler](#örnekler)
- [Best Practices](#best-practices)
- [API Referansı](#api-referansı)

## Giriş

Panel.go, 30+ field type için varsayılan OpenAPI mapping'leri sağlar. Ancak:

- Özel field type'larınız varsa
- Varsayılan mapping'leri değiştirmek istiyorsanız
- Belirli field'lar için özel schema tanımlamak istiyorsanız

Custom mapping kullanabilirsiniz.

## Mapping Seviyeleri

Custom mapping 3 seviyede yapılabilir:

### 1. Global Field Type Mapping

Tüm field type'ları için geçerli olan mapping'ler.

**Kullanım Alanı**: Özel field type'ları veya varsayılan mapping'leri değiştirmek

**Öncelik**: En düşük (diğer mapping'ler bunu override eder)

### 2. Specific Field Mapping

Belirli bir resource'daki belirli bir field için mapping.

**Kullanım Alanı**: Tek bir field'ın özel davranışı

**Öncelik**: En yüksek (diğer tüm mapping'leri override eder)

### 3. Resource-Level Mapping

Bir resource'daki tüm field'lar için mapping.

**Kullanım Alanı**: Resource-spesifik davranışlar

**Öncelik**: Orta (global'i override eder, specific'i override etmez)

## Global Field Type Mapping

Tüm field type'ları için geçerli olan mapping'ler tanımlayın.

### Kullanım

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/openapi"
    "github.com/ferdiunal/panel.go/pkg/core"
)

func main() {
    p := panel.New(config)

    // Global field type mapping
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

### Örnek: Custom Field Type

```go
// Özel field type tanımla
const TYPE_MARKDOWN fields.ElementType = "markdown"

// Mapping ekle
p.OpenAPI().MapFieldType(TYPE_MARKDOWN, func(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "string",
        Format: "markdown",
        Description: "Markdown formatted text",
        Example: "# Hello World\n\nThis is **markdown**",
    }
})
```

### Örnek: Varsayılan Mapping'i Override Etme

```go
// TYPE_RICHTEXT için varsayılan mapping'i değiştir
p.OpenAPI().MapFieldType(fields.TYPE_RICHTEXT, func(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "object",
        Properties: map[string]*openapi.Schema{
            "html": {
                Type: "string",
                Description: "HTML content",
            },
            "text": {
                Type: "string",
                Description: "Plain text content",
            },
        },
        Required: []string{"html"},
    }
})
```

## Specific Field Mapping

Belirli bir resource'daki belirli bir field için mapping tanımlayın.

### Kullanım

```go
// users resource'undaki avatar field'ı için özel mapping
p.OpenAPI().MapField("users", "avatar", func(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "string",
        Format: "uri",
        Description: "User avatar URL",
        Example: "https://example.com/avatars/user123.jpg",
        Pattern: "^https?://.*\\.(jpg|jpeg|png|gif)$",
    }
})
```

### Örnek: Enum Field

```go
// users resource'undaki role field'ı için enum mapping
p.OpenAPI().MapField("users", "role", func(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "string",
        Enum: []interface{}{"admin", "editor", "viewer"},
        Description: "User role",
        Default: "viewer",
    }
})
```

### Örnek: Nested Object

```go
// users resource'undaki address field'ı için nested object mapping
p.OpenAPI().MapField("users", "address", func(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "object",
        Properties: map[string]*openapi.Schema{
            "street": {
                Type: "string",
                Description: "Street address",
            },
            "city": {
                Type: "string",
                Description: "City",
            },
            "country": {
                Type: "string",
                Description: "Country",
            },
            "postal_code": {
                Type: "string",
                Description: "Postal code",
                Pattern: "^[0-9]{5}$",
            },
        },
        Required: []string{"street", "city", "country"},
    }
})
```

## Resource-Level Mapping

Bir resource'daki tüm field'lar için mapping tanımlayın.

### Kullanım

```go
// users resource'undaki tüm field'lar için özel mapping
p.OpenAPI().MapResource("users", func(element core.Element) *openapi.Schema {
    // Element'in key'ine göre özel mapping
    switch element.GetKey() {
    case "email":
        return &openapi.Schema{
            Type: "string",
            Format: "email",
            Description: "User email address",
            Example: "user@example.com",
        }
    case "phone":
        return &openapi.Schema{
            Type: "string",
            Format: "tel",
            Description: "User phone number",
            Pattern: "^\\+?[1-9]\\d{1,14}$",
        }
    default:
        // Varsayılan mapping kullan
        return nil
    }
})
```

### Örnek: Conditional Mapping

```go
// products resource'undaki field'lar için koşullu mapping
p.OpenAPI().MapResource("products", func(element core.Element) *openapi.Schema {
    key := element.GetKey()

    // Fiyat field'ları için özel format
    if strings.HasSuffix(key, "_price") || strings.HasSuffix(key, "_cost") {
        return &openapi.Schema{
            Type: "number",
            Format: "decimal",
            Minimum: 0,
            Description: "Price in USD",
            Example: 99.99,
        }
    }

    // Varsayılan mapping kullan
    return nil
})
```

## Örnekler

### Örnek 1: E-Ticaret Uygulaması

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/openapi"
)

func main() {
    p := panel.New(config)

    // Ürün fiyatları için özel mapping
    p.OpenAPI().MapField("products", "price", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "number",
            Format: "decimal",
            Minimum: 0,
            Maximum: 999999.99,
            Description: "Product price in USD",
            Example: 29.99,
        }
    })

    // Ürün stok durumu için enum
    p.OpenAPI().MapField("products", "stock_status", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "string",
            Enum: []interface{}{"in_stock", "out_of_stock", "pre_order"},
            Description: "Product stock status",
            Default: "in_stock",
        }
    })

    // Ürün kategorisi için referans
    p.OpenAPI().MapField("products", "category_id", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "integer",
            Description: "Product category ID",
            Example: 1,
        }
    })

    p.RefreshOpenAPISpec()
    p.Start()
}
```

### Örnek 2: Blog Uygulaması

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/openapi"
)

func main() {
    p := panel.New(config)

    // Blog post content için markdown mapping
    p.OpenAPI().MapField("posts", "content", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "string",
            Format: "markdown",
            Description: "Post content in Markdown format",
            MinLength: 100,
            MaxLength: 50000,
        }
    })

    // Post status için enum
    p.OpenAPI().MapField("posts", "status", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "string",
            Enum: []interface{}{"draft", "published", "archived"},
            Description: "Post publication status",
            Default: "draft",
        }
    })

    // Post tags için array
    p.OpenAPI().MapField("posts", "tags", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "array",
            Items: &openapi.Schema{
                Type: "string",
            },
            Description: "Post tags",
            Example: []interface{}{"golang", "web", "api"},
        }
    })

    p.RefreshOpenAPISpec()
    p.Start()
}
```

### Örnek 3: Kullanıcı Yönetimi

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/openapi"
)

func main() {
    p := panel.New(config)

    // Kullanıcı avatar'ı için URL validation
    p.OpenAPI().MapField("users", "avatar", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "string",
            Format: "uri",
            Description: "User avatar URL",
            Pattern: "^https?://.*\\.(jpg|jpeg|png|gif|webp)$",
            Example: "https://example.com/avatars/user123.jpg",
        }
    })

    // Kullanıcı permissions için array of enums
    p.OpenAPI().MapField("users", "permissions", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "array",
            Items: &openapi.Schema{
                Type: "string",
                Enum: []interface{}{
                    "users.view",
                    "users.create",
                    "users.update",
                    "users.delete",
                    "posts.view",
                    "posts.create",
                    "posts.update",
                    "posts.delete",
                },
            },
            Description: "User permissions",
            UniqueItems: true,
        }
    })

    // Kullanıcı metadata için free-form object
    p.OpenAPI().MapField("users", "metadata", func(element core.Element) *openapi.Schema {
        return &openapi.Schema{
            Type: "object",
            AdditionalProperties: &openapi.Schema{
                Type: "string",
            },
            Description: "User metadata (key-value pairs)",
            Example: map[string]interface{}{
                "theme": "dark",
                "language": "en",
            },
        }
    })

    p.RefreshOpenAPISpec()
    p.Start()
}
```

## Best Practices

### 1. Spec'i Yenileme

Custom mapping'ler ekledikten sonra mutlaka spec'i yenileyin:

```go
p.OpenAPI().MapFieldType(...)
p.OpenAPI().MapField(...)
p.OpenAPI().MapResource(...)

// Spec'i yenile
p.RefreshOpenAPISpec()
```

### 2. Nil Döndürme

Varsayılan mapping'i kullanmak istiyorsanız `nil` döndürün:

```go
p.OpenAPI().MapResource("users", func(element core.Element) *openapi.Schema {
    if element.GetKey() == "special_field" {
        return &openapi.Schema{...}
    }
    // Diğer field'lar için varsayılan mapping kullan
    return nil
})
```

### 3. Validation Rules

OpenAPI schema'da validation rules ekleyin:

```go
&openapi.Schema{
    Type: "string",
    MinLength: 3,
    MaxLength: 50,
    Pattern: "^[a-zA-Z0-9_]+$",
}
```

### 4. Examples

Her zaman example değerler ekleyin:

```go
&openapi.Schema{
    Type: "string",
    Format: "email",
    Example: "user@example.com",
}
```

### 5. Descriptions

Açıklayıcı description'lar ekleyin:

```go
&openapi.Schema{
    Type: "string",
    Description: "User email address. Must be unique and valid.",
}
```

### 6. Format Kullanımı

Standart OpenAPI format'larını kullanın:

```go
// Standart format'lar
Format: "date"       // 2024-01-01
Format: "date-time"  // 2024-01-01T12:00:00Z
Format: "email"      // user@example.com
Format: "uri"        // https://example.com
Format: "uuid"       // 123e4567-e89b-12d3-a456-426614174000
Format: "binary"     // Base64 encoded binary data
Format: "byte"       // Base64 encoded string
Format: "password"   // Password (masked in UI)
```

### 7. Enum Kullanımı

Sabit değerler için enum kullanın:

```go
&openapi.Schema{
    Type: "string",
    Enum: []interface{}{"active", "inactive", "pending"},
    Default: "pending",
}
```

### 8. Array Validation

Array field'ları için validation ekleyin:

```go
&openapi.Schema{
    Type: "array",
    Items: &openapi.Schema{
        Type: "string",
    },
    MinItems: 1,
    MaxItems: 10,
    UniqueItems: true,
}
```

## API Referansı

### OpenAPI() Metodu

Panel instance'ından custom mapping registry'sine erişim sağlar.

```go
func (p *Panel) OpenAPI() *openapi.CustomMappingRegistry
```

**Döndürür**: Custom mapping registry pointer'ı

### MapFieldType()

Global field type mapping ekler.

```go
func (r *CustomMappingRegistry) MapFieldType(
    fieldType fields.ElementType,
    mapper func(element core.Element) *openapi.Schema,
)
```

**Parametreler**:
- `fieldType`: Field type (örn: `fields.TYPE_TEXT`)
- `mapper`: Mapping fonksiyonu

### MapField()

Specific field mapping ekler.

```go
func (r *CustomMappingRegistry) MapField(
    resourceSlug string,
    fieldKey string,
    mapper func(element core.Element) *openapi.Schema,
)
```

**Parametreler**:
- `resourceSlug`: Resource slug (örn: "users")
- `fieldKey`: Field key (örn: "email")
- `mapper`: Mapping fonksiyonu

### MapResource()

Resource-level mapping ekler.

```go
func (r *CustomMappingRegistry) MapResource(
    resourceSlug string,
    mapper func(element core.Element) *openapi.Schema,
)
```

**Parametreler**:
- `resourceSlug`: Resource slug (örn: "users")
- `mapper`: Mapping fonksiyonu

### RefreshOpenAPISpec()

OpenAPI spec cache'ini temizler.

```go
func (p *Panel) RefreshOpenAPISpec()
```

## Sorun Giderme

### Custom Mapping Çalışmıyor

**Sorun**: Custom mapping eklediğiniz halde spec'te görünmüyor.

**Çözüm**:
1. `RefreshOpenAPISpec()` metodunu çağırdığınızdan emin olun
2. Mapping fonksiyonunun `nil` döndürmediğinden emin olun
3. Resource slug ve field key'in doğru olduğundan emin olun

### Spec Cache'lenmiş

**Sorun**: Değişiklikler spec'te görünmüyor.

**Çözüm**:
```go
p.RefreshOpenAPISpec()
```

### Mapping Önceliği

**Sorun**: Hangi mapping'in kullanılacağını bilmiyorum.

**Çözüm**: Öncelik sırası:
1. Specific Field Mapping (en yüksek)
2. Resource-Level Mapping
3. Global Field Type Mapping (en düşük)

## İleri Seviye

### Dynamic Mapping

Runtime'da dynamic mapping yapabilirsiniz:

```go
// Veritabanından enum değerlerini al
func getStatusEnums(db *gorm.DB) []interface{} {
    var statuses []string
    db.Model(&Status{}).Pluck("name", &statuses)

    enums := make([]interface{}, len(statuses))
    for i, s := range statuses {
        enums[i] = s
    }
    return enums
}

// Dynamic enum mapping
p.OpenAPI().MapField("orders", "status", func(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "string",
        Enum: getStatusEnums(db),
        Description: "Order status",
    }
})
```

### Conditional Mapping

Koşullu mapping yapabilirsiniz:

```go
p.OpenAPI().MapResource("products", func(element core.Element) *openapi.Schema {
    key := element.GetKey()

    // Fiyat field'ları
    if strings.Contains(key, "price") {
        return &openapi.Schema{
            Type: "number",
            Format: "decimal",
            Minimum: 0,
        }
    }

    // Boolean field'lar
    if strings.HasPrefix(key, "is_") || strings.HasPrefix(key, "has_") {
        return &openapi.Schema{
            Type: "boolean",
        }
    }

    return nil
})
```

### Reusable Mappers

Tekrar kullanılabilir mapper'lar oluşturabilirsiniz:

```go
// Reusable mapper fonksiyonları
func priceMapper(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "number",
        Format: "decimal",
        Minimum: 0,
        Description: "Price in USD",
    }
}

func emailMapper(element core.Element) *openapi.Schema {
    return &openapi.Schema{
        Type: "string",
        Format: "email",
        Description: "Email address",
    }
}

// Kullanım
p.OpenAPI().MapField("products", "price", priceMapper)
p.OpenAPI().MapField("users", "email", emailMapper)
```

## Kaynaklar

- [OpenAPI Specification](https://swagger.io/specification/)
- [JSON Schema](https://json-schema.org/)
- [OpenAPI Data Types](https://swagger.io/docs/specification/data-models/data-types/)
- [API-OPENAPI.md](./API-OPENAPI.md)

## Sonraki Adımlar

- [OpenAPI Dokümantasyonu](./API-OPENAPI.md)
- [Field Dokümantasyonu](./Fields.md)
- [Resource Dokümantasyonu](./Resources.md)
