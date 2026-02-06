package blog

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type CommentResource struct {
	resource.OptimizedBase
}

func NewCommentResource() *CommentResource {
	r := &CommentResource{}
	r.SetModel(&Comment{})
	r.SetSlug("comments")
	r.SetTitle("Comments")
	r.SetIcon("message-circle")
	r.SetGroup("Blog")
	r.SetVisible(true)

	r.SetFieldResolver(&CommentFieldResolver{})
	return r
}

type CommentFieldResolver struct{}

func (r *CommentFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Content", "content").Sortable().Required(),

		// MorphTo Relationship: Comment -> Commentable (Post, Video, etc.)
		fields.NewMorphTo("Commentable", "commentable").
			Types(map[string]string{
				"posts": "posts", // Database Type => Resource Slug
			}),

		fields.DateTime("Created At", "createdAt").ReadOnly().OnList(),
	}
}
