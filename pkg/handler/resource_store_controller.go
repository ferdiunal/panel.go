package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceStore creates a new resource.
// It handles POST requests to /api/resource/:resource endpoint.
// Supports both multipart/form-data and application/json content types.
func HandleResourceStore(h *FieldHandler, c *context.Context) error {
	data, err := h.parseBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if h.Policy != nil && !h.Policy.Create(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	result, err := h.Provider.Create(c, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": h.resolveResourceFields(c.Ctx, c.Resource(), result, h.Elements)})
}
