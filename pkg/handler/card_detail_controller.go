package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleCardDetail handles retrieving a single card by index.
// It validates the index and resolves the card's data.
func HandleCardDetail(h *FieldHandler, c *context.Context) error {
	index, err := c.ParamsInt("index")
	if err != nil || index < 0 || index >= len(h.Cards) {
		return c.Status(404).JSON(fiber.Map{"error": "Card not found"})
	}

	w := h.Cards[index]

	// Resolve data
	data, err := w.Resolve(c, h.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
