package notification

type BaseNotification struct{}

func (b *BaseNotification) Via(notifiable Notifiable) []string {
	return []string{"mail"}
}

func (b *BaseNotification) ToMail(notifiable Notifiable) *MailMessage {
	return nil
}

func (b *BaseNotification) ToSMS(notifiable Notifiable) *SMSMessage {
	return nil
}

func (b *BaseNotification) ShouldSend(notifiable Notifiable, channel string) bool {
	return true
}

type SimpleNotifiable struct {
	Email string
	Phone string
}

func NewSimpleNotifiable(email, phone string) *SimpleNotifiable {
	return &SimpleNotifiable{
		Email: email,
		Phone: phone,
	}
}

func (s *SimpleNotifiable) RouteNotificationForMail() string {
	return s.Email
}

func (s *SimpleNotifiable) RouteNotificationForSMS() string {
	return s.Phone
}

func (s *SimpleNotifiable) RouteNotificationFor(channel string) string {
	switch channel {
	case "mail":
		return s.RouteNotificationForMail()
	case "sms":
		return s.RouteNotificationForSMS()
	default:
		return ""
	}
}