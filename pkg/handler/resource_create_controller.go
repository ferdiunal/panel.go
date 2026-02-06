package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceCreate returns form fields for resource creation.
// It handles GET requests to /api/resource/:resource/create endpoint.
func HandleResourceCreate(h *FieldHandler, c *context.Context) error {
	if h.Policy != nil && !h.Policy.Create(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	createFields := make([]map[string]interface{}, 0)
	for _, element := range h.Elements {
		if !element.IsVisible(c.Resource()) {
			continue
		}

		serialized := element.JsonSerialize()

		// Resolve options (AutoOptions and callbacks)
		h.ResolveFieldOptions(element, serialized, nil)

		ctxStr := element.GetContext()

		if ctxStr != fields.HIDE_ON_CREATE &&
			ctxStr != fields.ONLY_ON_LIST &&
			ctxStr != fields.ONLY_ON_DETAIL &&
			ctxStr != fields.ONLY_ON_UPDATE {
			createFields = append(createFields, serialized)
		}
	}

	return c.JSON(fiber.Map{
		"fields": createFields,
	})
}
