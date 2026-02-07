package products

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type ProductResource struct {
	resource.OptimizedBase
}

type ProductFieldResolver struct{}

func NewProductResource() *ProductResource {
	r := &ProductResource{}
	// Model yok - field-based migration kullanılacak
	r.SetSlug("products")
	r.SetTitle("Products")
	r.SetIcon("dashboard")
	r.SetGroup("Products")
	r.SetFieldResolver(&ProductFieldResolver{})
	return r
}

func (r *ProductFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Name", "name").Required().Sortable().Searchable(),
		fields.Textarea("Description", "description").Searchable(),
		fields.RichText("Details", "details"),
		fields.Number("Price", "price").Required(),
		fields.Number("Stock", "stock"),
		fields.DateTime("Created At", "created_at").ReadOnly(),
		fields.DateTime("Updated At", "updated_at").ReadOnly(),
	}
}
