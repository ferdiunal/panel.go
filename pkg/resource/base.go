package resource

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ferdiunal/panel.go/pkg/auth"
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// / # Base Yapısı
// /
// / Bu yapı, temel Resource implementasyonunu sağlar ve `resource.Resource` arayüzünü implemente eder.
// / Özel resource implementasyonları oluştururken bu yapı gömülü (embedded) olarak kullanılabilir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Basit Resource Oluşturma**: Hızlı prototipleme için doğrudan kullanılabilir
// / 2. **Özel Resource Geliştirme**: Embedding yoluyla özel davranışlar eklenebilir
// / 3. **CRUD İşlemleri**: Veritabanı işlemleri için temel yapı sağlar
// / 4. **Admin Panel Entegrasyonu**: Panel.go admin arayüzü ile tam entegrasyon
// /
// / ## Özellikler
// /
// / - **Veri Modeli Yönetimi**: Herhangi bir Go struct'ı ile çalışabilir
// / - **Alan Yönetimi**: Form alanları ve görüntüleme alanları için destek (bkz. [docs/Fields.md](../../docs/Fields.md))
// / - **İlişki Desteği**: GORM ilişkileri için tam destek (bkz. [docs/Relationships.md](../../docs/Relationships.md))
// / - **Yetkilendirme**: Policy tabanlı erişim kontrolü
// / - **Dosya Yükleme**: Özelleştirilebilir dosya yükleme işleyicisi
// / - **Widget Desteği**: Dashboard widget'ları için destek
// / - **Filtreleme ve Sıralama**: Gelişmiş veri filtreleme ve sıralama özellikleri
// /
// / ## Örnek Kullanım
// /
// / ```go
// / type UserResource struct {
// /     resource.Base
// / }
// /
// / func NewUserResource() *UserResource {
// /     return &UserResource{
// /         Base: resource.Base{
// /             DataModel:  &models.User{},
// /             Identifier: "users",
// /             Label:      "Kullanıcılar",
// /             IconName:   "users",
// /             GroupName:  "Yönetim",
// /             FieldsVal: []fields.Element{
// /                 fields.NewText("name").SetLabel("Ad"),
// /                 fields.NewEmail("email").SetLabel("E-posta"),
// /             },
// /             PolicyVal: &UserPolicy{},
// /         },
// /     }
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Hızlı Geliştirme**: Temel CRUD işlemleri için hazır implementasyon
// / - **Esneklik**: Embedding ile kolayca özelleştirilebilir
// / - **Tip Güvenliği**: Go'nun tip sistemi ile tam uyumlu
// / - **Modülerlik**: Her bileşen bağımsız olarak yapılandırılabilir
// /
// / ## Dikkat Edilmesi Gerekenler
// /
// / - `DataModel` alanı mutlaka doldurulmalıdır
// / - `Identifier` benzersiz olmalı ve URL-safe karakterler içermelidir
// / - `FieldsVal` alanları model yapısı ile uyumlu olmalıdır
// / - Policy tanımlanmazsa varsayılan olarak tüm işlemler reddedilir
// /
// / ## İlgili Tipler
// /
// / - `Resource`: Ana resource arayüzü
// / - `fields.Element`: Alan tanımlamaları için
// / - `auth.Policy`: Yetkilendirme politikaları için
// / - `widget.Card`: Dashboard widget'ları için
type Base struct {
	/// Veritabanı modeli - Herhangi bir Go struct'ı olabilir (genellikle GORM modeli)
	DataModel any

	/// URL tanımlayıcısı - Benzersiz, URL-safe string (örn: "users", "blog-posts")
	Identifier string

	/// Görünen başlık - Kullanıcı arayüzünde gösterilecek insan okunabilir isim
	Label string

	/// Menü ikonu - Icon kütüphanesinden icon adı (örn: "users", "settings")
	IconName string

	/// Menü grubu - İlgili resource'ları gruplamak için kullanılır
	GroupName string

	/// Dinamik başlık fonksiyonu - i18n desteği için kullanılır
	titleFunc func(*fiber.Ctx) string

	/// Dinamik grup fonksiyonu - i18n desteği için kullanılır
	groupFunc func(*fiber.Ctx) string

	/// Alan tanımlamaları - Form ve liste görünümlerinde kullanılacak alanlar
	/// Detaylı bilgi için bkz: docs/Fields.md
	FieldsVal []fields.Element

	/// Widget tanımlamaları - Dashboard'da gösterilecek widget'lar
	WidgetsVal []widget.Card

	/// Sıralama ayarları - Varsayılan sıralama kuralları
	Sortable []Sortable

	/// Yetkilendirme politikası - CRUD işlemleri için erişim kontrolü
	PolicyVal auth.Policy

	/// Diyalog tipi - Oluşturma/düzenleme işlemleri için kullanılacak diyalog türü
	DialogType DialogType

	/// Dosya yükleme işleyicisi - Özel dosya yükleme mantığı için
	/// nil ise varsayılan yükleme mantığı kullanılır
	UploadHandler func(c *appContext.Context, file *multipart.FileHeader) (string, error)

	/// Seed ayarları - Veritabanı seed işlemleri için yapılandırma
	Seed SettingsSeed

	/// Özel işlemler - Resource üzerinde çalıştırılabilecek özel aksiyonlar
	ActionsVal []Action

	/// Filtre tanımlamaları - Liste görünümünde kullanılacak filtreler
	FiltersVal []Filter

	/// OpenAPI görünürlük kontrolü - false ise OpenAPI spec'te görünür (varsayılan)
	openAPIDisabled bool

	/// Kayıt başlığı için kullanılacak field adı - İlişki fieldlarında gösterilir
	recordTitleKey string

	/// Kayıt başlığını özel fonksiyon ile hesaplamak için kullanılır
	recordTitleFunc func(record any) string
}

// / # SettingsSeed Yapısı
// /
// / Bu yapı, veritabanı seed işlemleri veya ayar gruplandırması için yardımcı bir yapıdır.
// / Resource'ların başlangıç verilerini tanımlamak için kullanılır.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Veritabanı Seed**: İlk kurulumda örnek verileri eklemek için
// / 2. **Varsayılan Ayarlar**: Resource için varsayılan yapılandırma değerleri
// / 3. **Test Verileri**: Test ortamları için tutarlı veri setleri
// /
// / ## Örnek Kullanım
// /
// / ```go
// / seed := SettingsSeed{
// /     Key: "default_users",
// /     Value: map[string]any{
// /         "admin": map[string]any{
// /             "email": "admin@example.com",
// /             "role":  "admin",
// /         },
// /         "user": map[string]any{
// /             "email": "user@example.com",
// /             "role":  "user",
// /         },
// /     },
// / }
// / ```
// /
// / ## Alanlar
// /
// / - `Key`: Seed grubunu tanımlayan benzersiz anahtar
// / - `Value`: Seed verilerini içeren map yapısı
type SettingsSeed struct {
	/// Seed grubunun benzersiz tanımlayıcısı
	Key string

	/// Seed verilerini içeren esnek map yapısı
	/// Herhangi bir veri yapısını destekler
	Value map[string]any
}

// / # Model Metodu
// /
// / Bu fonksiyon, resource'a bağlı veri modelini döndürür.
// / Genellikle GORM modeli veya herhangi bir Go struct'ı olabilir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Veritabanı İşlemleri**: GORM sorguları için model tipini belirlemek
// / 2. **Reflection**: Model yapısını dinamik olarak incelemek
// / 3. **Validasyon**: Model alanlarını doğrulamak
// / 4. **Migration**: Veritabanı şeması oluşturmak
// /
// / ## Döndürür
// /
// / - `any`: Resource'ın veri modeli (genellikle struct pointer'ı)
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := NewUserResource()
// / model := resource.Model() // &models.User{} döner
// /
// / // GORM ile kullanım
// / db.Model(model).Find(&users)
// / ```
// /
// / ## Notlar
// /
// / - Model nil olmamalıdır
// / - Genellikle struct pointer'ı olarak tanımlanır
// / - GORM işlemleri için tip bilgisi sağlar
func (r Base) Model() any {
	return r.DataModel
}

// / # Slug Metodu
// /
// / Bu fonksiyon, resource'ın URL-safe benzersiz tanımlayıcısını döndürür.
// / API endpoint'leri ve routing için kullanılır.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **API Routing**: `/api/resources/{slug}` formatında URL oluşturma
// / 2. **Menü Oluşturma**: Navigasyon menüsünde benzersiz tanımlama
// / 3. **Cache Key**: Resource bazlı cache anahtarları oluşturma
// / 4. **Permission Check**: Yetkilendirme kontrollerinde resource tanımlama
// /
// / ## Döndürür
// /
// / - `string`: URL-safe benzersiz tanımlayıcı (örn: "users", "blog-posts")
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := NewUserResource()
// / slug := resource.Slug() // "users" döner
// /
// / // API endpoint oluşturma
// / endpoint := fmt.Sprintf("/api/%s", slug) // "/api/users"
// / ```
// /
// / ## Önemli Notlar
// /
// / - Slug benzersiz olmalıdır
// / - Sadece küçük harf, rakam ve tire (-) içermelidir
// / - Boşluk veya özel karakter içermemelidir
// / - Genellikle çoğul isim kullanılır (users, posts, categories)
func (r Base) Slug() string {
	return r.Identifier
}

// / # Title Metodu
// /
// / Bu fonksiyon, resource'ın kullanıcı arayüzünde gösterilecek insan okunabilir başlığını döndürür.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Menü Başlıkları**: Navigasyon menüsünde gösterim
// / 2. **Sayfa Başlıkları**: Liste ve detay sayfalarında başlık
// / 3. **Breadcrumb**: Sayfa yolu gösteriminde
// / 4. **Bildirimler**: Kullanıcı bildirimlerinde resource adı
// /
// / ## Döndürür
// /
// / - `string`: İnsan okunabilir başlık (örn: "Kullanıcılar", "Blog Yazıları")
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := NewUserResource()
// / title := resource.Title() // "Kullanıcılar" döner
// /
// / // Sayfa başlığı oluşturma
// / pageTitle := fmt.Sprintf("%s Listesi", title) // "Kullanıcılar Listesi"
// / ```
// /
// / ## Notlar
// /
// / - Genellikle çoğul isim kullanılır
// / - Türkçe karakter içerebilir
// / - Büyük harfle başlamalıdır
func (r Base) Title() string {
	return r.Label
}

// / TitleWithContext, kaynağın kullanıcı arayüzünde görünecek başlığını döner.
// /
// / Bu metod, SetTitleFunc ile ayarlanan dinamik başlık fonksiyonunu kullanır.
// / Eğer titleFunc ayarlanmamışsa, Title() metodunu fallback olarak kullanır.
// /
// / ## Parametreler
// / - `ctx`: Fiber context (i18n için gerekli)
// /
// / ## Döndürür
// / - `string`: Kullanıcı dostu başlık
// /
// / ## Örnek
// / ```go
// / title := resource.TitleWithContext(c.Ctx)
// / ```
func (r Base) TitleWithContext(ctx *fiber.Ctx) string {
	if r.titleFunc != nil && ctx != nil {
		return r.titleFunc(ctx)
	}
	if ctx != nil {
		return i18n.Trans(ctx, r.Title())
	}
	return r.Title()
}

// / # Icon Metodu
// /
// / Bu fonksiyon, resource'ın menüde gösterilecek ikon adını döndürür.
// / Icon kütüphanesinden (örn: Heroicons, FontAwesome) ikon adı kullanılır.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Navigasyon Menüsü**: Menü öğelerinde ikon gösterimi
// / 2. **Dashboard**: Dashboard widget'larında ikon
// / 3. **Butonlar**: İşlem butonlarında ikon
// / 4. **Breadcrumb**: Sayfa yolunda görsel gösterim
// /
// / ## Döndürür
// /
// / - `string`: Icon adı (örn: "users", "document-text", "cog")
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := NewUserResource()
// / icon := resource.Icon() // "users" döner
// /
// / // Frontend'de kullanım
// / // <Icon name={icon} />
// / ```
// /
// / ## Desteklenen Icon Kütüphaneleri
// /
// / - Heroicons (varsayılan)
// / - FontAwesome
// / - Material Icons
// / - Custom SVG icons
// /
// / ## Notlar
// /
// / - Icon adı kütüphane ile uyumlu olmalıdır
// / - Boş string döndürülürse varsayılan ikon kullanılır
func (r Base) Icon() string {
	return r.IconName
}

// / # Group Metodu
// /
// / Bu fonksiyon, resource'ın menüde hangi gruba ait olduğunu belirtir.
// / İlgili resource'ları gruplamak için kullanılır.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Menü Organizasyonu**: İlgili resource'ları gruplamak
// / 2. **Yetkilendirme**: Grup bazlı erişim kontrolü
// / 3. **Dashboard**: Grup bazlı widget organizasyonu
// / 4. **Raporlama**: Grup bazlı raporlar
// /
// / ## Döndürür
// /
// / - `string`: Grup adı (örn: "Yönetim", "İçerik", "Ayarlar")
// /
// / ## Örnek Kullanım
// /
// / ```go
// / userResource := NewUserResource()
// / group := userResource.Group() // "Yönetim" döner
// /
// / roleResource := NewRoleResource()
// / roleGroup := roleResource.Group() // "Yönetim" döner
// /
// / // Aynı gruptaki resource'lar menüde birlikte gösterilir
// / ```
// /
// / ## Grup Örnekleri
// /
// / - **Yönetim**: Kullanıcılar, Roller, İzinler
// / - **İçerik**: Blog, Sayfalar, Medya
// / - **E-ticaret**: Ürünler, Siparişler, Kategoriler
// / - **Ayarlar**: Genel Ayarlar, E-posta, Entegrasyonlar
// /
// / ## Notlar
// /
// / - Boş string döndürülürse grupsuz gösterilir
// / - Grup adları tutarlı olmalıdır
func (r Base) Group() string {
	return r.GroupName
}

// / GroupWithContext, kaynağın menüde hangi grup altında listeleneceğini belirler.
// /
// / Bu metod, SetGroupFunc ile ayarlanan dinamik grup fonksiyonunu kullanır.
// / Eğer groupFunc ayarlanmamışsa, Group() metodunu fallback olarak kullanır.
// /
// / ## Parametreler
// / - `ctx`: Fiber context (i18n için gerekli)
// /
// / ## Döndürür
// / - `string`: Grup adı
// /
// / ## Örnek
// / ```go
// / group := resource.GroupWithContext(c.Ctx)
// / ```
func (r Base) GroupWithContext(ctx *fiber.Ctx) string {
	if r.groupFunc != nil && ctx != nil {
		return r.groupFunc(ctx)
	}
	if ctx != nil {
		return i18n.Trans(ctx, r.Group())
	}
	return r.Group()
}

// / # Fields Metodu
// /
// / Bu fonksiyon, resource'ın tüm alan tanımlamalarını döndürür.
// / Form ve liste görünümlerinde kullanılacak alanları içerir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Form Oluşturma**: Oluşturma/düzenleme formlarında alan gösterimi
// / 2. **Liste Görünümü**: Tablo kolonlarını oluşturma
// / 3. **Detay Görünümü**: Kayıt detaylarını gösterme
// / 4. **Validasyon**: Alan bazlı doğrulama kuralları
// / 5. **Filtreleme**: Filtrelenebilir alanları belirleme
// /
// / ## Döndürür
// /
// / - `[]fields.Element`: Alan tanımlamaları dizisi
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := NewUserResource()
// / fields := resource.Fields()
// /
// / // Alanları döngüyle işleme
// / for _, field := range fields {
// /     fmt.Printf("Alan: %s, Tip: %s\n", field.GetKey(), field.GetType())
// / }
// / ```
// /
// / ## Alan Tipleri
// /
// / Detaylı bilgi için bkz: [docs/Fields.md](../../docs/Fields.md)
// /
// / - **Text**: Metin girişi
// / - **Email**: E-posta girişi
// / - **Password**: Şifre girişi
// / - **Number**: Sayı girişi
// / - **Date**: Tarih seçici
// / - **Select**: Açılır liste
// / - **BelongsTo**: İlişkili kayıt seçimi
// / - **HasMany**: Çoklu ilişki yönetimi
// /
// / ## Notlar
// /
// / - Alanlar sıralı olarak gösterilir
// / - Her alan benzersiz key'e sahip olmalıdır
// / - Alan tipleri model yapısı ile uyumlu olmalıdır
func (r Base) Fields() []fields.Element {
	return r.FieldsVal
}

// / # With Metodu
// /
// / Bu fonksiyon, GORM eager loading için yüklenecek ilişkileri belirtir.
// / Varsayılan implementasyon boş dizi döndürür.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **N+1 Problem Önleme**: İlişkili verileri tek sorguda yükleme
// / 2. **Performans Optimizasyonu**: Gereksiz sorgu sayısını azaltma
// / 3. **İlişki Gösterimi**: Liste görünümünde ilişkili verileri gösterme
// /
// / ## Döndürür
// /
// / - `[]string`: Eager loading yapılacak ilişki adları (varsayılan: boş)
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da override etme
// / func (r UserResource) With() []string {
// /     return []string{"Role", "Profile", "Posts"}
// / }
// /
// / // GORM ile kullanım
// / db.Preload(resource.With()...).Find(&users)
// / ```
// /
// / ## İlişki Örnekleri
// /
// / Detaylı bilgi için bkz: [docs/Relationships.md](../../docs/Relationships.md)
// /
// / - **BelongsTo**: `"Role"`, `"Category"`
// / - **HasOne**: `"Profile"`, `"Settings"`
// / - **HasMany**: `"Posts"`, `"Comments"`
// / - **ManyToMany**: `"Tags"`, `"Permissions"`
// /
// / ## Performans Notları
// /
// / - Sadece gerekli ilişkileri yükleyin
// / - Derin ilişkiler için nokta notasyonu kullanın: `"Posts.Comments"`
// / - Çok fazla ilişki performansı düşürebilir
// /
// / ## Önemli Uyarılar
// /
// / - İlişki adları model yapısındaki field adları ile eşleşmelidir
// / - Var olmayan ilişki adı hata verir
// / - Circular reference'lardan kaçının
func (r Base) With() []string {
	return []string{}
}

// / # Lenses Metodu
// /
// / Bu fonksiyon, resource için tanımlı özel görünümleri (lens) döndürür.
// / Varsayılan implementasyon boş dizi döndürür.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Özel Filtreler**: Önceden tanımlı filtre kombinasyonları
// / 2. **Özel Görünümler**: Farklı veri sunumları
// / 3. **Raporlar**: Özel raporlama görünümleri
// / 4. **Dashboard**: Özelleştirilmiş veri görünümleri
// /
// / ## Döndürür
// /
// / - `[]Lens`: Lens tanımlamaları dizisi (varsayılan: boş)
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da override etme
// / func (r UserResource) Lenses() []Lens {
// /     return []Lens{
// /         {
// /             Name:  "active-users",
// /             Label: "Aktif Kullanıcılar",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("status = ?", "active")
// /             },
// /         },
// /         {
// /             Name:  "premium-users",
// /             Label: "Premium Kullanıcılar",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("subscription_type = ?", "premium")
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ## Lens Özellikleri
// /
// / - **Name**: Benzersiz tanımlayıcı
// / - **Label**: Kullanıcı arayüzünde gösterilecek başlık
// / - **Query**: GORM sorgu modifikasyonu
// / - **Icon**: Opsiyonel ikon
// /
// / ## Avantajlar
// /
// / - Karmaşık filtreleri basitleştirir
// / - Kullanıcı deneyimini iyileştirir
// / - Tekrar kullanılabilir görünümler sağlar
// / - Performans optimizasyonu yapılabilir
// /
// / ## Notlar
// /
// / - Lens adları benzersiz olmalıdır
// / - Query fonksiyonu nil olmamalıdır
// / - Lens'ler menüde ayrı sekmeler olarak gösterilir
func (r Base) Lenses() []Lens {
	return []Lens{}
}

// / # Policy Metodu
// /
// / Bu fonksiyon, resource için yetkilendirme politikasını döndürür.
// / CRUD işlemleri için erişim kontrolü sağlar.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Erişim Kontrolü**: Kullanıcı bazlı işlem yetkilendirmesi
// / 2. **Rol Tabanlı Yetkilendirme**: Rol bazlı erişim kontrolü
// / 3. **Kayıt Bazlı Yetkilendirme**: Belirli kayıtlar için özel yetkiler
// / 4. **Alan Bazlı Yetkilendirme**: Belirli alanları gizleme/gösterme
// /
// / ## Döndürür
// /
// / - `auth.Policy`: Yetkilendirme politikası interface'i
// /
// / ## Örnek Kullanım
// /
// / ```go
// / type UserPolicy struct{}
// /
// / func (p *UserPolicy) ViewAny(ctx *appContext.Context) bool {
// /     return ctx.User.HasPermission("users.view")
// / }
// /
// / func (p *UserPolicy) View(ctx *appContext.Context, model any) bool {
// /     user := model.(*models.User)
// /     return ctx.User.ID == user.ID || ctx.User.IsAdmin()
// / }
// /
// / func (p *UserPolicy) Create(ctx *appContext.Context) bool {
// /     return ctx.User.HasPermission("users.create")
// / }
// /
// / func (p *UserPolicy) Update(ctx *appContext.Context, model any) bool {
// /     user := model.(*models.User)
// /     return ctx.User.ID == user.ID || ctx.User.IsAdmin()
// / }
// /
// / func (p *UserPolicy) Delete(ctx *appContext.Context, model any) bool {
// /     return ctx.User.IsAdmin()
// / }
// / ```
// /
// / ## Policy Metodları
// /
// / - **ViewAny**: Liste görünümü erişimi
// / - **View**: Detay görünümü erişimi
// / - **Create**: Oluşturma erişimi
// / - **Update**: Güncelleme erişimi
// / - **Delete**: Silme erişimi
// / - **Restore**: Geri yükleme erişimi (soft delete)
// / - **ForceDelete**: Kalıcı silme erişimi
// /
// / ## Güvenlik Notları
// /
// / - Policy tanımlanmazsa varsayılan olarak tüm işlemler reddedilir
// / - Her metod mutlaka implement edilmelidir
// / - Güvenlik açıklarına karşı dikkatli olun
// / - Kullanıcı girişlerini her zaman doğrulayın
// /
// / ## Önemli Uyarılar
// /
// / - Policy nil olmamalıdır
// / - Policy metodları panic fırlatmamalıdır
// / - Context her zaman kontrol edilmelidir
// / - Performans için cache kullanılabilir
func (r Base) Policy() auth.Policy {
	return r.PolicyVal
}

// / # GetSortable Metodu
// /
// / Bu fonksiyon, resource için varsayılan sıralama ayarlarını döndürür.
// / Liste görünümünde hangi alanlara göre sıralama yapılabileceğini belirtir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Liste Sıralama**: Kullanıcının verileri sıralamasına izin verme
// / 2. **Varsayılan Sıralama**: Sayfa yüklendiğinde varsayılan sıralama
// / 3. **Çoklu Sıralama**: Birden fazla alana göre sıralama
// /
// / ## Döndürür
// /
// / - `[]Sortable`: Sıralama yapılandırmaları dizisi
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Base yapısında tanımlama
// / resource := &Base{
// /     Sortable: []Sortable{
// /         {Field: "created_at", Direction: "desc", Default: true},
// /         {Field: "name", Direction: "asc"},
// /         {Field: "email", Direction: "asc"},
// /     },
// / }
// /
// / // Kullanım
// / sortable := resource.GetSortable()
// / for _, sort := range sortable {
// /     if sort.Default {
// /         db = db.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
// /     }
// / }
// / ```
// /
// / ## Sortable Özellikleri
// /
// / - **Field**: Sıralanacak alan adı (veritabanı kolonu)
// / - **Direction**: Sıralama yönü ("asc" veya "desc")
// / - **Default**: Varsayılan sıralama olup olmadığı
// /
// / ## Notlar
// /
// / - Alan adları veritabanı kolonları ile eşleşmelidir
// / - Sadece bir alan varsayılan olarak işaretlenmelidir
// / - İlişkili alanlarda nokta notasyonu kullanılabilir: "User.Name"
func (r Base) GetSortable() []Sortable {
	return r.Sortable
}

// / # GetDialogType Metodu
// /
// / Bu fonksiyon, oluşturma ve düzenleme işlemleri için kullanılacak diyalog tipini döndürür.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Modal Diyalog**: Küçük formlar için popup modal
// / 2. **Slide Over**: Yan taraftan açılan panel
// / 3. **Full Page**: Tam sayfa form
// / 4. **Drawer**: Alt veya üstten açılan çekmece
// /
// / ## Döndürür
// /
// / - `DialogType`: Diyalog tipi enum değeri
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := &Base{
// /     DialogType: DialogTypeModal,
// / }
// /
// / dialogType := resource.GetDialogType()
// / // Frontend'de kullanım için
// / ```
// /
// / ## Diyalog Tipleri
// /
// / - **DialogTypeModal**: Merkezi popup modal (varsayılan)
// / - **DialogTypeSlideOver**: Sağdan açılan panel
// / - **DialogTypeFullPage**: Tam sayfa görünümü
// / - **DialogTypeDrawer**: Alt/üst çekmece
// /
// / ## Seçim Kriterleri
// /
// / - **Modal**: Basit, az alanlı formlar için
// / - **SlideOver**: Orta karmaşıklıkta formlar için
// / - **FullPage**: Çok alanlı, karmaşık formlar için
// / - **Drawer**: Hızlı işlemler için
func (r Base) GetDialogType() DialogType {
	return r.DialogType
}

// / # SetDialogType Metodu
// /
// / Bu fonksiyon, diyalog tipini ayarlar ve yapılandırılmış resource'ı döndürür.
// / Method chaining için kullanılabilir.
// /
// / ## Parametreler
// /
// / - `dialogType`: Ayarlanacak diyalog tipi
// /
// / ## Döndürür
// /
// / - Yapılandırılmış Resource pointer'ı (method chaining için)
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := NewUserResource().
// /     SetDialogType(DialogTypeSlideOver)
// /
// / // Veya
// / resource.SetDialogType(DialogTypeModal)
// / ```
// /
// / ## Method Chaining
// /
// / Bu metod method chaining pattern'ini destekler:
// /
// / ```go
// / resource := NewUserResource().
// /     SetDialogType(DialogTypeSlideOver).
// /     // Diğer yapılandırma metodları...
// / ```
// /
// / ## Notlar
// /
// / - Method chaining için Resource interface'i döndürür
// / - Çalışma zamanında diyalog tipini değiştirmek için kullanılabilir
func (r *Base) SetDialogType(dialogType DialogType) Resource {
	r.DialogType = dialogType
	return r
}

// / # OpenAPIEnabled Metodu
// /
// / Bu fonksiyon, kaynağın OpenAPI spesifikasyonunda görünüp görünmeyeceğini döner.
// / Varsayılan olarak tüm kaynaklar OpenAPI'de görünür (true).
// /
// / ## Kullanım Senaryoları
// /
// / 1. **API Dokümantasyonu**: Hangi endpoint'lerin dokümante edileceğini kontrol etmek
// / 2. **Internal API'ler**: Dahili kullanım için API'leri gizlemek
// / 3. **Beta Özellikler**: Henüz hazır olmayan özellikleri gizlemek
// / 4. **Versiyonlama**: Eski API versiyonlarını gizlemek
// /
// / ## Döndürür
// /
// / - true: OpenAPI spec'te görünür (varsayılan)
// / - false: OpenAPI spec'te gizli
// /
// / ## Örnek Kullanım
// /
// /	enabled := r.OpenAPIEnabled()
func (b Base) OpenAPIEnabled() bool {
	return !b.openAPIDisabled
}

/// # SetOpenAPIEnabled Metodu
///
/// Bu fonksiyon, kaynağın OpenAPI görünürlüğünü ayarlar.
/// Method chaining desteği için Resource interface'i döndürür.
///
/// ## Parametreler
///
/// - enabled: true = OpenAPI'de görünür, false = OpenAPI'de gizli
///
/// ## Döndürür
///
/// - Resource pointer'ı (method chaining için)
///
/// ## Örnek Kullanım
///
///	r.SetOpenAPIEnabled(false) // OpenAPI'de gizle
///	r.SetOpenAPIEnabled(true)  // OpenAPI'de göster

// SetTitleFunc, resource'un başlığını dinamik olarak ayarlamak için bir fonksiyon belirler.
//
// Bu metod, i18n desteği için kullanılır. Başlık, kullanıcının diline göre
// otomatik olarak çevrilir.
//
// ## Parametreler
//
// - fn: Başlık döndüren fonksiyon (fiber.Ctx alır, string döndürür)
//
// ## Döndürür
//
// - Resource pointer'ı (method chaining için)
//
// ## Örnek Kullanım
//
//	r.SetTitleFunc(func(c *fiber.Ctx) string {
//	    return i18n.Trans(c, "resources.users.title")
//	})
func (b *Base) SetTitleFunc(fn func(*fiber.Ctx) string) Resource {
	b.titleFunc = fn
	return b
}

// SetGroupFunc, resource'un grubunu dinamik olarak ayarlamak için bir fonksiyon belirler.
//
// Bu metod, i18n desteği için kullanılır. Grup, kullanıcının diline göre
// otomatik olarak çevrilir.
//
// ## Parametreler
//
// - fn: Grup adı döndüren fonksiyon (fiber.Ctx alır, string döndürür)
//
// ## Döndürür
//
// - Resource pointer'ı (method chaining için)
//
// ## Örnek Kullanım
//
//	r.SetGroupFunc(func(c *fiber.Ctx) string {
//	    return i18n.Trans(c, "resources.groups.system")
//	})
func (b *Base) SetGroupFunc(fn func(*fiber.Ctx) string) Resource {
	b.groupFunc = fn
	return b
}

func (b *Base) SetOpenAPIEnabled(enabled bool) Resource {
	b.openAPIDisabled = !enabled
	return b
}

// / # Repository Metodu
// /
// / Bu fonksiyon, resource için özel veri sağlayıcı (data provider) döndürür.
// / Varsayılan implementasyon nil döndürür, özel repository kullanımı için override edilmelidir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Özel Sorgu Mantığı**: Karmaşık veritabanı sorguları için
// / 2. **Cache Entegrasyonu**: Veri cache mekanizması eklemek için
// / 3. **Çoklu Veritabanı**: Farklı veritabanlarından veri çekmek için
// / 4. **API Entegrasyonu**: Harici API'lerden veri almak için
// / 5. **Performans Optimizasyonu**: Özel indeksleme ve sorgu optimizasyonu
// /
// / ## Parametreler
// /
// / - `db`: GORM veritabanı bağlantısı
// /
// / ## Döndürür
// /
// / - `data.DataProvider`: Özel veri sağlayıcı interface'i (varsayılan: nil)
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel repository implementasyonu
// / type UserRepository struct {
// /     db    *gorm.DB
// /     cache *redis.Client
// / }
// /
// / func (r *UserRepository) Find(ctx context.Context, id any) (any, error) {
// /     // Cache'den kontrol et
// /     if cached, err := r.cache.Get(ctx, fmt.Sprintf("user:%v", id)).Result(); err == nil {
// /         return cached, nil
// /     }
// /
// /     // Veritabanından çek
// /     var user models.User
// /     if err := r.db.First(&user, id).Error; err != nil {
// /         return nil, err
// /     }
// /
// /     // Cache'e kaydet
// /     r.cache.Set(ctx, fmt.Sprintf("user:%v", id), user, time.Hour)
// /     return &user, nil
// / }
// /
// / // Resource'da kullanım
// / func (r UserResource) Repository(db *gorm.DB) data.DataProvider {
// /     return &UserRepository{
// /         db:    db,
// /         cache: redisClient,
// /     }
// / }
// / ```
// /
// / ## DataProvider Interface Metodları
// /
// / - **Find**: Tek kayıt getirme
// / - **FindAll**: Tüm kayıtları getirme
// / - **Create**: Yeni kayıt oluşturma
// / - **Update**: Kayıt güncelleme
// / - **Delete**: Kayıt silme
// / - **Count**: Toplam kayıt sayısı
// /
// / ## Avantajlar
// /
// / - **Esneklik**: Veri kaynağından bağımsız çalışma
// / - **Test Edilebilirlik**: Mock repository ile kolay test
// / - **Separation of Concerns**: İş mantığını veri erişiminden ayırma
// / - **Yeniden Kullanılabilirlik**: Repository'leri farklı resource'larda kullanma
// /
// / ## Notlar
// /
// / - nil döndürülürse varsayılan GORM provider kullanılır
// / - Repository thread-safe olmalıdır
// / - Context timeout'larına dikkat edilmelidir
func (r Base) Repository(client *gorm.DB) data.DataProvider {
	return nil
}

// / # Cards Metodu
// /
// / Bu fonksiyon, resource için tanımlı dashboard widget'larını (card'ları) döndürür.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Dashboard Metrikleri**: Toplam kullanıcı, sipariş sayısı gibi metrikler
// / 2. **Grafikler**: Zaman serisi grafikleri, pasta grafikleri
// / 3. **İstatistikler**: Özet istatistikler ve KPI'lar
// / 4. **Hızlı Erişim**: Sık kullanılan işlemler için kısayollar
// /
// / ## Döndürür
// /
// / - `[]widget.Card`: Widget tanımlamaları dizisi
// /
// / ## Örnek Kullanım
// /
// / ```go
// / resource := &Base{
// /     WidgetsVal: []widget.Card{
// /         widget.NewMetricCard("total_users", "Toplam Kullanıcı").
// /             SetValue(func(ctx *appContext.Context) any {
// /                 var count int64
// /                 ctx.DB.Model(&models.User{}).Count(&count)
// /                 return count
// /             }).
// /             SetIcon("users"),
// /
// /         widget.NewChartCard("user_growth", "Kullanıcı Artışı").
// /             SetChartType("line").
// /             SetData(func(ctx *appContext.Context) any {
// /                 // Grafik verilerini hazırla
// /                 return chartData
// /             }),
// /     },
// / }
// /
// / cards := resource.Cards()
// / ```
// /
// / ## Widget Tipleri
// /
// / - **MetricCard**: Tek sayısal değer gösterimi
// / - **ChartCard**: Grafik gösterimi (line, bar, pie, doughnut)
// / - **TableCard**: Tablo formatında veri
// / - **ListCard**: Liste formatında veri
// / - **CustomCard**: Özel HTML/React bileşeni
// /
// / ## Widget Özellikleri
// /
// / - **Title**: Widget başlığı
// / - **Value**: Gösterilecek değer (dinamik)
// / - **Icon**: Widget ikonu
// / - **Color**: Renk teması
// / - **Size**: Widget boyutu (small, medium, large)
// / - **Refresh**: Otomatik yenileme süresi
// /
// / ## Performans Notları
// /
// / - Widget verileri cache'lenebilir
// / - Ağır sorgular için background job kullanın
// / - Lazy loading ile performans artırılabilir
// /
// / ## Notlar
// /
// / - Widget'lar dashboard sayfasında gösterilir
// / - Her widget benzersiz ID'ye sahip olmalıdır
// / - Widget verileri asenkron yüklenebilir
func (r Base) Cards() []widget.Card {
	if r.WidgetsVal != nil {
		// Convert generic WidgetsVal to []widget.Card if necessary,
		// but since we updated the interface, we should update WidgetsVal type too.
		// For now, let's assume WidgetsVal holds []widget.Card
		return r.WidgetsVal
	}
	return []widget.Card{}
}

// / # StoreHandler Metodu
// /
// / Bu fonksiyon, dosya yükleme işlemlerini yönetir.
// / Özel yükleme mantığı tanımlanmamışsa varsayılan yerel depolama kullanır.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Yerel Depolama**: Dosyaları sunucu diskine kaydetme
// / 2. **Cloud Storage**: AWS S3, Google Cloud Storage, Azure Blob
// / 3. **CDN Entegrasyonu**: Dosyaları CDN'e yükleme
// / 4. **Görsel İşleme**: Yüklenen görselleri optimize etme, thumbnail oluşturma
// / 5. **Virus Tarama**: Yüklenen dosyaları güvenlik kontrolünden geçirme
// /
// / ## Parametreler
// /
// / - `c`: Uygulama context'i (kullanıcı, request bilgileri)
// / - `file`: Yüklenecek dosya header'ı
// / - `storagePath`: Dosyaların kaydedileceği yerel dizin (varsayılan: "./storage/public")
// / - `storageURL`: Dosyalara erişim için public URL prefix'i (varsayılan: "/storage")
// /
// / ## Döndürür
// /
// / - `string`: Yüklenen dosyanın public URL'i
// / - `error`: Hata durumunda hata mesajı
// /
// / ## Örnek Kullanım
// /
// / ### Varsayılan Kullanım (Yerel Depolama)
// /
// / ```go
// / resource := &Base{
// /     // UploadHandler tanımlanmamış, varsayılan kullanılır
// / }
// /
// / url, err := resource.StoreHandler(ctx, fileHeader, "./storage/public", "/storage")
// / // Dosya: ./storage/public/1234567890.jpg
// / // URL: /storage/1234567890.jpg
// / ```
// /
// / ### Özel S3 Yükleme
// /
// / ```go
// / resource := &Base{
// /     UploadHandler: func(c *appContext.Context, file *multipart.FileHeader) (string, error) {
// /         // S3'e yükle
// /         s3Client := getS3Client()
// /
// /         // Dosyayı aç
// /         src, err := file.Open()
// /         if err != nil {
// /             return "", err
// /         }
// /         defer src.Close()
// /
// /         // S3'e yükle
// /         key := fmt.Sprintf("uploads/%d/%s", time.Now().Unix(), file.Filename)
// /         _, err = s3Client.PutObject(&s3.PutObjectInput{
// /             Bucket: aws.String("my-bucket"),
// /             Key:    aws.String(key),
// /             Body:   src,
// /             ACL:    aws.String("public-read"),
// /         })
// /
// /         if err != nil {
// /             return "", err
// /         }
// /
// /         // Public URL döndür
// /         return fmt.Sprintf("https://cdn.example.com/%s", key), nil
// /     },
// / }
// / ```
// /
// / ### Görsel Optimizasyonu ile Yükleme
// /
// / ```go
// / resource := &Base{
// /     UploadHandler: func(c *appContext.Context, file *multipart.FileHeader) (string, error) {
// /         // Görseli aç
// /         src, _ := file.Open()
// /         defer src.Close()
// /
// /         // Görseli decode et
// /         img, _, err := image.Decode(src)
// /         if err != nil {
// /             return "", err
// /         }
// /
// /         // Resize et
// /         resized := resize.Resize(800, 0, img, resize.Lanczos3)
// /
// /         // Kaydet
// /         filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
// /         out, _ := os.Create(filepath.Join("./storage/public", filename))
// /         defer out.Close()
// /
// /         jpeg.Encode(out, resized, &jpeg.Options{Quality: 85})
// /
// /         return fmt.Sprintf("/storage/%s", filename), nil
// /     },
// / }
// / ```
// /
// / ## Varsayılan Davranış
// /
// / 1. Dosya adına timestamp eklenir (collision önleme)
// / 2. Dizin yoksa otomatik oluşturulur (0755 izinleri)
// / 3. Dosya uzantısı korunur
// / 4. Public URL döndürülür
// /
// / ## Güvenlik Notları
// /
// / - Dosya tipini her zaman doğrulayın
// / - Dosya boyutu limitlerini kontrol edin
// / - Dosya adlarını sanitize edin
// / - Virus taraması yapın (production'da)
// / - Kullanıcı bazlı upload limitleri uygulayın
// /
// / ## Önemli Uyarılar
// /
// / - Büyük dosyalar için streaming kullanın
// / - Disk alanını düzenli kontrol edin
// / - Yükleme işlemi sırasında timeout ayarlayın
// / - Başarısız yüklemelerde geçici dosyaları temizleyin
// /
// / ## Desteklenen Storage Çözümleri
// /
// / - **Yerel Disk**: Basit projeler için
// / - **AWS S3**: Ölçeklenebilir cloud storage
// / - **Google Cloud Storage**: Google Cloud entegrasyonu
// / - **Azure Blob Storage**: Microsoft Azure entegrasyonu
// / - **MinIO**: Self-hosted S3-compatible storage
// / - **Cloudinary**: Görsel optimizasyonu ve CDN
func (r Base) StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	if r.UploadHandler != nil {
		return r.UploadHandler(c, file)
	}

	// Default Storage Logic
	// Use configured storage path
	if storagePath == "" {
		storagePath = "./storage/public"
	}
	if storageURL == "" {
		storageURL = "/storage"
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	localPath := filepath.Join(storagePath, filename)

	// Ensure directory exists
	_ = os.MkdirAll(storagePath, 0755)

	if err := c.Ctx.SaveFile(file, localPath); err != nil {
		return "", err
	}
	// Public URL path
	// Ensure storageURL doesn't end with slash if filename starts with one, but filename is just name.
	// Handle slashes cleanly.
	return fmt.Sprintf("%s/%s", storageURL, filename), nil
}

// / # NavigationOrder Metodu
// /
// / Bu fonksiyon, resource'ın navigasyon menüsündeki sırasını belirler.
// / Düşük sayılar üstte, yüksek sayılar altta gösterilir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Menü Organizasyonu**: Önemli resource'ları üstte gösterme
// / 2. **Kullanıcı Deneyimi**: Sık kullanılan öğeleri erişilebilir yapma
// / 3. **Mantıksal Gruplama**: İlgili resource'ları yan yana gösterme
// /
// / ## Döndürür
// /
// / - `int`: Sıralama değeri (varsayılan: 99)
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da override etme
// / func (r DashboardResource) NavigationOrder() int {
// /     return 1 // En üstte göster
// / }
// /
// / func (r UserResource) NavigationOrder() int {
// /     return 10 // Üst sıralarda
// / }
// /
// / func (r SettingsResource) NavigationOrder() int {
// /     return 100 // Alt sıralarda
// / }
// / ```
// /
// / ## Sıralama Önerileri
// /
// / - **1-10**: Ana özellikler (Dashboard, Ana Sayfa)
// / - **11-30**: Sık kullanılan resource'lar (Kullanıcılar, İçerik)
// / - **31-60**: Orta öncelikli resource'lar (Raporlar, Analitik)
// / - **61-90**: Düşük öncelikli resource'lar (Loglar, Arşiv)
// / - **91-100**: Ayarlar ve yönetim
// /
// / ## Notlar
// /
// / - Varsayılan değer 99'dur
// / - Aynı sıra değerine sahip resource'lar alfabetik sıralanır
// / - Grup içinde sıralama yapılır
func (r Base) NavigationOrder() int {
	return 99
}

// / # Visible Metodu
// /
// / Bu fonksiyon, resource'ın navigasyon menüsünde görünür olup olmadığını belirler.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **API-Only Resource'lar**: Sadece API üzerinden erişilebilir resource'lar
// / 2. **Gizli Yönetim**: Admin panelinde gösterilmeyecek resource'lar
// / 3. **Koşullu Görünürlük**: Kullanıcı rolüne göre gösterim
// / 4. **Geliştirme Aşaması**: Henüz hazır olmayan resource'ları gizleme
// /
// / ## Döndürür
// /
// / - `bool`: Görünürlük durumu (varsayılan: true)
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Her zaman gizli resource
// / func (r InternalLogResource) Visible() bool {
// /     return false
// / }
// /
// / // Koşullu görünürlük
// / func (r AdminResource) Visible() bool {
// /     // Sadece admin kullanıcılar için görünür
// /     return currentUser.IsAdmin()
// / }
// /
// / // Geliştirme ortamında görünür
// / func (r BetaFeatureResource) Visible() bool {
// /     return os.Getenv("APP_ENV") == "development"
// / }
// / ```
// /
// / ## Kullanım Alanları
// /
// / - **API Endpoints**: Menüde gösterilmez ama API'den erişilebilir
// / - **Webhook Handlers**: Harici sistemler için endpoint'ler
// / - **Internal Tools**: Sadece sistem tarafından kullanılan resource'lar
// / - **Feature Flags**: Özellik açma/kapama için
// /
// / ## Notlar
// /
// / - false döndürülürse menüde gösterilmez
// / - API endpoint'leri hala erişilebilir olur
// / - Policy kontrolü yine de uygulanır
// / - Breadcrumb ve diğer UI öğelerinde de gizlenir
func (r Base) Visible() bool {
	return true
}

// / # GetFields Metodu
// /
// / Bu fonksiyon, belirli bir bağlam (context) içinde gösterilecek alanları döndürür.
// / Context'e göre dinamik alan filtreleme ve özelleştirme yapılabilir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Rol Bazlı Alan Gösterimi**: Kullanıcı rolüne göre farklı alanlar gösterme
// / 2. **Koşullu Alan Gösterimi**: Belirli koşullara göre alanları gizleme/gösterme
// / 3. **Dinamik Form Oluşturma**: Context'e göre form alanlarını değiştirme
// / 4. **İzin Bazlı Filtreleme**: Kullanıcının görebileceği alanları filtreleme
// /
// / ## Parametreler
// /
// / - `ctx`: Uygulama context'i (kullanıcı, request, veritabanı bilgileri içerir)
// /
// / ## Döndürür
// /
// / - `[]fields.Element`: Context'e göre filtrelenmiş alan tanımlamaları
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da override etme
// / func (r UserResource) GetFields(ctx *appContext.Context) []fields.Element {
// /     fields := r.FieldsVal
// /
// /     // Admin olmayan kullanıcılar için hassas alanları gizle
// /     if !ctx.User.IsAdmin() {
// /         filtered := []fields.Element{}
// /         for _, field := range fields {
// /             if field.GetKey() != "password" && field.GetKey() != "api_token" {
// /                 filtered = append(filtered, field)
// /             }
// /         }
// /         return filtered
// /     }
// /
// /     return fields
// / }
// /
// / // Kullanım
// / resource := NewUserResource()
// / fields := resource.GetFields(ctx)
// / ```
// /
// / ## Gelişmiş Örnekler
// /
// / ### Rol Bazlı Alan Filtreleme
// /
// / ```go
// / func (r OrderResource) GetFields(ctx *appContext.Context) []fields.Element {
// /     fields := r.FieldsVal
// /
// /     switch ctx.User.Role {
// /     case "customer":
// /         // Müşteriler sadece kendi siparişlerini görebilir
// /         return filterFields(fields, []string{"id", "items", "total", "status"})
// /     case "staff":
// /         // Personel daha fazla alan görebilir
// /         return filterFields(fields, []string{"id", "customer", "items", "total", "status", "notes"})
// /     case "admin":
// /         // Admin tüm alanları görebilir
// /         return fields
// /     default:
// /         return []fields.Element{}
// /     }
// / }
// / ```
// /
// / ### Dinamik Alan Ekleme
// /
// / ```go
// / func (r ProductResource) GetFields(ctx *appContext.Context) []fields.Element {
// /     fields := r.FieldsVal
// /
// /     // Premium kullanıcılar için ek alanlar ekle
// /     if ctx.User.IsPremium() {
// /         fields = append(fields,
// /             fields.NewText("internal_notes").SetLabel("Dahili Notlar"),
// /             fields.NewNumber("cost_price").SetLabel("Maliyet Fiyatı"),
// /         )
// /     }
// /
// /     return fields
// / }
// / ```
// /
// / ## Fields() ile Farkı
// /
// / - **Fields()**: Tüm alanları döndürür, context'ten bağımsız
// / - **GetFields()**: Context'e göre filtrelenmiş alanları döndürür
// /
// / ## Requirement 11.1
// /
// / Resource arayüzünü, alanları almak için metodlar içerecek şekilde genişletir.
// / Context-aware alan yönetimi sağlar.
// /
// / ## Notlar
// /
// / - Varsayılan implementasyon tüm alanları döndürür
// / - Özel filtreleme için override edilmelidir
// / - Context nil kontrolü yapılmalıdır
// / - Performans için cache kullanılabilir
func (r Base) GetFields(ctx *appContext.Context) []fields.Element {
	return r.FieldsVal
}

// / # GetCards Metodu
// /
// / Bu fonksiyon, belirli bir bağlam (context) içinde gösterilecek dashboard widget'larını döndürür.
// / Context'e göre dinamik widget filtreleme ve özelleştirme yapılabilir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Rol Bazlı Widget Gösterimi**: Kullanıcı rolüne göre farklı widget'lar gösterme
// / 2. **Koşullu Widget Gösterimi**: Belirli koşullara göre widget'ları gizleme/gösterme
// / 3. **Dinamik Dashboard**: Context'e göre dashboard içeriğini değiştirme
// / 4. **İzin Bazlı Filtreleme**: Kullanıcının görebileceği metrikleri filtreleme
// /
// / ## Parametreler
// /
// / - `ctx`: Uygulama context'i (kullanıcı, request, veritabanı bilgileri içerir)
// /
// / ## Döndürür
// /
// / - `[]widget.Card`: Context'e göre filtrelenmiş widget tanımlamaları
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da override etme
// / func (r UserResource) GetCards(ctx *appContext.Context) []widget.Card {
// /     cards := r.WidgetsVal
// /
// /     // Admin olmayan kullanıcılar için hassas widget'ları gizle
// /     if !ctx.User.IsAdmin() {
// /         filtered := []widget.Card{}
// /         for _, card := range cards {
// /             // Finansal widget'ları sadece admin görebilir
// /             if card.GetID() != "revenue" && card.GetID() != "profit" {
// /                 filtered = append(filtered, card)
// /             }
// /         }
// /         return filtered
// /     }
// /
// /     return cards
// / }
// /
// / // Kullanım
// / resource := NewUserResource()
// / cards := resource.GetCards(ctx)
// / ```
// /
// / ## Gelişmiş Örnekler
// /
// / ### Rol Bazlı Widget Filtreleme
// /
// / ```go
// / func (r SalesResource) GetCards(ctx *appContext.Context) []widget.Card {
// /     allCards := r.WidgetsVal
// /
// /     switch ctx.User.Role {
// /     case "sales_rep":
// /         // Satış temsilcileri sadece kendi metriklerini görebilir
// /         return []widget.Card{
// /             widget.NewMetricCard("my_sales", "Satışlarım"),
// /             widget.NewMetricCard("my_commission", "Komisyonum"),
// /         }
// /     case "sales_manager":
// /         // Satış müdürleri takım metriklerini görebilir
// /         return []widget.Card{
// /             widget.NewMetricCard("team_sales", "Takım Satışları"),
// /             widget.NewChartCard("team_performance", "Takım Performansı"),
// /         }
// /     case "admin":
// /         // Admin tüm widget'ları görebilir
// /         return allCards
// /     default:
// /         return []widget.Card{}
// /     }
// / }
// / ```
// /
// / ### Dinamik Widget Ekleme
// /
// / ```go
// / func (r AnalyticsResource) GetCards(ctx *appContext.Context) []widget.Card {
// /     cards := r.WidgetsVal
// /
// /     // Premium kullanıcılar için ek widget'lar ekle
// /     if ctx.User.IsPremium() {
// /         cards = append(cards,
// /             widget.NewChartCard("advanced_analytics", "Gelişmiş Analitik").
// /                 SetChartType("line").
// /                 SetData(func(c *appContext.Context) any {
// /                     return getAdvancedAnalytics(c)
// /                 }),
// /             widget.NewMetricCard("predictive_score", "Tahmin Skoru"),
// /         )
// /     }
// /
// /     return cards
// / }
// / ```
// /
// / ### Zaman Bazlı Widget Gösterimi
// /
// / ```go
// / func (r DashboardResource) GetCards(ctx *appContext.Context) []widget.Card {
// /     cards := r.WidgetsVal
// /     now := time.Now()
// /
// /     // Çalışma saatleri dışında bazı widget'ları gizle
// /     if now.Hour() < 9 || now.Hour() > 18 {
// /         filtered := []widget.Card{}
// /         for _, card := range cards {
// /             if card.GetID() != "live_support" {
// /                 filtered = append(filtered, card)
// /             }
// /         }
// /         return filtered
// /     }
// /
// /     return cards
// / }
// / ```
// /
// / ## Cards() ile Farkı
// /
// / - **Cards()**: Tüm widget'ları döndürür, context'ten bağımsız
// / - **GetCards()**: Context'e göre filtrelenmiş widget'ları döndürür
// /
// / ## Requirement 11.1
// /
// / Resource arayüzünü, card'ları almak için metodlar içerecek şekilde genişletir.
// / Context-aware widget yönetimi sağlar.
// /
// / ## Performans Notları
// /
// / - Widget verileri lazy load edilebilir
// / - Ağır hesaplamalar cache'lenmelidir
// / - Gereksiz widget'lar filtrelenerek performans artırılabilir
// /
// / ## Notlar
// /
// / - Varsayılan implementasyon tüm widget'ları döndürür
// / - Özel filtreleme için override edilmelidir
// / - Context nil kontrolü yapılmalıdır
// / - Widget ID'leri benzersiz olmalıdır
func (r Base) GetCards(ctx *appContext.Context) []widget.Card {
	if r.WidgetsVal != nil {
		return r.WidgetsVal
	}
	return []widget.Card{}
}

// / # GetLenses Metodu
// /
// / Bu fonksiyon, kaynağın tüm lens'lerini (özel görünümlerini) döndürür.
// / Lens'ler, önceden tanımlı filtre kombinasyonları ve özel veri görünümleridir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Hızlı Filtreler**: Sık kullanılan filtre kombinasyonlarına hızlı erişim
// / 2. **Özel Raporlar**: Belirli veri alt kümelerini görüntüleme
// / 3. **Kullanıcı Segmentasyonu**: Farklı kullanıcı gruplarını görüntüleme
// / 4. **Durum Bazlı Görünümler**: Kayıt durumlarına göre filtreleme
// /
// / ## Döndürür
// /
// / - `[]Lens`: Tüm lens tanımlamaları dizisi
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da lens tanımlama
// / func (r UserResource) Lenses() []Lens {
// /     return []Lens{
// /         {
// /             Name:  "active",
// /             Label: "Aktif Kullanıcılar",
// /             Icon:  "check-circle",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("status = ?", "active").
// /                     Where("deleted_at IS NULL")
// /             },
// /         },
// /         {
// /             Name:  "inactive",
// /             Label: "Pasif Kullanıcılar",
// /             Icon:  "x-circle",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("status = ?", "inactive")
// /             },
// /         },
// /         {
// /             Name:  "premium",
// /             Label: "Premium Üyeler",
// /             Icon:  "star",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("subscription_type = ?", "premium").
// /                     Where("subscription_expires_at > ?", time.Now())
// /             },
// /         },
// /     }
// / }
// /
// / // Kullanım
// / resource := NewUserResource()
// / lenses := resource.GetLenses()
// / ```
// /
// / ## Gelişmiş Örnekler
// /
// / ### E-ticaret Sipariş Lens'leri
// /
// / ```go
// / func (r OrderResource) Lenses() []Lens {
// /     return []Lens{
// /         {
// /             Name:  "pending",
// /             Label: "Bekleyen Siparişler",
// /             Icon:  "clock",
// /             Badge: func(db *gorm.DB) int64 {
// /                 var count int64
// /                 db.Model(&Order{}).Where("status = ?", "pending").Count(&count)
// /                 return count
// /             },
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("status = ?", "pending").
// /                     Order("created_at ASC")
// /             },
// /         },
// /         {
// /             Name:  "processing",
// /             Label: "İşleniyor",
// /             Icon:  "refresh",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("status = ?", "processing")
// /             },
// /         },
// /         {
// /             Name:  "completed",
// /             Label: "Tamamlanan",
// /             Icon:  "check",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("status = ?", "completed").
// /                     Order("completed_at DESC")
// /             },
// /         },
// /         {
// /             Name:  "high-value",
// /             Label: "Yüksek Değerli",
// /             Icon:  "currency-dollar",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("total >= ?", 1000).
// /                     Order("total DESC")
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ### Zaman Bazlı Lens'ler
// /
// / ```go
// / func (r ArticleResource) Lenses() []Lens {
// /     return []Lens{
// /         {
// /             Name:  "today",
// /             Label: "Bugün Yayınlanan",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 today := time.Now().Truncate(24 * time.Hour)
// /                 return db.Where("published_at >= ?", today)
// /             },
// /         },
// /         {
// /             Name:  "this-week",
// /             Label: "Bu Hafta",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
// /                 return db.Where("published_at >= ?", weekStart)
// /             },
// /         },
// /         {
// /             Name:  "popular",
// /             Label: "Popüler",
// /             Query: func(db *gorm.DB) *gorm.DB {
// /                 return db.Where("views > ?", 1000).
// /                     Order("views DESC")
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ## Lens Özellikleri
// /
// / - **Name**: Benzersiz tanımlayıcı (URL'de kullanılır)
// / - **Label**: Kullanıcı arayüzünde gösterilecek başlık
// / - **Icon**: Lens ikonu (opsiyonel)
// / - **Query**: GORM sorgu modifikasyonu fonksiyonu
// / - **Badge**: Kayıt sayısını gösteren badge (opsiyonel)
// / - **Description**: Lens açıklaması (opsiyonel)
// /
// / ## Lenses() ile Farkı
// /
// / - **Lenses()**: Lens'leri tanımlar (override edilebilir)
// / - **GetLenses()**: Tanımlı lens'leri döndürür (wrapper metod)
// /
// / ## Requirement 11.1
// /
// / Resource arayüzünü, lens'leri almak için metodlar içerecek şekilde genişletir.
// / Özel görünüm yönetimi sağlar.
// /
// / ## Kullanıcı Deneyimi
// /
// / - Lens'ler UI'da sekmeler olarak gösterilir
// / - Her lens için ayrı URL oluşturulur
// / - Badge'ler gerçek zamanlı güncellenebilir
// / - Lens seçimi URL'de saklanır
// /
// / ## Performans Notları
// /
// / - Query fonksiyonları verimli olmalıdır
// / - Badge hesaplamaları cache'lenebilir
// / - Ağır sorgular için indeks kullanın
// / - Gereksiz JOIN'lerden kaçının
// /
// / ## Notlar
// /
// / - Varsayılan implementasyon Lenses() metodunu çağırır
// / - Lens adları benzersiz olmalıdır
// / - Query fonksiyonu nil olmamalıdır
// / - Lens'ler alfabetik sıralanabilir
func (r Base) GetLenses() []Lens {
	return r.Lenses()
}

// / # GetPolicy Metodu
// /
// / Bu fonksiyon, kaynağın yetkilendirme politikasını döndürür.
// / Policy, CRUD işlemleri için erişim kontrolü sağlar.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Erişim Kontrolü**: Kullanıcının işlem yapma yetkisini kontrol etme
// / 2. **Middleware Entegrasyonu**: HTTP middleware'de yetki kontrolü
// / 3. **UI Rendering**: Butonları ve menüleri koşullu gösterme
// / 4. **API Güvenliği**: API endpoint'lerini koruma
// /
// / ## Döndürür
// /
// / - `auth.Policy`: Yetkilendirme politikası interface'i
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Policy kontrolü
// / resource := NewUserResource()
// / policy := resource.GetPolicy()
// /
// / // Liste görüntüleme yetkisi kontrolü
// / if policy.ViewAny(ctx) {
// /     // Kullanıcı listeyi görebilir
// /     users := fetchUsers()
// /     return users
// / } else {
// /     return errors.New("yetkisiz erişim")
// / }
// /
// / // Kayıt düzenleme yetkisi kontrolü
// / user := findUser(id)
// / if policy.Update(ctx, user) {
// /     // Kullanıcı bu kaydı düzenleyebilir
// /     updateUser(user, data)
// / } else {
// /     return errors.New("bu kaydı düzenleme yetkiniz yok")
// / }
// / ```
// /
// / ## Gelişmiş Örnekler
// /
// / ### Middleware'de Policy Kullanımı
// /
// / ```go
// / func ResourceAuthMiddleware(resource Resource) fiber.Handler {
// /     return func(c *fiber.Ctx) error {
// /         ctx := appContext.FromFiber(c)
// /         policy := resource.GetPolicy()
// /
// /         // HTTP metoduna göre policy kontrolü
// /         switch c.Method() {
// /         case "GET":
// /             if c.Params("id") != "" {
// /                 // Detay görüntüleme
// /                 if !policy.View(ctx, nil) {
// /                     return c.Status(403).JSON(fiber.Map{
// /                         "error": "Bu kaydı görüntüleme yetkiniz yok",
// /                     })
// /                 }
// /             } else {
// /                 // Liste görüntüleme
// /                 if !policy.ViewAny(ctx) {
// /                     return c.Status(403).JSON(fiber.Map{
// /                         "error": "Liste görüntüleme yetkiniz yok",
// /                     })
// /                 }
// /             }
// /         case "POST":
// /             if !policy.Create(ctx) {
// /                 return c.Status(403).JSON(fiber.Map{
// /                     "error": "Oluşturma yetkiniz yok",
// /                 })
// /             }
// /         case "PUT", "PATCH":
// /             if !policy.Update(ctx, nil) {
// /                 return c.Status(403).JSON(fiber.Map{
// /                     "error": "Güncelleme yetkiniz yok",
// /                 })
// /             }
// /         case "DELETE":
// /             if !policy.Delete(ctx, nil) {
// /                 return c.Status(403).JSON(fiber.Map{
// /                     "error": "Silme yetkiniz yok",
// /                 })
// /             }
// /         }
// /
// /         return c.Next()
// /     }
// / }
// / ```
// /
// / ### UI'da Koşullu Rendering
// /
// / ```go
// / func RenderResourceActions(ctx *appContext.Context, resource Resource, item any) []Action {
// /     policy := resource.GetPolicy()
// /     actions := []Action{}
// /
// /     // Görüntüleme butonu
// /     if policy.View(ctx, item) {
// /         actions = append(actions, Action{
// /             Name:  "view",
// /             Label: "Görüntüle",
// /             Icon:  "eye",
// /         })
// /     }
// /
// /     // Düzenleme butonu
// /     if policy.Update(ctx, item) {
// /         actions = append(actions, Action{
// /             Name:  "edit",
// /             Label: "Düzenle",
// /             Icon:  "pencil",
// /         })
// /     }
// /
// /     // Silme butonu
// /     if policy.Delete(ctx, item) {
// /         actions = append(actions, Action{
// /             Name:  "delete",
// /             Label: "Sil",
// /             Icon:  "trash",
// /             Color: "danger",
// /         })
// /     }
// /
// /     return actions
// / }
// / ```
// /
// / ### Toplu İşlem Yetkilendirmesi
// /
// / ```go
// / func BulkDeleteHandler(ctx *appContext.Context, resource Resource, ids []uint) error {
// /     policy := resource.GetPolicy()
// /     db := ctx.DB
// /
// /     // Her kayıt için yetki kontrolü
// /     for _, id := range ids {
// /         var item any
// /         if err := db.First(&item, id).Error; err != nil {
// /             return err
// /         }
// /
// /         if !policy.Delete(ctx, item) {
// /             return fmt.Errorf("ID %d için silme yetkiniz yok", id)
// /         }
// /     }
// /
// /     // Tüm yetkiler onaylandı, silme işlemini gerçekleştir
// /     return db.Delete(resource.Model(), ids).Error
// / }
// / ```
// /
// / ## Policy Interface Metodları
// /
// / - **ViewAny(ctx)**: Liste görüntüleme yetkisi
// / - **View(ctx, model)**: Tek kayıt görüntüleme yetkisi
// / - **Create(ctx)**: Yeni kayıt oluşturma yetkisi
// / - **Update(ctx, model)**: Kayıt güncelleme yetkisi
// / - **Delete(ctx, model)**: Kayıt silme yetkisi
// / - **Restore(ctx, model)**: Soft delete geri yükleme yetkisi
// / - **ForceDelete(ctx, model)**: Kalıcı silme yetkisi
// /
// / ## Policy() ile Farkı
// /
// / - **Policy()**: Policy'yi döndürür (Base struct'tan)
// / - **GetPolicy()**: Policy'yi döndürür (wrapper metod, aynı işlevi görür)
// /
// / ## Requirement 11.1
// /
// / Resource arayüzünü, politikaları almak için metodlar içerecek şekilde genişletir.
// / Tutarlı yetkilendirme yönetimi sağlar.
// /
// / ## Güvenlik En İyi Uygulamaları
// /
// / 1. **Her İşlem İçin Kontrol**: Her CRUD işleminden önce policy kontrolü yapın
// / 2. **Kayıt Bazlı Kontrol**: Genel yetki yeterli değilse kayıt bazlı kontrol yapın
// / 3. **Fail-Safe**: Policy tanımlı değilse varsayılan olarak erişimi reddedin
// / 4. **Audit Logging**: Yetki reddedilmelerini loglayın
// / 5. **Context Validation**: Context'in geçerli olduğundan emin olun
// /
// / ## Performans Notları
// /
// / - Policy kontrolleri cache'lenebilir (dikkatli kullanın)
// / - Veritabanı sorguları minimize edilmelidir
// / - Toplu işlemlerde batch kontrol yapın
// /
// / ## Notlar
// /
// / - Varsayılan implementasyon PolicyVal alanını döndürür
// / - Policy nil olmamalıdır (güvenlik riski)
// / - Her metod implement edilmelidir
// / - Context her zaman geçerli olmalıdır
func (r Base) GetPolicy() auth.Policy {
	return r.PolicyVal
}

// / # ResolveField Metodu
// /
// / Bu fonksiyon, bir alanın değerini dinamik olarak hesaplar ve dönüştürür.
// / Alan değerlerini model'den çıkarır, işler ve serileştirir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Dinamik Değer Hesaplama**: Hesaplanmış alanlar için değer üretme
// / 2. **Veri Dönüşümü**: Model verilerini UI formatına dönüştürme
// / 3. **İlişki Çözümleme**: İlişkili verileri yükleme ve formatlama
// / 4. **Özel Formatlama**: Tarih, para birimi gibi özel formatlamalar
// / 5. **API Serileştirme**: API response'ları için veri hazırlama
// /
// / ## Parametreler
// /
// / - `fieldName`: Çözümlenecek alanın adı (key)
// / - `item`: Değerin çıkarılacağı model instance'ı
// /
// / ## Döndürür
// /
// / - `any`: Çözümlenmiş ve serileştirilmiş alan değeri
// / - `error`: Alan bulunamazsa veya çözümleme başarısız olursa hata
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Basit kullanım
// / resource := NewUserResource()
// / user := &models.User{
// /     ID:    1,
// /     Name:  "Ahmet Yılmaz",
// /     Email: "ahmet@example.com",
// / }
// /
// / // Alan değerini çözümle
// / value, err := resource.ResolveField("name", user)
// / if err != nil {
// /     log.Fatal(err)
// / }
// / fmt.Println(value) // "Ahmet Yılmaz"
// /
// / // E-posta alanını çözümle
// / email, err := resource.ResolveField("email", user)
// / fmt.Println(email) // "ahmet@example.com"
// / ```
// /
// / ## Gelişmiş Örnekler
// /
// / ### Hesaplanmış Alan Çözümleme
// /
// / ```go
// / // Özel resource'da override etme
// / func (r UserResource) ResolveField(fieldName string, item any) (any, error) {
// /     user := item.(*models.User)
// /
// /     switch fieldName {
// /     case "full_name":
// /         // Hesaplanmış alan: ad + soyad
// /         return fmt.Sprintf("%s %s", user.FirstName, user.LastName), nil
// /
// /     case "age":
// /         // Yaş hesaplama
// /         if user.BirthDate != nil {
// /             age := time.Now().Year() - user.BirthDate.Year()
// /             return age, nil
// /         }
// /         return nil, nil
// /
// /     case "status_label":
// /         // Durum etiketi
// /         statusLabels := map[string]string{
// /             "active":   "Aktif",
// /             "inactive": "Pasif",
// /             "banned":   "Yasaklı",
// /         }
// /         return statusLabels[user.Status], nil
// /
// /     case "avatar_url":
// /         // Avatar URL oluşturma
// /         if user.Avatar != "" {
// /             return fmt.Sprintf("https://cdn.example.com/avatars/%s", user.Avatar), nil
// /         }
// /         return "https://cdn.example.com/avatars/default.png", nil
// /
// /     default:
// /         // Varsayılan çözümleme
// /         return r.Base.ResolveField(fieldName, item)
// /     }
// / }
// / ```
// /
// / ### İlişki Çözümleme
// /
// / ```go
// / func (r OrderResource) ResolveField(fieldName string, item any) (any, error) {
// /     order := item.(*models.Order)
// /
// /     switch fieldName {
// /     case "customer_name":
// /         // İlişkili müşteri adı
// /         if order.Customer != nil {
// /             return order.Customer.Name, nil
// /         }
// /         return "Bilinmeyen Müşteri", nil
// /
// /     case "items_count":
// /         // Ürün sayısı
// /         return len(order.Items), nil
// /
// /     case "total_formatted":
// /         // Formatlanmış toplam
// /         return fmt.Sprintf("₺%.2f", order.Total), nil
// /
// /     case "status_badge":
// /         // Durum badge'i
// /         badges := map[string]map[string]string{
// /             "pending": {
// /                 "label": "Bekliyor",
// /                 "color": "yellow",
// /             },
// /             "completed": {
// /                 "label": "Tamamlandı",
// /                 "color": "green",
// /             },
// /             "cancelled": {
// /                 "label": "İptal Edildi",
// /                 "color": "red",
// /             },
// /         }
// /         return badges[order.Status], nil
// /
// /     default:
// /         return r.Base.ResolveField(fieldName, item)
// /     }
// / }
// / ```
// /
// / ### Toplu Alan Çözümleme
// /
// / ```go
// / func ResolveMultipleFields(resource Resource, item any, fieldNames []string) (map[string]any, error) {
// /     result := make(map[string]any)
// /
// /     for _, fieldName := range fieldNames {
// /         value, err := resource.ResolveField(fieldName, item)
// /         if err != nil {
// /             return nil, fmt.Errorf("alan '%s' çözümlenemedi: %w", fieldName, err)
// /         }
// /         result[fieldName] = value
// /     }
// /
// /     return result, nil
// / }
// /
// / // Kullanım
// / fields := []string{"name", "email", "status_label", "created_at"}
// / values, err := ResolveMultipleFields(resource, user, fields)
// / // values = {"name": "Ahmet", "email": "ahmet@example.com", ...}
// / ```
// /
// / ### API Response Oluşturma
// /
// / ```go
// / func BuildAPIResponse(resource Resource, items []any) ([]map[string]any, error) {
// /     response := make([]map[string]any, 0, len(items))
// /
// /     fields := resource.Fields()
// /     fieldNames := make([]string, 0, len(fields))
// /     for _, field := range fields {
// /         fieldNames = append(fieldNames, field.GetKey())
// /     }
// /
// /     for _, item := range items {
// /         itemData := make(map[string]any)
// /         for _, fieldName := range fieldNames {
// /             value, err := resource.ResolveField(fieldName, item)
// /             if err != nil {
// /                 continue // Hatalı alanları atla
// /             }
// /             itemData[fieldName] = value
// /         }
// /         response = append(response, itemData)
// /     }
// /
// /     return response, nil
// / }
// / ```
// /
// / ## Varsayılan Davranış
// /
// / Varsayılan implementasyon şu adımları izler:
// /
// / 1. **Alan Arama**: FieldsVal içinde fieldName ile eşleşen alanı bulur
// / 2. **Değer Çıkarma**: field.Extract(item) ile model'den değeri çıkarır
// / 3. **Serileştirme**: field.JsonSerialize() ile değeri serileştirir
// / 4. **Değer Döndürme**: Serileştirilmiş değerin "value" anahtarını döndürür
// /
// / ## Requirement 11.2
// /
// / Resource'ların kendi alan çözümleme mantığını tanımlamasına izin verir.
// / Özelleştirilebilir veri dönüşümü sağlar.
// /
// / ## Hata Yönetimi
// /
// / ```go
// / value, err := resource.ResolveField("unknown_field", item)
// / if err != nil {
// /     if err.Error() == "field unknown_field not found" {
// /         // Alan bulunamadı
// /         log.Printf("Alan tanımlı değil: %v", err)
// /     } else {
// /         // Diğer hatalar
// /         log.Printf("Alan çözümleme hatası: %v", err)
// /     }
// / }
// / ```
// /
// / ## Performans Notları
// /
// / - Alan çözümleme işlemi her kayıt için çağrılır
// / - Ağır hesaplamalar cache'lenmelidir
// / - Veritabanı sorguları minimize edilmelidir
// / - Toplu işlemlerde N+1 problemine dikkat edin
// /
// / ## En İyi Uygulamalar
// /
// / 1. **Null Kontrolü**: Değerlerin nil olup olmadığını kontrol edin
// / 2. **Tip Dönüşümü**: Type assertion'ları güvenli yapın
// / 3. **Hata Yönetimi**: Anlamlı hata mesajları döndürün
// / 4. **Cache Kullanımı**: Tekrar eden hesaplamaları cache'leyin
// / 5. **Lazy Loading**: İlişkileri gerektiğinde yükleyin
// /
// / ## Önemli Uyarılar
// /
// / - Alan adı case-sensitive'dir
// / - item parametresi doğru tip olmalıdır
// / - Nil değerler düzgün handle edilmelidir
// / - Circular reference'lardan kaçının
// / - Panic yerine error döndürün
// /
// / ## Notlar
// /
// / - Varsayılan implementasyon temel alan tiplerini destekler
// / - Özel alan tipleri için override edilmelidir
// / - Extract ve JsonSerialize metodları field interface'inden gelir
// / - Serileştirilmiş değer map[string]any formatındadır
func (r Base) ResolveField(fieldName string, item any) (any, error) {
	// Varsayılan implementasyon: alanı bul ve değerini döndür
	for _, field := range r.FieldsVal {
		if field.GetKey() == fieldName {
			// Alan bulundu, değerini çöz
			// Extract the value from the item
			field.Extract(item)
			// Return the serialized value
			serialized := field.JsonSerialize()
			if val, ok := serialized["value"]; ok {
				return val, nil
			}
			return nil, nil
		}
	}
	// Alan bulunamadı
	return nil, fmt.Errorf("field %s not found", fieldName)
}

// / # GetActions Metodu
// /
// / Bu fonksiyon, kaynağın özel işlemlerini (action'larını) döndürür.
// / Action'lar, kayıtlar üzerinde gerçekleştirilebilecek özel işlemlerdir.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Toplu İşlemler**: Birden fazla kayıt üzerinde işlem yapma
// / 2. **Özel İş Mantığı**: Standart CRUD dışında özel işlemler
// / 3. **Durum Değişiklikleri**: Kayıt durumunu değiştirme işlemleri
// / 4. **Dışa Aktarma**: Verileri farklı formatlarda dışa aktarma
// / 5. **Bildirim Gönderme**: Seçili kayıtlara bildirim gönderme
// /
// / ## Döndürür
// /
// / - `[]Action`: Action tanımlamaları dizisi
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da action tanımlama
// / resource := &Base{
// /     ActionsVal: []Action{
// /         {
// /             Name:  "approve",
// /             Label: "Onayla",
// /             Icon:  "check",
// /             Color: "success",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 return ctx.DB.Model(&Model{}).
// /                     Where("id IN ?", ids).
// /                     Update("status", "approved").Error
// /             },
// /             Confirmation: &ActionConfirmation{
// /                 Title:   "Onaylamak istediğinize emin misiniz?",
// /                 Message: "Seçili kayıtlar onaylanacak.",
// /                 Type:    "warning",
// /             },
// /         },
// /         {
// /             Name:  "export",
// /             Label: "Excel'e Aktar",
// /             Icon:  "download",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 return exportToExcel(ctx, ids)
// /             },
// /         },
// /     },
// / }
// /
// / // Kullanım
// / actions := resource.GetActions()
// / ```
// /
// / ## Gelişmiş Örnekler
// /
// / ### E-posta Gönderme Action'ı
// /
// / ```go
// / func (r UserResource) GetActions() []Action {
// /     return []Action{
// /         {
// /             Name:  "send-welcome-email",
// /             Label: "Hoş Geldin E-postası Gönder",
// /             Icon:  "mail",
// /             Color: "primary",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 var users []models.User
// /                 if err := ctx.DB.Find(&users, ids).Error; err != nil {
// /                     return err
// /                 }
// /
// /                 for _, user := range users {
// /                     if err := sendWelcomeEmail(user.Email, user.Name); err != nil {
// /                         return fmt.Errorf("e-posta gönderilemedi: %w", err)
// /                     }
// /                 }
// /
// /                 return nil
// /             },
// /             Confirmation: &ActionConfirmation{
// /                 Title:   "E-posta Gönder",
// /                 Message: fmt.Sprintf("%d kullanıcıya hoş geldin e-postası gönderilecek.", len(ids)),
// /                 Type:    "info",
// /             },
// /             SuccessMessage: "E-postalar başarıyla gönderildi",
// /         },
// /     }
// / }
// / ```
// /
// / ### Durum Değiştirme Action'ları
// /
// / ```go
// / func (r OrderResource) GetActions() []Action {
// /     return []Action{
// /         {
// /             Name:  "mark-as-shipped",
// /             Label: "Kargoya Verildi Olarak İşaretle",
// /             Icon:  "truck",
// /             Color: "info",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 return ctx.DB.Model(&models.Order{}).
// /                     Where("id IN ?", ids).
// /                     Updates(map[string]any{
// /                         "status":      "shipped",
// /                         "shipped_at":  time.Now(),
// /                     }).Error
// /             },
// /             Visible: func(ctx *appContext.Context) bool {
// /                 // Sadece yetkili kullanıcılar görebilir
// /                 return ctx.User.HasPermission("orders.ship")
// /             },
// /         },
// /         {
// /             Name:  "cancel-orders",
// /             Label: "Siparişleri İptal Et",
// /             Icon:  "x-circle",
// /             Color: "danger",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 // İptal edilebilir mi kontrol et
// /                 var orders []models.Order
// /                 ctx.DB.Find(&orders, ids)
// /
// /                 for _, order := range orders {
// /                     if order.Status == "shipped" || order.Status == "delivered" {
// /                         return fmt.Errorf("sipariş #%d iptal edilemez (durum: %s)", order.ID, order.Status)
// /                     }
// /                 }
// /
// /                 // İptal et
// /                 return ctx.DB.Model(&models.Order{}).
// /                     Where("id IN ?", ids).
// /                     Update("status", "cancelled").Error
// /             },
// /             Confirmation: &ActionConfirmation{
// /                 Title:   "Siparişleri İptal Et",
// /                 Message: "Seçili siparişler iptal edilecek. Bu işlem geri alınamaz!",
// /                 Type:    "danger",
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ### Dışa Aktarma Action'ları
// /
// / ```go
// / func (r ReportResource) GetActions() []Action {
// /     return []Action{
// /         {
// /             Name:  "export-excel",
// /             Label: "Excel'e Aktar",
// /             Icon:  "document-download",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 var records []models.Report
// /                 ctx.DB.Find(&records, ids)
// /
// /                 file, err := generateExcel(records)
// /                 if err != nil {
// /                     return err
// /                 }
// /
// /                 // Dosyayı kullanıcıya indir
// /                 ctx.Ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// /                 ctx.Ctx.Set("Content-Disposition", "attachment; filename=rapor.xlsx")
// /                 return ctx.Ctx.Send(file)
// /             },
// /         },
// /         {
// /             Name:  "export-pdf",
// /             Label: "PDF'e Aktar",
// /             Icon:  "document",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 var records []models.Report
// /                 ctx.DB.Find(&records, ids)
// /
// /                 pdf, err := generatePDF(records)
// /                 if err != nil {
// /                     return err
// /                 }
// /
// /                 ctx.Ctx.Set("Content-Type", "application/pdf")
// /                 ctx.Ctx.Set("Content-Disposition", "attachment; filename=rapor.pdf")
// /                 return ctx.Ctx.Send(pdf)
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ### Koşullu Action Görünürlüğü
// /
// / ```go
// / func (r ArticleResource) GetActions() []Action {
// /     return []Action{
// /         {
// /             Name:  "publish",
// /             Label: "Yayınla",
// /             Icon:  "eye",
// /             Color: "success",
// /             Handler: func(ctx *appContext.Context, ids []uint) error {
// /                 return ctx.DB.Model(&models.Article{}).
// /                     Where("id IN ?", ids).
// /                     Updates(map[string]any{
// /                         "status":       "published",
// /                         "published_at": time.Now(),
// /                     }).Error
// /             },
// /             Visible: func(ctx *appContext.Context) bool {
// /                 // Sadece editör ve admin görebilir
// /                 return ctx.User.HasRole("editor") || ctx.User.HasRole("admin")
// /             },
// /             Enabled: func(ctx *appContext.Context, ids []uint) bool {
// /                 // Sadece draft durumundaki makaleler yayınlanabilir
// /                 var count int64
// /                 ctx.DB.Model(&models.Article{}).
// /                     Where("id IN ?", ids).
// /                     Where("status = ?", "draft").
// /                     Count(&count)
// /                 return count == int64(len(ids))
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ## Action Özellikleri
// /
// / - **Name**: Benzersiz action tanımlayıcısı
// / - **Label**: Kullanıcı arayüzünde gösterilecek başlık
// / - **Icon**: Action ikonu
// / - **Color**: Buton rengi (primary, success, danger, warning, info)
// / - **Handler**: Action işleyici fonksiyonu
// / - **Confirmation**: Onay diyalogu ayarları (opsiyonel)
// / - **Visible**: Action'ın görünür olup olmadığını belirleyen fonksiyon
// / - **Enabled**: Action'ın aktif olup olmadığını belirleyen fonksiyon
// / - **SuccessMessage**: Başarılı işlem mesajı
// / - **ErrorMessage**: Hata mesajı
// /
// / ## Requirement 11.4
// /
// / Resource arayüzünü, işlemleri almak için metodlar içerecek şekilde genişletir.
// / Özel iş mantığı yönetimi sağlar.
// /
// / ## Action Tipleri
// /
// / 1. **Bulk Actions**: Birden fazla kayıt üzerinde çalışır
// / 2. **Single Actions**: Tek kayıt üzerinde çalışır
// / 3. **Standalone Actions**: Kayıt seçimi gerektirmez
// /
// / ## Güvenlik Notları
// /
// / - Her action için yetki kontrolü yapın
// / - Kritik işlemler için onay diyalogu kullanın
// / - Handler fonksiyonlarında input validasyonu yapın
// / - Transaction kullanarak veri tutarlılığını sağlayın
// / - Hata durumlarını düzgün handle edin
// /
// / ## Performans Notları
// /
// / - Toplu işlemlerde batch update kullanın
// / - Ağır işlemler için background job kullanın
// / - Timeout ayarları yapın
// / - Progress bar gösterin (uzun işlemler için)
// /
// / ## Kullanıcı Deneyimi
// /
// / - Action'lar buton veya dropdown olarak gösterilir
// / - Confirmation diyalogları kullanıcıyı bilgilendirir
// / - Success/error mesajları gösterilir
// / - Loading state'leri gösterilir
// /
// / ## Notlar
// /
// / - Varsayılan implementasyon ActionsVal alanını döndürür
// / - Action adları benzersiz olmalıdır
// / - Handler fonksiyonu nil olmamalıdır
// / - Boş dizi döndürülürse action gösterilmez
func (r Base) GetActions() []Action {
	if r.ActionsVal != nil {
		return r.ActionsVal
	}
	return []Action{}
}

// / # GetFilters Metodu
// /
// / Bu fonksiyon, kaynağın filtreleme seçeneklerini döndürür.
// / Filtreler, kullanıcıların veri listesini daraltmasına ve arama yapmasına olanak tanır.
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Veri Filtreleme**: Kullanıcıların belirli kriterlere göre veri filtrelemesi
// / 2. **Gelişmiş Arama**: Çoklu kriter ile arama yapma
// / 3. **Durum Filtreleme**: Kayıt durumlarına göre filtreleme
// / 4. **Tarih Aralığı**: Belirli tarih aralıklarında filtreleme
// / 5. **İlişki Filtreleme**: İlişkili kayıtlara göre filtreleme
// /
// / ## Döndürür
// /
// / - `[]Filter`: Filter tanımlamaları dizisi
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Özel resource'da filter tanımlama
// / resource := &Base{
// /     FiltersVal: []Filter{
// /         {
// /             Name:  "status",
// /             Label: "Durum",
// /             Type:  "select",
// /             Options: []FilterOption{
// /                 {Label: "Aktif", Value: "active"},
// /                 {Label: "Pasif", Value: "inactive"},
// /                 {Label: "Beklemede", Value: "pending"},
// /             },
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if v, ok := value.(string); ok && v != "" {
// /                     return db.Where("status = ?", v)
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "created_at",
// /             Label: "Oluşturma Tarihi",
// /             Type:  "date_range",
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if dateRange, ok := value.(map[string]string); ok {
// /                     if start := dateRange["start"]; start != "" {
// /                         db = db.Where("created_at >= ?", start)
// /                     }
// /                     if end := dateRange["end"]; end != "" {
// /                         db = db.Where("created_at <= ?", end)
// /                     }
// /                 }
// /                 return db
// /             },
// /         },
// /     },
// / }
// /
// / // Kullanım
// / filters := resource.GetFilters()
// / ```
// /
// / ## Gelişmiş Örnekler
// /
// / ### E-ticaret Ürün Filtreleri
// /
// / ```go
// / func (r ProductResource) GetFilters() []Filter {
// /     return []Filter{
// /         {
// /             Name:  "category",
// /             Label: "Kategori",
// /             Type:  "select",
// /             Options: func(ctx *appContext.Context) []FilterOption {
// /                 var categories []models.Category
// /                 ctx.DB.Find(&categories)
// /
// /                 options := []FilterOption{}
// /                 for _, cat := range categories {
// /                     options = append(options, FilterOption{
// /                         Label: cat.Name,
// /                         Value: cat.ID,
// /                     })
// /                 }
// /                 return options
// /             },
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if categoryID, ok := value.(uint); ok && categoryID > 0 {
// /                     return db.Where("category_id = ?", categoryID)
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "price_range",
// /             Label: "Fiyat Aralığı",
// /             Type:  "number_range",
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if priceRange, ok := value.(map[string]float64); ok {
// /                     if min := priceRange["min"]; min > 0 {
// /                         db = db.Where("price >= ?", min)
// /                     }
// /                     if max := priceRange["max"]; max > 0 {
// /                         db = db.Where("price <= ?", max)
// /                     }
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "in_stock",
// /             Label: "Stokta Var",
// /             Type:  "boolean",
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if inStock, ok := value.(bool); ok && inStock {
// /                     return db.Where("stock > ?", 0)
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "search",
// /             Label: "Ara",
// /             Type:  "text",
// /             Placeholder: "Ürün adı veya açıklama...",
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if search, ok := value.(string); ok && search != "" {
// /                     searchTerm := "%" + search + "%"
// /                     return db.Where("name LIKE ? OR description LIKE ?", searchTerm, searchTerm)
// /                 }
// /                 return db
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ### Kullanıcı Filtreleri
// /
// / ```go
// / func (r UserResource) GetFilters() []Filter {
// /     return []Filter{
// /         {
// /             Name:  "role",
// /             Label: "Rol",
// /             Type:  "select",
// /             Multiple: true, // Çoklu seçim
// /             Options: []FilterOption{
// /                 {Label: "Admin", Value: "admin"},
// /                 {Label: "Editör", Value: "editor"},
// /                 {Label: "Kullanıcı", Value: "user"},
// /             },
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if roles, ok := value.([]string); ok && len(roles) > 0 {
// /                     return db.Where("role IN ?", roles)
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "registration_date",
// /             Label: "Kayıt Tarihi",
// /             Type:  "date_range",
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if dateRange, ok := value.(map[string]time.Time); ok {
// /                     if !dateRange["start"].IsZero() {
// /                         db = db.Where("created_at >= ?", dateRange["start"])
// /                     }
// /                     if !dateRange["end"].IsZero() {
// /                         db = db.Where("created_at <= ?", dateRange["end"])
// /                     }
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "verified",
// /             Label: "E-posta Doğrulanmış",
// /             Type:  "boolean",
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if verified, ok := value.(bool); ok {
// /                     if verified {
// /                         return db.Where("email_verified_at IS NOT NULL")
// /                     } else {
// /                         return db.Where("email_verified_at IS NULL")
// /                     }
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "subscription",
// /             Label: "Abonelik Durumu",
// /             Type:  "select",
// /             Options: []FilterOption{
// /                 {Label: "Tümü", Value: ""},
// /                 {Label: "Aktif Abonelik", Value: "active"},
// /                 {Label: "Süresi Dolmuş", Value: "expired"},
// /                 {Label: "İptal Edilmiş", Value: "cancelled"},
// /             },
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if status, ok := value.(string); ok && status != "" {
// /                     switch status {
// /                     case "active":
// /                         return db.Where("subscription_expires_at > ?", time.Now())
// /                     case "expired":
// /                         return db.Where("subscription_expires_at <= ?", time.Now()).
// /                             Where("subscription_expires_at IS NOT NULL")
// /                     case "cancelled":
// /                         return db.Where("subscription_cancelled_at IS NOT NULL")
// /                     }
// /                 }
// /                 return db
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ### Sipariş Filtreleri
// /
// / ```go
// / func (r OrderResource) GetFilters() []Filter {
// /     return []Filter{
// /         {
// /             Name:  "status",
// /             Label: "Sipariş Durumu",
// /             Type:  "select",
// /             Options: []FilterOption{
// /                 {Label: "Bekliyor", Value: "pending", Color: "yellow"},
// /                 {Label: "İşleniyor", Value: "processing", Color: "blue"},
// /                 {Label: "Kargoda", Value: "shipped", Color: "purple"},
// /                 {Label: "Teslim Edildi", Value: "delivered", Color: "green"},
// /                 {Label: "İptal Edildi", Value: "cancelled", Color: "red"},
// /             },
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if status, ok := value.(string); ok && status != "" {
// /                     return db.Where("status = ?", status)
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "total_amount",
// /             Label: "Toplam Tutar",
// /             Type:  "number_range",
// /             Min:   0,
// /             Max:   100000,
// /             Step:  100,
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if amountRange, ok := value.(map[string]float64); ok {
// /                     if min := amountRange["min"]; min > 0 {
// /                         db = db.Where("total >= ?", min)
// /                     }
// /                     if max := amountRange["max"]; max > 0 {
// /                         db = db.Where("total <= ?", max)
// /                     }
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "customer",
// /             Label: "Müşteri",
// /             Type:  "search_select", // Arama yapılabilir select
// /             SearchEndpoint: "/api/customers/search",
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if customerID, ok := value.(uint); ok && customerID > 0 {
// /                     return db.Where("customer_id = ?", customerID)
// /                 }
// /                 return db
// /             },
// /         },
// /         {
// /             Name:  "payment_method",
// /             Label: "Ödeme Yöntemi",
// /             Type:  "select",
// /             Options: []FilterOption{
// /                 {Label: "Kredi Kartı", Value: "credit_card", Icon: "credit-card"},
// /                 {Label: "Banka Transferi", Value: "bank_transfer", Icon: "bank"},
// /                 {Label: "Kapıda Ödeme", Value: "cash_on_delivery", Icon: "cash"},
// /             },
// /             Apply: func(db *gorm.DB, value any) *gorm.DB {
// /                 if method, ok := value.(string); ok && method != "" {
// /                     return db.Where("payment_method = ?", method)
// /                 }
// /                 return db
// /             },
// /         },
// /     }
// / }
// / ```
// /
// / ### Filtre Uygulama
// /
// / ```go
// / func ApplyFilters(db *gorm.DB, resource Resource, filterValues map[string]any) *gorm.DB {
// /     filters := resource.GetFilters()
// /
// /     for _, filter := range filters {
// /         if value, exists := filterValues[filter.Name]; exists {
// /             // Filtreyi uygula
// /             db = filter.Apply(db, value)
// /         }
// /     }
// /
// /     return db
// / }
// /
// / // Kullanım
// / filterValues := map[string]any{
// /     "status": "active",
// /     "created_at": map[string]string{
// /         "start": "2024-01-01",
// /         "end":   "2024-12-31",
// /     },
// / }
// /
// / db := ctx.DB.Model(&models.User{})
// / db = ApplyFilters(db, resource, filterValues)
// / db.Find(&users)
// / ```
// /
// / ## Filter Özellikleri
// /
// / - **Name**: Benzersiz filter tanımlayıcısı
// / - **Label**: Kullanıcı arayüzünde gösterilecek başlık
// / - **Type**: Filter tipi (select, text, number_range, date_range, boolean)
// / - **Options**: Seçim listesi için seçenekler
// / - **Apply**: Filtreyi GORM sorgusuna uygulayan fonksiyon
// / - **Multiple**: Çoklu seçim desteği (select için)
// / - **Placeholder**: Placeholder metni (text için)
// / - **Min/Max/Step**: Sayısal filtreler için sınırlar
// / - **Default**: Varsayılan değer
// / - **Visible**: Filter'ın görünür olup olmadığını belirleyen fonksiyon
// /
// / ## Filter Tipleri
// /
// / 1. **select**: Açılır liste (tek veya çoklu seçim)
// / 2. **text**: Metin girişi (arama için)
// / 3. **number**: Sayı girişi
// / 4. **number_range**: Sayı aralığı (min-max)
// / 5. **date**: Tarih seçici
// / 6. **date_range**: Tarih aralığı
// / 7. **boolean**: Checkbox (evet/hayır)
// / 8. **search_select**: Arama yapılabilir açılır liste
// /
// / ## Requirement 11.4
// /
// / Resource arayüzünü, filtreleri almak için metodlar içerecek şekilde genişletir.
// / Gelişmiş veri filtreleme yönetimi sağlar.
// /
// / ## Kullanıcı Deneyimi
// /
// / - Filtreler sidebar veya üst barda gösterilir
// / - Aktif filtreler badge olarak gösterilir
// / - Filtreler URL'de saklanır (bookmark yapılabilir)
// / - Filtre temizleme butonu sağlanır
// / - Filtre sayısı gösterilir
// /
// / ## Performans Notları
// /
// / - Filter sorguları verimli olmalıdır
// / - İndeksli kolonlar kullanın
// / - LIKE sorguları için full-text search düşünün
// / - Çok fazla filter performansı düşürebilir
// / - Filter değerleri cache'lenebilir
// /
// / ## En İyi Uygulamalar
// /
// / 1. **Anlamlı İsimler**: Filter adları açıklayıcı olmalı
// / 2. **Varsayılan Değerler**: Sık kullanılan filtrelere varsayılan değer verin
// / 3. **Validasyon**: Filter değerlerini doğrulayın
// / 4. **Null Kontrolü**: Boş değerleri düzgün handle edin
// / 5. **SQL Injection**: Parametreli sorgular kullanın
// / 6. **Kullanıcı Dostu**: Açık etiketler ve placeholder'lar kullanın
// /
// / ## Güvenlik Notları
// /
// / - Filter değerlerini her zaman sanitize edin
// / - SQL injection'a karşı parametreli sorgular kullanın
// / - Kullanıcı girişlerini doğrulayın
// / - Hassas verileri filtrelemeye izin vermeyin
// / - Rate limiting uygulayın (arama filtreleri için)
// /
// / ## Önemli Uyarılar
// /
// / - Apply fonksiyonu nil olmamalıdır
// / - Filter adları benzersiz olmalıdır
// / - Type değeri geçerli olmalıdır
// / - Options dinamik olarak yüklenebilir
// / - Circular reference'lardan kaçının
// /
// / ## Notlar
// /
// / - Varsayılan implementasyon FiltersVal alanını döndürür
// / - Boş dizi döndürülürse filter gösterilmez
// / - Filtreler sıralı olarak uygulanır
// / - Her filter bağımsız çalışmalıdır
func (r Base) GetFilters() []Filter {
	if r.FiltersVal != nil {
		return r.FiltersVal
	}
	return []Filter{}
}

// SetRecordTitleKey, kayıt başlığı için kullanılacak field adını ayarlar.
//
// Bu metod, RecordTitle metodunun hangi field'ı kullanacağını belirler.
// İlişki fieldlarında kayıtların okunabilir şekilde gösterilmesi için kullanılır.
//
// Parametreler:
// - key: Başlık için kullanılacak field adı (örn: "name", "title", "email")
//
// Döndürür:
// - Resource pointer'ı (method chaining için)
//
// Örnek:
//
//	r.SetRecordTitleKey("name") // User kayıtları için "name" field'ını kullan
//	r.SetRecordTitleKey("title") // Post kayıtları için "title" field'ını kullan
func (b *Base) SetRecordTitleKey(key string) Resource {
	b.recordTitleKey = key
	return b
}

// GetRecordTitleKey, kayıt başlığı için kullanılacak field adını döndürür.
//
// Bu metod, RecordTitle metodunun hangi field'ı kullanacağını belirler.
// Eğer SetRecordTitleKey ile bir değer ayarlanmamışsa varsayılan olarak "id" döner.
//
// Döndürür:
// - string: Başlık için kullanılacak field adı (varsayılan: "id")
//
// Örnek:
//
//	key := r.GetRecordTitleKey() // "name" veya varsayılan "id"
func (b *Base) GetRecordTitleKey() string {
	if b.recordTitleKey == "" {
		return "id"
	}
	return b.recordTitleKey
}

// SetRecordTitleFunc, kayıt başlığını özel bir fonksiyon ile hesaplamak için kullanılır.
//
// Bu metod, karmaşık başlık formatları için kullanılır. Örneğin, birden fazla field'ı
// birleştirerek başlık oluşturmak için kullanılabilir.
//
// Parametreler:
// - fn: Kayıt alıp başlık döndüren fonksiyon
//
// Döndürür:
// - Resource pointer'ı (method chaining için)
//
// Örnek:
//
//	r.SetRecordTitleFunc(func(record any) string {
//	    user := record.(*User)
//	    return user.FirstName + " " + user.LastName
//	})
func (b *Base) SetRecordTitleFunc(fn func(record any) string) Resource {
	b.recordTitleFunc = fn
	return b
}

// RecordTitle, bir kayıt için okunabilir başlık döndürür.
//
// Bu metod, ilişki fieldlarında kayıtların kullanıcı dostu şekilde gösterilmesi için kullanılır.
// Önce SetRecordTitleFunc ile ayarlanmış özel fonksiyonu kontrol eder, yoksa
// GetRecordTitleKey ile belirtilen field'ın değerini reflection ile alır.
//
// Parametreler:
// - record: Başlığı alınacak kayıt (genellikle model instance'ı)
//
// Döndürür:
// - string: Kaydın okunabilir başlığı
//
// Örnek:
//
//	user := &User{ID: 1, Name: "John Doe"}
//	title := resource.RecordTitle(user) // "John Doe"
//
//	post := &Post{ID: 1, Title: "Hello World"}
//	title := resource.RecordTitle(post) // "Hello World"
func (b *Base) RecordTitle(record any) string {
	// Özel fonksiyon varsa onu kullan
	if b.recordTitleFunc != nil {
		return b.recordTitleFunc(record)
	}

	// Reflection ile field değerini al
	titleKey := b.GetRecordTitleKey()

	// Nil kontrolü
	if record == nil {
		return ""
	}

	// Reflection ile field'a eriş
	v := reflect.ValueOf(record)

	// Pointer ise dereference et
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	// Struct değilse boş string döndür
	if v.Kind() != reflect.Struct {
		return ""
	}

	// Field'ı bul (case-insensitive)
	var fieldValue reflect.Value
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if strings.EqualFold(field.Name, titleKey) {
			fieldValue = v.Field(i)
			break
		}
	}

	// Field bulunamadıysa boş string döndür
	if !fieldValue.IsValid() {
		return ""
	}

	// Field değerini string'e çevir
	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", fieldValue.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", fieldValue.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", fieldValue.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", fieldValue.Bool())
	default:
		return fmt.Sprintf("%v", fieldValue.Interface())
	}
}
