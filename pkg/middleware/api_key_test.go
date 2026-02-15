package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAPIKeyAuthMiddleware_ValidKey(t *testing.T) {
	auth := NewAPIKeyAuth(true, "", []string{"secret-key"})

	app := fiber.New()
	app.Use(auth.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		if authed, ok := c.Locals(APIKeyAuthenticatedLocalKey).(bool); !ok || !authed {
			return c.Status(fiber.StatusUnauthorized).SendString("missing api key auth marker")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "secret-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAPIKeyAuthMiddleware_InvalidKey(t *testing.T) {
	auth := NewAPIKeyAuth(true, "", []string{"secret-key"})

	app := fiber.New()
	app.Use(auth.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "wrong-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestAPIKeyAuthMiddleware_MissingKeyFallsBack(t *testing.T) {
	auth := NewAPIKeyAuth(true, "", []string{"secret-key"})

	app := fiber.New()
	app.Use(auth.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		if _, ok := c.Locals(APIKeyAuthenticatedLocalKey).(bool); ok {
			t.Fatalf("api key marker should not be set when header is missing")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAPIKeyAuthMiddleware_Disabled(t *testing.T) {
	auth := NewAPIKeyAuth(false, "", []string{"secret-key"})

	app := fiber.New()
	app.Use(auth.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "wrong-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAPIKeyAuthMiddleware_SetConfigAtRuntime(t *testing.T) {
	auth := NewAPIKeyAuth(true, "", []string{"secret-key"})

	app := fiber.New()
	app.Use(auth.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		if authed, ok := c.Locals(APIKeyAuthenticatedLocalKey).(bool); !ok || !authed {
			return c.Status(fiber.StatusUnauthorized).SendString("missing api key auth marker")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	auth.SetConfig(true, "X-App-Key", []string{"new-secret"})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("X-App-Key", "new-secret")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// Old header/key pair should not work after runtime update.
	oldReq := httptest.NewRequest(fiber.MethodGet, "/", nil)
	oldReq.Header.Set("X-API-Key", "secret-key")

	oldResp, err := app.Test(oldReq)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	if oldResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401 when old header no longer authenticates, got %d", oldResp.StatusCode)
	}
}

func TestAPIKeyAuthMiddleware_DynamicValidator(t *testing.T) {
	auth := NewAPIKeyAuth(true, "", nil)
	auth.SetDynamicValidator(func(c *fiber.Ctx, incoming string) bool {
		return incoming == "managed-key"
	})

	app := fiber.New()
	app.Use(auth.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		if authed, ok := c.Locals(APIKeyAuthenticatedLocalKey).(bool); !ok || !authed {
			return c.Status(fiber.StatusUnauthorized).SendString("missing api key auth marker")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "managed-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	invalidReq := httptest.NewRequest(fiber.MethodGet, "/", nil)
	invalidReq.Header.Set("X-API-Key", "not-valid")

	invalidResp, err := app.Test(invalidReq)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if invalidResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", invalidResp.StatusCode)
	}
}

func TestAPIKeyAuthMiddleware_AtomicSnapshotMode(t *testing.T) {
	auth := NewAPIKeyAuth(true, "", []string{"secret-key"})
	auth.SetAtomicSnapshotEnabled(true)

	app := fiber.New()
	app.Use(auth.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		if authed, ok := c.Locals(APIKeyAuthenticatedLocalKey).(bool); !ok || !authed {
			return c.Status(fiber.StatusUnauthorized).SendString("missing api key auth marker")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "secret-key")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	auth.SetConfig(true, "X-App-Key", []string{"next-key"})

	oldReq := httptest.NewRequest(fiber.MethodGet, "/", nil)
	oldReq.Header.Set("X-API-Key", "secret-key")
	oldResp, err := app.Test(oldReq)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if oldResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", oldResp.StatusCode)
	}

	newReq := httptest.NewRequest(fiber.MethodGet, "/", nil)
	newReq.Header.Set("X-App-Key", "next-key")
	newResp, err := app.Test(newReq)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if newResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", newResp.StatusCode)
	}
}
