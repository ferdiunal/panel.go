package fields

import (
	"testing"
)

// TestBelongsToManyCreation tests BelongsToMany field creation
func TestBelongsToManyCreation(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")

	if field.Name != "Tags" {
		t.Errorf("Expected name 'Tags', got '%s'", field.Name)
	}

	if field.Key != "tags" {
		t.Errorf("Expected key 'tags', got '%s'", field.Key)
	}

	if field.RelatedResourceSlug != "tags" {
		t.Errorf("Expected related resource 'tags', got '%s'", field.RelatedResourceSlug)
	}

	if field.View != "belongs-to-many-field" {
		t.Errorf("Expected view 'belongs-to-many-field', got '%s'", field.View)
	}

	if field.LoadingStrategy != EAGER_LOADING {
		t.Errorf("Expected default loading strategy 'eager', got '%s'", field.LoadingStrategy)
	}
}

// TestBelongsToManyPivotTable tests PivotTable method
func TestBelongsToManyPivotTable(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	result := field.PivotTable("post_tag")

	if result != field {
		t.Error("PivotTable should return the field for chaining")
	}

	if field.PivotTableName != "post_tag" {
		t.Errorf("Expected pivot table 'post_tag', got '%s'", field.PivotTableName)
	}
}

// TestBelongsToManyForeignKey tests ForeignKey method
func TestBelongsToManyForeignKey(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	result := field.ForeignKey("post_id")

	if result != field {
		t.Error("ForeignKey should return the field for chaining")
	}

	if field.ForeignKeyColumn != "post_id" {
		t.Errorf("Expected foreign key 'post_id', got '%s'", field.ForeignKeyColumn)
	}
}

// TestBelongsToManyRelatedKey tests RelatedKey method
func TestBelongsToManyRelatedKey(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	result := field.RelatedKey("tag_id")

	if result != field {
		t.Error("RelatedKey should return the field for chaining")
	}

	if field.RelatedKeyColumn != "tag_id" {
		t.Errorf("Expected related key 'tag_id', got '%s'", field.RelatedKeyColumn)
	}
}

// TestBelongsToManyQuery tests Query method
func TestBelongsToManyQuery(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")

	queryFn := func(q interface{}) interface{} {
		return q
	}

	result := field.Query(queryFn)

	if result != field {
		t.Error("Query should return the field for chaining")
	}

	if field.QueryCallback == nil {
		t.Error("Query callback should not be nil")
	}
}

// TestBelongsToManyLoadingStrategies tests loading strategy methods
func TestBelongsToManyLoadingStrategies(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")

	// Test WithEagerLoad
	result := field.WithEagerLoad()
	if result != field {
		t.Error("WithEagerLoad should return the field for chaining")
	}
	if field.LoadingStrategy != EAGER_LOADING {
		t.Errorf("Expected loading strategy 'eager', got '%s'", field.LoadingStrategy)
	}

	// Test WithLazyLoad
	result = field.WithLazyLoad()
	if result != field {
		t.Error("WithLazyLoad should return the field for chaining")
	}
	if field.LoadingStrategy != LAZY_LOADING {
		t.Errorf("Expected loading strategy 'lazy', got '%s'", field.LoadingStrategy)
	}
}

// TestBelongsToManyGetRelationshipType tests GetRelationshipType method
func TestBelongsToManyGetRelationshipType(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	relType := field.GetRelationshipType()

	if relType != "belongsToMany" {
		t.Errorf("Expected relationship type 'belongsToMany', got '%s'", relType)
	}
}

// TestBelongsToManyGetRelatedResource tests GetRelatedResource method
func TestBelongsToManyGetRelatedResource(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	resource := field.GetRelatedResource()

	if resource != "tags" {
		t.Errorf("Expected related resource 'tags', got '%s'", resource)
	}
}

// TestBelongsToManyGetRelationshipName tests GetRelationshipName method
func TestBelongsToManyGetRelationshipName(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	name := field.GetRelationshipName()

	if name != "Tags" {
		t.Errorf("Expected relationship name 'Tags', got '%s'", name)
	}
}

// TestBelongsToManyResolveRelationship tests ResolveRelationship method
func TestBelongsToManyResolveRelationship(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")

	// Test with nil item
	result, err := field.ResolveRelationship(nil)
	if err != nil {
		t.Errorf("Expected no error for nil item, got %v", err)
	}
	if result == nil {
		t.Error("Expected empty slice, got nil")
	}

	// Test with valid item
	post := &Post{ID: 1, UserID: 1, Title: "Test"}
	result, err = field.ResolveRelationship(post)
	if err != nil {
		t.Errorf("Expected no error for valid item, got %v", err)
	}
	if result == nil {
		t.Error("Expected empty slice, got nil")
	}
}

// TestBelongsToManyValidateRelationship tests ValidateRelationship method
func TestBelongsToManyValidateRelationship(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")

	// Test with nil value
	err := field.ValidateRelationship(nil)
	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}

	// Test with valid value
	err = field.ValidateRelationship([]interface{}{})
	if err != nil {
		t.Errorf("Expected no error for valid value, got %v", err)
	}
}

// TestBelongsToManyGetQueryCallback tests GetQueryCallback method
func TestBelongsToManyGetQueryCallback(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")

	// Test with nil callback
	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test callback returns query unchanged
	query := "SELECT * FROM tags"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestBelongsToManyGetLoadingStrategy tests GetLoadingStrategy method
func TestBelongsToManyGetLoadingStrategy(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	strategy := field.GetLoadingStrategy()

	if strategy != EAGER_LOADING {
		t.Errorf("Expected loading strategy 'eager', got '%s'", strategy)
	}

	field.WithLazyLoad()
	strategy = field.GetLoadingStrategy()

	if strategy != LAZY_LOADING {
		t.Errorf("Expected loading strategy 'lazy', got '%s'", strategy)
	}
}

// TestBelongsToManyFluentChaining tests fluent API chaining
func TestBelongsToManyFluentChaining(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags").
		PivotTable("post_tag").
		ForeignKey("post_id").
		RelatedKey("tag_id").
		WithLazyLoad()

	if field.PivotTableName != "post_tag" {
		t.Errorf("Expected pivot table 'post_tag', got '%s'", field.PivotTableName)
	}

	if field.ForeignKeyColumn != "post_id" {
		t.Errorf("Expected foreign key 'post_id', got '%s'", field.ForeignKeyColumn)
	}

	if field.RelatedKeyColumn != "tag_id" {
		t.Errorf("Expected related key 'tag_id', got '%s'", field.RelatedKeyColumn)
	}

	if field.LoadingStrategy != LAZY_LOADING {
		t.Errorf("Expected loading strategy 'lazy', got '%s'", field.LoadingStrategy)
	}
}

// TestBelongsToManyCount tests Count method
func TestBelongsToManyCount(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	count := field.Count()

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestBelongsToManyNilQueryCallback tests nil query callback handling
func TestBelongsToManyNilQueryCallback(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	field.QueryCallback = nil

	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test that default callback returns query unchanged
	query := "SELECT * FROM tags"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestBelongsToManyImplementsRelationshipField tests that BelongsToMany implements RelationshipField interface
func TestBelongsToManyImplementsRelationshipField(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")

	// Test relationship field methods
	if field.GetRelationshipType() != "belongsToMany" {
		t.Errorf("Expected relationship type 'belongsToMany', got '%s'", field.GetRelationshipType())
	}

	if field.GetRelatedResource() != "tags" {
		t.Errorf("Expected related resource 'tags', got '%s'", field.GetRelatedResource())
	}

	if field.GetRelationshipName() != "Tags" {
		t.Errorf("Expected relationship name 'Tags', got '%s'", field.GetRelationshipName())
	}
}

// TestBelongsToManyPivotTableNaming tests pivot table naming convention
func TestBelongsToManyPivotTableNaming(t *testing.T) {
	field := BelongsToMany("Tags", "posts", "tags")

	// Pivot table should be alphabetically sorted
	if field.PivotTableName != "posts_tags" {
		t.Errorf("Expected pivot table 'posts_tags', got '%s'", field.PivotTableName)
	}
}
