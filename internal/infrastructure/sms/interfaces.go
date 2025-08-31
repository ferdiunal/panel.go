package sms

import "context"

type SMSProvider interface {
	SendSMS(ctx context.Context, req *SMSRequest) error
	ValidateConfig() error
}

type SMSRequest struct {
	To   string
	Body string
	From string
}

type SMSService struct {
	provider SMSProvider
}

func NewSMSService(provider SMSProvider) *SMSService {
	return &SMSService{
		provider: provider,
	}
}

func (s *SMSService) SendSMS(ctx context.Context, req *SMSRequest) error {
	if err := s.provider.ValidateConfig(); err != nil {
		return err
	}
	return s.provider.SendSMS(ctx, req)
}