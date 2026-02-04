package page

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// FieldResolver, sayfanın alanlarını dinamik olarak çözen interface.
type FieldResolver interface {
	ResolveFields(ctx *context.Context) []fields.Element
}

// CardResolver, sayfanın card'larını dinamik olarak çözen interface.
type CardResolver interface {
	ResolveCards(ctx *context.Context) []widget.Card
}

// Resolvable, sayfa çözümleme işlevselliğini sağlayan mixin.
type Resolvable struct {
	fieldResolver FieldResolver
	cardResolver  CardResolver
}

// SetFieldResolver, field resolver'ı ayarlar.
func (r *Resolvable) SetFieldResolver(fr FieldResolver) {
	r.fieldResolver = fr
}

// SetCardResolver, card resolver'ı ayarlar.
func (r *Resolvable) SetCardResolver(cr CardResolver) {
	r.cardResolver = cr
}

// ResolveFields, alanları çözer.
func (r *Resolvable) ResolveFields(ctx *context.Context) []fields.Element {
	if r.fieldResolver != nil {
		return r.fieldResolver.ResolveFields(ctx)
	}
	return []fields.Element{}
}

// ResolveCards, card'ları çözer.
func (r *Resolvable) ResolveCards(ctx *context.Context) []widget.Card {
	if r.cardResolver != nil {
		return r.cardResolver.ResolveCards(ctx)
	}
	return []widget.Card{}
}

// Navigable, sayfa navigasyon işlevselliğini sağlayan mixin.
type Navigable struct {
	icon            string
	group           string
	navigationOrder int
	visible         bool
}

// SetIcon, ikon adını ayarlar.
func (n *Navigable) SetIcon(icon string) {
	n.icon = icon
}

// GetIcon, ikon adını döner.
func (n *Navigable) GetIcon() string {
	return n.icon
}

// SetGroup, grup adını ayarlar.
func (n *Navigable) SetGroup(group string) {
	n.group = group
}

// GetGroup, grup adını döner.
func (n *Navigable) GetGroup() string {
	return n.group
}

// SetNavigationOrder, navigasyon sırasını ayarlar.
func (n *Navigable) SetNavigationOrder(order int) {
	n.navigationOrder = order
}

// GetNavigationOrder, navigasyon sırasını döner.
func (n *Navigable) GetNavigationOrder() int {
	return n.navigationOrder
}

// SetVisible, görünürlüğü ayarlar.
func (n *Navigable) SetVisible(visible bool) {
	n.visible = visible
}

// IsVisible, görünürlüğü döner.
func (n *Navigable) IsVisible() bool {
	return n.visible
}

// OptimizedBase, Page arayüzünü implement eden temel struct.
// Embedding için kullanılabilir.
type OptimizedBase struct {
	Resolvable
	Navigable
	slug  string
	title string
}

// SetSlug, slug'ı ayarlar.
func (b *OptimizedBase) SetSlug(s string) {
	b.slug = s
}

// Slug, slug'ı döner.
func (b *OptimizedBase) Slug() string {
	return b.slug
}

// SetTitle, başlığı ayarlar.
func (b *OptimizedBase) SetTitle(t string) {
	b.title = t
}

// Title, başlığı döner.
func (b *OptimizedBase) Title() string {
	return b.title
}

// Fields, alanları döner (ResolveFields'i kullanır).
func (b *OptimizedBase) Fields() []fields.Element {
	return b.ResolveFields(nil)
}

// Cards, card'ları döner (ResolveCards'i kullanır).
func (b *OptimizedBase) Cards() []widget.Card {
	return b.ResolveCards(nil)
}

// Icon, ikon adını döner.
func (b *OptimizedBase) Icon() string {
	return b.GetIcon()
}

// Group, grup adını döner.
func (b *OptimizedBase) Group() string {
	return b.GetGroup()
}

// NavigationOrder, navigasyon sırasını döner.
func (b *OptimizedBase) NavigationOrder() int {
	return b.GetNavigationOrder()
}

// Visible, görünürlüğü döner.
func (b *OptimizedBase) Visible() bool {
	return b.IsVisible()
}
