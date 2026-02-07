package fields

import (
	"github.com/ferdiunal/panel.go/pkg/core"
)

// Panel represents a container for grouping fields into sections/cards
type PanelField struct {
	*Schema
	Fields []core.Element
}

// NewPanel creates a new panel/section for grouping fields
func Panel(title string, fields ...core.Element) *PanelField {
	schema := NewField(title)
	schema.View = "panel-field"
	schema.Type = TYPE_PANEL

	return &PanelField{
		Schema: schema,
		Fields: fields,
	}
}

// WithDescription adds a description to the panel
func (p *PanelField) WithDescription(description string) *PanelField {
	p.Props["description"] = description
	return p
}

// WithColumns sets the grid columns for the panel (1-4)
func (p *PanelField) WithColumns(columns int) *PanelField {
	if columns < 1 {
		columns = 1
	}
	if columns > 4 {
		columns = 4
	}
	p.Props["columns"] = columns
	return p
}

// Collapsible makes the panel collapsible
func (p *PanelField) Collapsible() *PanelField {
	p.Props["collapsible"] = true
	return p
}

// DefaultCollapsed sets the panel to be collapsed by default
func (p *PanelField) DefaultCollapsed() *PanelField {
	p.Props["collapsible"] = true
	p.Props["defaultCollapsed"] = true
	return p
}

// GetFields returns the fields in this panel
func (p *PanelField) GetFields() []core.Element {
	return p.Fields
}
