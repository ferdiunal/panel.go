package notification

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

type TemplateData map[string]any

type TemplateRenderer interface {
	RenderHTML(ctx context.Context, templateName string, data TemplateData) (string, error)
	RenderText(ctx context.Context, templateName string, data TemplateData) (string, error)
}

type TemplComponent interface {
	Render(ctx context.Context, w io.Writer) error
}

type TemplRenderer struct {
	templates map[string]func(TemplateData) TemplComponent
}

func NewTemplRenderer() *TemplRenderer {
	return &TemplRenderer{
		templates: make(map[string]func(TemplateData) TemplComponent),
	}
}

func (tr *TemplRenderer) RegisterTemplate(name string, templateFunc func(TemplateData) TemplComponent) {
	tr.templates[name] = templateFunc
}

func (tr *TemplRenderer) RenderHTML(ctx context.Context, templateName string, data TemplateData) (string, error) {
	templateFunc, exists := tr.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template '%s' not found", templateName)
	}

	component := templateFunc(data)
	if component == nil {
		return "", fmt.Errorf("template function returned nil component")
	}

	var buf bytes.Buffer
	if err := component.Render(ctx, &buf); err != nil {
		return "", fmt.Errorf("failed to render template '%s': %w", templateName, err)
	}

	return buf.String(), nil
}

func (tr *TemplRenderer) RenderText(ctx context.Context, templateName string, data TemplateData) (string, error) {
	textTemplateName := templateName + "_text"
	templateFunc, exists := tr.templates[textTemplateName]
	if !exists {
		return "", fmt.Errorf("text template '%s' not found", textTemplateName)
	}

	component := templateFunc(data)
	if component == nil {
		return "", fmt.Errorf("text template function returned nil component")
	}

	var buf bytes.Buffer
	if err := component.Render(ctx, &buf); err != nil {
		return "", fmt.Errorf("failed to render text template '%s': %w", textTemplateName, err)
	}

	return buf.String(), nil
}

func (tr *TemplRenderer) HasTemplate(templateName string) bool {
	_, exists := tr.templates[templateName]
	return exists
}

func (tr *TemplRenderer) HasTextTemplate(templateName string) bool {
	_, exists := tr.templates[templateName+"_text"]
	return exists
}

func (tr *TemplRenderer) GetAvailableTemplates() []string {
	templates := make([]string, 0, len(tr.templates))
	for name := range tr.templates {
		if !contains(name, "_text") {
			templates = append(templates, name)
		}
	}
	return templates
}

func contains(slice string, item string) bool {
	return len(slice) >= len(item) && slice[len(slice)-len(item):] == item
}