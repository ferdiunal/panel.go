package fields

import (
	"testing"
)

// TestBelongsToCreation tests BelongsTo field creation
func TestBelongsToCreation(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")

	if field.Name != "Author" {
		t.Errorf("Expected name 'Author', got '%s'", field.Name)
	}

	if field.Key != "author_id" {
		t.Errorf("Expected key 'author_id', got '%s'", field.Key)
	}

	if field.RelatedResourceSlug != "authors" {
		t.Errorf("Expected related resource 'authors', got '%s'", field.RelatedResourceSlug)
	}

	if field.View != "belongs-to-field" {
		t.Errorf("Expected view 'belongs-to-field', got '%s'", field.View)
	}

	if field.DisplayKey != "name" {
		t.Errorf("Expected default display key 'name', got '%s'", field.DisplayKey)
	}

	if field.LoadingStrategy != EAGER_LOADING {
		t.Errorf("Expected default loading strategy 'eager', got '%s'", field.LoadingStrategy)
	}
}

// TestBelongsToDisplayUsing tests DisplayUsing method
func TestBelongsToDisplayUsing(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	result := field.DisplayUsing("email")

	if result != field {
		t.Error("DisplayUsing should return the field for chaining")
	}

	if field.DisplayKey != "email" {
		t.Errorf("Expected display key 'email', got '%s'", field.DisplayKey)
	}
}

// TestBelongsToSearchable tests Searchable method
func TestBelongsToSearchable(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	result := field.WithSearchableColumns("name", "email")

	if result != field {
		t.Error("WithSearchableColumns should return the field for chaining")
	}

	if len(field.SearchableColumns) != 2 {
		t.Errorf("Expected 2 searchable columns, got %d", len(field.SearchableColumns))
	}

	if field.SearchableColumns[0] != "name" || field.SearchableColumns[1] != "email" {
		t.Errorf("Expected searchable columns ['name', 'email'], got %v", field.SearchableColumns)
	}
}

// TestBelongsToQuery tests Query method
func TestBelongsToQuery(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")

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

// TestBelongsToLoadingStrategies tests loading strategy methods
func TestBelongsToLoadingStrategies(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")

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

// TestBelongsToGetRelationshipType tests GetRelationshipType method
func TestBelongsToGetRelationshipType(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	relType := field.GetRelationshipType()

	if relType != "belongsTo" {
		t.Errorf("Expected relationship type 'belongsTo', got '%s'", relType)
	}
}

// TestBelongsToGetRelatedResource tests GetRelatedResource method
func TestBelongsToGetRelatedResource(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	resource := field.GetRelatedResource()

	if resource != "authors" {
		t.Errorf("Expected related resource 'authors', got '%s'", resource)
	}
}

// TestBelongsToGetRelationshipName tests GetRelationshipName method
func TestBelongsToGetRelationshipName(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	name := field.GetRelationshipName()

	if name != "Author" {
		t.Errorf("Expected relationship name 'Author', got '%s'", name)
	}
}

// TestBelongsToGetDisplayKey tests GetDisplayKey method
func TestBelongsToGetDisplayKey(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	displayKey := field.GetDisplayKey()

	if displayKey != "name" {
		t.Errorf("Expected display key 'name', got '%s'", displayKey)
	}

	field.DisplayUsing("email")
	displayKey = field.GetDisplayKey()

	if displayKey != "email" {
		t.Errorf("Expected display key 'email', got '%s'", displayKey)
	}
}

// TestBelongsToGetSearchableColumns tests GetSearchableColumns method
func TestBelongsToGetSearchableColumns(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	columns := field.GetSearchableColumns()

	if len(columns) != 1 || columns[0] != "name" {
		t.Errorf("Expected searchable columns ['name'], got %v", columns)
	}

	field.WithSearchableColumns("name", "email")
	columns = field.GetSearchableColumns()

	if len(columns) != 2 || columns[0] != "name" || columns[1] != "email" {
		t.Errorf("Expected searchable columns ['name', 'email'], got %v", columns)
	}
}

// TestBelongsToGetQueryCallback tests GetQueryCallback method
func TestBelongsToGetQueryCallback(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")

	// Test with nil callback
	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test callback returns query unchanged
	query := "SELECT * FROM authors"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}

	// Test with custom callback
	customCallback := func(q interface{}) interface{} {
		return "MODIFIED"
	}
	field.Query(customCallback)
	callback = field.GetQueryCallback()
	result = callback("SELECT * FROM authors")
	if result != "MODIFIED" {
		t.Errorf("Expected 'MODIFIED', got %v", result)
	}
}

// TestBelongsToGetLoadingStrategy tests GetLoadingStrategy method
func TestBelongsToGetLoadingStrategy(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
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

// TestBelongsToImplementsRelationshipField tests that BelongsTo implements RelationshipField interface
func TestBelongsToImplementsRelationshipField(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")

	// Test relationship field methods
	if field.GetRelationshipType() != "belongsTo" {
		t.Errorf("Expected relationship type 'belongsTo', got '%s'", field.GetRelationshipType())
	}

	if field.GetRelatedResource() != "authors" {
		t.Errorf("Expected related resource 'authors', got '%s'", field.GetRelatedResource())
	}

	if field.GetRelationshipName() != "Author" {
		t.Errorf("Expected relationship name 'Author', got '%s'", field.GetRelationshipName())
	}
}

// TestBelongsToNilSearchableColumns tests nil searchable columns handling
func TestBelongsToNilSearchableColumns(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	field.SearchableColumns = nil

	columns := field.GetSearchableColumns()
	if columns == nil {
		t.Error("Expected empty slice, got nil")
	}

	if len(columns) != 0 {
		t.Errorf("Expected empty slice, got %v", columns)
	}
}

// TestBelongsToNilQueryCallback tests nil query callback handling
func TestBelongsToNilQueryCallback(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	field.QueryCallback = nil

	callback := field.GetQueryCallback()
	if callback == nil {
		t.Error("Expected non-nil callback")
	}

	// Test that default callback returns query unchanged
	query := "SELECT * FROM authors"
	result := callback(query)
	if result != query {
		t.Errorf("Expected query unchanged, got %v", result)
	}
}
