package panel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	"github.com/ferdiunal/panel.go/pkg/page"
	"github.com/ferdiunal/panel.go/pkg/permission"
	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceUser "github.com/ferdiunal/panel.go/pkg/resource/user"
	"github.com/ferdiunal/panel.go/pkg/service/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"gorm.io/gorm"
)

// Panel, uygulamanın ana yapısıdır.
// Veritabanı, Fiber uygulaması, yetkilendirme servisi ve kayıtlı kaynakları tutar.
type Panel struct {
	Config    Config
	Db        *gorm.DB
	Fiber     *fiber.App
	Auth      *auth.Service
	resources map[string]resource.Resource
	pages     map[string]page.Page
}

// New, yeni bir Panel örneği oluşturur ve başlatır.
// Veritabanı migration'ları, middleware kayıtları ve yönlendirmeleri yapılandırır.
func New(config Config) *Panel {
	app := fiber.New()
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
	app.Use(compress.New())

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
	if config.Environment != "test" {
		app.Use(csrf.New(csrf.Config{
			KeyLookup:      "header:X-CSRF-Token",
			CookieName:     "__Host-csrf-token",
			CookieSecure:   config.Environment == "production",
			CookieHTTPOnly: true,
			CookieSameSite: "Strict",
			Expiration:     24 * time.Hour,
		}))
	}

	app.Use(earlydata.New())
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

	if useEmbed && assetsFS != nil {
		// Use embedded assets (for SDK users and production)
		app.Use("/", filesystem.New(filesystem.Config{
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
		// Development mode with local path (for SDK developers only)
		if config.Storage.URL != "" && config.Storage.Path != "" {
			app.Static(config.Storage.URL, config.Storage.Path)
		} else {
			app.Static("/storage", "./storage/public")
		}

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
				fields.Text("Site Name").Default("Panel.go"),
				fields.Switch("Register").Default(true),
				fields.Switch("Forgot Password").Default(false),
			},
		})
	}

	// Register Dynamic Routes
	// /api/resource/:resource -> List/Index
	// /api/resource/:resource/:id -> Detail/Show/Update/Delete

	api := app.Group("/api")

	// Auth Routes
	authRoutes := api.Group("/auth")
	// SECURITY: Strict rate limiting for authentication endpoints (10 req/min)
	authRoutes.Use(middleware.AuthRateLimiter())
	authRoutes.Post("/sign-in/email", context.Wrap(authH.LoginEmail))
	authRoutes.Post("/sign-up/email", context.Wrap(authH.RegisterEmail))
	authRoutes.Post("/sign-out", context.Wrap(authH.SignOut))
	authRoutes.Post("/forgot-password", context.Wrap(authH.ForgotPassword))
	authRoutes.Get("/session", context.Wrap(authH.GetSession))

	api.Get("/init", context.Wrap(p.handleInit)) // App Initialization

	// Middleware
	api.Use(context.Wrap(authH.SessionMiddleware))
	// SECURITY: Rate limiting for general API endpoints (100 req/min)
	api.Use(middleware.APIRateLimiter())

	// Page Routes
	api.Get("/pages", context.Wrap(p.handlePages))
	api.Get("/pages/:slug", context.Wrap(p.handlePageDetail))
	api.Post("/pages/:slug", context.Wrap(p.handlePageSave))

	api.Get("/resource/:resource/cards", context.Wrap(p.handleResourceCards))
	api.Get("/resource/:resource/cards/:index", context.Wrap(p.handleResourceCard))
	api.Get("/resource/:resource/lenses", context.Wrap(p.handleResourceLenses))      // List available lenses
	api.Get("/resource/:resource/lens/:lens", context.Wrap(p.handleResourceLens))    // Lens data
	api.Get("/resource/:resource/morphable/:field", context.Wrap(p.handleMorphable)) // MorphTo field options
	api.Get("/resource/:resource/actions", context.Wrap(p.handleResourceActions))    // List available actions
	api.Post("/resource/:resource/actions/:action", context.Wrap(p.handleResourceActionExecute)) // Execute action
	api.Get("/resource/:resource", context.Wrap(p.handleResourceIndex))
	api.Post("/resource/:resource", context.Wrap(p.handleResourceStore))
	api.Get("/resource/:resource/create", context.Wrap(p.handleResourceCreate)) // New Route
	api.Get("/resource/:resource/:id", context.Wrap(p.handleResourceShow))
	api.Get("/resource/:resource/:id/detail", context.Wrap(p.handleResourceDetail))
	api.Get("/resource/:resource/:id/edit", context.Wrap(p.handleResourceEdit))
	api.Post("/resource/:resource/:id/fields/:field/resolve", context.Wrap(p.handleFieldResolve)) // Field resolver endpoint
	api.Put("/resource/:resource/:id", context.Wrap(p.handleResourceUpdate))
	api.Delete("/resource/:resource/:id", context.Wrap(p.handleResourceDestroy))
	api.Get("/navigation", context.Wrap(p.handleNavigation)) // Sidebar Navigation

	// Notification Routes
	notificationService := notification.NewService(db)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	api.Get("/notifications", context.Wrap(notificationHandler.HandleGetUnreadNotifications))
	api.Post("/notifications/:id/read", context.Wrap(notificationHandler.HandleMarkAsRead))
	api.Post("/notifications/read-all", context.Wrap(notificationHandler.HandleMarkAllAsRead))

	// SSE Notification Stream
	notificationSSEHandler := handler.NewNotificationSSEHandler(db)
	api.Get("/notifications/stream", context.Wrap(notificationSSEHandler.HandleNotificationStream))

	// /resolve endpoint for dynamic routing check
	api.Get("/resolve", context.Wrap(p.handleResolve))

	return p
}

// LoadSettings, veritabanından ayarları okur ve yapılandırmayı günceller.
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

func (p *Panel) Register(slug string, res resource.Resource) {
	res.SetDialogType(resource.DialogTypeSheet)
	p.resources[slug] = res
}

func (p *Panel) RegisterResource(res resource.Resource) {
	p.Register(res.Slug(), res)
}

func (p *Panel) RegisterPage(pg page.Page) {
	p.pages[pg.Slug()] = pg
}

func (p *Panel) Start() error {
	addr := fmt.Sprintf("%s:%s", p.Config.Server.Host, p.Config.Server.Port)
	return p.Fiber.Listen(addr)
}

// Helper to resolve resource and create handler
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

func (p *Panel) handleResourceIndex(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceIndex(h, c)
	})
}

func (p *Panel) handleResourceShow(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceShow(h, c)
	})
}

func (p *Panel) handleResourceDetail(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceDetail(h, c)
	})
}

func (p *Panel) handleResourceStore(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceStore(h, c)
	})
}

func (p *Panel) handleResourceCreate(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceCreate(h, c)
	})
}

func (p *Panel) handleResourceUpdate(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceUpdate(h, c)
	})
}

func (p *Panel) handleResourceDestroy(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceDestroy(h, c)
	})
}

func (p *Panel) handleResourceEdit(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleResourceEdit(h, c)
	})
}

func (p *Panel) handleFieldResolve(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleFieldResolve(h, c)
	})
}

func (p *Panel) handleResourceCards(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleCardList(h, c)
	})
}

func (p *Panel) handleResourceCard(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleCardDetail(h, c)
	})
}

func (p *Panel) handleResourceLenses(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleLensIndex(h, c)
	})
}

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

func (p *Panel) handleMorphable(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleMorphable(h, c)
	})
}

func (p *Panel) handleResourceActions(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleActionList(h, c)
	})
}

func (p *Panel) handleResourceActionExecute(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return handler.HandleActionExecute(h, c)
	})
}

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

func (p *Panel) handleInit(c *context.Context) error {
	fmt.Printf("DEBUG: handleInit called. Config: %+v\n", p.Config)
	fmt.Printf("DEBUG: SettingsValues: %+v\n", p.Config.SettingsValues)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in handleInit:", r)
		}
	}()

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
