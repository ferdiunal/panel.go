package products

import (
	"testing"

	"cargo.go/entity"
)

// TestProductResourceRecordTitle - RecordTitle fonksiyonu testi
func TestProductResourceRecordTitle(t *testing.T) {
	resource := NewProductResource()

	// Test case 1: Valid product
	product := &entity.Product{
		ID:   1,
		Name: "Test Product",
	}

	title := resource.RecordTitle(product)
	if title != "Test Product" {
		t.Errorf("Expected title 'Test Product', got '%s'", title)
	}

	// Test case 2: Invalid type
	invalidRecord := "not a product"
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

// TestProductResourceWith - With fonksiyonu testi
func TestProductResourceWith(t *testing.T) {
	resource := NewProductResource()

	relationships := resource.With()

	// Organization ve ShipmentRows ilişkileri yüklenmeli
	if len(relationships) != 2 {
		t.Errorf("Expected 2 relationships, got %d", len(relationships))
	}

	expectedRelationships := map[string]bool{
		"Organization": false,
		"ShipmentRows": false,
	}

	for _, rel := range relationships {
		if _, ok := expectedRelationships[rel]; ok {
			expectedRelationships[rel] = true
		}
	}

	for rel, found := range expectedRelationships {
		if !found {
			t.Errorf("Expected relationship '%s' not found", rel)
		}
	}
}

