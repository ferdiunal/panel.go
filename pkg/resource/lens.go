package resource

import (
	"github.com/ferdiunal/panel.go/pkg/fields"
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
	Fields() []fields.Element
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
