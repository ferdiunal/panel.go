package setting

import (
	"testing"

	domainSetting "github.com/ferdiunal/panel.go/pkg/domain/setting"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// TestNewSettingResource, yeni Setting resource'u oluşturur
func TestNewSettingResource(t *testing.T) {
	r := NewSettingResource()

	if r == nil {
		t.Error("Expected non-nil resource")
	}

	if r.Slug() != "settings" {
		t.Errorf("Expected slug 'settings', got '%s'", r.Slug())
	}

	if r.Title() != "Settings" {
		t.Errorf("Expected title 'Settings', got '%s'", r.Title())
	}
}

// TestSettingResourceModel, model'i test eder
func TestSettingResourceModel(t *testing.T) {
	r := NewSettingResource()

	model := r.Model()
	if model == nil {
		t.Error("Expected non-nil model")
	}

	_, ok := model.(*domainSetting.Setting)
	if !ok {
		t.Error("Expected Setting model")
	}
}

// TestSettingResourceFields, alanları test eder
func TestSettingResourceFields(t *testing.T) {
	r := NewSettingResource()

	fields := r.Fields()
	if len(fields) == 0 {
		t.Error("Expected at least one field")
	}
}

// TestSettingResourceImplementsResource, Resource interface'ini implement ettiğini test eder
func TestSettingResourceImplementsResource(t *testing.T) {
	var _ resource.Resource = (*SettingResource)(nil)
}

// TestSettingResourcePolicy, policy'yi test eder
func TestSettingResourcePolicy(t *testing.T) {
	r := NewSettingResource()

	policy := r.Policy()
	if policy == nil {
		t.Error("Expected non-nil policy")
	}
}

// TestSettingResourceSortable, sıralama ayarlarını test eder
func TestSettingResourceSortable(t *testing.T) {
	r := NewSettingResource()

	sortable := r.GetSortable()
	if len(sortable) == 0 {
		t.Error("Expected at least one sortable field")
	}

	if sortable[0].Column != "group" {
		t.Errorf("Expected 'group', got '%s'", sortable[0].Column)
	}

	if sortable[0].Direction != "asc" {
		t.Errorf("Expected 'asc', got '%s'", sortable[0].Direction)
	}
}
