package user

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/data/orm"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// UserResourceWrapper embeds GenericResource to override Repository
// Deprecated: Use NewUserResource() instead
type UserResourceWrapper struct {
	resource.Base
}

func (r UserResourceWrapper) Repository(db *gorm.DB) data.DataProvider {
	return orm.NewUserRepository(db)
}

// GetUserResource, Kullanıcı kaynağının (Resource) konfigürasyonunu döner.
// Deprecated: Use NewUserResource() instead
func GetUserResource() resource.Resource {
	return NewUserResource()
}
