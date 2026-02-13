package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipConstraints tests creating a new relationship constraints handler
func TestNewRelationshipConstraints(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	if constraints == nil {
		t.Error("Expected non-nil constraints handler")
	}

	if constraints.field != field {
		t.Error("Expected field to be set")
	}

	if constraints.limit != 0 {
		t.Errorf("Expected limit 0, got %d", constraints.limit)
	}

	if constraints.offset != 0 {
		t.Errorf("Expected offset 0, got %d", constraints.offset)
	}
}

// TestApplyLimit tests applying a limit constraint
func TestApplyLimit(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	results, err := constraints.ApplyLimit(ctx, 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if constraints.limit != 10 {
		t.Errorf("Expected limit 10, got %d", constraints.limit)
	}
}

// TestApplyLimitNegative tests applying a negative limit
func TestApplyLimitNegative(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	_, err := constraints.ApplyLimit(ctx, -5)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if constraints.limit != 0 {
		t.Errorf("Expected limit 0, got %d", constraints.limit)
	}
}

// TestApplyOffset tests applying an offset constraint
func TestApplyOffset(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	results, err := constraints.ApplyOffset(ctx, 20)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if constraints.offset != 20 {
		t.Errorf("Expected offset 20, got %d", constraints.offset)
	}
}

// TestApplyOffsetNegative tests applying a negative offset
func TestApplyOffsetNegative(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	_, err := constraints.ApplyOffset(ctx, -10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if constraints.offset != 0 {
		t.Errorf("Expected offset 0, got %d", constraints.offset)
	}
}

// TestApplyWhere tests applying a WHERE constraint
func TestApplyWhere(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	results, err := constraints.ApplyWhere(ctx, "status", "=", "active")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(constraints.constraints) != 1 {
		t.Errorf("Expected 1 constraint, got %d", len(constraints.constraints))
	}
}

// TestApplyWhereIn tests applying a WHERE IN constraint
func TestApplyWhereIn(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	values := []interface{}{1, 2, 3}
	results, err := constraints.ApplyWhereIn(ctx, "id", values)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if len(constraints.constraints) != 1 {
		t.Errorf("Expected 1 constraint, got %d", len(constraints.constraints))
	}
}

// TestGetLimit tests getting limit
func TestGetLimit(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	constraints.ApplyLimit(ctx, 15)

	if constraints.GetLimit() != 15 {
		t.Errorf("Expected limit 15, got %d", constraints.GetLimit())
	}
}

// TestGetOffset tests getting offset
func TestGetOffset(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	constraints.ApplyOffset(ctx, 25)

	if constraints.GetOffset() != 25 {
		t.Errorf("Expected offset 25, got %d", constraints.GetOffset())
	}
}

// TestGetConstraints tests getting constraints
func TestGetConstraints(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	constraints := NewRelationshipConstraints(field)

	ctx := context.Background()
	constraints.ApplyWhere(ctx, "status", "=", "active")

	constraintMap := constraints.GetConstraints()

	if len(constraintMap) != 1 {
		t.Errorf("Expected 1 constraint, got %d", len(constraintMap))
	}

	if _, ok := constraintMap["status"]; !ok {
		t.Error("Expected 'status' constraint to exist")
	}
}
