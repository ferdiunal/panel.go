package verification

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainVerification "github.com/ferdiunal/panel.go/pkg/domain/verification"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// VerificationResource, Verification entity'si için resource tanımı
type VerificationResource struct {
	resource.OptimizedBase
}

// NewVerificationResource, yeni bir Verification resource'u oluşturur
func NewVerificationResource() *VerificationResource {
	r := &VerificationResource{}

	r.SetModel(&domainVerification.Verification{})
	r.SetSlug("verifications")
	r.SetTitle("Verifications")
	r.SetIcon("shield-check")
	r.SetGroup("System")
	r.SetNavigationOrder(52)
	r.SetVisible(true)

	// Field resolver'ı ayarla
	r.SetFieldResolver(&VerificationFieldResolver{})

	// Card resolver'ı ayarla
	r.SetCardResolver(&VerificationCardResolver{})

	// Policy'yi ayarla
	r.SetPolicy(&VerificationPolicy{})

	return r
}

// Repository, Verification repository'sini döner
func (r *VerificationResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &domainVerification.Verification{})
}

// Lenses, özel görünümleri döner
func (r *VerificationResource) Lenses() []resource.Lens {
	return []resource.Lens{}
}

// GetActions, özel işlemleri döner
func (r *VerificationResource) GetActions() []resource.Action {
	return []resource.Action{}
}

// GetFilters, filtreleri döner
func (r *VerificationResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

// GetSortable, varsayılan sıralama ayarlarını döner
func (r *VerificationResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{
			Column:    "created_at",
			Direction: "desc",
		},
	}
}
