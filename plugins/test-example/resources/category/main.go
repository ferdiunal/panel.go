package category

import (
	"test-example/entity"
	"test-example/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("categories", func() resource.Resource {
		return NewCategoryResource()
	})
}

type CategoryResource struct {
	resource.OptimizedBase
}

type CategoryPolicy struct{}

func (p *CategoryPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *CategoryPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *CategoryPolicy) Create(ctx *context.Context) bool    { return true }
func (p *CategoryPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *CategoryPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type CategoryResolveFields struct{}

func NewCategoryResource() *CategoryResource {
	r := &CategoryResource{}
	r.SetSlug("categories")
	r.SetTitle("Categories")
	r.SetIcon("tag")
	r.SetGroup("Management")
	r.SetModel(&entity.Category{})
	r.SetFieldResolver(&CategoryResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&CategoryPolicy{})
	return r
}

func (r *CategoryResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name").Label("Kategori Adı").Placeholder("Kategori Adı").Required(),
		fields.BelongsToMany("Products", "products", resources.GetOrPanic("products")).
			PivotTable("product_categories").WithEagerLoad().Label("Ürünler").HideOnList(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *CategoryResource) RecordTitle(record interface{}) string {
	if cat, ok := record.(*entity.Category); ok {
		return cat.Name
	}
	return ""
}

func (r *CategoryResource) With() []string {
	return []string{"Products"}
}
