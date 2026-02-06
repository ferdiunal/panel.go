package user

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// UserResource, kullanıcı kaynağını temsil eder
type UserResource struct {
	resource.OptimizedBase
}

// NewUserResource, yeni bir UserResource oluşturur
func NewUserResource() *UserResource {
	r := &UserResource{}

	// Core ayarları yap
	r.SetModel(&domainUser.User{})
	r.SetSlug("users")
	r.SetTitle("Users")
	r.SetIcon("users")
	r.SetGroup("System")
	r.SetVisible(true)
	r.SetNavigationOrder(1)

	// Yetkilendirme politikasını ayarla
	r.SetPolicy(&UserPolicy{})

	// Resolver'ları ayarla
	r.SetFieldResolver(&UserFieldResolver{})
	r.SetCardResolver(&UserCardResolver{})

	return r
}

// Repository, UserResource için özel repository'yi döner
func (r *UserResource) Repository(db *gorm.DB) data.DataProvider {
	return NewUserDataProvider(db)
}
