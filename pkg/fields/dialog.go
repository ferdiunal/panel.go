package fields

import (
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/gofiber/fiber/v2"
)

// DialogContentType, dialog içeriğinin tipini belirler (form veya wizard)
type DialogContentType string

const (
	// DialogContentForm - Basit form içeriği
	DialogContentForm DialogContentType = "form"
	// DialogContentWizard - Multi-step wizard içeriği
	DialogContentWizard DialogContentType = "wizard"
)

// DialogStep, wizard mode için bir adımı temsil eder
//
// Kullanım Senaryosu:
// - Multi-step form wizard'ları için
// - Her adım kendi field'larına sahip
// - Adımlar atlanabilir veya zorunlu olabilir
//
// Örnek Kullanım:
//
//	step := DialogStep{
//	    Index:       0,
//	    Title:       "Kişisel Bilgiler",
//	    Description: "Lütfen kişisel bilgilerinizi girin",
//	    Fields: []core.Element{
//	        Text("Ad", "name").Required(),
//	        Email("Email", "email").Required(),
//	    },
//	    CanSkip: false,
//	}
type DialogStep struct {
	Index       int            `json:"index"`        // Adım sırası (0'dan başlar)
	Title       string         `json:"title"`        // Adım başlığı
	Description string         `json:"description"`  // Adım açıklaması
	Fields      []core.Element `json:"fields"`       // Adımdaki field'lar
	CanSkip     bool           `json:"can_skip"`     // Adım atlanabilir mi?
}

// DialogField, modal/dialog içinde form veya wizard gösteren bir field tipidir
//
// Kullanım Senaryosu:
// - Kullanıcıdan modal içinde veri toplamak için
// - Multi-step wizard formları için
// - Sayfa geçişlerinde kullanıcıyı bilgilendirmek için
//
// Özellikler:
// - Varsayılan açık veya buton ile tetiklenebilir
// - Basit form veya multi-step wizard mode
// - Özelleştirilebilir dialog boyutu
// - OnComplete ve OnSkip callback'leri
//
// Örnek Kullanım (Basit Form):
//
//	Dialog("Profil Tamamla", "profile_completion").
//	    DefaultOpen(true).
//	    DialogTitle("Profilinizi Tamamlayın").
//	    Content([]core.Element{
//	        Text("Telefon", "phone").Required(),
//	        Text("Adres", "address").Required(),
//	    }).
//	    OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
//	        // Veriyi kaydet
//	        return nil
//	    })
//
// Örnek Kullanım (Wizard):
//
//	Dialog("Onboarding", "onboarding_wizard").
//	    TriggerButton("Başlangıç Rehberini Başlat").
//	    DialogTitle("Hoş Geldiniz!").
//	    Wizard([]DialogStep{
//	        {
//	            Index:   0,
//	            Title:   "Kişisel Bilgiler",
//	            Fields:  []core.Element{Text("Ad", "name").Required()},
//	            CanSkip: false,
//	        },
//	        {
//	            Index:   1,
//	            Title:   "Tercihler",
//	            Fields:  []core.Element{Switch("Bildirimler", "notifications")},
//	            CanSkip: true,
//	        },
//	    })
type DialogField struct {
	*Schema

	// Trigger ayarları
	defaultOpen   bool   // Varsayılan açık mı?
	triggerButton string // Buton metni (boşsa varsayılan açık)
	triggerIcon   string // Buton ikonu

	// Content ayarları
	contentType DialogContentType // "form" veya "wizard"
	fields      []core.Element    // Basit form için fieldlar
	steps       []DialogStep      // Wizard için adımlar

	// Dialog ayarları
	dialogTitle string // Dialog başlığı
	dialogDesc  string // Dialog açıklaması
	dialogSize  string // Dialog boyutu: "sm", "md", "lg", "xl", "full"

	// Callbacks
	onComplete func(ctx *fiber.Ctx, data map[string]any) error // Tamamlandığında çağrılır
	onSkip     func(ctx *fiber.Ctx) error                      // Atlandığında çağrılır (wizard için)
}

// Dialog, yeni bir DialogField oluşturur
//
// Parametreler:
//   - name: Field'ın görüntü adı
//   - key: Veritabanı anahtarı (opsiyonel)
//
// Dönüş Değeri:
//   - *DialogField: Yapılandırılmış dialog field pointer'ı
//
// Örnek Kullanım:
//
//	field := Dialog("Profil Tamamla", "profile_completion")
func Dialog(name string, key ...string) *DialogField {
	schema := NewField(name, key...)
	schema.View = "dialog" // Frontend component adı

	return &DialogField{
		Schema:      schema,
		contentType: DialogContentForm, // Varsayılan: basit form
		dialogSize:  "md",              // Varsayılan: orta boyut
	}
}

// DefaultOpen, dialog'un varsayılan olarak açık olup olmayacağını ayarlar
//
// Parametreler:
//   - open: true ise dialog sayfa yüklendiğinde otomatik açılır
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").DefaultOpen(true)
func (f *DialogField) DefaultOpen(open bool) *DialogField {
	f.defaultOpen = open
	return f
}

// TriggerButton, dialog'u açacak butonun metnini ayarlar
//
// Parametreler:
//   - text: Buton metni
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").TriggerButton("Profili Düzenle")
func (f *DialogField) TriggerButton(text string) *DialogField {
	f.triggerButton = text
	return f
}

// TriggerIcon, dialog butonunun ikonunu ayarlar
//
// Parametreler:
//   - icon: İkon (emoji veya icon class)
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").TriggerIcon("✏️")
func (f *DialogField) TriggerIcon(icon string) *DialogField {
	f.triggerIcon = icon
	return f
}

// Content, basit form içeriğini ayarlar
//
// Parametreler:
//   - fields: Form field'ları
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").Content([]core.Element{
//	    Text("Ad", "name").Required(),
//	    Email("Email", "email").Required(),
//	})
func (f *DialogField) Content(fields []core.Element) *DialogField {
	f.contentType = DialogContentForm
	f.fields = fields
	return f
}

// Wizard, multi-step wizard içeriğini ayarlar
//
// Parametreler:
//   - steps: Wizard adımları
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Onboarding", "onboarding").Wizard([]DialogStep{
//	    {
//	        Index:   0,
//	        Title:   "Adım 1",
//	        Fields:  []core.Element{Text("Ad", "name")},
//	        CanSkip: false,
//	    },
//	})
func (f *DialogField) Wizard(steps []DialogStep) *DialogField {
	f.contentType = DialogContentWizard
	f.steps = steps
	return f
}

// DialogTitle, dialog başlığını ayarlar
//
// Parametreler:
//   - title: Dialog başlığı
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").DialogTitle("Profilinizi Düzenleyin")
func (f *DialogField) DialogTitle(title string) *DialogField {
	f.dialogTitle = title
	return f
}

// DialogDesc, dialog açıklamasını ayarlar
//
// Parametreler:
//   - desc: Dialog açıklaması
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").DialogDesc("Lütfen bilgilerinizi güncelleyin")
func (f *DialogField) DialogDesc(desc string) *DialogField {
	f.dialogDesc = desc
	return f
}

// DialogSize, dialog boyutunu ayarlar
//
// Parametreler:
//   - size: Dialog boyutu ("sm", "md", "lg", "xl", "full")
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").DialogSize("lg")
func (f *DialogField) DialogSize(size string) *DialogField {
	f.dialogSize = size
	return f
}

// OnComplete, dialog tamamlandığında çağrılacak callback'i ayarlar
//
// Parametreler:
//   - fn: Callback fonksiyonu
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Profil", "profile").OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
//	    // Veriyi kaydet
//	    return nil
//	})
func (f *DialogField) OnComplete(fn func(ctx *fiber.Ctx, data map[string]any) error) *DialogField {
	f.onComplete = fn
	return f
}

// OnSkip, dialog atlandığında çağrılacak callback'i ayarlar (wizard için)
//
// Parametreler:
//   - fn: Callback fonksiyonu
//
// Dönüş Değeri:
//   - *DialogField: Method chaining için
//
// Örnek Kullanım:
//
//	Dialog("Onboarding", "onboarding").OnSkip(func(ctx *fiber.Ctx) error {
//	    // Atlandı işaretini kaydet
//	    return nil
//	})
func (f *DialogField) OnSkip(fn func(ctx *fiber.Ctx) error) *DialogField {
	f.onSkip = fn
	return f
}

// JsonSerialize, DialogField'ı JSON uyumlu bir map'e serileştirir
//
// Bu metod, DialogField'ın tüm özelliklerini frontend'e göndermek için
// JSON formatına dönüştürür.
//
// Dönüş Değeri:
//   - map[string]any: JSON uyumlu map
//
// Örnek Kullanım:
//
//	field := Dialog("Profil", "profile").DefaultOpen(true)
//	json := field.JsonSerialize()
//	// json["view"] == "dialog"
//	// json["defaultOpen"] == true
func (f *DialogField) JsonSerialize() map[string]any {
	// Base schema'yı serialize et
	result := f.Schema.JsonSerialize()

	// DialogField özel özelliklerini ekle
	result["defaultOpen"] = f.defaultOpen
	result["triggerButton"] = f.triggerButton
	result["triggerIcon"] = f.triggerIcon
	result["contentType"] = f.contentType
	result["dialogTitle"] = f.dialogTitle
	result["dialogDesc"] = f.dialogDesc
	result["dialogSize"] = f.dialogSize

	// Fields'ı serialize et (basit form için)
	if f.contentType == DialogContentForm && len(f.fields) > 0 {
		serializedFields := make([]map[string]any, len(f.fields))
		for i, field := range f.fields {
			if serializer, ok := field.(interface{ JsonSerialize() map[string]any }); ok {
				serializedFields[i] = serializer.JsonSerialize()
			}
		}
		result["fields"] = serializedFields
	}

	// Steps'i serialize et (wizard için)
	if f.contentType == DialogContentWizard && len(f.steps) > 0 {
		serializedSteps := make([]map[string]any, len(f.steps))
		for i, step := range f.steps {
			serializedFields := make([]map[string]any, len(step.Fields))
			for j, field := range step.Fields {
				if serializer, ok := field.(interface{ JsonSerialize() map[string]any }); ok {
					serializedFields[j] = serializer.JsonSerialize()
				}
			}

			serializedSteps[i] = map[string]any{
				"index":       step.Index,
				"title":       step.Title,
				"description": step.Description,
				"fields":      serializedFields,
				"can_skip":    step.CanSkip,
			}
		}
		result["steps"] = serializedSteps
	}

	return result
}

// GetKey, DialogField'ın key'ini döndürür (core.Element interface implementasyonu)
func (f *DialogField) GetKey() string {
	return f.Schema.GetKey()
}

// GetView, DialogField'ın view tipini döndürür (core.Element interface implementasyonu)
func (f *DialogField) GetView() string {
	return f.Schema.GetView()
}

// GetContext, DialogField'ın context'ini döndürür (core.Element interface implementasyonu)
func (f *DialogField) GetContext() ElementContext {
	return f.Schema.GetContext()
}

// Extract, DialogField için veri çıkarır (core.Element interface implementasyonu)
func (f *DialogField) Extract(resource interface{}) {
	f.Schema.Extract(resource)
}

// Derleme zamanı kontrolü: DialogField'ın core.Element interface'ini implement ettiğinden emin olur
var _ core.Element = (*DialogField)(nil)
