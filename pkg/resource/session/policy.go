package session

import (
	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	domainSession "github.com/ferdiunal/panel.go/pkg/domain/session"
)

// SessionPolicy, Session entity'si için yetkilendirme politikası
type SessionPolicy struct{}

// ViewAny, tüm session'ları görme izni
func (p *SessionPolicy) ViewAny(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// View, belirli bir session'ı görme izni
func (p *SessionPolicy) View(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	session, ok := model.(*domainSession.Session)
	if !ok {
		return false
	}

	return session != nil
}

// Create, session oluşturma izni
func (p *SessionPolicy) Create(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	return true
}

// Update, session güncelleme izni
func (p *SessionPolicy) Update(ctx *context.Context, model any) bool {
	if ctx == nil {
		return false
	}

	session, ok := model.(*domainSession.Session)
	if !ok {
		return false
	}

	return session != nil
}

// Delete, session silme izni
func (p *SessionPolicy) Delete(ctx *context.Context, model any) bool {
	if ctx == nil {
		return true
	}

	session, ok := model.(*domainSession.Session)
	if !ok {
		return false
	}

	return session != nil
}

// Restore, session geri yükleme izni
func (p *SessionPolicy) Restore(ctx *context.Context, model any) bool {
	return false
}

// ForceDelete, session kalıcı silme izni
func (p *SessionPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false
}

var _ auth.Policy = (*SessionPolicy)(nil)
