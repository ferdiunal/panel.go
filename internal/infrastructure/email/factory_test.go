package email

import (
	"os"
	"testing"
)

func TestNewEmailProvider(t *testing.T) {
	tests := []struct {
		name        string
		config      *FactoryConfig
		expectedErr string
	}{
		{
			name: "valid SMTP config",
			config: &FactoryConfig{
				Provider: SMTP,
				SMTP: &SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					Username: "user@example.com",
					Password: "password",
					From:     "noreply@example.com",
				},
			},
			expectedErr: "",
		},
		{
			name: "valid Resend config",
			config: &FactoryConfig{
				Provider: Resend,
				Resend: &ResendConfig{
					APIKey: "test-api-key",
					From:   "noreply@example.com",
				},
			},
			expectedErr: "",
		},
		{
			name: "SMTP config missing",
			config: &FactoryConfig{
				Provider: SMTP,
			},
			expectedErr: "SMTP config is required",
		},
		{
			name: "Resend config missing",
			config: &FactoryConfig{
				Provider: Resend,
			},
			expectedErr: "Resend config is required",
		},
		{
			name: "unsupported provider",
			config: &FactoryConfig{
				Provider: ProviderType("unsupported"),
			},
			expectedErr: "unsupported email provider",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewEmailProvider(tt.config)
			
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("expected provider to be created")
				}
			} else {
				if err == nil {
					t.Error("expected error but got none")
				} else if !containsString(err.Error(), tt.expectedErr) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestNewEmailProviderFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedErr string
	}{
		{
			name: "SMTP provider from env",
			envVars: map[string]string{
				"EMAIL_PROVIDER": "smtp",
				"SMTP_HOST":      "smtp.example.com",
				"SMTP_PORT":      "587",
				"SMTP_USERNAME":  "user@example.com",
				"SMTP_PASSWORD":  "password",
				"SMTP_FROM":      "noreply@example.com",
			},
			expectedErr: "",
		},
		{
			name: "Resend provider from env",
			envVars: map[string]string{
				"EMAIL_PROVIDER": "resend",
				"RESEND_API_KEY": "test-api-key",
				"RESEND_FROM":    "noreply@example.com",
			},
			expectedErr: "",
		},
		{
			name: "default to SMTP when no provider specified",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_PORT":     "587",
				"SMTP_USERNAME": "user@example.com",
				"SMTP_PASSWORD": "password",
				"SMTP_FROM":     "noreply@example.com",
			},
			expectedErr: "",
		},
		{
			name: "missing SMTP host",
			envVars: map[string]string{
				"EMAIL_PROVIDER": "smtp",
			},
			expectedErr: "SMTP_HOST environment variable is required",
		},
		{
			name: "missing Resend API key",
			envVars: map[string]string{
				"EMAIL_PROVIDER": "resend",
			},
			expectedErr: "RESEND_API_KEY environment variable is required",
		},
		{
			name: "unsupported provider",
			envVars: map[string]string{
				"EMAIL_PROVIDER": "unsupported",
			},
			expectedErr: "unsupported email provider",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupEnv := setEnvVars(tt.envVars)
			defer cleanupEnv()
			
			provider, err := NewEmailProviderFromEnv()
			
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("expected provider to be created")
				}
			} else {
				if err == nil {
					t.Error("expected error but got none")
				} else if !containsString(err.Error(), tt.expectedErr) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestCreateSMTPProviderFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedErr string
		checkTLS    bool
		expectedTLS bool
	}{
		{
			name: "complete SMTP config",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_PORT":     "587",
				"SMTP_USERNAME": "user@example.com",
				"SMTP_PASSWORD": "password",
				"SMTP_FROM":     "noreply@example.com",
				"SMTP_USE_TLS":  "true",
			},
			expectedErr: "",
			checkTLS:    true,
			expectedTLS: true,
		},
		{
			name: "TLS set to false",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_PORT":     "587",
				"SMTP_USERNAME": "user@example.com",
				"SMTP_PASSWORD": "password",
				"SMTP_FROM":     "noreply@example.com",
				"SMTP_USE_TLS":  "false",
			},
			expectedErr: "",
			checkTLS:    true,
			expectedTLS: false,
		},
		{
			name: "default port when not specified",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_USERNAME": "user@example.com",
				"SMTP_PASSWORD": "password",
				"SMTP_FROM":     "noreply@example.com",
			},
			expectedErr: "",
		},
		{
			name: "invalid port",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_PORT":     "invalid",
				"SMTP_USERNAME": "user@example.com",
				"SMTP_PASSWORD": "password",
				"SMTP_FROM":     "noreply@example.com",
			},
			expectedErr: "invalid SMTP_PORT",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupEnv := setEnvVars(tt.envVars)
			defer cleanupEnv()
			
			provider, err := createSMTPProviderFromEnv()
			
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("expected provider to be created")
				}
				
				if tt.checkTLS {
					smtpProvider, ok := provider.(*SMTPProvider)
					if !ok {
						t.Error("expected SMTPProvider type")
					} else if smtpProvider.config.UseTLS != tt.expectedTLS {
						t.Errorf("expected UseTLS to be %v, got %v", tt.expectedTLS, smtpProvider.config.UseTLS)
					}
				}
			} else {
				if err == nil {
					t.Error("expected error but got none")
				} else if !containsString(err.Error(), tt.expectedErr) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestCreateResendProviderFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedErr string
	}{
		{
			name: "complete Resend config",
			envVars: map[string]string{
				"RESEND_API_KEY": "test-api-key",
				"RESEND_FROM":    "noreply@example.com",
			},
			expectedErr: "",
		},
		{
			name:        "missing API key",
			envVars:     map[string]string{},
			expectedErr: "RESEND_API_KEY environment variable is required",
		},
		{
			name: "missing from address",
			envVars: map[string]string{
				"RESEND_API_KEY": "test-api-key",
			},
			expectedErr: "RESEND_FROM environment variable is required",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupEnv := setEnvVars(tt.envVars)
			defer cleanupEnv()
			
			provider, err := createResendProviderFromEnv()
			
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("expected provider to be created")
				}
			} else {
				if err == nil {
					t.Error("expected error but got none")
				} else if !containsString(err.Error(), tt.expectedErr) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func setEnvVars(vars map[string]string) func() {
	originalVars := make(map[string]string)
	
	for key, value := range vars {
		originalVars[key] = os.Getenv(key)
		os.Setenv(key, value)
	}
	
	return func() {
		for key, originalValue := range originalVars {
			if originalValue == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, originalValue)
			}
		}
	}
}