package prices

import (
	"cargo.go/entity"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type PriceResource struct {
	resource.OptimizedBase
}

type PriceResolveFields struct{}

func NewPriceResource() *PriceResource {
	r := &PriceResource{}
	r.SetSlug("prices")
	r.SetTitle("Prices")
	r.SetIcon("prices")
	r.SetGroup("Operations")
	r.SetFieldResolver(&PriceResolveFields{})
	r.SetModel(&entity.Price{})
	r.SetVisible(true)
	return r
}

func (r *PriceResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("PriceList", "price_list_id", "price-lists").
			DisplayUsing("name").
			AutoOptions("name"),
		fields.Number("Price", "price"),
		fields.Number("DesiThreshold", "desi_threshold"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}
