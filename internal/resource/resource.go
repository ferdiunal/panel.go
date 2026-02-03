package resource

import (
	"mime/multipart"

	"github.com/ferdiunal/panel.go/internal/auth"
	appContext "github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/data"
	"github.com/ferdiunal/panel.go/internal/fields"
	"github.com/ferdiunal/panel.go/internal/widget"
	"gorm.io/gorm"
)

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
	Model() interface{}
	// Fields, kaynak form ve listelerinde gösterilecek alanları (fields) tanımlar.
	Fields() []fields.Element
	// With, ilişkili verilerin (eager loading) yüklenmesi için kullanılır.
	With() []string
	// Lenses, özel veri filtreleme görünümlerini tanımlar.
	Lenses() []Lens
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
	// StoreHandler, dosya yükleme işlemlerini yönetir.
	StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error)
	// NavigationOrder, menüdeki sıralama önceliğini döner. (Düşük sayı -> Üst sıra)
	NavigationOrder() int
	// Visible, kaynağın menüde görünüp görünmeyeceğini belirler.
	Visible() bool
}
