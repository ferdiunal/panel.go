package setting

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainSetting "github.com/ferdiunal/panel.go/pkg/domain/setting"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// SettingResource, Setting entity'si için resource tanımı
type SettingResource struct {
	resource.OptimizedBase
}

// NewSettingResource, yeni bir Setting resource'u oluşturur
func NewSettingResource() *SettingResource {
	r := &SettingResource{}

	r.SetModel(&domainSetting.Setting{})
	r.SetSlug("settings")
	r.SetTitle("Settings")
	r.SetIcon("settings")
	r.SetGroup("System")
	r.SetNavigationOrder(53)
	r.SetVisible(true)

	// Field resolver'ı ayarla
	r.SetFieldResolver(&SettingFieldResolver{})

	// Card resolver'ı ayarla
	r.SetCardResolver(&SettingCardResolver{})

	// Policy'yi ayarla
	r.SetPolicy(&SettingPolicy{})

	return r
}

// Repository, Setting repository'sini döner
func (r *SettingResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &domainSetting.Setting{})
}

// Lenses, özel görünümleri döner
func (r *SettingResource) Lenses() []resource.Lens {
	return []resource.Lens{}
}

// GetActions, özel işlemleri döner
func (r *SettingResource) GetActions() []resource.Action {
	return []resource.Action{}
}

// GetFilters, filtreleri döner
func (r *SettingResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

// GetSortable, varsayılan sıralama ayarlarını döner
func (r *SettingResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{
			Column:    "group",
			Direction: "asc",
		},
	}
}
