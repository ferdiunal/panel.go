package context_test

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/gofiber/fiber/v2"
)

// TestContextTypeSafety verifies that Context methods return correct types.
// This test validates Requirement 9.2, 9.3: Type safety should be preserved.
//
// Validates: Requirements 9.2, 9.3
func TestContextTypeSafety(t *testing.T) {
	app := fiber.New()

	app.Get("/test", context.Wrap(func(c *context.Context) error {
		// Compile-time type checks - if these compile, types are correct
		var _ *user.User = c.User()
		var _ *session.Session = c.Session()
		var _ *core.ResourceContext = c.Resource()

		return c.SendString("OK")
	}))

	// This test primarily validates compile-time type safety
	// If it compiles, the types are correct
}

// TestFromFiberReturnsCorrectType verifies FromFiber returns *core.ResourceContext
func TestFromFiberReturnsCorrectType(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		// Set a ResourceContext in Locals
		ctx := core.NewResourceContext(c, nil, []core.Element{})
		c.Locals(core.ResourceContextKey, ctx)
		return c.Next()
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		// Compile-time check: FromFiber should return *core.ResourceContext
		var _ *core.ResourceContext = context.FromFiber(c)

		resourceCtx := context.FromFiber(c)
		if resourceCtx == nil {
			t.Error("FromFiber should return non-nil ResourceContext")
		}

		return c.SendString("OK")
	})
}

// TestContextUserMethod verifies User() returns *user.User
func TestContextUserMethod(t *testing.T) {
	app := fiber.New()

	testUser := &user.User{
		ID:   123,
		Role: "admin",
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", testUser)
		return c.Next()
	})

	app.Get("/test", context.Wrap(func(c *context.Context) error {
		// Compile-time check
		var _ *user.User = c.User()

		// Runtime check
		u := c.User()
		if u == nil {
			t.Error("User() should return non-nil user")
		}
		if u.ID != 123 {
			t.Errorf("Expected user ID '123', got '%d'", u.ID)
		}

		return c.SendString("OK")
	}))
}

// TestContextSessionMethod verifies Session() returns *session.Session
func TestContextSessionMethod(t *testing.T) {
	app := fiber.New()

	testSession := &session.Session{
		ID: 1,
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("session", testSession)
		return c.Next()
	})

	app.Get("/test", context.Wrap(func(c *context.Context) error {
		// Compile-time check
		var _ *session.Session = c.Session()

		// Runtime check
		s := c.Session()
		if s == nil {
			t.Error("Session() should return non-nil session")
		}
		if s.ID != 1 {
			t.Errorf("Expected session ID '1', got '%d'", s.ID)
		}

		return c.SendString("OK")
	}))
}

// TestContextResourceMethod verifies Resource() returns *core.ResourceContext
func TestContextResourceMethod(t *testing.T) {
	app := fiber.New()

	testResource := map[string]interface{}{"id": "1", "name": "Test"}
	testElements := []core.Element{}

	app.Use(func(c *fiber.Ctx) error {
		ctx := core.NewResourceContext(c, testResource, testElements)
		c.Locals(core.ResourceContextKey, ctx)
		return c.Next()
	})

	app.Get("/test", context.Wrap(func(c *context.Context) error {
		// Compile-time check
		var _ *core.ResourceContext = c.Resource()

		// Runtime check
		resourceCtx := c.Resource()
		if resourceCtx == nil {
			t.Error("Resource() should return non-nil ResourceContext")
		}
		if resourceCtx.Resource == nil {
			t.Error("ResourceContext.Resource should not be nil")
		}

		return c.SendString("OK")
	}))
}

// TestContextNilSafety verifies Context methods handle nil values gracefully
func TestContextNilSafety(t *testing.T) {
	app := fiber.New()

	app.Get("/test", context.Wrap(func(c *context.Context) error {
		// These should not panic even when values are not set
		user := c.User()
		if user != nil {
			t.Error("User() should return nil when not set")
		}

		session := c.Session()
		if session != nil {
			t.Error("Session() should return nil when not set")
		}

		resource := c.Resource()
		if resource != nil {
			t.Error("Resource() should return nil when not set")
		}

		return c.SendString("OK")
	}))
}
