package resource

import (
	"fmt"
	"mime/multipart"

	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// OptimizedResource, Laravel Nova'nın trait pattern'ini Go'ya uyarlayan
// optimize edilmiş resource interface'i.
// Bu interface, daha az metod ile daha fazla işlevsellik sağlar.
type OptimizedResource interface {
	// Core Methods (8 metod)
	Model() any
	Fields() []fields.Element
	Slug() string
	Title() string
	Policy() auth.Policy
	Repository(db *gorm.DB) data.DataProvider
	Cards() []widget.Card
	Visible() bool
}

// FieldResolver, alanları dinamik olarak çözen interface.
type FieldResolver interface {
	ResolveFields(ctx *context.Context) []fields.Element
}

// CardResolver, card'ları dinamik olarak çözen interface.
type CardResolver interface {
	ResolveCards(ctx *context.Context) []widget.Card
}

// FilterResolver, filtreleri dinamik olarak çözen interface.
type FilterResolver interface {
	ResolveFilters(ctx *context.Context) []Filter
}

// LensResolver, lens'leri dinamik olarak çözen interface.
type LensResolver interface {
	ResolveLenses(ctx *context.Context) []Lens
}

// ActionResolver, işlemleri dinamik olarak çözen interface.
type ActionResolver interface {
	ResolveActions(ctx *context.Context) []Action
}

// Authorizable, yetkilendirme işlevselliğini sağlayan mixin.
type Authorizable struct {
	policy auth.Policy
}

// SetPolicy, yetkilendirme politikasını ayarlar.
func (a *Authorizable) SetPolicy(p auth.Policy) {
	a.policy = p
}

// GetPolicy, yetkilendirme politikasını döner.
func (a *Authorizable) GetPolicy() auth.Policy {
	return a.policy
}

// Resolvable, çözümleme işlevselliğini sağlayan mixin.
type Resolvable struct {
	fieldResolver  FieldResolver
	cardResolver   CardResolver
	filterResolver FilterResolver
	lensResolver   LensResolver
	actionResolver ActionResolver
}

// SetFieldResolver, field resolver'ı ayarlar.
func (r *Resolvable) SetFieldResolver(fr FieldResolver) {
	r.fieldResolver = fr
}

// SetCardResolver, card resolver'ı ayarlar.
func (r *Resolvable) SetCardResolver(cr CardResolver) {
	r.cardResolver = cr
}

// SetFilterResolver, filter resolver'ı ayarlar.
func (r *Resolvable) SetFilterResolver(fr FilterResolver) {
	r.filterResolver = fr
}

// SetLensResolver, lens resolver'ı ayarlar.
func (r *Resolvable) SetLensResolver(lr LensResolver) {
	r.lensResolver = lr
}

// SetActionResolver, action resolver'ı ayarlar.
func (r *Resolvable) SetActionResolver(ar ActionResolver) {
	r.actionResolver = ar
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

// ResolveFilters, filtreleri çözer.
func (r *Resolvable) ResolveFilters(ctx *context.Context) []Filter {
	if r.filterResolver != nil {
		return r.filterResolver.ResolveFilters(ctx)
	}
	return []Filter{}
}

// ResolveLenses, lens'leri çözer.
func (r *Resolvable) ResolveLenses(ctx *context.Context) []Lens {
	if r.lensResolver != nil {
		return r.lensResolver.ResolveLenses(ctx)
	}
	return []Lens{}
}

// ResolveActions, işlemleri çözer.
func (r *Resolvable) ResolveActions(ctx *context.Context) []Action {
	if r.actionResolver != nil {
		return r.actionResolver.ResolveActions(ctx)
	}
	return []Action{}
}

// Navigable, navigasyon işlevselliğini sağlayan mixin.
type Navigable struct {
	icon            string
	group           string
	navigationOrder int
	visible         bool
	dialogType      DialogType
	sortable        []Sortable
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

// SetDialogType, dialog tipini ayarlar.
func (n *Navigable) SetDialogType(dt DialogType) {
	n.dialogType = dt
}

// GetDialogType, dialog tipini döner.
func (n *Navigable) GetDialogType() DialogType {
	return n.dialogType
}

// SetSortable, sıralanabilir alanları ayarlar.
func (n *Navigable) SetSortable(sortable []Sortable) {
	n.sortable = sortable
}

// GetSortable, sıralanabilir alanları döner.
func (n *Navigable) GetSortable() []Sortable {
	return n.sortable
}

// OptimizedBase, OptimizedResource'u implement eden temel struct.
// Embedding için kullanılabilir.
type OptimizedBase struct {
	Authorizable
	Resolvable
	Navigable
	model      any
	slug       string
	title      string
	repository data.DataProvider
	cards      []widget.Card
}

// SetModel, model'i ayarlar.
func (b *OptimizedBase) SetModel(m any) {
	b.model = m
}

// Model, model'i döner.
func (b *OptimizedBase) Model() any {
	return b.model
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

// SetRepository, repository'yi ayarlar.
func (b *OptimizedBase) SetRepository(r data.DataProvider) {
	b.repository = r
}

// Repository, repository'yi döner.
func (b *OptimizedBase) Repository(db *gorm.DB) data.DataProvider {
	return b.repository
}

// SetCards, card'ları ayarlar.
func (b *OptimizedBase) SetCards(c []widget.Card) {
	b.cards = c
}

// Cards, card'ları döner.
func (b *OptimizedBase) Cards() []widget.Card {
	return b.cards
}

// Fields, alanları döner (ResolveFields'i kullanır).
func (b *OptimizedBase) Fields() []fields.Element {
	return b.ResolveFields(nil)
}

// Policy, politikayı döner.
func (b *OptimizedBase) Policy() auth.Policy {
	return b.GetPolicy()
}

// Visible, görünürlüğü döner.
func (b *OptimizedBase) Visible() bool {
	return b.IsVisible()
}

// With, eager loading yapılacak ilişkileri döner.
func (b *OptimizedBase) With() []string {
	return []string{}
}

// Lenses, tanımlı özel görünümleri döner.
func (b *OptimizedBase) Lenses() []Lens {
	return []Lens{}
}

// GetLenses, kaynağın tüm lens'lerini döner.
func (b *OptimizedBase) GetLenses() []Lens {
	return b.Lenses()
}

// Icon, menü ikonunu döner.
func (b *OptimizedBase) Icon() string {
	return b.GetIcon()
}

// Group, menü grubunu döner.
func (b *OptimizedBase) Group() string {
	return b.GetGroup()
}

// GetSortable, varsayılan sıralama ayarlarını döner.
func (b *OptimizedBase) GetSortable() []Sortable {
	return b.Navigable.sortable
}

// GetDialogType, diyalog tipini döner.
func (b *OptimizedBase) GetDialogType() DialogType {
	return b.Navigable.GetDialogType()
}

// SetDialogType, form görünüm tipini ayarlar.
func (b *OptimizedBase) SetDialogType(dt DialogType) Resource {
	b.Navigable.SetDialogType(dt)
	return b
}

// GetFields, belirli bir bağlamda gösterilecek alanları döner.
func (b *OptimizedBase) GetFields(ctx *context.Context) []fields.Element {
	return b.ResolveFields(ctx)
}

// GetCards, belirli bir bağlamda gösterilecek card'ları döner.
func (b *OptimizedBase) GetCards(ctx *context.Context) []widget.Card {
	return b.ResolveCards(ctx)
}

// GetPolicy, kaynağın yetkilendirme politikasını döner.
func (b *OptimizedBase) GetPolicy() auth.Policy {
	return b.Authorizable.GetPolicy()
}

// ResolveField, bir alanın değerini dinamik olarak hesaplayan ve dönüştüren fonksiyon.
func (b *OptimizedBase) ResolveField(fieldName string, item any) (any, error) {
	for _, field := range b.Fields() {
		if field.GetKey() == fieldName {
			field.Extract(item)
			serialized := field.JsonSerialize()
			if val, ok := serialized["value"]; ok {
				return val, nil
			}
			return nil, nil
		}
	}
	return nil, fmt.Errorf("field %s not found", fieldName)
}

// GetActions, kaynağın özel işlemlerini döner.
func (b *OptimizedBase) GetActions() []Action {
	return []Action{}
}

// GetFilters, kaynağın filtreleme seçeneklerini döner.
func (b *OptimizedBase) GetFilters() []Filter {
	return []Filter{}
}

// StoreHandler, dosya yükleme işlemlerini yönetir.
func (b *OptimizedBase) StoreHandler(c *context.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	return "", nil
}

// NavigationOrder, menüdeki sıralama önceliğini döner.
func (b *OptimizedBase) NavigationOrder() int {
	return b.GetNavigationOrder()
}
