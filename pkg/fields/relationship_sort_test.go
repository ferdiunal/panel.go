package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipSort tests creating a new relationship sort handler
func TestNewRelationshipSort(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	if sort == nil {
		t.Error("Expected non-nil sort handler")
	}

	if sort.field != field {
		t.Error("Expected field to be set")
	}

	if len(sort.sorts) != 0 {
		t.Error("Expected empty sorts")
	}
}

// TestApplySort tests applying a sort
func TestApplySort(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	results, err := sort.ApplySort(ctx, "name", "ASC")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(sort.sorts) != 1 {
		t.Errorf("Expected 1 sort, got %d", len(sort.sorts))
	}

	if sort.sorts["name"] != "ASC" {
		t.Errorf("Expected 'ASC', got '%s'", sort.sorts["name"])
	}
}

// TestApplySortEmptyColumn tests applying a sort with empty column
func TestApplySortEmptyColumn(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	_, err := sort.ApplySort(ctx, "", "ASC")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(sort.sorts) != 0 {
		t.Errorf("Expected 0 sorts, got %d", len(sort.sorts))
	}
}

// TestApplySortDescending tests applying a descending sort
func TestApplySortDescending(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	_, err := sort.ApplySort(ctx, "name", "DESC")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sort.sorts["name"] != "DESC" {
		t.Errorf("Expected 'DESC', got '%s'", sort.sorts["name"])
	}
}

// TestApplySortInvalidDirection tests applying a sort with invalid direction
func TestApplySortInvalidDirection(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	_, err := sort.ApplySort(ctx, "name", "INVALID")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sort.sorts["name"] != "ASC" {
		t.Errorf("Expected 'ASC' (default), got '%s'", sort.sorts["name"])
	}
}

// TestApplyMultipleSorts tests applying multiple sorts
func TestApplyMultipleSorts(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	sorts := map[string]string{
		"name":  "ASC",
		"email": "DESC",
	}

	results, err := sort.ApplyMultipleSorts(ctx, sorts)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(sort.sorts) != 2 {
		t.Errorf("Expected 2 sorts, got %d", len(sort.sorts))
	}
}

// TestApplyMultipleSortsEmpty tests applying multiple sorts with empty sorts
func TestApplyMultipleSortsEmpty(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	results, err := sort.ApplyMultipleSorts(ctx, map[string]string{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestRemoveSort tests removing sorts
func TestRemoveSort(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	sort.ApplySort(ctx, "name", "ASC")

	if len(sort.sorts) != 1 {
		t.Errorf("Expected 1 sort before removal, got %d", len(sort.sorts))
	}

	results, err := sort.RemoveSort(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(sort.sorts) != 0 {
		t.Errorf("Expected 0 sorts after removal, got %d", len(sort.sorts))
	}
}

// TestGetSorts tests getting sorts
func TestGetSorts(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	sort := NewRelationshipSort(field)

	ctx := context.Background()
	sort.ApplySort(ctx, "name", "ASC")

	sorts := sort.GetSorts()

	if len(sorts) != 1 {
		t.Errorf("Expected 1 sort, got %d", len(sorts))
	}

	if sorts["name"] != "ASC" {
		t.Errorf("Expected 'ASC', got '%s'", sorts["name"])
	}
}
