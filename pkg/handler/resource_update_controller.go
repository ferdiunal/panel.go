package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceUpdate updates an existing resource.
// It handles PUT requests to /api/resource/:resource/:id endpoint.
// Supports both multipart/form-data and application/json content types.
func HandleResourceUpdate(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")
	data, err := h.parseBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Fetch existing to check policy
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Update(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	result, err := h.Provider.Update(c, id, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": h.resolveResourceFields(c.Ctx, c.Resource(), result, h.Elements)})
}
