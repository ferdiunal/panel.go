package email

import (
	"context"
	"testing"
)

func TestNewSMTPProvider(t *testing.T) {
	config := &SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
		UseTLS:   true,
	}
	
	provider := NewSMTPProvider(config)
	if provider == nil {
		t.Fatal("expected provider to be created")
	}
	
	if provider.config != config {
		t.Error("expected config to be set correctly")
	}
}

func TestSMTPProvider_ValidateConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *SMTPConfig
		hasError bool
	}{
		{
			name: "valid config",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				Password: "password",
				From:     "noreply@example.com",
			},
			hasError: false,
		},
		{
			name: "missing host",
			config: &SMTPConfig{
				Port:     587,
				Username: "user@example.com",
				Password: "password",
				From:     "noreply@example.com",
			},
			hasError: true,
		},
		{
			name: "missing port",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Username: "user@example.com",
				Password: "password",
				From:     "noreply@example.com",
			},
			hasError: true,
		},
		{
			name: "missing username",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Password: "password",
				From:     "noreply@example.com",
			},
			hasError: true,
		},
		{
			name: "missing password",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				From:     "noreply@example.com",
			},
			hasError: true,
		},
		{
			name: "missing from",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				Password: "password",
			},
			hasError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewSMTPProvider(tt.config)
			err := provider.ValidateConfig()
			
			if tt.hasError && err == nil {
				t.Error("expected validation error but got none")
			}
			
			if !tt.hasError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestSMTPProvider_SendEmail_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      *SMTPConfig
		request     *EmailRequest
		expectedErr string
	}{
		{
			name: "invalid config",
			config: &SMTPConfig{
				Host: "smtp.example.com",
			},
			request: &EmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test",
				Text:    "Test message",
			},
			expectedErr: "invalid SMTP config",
		},
		{
			name: "no recipients",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				Password: "password",
				From:     "noreply@example.com",
			},
			request: &EmailRequest{
				Subject: "Test",
				Text:    "Test message",
			},
			expectedErr: "at least one recipient is required",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewSMTPProvider(tt.config)
			err := provider.SendEmail(context.Background(), tt.request)
			
			if err == nil {
				t.Error("expected error but got none")
			} else if tt.expectedErr != "" && !containsString(err.Error(), tt.expectedErr) {
				t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestSMTPProvider_SendTemplate(t *testing.T) {
	config := &SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
	}
	
	provider := NewSMTPProvider(config)
	req := &TemplateRequest{
		To:         []string{"test@example.com"},
		Subject:    "Test",
		TemplateID: "template-123",
	}
	
	err := provider.SendTemplate(context.Background(), req)
	if err == nil {
		t.Error("expected template sending to not be implemented")
	}
	
	if !containsString(err.Error(), "template sending not implemented") {
		t.Errorf("expected 'template sending not implemented' error, got: %v", err)
	}
}

func TestSMTPProvider_buildEmailBody(t *testing.T) {
	config := &SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
	}
	
	provider := NewSMTPProvider(config)
	
	tests := []struct {
		name     string
		request  *EmailRequest
		from     string
		contains []string
	}{
		{
			name: "text only email",
			request: &EmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Text:    "Test message",
			},
			from:     "sender@example.com",
			contains: []string{"From: sender@example.com", "To: test@example.com", "Subject: Test Subject", "Test message"},
		},
		{
			name: "html email",
			request: &EmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				HTML:    "<p>Test HTML</p>",
			},
			from:     "sender@example.com",
			contains: []string{"From: sender@example.com", "To: test@example.com", "Subject: Test Subject", "<p>Test HTML</p>", "multipart/alternative"},
		},
		{
			name: "email with CC and BCC",
			request: &EmailRequest{
				To:      []string{"test@example.com"},
				CC:      []string{"cc@example.com"},
				Subject: "Test Subject",
				Text:    "Test message",
			},
			from:     "sender@example.com",
			contains: []string{"From: sender@example.com", "To: test@example.com", "CC: cc@example.com", "Subject: Test Subject"},
		},
		{
			name: "email with reply-to",
			request: &EmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Text:    "Test message",
				ReplyTo: "reply@example.com",
			},
			from:     "sender@example.com",
			contains: []string{"From: sender@example.com", "Reply-To: reply@example.com", "Subject: Test Subject"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := provider.buildEmailBody(tt.request, tt.from)
			
			for _, expected := range tt.contains {
				if !containsString(body, expected) {
					t.Errorf("expected body to contain '%s', got: %s", expected, body)
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