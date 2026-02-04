package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceDetail handles retrieving detailed view of a single resource.
// It filters fields based on detail context and resolves their values.
func HandleResourceDetail(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.View(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Filter for Detail
	var detailElements []fields.Element
	for _, element := range h.Elements {
		if !element.IsVisible(c.Resource()) {
			continue
		}

		ctxStr := element.GetContext()

		// Skip if explicitly hidden on detail or restricted to other contexts
		if ctxStr == fields.HIDE_ON_DETAIL ||
			ctxStr == fields.ONLY_ON_LIST ||
			ctxStr == fields.ONLY_ON_FORM ||
			ctxStr == fields.HIDE_ON_UPDATE {
			continue
		}
		detailElements = append(detailElements, element)
	}

	// Resolve fields with values
	resolvedMap := h.resolveResourceFields(c.Ctx, c.Resource(), item, detailElements)

	// Convert map to ordered slice based on h.Elements order (preserving filtered list order)
	var orderedFields []map[string]interface{}
	// Iterate over detailElements to preserve order
	for _, element := range detailElements {
		if val, ok := resolvedMap[element.GetKey()]; ok {
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
