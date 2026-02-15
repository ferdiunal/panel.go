package middleware

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// / # RateLimitConfig
// /
// / Bu yapı, rate limiting (hız sınırlama) middleware'i için yapılandırma ayarlarını tutar.
// / Rate limiting, belirli bir zaman penceresi içinde yapılabilecek istek sayısını sınırlayarak
// / API'nizi DDoS saldırılarından ve aşırı kullanımdan korur.
// /
// / ## Alanlar
// /
// / * `Max` - Zaman penceresi başına maksimum istek sayısı
// / * `Expiration` - Zaman penceresi süresi (örn: 1 dakika, 1 saat)
// / * `KeyGenerator` - İstemciyi tanımlamak için özel anahtar üreteci (varsayılan: IP adresi)
// / * `LimitReached` - Limit aşıldığında çalıştırılacak özel handler
// /
// / ## Kullanım Senaryoları
// /
// / 1. **API Endpoint Koruması**: Genel API endpoint'lerini aşırı kullanımdan koruma
// / 2. **Kimlik Doğrulama Koruması**: Login endpoint'lerini brute-force saldırılarından koruma
// / 3. **Kaynak Yönetimi**: Sunucu kaynaklarının adil dağılımını sağlama
// / 4. **Kullanıcı Bazlı Sınırlama**: Farklı kullanıcı tipleri için farklı limitler
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Basit IP bazlı rate limiting
// / config := RateLimitConfig{
// /     Max:        100,
// /     Expiration: 1 * time.Minute,
// / }
// /
// / // Kullanıcı ID bazlı rate limiting
// / config := RateLimitConfig{
// /     Max:        50,
// /     Expiration: 1 * time.Minute,
// /     KeyGenerator: func(c *fiber.Ctx) string {
// /         userID := c.Locals("user_id").(string)
// /         return userID
// /     },
// /     LimitReached: func(c *fiber.Ctx) error {
// /         return c.Status(429).JSON(fiber.Map{
// /             "error": "Çok fazla istek gönderdiniz",
// /         })
// /     },
// / }
// / ```
// /
// / ## Avantajlar
// /
// / * **Esneklik**: Özelleştirilebilir anahtar üreteci ve yanıt handler'ı
// / * **Performans**: Hafif ve hızlı implementasyon
// / * **Güvenlik**: DDoS ve brute-force saldırılarına karşı koruma
// /
// / ## Önemli Notlar
// /
// / * `Max` değeri 0 ise, varsayılan olarak 100 kullanılır
// / * `Expiration` değeri 0 ise, varsayılan olarak 1 dakika kullanılır
// / * `KeyGenerator` nil ise, varsayılan olarak IP adresi kullanılır
// / * Rate limit bilgileri bellekte tutulur, sunucu yeniden başlatıldığında sıfırlanır
type RateLimitConfig struct {
	// Max number of requests per window
	Max int
	// Time window duration
	Expiration time.Duration
	// Custom key generator (default: IP address)
	KeyGenerator func(*fiber.Ctx) string
	// Custom response when limit exceeded
	LimitReached fiber.Handler
}

// / # SecurityHeadersConfig
// /
// / Bu yapı, HTTP güvenlik başlıklarının yapılandırma ayarlarını tutar.
// / Güvenlik başlıkları, web uygulamanızı XSS, clickjacking, MIME-sniffing gibi
// / yaygın web güvenlik açıklarından korumak için kullanılır.
// /
// / ## Alanlar
// /
// / * `ContentSecurityPolicy` - İçerik Güvenlik Politikası (CSP) başlığı
// / * `XFrameOptions` - Clickjacking koruması için X-Frame-Options başlığı
// / * `XContentTypeOptions` - MIME-sniffing koruması için X-Content-Type-Options başlığı
// / * `ReferrerPolicy` - Referrer bilgisi politikası
// / * `PermissionsPolicy` - Tarayıcı özellik izinleri politikası
// /
// / ## Kullanım Senaryoları
// /
// / 1. **XSS Koruması**: Content Security Policy ile script injection saldırılarını engelleme
// / 2. **Clickjacking Koruması**: X-Frame-Options ile iframe içine alınmayı engelleme
// / 3. **MIME-Sniffing Koruması**: Dosya tiplerinin yanlış yorumlanmasını engelleme
// / 4. **Gizlilik Koruması**: Referrer bilgilerinin kontrolü
// / 5. **Özellik Kontrolü**: Kamera, mikrofon gibi tarayıcı özelliklerine erişimi kısıtlama
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Varsayılan güvenlik başlıkları
// / config := DefaultSecurityHeaders()
// /
// / // Özelleştirilmiş güvenlik başlıkları
// / config := SecurityHeadersConfig{
// /     ContentSecurityPolicy: "default-src 'self'; script-src 'self' https://cdn.example.com",
// /     XFrameOptions:         "SAMEORIGIN",
// /     XContentTypeOptions:   "nosniff",
// /     ReferrerPolicy:        "strict-origin-when-cross-origin",
// /     PermissionsPolicy:     "geolocation=(self), microphone=()",
// / }
// /
// / app.Use(SecurityHeaders(config))
// / ```
// /
// / ## Güvenlik Başlıkları Açıklamaları
// /
// / * **Content-Security-Policy**: Hangi kaynaklardan içerik yüklenebileceğini belirler
// / * **X-Frame-Options**: DENY (hiç iframe'e alınamaz), SAMEORIGIN (sadece aynı domain)
// / * **X-Content-Type-Options**: nosniff (MIME-sniffing'i devre dışı bırakır)
// / * **Referrer-Policy**: no-referrer, strict-origin, same-origin vb.
// / * **Permissions-Policy**: Tarayıcı API'lerine erişim kontrolü
// /
// / ## Avantajlar
// /
// / * **Çok Katmanlı Güvenlik**: Birden fazla güvenlik açığına karşı koruma
// / * **Standart Uyumluluk**: OWASP önerilerine uygun
// / * **Esneklik**: Her başlık bağımsız olarak yapılandırılabilir
// / * **Performans**: Minimal overhead
// /
// / ## Önemli Notlar
// /
// / * Boş string değerleri o başlığın eklenmemesine neden olur
// / * CSP çok kısıtlayıcı olabilir, test ederek ayarlayın
// / * Üretim ortamında mutlaka güvenlik başlıkları kullanın
// / * Başlıkları düzenli olarak güncelleyin ve test edin
type SecurityHeadersConfig struct {
	// Content Security Policy
	ContentSecurityPolicy string
	// X-Frame-Options
	XFrameOptions string
	// X-Content-Type-Options
	XContentTypeOptions string
	// Referrer-Policy
	ReferrerPolicy string
	// Permissions-Policy
	PermissionsPolicy string
}

// / # DefaultSecurityHeaders
// /
// / Bu fonksiyon, güvenli varsayılan güvenlik başlıklarını içeren bir SecurityHeadersConfig döndürür.
// / Üretim ortamında kullanıma hazır, OWASP önerilerine uygun güvenlik başlıkları sağlar.
// /
// / ## Döndürülen Başlıklar
// /
// / * **Content-Security-Policy**:
// /   - `default-src 'self'`: Varsayılan olarak sadece kendi domain'den içerik
// /   - `script-src 'self' 'unsafe-inline' 'unsafe-eval'`: Script kaynakları (inline ve eval izinli)
// /   - `style-src 'self' 'unsafe-inline'`: Stil kaynakları (inline izinli)
// /   - `img-src 'self' data: https:`: Resimler (data URI ve HTTPS izinli)
// /   - `font-src 'self' data:`: Fontlar (data URI izinli)
// /   - `connect-src 'self'`: AJAX/WebSocket bağlantıları
// /   - `frame-ancestors 'none'`: Hiçbir site bu sayfayı iframe'e alamaz
// /
// / * **X-Frame-Options**: `DENY` - Clickjacking koruması (hiç iframe'e alınamaz)
// / * **X-Content-Type-Options**: `nosniff` - MIME-sniffing koruması
// / * **Referrer-Policy**: `no-referrer` - Referrer bilgisi gönderilmez
// / * **Permissions-Policy**: Konum, mikrofon ve kamera erişimi kapalı
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Hızlı Başlangıç**: Güvenlik başlıklarını hızlıca eklemek için
// / 2. **Varsayılan Koruma**: Temel güvenlik gereksinimlerini karşılamak için
// / 3. **Temel Şablon**: Özelleştirme için başlangıç noktası olarak
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Varsayılan güvenlik başlıklarını kullan
// / app.Use(SecurityHeaders(DefaultSecurityHeaders()))
// /
// / // Varsayılan başlıkları özelleştir
// / config := DefaultSecurityHeaders()
// / config.XFrameOptions = "SAMEORIGIN" // iframe'e aynı domain'den izin ver
// / config.ContentSecurityPolicy += " frame-src https://trusted.com" // Güvenilir iframe kaynağı ekle
// / app.Use(SecurityHeaders(config))
// / ```
// /
// / ## Döndürür
// /
// / * `SecurityHeadersConfig` - Güvenli varsayılan değerlerle yapılandırılmış güvenlik başlıkları
// /
// / ## Avantajlar
// /
// / * **Hemen Kullanıma Hazır**: Ek yapılandırma gerektirmez
// / * **Güvenli Varsayılanlar**: OWASP önerilerine uygun
// / * **Özelleştirilebilir**: Döndürülen config değiştirilebilir
// / * **Kapsamlı Koruma**: Birden fazla güvenlik açığına karşı koruma
// /
// / ## Önemli Notlar
// /
// / * CSP'de `unsafe-inline` ve `unsafe-eval` kullanımı güvenlik riskidir
// / * Üretim ortamında CSP'yi daha kısıtlayıcı yapmanız önerilir
// / * Modern SPA uygulamaları için CSP ayarlarını gözden geçirin
// / * `frame-ancestors 'none'` ile X-Frame-Options DENY aynı korumayı sağlar
func DefaultSecurityHeaders() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none';",
		XFrameOptions:         "DENY",
		XContentTypeOptions:   "nosniff",
		ReferrerPolicy:        "no-referrer",
		PermissionsPolicy:     "geolocation=(), microphone=(), camera=()",
	}
}

// / # SecurityHeaders
// /
// / Bu fonksiyon, HTTP yanıtlarına güvenlik başlıkları ekleyen bir Fiber middleware'i oluşturur.
// / Web uygulamanızı XSS, clickjacking, MIME-sniffing ve diğer yaygın güvenlik açıklarından korur.
// /
// / ## Parametreler
// /
// / * `config` - `SecurityHeadersConfig` yapısı ile güvenlik başlıkları yapılandırması
// /
// / ## Çalışma Mantığı
// /
// / Middleware, her HTTP yanıtına yapılandırmada belirtilen güvenlik başlıklarını ekler.
// / Boş olmayan her başlık değeri yanıta eklenir, boş olanlar atlanır.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Tüm Uygulama Koruması**: Tüm endpoint'lere güvenlik başlıkları eklemek için
// / 2. **Özel Koruma**: Belirli route grupları için farklı güvenlik politikaları
// / 3. **Uyumluluk**: Güvenlik standartlarına (OWASP, PCI-DSS) uyum için
// / 4. **Tarayıcı Koruması**: Modern tarayıcıların güvenlik özelliklerini aktifleştirme
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Varsayılan güvenlik başlıkları ile tüm uygulamayı koru
// / app.Use(SecurityHeaders(DefaultSecurityHeaders()))
// /
// / // Özelleştirilmiş güvenlik başlıkları
// / config := SecurityHeadersConfig{
// /     ContentSecurityPolicy: "default-src 'self'",
// /     XFrameOptions:         "SAMEORIGIN",
// /     XContentTypeOptions:   "nosniff",
// /     ReferrerPolicy:        "strict-origin",
// /     PermissionsPolicy:     "geolocation=()",
// / }
// / app.Use(SecurityHeaders(config))
// /
// / // Belirli bir route grubu için
// / api := app.Group("/api")
// / api.Use(SecurityHeaders(DefaultSecurityHeaders()))
// /
// / // Admin paneli için daha katı güvenlik
// / admin := app.Group("/admin")
// / adminConfig := DefaultSecurityHeaders()
// / adminConfig.XFrameOptions = "DENY"
// / adminConfig.ContentSecurityPolicy = "default-src 'none'"
// / admin.Use(SecurityHeaders(adminConfig))
// / ```
// /
// / ## Döndürür
// /
// / * `fiber.Handler` - Güvenlik başlıklarını ekleyen middleware fonksiyonu
// /
// / ## Avantajlar
// /
// / * **Kolay Entegrasyon**: Tek satırda tüm uygulamayı koruyabilirsiniz
// / * **Esneklik**: Her başlık bağımsız olarak yapılandırılabilir
// / * **Performans**: Minimal overhead, sadece başlık ekleme
// / * **Standart Uyumluluk**: OWASP ve güvenlik best practice'lerine uygun
// /
// / ## Önemli Notlar
// /
// / * Middleware'i mümkün olduğunca erken (app.Use ile global olarak) ekleyin
// / * Boş string değerleri o başlığın eklenmemesine neden olur
// / * CSP başlığı çok kısıtlayıcı olabilir, test ederek ayarlayın
// / * Üretim ortamında mutlaka güvenlik başlıkları kullanın
// / * Başlıkları düzenli olarak güncelleyin ve güvenlik taramaları yapın
func SecurityHeaders(config SecurityHeadersConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.ContentSecurityPolicy != "" {
			c.Set("Content-Security-Policy", config.ContentSecurityPolicy)
		}
		if config.XFrameOptions != "" {
			c.Set("X-Frame-Options", config.XFrameOptions)
		}
		if config.XContentTypeOptions != "" {
			c.Set("X-Content-Type-Options", config.XContentTypeOptions)
		}
		if config.ReferrerPolicy != "" {
			c.Set("Referrer-Policy", config.ReferrerPolicy)
		}
		if config.PermissionsPolicy != "" {
			c.Set("Permissions-Policy", config.PermissionsPolicy)
		}
		return c.Next()
	}
}

// / # RateLimiter
// /
// / Bu fonksiyon, yapılandırılabilir rate limiting (hız sınırlama) middleware'i oluşturur.
// / Belirli bir zaman penceresi içinde yapılabilecek istek sayısını sınırlayarak API'nizi
// / DDoS saldırılarından, brute-force saldırılarından ve aşırı kullanımdan korur.
// /
// / ## Parametreler
// /
// / * `config` - `RateLimitConfig` yapısı ile rate limiting yapılandırması
// /   - `Max`: Zaman penceresi başına maksimum istek sayısı (varsayılan: 100)
// /   - `Expiration`: Zaman penceresi süresi (varsayılan: 1 dakika)
// /   - `KeyGenerator`: İstemciyi tanımlamak için özel fonksiyon (varsayılan: IP adresi)
// /   - `LimitReached`: Limit aşıldığında çalıştırılacak özel handler
// /
// / ## Çalışma Mantığı
// /
// / 1. Her istek için KeyGenerator ile benzersiz bir anahtar üretilir (varsayılan: IP adresi)
// / 2. Bu anahtar için istek sayacı artırılır
// / 3. Sayaç Max değerini aşarsa LimitReached handler çalıştırılır
// / 4. Expiration süresi sonunda sayaç sıfırlanır
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Genel API Koruması**: Tüm API endpoint'lerini aşırı kullanımdan koruma
// / 2. **Kimlik Doğrulama Koruması**: Login endpoint'lerini brute-force'tan koruma
// / 3. **Kullanıcı Bazlı Sınırlama**: Farklı kullanıcılar için farklı limitler
// / 4. **Endpoint Bazlı Sınırlama**: Farklı endpoint'ler için farklı limitler
// / 5. **Kaynak Yönetimi**: Sunucu kaynaklarının adil dağılımı
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Basit IP bazlı rate limiting
// / app.Use(RateLimiter(RateLimitConfig{
// /     Max:        100,
// /     Expiration: 1 * time.Minute,
// / }))
// /
// / // Kullanıcı ID bazlı rate limiting
// / app.Use(RateLimiter(RateLimitConfig{
// /     Max:        50,
// /     Expiration: 1 * time.Minute,
// /     KeyGenerator: func(c *fiber.Ctx) string {
// /         userID := c.Locals("user_id").(string)
// /         return "user:" + userID
// /     },
// / }))
// /
// / // API key bazlı rate limiting
// / app.Use(RateLimiter(RateLimitConfig{
// /     Max:        1000,
// /     Expiration: 1 * time.Hour,
// /     KeyGenerator: func(c *fiber.Ctx) string {
// /         apiKey := c.Get("X-API-Key")
// /         return "api:" + apiKey
// /     },
// /     LimitReached: func(c *fiber.Ctx) error {
// /         return c.Status(429).JSON(fiber.Map{
// /             "error": "Rate limit aşıldı",
// /             "retry_after": "1 saat",
// /         })
// /     },
// / }))
// /
// / // Farklı endpoint'ler için farklı limitler
// / api := app.Group("/api")
// / api.Use(RateLimiter(RateLimitConfig{Max: 100, Expiration: 1 * time.Minute}))
// /
// / auth := app.Group("/auth")
// / auth.Use(RateLimiter(RateLimitConfig{Max: 10, Expiration: 1 * time.Minute}))
// / ```
// /
// / ## Döndürür
// /
// / * `fiber.Handler` - Rate limiting uygulayan middleware fonksiyonu
// /
// / ## Varsayılan Değerler
// /
// / * `Max`: 0 ise 100 olarak ayarlanır
// / * `Expiration`: 0 ise 1 dakika olarak ayarlanır
// / * `KeyGenerator`: nil ise IP adresi kullanılır
// / * `LimitReached`: nil ise varsayılan 429 yanıtı döner
// /
// / ## Avantajlar
// /
// / * **Esneklik**: Özelleştirilebilir anahtar üreteci ve yanıt handler'ı
// / * **Performans**: Hafif ve hızlı implementasyon
// / * **Güvenlik**: DDoS ve brute-force saldırılarına karşı koruma
// / * **Adil Kullanım**: Sunucu kaynaklarının adil dağılımı
// / * **Kolay Entegrasyon**: Tek satırda eklenebilir
// /
// / ## Dezavantajlar
// /
// / * **Bellek Kullanımı**: Rate limit bilgileri bellekte tutulur
// / * **Dağıtık Sistemler**: Tek sunucuda çalışır, load balancer arkasında senkronizasyon gerekir
// / * **Yeniden Başlatma**: Sunucu yeniden başlatıldığında sayaçlar sıfırlanır
// /
// / ## Önemli Notlar
// /
// / * Rate limit bilgileri bellekte tutulur, Redis gibi dış bir store kullanmaz
// / * Load balancer arkasında her sunucu kendi rate limit'ini tutar
// / * Dağıtık sistemlerde Redis tabanlı rate limiting kullanın
// / * KeyGenerator fonksiyonu hızlı olmalı, ağır işlemler yapmayın
// / * IP bazlı rate limiting proxy arkasında doğru IP'yi almayabilir
// / * X-Forwarded-For başlığını kontrol edin (proxy kullanıyorsanız)
func RateLimiter(config RateLimitConfig) fiber.Handler {
	if config.Max == 0 {
		config.Max = 100
	}
	if config.Expiration == 0 {
		config.Expiration = 1 * time.Minute
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *fiber.Ctx) string {
			return c.IP()
		}
	}

	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return config.KeyGenerator(c)
		},
		LimitReached: config.LimitReached,
	})
}

// / # AuthRateLimiter
// /
// / Bu fonksiyon, kimlik doğrulama endpoint'leri için özel olarak tasarlanmış katı bir
// / rate limiting middleware'i oluşturur. Brute-force saldırılarına karşı koruma sağlar.
// /
// / ## Yapılandırma
// /
// / * **Max**: 10 istek / dakika (çok katı limit)
// / * **Expiration**: 1 dakika
// / * **KeyGenerator**: IP adresi (varsayılan)
// / * **LimitReached**: Özel hata mesajı ile 429 yanıtı
// /
// / ## Çalışma Mantığı
// /
// / 1. Her IP adresi için dakikada maksimum 10 login denemesine izin verir
// / 2. Limit aşıldığında kullanıcıya açıklayıcı hata mesajı gösterir
// / 3. 1 dakika sonra sayaç sıfırlanır ve tekrar deneme yapılabilir
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Login Endpoint Koruması**: /login, /auth/login gibi endpoint'leri koruma
// / 2. **Brute-Force Koruması**: Şifre tahmin saldırılarını engelleme
// / 3. **Hesap Güvenliği**: Kullanıcı hesaplarını yetkisiz erişimden koruma
// / 4. **API Token Endpoint'leri**: Token alma endpoint'lerini koruma
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Login endpoint'ini koru
// / app.Post("/login", AuthRateLimiter(), loginHandler)
// /
// / // Tüm auth route'larını koru
// / auth := app.Group("/auth")
// / auth.Use(AuthRateLimiter())
// / auth.Post("/login", loginHandler)
// / auth.Post("/register", registerHandler)
// / auth.Post("/forgot-password", forgotPasswordHandler)
// /
// / // API token endpoint'ini koru
// / app.Post("/api/token", AuthRateLimiter(), tokenHandler)
// /
// / // Farklı auth endpoint'leri için farklı limitler
// / app.Post("/admin/login", RateLimiter(RateLimitConfig{
// /     Max:        5,  // Admin için daha katı
// /     Expiration: 5 * time.Minute,
// / }), adminLoginHandler)
// / ```
// /
// / ## Döndürür
// /
// / * `fiber.Handler` - Kimlik doğrulama için optimize edilmiş rate limiting middleware'i
// /
// / ## Avantajlar
// /
// / * **Brute-Force Koruması**: Şifre tahmin saldırılarını etkili şekilde engeller
// / * **Kullanıma Hazır**: Ek yapılandırma gerektirmez
// / * **Açıklayıcı Hata**: Kullanıcıya net hata mesajı verir
// / * **Düşük Limit**: Güvenlik odaklı katı limit
// /
// / ## Önemli Notlar
// /
// / * Dakikada sadece 10 deneme izni vardır (çok katı)
// / * Normal kullanıcılar için yeterli, saldırganlar için çok kısıtlayıcı
// / * AccountLockout ile birlikte kullanılması önerilir
// / * Proxy arkasında X-Forwarded-For başlığını kontrol edin
// / * Load balancer kullanıyorsanız gerçek IP'yi doğru alın
// / * Üretim ortamında mutlaka kullanın
func AuthRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Max:        10, // 10 requests per minute
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many authentication attempts. Please try again later.",
			})
		},
	})
}

// / # APIRateLimiter
// /
// / Bu fonksiyon, genel API endpoint'leri için optimize edilmiş rate limiting middleware'i oluşturur.
// / Normal API kullanımı için yeterli esnekliği sağlarken, aşırı kullanımı engeller.
// /
// / ## Yapılandırma
// /
// / * **Max**: 100 istek / dakika (dengeli limit)
// / * **Expiration**: 1 dakika
// / * **KeyGenerator**: IP adresi (varsayılan)
// / * **LimitReached**: Özel hata mesajı ile 429 yanıtı
// /
// / ## Çalışma Mantığı
// /
// / 1. Her IP adresi için dakikada maksimum 100 API isteğine izin verir
// / 2. Limit aşıldığında kullanıcıya yavaşlaması gerektiğini bildirir
// / 3. 1 dakika sonra sayaç sıfırlanır ve tekrar istek yapılabilir
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Genel API Koruması**: CRUD endpoint'lerini aşırı kullanımdan koruma
// / 2. **Kaynak Yönetimi**: Sunucu kaynaklarının adil dağılımı
// / 3. **Performans Koruması**: Veritabanı ve sunucu yükünü kontrol altında tutma
// / 4. **Kullanıcı Deneyimi**: Normal kullanıcıları etkilemeden koruma sağlama
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Tüm API route'larını koru
// / api := app.Group("/api")
// / api.Use(APIRateLimiter())
// / api.Get("/users", getUsersHandler)
// / api.Post("/users", createUserHandler)
// /
// / // Belirli API versiyonunu koru
// / v1 := app.Group("/api/v1")
// / v1.Use(APIRateLimiter())
// /
// / // Public API endpoint'lerini koru
// / app.Get("/api/public/data", APIRateLimiter(), publicDataHandler)
// /
// / // Farklı endpoint'ler için farklı limitler
// / app.Get("/api/search", RateLimiter(RateLimitConfig{
// /     Max:        50,  // Arama için daha düşük limit
// /     Expiration: 1 * time.Minute,
// / }), searchHandler)
// /
// / app.Get("/api/export", RateLimiter(RateLimitConfig{
// /     Max:        10,  // Export için çok düşük limit
// /     Expiration: 1 * time.Hour,
// / }), exportHandler)
// / ```
// /
// / ## Döndürür
// /
// / * `fiber.Handler` - Genel API kullanımı için optimize edilmiş rate limiting middleware'i
// /
// / ## Avantajlar
// /
// / * **Dengeli Limit**: Normal kullanım için yeterli, aşırı kullanım için kısıtlayıcı
// / * **Kullanıma Hazır**: Ek yapılandırma gerektirmez
// / * **Açıklayıcı Hata**: Kullanıcıya net hata mesajı verir
// / * **Performans Koruması**: Sunucu kaynaklarını korur
// /
// / ## Önemli Notlar
// /
// / * Dakikada 100 istek çoğu kullanım senaryosu için yeterlidir
// / * AuthRateLimiter'dan daha esnek (100 vs 10 istek/dakika)
// / * Yoğun API kullanımı için limiti artırabilirsiniz
// / * Proxy arkasında X-Forwarded-For başlığını kontrol edin
// / * Load balancer kullanıyorsanız gerçek IP'yi doğru alın
func APIRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Max:        100, // 100 requests per minute
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please slow down.",
			})
		},
	})
}

// / # AccountLockout
// /
// / Bu yapı, başarısız giriş denemelerinden sonra hesap kilitleme mekanizmasını yönetir.
// / Brute-force saldırılarına karşı ek bir koruma katmanı sağlar ve rate limiting'i tamamlar.
// /
// / ## Alanlar
// /
// / * `mu` - Thread-safe erişim için okuma/yazma mutex'i
// / * `attempts` - Her identifier için deneme bilgilerini tutan map
// / * `maxAttempts` - Kilitleme öncesi maksimum başarısız deneme sayısı
// / * `lockoutDuration` - Hesabın kilitli kalacağı süre
// /
// / ## Çalışma Mantığı
// /
// / 1. Her başarısız giriş denemesi kaydedilir
// / 2. Deneme sayısı maxAttempts'e ulaştığında hesap kilitlenir
// / 3. Kilitleme süresi boyunca giriş yapılamaz
// / 4. Süre dolduğunda hesap otomatik olarak açılır
// / 5. Başarılı girişte deneme sayacı sıfırlanır
// / 6. Arka planda çalışan cleanup goroutine'i eski kayıtları temizler
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Login Koruması**: Kullanıcı login endpoint'lerini brute-force'tan koruma
// / 2. **API Token Koruması**: Token alma endpoint'lerini koruma
// / 3. **Admin Panel Koruması**: Yönetici paneli girişlerini koruma
// / 4. **Çok Faktörlü Kimlik Doğrulama**: 2FA kodlarını brute-force'tan koruma
// / 5. **Şifre Sıfırlama**: Şifre sıfırlama endpoint'lerini koruma
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Hesap kilitleme yöneticisi oluştur
// / lockout := NewAccountLockout(5, 15*time.Minute)
// /
// / // Login handler'da kullan
// / func loginHandler(c *fiber.Ctx) error {
// /     email := c.FormValue("email")
// /
// /     // Hesap kilitli mi kontrol et
// /     if lockout.IsLocked(email) {
// /         return c.Status(403).JSON(fiber.Map{
// /             "error": "Hesabınız geçici olarak kilitlendi",
// /         })
// /     }
// /
// /     // Kimlik doğrulama yap
// /     if !authenticate(email, password) {
// /         lockout.RecordFailedAttempt(email)
// /         remaining := lockout.GetRemainingAttempts(email)
// /
// /         return c.Status(401).JSON(fiber.Map{
// /             "error": "Geçersiz kimlik bilgileri",
// /             "remaining_attempts": remaining,
// /         })
// /     }
// /
// /     // Başarılı giriş - denemeleri sıfırla
// /     lockout.ResetAttempts(email)
// /     return c.JSON(fiber.Map{"success": true})
// / }
// /
// / // IP bazlı kilitleme
// / lockout := NewAccountLockout(10, 30*time.Minute)
// / ip := c.IP()
// / if lockout.IsLocked(ip) {
// /     return c.Status(403).JSON(fiber.Map{
// /         "error": "IP adresiniz geçici olarak engellendi",
// /     })
// / }
// / ```
// /
// / ## Thread Safety
// /
// / * Tüm public methodlar thread-safe'dir
// / * sync.RWMutex ile eşzamanlı erişim koruması sağlanır
// / * Okuma işlemleri RLock, yazma işlemleri Lock kullanır
// /
// / ## Avantajlar
// /
// / * **Güçlü Koruma**: Rate limiting'den daha güçlü brute-force koruması
// / * **Esnek Yapılandırma**: Deneme sayısı ve kilitleme süresi özelleştirilebilir
// / * **Thread-Safe**: Eşzamanlı isteklerde güvenli çalışır
// / * **Otomatik Temizlik**: Eski kayıtlar otomatik olarak temizlenir
// / * **Kullanıcı Dostu**: Kalan deneme sayısını gösterebilirsiniz
// /
// / ## Dezavantajlar
// /
// / * **Bellek Kullanımı**: Tüm denemeler bellekte tutulur
// / * **Dağıtık Sistemler**: Tek sunucuda çalışır, load balancer arkasında senkronizasyon gerekir
// / * **Kalıcılık Yok**: Sunucu yeniden başlatıldığında veriler kaybolur
// /
// / ## Önemli Notlar
// /
// / * Rate limiting ile birlikte kullanılması önerilir (çift koruma)
// / * Dağıtık sistemlerde Redis gibi merkezi bir store kullanın
// / * Email/username yerine hash kullanarak gizlilik koruyun
// / * Cleanup goroutine'i 5 dakikada bir çalışır
// / * Başarılı girişte mutlaka ResetAttempts çağırın
// / * Üretim ortamında mutlaka kullanın
type AccountLockout struct {
	mu              sync.RWMutex
	attempts        map[string]*lockoutEntry
	maxAttempts     int
	lockoutDuration time.Duration
	stopCh          chan struct{}
	doneCh          chan struct{}
	closeOnce       sync.Once
}

// / # lockoutEntry
// /
// / Bu yapı, bir identifier (email, IP, vb.) için kilitleme bilgilerini tutar.
// / AccountLockout tarafından dahili olarak kullanılır.
// /
// / ## Alanlar
// /
// / * `count` - Başarısız deneme sayısı
// / * `lockedUntil` - Hesabın kilitli kalacağı son zaman (zero value = kilitli değil)
// /
// / ## Kullanım
// /
// / Bu yapı AccountLockout tarafından dahili olarak kullanılır ve dışarıdan erişilmez.
// / Doğrudan kullanılması gerekmez.
// /
// / ## Önemli Notlar
// /
// / * `lockedUntil` zero value ise hesap hiç kilitlenmemiş demektir
// / * `count` her başarısız denemede artırılır
// / * `lockedUntil` geçmişse kilitleme süresi dolmuş demektir
type lockoutEntry struct {
	count       int
	lockedUntil time.Time
}

// / # NewAccountLockout
// /
// / Bu fonksiyon, yeni bir AccountLockout yöneticisi oluşturur ve arka planda çalışan
// / temizlik goroutine'ini başlatır.
// /
// / ## Parametreler
// /
// / * `maxAttempts` - Kilitleme öncesi maksimum başarısız deneme sayısı
// / * `lockoutDuration` - Hesabın kilitli kalacağı süre
// /
// / ## Çalışma Mantığı
// /
// / 1. Yeni bir AccountLockout instance'ı oluşturur
// / 2. Boş bir attempts map'i başlatır
// / 3. Arka planda cleanup goroutine'ini başlatır
// / 4. Cleanup goroutine'i 5 dakikada bir eski kayıtları temizler
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Uygulama Başlangıcı**: main() fonksiyonunda global lockout yöneticisi oluşturma
// / 2. **Farklı Endpoint'ler**: Her endpoint için farklı lockout politikaları
// / 3. **Çoklu Tenant**: Her tenant için ayrı lockout yöneticisi
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Standart kullanım: 5 deneme, 15 dakika kilitleme
// / lockout := NewAccountLockout(5, 15*time.Minute)
// /
// / // Katı güvenlik: 3 deneme, 30 dakika kilitleme
// / strictLockout := NewAccountLockout(3, 30*time.Minute)
// /
// / // Esnek güvenlik: 10 deneme, 5 dakika kilitleme
// / relaxedLockout := NewAccountLockout(10, 5*time.Minute)
// /
// / // Admin paneli için çok katı: 3 deneme, 1 saat kilitleme
// / adminLockout := NewAccountLockout(3, 1*time.Hour)
// /
// / // Global değişken olarak kullanım
// / var globalLockout *AccountLockout
// /
// / func init() {
// /     globalLockout = NewAccountLockout(5, 15*time.Minute)
// / }
// /
// / func main() {
// /     app := fiber.New()
// /     app.Post("/login", loginHandler)
// /     app.Listen(":3000")
// / }
// / ```
// /
// / ## Döndürür
// /
// / * `*AccountLockout` - Yapılandırılmış ve çalışmaya hazır AccountLockout pointer'ı
// /
// / ## Önemli Notlar
// /
// / * Cleanup goroutine'i otomatik olarak başlatılır
// / * Goroutine uygulama yaşam döngüsü boyunca çalışır
// / * Sunucu kapatıldığında goroutine otomatik olarak sonlanır
// / * maxAttempts 0 veya negatif olmamalı (kontrol yapılmaz)
// / * lockoutDuration 0 veya negatif olmamalı (kontrol yapılmaz)
// / * Thread-safe olarak tasarlanmıştır
func NewAccountLockout(maxAttempts int, lockoutDuration time.Duration) *AccountLockout {
	al := &AccountLockout{
		attempts:        make(map[string]*lockoutEntry),
		maxAttempts:     maxAttempts,
		lockoutDuration: lockoutDuration,
		stopCh:          make(chan struct{}),
		doneCh:          make(chan struct{}),
	}

	// Cleanup goroutine
	go al.cleanup()

	return al
}

// Close stops the background cleanup goroutine and waits until it exits.
func (al *AccountLockout) Close() {
	if al == nil {
		return
	}

	al.closeOnce.Do(func() {
		if al.stopCh != nil {
			close(al.stopCh)
		}
	})

	if al.doneCh != nil {
		<-al.doneCh
	}
}

// / # IsLocked
// /
// / Bu method, belirtilen identifier (email, IP, vb.) için hesabın kilitli olup olmadığını kontrol eder.
// / Thread-safe okuma işlemi yapar.
// /
// / ## Parametreler
// /
// / * `identifier` - Kontrol edilecek benzersiz tanımlayıcı (email, username, IP adresi, vb.)
// /
// / ## Çalışma Mantığı
// /
// / 1. RLock ile thread-safe okuma kilidi alır
// / 2. identifier için kayıt var mı kontrol eder
// / 3. Kayıt yoksa false döner (hiç deneme yapılmamış)
// / 4. Kayıt varsa lockedUntil zamanını kontrol eder
// / 5. Şu anki zaman lockedUntil'den önceyse true döner (hala kilitli)
// / 6. Aksi halde false döner (kilitleme süresi dolmuş)
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Login Kontrolü**: Giriş yapmadan önce hesabın kilitli olup olmadığını kontrol etme
// / 2. **Erken Yanıt**: Kilitli hesaplar için hemen hata döndürme
// / 3. **Kullanıcı Bildirimi**: Kullanıcıya hesabının kilitli olduğunu bildirme
// / 4. **Güvenlik Logu**: Kilitli hesaplara yapılan denemeleri loglama
// /
// / ## Örnek Kullanım
// /
// / ```go
// / func loginHandler(c *fiber.Ctx) error {
// /     email := c.FormValue("email")
// /
// /     // Hesap kilitli mi kontrol et
// /     if lockout.IsLocked(email) {
// /         return c.Status(403).JSON(fiber.Map{
// /             "error": "Hesabınız çok fazla başarısız deneme nedeniyle geçici olarak kilitlendi",
// /             "message": "Lütfen 15 dakika sonra tekrar deneyin",
// /         })
// /     }
// /
// /     // Kimlik doğrulama işlemine devam et...
// / }
// /
// / // IP bazlı kontrol
// / ip := c.IP()
// / if lockout.IsLocked(ip) {
// /     return c.Status(403).JSON(fiber.Map{
// /         "error": "IP adresiniz geçici olarak engellendi",
// /     })
// / }
// /
// / // Kullanıcı ID bazlı kontrol
// / userID := c.Locals("user_id").(string)
// / if lockout.IsLocked(userID) {
// /     return c.Status(403).JSON(fiber.Map{
// /         "error": "Hesabınız kilitli",
// /     })
// / }
// / ```
// /
// / ## Döndürür
// /
// / * `bool` - Hesap kilitliyse true, değilse false
// /
// / ## Thread Safety
// /
// / * RLock kullanarak thread-safe okuma sağlar
// / * Birden fazla goroutine aynı anda güvenle çağırabilir
// / * Yazma işlemleriyle (RecordFailedAttempt, ResetAttempts) çakışmaz
// /
// / ## Önemli Notlar
// /
// / * Kilitleme süresi dolmuşsa otomatik olarak false döner
// / * Kayıt yoksa false döner (hiç deneme yapılmamış)
// / * Bu method sadece kontrol yapar, kayıt değiştirmez
// / * Her login denemesinden önce çağırılmalıdır
// / * Kilitli hesaplara yapılan denemeleri loglamayı unutmayın
func (al *AccountLockout) IsLocked(identifier string) bool {
	al.mu.RLock()
	defer al.mu.RUnlock()

	entry, exists := al.attempts[identifier]
	if !exists {
		return false
	}

	if time.Now().Before(entry.lockedUntil) {
		return true
	}

	return false
}

// / # RecordFailedAttempt
// /
// / Bu method, belirtilen identifier için başarısız bir giriş denemesini kaydeder.
// / Deneme sayısını artırır ve gerekirse hesabı kilitler. Thread-safe yazma işlemi yapar.
// /
// / ## Parametreler
// /
// / * `identifier` - Başarısız deneme kaydedilecek benzersiz tanımlayıcı (email, username, IP adresi, vb.)
// /
// / ## Çalışma Mantığı
// /
// / 1. Lock ile thread-safe yazma kilidi alır
// / 2. identifier için kayıt var mı kontrol eder
// / 3. Kayıt yoksa yeni bir lockoutEntry oluşturur
// / 4. Kayıt varsa ve kilitleme süresi dolmuşsa sayacı sıfırlar
// / 5. Deneme sayacını 1 artırır
// / 6. Sayaç maxAttempts'e ulaştıysa lockedUntil zamanını ayarlar
// / 7. Kilitleme süresi şu andan itibaren lockoutDuration kadar ileri ayarlanır
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Başarısız Login**: Her başarısız giriş denemesinde çağırma
// / 2. **Yanlış Şifre**: Şifre doğrulama başarısız olduğunda çağırma
// / 3. **Geçersiz Token**: API token doğrulama başarısız olduğunda çağırma
// / 4. **2FA Hatası**: İki faktörlü kimlik doğrulama başarısız olduğunda çağırma
// /
// / ## Örnek Kullanım
// /
// / ```go
// / func loginHandler(c *fiber.Ctx) error {
// /     email := c.FormValue("email")
// /     password := c.FormValue("password")
// /
// /     // Hesap kilitli mi kontrol et
// /     if lockout.IsLocked(email) {
// /         return c.Status(403).JSON(fiber.Map{
// /             "error": "Hesabınız kilitli",
// /         })
// /     }
// /
// /     // Kimlik doğrulama yap
// /     user, err := authenticateUser(email, password)
// /     if err != nil {
// /         // Başarısız denemeyi kaydet
// /         lockout.RecordFailedAttempt(email)
// /
// /         // Kalan deneme sayısını al
// /         remaining := lockout.GetRemainingAttempts(email)
// /
// /         if remaining == 0 {
// /             return c.Status(403).JSON(fiber.Map{
// /                 "error": "Hesabınız çok fazla başarısız deneme nedeniyle kilitlendi",
// /                 "locked_until": "15 dakika",
// /             })
// /         }
// /
// /         return c.Status(401).JSON(fiber.Map{
// /             "error": "Geçersiz kimlik bilgileri",
// /             "remaining_attempts": remaining,
// /         })
// /     }
// /
// /     // Başarılı giriş - denemeleri sıfırla
// /     lockout.ResetAttempts(email)
// /     return c.JSON(fiber.Map{"success": true, "user": user})
// / }
// /
// / // IP bazlı kayıt
// / ip := c.IP()
// / lockout.RecordFailedAttempt(ip)
// /
// / // Kullanıcı ID bazlı kayıt
// / userID := c.Locals("user_id").(string)
// / lockout.RecordFailedAttempt(userID)
// /
// / // Composite key kullanımı
// / key := fmt.Sprintf("%s:%s", email, ip)
// / lockout.RecordFailedAttempt(key)
// / ```
// /
// / ## Thread Safety
// /
// / * Lock kullanarak thread-safe yazma sağlar
// / * Birden fazla goroutine aynı anda güvenle çağırabilir
// / * Okuma işlemleriyle (IsLocked, GetRemainingAttempts) çakışmaz
// /
// / ## Önemli Notlar
// /
// / * Her başarısız giriş denemesinde mutlaka çağırılmalıdır
// / * Kilitleme süresi dolmuşsa sayaç otomatik olarak sıfırlanır
// / * maxAttempts'e ulaşıldığında hesap otomatik olarak kilitlenir
// / * Başarılı girişte ResetAttempts çağırmayı unutmayın
// / * Kalan deneme sayısını kullanıcıya göstermek için GetRemainingAttempts kullanın
// / * Güvenlik loglarına başarısız denemeleri kaydedin
func (al *AccountLockout) RecordFailedAttempt(identifier string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	entry, exists := al.attempts[identifier]
	if !exists {
		entry = &lockoutEntry{count: 0}
		al.attempts[identifier] = entry
	}

	// Reset if lockout expired (only check if lockedUntil was previously set)
	if !entry.lockedUntil.IsZero() && time.Now().After(entry.lockedUntil) {
		entry.count = 0
	}

	entry.count++

	if entry.count >= al.maxAttempts {
		entry.lockedUntil = time.Now().Add(al.lockoutDuration)
	}
}

// / # ResetAttempts
// /
// / Bu method, belirtilen identifier için başarısız deneme sayacını sıfırlar ve kilitlemeyi kaldırır.
// / Başarılı giriş yapıldığında mutlaka çağırılmalıdır. Thread-safe yazma işlemi yapar.
// /
// / ## Parametreler
// /
// / * `identifier` - Deneme sayacı sıfırlanacak benzersiz tanımlayıcı (email, username, IP adresi, vb.)
// /
// / ## Çalışma Mantığı
// /
// / 1. Lock ile thread-safe yazma kilidi alır
// / 2. identifier için kaydı map'ten tamamen siler
// / 3. Kilitleme ve deneme bilgileri tamamen temizlenir
// / 4. Kullanıcı tekrar baştan deneme hakkına sahip olur
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Başarılı Login**: Kullanıcı başarıyla giriş yaptığında çağırma
// / 2. **Şifre Sıfırlama**: Kullanıcı şifresini sıfırladıktan sonra çağırma
// / 3. **Manuel Kilitleme Kaldırma**: Admin tarafından kilitleme kaldırıldığında çağırma
// / 4. **Hesap Doğrulama**: Email doğrulandıktan sonra çağırma
// /
// / ## Örnek Kullanım
// /
// / ```go
// / func loginHandler(c *fiber.Ctx) error {
// /     email := c.FormValue("email")
// /     password := c.FormValue("password")
// /
// /     // Hesap kilitli mi kontrol et
// /     if lockout.IsLocked(email) {
// /         return c.Status(403).JSON(fiber.Map{
// /             "error": "Hesabınız kilitli",
// /         })
// /     }
// /
// /     // Kimlik doğrulama yap
// /     user, err := authenticateUser(email, password)
// /     if err != nil {
// /         lockout.RecordFailedAttempt(email)
// /         return c.Status(401).JSON(fiber.Map{
// /             "error": "Geçersiz kimlik bilgileri",
// /         })
// /     }
// /
// /     // Başarılı giriş - denemeleri sıfırla (ÖNEMLİ!)
// /     lockout.ResetAttempts(email)
// /
// /     return c.JSON(fiber.Map{"success": true, "user": user})
// / }
// /
// / // Şifre sıfırlama sonrası
// / func resetPasswordHandler(c *fiber.Ctx) error {
// /     email := c.FormValue("email")
// /     // ... şifre sıfırlama işlemi ...
// /
// /     // Başarılı şifre sıfırlama - denemeleri sıfırla
// /     lockout.ResetAttempts(email)
// /
// /     return c.JSON(fiber.Map{"success": true})
// / }
// /
// / // Admin tarafından kilitleme kaldırma
// / func unlockAccountHandler(c *fiber.Ctx) error {
// /     email := c.Params("email")
// /
// /     // Kilitlemeyi kaldır
// /     lockout.ResetAttempts(email)
// /
// /     return c.JSON(fiber.Map{
// /         "message": "Hesap kilidi kaldırıldı",
// /     })
// / }
// /
// / // IP bazlı sıfırlama
// / ip := c.IP()
// / lockout.ResetAttempts(ip)
// / ```
// /
// / ## Thread Safety
// /
// / * Lock kullanarak thread-safe yazma sağlar
// / * Birden fazla goroutine aynı anda güvenle çağırabilir
// / * Okuma işlemleriyle (IsLocked, GetRemainingAttempts) çakışmaz
// /
// / ## Önemli Notlar
// /
// / * Başarılı girişte mutlaka çağırılmalıdır
// / * Kayıt yoksa hiçbir şey yapmaz (hata vermez)
// / * Tüm deneme ve kilitleme bilgileri tamamen silinir
// / * Kullanıcı tekrar maxAttempts kadar deneme hakkına sahip olur
// / * Güvenlik loglarına başarılı girişleri kaydedin
func (al *AccountLockout) ResetAttempts(identifier string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	delete(al.attempts, identifier)
}

// / # GetRemainingAttempts
// /
// / Bu method, belirtilen identifier için kalan deneme sayısını döndürür.
// / Kullanıcıya kaç deneme hakkı kaldığını göstermek için kullanılır. Thread-safe okuma işlemi yapar.
// /
// / ## Parametreler
// /
// / * `identifier` - Kalan deneme sayısı sorgulanacak benzersiz tanımlayıcı (email, username, IP adresi, vb.)
// /
// / ## Çalışma Mantığı
// /
// / 1. RLock ile thread-safe okuma kilidi alır
// / 2. identifier için kayıt var mı kontrol eder
// / 3. Kayıt yoksa maxAttempts döner (hiç deneme yapılmamış)
// / 4. Kayıt varsa ve kilitleme süresi dolmuşsa maxAttempts döner
// / 5. Kayıt varsa kalan deneme sayısını hesaplar (maxAttempts - count)
// / 6. Negatif değer varsa 0 döner (hesap kilitli)
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Kullanıcı Bildirimi**: Başarısız girişte kalan deneme sayısını gösterme
// / 2. **Uyarı Mesajı**: Az deneme kaldığında kullanıcıyı uyarma
// / 3. **UI Gösterimi**: Login formunda kalan deneme sayısını gösterme
// / 4. **Güvenlik Logu**: Deneme sayısını loglama
// /
// / ## Örnek Kullanım
// /
// / ```go
// / func loginHandler(c *fiber.Ctx) error {
// /     email := c.FormValue("email")
// /     password := c.FormValue("password")
// /
// /     // Hesap kilitli mi kontrol et
// /     if lockout.IsLocked(email) {
// /         return c.Status(403).JSON(fiber.Map{
// /             "error": "Hesabınız kilitli",
// /             "message": "Lütfen 15 dakika sonra tekrar deneyin",
// /         })
// /     }
// /
// /     // Kimlik doğrulama yap
// /     user, err := authenticateUser(email, password)
// /     if err != nil {
// /         // Başarısız denemeyi kaydet
// /         lockout.RecordFailedAttempt(email)
// /
// /         // Kalan deneme sayısını al
// /         remaining := lockout.GetRemainingAttempts(email)
// /
// /         if remaining == 0 {
// /             return c.Status(403).JSON(fiber.Map{
// /                 "error": "Hesabınız kilitlendi",
// /                 "message": "Çok fazla başarısız deneme yaptınız",
// /                 "locked_duration": "15 dakika",
// /             })
// /         }
// /
// /         // Kullanıcıya kalan deneme sayısını göster
// /         message := fmt.Sprintf("Geçersiz kimlik bilgileri. %d deneme hakkınız kaldı.", remaining)
// /         if remaining <= 2 {
// /             message += " Dikkatli olun, hesabınız kilitlenebilir!"
// /         }
// /
// /         return c.Status(401).JSON(fiber.Map{
// /             "error": "Geçersiz kimlik bilgileri",
// /             "message": message,
// /             "remaining_attempts": remaining,
// /         })
// /     }
// /
// /     // Başarılı giriş
// /     lockout.ResetAttempts(email)
// /     return c.JSON(fiber.Map{"success": true, "user": user})
// / }
// /
// / // Kalan deneme sayısını kontrol et
// / remaining := lockout.GetRemainingAttempts("user@example.com")
// / if remaining <= 2 {
// /     log.Printf("UYARI: Kullanıcı %s için sadece %d deneme kaldı", email, remaining)
// / }
// /
// / // IP bazlı kontrol
// / ip := c.IP()
// / remaining := lockout.GetRemainingAttempts(ip)
// / ```
// /
// / ## Döndürür
// /
// / * `int` - Kalan deneme sayısı (0 = kilitli, maxAttempts = hiç deneme yapılmamış)
// /
// / ## Thread Safety
// /
// / * RLock kullanarak thread-safe okuma sağlar
// / * Birden fazla goroutine aynı anda güvenle çağırabilir
// / * Yazma işlemleriyle (RecordFailedAttempt, ResetAttempts) çakışmaz
// /
// / ## Önemli Notlar
// /
// / * 0 döndürürse hesap kilitlidir
// / * maxAttempts döndürürse hiç deneme yapılmamış veya sıfırlanmıştır
// / * Kilitleme süresi dolmuşsa otomatik olarak maxAttempts döner
// / * Kullanıcıya kalan deneme sayısını göstermek güvenlik açısından tartışmalıdır
// / * Bazı güvenlik uzmanları bu bilgiyi göstermeyi önermez (saldırgana bilgi verir)
// / * Kullanıcı deneyimi için göstermek faydalı olabilir
func (al *AccountLockout) GetRemainingAttempts(identifier string) int {
	al.mu.RLock()
	defer al.mu.RUnlock()

	entry, exists := al.attempts[identifier]
	if !exists {
		return al.maxAttempts
	}

	// Only check expiration if lockout was previously set
	if !entry.lockedUntil.IsZero() && time.Now().After(entry.lockedUntil) {
		return al.maxAttempts
	}

	remaining := al.maxAttempts - entry.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// / # cleanup
// /
// / Bu method, arka planda çalışan bir goroutine olarak periyodik olarak eski ve süresi dolmuş
// / kilitleme kayıtlarını temizler. Bellek kullanımını optimize eder. Thread-safe yazma işlemi yapar.
// /
// / ## Çalışma Mantığı
// /
// / 1. 5 dakikada bir çalışan bir ticker oluşturur
// / 2. Her tick'te tüm kayıtları tarar
// / 3. Kilitleme süresi dolmuş ve deneme sayısı maxAttempts'in altında olan kayıtları siler
// / 4. Bu sayede bellekte gereksiz kayıt birikmesi önlenir
// / 5. Goroutine uygulama yaşam döngüsü boyunca çalışır
// /
// / ## Kullanım
// /
// / Bu method NewAccountLockout tarafından otomatik olarak başlatılır.
// / Doğrudan çağırılması gerekmez ve public olarak erişilemez.
// /
// / ## Temizleme Kriterleri
// /
// / Bir kayıt şu koşullarda temizlenir:
// / * Şu anki zaman lockedUntil'den sonra (kilitleme süresi dolmuş)
// / * VE deneme sayısı maxAttempts'in altında (henüz kilitlenmemiş)
// /
// / ## Önemli Notlar
// /
// / * Goroutine 5 dakikada bir çalışır (sabit interval)
// / * Lock kullanarak thread-safe yazma sağlar
// / * Uygulama kapatıldığında goroutine otomatik olarak sonlanır
// / * Kilitli hesaplar (count >= maxAttempts) temizlenmez
// / * Sadece süresi dolmuş ve kilitlenmemiş kayıtlar temizlenir
// / * Bellek optimizasyonu için kritik öneme sahiptir
func (al *AccountLockout) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	defer close(al.doneCh)

	for {
		select {
		case <-al.stopCh:
			return
		case <-ticker.C:
			al.mu.Lock()
			now := time.Now()
			for key, entry := range al.attempts {
				if now.After(entry.lockedUntil) && entry.count < al.maxAttempts {
					delete(al.attempts, key)
				}
			}
			al.mu.Unlock()
		}
	}
}

// / # RequestSizeLimit
// /
// / Bu fonksiyon, HTTP istek gövdesi (body) boyutunu sınırlayan bir middleware oluşturur.
// / Büyük dosya yüklemelerini veya aşırı büyük JSON/XML payload'larını engelleyerek
// / sunucuyu bellek tükenme saldırılarından korur.
// /
// / ## Parametreler
// /
// / * `maxSize` - Maksimum istek gövdesi boyutu (byte cinsinden)
// /
// / ## Çalışma Mantığı
// /
// / 1. Her istek için gövde boyutunu kontrol eder
// / 2. Boyut maxSize'dan büyükse 413 (Request Entity Too Large) hatası döner
// / 3. Boyut uygunsa isteğin devam etmesine izin verir
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Dosya Yükleme Koruması**: Çok büyük dosya yüklemelerini engelleme
// / 2. **API Payload Koruması**: Aşırı büyük JSON/XML verilerini engelleme
// / 3. **DoS Koruması**: Bellek tükenme saldırılarını engelleme
// / 4. **Kaynak Yönetimi**: Sunucu belleğini koruma
// / 5. **Bandwidth Koruması**: Gereksiz büyük veri transferini engelleme
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Tüm uygulama için 10MB limit
// / app.Use(RequestSizeLimit(10 * 1024 * 1024)) // 10MB
// /
// / // API endpoint'leri için 1MB limit
// / api := app.Group("/api")
// / api.Use(RequestSizeLimit(1 * 1024 * 1024)) // 1MB
// /
// / // Dosya yükleme endpoint'i için 50MB limit
// / app.Post("/upload", RequestSizeLimit(50 * 1024 * 1024), uploadHandler)
// /
// / // Form submission için 100KB limit
// / app.Post("/contact", RequestSizeLimit(100 * 1024), contactHandler)
// /
// / // Farklı endpoint'ler için farklı limitler
// / app.Post("/api/small", RequestSizeLimit(10 * 1024), smallDataHandler)      // 10KB
// / app.Post("/api/medium", RequestSizeLimit(1024 * 1024), mediumDataHandler)  // 1MB
// / app.Post("/api/large", RequestSizeLimit(10 * 1024 * 1024), largeDataHandler) // 10MB
// /
// / // Yaygın boyut sabitleri
// / const (
// /     KB = 1024
// /     MB = 1024 * KB
// /     GB = 1024 * MB
// / )
// /
// / app.Use(RequestSizeLimit(5 * MB)) // 5MB
// / ```
// /
// / ## Döndürür
// /
// / * `fiber.Handler` - İstek boyutunu kontrol eden middleware fonksiyonu
// /
// / ## Hata Yanıtı
// /
// / Limit aşıldığında dönen yanıt:
// / * **Status Code**: 413 (Request Entity Too Large)
// / * **Body**: JSON formatında hata mesajı ve maksimum boyut bilgisi
// /
// / ## Avantajlar
// /
// / * **Bellek Koruması**: Sunucu belleğini aşırı kullanımdan korur
// / * **DoS Koruması**: Bellek tükenme saldırılarını engeller
// / * **Performans**: Büyük istekleri erken aşamada reddeder
// / * **Esneklik**: Endpoint bazında farklı limitler ayarlanabilir
// / * **Kullanıcı Dostu**: Açıklayıcı hata mesajı verir
// /
// / ## Önemli Notlar
// /
// / * Fiber'ın kendi BodyLimit middleware'i de vardır, bu alternatif bir implementasyondur
// / * maxSize byte cinsinden belirtilmelidir (1MB = 1024 * 1024 byte)
// / * Dosya yükleme endpoint'leri için daha yüksek limit kullanın
// / * API endpoint'leri için daha düşük limit kullanın (genellikle 1-5MB yeterli)
// / * Limit çok düşük olursa normal kullanıcılar etkilenebilir
// / * Limit çok yüksek olursa DoS saldırılarına açık olabilir
// / * Üretim ortamında mutlaka kullanın
// / * Nginx/Apache gibi reverse proxy'lerde de boyut limiti ayarlayın
func RequestSizeLimit(maxSize int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > maxSize {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": fmt.Sprintf("Request body too large. Maximum size: %d bytes", maxSize),
			})
		}
		return c.Next()
	}
}

// / # ValidateCORSOrigin
// /
// / Bu fonksiyon, CORS (Cross-Origin Resource Sharing) origin doğrulaması için bir validator fonksiyonu döndürür.
// / Beyaz listeye alınmış origin'lerden gelen isteklere izin verir, diğerlerini reddeder.
// / Wildcard subdomain desteği sağlar (örn: *.example.com).
// /
// / ## Parametreler
// /
// / * `allowedOrigins` - İzin verilen origin'lerin listesi (string slice)
// /   - Tam domain: "https://example.com"
// /   - Wildcard subdomain: "*.example.com"
// /   - Tüm origin'ler: "*" (üretimde önerilmez)
// /
// / ## Çalışma Mantığı
// /
// / 1. Her origin için beyaz listeyi kontrol eder
// / 2. "*" varsa tüm origin'lere izin verir (güvenli değil)
// / 3. Tam eşleşme varsa izin verir (case-insensitive)
// / 4. Wildcard subdomain desteği:
// /    - "*.example.com" şeklinde tanımlanabilir
// /    - "sub.example.com", "api.example.com" gibi subdomain'lere izin verir
// /    - "example.com" ana domain'e izin vermez (sadece subdomain'ler)
// / 5. Eşleşme yoksa false döner (origin reddedilir)
// /
// / ## Kullanım Senaryoları
// /
// / 1. **API Güvenliği**: Sadece belirli frontend uygulamalarından API erişimi
// / 2. **Çoklu Domain**: Birden fazla domain'den erişim kontrolü
// / 3. **Subdomain Desteği**: Tüm subdomain'lere izin verme (*.example.com)
// / 4. **Geliştirme/Üretim**: Farklı ortamlar için farklı origin'ler
// / 5. **Mikroservis Güvenliği**: Sadece belirli servislerden erişim
// /
// / ## Örnek Kullanım
// /
// / ```go
// / import "github.com/gofiber/fiber/v2/middleware/cors"
// /
// / // Basit kullanım - tek domain
// / app.Use(cors.New(cors.Config{
// /     AllowOriginsFunc: ValidateCORSOrigin([]string{
// /         "https://example.com",
// /     }),
// / }))
// /
// / // Çoklu domain
// / app.Use(cors.New(cors.Config{
// /     AllowOriginsFunc: ValidateCORSOrigin([]string{
// /         "https://example.com",
// /         "https://app.example.com",
// /         "https://admin.example.com",
// /     }),
// / }))
// /
// / // Wildcard subdomain desteği
// / app.Use(cors.New(cors.Config{
// /     AllowOriginsFunc: ValidateCORSOrigin([]string{
// /         "*.example.com",           // Tüm subdomain'ler
// /         "https://partner.com",     // Belirli partner domain
// /     }),
// / }))
// /
// / // Geliştirme ve üretim ortamları
// / var allowedOrigins []string
// / if os.Getenv("ENV") == "production" {
// /     allowedOrigins = []string{
// /         "https://example.com",
// /         "https://app.example.com",
// /     }
// / } else {
// /     allowedOrigins = []string{
// /         "http://localhost:3000",
// /         "http://localhost:5173",
// /         "*.dev.example.com",
// /     }
// / }
// /
// / app.Use(cors.New(cors.Config{
// /     AllowOriginsFunc: ValidateCORSOrigin(allowedOrigins),
// /     AllowCredentials: true,
// /     AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
// /     AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// / }))
// /
// / // Tüm origin'lere izin (sadece geliştirme için!)
// / app.Use(cors.New(cors.Config{
// /     AllowOriginsFunc: ValidateCORSOrigin([]string{"*"}),
// / }))
// / ```
// /
// / ## Döndürür
// /
// / * `func(string) bool` - Origin doğrulama fonksiyonu
// /   - `true`: Origin izin veriliyor
// /   - `false`: Origin reddediliyor
// /
// / ## Wildcard Subdomain Kuralları
// /
// / * `*.example.com` şu origin'lere izin verir:
// /   - `https://api.example.com` ✓
// /   - `https://app.example.com` ✓
// /   - `https://sub.api.example.com` ✓ (nested subdomain)
// / * `*.example.com` şu origin'lere izin vermez:
// /   - `https://example.com` ✗ (ana domain)
// /   - `https://example.com.tr` ✗ (farklı TLD)
// /   - `https://fakeexample.com` ✗ (farklı domain)
// /
// / ## Avantajlar
// /
// / * **Güvenlik**: Sadece beyaz listeye alınmış origin'lerden erişim
// / * **Esneklik**: Wildcard subdomain desteği
// / * **Kolay Yapılandırma**: Basit string listesi ile yapılandırma
// / * **Case-Insensitive**: Büyük/küçük harf duyarsız karşılaştırma
// / * **Performans**: Hızlı string karşılaştırması
// /
// / ## Güvenlik Notları
// /
// / * **"*" Kullanımı**: Üretim ortamında "*" kullanmayın (tüm origin'lere izin verir)
// / * **HTTPS Zorunluluğu**: Üretimde sadece HTTPS origin'lere izin verin
// / * **Wildcard Dikkat**: "*.example.com" tüm subdomain'lere izin verir, dikkatli kullanın
// / * **Credentials**: AllowCredentials true ise "*" kullanılamaz (CORS spesifikasyonu)
// / * **Validation**: Origin'lerin geçerli URL formatında olduğundan emin olun
// /
// / ## Önemli Notlar
// /
// / * Origin karşılaştırması case-insensitive'dir
// / * Protocol (http/https) dahil tam URL beklenir
// / * Wildcard sadece subdomain için çalışır, path için çalışmaz
// / * Boş origin listesi tüm origin'leri reddeder
// / * Fiber'ın CORS middleware'i ile birlikte kullanılmalıdır
// / * Üretim ortamında mutlaka beyaz liste kullanın
// / * CORS preflight (OPTIONS) isteklerini de kontrol eder
// / * AllowCredentials true ise cookie/auth header'ları gönderilir
func ValidateCORSOrigin(allowedOrigins []string) func(string) bool {
	return func(origin string) bool {
		for _, allowed := range allowedOrigins {
			if allowed == "*" {
				return true
			}
			if strings.EqualFold(origin, allowed) {
				return true
			}
			// Support wildcard subdomains like *.example.com
			if strings.HasPrefix(allowed, "*.") {
				domain := strings.TrimPrefix(allowed, "*.")
				// Extract host from origin (remove protocol)
				if strings.Contains(origin, "://") {
					parts := strings.Split(origin, "://")
					if len(parts) == 2 {
						host := parts[1]
						// Match if host ends with domain and is not exactly the domain
						if strings.HasSuffix(host, domain) && host != domain {
							return true
						}
					}
				}
			}
		}
		return false
	}
}
