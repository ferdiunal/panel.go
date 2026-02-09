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

/// # Action Interface
///
/// Bu interface, bir kaynak üzerinde gerçekleştirilebilecek özel toplu işlemleri (bulk actions) temsil eder.
/// Action'lar, kullanıcıların seçili kayıtlar üzerinde özel işlemler yapmasına olanak tanır.
///
/// ## Kullanım Senaryoları
///
/// - **Toplu Durum Değişikliği**: Seçili kayıtları aktif/pasif yapma
/// - **Toplu Silme**: Birden fazla kaydı aynı anda silme
/// - **E-posta Gönderimi**: Seçili kullanıcılara toplu e-posta gönderme
/// - **Dışa Aktarma**: Seçili kayıtları CSV/Excel formatında dışa aktarma
/// - **Toplu Güncelleme**: Seçili kayıtların belirli alanlarını güncelleme
/// - **Onay İşlemleri**: Bekleyen kayıtları toplu olarak onaylama/reddetme
///
/// ## Örnek Kullanım
///
/// ```go
/// type ActivateUsersAction struct{}
///
/// func (a *ActivateUsersAction) GetName() string {
///     return "Kullanıcıları Aktifleştir"
/// }
///
/// func (a *ActivateUsersAction) GetSlug() string {
///     return "activate-users"
/// }
///
/// func (a *ActivateUsersAction) GetIcon() string {
///     return "check-circle" // Lucide icon
/// }
///
/// func (a *ActivateUsersAction) Execute(ctx *appContext.Context, items []any) error {
///     for _, item := range items {
///         user := item.(*User)
///         user.IsActive = true
///         if err := ctx.DB.Save(user).Error; err != nil {
///             return err
///         }
///     }
///     return nil
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Action'lar transaction içinde çalıştırılmalıdır
/// - Hata durumunda tüm işlemler geri alınmalıdır (rollback)
/// - Kullanıcı yetkilendirmesi Execute metodunda kontrol edilmelidir
/// - İşlem sonucu kullanıcıya bildirim gösterilmelidir
///
/// ## Avantajlar
///
/// - Tekrarlayan işlemleri otomatikleştirir
/// - Kullanıcı deneyimini iyileştirir
/// - Kod tekrarını azaltır
/// - Merkezi hata yönetimi sağlar
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - Büyük veri setlerinde performans sorunları yaşanabilir
/// - İşlem süresi uzun olabilir, timeout ayarları yapılmalıdır
/// - Concurrent işlemlerde race condition'a dikkat edilmelidir
type Action interface {
	/// GetName, işlemin kullanıcı arayüzünde görünecek Türkçe adını döner.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: İşlemin görünen adı (örn: "Kullanıcıları Aktifleştir", "E-posta Gönder")
	///
	/// ## Örnek
	/// ```go
	/// func (a *ActivateAction) GetName() string {
	///     return "Aktifleştir"
	/// }
	/// ```
	GetName() string

	/// GetSlug, işlemin URL ve API endpoint'lerinde kullanılacak benzersiz tanımlayıcısını döner.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: URL-safe slug (örn: "activate-users", "send-email", "export-csv")
	///
	/// ## Örnek
	/// ```go
	/// func (a *ActivateAction) GetSlug() string {
	///     return "activate-users"
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - Slug kebab-case formatında olmalıdır
	/// - Özel karakter içermemelidir
	/// - Kaynak içinde benzersiz olmalıdır
	GetSlug() string

	/// GetIcon, işlemin kullanıcı arayüzünde gösterilecek ikon adını döner.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: Lucide icon set'inden ikon adı (örn: "check-circle", "mail", "trash")
	///
	/// ## Örnek
	/// ```go
	/// func (a *DeleteAction) GetIcon() string {
	///     return "trash-2"
	/// }
	/// ```
	///
	/// ## Referans
	/// Lucide icon listesi: https://lucide.dev/icons/
	GetIcon() string

	/// Execute, seçili kayıtlar üzerinde işlemi gerçekleştirir.
	///
	/// ## Parametreler
	/// - `ctx`: İstek bağlamı, veritabanı bağlantısı ve kullanıcı bilgilerini içerir
	/// - `items`: İşlem yapılacak kayıtların listesi (any tipinde, type assertion gereklidir)
	///
	/// ## Döndürür
	/// - `error`: İşlem başarısız olursa hata, başarılı olursa nil
	///
	/// ## Örnek
	/// ```go
	/// func (a *ActivateAction) Execute(ctx *appContext.Context, items []any) error {
	///     return ctx.DB.Transaction(func(tx *gorm.DB) error {
	///         for _, item := range items {
	///             user, ok := item.(*User)
	///             if !ok {
	///                 return fmt.Errorf("geçersiz kayıt tipi")
	///             }
	///
	///             user.IsActive = true
	///             if err := tx.Save(user).Error; err != nil {
	///                 return err
	///             }
	///         }
	///         return nil
	///     })
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - Transaction kullanımı önerilir
	/// - Type assertion ile kayıt tipini kontrol edin
	/// - Yetkilendirme kontrolü yapın
	/// - Hata mesajları kullanıcı dostu olmalıdır
	Execute(ctx *appContext.Context, items []any) error
}

/// # Filter Interface
///
/// Bu interface, kaynak listelerinde uygulanabilecek filtreleme seçeneklerini temsil eder.
/// Filtreler, kullanıcıların büyük veri setlerini daraltarak istedikleri kayıtları bulmalarını sağlar.
///
/// ## Kullanım Senaryoları
///
/// - **Durum Filtreleme**: Aktif/pasif, yayınlanmış/taslak kayıtları filtreleme
/// - **Tarih Aralığı**: Belirli tarih aralığındaki kayıtları gösterme
/// - **Kategori Filtreleme**: Belirli kategorilere ait kayıtları listeleme
/// - **Fiyat Aralığı**: Min-max fiyat aralığında ürün filtreleme
/// - **Kullanıcı Filtreleme**: Belirli kullanıcıya ait kayıtları gösterme
/// - **Boolean Filtreler**: Öne çıkan, onaylanmış vb. kayıtları filtreleme
///
/// ## Örnek Kullanım
///
/// ```go
/// type StatusFilter struct{}
///
/// func (f *StatusFilter) GetName() string {
///     return "Durum"
/// }
///
/// func (f *StatusFilter) GetSlug() string {
///     return "status"
/// }
///
/// func (f *StatusFilter) GetType() string {
///     return "select" // select, date, range, boolean
/// }
///
/// func (f *StatusFilter) GetOptions() map[string]string {
///     return map[string]string{
///         "active":   "Aktif",
///         "inactive": "Pasif",
///         "pending":  "Beklemede",
///     }
/// }
///
/// func (f *StatusFilter) Apply(db *gorm.DB, value any) *gorm.DB {
///     if status, ok := value.(string); ok && status != "" {
///         return db.Where("status = ?", status)
///     }
///     return db
/// }
/// ```
///
/// ## Filtre Tipleri
///
/// - **select**: Açılır liste (dropdown) ile tek seçim
/// - **multiselect**: Çoklu seçim listesi
/// - **date**: Tarih seçici
/// - **daterange**: Tarih aralığı seçici
/// - **range**: Sayısal aralık (min-max)
/// - **boolean**: Evet/Hayır seçimi
/// - **search**: Metin arama kutusu
///
/// ## Önemli Notlar
///
/// - Filtreler URL query parametrelerine yansıtılır
/// - Filtre değerleri sayfa yenilendiğinde korunur
/// - Birden fazla filtre aynı anda uygulanabilir
/// - Filtre değerleri validate edilmelidir
///
/// ## Avantajlar
///
/// - Kullanıcı deneyimini iyileştirir
/// - Büyük veri setlerinde gezinmeyi kolaylaştırır
/// - Performanslı sorgular oluşturur
/// - URL ile paylaşılabilir filtreler
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - SQL injection'a karşı parameterized query kullanın
/// - Geçersiz değerleri kontrol edin
/// - İndeksli sütunlarda filtreleme yapın
/// - Karmaşık filtrelerde performans testleri yapın
type Filter interface {
	/// GetName, filtrenin kullanıcı arayüzünde görünecek Türkçe adını döner.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: Filtrenin görünen adı (örn: "Durum", "Tarih Aralığı", "Kategori")
	///
	/// ## Örnek
	/// ```go
	/// func (f *StatusFilter) GetName() string {
	///     return "Durum"
	/// }
	/// ```
	GetName() string

	/// GetSlug, filtrenin URL query parametrelerinde kullanılacak benzersiz tanımlayıcısını döner.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: URL-safe slug (örn: "status", "date-range", "category")
	///
	/// ## Örnek
	/// ```go
	/// func (f *StatusFilter) GetSlug() string {
	///     return "status"
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - Slug kebab-case formatında olmalıdır
	/// - URL'de query parameter olarak kullanılır: ?status=active
	/// - Kaynak içinde benzersiz olmalıdır
	GetSlug() string

	/// GetType, filtrenin kullanıcı arayüzünde nasıl gösterileceğini belirten tip bilgisini döner.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: Filtre tipi (select, multiselect, date, daterange, range, boolean, search)
	///
	/// ## Örnek
	/// ```go
	/// func (f *StatusFilter) GetType() string {
	///     return "select"
	/// }
	///
	/// func (f *DateRangeFilter) GetType() string {
	///     return "daterange"
	/// }
	/// ```
	///
	/// ## Desteklenen Tipler
	/// - `select`: Tek seçimli dropdown
	/// - `multiselect`: Çoklu seçim
	/// - `date`: Tarih seçici
	/// - `daterange`: Başlangıç-bitiş tarihi
	/// - `range`: Min-max sayısal aralık
	/// - `boolean`: Evet/Hayır
	/// - `search`: Metin arama
	GetType() string

	/// GetOptions, select ve multiselect tipindeki filtreler için seçenekleri döner.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `map[string]string`: Key-value çiftleri (key: değer, value: görünen etiket)
	///
	/// ## Örnek
	/// ```go
	/// func (f *StatusFilter) GetOptions() map[string]string {
	///     return map[string]string{
	///         "draft":     "Taslak",
	///         "published": "Yayınlandı",
	///         "archived":  "Arşivlendi",
	///     }
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - Sadece select ve multiselect tipleri için gereklidir
	/// - Diğer tipler için boş map dönebilir
	/// - Options dinamik olarak veritabanından da yüklenebilir
	GetOptions() map[string]string

	/// Apply, filtre değerini GORM sorgu nesnesine uygular ve filtrelenmiş sorguyu döner.
	///
	/// ## Parametreler
	/// - `db`: GORM veritabanı sorgu nesnesi
	/// - `value`: Kullanıcının seçtiği filtre değeri (tip filtreye göre değişir)
	///
	/// ## Döndürür
	/// - `*gorm.DB`: Filtre uygulanmış GORM sorgu nesnesi
	///
	/// ## Örnek
	/// ```go
	/// // Select filtresi
	/// func (f *StatusFilter) Apply(db *gorm.DB, value any) *gorm.DB {
	///     if status, ok := value.(string); ok && status != "" {
	///         return db.Where("status = ?", status)
	///     }
	///     return db
	/// }
	///
	/// // Date range filtresi
	/// func (f *DateRangeFilter) Apply(db *gorm.DB, value any) *gorm.DB {
	///     if dateRange, ok := value.(map[string]string); ok {
	///         if start, ok := dateRange["start"]; ok && start != "" {
	///             db = db.Where("created_at >= ?", start)
	///         }
	///         if end, ok := dateRange["end"]; ok && end != "" {
	///             db = db.Where("created_at <= ?", end)
	///         }
	///     }
	///     return db
	/// }
	///
	/// // Range filtresi
	/// func (f *PriceRangeFilter) Apply(db *gorm.DB, value any) *gorm.DB {
	///     if priceRange, ok := value.(map[string]float64); ok {
	///         if min, ok := priceRange["min"]; ok {
	///             db = db.Where("price >= ?", min)
	///         }
	///         if max, ok := priceRange["max"]; ok {
	///             db = db.Where("price <= ?", max)
	///         }
	///     }
	///     return db
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - SQL injection'a karşı mutlaka parameterized query kullanın
	/// - Type assertion ile değer tipini kontrol edin
	/// - Geçersiz değerler için orijinal db nesnesini döndürün
	/// - İndeksli sütunlarda filtreleme yapın
	Apply(db *gorm.DB, value any) *gorm.DB
}

/// # Sortable Struct
///
/// Bu yapı, liste görünümlerinde varsayılan sıralama ayarlarını tanımlar.
/// Kaynak ilk yüklendiğinde hangi sütuna göre ve hangi yönde sıralanacağını belirler.
///
/// ## Kullanım Senaryoları
///
/// - **Tarih Sıralaması**: En yeni kayıtları üstte gösterme (created_at DESC)
/// - **Alfabetik Sıralama**: Kayıtları ada göre A-Z sıralama (name ASC)
/// - **Öncelik Sıralaması**: Önemli kayıtları üstte gösterme (priority DESC)
/// - **Fiyat Sıralaması**: Ürünleri fiyata göre sıralama (price ASC/DESC)
/// - **Durum Sıralaması**: Aktif kayıtları üstte gösterme (is_active DESC)
///
/// ## Örnek Kullanım
///
/// ```go
/// // Resource'da varsayılan sıralama tanımlama
/// func (r *PostResource) GetSortable() []resource.Sortable {
///     return []resource.Sortable{
///         {
///             Column:    "created_at",
///             Direction: "desc",
///         },
///     }
/// }
///
/// // Çoklu sıralama
/// func (r *UserResource) GetSortable() []resource.Sortable {
///     return []resource.Sortable{
///         {
///             Column:    "is_active",
///             Direction: "desc",
///         },
///         {
///             Column:    "created_at",
///             Direction: "desc",
///         },
///     }
/// }
///
/// // Alfabetik sıralama
/// func (r *CategoryResource) GetSortable() []resource.Sortable {
///     return []resource.Sortable{
///         {
///             Column:    "name",
///             Direction: "asc",
///         },
///     }
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Sıralama sütunu veritabanında mevcut olmalıdır
/// - Performans için sıralama sütunlarına indeks ekleyin
/// - Direction değeri "asc" veya "desc" olmalıdır
/// - Çoklu sıralama için birden fazla Sortable tanımlanabilir
/// - Kullanıcı arayüzden sıralamayı değiştirebilir
///
/// ## Avantajlar
///
/// - Tutarlı kullanıcı deneyimi sağlar
/// - En önemli kayıtlar önce gösterilir
/// - Performanslı sorgular oluşturur
/// - Kullanıcı beklentilerini karşılar
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - İndekslenmemiş sütunlarda sıralama yavaş olabilir
/// - Büyük veri setlerinde performans testleri yapın
/// - NULL değerleri olan sütunlarda sıralama davranışını kontrol edin
type Sortable struct {
	/// Column, sıralanacak veritabanı sütun adını belirtir.
	///
	/// ## Örnek Değerler
	/// - "created_at": Oluşturulma tarihine göre
	/// - "name": İsme göre
	/// - "price": Fiyata göre
	/// - "is_active": Aktiflik durumuna göre
	/// - "priority": Önceliğe göre
	///
	/// ## Önemli Notlar
	/// - Sütun adı veritabanı şemasıyla eşleşmelidir
	/// - Snake_case formatında olmalıdır
	/// - İndeksli sütunlar tercih edilmelidir
	Column string

	/// Direction, sıralama yönünü belirtir.
	///
	/// ## Geçerli Değerler
	/// - "asc": Artan sıralama (A-Z, 0-9, eski-yeni)
	/// - "desc": Azalan sıralama (Z-A, 9-0, yeni-eski)
	///
	/// ## Örnek
	/// ```go
	/// // En yeni kayıtlar üstte
	/// Direction: "desc"
	///
	/// // Alfabetik sıralama
	/// Direction: "asc"
	/// ```
	///
	/// ## Önemli Notlar
	/// - Küçük harfle yazılmalıdır
	/// - Geçersiz değerler varsayılan olarak "asc" kabul edilir
	Direction string
}

/// # Resource Interface
///
/// Bu interface, paneldeki her bir varlığı (örneğin Users, Posts, Products) temsil eder.
/// Resource, bir CRUD (Create, Read, Update, Delete) kaynağının tüm özelliklerini tanımlar:
/// veri modeli, alanlar, ilişkiler, yetkilendirme, görünüm ayarları ve özel işlemler.
///
/// ## Kullanım Senaryoları
///
/// - **Kullanıcı Yönetimi**: Kullanıcıları listeleme, ekleme, düzenleme, silme
/// - **İçerik Yönetimi**: Blog yazıları, sayfalar, yorumlar
/// - **E-Ticaret**: Ürünler, kategoriler, siparişler, müşteriler
/// - **Medya Yönetimi**: Görseller, videolar, dosyalar
/// - **Sistem Ayarları**: Yapılandırma, roller, izinler
///
/// ## Temel Örnek
///
/// ```go
/// type UserResource struct {
///     resource.OptimizedBase
/// }
///
/// func (r *UserResource) Model() any {
///     return &User{}
/// }
///
/// func (r *UserResource) Slug() string {
///     return "users"
/// }
///
/// func (r *UserResource) Title() string {
///     return "Kullanıcılar"
/// }
///
/// func (r *UserResource) Icon() string {
///     return "users"
/// }
///
/// func (r *UserResource) Fields() []fields.Element {
///     return []fields.Element{
///         fields.ID("ID").Sortable(),
///         fields.Text("Ad", "name").Required().Searchable(),
///         fields.Email("E-posta", "email").Required().Unique("users", "email"),
///         fields.DateTime("Kayıt Tarihi", "created_at").ReadOnly(),
///     }
/// }
/// ```
///
/// ## Referanslar
///
/// - **Alan Sistemi**: Detaylı bilgi için [docs/Fields.md](../../docs/Fields.md) dosyasına bakın
/// - **İlişkiler**: İlişki tanımlamaları için [docs/Relationships.md](../../docs/Relationships.md) dosyasına bakın
///
/// ## Önemli Notlar
///
/// - Her resource benzersiz bir slug'a sahip olmalıdır
/// - Model() metodu GORM model struct'ının pointer'ını dönmelidir
/// - Fields() metodu tüm alanları tanımlar, GetFields() ise bağlama göre filtreler
/// - Yetkilendirme için Policy() metodunu implement edin
/// - OptimizedBase kullanarak varsayılan implementasyonları kullanabilirsiniz
///
/// ## Avantajlar
///
/// - Tip güvenli CRUD işlemleri
/// - Otomatik API endpoint oluşturma
/// - Dinamik form ve liste görünümleri
/// - Merkezi yetkilendirme kontrolü
/// - Kolay özelleştirme ve genişletme
type Resource interface {
	/// Model, GORM model yapısının bir örneğini döner.
	///
	/// Bu metod, kaynağın veritabanı modelini tanımlar. GORM bu modeli kullanarak
	/// tablo yapısını, ilişkileri ve sorguları yönetir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `any`: GORM model struct'ının pointer'ı (örn: &User{}, &Post{})
	///
	/// ## Örnek
	/// ```go
	/// type User struct {
	///     ID        uint      `json:"id" gorm:"primaryKey"`
	///     Name      string    `json:"name" gorm:"size:255"`
	///     Email     string    `json:"email" gorm:"uniqueIndex;size:255"`
	///     CreatedAt time.Time `json:"createdAt"`
	/// }
	///
	/// func (r *UserResource) Model() any {
	///     return &User{}
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - Mutlaka pointer dönmelidir (&User{}, User{} değil)
	/// - GORM tag'leri ile veritabanı yapısını tanımlayın
	/// - JSON tag'leri ile API response formatını belirleyin
	/// - Soft delete için gorm.DeletedAt kullanın
	Model() any

	/// Fields, kaynak form ve listelerinde gösterilecek tüm alanları tanımlar.
	///
	/// Bu metod, kaynağın tüm alanlarını (fields) döner. Alanlar liste, detay ve form
	/// görünümlerinde kullanılır. Her alan kendi görünürlük ayarlarına sahiptir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]fields.Element`: Alan listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) Fields() []fields.Element {
	///     return []fields.Element{
	///         fields.ID("ID").Sortable().OnlyOnDetail(),
	///         fields.Text("Ad", "name").Required().Searchable().OnList().OnForm(),
	///         fields.Email("E-posta", "email").Required().OnList().OnForm(),
	///         fields.Password("Şifre", "password").Required().OnlyOnCreate(),
	///         fields.Switch("Aktif", "is_active").OnList().OnForm(),
	///         fields.DateTime("Kayıt Tarihi", "created_at").ReadOnly().OnList().OnDetail(),
	///     }
	/// }
	/// ```
	///
	/// ## Referans
	/// Alan tipleri ve özellikleri için [docs/Fields.md](../../docs/Fields.md) dosyasına bakın.
	///
	/// ## Önemli Notlar
	/// - OnList(), OnDetail(), OnForm() ile görünürlük kontrol edilir
	/// - Required(), Unique() gibi validasyon kuralları ekleyin
	/// - Searchable(), Sortable() ile liste özelliklerini belirleyin
	/// - ReadOnly() ile salt okunur alanlar tanımlayın
	Fields() []fields.Element

	/// GetFields, belirli bir bağlamda (context) gösterilecek alanları döner.
	///
	/// Bu metod, Fields() metodundan dönen alanları kullanıcı yetkilerine, istek tipine
	/// veya diğer bağlamsal faktörlere göre filtreler. Dinamik alan görünürlüğü için kullanılır.
	///
	/// ## Parametreler
	/// - `ctx`: İstek bağlamı, kullanıcı bilgileri ve yetkileri içerir
	///
	/// ## Döndürür
	/// - `[]fields.Element`: Filtrelenmiş alan listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) GetFields(ctx *appContext.Context) []fields.Element {
	///     allFields := r.Fields()
	///
	///     // Admin olmayan kullanıcılar için hassas alanları gizle
	///     if !ctx.User.IsAdmin() {
	///         return filterFields(allFields, func(f fields.Element) bool {
	///             return f.GetKey() != "salary" && f.GetKey() != "ssn"
	///         })
	///     }
	///
	///     return allFields
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - Kullanıcı rolüne göre alan görünürlüğü
	/// - Kayıt durumuna göre alan filtreleme
	/// - Özel izinlere göre hassas alanları gizleme
	/// - Dinamik form yapılandırması
	///
	/// ## Önemli Notlar
	/// - OptimizedBase varsayılan olarak tüm alanları döner
	/// - Özel filtreleme için bu metodu override edin
	/// - CanSee() callback'i ile alan bazında kontrol de yapılabilir
	GetFields(ctx *appContext.Context) []fields.Element

	/// With, ilişkili verilerin eager loading ile yüklenmesi için ilişki adlarını döner.
	///
	/// Bu metod, N+1 query problemini önlemek için GORM'un Preload özelliğini kullanır.
	/// Belirtilen ilişkiler sorgu sırasında otomatik olarak yüklenir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]string`: Eager load edilecek ilişki adları
	///
	/// ## Örnek
	/// ```go
	/// // Basit ilişkiler
	/// func (r *PostResource) With() []string {
	///     return []string{"Author", "Category", "Tags"}
	/// }
	///
	/// // İç içe ilişkiler
	/// func (r *OrderResource) With() []string {
	///     return []string{
	///         "Customer",
	///         "Items",
	///         "Items.Product",
	///         "Items.Product.Category",
	///     }
	/// }
	///
	/// // Koşullu eager loading
	/// func (r *UserResource) With() []string {
	///     return []string{
	///         "Profile",
	///         "Roles",
	///     }
	/// }
	/// ```
	///
	/// ## Referans
	/// İlişki tipleri için [docs/Relationships.md](../../docs/Relationships.md) dosyasına bakın.
	///
	/// ## Önemli Notlar
	/// - İlişki adları GORM model'deki field adlarıyla eşleşmelidir
	/// - Nokta notasyonu ile iç içe ilişkiler yüklenebilir
	/// - Gereksiz ilişkiler performansı düşürür, sadece gerekenleri ekleyin
	/// - Liste görünümünde kullanılan ilişkileri mutlaka ekleyin
	///
	/// ## Avantajlar
	/// - N+1 query problemini önler
	/// - Performansı önemli ölçüde artırır
	/// - Tek sorguda tüm ilişkili verileri yükler
	///
	/// ## Dikkat Edilmesi Gerekenler
	/// - Çok fazla ilişki yüklemek bellek kullanımını artırır
	/// - Büyük koleksiyonlarda sayfalama kullanın
	/// - Gereksiz iç içe ilişkilerden kaçının
	With() []string
	/// Lenses, özel veri filtreleme görünümlerini (lens) tanımlar.
	///
	/// Lens'ler, kaynağın verilerini önceden tanımlanmış filtrelerle görüntülemeye yarayan
	/// özel görünümlerdir. Kullanıcılar lens'ler arasında hızlıca geçiş yapabilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]Lens`: Lens listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) Lenses() []resource.Lens {
	///     return []resource.Lens{
	///         &ActiveUsersLens{},
	///         &InactiveUsersLens{},
	///         &RecentlyRegisteredLens{},
	///     }
	/// }
	///
	/// // Lens implementasyonu
	/// type ActiveUsersLens struct{}
	///
	/// func (l *ActiveUsersLens) Name() string {
	///     return "Aktif Kullanıcılar"
	/// }
	///
	/// func (l *ActiveUsersLens) Slug() string {
	///     return "active"
	/// }
	///
	/// func (l *ActiveUsersLens) Apply(db *gorm.DB) *gorm.DB {
	///     return db.Where("is_active = ?", true)
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - Aktif/pasif kayıtları görüntüleme
	/// - Onay bekleyen kayıtları listeleme
	/// - Son 30 gün içindeki kayıtları gösterme
	/// - Öne çıkan içerikleri filtreleme
	/// - Belirli durumdaki kayıtları gruplama
	///
	/// ## Önemli Notlar
	/// - Lens'ler URL'de query parameter olarak saklanır
	/// - Her lens benzersiz bir slug'a sahip olmalıdır
	/// - Lens filtreleri performanslı olmalıdır
	/// - Kullanıcı lens'ler arasında hızlıca geçiş yapabilir
	Lenses() []Lens

	/// GetLenses, kaynağın tüm lens'lerini döner.
	///
	/// Bu metod, Lenses() metodundan dönen lens'leri kullanıcı yetkilerine göre
	/// filtreleyebilir veya dinamik lens'ler ekleyebilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]Lens`: Filtrelenmiş lens listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) GetLenses() []resource.Lens {
	///     lenses := r.Lenses()
	///
	///     // Admin olmayan kullanıcılar için bazı lens'leri gizle
	///     if !ctx.User.IsAdmin() {
	///         return filterLenses(lenses, func(l resource.Lens) bool {
	///             return l.Slug() != "deleted"
	///         })
	///     }
	///
	///     return lenses
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - OptimizedBase varsayılan olarak tüm lens'leri döner
	/// - Özel filtreleme için bu metodu override edin
	GetLenses() []Lens

	/// Slug, kaynağın URL'de kullanılacak benzersiz tanımlayıcısını döner.
	///
	/// Slug, kaynağın tüm URL endpoint'lerinde kullanılır ve benzersiz olmalıdır.
	/// Genellikle model adının çoğul ve küçük harf halidir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: URL-safe slug (örn: "users", "blog-posts", "products")
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) Slug() string {
	///     return "users"
	/// }
	///
	/// func (r *BlogPostResource) Slug() string {
	///     return "blog-posts"
	/// }
	///
	/// func (r *ProductResource) Slug() string {
	///     return "products"
	/// }
	/// ```
	///
	/// ## URL Kullanımı
	/// Slug aşağıdaki endpoint'lerde kullanılır:
	/// - Liste: `/api/resources/{slug}`
	/// - Detay: `/api/resources/{slug}/{id}`
	/// - Oluştur: `/api/resources/{slug}`
	/// - Güncelle: `/api/resources/{slug}/{id}`
	/// - Sil: `/api/resources/{slug}/{id}`
	///
	/// ## Önemli Notlar
	/// - Slug kebab-case formatında olmalıdır
	/// - Özel karakter içermemelidir (sadece harf, rakam, tire)
	/// - Tüm panel içinde benzersiz olmalıdır
	/// - Değiştirilmemelidir (URL'ler bozulur)
	/// - Genellikle çoğul isim kullanılır
	Slug() string

	/// Title, kaynağın kullanıcı arayüzünde görünecek Türkçe başlığını döner.
	///
	/// Bu başlık menüde, sayfa başlığında ve breadcrumb'larda kullanılır.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: Kullanıcı dostu başlık (örn: "Kullanıcılar", "Blog Yazıları", "Ürünler")
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) Title() string {
	///     return "Kullanıcılar"
	/// }
	///
	/// func (r *BlogPostResource) Title() string {
	///     return "Blog Yazıları"
	/// }
	///
	/// func (r *ProductResource) Title() string {
	///     return "Ürünler"
	/// }
	/// ```
	///
	/// ## Kullanım Yerleri
	/// - Navigasyon menüsü
	/// - Sayfa başlığı (browser tab)
	/// - Breadcrumb navigasyonu
	/// - Bildirim mesajları
	///
	/// ## Önemli Notlar
	/// - Türkçe karakter kullanılabilir
	/// - Genellikle çoğul isim kullanılır
	/// - Kısa ve açıklayıcı olmalıdır
	/// - Büyük harfle başlamalıdır
	Title() string

	/// Icon, kaynağın menüde gösterilecek ikon adını döner.
	///
	/// İkon, Lucide icon set'inden seçilmelidir. Menü ve navigasyonda
	/// kaynağı görsel olarak temsil eder.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: Lucide icon adı (örn: "users", "file-text", "shopping-cart")
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) Icon() string {
	///     return "users"
	/// }
	///
	/// func (r *BlogPostResource) Icon() string {
	///     return "file-text"
	/// }
	///
	/// func (r *ProductResource) Icon() string {
	///     return "shopping-cart"
	/// }
	///
	/// func (r *SettingsResource) Icon() string {
	///     return "settings"
	/// }
	/// ```
	///
	/// ## Popüler İkonlar
	/// - `users`: Kullanıcılar
	/// - `file-text`: Dökümanlar, yazılar
	/// - `shopping-cart`: E-ticaret, ürünler
	/// - `image`: Medya, görseller
	/// - `mail`: E-posta, mesajlar
	/// - `settings`: Ayarlar
	/// - `folder`: Kategoriler, klasörler
	/// - `tag`: Etiketler
	/// - `calendar`: Takvim, etkinlikler
	/// - `package`: Paketler, siparişler
	///
	/// ## Referans
	/// Tüm ikonlar için: https://lucide.dev/icons/
	///
	/// ## Önemli Notlar
	/// - İkon adı kebab-case formatında olmalıdır
	/// - Lucide icon set'inde mevcut olmalıdır
	/// - Kaynağı görsel olarak temsil etmelidir
	/// - Tutarlı ikon kullanımı önerilir
	Icon() string

	/// Group, kaynağın menüde hangi grup altında listeleneceğini belirler.
	///
	/// Gruplar, ilgili kaynakları menüde organize etmek için kullanılır.
	/// Aynı gruptaki kaynaklar menüde birlikte gösterilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `string`: Grup adı (örn: "İçerik Yönetimi", "E-Ticaret", "Sistem")
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) Group() string {
	///     return "Kullanıcı Yönetimi"
	/// }
	///
	/// func (r *BlogPostResource) Group() string {
	///     return "İçerik Yönetimi"
	/// }
	///
	/// func (r *ProductResource) Group() string {
	///     return "E-Ticaret"
	/// }
	///
	/// func (r *SettingsResource) Group() string {
	///     return "Sistem"
	/// }
	///
	/// // Grupsuz kaynak (ana menüde gösterilir)
	/// func (r *DashboardResource) Group() string {
	///     return ""
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - **İçerik Yönetimi**: Blog yazıları, sayfalar, yorumlar
	/// - **E-Ticaret**: Ürünler, siparişler, kategoriler
	/// - **Kullanıcı Yönetimi**: Kullanıcılar, roller, izinler
	/// - **Medya**: Görseller, videolar, dosyalar
	/// - **Sistem**: Ayarlar, loglar, raporlar
	///
	/// ## Önemli Notlar
	/// - Boş string döndürülürse kaynak grupsuz olur
	/// - Aynı grup adı birden fazla kaynakta kullanılabilir
	/// - Grup adları Türkçe olabilir
	/// - Tutarlı gruplama kullanıcı deneyimini iyileştirir
	Group() string
	/// Policy, kaynak üzerindeki yetkilendirme (CRUD) kurallarını tanımlayan policy nesnesini döner.
	///
	/// Policy, kullanıcıların kaynak üzerinde hangi işlemleri yapabileceğini kontrol eder.
	/// Her CRUD işlemi (Create, Read, Update, Delete) için ayrı yetkilendirme kuralları tanımlanabilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `auth.Policy`: Yetkilendirme policy nesnesi
	///
	/// ## Örnek
	/// ```go
	/// type UserPolicy struct{}
	///
	/// func (p *UserPolicy) ViewAny(ctx *appContext.Context) bool {
	///     return ctx.User.Can("view-users")
	/// }
	///
	/// func (p *UserPolicy) View(ctx *appContext.Context, item any) bool {
	///     user := item.(*User)
	///     return ctx.User.Can("view-users") || ctx.User.ID == user.ID
	/// }
	///
	/// func (p *UserPolicy) Create(ctx *appContext.Context) bool {
	///     return ctx.User.Can("create-users")
	/// }
	///
	/// func (p *UserPolicy) Update(ctx *appContext.Context, item any) bool {
	///     user := item.(*User)
	///     return ctx.User.Can("update-users") || ctx.User.ID == user.ID
	/// }
	///
	/// func (p *UserPolicy) Delete(ctx *appContext.Context, item any) bool {
	///     user := item.(*User)
	///     return ctx.User.Can("delete-users") && ctx.User.ID != user.ID
	/// }
	///
	/// func (r *UserResource) Policy() auth.Policy {
	///     return &UserPolicy{}
	/// }
	/// ```
	///
	/// ## Policy Metodları
	/// - `ViewAny`: Liste görünümüne erişim kontrolü
	/// - `View`: Tek kayıt detayına erişim kontrolü
	/// - `Create`: Yeni kayıt oluşturma yetkisi
	/// - `Update`: Kayıt güncelleme yetkisi
	/// - `Delete`: Kayıt silme yetkisi
	///
	/// ## Kullanım Senaryoları
	/// - Rol tabanlı yetkilendirme (admin, editor, viewer)
	/// - Kayıt sahibi kontrolü (kullanıcı sadece kendi kayıtlarını düzenleyebilir)
	/// - Durum bazlı yetkilendirme (yayınlanmış içerik silinemez)
	/// - Özel izin kontrolleri (belirli alanları sadece admin görebilir)
	///
	/// ## Önemli Notlar
	/// - Policy metodları her istekte çalıştırılır
	/// - False döndürülürse 403 Forbidden hatası döner
	/// - Nil policy döndürülürse tüm işlemlere izin verilir
	/// - Policy kontrolü middleware seviyesinde yapılır
	Policy() auth.Policy

	/// GetPolicy, kaynağın yetkilendirme politikasını döner.
	///
	/// Bu metod, Policy() metodundan dönen policy'yi kullanıcı bağlamına göre
	/// özelleştirebilir veya dinamik policy'ler oluşturabilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `auth.Policy`: Yetkilendirme policy nesnesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) GetPolicy() auth.Policy {
	///     // Varsayılan policy'yi döndür
	///     return r.Policy()
	/// }
	///
	/// // Dinamik policy örneği
	/// func (r *PostResource) GetPolicy() auth.Policy {
	///     if someCondition {
	///         return &StrictPostPolicy{}
	///     }
	///     return &RelaxedPostPolicy{}
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - OptimizedBase varsayılan olarak Policy() metodunu çağırır
	/// - Özel policy mantığı için bu metodu override edin
	GetPolicy() auth.Policy

	/// GetSortable, liste görünümünde varsayılan sıralama ayarlarını döner.
	///
	/// Bu metod, kaynak ilk yüklendiğinde hangi sütuna göre ve hangi yönde
	/// sıralanacağını belirler. Kullanıcı daha sonra sıralamayı değiştirebilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]Sortable`: Sıralama ayarları listesi
	///
	/// ## Örnek
	/// ```go
	/// // Tek sıralama
	/// func (r *PostResource) GetSortable() []resource.Sortable {
	///     return []resource.Sortable{
	///         {Column: "created_at", Direction: "desc"},
	///     }
	/// }
	///
	/// // Çoklu sıralama (önce aktiflik, sonra tarih)
	/// func (r *UserResource) GetSortable() []resource.Sortable {
	///     return []resource.Sortable{
	///         {Column: "is_active", Direction: "desc"},
	///         {Column: "created_at", Direction: "desc"},
	///     }
	/// }
	///
	/// // Alfabetik sıralama
	/// func (r *CategoryResource) GetSortable() []resource.Sortable {
	///     return []resource.Sortable{
	///         {Column: "name", Direction: "asc"},
	///     }
	/// }
	///
	/// // Sıralama yok (veritabanı varsayılanı)
	/// func (r *LogResource) GetSortable() []resource.Sortable {
	///     return []resource.Sortable{}
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - Sıralama sütunları veritabanında mevcut olmalıdır
	/// - Performans için sıralama sütunlarına indeks ekleyin
	/// - Direction "asc" veya "desc" olmalıdır
	/// - Boş liste döndürülürse varsayılan sıralama kullanılır
	GetSortable() []Sortable

	/// GetDialogType, ekleme/düzenleme formunun hangi tipte açılacağını döner.
	///
	/// Dialog tipi, formun kullanıcı arayüzünde nasıl gösterileceğini belirler.
	/// Farklı tipler farklı kullanıcı deneyimleri sunar.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `DialogType`: Dialog tipi (Sheet, Drawer, Modal)
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) GetDialogType() resource.DialogType {
	///     return resource.DialogTypeSheet
	/// }
	///
	/// func (r *SettingsResource) GetDialogType() resource.DialogType {
	///     return resource.DialogTypeModal
	/// }
	///
	/// func (r *NotificationResource) GetDialogType() resource.DialogType {
	///     return resource.DialogTypeDrawer
	/// }
	/// ```
	///
	/// ## Dialog Tipleri
	/// - **Sheet**: Sayfanın altından yukarı açılan panel (mobil uyumlu)
	/// - **Drawer**: Sayfanın yanından açılan panel (geniş formlar için)
	/// - **Modal**: Sayfanın ortasında açılan popup (küçük formlar için)
	///
	/// ## Kullanım Senaryoları
	/// - **Sheet**: Hızlı düzenleme, mobil cihazlar
	/// - **Drawer**: Çok alanlı formlar, detaylı düzenleme
	/// - **Modal**: Basit formlar, onay diyalogları
	///
	/// ## Önemli Notlar
	/// - SetDialogType() ile runtime'da değiştirilebilir
	/// - Mobil cihazlarda Sheet önerilir
	/// - Çok alanlı formlarda Drawer kullanın
	GetDialogType() DialogType

	/// SetDialogType, form görünüm tipini ayarlar ve kaynağı döner (method chaining).
	///
	/// Bu metod, dialog tipini dinamik olarak değiştirmek için kullanılır.
	/// Method chaining desteği sayesinde zincirleme çağrılar yapılabilir.
	///
	/// ## Parametreler
	/// - `dialogType`: Yeni dialog tipi (Sheet, Drawer, Modal)
	///
	/// ## Döndürür
	/// - Yapılandırılmış Resource pointer'ı (method chaining için)
	///
	/// ## Örnek
	/// ```go
	/// // Method chaining ile kullanım
	/// resource := &UserResource{}
	/// resource.SetDialogType(resource.DialogTypeSheet)
	///
	/// // Zincirleme çağrılar
	/// resource.
	///     SetDialogType(resource.DialogTypeDrawer).
	///     // Diğer ayarlar...
	///
	/// // Koşullu ayarlama
	/// if isMobile {
	///     resource.SetDialogType(resource.DialogTypeSheet)
	/// } else {
	///     resource.SetDialogType(resource.DialogTypeModal)
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - Method chaining için Resource pointer'ı döner
	/// - Runtime'da dialog tipini değiştirebilir
	/// - OptimizedBase bu metodu implement eder
	SetDialogType(DialogType) Resource
	/// Repository, kaynağın veri erişim katmanını (DataProvider) döner.
	///
	/// Bu metod, varsayılan GORM DataProvider yerine özel bir repository implementasyonu
	/// kullanmak için override edilebilir. Özel veri kaynakları, cache katmanları veya
	/// karmaşık sorgular için kullanılır.
	///
	/// ## Parametreler
	/// - `db`: GORM veritabanı bağlantısı
	///
	/// ## Döndürür
	/// - `data.DataProvider`: Veri erişim katmanı implementasyonu
	///
	/// ## Örnek
	/// ```go
	/// // Varsayılan GORM provider kullanımı
	/// func (r *UserResource) Repository(db *gorm.DB) data.DataProvider {
	///     return data.NewGormDataProvider(db, r.Model())
	/// }
	///
	/// // Özel repository implementasyonu
	/// type CachedUserRepository struct {
	///     db    *gorm.DB
	///     cache *redis.Client
	/// }
	///
	/// func (r *CachedUserRepository) Find(id any) (any, error) {
	///     // Önce cache'den kontrol et
	///     if cached, err := r.cache.Get(fmt.Sprintf("user:%v", id)); err == nil {
	///         return cached, nil
	///     }
	///     // Cache'de yoksa veritabanından al
	///     user := &User{}
	///     err := r.db.First(user, id).Error
	///     if err == nil {
	///         r.cache.Set(fmt.Sprintf("user:%v", id), user, time.Hour)
	///     }
	///     return user, err
	/// }
	///
	/// func (r *UserResource) Repository(db *gorm.DB) data.DataProvider {
	///     return &CachedUserRepository{
	///         db:    db,
	///         cache: redisClient,
	///     }
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - Cache katmanı ekleme (Redis, Memcached)
	/// - Özel sorgular ve filtreleme mantığı
	/// - Çoklu veritabanı desteği
	/// - Audit logging ekleme
	/// - Soft delete özelleştirmesi
	/// - Performans optimizasyonları
	///
	/// ## Önemli Notlar
	/// - DataProvider interface'ini implement etmelidir
	/// - OptimizedBase varsayılan GormDataProvider kullanır
	/// - Özel repository thread-safe olmalıdır
	/// - Transaction desteği sağlanmalıdır
	Repository(db *gorm.DB) data.DataProvider

	/// Cards, kaynak dashboard'unda gösterilecek widget/card'ları döner.
	///
	/// Card'lar, kaynak hakkında özet bilgiler, istatistikler ve metrikler gösterir.
	/// Dashboard görünümünde kullanıcıya hızlı bilgi sunar.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]widget.Card`: Card/widget listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) Cards() []widget.Card {
	///     return []widget.Card{
	///         &widget.ValueCard{
	///             Title: "Toplam Kullanıcı",
	///             Value: func(ctx *appContext.Context) (any, error) {
	///                 var count int64
	///                 err := ctx.DB.Model(&User{}).Count(&count).Error
	///                 return count, err
	///             },
	///             Icon: "users",
	///         },
	///         &widget.ValueCard{
	///             Title: "Aktif Kullanıcılar",
	///             Value: func(ctx *appContext.Context) (any, error) {
	///                 var count int64
	///                 err := ctx.DB.Model(&User{}).Where("is_active = ?", true).Count(&count).Error
	///                 return count, err
	///             },
	///             Icon: "user-check",
	///             Color: "green",
	///         },
	///         &widget.TrendCard{
	///             Title: "Yeni Kayıtlar (30 Gün)",
	///             Value: func(ctx *appContext.Context) (any, error) {
	///                 var count int64
	///                 thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	///                 err := ctx.DB.Model(&User{}).
	///                     Where("created_at >= ?", thirtyDaysAgo).
	///                     Count(&count).Error
	///                 return count, err
	///             },
	///             Trend: "+12%",
	///             Icon: "trending-up",
	///         },
	///     }
	/// }
	/// ```
	///
	/// ## Card Tipleri
	/// - **ValueCard**: Tek bir değer gösterir (sayı, metin)
	/// - **TrendCard**: Değer ve trend gösterir (artış/azalış)
	/// - **ChartCard**: Grafik gösterir (line, bar, pie)
	/// - **ListCard**: Liste gösterir (son kayıtlar, aktiviteler)
	///
	/// ## Kullanım Senaryoları
	/// - Toplam kayıt sayısı
	/// - Aktif/pasif kayıt istatistikleri
	/// - Günlük/haftalık/aylık trendler
	/// - Gelir/gider özeti
	/// - Son aktiviteler
	/// - Performans metrikleri
	///
	/// ## Önemli Notlar
	/// - Card'lar her sayfa yüklendiğinde çalıştırılır
	/// - Performanslı sorgular kullanın
	/// - Cache kullanımı önerilir
	/// - Hata durumlarını handle edin
	Cards() []widget.Card

	/// GetCards, belirli bir bağlamda gösterilecek card'ları döner.
	///
	/// Bu metod, Cards() metodundan dönen card'ları kullanıcı yetkilerine göre
	/// filtreleyebilir veya dinamik card'ler ekleyebilir.
	///
	/// ## Parametreler
	/// - `ctx`: İstek bağlamı, kullanıcı bilgileri ve yetkileri içerir
	///
	/// ## Döndürür
	/// - `[]widget.Card`: Filtrelenmiş card listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) GetCards(ctx *appContext.Context) []widget.Card {
	///     cards := r.Cards()
	///
	///     // Admin olmayan kullanıcılar için bazı card'ları gizle
	///     if !ctx.User.IsAdmin() {
	///         return filterCards(cards, func(c widget.Card) bool {
	///             return c.GetTitle() != "Gelir Özeti"
	///         })
	///     }
	///
	///     // Dinamik card ekleme
	///     if ctx.User.Can("view-analytics") {
	///         cards = append(cards, &widget.ChartCard{
	///             Title: "Kullanıcı Aktivitesi",
	///             Type:  "line",
	///             Data:  getActivityData(ctx),
	///         })
	///     }
	///
	///     return cards
	/// }
	/// ```
	///
	/// ## Önemli Notlar
	/// - OptimizedBase varsayılan olarak tüm card'ları döner
	/// - Özel filtreleme için bu metodu override edin
	/// - Performans için gereksiz card'ları gizleyin
	GetCards(ctx *appContext.Context) []widget.Card

	/// ResolveField, bir alanın değerini dinamik olarak hesaplar ve dönüştürür.
	///
	/// Bu metod, veritabanında olmayan computed field'lar veya özel formatlama
	/// gerektiren alanlar için kullanılır. Alan değeri her okunduğunda çalıştırılır.
	///
	/// ## Parametreler
	/// - `fieldName`: Çözümlenecek alan adı
	/// - `item`: Kayıt nesnesi (any tipinde, type assertion gereklidir)
	///
	/// ## Döndürür
	/// - `any`: Hesaplanmış alan değeri
	/// - `error`: Hata durumunda hata mesajı
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) ResolveField(fieldName string, item any) (any, error) {
	///     user, ok := item.(*User)
	///     if !ok {
	///         return nil, fmt.Errorf("geçersiz kayıt tipi")
	///     }
	///
	///     switch fieldName {
	///     case "full_name":
	///         // Computed field: ad + soyad
	///         return user.FirstName + " " + user.LastName, nil
	///
	///     case "age":
	///         // Computed field: yaş hesaplama
	///         if user.BirthDate.IsZero() {
	///             return nil, nil
	///         }
	///         age := time.Now().Year() - user.BirthDate.Year()
	///         return age, nil
	///
	///     case "avatar_url":
	///         // URL oluşturma
	///         if user.Avatar == "" {
	///             return "/default-avatar.png", nil
	///         }
	///         return "/storage/" + user.Avatar, nil
	///
	///     case "status_label":
	///         // Durum etiketi
	///         if user.IsActive {
	///             return "Aktif", nil
	///         }
	///         return "Pasif", nil
	///
	///     case "order_count":
	///         // İlişkili kayıt sayısı
	///         var count int64
	///         db.Model(&Order{}).Where("user_id = ?", user.ID).Count(&count)
	///         return count, nil
	///     }
	///
	///     return nil, fmt.Errorf("alan bulunamadı: %s", fieldName)
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - **Computed Fields**: Birden fazla alandan hesaplanan değerler
	/// - **Formatlama**: Tarih, para birimi, yüzde formatlaması
	/// - **URL Oluşturma**: Dosya yollarını tam URL'ye çevirme
	/// - **İlişki Sayıları**: İlişkili kayıt sayılarını hesaplama
	/// - **Durum Etiketleri**: Kod değerlerini kullanıcı dostu metne çevirme
	/// - **Özel Hesaplamalar**: Yaş, süre, mesafe vb. hesaplamalar
	///
	/// ## Önemli Notlar
	/// - Her alan okunduğunda çalıştırılır, performansa dikkat edin
	/// - Type assertion ile kayıt tipini kontrol edin
	/// - Hata durumlarını handle edin
	/// - Veritabanı sorguları minimize edilmelidir
	/// - Cache kullanımı önerilir
	/// - Nil değerleri handle edin
	///
	/// ## Avantajlar
	/// - Veritabanında gereksiz sütun oluşturmaz
	/// - Dinamik değer hesaplama
	/// - Esnek formatlama
	/// - Merkezi hesaplama mantığı
	///
	/// ## Dikkat Edilmesi Gerekenler
	/// - N+1 query problemine dikkat edin
	/// - Ağır hesaplamalardan kaçının
	/// - Cache kullanın
	/// - Hata mesajları açıklayıcı olmalıdır
	ResolveField(fieldName string, item any) (any, error)
	/// GetActions, kaynağın özel toplu işlemlerini (bulk actions) döner.
	///
	/// Bu metod, kullanıcıların seçili kayıtlar üzerinde gerçekleştirebileceği
	/// özel işlemleri tanımlar. Action'lar liste görünümünde gösterilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]Action`: Action listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) GetActions() []resource.Action {
	///     return []resource.Action{
	///         &ActivateUsersAction{},
	///         &DeactivateUsersAction{},
	///         &SendEmailAction{},
	///         &ExportCSVAction{},
	///     }
	/// }
	///
	/// // Action implementasyonu
	/// type ActivateUsersAction struct{}
	///
	/// func (a *ActivateUsersAction) GetName() string {
	///     return "Aktifleştir"
	/// }
	///
	/// func (a *ActivateUsersAction) GetSlug() string {
	///     return "activate"
	/// }
	///
	/// func (a *ActivateUsersAction) GetIcon() string {
	///     return "check-circle"
	/// }
	///
	/// func (a *ActivateUsersAction) Execute(ctx *appContext.Context, items []any) error {
	///     return ctx.DB.Transaction(func(tx *gorm.DB) error {
	///         for _, item := range items {
	///             user := item.(*User)
	///             user.IsActive = true
	///             if err := tx.Save(user).Error; err != nil {
	///                 return err
	///             }
	///         }
	///         return nil
	///     })
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - Toplu durum değişikliği (aktif/pasif, onay/red)
	/// - Toplu silme işlemleri
	/// - E-posta gönderimi
	/// - Dışa aktarma (CSV, Excel, PDF)
	/// - Toplu güncelleme
	/// - Arşivleme/geri yükleme
	///
	/// ## Önemli Notlar
	/// - Action'lar kullanıcı yetkilerine göre filtrelenebilir
	/// - Transaction kullanımı önerilir
	/// - Hata durumlarını handle edin
	/// - İşlem sonucu kullanıcıya bildirim gösterin
	GetActions() []Action

	/// GetFilters, kaynağın filtreleme seçeneklerini döner.
	///
	/// Bu metod, kullanıcıların liste görünümünde kayıtları filtrelemek için
	/// kullanabileceği filtreleri tanımlar. Filtreler URL query parametrelerine yansıtılır.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `[]Filter`: Filter listesi
	///
	/// ## Örnek
	/// ```go
	/// func (r *UserResource) GetFilters() []resource.Filter {
	///     return []resource.Filter{
	///         &StatusFilter{},
	///         &RoleFilter{},
	///         &DateRangeFilter{},
	///         &ActiveFilter{},
	///     }
	/// }
	///
	/// // Filter implementasyonu
	/// type StatusFilter struct{}
	///
	/// func (f *StatusFilter) GetName() string {
	///     return "Durum"
	/// }
	///
	/// func (f *StatusFilter) GetSlug() string {
	///     return "status"
	/// }
	///
	/// func (f *StatusFilter) GetType() string {
	///     return "select"
	/// }
	///
	/// func (f *StatusFilter) GetOptions() map[string]string {
	///     return map[string]string{
	///         "active":   "Aktif",
	///         "inactive": "Pasif",
	///         "pending":  "Beklemede",
	///     }
	/// }
	///
	/// func (f *StatusFilter) Apply(db *gorm.DB, value any) *gorm.DB {
	///     if status, ok := value.(string); ok && status != "" {
	///         return db.Where("status = ?", status)
	///     }
	///     return db
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - Durum filtreleme (aktif/pasif, yayınlanmış/taslak)
	/// - Tarih aralığı filtreleme
	/// - Kategori/grup filtreleme
	/// - Fiyat aralığı filtreleme
	/// - Boolean filtreler (öne çıkan, onaylanmış)
	/// - Kullanıcı/yazar filtreleme
	///
	/// ## Önemli Notlar
	/// - Filtreler kullanıcı yetkilerine göre filtrelenebilir
	/// - İndeksli sütunlarda filtreleme yapın
	/// - SQL injection'a karşı parameterized query kullanın
	/// - Filtre değerleri URL'de saklanır
	GetFilters() []Filter

	/// StoreHandler, dosya yükleme işlemlerini yönetir ve dosya yolunu döner.
	///
	/// Bu metod, Image, Video, Audio, File gibi dosya alanlarında yüklenen
	/// dosyaların nasıl saklanacağını kontrol eder. Özel dosya işleme mantığı
	/// için override edilebilir.
	///
	/// ## Parametreler
	/// - `c`: İstek bağlamı
	/// - `file`: Yüklenen dosya header'ı
	/// - `storagePath`: Dosyanın kaydedileceği dizin yolu
	/// - `storageURL`: Dosyaya erişim için kullanılacak URL prefix'i
	///
	/// ## Döndürür
	/// - `string`: Kaydedilen dosyanın yolu/URL'i
	/// - `error`: Hata durumunda hata mesajı
	///
	/// ## Örnek
	/// ```go
	/// // Varsayılan implementasyon
	/// func (r *UserResource) StoreHandler(
	///     c *appContext.Context,
	///     file *multipart.FileHeader,
	///     storagePath string,
	///     storageURL string,
	/// ) (string, error) {
	///     // Dosya adını benzersiz yap
	///     filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	///     filepath := filepath.Join(storagePath, filename)
	///
	///     // Dosyayı kaydet
	///     if err := c.SaveFile(file, filepath); err != nil {
	///         return "", err
	///     }
	///
	///     // URL döndür
	///     return storageURL + "/" + filename, nil
	/// }
	///
	/// // S3 storage örneği
	/// func (r *ProductResource) StoreHandler(
	///     c *appContext.Context,
	///     file *multipart.FileHeader,
	///     storagePath string,
	///     storageURL string,
	/// ) (string, error) {
	///     // Dosyayı aç
	///     src, err := file.Open()
	///     if err != nil {
	///         return "", err
	///     }
	///     defer src.Close()
	///
	///     // S3'e yükle
	///     filename := fmt.Sprintf("%s/%d_%s", storagePath, time.Now().Unix(), file.Filename)
	///     url, err := s3Client.Upload(filename, src)
	///     if err != nil {
	///         return "", err
	///     }
	///
	///     return url, nil
	/// }
	///
	/// // Görsel işleme örneği
	/// func (r *UserResource) StoreHandler(
	///     c *appContext.Context,
	///     file *multipart.FileHeader,
	///     storagePath string,
	///     storageURL string,
	/// ) (string, error) {
	///     // Dosyayı aç
	///     src, err := file.Open()
	///     if err != nil {
	///         return "", err
	///     }
	///     defer src.Close()
	///
	///     // Görseli decode et
	///     img, _, err := image.Decode(src)
	///     if err != nil {
	///         return "", err
	///     }
	///
	///     // Yeniden boyutlandır
	///     resized := resize.Resize(800, 0, img, resize.Lanczos3)
	///
	///     // Kaydet
	///     filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	///     filepath := filepath.Join(storagePath, filename)
	///     out, err := os.Create(filepath)
	///     if err != nil {
	///         return "", err
	///     }
	///     defer out.Close()
	///
	///     jpeg.Encode(out, resized, &jpeg.Options{Quality: 85})
	///
	///     return storageURL + "/" + filename, nil
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - Yerel dosya sistemi storage
	/// - Cloud storage (S3, Google Cloud Storage, Azure Blob)
	/// - CDN entegrasyonu
	/// - Görsel işleme (resize, crop, watermark)
	/// - Video transcoding
	/// - Dosya sıkıştırma
	/// - EXIF data temizleme
	/// - Virus tarama
	///
	/// ## Önemli Notlar
	/// - Dosya boyutu kontrolü yapın
	/// - Dosya tipi validasyonu yapın
	/// - Benzersiz dosya adları kullanın
	/// - Hata durumlarını handle edin
	/// - Güvenlik kontrollerini atlamamayın
	/// - Temporary dosyaları temizleyin
	///
	/// ## Güvenlik Uyarıları
	/// - Dosya uzantısı kontrolü yapın
	/// - MIME type doğrulaması yapın
	/// - Dosya boyutu limiti koyun
	/// - Path traversal saldırılarına karşı koruma sağlayın
	/// - Executable dosyaları reddedin
	StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error)

	/// NavigationOrder, kaynağın menüdeki sıralama önceliğini döner.
	///
	/// Bu metod, kaynakların menüde hangi sırada gösterileceğini belirler.
	/// Düşük sayılar üst sırada, yüksek sayılar alt sırada gösterilir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `int`: Sıralama önceliği (düşük sayı = üst sıra)
	///
	/// ## Örnek
	/// ```go
	/// // Dashboard en üstte
	/// func (r *DashboardResource) NavigationOrder() int {
	///     return 1
	/// }
	///
	/// // Kullanıcılar ikinci sırada
	/// func (r *UserResource) NavigationOrder() int {
	///     return 10
	/// }
	///
	/// // Blog yazıları üçüncü sırada
	/// func (r *PostResource) NavigationOrder() int {
	///     return 20
	/// }
	///
	/// // Ayarlar en altta
	/// func (r *SettingsResource) NavigationOrder() int {
	///     return 1000
	/// }
	///
	/// // Varsayılan sıralama (0)
	/// func (r *CategoryResource) NavigationOrder() int {
	///     return 0
	/// }
	/// ```
	///
	/// ## Sıralama Stratejisi
	/// - **1-10**: Önemli kaynaklar (Dashboard, Ana Sayfa)
	/// - **10-100**: Sık kullanılan kaynaklar (Kullanıcılar, İçerik)
	/// - **100-500**: Normal kaynaklar (Kategoriler, Etiketler)
	/// - **500-1000**: Az kullanılan kaynaklar (Loglar, Raporlar)
	/// - **1000+**: Sistem kaynakları (Ayarlar, Yapılandırma)
	///
	/// ## Önemli Notlar
	/// - Aynı grup içindeki kaynaklar kendi aralarında sıralanır
	/// - 0 değeri varsayılan sıralamayı kullanır
	/// - Negatif değerler kullanılabilir (en üst sıra için)
	/// - Tutarlı aralıklar kullanın (10, 20, 30 gibi)
	NavigationOrder() int

	/// Visible, kaynağın menüde görünüp görünmeyeceğini belirler.
	///
	/// Bu metod, kaynağın navigasyon menüsünde gösterilip gösterilmeyeceğini
	/// kontrol eder. API endpoint'leri aktif kalır, sadece menü görünürlüğü etkilenir.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `bool`: true = menüde görünür, false = menüde gizli
	///
	/// ## Örnek
	/// ```go
	/// // Normal kaynak (menüde görünür)
	/// func (r *UserResource) Visible() bool {
	///     return true
	/// }
	///
	/// // API-only kaynak (menüde gizli)
	/// func (r *WebhookResource) Visible() bool {
	///     return false
	/// }
	///
	/// // Sistem kaynağı (menüde gizli)
	/// func (r *AuditLogResource) Visible() bool {
	///     return false
	/// }
	///
	/// // Koşullu görünürlük
	/// func (r *AdminResource) Visible() bool {
	///     // Sadece admin kullanıcılar için görünür
	///     // Not: Bu örnekte ctx erişimi yok, Policy kullanın
	///     return true
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - **API-Only Kaynaklar**: Sadece API üzerinden erişilen kaynaklar
	/// - **Sistem Kaynakları**: Audit log, sistem ayarları gibi teknik kaynaklar
	/// - **Gizli Kaynaklar**: Webhook, callback gibi arka plan kaynakları
	/// - **Geçici Kaynaklar**: Geliştirme aşamasındaki kaynaklar
	/// - **Yardımcı Kaynaklar**: Diğer kaynaklar tarafından kullanılan kaynaklar
	///
	/// ## Önemli Notlar
	/// - false döndürülse bile API endpoint'leri aktif kalır
	/// - Yetkilendirme için Policy kullanın, Visible() değil
	/// - Menü görünürlüğü kullanıcı bazlı değil, kaynak bazlıdır
	/// - Dinamik görünürlük için Policy kullanın
	///
	/// ## Avantajlar
	/// - Menü karmaşıklığını azaltır
	/// - Teknik kaynakları gizler
	/// - Kullanıcı deneyimini iyileştirir
	/// - API erişimini korur
	Visible() bool

	/// OpenAPIEnabled, kaynağın OpenAPI spesifikasyonunda görünüp görünmeyeceğini belirler.
	///
	/// Bu metod, kaynağın OpenAPI/Swagger dokümantasyonunda gösterilip gösterilmeyeceğini
	/// kontrol eder. false döndürülürse, kaynak OpenAPI spec'e dahil edilmez ve
	/// Swagger UI'da görünmez.
	///
	/// ## Parametreler
	/// Parametre almaz.
	///
	/// ## Döndürür
	/// - `bool`: true = OpenAPI spec'te görünür, false = OpenAPI spec'te gizli
	///
	/// ## Örnek
	/// ```go
	/// // Normal kaynak (OpenAPI'de görünür)
	/// func (r *UserResource) OpenAPIEnabled() bool {
	///     return true
	/// }
	///
	/// // Internal API kaynak (OpenAPI'de gizli)
	/// func (r *InternalResource) OpenAPIEnabled() bool {
	///     return false
	/// }
	///
	/// // Geliştirme aşamasındaki kaynak (OpenAPI'de gizli)
	/// func (r *BetaResource) OpenAPIEnabled() bool {
	///     return false
	/// }
	/// ```
	///
	/// ## Kullanım Senaryoları
	/// - **Internal API'ler**: Sadece sistem içi kullanım için olan kaynaklar
	/// - **Deprecated API'ler**: Kullanımdan kaldırılacak kaynaklar
	/// - **Beta/Experimental**: Henüz stabil olmayan kaynaklar
	/// - **Admin-Only**: Sadece admin kullanıcılar için olan kaynaklar
	/// - **Legacy API'ler**: Geriye dönük uyumluluk için tutulan kaynaklar
	///
	/// ## Önemli Notlar
	/// - false döndürülse bile API endpoint'leri aktif kalır
	/// - Sadece OpenAPI dokümantasyonu etkilenir
	/// - Yetkilendirme için Policy kullanın, OpenAPIEnabled() değil
	/// - Varsayılan değer true'dur (tüm kaynaklar OpenAPI'de görünür)
	///
	/// ## Avantajlar
	/// - Dokümantasyon karmaşıklığını azaltır
	/// - Internal API'leri gizler
	/// - Public API'yi daha temiz gösterir
	/// - API versiyonlama stratejisini destekler
	OpenAPIEnabled() bool

	/// RecordTitle, bir kayıt için okunabilir başlık döndürür.
	///
	/// Bu metod, ilişki fieldlarında kayıtların kullanıcı dostu şekilde gösterilmesi için kullanılır.
	/// Laravel Nova'nın title() pattern'ini takip eder.
	///
	/// Parametreler:
	/// - record: Başlığı alınacak kayıt (genellikle model instance'ı)
	///
	/// Dönüş:
	/// - string: Kaydın okunabilir başlığı
	///
	/// Örnek:
	///   user := &User{ID: 1, Name: "John Doe"}
	///   title := resource.RecordTitle(user) // "John Doe"
	RecordTitle(record any) string

	/// GetRecordTitleKey, kayıt başlığı için kullanılacak field adını döndürür.
	///
	/// Bu metod, RecordTitle metodunun hangi field'ı kullanacağını belirler.
	/// Varsayılan değer "id"'dir.
	///
	/// Dönüş:
	/// - string: Başlık için kullanılacak field adı
	///
	/// Örnek:
	///   key := resource.GetRecordTitleKey() // "name"
	GetRecordTitleKey() string

	/// SetRecordTitleKey, kayıt başlığı için kullanılacak field adını ayarlar.
	///
	/// Bu metod, RecordTitle metodunun hangi field'ı kullanacağını belirler.
	/// Fluent interface için Resource döndürür.
	///
	/// Parametreler:
	/// - key: Başlık için kullanılacak field adı
	///
	/// Dönüş:
	/// - Resource: Fluent interface için resource instance'ı
	///
	/// Örnek:
	///   resource.SetRecordTitleKey("name")
	SetRecordTitleKey(key string) Resource
}
