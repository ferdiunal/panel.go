package cargo_company

import (
	"testing"

	"cargo.go/entity"
)

// TestCargoCompanyResourceRecordTitle - RecordTitle fonksiyonu testi
func TestCargoCompanyResourceRecordTitle(t *testing.T) {
	resource := NewCargoCompanyResource()

	// Test case 1: Valid cargo company
	company := &entity.CargoCompany{
		ID:   1,
		Name: "Test Cargo Company",
	}

	title := resource.RecordTitle(company)
	if title != "Test Cargo Company" {
		t.Errorf("Expected title 'Test Cargo Company', got '%s'", title)
	}

	// Test case 2: Invalid type
	invalidRecord := "not a cargo company"
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

