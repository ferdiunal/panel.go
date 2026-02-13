package address

import (
	"test-example/entity"
	"test-example/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("addresses", func() resource.Resource {
		return NewAddressResource()
	})
}

type AddressResource struct {
	resource.OptimizedBase
}

type AddressPolicy struct{}

func (p *AddressPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *AddressPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *AddressPolicy) Create(ctx *context.Context) bool    { return true }
func (p *AddressPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *AddressPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type AddressResolveFields struct{}

func NewAddressResource() *AddressResource {
	r := &AddressResource{}
	r.SetSlug("addresses")
	r.SetTitle("Addresses")
	r.SetIcon("map-pin")
	r.SetGroup("Management")
	r.SetModel(&entity.Address{})
	r.SetFieldResolver(&AddressResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&AddressPolicy{})
	return r
}

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("Name", "name").Label("Adres Adı").Placeholder("Adres Adı").Required(),
		fields.Textarea("Address", "address").Label("Adres").Placeholder("Adres").Required(),
		fields.Text("City", "city").Label("Şehir").Placeholder("Şehir").Required(),
		fields.Text("Country", "country").Label("Ülke").Placeholder("Ülke").Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *AddressResource) RecordTitle(record interface{}) string {
	if addr, ok := record.(*entity.Address); ok {
		return addr.Name
	}
	return ""
}

func (r *AddressResource) With() []string {
	return []string{"Organization"}
}
