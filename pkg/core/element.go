// Package core, panel sisteminin temel interface'lerini ve tiplerini sağlar.
//
// # Genel Bakış
//
// Bu paket, tüm mimari için temel oluşturur ve dahili paketlere bağımlılığı yoktur.
// Diğer tüm paketler core'a bağımlı olmalıdır, tersi olmamalıdır.
//
// # Ana Bileşenler
//
// Core paketi şunları tanımlar:
//   - Element interface: Form ve liste alanları için ortak interface
//   - ResourceContext: Resource işlemleri sırasında kullanılan context
//   - ElementType ve ElementContext: Tip tanımları ve sabitler
//   - Callback fonksiyon tipleri: Görünürlük ve depolama işlemleri için
//
// # Mimari Prensipler
//
// 1. **Bağımlılık Yönü**: Core paketi hiçbir dahili pakete bağımlı değildir
// 2. **Interface Odaklı**: Tüm temel davranışlar interface'ler üzerinden tanımlanır
// 3. **Genişletilebilirlik**: Yeni alan tipleri Element interface'ini implement ederek eklenebilir
// 4. **Tip Güvenliği**: Compile-time tip kontrolü için güçlü tipler kullanılır
//
// # Kullanım Senaryoları
//
// Core paketi şu durumlarda kullanılır:
//   - Yeni alan tipleri oluşturma (Element interface'ini implement ederek)
//   - Resource context'e erişim gerektiren işlemler
//   - Görünürlük ve depolama callback'leri tanımlama
//   - Alan meta verilerine erişim
//
// # Önemli Notlar
//
// - Bu paket, panel sisteminin temel taşıdır ve değişiklikler tüm sistemi etkiler
// - Element interface'i, tüm alan tiplerinin ortak davranışlarını tanımlar
// - Callback fonksiyonları, dinamik davranış kontrolü için kullanılır
//
// # İlgili Dokümantasyon
//
// Detaylı bilgi için bakınız:
//   - docs/Fields.md: Alan sistemi ve kullanımı
//   - docs/Relationships.md: İlişki alanları ve yapılandırması
//
// # Örnek Kullanım
//
//	// Özel bir alan tipi oluşturma
//	type CustomField struct {
//	    *fields.Schema
//	}
//
//	func (f *CustomField) Extract(resource any) {
//	    // Özel veri çıkarma mantığı
//	}
//
//	func (f *CustomField) JsonSerialize() map[string]any {
//	    // Özel serileştirme mantığı
//	    return f.Schema.JsonSerialize()
//	}
package core

import (
	"github.com/gofiber/fiber/v2"
)

// AutoOptionsConfig, otomatik seçenek üretimi için yapılandırma bilgilerini tutar.
//
// # Genel Bakış
//
// Bu yapı, ilişki alanlarında (BelongsTo, HasOne, BelongsToMany) seçeneklerin
// veritabanından otomatik olarak yüklenmesi için kullanılır. Manuel olarak
// Options callback'i tanımlamaya gerek kalmadan, ilişkili kayıtlar otomatik
// olarak çekilir ve frontend'e gönderilir.
//
// # Kullanım Senaryoları
//
// AutoOptionsConfig şu durumlarda kullanılır:
//   - BelongsTo ilişkilerinde parent kayıtların listelenmesi
//   - HasOne ilişkilerinde boşta olan kayıtların listelenmesi
//   - BelongsToMany ilişkilerinde tüm ilişkili kayıtların listelenmesi
//
// # Alanlar
//
// Enabled: Otomatik seçenek üretiminin aktif olup olmadığını belirtir.
// true ise, backend otomatik olarak ilgili tablodan kayıtları çeker.
//
// DisplayField: İlişkili kayıtlarda hangi alanın label olarak gösterileceğini belirtir.
// Örneğin "name", "title", "email" gibi kullanıcı dostu bir alan seçilmelidir.
//
// # Önemli Notlar
//
// - HasOne ilişkilerinde, AutoOptions sadece "boşta olan" kayıtları getirir (foreign_key IS NULL)
// - BelongsTo ilişkilerinde, AutoOptions tüm parent kayıtları getirir
// - DisplayField, ilişkili tabloda mevcut bir sütun olmalıdır
// - Edit modunda, mevcut ilişkili kayıt otomatik olarak seçeneklere dahil edilir
//
// # Avantajlar
//
// - Manuel Options callback'i tanımlamaya gerek kalmaz
// - Kod tekrarını azaltır
// - Tutarlı davranış sağlar
// - Otomatik filtreleme (HasOne için)
//
// # Dezavantajlar
//
// - Karmaşık filtreleme gerektiren durumlarda yetersiz kalabilir
// - Çok sayıda kayıt varsa performans sorunu yaratabilir
// - Özel sıralama veya gruplama yapılamaz
//
// # İlgili Dokümantasyon
//
// Detaylı bilgi için bakınız:
//   - docs/Relationships.md: İlişki alanları ve AutoOptions kullanımı
//   - docs/Fields.md: Alan sistemi genel bakış
//
// # Örnek Kullanım
//
//	// BelongsTo ilişkisinde AutoOptions
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    AutoOptions("name") // Tüm author'ları getirir, name alanını gösterir
//
//	// HasOne ilişkisinde AutoOptions
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    AutoOptions("bio") // Sadece user_id'si boş olan profilleri getirir
//
//	// BelongsToMany ilişkisinde AutoOptions
//	field := fields.BelongsToMany("Roles", "roles", "roles").
//	    AutoOptions("name") // Tüm role'leri getirir, name alanını gösterir
type AutoOptionsConfig struct {
	Enabled      bool   // Otomatik seçenek üretimi aktif mi?
	DisplayField string // Hangi alan label olarak gösterilecek?
}

// Element, form ve liste görünümlerinde kullanılan alanlar için ortak interface'dir.
//
// # Genel Bakış
//
// Bu interface, panel sistemindeki tüm alan tiplerinin (Text, Number, Select, BelongsTo, vb.)
// uygulaması gereken temel davranışları tanımlar. Veri çıkarma, serileştirme, görünürlük
// kontrolü ve fluent yapılandırma metodları sağlar.
//
// # Temel Sorumluluklar
//
// Element interface şu sorumlulukları tanımlar:
//   - Veri Çıkarma: Resource'dan veri çıkarma ve element'e atama
//   - Serileştirme: Element verilerini JSON formatına dönüştürme
//   - Görünürlük: Element'in hangi context'lerde görüneceğini belirleme
//   - Yapılandırma: Fluent API ile element özelliklerini ayarlama
//   - Validasyon: Veri doğrulama kurallarını tanımlama ve uygulama
//   - Callback'ler: Özel davranışlar için callback fonksiyonları
//
// # Kullanım Senaryoları
//
// Element interface şu durumlarda kullanılır:
//   - Yeni alan tipleri oluşturma (interface'i implement ederek)
//   - Alan davranışlarını özelleştirme (callback'ler ile)
//   - Dinamik görünürlük kontrolü (CanSee ile)
//   - Veri transformasyonu (Resolve, Modify ile)
//   - Dosya yükleme işlemleri (StoreAs ile)
//
// # Method Chaining (Zincirleme Metodlar)
//
// Element interface, fluent API tasarımı kullanır. Metodlar Element döndürerek
// zincirleme çağrılara izin verir:
//
//	field := fields.Text("Name", "name").
//	    OnList().
//	    OnForm().
//	    Required().
//	    Searchable().
//	    Sortable()
//
// # Önemli Notlar
//
// - Tüm alan tipleri bu interface'i implement etmelidir
// - Fluent metodlar Element döndürerek zincirleme çağrılara izin verir
// - Callback fonksiyonları, dinamik davranış kontrolü sağlar
// - Görünürlük kontrolü, context bazlıdır (list, detail, form)
// - Validasyon kuralları, hem frontend hem backend'de uygulanır
//
// # İlgili Dokümantasyon
//
// Detaylı bilgi için bakınız:
//   - docs/Fields.md: Alan sistemi ve tüm alan tipleri
//   - docs/Relationships.md: İlişki alanları (BelongsTo, HasMany, vb.)
//
// # Örnek Implementasyon
//
//	type CustomField struct {
//	    *fields.Schema // Base implementation
//	}
//
//	func (f *CustomField) Extract(resource any) {
//	    // Özel veri çıkarma mantığı
//	    f.Schema.Extract(resource)
//	}
//
//	func (f *CustomField) JsonSerialize() map[string]any {
//	    // Özel serileştirme mantığı
//	    data := f.Schema.JsonSerialize()
//	    data["customProp"] = "customValue"
//	    return data
//	}
type Element interface {
	// ============================================================================
	// Temel Erişim Metodları (Basic Accessors)
	// ============================================================================
	//
	// Bu metodlar, element'in temel özelliklerine erişim sağlar.

	// GetKey, bu element için benzersiz tanımlayıcıyı döndürür.
	//
	// Key, element'i resource model'deki bir alana eşlemek için kullanılır.
	// Genellikle veritabanı sütun adı veya struct field adı ile eşleşir.
	//
	// Döndürür:
	//   - Element'in benzersiz key değeri (örn: "name", "email", "author_id")
	//
	// Örnek:
	//   field := fields.Text("Name", "name")
	//   key := field.GetKey() // "name"
	GetKey() string

	// GetView, bu element'in görünüm tipini döndürür.
	//
	// View tipi, element'in UI'da nasıl render edileceğini belirler.
	// Örneğin: "text", "number", "select", "belongsTo", "hasMany"
	//
	// Döndürür:
	//   - Element'in view tipi (örn: "text", "select", "belongsTo")
	//
	// Örnek:
	//   field := fields.Text("Name", "name")
	//   view := field.GetView() // "text"
	GetView() string

	// GetContext, bu element'in görüntülendiği context'i döndürür.
	//
	// Context, element'in nerede gösterileceğini belirtir:
	//   - List: Liste görünümünde
	//   - Detail: Detay görünümünde
	//   - Form: Form görünümünde (create/update)
	//
	// Döndürür:
	//   - Element'in görüntülenme context'i (ElementContext tipi)
	//
	// Örnek:
	//   field := fields.Text("Name", "name").OnList().OnForm()
	//   ctx := field.GetContext() // List ve Form context'lerini içerir
	GetContext() ElementContext

	// GetName, bu element'in görünen adını döndürür.
	//
	// Name, element'in kullanıcı arayüzünde gösterilecek insan okunabilir adıdır.
	// Genellikle form label'ı veya tablo başlığı olarak kullanılır.
	//
	// Döndürür:
	//   - Element'in görünen adı (örn: "Name", "Email", "Author")
	//
	// Örnek:
	//   field := fields.Text("Name", "name")
	//   name := field.GetName() // "Name"
	GetName() string

	// GetType, bu element'in veri tipini döndürür.
	//
	// Type, element'in veri tipini belirler ve validation, formatting ve
	// OpenAPI schema oluşturma için kullanılır.
	// Örneğin: TYPE_TEXT, TYPE_NUMBER, TYPE_BOOLEAN, TYPE_DATE
	//
	// Döndürür:
	//   - Element'in veri tipi (core.ElementType)
	//
	// Örnek:
	//   field := fields.Text("Name", "name")
	//   fieldType := field.GetType() // core.TYPE_TEXT
	GetType() ElementType

	// Rows, textarea field'ı için satır sayısını ayarlar.
	//
	// Bu metod, textarea field'larında görüntülenecek satır sayısını belirler.
	// Sadece textarea field'ları için geçerlidir, diğer field'lar için no-op olarak çalışır.
	//
	// Parametreler:
	//   - rows: Satır sayısı
	//
	// Döndürür:
	//   - Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Textarea("Description", "description").Rows(5)
	Rows(rows int) Element

	// ============================================================================
	// Veri İşleme Metodları (Data Processing)
	// ============================================================================
	//
	// Bu metodlar, veri çıkarma ve serileştirme işlemlerini yönetir.

	// Extract, verilen resource'dan veri çıkarır ve element'i doldurur.
	//
	// Bu metod, resource'dan (genellikle bir struct veya map) veriyi okur ve
	// element'in internal state'ine atar. Veri çıkarma işlemi, element'in
	// key'i kullanılarak yapılır.
	//
	// Parametreler:
	//   - resource: Veri çıkarılacak kaynak (struct, map, vb.)
	//
	// Örnek:
	//   type User struct {
	//       Name string `json:"name"`
	//   }
	//   user := User{Name: "John"}
	//   field := fields.Text("Name", "name")
	//   field.Extract(user) // field artık "John" değerini içerir
	Extract(resource any)

	// JsonSerialize, element'i JSON uyumlu bir map'e serileştirir.
	//
	// Bu metod, element'in tüm özelliklerini (name, key, view, props, vb.)
	// içeren bir map döndürür. Bu map, JSON encoding için hazırdır ve
	// frontend'e gönderilir.
	//
	// Döndürür:
	//   - JSON encoding için hazır map[string]any
	//
	// Örnek:
	//   field := fields.Text("Name", "name").Required()
	//   data := field.JsonSerialize()
	//   // data = {
	//   //   "name": "Name",
	//   //   "key": "name",
	//   //   "view": "text",
	//   //   "props": {"required": true},
	//   //   ...
	//   // }
	JsonSerialize() map[string]any

	// ============================================================================
	// Görünürlük Kontrolü (Visibility Control)
	// ============================================================================
	//
	// Bu metodlar, element'in görünürlüğünü kontrol eder.

	// IsVisible, element'in verilen context'te görünür olup olmadığını belirler.
	//
	// Bu metod, element'in görünürlük kurallarını (OnList, HideOnCreate, vb.)
	// ve CanSee callback'ini değerlendirerek görünürlük kararı verir.
	//
	// Parametreler:
	//   - ctx: Resource context (kullanıcı, item, vb. bilgilerini içerir)
	//
	// Döndürür:
	//   - true: Element görünür olmalı
	//   - false: Element gizli olmalı
	//
	// Örnek:
	//   field := fields.Text("Admin Note", "admin_note").
	//       CanSee(func(ctx *core.ResourceContext) bool {
	//           return ctx.User.Role == "admin"
	//       })
	//   visible := field.IsVisible(ctx) // Admin ise true
	IsVisible(ctx *ResourceContext) bool

	// IsSearchable, bu element'in aranabilir olup olmadığını döndürür.
	//
	// Aranabilir element'ler, global arama işlemlerinde kullanılır.
	// Searchable() metodu ile işaretlenen element'ler true döner.
	//
	// Döndürür:
	//   - true: Element aranabilir
	//   - false: Element aranabilir değil
	//
	// Örnek:
	//   field := fields.Text("Name", "name").Searchable()
	//   searchable := field.IsSearchable() // true
	IsSearchable() bool

	// ============================================================================
	// Fluent Setter'lar - Görünüm Kontrolü (View Control)
	// ============================================================================
	//
	// Bu metodlar, element'in hangi görünümlerde gösterileceğini kontrol eder.
	// Tüm metodlar Element döndürerek method chaining'e izin verir.

	// SetName, element'in görüntüleme adını ayarlar.
	//
	// Bu ad, UI'da alan etiketi (label) olarak gösterilir. Kullanıcı dostu,
	// açıklayıcı bir ad seçilmelidir.
	//
	// Parametreler:
	//   - name: Element'in görüntüleme adı (örn: "Kullanıcı Adı", "E-posta")
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("", "username").SetName("Kullanıcı Adı")
	SetName(name string) Element

	// SetKey, element'in benzersiz tanımlayıcısını ayarlar.
	//
	// Key, element'i resource'daki bir alana eşlemek için kullanılır.
	// Genellikle veritabanı sütun adı veya struct field adı ile eşleşir.
	//
	// Parametreler:
	//   - key: Element'in benzersiz key'i (örn: "username", "email")
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Name", "").SetKey("full_name")
	SetKey(key string) Element

	// OnList, element'i liste görünümünde görünür yapar.
	//
	// Bu metod eklemeli (additive) çalışır - element'i diğer görünümlerden gizlemez.
	// Birden fazla On* metodu zincirlenerek element birden fazla görünümde gösterilebilir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Name", "name").OnList().OnDetail()
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Görünürlük kontrolleri bölümü
	OnList() Element

	// OnDetail, element'i detay görünümünde görünür yapar.
	//
	// Bu metod eklemeli (additive) çalışır - element'i diğer görünümlerden gizlemez.
	// Detay görünümü, bir kaydın tüm bilgilerini gösterir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Description", "description").OnDetail().OnForm()
	OnDetail() Element

	// OnForm, element'i form görünümlerinde (create ve update) görünür yapar.
	//
	// Bu metod eklemeli (additive) çalışır - element'i diğer görünümlerden gizlemez.
	// Form görünümü, kayıt oluşturma ve güncelleme formlarını içerir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Name", "name").OnForm().OnList()
	OnForm() Element

	// HideOnList, element'i liste görünümünde gizler.
	//
	// Element diğer görünümlerde (detail, form) görünür olmaya devam eder.
	// Bu metod, liste görünümünü sadeleştirmek için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Textarea("Description", "description").
	//       OnDetail().OnForm().HideOnList()
	HideOnList() Element

	// HideOnDetail, element'i detay görünümünde gizler.
	//
	// Element diğer görünümlerde (list, form) görünür olmaya devam eder.
	// Bu metod, detay görünümünü sadeleştirmek için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Password("Password", "password").
	//       OnForm().HideOnDetail()
	HideOnDetail() Element

	// HideOnCreate, element'i oluşturma formunda gizler.
	//
	// Element güncelleme formunda ve diğer görünümlerde görünür olmaya devam eder.
	// Bu metod, sadece güncelleme sırasında düzenlenebilir alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Slug", "slug").
	//       OnForm().HideOnCreate() // Slug sadece update'te düzenlenebilir
	HideOnCreate() Element

	// HideOnUpdate, element'i güncelleme formunda gizler.
	//
	// Element oluşturma formunda ve diğer görünümlerde görünür olmaya devam eder.
	// Bu metod, sadece oluşturma sırasında ayarlanabilen alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Password("Password", "password").
	//       OnForm().HideOnUpdate() // Şifre sadece create'te ayarlanabilir
	HideOnUpdate() Element

	// OnlyOnList, element'i sadece liste görünümünde gösterir.
	//
	// Bu metod, element'i diğer tüm görünümlerden (detail, form) gizler.
	// Liste görünümüne özel alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Status Badge", "status").OnlyOnList()
	OnlyOnList() Element

	// OnlyOnDetail, element'i sadece detay görünümünde gösterir.
	//
	// Bu metod, element'i diğer tüm görünümlerden (list, form) gizler.
	// Detay görünümüne özel alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.DateTime("Created At", "created_at").
	//       OnlyOnDetail().ReadOnly()
	OnlyOnDetail() Element

	// OnlyOnCreate, element'i sadece oluşturma formunda gösterir.
	//
	// Bu metod, element'i diğer tüm görünümlerden (list, detail, update) gizler.
	// Sadece oluşturma sırasında gerekli alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Password("Password", "password").
	//       OnlyOnCreate().Required()
	OnlyOnCreate() Element

	// OnlyOnUpdate, element'i sadece güncelleme formunda gösterir.
	//
	// Bu metod, element'i diğer tüm görünümlerden (list, detail, create) gizler.
	// Sadece güncelleme sırasında düzenlenebilir alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Password("New Password", "password").
	//       OnlyOnUpdate().Nullable()
	OnlyOnUpdate() Element

	// OnlyOnForm, element'i sadece form görünümlerinde (create ve update) gösterir.
	//
	// Bu metod, element'i liste ve detay görünümlerinden gizler.
	// Sadece form'da düzenlenebilir alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Name", "name").OnlyOnForm().Required()
	OnlyOnForm() Element

	// ============================================================================
	// Fluent Setter'lar - Özellikler (Properties)
	// ============================================================================
	//
	// Bu metodlar, element'in davranış ve görünüm özelliklerini ayarlar.
	// Tüm metodlar Element döndürerek method chaining'e izin verir.

	// ReadOnly, element'i salt okunur olarak işaretler.
	//
	// Salt okunur element'ler görüntülenir ancak kullanıcılar tarafından
	// değiştirilemez. Genellikle sistem tarafından otomatik oluşturulan
	// veya hesaplanan alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Otomatik oluşturulan ID'ler
	//   - Sistem tarafından ayarlanan tarihler (created_at, updated_at)
	//   - Hesaplanan alanlar (toplam, ortalama, vb.)
	//
	// Örnek:
	//   field := fields.DateTime("Created At", "created_at").
	//       OnList().OnDetail().ReadOnly()
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Alan seçenekleri bölümü
	ReadOnly() Element

	// WithProps, element'e özel özellikler ekler.
	//
	// Özel özellikler, frontend'e ek veri göndermek için kullanılır.
	// Bu özellikler, frontend bileşenlerinin davranışını özelleştirmek
	// için kullanılabilir.
	//
	// Parametreler:
	//   - key: Özellik anahtarı
	//   - value: Özellik değeri (any tipi)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Color", "color").
	//       WithProps("type", "color").
	//       WithProps("showAlpha", true)
	WithProps(key string, value any) Element

	// Disabled, element'i devre dışı olarak işaretler.
	//
	// Devre dışı element'ler görünür ancak etkileşimli değildir.
	// Kullanıcılar bu element'lerle etkileşime geçemez.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Username", "username").
	//       OnForm().Disabled()
	Disabled() Element

	// Immutable, element'i değiştirilemez olarak işaretler.
	//
	// Değiştirilemez element'ler oluşturma sırasında ayarlanabilir
	// ancak güncelleme sırasında değiştirilemez. Bu, bir kez ayarlandıktan
	// sonra değişmemesi gereken alanlar için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Kullanıcı adları (username)
	//   - E-posta adresleri (bazı durumlarda)
	//   - Benzersiz tanımlayıcılar (slug, code)
	//
	// Örnek:
	//   field := fields.Text("Username", "username").
	//       OnForm().Required().Immutable()
	Immutable() Element

	// Required, element'i zorunlu olarak işaretler.
	//
	// Zorunlu element'ler, form gönderilebilmesi için bir değere sahip
	// olmalıdır. Frontend ve backend'de validasyon uygulanır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Name", "name").
	//       OnForm().Required()
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Validasyon bölümü
	Required() Element

	// Nullable, element'i nullable (boş bırakılabilir) olarak işaretler.
	//
	// Nullable element'ler null/nil değerlere sahip olabilir.
	// Veritabanında NULL değer olarak saklanabilir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Middle Name", "middle_name").
	//       OnForm().Nullable()
	Nullable() Element

	// Min, element için minimum sayısal değeri ayarlar.
	//
	// Bu metod, sayısal alanlar için minimum değeri belirler.
	// Frontend'de HTML 'min' özelliği olarak kullanılır.
	//
	// Parametreler:
	//   - value: Minimum değer (int, float, string)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Number("Age", "age").Min(18)
	Min(value any) Element

	// Max, element için maksimum sayısal değeri ayarlar.
	//
	// Bu metod, sayısal alanlar için maksimum değeri belirler.
	// Frontend'de HTML 'max' özelliği olarak kullanılır.
	//
	// Parametreler:
	//   - value: Maksimum değer (int, float, string)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Number("Quantity", "quantity").Max(100)
	Max(value any) Element

	// MinLength, element için minimum karakter uzunluğunu ayarlar.
	//
	// Bu metod, string alanlar için minimum karakter sayısını belirler.
	// Frontend'de HTML 'minlength' özelliği olarak kullanılır.
	//
	// Parametreler:
	//   - length: Minimum karakter sayısı
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Password("Password", "password").MinLength(8)
	MinLength(length int) Element

	// MaxLength, element için maksimum karakter uzunluğunu ayarlar.
	//
	// Bu metod, string alanlar için maksimum karakter sayısını belirler.
	// Frontend'de HTML 'maxlength' özelliği olarak kullanılır.
	//
	// Parametreler:
	//   - length: Maksimum karakter sayısı
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Title", "title").MaxLength(100)
	MaxLength(length int) Element

	// Placeholder, element için yer tutucu metni ayarlar.
	//
	// Yer tutucu metin, element boş olduğunda gösterilir ve
	// kullanıcıya ne girilmesi gerektiği hakkında ipucu verir.
	//
	// Parametreler:
	//   - placeholder: Yer tutucu metin
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Email", "email").
	//       OnForm().Placeholder("ornek@email.com")
	Placeholder(placeholder string) Element

	// Label, element için etiket metnini ayarlar.
	//
	// Etiket, element'in üstünde veya yanında görüntülenir.
	// SetName ile aynı işlevi görür.
	//
	// Parametreler:
	//   - label: Etiket metni
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("", "email").Label("E-posta Adresi")
	Label(label string) Element

	// HelpText, element için yardım metni ayarlar.
	//
	// Yardım metni, kullanıcılara ek rehberlik sağlar ve
	// element'in altında veya yanında gösterilir.
	//
	// Parametreler:
	//   - helpText: Yardım metni
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Password("Password", "password").
	//       OnForm().Required().
	//       HelpText("En az 8 karakter, bir büyük harf ve bir rakam içermelidir")
	HelpText(helpText string) Element

	// Filterable, element'i filtrelenebilir olarak işaretler.
	//
	// Filtrelenebilir element'ler, liste görünümlerini filtrelemek
	// için kullanılabilir. Kullanıcılar bu alanlara göre kayıtları
	// filtreleyebilir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Select("Status", "status").
	//       OnList().Filterable()
	Filterable() Element

	// Sortable, element'i sıralanabilir olarak işaretler.
	//
	// Sıralanabilir element'ler, liste görünümlerini sıralamak
	// için kullanılabilir. Kullanıcılar bu alanlara göre kayıtları
	// artan veya azalan sırada sıralayabilir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Önemli Not:
	//   - Sıralanabilir alanlar için veritabanı indeksi oluşturulmalıdır
	//
	// Örnek:
	//   field := fields.Text("Name", "name").
	//       OnList().Sortable().Searchable()
	Sortable() Element

	// Searchable, element'i aranabilir olarak işaretler.
	//
	// Aranabilir element'ler, global arama işlemlerinde kullanılır.
	// Kullanıcılar bu alanlarda arama yapabilir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Önemli Not:
	//   - Aranabilir alanlar için veritabanı indeksi oluşturulmalıdır
	//   - Fulltext search için özel indeks kullanılabilir
	//
	// Örnek:
	//   field := fields.Text("Name", "name").
	//       OnList().Searchable().Sortable()
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Arama ve sıralama bölümü
	Searchable() Element

	// Stacked, element'i yığılmış (tam genişlik) olarak işaretler.
	//
	// Yığılmış element'ler, container'ın tam genişliğini kaplar.
	// Genellikle uzun metin alanları veya özel bileşenler için kullanılır.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Textarea("Description", "description").
	//       OnForm().Stacked()
	Stacked() Element

	// SetTextAlign, element için metin hizalamasını ayarlar.
	//
	// Metin hizalama, liste görünümünde element'in içeriğinin
	// nasıl hizalanacağını belirler.
	//
	// Parametreler:
	//   - align: Hizalama değeri ("left", "center", "right")
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Sayılar için sağa hizalama ("right")
	//   - Boolean değerler için ortaya hizalama ("center")
	//   - Metin için sola hizalama ("left" - varsayılan)
	//
	// Örnek:
	//   field := fields.Number("Price", "price").
	//       OnList().SetTextAlign("right")
	SetTextAlign(align string) Element

	// ============================================================================
	// Callback Fonksiyonları (Callbacks)
	// ============================================================================
	//
	// Bu metodlar, element'in davranışını özelleştirmek için callback fonksiyonları
	// tanımlar ve yönetir. Callback'ler, dinamik görünürlük, veri transformasyonu
	// ve özel depolama mantığı için kullanılır.

	// CanSee, görünürlük callback fonksiyonunu ayarlar.
	//
	// Bu callback, element'in verilen resource context'te görünür olup olmadığını
	// belirler. Kullanıcı rolü, kayıt durumu veya diğer koşullara göre dinamik
	// görünürlük kontrolü sağlar.
	//
	// Parametreler:
	//   - fn: Görünürlük callback fonksiyonu (ResourceContext alır, bool döner)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Rol bazlı görünürlük kontrolü
	//   - Kayıt durumuna göre alan gösterme/gizleme
	//   - İzin bazlı erişim kontrolü
	//   - Koşullu alan görüntüleme
	//
	// Örnek:
	//   // Admin kullanıcılar için görünür
	//   field := fields.Text("Admin Note", "admin_note").
	//       OnForm().
	//       CanSee(func(ctx *core.ResourceContext) bool {
	//           return ctx.User.Role == "admin"
	//       })
	//
	//   // Yayınlanmış kayıtlar için görünür
	//   field := fields.DateTime("Published At", "published_at").
	//       OnDetail().
	//       CanSee(func(ctx *core.ResourceContext) bool {
	//           if post, ok := ctx.Item.(*Post); ok {
	//               return post.Status == "published"
	//           }
	//           return false
	//       })
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Görünürlük kontrolü bölümü
	CanSee(fn VisibilityFunc) Element

	// StoreAs, dosya yüklemeleri için özel depolama callback fonksiyonunu ayarlar.
	//
	// Bu callback, dosya yükleme işlemlerinde özel depolama mantığı uygulamak
	// için kullanılır. Dosyanın nereye ve nasıl kaydedileceğini kontrol eder.
	//
	// Parametreler:
	//   - fn: Depolama callback fonksiyonu (ResourceContext alır, dosya yolu döner)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Özel dosya adlandırma
	//   - Farklı depolama sistemleri (S3, local, vb.)
	//   - Dosya işleme (resize, compress, vb.)
	//   - Metadata ekleme
	//
	// Örnek:
	//   field := fields.Image("Avatar", "avatar").
	//       OnForm().
	//       StoreAs(func(ctx *core.ResourceContext) (interface{}, error) {
	//           // Özel dosya adı oluştur
	//           filename := fmt.Sprintf("user_%d_%d.jpg", ctx.User.ID, time.Now().Unix())
	//           // S3'e yükle
	//           path, err := uploadToS3(ctx.File, filename)
	//           return path, err
	//       })
	StoreAs(fn StorageCallbackFunc) Element

	// GetStorageCallback, depolama callback fonksiyonunu döndürür.
	//
	// Bu metod, element'e atanmış depolama callback'ini almak için kullanılır.
	// Callback atanmamışsa nil döner.
	//
	// Döndürür:
	//   - Depolama callback fonksiyonu veya nil
	//
	// Örnek:
	//   callback := field.GetStorageCallback()
	//   if callback != nil {
	//       // Callback kullan
	//   }
	GetStorageCallback() StorageCallbackFunc

	// Resolve, değeri görüntülemeden önce dönüştürmek için callback ayarlar.
	//
	// Bu callback, element'in değerini kullanıcıya gösterilmeden önce
	// dönüştürmek için kullanılır. Veri formatı değiştirme, hesaplama
	// veya ilişkili veri yükleme için kullanılabilir.
	//
	// Parametreler:
	//   - fn: Resolve callback fonksiyonu (value, item, context alır, any döner)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Hesaplanan alanlar (full_name = first_name + last_name)
	//   - Veri formatı dönüşümü (timestamp -> formatted date)
	//   - İlişkili veri yükleme
	//   - Özel görüntüleme mantığı
	//
	// Örnek:
	//   // Tam ad hesaplama
	//   field := fields.Text("Full Name", "full_name").
	//       OnList().
	//       Resolve(func(value any, item any, c *fiber.Ctx) any {
	//           if user, ok := item.(*User); ok {
	//               return user.FirstName + " " + user.LastName
	//           }
	//           return value
	//       })
	//
	//   // Fiyat formatı
	//   field := fields.Number("Price", "price").
	//       OnList().
	//       Resolve(func(value any, item any, c *fiber.Ctx) any {
	//           if price, ok := value.(float64); ok {
	//               return fmt.Sprintf("₺%.2f", price)
	//           }
	//           return value
	//       })
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Callback'ler bölümü
	Resolve(fn func(value any, item any, c *fiber.Ctx) any) Element

	// GetResolveCallback, resolve callback fonksiyonunu döndürür.
	//
	// Bu metod, element'e atanmış resolve callback'ini almak için kullanılır.
	// Callback atanmamışsa nil döner.
	//
	// Döndürür:
	//   - Resolve callback fonksiyonu veya nil
	//
	// Örnek:
	//   callback := field.GetResolveCallback()
	//   if callback != nil {
	//       resolvedValue := callback(value, item, ctx)
	//   }
	GetResolveCallback() func(value any, item any, c *fiber.Ctx) any

	// Modify, değeri kaydetmeden önce dönüştürmek için callback ayarlar.
	//
	// Bu callback, element'in değerini veritabanına kaydedilmeden önce
	// dönüştürmek için kullanılır. Veri temizleme, şifreleme, hashing
	// veya format dönüşümü için kullanılabilir.
	//
	// Parametreler:
	//   - fn: Modify callback fonksiyonu (value, context alır, any döner)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Şifre hashleme
	//   - Veri temizleme (trim, lowercase, vb.)
	//   - Şifreleme
	//   - Format dönüşümü
	//
	// Önemli Not:
	//   - Modify callback'i sadece form submit sırasında çalışır
	//   - Resolve callback'inden farklı olarak, sadece kaydetme işleminde kullanılır
	//
	// Örnek:
	//   // Şifre hashleme
	//   field := fields.Password("Password", "password").
	//       OnForm().
	//       Modify(func(value any, c *fiber.Ctx) any {
	//           if password, ok := value.(string); ok {
	//               hashed, _ := bcrypt.GenerateFromPassword(
	//                   []byte(password),
	//                   bcrypt.DefaultCost,
	//               )
	//               return string(hashed)
	//           }
	//           return value
	//       })
	//
	//   // E-posta temizleme
	//   field := fields.Email("Email", "email").
	//       OnForm().
	//       Modify(func(value any, c *fiber.Ctx) any {
	//           if email, ok := value.(string); ok {
	//               return strings.ToLower(strings.TrimSpace(email))
	//           }
	//           return value
	//       })
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Callback'ler bölümü
	Modify(fn func(value any, c *fiber.Ctx) any) Element

	// GetModifyCallback, modify callback fonksiyonunu döndürür.
	//
	// Bu metod, element'e atanmış modify callback'ini almak için kullanılır.
	// Callback atanmamışsa nil döner.
	//
	// Döndürür:
	//   - Modify callback fonksiyonu veya nil
	//
	// Örnek:
	//   callback := field.GetModifyCallback()
	//   if callback != nil {
	//       modifiedValue := callback(value, ctx)
	//   }
	GetModifyCallback() func(value any, c *fiber.Ctx) any

	// ============================================================================
	// Diğer Metodlar (Other)
	// ============================================================================
	//
	// Bu metodlar, element'in seçeneklerini, varsayılan değerlerini ve
	// otomatik yapılandırmasını yönetir.

	// Options, seçim tipi element'ler için mevcut seçenekleri ayarlar.
	//
	// Bu metod, Select, Combobox, Radio gibi seçim element'lerinde
	// kullanıcıya sunulacak seçenekleri tanımlar. Seçenekler, bir map
	// (key-value çiftleri) veya slice olarak verilebilir.
	//
	// Parametreler:
	//   - options: Seçenekler (map[string]string veya []string)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Statik seçenek listeleri (durum, kategori, vb.)
	//   - Sabit değer setleri
	//   - Enum değerleri
	//
	// Önemli Not:
	//   - İlişki alanlarında (BelongsTo, HasOne, vb.) AutoOptions kullanımı önerilir
	//   - Options ve AutoOptions birlikte kullanılamaz
	//
	// Örnek:
	//   // Map ile seçenekler
	//   field := fields.Select("Status", "status").
	//       OnForm().
	//       Options(map[string]string{
	//           "draft":     "Taslak",
	//           "published": "Yayınlandı",
	//           "archived":  "Arşivlendi",
	//       })
	//
	//   // Slice ile seçenekler
	//   field := fields.Select("Priority", "priority").
	//       OnForm().
	//       Options([]string{"Low", "Medium", "High"})
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Seçim alanları bölümü
	Options(options any) Element

	// GetAutoOptionsConfig, AutoOptions yapılandırmasını döndürür.
	//
	// Bu metod, element'in AutoOptions özelliğinin aktif olup olmadığını
	// ve hangi display field'ın kullanıldığını kontrol etmek için kullanılır.
	//
	// Döndürür:
	//   - AutoOptionsConfig yapısı (Enabled ve DisplayField içerir)
	//
	// Kullanım Senaryoları:
	//   - AutoOptions'ın aktif olup olmadığını kontrol etme
	//   - Display field'ı alma
	//   - Backend'de otomatik seçenek üretimi için yapılandırma okuma
	//
	// Örnek:
	//   config := field.GetAutoOptionsConfig()
	//   if config.Enabled {
	//       // Otomatik seçenek üretimi aktif
	//       displayField := config.DisplayField // "name", "title", vb.
	//   }
	//
	// İlgili Dokümantasyon:
	//   - docs/Relationships.md: AutoOptions kullanımı
	GetAutoOptionsConfig() AutoOptionsConfig

	// Default, element için varsayılan değeri ayarlar.
	//
	// Varsayılan değer, yeni kayıt oluştururken element'in başlangıç
	// değeri olarak kullanılır. Kullanıcı değeri değiştirmezse, bu
	// değer kaydedilir.
	//
	// Parametreler:
	//   - value: Varsayılan değer (any tipi)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Boolean alanlar için varsayılan durum (true/false)
	//   - Durum alanları için başlangıç durumu ("draft", "active", vb.)
	//   - Sayısal alanlar için başlangıç değeri (0, 1, vb.)
	//   - Tarih alanları için mevcut tarih
	//
	// Örnek:
	//   // Boolean varsayılan değer
	//   field := fields.Switch("Active", "is_active").
	//       OnForm().Default(true)
	//
	//   // String varsayılan değer
	//   field := fields.Select("Status", "status").
	//       OnForm().Default("draft")
	//
	//   // Sayısal varsayılan değer
	//   field := fields.Number("Quantity", "quantity").
	//       OnForm().Default(1)
	Default(value any) Element

	// ============================================================================
	// Genişletilmiş Alan Sistemi Metodları (Extended Field System Methods)
	// ============================================================================
	//
	// Bu metodlar, gelişmiş alan davranışları için kullanılır. Dinamik görünürlük,
	// bağımlılıklar, meta veriler ve özel görüntüleme mantığı sağlar.

	// IsHidden, element'in verilen görünürlük context'inde gizli olup olmadığını belirler.
	//
	// Bu metod, element'in belirli bir context'te (list, detail, create, update)
	// gizli olup olmadığını kontrol eder. HideOnList, HideOnCreate gibi metodlarla
	// ayarlanan görünürlük kurallarını değerlendirir.
	//
	// Parametreler:
	//   - ctx: Görünürlük context'i (list, detail, create, update)
	//
	// Döndürür:
	//   - true: Element bu context'te gizli
	//   - false: Element bu context'te görünür
	//
	// Örnek:
	//   field := fields.Password("Password", "password").HideOnUpdate()
	//   hidden := field.IsHidden(VisibilityContext{Context: "update"}) // true
	IsHidden(ctx VisibilityContext) bool

	// ResolveForDisplay, element'in değerini görüntüleme amaçlı çözümler.
	//
	// Bu metod, element'in değerini kullanıcıya gösterilmeden önce dönüştürür.
	// Resolve callback'inden farklı olarak, özellikle görüntüleme formatlaması
	// için tasarlanmıştır.
	//
	// Parametreler:
	//   - item: Değeri çözümlenecek kayıt
	//
	// Döndürür:
	//   - Çözümlenmiş değer (any tipi)
	//   - Hata (varsa)
	//
	// Kullanım Senaryoları:
	//   - Tarih formatı dönüşümü
	//   - Para birimi formatı
	//   - Hesaplanan alanlar
	//   - İlişkili veri gösterimi
	//
	// Örnek:
	//   value, err := field.ResolveForDisplay(user)
	//   if err != nil {
	//       // Hata yönetimi
	//   }
	ResolveForDisplay(item any) (any, error)

	// GetDependencies, bu element'in bağımlı olduğu alan adlarının listesini döndürür.
	//
	// Bağımlılıklar, alan görünürlüğünü ve validasyon sırasını belirlemek için
	// kullanılır. Bir alan, başka alanların değerlerine göre görünür veya
	// gizli olabilir.
	//
	// Döndürür:
	//   - Bağımlı olunan alan adları listesi ([]string)
	//
	// Kullanım Senaryoları:
	//   - Koşullu alan görünürlüğü
	//   - Validasyon sırası belirleme
	//   - Form bağımlılık grafiği oluşturma
	//
	// Örnek:
	//   field := fields.Text("Tax Number", "tax_number").
	//       DependsOn("is_company")
	//   deps := field.GetDependencies() // ["is_company"]
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Bağımlılıklar bölümü
	GetDependencies() []string

	// IsConditionallyVisible, element'in item'ın değerlerine göre görünür olup olmadığını belirler.
	//
	// Bu metod, element'in görünürlüğünün diğer alanların değerlerine bağlı
	// olup olmadığını kontrol eder. Dinamik form davranışı için kullanılır.
	//
	// Parametreler:
	//   - item: Değerleri kontrol edilecek kayıt
	//
	// Döndürür:
	//   - true: Element görünür olmalı
	//   - false: Element gizli olmalı
	//
	// Kullanım Senaryoları:
	//   - Koşullu alan gösterimi
	//   - Dinamik form yapısı
	//   - Bağımlı alan görünürlüğü
	//
	// Örnek:
	//   field := fields.Text("Tax Number", "tax_number").
	//       DependsOn("is_company").
	//       When("is_company", "=", true)
	//   visible := field.IsConditionallyVisible(company) // true veya false
	IsConditionallyVisible(item any) bool

	// GetMetadata, element hakkında meta verileri döndürür.
	//
	// Meta veriler, alan ilişkileri, kısıtlamalar ve özel özellikler hakkında
	// bilgi içerir. Bu bilgiler, frontend ve backend tarafından kullanılabilir.
	//
	// Döndürür:
	//   - Meta veri map'i (map[string]any)
	//
	// Meta Veri İçeriği:
	//   - Alan tipi bilgileri
	//   - İlişki yapılandırması
	//   - Validasyon kuralları
	//   - Özel özellikler
	//
	// Örnek:
	//   metadata := field.GetMetadata()
	//   // metadata = {
	//   //   "type": "belongsTo",
	//   //   "relatedResource": "users",
	//   //   "foreignKey": "author_id",
	//   //   ...
	//   // }
	GetMetadata() map[string]any

	// ============================================================================
	// Validasyon Metodları (Validation Methods)
	// ============================================================================
	//
	// Bu metodlar, element'in validasyon kurallarını yönetir ve değerleri doğrular.

	// GetValidationRules, bu element için tanımlanmış validasyon kurallarını döndürür.
	//
	// Validasyon kuralları, element'in değerinin geçerli olup olmadığını kontrol
	// etmek için kullanılır. Kurallar hem frontend hem backend'de uygulanır.
	//
	// Döndürür:
	//   - Validasyon kuralları listesi ([]interface{})
	//
	// Kural Tipleri:
	//   - Required: Zorunlu alan
	//   - MinLength/MaxLength: Minimum/maksimum uzunluk
	//   - Min/Max: Minimum/maksimum değer
	//   - Email: E-posta formatı
	//   - URL: URL formatı
	//   - Pattern: Regex deseni
	//   - Unique: Benzersizlik kontrolü
	//
	// Örnek:
	//   rules := field.GetValidationRules()
	//   // rules = [
	//   //   {type: "required"},
	//   //   {type: "minLength", value: 8},
	//   //   {type: "email"}
	//   // ]
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Validasyon bölümü
	GetValidationRules() []interface{}

	// AddValidationRule, bu element'e bir validasyon kuralı ekler.
	//
	// Bu metod, element'e dinamik olarak validasyon kuralı eklemek için
	// kullanılır. Mevcut kurallara ek olarak yeni kural tanımlanır.
	//
	// Parametreler:
	//   - rule: Eklenecek validasyon kuralı (interface{} tipi)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Örnek:
	//   field := fields.Text("Username", "username").
	//       AddValidationRule(map[string]interface{}{
	//           "type": "pattern",
	//           "value": "^[a-zA-Z0-9_]+$",
	//           "message": "Sadece harf, rakam ve alt çizgi kullanılabilir",
	//       })
	AddValidationRule(rule interface{}) Element

	// ValidateValue, bir değeri element'in validasyon kurallarına göre doğrular.
	//
	// Bu metod, verilen değerin element'in tüm validasyon kurallarını
	// karşılayıp karşılamadığını kontrol eder. Hata varsa detaylı hata
	// mesajı döner.
	//
	// Parametreler:
	//   - value: Doğrulanacak değer (interface{} tipi)
	//
	// Döndürür:
	//   - nil: Değer geçerli
	//   - error: Validasyon hatası (detaylı mesaj içerir)
	//
	// Kullanım Senaryoları:
	//   - Form submit öncesi validasyon
	//   - API endpoint'lerinde veri kontrolü
	//   - Özel validasyon mantığı
	//
	// Örnek:
	//   field := fields.Email("Email", "email").Required()
	//   err := field.ValidateValue("invalid-email")
	//   if err != nil {
	//       // Validasyon hatası: "Geçerli bir e-posta adresi giriniz"
	//   }
	ValidateValue(value interface{}) error

	// GetCustomValidators, bu element için özel validator fonksiyonlarını döndürür.
	//
	// Özel validator'lar, standart validasyon kurallarının ötesinde karmaşık
	// validasyon mantığı uygulamak için kullanılır.
	//
	// Döndürür:
	//   - Özel validator fonksiyonları listesi ([]interface{})
	//
	// Kullanım Senaryoları:
	//   - Veritabanı kontrolü gerektiren validasyonlar
	//   - Çoklu alan karşılaştırması
	//   - Karmaşık iş kuralları
	//   - Asenkron validasyonlar
	//
	// Örnek:
	//   validators := field.GetCustomValidators()
	//   for _, validator := range validators {
	//       if fn, ok := validator.(func(interface{}) error); ok {
	//           err := fn(value)
	//           if err != nil {
	//               // Validasyon hatası
	//           }
	//       }
	//   }
	GetCustomValidators() []interface{}

	// ============================================================================
	// Görüntüleme Metodları (Display Methods)
	// ============================================================================
	//
	// Bu metodlar, element'in değerlerinin nasıl görüntüleneceğini kontrol eder.

	// Display, element için özel görüntüleme callback'i ayarlar.
	//
	// Desteklenen callback imzaları:
	//   - func(value any) string
	//   - func(value any) any
	//   - func(value any, item any) string
	//   - func(value any, item any) any
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	Display(fn interface{}) Element

	// DisplayAs, görüntüleme format string'ini ayarlar.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	DisplayAs(format string) Element

	// DisplayUsingLabels, seçim alanlarında value yerine label göstermeyi etkinleştirir.
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	DisplayUsingLabels() Element

	// GetDisplayCallback, görüntüleme callback fonksiyonunu döndürür.
	//
	// Bu callback, element'in değerinin nasıl görüntüleneceğini belirler.
	// Callback, alan değeri ve ilgili kayıt modelini alır.
	//
	// Döndürür:
	//   - Görüntüleme callback fonksiyonu (value, item alır; any döner) veya nil
	//
	// Kullanım Senaryoları:
	//   - Özel format uygulama (para birimi, tarih, vb.)
	//   - Değer dönüşümü (boolean -> "Evet"/"Hayır")
	//   - Hesaplanan değerler
	//
	// Örnek:
	//   callback := field.GetDisplayCallback()
	//   if callback != nil {
	//       displayValue := callback(value, item) // "₺1,234.56"
	//   }
	GetDisplayCallback() func(value any, item any) any

	// GetDisplayedAs, görüntüleme format string'ini döndürür.
	//
	// Format string, değerin nasıl gösterileceğini belirten bir şablondur.
	// Printf-style formatlamayı destekler.
	//
	// Döndürür:
	//   - Format string (örn: "Durum: %s", "Fiyat: ₺%.2f")
	//
	// Örnek:
	//   format := field.GetDisplayedAs() // "Durum: %s"
	//   displayValue := fmt.Sprintf(format, value) // "Durum: Aktif"
	GetDisplayedAs() string

	// ShouldDisplayUsingLabels, değerlerin etiketler kullanılarak gösterilip gösterilmeyeceğini döndürür.
	//
	// Bu metod, özellikle Select ve Combobox gibi seçim alanlarında kullanılır.
	// true ise, değer yerine etiket (label) gösterilir.
	//
	// Döndürür:
	//   - true: Etiketleri kullan (örn: "1" yerine "Elektronik")
	//   - false: Ham değeri kullan
	//
	// Kullanım Senaryoları:
	//   - Select alanlarında kullanıcı dostu gösterim
	//   - Enum değerlerinin okunabilir hale getirilmesi
	//   - ID yerine isim gösterimi
	//
	// Örnek:
	//   field := fields.Select("Category", "category").
	//       Options(map[string]string{"1": "Elektronik", "2": "Giyim"}).
	//       DisplayUsingLabels()
	//   useLabels := field.ShouldDisplayUsingLabels() // true
	//   // Listede "1" yerine "Elektronik" gösterilir
	ShouldDisplayUsingLabels() bool

	// GetResolveHandle, client-side component etkileşimi için resolve handle'ını döndürür.
	//
	// Resolve handle, frontend bileşenlerinin element ile etkileşime geçmesi
	// için kullanılan benzersiz bir tanımlayıcıdır.
	//
	// Döndürür:
	//   - Resolve handle string'i
	//
	// Kullanım Senaryoları:
	//   - Frontend-backend senkronizasyonu
	//   - Dinamik bileşen yükleme
	//   - Client-side veri çözümleme
	//
	// Örnek:
	//   handle := field.GetResolveHandle() // "user_full_name_resolver"
	GetResolveHandle() string

	// ============================================================================
	// Bağımlılık Metodları (Dependency Methods)
	// ============================================================================
	//
	// Bu metodlar, alanlar arası bağımlılıkları yönetir ve koşullu görünürlük sağlar.

	// SetDependencies, bu element için alan bağımlılıklarını ayarlar.
	//
	// Bağımlılıklar, element'in görünürlüğünün veya davranışının diğer
	// alanların değerlerine bağlı olduğunu belirtir.
	//
	// Parametreler:
	//   - deps: Bağımlı olunan alan adları listesi ([]string)
	//
	// Döndürür:
	//   - Yapılandırılmış Element pointer'ı (method chaining için)
	//
	// Kullanım Senaryoları:
	//   - Koşullu alan gösterimi
	//   - Dinamik form yapısı
	//   - Bağımlı validasyon
	//   - Cascade değişiklikler
	//
	// Örnek:
	//   field := fields.Text("Tax Number", "tax_number").
	//       SetDependencies([]string{"is_company", "country"})
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Bağımlılıklar bölümü
	SetDependencies(deps []string) Element

	// GetDependencyRules, bu element için bağımlılık kurallarını döndürür.
	//
	// Bağımlılık kuralları, element'in hangi koşullarda görünür veya
	// gizli olacağını belirten kuralları içerir.
	//
	// Döndürür:
	//   - Bağımlılık kuralları map'i (map[string]interface{})
	//
	// Kural Formatı:
	//   {
	//     "field_name": {
	//       "operator": "=",
	//       "value": true
	//     }
	//   }
	//
	// Örnek:
	//   rules := field.GetDependencyRules()
	//   // rules = {
	//   //   "is_company": {"operator": "=", "value": true},
	//   //   "country": {"operator": "=", "value": "TR"}
	//   // }
	GetDependencyRules() map[string]interface{}

	// ResolveDependencies, bağımlılık kurallarını verilen context'e göre değerlendirir.
	//
	// Bu metod, element'in bağımlılık kurallarının karşılanıp karşılanmadığını
	// kontrol eder ve element'in görünür olup olmayacağını belirler.
	//
	// Parametreler:
	//   - context: Değerlendirme context'i (genellikle form verileri)
	//
	// Döndürür:
	//   - true: Bağımlılıklar karşılandı, element görünür olmalı
	//   - false: Bağımlılıklar karşılanmadı, element gizli olmalı
	//
	// Kullanım Senaryoları:
	//   - Form render sırasında alan görünürlüğü
	//   - Client-side dinamik form
	//   - Validasyon sırası belirleme
	//
	// Örnek:
	//   field := fields.Text("Tax Number", "tax_number").
	//       DependsOn("is_company").
	//       When("is_company", "=", true)
	//
	//   context := map[string]interface{}{"is_company": true}
	//   visible := field.ResolveDependencies(context) // true
	ResolveDependencies(context interface{}) bool

	// ============================================================================
	// Öneri Metodları (Suggestion Methods)
	// ============================================================================
	//
	// Bu metodlar, kullanıcıya dinamik öneriler sunmak için kullanılır.

	// GetSuggestionsCallback, öneri callback fonksiyonunu döndürür.
	//
	// Bu callback, kullanıcının girdiği sorguya göre öneri listesi oluşturur.
	// Autocomplete ve suggestion özellikleri için kullanılır.
	//
	// Döndürür:
	//   - Öneri callback fonksiyonu (string alır, []interface{} döner) veya nil
	//
	// Kullanım Senaryoları:
	//   - Statik öneri listeleri
	//   - Hesaplanan öneriler
	//   - Kullanıcı geçmişi bazlı öneriler
	//
	// Örnek:
	//   callback := field.GetSuggestionsCallback()
	//   if callback != nil {
	//       suggestions := callback("john") // ["john@gmail.com", "john@hotmail.com"]
	//   }
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Öneriler ve autocomplete bölümü
	GetSuggestionsCallback() func(string) []interface{}

	// GetAutoCompleteURL, autocomplete için API endpoint URL'ini döndürür.
	//
	// Bu URL, kullanıcının girdiği sorguya göre önerileri almak için
	// kullanılır. Backend'den dinamik öneri yükleme için kullanılır.
	//
	// Döndürür:
	//   - Autocomplete API endpoint URL'i (örn: "/api/cities/search")
	//
	// Kullanım Senaryoları:
	//   - Veritabanından öneri yükleme
	//   - Harici API'lerden öneri alma
	//   - Dinamik filtreleme
	//
	// Örnek:
	//   url := field.GetAutoCompleteURL() // "/api/cities/search"
	//   // Frontend: GET /api/cities/search?q=ist
	//   // Response: ["Istanbul", "Islamabad", "Isparta"]
	GetAutoCompleteURL() string

	// GetMinCharsForSuggestions, önerilerin gösterilmesi için gereken minimum karakter sayısını döndürür.
	//
	// Bu değer, kullanıcının kaç karakter girdikten sonra önerilerin
	// gösterileceğini belirler. Performans optimizasyonu için kullanılır.
	//
	// Döndürür:
	//   - Minimum karakter sayısı (genellikle 2 veya 3)
	//
	// Önemli Not:
	//   - Çok düşük değer (1) performans sorunlarına yol açabilir
	//   - Çok yüksek değer (5+) kullanıcı deneyimini olumsuz etkiler
	//   - Önerilen değer: 2-3 karakter
	//
	// Örnek:
	//   minChars := field.GetMinCharsForSuggestions() // 3
	//   // Kullanıcı "ist" yazdığında öneriler gösterilir
	GetMinCharsForSuggestions() int

	// GetSuggestions, verilen sorgu için öneri listesi döndürür.
	//
	// Bu metod, suggestions callback'ini veya autocomplete URL'ini kullanarak
	// öneri listesi oluşturur. Frontend tarafından kullanılır.
	//
	// Parametreler:
	//   - query: Arama sorgusu (kullanıcının girdiği metin)
	//
	// Döndürür:
	//   - Öneri listesi ([]interface{})
	//
	// Kullanım Senaryoları:
	//   - Autocomplete dropdown
	//   - Search suggestions
	//   - Quick filters
	//
	// Örnek:
	//   suggestions := field.GetSuggestions("john")
	//   // ["john@gmail.com", "john@hotmail.com", "john@outlook.com"]
	GetSuggestions(query string) []interface{}

	// ============================================================================
	// Dosya Ekleme Metodları (Attachment Methods)
	// ============================================================================
	//
	// Bu metodlar, dosya yükleme işlemlerini yönetir ve yapılandırır.

	// GetAcceptedMimeTypes, kabul edilen MIME tiplerinin listesini döndürür.
	//
	// MIME tipleri, hangi dosya formatlarının yüklenebileceğini belirler.
	// Frontend'de dosya seçici ve backend'de validasyon için kullanılır.
	//
	// Döndürür:
	//   - Kabul edilen MIME tipleri listesi ([]string)
	//
	// Yaygın MIME Tipleri:
	//   - Görseller: "image/jpeg", "image/png", "image/webp", "image/gif"
	//   - Videolar: "video/mp4", "video/webm", "video/quicktime"
	//   - Ses: "audio/mpeg", "audio/wav", "audio/ogg"
	//   - Dökümanlar: "application/pdf", "application/msword"
	//
	// Örnek:
	//   mimeTypes := field.GetAcceptedMimeTypes()
	//   // ["image/jpeg", "image/png", "image/webp"]
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Dosya yükleme yapılandırması bölümü
	GetAcceptedMimeTypes() []string

	// GetMaxFileSize, maksimum dosya boyutunu byte cinsinden döndürür.
	//
	// Bu değer, yüklenebilecek dosyanın maksimum boyutunu belirler.
	// Hem frontend hem backend'de validasyon için kullanılır.
	//
	// Döndürür:
	//   - Maksimum dosya boyutu (byte cinsinden, int64)
	//
	// Önerilen Boyutlar:
	//   - Profil fotoğrafları: 2-5 MB (2*1024*1024 - 5*1024*1024)
	//   - Dökümanlar: 10 MB (10*1024*1024)
	//   - Videolar: 100 MB (100*1024*1024)
	//   - Ses dosyaları: 50 MB (50*1024*1024)
	//
	// Önemli Not:
	//   - Server'ın upload limit'i de kontrol edilmelidir
	//   - Çok büyük dosyalar için chunk upload düşünülmelidir
	//
	// Örnek:
	//   maxSize := field.GetMaxFileSize() // 5242880 (5 MB)
	GetMaxFileSize() int64

	// GetStorageDisk, depolama diskinin adını döndürür.
	//
	// Depolama diski, dosyanın nereye kaydedileceğini belirler.
	// Genellikle "public", "private", "s3" gibi değerler kullanılır.
	//
	// Döndürür:
	//   - Depolama diski adı (örn: "public", "private", "s3")
	//
	// Disk Tipleri:
	//   - "public": Herkese açık dosyalar (görseller, vb.)
	//   - "private": Gizli dosyalar (dökümanlar, vb.)
	//   - "s3": AWS S3 depolama
	//   - "local": Yerel dosya sistemi
	//
	// Örnek:
	//   disk := field.GetStorageDisk() // "public"
	GetStorageDisk() string

	// GetStoragePath, depolama yolunu döndürür.
	//
	// Depolama yolu, dosyanın disk içindeki hangi klasöre kaydedileceğini
	// belirler. Organizasyon ve güvenlik için kullanılır.
	//
	// Döndürür:
	//   - Depolama yolu (örn: "avatars", "documents", "products/images")
	//
	// Önerilen Yapı:
	//   - Tip bazlı: "avatars", "documents", "videos"
	//   - Resource bazlı: "users/avatars", "products/images"
	//   - Tarih bazlı: "uploads/2024/01"
	//
	// Örnek:
	//   path := field.GetStoragePath() // "users/avatars"
	//   // Tam yol: /storage/public/users/avatars/filename.jpg
	GetStoragePath() string

	// ValidateAttachment, bir dosya ekini doğrular.
	//
	// Bu metod, dosya adı ve boyutunu kontrol ederek dosyanın
	// yüklenebilir olup olmadığını belirler. MIME tipi ve boyut
	// kontrolü yapar.
	//
	// Parametreler:
	//   - filename: Dosya adı (uzantı kontrolü için)
	//   - size: Dosya boyutu (byte cinsinden)
	//
	// Döndürür:
	//   - nil: Dosya geçerli
	//   - error: Validasyon hatası (detaylı mesaj içerir)
	//
	// Kontroller:
	//   - Dosya uzantısı kabul edilen MIME tiplerinde mi?
	//   - Dosya boyutu maksimum boyutun altında mı?
	//   - Dosya adı geçerli mi?
	//
	// Örnek:
	//   err := field.ValidateAttachment("photo.jpg", 3145728) // 3 MB
	//   if err != nil {
	//       // Validasyon hatası: "Dosya boyutu çok büyük"
	//   }
	ValidateAttachment(filename string, size int64) error

	// GetUploadCallback, dosya yükleme callback fonksiyonunu döndürür.
	//
	// Bu callback, dosya yükleme işleminde özel mantık uygulamak için
	// kullanılır. Dosya işleme, yeniden boyutlandırma, watermark ekleme
	// gibi işlemler için kullanılabilir.
	//
	// Döndürür:
	//   - Upload callback fonksiyonu (file, item alır, error döner) veya nil
	//
	// Kullanım Senaryoları:
	//   - Görsel yeniden boyutlandırma
	//   - Thumbnail oluşturma
	//   - Watermark ekleme
	//   - Video transcoding
	//   - Metadata çıkarma
	//
	// Örnek:
	//   callback := field.GetUploadCallback()
	//   if callback != nil {
	//       err := callback(file, item)
	//       if err != nil {
	//           // Upload hatası
	//       }
	//   }
	GetUploadCallback() func(interface{}, interface{}) error

	// ShouldRemoveEXIFData, EXIF verilerinin kaldırılıp kaldırılmayacağını döndürür.
	//
	// EXIF verileri, fotoğraflarda konum, kamera bilgisi gibi metadata içerir.
	// Gizlilik için bu verilerin kaldırılması önerilir.
	//
	// Döndürür:
	//   - true: EXIF verileri kaldırılmalı
	//   - false: EXIF verileri korunmalı
	//
	// EXIF Verileri İçerir:
	//   - GPS koordinatları (konum bilgisi)
	//   - Kamera modeli ve ayarları
	//   - Çekim tarihi ve saati
	//   - Yazılım bilgisi
	//
	// Önemli Not:
	//   - Gizlilik için EXIF verilerinin kaldırılması önerilir
	//   - Özellikle kullanıcı yüklü görsellerde kritik
	//
	// Örnek:
	//   shouldRemove := field.ShouldRemoveEXIFData() // true
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Görsel yükleme yapılandırması bölümü
	ShouldRemoveEXIFData() bool

	// RemoveEXIFData, bir dosyadan EXIF verilerini kaldırır.
	//
	// Bu metod, görsel dosyalardan EXIF metadata'sını temizler.
	// Gizlilik ve güvenlik için kullanılır.
	//
	// Parametreler:
	//   - ctx: İşlem context'i
	//   - file: EXIF verileri kaldırılacak dosya
	//
	// Döndürür:
	//   - nil: İşlem başarılı
	//   - error: İşlem hatası
	//
	// Desteklenen Formatlar:
	//   - JPEG/JPG
	//   - TIFF
	//   - PNG (sınırlı)
	//
	// Örnek:
	//   err := field.RemoveEXIFData(ctx, file)
	//   if err != nil {
	//       // EXIF kaldırma hatası
	//   }
	RemoveEXIFData(ctx interface{}, file interface{}) error

	// ============================================================================
	// Tekrarlayıcı Alan Metodları (Repeater Methods)
	// ============================================================================
	//
	// Bu metodlar, dinamik olarak tekrarlanan alan gruplarını yönetir.

	// IsRepeaterField, bu element'in bir tekrarlayıcı alan olup olmadığını döndürür.
	//
	// Tekrarlayıcı alanlar, kullanıcının dinamik olarak birden fazla değer
	// grubu eklemesine izin verir. Örneğin, telefon numaraları, adresler,
	// sosyal medya hesapları gibi.
	//
	// Döndürür:
	//   - true: Bu bir tekrarlayıcı alan
	//   - false: Bu normal bir alan
	//
	// Kullanım Senaryoları:
	//   - Çoklu telefon numaraları
	//   - Birden fazla adres
	//   - Sosyal medya hesapları listesi
	//   - Eğitim geçmişi
	//   - İş deneyimi
	//
	// Örnek:
	//   isRepeater := field.IsRepeaterField() // true
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Repeater fields bölümü
	IsRepeaterField() bool

	// GetRepeaterFields, tekrarlayıcı alan içindeki alt alanları döndürür.
	//
	// Her tekrar için gösterilecek alan listesini içerir. Bu alanlar,
	// her tekrarda aynı yapıda gösterilir.
	//
	// Döndürür:
	//   - Alt alan listesi ([]Element)
	//
	// Örnek:
	//   fields := field.GetRepeaterFields()
	//   // [
	//   //   Text("Type", "type"),
	//   //   Tel("Number", "number")
	//   // ]
	GetRepeaterFields() []Element

	// GetMinRepeats, minimum tekrar sayısını döndürür.
	//
	// Kullanıcının en az kaç tane değer grubu eklemesi gerektiğini belirler.
	// Validasyon için kullanılır.
	//
	// Döndürür:
	//   - Minimum tekrar sayısı (int)
	//
	// Önerilen Değerler:
	//   - Zorunlu alan: 1 veya daha fazla
	//   - İsteğe bağlı alan: 0
	//
	// Örnek:
	//   minRepeats := field.GetMinRepeats() // 1
	//   // Kullanıcı en az 1 telefon numarası eklemeli
	GetMinRepeats() int

	// GetMaxRepeats, maksimum tekrar sayısını döndürür.
	//
	// Kullanıcının en fazla kaç tane değer grubu ekleyebileceğini belirler.
	// Performans ve kullanılabilirlik için sınırlama koyar.
	//
	// Döndürür:
	//   - Maksimum tekrar sayısı (int)
	//
	// Önerilen Değerler:
	//   - Telefon numaraları: 3-5
	//   - Adresler: 3-5
	//   - Sosyal medya: 5-10
	//   - Eğitim/İş deneyimi: 10-20
	//
	// Önemli Not:
	//   - Çok yüksek değer performans sorunlarına yol açabilir
	//   - Kullanıcı deneyimi için makul bir limit belirleyin
	//
	// Örnek:
	//   maxRepeats := field.GetMaxRepeats() // 5
	//   // Kullanıcı en fazla 5 telefon numarası ekleyebilir
	GetMaxRepeats() int

	// ValidateRepeats, tekrar sayısını doğrular.
	//
	// Verilen tekrar sayısının minimum ve maksimum sınırlar içinde
	// olup olmadığını kontrol eder.
	//
	// Parametreler:
	//   - count: Kontrol edilecek tekrar sayısı
	//
	// Döndürür:
	//   - nil: Tekrar sayısı geçerli
	//   - error: Validasyon hatası (detaylı mesaj içerir)
	//
	// Örnek:
	//   err := field.ValidateRepeats(3)
	//   if err != nil {
	//       // Hata: "En az 1, en fazla 5 tekrar ekleyebilirsiniz"
	//   }
	ValidateRepeats(count int) error

	// ============================================================================
	// Zengin Metin Editörü Metodları (Rich Text Methods)
	// ============================================================================
	//
	// Bu metodlar, zengin metin editörü yapılandırmasını yönetir.

	// GetEditorType, editör tipini döndürür.
	//
	// Editör tipi, hangi WYSIWYG editörünün kullanılacağını belirler.
	// Farklı editörler farklı özellikler ve kullanıcı deneyimi sunar.
	//
	// Döndürür:
	//   - Editör tipi ("tiptap", "quill", "tinymce")
	//
	// Editör Tipleri:
	//   - "tiptap": Modern, modüler, Vue/React uyumlu
	//   - "quill": Hafif, hızlı, basit
	//   - "tinymce": Zengin özellikli, kurumsal
	//
	// Örnek:
	//   editorType := field.GetEditorType() // "tiptap"
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Zengin metin editörü yapılandırması bölümü
	GetEditorType() string

	// GetEditorLanguage, editör dilini döndürür.
	//
	// Editör arayüzünün hangi dilde gösterileceğini belirler.
	// Menüler, butonlar ve mesajlar bu dilde gösterilir.
	//
	// Döndürür:
	//   - Editör dili (örn: "tr", "en", "tr_TR")
	//
	// Desteklenen Diller:
	//   - "tr": Türkçe
	//   - "en": İngilizce
	//   - "de": Almanca
	//   - "fr": Fransızca
	//   - vb.
	//
	// Örnek:
	//   language := field.GetEditorLanguage() // "tr"
	GetEditorLanguage() string

	// GetEditorTheme, editör temasını döndürür.
	//
	// Editörün görsel temasını belirler. Tema, editörün renklerini,
	// stilini ve genel görünümünü etkiler.
	//
	// Döndürür:
	//   - Editör teması (örn: "snow", "bubble", "default")
	//
	// Yaygın Temalar:
	//   - "snow": Beyaz, temiz tema (Quill)
	//   - "bubble": Baloncuk tarzı toolbar (Quill)
	//   - "default": Varsayılan tema
	//   - "dark": Koyu tema
	//
	// Örnek:
	//   theme := field.GetEditorTheme() // "snow"
	GetEditorTheme() string

	// ============================================================================
	// Durum Metodları (Status Methods)
	// ============================================================================
	//
	// Bu metodlar, durum gösterimini ve renklendirilmesini yönetir.

	// GetStatusColors, durum renkleri eşleştirmesini döndürür.
	//
	// Her durum değeri için bir renk tanımlar. Liste görünümünde
	// durumları renkli badge'ler olarak göstermek için kullanılır.
	//
	// Döndürür:
	//   - Durum-renk eşleştirme map'i (map[string]string)
	//
	// Yaygın Renkler:
	//   - "green": Başarılı, aktif, onaylanmış
	//   - "red": Hata, pasif, reddedilmiş
	//   - "yellow"/"orange": Uyarı, beklemede
	//   - "blue": Bilgi, işlemde
	//   - "gray": Nötr, taslak, arşivlenmiş
	//
	// Örnek:
	//   colors := field.GetStatusColors()
	//   // {
	//   //   "draft": "gray",
	//   //   "published": "green",
	//   //   "archived": "red"
	//   // }
	//
	// İlgili Dokümantasyon:
	//   - docs/Fields.md: Durum renkleri ve badge'ler bölümü
	GetStatusColors() map[string]string

	// GetBadgeVariant, badge varyantını döndürür.
	//
	// Badge'in görsel stilini belirler. Solid, outline, subtle gibi
	// farklı varyantlar farklı görsel etkiler yaratır.
	//
	// Döndürür:
	//   - Badge varyantı ("solid", "outline", "subtle")
	//
	// Varyantlar:
	//   - "solid": Dolu, renkli arka plan
	//   - "outline": Çerçeveli, şeffaf arka plan
	//   - "subtle": Hafif renkli arka plan
	//
	// Örnek:
	//   variant := field.GetBadgeVariant() // "solid"
	GetBadgeVariant() string

	// ============================================================================
	// Pivot Metodları (Pivot Methods)
	// ============================================================================
	//
	// Bu metodlar, çoktan çoğa ilişkilerde pivot tablo alanlarını yönetir.

	// IsPivot, bu element'in bir pivot alan olup olmadığını döndürür.
	//
	// Pivot alanlar, çoktan çoğa ilişkilerde ara tabloda (pivot table)
	// saklanan ek bilgiler için kullanılır. Örneğin, kullanıcı-rol
	// ilişkisinde atanma tarihi gibi.
	//
	// Döndürür:
	//   - true: Bu bir pivot alan
	//   - false: Bu normal bir alan
	//
	// Kullanım Senaryoları:
	//   - Atanma tarihi (assigned_at)
	//   - Bitiş tarihi (expires_at)
	//   - Öncelik (priority)
	//   - Notlar (notes)
	//
	// Örnek:
	//   isPivot := field.IsPivot() // true
	//
	// İlgili Dokümantasyon:
	//   - docs/Relationships.md: BelongsToMany ve pivot fields bölümü
	//   - docs/Fields.md: Pivot fields bölümü
	IsPivot() bool

	// GetPivotResourceName, pivot resource adını döndürür.
	//
	// Pivot alanın hangi resource'a ait olduğunu belirtir. Bu bilgi,
	// pivot verilerini doğru tabloya kaydetmek için kullanılır.
	//
	// Döndürür:
	//   - Pivot resource adı (örn: "user_roles", "post_tags")
	//
	// Örnek:
	//   resourceName := field.GetPivotResourceName() // "user_roles"
	//   // Bu alan user_roles pivot tablosunda saklanır
	GetPivotResourceName() string
}
