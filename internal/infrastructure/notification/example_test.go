package notification_test

import (
	"context"
	"fmt"

	"panel.go/internal/infrastructure/email"
	"panel.go/internal/infrastructure/notification"
	"panel.go/internal/infrastructure/notification/channels"
	"panel.go/internal/infrastructure/sms"
)

type WelcomeNotification struct {
	notification.BaseNotification
	UserName string
}

func (w *WelcomeNotification) Via(notifiable notification.Notifiable) []string {
	return []string{"mail", "sms"}
}

func (w *WelcomeNotification) ToMail(notifiable notification.Notifiable) *notification.MailMessage {
	return notification.NewMailMessage().
		SetSubject("Welcome to our platform!").
		SetBody(fmt.Sprintf("Hello %s, welcome to our platform!", w.UserName)).
		SetHTMLBody(fmt.Sprintf("<h1>Hello %s</h1><p>Welcome to our platform!</p>", w.UserName))
}

func (w *WelcomeNotification) ToSMS(notifiable notification.Notifiable) *notification.SMSMessage {
	return notification.NewSMSMessage().
		SetBody(fmt.Sprintf("Hello %s! Welcome to our platform.", w.UserName))
}

func ExampleNotification() {
	emailProvider, _ := email.NewEmailProviderFromEnv()
	emailService := email.NewEmailService(emailProvider)
	
	smsProvider, _ := sms.NewSMSProviderFromEnv()
	smsService := sms.NewSMSService(smsProvider)
	
	manager := notification.NewNotificationManager()
	manager.AddChannel(channels.NewMailChannel(emailService))
	manager.AddChannel(channels.NewSMSChannel(smsService))
	
	dispatcher := notification.NewDispatcher(manager)
	
	user := notification.NewSimpleNotifiable("user@example.com", "+1234567890")
	welcome := &WelcomeNotification{UserName: "John Doe"}
	
	err := dispatcher.Send(context.Background(), user, welcome)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}
	
	fmt.Println("Notification sent successfully!")
}