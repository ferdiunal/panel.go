package shipment

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
	resources.Register("shipments", func() resource.Resource {
		return NewShipmentResource()
	})
}

type ShipmentResource struct {
	resource.OptimizedBase
}
type ShipmentResolveFields struct{}

func NewShipmentResource() *ShipmentResource {
	r := &ShipmentResource{}
	r.SetSlug("shipments")
	r.SetTitle("Shipments")
	r.SetIcon("shipments")
	r.SetGroup("Organization")
	r.SetFieldResolver(&ShipmentResolveFields{})
	r.SetModel(&entity.Shipment{})
	r.SetVisible(true)
	return r
}

func (r *ShipmentResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").
			WithSearchableColumns("name", "email").
			WithEagerLoad(),
		fields.BelongsTo("SenderAddress", "sender_address_id", resources.GetOrPanic("addresses")).
			DisplayUsing("name").
			WithSearchableColumns("name", "city", "country").
			WithEagerLoad(),
		fields.BelongsTo("ReceiverAddress", "receiver_address_id", resources.GetOrPanic("addresses")).
			DisplayUsing("name").
			WithSearchableColumns("name", "city", "country").
			WithEagerLoad(),
		// Registry pattern ile circular dependency çözüldü
		fields.HasMany("ShipmentRows", "shipment_rows", resources.GetOrPanic("shipment-rows")).
			ForeignKey("shipment_id").
			OwnerKey("id").
			WithEagerLoad(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for a shipment record
func (r *ShipmentResource) RecordTitle(record interface{}) string {
	if shipment, ok := record.(*entity.Shipment); ok {
		return shipment.Name
	}
	return ""
}

// With returns relationships to eager load (prevents N+1 queries)
func (r *ShipmentResource) With() []string {
	return []string{"Organization", "SenderAddress", "ReceiverAddress", "ShipmentRows"}
}
