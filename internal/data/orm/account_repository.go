package orm

import (
	"context"

	"github.com/ferdiunal/panel.go/internal/domain/account"
	"github.com/ferdiunal/panel.go/shared/uuid"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, a *account.Account) error {
	a.ID = uuid.NewUUID().String()
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (*account.Account, error) {
	var a account.Account
	if err := r.db.WithContext(ctx).Preload("User").First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountRepository) FindByProvider(ctx context.Context, providerID, accountID string) (*account.Account, error) {
	var a account.Account
	if err := r.db.WithContext(ctx).Preload("User").First(&a, "provider_id = ? AND account_id = ?", providerID, accountID).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountRepository) FindByUserID(ctx context.Context, userID string) ([]account.Account, error) {
	var accounts []account.Account
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) Update(ctx context.Context, a *account.Account) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&account.Account{}, "id = ?", id).Error
}
