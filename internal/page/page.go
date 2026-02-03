package page

import (
	"github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/fields"
	"github.com/ferdiunal/panel.go/internal/widget"
	"gorm.io/gorm"
)

// Page, panelde gösterilecek özel sayfaları (Dashboard vb.) tanımlar.
type Page interface {
	// Slug, sayfanın URL'deki tanımlayıcısıdır (örn: "dashboard").
	Slug() string
	// Title, menüde ve sayfada görünecek başlık.
	Title() string
	// Cards returns the cards/widgets for the resource dashboard
	Cards() []widget.Card
	// Fields, sayfada gösterilecek form alanlarını döner.
	Fields() []fields.Element
	// Save, sayfa formundan gelen verileri işler.
	Save(c *context.Context, db *gorm.DB, data map[string]interface{}) error
	// Icon, menüde görünecek ikon adı.
	Icon() string
	// Group, menüde hangi grup altında görüneceği.
	Group() string
	// NavigationOrder, menüdeki sıralama önceliğini döner.
	NavigationOrder() int
	// Visible, sayfanın menüde görünüp görünmeyeceğini belirler.
	Visible() bool
}

// Base, Page arayüzünü implemente eden temel yapı.
// Embedding için kullanılabilir.
type Base struct {
}

func (b Base) Slug() string {
	return ""
}

func (b Base) Title() string {
	return ""
}

func (b Base) Icon() string {
	return "circle" // Default icon
}

func (b Base) Group() string {
	return "Genel" // Default group
}

func (b Base) NavigationOrder() int {
	return 99
}

func (b Base) Visible() bool {
	return true
}

func (b Base) Cards() []widget.Card {
	return []widget.Card{}
}

func (b Base) Fields() []fields.Element {
	return []fields.Element{}
}

func (b Base) Save(c *context.Context, db *gorm.DB, data map[string]interface{}) error {
	return nil
}
