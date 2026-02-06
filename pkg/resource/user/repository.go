package user

import (
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/data/orm"
	"github.com/ferdiunal/panel.go/pkg/domain/account"
	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserDataProvider extends GormDataProvider to add custom Create logic
type UserDataProvider struct {
	*data.GormDataProvider
	accountRepo account.Repository
}

func NewUserDataProvider(db *gorm.DB) *UserDataProvider {
	return &UserDataProvider{
		GormDataProvider: data.NewGormDataProvider(db, &domainUser.User{}),
		accountRepo:      orm.NewAccountRepository(db),
	}
}

// Create overrides the default create to add Account creation
func (p *UserDataProvider) Create(ctx *context.Context, data map[string]interface{}) (interface{}, error) {
	// 1. Get password from data
	password, ok := data["password"].(string)
	if !ok || password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// 2. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Remove password from data before saving User
	delete(data, "password")

	// 3. Create User using default provider (GORM)
	result, err := p.GormDataProvider.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	user, ok := result.(*domainUser.User)
	if !ok {
		return result, nil
	}

	// 4. Create Account using AccountRepository
	// Note: AccountRepository handles ID generation
	acc := &account.Account{
		UserID:     user.ID,
		ProviderID: "credential",
		AccountID:  nil, // Credentials provider has no external account ID
		Password:   string(hashed),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Use standard context from wrapper
	stdCtx := ctx.Context()
	if err := p.accountRepo.Create(stdCtx, acc); err != nil {
		// Rollback user creation on failure
		_ = p.GormDataProvider.Delete(ctx, fmt.Sprint(user.ID))
		return nil, err
	}

	return user, nil
}
