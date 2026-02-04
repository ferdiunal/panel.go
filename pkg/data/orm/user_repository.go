package orm

import (
	"context"

	pkgContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/shared/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	*data.GormDataProvider
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		GormDataProvider: data.NewGormDataProvider(db, &user.User{}),
		db:               db,
	}
}

// CreateUser implements the domain repository interface (typed)
func (r *UserRepository) CreateUser(ctx context.Context, u *user.User) error {
	u.ID = uuid.NewUUID().String()
	return r.db.WithContext(ctx).Create(u).Error
}

// Create overrides GormDataProvider.Create (generic/resource)
func (r *UserRepository) Create(ctx *pkgContext.Context, data map[string]interface{}) (interface{}, error) {
	// Generate UUID if not present
	if _, ok := data["id"]; !ok || data["id"] == "" {
		data["id"] = uuid.NewUUID().String()
	}
	return r.GormDataProvider.Create(ctx, data)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepository) Update(ctx *pkgContext.Context, id string, data map[string]interface{}) (interface{}, error) {
	return r.GormDataProvider.Update(ctx, id, data)
}

func (r *UserRepository) Delete(ctx *pkgContext.Context, id string) error {
	return r.GormDataProvider.Delete(ctx, id)
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id).Error
}

func (r *UserRepository) Show(ctx *pkgContext.Context, id string) (interface{}, error) {
	return r.GormDataProvider.Show(ctx, id)
}

func (r *UserRepository) Index(ctx *pkgContext.Context, req data.QueryRequest) (*data.QueryResponse, error) {
	return r.GormDataProvider.Index(ctx, req)
}
