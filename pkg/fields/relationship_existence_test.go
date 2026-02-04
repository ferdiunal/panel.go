package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipExistence tests creating a new relationship existence handler
func TestNewRelationshipExistence(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	existence := NewRelationshipExistence(field)

	if existence == nil {
		t.Error("Expected non-nil existence handler")
	}

	if existence.field != field {
		t.Error("Expected field to be set")
	}
}

// TestExistsBelongsTo tests checking if BelongsTo relationship exists
func TestExistsBelongsTo(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	existence := NewRelationshipExistence(field)

	ctx := context.Background()
	exists, err := existence.Exists(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if exists != false {
		t.Errorf("Expected exists false, got %v", exists)
	}
}

// TestExistsHasMany tests checking if HasMany relationships exist
func TestExistsHasMany(t *testing.T) {
	field := NewHasMany("Posts", "posts", "posts")
	existence := NewRelationshipExistence(field)

	ctx := context.Background()
	exists, err := existence.Exists(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if exists != false {
		t.Errorf("Expected exists false, got %v", exists)
	}
}

// TestExistsHasOne tests checking if HasOne relationship exists
func TestExistsHasOne(t *testing.T) {
	field := NewHasOne("Profile", "profile", "profiles")
	existence := NewRelationshipExistence(field)

	ctx := context.Background()
	exists, err := existence.Exists(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if exists != false {
		t.Errorf("Expected exists false, got %v", exists)
	}
}

// TestExistsBelongsToMany tests checking if BelongsToMany relationships exist
func TestExistsBelongsToMany(t *testing.T) {
	field := NewBelongsToMany("Tags", "tags", "tags")
	existence := NewRelationshipExistence(field)

	ctx := context.Background()
	exists, err := existence.Exists(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if exists != false {
		t.Errorf("Expected exists false, got %v", exists)
	}
}

// TestExistsMorphTo tests checking if MorphTo relationship exists
func TestExistsMorphTo(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	existence := NewRelationshipExistence(field)

	ctx := context.Background()
	exists, err := existence.Exists(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if exists != false {
		t.Errorf("Expected exists false, got %v", exists)
	}
}

// TestDoesntExist tests checking if no related resources exist
func TestDoesntExist(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	existence := NewRelationshipExistence(field)

	ctx := context.Background()
	doesntExist, err := existence.DoesntExist(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if doesntExist != true {
		t.Errorf("Expected doesntExist true, got %v", doesntExist)
	}
}
