package guvenlik

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"panel.go/cmd/web"
	_err "panel.go/internal/errors"
	"panel.go/internal/interfaces/handler"
	"panel.go/shared/validate"
)

type ChangePasswordPayload struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8,max=32"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=32"`
}

func Get(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		return handler.WithAuthView(c, "Güvenlik", web.Guvenlik)
	}
}

func Update(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		var request ChangePasswordPayload

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

		fmt.Println(request)

		// Check if passwords match
		if request.NewPassword != request.ConfirmPassword {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					<p>Yeni parolalar eşleşmiyor</p>
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
				fmt.Println(field, actualMessage)
				switch field {
				case "currentPassword":
					errorHTML += `<p>Mevcut Parola: ` + actualMessage + `</p>`
				case "newPassword":
					errorHTML += `<p>Yeni Parola: ` + actualMessage + `</p>`
				case "confirmPassword":
					errorHTML += `<p>Parola Tekrar: ` + actualMessage + `</p>`
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

		success, err := options.Service.AuthService.ChangePassword(c, request.CurrentPassword, request.NewPassword)

		if errors.Is(err, _err.ErrAuthentication) {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					<p>Mevcut parola yanlış</p>
				</div>
			`)
		}

		if err != nil {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusInternalServerError).SendString(`
				<div class="error">
					<p>Parola değiştirme sırasında bir hata oluştu</p>
				</div>
			`)
		}

		if !success {
			c.Set("HX-Retarget", "#message")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusUnprocessableEntity).SendString(`
				<div class="error">
					<p>Parola değiştirilemedi</p>
				</div>
			`)
		}

		c.Set("HX-Retarget", "#save-status")
		c.Set("HX-Reswap", "innerHTML")

		return c.SendString(`
			<div class="success">
				Parola Değiştirildi
			</div>
			<script>
				setTimeout(function() {
					document.getElementById('save-status').innerHTML = '';
					// Clear form fields
					document.getElementById('currentPassword').value = '';
					document.getElementById('newPassword').value = '';
					document.getElementById('confirmPassword').value = '';
				}, 3000);
			</script>
		`)
	}
}
