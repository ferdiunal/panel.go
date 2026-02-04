package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipSearch tests creating a new relationship search handler
func TestNewRelationshipSearch(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	search := NewRelationshipSearch(field)

	if search == nil {
		t.Error("Expected non-nil search handler")
	}

	if search.field != field {
		t.Error("Expected field to be set")
	}
}

// TestSearchEmptyTerm tests search with empty term
func TestSearchEmptyTerm(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	search := NewRelationshipSearch(field)

	ctx := context.Background()
	results, err := search.Search(ctx, "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestSearchWithTerm tests search with term
func TestSearchWithTerm(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	field.WithSearchableColumns("name", "email")
	search := NewRelationshipSearch(field)

	ctx := context.Background()
	results, err := search.Search(ctx, "John")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}
}

// TestSearchInColumns tests search in specific columns
func TestSearchInColumns(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	search := NewRelationshipSearch(field)

	ctx := context.Background()
	results, err := search.SearchInColumns(ctx, "John", []string{"name", "email"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}
}

// TestSearchInColumnsEmpty tests search in columns with empty columns
func TestSearchInColumnsEmpty(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	search := NewRelationshipSearch(field)

	ctx := context.Background()
	results, err := search.SearchInColumns(ctx, "John", []string{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestGetSearchableColumns tests getting searchable columns
func TestGetSearchableColumns(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	field.WithSearchableColumns("name", "email")
	search := NewRelationshipSearch(field)

	columns := search.GetSearchableColumns()

	if len(columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(columns))
	}

	if columns[0] != "name" || columns[1] != "email" {
		t.Errorf("Expected ['name', 'email'], got %v", columns)
	}
}

// TestGetSearchableColumnsHasMany tests getting searchable columns for HasMany
func TestGetSearchableColumnsHasMany(t *testing.T) {
	field := NewHasMany("Posts", "posts", "posts")
	search := NewRelationshipSearch(field)

	columns := search.GetSearchableColumns()

	if len(columns) != 0 {
		t.Errorf("Expected 0 columns for HasMany, got %d", len(columns))
	}
}

// TestCaseInsensitiveSearch tests case-insensitive search
func TestCaseInsensitiveSearch(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	search := NewRelationshipSearch(field)

	ctx := context.Background()
	results, err := search.CaseInsensitiveSearch(ctx, "JOHN")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}
}
