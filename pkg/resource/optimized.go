package resource

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ferdiunal/panel.go/pkg/action"
	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// / # OptimizedResource Interface
// /
// / Bu interface, Laravel Nova'nın trait pattern'ini Go'ya uyarlayan optimize edilmiş
// / resource interface'idir. Geleneksel Resource interface'ine göre daha az metod ile
// / daha fazla işlevsellik sağlar ve mixin pattern kullanarak kod tekrarını azaltır.
// /
// / ## Tasarım Felsefesi
// /
// / OptimizedResource, "composition over inheritance" prensibini takip eder ve
// / Go'nun embedding özelliğini kullanarak trait benzeri davranış sağlar. Bu yaklaşım:
// /
// / - **Daha Az Kod**: Tekrarlayan kod yazmayı önler
// / - **Daha Fazla Esneklik**: Sadece ihtiyaç duyulan özellikleri implement edin
// / - **Daha İyi Bakım**: Merkezi mixin'ler sayesinde değişiklikler tek yerden yapılır
// / - **Tip Güvenliği**: Compile-time tip kontrolü sağlar
// /
// / ## Core Methods (8 Temel Metod)
// /
// / Interface, minimum 8 metod gerektirir:
// /
// / 1. `Model()` - Veritabanı model'ini döner
// / 2. `Fields()` - Resource alanlarını döner (bkz: [Fields.md](../../docs/Fields.md))
// / 3. `Slug()` - URL-friendly benzersiz tanımlayıcı
// / 4. `Title()` - İnsan-okunabilir başlık
// / 5. `Policy()` - Yetkilendirme politikası
// / 6. `Repository()` - Veri erişim katmanı
// / 7. `Cards()` - Dashboard widget'ları
// / 8. `Visible()` - Menüde görünürlük durumu
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Basit CRUD Resource
// /
// / ```go
// / type ProductResource struct {
// /     resource.OptimizedBase
// / }
// /
// / func NewProductResource() *ProductResource {
// /     r := &ProductResource{}
// /     r.SetModel(&models.Product{})
// /     r.SetSlug("products")
// /     r.SetTitle("Ürünler")
// /     r.SetVisible(true)
// /     r.SetFieldResolver(r)
// /     return r
// / }
// /
// / func (r *ProductResource) ResolveFields(ctx *context.Context) []fields.Element {
// /     return []fields.Element{
// /         fields.ID("ID").Sortable(),
// /         fields.Text("Ad", "name").Required().Searchable(),
// /         fields.Number("Fiyat", "price").Required(),
// /     }
// / }
// / ```
// /
// / ### Senaryo 2: İlişkili Resource
// /
// / ```go
// / type PostResource struct {
// /     resource.OptimizedBase
// / }
// /
// / func (r *PostResource) ResolveFields(ctx *context.Context) []fields.Element {
// /     return []fields.Element{
// /         fields.Text("Başlık", "title").Required(),
// /         fields.BelongsTo("Yazar", "author_id", "users").
// /             DisplayUsing("name").
// /             WithEagerLoad(),
// /         fields.HasMany("Yorumlar", "comments", "comments"),
// /     }
// / }
// / ```
// / (Detaylı ilişki örnekleri için bkz: [Relationships.md](../../docs/Relationships.md))
// /
// / ### Senaryo 3: Özel Yetkilendirme
// /
// / ```go
// / type AdminResource struct {
// /     resource.OptimizedBase
// / }
// /
// / func NewAdminResource() *AdminResource {
// /     r := &AdminResource{}
// /     r.SetPolicy(&AdminPolicy{})
// /     return r
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Performans**: Gereksiz metod çağrıları ve reflection kullanımı minimize edilmiştir
// / - **Basitlik**: Sadece 8 core metod implement etmeniz yeterlidir
// / - **Genişletilebilirlik**: Mixin'ler ile kolayca yeni özellikler eklenebilir
// / - **Test Edilebilirlik**: Her mixin bağımsız olarak test edilebilir
// / - **Tip Güvenliği**: Interface kontratı compile-time'da doğrulanır
// /
// / ## Dezavantajlar
// /
// / - **Öğrenme Eğrisi**: Mixin pattern'i anlamak zaman alabilir
// / - **Embedding Sınırlamaları**: Go'nun embedding kurallarına dikkat edilmelidir
// / - **Dokümantasyon**: Mixin'lerin nasıl çalıştığını anlamak için dokümantasyon gereklidir
// /
// / ## Önemli Notlar
// /
// / - **Thread Safety**: OptimizedResource implementasyonları thread-safe değildir.
// /   Concurrent kullanım için senkronizasyon mekanizmaları eklenmelidir.
// / - **Lazy Initialization**: Bazı alanlar (repository, policy) lazy olarak initialize edilir.
// / - **Context Kullanımı**: ResolveFields gibi metodlar context alır, bu sayede
// /   kullanıcı bazlı dinamik davranış sağlanabilir.
// /
// / ## İlgili Bileşenler
// /
// / - `OptimizedBase`: Bu interface'in varsayılan implementasyonu
// / - `Authorizable`: Yetkilendirme mixin'i
// / - `Resolvable`: Dinamik çözümleme mixin'i
// / - `Navigable`: Navigasyon ve görünürlük mixin'i
// /
// / ## Referanslar
// /
// / - [Fields Dokümantasyonu](../../docs/Fields.md)
// / - [Relationships Dokümantasyonu](../../docs/Relationships.md)
type OptimizedResource interface {
	// Model, resource'un temsil ettiği veritabanı model'ini döner.
	//
	// Bu metod, GORM ile çalışmak için kullanılır ve migration, query
	// oluşturma gibi işlemlerde model tipini belirler.
	//
	// Döndürür:
	// - Model instance'ı (genellikle pointer)
	//
	// Örnek:
	//   func (r *ProductResource) Model() any {
	//       return &models.Product{}
	//   }
	Model() any

	// Fields, resource'un tüm alanlarını döner.
	//
	// Bu metod, form, liste ve detay görünümlerinde hangi alanların
	// gösterileceğini belirler. Alanlar, görünürlük ayarlarına göre
	// filtrelenir (OnList, OnDetail, OnForm).
	//
	// Döndürür:
	// - Alan listesi (fields.Element slice)
	//
	// Not: Dinamik alan çözümlemesi için FieldResolver kullanın.
	//
	// Örnek:
	//   func (r *ProductResource) Fields() []fields.Element {
	//       return []fields.Element{
	//           fields.ID("ID"),
	//           fields.Text("Ad", "name"),
	//       }
	//   }
	//
	// Detaylı alan örnekleri için bkz: [Fields.md](../../docs/Fields.md)
	Fields() []fields.Element

	// Slug, resource'un URL-friendly benzersiz tanımlayıcısını döner.
	//
	// Slug, API endpoint'lerinde ve routing'de kullanılır.
	// Genellikle model adının çoğul, küçük harf ve tire ile ayrılmış halidir.
	//
	// Döndürür:
	// - URL-friendly string (örn: "products", "blog-posts")
	//
	// Örnek:
	//   func (r *ProductResource) Slug() string {
	//       return "products"
	//   }
	Slug() string

	// Title, resource'un insan-okunabilir başlığını döner.
	//
	// Bu başlık, menülerde, sayfa başlıklarında ve breadcrumb'larda
	// gösterilir. Genellikle çoğul formda ve büyük harfle başlar.
	//
	// Döndürür:
	// - İnsan-okunabilir başlık (örn: "Ürünler", "Blog Yazıları")
	//
	// Örnek:
	//   func (r *ProductResource) Title() string {
	//       return "Ürünler"
	//   }
	Title() string

	// Policy, resource'un yetkilendirme politikasını döner.
	//
	// Policy, kullanıcının resource üzerinde hangi işlemleri yapabileceğini
	// belirler (view, create, update, delete, vb.).
	//
	// Döndürür:
	// - Yetkilendirme politikası (auth.Policy interface)
	//
	// Örnek:
	//   func (r *ProductResource) Policy() auth.Policy {
	//       return &policies.ProductPolicy{}
	//   }
	Policy() auth.Policy

	// Repository, resource'un veri erişim katmanını döner.
	//
	// Repository, veritabanı işlemlerini (CRUD, query, pagination) yönetir.
	// GORM DB instance'ı parametre olarak alır.
	//
	// Parametreler:
	// - db: GORM database instance
	//
	// Döndürür:
	// - Veri sağlayıcı (data.DataProvider interface)
	//
	// Örnek:
	//   func (r *ProductResource) Repository(db *gorm.DB) data.DataProvider {
	//       return data.NewGormProvider(db, r.Model())
	//   }
	Repository(client interface{}) data.DataProvider

	// Cards, resource'un dashboard widget'larını döner.
	//
	// Card'lar, dashboard'da istatistik, grafik ve özet bilgiler
	// göstermek için kullanılır.
	//
	// Döndürür:
	// - Widget listesi (widget.Card slice)
	//
	// Örnek:
	//   func (r *ProductResource) Cards() []widget.Card {
	//       return []widget.Card{
	//           widget.Value("Toplam Ürün", "products.count"),
	//           widget.Trend("Aylık Satış", "sales.monthly"),
	//       }
	//   }
	Cards() []widget.Card

	// Visible, resource'un menüde görünür olup olmadığını döner.
	//
	// False döndürülürse, resource menüde gösterilmez ancak
	// API endpoint'leri hala erişilebilir olur.
	//
	// Döndürür:
	// - Görünürlük durumu (true: görünür, false: gizli)
	//
	// Örnek:
	//   func (r *ProductResource) Visible() bool {
	//       return true
	//   }
	Visible() bool
}

// / # FieldResolver Interface
// /
// / Bu interface, resource alanlarını dinamik olarak çözmek için kullanılır.
// / Context bazlı alan çözümlemesi sayesinde, kullanıcı rolüne, izinlerine veya
// / diğer runtime koşullarına göre farklı alanlar gösterilebilir.
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Rol Bazlı Alan Görünürlüğü
// /
// / ```go
// / func (r *UserResource) ResolveFields(ctx *context.Context) []fields.Element {
// /     baseFields := []fields.Element{
// /         fields.ID("ID"),
// /         fields.Text("Ad", "name"),
// /         fields.Email("E-posta", "email"),
// /     }
// /
// /     // Admin kullanıcılar için ek alanlar
// /     if ctx.User.IsAdmin() {
// /         baseFields = append(baseFields,
// /             fields.Text("IP Adresi", "last_ip"),
// /             fields.DateTime("Son Giriş", "last_login"),
// /         )
// /     }
// /
// /     return baseFields
// / }
// / ```
// /
// / ### Senaryo 2: Özellik Bazlı Alan Ekleme
// /
// / ```go
// / func (r *ProductResource) ResolveFields(ctx *context.Context) []fields.Element {
// /     fields := []fields.Element{
// /         fields.Text("Ürün Adı", "name"),
// /         fields.Number("Fiyat", "price"),
// /     }
// /
// /     // Stok takibi özelliği aktifse
// /     if ctx.App.HasFeature("inventory") {
// /         fields = append(fields,
// /             fields.Number("Stok", "stock"),
// /             fields.Text("Depo", "warehouse"),
// /         )
// /     }
// /
// /     return fields
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Dinamik Davranış**: Runtime'da alan listesi değiştirilebilir
// / - **Context Farkındalığı**: Kullanıcı, istek ve uygulama durumuna göre karar verilebilir
// / - **Güvenlik**: Hassas alanlar sadece yetkili kullanıcılara gösterilebilir
// / - **Esneklik**: Farklı senaryolar için farklı alan setleri
// /
// / ## Önemli Notlar
// /
// / - Context nil olabilir, bu durumu kontrol edin
// / - Performans için ağır hesaplamalardan kaçının
// / - Alan listesi cache'lenebilir ancak context'e bağlı olmalıdır
// /
// / Detaylı alan örnekleri için bkz: [Fields.md](../../docs/Fields.md)
type FieldResolver interface {
	// ResolveFields, verilen context'e göre resource alanlarını çözer ve döner.
	//
	// Bu metod, her alan listesi gerektiğinde çağrılır (liste, detay, form görünümleri).
	// Context parametresi, kullanıcı bilgisi, istek detayları ve uygulama durumu içerir.
	//
	// Parametreler:
	// - ctx: İstek context'i (nil olabilir)
	//
	// Döndürür:
	// - Dinamik olarak çözümlenmiş alan listesi
	//
	// Örnek:
	//   func (r *PostResource) ResolveFields(ctx *context.Context) []fields.Element {
	//       return []fields.Element{
	//           fields.Text("Başlık", "title").Required(),
	//           fields.Textarea("İçerik", "content"),
	//       }
	//   }
	ResolveFields(ctx *context.Context) []fields.Element
}

// / # CardResolver Interface
// /
// / Bu interface, dashboard widget'larını (card'ları) dinamik olarak çözmek için kullanılır.
// / Context bazlı widget çözümlemesi sayesinde, kullanıcıya özel istatistikler ve
// / metrikler gösterilebilir.
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Kullanıcı Bazlı İstatistikler
// /
// / ```go
// / func (r *OrderResource) ResolveCards(ctx *context.Context) []widget.Card {
// /     cards := []widget.Card{
// /         widget.Value("Toplam Sipariş", "orders.count"),
// /     }
// /
// /     // Satış yöneticileri için ek metrikler
// /     if ctx.User.HasRole("sales_manager") {
// /         cards = append(cards,
// /             widget.Trend("Aylık Gelir", "revenue.monthly"),
// /             widget.Partition("Kategori Dağılımı", "orders.by_category"),
// /         )
// /     }
// /
// /     return cards
// / }
// / ```
// /
// / ### Senaryo 2: Zaman Bazlı Widget'lar
// /
// / ```go
// / func (r *AnalyticsResource) ResolveCards(ctx *context.Context) []widget.Card {
// /     now := time.Now()
// /
// /     return []widget.Card{
// /         widget.Value("Bugün", fmt.Sprintf("analytics.today.%s", now.Format("2006-01-02"))),
// /         widget.Trend("Bu Hafta", "analytics.week"),
// /         widget.Chart("Aylık Trend", "analytics.monthly"),
// /     }
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Kişiselleştirme**: Her kullanıcı kendi metriklerini görür
// / - **Dinamik Metrikler**: Runtime'da hesaplanan değerler
// / - **Performans**: Sadece gerekli widget'lar yüklenir
// /
// / ## Önemli Notlar
// /
// / - Widget'lar lazy load edilebilir
// / - Ağır hesaplamalar background job'larda yapılmalıdır
// / - Cache stratejisi kullanın
type CardResolver interface {
	// ResolveCards, verilen context'e göre dashboard widget'larını çözer ve döner.
	//
	// Bu metod, dashboard görünümü yüklendiğinde çağrılır.
	// Context parametresi, kullanıcı bilgisi ve istek detayları içerir.
	//
	// Parametreler:
	// - ctx: İstek context'i (nil olabilir)
	//
	// Döndürür:
	// - Dinamik olarak çözümlenmiş widget listesi
	//
	// Örnek:
	//   func (r *ProductResource) ResolveCards(ctx *context.Context) []widget.Card {
	//       return []widget.Card{
	//           widget.Value("Toplam Ürün", "products.count"),
	//           widget.Trend("Aylık Satış", "sales.monthly"),
	//       }
	//   }
	ResolveCards(ctx *context.Context) []widget.Card
}

// / # FilterResolver Interface
// /
// / Bu interface, liste görünümündeki filtreleri dinamik olarak çözmek için kullanılır.
// / Context bazlı filtre çözümlemesi sayesinde, kullanıcıya özel filtreleme seçenekleri
// / sunulabilir.
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Rol Bazlı Filtreler
// /
// / ```go
// / func (r *OrderResource) ResolveFilters(ctx *context.Context) []Filter {
// /     filters := []Filter{
// /         NewSelectFilter("Durum", "status", map[string]string{
// /             "pending": "Beklemede",
// /             "completed": "Tamamlandı",
// /         }),
// /     }
// /
// /     // Admin kullanıcılar için ek filtreler
// /     if ctx.User.IsAdmin() {
// /         filters = append(filters,
// /             NewSelectFilter("Ödeme Yöntemi", "payment_method", paymentMethods),
// /             NewDateRangeFilter("Tarih Aralığı", "created_at"),
// /         )
// /     }
// /
// /     return filters
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Dinamik Filtreler**: Kullanıcıya özel filtreleme
// / - **Güvenlik**: Hassas filtreler sadece yetkili kullanıcılara
// / - **Esneklik**: Runtime'da filtre seçenekleri değiştirilebilir
type FilterResolver interface {
	// ResolveFilters, verilen context'e göre filtreleri çözer ve döner.
	//
	// Bu metod, liste görünümü yüklendiğinde çağrılır.
	// Context parametresi, kullanıcı bilgisi ve istek detayları içerir.
	//
	// Parametreler:
	// - ctx: İstek context'i (nil olabilir)
	//
	// Döndürür:
	// - Dinamik olarak çözümlenmiş filtre listesi
	ResolveFilters(ctx *context.Context) []Filter
}

// / # LensResolver Interface
// /
// / Bu interface, özel görünümleri (lens'leri) dinamik olarak çözmek için kullanılır.
// / Lens'ler, aynı resource'un farklı perspektiflerden görüntülenmesini sağlar.
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Özel Görünümler
// /
// / ```go
// / func (r *OrderResource) ResolveLenses(ctx *context.Context) []Lens {
// /     lenses := []Lens{
// /         NewLens("Bekleyen Siparişler", "pending").
// /             WithQuery(func(q *Query) *Query {
// /                 return q.Where("status", "=", "pending")
// /             }),
// /     }
// /
// /     // Muhasebe departmanı için özel lens
// /     if ctx.User.Department == "accounting" {
// /         lenses = append(lenses,
// /             NewLens("Ödeme Bekleyenler", "awaiting_payment"),
// /         )
// /     }
// /
// /     return lenses
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Özel Görünümler**: Farklı kullanım senaryoları için optimize edilmiş listeler
// / - **Performans**: Önceden tanımlı query'ler
// / - **Kullanıcı Deneyimi**: Hızlı erişim için kısayollar
type LensResolver interface {
	// ResolveLenses, verilen context'e göre lens'leri çözer ve döner.
	//
	// Bu metod, resource menüsü oluşturulurken çağrılır.
	// Context parametresi, kullanıcı bilgisi ve istek detayları içerir.
	//
	// Parametreler:
	// - ctx: İstek context'i (nil olabilir)
	//
	// Döndürür:
	// - Dinamik olarak çözümlenmiş lens listesi
	ResolveLenses(ctx *context.Context) []Lens
}

// / # ActionResolver Interface
// /
// / Bu interface, toplu işlemleri (action'ları) dinamik olarak çözmek için kullanılır.
// / Context bazlı action çözümlemesi sayesinde, kullanıcıya özel işlemler sunulabilir.
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Yetki Bazlı İşlemler
// /
// / ```go
// / func (r *OrderResource) ResolveActions(ctx *context.Context) []Action {
// /     actions := []Action{
// /         NewAction("Dışa Aktar", "export").
// /             WithHandler(exportHandler),
// /     }
// /
// /     // Admin kullanıcılar için ek işlemler
// /     if ctx.User.IsAdmin() {
// /         actions = append(actions,
// /             NewAction("Toplu Sil", "bulk_delete").
// /                 WithConfirmation("Emin misiniz?").
// /                 WithHandler(bulkDeleteHandler),
// /         )
// /     }
// /
// /     return actions
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Güvenlik**: Hassas işlemler sadece yetkili kullanıcılara
// / - **Esneklik**: Runtime'da işlem listesi değiştirilebilir
// / - **Kullanıcı Deneyimi**: Kullanıcıya özel işlemler
type ActionResolver interface {
	// ResolveActions, verilen context'e göre işlemleri çözer ve döner.
	//
	// Bu metod, liste görünümü yüklendiğinde çağrılır.
	// Context parametresi, kullanıcı bilgisi ve istek detayları içerir.
	//
	// Parametreler:
	// - ctx: İstek context'i (nil olabilir)
	//
	// Döndürür:
	// - Dinamik olarak çözümlenmiş işlem listesi
	ResolveActions(ctx *context.Context) []Action
}

// / # Authorizable Struct (Mixin)
// /
// / Bu yapı, resource'lara yetkilendirme işlevselliği ekleyen bir mixin'dir.
// / Go'nun embedding özelliği kullanılarak trait benzeri davranış sağlar.
// /
// / ## Mixin Pattern
// /
// / Mixin pattern, kod tekrarını önlemek ve ortak işlevselliği paylaşmak için
// / kullanılan bir tasarım desenidir. Go'da embedding ile implement edilir:
// /
// / ```go
// / type MyResource struct {
// /     resource.Authorizable  // Mixin embedding
// /     // ... diğer alanlar
// / }
// / ```
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Özel Policy Tanımlama
// /
// / ```go
// / type ProductResource struct {
// /     resource.OptimizedBase
// / }
// /
// / func NewProductResource() *ProductResource {
// /     r := &ProductResource{}
// /     r.SetPolicy(&policies.ProductPolicy{})
// /     return r
// / }
// / ```
// /
// / ### Senaryo 2: Dinamik Policy
// /
// / ```go
// / func NewUserResource(isAdmin bool) *UserResource {
// /     r := &UserResource{}
// /
// /     if isAdmin {
// /         r.SetPolicy(&policies.AdminUserPolicy{})
// /     } else {
// /         r.SetPolicy(&policies.RegularUserPolicy{})
// /     }
// /
// /     return r
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Kod Tekrarı Yok**: Policy yönetimi tek bir yerde
// / - **Kolay Kullanım**: Sadece embed edin ve kullanın
// / - **Esneklik**: Runtime'da policy değiştirilebilir
// / - **Tip Güvenliği**: Interface kontratı garanti edilir
// /
// / ## Önemli Notlar
// /
// / - Policy nil olabilir, bu durumu kontrol edin
// / - Thread-safe değildir, concurrent kullanım için mutex ekleyin
// / - Policy değişiklikleri mevcut istekleri etkilemez
type Authorizable struct {
	policy auth.Policy
}

// / SetPolicy, resource'un yetkilendirme politikasını ayarlar.
// /
// / Bu metod, resource'un hangi policy ile yetkilendirme yapacağını belirler.
// / Policy, kullanıcının resource üzerinde hangi işlemleri yapabileceğini kontrol eder.
// /
// / Parametreler:
// / - p: Yetkilendirme politikası (auth.Policy interface)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetPolicy(&policies.ProductPolicy{})
func (a *Authorizable) SetPolicy(p auth.Policy) {
	a.policy = p
}

// / GetPolicy, resource'un yetkilendirme politikasını döner.
// /
// / Bu metod, mevcut policy'yi almak için kullanılır.
// / Policy nil olabilir, bu durumu kontrol etmek önemlidir.
// /
// / Döndürür:
// / - Yetkilendirme politikası (nil olabilir)
// /
// / Örnek:
// /   policy := r.GetPolicy()
// /   if policy != nil && policy.CanView(user, item) {
// /       // İşlem yapılabilir
// /   }
func (a *Authorizable) GetPolicy() auth.Policy {
	return a.policy
}

// / # Resolvable Struct (Mixin)
// /
// / Bu yapı, resource'lara dinamik çözümleme işlevselliği ekleyen bir mixin'dir.
// / Field, card, filter, lens ve action resolver'ları yönetir.
// /
// / ## Resolver Pattern
// /
// / Resolver pattern, bileşenlerin dinamik olarak çözümlenmesini sağlar.
// / Bu sayede context bazlı davranış elde edilir:
// /
// / - **FieldResolver**: Kullanıcıya göre farklı alanlar
// / - **CardResolver**: Kullanıcıya göre farklı widget'lar
// / - **FilterResolver**: Kullanıcıya göre farklı filtreler
// / - **LensResolver**: Kullanıcıya göre farklı görünümler
// / - **ActionResolver**: Kullanıcıya göre farklı işlemler
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Field Resolver Kullanımı
// /
// / ```go
// / type ProductResource struct {
// /     resource.OptimizedBase
// / }
// /
// / func NewProductResource() *ProductResource {
// /     r := &ProductResource{}
// /     r.SetFieldResolver(r)  // Self-resolver
// /     return r
// / }
// /
// / func (r *ProductResource) ResolveFields(ctx *context.Context) []fields.Element {
// /     // Context bazlı alan çözümlemesi
// /     return []fields.Element{...}
// / }
// / ```
// /
// / ### Senaryo 2: Çoklu Resolver
// /
// / ```go
// / func NewOrderResource() *OrderResource {
// /     r := &OrderResource{}
// /     r.SetFieldResolver(r)
// /     r.SetCardResolver(r)
// /     r.SetFilterResolver(r)
// /     return r
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Dinamik Davranış**: Runtime'da bileşen çözümlemesi
// / - **Context Farkındalığı**: Kullanıcı ve istek bazlı karar verme
// / - **Modülerlik**: Her resolver bağımsız olarak yönetilebilir
// / - **Genişletilebilirlik**: Yeni resolver türleri kolayca eklenebilir
// /
// / ## Önemli Notlar
// /
// / - Resolver'lar nil olabilir, bu durumda boş slice döner
// / - Context nil olabilir, resolver implementasyonunda kontrol edin
// / - Resolver'lar lazy olarak çağrılır (sadece gerektiğinde)
type Resolvable struct {
	fieldResolver  FieldResolver
	cardResolver   CardResolver
	filterResolver FilterResolver
	lensResolver   LensResolver
	actionResolver ActionResolver
}

// / SetFieldResolver, field resolver'ı ayarlar.
// /
// / Bu metod, alanların dinamik olarak çözümlenmesi için kullanılacak
// / resolver'ı belirler. Genellikle resource kendisi resolver olarak kullanılır.
// /
// / Parametreler:
// / - fr: Field resolver (FieldResolver interface)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetFieldResolver(r)  // Self-resolver
func (r *Resolvable) SetFieldResolver(fr FieldResolver) {
	r.fieldResolver = fr
}

// / SetCardResolver, card resolver'ı ayarlar.
// /
// / Bu metod, dashboard widget'larının dinamik olarak çözümlenmesi için
// / kullanılacak resolver'ı belirler.
// /
// / Parametreler:
// / - cr: Card resolver (CardResolver interface)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetCardResolver(r)
func (r *Resolvable) SetCardResolver(cr CardResolver) {
	r.cardResolver = cr
}

// / SetFilterResolver, filter resolver'ı ayarlar.
// /
// / Bu metod, filtrelerin dinamik olarak çözümlenmesi için kullanılacak
// / resolver'ı belirler.
// /
// / Parametreler:
// / - fr: Filter resolver (FilterResolver interface)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetFilterResolver(r)
func (r *Resolvable) SetFilterResolver(fr FilterResolver) {
	r.filterResolver = fr
}

// / SetLensResolver, lens resolver'ı ayarlar.
// /
// / Bu metod, özel görünümlerin dinamik olarak çözümlenmesi için kullanılacak
// / resolver'ı belirler.
// /
// / Parametreler:
// / - lr: Lens resolver (LensResolver interface)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetLensResolver(r)
func (r *Resolvable) SetLensResolver(lr LensResolver) {
	r.lensResolver = lr
}

// / SetActionResolver, action resolver'ı ayarlar.
// /
// / Bu metod, toplu işlemlerin dinamik olarak çözümlenmesi için kullanılacak
// / resolver'ı belirler.
// /
// / Parametreler:
// / - ar: Action resolver (ActionResolver interface)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetActionResolver(r)
func (r *Resolvable) SetActionResolver(ar ActionResolver) {
	r.actionResolver = ar
}

// / ResolveFields, verilen context'e göre alanları çözer.
// /
// / Bu metod, field resolver varsa onun ResolveFields metodunu çağırır.
// / Resolver yoksa boş slice döner.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Çözümlenmiş alan listesi (boş olabilir)
// /
// / Örnek:
// /   fields := r.ResolveFields(ctx)
// /   for _, field := range fields {
// /       // Alan işlemleri
// /   }
func (r *Resolvable) ResolveFields(ctx *context.Context) []fields.Element {
	if r.fieldResolver != nil {
		return r.fieldResolver.ResolveFields(ctx)
	}
	return []fields.Element{}
}

// / ResolveCards, verilen context'e göre card'ları çözer.
// /
// / Bu metod, card resolver varsa onun ResolveCards metodunu çağırır.
// / Resolver yoksa boş slice döner.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Çözümlenmiş widget listesi (boş olabilir)
func (r *Resolvable) ResolveCards(ctx *context.Context) []widget.Card {
	if r.cardResolver != nil {
		return r.cardResolver.ResolveCards(ctx)
	}
	return []widget.Card{}
}

// / ResolveFilters, verilen context'e göre filtreleri çözer.
// /
// / Bu metod, filter resolver varsa onun ResolveFilters metodunu çağırır.
// / Resolver yoksa boş slice döner.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Çözümlenmiş filtre listesi (boş olabilir)
func (r *Resolvable) ResolveFilters(ctx *context.Context) []Filter {
	if r.filterResolver != nil {
		return r.filterResolver.ResolveFilters(ctx)
	}
	return []Filter{}
}

// / ResolveLenses, verilen context'e göre lens'leri çözer.
// /
// / Bu metod, lens resolver varsa onun ResolveLenses metodunu çağırır.
// / Resolver yoksa boş slice döner.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Çözümlenmiş lens listesi (boş olabilir)
func (r *Resolvable) ResolveLenses(ctx *context.Context) []Lens {
	if r.lensResolver != nil {
		return r.lensResolver.ResolveLenses(ctx)
	}
	return []Lens{}
}

// / ResolveActions, verilen context'e göre işlemleri çözer.
// /
// / Bu metod, action resolver varsa onun ResolveActions metodunu çağırır.
// / Resolver yoksa boş slice döner.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Çözümlenmiş işlem listesi (boş olabilir)
func (r *Resolvable) ResolveActions(ctx *context.Context) []Action {
	if r.actionResolver != nil {
		return r.actionResolver.ResolveActions(ctx)
	}
	return []Action{}
}

// / # Navigable Struct (Mixin)
// /
// / Bu yapı, resource'lara navigasyon ve görünürlük işlevselliği ekleyen bir mixin'dir.
// / Menü yapısı, ikon, grup, sıralama ve dialog tipi gibi UI ile ilgili özellikleri yönetir.
// /
// / ## UI Konfigürasyonu
// /
// / Navigable mixin, resource'un kullanıcı arayüzünde nasıl görüneceğini kontrol eder:
// /
// / - **Icon**: Menüde gösterilecek ikon
// / - **Group**: Menü gruplaması (örn: "İçerik Yönetimi", "Kullanıcılar")
// / - **NavigationOrder**: Menüdeki sıralama önceliği
// / - **Visible**: Menüde görünür olup olmadığı
// / - **DialogType**: Form görünüm tipi (modal, drawer, fullscreen)
// / - **Sortable**: Varsayılan sıralama ayarları
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Temel Navigasyon Ayarları
// /
// / ```go
// / func NewProductResource() *ProductResource {
// /     r := &ProductResource{}
// /     r.SetIcon("shopping-cart")
// /     r.SetGroup("E-Ticaret")
// /     r.SetNavigationOrder(10)
// /     r.SetVisible(true)
// /     return r
// / }
// / ```
// /
// / ### Senaryo 2: Gruplu Menü Yapısı
// /
// / ```go
// / // İçerik Yönetimi Grubu
// / postResource.SetGroup("İçerik Yönetimi")
// / postResource.SetNavigationOrder(1)
// /
// / categoryResource.SetGroup("İçerik Yönetimi")
// / categoryResource.SetNavigationOrder(2)
// /
// / tagResource.SetGroup("İçerik Yönetimi")
// / tagResource.SetNavigationOrder(3)
// /
// / // Kullanıcı Yönetimi Grubu
// / userResource.SetGroup("Kullanıcı Yönetimi")
// / userResource.SetNavigationOrder(1)
// /
// / roleResource.SetGroup("Kullanıcı Yönetimi")
// / roleResource.SetNavigationOrder(2)
// / ```
// /
// / ### Senaryo 3: Dialog Tipi Ayarları
// /
// / ```go
// / // Modal dialog (küçük formlar için)
// / quickNoteResource.SetDialogType(DialogTypeModal)
// /
// / // Drawer (yan panel, orta boyut formlar için)
// / userResource.SetDialogType(DialogTypeDrawer)
// /
// / // Fullscreen (büyük, karmaşık formlar için)
// / productResource.SetDialogType(DialogTypeFullscreen)
// / ```
// /
// / ### Senaryo 4: Varsayılan Sıralama
// /
// / ```go
// / r.SetSortable([]Sortable{
// /     {Field: "created_at", Direction: "desc"},
// /     {Field: "name", Direction: "asc"},
// / })
// / ```
// /
// / ## Dialog Tipleri
// /
// / - **Modal**: Küçük popup pencere, basit formlar için ideal
// / - **Drawer**: Yan panel, orta boyut formlar için uygun
// / - **Fullscreen**: Tam ekran, karmaşık formlar için önerilir
// /
// / ## Avantajlar
// /
// / - **Merkezi UI Yönetimi**: Tüm UI ayarları tek yerde
// / - **Tutarlı Görünüm**: Standart menü yapısı
// / - **Esneklik**: Runtime'da değiştirilebilir ayarlar
// / - **Kullanıcı Deneyimi**: Optimize edilmiş navigasyon
// /
// / ## Önemli Notlar
// /
// / - NavigationOrder küçük değerler önce gösterilir (1, 2, 3...)
// / - Visible false olsa bile API endpoint'leri erişilebilir kalır
// / - Icon isimleri frontend icon kütüphanesine göre ayarlanmalıdır
// / - Group boş bırakılırsa resource grupsuz gösterilir
type Navigable struct {
	icon            string
	group           string
	navigationOrder int
	visible         bool
	dialogType      DialogType
	dialogSize      DialogSize
	sortable        []Sortable
	rowClickAction  IndexRowClickAction
	paginationType  IndexPaginationType
	reorderEnabled  bool
	reorderColumn   string
}

// / SetIcon, menüde gösterilecek ikon adını ayarlar.
// /
// / İkon adı, frontend'de kullanılan icon kütüphanesine göre belirlenmelidir.
// / Genellikle Heroicons, Lucide veya Material Icons kullanılır.
// /
// / Parametreler:
// / - icon: İkon adı (örn: "shopping-cart", "users", "document-text")
// /
// / Örnek:
// /   r.SetIcon("shopping-cart")
// /   r.SetIcon("users")
// /   r.SetIcon("document-text")
func (n *Navigable) SetIcon(icon string) {
	n.icon = icon
}

// / GetIcon, menüde gösterilecek ikon adını döner.
// /
// / Döndürür:
// / - İkon adı (boş string olabilir)
// /
// / Örnek:
// /   icon := r.GetIcon()
// /   if icon == "" {
// /       icon = "default-icon"
// /   }
func (n *Navigable) GetIcon() string {
	return n.icon
}

// / SetGroup, resource'un menü grubunu ayarlar.
// /
// / Grup, ilgili resource'ları menüde bir arada göstermek için kullanılır.
// / Aynı gruptaki resource'lar birlikte listelenir ve grup başlığı altında gösterilir.
// /
// / Parametreler:
// / - group: Grup adı (örn: "İçerik Yönetimi", "Kullanıcılar", "E-Ticaret")
// /
// / Örnek:
// /   postResource.SetGroup("İçerik Yönetimi")
// /   userResource.SetGroup("Kullanıcılar")
// /   productResource.SetGroup("E-Ticaret")
func (n *Navigable) SetGroup(group string) {
	n.group = group
}

// / GetGroup, resource'un menü grubunu döner.
// /
// / Döndürür:
// / - Grup adı (boş string olabilir)
// /
// / Örnek:
// /   group := r.GetGroup()
// /   if group == "" {
// /       group = "Diğer"
// /   }
func (n *Navigable) GetGroup() string {
	return n.group
}

// / SetNavigationOrder, menüdeki sıralama önceliğini ayarlar.
// /
// / Küçük değerler önce gösterilir. Aynı gruptaki resource'lar kendi
// / aralarında bu değere göre sıralanır.
// /
// / Parametreler:
// / - order: Sıralama önceliği (1, 2, 3... küçük değerler önce)
// /
// / Örnek:
// /   postResource.SetNavigationOrder(1)      // İlk sırada
// /   categoryResource.SetNavigationOrder(2)  // İkinci sırada
// /   tagResource.SetNavigationOrder(3)       // Üçüncü sırada
func (n *Navigable) SetNavigationOrder(order int) {
	n.navigationOrder = order
}

// / GetNavigationOrder, menüdeki sıralama önceliğini döner.
// /
// / Döndürür:
// / - Sıralama önceliği (varsayılan: 0)
// /
// / Örnek:
// /   order := r.GetNavigationOrder()
func (n *Navigable) GetNavigationOrder() int {
	return n.navigationOrder
}

// / SetVisible, resource'un menüde görünür olup olmadığını ayarlar.
// /
// / False olarak ayarlanırsa, resource menüde gösterilmez ancak
// / API endpoint'leri hala erişilebilir olur. Bu, programatik erişim
// / için kullanılan ancak UI'da gösterilmesi gerekmeyen resource'lar için kullanışlıdır.
// /
// / Parametreler:
// / - visible: Görünürlük durumu (true: görünür, false: gizli)
// /
// / Örnek:
// /   // Menüde göster
// /   productResource.SetVisible(true)
// /
// /   // Menüde gizle (API hala erişilebilir)
// /   internalLogResource.SetVisible(false)
func (n *Navigable) SetVisible(visible bool) {
	n.visible = visible
}

// / IsVisible, resource'un menüde görünür olup olmadığını döner.
// /
// / Döndürür:
// / - Görünürlük durumu (true: görünür, false: gizli)
// /
// / Örnek:
// /   if r.IsVisible() {
// /       // Menüye ekle
// /   }
func (n *Navigable) IsVisible() bool {
	return n.visible
}

// / SetDialogType, form görünüm tipini ayarlar.
// /
// / Dialog tipi, create ve update formlarının nasıl gösterileceğini belirler:
// / - Modal: Küçük popup pencere (basit formlar)
// / - Drawer: Yan panel (orta boyut formlar)
// / - Fullscreen: Tam ekran (karmaşık formlar)
// /
// / Parametreler:
// / - dt: Dialog tipi (DialogTypeModal, DialogTypeDrawer, DialogTypeFullscreen)
// /
// / Örnek:
// /   // Basit form için modal
// /   quickNoteResource.SetDialogType(DialogTypeModal)
// /
// /   // Orta boyut form için drawer
// /   userResource.SetDialogType(DialogTypeDrawer)
// /
// /   // Karmaşık form için fullscreen
// /   productResource.SetDialogType(DialogTypeFullscreen)
func (n *Navigable) SetDialogType(dt DialogType) {
	n.dialogType = dt
}

// / GetDialogType, form görünüm tipini döner.
// /
// / Döndürür:
// / - Dialog tipi (varsayılan: DialogTypeModal)
// /
// / Örnek:
// /   dialogType := r.GetDialogType()
// /   switch dialogType {
// /   case DialogTypeModal:
// /       // Modal render et
// /   case DialogTypeDrawer:
// /       // Drawer render et
// /   case DialogTypeFullscreen:
// /       // Fullscreen render et
// /   }
func (n *Navigable) GetDialogType() DialogType {
	return n.dialogType
}

// SetDialogSize, form modal/sheet genişlik preset'ini ayarlar.
func (n *Navigable) SetDialogSize(ds DialogSize) {
	n.dialogSize = ds
}

// GetDialogSize, form modal/sheet genişlik preset'ini döner.
// Boş ise varsayılan olarak DialogSizeMD kullanılır.
func (n *Navigable) GetDialogSize() DialogSize {
	if n.dialogSize == "" {
		return DialogSizeMD
	}
	return n.dialogSize
}

// SetIndexRowClickAction, index satır tıklama aksiyonunu ayarlar.
func (n *Navigable) SetIndexRowClickAction(action IndexRowClickAction) {
	n.rowClickAction = NormalizeIndexRowClickAction(action)
}

// GetIndexRowClickAction, index satır tıklama aksiyonunu döner.
func (n *Navigable) GetIndexRowClickAction() IndexRowClickAction {
	return NormalizeIndexRowClickAction(n.rowClickAction)
}

// SetIndexPaginationType, index pagination tipini ayarlar.
func (n *Navigable) SetIndexPaginationType(paginationType IndexPaginationType) {
	n.paginationType = NormalizeIndexPaginationType(paginationType)
}

// GetIndexPaginationType, index pagination tipini döner.
// Varsayılan değer "links" olarak normalize edilir.
func (n *Navigable) GetIndexPaginationType() IndexPaginationType {
	return NormalizeIndexPaginationType(n.paginationType)
}

// SetIndexReorder, index reorder ayarlarını günceller.
func (n *Navigable) SetIndexReorder(enabled bool, column string) {
	normalizedColumn := NormalizeIndexReorderColumn(column)
	n.reorderEnabled = enabled && normalizedColumn != ""
	n.reorderColumn = normalizedColumn
}

// EnableIndexReorder, verilen kolon için reorder özelliğini etkinleştirir.
func (n *Navigable) EnableIndexReorder(column string) {
	n.SetIndexReorder(true, column)
}

// DisableIndexReorder, index reorder özelliğini devre dışı bırakır.
func (n *Navigable) DisableIndexReorder() {
	n.reorderEnabled = false
}

// GetIndexReorderConfig, index reorder ayarlarını döner.
func (n *Navigable) GetIndexReorderConfig() IndexReorderConfig {
	column := NormalizeIndexReorderColumn(n.reorderColumn)

	return IndexReorderConfig{
		Enabled: n.reorderEnabled && column != "",
		Column:  column,
	}
}

// / SetSortable, varsayılan sıralama ayarlarını belirler.
// /
// / Bu ayarlar, liste görünümü ilk yüklendiğinde uygulanacak
// / sıralama kurallarını tanımlar. Birden fazla alan belirtilebilir.
// /
// / Parametreler:
// / - sortable: Sıralama ayarları listesi
// /
// / Örnek:
// /   // Önce oluşturma tarihine göre azalan, sonra ada göre artan
// /   r.SetSortable([]Sortable{
// /       {Field: "created_at", Direction: "desc"},
// /       {Field: "name", Direction: "asc"},
// /   })
// /
// /   // Sadece ada göre artan
// /   r.SetSortable([]Sortable{
// /       {Field: "name", Direction: "asc"},
// /   })
func (n *Navigable) SetSortable(sortable []Sortable) {
	n.sortable = sortable
}

// / GetSortable, varsayılan sıralama ayarlarını döner.
// /
// / Döndürür:
// / - Sıralama ayarları listesi (boş olabilir)
// /
// / Örnek:
// /   sortable := r.GetSortable()
// /   for _, sort := range sortable {
// /       query = query.Order(sort.Field + " " + sort.Direction)
// /   }
func (n *Navigable) GetSortable() []Sortable {
	return n.sortable
}

// / # OptimizedBase Struct
// /
// / Bu yapı, OptimizedResource interface'ini implement eden temel struct'tır.
// / Tüm mixin'leri (Authorizable, Resolvable, Navigable) embed ederek trait benzeri
// / davranış sağlar ve resource oluşturmayı kolaylaştırır.
// /
// / ## Composition Pattern
// /
// / OptimizedBase, Go'nun embedding özelliğini kullanarak birden fazla mixin'i
// / bir araya getirir. Bu sayede:
// /
// / - **Authorizable**: Yetkilendirme işlevselliği
// / - **Resolvable**: Dinamik çözümleme işlevselliği
// / - **Navigable**: Navigasyon ve UI işlevselliği
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Basit Resource Oluşturma
// /
// / ```go
// / type ProductResource struct {
// /     resource.OptimizedBase
// / }
// /
// / func NewProductResource() *ProductResource {
// /     r := &ProductResource{}
// /     r.SetModel(&models.Product{})
// /     r.SetSlug("products")
// /     r.SetTitle("Ürünler")
// /     r.SetVisible(true)
// /     r.SetIcon("shopping-cart")
// /     r.SetGroup("E-Ticaret")
// /     r.SetFieldResolver(r)
// /     return r
// / }
// / ```
// /
// / ### Senaryo 2: Tam Özellikli Resource
// /
// / ```go
// / func NewOrderResource() *OrderResource {
// /     r := &OrderResource{}
// /
// /     // Model ve temel ayarlar
// /     r.SetModel(&models.Order{})
// /     r.SetSlug("orders")
// /     r.SetTitle("Siparişler")
// /
// /     // UI ayarları
// /     r.SetIcon("shopping-bag")
// /     r.SetGroup("E-Ticaret")
// /     r.SetNavigationOrder(2)
// /     r.SetVisible(true)
// /     r.SetDialogType(DialogTypeDrawer)
// /
// /     // Resolver'lar
// /     r.SetFieldResolver(r)
// /     r.SetCardResolver(r)
// /     r.SetFilterResolver(r)
// /
// /     // Policy
// /     r.SetPolicy(&policies.OrderPolicy{})
// /
// /     // Sıralama
// /     r.SetSortable([]Sortable{
// /         {Field: "created_at", Direction: "desc"},
// /     })
// /
// /     return r
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Hızlı Başlangıç**: Tüm temel işlevsellik hazır
// / - **Modüler Yapı**: Sadece ihtiyaç duyulan özellikleri kullanın
// / - **Tip Güvenliği**: Interface kontratı garanti edilir
// / - **Kolay Genişletme**: Yeni özellikler kolayca eklenebilir
// /
// / ## Önemli Notlar
// /
// / - Thread-safe değildir, concurrent kullanım için senkronizasyon ekleyin
// / - Repository lazy olarak initialize edilir
// / - Resolver'lar nil olabilir, bu durumda varsayılan davranış kullanılır
// /
// / Detaylı alan örnekleri için bkz: [Fields.md](../../docs/Fields.md)
// / Detaylı ilişki örnekleri için bkz: [Relationships.md](../../docs/Relationships.md)
type OptimizedBase struct {
	Authorizable
	Resolvable
	Navigable
	model           any
	slug            string
	title           string
	titleFunc       func(*fiber.Ctx) string
	groupFunc       func(*fiber.Ctx) string
	repository      data.DataProvider
	cards           []widget.Card
	openAPIDisabled bool
	recordTitleKey  string
	recordTitleFunc func(record any) string
}

// / SetModel, resource'un temsil ettiği veritabanı model'ini ayarlar.
// /
// / Bu metod, GORM ile çalışmak için kullanılacak model tipini belirler.
// / Model, migration, query oluşturma ve veri işleme için kullanılır.
// /
// / Parametreler:
// / - m: Model instance'ı (genellikle pointer)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetModel(&models.Product{})
func (b *OptimizedBase) SetModel(m any) {
	b.model = m
}

// / Model, resource'un temsil ettiği veritabanı model'ini döner.
// /
// / Bu metod, OptimizedResource interface'inin bir parçasıdır ve
// / GORM işlemleri için model tipini sağlar.
// /
// / Döndürür:
// / - Model instance'ı (any type)
// /
// / Örnek:
// /   model := r.Model()
// /   // GORM ile kullanım
// /   db.Model(model).Find(&results)
func (b *OptimizedBase) Model() any {
	return b.model
}

// / SetSlug, resource'un URL-friendly benzersiz tanımlayıcısını ayarlar.
// /
// / Slug, API endpoint'lerinde ve routing'de kullanılır.
// / Genellikle model adının çoğul, küçük harf ve tire ile ayrılmış halidir.
// /
// / Parametreler:
// / - s: URL-friendly string (örn: "products", "blog-posts")
// /
// / Örnek:
// /   r.SetSlug("products")
// /   r.SetSlug("blog-posts")
// /   r.SetSlug("user-profiles")
func (b *OptimizedBase) SetSlug(s string) {
	b.slug = s
}

// / Slug, resource'un URL-friendly benzersiz tanımlayıcısını döner.
// /
// / Bu metod, OptimizedResource interface'inin bir parçasıdır ve
// / routing için slug değerini sağlar.
// /
// / Döndürür:
// / - URL-friendly string
// /
// / Örnek:
// /   slug := r.Slug()
// /   // API endpoint: /api/resources/{slug}
// /   endpoint := fmt.Sprintf("/api/resources/%s", slug)
func (b *OptimizedBase) Slug() string {
	return b.slug
}

// / SetTitle, resource'un insan-okunabilir başlığını ayarlar.
// /
// / Bu başlık, menülerde, sayfa başlıklarında ve breadcrumb'larda
// / gösterilir. Genellikle çoğul formda ve büyük harfle başlar.
// /
// / Parametreler:
// / - t: İnsan-okunabilir başlık (örn: "Ürünler", "Blog Yazıları")
// /
// / Örnek:
// /   r.SetTitle("Ürünler")
// /   r.SetTitle("Blog Yazıları")
// /   r.SetTitle("Kullanıcı Profilleri")
func (b *OptimizedBase) SetTitle(t string) {
	b.title = t
}

// / Title, resource'un insan-okunabilir başlığını döner.
// /
// / Bu metod, OptimizedResource interface'inin bir parçasıdır ve
// / UI'da gösterilecek başlığı sağlar.
// /
// / Döndürür:
// / - İnsan-okunabilir başlık
// /
// / Örnek:
// /   title := r.Title()
// /   // Sayfa başlığı: {title} - Admin Panel
// /   pageTitle := fmt.Sprintf("%s - Admin Panel", title)
func (b *OptimizedBase) Title() string {
	return b.title
}

// TitleWithContext, kaynağın kullanıcı arayüzünde görünecek başlığını döner.
//
// Bu metod, SetTitleFunc ile ayarlanan dinamik başlık fonksiyonunu kullanır.
// Eğer titleFunc ayarlanmamışsa, Title() metodunu fallback olarak kullanır.
//
// Parametreler:
// - ctx: Fiber context (i18n için gerekli)
//
// Döndürür:
// - string: Kullanıcı dostu başlık
//
// Örnek:
//
//	title := resource.TitleWithContext(c.Ctx)
func (b *OptimizedBase) TitleWithContext(ctx *fiber.Ctx) string {
	if b.titleFunc != nil && ctx != nil {
		return b.titleFunc(ctx)
	}
	return b.Title()
}

// / SetRepository, resource'un veri erişim katmanını ayarlar.
// /
// / Repository, veritabanı işlemlerini (CRUD, query, pagination) yönetir.
// / Genellikle GORM provider kullanılır ancak özel implementasyonlar da mümkündür.
// /
// / Parametreler:
// / - r: Veri sağlayıcı (data.DataProvider interface)
// /
// / Örnek:
// /   r := &ProductResource{}
// /   r.SetRepository(data.NewGormProvider(db, &models.Product{}))
func (b *OptimizedBase) SetRepository(r data.DataProvider) {
	b.repository = r
}

// / Repository, resource'un veri erişim katmanını döner.
// /
// / Bu metod, OptimizedResource interface'inin bir parçasıdır ve
// / veritabanı işlemleri için repository sağlar.
// /
// / GORM Optimizasyonları:
// / - Otomatik olarak With() metodundan eager loading ilişkilerini yükler
// / - N+1 query problemini önler
// / - Her çağrıda yeni bir provider instance'ı oluşturur (thread-safe)
// /
// / Parametreler:
// / - db: GORM database instance
// /
// / Döndürür:
// / - Veri sağlayıcı (data.DataProvider interface)
// /
// / Örnek:
// /   repo := r.Repository(db)
// /   items, err := repo.List(ctx, query)
// /
// / Detaylı ilişki örnekleri için bkz: [Relationships.md](../../docs/Relationships.md)
func (b *OptimizedBase) Repository(client *gorm.DB) data.DataProvider {
	if client == nil {
		return nil
	}

	// Create new GormDataProvider instance
	provider := data.NewGormDataProvider(client, b.Model())

	// Automatically configure eager loading from With() method
	// This prevents N+1 query problems
	if relationships := b.With(); len(relationships) > 0 {
		provider.SetWith(relationships)
	}

	return provider
}

// / SetCards, dashboard widget'larını ayarlar.
// /
// / Card'lar, dashboard'da istatistik, grafik ve özet bilgiler
// / göstermek için kullanılır.
// /
// / Parametreler:
// / - c: Widget listesi (widget.Card slice)
// /
// / Örnek:
// /   r.SetCards([]widget.Card{
// /       widget.Value("Toplam Ürün", "products.count"),
// /       widget.Trend("Aylık Satış", "sales.monthly"),
// /   })
func (b *OptimizedBase) SetCards(c []widget.Card) {
	b.cards = c
}

// / Cards, dashboard widget'larını döner.
// /
// / Bu metod, OptimizedResource interface'inin bir parçasıdır ve
// / dashboard'da gösterilecek widget'ları sağlar.
// /
// / Döndürür:
// / - Widget listesi (boş olabilir)
// /
// / Örnek:
// /   cards := r.Cards()
// /   for _, card := range cards {
// /       // Card render işlemi
// /   }
func (b *OptimizedBase) Cards() []widget.Card {
	return b.cards
}

// / Fields, resource'un tüm alanlarını döner.
// /
// / Bu metod, ResolveFields metodunu nil context ile çağırır.
// / Dinamik alan çözümlemesi için FieldResolver kullanılır.
// /
// / Döndürür:
// / - Alan listesi (boş olabilir)
// /
// / Not: Context gerektiren dinamik davranış için GetFields(ctx) kullanın.
// /
// / Örnek:
// /   fields := r.Fields()
// /   for _, field := range fields {
// /       // Alan işlemleri
// /   }
// /
// / Detaylı alan örnekleri için bkz: [Fields.md](../../docs/Fields.md)
func (b *OptimizedBase) Fields() []fields.Element {
	return b.ResolveFields(nil)
}

// / Policy, resource'un yetkilendirme politikasını döner.
// /
// / Bu metod, OptimizedResource interface'inin bir parçasıdır ve
// / Authorizable mixin'in GetPolicy metodunu çağırır.
// /
// / Döndürür:
// / - Yetkilendirme politikası (nil olabilir)
// /
// / Örnek:
// /   policy := r.Policy()
// /   if policy != nil && policy.CanView(user, item) {
// /       // Görüntüleme izni var
// /   }
func (b *OptimizedBase) Policy() auth.Policy {
	return b.GetPolicy()
}

// / Visible, resource'un menüde görünür olup olmadığını döner.
// /
// / Bu metod, OptimizedResource interface'inin bir parçasıdır ve
// / Navigable mixin'in IsVisible metodunu çağırır.
// /
// / Döndürür:
// / - Görünürlük durumu (true: görünür, false: gizli)
// /
// / Örnek:
// /   if r.Visible() {
// /       // Menüye ekle
// /   }
func (b *OptimizedBase) Visible() bool {
	return b.IsVisible()
}

// / With, eager loading yapılacak ilişkileri döner.
// /
// / Bu metod, GORM'un Preload özelliği için kullanılacak ilişki
// / isimlerini belirtir. N+1 query problemini önlemek için kullanılır.
// /
// / Döndürür:
// / - İlişki isimleri listesi (boş slice)
// /
// / Not: Varsayılan implementasyon boş slice döner. Override edilebilir.
// /
// / Örnek:
// /   func (r *PostResource) With() []string {
// /       return []string{"Author", "Category", "Tags"}
// /   }
// /
// / Detaylı ilişki örnekleri için bkz: [Relationships.md](../../docs/Relationships.md)
func (b *OptimizedBase) With() []string {
	return []string{}
}

// / Lenses, resource'un özel görünümlerini döner.
// /
// / Lens'ler, aynı resource'un farklı perspektiflerden görüntülenmesini
// / sağlar (örn: "Bekleyen Siparişler", "Tamamlanan Siparişler").
// /
// / Döndürür:
// / - Lens listesi (boş slice)
// /
// / Not: Varsayılan implementasyon boş slice döner. Override edilebilir.
// /
// / Örnek:
// /   func (r *OrderResource) Lenses() []Lens {
// /       return []Lens{
// /           NewLens("Bekleyen", "pending"),
// /           NewLens("Tamamlanan", "completed"),
// /       }
// /   }
func (b *OptimizedBase) Lenses() []Lens {
	return []Lens{}
}

// / GetLenses, resource'un tüm lens'lerini döner.
// /
// / Bu metod, Lenses metodunu çağırır ve sonucu döner.
// / Dinamik lens çözümlemesi için LensResolver kullanılabilir.
// /
// / Döndürür:
// / - Lens listesi (boş olabilir)
// /
// / Örnek:
// /   lenses := r.GetLenses()
// /   for _, lens := range lenses {
// /       // Lens işlemleri
// /   }
func (b *OptimizedBase) GetLenses() []Lens {
	return b.Lenses()
}

// / Icon, menüde gösterilecek ikon adını döner.
// /
// / Bu metod, Navigable mixin'in GetIcon metodunu çağırır.
// /
// / Döndürür:
// / - İkon adı (boş string olabilir)
// /
// / Örnek:
// /   icon := r.Icon()
func (b *OptimizedBase) Icon() string {
	return b.GetIcon()
}

// / Group, menü grubunu döner.
// /
// / Bu metod, Navigable mixin'in GetGroup metodunu çağırır.
// /
// / Döndürür:
// / - Grup adı (boş string olabilir)
// /
// / Örnek:
// /   group := r.Group()
func (b *OptimizedBase) Group() string {
	return b.GetGroup()
}

// GroupWithContext, kaynağın menüde hangi grup altında listeleneceğini belirler.
//
// Bu metod, SetGroupFunc ile ayarlanan dinamik grup fonksiyonunu kullanır.
// Eğer groupFunc ayarlanmamışsa, GetGroup() metodunu fallback olarak kullanır.
//
// Parametreler:
// - ctx: Fiber context (i18n için gerekli)
//
// Döndürür:
// - string: Grup adı
//
// Örnek:
//
//	group := resource.GroupWithContext(c.Ctx)
func (b *OptimizedBase) GroupWithContext(ctx *fiber.Ctx) string {
	fmt.Printf("[DEBUG] GroupWithContext called - groupFunc nil: %v, ctx nil: %v\n", b.groupFunc == nil, ctx == nil)
	if b.groupFunc != nil && ctx != nil {
		result := b.groupFunc(ctx)
		fmt.Printf("[DEBUG] groupFunc returned: %s\n", result)
		return result
	}
	fallback := b.GetGroup()
	fmt.Printf("[DEBUG] Using fallback GetGroup(): %s\n", fallback)
	return fallback
}

// / GetSortable, varsayılan sıralama ayarlarını döner.
// /
// / Bu metod, Navigable mixin'in sortable alanını döner.
// /
// / Döndürür:
// / - Sıralama ayarları listesi (boş olabilir)
// /
// / Örnek:
// /   sortable := r.GetSortable()
func (b *OptimizedBase) GetSortable() []Sortable {
	return b.Navigable.sortable
}

// / GetDialogType, form görünüm tipini döner.
// /
// / Bu metod, Navigable mixin'in GetDialogType metodunu çağırır.
// /
// / Döndürür:
// / - Dialog tipi (DialogTypeModal, DialogTypeDrawer, DialogTypeFullscreen)
// /
// / Örnek:
// /   dialogType := r.GetDialogType()
func (b *OptimizedBase) GetDialogType() DialogType {
	return b.Navigable.GetDialogType()
}

// / SetDialogType, form görünüm tipini ayarlar ve resource'u döner.
// /
// / Bu metod, method chaining için resource pointer'ı döner.
// /
// / Parametreler:
// / - dt: Dialog tipi
// /
// / Döndürür:
// / - Yapılandırılmış resource pointer'ı (method chaining için)
// /
// / Örnek:
// /   r.SetDialogType(DialogTypeDrawer).SetVisible(true)
func (b *OptimizedBase) SetDialogType(dt DialogType) Resource {
	b.Navigable.SetDialogType(dt)
	return b
}

// GetDialogSize, form modal/sheet genişlik preset'ini döner.
func (b *OptimizedBase) GetDialogSize() DialogSize {
	return b.Navigable.GetDialogSize()
}

// SetDialogSize, form modal/sheet genişlik preset'ini ayarlar ve resource'u döner.
func (b *OptimizedBase) SetDialogSize(ds DialogSize) Resource {
	b.Navigable.SetDialogSize(ds)
	return b
}

// GetIndexRowClickAction, index satır tıklama aksiyonunu döner.
func (b *OptimizedBase) GetIndexRowClickAction() IndexRowClickAction {
	return b.Navigable.GetIndexRowClickAction()
}

// SetIndexRowClickAction, index satır tıklama aksiyonunu ayarlar.
func (b *OptimizedBase) SetIndexRowClickAction(action IndexRowClickAction) Resource {
	b.Navigable.SetIndexRowClickAction(action)
	return b
}

// GetIndexPaginationType, index sayfasında kullanılacak pagination tipini döner.
func (b *OptimizedBase) GetIndexPaginationType() IndexPaginationType {
	return b.Navigable.GetIndexPaginationType()
}

// SetIndexPaginationType, index sayfasında kullanılacak pagination tipini ayarlar.
func (b *OptimizedBase) SetIndexPaginationType(paginationType IndexPaginationType) Resource {
	b.Navigable.SetIndexPaginationType(paginationType)
	return b
}

// GetIndexReorderConfig, index reorder ayarlarını döner.
func (b *OptimizedBase) GetIndexReorderConfig() IndexReorderConfig {
	return b.Navigable.GetIndexReorderConfig()
}

// SetIndexReorder, index reorder ayarlarını toplu şekilde ayarlar.
func (b *OptimizedBase) SetIndexReorder(enabled bool, column string) Resource {
	b.Navigable.SetIndexReorder(enabled, column)
	return b
}

// EnableIndexReorder, verilen kolon için index reorder özelliğini etkinleştirir.
func (b *OptimizedBase) EnableIndexReorder(column string) Resource {
	b.Navigable.EnableIndexReorder(column)
	return b
}

// DisableIndexReorder, index reorder özelliğini devre dışı bırakır.
func (b *OptimizedBase) DisableIndexReorder() Resource {
	b.Navigable.DisableIndexReorder()
	return b
}

// / GetFields, belirli bir context'e göre alanları döner.
// /
// / Bu metod, Resolvable mixin'in ResolveFields metodunu çağırır.
// / Context bazlı dinamik alan çözümlemesi sağlar.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Çözümlenmiş alan listesi
// /
// / Örnek:
// /   fields := r.GetFields(ctx)
// /   // Context'e göre farklı alanlar döner
// /
// / Detaylı alan örnekleri için bkz: [Fields.md](../../docs/Fields.md)
func (b *OptimizedBase) GetFields(ctx *context.Context) []fields.Element {
	return b.ResolveFields(ctx)
}

// / GetFieldsWithContext, context ile field'ları resolve eder ve cache'ler.
// /
// / Bu metod, lazy loading yaklaşımı kullanarak field'ları context ile birlikte
// / resolve eder ve her request için cache'ler. Bu sayede i18n translation'lar
// / doğru context ile çalışır.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Cache'lenmiş field listesi
// /
// / Örnek:
// /   fields := r.GetFieldsWithContext(ctx)
// /   // Field'lar context ile resolve edilir ve cache'lenir
// /
// / İş Akışı:
// / 1. Context nil ise eski davranışa fallback (Fields() çağır)
// / 2. ResolveFields ile field'ları context ile resolve et
// / 3. Her field'ı clone edip context cache'ine ekle
// / 4. Cache'lenmiş field'ları döndür
func (b *OptimizedBase) GetFieldsWithContext(ctx *context.Context) []fields.Element {
	if ctx == nil {
		// Fallback to old behavior
		return b.Fields()
	}

	// ResolveFields ile field'ları al (context ile)
	return b.ResolveFields(ctx)
}

// / GetCards, belirli bir context'e göre widget'ları döner.
// /
// / Bu metod, Resolvable mixin'in ResolveCards metodunu çağırır.
// / Context bazlı dinamik widget çözümlemesi sağlar.
// /
// / Parametreler:
// / - ctx: İstek context'i (nil olabilir)
// /
// / Döndürür:
// / - Çözümlenmiş widget listesi
// /
// / Örnek:
// /   cards := r.GetCards(ctx)
// /   // Context'e göre farklı widget'lar döner
func (b *OptimizedBase) GetCards(ctx *context.Context) []widget.Card {
	return b.ResolveCards(ctx)
}

// / GetPolicy, resource'un yetkilendirme politikasını döner.
// /
// / Bu metod, Authorizable mixin'in GetPolicy metodunu çağırır.
// /
// / Döndürür:
// / - Yetkilendirme politikası (nil olabilir)
// /
// / Örnek:
// /   policy := r.GetPolicy()
func (b *OptimizedBase) GetPolicy() auth.Policy {
	return b.Authorizable.GetPolicy()
}

// / ResolveField, bir alanın değerini dinamik olarak hesaplar ve döner.
// /
// / Bu metod, belirtilen alan adına sahip field'ı bulur, item'dan değeri
// / extract eder ve serialize edilmiş değeri döner. Display callback'leri
// / ve diğer transformasyonlar uygulanır.
// /
// / Parametreler:
// / - fieldName: Alan adı (key)
// / - item: Değerin extract edileceği kayıt
// /
// / Döndürür:
// / - Çözümlenmiş ve serialize edilmiş alan değeri
// / - Hata (alan bulunamazsa)
// /
// / Örnek:
// /   value, err := r.ResolveField("price", product)
// /   if err != nil {
// /       // Alan bulunamadı
// /   }
// /   // value: formatlanmış fiyat değeri
// /
// / Kullanım Senaryoları:
// / - API response'larında özel formatlama
// / - Export işlemlerinde değer dönüşümü
// / - Dinamik alan değeri hesaplama
func (b *OptimizedBase) ResolveField(fieldName string, item any) (any, error) {
	for _, field := range b.Fields() {
		if field.GetKey() == fieldName {
			field.Extract(item)
			serialized := field.JsonSerialize()
			if val, ok := serialized["data"]; ok {
				return val, nil
			}
			if val, ok := serialized["value"]; ok {
				return val, nil
			}
			return nil, nil
		}
	}
	return nil, fmt.Errorf("field %s not found", fieldName)
}

// / GetDefaultActions, tüm resource'larda varsayılan olarak bulunan toplu işlemleri döner.
// /
// / Bu metod, her resource'da otomatik olarak kullanılabilir olan temel action'ları sağlar.
// / Resource'lar kendi özel action'larını eklemek için GetActions() metodunu override edebilir.
// /
// / Varsayılan Action'lar:
// / 1. **Seçilenleri Sil**: Checkbox ile seçilen kayıtları siler (destructive, transaction içinde)
// /
// / Döndürür:
// / - Varsayılan action listesi
// /
// / Örnek Kullanım:
// /   // Resource'da default action'ları kullan
// /   func (r *ProductResource) GetActions() []Action {
// /       return r.OptimizedBase.GetDefaultActions()
// /   }
// /
// /   // Default action'lara ek olarak özel action'lar ekle
// /   func (r *ProductResource) GetActions() []Action {
// /       actions := r.OptimizedBase.GetDefaultActions()
// /       actions = append(actions,
// /           action.New("Dışa Aktar").SetIcon("download")...,
// /       )
// /       return actions
// /   }
// /
// / Önemli Notlar:
// / - Default action'lar transaction içinde çalışır (rollback desteği)
// / - Destructive action'lar onay mesajı gerektirir
// / - Her resource için model tipine göre otomatik çalışır
func (b *OptimizedBase) GetDefaultActions() []Action {
	return []Action{
		// Seçilenleri Sil - Tüm resource'larda varsayılan bulk delete action
		action.New("Seçilenleri Sil").
			SetIcon("trash-2").
			SetSlug("delete-selected").
			Destructive().
			Confirm("Seçili kayıtları silmek istediğinizden emin misiniz? Bu işlem geri alınamaz.").
			ConfirmButton("Evet, Sil").
			CancelButton("İptal").
			ShowOnlyOnIndex().
			Handle(func(ctx *action.ActionContext) error {
				// Transaction içinde çalış (hata durumunda rollback)
				return ctx.DB.Transaction(func(tx *gorm.DB) error {
					for _, item := range ctx.Models {
						// Her kaydı sil
						if err := tx.Delete(item).Error; err != nil {
							return fmt.Errorf("kayıt silinirken hata oluştu: %w", err)
						}
					}
					return nil
				})
			}),
	}
}

// / GetActions, resource'un toplu işlemlerini döner.
// /
// / Bu metod, varsayılan olarak GetDefaultActions() metodunu çağırır ve
// / tüm resource'larda bulunan temel action'ları döndürür.
// /
// / Resource'lar özel action'lar eklemek için bu metodu override edebilir:
// /
// / Döndürür:
// / - İşlem listesi (varsayılan: bulk delete action)
// /
// / Örnek 1: Sadece default action'ları kullan
// /   func (r *ProductResource) GetActions() []Action {
// /       return r.OptimizedBase.GetDefaultActions()
// /   }
// /
// / Örnek 2: Default action'lara ek olarak özel action'lar ekle
// /   func (r *ProductResource) GetActions() []Action {
// /       actions := r.OptimizedBase.GetDefaultActions()
// /       actions = append(actions,
// /           action.New("Dışa Aktar").SetIcon("download").Handle(...),
// /           action.New("Toplu Güncelle").SetIcon("edit").Handle(...),
// /       )
// /       return actions
// /   }
// /
// / Örnek 3: Sadece özel action'lar kullan (default action'ları kullanma)
// /   func (r *ProductResource) GetActions() []Action {
// /       return []Action{
// /           action.New("Özel İşlem").SetIcon("star").Handle(...),
// /       }
// /   }
func (b *OptimizedBase) GetActions() []Action {
	return b.GetDefaultActions()
}

// / GetFilters, resource'un filtreleme seçeneklerini döner.
// /
// / Filter'lar, liste görünümünde kayıtları filtrelemek için kullanılır
// / (örn: durum, kategori, tarih aralığı).
// /
// / Döndürür:
// / - Filtre listesi (boş slice)
// /
// / Not: Varsayılan implementasyon boş slice döner. Override edilebilir.
// /
// / Örnek:
// /   func (r *OrderResource) GetFilters() []Filter {
// /       return []Filter{
// /           NewSelectFilter("Durum", "status", statusOptions),
// /           NewDateRangeFilter("Tarih", "created_at"),
// /       }
// /   }
func (b *OptimizedBase) GetFilters() []Filter {
	return []Filter{}
}

// / StoreHandler, dosya yükleme işlemlerini yönetir.
// /
// / Bu metod, form'dan gelen dosyaları işlemek ve saklamak için kullanılır.
// / Varsayılan implementasyon boş string döner, özel dosya işleme için
// / override edilmelidir.
// /
// / Parametreler:
// / - c: İstek context'i
// / - file: Yüklenen dosya header'ı
// / - storagePath: Dosyanın saklanacağı yol
// / - storageURL: Dosyanın erişim URL'i
// /
// / Döndürür:
// / - Dosyanın kaydedildiği yol veya URL
// / - Hata (işlem başarısızsa)
// /
// / Örnek:
// /   func (r *ProductResource) StoreHandler(
// /       c *context.Context,
// /       file *multipart.FileHeader,
// /       storagePath string,
// /       storageURL string,
// /   ) (string, error) {
// /       // Dosyayı S3'e yükle
// /       url, err := s3.Upload(file, storagePath)
// /       if err != nil {
// /           return "", err
// /       }
// /       return url, nil
// /   }
// /
// / Kullanım Senaryoları:
// / - Görsel yükleme ve işleme
// / - Dosya validasyonu
// / - Cloud storage entegrasyonu
// / - Thumbnail oluşturma
func (b *OptimizedBase) StoreHandler(c *context.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	if storagePath == "" {
		storagePath = "./storage/public"
	}
	if storageURL == "" {
		storageURL = "/storage"
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	localPath := filepath.Join(storagePath, filename)

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return "", err
	}

	if err := c.Ctx.SaveFile(file, localPath); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(storageURL, "/"), filename), nil
}

// / NavigationOrder, menüdeki sıralama önceliğini döner.
// /
// / Bu metod, Navigable mixin'in GetNavigationOrder metodunu çağırır.
// /
// / Döndürür:
// / - Sıralama önceliği (varsayılan: 0)
// /
// / Örnek:
// /   order := r.NavigationOrder()
func (b *OptimizedBase) NavigationOrder() int {
	return b.GetNavigationOrder()
}

// OpenAPIEnabled, kaynağın OpenAPI spesifikasyonunda görünüp görünmeyeceğini döner.
//
// Bu metod, kaynağın OpenAPI/Swagger dokümantasyonunda gösterilip gösterilmeyeceğini
// kontrol eder. Varsayılan değer true'dur (tüm kaynaklar OpenAPI'de görünür).
//
// Döndürür:
// - true: OpenAPI spec'te görünür (varsayılan)
// - false: OpenAPI spec'te gizli
//
// Örnek:
//
//	enabled := r.OpenAPIEnabled()
func (b *OptimizedBase) OpenAPIEnabled() bool {
	return !b.openAPIDisabled
}

// SetOpenAPIEnabled, kaynağın OpenAPI görünürlüğünü ayarlar.
//
// Bu metod, kaynağın OpenAPI/Swagger dokümantasyonunda gösterilip gösterilmeyeceğini
// belirler. Method chaining desteği için resource pointer'ı döner.
//
// Parametreler:
// - enabled: true = OpenAPI'de görünür, false = OpenAPI'de gizli
//
// Döndürür:
// - Resource pointer'ı (method chaining için)
//
// Örnek:
//
//	r.SetOpenAPIEnabled(false) // OpenAPI'de gizle
//	r.SetOpenAPIEnabled(true)  // OpenAPI'de göster
func (b *OptimizedBase) SetOpenAPIEnabled(enabled bool) Resource {
	b.openAPIDisabled = !enabled
	return b
}

// SetTitleFunc, resource'un başlığını dinamik olarak ayarlamak için bir fonksiyon belirler.
//
// Bu metod, i18n desteği için kullanılır. Başlık, kullanıcının diline göre
// otomatik olarak çevrilir.
//
// Parametreler:
// - fn: Başlık döndüren fonksiyon (fiber.Ctx alır, string döndürür)
//
// Döndürür:
// - Resource pointer'ı (method chaining için)
//
// Örnek:
//
//	r.SetTitleFunc(func(c *fiber.Ctx) string {
//	    return i18n.Trans(c, "resources.users.title")
//	})
func (o *OptimizedBase) SetTitleFunc(fn func(*fiber.Ctx) string) Resource {
	o.titleFunc = fn
	return o
}

// SetGroupFunc, resource'un grubunu dinamik olarak ayarlamak için bir fonksiyon belirler.
//
// Bu metod, i18n desteği için kullanılır. Grup, kullanıcının diline göre
// otomatik olarak çevrilir.
//
// Parametreler:
// - fn: Grup adı döndüren fonksiyon (fiber.Ctx alır, string döndürür)
//
// Döndürür:
// - Resource pointer'ı (method chaining için)
//
// Örnek:
//
//	r.SetGroupFunc(func(c *fiber.Ctx) string {
//	    return i18n.Trans(c, "resources.groups.system")
//	})
func (o *OptimizedBase) SetGroupFunc(fn func(*fiber.Ctx) string) Resource {
	o.groupFunc = fn
	return o
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
func (o *OptimizedBase) SetRecordTitleKey(key string) Resource {
	o.recordTitleKey = key
	return o
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
func (o *OptimizedBase) GetRecordTitleKey() string {
	if o.recordTitleKey == "" {
		return "id"
	}
	return o.recordTitleKey
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
func (o *OptimizedBase) SetRecordTitleFunc(fn func(record any) string) Resource {
	o.recordTitleFunc = fn
	return o
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
func (o *OptimizedBase) RecordTitle(record any) string {
	// Özel fonksiyon varsa onu kullan
	if o.recordTitleFunc != nil {
		return o.recordTitleFunc(record)
	}

	// Reflection ile field değerini al
	titleKey := o.GetRecordTitleKey()

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
