package blog

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type ProfileResource struct {
	resource.OptimizedBase
}

func NewProfileResource() *ProfileResource {
	r := &ProfileResource{}
	r.SetModel(&Profile{})
	r.SetSlug("profiles")
	r.SetTitle("Profiles")
	r.SetIcon("id-card")
	r.SetGroup("Blog")
	r.SetVisible(true)

	r.SetFieldResolver(&ProfileFieldResolver{})
	return r
}

func (r *ProfileResource) With() []string {
	return []string{"Author"}
}

type ProfileFieldResolver struct{}

func (r *ProfileFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),

		// BelongsTo Relationship: Profile -> Author
		fields.NewBelongsTo("Author", "author_id", "authors").AutoOptions("name"),

		fields.Text("Bio", "bio").Sortable(),
		fields.Text("Website", "website").Sortable(),

		fields.DateTime("Created At", "createdAt").ReadOnly().OnList(),
	}
}
