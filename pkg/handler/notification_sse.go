package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	notificationDomain "github.com/ferdiunal/panel.go/pkg/domain/notification"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// NotificationSSEHandler handles SSE connections for notifications
type NotificationSSEHandler struct {
	db *gorm.DB
}

// NewNotificationSSEHandler creates a new SSE handler
func NewNotificationSSEHandler(db *gorm.DB) *NotificationSSEHandler {
	return &NotificationSSEHandler{db: db}
}

// HandleNotificationStream streams notifications via SSE
func (h *NotificationSSEHandler) HandleNotificationStream(c *context.Context) error {
	// Get user ID from context
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no") // Disable nginx buffering

	// Send initial notifications
	var notifications []notificationDomain.Notification
	h.db.Where("user_id = ? AND read = ?", userID, false).
		Order("created_at DESC").
		Limit(50).
		Find(&notifications)

	if len(notifications) > 0 {
		data, _ := json.Marshal(notifications)
		fmt.Fprintf(c, "data: %s\n\n", data)
		c.Flush()
	}

	// Poll for new notifications every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	lastCheck := time.Now()

	for {
		select {
		case <-ticker.C:
			// Check for new notifications since last check
			var newNotifications []notificationDomain.Notification
			h.db.Where("user_id = ? AND read = ? AND created_at > ?",
				userID, false, lastCheck).
				Order("created_at DESC").
				Find(&newNotifications)

			if len(newNotifications) > 0 {
				data, _ := json.Marshal(newNotifications)
				fmt.Fprintf(c, "data: %s\n\n", data)
				c.Flush()
				lastCheck = time.Now()
			}

		case <-c.Context().Done():
			// Client disconnected
			return nil
		}
	}
}
