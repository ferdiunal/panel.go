package session

import (
	"testing"

	domainSession "github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// TestNewSessionResource, yeni Session resource'u oluşturur
func TestNewSessionResource(t *testing.T) {
	r := NewSessionResource()

	if r == nil {
		t.Error("Expected non-nil resource")
	}

	if r.Slug() != "sessions" {
		t.Errorf("Expected slug 'sessions', got '%s'", r.Slug())
	}

	if r.Title() != "Sessions" {
		t.Errorf("Expected title 'Sessions', got '%s'", r.Title())
	}
}

// TestSessionResourceModel, model'i test eder
func TestSessionResourceModel(t *testing.T) {
	r := NewSessionResource()

	model := r.Model()
	if model == nil {
		t.Error("Expected non-nil model")
	}

	_, ok := model.(*domainSession.Session)
	if !ok {
		t.Error("Expected Session model")
	}
}

// TestSessionResourceFields, alanları test eder
func TestSessionResourceFields(t *testing.T) {
	r := NewSessionResource()

	fields := r.Fields()
	if len(fields) == 0 {
		t.Error("Expected at least one field")
	}
}

// TestSessionResourceImplementsResource, Resource interface'ini implement ettiğini test eder
func TestSessionResourceImplementsResource(t *testing.T) {
	var _ resource.Resource = (*SessionResource)(nil)
}

// TestSessionResourcePolicy, policy'yi test eder
func TestSessionResourcePolicy(t *testing.T) {
	r := NewSessionResource()

	policy := r.Policy()
	if policy == nil {
		t.Error("Expected non-nil policy")
	}
}

// TestSessionResourceSortable, sıralama ayarlarını test eder
func TestSessionResourceSortable(t *testing.T) {
	r := NewSessionResource()

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
