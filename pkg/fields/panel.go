package fields

import (
	"github.com/ferdiunal/panel.go/pkg/core"
)

// PanelField, alanları bölümlere/kartlara gruplamak için bir konteyner temsil eder.
//
// Panel, ilgili alanları görsel olarak gruplamak ve organize etmek için kullanılır.
// Form sayfalarında alanları mantıksal bölümlere ayırmak için idealdir.
//
// # Kullanım Senaryoları
//
// - **Profil Bilgileri**: Ad, soyad, e-posta gibi kişisel bilgileri grupla
// - **Adres Bilgileri**: Şehir, ilçe, posta kodu gibi adres alanlarını grupla
// - **Güvenlik Ayarları**: Şifre, 2FA gibi güvenlik alanlarını grupla
//
// # Özellikler
//
// - **Başlık**: Panel başlığı
// - **Açıklama**: Panel açıklaması (opsiyonel)
// - **Sütun Sayısı**: Grid sütun sayısı (1-4)
// - **Daraltılabilir**: Panel daraltılabilir/genişletilebilir
// - **Varsayılan Daraltılmış**: Panel varsayılan olarak daraltılmış
//
// # Kullanım Örneği
//
//	panel := fields.Panel("Kişisel Bilgiler",
//	    fields.Text("Ad", "first_name"),
//	    fields.Text("Soyad", "last_name"),
//	    fields.Email("E-posta", "email"),
//	).WithDescription("Kullanıcının kişisel bilgileri").
//	  WithColumns(2).
//	  Collapsible()
type PanelField struct {
	*Schema
	Fields []core.Element
}

// Panel, alanları gruplamak için yeni bir panel/bölüm oluşturur.
//
// Bu fonksiyon, ilgili alanları görsel olarak gruplamak için bir konteyner oluşturur.
// Form sayfalarında alanları mantıksal bölümlere ayırmak için kullanılır.
//
// # Parametreler
//
// - **title**: Panel başlığı (örn. "Kişisel Bilgiler", "Adres Bilgileri")
// - **fields**: Panel içinde görüntülenecek alanlar
//
// # Kullanım Örneği
//
//	panel := fields.Panel("Kişisel Bilgiler",
//	    fields.Text("Ad", "first_name"),
//	    fields.Text("Soyad", "last_name"),
//	    fields.Email("E-posta", "email"),
//	)
//
// Döndürür:
//   - Yapılandırılmış PanelField pointer'ı
func Panel(title string, fields ...core.Element) *PanelField {
	schema := NewField(title)
	schema.View = "panel-field"
	schema.Type = TYPE_PANEL

	return &PanelField{
		Schema: schema,
		Fields: fields,
	}
}

// WithDescription, panel'e açıklama ekler.
//
// Bu metod, panel başlığının altında görüntülenecek açıklayıcı bir metin ekler.
// Kullanıcılara panel içeriği hakkında bilgi vermek için kullanılır.
//
// # Parametreler
//
// - **description**: Panel açıklaması
//
// # Kullanım Örneği
//
//	panel.WithDescription("Kullanıcının kişisel bilgilerini girin")
//
// Döndürür:
//   - PanelField pointer'ı (method chaining için)
func (p *PanelField) WithDescription(description string) *PanelField {
	p.Props["description"] = description
	return p
}

// WithColumns, panel için grid sütun sayısını ayarlar (1-4).
//
// Bu metod, panel içindeki alanların kaç sütunda görüntüleneceğini belirler.
// Sütun sayısı 1 ile 4 arasında olmalıdır.
//
// # Parametreler
//
// - **columns**: Sütun sayısı (1-4 arası)
//
// # Kullanım Örneği
//
//	panel.WithColumns(2)  // Alanlar 2 sütunda görüntülenir
//	panel.WithColumns(3)  // Alanlar 3 sütunda görüntülenir
//
// # Önemli Notlar
//
// - Sütun sayısı 1'den küçükse 1'e ayarlanır
// - Sütun sayısı 4'ten büyükse 4'e ayarlanır
//
// Döndürür:
//   - PanelField pointer'ı (method chaining için)
func (p *PanelField) WithColumns(columns int) *PanelField {
	if columns < 1 {
		columns = 1
	}
	if columns > 4 {
		columns = 4
	}
	p.Props["columns"] = columns
	return p
}

// Collapsible, panel'i daraltılabilir yapar.
//
// Bu metod, panel başlığına tıklandığında panel içeriğinin daraltılıp genişletilebilmesini sağlar.
// Uzun formlarda alan tasarrufu sağlamak için kullanılır.
//
// # Kullanım Örneği
//
//	panel.Collapsible()
//	// Panel başlığına tıklandığında içerik daraltılır/genişletilir
//
// Döndürür:
//   - PanelField pointer'ı (method chaining için)
func (p *PanelField) Collapsible() *PanelField {
	p.Props["collapsible"] = true
	return p
}

// DefaultCollapsed, panel'in varsayılan olarak daraltılmış olmasını sağlar.
//
// Bu metod, panel'in sayfa yüklendiğinde daraltılmış olarak görüntülenmesini sağlar.
// Kullanıcı panel başlığına tıklayarak içeriği genişletebilir.
//
// # Kullanım Örneği
//
//	panel.DefaultCollapsed()
//	// Panel sayfa yüklendiğinde daraltılmış olarak görüntülenir
//
// # Önemli Notlar
//
// - Bu metod otomatik olarak Collapsible() özelliğini de aktif eder
//
// Döndürür:
//   - PanelField pointer'ı (method chaining için)
func (p *PanelField) DefaultCollapsed() *PanelField {
	p.Props["collapsible"] = true
	p.Props["defaultCollapsed"] = true
	return p
}

// GetFields, bu panel içindeki alanları döndürür.
//
// Döndürür:
//   - Panel içindeki alanların listesi
func (p *PanelField) GetFields() []core.Element {
	return p.Fields
}
