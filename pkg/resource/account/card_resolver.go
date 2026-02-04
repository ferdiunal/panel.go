package account

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// AccountCardResolver, Account card'larını çözer
type AccountCardResolver struct{}

// ResolveCards, Account card'larını döner
func (r *AccountCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
	return []widget.Card{}
}
