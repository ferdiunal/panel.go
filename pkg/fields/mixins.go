package fields

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"gorm.io/gorm"
)

// Searchable, alanın aranabilir olmasını sağlayan mixin.
type Searchable struct {
	searchableColumns []string
	searchCallback    func(*gorm.DB, string) *gorm.DB
}

// SetSearchableColumns, aranabilir sütunları ayarlar.
func (s *Searchable) SetSearchableColumns(columns []string) {
	s.searchableColumns = columns
}

// GetSearchableColumns, aranabilir sütunları döner.
func (s *Searchable) GetSearchableColumns() []string {
	if s.searchableColumns == nil {
		return []string{}
	}
	return s.searchableColumns
}

// SetSearchCallback, özel arama callback'ini ayarlar.
func (s *Searchable) SetSearchCallback(cb func(*gorm.DB, string) *gorm.DB) {
	s.searchCallback = cb
}

// GetSearchCallback, arama callback'ini döner.
func (s *Searchable) GetSearchCallback() func(*gorm.DB, string) *gorm.DB {
	return s.searchCallback
}

// Sortable, alanın sıralanabilir olmasını sağlayan mixin.
type Sortable struct {
	sortable      bool
	sortDirection string
	sortCallback  func(*gorm.DB, string) *gorm.DB
}

// SetSortable, alanın sıralanabilir olup olmadığını ayarlar.
func (s *Sortable) SetSortable(sortable bool) {
	s.sortable = sortable
}

// IsSortable, alanın sıralanabilir olup olmadığını döner.
func (s *Sortable) IsSortable() bool {
	return s.sortable
}

// SetSortDirection, varsayılan sıralama yönünü ayarlar.
func (s *Sortable) SetSortDirection(direction string) {
	s.sortDirection = direction
}

// GetSortDirection, varsayılan sıralama yönünü döner.
func (s *Sortable) GetSortDirection() string {
	if s.sortDirection == "" {
		return "asc"
	}
	return s.sortDirection
}

// SetSortCallback, özel sıralama callback'ini ayarlar.
func (s *Sortable) SetSortCallback(cb func(*gorm.DB, string) *gorm.DB) {
	s.sortCallback = cb
}

// GetSortCallback, sıralama callback'ini döner.
func (s *Sortable) GetSortCallback() func(*gorm.DB, string) *gorm.DB {
	return s.sortCallback
}

// Filterable, alanın filtrelenebilir olmasını sağlayan mixin.
type Filterable struct {
	filterable     bool
	filterCallback func(*gorm.DB, any) *gorm.DB
	filterOptions  map[string]string
}

// SetFilterable, alanın filtrelenebilir olup olmadığını ayarlar.
func (f *Filterable) SetFilterable(filterable bool) {
	f.filterable = filterable
}

// IsFilterable, alanın filtrelenebilir olup olmadığını döner.
func (f *Filterable) IsFilterable() bool {
	return f.filterable
}

// SetFilterCallback, özel filtreleme callback'ini ayarlar.
func (f *Filterable) SetFilterCallback(cb func(*gorm.DB, any) *gorm.DB) {
	f.filterCallback = cb
}

// GetFilterCallback, filtreleme callback'ini döner.
func (f *Filterable) GetFilterCallback() func(*gorm.DB, any) *gorm.DB {
	return f.filterCallback
}

// SetFilterOptions, filtreleme seçeneklerini ayarlar.
func (f *Filterable) SetFilterOptions(options map[string]string) {
	f.filterOptions = options
}

// GetFilterOptions, filtreleme seçeneklerini döner.
func (f *Filterable) GetFilterOptions() map[string]string {
	if f.filterOptions == nil {
		return make(map[string]string)
	}
	return f.filterOptions
}

// Validatable, alanın doğrulanabilir olmasını sağlayan mixin.
type Validatable struct {
	rules      []string
	validators []func(any) error
}

// SetRules, doğrulama kurallarını ayarlar.
func (v *Validatable) SetRules(rules []string) {
	v.rules = rules
}

// GetRules, doğrulama kurallarını döner.
func (v *Validatable) GetRules() []string {
	if v.rules == nil {
		return []string{}
	}
	return v.rules
}

// AddValidator, özel doğrulayıcı ekler.
func (v *Validatable) AddValidator(validator func(any) error) {
	v.validators = append(v.validators, validator)
}

// GetValidators, doğrulayıcıları döner.
func (v *Validatable) GetValidators() []func(any) error {
	if v.validators == nil {
		return []func(any) error{}
	}
	return v.validators
}

// Displayable, alanın görüntülenme şeklini özelleştiren mixin.
type Displayable struct {
	displayCallback func(*context.Context, any) string
	displayFormat   string
}

// SetDisplayCallback, özel görüntüleme callback'ini ayarlar.
func (d *Displayable) SetDisplayCallback(cb func(*context.Context, any) string) {
	d.displayCallback = cb
}

// GetDisplayCallback, görüntüleme callback'ini döner.
func (d *Displayable) GetDisplayCallback() func(*context.Context, any) string {
	return d.displayCallback
}

// SetDisplayFormat, görüntüleme formatını ayarlar.
func (d *Displayable) SetDisplayFormat(format string) {
	d.displayFormat = format
}

// GetDisplayFormat, görüntüleme formatını döner.
func (d *Displayable) GetDisplayFormat() string {
	return d.displayFormat
}

// Hideable, alanın gizlenebilir olmasını sağlayan mixin.
type Hideable struct {
	hidden       bool
	hideCallback func(*context.Context) bool
	showOnIndex  bool
	showOnDetail bool
	showOnCreate bool
	showOnUpdate bool
}

// SetHidden, alanın gizli olup olmadığını ayarlar.
func (h *Hideable) SetHidden(hidden bool) {
	h.hidden = hidden
}

// IsHidden, alanın gizli olup olmadığını döner.
func (h *Hideable) IsHidden() bool {
	return h.hidden
}

// SetHideCallback, özel gizleme callback'ini ayarlar.
func (h *Hideable) SetHideCallback(cb func(*context.Context) bool) {
	h.hideCallback = cb
}

// GetHideCallback, gizleme callback'ini döner.
func (h *Hideable) GetHideCallback() func(*context.Context) bool {
	return h.hideCallback
}

// SetShowOnIndex, alanın liste görünümünde gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnIndex(show bool) {
	h.showOnIndex = show
}

// ShowOnIndex, alanın liste görünümünde gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnIndex() bool {
	return h.showOnIndex
}

// SetShowOnDetail, alanın detay görünümünde gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnDetail(show bool) {
	h.showOnDetail = show
}

// ShowOnDetail, alanın detay görünümünde gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnDetail() bool {
	return h.showOnDetail
}

// SetShowOnCreate, alanın oluşturma formunda gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnCreate(show bool) {
	h.showOnCreate = show
}

// ShowOnCreate, alanın oluşturma formunda gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnCreate() bool {
	return h.showOnCreate
}

// SetShowOnUpdate, alanın güncelleme formunda gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnUpdate(show bool) {
	h.showOnUpdate = show
}

// ShowOnUpdate, alanın güncelleme formunda gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnUpdate() bool {
	return h.showOnUpdate
}
