package commissions

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
	resources.Register("commissions", func() resource.Resource {
		return NewCommissionResource()
	})
}

type CommissionResource struct {
	resource.OptimizedBase
}

type CommissionResolveFields struct{}

func NewCommissionResource() *CommissionResource {
	r := &CommissionResource{}
	r.SetSlug("commissions")
	r.SetTitle("Commissions")
	r.SetIcon("commissions")
	r.SetGroup("Operations")
	r.SetFieldResolver(&CommissionResolveFields{})
	r.SetModel(&entity.Commission{})
	r.SetVisible(true)
	return r
}

func (r *CommissionResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").
			WithSearchableColumns("name", "email").
			WithEagerLoad(),
		fields.Text("Name", "name"),
		fields.BelongsTo("PriceList", "price_list_id", resources.GetOrPanic("price-lists")).
			DisplayUsing("name").
			WithSearchableColumns("name").
			WithEagerLoad(),
		fields.Number("Commission", "commission"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for a commission record
func (r *CommissionResource) RecordTitle(record interface{}) string {
	if commission, ok := record.(*entity.Commission); ok {
		return commission.Name
	}
	return ""
}
