package fields

import (
	"testing"
)

// TestHasOneCreation tests HasOne field creation
func TestHasOneCreation(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")

	if field.Name != "Profile" {
		t.Errorf("Expected name 'Profile', got '%s'", field.Name)
	}

	if field.Key != "profile" {
		t.Errorf("Expected key 'profile', got '%s'", field.Key)
	}

	if field.RelatedResourceSlug != "profiles" {
		t.Errorf("Expected related resource 'profiles', got '%s'", field.RelatedResourceSlug)
	}

	if field.View != "has-one-field" {
		t.Errorf("Expected view 'has-one-field', got '%s'", field.View)
	}

	if field.LoadingStrategy != EAGER_LOADING {
		t.Errorf("Expected default loading strategy 'eager', got '%s'", field.LoadingStrategy)
	}
}

// TestHasOneForeignKey tests ForeignKey method
func TestHasOneForeignKey(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
	result := field.ForeignKey("user_id")

	if result != field {
		t.Error("ForeignKey should return the field for chaining")
	}

	if field.ForeignKeyColumn != "user_id" {
		t.Errorf("Expected foreign key 'user_id', got '%s'", field.ForeignKeyColumn)
	}
}

// TestHasOneOwnerKey tests OwnerKey method
func TestHasOneOwnerKey(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
	result := field.OwnerKey("id")

	if result != field {
		t.Error("OwnerKey should return the field for chaining")
	}

	if field.OwnerKeyColumn != "id" {
		t.Errorf("Expected owner key 'id', got '%s'", field.OwnerKeyColumn)
	}
}

// TestHasOneQuery tests Query method
func TestHasOneQuery(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")

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

// TestHasOneLoadingStrategies tests loading strategy methods
func TestHasOneLoadingStrategies(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")

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

// TestHasOneGetRelationshipType tests GetRelationshipType method
func TestHasOneGetRelationshipType(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
	relType := field.GetRelationshipType()

	if relType != "hasOne" {
		t.Errorf("Expected relationship type 'hasOne', got '%s'", relType)
	}
}

// TestHasOneGetRelatedResource tests GetRelatedResource method
func TestHasOneGetRelatedResource(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
	resource := field.GetRelatedResource()

	if resource != "profiles" {
		t.Errorf("Expected related resource 'profiles', got '%s'", resource)
	}
}

// TestHasOneGetRelationshipName tests GetRelationshipName method
func TestHasOneGetRelationshipName(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
	name := field.GetRelationshipName()

	if name != "Profile" {
		t.Errorf("Expected relationship name 'Profile', got '%s'", name)
	}
}

// TestHasOneResolveRelationship tests ResolveRelationship method
func TestHasOneResolveRelationship(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")

	// Test with nil item
	result, err := field.ResolveRelationship(nil)
	if err != nil {
		t.Errorf("Expected no error for nil item, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for nil item, got %v", result)
	}

	// Test with valid item
	user := &User{ID: 1, Name: "John", Email: "john@example.com"}
	result, err = field.ResolveRelationship(user)
	if err != nil {
		t.Errorf("Expected no error for valid item, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for valid item, got %v", result)
	}
}

// TestHasOneValidateRelationship tests ValidateRelationship method
func TestHasOneValidateRelationship(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")

	// Test with nil value
	err := field.ValidateRelationship(nil)
	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}

	// Test with valid value
	err = field.ValidateRelationship(&Profile{ID: 1, UserID: 1})
	if err != nil {
		t.Errorf("Expected no error for valid value, got %v", err)
	}
}

// TestHasOneGetQueryCallback tests GetQueryCallback method
func TestHasOneGetQueryCallback(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")

	// Test with nil callback
	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test callback returns query unchanged
	query := "SELECT * FROM profiles"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestHasOneGetLoadingStrategy tests GetLoadingStrategy method
func TestHasOneGetLoadingStrategy(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
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

// TestHasOneFluentChaining tests fluent API chaining
func TestHasOneFluentChaining(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles").
		ForeignKey("user_id").
		OwnerKey("id").
		WithLazyLoad()

	if field.ForeignKeyColumn != "user_id" {
		t.Errorf("Expected foreign key 'user_id', got '%s'", field.ForeignKeyColumn)
	}

	if field.OwnerKeyColumn != "id" {
		t.Errorf("Expected owner key 'id', got '%s'", field.OwnerKeyColumn)
	}

	if field.LoadingStrategy != LAZY_LOADING {
		t.Errorf("Expected loading strategy 'lazy', got '%s'", field.LoadingStrategy)
	}
}

// TestHasOneNilQueryCallback tests nil query callback handling
func TestHasOneNilQueryCallback(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
	field.QueryCallback = nil

	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test that default callback returns query unchanged
	query := "SELECT * FROM profiles"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}

// TestHasOneImplementsRelationshipField tests that HasOne implements RelationshipField interface
func TestHasOneImplementsRelationshipField(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")

	// Test relationship field methods
	if field.GetRelationshipType() != "hasOne" {
		t.Errorf("Expected relationship type 'hasOne', got '%s'", field.GetRelationshipType())
	}

	if field.GetRelatedResource() != "profiles" {
		t.Errorf("Expected related resource 'profiles', got '%s'", field.GetRelatedResource())
	}

	if field.GetRelationshipName() != "Profile" {
		t.Errorf("Expected relationship name 'Profile', got '%s'", field.GetRelationshipName())
	}
}
