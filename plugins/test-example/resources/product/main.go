package product

import (
	"test-example/entity"
	"test-example/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("products", func() resource.Resource {
		return NewProductResource()
	})
}

type ProductResource struct {
	resource.OptimizedBase
}

type ProductPolicy struct{}

func (p *ProductPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *ProductPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *ProductPolicy) Create(ctx *context.Context) bool    { return true }
func (p *ProductPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *ProductPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type ProductResolveFields struct{}

func NewProductResource() *ProductResource {
	r := &ProductResource{}
	r.SetSlug("products")
	r.SetTitle("Products")
	r.SetIcon("package")
	r.SetGroup("Management")
	r.SetModel(&entity.Product{})
	r.SetFieldResolver(&ProductResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&ProductPolicy{})
	return r
}

func (r *ProductResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("Name", "name").Label("Ürün Adı").Placeholder("Ürün Adı").Required(),
		fields.BelongsToMany("Categories", "categories", resources.GetOrPanic("categories")).
			PivotTable("product_categories").WithEagerLoad().Label("Kategoriler").HideOnList(),
		fields.HasMany("ShipmentRows", "shipment_rows", resources.GetOrPanic("shipment-rows")).
			ForeignKey("product_id").OwnerKey("id").WithEagerLoad().Label("Gönderi Satırları").HideOnList().HideOnCreate(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *ProductResource) RecordTitle(record interface{}) string {
	if product, ok := record.(*entity.Product); ok {
		return product.Name
	}
	return ""
}

func (r *ProductResource) With() []string {
	return []string{"Organization", "Categories", "ShipmentRows"}
}
