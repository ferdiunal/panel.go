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

	// Add default success notification if none exists
	if c.Resource() != nil {
		notifications := c.Resource().GetNotifications()
		if len(notifications) == 0 {
			c.Resource().NotifySuccess("Record deleted successfully")
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

	return c.JSON(fiber.Map{
		"message":       "Deleted successfully",
		"notifications": notificationsResponse,
	})
}
