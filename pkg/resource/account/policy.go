package account

import (
	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	domainAccount "github.com/ferdiunal/panel.go/pkg/domain/account"
)

// AccountPolicy, Account entity'si için yetkilendirme politikası
type AccountPolicy struct{}

// ViewAny, tüm account'ları görme izni
func (p *AccountPolicy) ViewAny(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// View, belirli bir account'ı görme izni
func (p *AccountPolicy) View(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	account, ok := model.(*domainAccount.Account)
	if !ok {
		return false
	}

	return account != nil
}

// Create, account oluşturma izni
func (p *AccountPolicy) Create(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// Update, account güncelleme izni
func (p *AccountPolicy) Update(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	account, ok := model.(*domainAccount.Account)
	if !ok {
		return false
	}

	return account != nil
}

// Delete, account silme izni
func (p *AccountPolicy) Delete(ctx *context.Context, model any) bool {
	if ctx == nil {
		return true
	}

	account, ok := model.(*domainAccount.Account)
	if !ok {
		return false
	}

	return account != nil
}

// Restore, account geri yükleme izni
func (p *AccountPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

// ForceDelete, account kalıcı silme izni
func (p *AccountPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}

var _ auth.Policy = (*AccountPolicy)(nil)
