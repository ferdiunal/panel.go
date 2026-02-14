/// # Panel.go - Ana Uygulama Paketi
///
/// Bu paket, Panel.go yönetim panelinin temel yapısını ve işlevselliğini sağlar.
/// Veritabanı bağlantısı, HTTP sunucusu, kimlik doğrulama, kaynaklar ve sayfaları yönetir.
///
/// ## Temel Özellikler
/// - Fiber web framework entegrasyonu
/// - GORM ORM desteği
/// - Kimlik doğrulama ve oturum yönetimi
/// - Dinamik kaynak (Resource) ve sayfa (Page) yönetimi
/// - Güvenlik middleware'leri (CORS, CSRF, Helmet, vb.)
/// - Bildirim sistemi
/// - İzin yönetimi
///
/// ## Kullanım Örneği
/// ```go
/// config := panel.Config{
///     Database: DatabaseConfig{Instance: db},
///     Server: ServerConfig{Host: "localhost", Port: "8080"},
///     Environment: "development",
/// }
/// p := panel.New(config)
/// p.Start()
/// ```

package panel

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/data/orm"
	"github.com/ferdiunal/panel.go/pkg/domain/account"
	notificationDomain "github.com/ferdiunal/panel.go/pkg/domain/notification"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/setting"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/domain/verification"
	"github.com/ferdiunal/panel.go/pkg/handler"
	authHandler "github.com/ferdiunal/panel.go/pkg/handler/auth"
	"github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/ferdiunal/panel.go/pkg/notification"
	"github.com/ferdiunal/panel.go/pkg/openapi"
	"github.com/ferdiunal/panel.go/pkg/page"
	"github.com/ferdiunal/panel.go/pkg/permission"
	"github.com/ferdiunal/panel.go/pkg/plugin"
	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceUser "github.com/ferdiunal/panel.go/pkg/resource/user"
	"github.com/ferdiunal/panel.go/pkg/service/auth"
	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// / # Panel Yapısı
// /
// / Panel, Panel.go yönetim panelinin ana yapısıdır. Uygulamanın tüm bileşenlerini
// / (veritabanı, web sunucusu, kimlik doğrulama, kaynaklar ve sayfalar) yönetir.
// /
// / ## Alanlar
// / - `Config`: Uygulama yapılandırması (veritabanı, sunucu, ortam ayarları)
// / - `Db`: GORM veritabanı bağlantısı
// / - `Fiber`: Fiber web framework uygulaması
// / - `Auth`: Kimlik doğrulama servisi
// / - `resources`: Kayıtlı kaynakların haritası (slug -> Resource)
// / - `pages`: Kayıtlı sayfaların haritası (slug -> Page)
// /
// / ## Avantajlar
// / - Merkezi yönetim: Tüm bileşenler tek bir yapıda
// / - Kolay genişletme: Yeni kaynaklar ve sayfalar kolayca eklenebilir
// / - Güvenlik: Yerleşik güvenlik middleware'leri
// / - Esneklik: Özelleştirilebilir yapılandırma
// /
// / ## Uyarılar
// / - Panel örneği oluşturulduktan sonra değiştirilmemelidir
// / - Kaynaklar ve sayfalar Start() çağrılmadan önce kayıtlı olmalıdır
type Panel struct {
	Config             Config
	Db                 *gorm.DB
	Fiber              *fiber.App
	Auth               *auth.Service
	resources          map[string]resource.Resource
	pages              map[string]page.Page
	plugins            []interface{} // Registered plugins
	http2PushResources []string      // Cached list of critical assets for HTTP/2 push
	openAPIHandler     *handler.OpenAPIHandler
}

// / # New Fonksiyonu
// /
// / Yeni bir Panel örneği oluşturur ve başlatır. Veritabanı migration'ları,
// / middleware kayıtları, güvenlik ayarları ve API yönlendirmelerini yapılandırır.
// /
// / ## Parametreler
// / - `config`: Panel yapılandırması (veritabanı, sunucu, ortam, CORS, OAuth vb.)
// /
// / ## Dönüş Değeri
// / Tamamen yapılandırılmış ve başlatılmış Panel örneğine işaretçi
// /
// / ## Yapılandırılan Bileşenler
// / 1. **Kimlik Doğrulama**: Kullanıcı, oturum ve hesap yönetimi
// / 2. **Veritabanı**: Otomatik migration'lar
// / 3. **Middleware'ler**:
// /    - Sıkıştırma (Compression)
// /    - CORS (Cross-Origin Resource Sharing)
// /    - CSRF Koruması
// /    - Güvenlik Başlıkları (Helmet)
// /    - İstek Boyutu Sınırı
// /    - Denetim Günlüğü (Audit Logging)
// / 4. **Statik Dosyalar**: Gömülü veya yerel UI dosyaları
// / 5. **İzinler**: İzin dosyasından yükleme
// / 6. **Kaynaklar ve Sayfalar**: Varsayılan kaynaklar ve sayfalar
// / 7. **API Yönlendirmeleri**: Kimlik doğrulama, kaynaklar, sayfalar, bildirimler
// /
// / ## Güvenlik Özellikleri
// / - CORS: Yapılandırılabilir izin verilen kaynaklar
// / - CSRF: X-CSRF-Token başlığı ile korumalı
// / - Helmet: Güvenlik başlıkları
// / - Content Security Policy: Komut dosyası ve stil kaynakları sınırlandırılmış
// / - Rate Limiting: Kimlik doğrulama (10 req/min), API (100 req/min)
// / - Hesap Kilitleme: 5 başarısız denemeden sonra 15 dakika kilitleme
// / - İstek Boyutu Sınırı: 10MB
// /
// / ## Kullanım Örneği
// / ```go
// / config := panel.Config{
// /     Database: DatabaseConfig{Instance: db},
// /     Server: ServerConfig{Host: "localhost", Port: "8080"},
// /     Environment: "development",
// /     CORS: CORSConfig{AllowedOrigins: []string{"http://localhost:3000"}},
// / }
// / p := panel.New(config)
// / if err := p.Start(); err != nil {
// /     log.Fatal(err)
// / }
// / ```
// /
// / ## Uyarılar
// / - Üretim ortamında CORS ayarlarını dikkatli yapılandırın
// / - İzin dosyası yüklenemezse panic oluşturulur
// / - Veritabanı bağlantısı başarısız olursa panic oluşturulur
// / - Geliştirme ortamında yerel UI dosyaları kullanılabilir
// /
// / ## Önemli Notlar
// / - Tüm kaynaklar ve sayfalar Start() çağrılmadan önce kayıtlı olmalıdır
// / - Middleware sırası önemlidir ve değiştirilmemelidir
// / - CSRF koruması test ortamında devre dışı bırakılır

// / # langMiddleware Fonksiyonu
// /
// / URL parametresinden veya request'ten dil bilgisini alır ve c.Locals("lang") ile set eder.
// / Bu middleware, i18n URL prefix özelliği için kullanılır.
// /
// / ## Parametreler
// / - `config`: Panel yapılandırması (I18n ayarları için)
// /
// / ## Dönüş Değeri
// / - `fiber.Handler`: Fiber middleware handler
// /
// / ## Davranış
// / 1. URL parametresinden lang değerini alır (c.Params("lang"))
// / 2. Eğer lang yoksa, i18n.GetLocale(c) kullanır (fiberi18n middleware'inin set ettiği değer)
// / 3. Desteklenen diller arasında mı kontrol eder
// / 4. Desteklenmiyorsa varsayılan dili kullanır
// / 5. c.Locals("lang", lang) ile set eder
// /
// / ## Kullanım Örneği
// / ```go
// / api := app.Group("/api")
// / api.Use(langMiddleware(config))
// / ```
func langMiddleware(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.Params("lang", "")

		if lang == "" {
			// Varsayılan dil için query param > header > cookie > config
			if config.I18n.Enabled {
				lang = i18n.GetLocale(c)
			} else {
				lang = config.I18n.DefaultLanguage.String()
			}
		}

		// Desteklenen diller arasında mı kontrol et
		supported := false
		for _, acceptedLang := range config.I18n.AcceptLanguages {
			if lang == acceptedLang.String() {
				supported = true
				break
			}
		}

		if !supported {
			lang = config.I18n.DefaultLanguage.String()
		}

		c.Locals("lang", lang)
		return c.Next()
	}
}

func New(config Config) *Panel {
	// configRef, closure'ların her zaman güncel config'i okumasını sağlar.
	// Panel oluşturulduktan sonra p.Config'e yeniden atanır.
	configRef := &config

	// SECURITY: Configure Fiber with TrustProxy for production deployments behind reverse proxy
	// This is REQUIRED for earlydata middleware to work securely
	app := fiber.New(fiber.Config{
		// Enable trusted proxy check for production environments
		EnableTrustedProxyCheck: config.Environment == "production",
		// Trust common reverse proxy IPs (nginx, cloudflare, etc.)
		// In production, configure this based on your infrastructure
		TrustedProxies: []string{"127.0.0.1", "::1"},
		// Use X-Forwarded-For header to get real client IP
		ProxyHeader: fiber.HeaderXForwardedFor,
	})
	db := config.Database.Instance

	// Inject DB to Context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// Auth Components
	userRepo := orm.NewUserRepository(db)
	sessionRepo := orm.NewSessionRepository(db)
	accountRepo := orm.NewAccountRepository(db)

	authService := auth.NewService(userRepo, sessionRepo, accountRepo)
	// SECURITY: Account lockout after 5 failed attempts, 15 minute lockout duration
	accountLockout := middleware.NewAccountLockout(5, 15*time.Minute)
	authH := authHandler.NewHandler(authService, accountLockout, config.Environment)

	// Auto Migrate Auth Domains
	db.AutoMigrate(&user.User{}, &session.Session{}, &account.Account{}, &verification.Verification{}, &setting.Setting{}, &notificationDomain.Notification{})

	// Middleware Registration
	// SECURITY: EncryptCookie middleware - MUST be registered BEFORE other cookie middleware
	// Encrypts cookie values using AES-GCM (cookie names remain unencrypted)
	// CRITICAL: Requires stable encryption key (changing key makes existing cookies unreadable)
	// CSRF and session token cookies are excluded from encryption for compatibility
	if config.EncryptionCookie.Key != "" {
		app.Use(encryptcookie.New(config.EncryptionCookie))
	}

	// I18n Middleware
	if config.I18n.Enabled {
		app.Use(fiberi18n.New(&fiberi18n.Config{
			RootPath:         config.I18n.RootPath,
			AcceptLanguages:  config.I18n.AcceptLanguages,
			DefaultLanguage:  config.I18n.DefaultLanguage,
			FormatBundleFile: config.I18n.FormatBundleFile,
		}))
		app.Use(langMiddleware(config))
	}

	// PERFORMANCE: Compress middleware with optimized settings for API responses
	// LevelBestSpeed prioritizes latency over compression ratio (ideal for APIs)
	// Bodies < 200 bytes are automatically skipped (compression would increase size)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Optimize for low latency
	}))

	// SECURITY: CORS Configuration - NEVER use "*" in production
	// Configure allowed origins in your config
	allowedOrigins := config.CORS.AllowedOrigins
	if len(allowedOrigins) == 0 {
		// Default to localhost for development
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,X-CSRF-Token",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600,
	}))

	// SECURITY: CSRF Protection - ALWAYS enabled (except in test environment)
	// Uses Double Submit Cookie pattern: token in cookie + header
	// Cookie name is consistent across environments for frontend compatibility
	if config.Environment != "test" {
		app.Use(csrf.New(csrf.Config{
			KeyLookup:      "header:X-CSRF-Token",              // Extract token from header
			CookieName:     "csrf_token",                       // Must match frontend axios xsrfCookieName
			CookieSecure:   config.Environment == "production", // HTTPS only in production
			CookieHTTPOnly: false,                              // CRITICAL: Must be false for SPA (JavaScript needs to read cookie)
			CookieSameSite: "Lax",                              // Lax allows GET requests from external sites (better UX than Strict)
			Expiration:     24 * time.Hour,
		}))
	}

	// SECURITY: EarlyData (TLS 1.3 0-RTT) middleware
	// CRITICAL: Only enabled in production with TrustProxy enabled
	// EarlyData allows replay attacks, so it's only safe behind a reverse proxy
	// Default behavior: Only allows safe HTTP methods (GET, HEAD, OPTIONS, TRACE)
	// See RFC 8446 Section 8 for security implications
	if config.Environment == "production" {
		app.Use(earlydata.New(earlydata.Config{
			Error: fiber.ErrTooEarly, // Return 425 Too Early for rejected requests
		}))
	}
	// Note: In development, earlydata is disabled for security
	// Enable TrustProxy and use a reverse proxy in production to use this feature

	app.Use(etag.New())

	// SECURITY: Enhanced security headers
	app.Use(helmet.New(helmet.Config{
		CrossOriginResourcePolicy: "cross-origin",
	}))

	// SECURITY: Additional security headers (helmet doesn't support all of these)
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none';")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		return c.Next()
	})

	// SECURITY: Request size limits to prevent DoS attacks
	app.Use(middleware.RequestSizeLimit(10 * 1024 * 1024)) // 10MB limit

	// SECURITY: Audit logging for security events
	auditLogger := &middleware.ConsoleAuditLogger{}
	app.Use(middleware.AuditMiddleware(auditLogger))

	// I18N: Çoklu dil desteği (Internationalization)
	// Uygulamanın farklı dillerde gösterilmesini sağlar
	// Dil seçimi: 1) Query param (?lang=tr), 2) Accept-Language header, 3) DefaultLanguage
	if config.I18n.Enabled {
		// Varsayılan değerleri ayarla
		rootPath := config.I18n.RootPath
		if rootPath == "" {
			rootPath = "./locales"
		}

		acceptLanguages := config.I18n.AcceptLanguages
		if len(acceptLanguages) == 0 {
			acceptLanguages = []language.Tag{language.Turkish, language.English}
		}

		defaultLanguage := config.I18n.DefaultLanguage
		if defaultLanguage == language.Und {
			defaultLanguage = language.Turkish
		}

		formatBundleFile := config.I18n.FormatBundleFile
		if formatBundleFile == "" {
			formatBundleFile = "yaml"
		}

		app.Use(fiberi18n.New(&fiberi18n.Config{
			RootPath:         rootPath,
			AcceptLanguages:  acceptLanguages,
			DefaultLanguage:  defaultLanguage,
			FormatBundleFile: formatBundleFile,
		}))
	}

	// Static file serving with priority:
	// 1. assets/ui/ (project-specific build from plugin:build)
	// 2. pkg/panel/ui/ (SDK development mode)
	// 3. Embedded assets (backward compatibility)
	//
	// HTML files are served with dynamic injection (RTL, theme, lang)
	// Static assets (JS, CSS, images) are served directly

	// PERFORMANCE: High compression for static assets
	staticCompress := compress.New(compress.Config{
		Level: compress.LevelBestCompression,
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api")
		},
	})

	// Storage serving (common for all priorities)
	if config.Storage.URL != "" && config.Storage.Path != "" {
		app.Use(config.Storage.URL, staticCompress)
		app.Static(config.Storage.URL, config.Storage.Path)
	} else {
		app.Use("/storage", staticCompress)
		app.Static("/storage", "./storage/public")
	}

	// Priority 1: Check assets/ui/index.html (project-specific build)
	assetsUIPath := "assets/ui/index.html"
	if _, err := os.Stat(assetsUIPath); err == nil {
		fmt.Println("✓ Serving UI from: assets/ui/ (project-specific build)")

		// Static assets (JS, CSS, images, fonts)
		app.Use("/assets", staticCompress)
		app.Static("/assets", "./assets/ui/assets")

		// HTML files with dynamic injection
		app.Get("/*", func(c *fiber.Ctx) error {
			// Skip API routes
			if strings.HasPrefix(c.Path(), "/api") {
				return c.Next()
			}

			// Serve static assets directly
			path := c.Path()
			if strings.HasPrefix(path, "/assets/") {
				return c.Next()
			}

			// Check if file exists and is not HTML
			if path != "/" && strings.Contains(path, ".") {
				filePath := filepath.Join("assets/ui", path)
				if _, err := os.Stat(filePath); err == nil {
					return c.SendFile(filePath)
				}
			}

			// Serve HTML with injection
			return ServeHTML(c, assetsUIPath, *configRef)
		})
	} else if config.Environment == "development" {
		// Priority 2: Check pkg/panel/ui/index.html (SDK development)
		localUIPath := "pkg/panel/ui/index.html"
		if _, err := os.Stat(localUIPath); err == nil {
			fmt.Println("✓ Serving UI from: pkg/panel/ui/ (SDK development)")

			// Static assets
			app.Use("/assets", staticCompress)
			app.Static("/assets", "./pkg/panel/ui/assets")

			// HTML files with dynamic injection
			app.Get("/*", func(c *fiber.Ctx) error {
				if strings.HasPrefix(c.Path(), "/api") {
					return c.Next()
				}

				path := c.Path()
				if strings.HasPrefix(path, "/assets/") {
					return c.Next()
				}

				if path != "/" && strings.Contains(path, ".") {
					filePath := filepath.Join("pkg/panel/ui", path)
					if _, err := os.Stat(filePath); err == nil {
						return c.SendFile(filePath)
					}
				}

				return ServeHTML(c, localUIPath, *configRef)
			})
		} else {
			// Priority 3: Embedded assets (backward compatibility)
			fmt.Println("✓ Serving UI from: embedded assets (backward compatibility)")
			assetsFS, err := GetFileSystem(true)
			if err != nil {
				fmt.Println("Warning: Failed to load embedded assets:", err)
			}

			if assetsFS != nil {
				// Static assets from embedded FS
				app.Use("/", staticCompress, filesystem.New(filesystem.Config{
					Root:   http.FS(assetsFS),
					Browse: false,
					Next: func(c *fiber.Ctx) bool {
						// Skip API routes and HTML files
						return strings.HasPrefix(c.Path(), "/api") ||
							c.Path() == "/" ||
							!strings.Contains(c.Path(), ".")
					},
				}))

				// HTML with injection (read from embedded FS)
				app.Get("/*", func(c *fiber.Ctx) error {
					if strings.HasPrefix(c.Path(), "/api") {
						return c.Next()
					}

					path := c.Path()
					if path != "/" && strings.Contains(path, ".") && !strings.HasSuffix(path, ".html") {
						return c.Next()
					}

					// Read index.html from embedded FS
					htmlBytes, err := fs.ReadFile(assetsFS, "index.html")
					if err != nil {
						return c.Status(500).SendString("Failed to load UI")
					}

					// Inject and serve
					data := GetHTMLInjectionData(c, *configRef)
					initData := GetInitData(c, *configRef, data)
					initJSON, _ := json.Marshal(initData)
					html := InjectHTML(string(htmlBytes), data, string(initJSON))

					c.Set("Content-Type", "text/html; charset=utf-8")
					c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
					return c.SendString(html)
				})
			}
		}
	} else {
		// Priority 3: Embedded assets (backward compatibility)
		fmt.Println("✓ Serving UI from: embedded assets (backward compatibility)")
		assetsFS, err := GetFileSystem(true)
		if err != nil {
			fmt.Println("Warning: Failed to load embedded assets:", err)
		}

		if assetsFS != nil {
			// Static assets from embedded FS
			app.Use("/", staticCompress, filesystem.New(filesystem.Config{
				Root:   http.FS(assetsFS),
				Browse: false,
				Next: func(c *fiber.Ctx) bool {
					return strings.HasPrefix(c.Path(), "/api") ||
						c.Path() == "/" ||
						!strings.Contains(c.Path(), ".")
				},
			}))

			// HTML with injection
			app.Get("/*", func(c *fiber.Ctx) error {
				if strings.HasPrefix(c.Path(), "/api") {
					return c.Next()
				}

				path := c.Path()
				if path != "/" && strings.Contains(path, ".") && !strings.HasSuffix(path, ".html") {
					return c.Next()
				}

				htmlBytes, err := fs.ReadFile(assetsFS, "index.html")
				if err != nil {
					return c.Status(500).SendString("Failed to load UI")
				}

				data := GetHTMLInjectionData(c, *configRef)
				initData := GetInitData(c, *configRef, data)
				initJSON, _ := json.Marshal(initData)
				html := InjectHTML(string(htmlBytes), data, string(initJSON))

				c.Set("Content-Type", "text/html; charset=utf-8")
				c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
				return c.SendString(html)
			})
		}
	}

	// İzinleri yükle
	if config.Permissions.Path != "" {
		if _, err := permission.Load(config.Permissions.Path); err != nil {
			// İzin dosyası yüklenemezse panic oluşturabilir veya loglayabiliriz.
			// Şimdilik panic yapalım ki geliştirici fark etsin.
			panic(fmt.Errorf("izin dosyası yüklenemedi: %w", err))
		}
	} else {
		// Varsayılan bir yol deneyebiliriz veya boş bırakabiliriz.
		// Örneğin "permissions.toml" var mı diye bakabiliriz.
		if _, err := os.Stat("permissions.toml"); err == nil {
			_, _ = permission.Load("permissions.toml")
		}
	}

	p := &Panel{
		Config:    config,
		Db:        db,
		Fiber:     app,
		Auth:      authService,
		resources: make(map[string]resource.Resource),
		pages:     make(map[string]page.Page),
		plugins:   make([]interface{}, 0),
	}

	// Load Dynamic Settings (veritabanından)
	_ = p.LoadSettings()

	// configRef'i p.Config'e yönlendir — artık closure'lar
	// veritabanından yüklenen güncel settings'i görür
	configRef = &p.Config

	// Plugin System: Auto-discovery (optional)
	if config.Plugins.AutoDiscover && config.Plugins.Path != "" {
		// Import plugin package for auto-discovery
		// Note: This is optional, manual import is preferred
		fmt.Printf("Plugin auto-discovery enabled, path: %s\n", config.Plugins.Path)
		// Auto-discovery implementation would go here
		// For now, we rely on manual import via init() functions
	}

	// Register Default Resources
	if p.Config.UserResource != nil {
		p.RegisterResource(p.Config.UserResource)
	} else {
		p.RegisterResource(resourceUser.GetUserResource())
	}

	// Register Additional Resources
	for _, res := range p.Config.Resources {
		p.RegisterResource(res)
	}

	// Auto-discover resources from global registry
	// Requires resources to register themselves via init() functions
	for _, res := range resource.List() {
		p.RegisterResource(res)
	}

	// Register Pages from Config
	// Kullanıcı tarafından tanımlanan sayfaları kaydet
	for _, pg := range p.Config.Pages {
		p.RegisterPage(pg)
	}

	// Register Dynamic Routes
	// /api/resource/:resource -> List/Index
	// /api/resource/:resource/:id -> Detail/Show/Update/Delete

	api := app.Group("/api")

	// RESILIENCE: Circuit Breaker middleware
	// Servis hatalarını yönetir ve sistem çökmelerini önler
	// Üç durum: Closed (Normal), Open (Devre Dışı), Half-Open (Test)
	if config.CircuitBreaker.Enabled {
		// Varsayılan değerleri ayarla
		failureThreshold := config.CircuitBreaker.FailureThreshold
		if failureThreshold == 0 {
			failureThreshold = 5 // 5 ardışık hata sonrası devre aç
		}

		timeout := config.CircuitBreaker.Timeout
		if timeout == 0 {
			timeout = 10 * time.Second // 10 saniye bekle
		}

		successThreshold := config.CircuitBreaker.SuccessThreshold
		if successThreshold == 0 {
			successThreshold = 5 // 5 başarılı istek sonrası devre kapat
		}

		halfOpenMaxConcurrent := config.CircuitBreaker.HalfOpenMaxConcurrent
		if halfOpenMaxConcurrent == 0 {
			halfOpenMaxConcurrent = 1 // Half-open'da 1 eşzamanlı istek
		}

		// Circuit Breaker oluştur
		cb := circuitbreaker.New(circuitbreaker.Config{
			FailureThreshold:      failureThreshold,
			Timeout:               timeout,
			SuccessThreshold:      successThreshold,
			HalfOpenMaxConcurrent: halfOpenMaxConcurrent,
			// IsFailure: Hangi hataların sayılacağını belirler (varsayılan: status >= 500)
			IsFailure: func(c *fiber.Ctx, err error) bool {
				// 500+ status kodları hata olarak sayılır
				return err != nil
			},
			// OnOpen: Devre açıldığında çağrılır (503 Service Unavailable)
			OnOpen: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
					"error":   "Service temporarily unavailable",
					"message": "The service is experiencing high failure rates. Please try again later.",
					"code":    "CIRCUIT_BREAKER_OPEN",
				})
			},
			// OnHalfOpen: Half-open durumunda çağrılır (429 Too Many Requests)
			OnHalfOpen: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error":   "Service is recovering",
					"message": "The service is testing recovery. Please wait.",
					"code":    "CIRCUIT_BREAKER_HALF_OPEN",
				})
			},
		})

		// Circuit Breaker'ı API route'larına uygula
		api.Use(circuitbreaker.Middleware(cb))
	}

	// langMiddleware: URL parametresinden veya request'ten dil bilgisini al
	api.Use(langMiddleware(config))

	// registerAPIRoutes: Tüm API route'larını kaydet
	// Bu closure fonksiyon, dual route registration için kullanılır
	registerAPIRoutes := func(apiGroup fiber.Router) {
		// Auth Routes
		authRoutes := apiGroup.Group("/auth")
		// SECURITY: Strict rate limiting for authentication endpoints (10 req/min)
		// TEMPORARY: Rate limiting disabled for development
		// authRoutes.Use(middleware.AuthRateLimiter())
		authRoutes.Post("/sign-in/email", context.Wrap(authH.LoginEmail))
		authRoutes.Post("/sign-up/email", context.Wrap(authH.RegisterEmail))
		authRoutes.Post("/sign-out", context.Wrap(authH.SignOut))
		authRoutes.Post("/forgot-password", context.Wrap(authH.ForgotPassword))
		authRoutes.Get("/session", context.Wrap(authH.GetSession))

		apiGroup.Get("/init", context.Wrap(p.handleInit)) // App Initialization

		// Middleware
		apiGroup.Use(context.Wrap(authH.SessionMiddleware))
		// SECURITY: Rate limiting for general API endpoints (100 req/min)
		// TEMPORARY: Rate limiting disabled for development
		// apiGroup.Use(middleware.APIRateLimiter())

		// Page Routes
		apiGroup.Get("/pages", context.Wrap(p.handlePages))
		apiGroup.Get("/pages/:slug", context.Wrap(p.handlePageDetail))
		apiGroup.Post("/pages/:slug", context.Wrap(p.handlePageSave))

		// Resource Routes
		apiGroup.Get("/resource/:resource/cards", context.Wrap(p.handleResourceCards))
		apiGroup.Get("/resource/:resource/cards/:index", context.Wrap(p.handleResourceCard))
		apiGroup.Get("/resource/:resource/lenses", context.Wrap(p.handleResourceLenses))                  // List available lenses
		apiGroup.Get("/resource/:resource/lens/:lens/cards", context.Wrap(p.handleResourceLensCards))     // Lens cards
		apiGroup.Get("/resource/:resource/lens/:lens", context.Wrap(p.handleResourceLens))                // Lens data
		apiGroup.Get("/resource/:resource/lens/:lens/actions", context.Wrap(p.handleResourceLensActions)) // Lens actions
		apiGroup.Post("/resource/:resource/lens/:lens/actions/:action", context.Wrap(p.handleResourceLensActionExecute))
		apiGroup.Get("/resource/:resource/morphable/:field", context.Wrap(p.handleMorphable))             // MorphTo field options
		apiGroup.Get("/resource/:resource/actions", context.Wrap(p.handleResourceActions))                // List available actions
		apiGroup.Post("/resource/:resource/actions/:action", context.Wrap(p.handleResourceActionExecute)) // Execute action
		apiGroup.Get("/resource/:resource", context.Wrap(p.handleResourceIndex))
		apiGroup.Post("/resource/:resource", context.Wrap(p.handleResourceStore))
		apiGroup.Get("/resource/:resource/create", context.Wrap(p.handleResourceCreate)) // New Route
		apiGroup.Get("/resource/:resource/:id", context.Wrap(p.handleResourceShow))
		apiGroup.Get("/resource/:resource/:id/detail", context.Wrap(p.handleResourceDetail))
		apiGroup.Get("/resource/:resource/:id/edit", context.Wrap(p.handleResourceEdit))
		apiGroup.Post("/resource/:resource/:id/fields/:field/resolve", context.Wrap(p.handleFieldResolve))          // Field resolver endpoint
		apiGroup.Post("/resource/:resource/fields/resolve-dependencies", context.Wrap(p.handleResolveDependencies)) // Dependency resolver endpoint

		// Hover card resolver endpoints - Support GET, POST, PATCH, DELETE
		apiGroup.Get("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve))    // Hover card resolver (GET)
		apiGroup.Post("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve))   // Hover card resolver (POST)
		apiGroup.Patch("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve))  // Hover card resolver (PATCH)
		apiGroup.Delete("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve)) // Hover card resolver (DELETE)
		apiGroup.Put("/resource/:resource/:id", context.Wrap(p.handleResourceUpdate))
		apiGroup.Delete("/resource/:resource/:id", context.Wrap(p.handleResourceDestroy))
		apiGroup.Get("/navigation", context.Wrap(p.handleNavigation)) // Sidebar Navigation

		// /resolve endpoint for dynamic routing check
		apiGroup.Get("/resolve", context.Wrap(p.handleResolve))
	}

	// Prefix'siz route'lar (varsayılan dil veya URL prefix kapalı)
	registerAPIRoutes(api)

	// Prefix'li route'lar (URL prefix açıksa)
	if config.I18n.Enabled && config.I18n.UseURLPrefix {
		apiLang := app.Group("/api/:lang")
		apiLang.Use(langMiddleware(config))
		registerAPIRoutes(apiLang)
	}

	// Notification Routes (dual route registration dışında)
	notificationProvider := data.NewGormDataProvider(db, &notificationDomain.Notification{})
	notificationService := notification.NewService(notificationProvider)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	api.Get("/notifications", context.Wrap(notificationHandler.HandleGetUnreadNotifications))
	api.Post("/notifications/:id/read", context.Wrap(notificationHandler.HandleMarkAsRead))
	api.Post("/notifications/read-all", context.Wrap(notificationHandler.HandleMarkAllAsRead))

	// OpenAPI Routes (dual route registration dışında)
	openAPIConfig := openapi.SpecGeneratorConfig{
		Title:       "Panel.go Admin Panel",
		Version:     "1.0.0",
		Description: "Panel.go Admin Panel API",
		ServerURL:   "",
	}
	p.openAPIHandler = handler.NewOpenAPIHandler(p.resources, openAPIConfig)
	api.Get("/openapi.json", p.openAPIHandler.GetSpec)
	api.Get("/docs", p.openAPIHandler.SwaggerUI)
	api.Get("/docs/redoc", p.openAPIHandler.ReDocUI)
	api.Get("/docs/rapidoc", p.openAPIHandler.RapidocUI)

	// Boot Plugins: Load plugins from global registry and boot them
	// Plugins register themselves via init() functions
	// This must be called after all routes are registered
	if err := p.bootPluginsFromRegistry(); err != nil {
		fmt.Printf("Warning: Plugin boot failed: %v\n", err)
	}

	return p
}

// / # bootPluginsFromRegistry Metodu
// /
// / Global registry'den plugin'leri alır ve boot eder.
// / Bu metod New() fonksiyonunda otomatik olarak çağrılır.
// /
// / ## Parametreler
// / Yok (alıcı: *Panel)
// /
// / ## Dönüş Değeri
// / - `error`: Plugin boot işlemi başarısızsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Global registry'den tüm plugin'leri alır
// / 2. Her plugin'i RegisterPlugin ile kaydeder
// / 3. BootPlugins ile tüm plugin'leri boot eder
// /
// / ## Önemli Notlar
// / - Bu metod New() fonksiyonunda otomatik olarak çağrılır
// / - Plugin'ler init() fonksiyonunda global registry'ye kaydedilmelidir
func (p *Panel) bootPluginsFromRegistry() error {
	// Global registry'den tüm plugin'leri al
	plugins := plugin.All()

	if len(plugins) == 0 {
		return nil
	}

	fmt.Printf("Booting %d plugins from registry...\n", len(plugins))

	// Her plugin'i kaydet
	for _, plg := range plugins {
		if err := p.RegisterPlugin(plg); err != nil {
			return fmt.Errorf("failed to register plugin '%s': %w", plg.Name(), err)
		}
	}

	// Tüm plugin'leri boot et
	if err := p.BootPlugins(); err != nil {
		return fmt.Errorf("failed to boot plugins: %w", err)
	}

	fmt.Printf("Successfully booted %d plugins\n", len(plugins))
	return nil
}

// / # LoadSettings Metodu
// /
// / Veritabanından dinamik ayarları okur ve Panel yapılandırmasını günceller.
// / Ayarlar JSON formatında depolanır ve otomatik olarak ayrıştırılır.
// /
// / ## Parametreler
// / Yok (alıcı: *Panel)
// /
// / ## Dönüş Değeri
// / - `error`: Veritabanı hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Setting tablosunun var olup olmadığını kontrol eder
// / 2. Tüm ayarları veritabanından okur
// / 3. JSON değerleri ayrıştırır (başarısız olursa string olarak işler)
// / 4. Ayarları yapılandırmaya yükler
// / 5. Özel ayarları (site_name, register, forgot_password) işler
// /
// / ## Desteklenen Ayarlar
// / - `site_name`: Sitenin adı (string)
// / - `register`: Kayıt özelliğinin etkin olup olmadığı (boolean)
// / - `forgot_password`: Şifremi unuttum özelliğinin etkin olup olmadığı (boolean)
// /
// / ## Kullanım Örneği
// / ```go
// / p := panel.New(config)
// / if err := p.LoadSettings(); err != nil {
// /     log.Printf("Ayarlar yüklenemedi: %v", err)
// / }
// / ```
// /
// / ## Uyarılar
// / - İlk çalıştırmada Setting tablosu henüz oluşturulmamış olabilir
// / - JSON ayrıştırma başarısız olursa değer string olarak işlenir
// / - Ayarlar başlatıldıktan sonra yüklenmelidir
// /
// / ## Önemli Notlar
// / - Tablo yoksa hata döndürmez, sessizce devam eder
// / - Ayarlar yapılandırmayı geçersiz kılar
// / - Dinamik ayarlar üretim ortamında kullanılabilir
func (p *Panel) LoadSettings() error {
	var settings []setting.Setting
	// Tablo yoksa henüz hata vermesin (ilk çalıştırma)
	if !p.Db.Migrator().HasTable(&setting.Setting{}) {
		return nil
	}

	// Features config'ini settings tablosu ile senkronize et.
	// Eğer key yoksa config'deki varsayılan değerle oluştur.
	featureDefaults := map[string]string{
		"registration_enabled":    fmt.Sprintf("%v", p.Config.Features.Register),
		"forgot_password_enabled": fmt.Sprintf("%v", p.Config.Features.ForgotPassword),
	}
	for key, defaultVal := range featureDefaults {
		var count int64
		p.Db.Model(&setting.Setting{}).Where("key = ?", key).Count(&count)
		if count == 0 {
			p.Db.Create(&setting.Setting{Key: key, Value: defaultVal})
		}
	}

	if err := p.Db.Find(&settings).Error; err != nil {
		return err
	}

	config := &p.Config.SettingsValues
	if config.Values == nil {
		config.Values = make(map[string]interface{})
	}

	for _, s := range settings {
		// Parse JSON value
		var val interface{}
		if err := json.Unmarshal([]byte(s.Value), &val); err != nil {
			// If not JSON, treat as string
			val = s.Value
		}

		config.Values[s.Key] = val

		switch s.Key {
		case "site_name":
			if v, ok := val.(string); ok {
				config.SiteName = v
			}
		case "registration_enabled":
			boolVal := parseBoolValue(val)
			config.Register = boolVal
			p.Config.Features.Register = boolVal
		case "forgot_password_enabled":
			boolVal := parseBoolValue(val)
			config.ForgotPassword = boolVal
			p.Config.Features.ForgotPassword = boolVal
		}
	}
	return nil
}

// parseBoolValue, çeşitli tiplerdeki değerleri bool'a çevirir.
// String "true"/"false", bool true/false ve float64 1/0 destekler.
func parseBoolValue(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1"
	case float64:
		return v != 0
	}
	return false
}

// / # Register Metodu
// /
// / Verilen slug ve kaynak (Resource) çiftini Panel'e kaydeder.
// / Kaynağın diyalog türünü Sheet olarak ayarlar.
// /
// / ## Parametreler
// / - `slug`: Kaynağın benzersiz tanımlayıcısı (örn: "users", "products")
// / - `res`: Kayıtlı edilecek Resource nesnesi
// /
// / ## Davranış
// / 1. Kaynağın diyalog türünü Sheet olarak ayarlar
// / 2. Kaynağı resources haritasına ekler
// /
// / ## Kullanım Örneği
// / ```go
// / p := panel.New(config)
// / userResource := resourceUser.GetUserResource()
// / p.Register("users", userResource)
// / ```
// /
// / ## Uyarılar
// / - Aynı slug ile birden fazla kaynak kaydedilirse, son kayıt öncekini geçersiz kılar
// / - Kaynaklar Start() çağrılmadan önce kayıtlı olmalıdır
// /
// / ## Önemli Notlar
// / - Diyalog türü her zaman Sheet olarak ayarlanır
// / - Slug benzersiz olmalıdır
func (p *Panel) Register(slug string, res resource.Resource) {
	res.SetDialogType(resource.DialogTypeSheet)
	p.resources[slug] = res
}

// / # RegisterResource Metodu
// /
// / Kaynağın kendi slug'ını kullanarak Panel'e kaydeder.
// / Register() metodunun bir sarmalayıcısıdır.
// /
// / ## Parametreler
// / - `res`: Kayıtlı edilecek Resource nesnesi
// /
// / ## Davranış
// / Resource'un Slug() metodunu çağırarak slug'ı alır ve Register() metodunu çağırır.
// /
// / ## Kullanım Örneği
// / ```go
// / p := panel.New(config)
// / userResource := resourceUser.GetUserResource()
// / p.RegisterResource(userResource)
// / ```
// /
// / ## Avantajlar
// / - Daha temiz ve daha az hata yapma olasılığı
// / - Resource'un kendi slug'ını kullanır
// /
// / ## Uyarılar
// / - Resource'un Slug() metodu doğru değer döndürmelidir
func (p *Panel) RegisterResource(res resource.Resource) {
	p.Register(res.Slug(), res)
}

// / # RegisterPage Metodu
// /
// / Verilen sayfayı (Page) Panel'e kaydeder.
// / Sayfanın slug'ını kullanarak pages haritasına ekler.
// /
// / ## Parametreler
// / - `pg`: Kayıtlı edilecek Page nesnesi
// /
// / ## Davranış
// / Page'in Slug() metodunu çağırarak slug'ı alır ve pages haritasına ekler.
// /
// / ## Kullanım Örneği
// / ```go
// / p := panel.New(config)
// / dashboardPage := &page.Dashboard{}
// / p.RegisterPage(dashboardPage)
// / ```
// /
// / ## Desteklenen Sayfalar
// / - Dashboard: Ana kontrol paneli
// / - Settings: Sistem ayarları
// / - Özel sayfalar: Kullanıcı tarafından tanımlanan sayfalar
// /
// / ## Uyarılar
// / - Aynı slug ile birden fazla sayfa kaydedilirse, son kayıt öncekini geçersiz kılar
// / - Sayfalar Start() çağrılmadan önce kayıtlı olmalıdır
// /
// / ## Önemli Notlar
// / - Sayfanın Slug() metodu benzersiz bir değer döndürmelidir
// / - Sayfalar API yönlendirmelerine otomatik olarak eklenir
func (p *Panel) RegisterPage(pg page.Page) {
	p.pages[pg.Slug()] = pg
}

// / # Start Metodu
// /
// / Panel'i yapılandırılan host ve port'ta başlatır.
// / HTTP sunucusunu başlatır ve gelen istekleri dinlemeye başlar.
// /
// / ## Parametreler
// / Yok (alıcı: *Panel)
// /
// / ## Dönüş Değeri
// / - `error`: Sunucu başlatılamadığında hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Host ve port'tan adres oluşturur (örn: "localhost:8080")
// / 2. Fiber uygulamasını Listen() metoduyla başlatır
// / 3. Gelen HTTP isteklerini dinlemeye başlar
// /
// / ## Kullanım Örneği
// / ```go
// / p := panel.New(config)
// / if err := p.Start(); err != nil {
// /     log.Fatal(err)
// / }
// / ```
// /
// / ## Uyarılar
// / - Bu metod bloklanır ve sunucu kapatılana kadar döndürmez
// / - Kaynaklar ve sayfalar Start() çağrılmadan önce kayıtlı olmalıdır
// / - Port zaten kullanımda ise hata döndürülür
// /
// / ## Önemli Notlar
// / - Üretim ortamında uygun host ve port yapılandırması yapılmalıdır
// / - Sunucuyu durdurmak için SIGINT (Ctrl+C) sinyali gönderilebilir
// / - Tüm middleware'ler ve yönlendirmeler bu noktada aktif hale gelir
func (p *Panel) Start() error {
	addr := fmt.Sprintf("%s:%s", p.Config.Server.Host, p.Config.Server.Port)
	return p.Fiber.Listen(addr)
}

// / # withResourceHandler Metodu
// /
// / Kaynağı çözer ve bir FieldHandler oluşturarak verilen fonksiyonu çalıştırır.
// / Tüm kaynak işleme endpoint'leri için kullanılan yardımcı metod.
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// / - `fn`: Çalıştırılacak fonksiyon (FieldHandler alır ve error döndürür)
// /
// / ## Dönüş Değeri
// / - `error`: Kaynak bulunamadığında veya fonksiyon hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. URL parametrelerinden "resource" slug'ını alır
// / 2. Kaynağı resources haritasında arar
// / 3. Kaynak bulunamazsa 404 hatası döndürür
// / 4. FieldHandler oluşturur
// / 5. Verilen fonksiyonu FieldHandler ile çalıştırır
// /
// / ## Kullanım Örneği
// / ```go
// / func (p *Panel) handleResourceIndex(c *context.Context) error {
// /     return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
// /         return handler.HandleResourceIndex(h, c)
// /     })
// / }
// / ```
// /
// / ## Avantajlar
// / - Kod tekrarını azaltır
// / - Tutarlı hata işleme
// / - Merkezi kaynak çözümleme
// /
// / ## Uyarılar
// / - Kaynak slug'ı URL parametresi olarak geçilmelidir
// / - Kaynak önceden kayıtlı olmalıdır
// /
// / ## Önemli Notlar
// / - Bu metod tüm kaynak endpoint'leri tarafından kullanılır
// / - Hata işleme otomatik olarak yapılır
func (p *Panel) withResourceHandler(c *context.Context, fn func(*handler.FieldHandler) error) error {
	slug := c.Params("resource")
	res, ok := p.resources[slug]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}
	h := handler.NewResourceHandler(p.Db, res, p.Config.Storage.Path, p.Config.Storage.URL)
	return fn(h)
}

// / # handleResourceIndex Metodu
// /
// / Kaynağın tüm kayıtlarını listeler. Filtreleme, sıralama ve sayfalama destekler.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceIndex handler'ını çalıştırır
// /
// / ## Kullanım Örneği
// / ```
// / GET /api/resource/users
// / ```
// /
// / ## Yanıt Örneği
// / ```json
// / {
// /   "data": [
// /     {"id": 1, "name": "John", "email": "john@example.com"},
// /     {"id": 2, "name": "Jane", "email": "jane@example.com"}
// /   ],
// /   "meta": {
// /     "total": 2,
// /     "per_page": 15,
// /     "current_page": 1
// /   }
// / }
// / ```
func (p *Panel) handleResourceIndex(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceIndex(h, c)
	})
}

// / # handleResourceShow Metodu
// /
// / Kaynağın tek bir kaydını gösterir. Temel bilgileri döndürür.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/:id`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceShow handler'ını çalıştırır
// /
// / ## Kullanım Örneği
// / ```
// / GET /api/resource/users/1
// / ```
func (p *Panel) handleResourceShow(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceShow(h, c)
	})
}

// / # handleResourceDetail Metodu
// /
// / Kaynağın detaylı bilgilerini gösterir. Tüm alanları ve ilişkileri içerir.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/:id/detail`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceDetail handler'ını çalıştırır
func (p *Panel) handleResourceDetail(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceDetail(h, c)
	})
}

// / # handleResourceStore Metodu
// /
// / Yeni bir kayıt oluşturur. POST isteği ile gönderilen verileri veritabanına kaydeder.
// /
// / ## HTTP Endpoint
// / `POST /api/resource/:resource`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## İstek Gövdesi Örneği
// / ```json
// / {
// /   "name": "John Doe",
// /   "email": "john@example.com",
// /   "role": "admin"
// / }
// / ```
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceStore handler'ını çalıştırır
// / 4. Yeni kaydı veritabanına kaydeder
func (p *Panel) handleResourceStore(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceStore(h, c)
	})
}

// / # handleResourceCreate Metodu
// /
// / Yeni kayıt oluşturma formunun alanlarını döndürür.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/create`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceCreate handler'ını çalıştırır
// / 4. Oluşturma formunun alanlarını döndürür
func (p *Panel) handleResourceCreate(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceCreate(h, c)
	})
}

// / # handleResourceUpdate Metodu
// /
// / Mevcut bir kaydı günceller. PUT isteği ile gönderilen verileri veritabanında günceller.
// /
// / ## HTTP Endpoint
// / `PUT /api/resource/:resource/:id`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## İstek Gövdesi Örneği
// / ```json
// / {
// /   "name": "Jane Doe",
// /   "email": "jane@example.com"
// / }
// / ```
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceUpdate handler'ını çalıştırır
// / 4. Kaydı veritabanında günceller
func (p *Panel) handleResourceUpdate(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceUpdate(h, c)
	})
}

// / # handleResourceDestroy Metodu
// /
// / Bir kaydı siler. DELETE isteği ile kaydı veritabanından kaldırır.
// /
// / ## HTTP Endpoint
// / `DELETE /api/resource/:resource/:id`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceDestroy handler'ını çalıştırır
// / 4. Kaydı veritabanından siler
// /
// / ## Uyarılar
// / - Bu işlem geri alınamaz
// / - İlişkili kayıtlar etkilenebilir
func (p *Panel) handleResourceDestroy(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceDestroy(h, c)
	})
}

// / # handleResourceEdit Metodu
// /
// / Kaydı düzenleme formunun alanlarını döndürür.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/:id/edit`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResourceEdit handler'ını çalıştırır
// / 4. Düzenleme formunun alanlarını döndürür
func (p *Panel) handleResourceEdit(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceEdit(h, c)
	})
}

// / # handleFieldResolve Metodu
// /
// / Dinamik alan değerlerini çözer. Bağımlı alanlar için seçenekleri döndürür.
// /
// / ## HTTP Endpoint
// / `POST /api/resource/:resource/:id/fields/:field/resolve`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleFieldResolve handler'ını çalıştırır
// / 4. Alan değerlerini çözer
func (p *Panel) handleFieldResolve(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleFieldResolve(h, c)
	})
}

// / # handleHoverCardResolve Metodu
// /
// / Hover card verilerini çözmek için HTTP handler wrapper'ı.
// /
// / ## Parametreler
// / - `c *context.Context`: Fiber context wrapper'ı
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleHoverCardResolve handler'ını çalıştırır
// / 4. Hover card verilerini döndürür
// /
// / ## Endpoint
// / - GET /api/resource/:resource/resolver/:field
// / - POST /api/resource/:resource/resolver/:field
// / - PATCH /api/resource/:resource/resolver/:field
// / - DELETE /api/resource/:resource/resolver/:field
func (p *Panel) handleHoverCardResolve(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleHoverCardResolve(h)(c.Ctx)
	})
}

// / # handleResolveDependencies Metodu
// /
// / Alan bağımlılıklarını çözer. Bağımlı alanların seçeneklerini döndürür.
// /
// / ## HTTP Endpoint
// / `POST /api/resource/:resource/fields/resolve-dependencies`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleResolveDependencies handler'ını çalıştırır
// / 4. Bağımlılıkları çözer
func (p *Panel) handleResolveDependencies(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResolveDependencies(h, c)
	})
}

// / # handleResourceCards Metodu
// /
// / Kaynağın tüm kartlarını listeler. Kartlar, kaynağın özet görünümüdür.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/cards`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleCardList handler'ını çalıştırır
// / 4. Tüm kartları döndürür
func (p *Panel) handleResourceCards(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleCardList(h, c)
	})
}

// / # handleResourceCard Metodu
// /
// / Kaynağın belirli bir kartını gösterir.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/cards/:index`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleCardDetail handler'ını çalıştırır
// / 4. Belirli kartın detaylarını döndürür
func (p *Panel) handleResourceCard(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleCardDetail(h, c)
	})
}

// / # handleResourceLenses Metodu
// /
// / Kaynağın tüm lens'lerini listeler. Lens'ler, kaynağın alternatif görünümleridir.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/lenses`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleLensIndex handler'ını çalıştırır
// / 4. Tüm lens'leri döndürür
// /
// / ## Lens Nedir?
// / Lens'ler, kaynağın farklı perspektiflerden görüntülenmesini sağlar.
// / Örneğin, kullanıcılar kaynağını "Aktif Kullanıcılar" ve "Pasif Kullanıcılar" lens'leriyle görebilir.
func (p *Panel) handleResourceLenses(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleLensIndex(h, c)
	})
}

// / # handleResourceLens Metodu
// /
// / Kaynağın belirli bir lens'ini gösterir. Lens'e göre filtrelenmiş verileri döndürür.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/lens/:lens`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. URL parametrelerinden resource ve lens slug'larını alır
// / 2. Kaynağı resources haritasında arar
// / 3. Lens'i kaynağın lens'leri arasında arar
// / 4. LensHandler oluşturur
// / 5. HandleLens handler'ını çalıştırır
// / 6. Lens'e göre filtrelenmiş verileri döndürür
// /
// / ## Kullanım Örneği
// / ```
// / GET /api/resource/users/lens/active-users
// / ```
func (p *Panel) handleResourceLens(c *context.Context) error {
	return p.withLensHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleLens(h, c)
	})
}

func (p *Panel) withLensHandler(c *context.Context, fn func(*handler.FieldHandler) error) error {
	slug := c.Params("resource")
	lensSlug := c.Params("lens")

	res, ok := p.resources[slug]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}

	// Find Lens
	var targetLens resource.Lens
	for _, l := range res.GetLenses() {
		if l.Slug() == lensSlug {
			targetLens = l
			break
		}
	}

	if targetLens == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lens not found",
		})
	}

	// Create Handler for Lens
	h := handler.NewLensHandler(p.Db, res, targetLens)

	return fn(h)
}

// handleResourceLensCards returns cards for the selected lens.
func (p *Panel) handleResourceLensCards(c *context.Context) error {
	return p.withLensHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleLensCards(h, c)
	})
}

// handleResourceLensActions lists available actions for the selected lens.
func (p *Panel) handleResourceLensActions(c *context.Context) error {
	return p.withLensHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleActionList(h, c)
	})
}

// handleResourceLensActionExecute executes an action in lens context.
func (p *Panel) handleResourceLensActionExecute(c *context.Context) error {
	return p.withLensHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleActionExecute(h, c)
	})
}

// / # handleMorphable Metodu
// /
// / MorphTo alanı için seçenekleri döndürür. Polimorfik ilişkiler için kullanılır.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/morphable/:field`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleMorphable handler'ını çalıştırır
// / 4. MorphTo alanı için seçenekleri döndürür
// /
// / ## MorphTo Nedir?
// / MorphTo, bir modelin birden fazla model türüne ait olabileceği polimorfik ilişkidir.
// / Örneğin, bir Yorum (Comment) modeli, Yazı (Post) veya Video (Video) modeline ait olabilir.
func (p *Panel) handleMorphable(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleMorphable(h, c)
	})
}

// / # handleResourceActions Metodu
// /
// / Kaynağın tüm eylemlerini (action) listeler. Eylemler, toplu işlemler için kullanılır.
// /
// / ## HTTP Endpoint
// / `GET /api/resource/:resource/actions`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleActionList handler'ını çalıştırır
// / 4. Tüm eylemleri döndürür
// /
// / ## Eylem Nedir?
// / Eylemler, seçilen kayıtlar üzerinde toplu işlemler gerçekleştirmek için kullanılır.
// / Örneğin, "Sil", "Yayınla", "Arşivle" gibi eylemler.
func (p *Panel) handleResourceActions(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleActionList(h, c)
	})
}

// / # handleResourceActionExecute Metodu
// /
// / Belirli bir eylemi (action) çalıştırır. Seçilen kayıtlar üzerinde işlem gerçekleştirir.
// /
// / ## HTTP Endpoint
// / `POST /api/resource/:resource/actions/:action`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## İstek Gövdesi Örneği
// / ```json
// / {
// /   "ids": [1, 2, 3],
// /   "fields": {
// /     "status": "published"
// /   }
// / }
// / ```
// /
// / ## Davranış
// / 1. Kaynağı çözer
// / 2. FieldHandler oluşturur
// / 3. HandleActionExecute handler'ını çalıştırır
// / 4. Eylemi seçilen kayıtlar üzerinde çalıştırır
func (p *Panel) handleResourceActionExecute(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleActionExecute(h, c)
	})
}

// / # handleNavigation Metodu
// /
// / Yan menü (sidebar) için navigasyon öğelerini döndürür.
// / Kayıtlı kaynaklar ve sayfaları hiyerarşik olarak organize eder.
// /
// / ## HTTP Endpoint
// / `GET /api/navigation`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Tüm kayıtlı kaynakları ve sayfaları toplar
// / 2. Görünür olanları filtreler (Visible() == true)
// / 3. Her öğe için NavItem oluşturur
// / 4. Öğeleri NavigationOrder'a göre sıralar
// / 5. Aynı order'a sahip öğeleri başlığa göre alfabetik sıralar
// / 6. JSON formatında döndürür
// /
// / ## Yanıt Örneği
// / ```json
// / {
// /   "data": [
// /     {
// /       "slug": "users",
// /       "title": "Kullanıcılar",
// /       "icon": "users",
// /       "group": "Yönetim",
// /       "type": "resource",
// /       "order": 1
// /     },
// /     {
// /       "slug": "dashboard",
// /       "title": "Kontrol Paneli",
// /       "icon": "dashboard",
// /       "group": "",
// /       "type": "page",
// /       "order": 0
// /     }
// /   ]
// / }
// / ```
// /
// / ## NavItem Yapısı
// / - `slug`: Kaynağın veya sayfanın benzersiz tanımlayıcısı
// / - `title`: Görüntülenecek başlık
// / - `icon`: İkon adı (CSS sınıfı veya ikon kütüphanesi)
// / - `group`: Menü grubu (örn: "Yönetim", "Ayarlar")
// / - `type`: "resource" veya "page"
// / - `order`: Sıralama için kullanılan sayı (dahili kullanım)
// /
// / ## Kullanım Örneği
// / ```
// / GET /api/navigation
// / ```
func (p *Panel) handleNavigation(c *context.Context) error {
	type NavItem struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
		Icon  string `json:"icon"`
		Group string `json:"group"`
		Type  string `json:"type"`  // "resource" or "page"
		Order int    `json:"order"` // Internal use for sorting
		URL   string `json:"url"`   // Full URL with language prefix
	}

	// Dil bilgisini al
	lang := i18n.GetLocale(c.Ctx)
	defaultLang := p.Config.I18n.DefaultLanguage.String()

	// URL prefix'ini hesapla
	urlPrefix := ""
	if p.Config.I18n.Enabled && p.Config.I18n.UseURLPrefix {
		if !p.Config.I18n.URLPrefixOptional || lang != defaultLang {
			urlPrefix = "/" + lang
		}
	}

	items := []NavItem{}
	for slug, res := range p.resources {
		if !res.Visible() {
			continue
		}
		items = append(items, NavItem{
			Slug:  slug,
			Title: res.TitleWithContext(c.Ctx),
			Icon:  res.Icon(),
			Group: res.GroupWithContext(c.Ctx),
			Type:  "resource",
			Order: res.NavigationOrder(),
			URL:   urlPrefix + "/resource/" + slug,
		})
	}

	for slug, pg := range p.pages {
		if !pg.Visible() {
			continue
		}
		items = append(items, NavItem{
			Slug:  slug,
			Title: i18n.Trans(c.Ctx, pg.Title()),
			Icon:  pg.Icon(),
			Group: pg.Group(),
			Type:  "page",
			Order: pg.NavigationOrder(),
			URL:   urlPrefix + "/page/" + slug,
		})
	}

	// Sort Items: Order (asc), then Title (asc)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Order != items[j].Order {
			return items[i].Order < items[j].Order
		}
		return items[i].Title < items[j].Title
	})

	return c.JSON(fiber.Map{
		"data": items,
	})
}

// / # handleInit Metodu
// /
// / Uygulamanın başlatılması için gerekli bilgileri döndürür.
// / Özellikler, OAuth ayarları, sürüm ve dinamik ayarları içerir.
// /
// / ## HTTP Endpoint
// / `GET /api/init`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Yapılandırmadan özellik ayarlarını alır
// / 2. Veritabanı ayarlarını kontrol eder (varsa)
// / 3. OAuth ayarlarını kontrol eder
// / 4. Tüm bilgileri JSON formatında döndürür
// /
// / ## Yanıt Örneği
// / ```json
// / {
// /   "features": {
// /     "register": true,
// /     "forgot_password": false
// /   },
// /   "oauth": {
// /     "google": true
// /   },
// /   "version": "1.0.0",
// /   "settings": {
// /     "site_name": "Panel.go",
// /     "register": true,
// /     "forgot_password": false
// /   }
// / }
// / ```
// /
// / ## Özellikler
// / - `register`: Kullanıcı kaydı etkin mi?
// / - `forgot_password`: Şifremi unuttum özelliği etkin mi?
// /
// / ## Kullanım Örneği
// / ```
// / GET /api/init
// / ```
// /
// / ## Uyarılar
// / - Bu endpoint kimlik doğrulama gerektirmez
// / - Hassas bilgiler döndürülmemelidir
// / - Ayarlar veritabanından yüklenir (varsa)
// / # getLanguageName Fonksiyonu
// /
// / Dil kodundan (language.Tag) dil adını döndürür.
// / i18n desteklenen diller listesi için kullanılır.
// /
// / ## Parametreler
// / - `lang`: Dil kodu (language.Tag)
// /
// / ## Dönüş Değeri
// / - `string`: Dil adı (örn: "Türkçe", "English")
// /
// / ## Davranış
// / 1. Dil kodunu string'e çevirir
// / 2. Bilinen diller haritasında arar
// / 3. Bulunursa dil adını döndürür
// / 4. Bulunamazsa dil kodunu döndürür (fallback)
// /
// / ## Desteklenen Diller
// / - tr: Türkçe
// / - en: English
// / - de: Deutsch
// / - fr: Français
// / - es: Español
// / - ar: العربية
// /
// / ## Kullanım Örneği
// / ```go
// / name := getLanguageName(language.Turkish)
// / // Çıktı: "Türkçe"
// / ```
func getLanguageName(lang language.Tag) string {
	names := map[string]string{
		"tr": "Türkçe",
		"en": "English",
		"de": "Deutsch",
		"fr": "Français",
		"es": "Español",
		"ar": "العربية",
	}

	if name, ok := names[lang.String()]; ok {
		return name
	}
	return lang.String()
}

func (p *Panel) handleInit(c *context.Context) error {
	// Get CSRF token from context (set by CSRF middleware) and send in response header
	if csrfToken := c.Locals("csrf"); csrfToken != nil {
		c.Set("X-CSRF-Token", csrfToken.(string))
	}

	injectionData := GetHTMLInjectionData(c.Ctx, p.Config)
	initData := GetInitData(c.Ctx, p.Config, injectionData)

	return c.JSON(initData)
}

// / # handleResolve Metodu
// /
// / Verilen yolu (path) çözer ve karşılık gelen kaynağı veya sayfayı bulur.
// / Dinamik yönlendirme için kullanılır.
// /
// / ## HTTP Endpoint
// / `GET /api/resolve?path=/users`
// /
// / ## Parametreler
// / - `c`: İstek bağlamı (Context)
// /
// / ## Query Parametreleri
// / - `path`: Çözülecek yol (örn: "/users", "users")
// /
// / ## Dönüş Değeri
// / - `error`: İşlem hatası varsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Query parametrelerinden "path" değerini alır
// / 2. Başında "/" varsa kaldırır
// / 3. Kaynaklar arasında arar
// / 4. Bulunursa kaynak bilgilerini döndürür
// / 5. Bulunamazsa 404 hatası döndürür
// /
// / ## Yanıt Örneği (Başarılı)
// / ```json
// / {
// /   "type": "resource",
// /   "slug": "users",
// /   "meta": {
// /     "title": "Kullanıcılar",
// /     "icon": "users",
// /     "group": "Yönetim"
// /   }
// / }
// / ```
// /
// / ## Yanıt Örneği (Başarısız)
// / ```json
// / {
// /   "error": "Page not found"
// / }
// / ```
// /
// / ## Kullanım Örneği
// / ```
// / GET /api/resolve?path=/users
// / GET /api/resolve?path=products
// / ```
// /
// / ## Önemli Notlar
// / - Şu anda sadece kaynaklar desteklenir
// / - Gelecekte özel sayfalar ve veritabanı tabanlı sayfalar eklenebilir
// / - Yol normalleştirilir (başında "/" kaldırılır)
func (p *Panel) handleResolve(c *context.Context) error {
	path := c.Query("path")
	// Simple Logic: Check if path matches a known resource slug
	// E.g. path "/users" -> Resource "users"
	// We might strip leading "/"
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// Check Resources
	if res, ok := p.resources[path]; ok {
		return c.JSON(fiber.Map{
			"type": "resource",
			"slug": path,
			"meta": fiber.Map{
				"title": res.TitleWithContext(c.Ctx),
				"icon":  res.Icon(),
				"group": res.GroupWithContext(c.Ctx),
			},
		})
	}

	// Future: Check custom pages, database driven pages etc.

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Page not found",
	})
}

// / # RegisterPlugin Metodu
// /
// / Plugin'i Panel'e kaydeder ve plugin'in Register() metodunu çağırır.
// /
// / ## Parametreler
// / - `p`: Kaydedilecek plugin
// /
// / ## Dönüş Değeri
// / - `error`: Plugin kaydı başarısızsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Plugin'in Register() metodunu çağırır
// / 2. Plugin'i plugins listesine ekler
// /
// / ## Kullanım Örneği
// / ```go
// / p := panel.New(config)
// / if err := p.RegisterPlugin(&MyPlugin{}); err != nil {
// /     log.Fatal(err)
// / }
// / ```
// /
// / ## Önemli Notlar
// / - Bu metod Start() çağrılmadan önce çağrılmalıdır
// / - Plugin'in Register() metodu hata döndürürse kayıt yapılmaz
func (p *Panel) RegisterPlugin(plugin interface{}) error {
	// Type assertion: Plugin interface'ini kontrol et
	if pluginImpl, ok := plugin.(interface {
		Register(panel interface{}) error
	}); ok {
		if err := pluginImpl.Register(p); err != nil {
			return fmt.Errorf("plugin registration failed: %w", err)
		}
		p.plugins = append(p.plugins, plugin)
		return nil
	}
	return fmt.Errorf("plugin does not implement Register method")
}

type pluginPageAdapter struct {
	page.Base
	src plugin.Page
}

func newPluginPageAdapter(src plugin.Page) page.Page {
	return &pluginPageAdapter{src: src}
}

func (a *pluginPageAdapter) Slug() string {
	return a.src.Slug()
}

func (a *pluginPageAdapter) Title() string {
	return a.src.Title()
}

func (a *pluginPageAdapter) Icon() string {
	return a.src.Icon()
}

func (a *pluginPageAdapter) Group() string {
	return a.src.Group()
}

func (a *pluginPageAdapter) Visible() bool {
	return a.src.Visible()
}

func (a *pluginPageAdapter) NavigationOrder() int {
	return a.src.NavigationOrder()
}

// / # BootPlugins Metodu
// /
// / Tüm kayıtlı plugin'leri boot eder. Plugin'lerin Boot() metodunu çağırır.
// /
// / ## Parametreler
// / Yok (alıcı: *Panel)
// /
// / ## Dönüş Değeri
// / - `error`: Plugin boot işlemi başarısızsa hata, aksi takdirde nil
// /
// / ## Davranış
// / 1. Tüm kayıtlı plugin'leri dolaşır
// / 2. Her plugin'in Boot() metodunu çağırır
// / 3. Plugin'lerin resource, page, middleware vb. eklemelerini yapar
// /
// / ## Kullanım Örneği
// / ```go
// / p := panel.New(config)
// / p.RegisterPlugin(&MyPlugin{})
// / if err := p.BootPlugins(); err != nil {
// /     log.Fatal(err)
// / }
// / p.Start()
// / ```
// /
// / ## Önemli Notlar
// / - Bu metod Start() çağrılmadan önce çağrılmalıdır
// / - Plugin'lerin Boot() metodu hata döndürürse boot işlemi durur
// / - Plugin'ler sırayla boot edilir
func (p *Panel) BootPlugins() error {
	for _, plg := range p.plugins {
		// Type assertion: Plugin interface'ini kontrol et
		if pluginImpl, ok := plg.(interface {
			Boot(panel interface{}) error
			Name() string
		}); ok {
			if err := pluginImpl.Boot(p); err != nil {
				return fmt.Errorf("plugin '%s' boot failed: %w", pluginImpl.Name(), err)
			}

			// Plugin'in resource'larını kaydet
			if resourceProvider, ok := plg.(interface {
				Resources() []resource.Resource
			}); ok {
				if resources := resourceProvider.Resources(); resources != nil {
					for _, res := range resources {
						p.RegisterResource(res)
					}
				}
			}

			// Plugin'in page'lerini kaydet
			if pageProvider, ok := plg.(interface {
				Pages() []plugin.Page
			}); ok {
				if pluginPages := pageProvider.Pages(); pluginPages != nil {
					for _, pg := range pluginPages {
						if pg == nil {
							continue
						}
						// Full page.Page implementasyonu varsa doğrudan kaydet.
						if fullPage, ok := pg.(page.Page); ok {
							p.RegisterPage(fullPage)
							continue
						}
						// Plugin page alt kümesini page.Page'e adapt et.
						p.RegisterPage(newPluginPageAdapter(pg))
					}
				}
			}

			// Plugin'in middleware'lerini kaydet
			if middlewareProvider, ok := plg.(interface {
				Middleware() []fiber.Handler
			}); ok {
				if middlewares := middlewareProvider.Middleware(); middlewares != nil {
					for _, mw := range middlewares {
						p.Fiber.Use(mw)
					}
				}
			}

			// Plugin'in route'larını kaydet
			if routeProvider, ok := plg.(interface {
				Routes(router fiber.Router)
			}); ok {
				routeProvider.Routes(p.Fiber)
			}

			// Plugin'in migration'larını çalıştır
			if migrationProvider, ok := plg.(interface {
				Migrations() []plugin.Migration
			}); ok {
				if migrations := migrationProvider.Migrations(); migrations != nil {
					for _, migration := range migrations {
						if migration == nil {
							continue
						}
						if err := migration.Up(p.Db); err != nil {
							return fmt.Errorf("plugin '%s' migration '%s' failed: %w", pluginImpl.Name(), migration.Name(), err)
						}
					}
				}
			}
		}
	}
	return nil
}
