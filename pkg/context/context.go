package context

import (
	stdcontext "context"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/user"

	"github.com/gofiber/fiber/v2"
)

// FromFiber retrieves the ResourceContext from fiber.Locals
func FromFiber(c *fiber.Ctx) *core.ResourceContext {
	val := c.Locals(core.ResourceContextKey)
	if val == nil {
		return nil
	}
	return val.(*core.ResourceContext)
}

// Context wraps fiber.Ctx to provide type-safe access to Locals
type Context struct {
	*fiber.Ctx
}

// Handler is a custom handler type that uses our typed Context
type Handler func(*Context) error

// User retrieves the authenticated user from Locals
func (c *Context) User() *user.User {
	if u, ok := c.Locals("user").(*user.User); ok {
		return u
	}
	return nil
}

// Session retrieves the active session from Locals
func (c *Context) Session() *session.Session {
	if s, ok := c.Locals("session").(*session.Session); ok {
		return s
	}
	return nil
}

// Resource retrieves the ResourceContext from Locals
func (c *Context) Resource() *core.ResourceContext {
	return FromFiber(c.Ctx)
}

// Wrap converts our custom Handler to a standard fiber.Handler
func Wrap(h Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return h(&Context{Ctx: c})
	}
}

// HasRole checks if the authenticated user has the given role
func (c *Context) HasRole(role string) bool {
	u := c.User()
	if u == nil {
		return false
	}
	return u.Role == role || u.Role == "admin"
}

// HasPermission checks if the authenticated user has permission for the action
func (c *Context) HasPermission(action string) bool {
	u := c.User()
	if u == nil {
		return false
	}
	if u.Role == "admin" {
		return true
	}
	// TODO: Integrate actual permission logic
	return true
}

// Context returns the underlying context.Context from fiber.Ctx
func (c *Context) Context() stdcontext.Context {
	return c.Ctx.Context()
}
