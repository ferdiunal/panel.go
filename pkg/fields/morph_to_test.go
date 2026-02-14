package fields

import (
	"testing"
)

// TestMorphToCreation tests MorphTo field creation
func TestMorphToCreation(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	if field.Name != "Commentable" {
		t.Errorf("Expected name 'Commentable', got '%s'", field.Name)
	}

	if field.Key != "commentable" {
		t.Errorf("Expected key 'commentable', got '%s'", field.Key)
	}

	if field.View != "morph-to-field" {
		t.Errorf("Expected view 'morph-to-field', got '%s'", field.View)
	}

	if field.LoadingStrategy != EAGER_LOADING {
		t.Errorf("Expected default loading strategy 'eager', got '%s'", field.LoadingStrategy)
	}

	if len(field.TypeMappings) != 0 {
		t.Errorf("Expected empty type mappings, got %v", field.TypeMappings)
	}
}

// TestMorphToTypes tests Types method
func TestMorphToTypes(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	typeMappings := map[string]string{
		"post":  "posts",
		"video": "videos",
	}

	result := field.Types(typeMappings)

	if result != field {
		t.Error("Types should return the field for chaining")
	}

	if len(field.TypeMappings) != 2 {
		t.Errorf("Expected 2 type mappings, got %d", len(field.TypeMappings))
	}

	if field.TypeMappings["post"] != "posts" {
		t.Errorf("Expected 'post' -> 'posts', got '%s'", field.TypeMappings["post"])
	}

	if field.TypeMappings["video"] != "videos" {
		t.Errorf("Expected 'video' -> 'videos', got '%s'", field.TypeMappings["video"])
	}
}

// TestMorphToQuery tests Query method
func TestMorphToQuery(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

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

// TestMorphToLoadingStrategies tests loading strategy methods
func TestMorphToLoadingStrategies(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

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

// TestMorphToGetRelationshipType tests GetRelationshipType method
func TestMorphToGetRelationshipType(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	relType := field.GetRelationshipType()

	if relType != "morphTo" {
		t.Errorf("Expected relationship type 'morphTo', got '%s'", relType)
	}
}

// TestMorphToGetRelatedResource tests GetRelatedResourceSlug method
func TestMorphToGetRelatedResource(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	resource := field.GetRelatedResourceSlug()

	if resource != "" {
		t.Errorf("Expected empty related resource for MorphTo, got '%s'", resource)
	}
}

// TestMorphToGetRelationshipName tests GetRelationshipName method
func TestMorphToGetRelationshipName(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	name := field.GetRelationshipName()

	if name != "Commentable" {
		t.Errorf("Expected relationship name 'Commentable', got '%s'", name)
	}
}

// TestMorphToResolveRelationship tests ResolveRelationship method
func TestMorphToResolveRelationship(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	// Test with nil item
	result, err := field.ResolveRelationship(nil)
	if err != nil {
		t.Errorf("Expected no error for nil item, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for nil item, got %v", result)
	}

	// Test with valid item
	comment := &Comment{ID: 1, PostID: 1, UserID: 1, Content: "Test"}
	result, err = field.ResolveRelationship(comment)
	if err != nil {
		t.Errorf("Expected no error for valid item, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for valid item, got %v", result)
	}
}

// TestMorphToValidateRelationship tests ValidateRelationship method
func TestMorphToValidateRelationship(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	// Test with nil value
	err := field.ValidateRelationship(nil)
	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}

	// Test with valid value
	err = field.ValidateRelationship(&Taggable{ID: 1, TaggableID: 1, TaggableType: "post"})
	if err != nil {
		t.Errorf("Expected no error for valid value, got %v", err)
	}
}

// TestMorphToGetQueryCallback tests GetQueryCallback method
func TestMorphToGetQueryCallback(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	// Test with nil callback
	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test callback returns query unchanged
	query := "SELECT * FROM comments"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestMorphToGetLoadingStrategy tests GetLoadingStrategy method
func TestMorphToGetLoadingStrategy(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
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

// TestMorphToGetTypes tests GetTypes method
func TestMorphToGetTypes(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	typeMappings := map[string]string{
		"post":  "posts",
		"video": "videos",
	}

	field.Types(typeMappings)
	types := field.GetTypes()

	if len(types) != 2 {
		t.Errorf("Expected 2 type mappings, got %d", len(types))
	}

	if types["post"] != "posts" {
		t.Errorf("Expected 'post' -> 'posts', got '%s'", types["post"])
	}

	if types["video"] != "videos" {
		t.Errorf("Expected 'video' -> 'videos', got '%s'", types["video"])
	}
}

// TestMorphToGetResourceForType tests GetResourceForType method
func TestMorphToGetResourceForType(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	typeMappings := map[string]string{
		"post":  "posts",
		"video": "videos",
	}

	field.Types(typeMappings)

	// Test valid type
	resource, err := field.GetResourceForType("post")
	if err != nil {
		t.Errorf("Expected no error for valid type, got %v", err)
	}
	if resource != "posts" {
		t.Errorf("Expected resource 'posts', got '%s'", resource)
	}

	// Test invalid type
	resource, err = field.GetResourceForType("invalid")
	if err == nil {
		t.Error("Expected error for invalid type")
	}
	if resource != "" {
		t.Errorf("Expected empty resource for invalid type, got '%s'", resource)
	}
}

// TestMorphToFluentChaining tests fluent API chaining
func TestMorphToFluentChaining(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable").
		Types(map[string]string{
			"post":  "posts",
			"video": "videos",
		}).
		WithLazyLoad()

	if len(field.TypeMappings) != 2 {
		t.Errorf("Expected 2 type mappings, got %d", len(field.TypeMappings))
	}

	if field.LoadingStrategy != LAZY_LOADING {
		t.Errorf("Expected loading strategy 'lazy', got '%s'", field.LoadingStrategy)
	}
}

// TestMorphToNilQueryCallback tests nil query callback handling
func TestMorphToNilQueryCallback(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	field.QueryCallback = nil

	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test that default callback returns query unchanged
	query := "SELECT * FROM comments"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestMorphToNilTypeMappings tests nil type mappings handling
func TestMorphToNilTypeMappings(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	field.TypeMappings = nil

	types := field.GetTypes()
	if types == nil {
		t.Error("Expected empty map, got nil")
	}

	if len(types) != 0 {
		t.Errorf("Expected empty map, got %v", types)
	}
}

// TestMorphToImplementsRelationshipField tests that MorphTo implements RelationshipField interface
func TestMorphToImplementsRelationshipField(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")

	// Test relationship field methods
	if field.GetRelationshipType() != "morphTo" {
		t.Errorf("Expected relationship type 'morphTo', got '%s'", field.GetRelationshipType())
	}

	if field.GetRelationshipName() != "Commentable" {
		t.Errorf("Expected relationship name 'Commentable', got '%s'", field.GetRelationshipName())
	}
}
