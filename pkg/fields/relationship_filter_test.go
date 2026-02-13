package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipFilter tests creating a new relationship filter handler
func TestNewRelationshipFilter(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	filter := NewRelationshipFilter(field)

	if filter == nil {
		t.Error("Expected non-nil filter handler")
	}

	if filter.field != field {
		t.Error("Expected field to be set")
	}

	if len(filter.filters) != 0 {
		t.Error("Expected empty filters")
	}
}

// TestApplyFilter tests applying a filter
func TestApplyFilter(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	filter := NewRelationshipFilter(field)

	ctx := context.Background()
	results, err := filter.ApplyFilter(ctx, "status", "=", "active")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(filter.filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(filter.filters))
	}
}

// TestApplyFilterEmptyColumn tests applying a filter with empty column
func TestApplyFilterEmptyColumn(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	filter := NewRelationshipFilter(field)

	ctx := context.Background()
	results, err := filter.ApplyFilter(ctx, "", "=", "active")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestApplyMultipleFilters tests applying multiple filters
func TestApplyMultipleFilters(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	filter := NewRelationshipFilter(field)

	ctx := context.Background()
	filters := map[string]interface{}{
		"status": map[string]interface{}{
			"operator": "=",
			"value":    "active",
		},
		"role": map[string]interface{}{
			"operator": "=",
			"value":    "admin",
		},
	}

	results, err := filter.ApplyMultipleFilters(ctx, filters)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(filter.filters) != 2 {
		t.Errorf("Expected 2 filters, got %d", len(filter.filters))
	}
}

// TestApplyMultipleFiltersEmpty tests applying multiple filters with empty filters
func TestApplyMultipleFiltersEmpty(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	filter := NewRelationshipFilter(field)

	ctx := context.Background()
	results, err := filter.ApplyMultipleFilters(ctx, map[string]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestRemoveFilter tests removing filters
func TestRemoveFilter(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	filter := NewRelationshipFilter(field)

	ctx := context.Background()
	filter.ApplyFilter(ctx, "status", "=", "active")

	if len(filter.filters) != 1 {
		t.Errorf("Expected 1 filter before removal, got %d", len(filter.filters))
	}

	results, err := filter.RemoveFilter(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(filter.filters) != 0 {
		t.Errorf("Expected 0 filters after removal, got %d", len(filter.filters))
	}
}

// TestGetFilters tests getting filters
func TestGetFilters(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	filter := NewRelationshipFilter(field)

	ctx := context.Background()
	filter.ApplyFilter(ctx, "status", "=", "active")

	filters := filter.GetFilters()

	if len(filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(filters))
	}

	if _, ok := filters["status"]; !ok {
		t.Error("Expected 'status' filter to exist")
	}
}
