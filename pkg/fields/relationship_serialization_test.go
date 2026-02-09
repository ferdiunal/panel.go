package fields

import (
	"encoding/json"
	"testing"
)

// TestNewRelationshipSerialization tests creating a new relationship serialization handler
func TestNewRelationshipSerialization(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	if serialization == nil {
		t.Error("Expected non-nil serialization handler")
	}

	if serialization.field != field {
		t.Error("Expected field to be set")
	}
}

// TestSerializeRelationshipNil tests serializing a nil relationship
func TestSerializeRelationshipNil(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	jsonData, err := serialization.SerializeRelationship(nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if jsonData["value"] != nil {
		t.Errorf("Expected nil value, got %v", jsonData["value"])
	}
}

// TestSerializeRelationship tests serializing a relationship
func TestSerializeRelationship(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	jsonData, err := serialization.SerializeRelationship(map[string]interface{}{"id": 1, "name": "John"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if jsonData["type"] != "belongsTo" {
		t.Errorf("Expected type 'belongsTo', got %v", jsonData["type"])
	}

	if jsonData["name"] != "Author" {
		t.Errorf("Expected name 'Author', got %v", jsonData["name"])
	}

	if jsonData["resource"] != "authors" {
		t.Errorf("Expected resource 'authors', got %v", jsonData["resource"])
	}
}

// TestSerializeRelationships tests serializing multiple relationships
func TestSerializeRelationships(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	items := []interface{}{
		map[string]interface{}{"id": 1, "name": "John"},
		map[string]interface{}{"id": 2, "name": "Jane"},
	}

	jsonData, err := serialization.SerializeRelationships(items)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(jsonData) != 2 {
		t.Errorf("Expected 2 items, got %d", len(jsonData))
	}
}

// TestSerializeRelationshipsEmpty tests serializing empty relationships
func TestSerializeRelationshipsEmpty(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	jsonData, err := serialization.SerializeRelationships([]interface{}{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(jsonData) != 0 {
		t.Errorf("Expected 0 items, got %d", len(jsonData))
	}
}

// TestSerializeRelationshipsNil tests serializing nil relationships
func TestSerializeRelationshipsNil(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	jsonData, err := serialization.SerializeRelationships(nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(jsonData) != 0 {
		t.Errorf("Expected 0 items, got %d", len(jsonData))
	}
}

// TestToJSON tests converting relationship to JSON string
func TestToJSON(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	jsonStr, err := serialization.ToJSON(map[string]interface{}{"id": 1, "name": "John"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	// Verify it's valid JSON
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}
}

// TestToJSONArray tests converting relationships to JSON array string
func TestToJSONArray(t *testing.T) {
	field := BelongsTo("Author", "author_id", "authors")
	serialization := NewRelationshipSerialization(field)

	items := []interface{}{
		map[string]interface{}{"id": 1, "name": "John"},
		map[string]interface{}{"id": 2, "name": "Jane"},
	}

	jsonStr, err := serialization.ToJSONArray(items)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	// Verify it's valid JSON
	var data []map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 items in JSON array, got %d", len(data))
	}
}
