package handler

import (
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleCardList handles listing all cards for a resource.
// It resolves each card's data and returns them with their metadata.
func HandleCardList(h *FieldHandler, c *context.Context) error {
	resp := make([]map[string]interface{}, 0)
	for i, w := range h.Cards {
		// Serialize base properties (component, width, etc.)
		serialized := w.JsonSerialize()
		serialized["index"] = i
		serialized["name"] = w.Name()
		serialized["component"] = w.Component()
		serialized["width"] = w.Width()

		// Resolve data
		data, err := w.Resolve(c, h.DB)
		if err != nil {
			fmt.Printf("Error resolving card %s: %v\n", w.Name(), err)
			serialized["error"] = err.Error()
		} else {
			// Assign resolved data to "data" key
			serialized["data"] = data
		}

		resp = append(resp, serialized)
	}

	return c.JSON(fiber.Map{
		"data": resp,
	})
}
