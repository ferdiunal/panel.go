package web

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func LoginHandler(c *fiber.Ctx) error {
	// Parse form data
	if err := c.BodyParser(c); err != nil {
		innerErr := c.Status(fiber.StatusBadRequest).SendString("Bad Request")
		if innerErr != nil {
			log.Fatalf("Could not send error in HelloWebHandler: %e", innerErr)
		}
	}

	return nil
}
