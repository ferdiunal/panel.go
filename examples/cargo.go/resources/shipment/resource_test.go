package shipment

import (
	"testing"

	"cargo.go/entity"
)

// TestShipmentResourceRecordTitle - RecordTitle fonksiyonu testi
func TestShipmentResourceRecordTitle(t *testing.T) {
	resource := NewShipmentResource()

	// Test case 1: Valid shipment
	shipment := &entity.Shipment{
		ID:   1,
		Name: "Test Shipment",
	}

	title := resource.RecordTitle(shipment)
	if title != "Test Shipment" {
		t.Errorf("Expected title 'Test Shipment', got '%s'", title)
	}

	// Test case 2: Invalid type
	invalidRecord := "not a shipment"
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

// TestShipmentResourceWith - With fonksiyonu testi
func TestShipmentResourceWith(t *testing.T) {
	resource := NewShipmentResource()

	relationships := resource.With()

	// Organization, SenderAddress, ReceiverAddress, ShipmentRows ilişkileri yüklenmeli
	if len(relationships) != 4 {
		t.Errorf("Expected 4 relationships, got %d", len(relationships))
	}

	expectedRelationships := map[string]bool{
		"Organization":    false,
		"SenderAddress":   false,
		"ReceiverAddress": false,
		"ShipmentRows":    false,
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

