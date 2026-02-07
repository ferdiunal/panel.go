package fields

// HasMany represents a one-to-many relationship (e.g., Author -> Posts)
type HasManyField struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	ForeignKeyColumn    string
	OwnerKeyColumn      string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
}

// NewHasMany creates a new HasMany relationship field.
// Accepts either a string (resource slug) or a resource instance.
//
// Örnek kullanım:
//   // String ile:
//   fields.NewHasMany("Posts", "posts", "posts")
//
//   // Resource instance ile:
//   fields.NewHasMany("Posts", "posts", blog.NewPostResource())
func HasMany(name, key string, relatedResource interface{}) *HasManyField {
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

	h := &HasManyField{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "has-many-field",
			Type: TYPE_RELATIONSHIP,
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
	h.WithProps("related_resource", slug)
	if resourceInstance != nil {
		h.WithProps("related_resource_instance", resourceInstance)
	}
	return h
}

// ForeignKey sets the foreign key column name
func (h *HasManyField) ForeignKey(key string) *HasManyField {
	h.ForeignKeyColumn = key
	return h
}

// OwnerKey sets the owner key column name
func (h *HasManyField) OwnerKey(key string) *HasManyField {
	h.OwnerKeyColumn = key
	return h
}

// Query sets the query callback for customizing relationship query
func (h *HasManyField) Query(fn func(interface{}) interface{}) *HasManyField {
	h.QueryCallback = fn
	return h
}

// WithEagerLoad sets the loading strategy to eager loading
func (h *HasManyField) WithEagerLoad() *HasManyField {
	h.LoadingStrategy = EAGER_LOADING
	return h
}

// WithLazyLoad sets the loading strategy to lazy loading
func (h *HasManyField) WithLazyLoad() *HasManyField {
	h.LoadingStrategy = LAZY_LOADING
	return h
}

// GetRelationshipType returns the relationship type
func (h *HasManyField) GetRelationshipType() string {
	return "hasMany"
}

// GetRelatedResource returns the related resource slug
func (h *HasManyField) GetRelatedResource() string {
	return h.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (h *HasManyField) GetRelationshipName() string {
	return h.Name
}

// ResolveRelationship resolves the relationship by loading all related resources
func (h *HasManyField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}

	// In a real implementation, this would query the database
	// For now, return empty slice
	return []interface{}{}, nil
}

// ValidateRelationship validates the relationship
func (h *HasManyField) ValidateRelationship(value interface{}) error {
	// Validate that foreign key references are valid
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for HasMany)
func (h *HasManyField) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for HasMany)
func (h *HasManyField) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (h *HasManyField) GetQueryCallback() func(interface{}) interface{} {
	if h.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return h.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (h *HasManyField) GetLoadingStrategy() LoadingStrategy {
	if h.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return h.LoadingStrategy
}

// Searchable marks the element as searchable (implements Element interface)
func (h *HasManyField) Searchable() Element {
	h.GlobalSearch = true
	return h
}

// Count returns the count of related resources
func (h *HasManyField) Count() int64 {
	// In a real implementation, this would execute a COUNT query
	return 0
}

// IsRequired returns whether the field is required
func (h *HasManyField) IsRequired() bool {
	return h.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for HasMany)
func (h *HasManyField) GetTypes() map[string]string {
	return make(map[string]string)
}
