package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"panel.go/cmd/web"
	"panel.go/internal/handler/login"
	"panel.go/internal/interfaces/handler"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-Csrf-Token",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	handleOptions := &handler.Options{
		Ent:   s.Ent,
		Store: s.Store,
	}

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	s.App.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	s.App.Get("/login", login.Get(handleOptions))
	s.App.Post("/login", login.Post(handleOptions))
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.Db.Health())
}
