package resource

import (
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// / Bu interface, veri tabanı sorgularını özelleştirerek belirli görünümler (segmentler) oluşturmak için kullanılır.
// /
// / # Genel Bakış
// /
// / Lens'ler, kaynak verilerinin farklı perspektiflerden görüntülenmesini sağlar. Örneğin, bir kullanıcı kaynağında
// / "Aktif Kullanıcılar", "Pasif Kullanıcılar", "Yöneticiler" gibi farklı lens'ler tanımlayabilirsiniz.
// / Her lens, kendi sorgu mantığı, alanları ve widget'larına sahip olabilir.
// /
// / # Kullanım Senaryoları
// /
// / - **Veri Segmentasyonu**: Verileri belirli kriterlere göre filtreleyerek farklı görünümler oluşturma
// / - **Özel Raporlama**: Belirli kullanıcı grupları için özelleştirilmiş veri görünümleri
// / - **İş Akışı Yönetimi**: Farklı iş akışı aşamalarındaki kayıtları görüntüleme
// / - **Performans Optimizasyonu**: Sık kullanılan filtreleri önceden tanımlayarak sorgu performansını artırma
// /
// / # Avantajlar
// /
// / - **Esneklik**: Her lens için farklı alan setleri ve widget'lar tanımlayabilirsiniz
// / - **Yeniden Kullanılabilirlik**: Karmaşık sorguları tekrar tekrar yazmak yerine lens olarak tanımlayın
// / - **Kullanıcı Deneyimi**: Kullanıcıların ihtiyaç duydukları verilere hızlıca erişmesini sağlar
// / - **Bakım Kolaylığı**: Sorgu mantığı merkezi bir yerde tanımlanır
// /
// / # Önemli Notlar
// /
// / - Lens'ler, temel kaynak sorgusunu modifiye eder, tamamen yeni bir sorgu oluşturmaz
// / - Her lens'in benzersiz bir slug'ı olmalıdır (URL'de kullanılır)
// / - Lens-spesifik alanlar tanımlanmazsa, kaynağın varsayılan alanları kullanılır
// / - GetFields ve Fields metodları birlikte çalışır: Fields() statik alanları, GetFields() dinamik alanları döner
// /
// / # Kullanım Örneği
// /
// / ```go
// / type ActiveUsersLens struct{}
// /
// / func (l *ActiveUsersLens) Name() string {
// /     return "Aktif Kullanıcılar"
// / }
// /
// / func (l *ActiveUsersLens) Slug() string {
// /     return "active-users"
// / }
// /
// / func (l *ActiveUsersLens) Query(db *gorm.DB) *gorm.DB {
// /     return db.Where("status = ?", "active").Where("last_login > ?", time.Now().AddDate(0, -1, 0))
// / }
// /
// / func (l *ActiveUsersLens) Fields() []fields.Element {
// /     return []fields.Element{
// /         fields.NewText("name").SetLabel("İsim"),
// /         fields.NewText("email").SetLabel("E-posta"),
// /         fields.NewDateTime("last_login").SetLabel("Son Giriş"),
// /     }
// / }
// /
// / func (l *ActiveUsersLens) GetFields(ctx *appContext.Context) []fields.Element {
// /     // Dinamik alan yapılandırması
// /     return l.Fields()
// / }
// /
// / func (l *ActiveUsersLens) GetCards(ctx *appContext.Context) []widget.Card {
// /     return []widget.Card{
// /         widget.NewMetricCard("Toplam Aktif", "1,234"),
// /     }
// / }
// /
// / // Resource'a lens ekleme
// / func (r *UserResource) Lenses() []resource.Lens {
// /     return []resource.Lens{
// /         &ActiveUsersLens{},
// /     }
// / }
// / ```
// /
// / # İlgili Dokümantasyon
// /
// / - Alan tanımlamaları için: docs/Fields.md
// / - İlişki yönetimi için: docs/Relationships.md
// /
// / # Gereksinimler
// /
// / - Requirement 13.1: Lens'lerin özel sorgu mantığını tanımlamasına izin ver
// / - Requirement 13.2: Lens'lerin lens-spesifik alanları filtrelemesine izin ver
// / - Requirement 13.3: Lens'lerin lens-spesifik işlemler ve card'lar tanımlamasına izin ver
type Lens interface {
	/// Bu fonksiyon, Lens'in kullanıcı arayüzünde görüntülenecek adını döner.
	///
	/// # Parametreler
	///
	/// Parametre almaz.
	///
	/// # Döndürür
	///
	/// - `string`: Lens'in görünen adı (örn: "Aktif Kullanıcılar", "Bekleyen Siparişler")
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (l *ActiveUsersLens) Name() string {
	///     return "Aktif Kullanıcılar"
	/// }
	/// ```
	///
	/// # Önemli Notlar
	///
	/// - Ad, kullanıcı dostu ve açıklayıcı olmalıdır
	/// - Türkçe karakter kullanımı desteklenir
	/// - Benzersiz olması önerilir ancak zorunlu değildir
	Name() string

	/// Bu fonksiyon, Lens'in URL'de kullanılacak benzersiz tanımlayıcısını döner.
	///
	/// # Parametreler
	///
	/// Parametre almaz.
	///
	/// # Döndürür
	///
	/// - `string`: URL-safe slug (örn: "active-users", "pending-orders")
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (l *ActiveUsersLens) Slug() string {
	///     return "active-users"
	/// }
	/// ```
	///
	/// # Önemli Notlar
	///
	/// - Slug, URL-safe olmalıdır (küçük harf, tire ile ayrılmış)
	/// - Kaynak içinde benzersiz olmalıdır
	/// - Türkçe karakter kullanmayın
	/// - Boşluk yerine tire (-) kullanın
	///
	/// # Uyarılar
	///
	/// ⚠️ Aynı slug'a sahip iki lens tanımlanırsa, beklenmeyen davranışlar oluşabilir
	Slug() string

	/// Bu fonksiyon, temel GORM sorgusunu modifiye ederek lens-spesifik filtreleme uygular.
	///
	/// # Parametreler
	///
	/// - `db`: GORM veritabanı bağlantısı pointer'ı
	///
	/// # Döndürür
	///
	/// - `*gorm.DB`: Modifiye edilmiş GORM sorgu builder'ı
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (l *ActiveUsersLens) Query(db *gorm.DB) *gorm.DB {
	///     // Basit filtreleme
	///     return db.Where("status = ?", "active")
	/// }
	///
	/// func (l *RecentOrdersLens) Query(db *gorm.DB) *gorm.DB {
	///     // Karmaşık filtreleme
	///     return db.Where("created_at > ?", time.Now().AddDate(0, 0, -7)).
	///         Where("status IN ?", []string{"pending", "processing"}).
	///         Order("created_at DESC")
	/// }
	///
	/// func (l *HighValueCustomersLens) Query(db *gorm.DB) *gorm.DB {
	///     // Join ile filtreleme
	///     return db.Joins("LEFT JOIN orders ON orders.customer_id = customers.id").
	///         Group("customers.id").
	///         Having("SUM(orders.total) > ?", 10000)
	/// }
	/// ```
	///
	/// # Önemli Notlar
	///
	/// - Bu fonksiyon, mevcut sorguyu modifiye eder, yeni bir sorgu oluşturmaz
	/// - GORM'un tüm sorgu metodları kullanılabilir (Where, Join, Order, Group, vb.)
	/// - Performans için index'li kolonlar üzerinde filtreleme yapın
	/// - SQL injection'a karşı her zaman parametreli sorgular kullanın
	///
	/// # Uyarılar
	///
	/// ⚠️ Raw SQL kullanırken SQL injection riskine dikkat edin
	/// ⚠️ Karmaşık join'ler performansı etkileyebilir
	Query(db *gorm.DB) *gorm.DB

	/// Bu fonksiyon, lens görünümünde kullanılacak statik alan listesini döner.
	///
	/// # Parametreler
	///
	/// Parametre almaz.
	///
	/// # Döndürür
	///
	/// - `[]fields.Element`: Alan elemanları dizisi (boş dizi dönerse kaynak alanları kullanılır)
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (l *ActiveUsersLens) Fields() []fields.Element {
	///     return []fields.Element{
	///         fields.NewText("name").SetLabel("İsim"),
	///         fields.NewText("email").SetLabel("E-posta"),
	///         fields.NewDateTime("last_login").SetLabel("Son Giriş"),
	///         fields.NewBadge("status").SetLabel("Durum"),
	///     }
	/// }
	///
	/// // Boş dizi dönerek kaynak alanlarını kullanma
	/// func (l *AllUsersLens) Fields() []fields.Element {
	///     return []fields.Element{}
	/// }
	/// ```
	///
	/// # Önemli Notlar
	///
	/// - Boş dizi dönülürse, kaynağın varsayılan alanları kullanılır
	/// - GetFields() ile birlikte çalışır: Fields() statik, GetFields() dinamik alanlar için
	/// - Alan sıralaması, tablodaki görünüm sırasını belirler
	///
	/// # İlgili Dokümantasyon
	///
	/// - Alan tipleri ve kullanımı için: docs/Fields.md
	Fields() []fields.Element

	/// Bu fonksiyon, belirli bir bağlamda gösterilecek dinamik lens-spesifik alanları döner.
	///
	/// # Parametreler
	///
	/// - `ctx`: Uygulama bağlamı (kullanıcı bilgisi, izinler, vb. içerir)
	///
	/// # Döndürür
	///
	/// - `[]fields.Element`: Bağlama göre yapılandırılmış alan elemanları dizisi
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (l *ActiveUsersLens) GetFields(ctx *appContext.Context) []fields.Element {
	///     baseFields := l.Fields()
	///
	///     // Yönetici ise ek alanlar göster
	///     if ctx.User.IsAdmin() {
	///         baseFields = append(baseFields,
	///             fields.NewText("ip_address").SetLabel("IP Adresi"),
	///             fields.NewDateTime("created_at").SetLabel("Kayıt Tarihi"),
	///         )
	///     }
	///
	///     return baseFields
	/// }
	///
	/// // Kullanıcı diline göre alan etiketleri
	/// func (l *ProductsLens) GetFields(ctx *appContext.Context) []fields.Element {
	///     locale := ctx.GetLocale()
	///
	///     return []fields.Element{
	///         fields.NewText("name").SetLabel(translate("product.name", locale)),
	///         fields.NewNumber("price").SetLabel(translate("product.price", locale)),
	///     }
	/// }
	/// ```
	///
	/// # Önemli Notlar
	///
	/// - Context üzerinden kullanıcı bilgilerine, izinlere ve diğer bağlamsal verilere erişilebilir
	/// - Fields() metodundan farklı olarak, her istekte çalışır ve dinamik yapılandırma sağlar
	/// - Performans için ağır işlemleri cache'leyin
	///
	/// # İlgili Dokümantasyon
	///
	/// - Alan tipleri ve kullanımı için: docs/Fields.md
	GetFields(ctx *appContext.Context) []fields.Element

	/// Bu fonksiyon, lens görünümünde gösterilecek widget/card bileşenlerini döner.
	///
	/// # Parametreler
	///
	/// - `ctx`: Uygulama bağlamı (kullanıcı bilgisi, izinler, vb. içerir)
	///
	/// # Döndürür
	///
	/// - `[]widget.Card`: Widget/Card bileşenleri dizisi
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (l *ActiveUsersLens) GetCards(ctx *appContext.Context) []widget.Card {
	///     // Veritabanından istatistik çekme
	///     var totalActive int64
	///     ctx.DB.Model(&User{}).Where("status = ?", "active").Count(&totalActive)
	///
	///     var avgLoginTime float64
	///     ctx.DB.Model(&User{}).
	///         Where("status = ?", "active").
	///         Select("AVG(TIMESTAMPDIFF(HOUR, last_login, NOW()))").
	///         Scan(&avgLoginTime)
	///
	///     return []widget.Card{
	///         widget.NewMetricCard("Toplam Aktif Kullanıcı", fmt.Sprintf("%d", totalActive)),
	///         widget.NewMetricCard("Ort. Son Giriş", fmt.Sprintf("%.1f saat önce", avgLoginTime)),
	///         widget.NewTrendCard("Haftalık Artış", "+12%", "up"),
	///     }
	/// }
	///
	/// // Boş card listesi
	/// func (l *SimpleLens) GetCards(ctx *appContext.Context) []widget.Card {
	///     return []widget.Card{}
	/// }
	/// ```
	///
	/// # Önemli Notlar
	///
	/// - Card'lar, lens görünümünün üst kısmında gösterilir
	/// - Metrik, trend, grafik gibi farklı card tipleri kullanılabilir
	/// - Performans için ağır sorguları cache'leyin veya arka planda hesaplayın
	/// - Boş dizi dönülebilir (card gösterilmez)
	///
	/// # Uyarılar
	///
	/// ⚠️ Her istekte çalıştığı için performans kritiktir
	/// ⚠️ Karmaşık hesaplamalar sayfa yükleme süresini artırabilir
	GetCards(ctx *appContext.Context) []widget.Card

	/// Bu fonksiyon, Lens'in görünen adını döner (Name() ile aynı).
	///
	/// # Parametreler
	///
	/// Parametre almaz.
	///
	/// # Döndürür
	///
	/// - `string`: Lens'in görünen adı
	///
	/// # Önemli Notlar
	///
	/// - Bu metod, Name() metodunun alternatif bir versiyonudur
	/// - Genellikle Name() kullanımı tercih edilir
	GetName() string

	/// Bu fonksiyon, Lens'in URL tanımlayıcısını döner (Slug() ile aynı).
	///
	/// # Parametreler
	///
	/// Parametre almaz.
	///
	/// # Döndürür
	///
	/// - `string`: URL-safe slug
	///
	/// # Önemli Notlar
	///
	/// - Bu metod, Slug() metodunun alternatif bir versiyonudur
	/// - Genellikle Slug() kullanımı tercih edilir
	GetSlug() string

	/// Bu fonksiyon, lens-spesifik sorgu mantığını fonksiyon olarak döner.
	///
	/// # Parametreler
	///
	/// Parametre almaz.
	///
	/// # Döndürür
	///
	/// - `func(*gorm.DB) *gorm.DB`: Sorgu modifikasyon fonksiyonu
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (l *ActiveUsersLens) GetQuery() func(*gorm.DB) *gorm.DB {
	///     return func(db *gorm.DB) *gorm.DB {
	///         return db.Where("status = ?", "active")
	///     }
	/// }
	///
	/// // Query() metodunu yeniden kullanma
	/// func (l *ActiveUsersLens) GetQuery() func(*gorm.DB) *gorm.DB {
	///     return l.Query
	/// }
	/// ```
	///
	/// # Önemli Notlar
	///
	/// - Bu metod, Query() metodunun fonksiyonel bir versiyonudur
	/// - Genellikle Query() metodunu doğrudan döndürmek yeterlidir
	/// - Dinamik sorgu oluşturma senaryolarında kullanışlıdır
	GetQuery() func(*gorm.DB) *gorm.DB
}

// / Bu yapı, kaynak formlarının (ekleme/düzenleme/detay) kullanıcı arayüzünde nasıl sunulacağını belirler.
// /
// / # Genel Bakış
// /
// / DialogType, panel.go'da form görüntüleme modunu kontrol eden bir string sabiti türüdür.
// / Farklı cihaz tipleri ve kullanıcı deneyimi senaryoları için optimize edilmiş üç farklı sunum modu sunar.
// /
// / # Kullanım Senaryoları
// /
// / - **Masaüstü Uygulamalar**: Sheet veya Modal tercih edilir (geniş ekran alanı)
// / - **Mobil Uygulamalar**: Drawer tercih edilir (dokunmatik etkileşim için optimize)
// / - **Hızlı Düzenleme**: Sheet, yan panel olarak hızlı erişim sağlar
// / - **Odaklanma Gerektiren Formlar**: Modal, kullanıcının dikkatini forma yönlendirir
// / - **Çok Adımlı Formlar**: Drawer veya Modal, daha fazla alan sağlar
// /
// / # Avantajlar
// /
// / - **Esneklik**: Her kaynak için farklı dialog tipi seçilebilir
// / - **Responsive Tasarım**: Cihaz tipine göre en uygun görünüm
// / - **Kullanıcı Deneyimi**: Her senaryo için optimize edilmiş etkileşim
// / - **Tutarlılık**: Tüm kaynaklarda standart sunum modları
// /
// / # Dezavantajlar
// /
// / - **Mobil Sınırlamalar**: Sheet, küçük ekranlarda kullanışsız olabilir
// / - **Modal Engelleme**: Modal, arka plan etkileşimini tamamen engeller
// / - **Drawer Erişim**: Drawer, üst kısımdaki içeriğe erişimi zorlaştırabilir
// /
// / # Önemli Notlar
// /
// / - Varsayılan değer: DialogTypeSheet
// / - Her kaynak için SetDialogType() metodu ile özelleştirilebilir
// / - Dialog tipi, tüm form işlemleri için geçerlidir (create, edit, detail)
// / - Frontend tarafında otomatik olarak uygun bileşen render edilir
// /
// / # Kullanım Örneği
// /
// / ```go
// / // Resource tanımında dialog tipi belirleme
// / func (r *UserResource) Configure() {
// /     r.SetDialogType(resource.DialogTypeModal)
// / }
// /
// / // Koşullu dialog tipi seçimi
// / func (r *ProductResource) Configure() {
// /     if r.IsMobile() {
// /         r.SetDialogType(resource.DialogTypeDrawer)
// /     } else {
// /         r.SetDialogType(resource.DialogTypeSheet)
// /     }
// / }
// /
// / // Varsayılan değer kullanımı (Sheet)
// / func (r *OrderResource) Configure() {
// /     // SetDialogType çağrılmazsa DialogTypeSheet kullanılır
// / }
// / ```
// /
// / # Performans Notları
// /
// / - Dialog tipi seçimi, render performansını etkilemez
// / - Her dialog tipi, lazy loading destekler
// / - Modal ve Sheet, overlay rendering kullanır
// / - Drawer, transform animasyonları kullanır (GPU hızlandırmalı)
// /
// / # Uyarılar
// /
// / ⚠️ Dialog tipi değişikliği, mevcut açık formları etkilemez (sayfa yenileme gerekir)
// / ⚠️ Çok büyük formlar için Modal veya Drawer tercih edilmelidir (Sheet sınırlı genişliğe sahiptir)
type DialogType string

const (
	/// Bu sabit, formu sağdan açılan bir panel (Sheet) içinde gösterir.
	///
	/// # Özellikler
	///
	/// - **Konum**: Ekranın sağ tarafından açılır
	/// - **Genişlik**: Sabit genişlik (genellikle ekranın %30-40'ı)
	/// - **Arka Plan**: Yarı saydam overlay ile arka plan görünür kalır
	/// - **Etkileşim**: Arka plana tıklayarak kapatılabilir
	/// - **Animasyon**: Sağdan sola kayma animasyonu
	///
	/// # Kullanım Senaryoları
	///
	/// - Hızlı düzenleme işlemleri
	/// - Yan bilgi panelleri
	/// - Masaüstü uygulamalar
	/// - Liste görünümü ile birlikte çalışma
	/// - Çoklu form açma senaryoları
	///
	/// # Avantajlar
	///
	/// - Arka plan içeriği görünür kalır
	/// - Hızlı erişim ve kapatma
	/// - Çoklu sheet açılabilir (stack)
	/// - Liste ile birlikte görüntüleme
	/// - Minimal dikkat dağınıklığı
	///
	/// # Dezavantajlar
	///
	/// - Sınırlı genişlik (geniş formlar için uygun değil)
	/// - Mobil cihazlarda kullanışsız
	/// - Çok fazla alan gerektiren formlar için yetersiz
	///
	/// # Önemli Notlar
	///
	/// - Bu, varsayılan dialog tipidir
	/// - Masaüstü kullanıcı deneyimi için optimize edilmiştir
	/// - Responsive tasarımda mobilde otomatik olarak Drawer'a dönüşebilir
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (r *UserResource) Configure() {
	///     r.SetDialogType(resource.DialogTypeSheet)
	/// }
	/// ```
	///
	/// # Uyarılar
	///
	/// ⚠️ Çok geniş formlar için Modal veya Drawer tercih edin
	/// ⚠️ Mobil cihazlarda kullanıcı deneyimi düşük olabilir
	DialogTypeSheet DialogType = "sheet"

	/// Bu sabit, formu alttan açılan bir çekmece (Drawer) içinde gösterir.
	///
	/// # Özellikler
	///
	/// - **Konum**: Ekranın alt tarafından açılır
	/// - **Yükseklik**: Dinamik yükseklik (içeriğe göre veya tam ekran)
	/// - **Arka Plan**: Yarı saydam overlay ile arka plan görünür kalır
	/// - **Etkileşim**: Aşağı kaydırarak veya arka plana tıklayarak kapatılabilir
	/// - **Animasyon**: Alttan yukarı kayma animasyonu
	///
	/// # Kullanım Senaryoları
	///
	/// - Mobil uygulamalar (öncelikli)
	/// - Dokunmatik ekran cihazları
	/// - Hızlı işlemler (filtreleme, sıralama)
	/// - Kısa formlar
	/// - Native app benzeri deneyim
	///
	/// # Avantajlar
	///
	/// - Mobil için optimize edilmiş
	/// - Dokunmatik etkileşim dostu
	/// - Kaydırma ile kapatma (swipe to dismiss)
	/// - Tam ekran genişliği
	/// - Native app hissi
	/// - Erişilebilirlik (thumb-friendly)
	///
	/// # Dezavantajlar
	///
	/// - Masaüstünde Sheet kadar kullanışlı değil
	/// - Çok uzun formlar için scroll gerektirir
	/// - Üst kısımdaki içeriğe erişim zorlaşır
	///
	/// # Önemli Notlar
	///
	/// - Mobil cihazlar için en uygun seçenektir
	/// - iOS ve Android native drawer davranışını taklit eder
	/// - Responsive tasarımda otomatik olarak seçilebilir
	/// - Tam ekran moda geçiş desteklenir
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (r *MobileResource) Configure() {
	///     r.SetDialogType(resource.DialogTypeDrawer)
	/// }
	///
	/// // Cihaz tipine göre otomatik seçim
	/// func (r *ResponsiveResource) Configure() {
	///     if r.IsMobileDevice() {
	///         r.SetDialogType(resource.DialogTypeDrawer)
	///     } else {
	///         r.SetDialogType(resource.DialogTypeSheet)
	///     }
	/// }
	/// ```
	///
	/// # Uyarılar
	///
	/// ⚠️ Çok uzun formlar için tam ekran moda geçiş düşünün
	/// ⚠️ Masaüstü kullanıcıları için Sheet veya Modal daha uygun olabilir
	DialogTypeDrawer DialogType = "drawer"

	/// Bu sabit, formu ekranın ortasında klasik bir modal pencere içinde gösterir.
	///
	/// # Özellikler
	///
	/// - **Konum**: Ekranın tam ortasında
	/// - **Boyut**: Sabit veya dinamik boyut (içeriğe göre)
	/// - **Arka Plan**: Koyu overlay ile arka plan tamamen engellenir
	/// - **Etkileşim**: Kapatma butonu veya ESC tuşu ile kapatılır
	/// - **Animasyon**: Fade-in ve scale animasyonu
	///
	/// # Kullanım Senaryoları
	///
	/// - Önemli işlemler (silme, onay)
	/// - Odaklanma gerektiren formlar
	/// - Çok adımlı wizard formları
	/// - Kritik veri girişi
	/// - Klasik web uygulamaları
	///
	/// # Avantajlar
	///
	/// - Kullanıcının tam dikkatini çeker
	/// - Arka plan etkileşimini engeller
	/// - Tanıdık kullanıcı deneyimi
	/// - Esnek boyutlandırma
	/// - Tüm platformlarda tutarlı görünüm
	/// - Kritik işlemler için ideal
	///
	/// # Dezavantajlar
	///
	/// - Arka plan içeriği görünmez
	/// - Çoklu modal açılamaz (stack problemi)
	/// - Mobilde ekran alanını fazla kaplar
	/// - Liste ile birlikte görüntüleme yapılamaz
	///
	/// # Önemli Notlar
	///
	/// - Klasik web uygulamaları için tanıdık bir deneyim sunar
	/// - Kritik işlemler için kullanıcının dikkatini çeker
	/// - ESC tuşu ile kapatma desteklenir
	/// - Arka plan scroll'u otomatik olarak devre dışı bırakılır
	///
	/// # Kullanım Örneği
	///
	/// ```go
	/// func (r *ImportantResource) Configure() {
	///     r.SetDialogType(resource.DialogTypeModal)
	/// }
	///
	/// // Kritik işlemler için modal kullanımı
	/// func (r *PaymentResource) Configure() {
	///     // Ödeme formları için kullanıcının tam dikkatini çekmek
	///     r.SetDialogType(resource.DialogTypeModal)
	/// }
	///
	/// // Çok adımlı formlar için
	/// func (r *WizardResource) Configure() {
	///     r.SetDialogType(resource.DialogTypeModal)
	/// }
	/// ```
	///
	/// # Uyarılar
	///
	/// ⚠️ Çok sık kullanımı kullanıcı deneyimini olumsuz etkileyebilir
	/// ⚠️ Mobil cihazlarda ekran alanını fazla kaplar
	/// ⚠️ Çoklu modal açma senaryolarından kaçının (UX problemi)
	/// ⚠️ Uzun formlar için scroll yönetimi gerekebilir
	DialogTypeModal DialogType = "modal"
)

// DialogSize, resource formlarında kullanılacak modal/sheet genişlik preset'ini temsil eder.
type DialogSize string

const (
	DialogSizeSM   DialogSize = "sm"
	DialogSizeMD   DialogSize = "md"
	DialogSizeLG   DialogSize = "lg"
	DialogSizeXL   DialogSize = "xl"
	DialogSize2XL  DialogSize = "2xl"
	DialogSize3XL  DialogSize = "3xl"
	DialogSize4XL  DialogSize = "4xl"
	DialogSize5XL  DialogSize = "5xl"
	DialogSizeFull DialogSize = "full"
)
