package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuditEvent represents a security audit event
type AuditEvent struct {
	Timestamp    time.Time              `json:"timestamp"`
	EventType    string                 `json:"event_type"`
	UserID       string                 `json:"user_id,omitempty"`
	Email        string                 `json:"email,omitempty"`
	IP           string                 `json:"ip"`
	UserAgent    string                 `json:"user_agent"`
	Method       string                 `json:"method"`
	Path         string                 `json:"path"`
	StatusCode   int                    `json:"status_code"`
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Resource     string                 `json:"resource,omitempty"`
	Action       string                 `json:"action,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AuditLogger interface for audit logging implementations
type AuditLogger interface {
	Log(event AuditEvent) error
}

// ConsoleAuditLogger logs to console (for development)
type ConsoleAuditLogger struct{}

func (l *ConsoleAuditLogger) Log(event AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	fmt.Printf("[AUDIT] %s\n", string(data))
	return nil
}

// FileAuditLogger logs to file (for production)
type FileAuditLogger struct {
	// TODO: Implement file-based logging
}

func (l *FileAuditLogger) Log(event AuditEvent) error {
	// TODO: Implement file writing with rotation
	return nil
}

// AuditMiddleware creates an audit logging middleware
func AuditMiddleware(logger AuditLogger) fiber.Handler {
	if logger == nil {
		logger = &ConsoleAuditLogger{}
	}

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request completes
		event := AuditEvent{
			Timestamp:  start,
			Method:     c.Method(),
			Path:       c.Path(),
			IP:         c.IP(),
			UserAgent:  c.Get("User-Agent"),
			StatusCode: c.Response().StatusCode(),
			Success:    c.Response().StatusCode() < 400,
		}

		// Extract user info from context if available
		if user := c.Locals("user"); user != nil {
			// Type assertion to extract user details
			// This depends on your user struct
			event.Metadata = map[string]interface{}{
				"user": user,
			}
		}

		// Determine event type based on path and method
		event.EventType = determineEventType(c.Method(), c.Path(), c.Response().StatusCode())

		// Log the event
		if logErr := logger.Log(event); logErr != nil {
			fmt.Printf("[AUDIT ERROR] Failed to log event: %v\n", logErr)
		}

		return err
	}
}

// LogAuthEvent logs authentication-related events
func LogAuthEvent(logger AuditLogger, eventType, email, ip, userAgent string, success bool, errorMsg string) {
	event := AuditEvent{
		Timestamp:    time.Now(),
		EventType:    eventType,
		Email:        email,
		IP:           ip,
		UserAgent:    userAgent,
		Success:      success,
		ErrorMessage: errorMsg,
	}

	if err := logger.Log(event); err != nil {
		fmt.Printf("[AUDIT ERROR] Failed to log auth event: %v\n", err)
	}
}

// LogPermissionCheck logs permission check events
func LogPermissionCheck(logger AuditLogger, userID, resource, action string, granted bool) {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "permission_check",
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Success:   granted,
	}

	if err := logger.Log(event); err != nil {
		fmt.Printf("[AUDIT ERROR] Failed to log permission check: %v\n", err)
	}
}

// LogDataAccess logs data access events
func LogDataAccess(logger AuditLogger, userID, resource, action string, recordID string) {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "data_access",
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Metadata: map[string]interface{}{
			"record_id": recordID,
		},
	}

	if err := logger.Log(event); err != nil {
		fmt.Printf("[AUDIT ERROR] Failed to log data access: %v\n", err)
	}
}

// determineEventType determines the event type based on request details
func determineEventType(method, path string, statusCode int) string {
	// Authentication endpoints
	if contains(path, "/auth/sign-in") {
		if statusCode < 400 {
			return "login_success"
		}
		return "login_failure"
	}
	if contains(path, "/auth/sign-up") {
		return "registration"
	}
	if contains(path, "/auth/sign-out") {
		return "logout"
	}
	if contains(path, "/auth/forgot-password") {
		return "password_reset_request"
	}

	// Resource operations
	if contains(path, "/resource/") {
		switch method {
		case "GET":
			return "resource_read"
		case "POST":
			return "resource_create"
		case "PUT", "PATCH":
			return "resource_update"
		case "DELETE":
			return "resource_delete"
		}
	}

	// Settings operations
	if contains(path, "/pages/settings") {
		return "settings_change"
	}

	return "api_request"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
	       len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
	       len(s) > len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
