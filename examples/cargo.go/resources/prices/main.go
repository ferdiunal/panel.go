package prices

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
	resources.Register("price-lists", func() resource.Resource {
		return NewPriceListResource()
	})
}

type PriceListResource struct {
	resource.OptimizedBase
}

type PriceListResolveFields struct{}

func NewPriceListResource() *PriceListResource {
	r := &PriceListResource{}
	r.SetSlug("price-lists")
	r.SetTitle("Price Lists")
	r.SetIcon("price_lists")
	r.SetGroup("Operations")
	r.SetFieldResolver(&PriceListResolveFields{})
	r.SetModel(&entity.PriceList{})
	r.SetVisible(true)
	return r
}

func (r *PriceListResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("CargoCompany", "cargo_company_id", resources.GetOrPanic("cargo-companies")).
			DisplayUsing("name").
			WithSearchableColumns("name").
			AutoOptions("name").
			WithEagerLoad(),
		fields.Text("Name", "name"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// RecordTitle returns the title for a price list record
func (r *PriceListResource) RecordTitle(record interface{}) string {
	if priceList, ok := record.(*entity.PriceList); ok {
		return priceList.Name
	}
	return ""
}
