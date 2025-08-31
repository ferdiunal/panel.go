package notification_test

import (
	"context"
	"fmt"
	"io"

	"panel.go/internal/infrastructure/email"
	"panel.go/internal/infrastructure/notification"
	"panel.go/internal/infrastructure/notification/channels"
)

type WelcomeTemplateNotification struct {
	notification.BaseNotification
	renderer notification.TemplateRenderer
	userName string
	email    string
}

func (w *WelcomeTemplateNotification) Via(notifiable notification.Notifiable) []string {
	return []string{"mail"}
}

func (w *WelcomeTemplateNotification) ToMail(notifiable notification.Notifiable) *notification.MailMessage {
	templateMail := notification.NewTemplateMailMessage(w.renderer)
	templateMail.SetTemplate("welcome").
		SetData("userName", w.userName).
		SetData("email", w.email).
		SetSubject("Welcome to Panel.GO!")

	if err := templateMail.RenderContent(context.Background()); err != nil {
		return nil
	}

	return templateMail.MailMessage
}

func ExampleTemplateNotification() {
	renderer := notification.NewEmailTemplateRenderer()
	
	renderer.RegisterTemplate("welcome", func(data notification.TemplateData) notification.TemplComponent {
		return &mockComponent{
			html: fmt.Sprintf("<h1>Welcome %s!</h1><p>Email: %s</p>", 
				data["userName"], data["email"]),
		}
	})
	
	renderer.RegisterTemplate("welcome_text", func(data notification.TemplateData) notification.TemplComponent {
		return &mockComponent{
			html: fmt.Sprintf("Welcome %s! Email: %s", 
				data["userName"], data["email"]),
		}
	})
	
	emailProvider, _ := email.NewEmailProviderFromEnv()
	emailService := email.NewEmailService(emailProvider)
	
	manager := notification.NewNotificationManager()
	manager.AddChannel(channels.NewMailChannel(emailService))
	
	dispatcher := notification.NewDispatcher(manager)
	
	user := notification.NewSimpleNotifiable("user@example.com", "+1234567890")
	welcome := &WelcomeTemplateNotification{
		renderer: renderer,
		userName: "John Doe", 
		email:    "user@example.com",
	}
	
	err := dispatcher.Send(context.Background(), user, welcome)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
		return
	}
	
	fmt.Println("Template notification sent successfully!")
}

type mockComponent struct {
	html string
}

func (m *mockComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(m.html))
	return err
}