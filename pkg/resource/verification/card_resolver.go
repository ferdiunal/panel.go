package verification

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// VerificationCardResolver, Verification card'larını çözer
type VerificationCardResolver struct{}

// ResolveCards, Verification card'larını döner
func (r *VerificationCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
	return []widget.Card{}
}
