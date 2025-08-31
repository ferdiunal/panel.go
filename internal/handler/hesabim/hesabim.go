package hesabim

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"panel.go/cmd/web"
	_err "panel.go/internal/errors"
	"panel.go/internal/interfaces/handler"
	"panel.go/internal/repository"
	"panel.go/shared/validate"
)

func Get(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		return handler.WithAuthView(c, "Hesabım", web.Hesabim)
	}
}

func Update(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		var request repository.UserCreatePayload

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
				case "name":
					errorHTML += `<p>Ad: ` + actualMessage + `</p>`
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

		success, err := options.Service.AuthService.UpdateAccount(c, &request)

		if errors.Is(err, _err.ErrUniqueEmail) {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					<p>E-posta zaten kullanılıyor</p>
				</div>
			`)
		}

		if err != nil {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusInternalServerError).SendString(`
				<div class="error">
					<p>Güncelleme sırasında bir hata oluştu</p>
				</div>
			`)
		}

		if !success {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					<p>Hesap güncellenemedi</p>
				</div>
			`)
		}

		c.Set("HX-Retarget", "#save-status")
		c.Set("HX-Reswap", "innerHTML")

		return c.SendString(`
			<div class="badge badge-success gap-2 animate-pulse">
				<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-4 h-4 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
				Kaydedildi
			</div>
			<script>
				setTimeout(function() {
					document.getElementById('save-status').innerHTML = '';
				}, 3000);
			</script>
		`)
	}
}
