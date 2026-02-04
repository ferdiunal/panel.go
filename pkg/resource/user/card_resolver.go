package user

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// UserCardResolver, kullanıcı card'larını çözer
type UserCardResolver struct{}

// ResolveCards, kullanıcı card'larını döner
func (r *UserCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
	return []widget.Card{
		// Card'lar burada tanımlanabilir
		// Örnek: widget.NewValueCard("Toplam Kullanıcı", "1,234")
	}
}
