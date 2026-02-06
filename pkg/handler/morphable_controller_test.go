package handler

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/fields"
)

// TestGetDisplayValue tests the getDisplayValue function
func TestGetDisplayValue(t *testing.T) {
	tests := []struct {
		name     string
		record   map[string]interface{}
		expected string
	}{
		{
			name: "name field present",
			record: map[string]interface{}{
				"id":   1,
				"name": "Test User",
			},
			expected: "Test User",
		},
		{
			name: "title field present",
			record: map[string]interface{}{
				"id":    1,
				"title": "Test Post",
			},
			expected: "Test Post",
		},
		{
			name: "email field present",
			record: map[string]interface{}{
				"id":    1,
				"email": "test@example.com",
			},
			expected: "test@example.com",
		},
		{
			name: "username field present",
			record: map[string]interface{}{
				"id":       1,
				"username": "testuser",
			},
			expected: "testuser",
		},
		{
			name: "fallback to id",
			record: map[string]interface{}{
				"id": 42,
			},
			expected: "#42",
		},
		{
			name: "name takes priority over title",
			record: map[string]interface{}{
				"id":    1,
				"name":  "User Name",
				"title": "User Title",
			},
			expected: "User Name",
		},
		{
			name:     "empty record returns Unknown",
			record:   map[string]interface{}{},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDisplayValue(tt.record)
			if result != tt.expected {
				t.Errorf("getDisplayValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMorphableOption tests the MorphableOption struct
func TestMorphableOption(t *testing.T) {
	opt := MorphableOption{
		Value:    1,
		Display:  "Test User",
		Avatar:   "https://example.com/avatar.png",
		Subtitle: "Admin",
	}

	if opt.Value != 1 {
		t.Errorf("Expected Value 1, got %v", opt.Value)
	}
	if opt.Display != "Test User" {
		t.Errorf("Expected Display 'Test User', got %v", opt.Display)
	}
	if opt.Avatar != "https://example.com/avatar.png" {
		t.Errorf("Expected Avatar URL, got %v", opt.Avatar)
	}
	if opt.Subtitle != "Admin" {
		t.Errorf("Expected Subtitle 'Admin', got %v", opt.Subtitle)
	}
}

// TestMorphableControllerCreation tests NewMorphableController
func TestMorphableControllerCreation(t *testing.T) {
	resources := make(map[string]interface{})
	resources["users"] = "users resource"

	controller := NewMorphableController(nil, resources)

	if controller == nil {
		t.Error("Expected controller to be created")
	}
	if controller.Resources == nil {
		t.Error("Expected Resources to be initialized")
	}
	if len(controller.Resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(controller.Resources))
	}
}

// TestQueryMorphableResourcesNilDB tests queryMorphableResources with nil DB
func TestQueryMorphableResourcesNilDB(t *testing.T) {
	results, err := queryMorphableResources(nil, "users", "", 10, "")

	if err != nil {
		t.Errorf("Expected no error for nil DB, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected empty results for nil DB, got %d results", len(results))
	}
}

// TestMorphToFieldIntegration tests MorphTo field with HandleMorphable
func TestMorphToFieldIntegration(t *testing.T) {
	// Create a MorphTo field
	field := fields.NewMorphTo("Commentable", "commentable").
		Types(map[string]string{
			"posts":  "posts",
			"videos": "videos",
		})

	// Verify types are properly set
	types := field.GetTypes()
	if len(types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(types))
	}

	if types["posts"] != "posts" {
		t.Errorf("Expected posts type mapping, got %v", types["posts"])
	}

	if types["videos"] != "videos" {
		t.Errorf("Expected videos type mapping, got %v", types["videos"])
	}

	// Verify GetResourceForType
	slug, err := field.GetResourceForType("posts")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if slug != "posts" {
		t.Errorf("Expected 'posts' slug, got %v", slug)
	}

	// Verify unknown type returns error
	_, err = field.GetResourceForType("unknown")
	if err == nil {
		t.Error("Expected error for unknown type")
	}
}
