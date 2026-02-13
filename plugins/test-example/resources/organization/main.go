package organization

import (
	"test-example/entity"
	"test-example/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("organizations", func() resource.Resource {
		return NewOrganizationResource()
	})
}

type OrganizationResource struct {
	resource.OptimizedBase
}

type OrganizationPolicy struct{}

func (p *OrganizationPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *OrganizationPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *OrganizationPolicy) Create(ctx *context.Context) bool    { return true }
func (p *OrganizationPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *OrganizationPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type OrganizationResolveFields struct{}

func NewOrganizationResource() *OrganizationResource {
	r := &OrganizationResource{}
	r.SetSlug("organizations")
	r.SetTitle("Organizations")
	r.SetIcon("building")
	r.SetGroup("Management")
	r.SetModel(&entity.Organization{})
	r.SetFieldResolver(&OrganizationResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&OrganizationPolicy{})
	return r
}

func (r *OrganizationResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name").Label("Organizasyon Adı").Placeholder("Organizasyon Adı").Required(),
		fields.Email("Email", "email").Label("E-posta").Placeholder("E-posta").Required(),
		fields.Tel("Phone", "phone").Label("Telefon").Placeholder("Telefon").Required(),
		fields.HasOne("BillingInfo", "billing_info", resources.GetOrPanic("billing-info")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Fatura Bilgisi").HideOnList().HideOnCreate(),
		fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Adresler").HideOnList().HideOnCreate(),
		fields.HasMany("Products", "products", resources.GetOrPanic("products")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Ürünler").HideOnList().HideOnCreate(),
		fields.HasMany("Shipments", "shipments", resources.GetOrPanic("shipments")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Gönderiler").HideOnList().HideOnCreate(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *OrganizationResource) RecordTitle(record interface{}) string {
	if org, ok := record.(*entity.Organization); ok {
		return org.Name
	}
	return ""
}

func (r *OrganizationResource) With() []string {
	return []string{"BillingInfo", "Addresses", "Products", "Shipments"}
}
