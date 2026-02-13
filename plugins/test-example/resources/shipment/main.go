package shipment

import (
	"test-example/entity"
	"test-example/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("shipments", func() resource.Resource {
		return NewShipmentResource()
	})
}

type ShipmentResource struct {
	resource.OptimizedBase
}

type ShipmentPolicy struct{}

func (p *ShipmentPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *ShipmentPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentPolicy) Create(ctx *context.Context) bool    { return true }
func (p *ShipmentPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type ShipmentResolveFields struct{}

func NewShipmentResource() *ShipmentResource {
	r := &ShipmentResource{}
	r.SetSlug("shipments")
	r.SetTitle("Shipments")
	r.SetIcon("truck")
	r.SetGroup("Operations")
	r.SetModel(&entity.Shipment{})
	r.SetFieldResolver(&ShipmentResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&ShipmentPolicy{})
	return r
}

func (r *ShipmentResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("Name", "name").Label("Gönderi Adı").Placeholder("Gönderi Adı").Required(),
		fields.HasMany("ShipmentRows", "shipment_rows", resources.GetOrPanic("shipment-rows")).
			ForeignKey("shipment_id").OwnerKey("id").WithEagerLoad().Label("Gönderi Satırları").HideOnList().HideOnCreate(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *ShipmentResource) RecordTitle(record interface{}) string {
	if shipment, ok := record.(*entity.Shipment); ok {
		return shipment.Name
	}
	return ""
}

func (r *ShipmentResource) With() []string {
	return []string{"Organization", "ShipmentRows"}
}
