package resource

import (
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// Lens, veri tabanı sorgularını özelleştirerek belirli görünümler (segmentler) oluşturmak için kullanılan arayüzdür.
type Lens interface {
	// Name, Lens'in görünen adı.
	Name() string
	// Slug, URL tanımlayıcısı.
	Slug() string
	// Query, temel sorguyu modifiye eden fonksiyon.
	Query(db *gorm.DB) *gorm.DB
	// Fields, bu lens görünümünde kullanılacak alanlar (opsiyonel).
	// Requirement 13.2: Lens'lerin lens-spesifik alanları filtrelemesine izin ver
	Fields() []fields.Element
	// GetFields, belirli bir bağlamda gösterilecek lens-spesifik alanları döner.
	// Requirement 13.2: Lens'lerin lens-spesifik alanları filtrelemesine izin ver
	GetFields(ctx *appContext.Context) []fields.Element
	// GetCards, lens-spesifik card'ları döner.
	// Requirement 13.3: Lens'lerin lens-spesifik işlemler ve card'lar tanımlamasına izin ver
	GetCards(ctx *appContext.Context) []widget.Card
	// GetName, Lens'in görünen adını döner.
	// Requirement 13.1: Lens'lerin özel sorgu mantığını tanımlamasına izin ver
	GetName() string
	// GetSlug, Lens'in URL tanımlayıcısını döner.
	// Requirement 13.1: Lens'lerin özel sorgu mantığını tanımlamasına izin ver
	GetSlug() string
	// GetQuery, lens-spesifik sorgu mantığını döner.
	// Requirement 13.1: Lens'lerin özel sorgu mantığını tanımlamasına izin ver
	GetQuery() func(*gorm.DB) *gorm.DB
}

// DialogType, kaynak formlarının (ekleme/düzenleme/detay) sunum şeklini belirler.
type DialogType string

const (
	// DialogTypeSheet, formu sağdan açılan bir panel (Sheet) içinde gösterir. (Varsayılan)
	DialogTypeSheet DialogType = "sheet"
	// DialogTypeDrawer, formu alttan açılan bir çekmece (Drawer) içinde gösterir. Özellikle mobil için uygundur.
	DialogTypeDrawer DialogType = "drawer"
	// DialogTypeModal, formu ekranın ortasında klasik bir modal pencere içinde gösterir.
	DialogTypeModal DialogType = "modal"
)
