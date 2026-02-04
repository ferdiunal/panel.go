package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleFieldList handles listing fields for a resource context.
// It retrieves fields from the context, applies visibility filtering,
// and resolves values if a resource instance is present.
func HandleFieldList(h *FieldHandler, c *context.Context) error {
	ctx := c.Resource()
	if ctx == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Field context not found",
		})
	}

	response := make([]map[string]interface{}, 0)

	for _, element := range ctx.Elements {
		// Resolve value if resource is present
		if ctx.Resource != nil {
			element.Extract(ctx.Resource)
		}
		response = append(response, element.JsonSerialize())
	}

	return c.JSON(response)
}
