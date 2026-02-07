package domain

import (
	"time"

	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeSuccess NotificationType = "success"
	NotificationTypeError   NotificationType = "error"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeInfo    NotificationType = "info"
)

// Notification represents a user notification stored in the database
type Notification struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    *uint            `gorm:"index" json:"user_id"`
	Message   string           `gorm:"type:text;not null" json:"message"`
	Type      NotificationType `gorm:"type:varchar(20);not null;default:'info'" json:"type"`
	Duration  int              `gorm:"default:3000" json:"duration"` // Duration in milliseconds
	Read      bool             `gorm:"default:false" json:"read"`
	ReadAt    *time.Time       `json:"read_at"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName specifies the table name for the Notification model
func (Notification) TableName() string {
	return "notifications"
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead(db *gorm.DB) error {
	now := time.Now()
	n.Read = true
	n.ReadAt = &now
	return db.Save(n).Error
}
