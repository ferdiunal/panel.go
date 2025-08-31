package notification

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

type mockTemplComponent struct {
	content string
	err     error
}

func (m *mockTemplComponent) Render(ctx context.Context, w io.Writer) error {
	if m.err != nil {
		return m.err
	}
	_, err := w.Write([]byte(m.content))
	return err
}

func TestNewTemplRenderer(t *testing.T) {
	renderer := NewTemplRenderer()
	
	if renderer == nil {
		t.Fatal("expected renderer to be created")
	}
	
	if renderer.templates == nil {
		t.Error("expected templates map to be initialized")
	}
}

func TestTemplRenderer_RegisterTemplate(t *testing.T) {
	renderer := NewTemplRenderer()
	
	templateFunc := func(data TemplateData) TemplComponent {
		return &mockTemplComponent{content: "test content"}
	}
	
	renderer.RegisterTemplate("test", templateFunc)
	
	if len(renderer.templates) != 1 {
		t.Errorf("expected 1 template, got %d", len(renderer.templates))
	}
	
	if !renderer.HasTemplate("test") {
		t.Error("expected template to be registered")
	}
}

func TestTemplRenderer_RenderHTML(t *testing.T) {
	tests := []struct {
		name           string
		templateName   string
		registerTemplate bool
		content        string
		renderError    error
		expectedError  string
		expectedContent string
	}{
		{
			name:            "successful render",
			templateName:    "test",
			registerTemplate: true,
			content:         "<h1>Hello World</h1>",
			expectedContent: "<h1>Hello World</h1>",
		},
		{
			name:          "template not found",
			templateName:  "nonexistent",
			expectedError: "template 'nonexistent' not found",
		},
		{
			name:            "render error",
			templateName:    "error",
			registerTemplate: true,
			renderError:     errors.New("render failed"),
			expectedError:   "failed to render template 'error'",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewTemplRenderer()
			
			if tt.registerTemplate {
				renderer.RegisterTemplate(tt.templateName, func(data TemplateData) TemplComponent {
					return &mockTemplComponent{
						content: tt.content,
						err:     tt.renderError,
					}
				})
			}
			
			result, err := renderer.RenderHTML(context.Background(), tt.templateName, TemplateData{})
			
			if tt.expectedError != "" {
				if err == nil {
					t.Error("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expectedContent {
					t.Errorf("expected content '%s', got '%s'", tt.expectedContent, result)
				}
			}
		})
	}
}

func TestTemplRenderer_RenderText(t *testing.T) {
	renderer := NewTemplRenderer()
	
	renderer.RegisterTemplate("test_text", func(data TemplateData) TemplComponent {
		return &mockTemplComponent{content: "Hello World"}
	})
	
	result, err := renderer.RenderText(context.Background(), "test", TemplateData{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if result != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", result)
	}
}

func TestNewTemplateMailMessage(t *testing.T) {
	renderer := NewTemplRenderer()
	tmm := NewTemplateMailMessage(renderer)
	
	if tmm == nil {
		t.Fatal("expected template mail message to be created")
	}
	
	if tmm.MailMessage == nil {
		t.Error("expected mail message to be initialized")
	}
	
	if tmm.renderer != renderer {
		t.Error("expected renderer to be set correctly")
	}
}

func TestTemplateMailMessage_SetTemplate(t *testing.T) {
	renderer := NewTemplRenderer()
	tmm := NewTemplateMailMessage(renderer)
	
	result := tmm.SetTemplate("test")
	
	if result != tmm {
		t.Error("expected fluent interface to return self")
	}
	
	if tmm.TemplateName != "test" {
		t.Errorf("expected template name 'test', got '%s'", tmm.TemplateName)
	}
}

func TestTemplateMailMessage_SetData(t *testing.T) {
	renderer := NewTemplRenderer()
	tmm := NewTemplateMailMessage(renderer)
	
	result := tmm.SetData("key", "value")
	
	if result != tmm {
		t.Error("expected fluent interface to return self")
	}
	
	if tmm.TemplateData["key"] != "value" {
		t.Error("expected template data to be set correctly")
	}
}

func TestTemplateMailMessage_RenderContent(t *testing.T) {
	tests := []struct {
		name           string
		templateName   string
		setupRenderer  bool
		expectedError  string
	}{
		{
			name:          "missing template name",
			expectedError: "template name is required",
		},
		{
			name:         "missing renderer",
			templateName: "test",
			setupRenderer: false,
			expectedError: "template renderer is required",
		},
		{
			name:          "successful render",
			templateName:  "test",
			setupRenderer: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var renderer *TemplRenderer
			if tt.setupRenderer {
				renderer = NewTemplRenderer()
				renderer.RegisterTemplate("test", func(data TemplateData) TemplComponent {
					return &mockTemplComponent{content: "<h1>Test</h1>"}
				})
				renderer.RegisterTemplate("test_text", func(data TemplateData) TemplComponent {
					return &mockTemplComponent{content: "Test"}
				})
			}
			
			var tmm *TemplateMailMessage
			if tt.setupRenderer {
				tmm = NewTemplateMailMessage(renderer)
			} else {
				tmm = &TemplateMailMessage{
					MailMessage:  NewMailMessage(),
					TemplateData: make(TemplateData),
				}
			}
			tmm.TemplateName = tt.templateName
			
			err := tmm.RenderContent(context.Background())
			
			if tt.expectedError != "" {
				if err == nil {
					t.Error("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tmm.HTMLBody != "<h1>Test</h1>" {
					t.Errorf("expected HTML body '<h1>Test</h1>', got '%s'", tmm.HTMLBody)
				}
			}
		})
	}
}

func TestNewTemplateNotification(t *testing.T) {
	renderer := NewTemplRenderer()
	tn := NewTemplateNotification(renderer, "test")
	
	if tn == nil {
		t.Fatal("expected template notification to be created")
	}
	
	if tn.templateRenderer != renderer {
		t.Error("expected renderer to be set correctly")
	}
	
	if tn.templateName != "test" {
		t.Errorf("expected template name 'test', got '%s'", tn.templateName)
	}
}

func TestTemplateNotification_FluentInterface(t *testing.T) {
	renderer := NewTemplRenderer()
	tn := NewTemplateNotification(renderer, "test")
	
	result := tn.SetData("key", "value").
		SetSubject("Test Subject").
		SetChannels([]string{"mail", "sms"})
	
	if result != tn {
		t.Error("expected fluent interface to return self")
	}
	
	if tn.templateData["key"] != "value" {
		t.Error("expected template data to be set")
	}
	
	if tn.subject != "Test Subject" {
		t.Error("expected subject to be set")
	}
	
	if len(tn.channels) != 2 || tn.channels[0] != "mail" || tn.channels[1] != "sms" {
		t.Error("expected channels to be set correctly")
	}
}