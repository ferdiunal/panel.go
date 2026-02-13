package panel

import (
	"fmt"
	"sort"
	"time"

	"github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/data/orm"
	"github.com/ferdiunal/panel.go/internal/domain/account"
	"github.com/ferdiunal/panel.go/internal/domain/audit"
	"github.com/ferdiunal/panel.go/internal/domain/session"
	"github.com/ferdiunal/panel.go/internal/domain/setting"
	"github.com/ferdiunal/panel.go/internal/domain/user"
	"github.com/ferdiunal/panel.go/internal/domain/verification"
	"github.com/ferdiunal/panel.go/internal/handler"
	authHandler "github.com/ferdiunal/panel.go/internal/handler/auth"
	obs "github.com/ferdiunal/panel.go/internal/observability"
	"github.com/ferdiunal/panel.go/internal/page"
	"github.com/ferdiunal/panel.go/internal/resource"
	resourceAudit "github.com/ferdiunal/panel.go/internal/resource/audit"
	resourceUser "github.com/ferdiunal/panel.go/internal/resource/user"
	"github.com/ferdiunal/panel.go/internal/service/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
	if db == nil {
		panic("panel: database instance is required")
	}
	if (config.Storage.Path == "" && config.Storage.URL != "") || (config.Storage.Path != "" && config.Storage.URL == "") {
		panic("panel: storage path and storage URL must be set together")
	}
	if config.Environment == "" {
		config.Environment = "development"
	}

	// Auth Components
	userRepo := orm.NewUserRepository(db)
	sessionRepo := orm.NewSessionRepository(db)
	accountRepo := orm.NewAccountRepository(db)

	authService := auth.NewService(db, userRepo, sessionRepo, accountRepo)
	authH := authHandler.NewHandler(authService)

	// Auto Migrate Auth Domains
	db.AutoMigrate(&user.User{}, &session.Session{}, &account.Account{}, &verification.Verification{}, &setting.Setting{}, &audit.Log{})

	// Middleware Registration
	app.Use(requestid.New())
	app.Use(obs.RequestMetricsMiddleware())
	app.Use(obs.RequestLoggerMiddleware())
	app.Use(compress.New())
	app.Use(cors.New())
	if config.Environment == "production" {
		app.Use(csrf.New())
		// app.Use(circuitbreaker.New(circuitbreaker.Config{
		// 	FailureThreshold: 3,
		// }))
	}
	app.Use(earlydata.New())
	app.Use(etag.New())
	app.Use(helmet.New(helmet.Config{
		CrossOriginResourcePolicy: "cross-origin",
	}))
	// Static file serving using configured paths
	// We need config first, but config is in Config struct which is in Panel, but Panel not created yet.
	// Config IS available as 'config' argument to New().
	if config.Storage.URL != "" && config.Storage.Path != "" {
		app.Static(config.Storage.URL, config.Storage.Path)
	} else {
		// Fallback or explicit check in main.go ensures they are set.
		app.Static("/storage", "./storage/public")
	}

	app.Get("/health", obs.HealthHandler())
	app.Get("/ready", obs.ReadyHandler(db))
	app.Get("/metrics", obs.MetricsHandler())

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

	// Register Audit Log Resource (readonly, admin only)
	p.RegisterResource(resourceAudit.GetAuditLogResource())

	// Register Pages from Config
	if p.Config.DashboardPage != nil {
		p.RegisterPage(p.Config.DashboardPage)
	}
	if p.Config.SettingsPage != nil {
		p.RegisterPage(p.Config.SettingsPage)
	}

	// Register Dynamic Routes
	// /api/resource/:resource -> List/Index
	// /api/resource/:resource/:id -> Detail/Show/Update/Delete

	api := app.Group("/api")

	// Auth Routes
	authRoutes := api.Group("/auth")
	authLoginLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "too many login requests",
			})
		},
	})
	authRoutes.Post("/sign-in/email", authLoginLimiter, context.Wrap(authH.LoginEmail))
	authRoutes.Post("/sign-up/email", context.Wrap(authH.RegisterEmail))
	authRoutes.Post("/sign-out", context.Wrap(authH.SignOut))
	authRoutes.Get("/session", context.Wrap(authH.GetSession))

	api.Get("/init", context.Wrap(p.handleInit)) // App Initialization

	// Middleware
	api.Use(context.Wrap(authH.SessionMiddleware))
	api.Use(context.Wrap(obs.AuditMiddleware(db)))

	// Page Routes
	api.Get("/pages", context.Wrap(p.handlePages))
	api.Get("/pages/:slug", context.Wrap(p.handlePageDetail))
	api.Post("/pages/:slug", context.Wrap(p.handlePageSave))

	api.Get("/resource/:resource/cards", context.Wrap(p.handleResourceCards))
	api.Get("/resource/:resource/cards/:index", context.Wrap(p.handleResourceCard))
	api.Get("/resource/:resource", context.Wrap(p.handleResourceIndex))
	api.Post("/resource/:resource", context.Wrap(p.handleResourceStore))
	api.Get("/resource/:resource/create", context.Wrap(p.handleResourceCreate)) // New Route
	api.Get("/resource/:resource/:id", context.Wrap(p.handleResourceShow))
	api.Get("/resource/:resource/:id/detail", context.Wrap(p.handleResourceDetail))
	api.Get("/resource/:resource/:id/edit", context.Wrap(p.handleResourceEdit))
	api.Put("/resource/:resource/:id", context.Wrap(p.handleResourceUpdate))
	api.Delete("/resource/:resource/:id", context.Wrap(p.handleResourceDestroy))
	api.Get("/resource/:resource/lens/:lens", context.Wrap(p.handleResourceLens)) // Lens Route
	api.Get("/navigation", context.Wrap(p.handleNavigation))                      // Sidebar Navigation

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
		val := s.Value["value"]
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
		return h.Index(c)
	})
}

func (p *Panel) handleResourceShow(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.Show(c)
	})
}

func (p *Panel) handleResourceDetail(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.Detail(c)
	})
}

func (p *Panel) handleResourceStore(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.Store(c)
	})
}

func (p *Panel) handleResourceCreate(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.Create(c)
	})
}

func (p *Panel) handleResourceUpdate(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.Update(c)
	})
}

func (p *Panel) handleResourceDestroy(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.Destroy(c)
	})
}

func (p *Panel) handleResourceEdit(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.Edit(c)
	})
}

func (p *Panel) handleResourceCards(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.ListCards(c)
	})
}

func (p *Panel) handleResourceCard(c *context.Context) error {
	return p.withResourceHandler(c, func(h *handler.FieldHandler) error {
		return h.GetCard(c)
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

	// Lens inherently implies a List/Index view
	return h.Index(c)
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
	return c.JSON(fiber.Map{
		"features": fiber.Map{
			"register":        p.Config.Features.Register,
			"forgot_password": p.Config.Features.ForgotPassword,
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
