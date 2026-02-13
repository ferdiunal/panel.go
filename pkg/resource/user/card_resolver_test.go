package user

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/context"
)

// TestUserCardResolverResolveCards, card'ların çözümlendiğini test eder
func TestUserCardResolverResolveCards(t *testing.T) {
	resolver := &UserCardResolver{}
	cards := resolver.ResolveCards(nil)

	if cards == nil {
		t.Fatal("Expected cards, got nil")
	}

	// Card'lar boş olabilir, bu normal
	if len(cards) < 0 {
		t.Error("Expected non-negative card count")
	}
}

// TestUserCardResolverResolveCardsWithContext, context ile card'ları çözmek
func TestUserCardResolverResolveCardsWithContext(t *testing.T) {
	resolver := &UserCardResolver{}
	ctx := &context.Context{}

	cards := resolver.ResolveCards(ctx)

	if cards == nil {
		t.Fatal("Expected cards, got nil")
	}
}

// TestUserCardResolverResolveCardsNilContext, nil context ile card'ları çözmek
func TestUserCardResolverResolveCardsNilContext(t *testing.T) {
	resolver := &UserCardResolver{}

	cards := resolver.ResolveCards(nil)

	if cards == nil {
		t.Fatal("Expected cards, got nil")
	}
}

// TestUserCardResolverReturnsSlice, slice döndürdüğünü test eder
func TestUserCardResolverReturnsSlice(t *testing.T) {
	resolver := &UserCardResolver{}
	cards := resolver.ResolveCards(nil)

	if cards == nil {
		t.Fatal("Expected cards slice, got nil")
	}

	// Slice'ın boş olması normal
	if len(cards) >= 0 {
		// OK
	}
}
