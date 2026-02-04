package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleLensIndex lists all available lenses for a resource.
// This corresponds to Laravel Nova's LensController@index method.
func HandleLensIndex(h *FieldHandler, c *context.Context) error {
	if h.Resource == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}

	lenses := h.Resource.Lenses()
	response := make([]map[string]interface{}, 0)

	for _, lens := range lenses {
		response = append(response, map[string]interface{}{
			"name": lens.Name(),
			"slug": lens.Slug(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

// HandleLens handles lens-based resource listing with filtered data.
// This corresponds to Laravel Nova's LensController@show method.
// It uses the lens handler which has already been configured with the lens query.
func HandleLens(h *FieldHandler, c *context.Context) error {
	// Lens handler is already configured with filtered query via NewLensHandler
	// We can directly use the Index logic but with the lens's filtered dataset
	return HandleResourceIndex(h, c)
}
