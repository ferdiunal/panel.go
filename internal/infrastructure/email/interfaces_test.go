package email

import (
	"context"
	"errors"
	"testing"
)

type mockEmailProvider struct {
	validateError error
	sendError     error
	sendCalled    bool
	templateCalled bool
}

func (m *mockEmailProvider) ValidateConfig() error {
	return m.validateError
}

func (m *mockEmailProvider) SendEmail(ctx context.Context, req *EmailRequest) error {
	m.sendCalled = true
	return m.sendError
}

func (m *mockEmailProvider) SendTemplate(ctx context.Context, req *TemplateRequest) error {
	m.templateCalled = true
	return m.sendError
}

func TestNewEmailService(t *testing.T) {
	provider := &mockEmailProvider{}
	service := NewEmailService(provider)
	
	if service == nil {
		t.Fatal("expected service to be created")
	}
	
	if service.provider != provider {
		t.Error("expected provider to be set correctly")
	}
}

func TestEmailService_SendEmail(t *testing.T) {
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
			provider := &mockEmailProvider{
				validateError: tt.validateError,
				sendError:     tt.sendError,
			}
			service := NewEmailService(provider)
			
			req := &EmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test",
				Text:    "Test message",
			}
			
			err := service.SendEmail(context.Background(), req)
			
			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			if tt.validateError == nil && !provider.sendCalled {
				t.Error("expected SendEmail to be called on provider")
			}
		})
	}
}

func TestEmailService_SendTemplate(t *testing.T) {
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
			provider := &mockEmailProvider{
				validateError: tt.validateError,
				sendError:     tt.sendError,
			}
			service := NewEmailService(provider)
			
			req := &TemplateRequest{
				To:         []string{"test@example.com"},
				Subject:    "Test",
				TemplateID: "template-123",
			}
			
			err := service.SendTemplate(context.Background(), req)
			
			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			if tt.validateError == nil && !provider.templateCalled {
				t.Error("expected SendTemplate to be called on provider")
			}
		})
	}
}