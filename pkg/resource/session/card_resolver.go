package session

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// SessionCardResolver, Session card'larını çözer
type SessionCardResolver struct{}

// ResolveCards, Session card'larını döner
func (r *SessionCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
	return []widget.Card{}
}
