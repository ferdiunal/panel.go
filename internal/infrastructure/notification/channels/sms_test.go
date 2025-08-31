package channels

import (
	"context"
	"errors"
	"testing"

	"panel.go/internal/infrastructure/notification"
	"panel.go/internal/infrastructure/sms"
)

type mockSMSService struct {
	sendError  error
	sendCalled bool
	lastReq    *sms.SMSRequest
}

func (m *mockSMSService) SendSMS(ctx context.Context, req *sms.SMSRequest) error {
	m.sendCalled = true
	m.lastReq = req
	return m.sendError
}

type testSMSNotification struct {
	notification.BaseNotification
	smsMessage *notification.SMSMessage
}

func (t *testSMSNotification) Via(notifiable notification.Notifiable) []string {
	return []string{"sms"}
}

func (t *testSMSNotification) ToSMS(notifiable notification.Notifiable) *notification.SMSMessage {
	return t.smsMessage
}

func TestNewSMSChannel(t *testing.T) {
	smsService := &mockSMSService{}
	channel := NewSMSChannel(smsService)
	
	if channel == nil {
		t.Fatal("expected channel to be created")
	}
	
	if channel.Name() != SMSChannelName {
		t.Errorf("expected channel name to be %s, got %s", SMSChannelName, channel.Name())
	}
}

func TestSMSChannel_Send(t *testing.T) {
	tests := []struct {
		name          string
		smsMessage    *notification.SMSMessage
		phoneNumber   string
		sendError     error
		expectedError string
	}{
		{
			name: "successful send",
			smsMessage: notification.NewSMSMessage().
				SetBody("Test SMS message"),
			phoneNumber:   "+1234567890",
			sendError:     nil,
			expectedError: "",
		},
		{
			name:          "no SMS message",
			smsMessage:    nil,
			phoneNumber:   "+1234567890",
			sendError:     nil,
			expectedError: "notification does not support SMS channel",
		},
		{
			name: "no phone number",
			smsMessage: notification.NewSMSMessage().
				SetBody("Test SMS message"),
			phoneNumber:   "",
			sendError:     nil,
			expectedError: "no phone number found for notifiable",
		},
		{
			name: "SMS service error",
			smsMessage: notification.NewSMSMessage().
				SetBody("Test SMS message"),
			phoneNumber:   "+1234567890",
			sendError:     errors.New("send failed"),
			expectedError: "send failed",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smsService := &mockSMSService{sendError: tt.sendError}
			channel := NewSMSChannel(smsService)
			
			notifiable := notification.NewSimpleNotifiable("test@example.com", tt.phoneNumber)
			notif := &testSMSNotification{smsMessage: tt.smsMessage}
			
			err := channel.Send(context.Background(), notifiable, notif)
			
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.smsMessage != nil && !smsService.sendCalled {
					t.Error("expected SMS service to be called")
				}
				if tt.smsMessage != nil && smsService.lastReq != nil {
					if smsService.lastReq.To != tt.phoneNumber {
						t.Errorf("expected SMS to be sent to %s, got %s", tt.phoneNumber, smsService.lastReq.To)
					}
					if smsService.lastReq.Body != tt.smsMessage.Body {
						t.Errorf("expected body %s, got %s", tt.smsMessage.Body, smsService.lastReq.Body)
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