package handler

import (
	panelcontext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
)

// ensureResourceContext guarantees that request-scoped ResourceContext exists.
// Visibility rules are evaluated against this context in all field IsVisible checks.
func ensureResourceContext(
	c *panelcontext.Context,
	resource any,
	lens any,
	visibilityCtx core.VisibilityContext,
) *core.ResourceContext {
	if c == nil || c.Ctx == nil {
		return nil
	}

	resourceCtx := c.Resource()
	if resourceCtx == nil {
		resourceCtx = core.NewResourceContextWithVisibility(
			c.Ctx,
			resource,
			lens,
			visibilityCtx,
			nil,
			c.User(),
			nil,
		)
		c.Locals(core.ResourceContextKey, resourceCtx)
		return resourceCtx
	}

	if visibilityCtx != "" {
		resourceCtx.VisibilityCtx = visibilityCtx
	}
	if resourceCtx.Resource == nil && resource != nil {
		resourceCtx.Resource = resource
	}
	if resourceCtx.Lens == nil && lens != nil {
		resourceCtx.Lens = lens
	}
	if resourceCtx.Request == nil {
		resourceCtx.Request = c.Ctx
	}
	if resourceCtx.User == nil {
		if user := c.User(); user != nil {
			resourceCtx.User = user
		}
	}

	return resourceCtx
}
