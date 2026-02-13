package fields

import (
	"github.com/ferdiunal/panel.go/pkg/core"
)

// Tab, bir tab'ın yapısını temsil eder.
//
// Her tab, benzersiz bir değer (value), görünen bir etiket (label) ve
// içerdiği alanlardan (fields) oluşur.
//
// # Kullanım Örneği
//
//	tab := Tab{
//	    Value: "tr",
//	    Label: "Türkçe",
//	    Fields: []core.Element{
//	        fields.Text("Başlık", "title_tr"),
//	        fields.Textarea("Açıklama", "description_tr"),
//	    },
//	}
type Tab struct {
	// Value, tab'ın benzersiz tanımlayıcısıdır (örn: "tr", "en", "general")
	Value string `json:"value"`

	// Label, tab'ın görünen adıdır (örn: "Türkçe", "English", "Genel Bilgiler")
	Label string `json:"label"`

	// Fields, tab içinde görüntülenecek alanların listesidir
	Fields []core.Element `json:"fields"`
}

// TabsField, alanları tab'lara ayırmak için bir konteyner temsil eder.
//
// TabsField, ilgili alanları tab'lar halinde organize etmek için kullanılır.
// Her tab, kendi başlığı ve içeriği ile ayrı bir bölüm oluşturur.
//
// # Kullanım Senaryoları
//
// - **Çoklu Dil Desteği**: Türkçe, İngilizce, vb. tab'ları ile çeviri alanları
// - **Kategorize Edilmiş Formlar**: Genel Bilgiler, Adres, İletişim tab'ları
// - **Karmaşık Form Organizasyonu**: Uzun formları mantıksal bölümlere ayırma
// - **İlgili Alan Grupları**: Benzer alanları bir arada gösterme
//
// # Özellikler
//
// - **Çoklu Tab**: Birden fazla tab ekleme
// - **Tab Pozisyonu**: Tab'ların konumu (top, bottom, left, right)
// - **Tab Variant**: Tab görünümü (default, line)
// - **Varsayılan Tab**: Sayfa yüklendiğinde aktif olacak tab
//
// # Kullanım Örneği
//
//	tabs := fields.Tabs("Ürün Bilgileri").
//	    AddTab("tr", "Türkçe",
//	        fields.Text("Başlık", "title_tr"),
//	        fields.Textarea("Açıklama", "description_tr"),
//	    ).
//	    AddTab("en", "English",
//	        fields.Text("Title", "title_en"),
//	        fields.Textarea("Description", "description_en"),
//	    ).
//	    WithSide("top").
//	    WithVariant("line").
//	    WithDefaultTab("tr")
type TabsField struct {
	*Schema
	Tabs []Tab
}

// Tabs, alanları tab'lara ayırmak için yeni bir tabs konteyner oluşturur.
//
// Bu fonksiyon, ilgili alanları tab'lar halinde organize etmek için bir konteyner oluşturur.
// Form sayfalarında alanları mantıksal tab'lara ayırmak için kullanılır.
//
// # Parametreler
//
// - **title**: Tabs konteyner başlığı (örn. "Ürün Bilgileri", "Çeviriler")
//
// # Kullanım Örneği
//
//	tabs := fields.Tabs("Ürün Bilgileri")
//
// Döndürür:
//   - Yapılandırılmış TabsField pointer'ı
func Tabs(title string) *TabsField {
	schema := NewField(title)
	schema.View = "tabs-field"
	schema.Type = TYPE_TABS

	return &TabsField{
		Schema: schema,
		Tabs:   []Tab{},
	}
}

// AddTab, tabs konteyner'a yeni bir tab ekler.
//
// Bu metod, tab'a benzersiz bir değer (value), görünen bir etiket (label) ve
// içerdiği alanları (fields) ekler.
//
// # Parametreler
//
// - **value**: Tab'ın benzersiz tanımlayıcısı (örn: "tr", "en", "general")
// - **label**: Tab'ın görünen adı (örn: "Türkçe", "English", "Genel Bilgiler")
// - **fields**: Tab içinde görüntülenecek alanlar
//
// # Kullanım Örneği
//
//	tabs.AddTab("tr", "Türkçe",
//	    fields.Text("Başlık", "title_tr"),
//	    fields.Textarea("Açıklama", "description_tr"),
//	)
//
// Döndürür:
//   - TabsField pointer'ı (method chaining için)
func (t *TabsField) AddTab(value, label string, fields ...core.Element) *TabsField {
	t.Tabs = append(t.Tabs, Tab{
		Value:  value,
		Label:  label,
		Fields: fields,
	})
	return t
}

// WithSide, tab'ların pozisyonunu ayarlar.
//
// Bu metod, tab'ların nerede görüntüleneceğini belirler.
// Geçerli değerler: "top", "bottom", "left", "right"
//
// # Parametreler
//
// - **side**: Tab pozisyonu ("top", "bottom", "left", "right")
//
// # Kullanım Örneği
//
//	tabs.WithSide("top")    // Tab'lar üstte
//	tabs.WithSide("left")   // Tab'lar solda
//
// Döndürür:
//   - TabsField pointer'ı (method chaining için)
func (t *TabsField) WithSide(side string) *TabsField {
	t.Props["side"] = side
	return t
}

// WithVariant, tab'ların görünüm stilini ayarlar.
//
// Bu metod, tab'ların nasıl görüntüleneceğini belirler.
// Geçerli değerler: "default", "line"
//
// # Parametreler
//
// - **variant**: Tab görünümü ("default", "line")
//
// # Kullanım Örneği
//
//	tabs.WithVariant("default")  // Varsayılan tab görünümü
//	tabs.WithVariant("line")     // Çizgi altı tab görünümü
//
// Döndürür:
//   - TabsField pointer'ı (method chaining için)
func (t *TabsField) WithVariant(variant string) *TabsField {
	t.Props["variant"] = variant
	return t
}

// WithDefaultTab, sayfa yüklendiğinde aktif olacak tab'ı ayarlar.
//
// Bu metod, varsayılan olarak hangi tab'ın açık olacağını belirler.
// Belirtilen value, AddTab ile eklenen tab'lardan birinin value'su olmalıdır.
//
// # Parametreler
//
// - **value**: Varsayılan aktif tab'ın value'su
//
// # Kullanım Örneği
//
//	tabs.WithDefaultTab("tr")  // Türkçe tab'ı varsayılan olarak açık
//
// Döndürür:
//   - TabsField pointer'ı (method chaining için)
func (t *TabsField) WithDefaultTab(value string) *TabsField {
	t.Props["defaultTab"] = value
	return t
}

// GetFields, tüm tab'lardaki alanları düz bir liste olarak döndürür.
//
// Bu metod, tabs konteyner içindeki tüm alanları toplar ve tek bir liste halinde döndürür.
// Backend'de alan validasyonu ve veri işleme için kullanılır.
//
// # Kullanım Örneği
//
//	allFields := tabs.GetFields()
//	// Tüm tab'lardaki field'ları içerir
//
// Döndürür:
//   - Tüm tab'lardaki alanların birleştirilmiş listesi
func (t *TabsField) GetFields() []core.Element {
	var fields []core.Element
	for _, tab := range t.Tabs {
		fields = append(fields, tab.Fields...)
	}
	return fields
}
