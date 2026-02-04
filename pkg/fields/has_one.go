package fields

// HasOne represents a one-to-one relationship (e.g., User -> Profile)
type HasOne struct {
	Schema
	RelatedResourceSlug string
	ForeignKeyColumn    string
	OwnerKeyColumn      string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
}

// NewHasOne creates a new HasOne relationship field
func NewHasOne(name, key, relatedResource string) *HasOne {
	return &HasOne{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "has-one-field",
			Type: TYPE_RELATIONSHIP,
		},
		RelatedResourceSlug: relatedResource,
		ForeignKeyColumn:    relatedResource + "_id",
		OwnerKeyColumn:      "id",
		LoadingStrategy:     EAGER_LOADING,
	}
}

// ForeignKey sets the foreign key column name
func (h *HasOne) ForeignKey(key string) *HasOne {
	h.ForeignKeyColumn = key
	return h
}

// OwnerKey sets the owner key column name
func (h *HasOne) OwnerKey(key string) *HasOne {
	h.OwnerKeyColumn = key
	return h
}

// Query sets the query callback for customizing relationship query
func (h *HasOne) Query(fn func(interface{}) interface{}) *HasOne {
	h.QueryCallback = fn
	return h
}

// WithEagerLoad sets the loading strategy to eager loading
func (h *HasOne) WithEagerLoad() *HasOne {
	h.LoadingStrategy = EAGER_LOADING
	return h
}

// WithLazyLoad sets the loading strategy to lazy loading
func (h *HasOne) WithLazyLoad() *HasOne {
	h.LoadingStrategy = LAZY_LOADING
	return h
}

// GetRelationshipType returns the relationship type
func (h *HasOne) GetRelationshipType() string {
	return "hasOne"
}

// GetRelatedResource returns the related resource slug
func (h *HasOne) GetRelatedResource() string {
	return h.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (h *HasOne) GetRelationshipName() string {
	return h.Name
}

// ResolveRelationship resolves the relationship by loading single related resource
func (h *HasOne) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// In a real implementation, this would query the database
	// For now, return nil
	return nil, nil
}

// ValidateRelationship validates the relationship
func (h *HasOne) ValidateRelationship(value interface{}) error {
	// Validate that at most one related resource exists
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for HasOne)
func (h *HasOne) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for HasOne)
func (h *HasOne) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (h *HasOne) GetQueryCallback() func(interface{}) interface{} {
	if h.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return h.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (h *HasOne) GetLoadingStrategy() LoadingStrategy {
	if h.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return h.LoadingStrategy
}

// Searchable marks the element as searchable (implements Element interface)
func (h *HasOne) Searchable() Element {
	h.GlobalSearch = true
	return h
}

// IsRequired returns whether the field is required
func (h *HasOne) IsRequired() bool {
	return h.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for HasOne)
func (h *HasOne) GetTypes() map[string]string {
	return make(map[string]string)
}
