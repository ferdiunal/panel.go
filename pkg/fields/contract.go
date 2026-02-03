package fields

import (
	"context"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

// Element, form ve liste görünümlerinde kullanılan alanların (fields) ortak arayüzüdür.
type Element interface {
	// GetKey, alanın veritabanı veya JSON anahtarını döner.
	GetKey() string
	// GetView, alanın frontend tarafındaki bileşen adını döner.
	GetView() string
	// Extract, verilen kaynaktan (modelden) veriyi çıkarır ve alana atar.
	Extract(resource interface{})
	// JsonSerialize, alanın frontend'e gönderilecek JSON temsilini oluşturur.
	JsonSerialize() map[string]interface{}
	// GetContext, alanın hangi bağlamda (form, liste, detay) olduğunu döner.
	GetContext() ElementContext

	// Fluid Setters (Zincirleme Metodlar)

	// SetName, alanın görünen adını belirler.
	SetName(name string) Element
	// SetKey, alanın veri anahtarını belirler.
	SetKey(key string) Element
	// OnList, alanı liste görünümünde aktif eder.
	OnList() Element
	// OnDetail, alanı detay görünümünde aktif eder.
	OnDetail() Element
	// OnForm, alanı form (ekleme/düzenleme) görünümünde aktif eder.
	OnForm() Element
	// HideOnList, alanı liste görünümünde gizler.
	HideOnList() Element
	// HideOnDetail, alanı detay görünümünde gizler.
	HideOnDetail() Element
	// HideOnCreate, alanı ekleme formunda gizler.
	HideOnCreate() Element
	// HideOnUpdate, alanı güncelleme formunda gizler.
	HideOnUpdate() Element
	// OnlyOnList, alanı sadece liste görünümünde gösterir.
	OnlyOnList() Element
	// OnlyOnDetail, alanı sadece detay görünümünde gösterir.
	OnlyOnDetail() Element
	// OnlyOnCreate, alanı sadece ekleme formunda gösterir.
	OnlyOnCreate() Element
	// OnlyOnUpdate, alanı sadece güncelleme formunda gösterir.
	OnlyOnUpdate() Element
	// OnlyOnForm, alanı sadece formlarda gösterir.
	OnlyOnForm() Element
	// ReadOnly, alanı salt okunur yapar.
	ReadOnly() Element
	// WithProps, frontend bileşenine ekstra parametreler (props) geçer.
	WithProps(key string, value interface{}) Element
	// Disabled, alanı pasif (kullanılamaz) hale getirir.
	Disabled() Element
	// Immutable, alanın değerinin değiştirilmesini engeller.
	Immutable() Element
	// Required, alanın zorunlu olduğunu belirtir.
	Required() Element
	// Nullable, alanın boş (null) olabileceğini belirtir.
	Nullable() Element
	// Placeholder, input alanında görünecek yer tutucu metni belirler.
	Placeholder(placeholder string) Element
	// Label, input etiketi metnini belirler.
	Label(label string) Element
	// HelpText, alanın altında görünecek yardımcı metni belirler.
	HelpText(helpText string) Element
	// Filterable, alanın filtrelenebilir olduğunu belirtir.
	Filterable() Element
	// Sortable, alanın sıralanabilir olduğunu belirtir.
	Sortable() Element
	// Searchable, alanın aranabilir olduğunu belirtir.
	Searchable() Element
	// IsSearchable, alanın aranabilir olup olmadığını döner.
	IsSearchable() bool
	// Stacked, alanın yığın (stacked) görünümde olup olmadığını belirler (sadece UI için).
	Stacked() Element
	// SetTextAlign, metin hizalamasını ayarlar (left, center, right).
	SetTextAlign(align string) Element
	// IsVisible, alanın verilen context içinde görünür olup olmadığını kontrol eder.
	IsVisible(ctx context.Context) bool
	// CanSee, alanın görünürlüğünü dinamik olarak belirleyen fonksiyonu tanımlar.
	CanSee(fn VisibilityFunc) Element

	// Storage

	// StoreAs, dosya yükleme işlemleri için özel kayıt mantığını belirler.
	StoreAs(fn StorageCallbackFunc) Element
	// GetStorageCallback, tanımlı dosya kayıt fonksiyonunu döner.
	GetStorageCallback() StorageCallbackFunc

	// Transform

	// Resolve, verinin frontend'e gitmeden önce işlenmesini (formatlanmasını) sağlar.
	Resolve(fn func(value interface{}) interface{}) Element
	// GetResolveCallback, dönüştürme fonksiyonunu döner.
	GetResolveCallback() func(value interface{}) interface{}
	// Modify, verinin veritabanına kaydedilmeden önce işlenmesini sağlar.
	Modify(fn func(value interface{}) interface{}) Element
	// ModifyCallback, modifikasyon fonksiyonunu döner.
	GetModifyCallback() func(value interface{}) interface{}

	// Options, Combobox veya Select gibi alanlar için seçenekleri belirler.
	Options(options interface{}) Element
	// Default, alanın varsayılan değerini belirler.
	Default(value interface{}) Element
}

// VisibilityFunc, görünürlük kontrol fonksiyonu tipidir.
type VisibilityFunc func(ctx context.Context) bool

// StorageCallbackFunc, dosya kayıt fonksiyonu tipidir.
type StorageCallbackFunc func(c *fiber.Ctx, file *multipart.FileHeader) (string, error)
