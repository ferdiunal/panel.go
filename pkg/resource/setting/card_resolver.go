package setting

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// SettingCardResolver, Setting card'larını çözer
type SettingCardResolver struct{}

// ResolveCards, Setting card'larını döner
func (r *SettingCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
	return []widget.Card{}
}
