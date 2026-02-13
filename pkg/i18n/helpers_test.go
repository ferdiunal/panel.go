package i18n

import (
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	// i18n middleware'ini ekle
	app.Use(fiberi18n.New(&fiberi18n.Config{
		RootPath:         "../../locales",
		AcceptLanguages:  []language.Tag{language.Turkish, language.English},
		DefaultLanguage:  language.Turkish,
		FormatBundleFile: "yaml",
	}))

	return app
}

func TestTrans(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		message := Trans(c, "welcome")
		return c.SendString(message)
	})

	req, _ := app.Test(fiber.NewRequest("GET", "/test?lang=tr"))
	assert.Equal(t, 200, req.StatusCode)
}

func TestTransWithTemplate(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		message := Trans(c, "welcomeWithName", map[string]interface{}{
			"Name": "Ahmet",
		})
		return c.SendString(message)
	})

	req, _ := app.Test(fiber.NewRequest("GET", "/test?lang=tr"))
	assert.Equal(t, 200, req.StatusCode)
}

func TestGetLocale(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		locale := GetLocale(c)
		return c.SendString(locale)
	})

	req, _ := app.Test(fiber.NewRequest("GET", "/test?lang=tr"))
	assert.Equal(t, 200, req.StatusCode)
}

func TestHasTranslation(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		exists := HasTranslation(c, "welcome")
		if exists {
			return c.SendString("exists")
		}
		return c.SendString("not exists")
	})

	req, _ := app.Test(fiber.NewRequest("GET", "/test?lang=tr"))
	assert.Equal(t, 200, req.StatusCode)
}

func TestTransWithFallback(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		message := TransWithFallback(c, "unknown.key", "Fallback Message")
		return c.SendString(message)
	})

	req, _ := app.Test(fiber.NewRequest("GET", "/test?lang=tr"))
	assert.Equal(t, 200, req.StatusCode)
}
