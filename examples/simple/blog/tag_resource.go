package blog

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type TagResource struct {
	resource.OptimizedBase
}

func NewTagResource() *TagResource {
	r := &TagResource{}
	r.SetModel(&Tag{})
	r.SetSlug("tags")
	r.SetTitle("Tags")
	r.SetIcon("tag")
	r.SetGroup("Blog")
	r.SetVisible(true)

	r.SetFieldResolver(&TagFieldResolver{})
	return r
}

type TagFieldResolver struct{}

func (r *TagFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Name", "name").Sortable().Required(),

		fields.DateTime("Created At", "createdAt").ReadOnly().OnList(),
	}
}
