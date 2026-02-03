package widget

import (
	"github.com/ferdiunal/panel.go/internal/context"
	"gorm.io/gorm"
)

// Card matches Laravel Nova's Card concept.
// Metrics (Value, Trend) are just specific types of Cards.
type Card interface {
	Name() string
	Component() string // Frontend component name (e.g. "value-metric", "custom-card")
	Width() string     // "1/3", "1/2", "full", etc.
	Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
	JsonSerialize() map[string]interface{}
}

// BaseCard provides common fields for all cards
type BaseCard struct {
	TitleStr     string
	ComponentStr string
	WidthStr     string
}

func (c *BaseCard) Name() string {
	return c.TitleStr
}

func (c *BaseCard) Component() string {
	if c.ComponentStr == "" {
		return "card" // Default component
	}
	return c.ComponentStr
}

func (c *BaseCard) Width() string {
	if c.WidthStr == "" {
		return "1/3"
	}
	return c.WidthStr
}

// CustomCard allows creating arbitrary cards from outside
type CustomCard struct {
	BaseCard
	Content     interface{} // Static content or initial data
	ResolveFunc func(ctx *context.Context, db *gorm.DB) (interface{}, error)
}

func (c *CustomCard) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if c.ResolveFunc != nil {
		return c.ResolveFunc(ctx, db)
	}
	return c.Content, nil
}

func (c *CustomCard) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": c.Component(),
		"title":     c.Name(),
		"width":     c.Width(),
		"content":   c.Content,
	}
}

// NewCard creates a simple custom card
func NewCard(title, component string) *CustomCard {
	return &CustomCard{
		BaseCard: BaseCard{
			TitleStr:     title,
			ComponentStr: component,
			WidthStr:     "1/3",
		},
	}
}

// Fluent setters for CustomCard
func (c *CustomCard) SetWidth(w string) *CustomCard {
	c.WidthStr = w
	return c
}

func (c *CustomCard) SetContent(content interface{}) *CustomCard {
	c.Content = content
	return c
}
