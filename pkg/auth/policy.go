package auth

import "github.com/ferdiunal/panel.go/pkg/context"

type Policy interface {
	ViewAny(ctx *context.Context) bool
	View(ctx *context.Context, model interface{}) bool
	Create(ctx *context.Context) bool
	Update(ctx *context.Context, model interface{}) bool
	Delete(ctx *context.Context, model interface{}) bool
}
