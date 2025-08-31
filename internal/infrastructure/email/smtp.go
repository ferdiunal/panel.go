package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

type SMTPProvider struct {
	config *SMTPConfig
}

func NewSMTPProvider(config *SMTPConfig) *SMTPProvider {
	return &SMTPProvider{
		config: config,
	}
}

func (s *SMTPProvider) ValidateConfig() error {
	if s.config.Host == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if s.config.Port == 0 {
		return fmt.Errorf("SMTP port is required")
	}
	if s.config.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if s.config.Password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if s.config.From == "" {
		return fmt.Errorf("SMTP from address is required")
	}
	return nil
}

func (s *SMTPProvider) SendEmail(ctx context.Context, req *EmailRequest) error {
	if err := s.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid SMTP config: %w", err)
	}

	if len(req.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	from := req.From
	if from == "" {
		from = s.config.From
	}

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	body := s.buildEmailBody(req, from)

	recipients := make([]string, 0, len(req.To)+len(req.CC)+len(req.BCC))
	recipients = append(recipients, req.To...)
	recipients = append(recipients, req.CC...)
	recipients = append(recipients, req.BCC...)

	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, from, recipients, body)
	}

	return smtp.SendMail(addr, auth, from, recipients, []byte(body))
}

func (s *SMTPProvider) SendTemplate(ctx context.Context, req *TemplateRequest) error {
	return fmt.Errorf("template sending not implemented for SMTP provider")
}

func (s *SMTPProvider) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, body string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	if s.config.UseTLS {
		tlsConfig := &tls.Config{
			ServerName: s.config.Host,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write([]byte(body))
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}

func (s *SMTPProvider) buildEmailBody(req *EmailRequest, from string) string {
	var body strings.Builder

	body.WriteString(fmt.Sprintf("From: %s\r\n", from))
	body.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(req.To, ", ")))
	
	if len(req.CC) > 0 {
		body.WriteString(fmt.Sprintf("CC: %s\r\n", strings.Join(req.CC, ", ")))
	}
	
	if req.ReplyTo != "" {
		body.WriteString(fmt.Sprintf("Reply-To: %s\r\n", req.ReplyTo))
	}
	
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", req.Subject))
	body.WriteString("MIME-Version: 1.0\r\n")

	if req.HTML != "" {
		body.WriteString("Content-Type: multipart/alternative; boundary=boundary123\r\n\r\n")
		
		if req.Text != "" {
			body.WriteString("--boundary123\r\n")
			body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
			body.WriteString(req.Text)
			body.WriteString("\r\n\r\n")
		}
		
		body.WriteString("--boundary123\r\n")
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		body.WriteString(req.HTML)
		body.WriteString("\r\n\r\n--boundary123--\r\n")
	} else {
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		body.WriteString(req.Text)
	}

	return body.String()
}