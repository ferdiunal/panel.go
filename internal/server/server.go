package server

import (
	"github.com/gofiber/fiber/v2"

	"panel.go/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "panel.go",
			AppName:      "panel.go",
		}),

		db: database.New(),
	}

	return server
}
