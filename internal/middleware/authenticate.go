package middleware

import (
	"github.com/gofiber/fiber/v2"
	"panel.go/internal/service"
)

func Authenticate(service *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := service.VerifyToken(c, c.Cookies("access_token"))
		if err != nil {
			return c.Redirect("/giris")
		}
		return c.Next()
	}
}
