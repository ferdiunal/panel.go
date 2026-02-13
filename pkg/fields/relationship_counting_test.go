package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipCounting tests creating a new relationship counting handler
func TestNewRelationshipCounting(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	counting := NewRelationshipCounting(field)

	if counting == nil {
		t.Error("Expected non-nil counting handler")
	}

	if counting.field != field {
		t.Error("Expected field to be set")
	}
}

// TestCountBelongsTo tests counting BelongsTo relationships
func TestCountBelongsTo(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	counting := NewRelationshipCounting(field)

	ctx := context.Background()
	count, err := counting.Count(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestCountHasMany tests counting HasMany relationships
func TestCountHasMany(t *testing.T) {
	field := HasMany("Posts", "posts", "posts")
	counting := NewRelationshipCounting(field)

	ctx := context.Background()
	count, err := counting.Count(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestCountHasOne tests counting HasOne relationships
func TestCountHasOne(t *testing.T) {
	field := HasOne("Profile", "profile", "profiles")
	counting := NewRelationshipCounting(field)

	ctx := context.Background()
	count, err := counting.Count(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestCountBelongsToMany tests counting BelongsToMany relationships
func TestCountBelongsToMany(t *testing.T) {
	field := BelongsToMany("Tags", "tags", "tags")
	counting := NewRelationshipCounting(field)

	ctx := context.Background()
	count, err := counting.Count(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestCountMorphTo tests counting MorphTo relationships
func TestCountNewMorphTo(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	counting := NewRelationshipCounting(field)

	ctx := context.Background()
	count, err := counting.Count(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}
