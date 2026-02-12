package products

import (
	"cargo.go/entity"
	"cargo.go/resources" // Registry import
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// init fonksiyonu ile resource'u registry'ye kaydet
func init() {
	resources.Register("products", func() resource.Resource {
		return NewProductResource()
	})
}

type ProductResource struct {
	resource.OptimizedBase
}

type ProductResolveFields struct{}

func NewProductResource() *ProductResource {
	r := &ProductResource{}
	r.SetSlug("products")
	r.SetTitle("Products")
	r.SetIcon("products")
	r.SetGroup("Organization")
	r.SetFieldResolver(&ProductResolveFields{})
	r.SetModel(&entity.Product{})
	r.SetVisible(true)
	return r
}

func (r *ProductResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").
			WithSearchableColumns("name", "email").
			WithEagerLoad(),
		fields.Text("Name", "name"),
		// Registry pattern ile circular dependency çözüldü
		fields.HasMany("ShipmentRows", "shipment_rows", resources.GetOrPanic("shipment-rows")).
			ForeignKey("product_id").
			OwnerKey("id").
			WithEagerLoad(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for a product record
func (r *ProductResource) RecordTitle(record interface{}) string {
	if product, ok := record.(*entity.Product); ok {
		return product.Name
	}
	return ""
}

// With returns relationships to eager load (prevents N+1 queries)
func (r *ProductResource) With() []string {
	return []string{"Organization", "ShipmentRows"}
}
