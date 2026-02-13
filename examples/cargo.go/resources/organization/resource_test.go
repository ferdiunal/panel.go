package organization

import (
	"testing"

	"cargo.go/entity"
)

// TestOrganizationResourceRecordTitle - RecordTitle fonksiyonu testi
func TestOrganizationResourceRecordTitle(t *testing.T) {
	resource := NewOrganizationResource()

	// Test case 1: Valid organization
	org := &entity.Organization{
		ID:   1,
		Name: "Test Organization",
	}

	title := resource.RecordTitle(org)
	if title != "Test Organization" {
		t.Errorf("Expected title 'Test Organization', got '%s'", title)
	}

	// Test case 2: Invalid type
	invalidRecord := "not an organization"
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

// TestOrganizationResourceWith - With fonksiyonu testi
func TestOrganizationResourceWith(t *testing.T) {
	resource := NewOrganizationResource()

	relationships := resource.With()

	// Addresses ilişkisi yüklenmeli
	if len(relationships) != 1 {
		t.Errorf("Expected 1 relationship, got %d", len(relationships))
	}

	if relationships[0] != "Addresses" {
		t.Errorf("Expected 'Addresses' relationship, got '%s'", relationships[0])
	}
}


// TestOrganizationPolicy - Policy test'leri
func TestOrganizationPolicy(t *testing.T) {
	policy := &OrganizationPolicy{}

	// ViewAny her zaman true dönmeli
	if !policy.ViewAny(nil) {
		t.Error("Expected ViewAny to return true")
	}

	// View her zaman true dönmeli
	if !policy.View(nil, nil) {
		t.Error("Expected View to return true")
	}

	// Create her zaman true dönmeli
	if !policy.Create(nil) {
		t.Error("Expected Create to return true")
	}

	// Update her zaman true dönmeli
	if !policy.Update(nil, nil) {
		t.Error("Expected Update to return true")
	}

	// Delete her zaman true dönmeli
	if !policy.Delete(nil, nil) {
		t.Error("Expected Delete to return true")
	}
}
