package commissions

import (
	"testing"

	"cargo.go/entity"
)

// TestCommissionResourceRecordTitle - RecordTitle fonksiyonu testi
func TestCommissionResourceRecordTitle(t *testing.T) {
	resource := NewCommissionResource()

	// Test case 1: Valid commission
	commission := &entity.Commission{
		ID:   1,
		Name: "Test Commission",
	}

	title := resource.RecordTitle(commission)
	if title != "Test Commission" {
		t.Errorf("Expected title 'Test Commission', got '%s'", title)
	}

	// Test case 2: Invalid type
	invalidRecord := "not a commission"
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

