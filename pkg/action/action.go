package action

import (
	"fmt"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"gorm.io/gorm"
)

// Bu interface, panel uygulamasında kaynaklar (resources) üzerinde gerçekleştirilebilecek
// toplu işlemleri (bulk operations) tanımlar. Kullanıcılar seçili kaynaklar üzerinde
// silme, güncelleme, dışa aktarma gibi işlemleri gerçekleştirebilirler.
//
// Kullanım Senaryoları:
// - Seçili ürünleri silme
// - Seçili kullanıcıları devre dışı bırakma
// - Seçili siparişleri dışa aktarma
// - Seçili öğeleri toplu güncelleme
//
// Örnek Kullanım:
//   action := action.New("Delete Products").
//       Destructive().
//       Handle(func(ctx *action.ActionContext) error {
//           // Seçili ürünleri sil
//           return ctx.DB.Delete(ctx.Models).Error
//       })
type Action interface {
	// Bu metod, kullanıcı arayüzünde gösterilecek aksiyon adını döndürür.
	// Örneğin: "Delete", "Export", "Archive"
	// Döndürür: Aksiyon adı (string)
	GetName() string

	// Bu metod, aksiyonun URL-güvenli tanımlayıcısını döndürür.
	// Slug, API çağrılarında ve URL'lerde kullanılır.
	// Örneğin: "delete-products", "export-users"
	// Döndürür: URL-güvenli slug (string)
	GetSlug() string

	// Bu metod, aksiyonun simgesini döndürür.
	// Simge adı, ön uçta görüntülenmek üzere kullanılır.
	// Örneğin: "trash", "download", "archive"
	// Döndürür: Simge adı (string)
	GetIcon() string

	// Bu metod, aksiyonu gerçekleştirmeden önce gösterilecek
	// onay mesajını döndürür.
	// Örneğin: "Bu ürünleri silmek istediğinizden emin misiniz?"
	// Döndürür: Onay mesajı (string)
	GetConfirmText() string

	// Bu metod, onay iletişim kutusundaki "Onayla" düğmesinin
	// metni döndürür.
	// Örneğin: "Sil", "Onayla", "Devam Et"
	// Döndürür: Düğme metni (string)
	GetConfirmButtonText() string

	// Bu metod, onay iletişim kutusundaki "İptal" düğmesinin
	// metni döndürür.
	// Örneğin: "İptal", "Vazgeç"
	// Döndürür: Düğme metni (string)
	GetCancelButtonText() string

	// Bu metod, aksiyonun yıkıcı (destructive) olup olmadığını belirtir.
	// Yıkıcı aksiyonlar (silme, kalıcı değişiklik) kırmızı renkle gösterilir.
	// Döndürür: true ise yıkıcı, false ise güvenli
	IsDestructive() bool

	// Bu metod, aksiyonun sadece liste (index) görünümünde
	// kullanılabilir olup olmadığını belirtir.
	// Döndürür: true ise sadece liste görünümünde göster
	OnlyOnIndex() bool

	// Bu metod, aksiyonun sadece detay görünümünde
	// kullanılabilir olup olmadığını belirtir.
	// Döndürür: true ise sadece detay görünümünde göster
	OnlyOnDetail() bool

	// Bu metod, aksiyonun satır içi (inline) olarak gösterilip gösterilmeyeceğini belirtir.
	// Satır içi aksiyonlar, her satırda ayrı ayrı gösterilir.
	// Döndürür: true ise satır içi göster
	ShowInline() bool

	// Bu metod, aksiyonu gerçekleştirmek için gerekli olan
	// form alanlarını döndürür.
	// Örneğin: kategori seçimi, neden metni, hedef klasör
	// Döndürür: core.Element alanlarının dilimi
	GetFields() []core.Element

	// Bu metod, aksiyonu gerçekleştirmek için gerekli olan
	// işlemi yürütür. Seçili kaynaklar ve form verileri ile çalışır.
	//
	// Parametreler:
	//   - ctx: Panel bağlamı (kullanıcı, veritabanı, istek bilgileri)
	//   - items: Seçili kaynaklar (modeller)
	//
	// Döndürür: Hata varsa error, başarılı ise nil
	//
	// Önemli Notlar:
	// - ctx.Locals("action_fields") ile form verilerine erişebilirsiniz
	// - ctx.Locals("db") ile veritabanı bağlantısını alabilirsiniz
	// - ctx.Locals("user") ile mevcut kullanıcıyı alabilirsiniz
	Execute(ctx *context.Context, items []any) error

	// Bu metod, aksiyonun belirli bir bağlamda çalıştırılabilir olup olmadığını
	// kontrol eder. Yetkilendirme ve izin kontrolü için kullanılır.
	//
	// Parametreler:
	//   - ctx: Aksiyon bağlamı (modeller, alanlar, kullanıcı, veritabanı)
	//
	// Döndürür: true ise aksiyon çalıştırılabilir, false ise gizle
	//
	// Örnek Kullanım:
	//   CanRun: func(ctx *ActionContext) bool {
	//       // Sadece yöneticiler bu aksiyonu çalıştırabilir
	//       return ctx.User.IsAdmin
	//   }
	CanRun(ctx *ActionContext) bool
}

// Bu yapı, Action interface'inin temel bir uygulamasını sağlar.
// Özel aksiyonlar oluşturmak için bu yapıyı gömüp (embed) varsayılan davranışı
// miras alabilirsiniz. Fluent API deseni kullanarak aksiyonları kolayca yapılandırabilirsiniz.
//
// Kullanım Senaryoları:
// - Silme aksiyonları
// - Toplu güncelleme aksiyonları
// - Dışa aktarma aksiyonları
// - Özel iş mantığı aksiyonları
//
// Örnek Kullanım:
//   deleteAction := action.New("Delete Selected").
//       SetIcon("trash").
//       Destructive().
//       Confirm("Bu öğeleri silmek istediğinizden emin misiniz?").
//       Handle(func(ctx *action.ActionContext) error {
//           return ctx.DB.Delete(ctx.Models).Error
//       }).
//       AuthorizeUsing(func(ctx *action.ActionContext) bool {
//           return ctx.User.IsAdmin
//       })
//
// Alan Açıklamaları:
// - Name: Kullanıcı arayüzünde gösterilecek aksiyon adı
// - Slug: URL-güvenli tanımlayıcı (otomatik olarak oluşturulur)
// - Icon: Simge adı (örneğin: "trash", "download", "archive")
// - ConfirmText: Onay iletişim kutusunda gösterilecek mesaj
// - ConfirmButtonText: Onay düğmesinin metni (varsayılan: "Confirm")
// - CancelButtonText: İptal düğmesinin metni (varsayılan: "Cancel")
// - DestructiveType: Aksiyon yıkıcı ise true (kırmızı renkle gösterilir)
// - OnlyOnIndexFlag: Sadece liste görünümünde gösterilecekse true
// - OnlyOnDetailFlag: Sadece detay görünümünde gösterilecekse true
// - ShowInlineFlag: Satır içi gösterilecekse true
// - Fields: Aksiyon için gerekli form alanları
// - HandleFunc: Aksiyonu gerçekleştiren işlev
// - CanRunFunc: Aksiyonun çalıştırılabilir olup olmadığını kontrol eden işlev
type BaseAction struct {
	// Kullanıcı arayüzünde gösterilecek aksiyon adı
	Name string

	// URL-güvenli tanımlayıcı (otomatik olarak oluşturulur)
	Slug string

	// Simge adı (ön uçta görüntülenmek üzere)
	Icon string

	// Onay iletişim kutusunda gösterilecek mesaj
	ConfirmText string

	// Onay düğmesinin metni
	ConfirmButtonText string

	// İptal düğmesinin metni
	CancelButtonText string

	// Aksiyon yıkıcı ise true (kırmızı renkle gösterilir)
	DestructiveType bool

	// Sadece liste görünümünde gösterilecekse true
	OnlyOnIndexFlag bool

	// Sadece detay görünümünde gösterilecekse true
	OnlyOnDetailFlag bool

	// Satır içi gösterilecekse true
	ShowInlineFlag bool

	// Aksiyon için gerekli form alanları
	Fields []core.Element

	// Aksiyonu gerçekleştiren işlev
	HandleFunc func(ctx *ActionContext) error

	// Aksiyonun çalıştırılabilir olup olmadığını kontrol eden işlev
	CanRunFunc func(ctx *ActionContext) bool
}

// Bu fonksiyon, verilen ad ile yeni bir BaseAction oluşturur.
// Slug otomatik olarak addan oluşturulur (boşluklar tire ile değiştirilir, küçük harfe çevrilir).
//
// Parametreler:
//   - name: Aksiyon adı (örneğin: "Delete Products", "Export Users")
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı
//
// Varsayılan Değerler:
//   - ConfirmButtonText: "Confirm"
//   - CancelButtonText: "Cancel"
//   - Diğer alanlar: boş/false
//
// Örnek Kullanım:
//   action := action.New("Delete Products")
//   // Slug otomatik olarak "delete-products" olur
//
// Önemli Notlar:
// - Slug otomatik oluşturulur ancak SetSlug() ile değiştirilebilir
// - Fluent API ile yapılandırma yapabilirsiniz
func New(name string) *BaseAction {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return &BaseAction{
		Name:              name,
		Slug:              slug,
		ConfirmButtonText: "Confirm",
		CancelButtonText:  "Cancel",
	}
}

// ============================================================================
// Fluent API Metodları - Aksiyon Yapılandırması
// ============================================================================
// Bu metodlar, aksiyon özelliklerini yapılandırmak için fluent API deseni
// kullanır. Her metod, method chaining'i desteklemek için BaseAction pointer'ı
// döndürür.
//
// Örnek Kullanım:
//   action := action.New("Delete").
//       SetIcon("trash").
//       Destructive().
//       Confirm("Silmek istediğinizden emin misiniz?").
//       Handle(deleteHandler).
//       AuthorizeUsing(isAdmin)

// Bu metod, aksiyonun görüntü adını ayarlar.
//
// Parametreler:
//   - name: Yeni aksiyon adı
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.SetName("Delete Products")
func (a *BaseAction) SetName(name string) *BaseAction {
	a.Name = name
	return a
}

// Bu metod, aksiyonun URL-güvenli tanımlayıcısını ayarlar.
// Slug, API çağrılarında ve URL'lerde kullanılır.
//
// Parametreler:
//   - slug: URL-güvenli tanımlayıcı (örneğin: "delete-products")
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.SetSlug("custom-delete-action")
//
// Önemli Notlar:
// - Slug otomatik olarak New() fonksiyonunda oluşturulur
// - Bu metod ile varsayılan slug'ı geçersiz kılabilirsiniz
func (a *BaseAction) SetSlug(slug string) *BaseAction {
	a.Slug = slug
	return a
}

// Bu metod, aksiyonun simgesini ayarlar.
// Simge adı, ön uçta görüntülenmek üzere kullanılır.
//
// Parametreler:
//   - icon: Simge adı (örneğin: "trash", "download", "archive")
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.SetIcon("trash")
//   action.SetIcon("download")
//   action.SetIcon("archive")
func (a *BaseAction) SetIcon(icon string) *BaseAction {
	a.Icon = icon
	return a
}

// Bu metod, aksiyonu gerçekleştirmeden önce gösterilecek
// onay mesajını ayarlar.
//
// Parametreler:
//   - text: Onay mesajı
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.Confirm("Bu ürünleri silmek istediğinizden emin misiniz?")
//   action.Confirm("Bu işlem geri alınamaz. Devam etmek istiyor musunuz?")
//
// Önemli Notlar:
// - Onay mesajı boş bırakılırsa, onay iletişim kutusu gösterilmez
// - Yıkıcı aksiyonlar için onay mesajı ayarlamanız önerilir
func (a *BaseAction) Confirm(text string) *BaseAction {
	a.ConfirmText = text
	return a
}

// Bu metod, onay iletişim kutusundaki "Onayla" düğmesinin
// metni ayarlar.
//
// Parametreler:
//   - text: Düğme metni
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Varsayılan Değer: "Confirm"
//
// Örnek Kullanım:
//   action.ConfirmButton("Sil")
//   action.ConfirmButton("Onayla")
//   action.ConfirmButton("Devam Et")
func (a *BaseAction) ConfirmButton(text string) *BaseAction {
	a.ConfirmButtonText = text
	return a
}

// Bu metod, onay iletişim kutusundaki "İptal" düğmesinin
// metni ayarlar.
//
// Parametreler:
//   - text: Düğme metni
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Varsayılan Değer: "Cancel"
//
// Örnek Kullanım:
//   action.CancelButton("İptal")
//   action.CancelButton("Vazgeç")
func (a *BaseAction) CancelButton(text string) *BaseAction {
	a.CancelButtonText = text
	return a
}

// Bu metod, aksiyonu yıkıcı (destructive) olarak işaretler.
// Yıkıcı aksiyonlar, kullanıcı arayüzünde kırmızı renkle gösterilir
// ve ek onay gerektirir.
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.Destructive()
//
// Önemli Notlar:
// - Silme, kalıcı değişiklik gibi işlemler için kullanılır
// - Yıkıcı aksiyonlar için Confirm() ile onay mesajı ayarlamanız önerilir
// - Ön uçta kırmızı renkle gösterilir
func (a *BaseAction) Destructive() *BaseAction {
	a.DestructiveType = true
	return a
}

// Bu metod, aksiyonu sadece liste (index) görünümünde
// kullanılabilir olarak işaretler.
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.ShowOnlyOnIndex()
//
// Önemli Notlar:
// - ShowOnlyOnDetail() ile birlikte kullanılamaz
// - Detay görünümünde bu aksiyon gizlenecektir
func (a *BaseAction) ShowOnlyOnIndex() *BaseAction {
	a.OnlyOnIndexFlag = true
	return a
}

// Bu metod, aksiyonu sadece detay görünümünde
// kullanılabilir olarak işaretler.
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.ShowOnlyOnDetail()
//
// Önemli Notlar:
// - ShowOnlyOnIndex() ile birlikte kullanılamaz
// - Liste görünümünde bu aksiyon gizlenecektir
func (a *BaseAction) ShowOnlyOnDetail() *BaseAction {
	a.OnlyOnDetailFlag = true
	return a
}

// Bu metod, aksiyonu satır içi (inline) olarak gösterilecek şekilde işaretler.
// Satır içi aksiyonlar, her satırda ayrı ayrı gösterilir.
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.ShowInlineAction()
//
// Önemli Notlar:
// - Satır içi aksiyonlar, her satırda ayrı düğme olarak gösterilir
// - Toplu işlem aksiyonları ile birlikte kullanılabilir
func (a *BaseAction) ShowInlineAction() *BaseAction {
	a.ShowInlineFlag = true
	return a
}

// Bu metod, aksiyonu gerçekleştirmek için gerekli olan
// form alanlarını ayarlar.
//
// Parametreler:
//   - fields: core.Element alanlarının değişken sayıda argümanı
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.WithFields(
//       &fields.Text{Name: "reason", Label: "Neden"},
//       &fields.Select{Name: "category", Label: "Kategori"},
//   )
//
// Önemli Notlar:
// - Alanlar, aksiyon gerçekleştirilmeden önce kullanıcıdan toplanır
// - Alanlar, ActionContext.Fields haritasında kullanılabilir
func (a *BaseAction) WithFields(fields ...core.Element) *BaseAction {
	a.Fields = fields
	return a
}

// Bu metod, aksiyonu gerçekleştiren işlevi ayarlar.
// Bu işlev, aksiyon tetiklendiğinde çağrılır.
//
// Parametreler:
//   - fn: Aksiyon işlevi (ActionContext alır, error döndürür)
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.Handle(func(ctx *ActionContext) error {
//       // Seçili modelleri işle
//       for _, model := range ctx.Models {
//           // İşlem yap
//       }
//       return nil
//   })
//
// Önemli Notlar:
// - Bu metod zorunludur, aksi takdirde Execute() hata döndürür
// - ActionContext, modeller, alanlar, kullanıcı ve veritabanı bilgisi içerir
// - Hata döndürürseniz, işlem başarısız olur
func (a *BaseAction) Handle(fn func(ctx *ActionContext) error) *BaseAction {
	a.HandleFunc = fn
	return a
}

// Bu metod, aksiyonun belirli bir bağlamda çalıştırılabilir olup olmadığını
// kontrol eden işlevi ayarlar. Yetkilendirme ve izin kontrolü için kullanılır.
//
// Parametreler:
//   - fn: Yetkilendirme işlevi (ActionContext alır, bool döndürür)
//
// Döndürür: Yapılandırılmış BaseAction pointer'ı (method chaining için)
//
// Örnek Kullanım:
//   action.AuthorizeUsing(func(ctx *ActionContext) bool {
//       // Sadece yöneticiler bu aksiyonu çalıştırabilir
//       user := ctx.User.(*User)
//       return user.IsAdmin
//   })
//
// Önemli Notlar:
// - Bu metod ayarlanmazsa, aksiyon her zaman çalıştırılabilir
// - Yetkilendirme başarısız olursa, aksiyon gizlenir
// - Kullanıcı rolü, izin ve diğer kontroller burada yapılır
func (a *BaseAction) AuthorizeUsing(fn func(ctx *ActionContext) bool) *BaseAction {
	a.CanRunFunc = fn
	return a
}

// ============================================================================
// Interface Implementation Metodları - Action Interface'i Gerçekleştirme
// ============================================================================
// Bu metodlar, Action interface'inin tüm metodlarını gerçekleştirir.
// BaseAction, Action interface'ini tam olarak uygular.

// Bu metod, aksiyonun görüntü adını döndürür.
// Bu metod, Action interface'inin GetName() metodunu gerçekleştirir.
//
// Döndürür: Aksiyonun adı (string)
//
// Örnek Kullanım:
//   action := action.New("Delete Products")
//   name := action.GetName() // "Delete Products"
func (a *BaseAction) GetName() string {
	return a.Name
}

// Bu metod, aksiyonun URL-güvenli tanımlayıcısını döndürür.
// Bu metod, Action interface'inin GetSlug() metodunu gerçekleştirir.
//
// Döndürür: URL-güvenli slug (string)
//
// Örnek Kullanım:
//   action := action.New("Delete Products")
//   slug := action.GetSlug() // "delete-products"
func (a *BaseAction) GetSlug() string {
	return a.Slug
}

// Bu metod, aksiyonun simgesini döndürür.
// Bu metod, Action interface'inin GetIcon() metodunu gerçekleştirir.
//
// Döndürür: Simge adı (string)
//
// Örnek Kullanım:
//   action.SetIcon("trash")
//   icon := action.GetIcon() // "trash"
func (a *BaseAction) GetIcon() string {
	return a.Icon
}

// Bu metod, aksiyonun onay mesajını döndürür.
// Bu metod, Action interface'inin GetConfirmText() metodunu gerçekleştirir.
//
// Döndürür: Onay mesajı (string)
//
// Örnek Kullanım:
//   action.Confirm("Bu işlemi gerçekleştirmek istediğinizden emin misiniz?")
//   text := action.GetConfirmText()
//
// Önemli Notlar:
// - Onay mesajı boş ise, onay iletişim kutusu gösterilmez
func (a *BaseAction) GetConfirmText() string {
	return a.ConfirmText
}

// Bu metod, onay düğmesinin metnini döndürür.
// Bu metod, Action interface'inin GetConfirmButtonText() metodunu gerçekleştirir.
//
// Döndürür: Onay düğmesinin metni (string)
//
// Varsayılan Değer: "Confirm"
//
// Örnek Kullanım:
//   action.ConfirmButton("Sil")
//   text := action.GetConfirmButtonText() // "Sil"
func (a *BaseAction) GetConfirmButtonText() string {
	return a.ConfirmButtonText
}

// Bu metod, iptal düğmesinin metnini döndürür.
// Bu metod, Action interface'inin GetCancelButtonText() metodunu gerçekleştirir.
//
// Döndürür: İptal düğmesinin metni (string)
//
// Varsayılan Değer: "Cancel"
//
// Örnek Kullanım:
//   action.CancelButton("İptal")
//   text := action.GetCancelButtonText() // "İptal"
func (a *BaseAction) GetCancelButtonText() string {
	return a.CancelButtonText
}

// Bu metod, aksiyonun yıkıcı olup olmadığını döndürür.
// Bu metod, Action interface'inin IsDestructive() metodunu gerçekleştirir.
//
// Döndürür: true ise yıkıcı, false ise güvenli
//
// Örnek Kullanım:
//   action.Destructive()
//   isDestructive := action.IsDestructive() // true
//
// Önemli Notlar:
// - Yıkıcı aksiyonlar, ön uçta kırmızı renkle gösterilir
// - Silme, kalıcı değişiklik gibi işlemler için kullanılır
func (a *BaseAction) IsDestructive() bool {
	return a.DestructiveType
}

// Bu metod, aksiyonun sadece liste görünümünde kullanılabilir olup olmadığını döndürür.
// Bu metod, Action interface'inin OnlyOnIndex() metodunu gerçekleştirir.
//
// Döndürür: true ise sadece liste görünümünde göster
//
// Örnek Kullanım:
//   action.ShowOnlyOnIndex()
//   onlyIndex := action.OnlyOnIndex() // true
//
// Önemli Notlar:
// - true ise, aksiyon detay görünümünde gizlenir
func (a *BaseAction) OnlyOnIndex() bool {
	return a.OnlyOnIndexFlag
}

// Bu metod, aksiyonun sadece detay görünümünde kullanılabilir olup olmadığını döndürür.
// Bu metod, Action interface'inin OnlyOnDetail() metodunu gerçekleştirir.
//
// Döndürür: true ise sadece detay görünümünde göster
//
// Örnek Kullanım:
//   action.ShowOnlyOnDetail()
//   onlyDetail := action.OnlyOnDetail() // true
//
// Önemli Notlar:
// - true ise, aksiyon liste görünümünde gizlenir
func (a *BaseAction) OnlyOnDetail() bool {
	return a.OnlyOnDetailFlag
}

// Bu metod, aksiyonun satır içi (inline) olarak gösterilip gösterilmeyeceğini döndürür.
// Bu metod, Action interface'inin ShowInline() metodunu gerçekleştirir.
//
// Döndürür: true ise satır içi göster
//
// Örnek Kullanım:
//   action.ShowInlineAction()
//   inline := action.ShowInline() // true
//
// Önemli Notlar:
// - true ise, aksiyon her satırda ayrı düğme olarak gösterilir
func (a *BaseAction) ShowInline() bool {
	return a.ShowInlineFlag
}

// Bu metod, aksiyonu gerçekleştirmek için gerekli olan form alanlarını döndürür.
// Bu metod, Action interface'inin GetFields() metodunu gerçekleştirir.
//
// Döndürür: core.Element alanlarının dilimi
//
// Örnek Kullanım:
//   fields := action.GetFields()
//   for _, field := range fields {
//       // Her alanı işle
//   }
//
// Önemli Notlar:
// - Alanlar, aksiyon gerçekleştirilmeden önce kullanıcıdan toplanır
// - Alanlar, ActionContext.Fields haritasında kullanılabilir
func (a *BaseAction) GetFields() []core.Element {
	return a.Fields
}

// Bu metod, aksiyonu gerçekleştirmek için gerekli olan işlemi yürütür.
// Bu metod, Action interface'inin Execute() metodunu gerçekleştirir.
//
// Parametreler:
//   - ctx: Panel bağlamı (kullanıcı, veritabanı, istek bilgileri)
//   - items: Seçili kaynaklar (modeller)
//
// Döndürür: Hata varsa error, başarılı ise nil
//
// İşlem Akışı:
//   1. HandleFunc'in tanımlanıp tanımlanmadığını kontrol et
//   2. ctx.Locals("action_fields") ile form verilerini al
//   3. ctx.Locals("db") ile veritabanı bağlantısını al
//   4. ActionContext oluştur
//   5. HandleFunc'i çağır
//
// Örnek Kullanım:
//   err := action.Execute(ctx, selectedItems)
//   if err != nil {
//       // Hata işle
//   }
//
// Önemli Notlar:
// - HandleFunc ayarlanmamışsa, hata döndürür
// - ctx.Locals("action_fields") ile form verilerine erişebilirsiniz
// - ctx.Locals("db") ile veritabanı bağlantısını alabilirsiniz
// - ctx.Locals("user") ile mevcut kullanıcıyı alabilirsiniz
// - ctx.Params("resource") ile kaynak adını alabilirsiniz
func (a *BaseAction) Execute(ctx *context.Context, items []any) error {
	if a.HandleFunc == nil {
		return fmt.Errorf("action handler not defined")
	}

	// ctx.Locals("action_fields") ile form verilerini al
	// Bu veriler, aksiyon gerçekleştirilmeden önce kullanıcı tarafından doldurulur
	fields := make(map[string]interface{})
	if actionFields := ctx.Locals("action_fields"); actionFields != nil {
		if f, ok := actionFields.(map[string]interface{}); ok {
			fields = f
		}
	}

	// ctx.Locals("db") ile veritabanı bağlantısını al
	// Bu bağlantı, aksiyon işlevinde modelleri sorgulamak ve güncellemek için kullanılır
	var db *gorm.DB
	if dbVal := ctx.Locals("db"); dbVal != nil {
		if d, ok := dbVal.(*gorm.DB); ok {
			db = d
		}
	}

	// ActionContext oluştur ve HandleFunc'i çağır
	// ActionContext, aksiyon işlevine gerekli tüm bilgileri sağlar
	actionCtx := &ActionContext{
		Models:   items,
		Fields:   fields,
		User:     ctx.Locals("user"),
		Resource: ctx.Params("resource"),
		DB:       db,
		Ctx:      ctx.Ctx,
	}

	return a.HandleFunc(actionCtx)
}

// Bu metod, aksiyonun belirli bir bağlamda çalıştırılabilir olup olmadığını kontrol eder.
// Bu metod, Action interface'inin CanRun() metodunu gerçekleştirir.
//
// Parametreler:
//   - ctx: Aksiyon bağlamı (modeller, alanlar, kullanıcı, veritabanı)
//
// Döndürür: true ise aksiyon çalıştırılabilir, false ise gizle
//
// İşlem Akışı:
//   1. CanRunFunc'in tanımlanıp tanımlanmadığını kontrol et
//   2. Tanımlanmamışsa, true döndür (varsayılan olarak çalıştırılabilir)
//   3. Tanımlanmışsa, CanRunFunc'i çağır ve sonucunu döndür
//
// Örnek Kullanım:
//   if action.CanRun(actionCtx) {
//       // Aksiyonu çalıştır
//   }
//
// Önemli Notlar:
// - CanRunFunc ayarlanmamışsa, aksiyon her zaman çalıştırılabilir
// - Yetkilendirme başarısız olursa, aksiyon gizlenir
// - Kullanıcı rolü, izin ve diğer kontroller burada yapılır
func (a *BaseAction) CanRun(ctx *ActionContext) bool {
	if a.CanRunFunc == nil {
		return true
	}
	return a.CanRunFunc(ctx)
}
