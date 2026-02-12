package address

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
	resources.Register("addresses", func() resource.Resource {
		return NewAddressResource()
	})
}

type AddressResource struct {
	resource.OptimizedBase
}
type AddressResolveFields struct{}

func NewAddressResource() *AddressResource {
	r := &AddressResource{}
	r.SetSlug("addresses")
	r.SetTitle("Addresses")
	r.SetIcon("addresses")
	r.SetGroup("Organization")
	r.SetFieldResolver(&AddressResolveFields{})
	r.SetModel(&entity.Address{})
	r.SetVisible(true)
	r.SetPolicy(&AddressPolicy{})
	r.SetRecordTitleKey("name")
	return r
}

type AddressPolicy struct {
}

func (r *AddressPolicy) ViewAny(ctx *context.Context) bool {
	if ctx.HasRole("organization") {
		return true
	}
	return false
}

func (r *AddressPolicy) View(ctx *context.Context, model interface{}) bool {
	return true
}

func (r *AddressPolicy) Create(ctx *context.Context) bool {
	return true
}

func (r *AddressPolicy) Update(ctx *context.Context, model interface{}) bool {
	return true
}

func (r *AddressPolicy) Delete(ctx *context.Context, model interface{}) bool {
	return true
}

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").
			WithSearchableColumns("name", "email", "phone").
			AutoOptions("name").
			WithEagerLoad().
			Label("Organizasyon"),
		fields.Text("Name", "name"),
		fields.Select("Type", "type").Options([]map[string]interface{}{
			{"value": entity.AddressTypeSender, "label": "Sender"},
			{"value": entity.AddressTypeReceiver, "label": "Receiver"},
			{"value": entity.AddressTypeInvoice, "label": "Invoice"},
		}),
		fields.Textarea("Address", "address"),
		fields.Text("City", "city"),
		fields.Text("State", "state"),
		fields.Text("ZipCode", "zip_code"),
		fields.Text("Country", "country"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for an address record
func (r *AddressResource) RecordTitle(record interface{}) string {
	if addr, ok := record.(*entity.Address); ok {
		return addr.Name
	}
	return ""
}

// With returns relationships to eager load (prevents N+1 queries)
func (r *AddressResource) With() []string {
	return []string{"Organization"}
}
