package fields

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/context"
	"gorm.io/gorm"
)

// TestSearchableMixin, Searchable mixin'ini test eder
func TestSearchableMixin(t *testing.T) {
	s := &Searchable{}

	// Varsayılan değer
	cols := s.GetSearchableColumns()
	if len(cols) != 0 {
		t.Error("Expected empty searchable columns by default")
	}

	// Sütunları ayarla
	s.SetSearchableColumns([]string{"name", "email"})
	cols = s.GetSearchableColumns()

	if len(cols) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(cols))
	}

	if cols[0] != "name" || cols[1] != "email" {
		t.Error("Expected columns to be set correctly")
	}
}

// TestSearchableCallback, Searchable callback'ini test eder
func TestSearchableCallback(t *testing.T) {
	s := &Searchable{}

	// Varsayılan callback nil
	if s.GetSearchCallback() != nil {
		t.Error("Expected nil callback by default")
	}

	// Callback ayarla
	callback := func(db *gorm.DB, term string) *gorm.DB {
		return db
	}
	s.SetSearchCallback(callback)

	if s.GetSearchCallback() == nil {
		t.Error("Expected callback to be set")
	}
}

// TestSortableMixin, Sortable mixin'ini test eder
func TestSortableMixin(t *testing.T) {
	s := &Sortable{}

	// Varsayılan değer
	if s.IsSortable() {
		t.Error("Expected not sortable by default")
	}

	// Sıralanabilir yap
	s.SetSortable(true)
	if !s.IsSortable() {
		t.Error("Expected sortable after setting")
	}

	// Sıralama yönü
	if s.GetSortDirection() != "asc" {
		t.Errorf("Expected 'asc', got '%s'", s.GetSortDirection())
	}

	s.SetSortDirection("desc")
	if s.GetSortDirection() != "desc" {
		t.Errorf("Expected 'desc', got '%s'", s.GetSortDirection())
	}
}

// TestSortableCallback, Sortable callback'ini test eder
func TestSortableCallback(t *testing.T) {
	s := &Sortable{}

	// Varsayılan callback nil
	if s.GetSortCallback() != nil {
		t.Error("Expected nil callback by default")
	}

	// Callback ayarla
	callback := func(db *gorm.DB, direction string) *gorm.DB {
		return db
	}
	s.SetSortCallback(callback)

	if s.GetSortCallback() == nil {
		t.Error("Expected callback to be set")
	}
}

// TestFilterableMixin, Filterable mixin'ini test eder
func TestFilterableMixin(t *testing.T) {
	f := &Filterable{}

	// Varsayılan değer
	if f.IsFilterable() {
		t.Error("Expected not filterable by default")
	}

	// Filtrelenebilir yap
	f.SetFilterable(true)
	if !f.IsFilterable() {
		t.Error("Expected filterable after setting")
	}

	// Seçenekler
	options := f.GetFilterOptions()
	if len(options) != 0 {
		t.Error("Expected empty options by default")
	}

	f.SetFilterOptions(map[string]string{
		"active":   "Aktif",
		"inactive": "İnaktif",
	})

	options = f.GetFilterOptions()
	if len(options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(options))
	}
}

// TestFilterableCallback, Filterable callback'ini test eder
func TestFilterableCallback(t *testing.T) {
	f := &Filterable{}

	// Varsayılan callback nil
	if f.GetFilterCallback() != nil {
		t.Error("Expected nil callback by default")
	}

	// Callback ayarla
	callback := func(db *gorm.DB, value any) *gorm.DB {
		return db
	}
	f.SetFilterCallback(callback)

	if f.GetFilterCallback() == nil {
		t.Error("Expected callback to be set")
	}
}

// TestValidatableMixin, Validatable mixin'ini test eder
func TestValidatableMixin(t *testing.T) {
	v := &Validatable{}

	// Varsayılan değer
	rules := v.GetRules()
	if len(rules) != 0 {
		t.Error("Expected empty rules by default")
	}

	// Kuralları ayarla
	v.SetRules([]string{"required", "email"})
	rules = v.GetRules()

	if len(rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(rules))
	}

	// Doğrulayıcı ekle
	v.AddValidator(func(value any) error {
		return nil
	})

	validators := v.GetValidators()
	if len(validators) != 1 {
		t.Errorf("Expected 1 validator, got %d", len(validators))
	}
}

// TestDisplayableMixin, Displayable mixin'ini test eder
func TestDisplayableMixin(t *testing.T) {
	d := &Displayable{}

	// Varsayılan değer
	if d.GetDisplayCallback() != nil {
		t.Error("Expected nil callback by default")
	}

	if d.GetDisplayFormat() != "" {
		t.Error("Expected empty format by default")
	}

	// Format ayarla
	d.SetDisplayFormat("uppercase")
	if d.GetDisplayFormat() != "uppercase" {
		t.Errorf("Expected 'uppercase', got '%s'", d.GetDisplayFormat())
	}

	// Callback ayarla
	callback := func(ctx *context.Context, value any) string {
		return "formatted"
	}
	d.SetDisplayCallback(callback)

	if d.GetDisplayCallback() == nil {
		t.Error("Expected callback to be set")
	}
}

// TestHideableMixin, Hideable mixin'ini test eder
func TestHideableMixin(t *testing.T) {
	h := &Hideable{}

	// Varsayılan değer
	if h.IsHidden() {
		t.Error("Expected not hidden by default")
	}

	// Gizle
	h.SetHidden(true)
	if !h.IsHidden() {
		t.Error("Expected hidden after setting")
	}

	// Görünürlük ayarları
	h.SetShowOnIndex(true)
	if !h.ShowOnIndex() {
		t.Error("Expected show on index")
	}

	h.SetShowOnDetail(false)
	if h.ShowOnDetail() {
		t.Error("Expected not show on detail")
	}

	h.SetShowOnCreate(true)
	if !h.ShowOnCreate() {
		t.Error("Expected show on create")
	}

	h.SetShowOnUpdate(false)
	if h.ShowOnUpdate() {
		t.Error("Expected not show on update")
	}
}

// TestHideableCallback, Hideable callback'ini test eder
func TestHideableCallback(t *testing.T) {
	h := &Hideable{}

	// Varsayılan callback nil
	if h.GetHideCallback() != nil {
		t.Error("Expected nil callback by default")
	}

	// Callback ayarla
	callback := func(ctx *context.Context) bool {
		return true
	}
	h.SetHideCallback(callback)

	if h.GetHideCallback() == nil {
		t.Error("Expected callback to be set")
	}
}

// TestMixinCombination, mixin'leri birlikte kullanmak
func TestMixinCombination(t *testing.T) {
	type AdvancedField struct {
		*Schema
		Searchable
		Sortable
		Filterable
		Validatable
		Displayable
		Hideable
	}

	f := &AdvancedField{
		Schema: &Schema{
			Key:   "test",
			Name:  "Test Field",
			Props: make(map[string]interface{}),
		},
	}

	// Searchable
	f.SetSearchableColumns([]string{"name"})
	if len(f.GetSearchableColumns()) != 1 {
		t.Error("Expected searchable columns")
	}

	// Sortable
	f.Sortable.SetSortable(true)
	if !f.Sortable.IsSortable() {
		t.Error("Expected sortable")
	}

	// Filterable
	f.Filterable.SetFilterable(true)
	if !f.Filterable.IsFilterable() {
		t.Error("Expected filterable")
	}

	// Validatable
	f.SetRules([]string{"required"})
	if len(f.GetRules()) != 1 {
		t.Error("Expected rules")
	}

	// Displayable
	f.SetDisplayFormat("uppercase")
	if f.GetDisplayFormat() != "uppercase" {
		t.Error("Expected display format")
	}

	// Hideable
	f.SetShowOnIndex(true)
	if !f.ShowOnIndex() {
		t.Error("Expected show on index")
	}
}
