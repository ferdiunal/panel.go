package notification

import (
	"context"
	"errors"
	"testing"
)

func TestNewDispatcher(t *testing.T) {
	manager := NewNotificationManager()
	dispatcher := NewDispatcher(manager)
	
	if dispatcher == nil {
		t.Fatal("expected dispatcher to be created")
	}
	
	if dispatcher.manager != manager {
		t.Error("expected manager to be set correctly")
	}
}

func TestDispatcher_Send(t *testing.T) {
	manager := NewNotificationManager()
	channel := &mockChannel{name: "mail"}
	manager.AddChannel(channel)
	
	dispatcher := NewDispatcher(manager)
	notifiable := NewSimpleNotifiable("test@example.com", "+1234567890")
	notification := &testNotification{
		channels: []string{"mail"},
	}
	
	err := dispatcher.Send(context.Background(), notifiable, notification)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if !channel.sendCalled {
		t.Error("expected channel to be called")
	}
}

func TestDispatcher_SendVia(t *testing.T) {
	manager := NewNotificationManager()
	mailChannel := &mockChannel{name: "mail"}
	smsChannel := &mockChannel{name: "sms"}
	manager.AddChannel(mailChannel)
	manager.AddChannel(smsChannel)
	
	dispatcher := NewDispatcher(manager)
	notifiable := NewSimpleNotifiable("test@example.com", "+1234567890")
	notification := &testNotification{
		channels: []string{"mail", "sms"},
	}
	
	err := dispatcher.SendVia(context.Background(), []string{"sms"}, notifiable, notification)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if mailChannel.sendCalled {
		t.Error("expected mail channel not to be called")
	}
	
	if !smsChannel.sendCalled {
		t.Error("expected SMS channel to be called")
	}
}

func TestDispatcher_SendToMany(t *testing.T) {
	tests := []struct {
		name          string
		notifiables   int
		channelError  error
		expectedError bool
	}{
		{
			name:          "successful send to many",
			notifiables:   3,
			channelError:  nil,
			expectedError: false,
		},
		{
			name:          "partial failure",
			notifiables:   3,
			channelError:  errors.New("send failed"),
			expectedError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewNotificationManager()
			channel := &mockChannel{
				name:      "mail",
				sendError: tt.channelError,
			}
			manager.AddChannel(channel)
			
			dispatcher := NewDispatcher(manager)
			
			notifiables := make([]Notifiable, tt.notifiables)
			for i := 0; i < tt.notifiables; i++ {
				notifiables[i] = NewSimpleNotifiable("test@example.com", "+1234567890")
			}
			
			notification := &testNotification{
				channels: []string{"mail"},
			}
			
			err := dispatcher.SendToMany(context.Background(), notifiables, notification)
			
			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestDispatcher_SendToManyVia(t *testing.T) {
	manager := NewNotificationManager()
	mailChannel := &mockChannel{name: "mail"}
	smsChannel := &mockChannel{name: "sms"}
	manager.AddChannel(mailChannel)
	manager.AddChannel(smsChannel)
	
	dispatcher := NewDispatcher(manager)
	
	notifiables := []Notifiable{
		NewSimpleNotifiable("test1@example.com", "+1234567890"),
		NewSimpleNotifiable("test2@example.com", "+1234567891"),
	}
	
	notification := &testNotification{
		channels: []string{"mail", "sms"},
	}
	
	err := dispatcher.SendToManyVia(context.Background(), []string{"mail"}, notifiables, notification)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if !mailChannel.sendCalled {
		t.Error("expected mail channel to be called")
	}
	
	if smsChannel.sendCalled {
		t.Error("expected SMS channel not to be called")
	}
}

func TestDispatcher_GetAvailableChannels(t *testing.T) {
	manager := NewNotificationManager()
	manager.AddChannel(&mockChannel{name: "mail"})
	manager.AddChannel(&mockChannel{name: "sms"})
	
	dispatcher := NewDispatcher(manager)
	channels := dispatcher.GetAvailableChannels()
	
	if len(channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(channels))
	}
}