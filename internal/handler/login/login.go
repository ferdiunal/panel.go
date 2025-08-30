package login

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"panel.go/cmd/web"
	"panel.go/internal/interfaces/handler"
	"panel.go/internal/service"
	"panel.go/shared/validate"
)

func Get(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		return handler.View(c, "Giriş Yap", web.LoginForm)
	}
}

func Post(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		var request service.LoginPayload

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

		fmt.Println(c.Context())

		token, err := options.Service.AuthService.Login(
			c,
			request.Email,
			request.Password,
		)

		if err != nil {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="text-error text-sm mt-2">
					<p>` + err.Error() + `</p>
				</div>
			`)
		}

		// Set cookie with expire time
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    token.AccessToken,
			Expires:  token.ExpiresAt,
			HTTPOnly: true,
			Secure:   options.Prod,
			SameSite: "Lax",
			Path:     "/",
		})

		// Success - redirect to dashboard
		c.Set("HX-Redirect", "/dashboard")
		return c.SendStatus(fiber.StatusOK)
	}
}
