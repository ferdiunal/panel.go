package handler

import (
	"strconv"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/notification"
	"github.com/gofiber/fiber/v2"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	Service *notification.Service
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(service *notification.Service) *NotificationHandler {
	return &NotificationHandler{
		Service: service,
	}
}

// HandleGetUnreadNotifications returns unread notifications for the current user
func (h *NotificationHandler) HandleGetUnreadNotifications(c *context.Context) error {
	// Get user from context
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Extract user ID
	var userID uint
	if u, ok := user.(interface{ GetID() uint }); ok {
		userID = u.GetID()
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user"})
	}

	// Get unread notifications
	notifications, err := h.Service.GetUnreadNotifications(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": notifications,
	})
}

// HandleMarkAsRead marks a notification as read
func (h *NotificationHandler) HandleMarkAsRead(c *context.Context) error {
	// Get notification ID from params
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	// Mark as read
	if err := h.Service.MarkAsRead(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// HandleMarkAllAsRead marks all notifications as read for the current user
func (h *NotificationHandler) HandleMarkAllAsRead(c *context.Context) error {
	// Get user from context
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Extract user ID
	var userID uint
	if u, ok := user.(interface{ GetID() uint }); ok {
		userID = u.GetID()
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user"})
	}

	// Mark all as read
	if err := h.Service.MarkAllAsRead(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications marked as read",
	})
}
