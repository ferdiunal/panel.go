package server

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"panel.go/internal/ent"
	"panel.go/internal/ent/migrate"
	"panel.go/internal/interfaces/handler"
	"panel.go/internal/repository"
	"panel.go/internal/service"
	"panel.go/shared/encrypt"
	"panel.go/shared/uuid"

	entsql "entgo.io/ent/dialect/sql"
	db "github.com/ferdiunal/go.utils/database"
	goutils "github.com/ferdiunal/go.utils/database/interfaces"
)

type FiberServer struct {
	*fiber.App

	Db      goutils.DatabaseService
	Ent     *ent.Client
	Store   *session.Store
	Encrypt encrypt.Crypt
	Service *handler.Services
}

func New() *FiberServer {
	db, err := db.New()
	if err != nil {
		log.Fatal("Failed to initialize database")
	}

	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	drv := entsql.OpenDB("postgres", db.Db())
	client := ent.NewClient(ent.Driver(drv))

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithForeignKeys(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	store := session.New()

	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	encrypt := encrypt.NewCrypt(encryptionKey)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "panel.go",
			AppName:      "panel.go",
		}),

		Db: db,

		Store: store,

		Ent: client,

		Encrypt: encrypt,
	}

	appKey := os.Getenv("APP_KEY")
	headerName := "X-Csrf-Token"
	cookieName := "__Host-csrf_"
	server.RegisterServices()
	server.App.Use(logger.New())
	server.App.Use(helmet.New())

	server.App.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		},
		Storage: server.Store.Storage,
	}))

	server.App.Use(encryptcookie.New(encryptcookie.Config{
		Key: appKey,
		Encryptor: func(decryptedString, key string) (string, error) {
			return server.Encrypt.Encrypt(decryptedString)
		},
		Decryptor: func(encryptedString, key string) (string, error) {
			return server.Encrypt.Decrypt(encryptedString)
		},
	}))

	server.App.Use(csrf.New(csrf.Config{
		KeyLookup:         "header:" + headerName,
		CookieName:        cookieName,
		CookieSameSite:    "Lax",
		CookieSecure:      true,
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		ContextKey:        "csrf",
		KeyGenerator:      func() string { return uuid.NewUUID().String() },
		SingleUseToken:    true,
		Extractor:         csrf.CsrfFromCookie(cookieName),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			c.Set("HX-Refresh", "true")
			return c.Status(419).SendString("Page Expired")
		},
	}))

	return server
}

func (s *FiberServer) RegisterServices() {
	s.Service = &handler.Services{
		AuthService: service.NewAuthService(
			repository.NewAccountRepository(s.Ent),
			repository.NewUserRepository(s.Ent),
			repository.NewSessionRepository(s.Ent),
			s.Encrypt,
		),
	}
}
