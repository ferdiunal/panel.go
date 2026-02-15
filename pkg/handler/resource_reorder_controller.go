package handler

import (
	"fmt"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
)

type resourceReorderRequest struct {
	IDs []any `json:"ids"`
}

// HandleResourceReorder reorders records based on incoming ID sequence and updates
// the configured order column for each item.
func HandleResourceReorder(h *FieldHandler, c *context.Context) error {
	if h.Policy != nil && !h.Policy.Update(c, nil) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	reorderCfg := h.IndexReorderConfig
	reorderCfg.Column = resource.NormalizeIndexReorderColumn(reorderCfg.Column)
	if !reorderCfg.Enabled || reorderCfg.Column == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reorder is not enabled for this resource",
		})
	}

	var req resourceReorderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ids must include at least one record",
		})
	}

	txProvider, err := h.Provider.BeginTx(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	committed := false
	defer func() {
		if !committed {
			_ = txProvider.Rollback()
		}
	}()

	normalizedIDs := make([]string, 0, len(req.IDs))
	for index, rawID := range req.IDs {
		id := strings.TrimSpace(fmt.Sprint(rawID))
		if id == "" || id == "<nil>" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ids contains an invalid value",
			})
		}

		normalizedIDs = append(normalizedIDs, id)

		if _, err := txProvider.Update(c, id, map[string]interface{}{
			reorderCfg.Column: index + 1,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if err := txProvider.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	committed = true

	return c.JSON(fiber.Map{
		"success": true,
		"column":  reorderCfg.Column,
		"ids":     normalizedIDs,
	})
}
