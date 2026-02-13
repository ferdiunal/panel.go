package shipment

import (
	"cargo.go/entity"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type ShipmentRowResource struct {
	resource.OptimizedBase
}

type ShipmentRowResolveFields struct{}

func NewShipmentRowResource() *ShipmentRowResource {
	r := &ShipmentRowResource{}
	r.SetSlug("shipment-rows")
	r.SetTitle("Shipment Rows")
	r.SetIcon("shipment_rows")
	r.SetGroup("Organization")
	r.SetFieldResolver(&ShipmentRowResolveFields{})
	r.SetModel(&entity.ShipmentRow{})
	r.SetVisible(true)
	return r
}

func (r *ShipmentRowResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Shipment", "shipment_id", "shipments"),
		fields.BelongsTo("Product", "product_id", "products"),
		fields.Text("Quantity", "quantity"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}
