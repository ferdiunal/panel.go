package fields

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
)

// Panel represents a container for grouping fields into sections/cards
type Panel struct {
	*Schema
	Fields []core.Element
}

// NewPanel creates a new panel/section for grouping fields
func NewPanel(title string, fields ...core.Element) *Panel {
	schema := NewField(title)
	schema.View = "panel-field"
	schema.Type = TYPE_PANEL

	return &Panel{
		Schema: schema,
		Fields: fields,
	}
}

// WithDescription adds a description to the panel
func (p *Panel) WithDescription(description string) *Panel {
	p.Props["description"] = description
	return p
}

// WithColumns sets the grid columns for the panel (1-4)
func (p *Panel) WithColumns(columns int) *Panel {
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
func (p *Panel) Collapsible() *Panel {
	p.Props["collapsible"] = true
	return p
}

// DefaultCollapsed sets the panel to be collapsed by default
func (p *Panel) DefaultCollapsed() *Panel {
	p.Props["collapsible"] = true
	p.Props["defaultCollapsed"] = true
	return p
}

// GetFields returns the fields in this panel
func (p *Panel) GetFields() []core.Element {
	return p.Fields
}
