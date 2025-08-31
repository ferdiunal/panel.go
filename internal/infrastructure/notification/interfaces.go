package notification

import "context"

type Notifiable interface {
	RouteNotificationForMail() string
	RouteNotificationForSMS() string
	RouteNotificationFor(channel string) string
}

type Notification interface {
	Via(notifiable Notifiable) []string
	ToMail(notifiable Notifiable) *MailMessage
	ToSMS(notifiable Notifiable) *SMSMessage
	ShouldSend(notifiable Notifiable, channel string) bool
}

type Channel interface {
	Send(ctx context.Context, notifiable Notifiable, notification Notification) error
	Name() string
}

type NotificationManager struct {
	channels map[string]Channel
}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		channels: make(map[string]Channel),
	}
}

func (nm *NotificationManager) AddChannel(channel Channel) {
	nm.channels[channel.Name()] = channel
}

func (nm *NotificationManager) Send(ctx context.Context, notifiable Notifiable, notification Notification) error {
	channels := notification.Via(notifiable)
	
	for _, channelName := range channels {
		if channel, exists := nm.channels[channelName]; exists {
			if notification.ShouldSend(notifiable, channelName) {
				if err := channel.Send(ctx, notifiable, notification); err != nil {
					return err
				}
			}
		}
	}
	
	return nil
}

func (nm *NotificationManager) SendVia(ctx context.Context, channels []string, notifiable Notifiable, notification Notification) error {
	for _, channelName := range channels {
		if channel, exists := nm.channels[channelName]; exists {
			if notification.ShouldSend(notifiable, channelName) {
				if err := channel.Send(ctx, notifiable, notification); err != nil {
					return err
				}
			}
		}
	}
	
	return nil
}

func (nm *NotificationManager) GetChannel(name string) (Channel, bool) {
	channel, exists := nm.channels[name]
	return channel, exists
}

func (nm *NotificationManager) GetAvailableChannels() []string {
	channels := make([]string, 0, len(nm.channels))
	for name := range nm.channels {
		channels = append(channels, name)
	}
	return channels
}