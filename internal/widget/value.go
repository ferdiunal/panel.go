package widget

import (
	"github.com/ferdiunal/panel.go/internal/context"
	"gorm.io/gorm"
)

type Value struct {
	Title     string
	QueryFunc func(ctx *context.Context, db *gorm.DB) (int64, error)
}

func (w *Value) Name() string {
	return w.Title
}

func (w *Value) Component() string {
	return "value-metric"
}

func (w *Value) Width() string {
	return "1/3"
}

func (w *Value) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	val, err := w.QueryFunc(ctx, db)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"value": val,
		"title": w.Title,
	}, nil
}

func (w *Value) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": "value-metric",
		"title":     w.Title,
		"width":     "1/3",
	}
}

// Helpers

func NewCountWidget(title string, model interface{}) *Value {
	return &Value{
		Title: title,
		QueryFunc: func(ctx *context.Context, db *gorm.DB) (int64, error) {
			var total int64
			if err := db.Model(model).Count(&total).Error; err != nil {
				return 0, err
			}
			return total, nil
		},
	}
}
