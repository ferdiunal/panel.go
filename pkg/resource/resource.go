package resource

import (
	"mime/multipart"

	"github.com/ferdiunal/panel.go/pkg/auth"
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// Action, kaynakta gerçekleştirilebilecek özel işlemleri temsil eder.
// Örneğin: "Activate", "Deactivate", "Send Email" vb.
type Action interface {
	// GetName, işlemin görünen adını döner.
	GetName() string
	// GetSlug, işlemin URL tanımlayıcısını döner.
	GetSlug() string
	// GetIcon, işlemin ikonunu döner.
	GetIcon() string
	// Execute, işlemi gerçekleştirir.
	Execute(ctx *appContext.Context, items []any) error
}

// Filter, kaynakta uygulanabilecek filtreleme seçeneklerini temsil eder.
// Örneğin: "Status", "Date Range", "Category" vb.
type Filter interface {
	// GetName, filtrenin görünen adını döner.
	GetName() string
	// GetSlug, filtrenin URL tanımlayıcısını döner.
	GetSlug() string
	// GetType, filtrenin tipini döner (select, date, range, vb.).
	GetType() string
	// GetOptions, filtrenin seçeneklerini döner.
	GetOptions() map[string]string
	// Apply, filtreyi sorguya uygular.
	Apply(db *gorm.DB, value any) *gorm.DB
}

// Sortable, liste görünümlerinde varsayılan sıralama ayarlarını tanımlar.
type Sortable struct {
	// Column, sıralanacak veritabanı sütun adı.
	Column string
	// Direction, sıralama yönü (asc veya desc).
	Direction string
}

// Resource, paneldeki her bir varlığı (örneğin Users, Posts) temsil eden arayüzdür.
// Bu arayüz, kaynağın veri modelini, alanlarını, görünüm ayarlarını ve yetkilendirme kurallarını tanımlar.
type Resource interface {
	// Model, GORM model yapısının bir örneğini döner.
	Model() any
	// Fields, kaynak form ve listelerinde gösterilecek alanları (fields) tanımlar.
	Fields() []fields.Element
	// GetFields, belirli bir bağlamda gösterilecek alanları döner.
	// Requirement 11.1: Resource arayüzünü, alanları almak için metodlar içerecek şekilde genişlet
	GetFields(ctx *appContext.Context) []fields.Element
	// With, ilişkili verilerin (eager loading) yüklenmesi için kullanılır.
	With() []string
	// Lenses, özel veri filtreleme görünümlerini tanımlar.
	Lenses() []Lens
	// GetLenses, kaynağın tüm lens'lerini döner.
	// Requirement 11.1: Resource arayüzünü, lens'leri almak için metodlar içerecek şekilde genişlet
	GetLenses() []Lens
	// Slug, kaynağın URL'deki tanımlayıcısıdır (örn: "users").
	Slug() string
	// Title, panelde görünecek başlık (örn: "Kullanıcılar").
	Title() string
	// Icon, menüde kullanılacak ikon adı (Lucide ikon seti).
	Icon() string
	// Group, kaynağın menüde hangi grup altında listeleneceğini belirler.
	Group() string
	// Policy, kaynak üzerindeki yetkilendirme (CRUD) kurallarını döner.
	Policy() auth.Policy
	// GetPolicy, kaynağın yetkilendirme politikasını döner.
	// Requirement 11.1: Resource arayüzünü, politikaları almak için metodlar içerecek şekilde genişlet
	GetPolicy() auth.Policy
	// GetSortable, varsayılan sıralama ayarlarını döner.
	GetSortable() []Sortable
	// GetDialogType, ekleme/düzenleme formunun hangi tipte açılacağını (Sheet, Drawer, Modal) döner.
	GetDialogType() DialogType
	// SetDialogType, form görünüm tipini ayarlar (Zincirleme metod kullanımına uygun).
	SetDialogType(DialogType) Resource
	// Repository, kaynağın veri erişim katmanını (DataProvider) döner.
	// Varsayılan GormDataProvider yerine özel bir repository kullanılmak istenirse bu metoddan dönülmelidir.
	Repository(db *gorm.DB) data.DataProvider
	// Widgets, resource üzerinde gösterilecek widget'ları döner.
	// Cards returns the cards/widgets for the resource dashboard
	Cards() []widget.Card
	// GetCards, belirli bir bağlamda gösterilecek card'ları döner.
	// Requirement 11.1: Resource arayüzünü, card'ları almak için metodlar içerecek şekilde genişlet
	GetCards(ctx *appContext.Context) []widget.Card
	// ResolveField, bir alanın değerini dinamik olarak hesaplayan ve dönüştüren fonksiyon.
	// Requirement 11.2: Resource'ların kendi alan çözümleme mantığını tanımlamasına izin ver
	ResolveField(fieldName string, item any) (any, error)
	// GetActions, kaynağın özel işlemlerini döner.
	// Requirement 11.4: Resource arayüzünü, işlemleri almak için metodlar içerecek şekilde genişlet
	GetActions() []Action
	// GetFilters, kaynağın filtreleme seçeneklerini döner.
	// Requirement 11.4: Resource arayüzünü, filtreleri almak için metodlar içerecek şekilde genişlet
	GetFilters() []Filter
	// StoreHandler, dosya yükleme işlemlerini yönetir.
	StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error)
	// NavigationOrder, menüdeki sıralama önceliğini döner. (Düşük sayı -> Üst sıra)
	NavigationOrder() int
	// Visible, kaynağın menüde görünüp görünmeyeceğini belirler.
	Visible() bool
}
