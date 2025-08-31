package email

import (
	"fmt"
	"os"
	"strconv"
)

type ProviderType string

const (
	SMTP   ProviderType = "smtp"
	Resend ProviderType = "resend"
)

type FactoryConfig struct {
	Provider ProviderType
	SMTP     *SMTPConfig
	Resend   *ResendConfig
}

func NewEmailProviderFromEnv() (EmailProvider, error) {
	providerType := os.Getenv("EMAIL_PROVIDER")
	if providerType == "" {
		providerType = "smtp"
	}

	switch ProviderType(providerType) {
	case SMTP:
		return createSMTPProviderFromEnv()
	case Resend:
		return createResendProviderFromEnv()
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", providerType)
	}
}

func NewEmailProvider(config *FactoryConfig) (EmailProvider, error) {
	switch config.Provider {
	case SMTP:
		if config.SMTP == nil {
			return nil, fmt.Errorf("SMTP config is required")
		}
		return NewSMTPProvider(config.SMTP), nil
	case Resend:
		if config.Resend == nil {
			return nil, fmt.Errorf("Resend config is required")
		}
		return NewResendProvider(config.Resend), nil
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", config.Provider)
	}
}

func createSMTPProviderFromEnv() (EmailProvider, error) {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return nil, fmt.Errorf("SMTP_HOST environment variable is required")
	}

	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	username := os.Getenv("SMTP_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("SMTP_USERNAME environment variable is required")
	}

	password := os.Getenv("SMTP_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD environment variable is required")
	}

	from := os.Getenv("SMTP_FROM")
	if from == "" {
		return nil, fmt.Errorf("SMTP_FROM environment variable is required")
	}

	useTLSStr := os.Getenv("SMTP_USE_TLS")
	useTLS := useTLSStr == "true" || useTLSStr == "1"

	config := &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		UseTLS:   useTLS,
	}

	return NewSMTPProvider(config), nil
}

func createResendProviderFromEnv() (EmailProvider, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY environment variable is required")
	}

	from := os.Getenv("RESEND_FROM")
	if from == "" {
		return nil, fmt.Errorf("RESEND_FROM environment variable is required")
	}

	config := &ResendConfig{
		APIKey: apiKey,
		From:   from,
	}

	return NewResendProvider(config), nil
}