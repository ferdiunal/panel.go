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
	// Model ile migration - GORM AutoMigrate kullanılacak
	r.SetModel(&Product{})
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

		// Basic Information Panel
		fields.Panel("Basic Information",
			fields.Text("Name", "name").Required().Sortable().Searchable(),
			fields.Number("Price", "price").Required(),
			fields.Number("Stock", "stock"),
		).WithDescription("Product basic details").WithColumns(2),

		// Description Panel
		fields.Panel("Description",
			fields.Textarea("Short Description", "description").Searchable(),
			fields.RichText("Full Details", "details"),
		).WithDescription("Product descriptions and details").Collapsible(),

		// Metadata Panel
		fields.Panel("Metadata",
			fields.DateTime("Created At", "created_at").ReadOnly(),
			fields.DateTime("Updated At", "updated_at").ReadOnly(),
		).WithDescription("System information").WithColumns(2).DefaultCollapsed(),
	}
}
