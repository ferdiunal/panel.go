package fields

// HasMany represents a one-to-many relationship (e.g., Author -> Posts)
type HasMany struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	ForeignKeyColumn    string
	OwnerKeyColumn      string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
}

// NewHasMany creates a new HasMany relationship field
func NewHasMany(name, key, relatedResource string) *HasMany {
	h := &HasMany{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "has-many-field",
			Type: TYPE_RELATIONSHIP,
		},
		RelatedResourceSlug: relatedResource,
		ForeignKeyColumn:    relatedResource + "_id",
		OwnerKeyColumn:      "id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithForeignKey(relatedResource + "_id").
			WithReferences("id"),
	}
	h.WithProps("related_resource", relatedResource)
	return h
}

// NewHasManyResource, resource instance kullanarak HasMany ilişkisi oluşturur.
// Bu metod, resource referansı kullanarak tip güvenli ilişki tanımlaması sağlar.
//
// Örnek kullanım:
//   fields.NewHasManyResource("Posts", "posts", blog.NewPostResource())
func NewHasManyResource(name, key string, relatedResource interface{}) *HasMany {
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

	h := &HasMany{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "has-many-field",
			Type: TYPE_RELATIONSHIP,
		},
		RelatedResourceSlug: slug,
		RelatedResource:     relatedResource,
		ForeignKeyColumn:    slug + "_id",
		OwnerKeyColumn:      "id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithForeignKey(slug + "_id").
			WithReferences("id"),
	}
	h.WithProps("related_resource", slug)
	h.WithProps("related_resource_instance", relatedResource)
	return h
}

// ForeignKey sets the foreign key column name
func (h *HasMany) ForeignKey(key string) *HasMany {
	h.ForeignKeyColumn = key
	return h
}

// OwnerKey sets the owner key column name
func (h *HasMany) OwnerKey(key string) *HasMany {
	h.OwnerKeyColumn = key
	return h
}

// Query sets the query callback for customizing relationship query
func (h *HasMany) Query(fn func(interface{}) interface{}) *HasMany {
	h.QueryCallback = fn
	return h
}

// WithEagerLoad sets the loading strategy to eager loading
func (h *HasMany) WithEagerLoad() *HasMany {
	h.LoadingStrategy = EAGER_LOADING
	return h
}

// WithLazyLoad sets the loading strategy to lazy loading
func (h *HasMany) WithLazyLoad() *HasMany {
	h.LoadingStrategy = LAZY_LOADING
	return h
}

// GetRelationshipType returns the relationship type
func (h *HasMany) GetRelationshipType() string {
	return "hasMany"
}

// GetRelatedResource returns the related resource slug
func (h *HasMany) GetRelatedResource() string {
	return h.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (h *HasMany) GetRelationshipName() string {
	return h.Name
}

// ResolveRelationship resolves the relationship by loading all related resources
func (h *HasMany) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}

	// In a real implementation, this would query the database
	// For now, return empty slice
	return []interface{}{}, nil
}

// ValidateRelationship validates the relationship
func (h *HasMany) ValidateRelationship(value interface{}) error {
	// Validate that foreign key references are valid
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for HasMany)
func (h *HasMany) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for HasMany)
func (h *HasMany) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (h *HasMany) GetQueryCallback() func(interface{}) interface{} {
	if h.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return h.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (h *HasMany) GetLoadingStrategy() LoadingStrategy {
	if h.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return h.LoadingStrategy
}

// Searchable marks the element as searchable (implements Element interface)
func (h *HasMany) Searchable() Element {
	h.GlobalSearch = true
	return h
}

// Count returns the count of related resources
func (h *HasMany) Count() int64 {
	// In a real implementation, this would execute a COUNT query
	return 0
}

// IsRequired returns whether the field is required
func (h *HasMany) IsRequired() bool {
	return h.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for HasMany)
func (h *HasMany) GetTypes() map[string]string {
	return make(map[string]string)
}
