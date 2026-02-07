package notification

import (
	"github.com/ferdiunal/panel.go/pkg/core"
	notificationDomain "github.com/ferdiunal/panel.go/pkg/domain/notification"
	"gorm.io/gorm"
)

// Service handles notification operations
type Service struct {
	db *gorm.DB
}

// NewService creates a new notification service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// SaveNotifications saves notifications from context to database
func (s *Service) SaveNotifications(ctx *core.ResourceContext) error {
	notifications := ctx.GetNotifications()
	if len(notifications) == 0 {
		return nil
	}

	// Get user ID from context if available
	var userID *uint
	if ctx.User != nil {
		if user, ok := ctx.User.(interface{ GetID() uint }); ok {
			id := user.GetID()
			userID = &id
		}
	}

	// Convert context notifications to domain notifications
	for _, notif := range notifications {
		dbNotif := &notificationDomain.Notification{
			UserID:   userID,
			Message:  notif.Message,
			Type:     notificationDomain.NotificationType(notif.Type),
			Duration: notif.Duration,
			Read:     false,
		}

		if err := s.db.Create(dbNotif).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetUnreadNotifications retrieves unread notifications for a user
func (s *Service) GetUnreadNotifications(userID uint) ([]notificationDomain.Notification, error) {
	var notifications []notificationDomain.Notification
	err := s.db.Where("user_id = ? AND read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// MarkAsRead marks a notification as read
func (s *Service) MarkAsRead(notificationID uint) error {
	var notif notificationDomain.Notification
	if err := s.db.First(&notif, notificationID).Error; err != nil {
		return err
	}
	return notif.MarkAsRead(s.db)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *Service) MarkAllAsRead(userID uint) error {
	return s.db.Model(&notificationDomain.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}
