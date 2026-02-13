package shipment_row

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
	resources.Register("shipment-rows", func() resource.Resource {
		return NewShipmentRowResource()
	})
}

// ShipmentRowResource - Gönderi satırı resource'u
// BelongsTo: Shipment (bir gönderi satırı bir gönderiye aittir)
// BelongsTo: Product (bir gönderi satırı bir ürüne aittir)
type ShipmentRowResource struct {
	resource.OptimizedBase
}

type ShipmentRowResolveFields struct{}

// NewShipmentRowResource - Yeni ShipmentRow resource'u oluşturur
func NewShipmentRowResource() *ShipmentRowResource {
	r := &ShipmentRowResource{}
	r.SetSlug("shipment-rows")
	r.SetTitle("Shipment Rows")
	r.SetIcon("shipment_rows")
	r.SetGroup("Operations")
	r.SetFieldResolver(&ShipmentRowResolveFields{})
	r.SetModel(&entity.ShipmentRow{})
	r.SetVisible(true)
	return r
}

// ResolveFields - ShipmentRow için field'ları tanımlar
func (r *ShipmentRowResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		// Registry pattern ile circular dependency çözüldü
		fields.BelongsTo("Shipment", "shipment_id", resources.GetOrPanic("shipments")).
			DisplayUsing("name").
			WithSearchableColumns("name").
			WithEagerLoad(),
		fields.BelongsTo("Product", "product_id", resources.GetOrPanic("products")).
			DisplayUsing("name").
			WithSearchableColumns("name").
			WithEagerLoad(),
		fields.Number("Quantity", "quantity").
			Label("Miktar").
			Placeholder("Miktar").
			Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for a shipment row record
func (r *ShipmentRowResource) RecordTitle(record interface{}) string {
	if row, ok := record.(*entity.ShipmentRow); ok {
		if row.Product != nil {
			return row.Product.Name
		}
		if row.Shipment != nil {
			return row.Shipment.Name
		}
	}
	return ""
}

// With - Eager load edilecek ilişkileri döndürür (N+1 query problemini önler)
func (r *ShipmentRowResource) With() []string {
	return []string{"Shipment", "Product"}
}
