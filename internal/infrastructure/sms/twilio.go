package sms

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	From       string
}

type TwilioProvider struct {
	config     *TwilioConfig
	httpClient *http.Client
	baseURL    string
}

type twilioResponse struct {
	SID         string `json:"sid"`
	ErrorCode   int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewTwilioProvider(config *TwilioConfig) *TwilioProvider {
	return &TwilioProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.twilio.com/2010-04-01",
	}
}

func (t *TwilioProvider) ValidateConfig() error {
	if t.config.AccountSID == "" {
		return fmt.Errorf("Twilio Account SID is required")
	}
	if t.config.AuthToken == "" {
		return fmt.Errorf("Twilio Auth Token is required")
	}
	if t.config.From == "" {
		return fmt.Errorf("Twilio From number is required")
	}
	return nil
}

func (t *TwilioProvider) SendSMS(ctx context.Context, req *SMSRequest) error {
	if err := t.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid Twilio config: %w", err)
	}

	if req.To == "" {
		return fmt.Errorf("recipient phone number is required")
	}

	if req.Body == "" {
		return fmt.Errorf("SMS body is required")
	}

	from := req.From
	if from == "" {
		from = t.config.From
	}

	data := url.Values{}
	data.Set("To", req.To)
	data.Set("From", from)
	data.Set("Body", req.Body)

	endpoint := fmt.Sprintf("%s/Accounts/%s/Messages.json", t.baseURL, t.config.AccountSID)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	auth := base64.StdEncoding.EncodeToString([]byte(t.config.AccountSID + ":" + t.config.AuthToken))
	httpReq.Header.Set("Authorization", "Basic "+auth)

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	var twilioResp twilioResponse
	if err := json.NewDecoder(resp.Body).Decode(&twilioResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode >= 400 || twilioResp.ErrorCode != 0 {
		return fmt.Errorf("Twilio API error (%d): %s", twilioResp.ErrorCode, twilioResp.ErrorMessage)
	}

	return nil
}