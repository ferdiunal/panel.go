package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipValidator tests creating a new relationship validator
func TestNewRelationshipValidator(t *testing.T) {
	validator := NewRelationshipValidator()

	if validator == nil {
		t.Error("Expected non-nil validator")
	}
}

// TestValidateExistsNilValue tests ValidateExists with nil value
func TestValidateExistsNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewBelongsTo("Author", "author_id", "authors")

	ctx := context.Background()
	err := validator.ValidateExists(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value on nullable field, got %v", err)
	}
}

// TestValidateExistsRequiredField tests ValidateExists with required field
func TestValidateExistsRequiredField(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewBelongsTo("Author", "author_id", "authors")
	field.Required()

	ctx := context.Background()
	err := validator.ValidateExists(ctx, nil, field)

	if err == nil {
		t.Error("Expected error for nil value on required field")
	}
}

// TestValidateForeignKeyNilValue tests ValidateForeignKey with nil value
func TestValidateForeignKeyNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewBelongsTo("Author", "author_id", "authors")

	ctx := context.Background()
	err := validator.ValidateForeignKey(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}
}

// TestValidatePivotNilValue tests ValidatePivot with nil value
func TestValidatePivotNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewBelongsToMany("Tags", "tags", "tags")

	ctx := context.Background()
	err := validator.ValidatePivot(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}
}

// TestValidateMorphTypeNilValue tests ValidateMorphType with nil value
func TestValidateMorphTypeNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewMorphTo("Commentable", "commentable")

	ctx := context.Background()
	err := validator.ValidateMorphType(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}
}

// TestValidateMorphTypeNoTypesRegistered tests ValidateMorphType with no types registered
func TestValidateMorphTypeNoTypesRegistered(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewMorphTo("Commentable", "commentable")

	ctx := context.Background()
	err := validator.ValidateMorphType(ctx, &Taggable{}, field)

	if err == nil {
		t.Error("Expected error for no types registered")
	}
}

// TestValidateBelongsToNilValue tests ValidateBelongsTo with nil value
func TestValidateBelongsToNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewBelongsTo("Author", "author_id", "authors")

	ctx := context.Background()
	err := validator.ValidateBelongsTo(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value on nullable field, got %v", err)
	}
}

// TestValidateBelongsToRequiredField tests ValidateBelongsTo with required field
func TestValidateBelongsToRequiredField(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewBelongsTo("Author", "author_id", "authors")
	field.Required()

	ctx := context.Background()
	err := validator.ValidateBelongsTo(ctx, nil, field)

	if err == nil {
		t.Error("Expected error for nil value on required field")
	}
}

// TestValidateHasManyNilValue tests ValidateHasMany with nil value
func TestValidateHasManyNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewHasMany("Posts", "posts", "posts")

	ctx := context.Background()
	err := validator.ValidateHasMany(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}
}

// TestValidateHasOneNilValue tests ValidateHasOne with nil value
func TestValidateHasOneNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewHasOne("Profile", "profile", "profiles")

	ctx := context.Background()
	err := validator.ValidateHasOne(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}
}

// TestValidateBelongsToManyNilValue tests ValidateBelongsToMany with nil value
func TestValidateBelongsToManyNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewBelongsToMany("Tags", "tags", "tags")

	ctx := context.Background()
	err := validator.ValidateBelongsToMany(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}
}

// TestValidateMorphToNilValue tests ValidateMorphTo with nil value
func TestValidateMorphToNilValue(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewMorphTo("Commentable", "commentable")

	ctx := context.Background()
	err := validator.ValidateMorphTo(ctx, nil, field)

	if err != nil {
		t.Errorf("Expected no error for nil value, got %v", err)
	}
}

// TestValidateMorphToNoTypesRegistered tests ValidateMorphTo with no types registered
func TestValidateMorphToNoTypesRegistered(t *testing.T) {
	validator := NewRelationshipValidator()
	field := NewMorphTo("Commentable", "commentable")

	ctx := context.Background()
	err := validator.ValidateMorphTo(ctx, &Taggable{}, field)

	if err == nil {
		t.Error("Expected error for no types registered")
	}
}
