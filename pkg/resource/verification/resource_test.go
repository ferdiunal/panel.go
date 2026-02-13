package verification

import (
	"testing"

	domainVerification "github.com/ferdiunal/panel.go/pkg/domain/verification"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// TestNewVerificationResource, yeni Verification resource'u oluşturur
func TestNewVerificationResource(t *testing.T) {
	r := NewVerificationResource()

	if r == nil {
		t.Error("Expected non-nil resource")
	}

	if r.Slug() != "verifications" {
		t.Errorf("Expected slug 'verifications', got '%s'", r.Slug())
	}

	if r.Title() != "Verifications" {
		t.Errorf("Expected title 'Verifications', got '%s'", r.Title())
	}
}

// TestVerificationResourceModel, model'i test eder
func TestVerificationResourceModel(t *testing.T) {
	r := NewVerificationResource()

	model := r.Model()
	if model == nil {
		t.Error("Expected non-nil model")
	}

	_, ok := model.(*domainVerification.Verification)
	if !ok {
		t.Error("Expected Verification model")
	}
}

// TestVerificationResourceFields, alanları test eder
func TestVerificationResourceFields(t *testing.T) {
	r := NewVerificationResource()

	fields := r.Fields()
	if len(fields) == 0 {
		t.Error("Expected at least one field")
	}
}

// TestVerificationResourceImplementsResource, Resource interface'ini implement ettiğini test eder
func TestVerificationResourceImplementsResource(t *testing.T) {
	var _ resource.Resource = (*VerificationResource)(nil)
}

// TestVerificationResourcePolicy, policy'yi test eder
func TestVerificationResourcePolicy(t *testing.T) {
	r := NewVerificationResource()

	policy := r.Policy()
	if policy == nil {
		t.Error("Expected non-nil policy")
	}
}

// TestVerificationResourceSortable, sıralama ayarlarını test eder
func TestVerificationResourceSortable(t *testing.T) {
	r := NewVerificationResource()

	sortable := r.GetSortable()
	if len(sortable) == 0 {
		t.Error("Expected at least one sortable field")
	}

	if sortable[0].Column != "created_at" {
		t.Errorf("Expected 'created_at', got '%s'", sortable[0].Column)
	}

	if sortable[0].Direction != "desc" {
		t.Errorf("Expected 'desc', got '%s'", sortable[0].Direction)
	}
}
