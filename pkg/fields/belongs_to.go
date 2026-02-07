package fields

import (
	"fmt"
	"reflect"
)

// BelongsTo represents a one-to-one inverse relationship (e.g., Post -> Author)
type BelongsToField struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	DisplayKey          string
	SearchableColumns   []string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
}

// NewBelongsTo creates a new BelongsTo relationship field.
// Accepts either a string (resource slug) or a resource instance.
//
// Örnek kullanım:
//   // String ile:
//   fields.NewBelongsTo("Author", "author_id", "authors")
//
//   // Resource instance ile:
//   fields.NewBelongsTo("Author", "author_id", blog.NewAuthorResource())
func BelongsTo(name, key string, relatedResource interface{}) *BelongsToField {
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

	b := &BelongsToField{
		Schema: Schema{
			Name:  name,
			Key:   key,
			View:  "belongs-to-field",
			Type:  TYPE_RELATIONSHIP,
			Props: make(map[string]interface{}),
		},
		RelatedResourceSlug: slug,
		RelatedResource:     resourceInstance,
		DisplayKey:          "name",
		SearchableColumns:   []string{"name"},
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithForeignKey(key). // Don't add "_id" suffix - key should already include it
			WithReferences("id"),
	}
	b.WithProps("related_resource", slug)
	if resourceInstance != nil {
		b.WithProps("related_resource_instance", resourceInstance)
	}
	return b
}

// AutoOptions enables automatic options generation from the related table.
// displayField is the column name to use for the option label.
func (b *BelongsToField) AutoOptions(displayField string) *BelongsToField {
	b.Schema.AutoOptions(displayField)
	return b
}

// DisplayUsing sets the display key for showing related resource

func (b *BelongsToField) DisplayUsing(key string) *BelongsToField {
	b.DisplayKey = key
	return b
}

// WithSearchableColumns sets the searchable columns for BelongsTo
func (b *BelongsToField) WithSearchableColumns(columns ...string) *BelongsToField {
	b.SearchableColumns = columns
	return b
}

// Searchable marks the element as searchable (implements Element interface)
func (b *BelongsToField) Searchable() Element {
	b.GlobalSearch = true
	return b
}

// Query sets the query callback for customizing relationship query
func (b *BelongsToField) Query(fn func(interface{}) interface{}) *BelongsToField {
	b.QueryCallback = fn
	return b
}

// WithEagerLoad sets the loading strategy to eager loading
func (b *BelongsToField) WithEagerLoad() *BelongsToField {
	b.LoadingStrategy = EAGER_LOADING
	return b
}

// WithLazyLoad sets the loading strategy to lazy loading
func (b *BelongsToField) WithLazyLoad() *BelongsToField {
	b.LoadingStrategy = LAZY_LOADING
	return b
}

// GetRelationshipType returns the relationship type
func (b *BelongsToField) GetRelationshipType() string {
	return "belongsTo"
}

// GetRelatedResource returns the related resource slug
func (b *BelongsToField) GetRelatedResource() string {
	return b.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (b *BelongsToField) GetRelationshipName() string {
	return b.Name
}

// ResolveRelationship resolves the relationship value using reflection
func (b *BelongsToField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		if b.Schema.IsRequired {
			return nil, &RelationshipError{
				FieldName:        b.Name,
				RelationshipType: "belongsTo",
				Message:          "Related resource is required",
				Context: map[string]interface{}{
					"related_resource": b.RelatedResourceSlug,
					"display_key":      b.DisplayKey,
				},
			}
		}
		return nil, nil
	}

	// Extract the relationship value using reflection
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Try to get the field by key name
	fieldVal := v.FieldByName(b.Key)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("field %s not found in %T", b.Key, item)
	}

	return fieldVal.Interface(), nil
}

// ValidateRelationship validates the relationship
func (b *BelongsToField) ValidateRelationship(value interface{}) error {
	if value == nil {
		if b.Schema.IsRequired {
			return &RelationshipError{
				FieldName:        b.Name,
				RelationshipType: "belongsTo",
				Message:          "Related resource is required",
				Context: map[string]interface{}{
					"related_resource": b.RelatedResourceSlug,
				},
			}
		}
		return nil
	}

	// Additional validation logic can be added here
	return nil
}

// GetDisplayKey returns the display key
func (b *BelongsToField) GetDisplayKey() string {
	if b.DisplayKey == "" {
		return "name"
	}
	return b.DisplayKey
}

// GetSearchableColumns returns the searchable columns
func (b *BelongsToField) GetSearchableColumns() []string {
	if b.SearchableColumns == nil {
		return []string{}
	}
	return b.SearchableColumns
}

// GetQueryCallback returns the query callback
func (b *BelongsToField) GetQueryCallback() func(interface{}) interface{} {
	if b.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return b.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (b *BelongsToField) GetLoadingStrategy() LoadingStrategy {
	if b.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return b.LoadingStrategy
}

// IsRequired returns whether the field is required
func (b *BelongsToField) IsRequired() bool {
	return b.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for BelongsTo)
func (b *BelongsToField) GetTypes() map[string]string {
	return make(map[string]string)
}
