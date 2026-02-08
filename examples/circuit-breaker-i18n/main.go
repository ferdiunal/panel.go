package main

import (
	"log"
	"time"

	"github.com/ferdiunal/panel.go/pkg/panel"
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Veritabanı bağlantısı
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanı bağlantısı başarısız:", err)
	}

	// Panel yapılandırması
	config := panel.Config{
		Server: panel.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Database: panel.DatabaseConfig{
			Instance: db,
		},
		Environment: "development",

		// Circuit Breaker - Etkin
		CircuitBreaker: panel.CircuitBreakerConfig{
			Enabled:                true,
			FailureThreshold:       3,      // Test için düşük değer
			Timeout:                5 * time.Second,
			SuccessThreshold:       2,
			HalfOpenMaxConcurrent:  1,
		},

		// i18n - Etkin
		I18n: panel.I18nConfig{
			Enabled:          true,
			RootPath:         "./locales",
			AcceptLanguages:  []language.Tag{language.Turkish, language.English},
			DefaultLanguage:  language.Turkish,
			FormatBundleFile: "yaml",
		},
	}

	// Panel oluştur
	p := panel.New(config)

	// Test endpoint'leri ekle
	setupTestEndpoints(p)

	// Sunucuyu başlat
	log.Println("🚀 Panel.go başlatılıyor...")
	log.Println("📍 Adres: http://localhost:8080")
	log.Println("🌍 i18n: Etkin (tr, en)")
	log.Println("🔌 Circuit Breaker: Etkin")
	log.Println("")
	log.Println("Test Endpoint'leri:")
	log.Println("  GET  /api/test/welcome          - Basit çeviri")
	log.Println("  GET  /api/test/welcome/:name    - Template ile çeviri")
	log.Println("  GET  /api/test/error            - Circuit breaker testi (hata)")
	log.Println("  GET  /api/test/success          - Circuit breaker testi (başarılı)")
	log.Println("")
	log.Println("Dil değiştirme:")
	log.Println("  ?lang=tr  - Türkçe")
	log.Println("  ?lang=en  - İngilizce")
	log.Println("")

	if err := p.Start(); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}

func setupTestEndpoints(p *panel.Panel) {
	// Test endpoint'leri için grup
	test := p.Fiber.Group("/api/test")

	// 1. Basit çeviri testi
	test.Get("/welcome", func(c *fiber.Ctx) error {
		message := fiberi18n.MustLocalize(c, "welcome")

		return c.JSON(fiber.Map{
			"message": message,
			"lang":    getLang(c),
		})
	})

	// 2. Template değişkenleri ile çeviri testi
	test.Get("/welcome/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")

		message := fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
			MessageID: "welcomeWithName",
			TemplateData: map[string]string{
				"Name": name,
			},
		})

		return c.JSON(fiber.Map{
			"message": message,
			"name":    name,
			"lang":    getLang(c),
		})
	})

	// 3. Circuit breaker testi - Hata simülasyonu
	errorCount := 0
	test.Get("/error", func(c *fiber.Ctx) error {
		errorCount++

		// İlk 5 istekte hata döndür (circuit breaker'ı tetikle)
		if errorCount <= 5 {
			log.Printf("❌ Hata simülasyonu: %d/5", errorCount)
			return c.Status(500).JSON(fiber.Map{
				"error":   "Simulated error",
				"count":   errorCount,
				"message": "Bu hata circuit breaker'ı tetiklemek için simüle edildi",
			})
		}

		// 5 hatadan sonra başarılı yanıt döndür
		log.Printf("✅ Başarılı yanıt: %d", errorCount)
		return c.JSON(fiber.Map{
			"success": true,
			"count":   errorCount,
			"message": "Servis kurtarıldı",
		})
	})

	// 4. Circuit breaker testi - Başarılı istek
	test.Get("/success", func(c *fiber.Ctx) error {
		message := fiberi18n.MustLocalize(c, "success.created")

		return c.JSON(fiber.Map{
			"success": true,
			"message": message,
			"lang":    getLang(c),
		})
	})

	// 5. Tüm çevirileri listele
	test.Get("/translations", func(c *fiber.Ctx) error {
		lang := getLang(c)

		translations := map[string]string{
			"welcome":                    fiberi18n.MustLocalize(c, "welcome"),
			"error.notFound":             fiberi18n.MustLocalize(c, "error.notFound"),
			"error.unauthorized":         fiberi18n.MustLocalize(c, "error.unauthorized"),
			"error.serverError":          fiberi18n.MustLocalize(c, "error.serverError"),
			"circuitBreaker.open":        fiberi18n.MustLocalize(c, "circuitBreaker.open"),
			"success.created":            fiberi18n.MustLocalize(c, "success.created"),
			"button.save":                fiberi18n.MustLocalize(c, "button.save"),
			"navigation.dashboard":       fiberi18n.MustLocalize(c, "navigation.dashboard"),
		}

		return c.JSON(fiber.Map{
			"lang":         lang,
			"translations": translations,
		})
	})
}

// getLang, fiber context'ten dil bilgisini güvenli bir şekilde alır.
// Type assertion panic riskini önler.
func getLang(c *fiber.Ctx) string {
	if lang, ok := c.Locals("lang").(string); ok && lang != "" {
		return lang
	}
	return "en" // fallback to default
}
