package user

import (
	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/resource"
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
