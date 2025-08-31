package notification

import (
	"context"
	"errors"
	"testing"
)

type mockChannel struct {
	name      string
	sendError error
	sendCalled bool
}

func (m *mockChannel) Name() string {
	return m.name
}

func (m *mockChannel) Send(ctx context.Context, notifiable Notifiable, notification Notification) error {
	m.sendCalled = true
	return m.sendError
}

type testNotification struct {
	BaseNotification
	channels    []string
	shouldSend  map[string]bool
	mailMessage *MailMessage
	smsMessage  *SMSMessage
}

func (t *testNotification) Via(notifiable Notifiable) []string {
	return t.channels
}

func (t *testNotification) ToMail(notifiable Notifiable) *MailMessage {
	return t.mailMessage
}

func (t *testNotification) ToSMS(notifiable Notifiable) *SMSMessage {
	return t.smsMessage
}

func (t *testNotification) ShouldSend(notifiable Notifiable, channel string) bool {
	if t.shouldSend == nil {
		return true
	}
	return t.shouldSend[channel]
}

func TestNewNotificationManager(t *testing.T) {
	manager := NewNotificationManager()
	
	if manager == nil {
		t.Fatal("expected manager to be created")
	}
	
	if manager.channels == nil {
		t.Error("expected channels map to be initialized")
	}
}

func TestNotificationManager_AddChannel(t *testing.T) {
	manager := NewNotificationManager()
	channel := &mockChannel{name: "test"}
	
	manager.AddChannel(channel)
	
	if len(manager.channels) != 1 {
		t.Errorf("expected 1 channel, got %d", len(manager.channels))
	}
	
	if manager.channels["test"] != channel {
		t.Error("channel not added correctly")
	}
}

func TestNotificationManager_Send(t *testing.T) {
	tests := []struct {
		name           string
		channels       []string
		availChannels  []string
		shouldSend     map[string]bool
		channelError   error
		expectedError  bool
		expectedCalls  int
	}{
		{
			name:          "successful send",
			channels:      []string{"mail"},
			availChannels: []string{"mail"},
			expectedError: false,
			expectedCalls: 1,
		},
		{
			name:          "channel not available",
			channels:      []string{"unavailable"},
			availChannels: []string{"mail"},
			expectedError: false,
			expectedCalls: 0,
		},
		{
			name:          "should not send",
			channels:      []string{"mail"},
			availChannels: []string{"mail"},
			shouldSend:    map[string]bool{"mail": false},
			expectedError: false,
			expectedCalls: 0,
		},
		{
			name:          "channel send error",
			channels:      []string{"mail"},
			availChannels: []string{"mail"},
			channelError:  errors.New("send failed"),
			expectedError: true,
			expectedCalls: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewNotificationManager()
			
			channelMocks := make(map[string]*mockChannel)
			for _, chName := range tt.availChannels {
				mock := &mockChannel{
					name:      chName,
					sendError: tt.channelError,
				}
				channelMocks[chName] = mock
				manager.AddChannel(mock)
			}
			
			notifiable := NewSimpleNotifiable("test@example.com", "+1234567890")
			notification := &testNotification{
				channels:   tt.channels,
				shouldSend: tt.shouldSend,
			}
			
			err := manager.Send(context.Background(), notifiable, notification)
			
			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			totalCalls := 0
			for _, mock := range channelMocks {
				if mock.sendCalled {
					totalCalls++
				}
			}
			
			if totalCalls != tt.expectedCalls {
				t.Errorf("expected %d channel calls, got %d", tt.expectedCalls, totalCalls)
			}
		})
	}
}

func TestNotificationManager_GetChannel(t *testing.T) {
	manager := NewNotificationManager()
	channel := &mockChannel{name: "test"}
	manager.AddChannel(channel)
	
	foundChannel, exists := manager.GetChannel("test")
	if !exists {
		t.Error("expected channel to exist")
	}
	
	if foundChannel != channel {
		t.Error("expected correct channel to be returned")
	}
	
	_, exists = manager.GetChannel("nonexistent")
	if exists {
		t.Error("expected nonexistent channel to not exist")
	}
}

func TestNotificationManager_GetAvailableChannels(t *testing.T) {
	manager := NewNotificationManager()
	
	channels := manager.GetAvailableChannels()
	if len(channels) != 0 {
		t.Errorf("expected 0 channels, got %d", len(channels))
	}
	
	manager.AddChannel(&mockChannel{name: "mail"})
	manager.AddChannel(&mockChannel{name: "sms"})
	
	channels = manager.GetAvailableChannels()
	if len(channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(channels))
	}
}