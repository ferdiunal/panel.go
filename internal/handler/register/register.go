package register

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"panel.go/cmd/web"
	_err "panel.go/internal/errors"
	"panel.go/internal/interfaces/handler"
	"panel.go/internal/service"
	"panel.go/shared/validate"
)

func Get(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		return handler.View(c, "Kayıt Ol", web.RegisterForm)
	}
}

func Post(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		var request service.RegisterPayload

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
				case "name":
					errorHTML += `<p>Adınız: ` + actualMessage + `</p>`
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

		success, err := options.Service.AuthService.Register(
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

		if !success {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			if errors.Is(err, _err.ErrUserExists) {
				return c.Status(fiber.StatusUnprocessableEntity).SendString(`
					<div class="error">
						<p>E-posta zaten kullanılıyor</p>
					</div>
				`)
			}

			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					<p>Bir hata oluştu</p>
				</div>
			`)
		}

		// Success - redirect to dashboard
		c.Set("HX-Redirect", "/giris")
		return c.Status(fiber.StatusOK).SendString(`
			<div class="success">
				<p>Kayıt başarılı</p>
			</div>
		`)
	}
}
