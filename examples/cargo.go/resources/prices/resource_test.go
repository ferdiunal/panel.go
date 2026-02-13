package prices

import (
	"testing"

	"cargo.go/entity"
)

// TestPriceListResourceRecordTitle - RecordTitle fonksiyonu testi
func TestPriceListResourceRecordTitle(t *testing.T) {
	resource := NewPriceListResource()

	// Test case 1: Valid price list
	priceList := &entity.PriceList{
		ID:   1,
		Name: "Test Price List",
	}

	title := resource.RecordTitle(priceList)
	if title != "Test Price List" {
		t.Errorf("Expected title 'Test Price List', got '%s'", title)
	}

	// Test case 2: Invalid type
	invalidRecord := "not a price list"
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

