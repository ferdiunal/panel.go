package channels

import (
	"context"
	"fmt"

	"panel.go/internal/infrastructure/notification"
	"panel.go/internal/infrastructure/sms"
)

const SMSChannelName = "sms"

type SMSSender interface {
	SendSMS(ctx context.Context, req *sms.SMSRequest) error
}

type SMSChannel struct {
	smsSender SMSSender
}

func NewSMSChannel(smsSender SMSSender) *SMSChannel {
	return &SMSChannel{
		smsSender: smsSender,
	}
}

func (c *SMSChannel) Name() string {
	return SMSChannelName
}

func (c *SMSChannel) Send(ctx context.Context, notifiable notification.Notifiable, notif notification.Notification) error {
	smsMessage := notif.ToSMS(notifiable)
	if smsMessage == nil {
		return fmt.Errorf("notification does not support SMS channel")
	}

	to := notifiable.RouteNotificationForSMS()
	if to == "" {
		return fmt.Errorf("no phone number found for notifiable")
	}

	smsReq := &sms.SMSRequest{
		To:   to,
		Body: smsMessage.Body,
		From: smsMessage.From,
	}

	return c.smsSender.SendSMS(ctx, smsReq)
}