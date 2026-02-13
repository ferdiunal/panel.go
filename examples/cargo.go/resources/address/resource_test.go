package address

import (
	"testing"

	"cargo.go/entity"
)

// TestAddressResourceRecordTitle - RecordTitle fonksiyonu testi
func TestAddressResourceRecordTitle(t *testing.T) {
	resource := NewAddressResource()

	// Test case 1: Valid address
	addr := &entity.Address{
		ID:   1,
		Name: "Test Address",
	}

	title := resource.RecordTitle(addr)
	if title != "Test Address" {
		t.Errorf("Expected title 'Test Address', got '%s'", title)
	}

	// Test case 2: Invalid type
	invalidRecord := "not an address"
	title = resource.RecordTitle(invalidRecord)
	if title != "" {
		t.Errorf("Expected empty title for invalid type, got '%s'", title)
	}

	// Test case 3: Nil record
	title = resource.RecordTitle(nil)
	if title != "" {
		t.Errorf("Expected empty title for nil record, got '%s'", title)
	}
}

// TestAddressResourceWith - With fonksiyonu testi
func TestAddressResourceWith(t *testing.T) {
	resource := NewAddressResource()

	relationships := resource.With()

	// Organization ilişkisi yüklenmeli
	if len(relationships) != 1 {
		t.Errorf("Expected 1 relationship, got %d", len(relationships))
	}

	if relationships[0] != "Organization" {
		t.Errorf("Expected 'Organization' relationship, got '%s'", relationships[0])
	}
}

