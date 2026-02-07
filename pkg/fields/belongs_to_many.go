package fields

import (
	"strings"
)

// BelongsToMany represents a many-to-many relationship (e.g., User -> Roles)
type BelongsToManyField struct {
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

// NewBelongsToMany creates a new BelongsToMany relationship field.
// Accepts either a string (resource slug) or a resource instance.
//
// Örnek kullanım:
//   // String ile:
//   fields.NewBelongsToMany("Tags", "tags", "tags")
//
//   // Resource instance ile:
//   fields.NewBelongsToMany("Tags", "tags", blog.NewTagResource())
func BelongsToMany(name, key string, relatedResource interface{}) *BelongsToManyField {
	// Resource interface'inden slug'ı al
	type resourceSlugger interface {
		Slug() string
	}

	var slug string
	var resourceInstance interface{}

	// Check if relatedResource is a string or a resource instance
	if slugStr, ok := relatedResource.(string); ok {
		// String slug provided
		slug = slugStr
	} else if res, ok := relatedResource.(resourceSlugger); ok {
		// Resource instance provided
		slug = res.Slug()
		resourceInstance = relatedResource
	} else {
		// Fallback: empty slug
		slug = ""
	}

	// Generate pivot table name using convention
	pivotTable := generatePivotTableName(key, slug)

	b := &BelongsToManyField{
		Schema: Schema{
			Name:  name,
			Key:   key,
			View:  "belongs-to-many-field",
			Type:  TYPE_RELATIONSHIP,
			Props: make(map[string]interface{}),
		},
		RelatedResourceSlug: slug,
		RelatedResource:     resourceInstance,
		PivotTableName:      pivotTable,
		ForeignKeyColumn:    "user_id",
		RelatedKeyColumn:    slug + "_id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithPivotTable(pivotTable, "user_id", slug+"_id"),
	}
	b.WithProps("related_resource", slug)
	if resourceInstance != nil {
		b.WithProps("related_resource_instance", resourceInstance)
	}
	return b
}

// AutoOptions enables automatic options generation from the related table.
// displayField is the column name to use for the option label.
func (b *BelongsToManyField) AutoOptions(displayField string) *BelongsToManyField {
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
func (b *BelongsToManyField) PivotTable(table string) *BelongsToManyField {
	b.PivotTableName = table
	return b
}

// ForeignKey sets the foreign key in pivot table
func (b *BelongsToManyField) ForeignKey(key string) *BelongsToManyField {
	b.ForeignKeyColumn = key
	return b
}

// RelatedKey sets the related key in pivot table
func (b *BelongsToManyField) RelatedKey(key string) *BelongsToManyField {
	b.RelatedKeyColumn = key
	return b
}

// Query sets the query callback for customizing relationship query
func (b *BelongsToManyField) Query(fn func(interface{}) interface{}) *BelongsToManyField {
	b.QueryCallback = fn
	return b
}

// WithEagerLoad sets the loading strategy to eager loading
func (b *BelongsToManyField) WithEagerLoad() *BelongsToManyField {
	b.LoadingStrategy = EAGER_LOADING
	return b
}

// WithLazyLoad sets the loading strategy to lazy loading
func (b *BelongsToManyField) WithLazyLoad() *BelongsToManyField {
	b.LoadingStrategy = LAZY_LOADING
	return b
}

// GetRelationshipType returns the relationship type
func (b *BelongsToManyField) GetRelationshipType() string {
	return "belongsToMany"
}

// GetRelatedResource returns the related resource slug
func (b *BelongsToManyField) GetRelatedResource() string {
	return b.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (b *BelongsToManyField) GetRelationshipName() string {
	return b.Name
}

// ResolveRelationship resolves the relationship by loading through pivot table
func (b *BelongsToManyField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}

	// In a real implementation, this would query the database through pivot table
	// For now, return empty slice
	return []interface{}{}, nil
}

// ValidateRelationship validates the relationship
func (b *BelongsToManyField) ValidateRelationship(value interface{}) error {
	// Validate that pivot table entries are valid
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for BelongsToMany)
func (b *BelongsToManyField) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for BelongsToMany)
func (b *BelongsToManyField) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (b *BelongsToManyField) GetQueryCallback() func(interface{}) interface{} {
	if b.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return b.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (b *BelongsToManyField) GetLoadingStrategy() LoadingStrategy {
	if b.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return b.LoadingStrategy
}

// Searchable marks the element as searchable (implements Element interface)
func (b *BelongsToManyField) Searchable() Element {
	b.GlobalSearch = true
	return b
}

// Count returns the count of related resources
func (b *BelongsToManyField) Count() int64 {
	// In a real implementation, this would execute a COUNT query on pivot table
	return 0
}

// IsRequired returns whether the field is required
func (b *BelongsToManyField) IsRequired() bool {
	return b.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for BelongsToMany)
func (b *BelongsToManyField) GetTypes() map[string]string {
	return make(map[string]string)
}
