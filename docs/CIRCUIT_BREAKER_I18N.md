# Circuit Breaker ve i18n Middleware Kullanım Kılavuzu

Bu dokümantasyon, Panel.go'da Circuit Breaker ve i18n (Internationalization) middleware'lerinin nasıl kullanılacağını açıklar.

## İçindekiler

1. [Circuit Breaker Middleware](#circuit-breaker-middleware)
2. [i18n Middleware](#i18n-middleware)
3. [Yapılandırma Örnekleri](#yapılandırma-örnekleri)
4. [Kullanım Örnekleri](#kullanım-örnekleri)

---

## Circuit Breaker Middleware

### Genel Bakış

Circuit Breaker, servis hatalarını yönetmek ve sistem çökmelerini önlemek için kullanılan bir dayanıklılık (resilience) desenidir. Üç durum arasında geçiş yapar:

- **Closed (Kapalı)**: Normal çalışma, istekler geçer, hatalar sayılır
- **Open (Açık)**: Devre açık, istekler hemen reddedilir (503 Service Unavailable)
- **Half-Open (Yarı Açık)**: Test modu, sınırlı sayıda istek geçer

### Yapılandırma

```go
config := panel.Config{
    // ... diğer yapılandırmalar
    CircuitBreaker: panel.CircuitBreakerConfig{
        Enabled:                true,
        FailureThreshold:       5,      // 5 ardışık hata sonrası devre aç
        Timeout:                10 * time.Second,  // 10 saniye bekle
        SuccessThreshold:       5,      // 5 başarılı istek sonrası devre kapat
        HalfOpenMaxConcurrent:  1,      // Half-open'da 1 eşzamanlı istek
    },
}
```

### Parametreler

| Parametre | Tip | Varsayılan | Açıklama |
|-----------|-----|------------|----------|
| `Enabled` | bool | false | Circuit Breaker'ı etkinleştirir |
| `FailureThreshold` | int | 5 | Devre açılmadan önce kaç ardışık hata olması gerektiği |
| `Timeout` | time.Duration | 10s | Devre açıldıktan sonra ne kadar süre bekleneceği |
| `SuccessThreshold` | int | 5 | Devre kapanmadan önce kaç başarılı istek olması gerektiği |
| `HalfOpenMaxConcurrent` | int | 1 | Half-open durumunda kaç eşzamanlı istek izin verileceği |

### Durum Geçişleri

```
Closed (Normal)
    |
    | FailureThreshold aşıldı
    v
Open (Devre Açık)
    |
    | Timeout süresi doldu
    v
Half-Open (Test)
    |
    | SuccessThreshold başarılı istek
    v
Closed (Normal)
```

### Hata Yanıtları

**Devre Açık (Open):**
```json
{
  "error": "Service temporarily unavailable",
  "message": "The service is experiencing high failure rates. Please try again later.",
  "code": "CIRCUIT_BREAKER_OPEN"
}
```
HTTP Status: 503 Service Unavailable

**Half-Open (Test):**
```json
{
  "error": "Service is recovering",
  "message": "The service is testing recovery. Please wait.",
  "code": "CIRCUIT_BREAKER_HALF_OPEN"
}
```
HTTP Status: 429 Too Many Requests

### Best Practices

1. **Kritik Servislere Uygulayın**: Dış API'ler, veritabanı bağlantıları
2. **Timeout Değerini Ayarlayın**: Servis yanıt süresine göre optimize edin
3. **Monitoring Ekleyin**: Circuit breaker durumlarını izleyin
4. **Fallback Mekanizmaları**: Alternatif yanıtlar tanımlayın
5. **Test Edin**: Farklı hata senaryolarını test edin

---

## i18n Middleware

### Genel Bakış

i18n (Internationalization) middleware, uygulamanın farklı dillerde gösterilmesini sağlar. go-i18n kütüphanesi kullanılarak mesajların çevirileri yönetilir.

### Yapılandırma

```go
import "golang.org/x/text/language"

config := panel.Config{
    // ... diğer yapılandırmalar
    I18n: panel.I18nConfig{
        Enabled:          true,
        RootPath:         "./locales",
        AcceptLanguages:  []language.Tag{language.Turkish, language.English},
        DefaultLanguage:  language.Turkish,
        FormatBundleFile: "yaml",
    },
}
```

### Parametreler

| Parametre | Tip | Varsayılan | Açıklama |
|-----------|-----|------------|----------|
| `Enabled` | bool | false | i18n'i etkinleştirir |
| `RootPath` | string | "./locales" | Dil dosyalarının bulunduğu dizin |
| `AcceptLanguages` | []language.Tag | [tr, en] | Desteklenen diller listesi |
| `DefaultLanguage` | language.Tag | Turkish | Varsayılan dil (fallback) |
| `FormatBundleFile` | string | "yaml" | Dil dosyası formatı (yaml, json, toml) |

### Dil Dosyası Yapısı

**Dizin Yapısı:**
```
locales/
├── tr.yaml
└── en.yaml
```

**Örnek Dil Dosyası (locales/tr.yaml):**
```yaml
# Genel Mesajlar
welcome:
  other: "Hoş geldiniz"

welcomeWithName:
  other: "Hoş geldiniz, {{.Name}}"

# Hata Mesajları
error:
  notFound:
    other: "Kayıt bulunamadı"
  unauthorized:
    other: "Bu işlem için yetkiniz yok"

# Circuit Breaker Mesajları
circuitBreaker:
  open:
    other: "Servis geçici olarak kullanılamıyor. Lütfen birkaç dakika sonra tekrar deneyin."
```

### Dil Seçimi

Dil, şu sırayla belirlenir:

1. **Query Parametresi**: `?lang=tr`
2. **Accept-Language Header**: `Accept-Language: tr-TR,tr;q=0.9,en;q=0.8`
3. **DefaultLanguage**: Fallback dil

### Kullanım

**Handler'da Çeviri Kullanımı:**

```go
import (
    "github.com/gofiber/contrib/fiberi18n/v2"
    "github.com/gofiber/fiber/v2"
)

func MyHandler(c *fiber.Ctx) error {
    // Basit çeviri
    message := fiberi18n.MustLocalize(c, "welcome")
    // Çıktı: "Hoş geldiniz"

    // Template değişkenleri ile çeviri
    message := fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
        MessageID: "welcomeWithName",
        TemplateData: map[string]string{
            "Name": "Ahmet",
        },
    })
    // Çıktı: "Hoş geldiniz, Ahmet"

    return c.JSON(fiber.Map{
        "message": message,
    })
}
```

**Dil Değiştirme:**

```bash
# Query parametresi ile
curl http://localhost:8080/api/resource/users?lang=en

# Header ile
curl -H "Accept-Language: en-US,en;q=0.9" http://localhost:8080/api/resource/users
```

---

## Yapılandırma Örnekleri

### Minimal Yapılandırma

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "golang.org/x/text/language"
    "time"
)

func main() {
    config := panel.Config{
        Server: panel.ServerConfig{
            Host: "localhost",
            Port: "8080",
        },
        Database: panel.DatabaseConfig{
            Instance: db, // GORM instance
        },
        Environment: "development",

        // Circuit Breaker - Varsayılan değerlerle
        CircuitBreaker: panel.CircuitBreakerConfig{
            Enabled: true,
        },

        // i18n - Varsayılan değerlerle
        I18n: panel.I18nConfig{
            Enabled: true,
        },
    }

    p := panel.New(config)
    p.Start()
}
```

### Üretim Yapılandırması

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "golang.org/x/text/language"
    "time"
)

func main() {
    config := panel.Config{
        Server: panel.ServerConfig{
            Host: "0.0.0.0",
            Port: "8080",
        },
        Database: panel.DatabaseConfig{
            Instance: db,
        },
        Environment: "production",

        // Circuit Breaker - Üretim ayarları
        CircuitBreaker: panel.CircuitBreakerConfig{
            Enabled:                true,
            FailureThreshold:       10,     // Daha yüksek eşik
            Timeout:                30 * time.Second,  // Daha uzun timeout
            SuccessThreshold:       10,     // Daha fazla başarılı istek
            HalfOpenMaxConcurrent:  2,      // Daha fazla test isteği
        },

        // i18n - Çoklu dil desteği
        I18n: panel.I18nConfig{
            Enabled:          true,
            RootPath:         "/etc/panel/locales",
            AcceptLanguages:  []language.Tag{
                language.Turkish,
                language.English,
                language.German,
                language.French,
            },
            DefaultLanguage:  language.English,
            FormatBundleFile: "yaml",
        },
    }

    p := panel.New(config)
    p.Start()
}
```

---

## Kullanım Örnekleri

### Örnek 1: Circuit Breaker ile Dış API Çağrısı

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/gofiber/fiber/v2"
    "time"
)

func main() {
    config := panel.Config{
        // ... temel yapılandırma
        CircuitBreaker: panel.CircuitBreakerConfig{
            Enabled:                true,
            FailureThreshold:       3,
            Timeout:                5 * time.Second,
            SuccessThreshold:       2,
            HalfOpenMaxConcurrent:  1,
        },
    }

    p := panel.New(config)

    // Dış API çağrısı yapan endpoint
    p.Fiber.Get("/api/external", func(c *fiber.Ctx) error {
        // Circuit breaker otomatik olarak çalışır
        // Eğer dış API çökerse, circuit breaker devreye girer

        resp, err := http.Get("https://external-api.com/data")
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "error": "External API error",
            })
        }

        return c.JSON(resp)
    })

    p.Start()
}
```

### Örnek 2: i18n ile Çoklu Dil Desteği

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/gofiber/contrib/fiberi18n/v2"
    "github.com/gofiber/fiber/v2"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
)

func main() {
    config := panel.Config{
        // ... temel yapılandırma
        I18n: panel.I18nConfig{
            Enabled:          true,
            RootPath:         "./locales",
            AcceptLanguages:  []language.Tag{language.Turkish, language.English},
            DefaultLanguage:  language.Turkish,
            FormatBundleFile: "yaml",
        },
    }

    p := panel.New(config)

    // Çeviri kullanan endpoint
    p.Fiber.Get("/api/welcome", func(c *fiber.Ctx) error {
        // Basit çeviri
        message := fiberi18n.MustLocalize(c, "welcome")

        return c.JSON(fiber.Map{
            "message": message,
        })
    })

    // Template değişkenleri ile çeviri
    p.Fiber.Get("/api/welcome/:name", func(c *fiber.Ctx) error {
        name := c.Params("name")

        message := fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
            MessageID: "welcomeWithName",
            TemplateData: map[string]string{
                "Name": name,
            },
        })

        return c.JSON(fiber.Map{
            "message": message,
        })
    })

    p.Start()
}
```

### Örnek 3: Her İkisini Birlikte Kullanma

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/gofiber/contrib/fiberi18n/v2"
    "github.com/gofiber/fiber/v2"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
    "time"
)

func main() {
    config := panel.Config{
        Server: panel.ServerConfig{
            Host: "localhost",
            Port: "8080",
        },
        Database: panel.DatabaseConfig{
            Instance: db,
        },
        Environment: "production",

        // Circuit Breaker etkin
        CircuitBreaker: panel.CircuitBreakerConfig{
            Enabled:                true,
            FailureThreshold:       5,
            Timeout:                10 * time.Second,
            SuccessThreshold:       5,
            HalfOpenMaxConcurrent:  1,
        },

        // i18n etkin
        I18n: panel.I18nConfig{
            Enabled:          true,
            RootPath:         "./locales",
            AcceptLanguages:  []language.Tag{language.Turkish, language.English},
            DefaultLanguage:  language.Turkish,
            FormatBundleFile: "yaml",
        },
    }

    p := panel.New(config)

    // Circuit breaker korumalı ve çoklu dil destekli endpoint
    p.Fiber.Get("/api/data", func(c *fiber.Ctx) error {
        // Dış API çağrısı (circuit breaker korumalı)
        data, err := fetchExternalData()
        if err != nil {
            // Hata mesajını kullanıcının diline göre döndür
            errorMsg := fiberi18n.MustLocalize(c, "error.serviceUnavailable")
            return c.Status(503).JSON(fiber.Map{
                "error": errorMsg,
            })
        }

        // Başarı mesajını kullanıcının diline göre döndür
        successMsg := fiberi18n.MustLocalize(c, "success.dataFetched")

        return c.JSON(fiber.Map{
            "message": successMsg,
            "data":    data,
        })
    })

    p.Start()
}

func fetchExternalData() (interface{}, error) {
    // Dış API çağrısı simülasyonu
    // Circuit breaker bu fonksiyonun hatalarını yönetir
    return nil, nil
}
```

---

## Monitoring ve Debugging

### Circuit Breaker Durumunu İzleme

Circuit breaker durumunu izlemek için log'ları kontrol edin:

```bash
# Devre açıldığında
[CIRCUIT_BREAKER] State changed: CLOSED -> OPEN (failures: 5)

# Half-open durumuna geçildiğinde
[CIRCUIT_BREAKER] State changed: OPEN -> HALF_OPEN (timeout: 10s)

# Devre kapandığında
[CIRCUIT_BREAKER] State changed: HALF_OPEN -> CLOSED (successes: 5)
```

### i18n Debugging

Dil seçimini debug etmek için:

```go
// Handler'da mevcut dili kontrol et
func MyHandler(c *fiber.Ctx) error {
    lang := fiberi18n.MustGetLocale(c)
    fmt.Printf("Current language: %s\n", lang)

    // ...
}
```

---

## Sorun Giderme

### Circuit Breaker Sorunları

**Problem**: Circuit breaker çok sık açılıyor
- **Çözüm**: `FailureThreshold` değerini artırın
- **Çözüm**: `Timeout` değerini artırın

**Problem**: Circuit breaker hiç açılmıyor
- **Çözüm**: `Enabled` değerinin `true` olduğunu kontrol edin
- **Çözüm**: Hata yanıtlarının 500+ status kodu döndürdüğünü kontrol edin

### i18n Sorunları

**Problem**: Çeviriler gösterilmiyor
- **Çözüm**: `Enabled` değerinin `true` olduğunu kontrol edin
- **Çözüm**: Dil dosyalarının doğru dizinde olduğunu kontrol edin
- **Çözüm**: YAML formatının doğru olduğunu kontrol edin

**Problem**: Yanlış dil gösteriliyor
- **Çözüm**: `Accept-Language` header'ını kontrol edin
- **Çözüm**: `?lang=` query parametresini kullanın
- **Çözüm**: `DefaultLanguage` ayarını kontrol edin

---

## Kaynaklar

- [i18n Helper Fonksiyonları Kullanım Kılavuzu](I18N_HELPERS.md)
- [Fiber Circuit Breaker Dokümantasyonu](https://docs.gofiber.io/contrib/circuitbreaker/)
- [Fiber i18n Dokümantasyonu](https://docs.gofiber.io/contrib/fiberi18n/)
- [go-i18n Kütüphanesi](https://github.com/nicksnyder/go-i18n)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
