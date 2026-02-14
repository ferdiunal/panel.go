package i18n

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func setupTestApp(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New()
	localesDir := t.TempDir()

	trFile := filepath.Join(localesDir, "tr.yaml")
	enFile := filepath.Join(localesDir, "en.yaml")

	// Keep bundles minimal and deterministic for tests.
	trContent := "welcome: \"Hoş geldiniz\"\nwelcomeWithName: \"Hoş geldiniz, {{.Name}}\"\n"
	enContent := "welcome: \"Welcome\"\nwelcomeWithName: \"Welcome, {{.Name}}\"\n"

	if err := os.WriteFile(trFile, []byte(trContent), 0644); err != nil {
		t.Fatalf("write tr locale: %v", err)
	}
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatalf("write en locale: %v", err)
	}

	// i18n middleware'ini ekle
	app.Use(fiberi18n.New(&fiberi18n.Config{
		RootPath:         localesDir,
		AcceptLanguages:  []language.Tag{language.Turkish, language.English},
		DefaultLanguage:  language.Turkish,
		FormatBundleFile: "yaml",
	}))

	return app
}

func TestTrans(t *testing.T) {
	app := setupTestApp(t)

	app.Get("/test", func(c *fiber.Ctx) error {
		message := Trans(c, "welcome")
		return c.SendString(message)
	})

	req, _ := app.Test(httptest.NewRequest("GET", "/test?lang=tr", nil))
	assert.Equal(t, 200, req.StatusCode)
}

func TestTransWithTemplate(t *testing.T) {
	app := setupTestApp(t)

	app.Get("/test", func(c *fiber.Ctx) error {
		message := Trans(c, "welcomeWithName", map[string]interface{}{
			"Name": "Ahmet",
		})
		return c.SendString(message)
	})

	req, _ := app.Test(httptest.NewRequest("GET", "/test?lang=tr", nil))
	assert.Equal(t, 200, req.StatusCode)
}

func TestGetLocale(t *testing.T) {
	app := setupTestApp(t)

	app.Get("/test", func(c *fiber.Ctx) error {
		locale := GetLocale(c)
		return c.SendString(locale)
	})

	req, _ := app.Test(httptest.NewRequest("GET", "/test?lang=tr", nil))
	assert.Equal(t, 200, req.StatusCode)
}

func TestHasTranslation(t *testing.T) {
	app := setupTestApp(t)

	app.Get("/test", func(c *fiber.Ctx) error {
		exists := HasTranslation(c, "welcome")
		if exists {
			return c.SendString("exists")
		}
		return c.SendString("not exists")
	})

	req, _ := app.Test(httptest.NewRequest("GET", "/test?lang=tr", nil))
	assert.Equal(t, 200, req.StatusCode)
}

func TestTransWithFallback(t *testing.T) {
	app := setupTestApp(t)

	app.Get("/test", func(c *fiber.Ctx) error {
		message := TransWithFallback(c, "unknown.key", "Fallback Message")
		return c.SendString(message)
	})

	req, _ := app.Test(httptest.NewRequest("GET", "/test?lang=tr", nil))
	assert.Equal(t, 200, req.StatusCode)
}
