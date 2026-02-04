package fields

import (
	"fmt"
)

// MorphTo represents a polymorphic relationship (e.g., Comment -> Commentable)
type MorphTo struct {
	Schema
	TypeMappings    map[string]string // Type => Resource slug mapping
	QueryCallback   func(query interface{}) interface{}
	LoadingStrategy LoadingStrategy
}

// NewMorphTo creates a new MorphTo relationship field
func NewMorphTo(name, key string) *MorphTo {
	return &MorphTo{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "morph-to-field",
			Type: TYPE_RELATIONSHIP,
		},
		TypeMappings:    make(map[string]string),
		LoadingStrategy: EAGER_LOADING,
	}
}

// Types sets the type mappings for polymorphic relationship
func (m *MorphTo) Types(types map[string]string) *MorphTo {
	m.TypeMappings = types
	return m
}

// Query sets the query callback for customizing relationship query
func (m *MorphTo) Query(fn func(interface{}) interface{}) *MorphTo {
	m.QueryCallback = fn
	return m
}

// WithEagerLoad sets the loading strategy to eager loading
func (m *MorphTo) WithEagerLoad() *MorphTo {
	m.LoadingStrategy = EAGER_LOADING
	return m
}

// WithLazyLoad sets the loading strategy to lazy loading
func (m *MorphTo) WithLazyLoad() *MorphTo {
	m.LoadingStrategy = LAZY_LOADING
	return m
}

// GetRelationshipType returns the relationship type
func (m *MorphTo) GetRelationshipType() string {
	return "morphTo"
}

// GetRelatedResource returns the related resource slug (not applicable for MorphTo)
func (m *MorphTo) GetRelatedResource() string {
	return ""
}

// GetRelationshipName returns the relationship name
func (m *MorphTo) GetRelationshipName() string {
	return m.Name
}

// ResolveRelationship resolves the relationship by loading based on morph type
func (m *MorphTo) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// In a real implementation, this would:
	// 1. Extract the morph type from the item
	// 2. Look up the resource slug from the type mapping
	// 3. Query the appropriate resource table
	// For now, return nil
	return nil, nil
}

// ValidateRelationship validates the relationship
func (m *MorphTo) ValidateRelationship(value interface{}) error {
	// Validate that morph type is registered
	// In a real implementation, this would check that the type exists in the mapping
	return nil
}

// GetDisplayKey returns the display key (not used for MorphTo)
func (m *MorphTo) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for MorphTo)
func (m *MorphTo) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (m *MorphTo) GetQueryCallback() func(interface{}) interface{} {
	if m.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return m.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (m *MorphTo) GetLoadingStrategy() LoadingStrategy {
	if m.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return m.LoadingStrategy
}

// GetTypes returns the type mappings
func (m *MorphTo) GetTypes() map[string]string {
	if m.TypeMappings == nil {
		return make(map[string]string)
	}
	return m.TypeMappings
}

// GetResourceForType returns the resource slug for a given type
func (m *MorphTo) GetResourceForType(morphType string) (string, error) {
	resource, ok := m.TypeMappings[morphType]
	if !ok {
		return "", fmt.Errorf("morph type '%s' is not registered", morphType)
	}
	return resource, nil
}

// Searchable marks the element as searchable (implements Element interface)
func (m *MorphTo) Searchable() Element {
	m.GlobalSearch = true
	return m
}

// IsRequired returns whether the field is required
func (m *MorphTo) IsRequired() bool {
	return m.Schema.IsRequired
}
