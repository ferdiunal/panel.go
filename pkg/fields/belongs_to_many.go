package fields

import (
	"strings"
)

// BelongsToMany represents a many-to-many relationship (e.g., User -> Roles)
type BelongsToMany struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	PivotTableName      string
	ForeignKeyColumn    string
	RelatedKeyColumn    string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
}

// NewBelongsToMany creates a new BelongsToMany relationship field
func NewBelongsToMany(name, key, relatedResource string) *BelongsToMany {
	// Generate pivot table name using convention
	pivotTable := generatePivotTableName(key, relatedResource)

	b := &BelongsToMany{
		Schema: Schema{
			Name:  name,
			Key:   key,
			View:  "belongs-to-many-field",
			Type:  TYPE_RELATIONSHIP,
			Props: make(map[string]interface{}),
		},
		RelatedResourceSlug: relatedResource,
		PivotTableName:      pivotTable,
		ForeignKeyColumn:    "user_id",
		RelatedKeyColumn:    relatedResource + "_id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithPivotTable(pivotTable, "user_id", relatedResource+"_id"),
	}
	b.WithProps("related_resource", relatedResource)
	return b
}

// NewBelongsToManyResource, resource instance kullanarak BelongsToMany ilişkisi oluşturur.
// Bu metod, resource referansı kullanarak tip güvenli ilişki tanımlaması sağlar.
//
// Örnek kullanım:
//   fields.NewBelongsToManyResource("Tags", "tags", blog.NewTagResource())
func NewBelongsToManyResource(name, key string, relatedResource interface{}) *BelongsToMany {
	// Resource interface'inden slug'ı al
	type resourceSlugger interface {
		Slug() string
	}

	var slug string
	if res, ok := relatedResource.(resourceSlugger); ok {
		slug = res.Slug()
	} else {
		slug = ""
	}

	// Generate pivot table name using convention
	pivotTable := generatePivotTableName(key, slug)

	b := &BelongsToMany{
		Schema: Schema{
			Name:  name,
			Key:   key,
			View:  "belongs-to-many-field",
			Type:  TYPE_RELATIONSHIP,
			Props: make(map[string]interface{}),
		},
		RelatedResourceSlug: slug,
		RelatedResource:     relatedResource,
		PivotTableName:      pivotTable,
		ForeignKeyColumn:    "user_id",
		RelatedKeyColumn:    slug + "_id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithPivotTable(pivotTable, "user_id", slug+"_id"),
	}
	b.WithProps("related_resource", slug)
	b.WithProps("related_resource_instance", relatedResource)
	return b
}

// AutoOptions enables automatic options generation from the related table.
// displayField is the column name to use for the option label.
func (b *BelongsToMany) AutoOptions(displayField string) *BelongsToMany {
	b.Schema.AutoOptions(displayField)
	return b
}

// generatePivotTableName generates a pivot table name using convention
func generatePivotTableName(key, relatedResource string) string {
	// Convert to snake_case and sort alphabetically
	parts := []string{key, relatedResource}
	// Simple sort: if first > second, swap
	if parts[0] > parts[1] {
		parts[0], parts[1] = parts[1], parts[0]
	}
	return strings.Join(parts, "_")
}

// PivotTable sets the pivot table name
func (b *BelongsToMany) PivotTable(table string) *BelongsToMany {
	b.PivotTableName = table
	return b
}

// ForeignKey sets the foreign key in pivot table
func (b *BelongsToMany) ForeignKey(key string) *BelongsToMany {
	b.ForeignKeyColumn = key
	return b
}

// RelatedKey sets the related key in pivot table
func (b *BelongsToMany) RelatedKey(key string) *BelongsToMany {
	b.RelatedKeyColumn = key
	return b
}

// Query sets the query callback for customizing relationship query
func (b *BelongsToMany) Query(fn func(interface{}) interface{}) *BelongsToMany {
	b.QueryCallback = fn
	return b
}

// WithEagerLoad sets the loading strategy to eager loading
func (b *BelongsToMany) WithEagerLoad() *BelongsToMany {
	b.LoadingStrategy = EAGER_LOADING
	return b
}

// WithLazyLoad sets the loading strategy to lazy loading
func (b *BelongsToMany) WithLazyLoad() *BelongsToMany {
	b.LoadingStrategy = LAZY_LOADING
	return b
}

// GetRelationshipType returns the relationship type
func (b *BelongsToMany) GetRelationshipType() string {
	return "belongsToMany"
}

// GetRelatedResource returns the related resource slug
func (b *BelongsToMany) GetRelatedResource() string {
	return b.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (b *BelongsToMany) GetRelationshipName() string {
	return b.Name
}

// ResolveRelationship resolves the relationship by loading through pivot table
func (b *BelongsToMany) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}

	// In a real implementation, this would query the database through pivot table
	// For now, return empty slice
	return []interface{}{}, nil
}

// ValidateRelationship validates the relationship
func (b *BelongsToMany) ValidateRelationship(value interface{}) error {
	// Validate that pivot table entries are valid
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for BelongsToMany)
func (b *BelongsToMany) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for BelongsToMany)
func (b *BelongsToMany) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (b *BelongsToMany) GetQueryCallback() func(interface{}) interface{} {
	if b.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return b.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (b *BelongsToMany) GetLoadingStrategy() LoadingStrategy {
	if b.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return b.LoadingStrategy
}

// Searchable marks the element as searchable (implements Element interface)
func (b *BelongsToMany) Searchable() Element {
	b.GlobalSearch = true
	return b
}

// Count returns the count of related resources
func (b *BelongsToMany) Count() int64 {
	// In a real implementation, this would execute a COUNT query on pivot table
	return 0
}

// IsRequired returns whether the field is required
func (b *BelongsToMany) IsRequired() bool {
	return b.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for BelongsToMany)
func (b *BelongsToMany) GetTypes() map[string]string {
	return make(map[string]string)
}
