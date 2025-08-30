package login

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"panel.go/cmd/web"
	"panel.go/internal/constants"
	_err "panel.go/internal/errors"
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
				<div class="error">
					<p>Geçersiz form verisi</p>
				</div>
			`)
		}

		// Validate struct
		_errors := validate.ValidateStruct(request)
		if len(_errors) > 0 {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")

			var errorHTML string
			for field, messageMap := range _errors {
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

		token, err := options.Service.AuthService.Login(
			c,
			&request,
		)

		if errors.Is(err, _err.ErrAuthentication) {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					<p>E-posta veya parola hatalı</p>
				</div>
			`)
		}

		// Set cookie with expire time
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    token.AccessToken,
			Expires:  token.ExpiresAt,
			HTTPOnly: true,
			Secure:   constants.APP_ENV == "production",
			SameSite: "Lax",
			Path:     "/",
		})

		// Success - redirect to dashboard
		c.Set("HX-Redirect", "/dashboard")
		return c.SendStatus(fiber.StatusOK)
	}
}
