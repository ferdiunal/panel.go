package cargo_company

import (
	"cargo.go/entity"
	"cargo.go/resources"        // Registry import
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// init fonksiyonu ile resource'u registry'ye kaydet
func init() {
	resources.Register("cargo-companies", func() resource.Resource {
		return NewCargoCompanyResource()
	})
}

type CargoCompanyResource struct {
	resource.OptimizedBase
}
type CargoCompanyResolveFields struct{}

func NewCargoCompanyResource() *CargoCompanyResource {
	r := &CargoCompanyResource{}
	r.SetSlug("cargo-companies")
	r.SetTitle("Cargo Companies")
	r.SetIcon("cargo_companies")
	r.SetGroup("Organization")
	r.SetModel(&entity.CargoCompany{})
	r.SetFieldResolver(&CargoCompanyResolveFields{})
	r.SetVisible(true)
	return r
}

func (r *CargoCompanyResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for a cargo company record
func (r *CargoCompanyResource) RecordTitle(record interface{}) string {
	if company, ok := record.(*entity.CargoCompany); ok {
		return company.Name
	}
	return ""
}
