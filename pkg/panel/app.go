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
	"github.com/ferdiunal/panel.go/pkg/data/orm"
	"github.com/ferdiunal/panel.go/pkg/domain/account"
	notificationDomain "github.com/ferdiunal/panel.go/pkg/domain/notification"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/setting"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/domain/verification"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/handler"
	authHandler "github.com/ferdiunal/panel.go/pkg/handler/auth"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/ferdiunal/panel.go/pkg/notification"
	"github.com/ferdiunal/panel.go/pkg/openapi"
	"github.com/ferdiunal/panel.go/pkg/page"
	"github.com/ferdiunal/panel.go/pkg/permission"
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
	http2PushResources []string // Cached list of critical assets for HTTP/2 push
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
func New(config Config) *Panel {
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

	// Static file serving
	// For SDK users: always use embedded assets
	// For SDK developers: use local path if available (development mode only)
	useEmbed := true
	localUIPath := "./pkg/panel/ui/index.html"

	// Check if we're in SDK development mode (local UI files exist)
	if config.Environment == "development" {
		if _, statErr := os.Stat(localUIPath); statErr == nil {
			useEmbed = false
		}
	}

	assetsFS, err := GetFileSystem(useEmbed)
	if err != nil {
		fmt.Println("Warning: Failed to load embedded assets:", err)
	}

	// PERFORMANCE: Discover critical assets for HTTP/2 push (if enabled)
	// Automatically find JS, CSS, and font files in assets directory
	var http2PushResources []string
	if config.EnableHTTP2Push && assetsFS != nil {
		fs.WalkDir(assetsFS, "assets", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			// Push critical resources: JS, CSS, and fonts
			ext := filepath.Ext(path)
			if ext == ".js" || ext == ".css" || ext == ".woff" || ext == ".woff2" || ext == ".ttf" || ext == ".otf" {
				// Convert to web path: assets/index.js -> /assets/index.js
				webPath := "/" + filepath.ToSlash(path)
				http2PushResources = append(http2PushResources, webPath)
			}
			return nil
		})

		if len(http2PushResources) > 0 {
			fmt.Printf("HTTP/2 Push enabled for %d resources: %v\n", len(http2PushResources), http2PushResources)
		}
	}

	if useEmbed && assetsFS != nil {
		// PERFORMANCE: High compression for static assets (cached, compressed once)
		// Use LevelBestCompression for static files since they're cached and served repeatedly
		staticCompress := compress.New(compress.Config{
			Level: compress.LevelBestCompression, // Maximum compression for static assets
			Next: func(c *fiber.Ctx) bool {
				// Skip compression for API routes
				return strings.HasPrefix(c.Path(), "/api")
			},
		})

		// PERFORMANCE: HTTP/2 Server Push for critical resources (optional, config-controlled)
		// Push critical assets (JS, CSS, fonts) proactively to reduce round-trip latency
		// Uses Link header with rel=preload (works with both HTTP/1.1 and HTTP/2)
		// IMPORTANT: Only enabled if config.EnableHTTP2Push is true and resources are discovered
		if config.EnableHTTP2Push && len(http2PushResources) > 0 {
			app.Use("/", func(c *fiber.Ctx) error {
				// Only push on initial page load (HTML requests)
				// Skip for API routes and asset requests
				if c.Path() == "/" || (!strings.HasPrefix(c.Path(), "/api") && !strings.HasPrefix(c.Path(), "/assets")) {
					// Use Link header for HTTP/2 Server Push
					// This works with both HTTP/1.1 (as preload hint) and HTTP/2 (as server push)
					var links []string
					for _, resource := range http2PushResources {
						// Determine resource type from extension
						ext := filepath.Ext(resource)
						var asType string
						switch ext {
						case ".js":
							asType = "script"
						case ".css":
							asType = "style"
						case ".woff", ".woff2", ".ttf", ".otf":
							asType = "font"
						default:
							asType = "fetch"
						}
						links = append(links, fmt.Sprintf("<%s>; rel=preload; as=%s", resource, asType))
					}
					if len(links) > 0 {
						c.Set("Link", strings.Join(links, ", "))
					}
				}
				return c.Next()
			})
		}

		// Use embedded assets (for SDK users and production)
		app.Use("/", staticCompress, filesystem.New(filesystem.Config{
			Root:         http.FS(assetsFS),
			Browse:       false,
			Index:        "index.html",
			NotFoundFile: "index.html",
			MaxAge:       3600,
			Next: func(c *fiber.Ctx) bool {
				return strings.HasPrefix(c.Path(), "/api")
			},
		}))
	} else {
		// PERFORMANCE: High compression for static assets in development mode
		staticCompress := compress.New(compress.Config{
			Level: compress.LevelBestCompression,
			Next: func(c *fiber.Ctx) bool {
				return strings.HasPrefix(c.Path(), "/api")
			},
		})

		// Development mode with local path (for SDK developers only)
		if config.Storage.URL != "" && config.Storage.Path != "" {
			app.Use(config.Storage.URL, staticCompress)
			app.Static(config.Storage.URL, config.Storage.Path)
		} else {
			app.Use("/storage", staticCompress)
			app.Static("/storage", "./storage/public")
		}

		app.Use("/", staticCompress)
		app.Static("/", "./pkg/panel/ui")
		app.Get("*", func(c *fiber.Ctx) error {
			// Skip API routes
			if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
				return c.Next()
			}
			return c.SendFile(localUIPath)
		})
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
	}

	// Load Dynamic Settings
	_ = p.LoadSettings()

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

	// Register Pages from Config
	if p.Config.DashboardPage != nil {
		p.RegisterPage(p.Config.DashboardPage)
	} else {
		p.RegisterPage(&page.Dashboard{})
	}
	if p.Config.SettingsPage != nil {
		p.RegisterPage(p.Config.SettingsPage)
	} else {
		// Default Settings Page
		p.RegisterPage(&page.Settings{
			Elements: []fields.Element{
				fields.Text("Site Name", "site_name").
					Label("Site Name").
					Placeholder("Enter site name").
					Default("Panel.go").
					Required(),
				fields.Text("Site URL", "site_url").
					Label("Site URL").
					Placeholder("https://example.com").
					Required(),
				fields.Textarea("Site Description", "site_description").
					Label("Site Description").
					Placeholder("Enter site description").
					Rows(3),
				fields.Email("Contact Email", "contact_email").
					Label("Contact Email").
					Placeholder("contact@example.com"),
				fields.Tel("Contact Phone", "contact_phone").
					Label("Contact Phone").
					Placeholder("+90 555 123 4567"),
				fields.Textarea("Contact Address", "contact_address").
					Label("Contact Address").
					Placeholder("Enter contact address").
					Rows(2),
				fields.Switch("Register Enable", "register_enable").
					Label("User Registration").
					HelpText("Allow new users to register").
					Default(true),
				fields.Switch("Forgot Password Enable", "forgot_password_enable").
					Label("Forgot Password").
					HelpText("Enable password reset functionality").
					Default(false),
				fields.Switch("Maintenance Mode", "maintenance_mode").
					Label("Maintenance Mode").
					HelpText("Put the site in maintenance mode").
					Default(false),
				fields.Switch("Debug Mode", "debug_mode").
					Label("Debug Mode").
					HelpText("Enable debug mode (development only)").
					Default(false),
			},
		})
	}

	// Register Account Page
	if p.Config.AccountPage != nil {
		p.RegisterPage(p.Config.AccountPage)
	} else {
		// Default Account Page
		p.RegisterPage(&page.Account{
			Elements: []fields.Element{
				fields.Text("Name", "name").
					Label("Full Name").
					Placeholder("Enter your full name").
					Required(),
				fields.Email("Email", "email").
					Label("Email Address").
					Placeholder("your@email.com").
					Required(),
				fields.Image("Image", "image").
					Label("Profile Picture").
					HelpText("Upload your profile picture"),
				fields.Password("Current Password", "current_password").
					Label("Current Password").
					Placeholder("Enter current password").
					HelpText("Required to change password"),
				fields.Password("New Password", "new_password").
					Label("New Password").
					Placeholder("Enter new password").
					HelpText("Leave blank to keep current password"),
				fields.Password("Confirm Password", "confirm_password").
					Label("Confirm Password").
					Placeholder("Confirm new password"),
				fields.Switch("Email Notifications", "email_notifications").
					Label("Email Notifications").
					HelpText("Receive notifications via email").
					Default(true),
				fields.Switch("SMS Notifications", "sms_notifications").
					Label("SMS Notifications").
					HelpText("Receive notifications via SMS").
					Default(false),
				fields.Select("Language", "language").
					Label("Language").
					Placeholder("Select language").
					Options(map[string]string{
						"en": "English",
						"tr": "Türkçe",
					}).
					Default("en"),
				fields.Select("Theme", "theme").
					Label("Theme").
					Placeholder("Select theme").
					Options(map[string]string{
						"light": "Light",
						"dark":  "Dark",
						"auto":  "Auto",
					}).
					Default("auto"),
			},
		})
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
			FailureThreshold:       failureThreshold,
			Timeout:                timeout,
			SuccessThreshold:       successThreshold,
			HalfOpenMaxConcurrent:  halfOpenMaxConcurrent,
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

	// Auth Routes
	authRoutes := api.Group("/auth")
	// SECURITY: Strict rate limiting for authentication endpoints (10 req/min)
	// TEMPORARY: Rate limiting disabled for development
	// authRoutes.Use(middleware.AuthRateLimiter())
	authRoutes.Post("/sign-in/email", context.Wrap(authH.LoginEmail))
	authRoutes.Post("/sign-up/email", context.Wrap(authH.RegisterEmail))
	authRoutes.Post("/sign-out", context.Wrap(authH.SignOut))
	authRoutes.Post("/forgot-password", context.Wrap(authH.ForgotPassword))
	authRoutes.Get("/session", context.Wrap(authH.GetSession))

	api.Get("/init", context.Wrap(p.handleInit)) // App Initialization

	// Middleware
	api.Use(context.Wrap(authH.SessionMiddleware))
	// SECURITY: Rate limiting for general API endpoints (100 req/min)
	// TEMPORARY: Rate limiting disabled for development
	// api.Use(middleware.APIRateLimiter())

	// Page Routes
	api.Get("/pages", context.Wrap(p.handlePages))
	api.Get("/pages/:slug", context.Wrap(p.handlePageDetail))
	api.Post("/pages/:slug", context.Wrap(p.handlePageSave))

	api.Get("/resource/:resource/cards", context.Wrap(p.handleResourceCards))
	api.Get("/resource/:resource/cards/:index", context.Wrap(p.handleResourceCard))
	api.Get("/resource/:resource/lenses", context.Wrap(p.handleResourceLenses))                  // List available lenses
	api.Get("/resource/:resource/lens/:lens", context.Wrap(p.handleResourceLens))                // Lens data
	api.Get("/resource/:resource/morphable/:field", context.Wrap(p.handleMorphable))             // MorphTo field options
	api.Get("/resource/:resource/actions", context.Wrap(p.handleResourceActions))                // List available actions
	api.Post("/resource/:resource/actions/:action", context.Wrap(p.handleResourceActionExecute)) // Execute action
	api.Get("/resource/:resource", context.Wrap(p.handleResourceIndex))
	api.Post("/resource/:resource", context.Wrap(p.handleResourceStore))
	api.Get("/resource/:resource/create", context.Wrap(p.handleResourceCreate)) // New Route
	api.Get("/resource/:resource/:id", context.Wrap(p.handleResourceShow))
	api.Get("/resource/:resource/:id/detail", context.Wrap(p.handleResourceDetail))
	api.Get("/resource/:resource/:id/edit", context.Wrap(p.handleResourceEdit))
	api.Post("/resource/:resource/:id/fields/:field/resolve", context.Wrap(p.handleFieldResolve))          // Field resolver endpoint
	api.Post("/resource/:resource/fields/resolve-dependencies", context.Wrap(p.handleResolveDependencies)) // Dependency resolver endpoint

	// Hover card resolver endpoints - Support GET, POST, PATCH, DELETE
	api.Get("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve))    // Hover card resolver (GET)
	api.Post("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve))   // Hover card resolver (POST)
	api.Patch("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve))  // Hover card resolver (PATCH)
	api.Delete("/resource/:resource/resolver/:field", context.Wrap(p.handleHoverCardResolve)) // Hover card resolver (DELETE)
	api.Put("/resource/:resource/:id", context.Wrap(p.handleResourceUpdate))
	api.Delete("/resource/:resource/:id", context.Wrap(p.handleResourceDestroy))
	api.Get("/navigation", context.Wrap(p.handleNavigation)) // Sidebar Navigation

	// Notification Routes
	notificationService := notification.NewService(db)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	api.Get("/notifications", context.Wrap(notificationHandler.HandleGetUnreadNotifications))
	api.Post("/notifications/:id/read", context.Wrap(notificationHandler.HandleMarkAsRead))
	api.Post("/notifications/read-all", context.Wrap(notificationHandler.HandleMarkAllAsRead))

	// OpenAPI Routes
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

	// /resolve endpoint for dynamic routing check
	api.Get("/resolve", context.Wrap(p.handleResolve))

	return p
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
		case "register":
			if v, ok := val.(bool); ok {
				config.Register = v
				p.Config.Features.Register = v
			}
		case "forgot_password":
			if v, ok := val.(bool); ok {
				config.ForgotPassword = v
				p.Config.Features.ForgotPassword = v
			}
		}
	}
	return nil
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
	for _, l := range res.Lenses() {
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

	// Use the lens controller
	return handler.HandleLens(h, c)
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
	}

	items := []NavItem{}
	for slug, res := range p.resources {
		if !res.Visible() {
			continue
		}
		items = append(items, NavItem{
			Slug:  slug,
			Title: res.Title(),
			Icon:  res.Icon(),
			Group: res.Group(),
			Type:  "resource",
			Order: res.NavigationOrder(),
		})
	}

	for slug, pg := range p.pages {
		if !pg.Visible() {
			continue
		}
		items = append(items, NavItem{
			Slug:  slug,
			Title: pg.Title(),
			Icon:  pg.Icon(),
			Group: pg.Group(),
			Type:  "page",
			Order: pg.NavigationOrder(),
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
func (p *Panel) handleInit(c *context.Context) error {
	fmt.Printf("DEBUG: handleInit called. Config: %+v\n", p.Config)
	fmt.Printf("DEBUG: SettingsValues: %+v\n", p.Config.SettingsValues)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in handleInit:", r)
		}
	}()

	// Get CSRF token from context (set by CSRF middleware) and send in response header
	// This allows JavaScript to read the token (since the cookie is HTTPOnly)
	if csrfToken := c.Locals("csrf"); csrfToken != nil {
		c.Set("X-CSRF-Token", csrfToken.(string))
	}

	// Get features from settings or use config defaults
	registerEnabled := p.Config.Features.Register
	forgotPasswordEnabled := p.Config.Features.ForgotPassword

	// Check if settings have override values
	if settings := p.Config.SettingsValues.Values; settings != nil {
		if registerVal, ok := settings["register"]; ok {
			if boolVal, ok := registerVal.(bool); ok {
				registerEnabled = boolVal
			}
		}
		if forgotVal, ok := settings["forgot_password"]; ok {
			if boolVal, ok := forgotVal.(bool); ok {
				forgotPasswordEnabled = boolVal
			}
		}
	}

	return c.JSON(fiber.Map{
		"features": fiber.Map{
			"register":        registerEnabled,
			"forgot_password": forgotPasswordEnabled,
		},
		"oauth": fiber.Map{
			"google": p.Config.OAuth.Google.Enabled(),
		},
		"version":  "1.0.0",
		"settings": p.Config.SettingsValues.Values,
	})
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
				"title": res.Title(),
				"icon":  res.Icon(),
				"group": res.Group(),
			},
		})
	}

	// Future: Check custom pages, database driven pages etc.

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Page not found",
	})
}

