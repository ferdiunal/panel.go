package channels

import (
	"context"
	"fmt"

	"panel.go/internal/infrastructure/email"
	"panel.go/internal/infrastructure/notification"
)

const MailChannelName = "mail"

type EmailSender interface {
	SendEmail(ctx context.Context, req *email.EmailRequest) error
}

type MailChannel struct {
	emailSender EmailSender
}

func NewMailChannel(emailSender EmailSender) *MailChannel {
	return &MailChannel{
		emailSender: emailSender,
	}
}

func (c *MailChannel) Name() string {
	return MailChannelName
}

func (c *MailChannel) Send(ctx context.Context, notifiable notification.Notifiable, notif notification.Notification) error {
	mailMessage := notif.ToMail(notifiable)
	if mailMessage == nil {
		return fmt.Errorf("notification does not support mail channel")
	}

	to := notifiable.RouteNotificationForMail()
	if to == "" {
		return fmt.Errorf("no mail address found for notifiable")
	}

	emailReq := &email.EmailRequest{
		To:      []string{to},
		Subject: mailMessage.Subject,
		Text:    mailMessage.Body,
		HTML:    mailMessage.HTMLBody,
		From:    mailMessage.From,
		ReplyTo: mailMessage.ReplyTo,
		CC:      mailMessage.CC,
		BCC:     mailMessage.BCC,
	}

	return c.emailSender.SendEmail(ctx, emailReq)
}