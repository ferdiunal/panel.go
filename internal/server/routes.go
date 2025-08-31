package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"panel.go/cmd/web"
	"panel.go/internal/handler/avatar"
	"panel.go/internal/handler/dashboard"
	"panel.go/internal/handler/guvenlik"
	"panel.go/internal/handler/hesabim"
	"panel.go/internal/handler/login"
	"panel.go/internal/handler/register"
	"panel.go/internal/interfaces/handler"
	"panel.go/internal/middleware"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080,http://127.0.0.1:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-Csrf-Token",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	handleOptions := &handler.Options{
		Store:   s.Store,
		Service: s.Service,
	}

	s.App.Get("/health", s.healthHandler)

	s.App.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	s.App.Get("/giris", login.Get(handleOptions))
	s.App.Post("/giris", login.Post(handleOptions))

	s.App.Get("/avatar/:avatar", avatar.Get(handleOptions))

	s.App.Get("/kayit", register.Get(handleOptions))
	s.App.Post("/kayit", register.Post(handleOptions))

	s.App.Delete("/logout", middleware.Authenticate(s.Service.AuthService), func(c *fiber.Ctx) error {
		err := s.Service.AuthService.Logout(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		c.Set("HX-Redirect", "/giris")
		return c.SendStatus(fiber.StatusNoContent)
	})
	s.App.Get("/dashboard", middleware.Authenticate(s.Service.AuthService), dashboard.Get(handleOptions))
	s.App.Get("/hesabim", middleware.Authenticate(s.Service.AuthService), hesabim.Get(handleOptions))
	s.App.Put("/hesabim", middleware.Authenticate(s.Service.AuthService), hesabim.Update(handleOptions))
	s.App.Get("/guvenlik", middleware.Authenticate(s.Service.AuthService), guvenlik.Get(handleOptions))
	s.App.Put("/guvenlik", middleware.Authenticate(s.Service.AuthService), guvenlik.Update(handleOptions))
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.Db.Health())
}
