package blog

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

type AuthorResource struct {
	resource.OptimizedBase
}

func NewAuthorResource() *AuthorResource {
	r := &AuthorResource{}
	r.SetModel(&Author{})
	r.SetSlug("authors")
	r.SetTitle("Authors")
	r.SetIcon("user")
	r.SetGroup("Blog")
	r.SetVisible(true)

	r.SetFieldResolver(&AuthorFieldResolver{})
	return r
}

func (r *AuthorResource) With() []string {
	return []string{"Profile"}
}

type AuthorFieldResolver struct{}

func (r *AuthorFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Name", "name").Sortable().Required(),
		fields.Email("Email", "email").Sortable().Required(),

		// HasOne Relationship: Author -> Profile
		fields.HasOne("Profile", "profile", "profiles").
			ForeignKey("author_id").
			AutoOptions("bio"),
		fields.DateTime("Created At", "createdAt").ReadOnly().OnList(),
	}
}
