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

	// Add default success notification if none exists
	if c.Resource() != nil {
		notifications := c.Resource().GetNotifications()
		if len(notifications) == 0 {
			c.Resource().NotifySuccess("Record created successfully")
		}
	}

	// Save notifications to database
	if c.Resource() != nil && h.NotificationService != nil {
		if err := h.NotificationService.SaveNotifications(c.Resource()); err != nil {
			// Log error but don't fail the request
			// fmt.Printf("Failed to save notifications: %v\n", err)
		}
	}

	// Get notifications for response
	var notificationsResponse []map[string]interface{}
	if c.Resource() != nil {
		for _, notif := range c.Resource().GetNotifications() {
			notificationsResponse = append(notificationsResponse, map[string]interface{}{
				"message":  notif.Message,
				"type":     notif.Type,
				"duration": notif.Duration,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":          h.resolveResourceFields(c.Ctx, c.Resource(), result, h.Elements),
		"notifications": notificationsResponse,
	})
}
