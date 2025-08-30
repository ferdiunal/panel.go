package hesabim

import (
	"github.com/gofiber/fiber/v2"
	"panel.go/cmd/web"
	"panel.go/internal/interfaces/handler"
)

func Get(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		return handler.WithAuthView(c, "Hesabım", web.Hesabim)
	}
}
