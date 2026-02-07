package notification

// Type represents the notification type
type Type string

const (
	TypeSuccess Type = "success"
	TypeError   Type = "error"
	TypeWarning Type = "warning"
	TypeInfo    Type = "info"
)

// Notification represents a user notification
type Notification struct {
	Message  string `json:"message"`
	Type     Type   `json:"type"`
	Duration int    `json:"duration"` // Duration in milliseconds
}

// New creates a new notification
func New(message string, notifType Type) *Notification {
	return &Notification{
		Message:  message,
		Type:     notifType,
		Duration: 3000, // Default 3 seconds
	}
}

// Success creates a success notification
func Success(message string) *Notification {
	return New(message, TypeSuccess)
}

// Error creates an error notification
func Error(message string) *Notification {
	return New(message, TypeError)
}

// Warning creates a warning notification
func Warning(message string) *Notification {
	return New(message, TypeWarning)
}

// Info creates an info notification
func Info(message string) *Notification {
	return New(message, TypeInfo)
}

// SetDuration sets the notification duration
func (n *Notification) SetDuration(duration int) *Notification {
	n.Duration = duration
	return n
}
