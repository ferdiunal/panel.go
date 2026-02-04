package user

import (
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
)

// UserPolicy, Kullanıcı yönetimi için yetkilendirme kurallarını belirler.
// Permissions examples/simple/permissions.toml'den yönetilir.
type UserPolicy struct{}

// ViewAny, kullanıcının listeyi görüntüleyip görüntüleyemeyeceğini belirler.
func (p UserPolicy) ViewAny(ctx *appContext.Context) bool {
	return ctx.HasPermission("users.view_any")
}

// View, kullanıcının belirli bir kaydı görüntüleyip görüntüleyemeyeceğini belirler.
func (p UserPolicy) View(ctx *appContext.Context, model any) bool {
	return ctx.HasPermission("users.view")
}

// Create, yeni bir kullanıcı oluşturma yetkisini kontrol eder.
func (p UserPolicy) Create(ctx *appContext.Context) bool {
	return ctx.HasPermission("users.create")
}

// Update, mevcut bir kullanıcıyı güncelleme yetkisini kontrol eder.
func (p UserPolicy) Update(ctx *appContext.Context, model any) bool {
	return ctx.HasPermission("users.update")
}

// Delete, bir kullanıcıyı silme yetkisini kontrol eder.
// Kendini silmeyi engeller.
func (p UserPolicy) Delete(ctx *appContext.Context, model any) bool {
	// Genel yetki kontrolü (model nil ise)
	if model == nil {
		return true
	}

	userModel, ok := model.(*domainUser.User)
	if !ok {
		return false
	}

	// Context nil ise false döner
	if ctx == nil {
		return false
	}

	authUser := ctx.User()
	if authUser == nil {
		return false
	}

	// Kendini silmeyi engelle
	if userModel.ID == authUser.ID {
		return false
	}

	return true
}
