package observability

import (
	stdContext "context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	}
}

func ReadyHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if db == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not_ready",
				"error":  "database is nil",
			})
		}

		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not_ready",
				"error":  err.Error(),
			})
		}

		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not_ready",
				"error":  err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "ready",
		})
	}
}
