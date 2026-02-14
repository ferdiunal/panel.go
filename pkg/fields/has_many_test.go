package fields

import (
	"fmt"
	"testing"
)

type hasManyRecordTitleResource struct {
	slug string
}

func (r hasManyRecordTitleResource) Slug() string {
	return r.slug
}

func (r hasManyRecordTitleResource) RecordTitle(record any) string {
	return fmt.Sprintf("record-%v", record)
}

// TestHasManyCreation tests HasMany field creation
func TestHasManyCreation(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")

	if field.Name != "Posts" {
		t.Errorf("Expected name 'Posts', got '%s'", field.Name)
	}

	if field.Key != "posts" {
		t.Errorf("Expected key 'posts', got '%s'", field.Key)
	}

	if field.RelatedResourceSlug != "posts" {
		t.Errorf("Expected related resource 'posts', got '%s'", field.RelatedResourceSlug)
	}

	if field.View != "has-many-field" {
		t.Errorf("Expected view 'has-many-field', got '%s'", field.View)
	}

	if field.LoadingStrategy != EAGER_LOADING {
		t.Errorf("Expected default loading strategy 'eager', got '%s'", field.LoadingStrategy)
	}
}

// TestHasManyForeignKey tests ForeignKey method
func TestHasManyForeignKey(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	result := field.ForeignKey("author_id")

	if result != field {
		t.Error("ForeignKey should return the field for chaining")
	}

	if field.ForeignKeyColumn != "author_id" {
		t.Errorf("Expected foreign key 'author_id', got '%s'", field.ForeignKeyColumn)
	}
}

// TestHasManyOwnerKey tests OwnerKey method
func TestHasManyOwnerKey(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	result := field.OwnerKey("user_id")

	if result != field {
		t.Error("OwnerKey should return the field for chaining")
	}

	if field.OwnerKeyColumn != "user_id" {
		t.Errorf("Expected owner key 'user_id', got '%s'", field.OwnerKeyColumn)
	}
}

// TestHasManyQuery tests Query method
func TestHasManyQuery(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")

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

// TestHasManyLoadingStrategies tests loading strategy methods
func TestHasManyLoadingStrategies(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")

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

// TestHasManyGetRelationshipType tests GetRelationshipType method
func TestHasManyGetRelationshipType(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	relType := field.GetRelationshipType()

	if relType != "hasMany" {
		t.Errorf("Expected relationship type 'hasMany', got '%s'", relType)
	}
}

// TestHasManyGetRelatedResource tests GetRelatedResourceSlug method
func TestHasManyGetRelatedResource(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	resource := field.GetRelatedResourceSlug()

	if resource != "posts" {
		t.Errorf("Expected related resource 'posts', got '%s'", resource)
	}
}

// TestHasManyGetRelationshipName tests GetRelationshipName method
func TestHasManyGetRelationshipName(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	name := field.GetRelationshipName()

	if name != "Posts" {
		t.Errorf("Expected relationship name 'Posts', got '%s'", name)
	}
}

// TestHasManyResolveRelationship tests ResolveRelationship method
func TestHasManyResolveRelationship(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")

	// Test with nil item
	result, err := field.ResolveRelationship(nil)
	if err != nil {
		t.Errorf("Expected no error for nil item, got %v", err)
	}
	if result == nil {
		t.Error("Expected empty slice, got nil")
	}

	// Test with valid item
	user := &User{ID: 1, Name: "John", Email: "john@example.com"}
	result, err = field.ResolveRelationship(user)
	if err != nil {
		t.Errorf("Expected no error for valid item, got %v", err)
	}
	if result == nil {
		t.Error("Expected empty slice, got nil")
	}
}

// TestHasManyValidateRelationship tests ValidateRelationship method
func TestHasManyValidateRelationship(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")

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

// TestHasManyGetQueryCallback tests GetQueryCallback method
func TestHasManyGetQueryCallback(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")

	// Test with nil callback
	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test callback returns query unchanged
	query := "SELECT * FROM posts"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestHasManyGetLoadingStrategy tests GetLoadingStrategy method
func TestHasManyGetLoadingStrategy(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
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

// TestHasManyFluentChaining tests fluent API chaining
func TestHasManyFluentChaining(t *testing.T) {
	field := HasMany("Posts", "posts", "posts").
		ForeignKey("author_id").
		OwnerKey("id").
		WithLazyLoad()

	if field.ForeignKeyColumn != "author_id" {
		t.Errorf("Expected foreign key 'author_id', got '%s'", field.ForeignKeyColumn)
	}

	if field.OwnerKeyColumn != "id" {
		t.Errorf("Expected owner key 'id', got '%s'", field.OwnerKeyColumn)
	}

	if field.LoadingStrategy != LAZY_LOADING {
		t.Errorf("Expected loading strategy 'lazy', got '%s'", field.LoadingStrategy)
	}
}

// TestHasManyCount tests Count method
func TestHasManyCount(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	count := field.Count()

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestHasManyNilQueryCallback tests nil query callback handling
func TestHasManyNilQueryCallback(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	field.QueryCallback = nil

	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test that default callback returns query unchanged
	query := "SELECT * FROM posts"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestHasManyImplementsRelationshipField tests that HasMany implements RelationshipField interface
func TestHasManyImplementsRelationshipField(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")

	// Test relationship field methods
	if field.GetRelationshipType() != "hasMany" {
		t.Errorf("Expected relationship type 'hasMany', got '%s'", field.GetRelationshipType())
	}

	if field.GetRelatedResourceSlug() != "posts" {
		t.Errorf("Expected related resource 'posts', got '%s'", field.GetRelatedResourceSlug())
	}

	if field.GetRelationshipName() != "Posts" {
		t.Errorf("Expected relationship name 'Posts', got '%s'", field.GetRelationshipName())
	}
}

func TestHasManyLabelPlaceholderChainingKeepsConcreteType(t *testing.T) {
	element := HasMany(
		"ProductVariants",
		"product_variants",
		hasManyRecordTitleResource{slug: "product-variants"},
	).
		Label("Urun Varyantlari").
		Placeholder("Urun varyanti secin").
		Searchable()

	hasManyField, ok := element.(*HasManyField)
	if !ok {
		t.Fatalf("expected chained element to remain *HasManyField, got %T", element)
	}

	if hasManyField.GetRelatedResourceSlug() != "product-variants" {
		t.Fatalf("expected related resource slug to be preserved, got %s", hasManyField.GetRelatedResourceSlug())
	}

	if relField, ok := IsRelationshipField(element); !ok || relField == nil {
		t.Fatalf("expected chained element to be detected as concrete relationship field")
	}
}
