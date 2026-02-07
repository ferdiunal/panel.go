package fields

import (
	"reflect"
)

// HasOne represents a one-to-one relationship (e.g., User -> Profile)
type HasOneField struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	ForeignKeyColumn    string
	OwnerKeyColumn      string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
}

// NewHasOne creates a new HasOne relationship field.
// Accepts either a string (resource slug) or a resource instance.
//
// Örnek kullanım:
//   // String ile:
//   fields.NewHasOne("Profile", "profile", "profiles")
//
//   // Resource instance ile:
//   fields.NewHasOne("Profile", "profile", blog.NewProfileResource())
func HasOne(name, key string, relatedResource interface{}) *HasOneField {
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

	h := &HasOneField{
		Schema: Schema{
			Name:  name,
			Key:   key,
			View:  "has-one-field",
			Type:  TYPE_RELATIONSHIP,
			Props: make(map[string]interface{}),
		},
		RelatedResourceSlug: slug,
		RelatedResource:     resourceInstance,
		ForeignKeyColumn:    slug + "_id",
		OwnerKeyColumn:      "id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithForeignKey(slug + "_id").
			WithReferences("id"),
	}
	// Store relationship details in props for generic access (when Schema interface is used)
	h.WithProps("related_resource", slug)
	if resourceInstance != nil {
		h.WithProps("related_resource_instance", resourceInstance)
	}
	h.WithProps("foreign_key", h.ForeignKeyColumn)
	return h
}

// AutoOptions enables automatic options generation from the related table.
// displayField is the column name to use for the option label.
func (h *HasOneField) AutoOptions(displayField string) *HasOneField {
	h.Schema.AutoOptions(displayField)
	return h
}

// ForeignKey sets the foreign key column name
func (h *HasOneField) ForeignKey(key string) *HasOneField {
	h.ForeignKeyColumn = key
	h.WithProps("foreign_key", key)
	return h
}

// OwnerKey sets the owner key column name
func (h *HasOneField) OwnerKey(key string) *HasOneField {
	h.OwnerKeyColumn = key
	return h
}

// Query sets the query callback for customizing relationship query
func (h *HasOneField) Query(fn func(interface{}) interface{}) *HasOneField {
	h.QueryCallback = fn
	return h
}

// WithEagerLoad sets the loading strategy to eager loading
func (h *HasOneField) WithEagerLoad() *HasOneField {
	h.LoadingStrategy = EAGER_LOADING
	return h
}

// WithLazyLoad sets the loading strategy to lazy loading
func (h *HasOneField) WithLazyLoad() *HasOneField {
	h.LoadingStrategy = LAZY_LOADING
	return h
}

// Extract extracts the value from the resource.
// For HasOne, we want to extract the ID of the related resource if it's a struct.
func (h *HasOneField) Extract(resource interface{}) {
	h.Schema.Extract(resource)

	if h.Schema.Data != nil {
		v := reflect.ValueOf(h.Schema.Data)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				h.Schema.Data = nil
				return
			}
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			// Try to find ID or Id field
			idField := v.FieldByName("ID")
			if !idField.IsValid() {
				idField = v.FieldByName("Id")
			}

			if idField.IsValid() && idField.CanInterface() {
				h.Schema.Data = idField.Interface()
			}
		}
	}
}

// GetRelationshipType returns the relationship type
func (h *HasOneField) GetRelationshipType() string {
	return "hasOne"
}

// GetRelatedResource returns the related resource slug
func (h *HasOneField) GetRelatedResource() string {
	return h.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (h *HasOneField) GetRelationshipName() string {
	return h.Name
}

// ResolveRelationship resolves the relationship by loading single related resource
func (h *HasOneField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// In a real implementation, this would query the database
	// For now, return nil
	return nil, nil
}

// ValidateRelationship validates the relationship
func (h *HasOneField) ValidateRelationship(value interface{}) error {
	// Validate that at most one related resource exists
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for HasOne)
func (h *HasOneField) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for HasOne)
func (h *HasOneField) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (h *HasOneField) GetQueryCallback() func(interface{}) interface{} {
	if h.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return h.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (h *HasOneField) GetLoadingStrategy() LoadingStrategy {
	if h.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return h.LoadingStrategy
}

// Searchable marks the element as searchable (implements Element interface)
func (h *HasOneField) Searchable() Element {
	h.GlobalSearch = true
	return h
}

// IsRequired returns whether the field is required
func (h *HasOneField) IsRequired() bool {
	return h.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for HasOne)
func (h *HasOneField) GetTypes() map[string]string {
	return make(map[string]string)
}
