package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ResendConfig struct {
	APIKey string
	From   string
}

type ResendProvider struct {
	config     *ResendConfig
	httpClient *http.Client
	baseURL    string
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Text    string   `json:"text,omitempty"`
	HTML    string   `json:"html,omitempty"`
	ReplyTo string   `json:"reply_to,omitempty"`
}

type resendTemplateRequest struct {
	From         string                 `json:"from"`
	To           []string               `json:"to"`
	CC           []string               `json:"cc,omitempty"`
	BCC          []string               `json:"bcc,omitempty"`
	Subject      string                 `json:"subject"`
	TemplateID   string                 `json:"template_id"`
	TemplateData map[string]any `json:"template_data,omitempty"`
	ReplyTo      string                 `json:"reply_to,omitempty"`
}


type resendErrorResponse struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

func NewResendProvider(config *ResendConfig) *ResendProvider {
	return &ResendProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.resend.com",
	}
}

func (r *ResendProvider) ValidateConfig() error {
	if r.config.APIKey == "" {
		return fmt.Errorf("Resend API key is required")
	}
	if r.config.From == "" {
		return fmt.Errorf("Resend from address is required")
	}
	return nil
}

func (r *ResendProvider) SendEmail(ctx context.Context, req *EmailRequest) error {
	if err := r.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid Resend config: %w", err)
	}

	if len(req.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	from := req.From
	if from == "" {
		from = r.config.From
	}

	payload := resendEmailRequest{
		From:    from,
		To:      req.To,
		CC:      req.CC,
		BCC:     req.BCC,
		Subject: req.Subject,
		Text:    req.Text,
		HTML:    req.HTML,
		ReplyTo: req.ReplyTo,
	}

	return r.sendRequest(ctx, "/emails", payload)
}

func (r *ResendProvider) SendTemplate(ctx context.Context, req *TemplateRequest) error {
	if err := r.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid Resend config: %w", err)
	}

	if len(req.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	if req.TemplateID == "" {
		return fmt.Errorf("template ID is required")
	}

	from := req.From
	if from == "" {
		from = r.config.From
	}

	payload := resendTemplateRequest{
		From:         from,
		To:           req.To,
		CC:           req.CC,
		BCC:          req.BCC,
		Subject:      req.Subject,
		TemplateID:   req.TemplateID,
		TemplateData: req.TemplateData,
		ReplyTo:      req.ReplyTo,
	}

	return r.sendRequest(ctx, "/emails", payload)
}

func (r *ResendProvider) sendRequest(ctx context.Context, endpoint string, payload any) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.config.APIKey)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp resendErrorResponse
		if json.Unmarshal(body, &errorResp) == nil {
			return fmt.Errorf("Resend API error (%d): %s", resp.StatusCode, errorResp.Message)
		}
		return fmt.Errorf("Resend API error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}