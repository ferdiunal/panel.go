package shipment_row

import (
	"fmt"

	"test-example/entity"
	"test-example/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("shipment-rows", func() resource.Resource {
		return NewShipmentRowResource()
	})
}

type ShipmentRowResource struct {
	resource.OptimizedBase
}

type ShipmentRowPolicy struct{}

func (p *ShipmentRowPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *ShipmentRowPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentRowPolicy) Create(ctx *context.Context) bool    { return true }
func (p *ShipmentRowPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentRowPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type ShipmentRowResolveFields struct{}

func NewShipmentRowResource() *ShipmentRowResource {
	r := &ShipmentRowResource{}
	r.SetSlug("shipment-rows")
	r.SetTitle("Shipment Rows")
	r.SetIcon("list")
	r.SetGroup("Operations")
	r.SetModel(&entity.ShipmentRow{})
	r.SetFieldResolver(&ShipmentRowResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&ShipmentRowPolicy{})
	return r
}

func (r *ShipmentRowResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Shipment", "shipment_id", resources.GetOrPanic("shipments")).
			DisplayUsing("name").WithSearchableColumns("name").WithEagerLoad().Label("Gönderi").Required(),
		fields.BelongsTo("Product", "product_id", resources.GetOrPanic("products")).
			DisplayUsing("name").WithSearchableColumns("name").WithEagerLoad().Label("Ürün").Required(),
		fields.Number("Quantity", "quantity").Label("Miktar").Placeholder("Miktar").Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *ShipmentRowResource) RecordTitle(record interface{}) string {
	if row, ok := record.(*entity.ShipmentRow); ok {
		return fmt.Sprintf("Row #%!d(MISSING)", row.ID)
	}
	return ""
}

func (r *ShipmentRowResource) With() []string {
	return []string{"Shipment", "Product"}
}
