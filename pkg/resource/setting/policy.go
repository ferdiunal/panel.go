package setting

import (
	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	domainSetting "github.com/ferdiunal/panel.go/pkg/domain/setting"
)

// SettingPolicy, Setting entity'si için yetkilendirme politikası
type SettingPolicy struct{}

// ViewAny, tüm setting'leri görme izni
func (p *SettingPolicy) ViewAny(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// View, belirli bir setting'i görme izni
func (p *SettingPolicy) View(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	setting, ok := model.(*domainSetting.Setting)
	if !ok {
		return false
	}

	return setting != nil
}

// Create, setting oluşturma izni
func (p *SettingPolicy) Create(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// Update, setting güncelleme izni
func (p *SettingPolicy) Update(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	setting, ok := model.(*domainSetting.Setting)
	if !ok {
		return false
	}

	return setting != nil
}

// Delete, setting silme izni
func (p *SettingPolicy) Delete(ctx *context.Context, model any) bool {
	if ctx == nil {
		return true
	}

	setting, ok := model.(*domainSetting.Setting)
	if !ok {
		return false
	}

	return setting != nil
}

// Restore, setting geri yükleme izni
func (p *SettingPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

// ForceDelete, setting kalıcı silme izni
func (p *SettingPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}

var _ auth.Policy = (*SettingPolicy)(nil)
