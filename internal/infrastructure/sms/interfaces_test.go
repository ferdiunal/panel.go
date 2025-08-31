package sms

import (
	"context"
	"errors"
	"testing"
)

type mockSMSProvider struct {
	validateError error
	sendError     error
	sendCalled    bool
}

func (m *mockSMSProvider) ValidateConfig() error {
	return m.validateError
}

func (m *mockSMSProvider) SendSMS(ctx context.Context, req *SMSRequest) error {
	m.sendCalled = true
	return m.sendError
}

func TestNewSMSService(t *testing.T) {
	provider := &mockSMSProvider{}
	service := NewSMSService(provider)
	
	if service == nil {
		t.Fatal("expected service to be created")
	}
	
	if service.provider != provider {
		t.Error("expected provider to be set correctly")
	}
}

func TestSMSService_SendSMS(t *testing.T) {
	tests := []struct {
		name          string
		validateError error
		sendError     error
		expectedError bool
	}{
		{
			name:          "successful send",
			validateError: nil,
			sendError:     nil,
			expectedError: false,
		},
		{
			name:          "validation error",
			validateError: errors.New("validation failed"),
			sendError:     nil,
			expectedError: true,
		},
		{
			name:          "send error",
			validateError: nil,
			sendError:     errors.New("send failed"),
			expectedError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &mockSMSProvider{
				validateError: tt.validateError,
				sendError:     tt.sendError,
			}
			service := NewSMSService(provider)
			
			req := &SMSRequest{
				To:   "+1234567890",
				Body: "Test message",
			}
			
			err := service.SendSMS(context.Background(), req)
			
			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			if tt.validateError == nil && !provider.sendCalled {
				t.Error("expected SendSMS to be called on provider")
			}
		})
	}
}