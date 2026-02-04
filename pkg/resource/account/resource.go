package account

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainAccount "github.com/ferdiunal/panel.go/pkg/domain/account"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// AccountResource, Account entity'si için resource tanımı
type AccountResource struct {
	resource.OptimizedBase
}

// NewAccountResource, yeni bir Account resource'u oluşturur
func NewAccountResource() *AccountResource {
	r := &AccountResource{}

	r.SetModel(&domainAccount.Account{})
	r.SetSlug("accounts")
	r.SetTitle("Accounts")
	r.SetIcon("key")
	r.SetGroup("System")
	r.SetNavigationOrder(50)
	r.SetVisible(true)

	// Field resolver'ı ayarla
	r.SetFieldResolver(&AccountFieldResolver{})

	// Card resolver'ı ayarla
	r.SetCardResolver(&AccountCardResolver{})

	// Policy'yi ayarla
	r.SetPolicy(&AccountPolicy{})

	return r
}

// Repository, Account repository'sini döner
func (r *AccountResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &domainAccount.Account{})
}

// With, eager loading yapılacak ilişkileri döner
func (r *AccountResource) With() []string {
	return []string{"User"}
}

// Lenses, özel görünümleri döner
func (r *AccountResource) Lenses() []resource.Lens {
	return []resource.Lens{}
}

// GetActions, özel işlemleri döner
func (r *AccountResource) GetActions() []resource.Action {
	return []resource.Action{}
}

// GetFilters, filtreleri döner
func (r *AccountResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

// GetSortable, varsayılan sıralama ayarlarını döner
func (r *AccountResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{
			Column:    "created_at",
			Direction: "desc",
		},
	}
}
