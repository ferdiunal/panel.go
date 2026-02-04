package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceEdit returns form fields for resource editing with current values.
// It handles GET requests to /api/resource/:resource/:id/edit endpoint.
func HandleResourceEdit(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Update(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Filter for Update
	var updateElements []fields.Element
	for _, element := range h.Elements {
		if !element.IsVisible(c.Resource()) {
			continue
		}

		ctxStr := element.GetContext()
		if ctxStr != fields.HIDE_ON_UPDATE &&
			ctxStr != fields.ONLY_ON_LIST &&
			ctxStr != fields.ONLY_ON_DETAIL &&
			ctxStr != fields.ONLY_ON_CREATE {
			updateElements = append(updateElements, element)
		}
	}

	// Resolve fields with values
	resolvedMap := h.resolveResourceFields(c.Ctx, c.Resource(), item, updateElements)

	// Convert map to ordered slice based on h.Elements order
	var orderedFields []map[string]interface{}
	for _, element := range updateElements {
		if val, ok := resolvedMap[element.GetKey()]; ok {
			// Cast val to map
			if fieldMap, ok := val.(map[string]interface{}); ok {
				orderedFields = append(orderedFields, fieldMap)
			}
		}
	}

	return c.JSON(fiber.Map{
		"fields": orderedFields,
		"meta": fiber.Map{
			"title": h.Title,
		},
	})
}
