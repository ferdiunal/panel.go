package sms

import (
	"fmt"
	"os"
)

type ProviderType string

const (
	Twilio ProviderType = "twilio"
)

type FactoryConfig struct {
	Provider ProviderType
	Twilio   *TwilioConfig
}

func NewSMSProviderFromEnv() (SMSProvider, error) {
	providerType := os.Getenv("SMS_PROVIDER")
	if providerType == "" {
		providerType = "twilio"
	}

	switch ProviderType(providerType) {
	case Twilio:
		return createTwilioProviderFromEnv()
	default:
		return nil, fmt.Errorf("unsupported SMS provider: %s", providerType)
	}
}

func NewSMSProvider(config *FactoryConfig) (SMSProvider, error) {
	switch config.Provider {
	case Twilio:
		if config.Twilio == nil {
			return nil, fmt.Errorf("Twilio config is required")
		}
		return NewTwilioProvider(config.Twilio), nil
	default:
		return nil, fmt.Errorf("unsupported SMS provider: %s", config.Provider)
	}
}

func createTwilioProviderFromEnv() (SMSProvider, error) {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if accountSID == "" {
		return nil, fmt.Errorf("TWILIO_ACCOUNT_SID environment variable is required")
	}

	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("TWILIO_AUTH_TOKEN environment variable is required")
	}

	from := os.Getenv("TWILIO_FROM")
	if from == "" {
		return nil, fmt.Errorf("TWILIO_FROM environment variable is required")
	}

	config := &TwilioConfig{
		AccountSID: accountSID,
		AuthToken:  authToken,
		From:       from,
	}

	return NewTwilioProvider(config), nil
}