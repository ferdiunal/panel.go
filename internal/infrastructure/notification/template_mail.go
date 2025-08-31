package notification

import (
	"context"
	"fmt"
)

type TemplateMailMessage struct {
	*MailMessage
	TemplateName string
	TemplateData TemplateData
	renderer     TemplateRenderer
}

func NewTemplateMailMessage(renderer TemplateRenderer) *TemplateMailMessage {
	return &TemplateMailMessage{
		MailMessage:  NewMailMessage(),
		TemplateData: make(TemplateData),
		renderer:     renderer,
	}
}

func (tmm *TemplateMailMessage) SetTemplate(templateName string) *TemplateMailMessage {
	tmm.TemplateName = templateName
	return tmm
}

func (tmm *TemplateMailMessage) SetData(key string, value any) *TemplateMailMessage {
	tmm.TemplateData[key] = value
	return tmm
}

func (tmm *TemplateMailMessage) SetAllData(data TemplateData) *TemplateMailMessage {
	tmm.TemplateData = data
	return tmm
}

func (tmm *TemplateMailMessage) RenderContent(ctx context.Context) error {
	if tmm.TemplateName == "" {
		return fmt.Errorf("template name is required")
	}

	if tmm.renderer == nil {
		return fmt.Errorf("template renderer is required")
	}

	htmlContent, err := tmm.renderer.RenderHTML(ctx, tmm.TemplateName, tmm.TemplateData)
	if err != nil {
		return fmt.Errorf("failed to render HTML template: %w", err)
	}
	tmm.HTMLBody = htmlContent

	textContent, err := tmm.renderer.RenderText(ctx, tmm.TemplateName, tmm.TemplateData)
	if err != nil {
		tmm.Body = ""
	} else {
		tmm.Body = textContent
	}

	return nil
}

func (tmm *TemplateMailMessage) GetTemplateData() TemplateData {
	return tmm.TemplateData
}

type TemplateNotification struct {
	BaseNotification
	templateRenderer TemplateRenderer
	templateName     string
	templateData     TemplateData
	channels         []string
	subject          string
}

func NewTemplateNotification(renderer TemplateRenderer, templateName string) *TemplateNotification {
	return &TemplateNotification{
		templateRenderer: renderer,
		templateName:     templateName,
		templateData:     make(TemplateData),
		channels:         []string{"mail"},
	}
}

func (tn *TemplateNotification) SetData(key string, value any) *TemplateNotification {
	tn.templateData[key] = value
	return tn
}

func (tn *TemplateNotification) SetAllData(data TemplateData) *TemplateNotification {
	tn.templateData = data
	return tn
}

func (tn *TemplateNotification) SetSubject(subject string) *TemplateNotification {
	tn.subject = subject
	return tn
}

func (tn *TemplateNotification) SetChannels(channels []string) *TemplateNotification {
	tn.channels = channels
	return tn
}

func (tn *TemplateNotification) Via(notifiable Notifiable) []string {
	return tn.channels
}

func (tn *TemplateNotification) ToMail(notifiable Notifiable) *MailMessage {
	templateMail := NewTemplateMailMessage(tn.templateRenderer)
	templateMail.SetTemplate(tn.templateName)
	templateMail.SetAllData(tn.templateData)

	if tn.subject != "" {
		templateMail.SetSubject(tn.subject)
	}

	if err := templateMail.RenderContent(context.Background()); err != nil {
		return nil
	}

	return templateMail.MailMessage
}

func (tn *TemplateNotification) GetTemplateName() string {
	return tn.templateName
}

func (tn *TemplateNotification) GetTemplateData() TemplateData {
	return tn.templateData
}