package fields

import (
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

// MorphTo represents a polymorphic relationship (e.g., Comment -> Commentable)
type MorphTo struct {
	Schema
	TypeMappings       map[string]string // Type => Resource slug mapping
	DisplayMappings    map[string]string // Type => Display field name
	QueryCallback      func(query interface{}) interface{}
	LoadingStrategy    LoadingStrategy
	GormRelationConfig *RelationshipGormConfig
}

// NewMorphTo creates a new MorphTo relationship field
func NewMorphTo(name, key string) *MorphTo {
	return &MorphTo{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "morph-to-field",
			Type: TYPE_RELATIONSHIP,
			Props: map[string]interface{}{
				"types":    []map[string]string{},
				"displays": map[string]string{},
			},
		},
		TypeMappings:    make(map[string]string),
		DisplayMappings: make(map[string]string),
		LoadingStrategy: EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithPolymorphic(key+"_type", key+"_id"),
	}
}

// Types sets the type mappings for polymorphic relationship
func (m *MorphTo) Types(types map[string]string) *MorphTo {
	m.TypeMappings = types
	m.Props["types"] = m.formatTypesForFrontend(types)
	return m
}

// Displays sets the display field for each type (Type => Field Name)
func (m *MorphTo) Displays(displays map[string]string) *MorphTo {
	m.DisplayMappings = displays
	m.Props["displays"] = displays
	return m
}

// formatTypesForFrontend converts type mappings to frontend select options
func (m *MorphTo) formatTypesForFrontend(types map[string]string) []map[string]string {
	var options []map[string]string
	for dbType, resourceSlug := range types {
		label := resourceSlug
		if len(resourceSlug) > 0 {
			label = string(resourceSlug[0]-32) + resourceSlug[1:]
		}

		options = append(options, map[string]string{
			"label": label,
			"value": dbType,
			"slug":  resourceSlug,
		})
	}
	return options
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

	val := reflect.ValueOf(item)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, nil
	}

	// Calculate field names
	// Key is usually camelCase or snake_case in JS, but we want the struct field name prefix.
	// e.g. "commentable" -> "Commentable"
	baseName := strcase.ToCamel(m.Key)

	typeFieldNames := []string{baseName + "Type", baseName + "_Type"}
	idFieldNames := []string{baseName + "ID", baseName + "Id", baseName + "_ID", baseName + "_Id"}

	var typeVal string
	var idVal interface{}
	var foundType, foundID bool

	// Find Type
	for _, name := range typeFieldNames {
		f := val.FieldByName(name)
		if f.IsValid() {
			typeVal = f.String()
			foundType = true
			break
		}
	}

	// Find ID
	for _, name := range idFieldNames {
		f := val.FieldByName(name)
		if f.IsValid() {
			idVal = f.Interface()
			foundID = true
			break
		}
	}

	if !foundType && !foundID {
		return nil, nil
	}

	return map[string]interface{}{
		"type":        typeVal,
		"id":          idVal,
		"morphToType": typeVal,
		"morphToId":   idVal,
	}, nil
}

// Extract overrides Schema.Extract to handle MorphTo specific extraction
func (m *MorphTo) Extract(resource interface{}) {
	// Don't call Schema.Extract because MorphTo doesn't have a direct field in the struct
	// Instead, directly resolve the polymorphic relationship from type and id fields
	resolved, _ := m.ResolveRelationship(resource)
	m.Data = resolved
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
