package fields

import (
	"fmt"
	"reflect"
)

// BelongsTo represents a one-to-one inverse relationship (e.g., Post -> Author)
type BelongsTo struct {
	Schema
	RelatedResourceSlug string
	DisplayKey          string
	SearchableColumns   []string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
}

// NewBelongsTo creates a new BelongsTo relationship field
func NewBelongsTo(name, key, relatedResource string) *BelongsTo {
	return &BelongsTo{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "belongs-to-field",
			Type: TYPE_RELATIONSHIP,
		},
		RelatedResourceSlug: relatedResource,
		DisplayKey:          "name",
		SearchableColumns:   []string{"name"},
		LoadingStrategy:     EAGER_LOADING,
	}
}

// DisplayUsing sets the display key for showing related resource
func (b *BelongsTo) DisplayUsing(key string) *BelongsTo {
	b.DisplayKey = key
	return b
}

// WithSearchableColumns sets the searchable columns for BelongsTo
func (b *BelongsTo) WithSearchableColumns(columns ...string) *BelongsTo {
	b.SearchableColumns = columns
	return b
}

// Searchable marks the element as searchable (implements Element interface)
func (b *BelongsTo) Searchable() Element {
	b.GlobalSearch = true
	return b
}

// Query sets the query callback for customizing relationship query
func (b *BelongsTo) Query(fn func(interface{}) interface{}) *BelongsTo {
	b.QueryCallback = fn
	return b
}

// WithEagerLoad sets the loading strategy to eager loading
func (b *BelongsTo) WithEagerLoad() *BelongsTo {
	b.LoadingStrategy = EAGER_LOADING
	return b
}

// WithLazyLoad sets the loading strategy to lazy loading
func (b *BelongsTo) WithLazyLoad() *BelongsTo {
	b.LoadingStrategy = LAZY_LOADING
	return b
}

// GetRelationshipType returns the relationship type
func (b *BelongsTo) GetRelationshipType() string {
	return "belongsTo"
}

// GetRelatedResource returns the related resource slug
func (b *BelongsTo) GetRelatedResource() string {
	return b.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (b *BelongsTo) GetRelationshipName() string {
	return b.Name
}

// ResolveRelationship resolves the relationship value using reflection
func (b *BelongsTo) ResolveRelationship(item interface{}) (interface{}, error) {
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
func (b *BelongsTo) ValidateRelationship(value interface{}) error {
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
func (b *BelongsTo) GetDisplayKey() string {
	if b.DisplayKey == "" {
		return "name"
	}
	return b.DisplayKey
}

// GetSearchableColumns returns the searchable columns
func (b *BelongsTo) GetSearchableColumns() []string {
	if b.SearchableColumns == nil {
		return []string{}
	}
	return b.SearchableColumns
}

// GetQueryCallback returns the query callback
func (b *BelongsTo) GetQueryCallback() func(interface{}) interface{} {
	if b.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return b.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (b *BelongsTo) GetLoadingStrategy() LoadingStrategy {
	if b.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return b.LoadingStrategy
}

// IsRequired returns whether the field is required
func (b *BelongsTo) IsRequired() bool {
	return b.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for BelongsTo)
func (b *BelongsTo) GetTypes() map[string]string {
	return make(map[string]string)
}
