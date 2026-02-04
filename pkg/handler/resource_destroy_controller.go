package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceDestroy deletes a resource.
// It handles DELETE requests to /api/resource/:resource/:id endpoint.
func HandleResourceDestroy(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")

	// Fetch for Policy Check
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Delete(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.Provider.Delete(c, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Deleted successfully"})
}
