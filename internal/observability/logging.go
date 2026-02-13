package observability

import (
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

var requestLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func RequestLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		requestLogger.Info("http_request",
			"request_id", c.GetRespHeader(fiber.HeaderXRequestID),
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP(),
			"user_agent", c.Get(fiber.HeaderUserAgent),
		)
		return err
	}
}
