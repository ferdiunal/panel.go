package account

import (
	"testing"

	domainAccount "github.com/ferdiunal/panel.go/pkg/domain/account"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// TestNewAccountResource, yeni Account resource'u oluşturur
func TestNewAccountResource(t *testing.T) {
	r := NewAccountResource()

	if r == nil {
		t.Error("Expected non-nil resource")
	}

	if r.Slug() != "accounts" {
		t.Errorf("Expected slug 'accounts', got '%s'", r.Slug())
	}

	if r.Title() != "Accounts" {
		t.Errorf("Expected title 'Accounts', got '%s'", r.Title())
	}
}

// TestAccountResourceModel, model'i test eder
func TestAccountResourceModel(t *testing.T) {
	r := NewAccountResource()

	model := r.Model()
	if model == nil {
		t.Error("Expected non-nil model")
	}

	_, ok := model.(*domainAccount.Account)
	if !ok {
		t.Error("Expected Account model")
	}
}

// TestAccountResourceFields, alanları test eder
func TestAccountResourceFields(t *testing.T) {
	r := NewAccountResource()

	fields := r.Fields()
	if len(fields) == 0 {
		t.Error("Expected at least one field")
	}
}

// TestAccountResourceImplementsResource, Resource interface'ini implement ettiğini test eder
func TestAccountResourceImplementsResource(t *testing.T) {
	var _ resource.Resource = (*AccountResource)(nil)
}

// TestAccountResourcePolicy, policy'yi test eder
func TestAccountResourcePolicy(t *testing.T) {
	r := NewAccountResource()

	policy := r.Policy()
	if policy == nil {
		t.Error("Expected non-nil policy")
	}
}

// TestAccountResourceSortable, sıralama ayarlarını test eder
func TestAccountResourceSortable(t *testing.T) {
	r := NewAccountResource()

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
