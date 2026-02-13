package shipment_row

import (
	"testing"

	"cargo.go/entity"
)

// TestShipmentRowResourceRecordTitle - RecordTitle fonksiyonu testi
func TestShipmentRowResourceRecordTitle(t *testing.T) {
	resource := NewShipmentRowResource()

	// Test case 1: ShipmentRow with Product
	product := &entity.Product{
		ID:   1,
		Name: "Test Product",
	}

	row := &entity.ShipmentRow{
		ID:        1,
		ProductID: 1,
		Product:   product,
		Quantity:  5,
	}

	title := resource.RecordTitle(row)
	if title != "Test Product" {
		t.Errorf("Expected title 'Test Product', got '%s'", title)
	}

	// Test case 2: ShipmentRow with Shipment (no Product)
	shipment := &entity.Shipment{
		ID:   1,
		Name: "Test Shipment",
	}

	row2 := &entity.ShipmentRow{
		ID:         2,
		ShipmentID: 1,
		Shipment:   shipment,
		Quantity:   10,
	}

	title = resource.RecordTitle(row2)
	if title != "Test Shipment" {
		t.Errorf("Expected title 'Test Shipment', got '%s'", title)
	}

	// Test case 3: ShipmentRow without Product and Shipment
	row3 := &entity.ShipmentRow{
		ID:       3,
		Quantity: 15,
	}

	title = resource.RecordTitle(row3)
	if title != "" {
		t.Errorf("Expected empty title, got '%s'", title)
	}

	// Test case 4: Invalid type
	invalidRecord := "not a shipment row"
	title = resource.RecordTitle(invalidRecord)
	if title != "" {
		t.Errorf("Expected empty title for invalid type, got '%s'", title)
	}

	// Test case 5: Nil record
	title = resource.RecordTitle(nil)
	if title != "" {
		t.Errorf("Expected empty title for nil record, got '%s'", title)
	}
}

// TestShipmentRowResourceWith - With fonksiyonu testi
func TestShipmentRowResourceWith(t *testing.T) {
	resource := NewShipmentRowResource()

	relationships := resource.With()

	// Shipment ve Product ilişkileri yüklenmeli
	if len(relationships) != 2 {
		t.Errorf("Expected 2 relationships, got %d", len(relationships))
	}

	expectedRelationships := map[string]bool{
		"Shipment": false,
		"Product":  false,
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

