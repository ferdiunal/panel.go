package email

import "context"

type EmailProvider interface {
	SendEmail(ctx context.Context, req *EmailRequest) error
	SendTemplate(ctx context.Context, req *TemplateRequest) error
	ValidateConfig() error
}

type EmailRequest struct {
	To      []string
	CC      []string
	BCC     []string
	Subject string
	Text    string
	HTML    string
	From    string
	ReplyTo string
}

type TemplateRequest struct {
	To           []string
	CC           []string
	BCC          []string
	Subject      string
	TemplateID   string
	TemplateData map[string]any
	From         string
	ReplyTo      string
}

type EmailService struct {
	provider EmailProvider
}

func NewEmailService(provider EmailProvider) *EmailService {
	return &EmailService{
		provider: provider,
	}
}

func (s *EmailService) SendEmail(ctx context.Context, req *EmailRequest) error {
	if err := s.provider.ValidateConfig(); err != nil {
		return err
	}
	return s.provider.SendEmail(ctx, req)
}

func (s *EmailService) SendTemplate(ctx context.Context, req *TemplateRequest) error {
	if err := s.provider.ValidateConfig(); err != nil {
		return err
	}
	return s.provider.SendTemplate(ctx, req)
}