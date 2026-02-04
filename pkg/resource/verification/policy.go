package verification

import (
	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	domainVerification "github.com/ferdiunal/panel.go/pkg/domain/verification"
)

// VerificationPolicy, Verification entity'si için yetkilendirme politikası
type VerificationPolicy struct{}

// ViewAny, tüm verification'ları görme izni
func (p *VerificationPolicy) ViewAny(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// View, belirli bir verification'ı görme izni
func (p *VerificationPolicy) View(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	verification, ok := model.(*domainVerification.Verification)
	if !ok {
		return false
	}

	return verification != nil
}

// Create, verification oluşturma izni
func (p *VerificationPolicy) Create(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// Update, verification güncelleme izni
func (p *VerificationPolicy) Update(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	verification, ok := model.(*domainVerification.Verification)
	if !ok {
		return false
	}

	return verification != nil
}

// Delete, verification silme izni
func (p *VerificationPolicy) Delete(ctx *context.Context, model any) bool {
	if ctx == nil {
		return true
	}

	verification, ok := model.(*domainVerification.Verification)
	if !ok {
		return false
	}

	return verification != nil
}

// Restore, verification geri yükleme izni
func (p *VerificationPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

// ForceDelete, verification kalıcı silme izni
func (p *VerificationPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}

var _ auth.Policy = (*VerificationPolicy)(nil)
