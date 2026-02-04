package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceShow retrieves a single resource by ID.
// It handles GET requests to /api/resource/:resource/:id endpoint.
func HandleResourceShow(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.View(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Determine elements to use
	// For Show, we generally use the handler's elements, but context overrides could happen if middleware sets them
	// For now, simpler to use h.Elements as middleware might be Index specific or general.
	// But let's check context just in case consistent with Index.
	ctx := context.FromFiber(c.Ctx)
	var elements []fields.Element
	if ctx != nil && len(ctx.Elements) > 0 {
		elements = ctx.Elements
	} else {
		elements = h.Elements
	}

	return c.JSON(fiber.Map{
		"data": h.resolveResourceFields(c.Ctx, c.Resource(), item, elements),
		"meta": fiber.Map{
			"title": h.Resource.Title(),
			"policy": fiber.Map{
				"view":   h.Policy == nil || h.Policy.View(c, item),
				"update": h.Policy == nil || h.Policy.Update(c, item),
				"delete": h.Policy == nil || h.Policy.Delete(c, item),
			},
		},
	})
}
