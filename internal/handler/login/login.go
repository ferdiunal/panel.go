package login

import (
	"github.com/gofiber/fiber/v2"
	"panel.go/cmd/web"
	"panel.go/internal/interfaces/handler"
	"panel.go/shared/validate"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func Get(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		return handler.View(c, "Giriş Yap", web.LoginForm)
	}
}

func Post(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		var request LoginRequest

		// Parse form data
		if err := c.BodyParser(&request); err != nil {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="text-error text-sm mt-2">
					<p>Geçersiz form verisi</p>
				</div>
			`)
		}

		// Validate struct
		errors := validate.ValidateStruct(request)
		if len(errors) > 0 {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")

			var errorHTML string
			for field, messageMap := range errors {
				actualMessage := messageMap["message"]
				switch field {
				case "email":
					errorHTML += `<p>E-posta: ` + actualMessage + `</p>`
				case "password":
					errorHTML += `<p>Parola: ` + actualMessage + `</p>`
				default:
					errorHTML += `<p>` + actualMessage + `</p>`
				}
			}

			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					` + errorHTML + `
				</div>
			`)
		}

		// Simulate authentication (replace with real auth logic)
		if request.Email != "test@test.com" || request.Password != "password" {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnauthorized).SendString(`
				<div class="error">
					<p>E-posta veya parola hatalı</p>
				</div>
			`)
		}

		// Success - redirect to dashboard
		c.Set("HX-Redirect", "/dashboard")
		return c.SendStatus(fiber.StatusOK)
	}
}
