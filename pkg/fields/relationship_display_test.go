package fields

import (
	"testing"
)

// TestNewRelationshipDisplay tests creating a new relationship display handler
func TestNewRelationshipDisplay(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	display := NewRelationshipDisplay(field)

	if display == nil {
		t.Error("Expected non-nil display handler")
	}

	if display.field != field {
		t.Error("Expected field to be set")
	}
}

// TestDisplayBelongsTo tests displaying a BelongsTo relationship
func TestDisplayBelongsTo(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	display := NewRelationshipDisplay(field)

	displayValue, err := display.GetDisplayValue(nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if displayValue != "" {
		t.Errorf("Expected empty display value for nil, got '%s'", displayValue)
	}
}

// TestDisplayHasMany tests displaying a HasMany relationship
func TestDisplayHasMany(t *testing.T) {
	field := NewHasMany("Posts", "posts", "posts")
	display := NewRelationshipDisplay(field)

	displayValue, err := display.GetDisplayValue(map[string]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if displayValue == "" {
		t.Error("Expected non-empty display value")
	}
}

// TestDisplayHasOne tests displaying a HasOne relationship
func TestDisplayHasOne(t *testing.T) {
	field := NewHasOne("Profile", "profile", "profiles")
	display := NewRelationshipDisplay(field)

	displayValue, err := display.GetDisplayValue(map[string]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if displayValue != "Related resource" {
		t.Errorf("Expected 'Related resource', got '%s'", displayValue)
	}
}

// TestDisplayBelongsToMany tests displaying a BelongsToMany relationship
func TestDisplayBelongsToMany(t *testing.T) {
	field := NewBelongsToMany("Tags", "tags", "tags")
	display := NewRelationshipDisplay(field)

	displayValue, err := display.GetDisplayValue(map[string]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if displayValue == "" {
		t.Error("Expected non-empty display value")
	}
}

// TestDisplayMorphTo tests displaying a MorphTo relationship
func TestDisplayMorphTo(t *testing.T) {
	field := NewMorphTo("Commentable", "commentable")
	display := NewRelationshipDisplay(field)

	displayValue, err := display.GetDisplayValue(map[string]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if displayValue == "" {
		t.Error("Expected non-empty display value")
	}
}

// TestGetDisplayValues tests getting display values for multiple items
func TestGetDisplayValues(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	display := NewRelationshipDisplay(field)

	items := []interface{}{nil, nil, nil}
	displayValues, err := display.GetDisplayValues(items)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(displayValues) != 3 {
		t.Errorf("Expected 3 display values, got %d", len(displayValues))
	}
}

// TestGetDisplayValuesEmpty tests getting display values for empty items
func TestGetDisplayValuesEmpty(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	display := NewRelationshipDisplay(field)

	displayValues, err := display.GetDisplayValues([]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(displayValues) != 0 {
		t.Errorf("Expected 0 display values, got %d", len(displayValues))
	}
}

// TestGetDisplayValuesNil tests getting display values for nil items
func TestGetDisplayValuesNil(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	display := NewRelationshipDisplay(field)

	displayValues, err := display.GetDisplayValues(nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(displayValues) != 0 {
		t.Errorf("Expected 0 display values, got %d", len(displayValues))
	}
}

// TestFormatDisplayValue tests formatting a display value
func TestFormatDisplayValue(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	display := NewRelationshipDisplay(field)

	formattedValue := display.FormatDisplayValue("John Doe")

	if formattedValue != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", formattedValue)
	}
}

// TestFormatDisplayValueNil tests formatting a nil display value
func TestFormatDisplayValueNil(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	display := NewRelationshipDisplay(field)

	formattedValue := display.FormatDisplayValue(nil)

	if formattedValue != "" {
		t.Errorf("Expected empty string, got '%s'", formattedValue)
	}
}

// TestDisplayWithCustomDisplayKey tests displaying with custom display key
func TestDisplayWithCustomDisplayKey(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	field.DisplayUsing("email")
	display := NewRelationshipDisplay(field)

	displayValue, err := display.GetDisplayValue(map[string]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if displayValue == "" {
		t.Error("Expected non-empty display value")
	}
}
