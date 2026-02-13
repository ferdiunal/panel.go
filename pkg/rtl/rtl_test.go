package rtl

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestIsRTL(t *testing.T) {
	tests := []struct {
		name     string
		lang     language.Tag
		expected bool
	}{
		{"Arabic is RTL", language.Arabic, true},
		{"Hebrew is RTL", language.Hebrew, true},
		{"Persian is RTL", language.Persian, true},
		{"Urdu is RTL", language.Urdu, true},
		{"English is LTR", language.English, false},
		{"Turkish is LTR", language.Turkish, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRTL(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsRTLString(t *testing.T) {
	tests := []struct {
		name     string
		langCode string
		expected bool
	}{
		{"ar is RTL", "ar", true},
		{"he is RTL", "he", true},
		{"fa is RTL", "fa", true},
		{"ur is RTL", "ur", true},
		{"en is LTR", "en", false},
		{"tr is LTR", "tr", false},
		{"invalid code is LTR", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRTLString(tt.langCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDirection(t *testing.T) {
	tests := []struct {
		name     string
		lang     language.Tag
		expected string
	}{
		{"Arabic direction", language.Arabic, "rtl"},
		{"Hebrew direction", language.Hebrew, "rtl"},
		{"English direction", language.English, "ltr"},
		{"Turkish direction", language.Turkish, "ltr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDirection(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDirectionString(t *testing.T) {
	tests := []struct {
		name     string
		langCode string
		expected string
	}{
		{"ar direction", "ar", "rtl"},
		{"he direction", "he", "rtl"},
		{"en direction", "en", "ltr"},
		{"tr direction", "tr", "ltr"},
		{"invalid direction", "invalid", "ltr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDirectionString(tt.langCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDirectionFromContext(t *testing.T) {
	app := fiber.New()

	// Test with query parameter
	app.Get("/test-query", func(c *fiber.Ctx) error {
		dir := GetDirectionFromContext(c)
		return c.SendString(dir)
	})

	// Test with Accept-Language header
	app.Get("/test-header", func(c *fiber.Ctx) error {
		dir := GetDirectionFromContext(c)
		return c.SendString(dir)
	})

	// Test query parameter with RTL language
	req := fiber.AcquireRequest()
	req.SetRequestURI("/test-query?lang=ar")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test query parameter with LTR language
	req = fiber.AcquireRequest()
	req.SetRequestURI("/test-query?lang=en")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(Middleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Test with RTL language
	req := fiber.AcquireRequest()
	req.SetRequestURI("/test?lang=ar")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "rtl", string(resp.Header.Peek("X-Text-Direction")))

	// Test with LTR language
	req = fiber.AcquireRequest()
	req.SetRequestURI("/test?lang=en")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "ltr", string(resp.Header.Peek("X-Text-Direction")))
}
