package auth

import (
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func TestSessionMiddleware_AllowsAPIKeyAuthenticatedRequest(t *testing.T) {
	h := &Handler{environment: "test"}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.APIKeyAuthenticatedLocalKey, true)
		return c.Next()
	})
	app.Get("/", appContext.Wrap(h.SessionMiddleware), appContext.Wrap(func(c *appContext.Context) error {
		if c.User() == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("user not set")
		}
		if c.Session() == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("session not set")
		}
		return c.SendStatus(fiber.StatusOK)
	}))

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestSessionMiddleware_RejectsWithoutSessionOrAPIKey(t *testing.T) {
	h := &Handler{environment: "test"}

	app := fiber.New()
	app.Get("/", appContext.Wrap(h.SessionMiddleware))

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}
