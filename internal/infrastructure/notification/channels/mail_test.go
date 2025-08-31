package channels

import (
	"context"
	"errors"
	"testing"

	"panel.go/internal/infrastructure/email"
	"panel.go/internal/infrastructure/notification"
)

type mockEmailService struct {
	sendError  error
	sendCalled bool
	lastReq    *email.EmailRequest
}

func (m *mockEmailService) SendEmail(ctx context.Context, req *email.EmailRequest) error {
	m.sendCalled = true
	m.lastReq = req
	return m.sendError
}

func (m *mockEmailService) SendTemplate(ctx context.Context, req *email.TemplateRequest) error {
	return errors.New("not implemented")
}

type testMailNotification struct {
	notification.BaseNotification
	mailMessage *notification.MailMessage
}

func (t *testMailNotification) Via(notifiable notification.Notifiable) []string {
	return []string{"mail"}
}

func (t *testMailNotification) ToMail(notifiable notification.Notifiable) *notification.MailMessage {
	return t.mailMessage
}

func TestNewMailChannel(t *testing.T) {
	emailService := &mockEmailService{}
	channel := NewMailChannel(emailService)
	
	if channel == nil {
		t.Fatal("expected channel to be created")
	}
	
	if channel.Name() != MailChannelName {
		t.Errorf("expected channel name to be %s, got %s", MailChannelName, channel.Name())
	}
}

func TestMailChannel_Send(t *testing.T) {
	tests := []struct {
		name          string
		mailMessage   *notification.MailMessage
		emailAddr     string
		sendError     error
		expectedError string
	}{
		{
			name: "successful send",
			mailMessage: notification.NewMailMessage().
				SetSubject("Test").
				SetBody("Test body"),
			emailAddr:     "test@example.com",
			sendError:     nil,
			expectedError: "",
		},
		{
			name:          "no mail message",
			mailMessage:   nil,
			emailAddr:     "test@example.com",
			sendError:     nil,
			expectedError: "notification does not support mail channel",
		},
		{
			name: "no email address",
			mailMessage: notification.NewMailMessage().
				SetSubject("Test").
				SetBody("Test body"),
			emailAddr:     "",
			sendError:     nil,
			expectedError: "no mail address found for notifiable",
		},
		{
			name: "email service error",
			mailMessage: notification.NewMailMessage().
				SetSubject("Test").
				SetBody("Test body"),
			emailAddr:     "test@example.com",
			sendError:     errors.New("send failed"),
			expectedError: "send failed",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailService := &mockEmailService{sendError: tt.sendError}
			channel := NewMailChannel(emailService)
			
			notifiable := notification.NewSimpleNotifiable(tt.emailAddr, "+1234567890")
			notif := &testMailNotification{mailMessage: tt.mailMessage}
			
			err := channel.Send(context.Background(), notifiable, notif)
			
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.mailMessage != nil && !emailService.sendCalled {
					t.Error("expected email service to be called")
				}
				if tt.mailMessage != nil && emailService.lastReq != nil {
					if len(emailService.lastReq.To) == 0 || emailService.lastReq.To[0] != tt.emailAddr {
						t.Errorf("expected email to be sent to %s, got %v", tt.emailAddr, emailService.lastReq.To)
					}
					if emailService.lastReq.Subject != tt.mailMessage.Subject {
						t.Errorf("expected subject %s, got %s", tt.mailMessage.Subject, emailService.lastReq.Subject)
					}
				}
			} else {
				if err == nil {
					t.Error("expected error but got none")
				} else if !containsString(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(substr) <= len(s) && func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()))
}