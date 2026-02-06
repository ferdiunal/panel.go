package blog

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type PostResource struct {
	resource.OptimizedBase
}

func NewPostResource() *PostResource {
	r := &PostResource{}
	r.SetModel(&Post{})
	r.SetSlug("posts")
	r.SetTitle("Posts")
	r.SetIcon("file-text")
	r.SetGroup("Blog")
	r.SetVisible(true)

	r.SetFieldResolver(&PostFieldResolver{})
	return r
}

func (r *PostResource) With() []string {
	return []string{"Author", "Tags"}
}

type PostFieldResolver struct{}

func (r *PostFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Title", "title").Sortable().Required(),
		fields.Text("Content", "content").HideOnList(),

		// BelongsTo Relationship: Post -> Author
		fields.NewBelongsTo("Author", "author_id", "authors").AutoOptions("name").Required(),

		// BelongsToMany Relationship: Post <-> Tag
		fields.NewBelongsToMany("Tags", "tags", "tags").AutoOptions("name"),

		fields.DateTime("Created At", "createdAt").ReadOnly().OnList(),
	}
}
