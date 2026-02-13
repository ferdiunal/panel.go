package organization

import (
	"cargo.go/entity"
	"cargo.go/resources" // Registry import
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// init fonksiyonu ile resource'u registry'ye kaydet
func init() {
	resources.Register("organizations", func() resource.Resource {
		return NewOrganizationResource()
	})
}

type OrganizationResource struct {
	resource.OptimizedBase
}

type OrganizationPolicy struct{}

func (p *OrganizationPolicy) ViewAny(ctx *context.Context) bool {
	return true
}

func (p *OrganizationPolicy) View(ctx *context.Context, model interface{}) bool {
	return true
}

func (p *OrganizationPolicy) Create(ctx *context.Context) bool {
	return true
}

func (p *OrganizationPolicy) Update(ctx *context.Context, model interface{}) bool {
	return true
}

func (p *OrganizationPolicy) Delete(ctx *context.Context, model interface{}) bool {
	return true
}

type OrganizationResolveFields struct{}

func NewOrganizationResource() *OrganizationResource {
	r := &OrganizationResource{}
	r.SetSlug("organizations")
	r.SetTitle("Organizations")
	r.SetIcon("organizations")
	r.SetGroup("Operations")
	r.SetModel(&entity.Organization{})
	r.SetFieldResolver(&OrganizationResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&OrganizationPolicy{})
	return r
}

func (r *OrganizationResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name").
			Label("Şirket Adı").
			Placeholder("Şirket Adı").
			Required(),
		fields.Email("Email", "email").
			Label("Şirket E-posta").
			Placeholder("Şirket E-posta").
			Required(),
		fields.Tel("Phone", "phone").
			Label("Şirket Telefon").
			Placeholder("Şirket Telefon").
			Required(),
		// Registry pattern ile circular dependency çözüldü
		fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses")).
			ForeignKey("organization_id").
			OwnerKey("id").
			WithEagerLoad().
			Label("Fatura Adresi").
			HideOnList().
			HideOnCreate(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for an organization record
func (r *OrganizationResource) RecordTitle(record interface{}) string {
	if org, ok := record.(*entity.Organization); ok {
		return org.Name
	}
	return ""
}

// With returns relationships to eager load (prevents N+1 queries)
func (r *OrganizationResource) With() []string {
	return []string{"Addresses"}
}
