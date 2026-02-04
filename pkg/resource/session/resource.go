package session

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainSession "github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// SessionResource, Session entity'si için resource tanımı
type SessionResource struct {
	resource.OptimizedBase
}

// NewSessionResource, yeni bir Session resource'u oluşturur
func NewSessionResource() *SessionResource {
	r := &SessionResource{}

	r.SetModel(&domainSession.Session{})
	r.SetSlug("sessions")
	r.SetTitle("Sessions")
	r.SetIcon("clock")
	r.SetGroup("System")
	r.SetNavigationOrder(51)
	r.SetVisible(true)

	// Field resolver'ı ayarla
	r.SetFieldResolver(&SessionFieldResolver{})

	// Card resolver'ı ayarla
	r.SetCardResolver(&SessionCardResolver{})

	// Policy'yi ayarla
	r.SetPolicy(&SessionPolicy{})

	return r
}

// Repository, Session repository'sini döner
func (r *SessionResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &domainSession.Session{})
}

// With, eager loading yapılacak ilişkileri döner
func (r *SessionResource) With() []string {
	return []string{"User"}
}

// Lenses, özel görünümleri döner
func (r *SessionResource) Lenses() []resource.Lens {
	return []resource.Lens{}
}

// GetActions, özel işlemleri döner
func (r *SessionResource) GetActions() []resource.Action {
	return []resource.Action{}
}

// GetFilters, filtreleri döner
func (r *SessionResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

// GetSortable, varsayılan sıralama ayarlarını döner
func (r *SessionResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{
			Column:    "created_at",
			Direction: "desc",
		},
	}
}
