package billing_info

import (
	"test-example/entity"
	"test-example/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("billing-info", func() resource.Resource {
		return NewBillingInfoResource()
	})
}

type BillingInfoResource struct {
	resource.OptimizedBase
}

type BillingInfoPolicy struct{}

func (p *BillingInfoPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *BillingInfoPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *BillingInfoPolicy) Create(ctx *context.Context) bool    { return true }
func (p *BillingInfoPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *BillingInfoPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type BillingInfoResolveFields struct{}

func NewBillingInfoResource() *BillingInfoResource {
	r := &BillingInfoResource{}
	r.SetSlug("billing-info")
	r.SetTitle("Billing Info")
	r.SetIcon("receipt")
	r.SetGroup("Management")
	r.SetModel(&entity.BillingInfo{})
	r.SetFieldResolver(&BillingInfoResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&BillingInfoPolicy{})
	return r
}

func (r *BillingInfoResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("TaxNumber", "tax_number").Label("Vergi Numarası").Placeholder("Vergi Numarası").Required(),
		fields.Text("TaxOffice", "tax_office").Label("Vergi Dairesi").Placeholder("Vergi Dairesi").Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *BillingInfoResource) RecordTitle(record interface{}) string {
	if info, ok := record.(*entity.BillingInfo); ok {
		return info.TaxNumber
	}
	return ""
}

func (r *BillingInfoResource) With() []string {
	return []string{"Organization"}
}
