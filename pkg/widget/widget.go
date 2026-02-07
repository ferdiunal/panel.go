package widget

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"gorm.io/gorm"
)

// CardType represents the type of card being displayed
type CardType string

const (
	CardTypeValue     CardType = "value"
	CardTypeTrend     CardType = "trend"
	CardTypeTable     CardType = "table"
	CardTypePartition CardType = "partition"
	CardTypeProgress  CardType = "progress"
)

// Card matches Laravel Nova's Card concept.
// Metrics (Value, Trend) are just specific types of Cards.
type Card interface {
	Name() string
	Component() string // Frontend component name (e.g. "value-metric", "custom-card")
	Width() string     // "1/3", "1/2", "full", etc.
	GetType() CardType // Returns the type of card
	Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
	HandleError(err error) map[string]interface{} // Handles errors and returns error response
	GetMetadata() map[string]interface{}          // Returns card metadata
	JsonSerialize() map[string]interface{}
}

// CardResolver defines the interface for resolving card data
type CardResolver interface {
	Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
}

// BaseCard provides common fields for all cards
type BaseCard struct {
	TitleStr     string
	ComponentStr string
	WidthStr     string
	CardTypeVal  CardType
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

func (c *BaseCard) GetType() CardType {
	if c.CardTypeVal == "" {
		return CardTypeValue // Default type
	}
	return c.CardTypeVal
}

func (c *BaseCard) HandleError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
		"title": c.TitleStr,
	}
}

func (c *BaseCard) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":      c.TitleStr,
		"component": c.Component(),
		"width":     c.Width(),
		"type":      c.GetType(),
	}
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
		"type":      c.GetType(),
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
			CardTypeVal:  CardTypeValue,
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
