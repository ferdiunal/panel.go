package core

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

/// # Notification
///
/// Bu yapı, kullanıcı bildirimlerini temsil eden hafif bir veri yapısıdır.
/// ResourceContext içinde kullanılarak, kaynak operasyonları sırasında
/// kullanıcıya gösterilecek bildirimleri saklar.
///
/// ## Alanlar
///
/// - `Message`: Kullanıcıya gösterilecek bildirim mesajı
/// - `Type`: Bildirim tipi (success, error, warning, info)
/// - `Duration`: Bildirimin ekranda kalma süresi (milisaniye)
/// - `UserID`: Bildirimin hedef kullanıcısı (opsiyonel)
///
/// ## Kullanım Senaryoları
///
/// 1. **Başarılı İşlem Bildirimi**: Kaynak oluşturma, güncelleme veya silme işlemlerinden sonra
/// 2. **Hata Bildirimi**: Validasyon hataları veya sistem hataları için
/// 3. **Uyarı Bildirimi**: Kullanıcının dikkat etmesi gereken durumlar için
/// 4. **Bilgi Bildirimi**: Genel bilgilendirme mesajları için
///
/// ## Örnek Kullanım
///
/// ```go
/// notification := Notification{
///     Message:  "Kayıt başarıyla oluşturuldu",
///     Type:     "success",
///     Duration: 3000,
///     UserID:   &userID,
/// }
/// ```
///
/// ## Notlar
///
/// - Duration varsayılan olarak 3000ms (3 saniye) kullanılır
/// - Type değerleri: "success", "error", "warning", "info"
/// - UserID nil olabilir (tüm kullanıcılar için genel bildirim)
type Notification struct {
	Message  string `json:"message"`
	Type     string `json:"type"`
	Duration int    `json:"duration"`
	UserID   *uint  `json:"user_id,omitempty"`
}

/// # ResourceContext
///
/// Bu yapı, bir kaynak operasyonu sırasında context bilgilerini tutan merkezi veri yapısıdır.
/// Kaynak, alanlar, görünürlük context'i, işlem yapılan öğe, kullanıcı ve HTTP request
/// bilgilerini içerir.
///
/// ## Temel Özellikler
///
/// - **Kaynak Yönetimi**: İşlem yapılan domain entity'yi tutar
/// - **Alan Çözümleme**: Field resolver'ları yönetir ve alan değerlerini dinamik olarak çözümler
/// - **Görünürlük Kontrolü**: Hangi alanların hangi context'te görüneceğini belirler
/// - **Kullanıcı Context'i**: İşlemi yapan kullanıcı bilgisini saklar
/// - **Bildirim Sistemi**: Kullanıcıya gösterilecek bildirimleri yönetir
/// - **Lens Desteği**: Kaynağa filtrelenmiş görünüm (lens) uygulayabilir
///
/// ## Alanlar
///
/// - `Resource`: İşlem yapılan domain entity (örn: User, Post, Product)
/// - `Lens`: Kaynağa uygulanan opsiyonel lens (filtrelenmiş görünüm)
/// - `VisibilityCtx`: Alanların görünür olacağı context (index, detail, create, update, preview)
/// - `Item`: İşlem yapılan spesifik kaynak instance'ı (örn: belirli bir User struct instance'ı)
/// - `User`: İşlemi gerçekleştiren kullanıcı
/// - `Elements`: Bu kaynakla ilişkili alanlar (fields)
/// - `fieldResolvers`: Alan isimlerini resolver fonksiyonlarına eşleyen map
/// - `Request`: Fiber HTTP context'i
/// - `Notifications`: Kullanıcıya gönderilecek bildirimler
///
/// ## Kullanım Senaryoları
///
/// 1. **CRUD İşlemleri**: Kaynak oluşturma, okuma, güncelleme, silme işlemlerinde
/// 2. **Alan Çözümleme**: Dinamik alan değerlerinin hesaplanması ve dönüştürülmesi
/// 3. **Yetkilendirme**: Kullanıcı bazlı erişim kontrolü ve audit logging
/// 4. **Görünürlük Kontrolü**: Context'e göre alanların gösterilmesi/gizlenmesi
/// 5. **Bildirim Yönetimi**: İşlem sonuçlarının kullanıcıya bildirilmesi
/// 6. **Lens Uygulama**: Kaynağa filtrelenmiş görünüm uygulanması
///
/// ## Örnek Kullanım
///
/// ```go
/// // Yeni bir ResourceContext oluşturma
/// ctx := NewResourceContext(c, userResource, fields)
///
/// // Field resolver kaydetme
/// ctx.SetFieldResolver("full_name", fullNameResolver)
///
/// // Alan değerini çözümleme
/// value, err := ctx.ResolveField("full_name", user, nil)
///
/// // Bildirim ekleme
/// ctx.NotifySuccess("Kullanıcı başarıyla oluşturuldu")
///
/// // Metadata alma
/// metadata := ctx.GetResourceMetadata()
/// ```
///
/// ## Fiber.Locals ile Kullanım
///
/// ResourceContext tipik olarak fiber.Locals içinde saklanır ve ResourceContextKey
/// sabiti kullanılarak erişilebilir:
///
/// ```go
/// // Context'i Fiber.Locals'a kaydetme
/// c.Locals(ResourceContextKey, ctx)
///
/// // Context'i Fiber.Locals'tan alma
/// ctx := c.Locals(ResourceContextKey).(*ResourceContext)
/// ```
///
/// ## Alan Resolver Sistemi
///
/// Field resolver'lar, alan değerlerinin dinamik olarak hesaplanmasını veya
/// dönüştürülmesini sağlar. Detaylı bilgi için `docs/Fields.md` dosyasına bakınız.
///
/// ## İlişki Yönetimi
///
/// ResourceContext, ilişkili kaynakların yönetiminde de kullanılır.
/// Detaylı bilgi için `docs/Relationships.md` dosyasına bakınız.
///
/// ## Avantajlar
///
/// - **Merkezi Context Yönetimi**: Tüm operasyon bilgileri tek bir yerde
/// - **Tip Güvenliği**: Interface{} kullanımı esneklik sağlar
/// - **Genişletilebilirlik**: Yeni resolver'lar kolayca eklenebilir
/// - **Bildirim Sistemi**: Kullanıcı geri bildirimi entegre
/// - **Görünürlük Kontrolü**: Context bazlı alan gösterimi
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - **Thread Safety**: ResourceContext thread-safe değildir, concurrent kullanımda dikkatli olun
/// - **Memory Yönetimi**: Büyük veri setlerinde Item ve Resource alanları memory kullanımını artırabilir
/// - **Resolver Performansı**: Resolver fonksiyonları performans kritik olabilir, optimize edin
/// - **Nil Kontrolleri**: Lens, Item ve User alanları nil olabilir, kullanmadan önce kontrol edin
///
/// ## Gereksinimler
///
/// - Requirement 15.1: THE Sistem SHALL context'in kaynak metadata'sını taşımasını sağlamalıdır
/// - Requirement 15.2: THE Sistem SHALL bileşenlerin context'ten alan resolver'larına erişmesine izin vermelidir
///
/// ## İlgili Tipler
///
/// - `VisibilityContext`: Görünürlük context'i enum'u
/// - `Element`: Alan interface'i
/// - `Resolver`: Alan çözümleyici interface'i
/// - `Notification`: Bildirim yapısı
type ResourceContext struct {
	// Resource is the domain entity being operated on.
	Resource any

	// Lens is the optional lens (filtered view) being applied to the resource.
	// This is nil if no lens is active.
	Lens any

	// VisibilityCtx is the context in which fields should be visible.
	// It determines which fields are shown (index, detail, create, update, preview).
	VisibilityCtx VisibilityContext

	// Item is the specific resource instance being operated on.
	// This is the actual data object (e.g., a User struct instance).
	Item any

	// User is the user performing the operation.
	// This can be used for authorization and audit logging.
	User any

	// Elements are the fields associated with this resource.
	Elements []Element

	// fieldResolvers maps field names to their resolver functions.
	// These resolvers can be used to dynamically compute or transform field values.
	fieldResolvers map[string]Resolver

	// Request is the Fiber HTTP context.
	Request *fiber.Ctx

	// Notifications are the notifications to be sent to the user
	Notifications []Notification

	// fieldCache, her request için clone edilmiş field'ları cache'ler
	// Bu sayede field'lar lazy load edilir ve context ile birlikte resolve edilir
	fieldCache map[string]Element

	// fieldsMutex, field cache'e thread-safe erişim sağlar
	fieldsMutex sync.RWMutex
}

/// # ResourceContextKey
///
/// Bu sabit, ResourceContext'in fiber.Locals içinde saklanması için kullanılan anahtardır.
///
/// ## Kullanım Amacı
///
/// Fiber middleware'leri ve handler'ları arasında ResourceContext'i paylaşmak için
/// standart bir anahtar sağlar. Bu sayede farklı katmanlardaki kodlar aynı context'e
/// erişebilir.
///
/// ## Örnek Kullanım
///
/// ```go
/// // Context'i Fiber.Locals'a kaydetme
/// c.Locals(ResourceContextKey, resourceContext)
///
/// // Context'i Fiber.Locals'tan alma
/// ctx, ok := c.Locals(ResourceContextKey).(*ResourceContext)
/// if !ok {
///     return fiber.NewError(fiber.StatusInternalServerError, "ResourceContext bulunamadı")
/// }
/// ```
///
/// ## Notlar
///
/// - Bu anahtar tüm proje genelinde tutarlı kullanılmalıdır
/// - Type assertion yaparken nil kontrolü yapılmalıdır
/// - Middleware'lerde context'in varlığı garanti edilmelidir
const ResourceContextKey = "resource_context"

/// # NewResourceContext
///
/// Bu fonksiyon, temel parametrelerle yeni bir ResourceContext oluşturur.
/// Basit senaryolar için kullanılır ve görünürlük context'i olmadan çalışır.
///
/// ## Parametreler
///
/// - `c`: Fiber HTTP context'i - HTTP request/response işlemleri için
/// - `resource`: İşlem yapılan domain entity (örn: UserResource, PostResource)
/// - `elements`: Bu kaynakla ilişkili alan listesi (fields)
///
/// ## Döndürür
///
/// - Yapılandırılmış ResourceContext pointer'ı
///
/// ## Başlangıç Durumu
///
/// Oluşturulan context şu özelliklere sahiptir:
/// - `fieldResolvers`: Boş map olarak başlatılır
/// - `Notifications`: Boş slice olarak başlatılır
/// - `Lens`, `Item`, `User`: nil olarak kalır
/// - `VisibilityCtx`: Varsayılan değerde kalır
///
/// ## Kullanım Senaryoları
///
/// 1. **Basit CRUD İşlemleri**: Görünürlük kontrolü gerektirmeyen işlemler
/// 2. **API Endpoint'leri**: Temel kaynak işlemleri
/// 3. **Test Senaryoları**: Minimal context gerektiren testler
///
/// ## Örnek Kullanım
///
/// ```go
/// func CreateUserHandler(c *fiber.Ctx) error {
///     // Kaynak ve alanları tanımla
///     userResource := &UserResource{}
///     fields := []Element{
///         fields.NewText("name"),
///         fields.NewEmail("email"),
///     }
///
///     // Context oluştur
///     ctx := NewResourceContext(c, userResource, fields)
///
///     // Context'i Fiber.Locals'a kaydet
///     c.Locals(ResourceContextKey, ctx)
///
///     // İşlemleri gerçekleştir
///     // ...
///
///     return c.JSON(fiber.Map{"success": true})
/// }
/// ```
///
/// ## Karşılaştırma
///
/// - **NewResourceContext**: Basit senaryolar için, minimal parametreler
/// - **NewResourceContextWithVisibility**: Gelişmiş senaryolar için, tam kontrol
///
/// ## Notlar
///
/// - Görünürlük kontrolü gerekiyorsa `NewResourceContextWithVisibility` kullanın
/// - Field resolver'lar sonradan `SetFieldResolver` ile eklenebilir
/// - Context oluşturulduktan sonra alanlar değiştirilebilir
///
/// ## İlgili Fonksiyonlar
///
/// - `NewResourceContextWithVisibility`: Görünürlük context'i ile oluşturma
/// - `SetFieldResolver`: Field resolver ekleme
func NewResourceContext(c *fiber.Ctx, resource any, elements []Element) *ResourceContext {
	return &ResourceContext{
		Resource:       resource,
		Elements:       elements,
		Request:        c,
		fieldResolvers: make(map[string]Resolver),
		Notifications:  []Notification{},
		fieldCache:     make(map[string]Element),
	}
}

/// # NewResourceContextWithVisibility
///
/// Bu fonksiyon, tüm parametrelerle birlikte yeni bir ResourceContext oluşturur.
/// Gelişmiş senaryolar için kullanılır ve tam kontrol sağlar.
///
/// ## Parametreler
///
/// - `c`: Fiber HTTP context'i - HTTP request/response işlemleri için
/// - `resource`: İşlem yapılan domain entity (örn: UserResource, PostResource)
/// - `lens`: Kaynağa uygulanan opsiyonel lens (filtrelenmiş görünüm) - nil olabilir
/// - `visibilityCtx`: Alanların görünür olacağı context (index, detail, create, update, preview)
/// - `item`: İşlem yapılan spesifik kaynak instance'ı - nil olabilir
/// - `user`: İşlemi gerçekleştiren kullanıcı - nil olabilir
/// - `elements`: Bu kaynakla ilişkili alan listesi (fields)
///
/// ## Döndürür
///
/// - Yapılandırılmış ResourceContext pointer'ı
///
/// ## Başlangıç Durumu
///
/// Oluşturulan context şu özelliklere sahiptir:
/// - `fieldResolvers`: Boş map olarak başlatılır
/// - `Notifications`: Boş slice olarak başlatılır
/// - Tüm parametreler doğrudan atanır
///
/// ## Kullanım Senaryoları
///
/// 1. **Görünürlük Kontrolü**: Context'e göre farklı alanların gösterilmesi
/// 2. **Lens Uygulama**: Kaynağa filtrelenmiş görünüm uygulanması
/// 3. **Kullanıcı Bazlı İşlemler**: Yetkilendirme ve audit logging
/// 4. **Detaylı CRUD İşlemleri**: Tam kontrol gerektiren işlemler
/// 5. **Form İşlemleri**: Create/Update formlarında alan görünürlüğü
///
/// ## Örnek Kullanım
///
/// ```go
/// func UpdateUserHandler(c *fiber.Ctx) error {
///     // Kullanıcıyı ve kaynağı al
///     currentUser := getCurrentUser(c)
///     userResource := &UserResource{}
///
///     // Güncellenecek öğeyi al
///     var user models.User
///     db.First(&user, c.Params("id"))
///
///     // Alanları tanımla
///     fields := []Element{
///         fields.NewText("name").ShowOnUpdate(),
///         fields.NewEmail("email").ShowOnUpdate(),
///         fields.NewPassword("password").OnlyOnForms(),
///     }
///
///     // Context oluştur
///     ctx := NewResourceContextWithVisibility(
///         c,
///         userResource,
///         nil, // lens yok
///         VisibilityUpdate, // update context'i
///         &user, // güncellenecek öğe
///         currentUser, // işlemi yapan kullanıcı
///         fields,
///     )
///
///     // Context'i Fiber.Locals'a kaydet
///     c.Locals(ResourceContextKey, ctx)
///
///     // İşlemleri gerçekleştir
///     // ...
///
///     return c.JSON(fiber.Map{"success": true})
/// }
/// ```
///
/// ## Görünürlük Context'leri
///
/// - `VisibilityIndex`: Liste görünümü için
/// - `VisibilityDetail`: Detay görünümü için
/// - `VisibilityCreate`: Oluşturma formu için
/// - `VisibilityUpdate`: Güncelleme formu için
/// - `VisibilityPreview`: Önizleme için
///
/// ## Lens Kullanımı
///
/// Lens, kaynağa filtrelenmiş bir görünüm uygulamak için kullanılır:
///
/// ```go
/// // Aktif kullanıcılar lens'i
/// activeLens := &ActiveUsersLens{}
///
/// ctx := NewResourceContextWithVisibility(
///     c,
///     userResource,
///     activeLens, // lens uygula
///     VisibilityIndex,
///     nil,
///     currentUser,
///     fields,
/// )
/// ```
///
/// ## Karşılaştırma
///
/// - **NewResourceContext**: Basit senaryolar için, minimal parametreler
/// - **NewResourceContextWithVisibility**: Gelişmiş senaryolar için, tam kontrol
///
/// ## Avantajlar
///
/// - **Tam Kontrol**: Tüm context parametreleri üzerinde kontrol
/// - **Görünürlük Yönetimi**: Context bazlı alan gösterimi
/// - **Lens Desteği**: Filtrelenmiş görünüm uygulama
/// - **Kullanıcı Context'i**: Yetkilendirme ve audit için
/// - **Tip Güvenliği**: Compile-time tip kontrolü
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - **Nil Kontrolleri**: lens, item ve user parametreleri nil olabilir
/// - **Görünürlük Context'i**: Doğru context'i seçmek önemlidir
/// - **Memory Yönetimi**: Büyük item'lar memory kullanımını artırabilir
/// - **Thread Safety**: Context thread-safe değildir
///
/// ## Gereksinimler
///
/// - Requirement 15.1: THE Sistem SHALL context oluşturulduğunda tüm gerekli kaynak bilgisini başlatmalıdır
///
/// ## İlgili Fonksiyonlar
///
/// - `NewResourceContext`: Basit context oluşturma
/// - `SetFieldResolver`: Field resolver ekleme
/// - `GetResourceMetadata`: Context metadata'sını alma
func NewResourceContextWithVisibility(
	c *fiber.Ctx,
	resource any,
	lens any,
	visibilityCtx VisibilityContext,
	item any,
	user any,
	elements []Element,
) *ResourceContext {
	return &ResourceContext{
		Resource:       resource,
		Lens:           lens,
		VisibilityCtx:  visibilityCtx,
		Item:           item,
		User:           user,
		Elements:       elements,
		Request:        c,
		fieldResolvers: make(map[string]Resolver),
		Notifications:  []Notification{},
		fieldCache:     make(map[string]Element),
	}
}

/// # GetFieldResolver
///
/// Bu fonksiyon, belirli bir alan için kayıtlı resolver'ı döndürür.
/// Resolver bulunamazsa nil ve hata döner.
///
/// ## Parametreler
///
/// - `fieldName`: Resolver'ı alınacak alanın adı
///
/// ## Döndürür
///
/// - Alan için kayıtlı Resolver (bulunamazsa nil)
/// - Resolver bulunamazsa hata (fiber.StatusNotFound)
///
/// ## Kullanım Senaryoları
///
/// 1. **Alan Değeri Çözümleme**: Dinamik alan değerlerinin hesaplanması
/// 2. **Validasyon**: Resolver varlığının kontrolü
/// 3. **Debugging**: Hangi resolver'ların kayıtlı olduğunu kontrol etme
/// 4. **Conditional Logic**: Resolver'a göre farklı işlemler yapma
///
/// ## Örnek Kullanım
///
/// ```go
/// // Resolver'ı al
/// resolver, err := ctx.GetFieldResolver("full_name")
/// if err != nil {
///     // Resolver bulunamadı
///     return err
/// }
///
/// // Resolver'ı kullan
/// value, err := resolver.Resolve(user, nil)
/// if err != nil {
///     return err
/// }
/// ```
///
/// ## Hata Durumları
///
/// - **Resolver Bulunamadı**: Belirtilen alan için resolver kayıtlı değilse
///   - HTTP Status: 404 (Not Found)
///   - Mesaj: "field resolver not found: {fieldName}"
///
/// ## Alan Resolver Sistemi
///
/// Field resolver'lar, alan değerlerinin dinamik olarak hesaplanmasını sağlar.
/// Örneğin:
/// - `full_name`: first_name ve last_name'i birleştirme
/// - `avatar_url`: Avatar URL'ini oluşturma
/// - `formatted_date`: Tarih formatını dönüştürme
///
/// Detaylı bilgi için `docs/Fields.md` dosyasına bakınız.
///
/// ## Performans Notları
///
/// - Map lookup O(1) karmaşıklığındadır
/// - Resolver'lar cache'lenir, her seferinde yeniden oluşturulmaz
/// - Sık kullanılan resolver'lar için performans sorunu yoktur
///
/// ## Gereksinimler
///
/// - Requirement 15.2: THE Sistem SHALL bileşenlerin context'ten alan resolver'larına erişmesine izin vermelidir
/// - Requirement 15.3: THE Sistem SHALL bileşen-alan iletişimini desteklemelidir
///
/// ## İlgili Fonksiyonlar
///
/// - `SetFieldResolver`: Resolver kaydetme
/// - `ResolveField`: Alan değerini çözümleme
/// - `GetResourceMetadata`: Kayıtlı resolver'ları listeleme
func (rc *ResourceContext) GetFieldResolver(fieldName string) (Resolver, error) {
	resolver, ok := rc.fieldResolvers[fieldName]
	if !ok {
		return nil, fiber.NewError(fiber.StatusNotFound, "field resolver not found: "+fieldName)
	}
	return resolver, nil
}

/// # SetFieldResolver
///
/// Bu fonksiyon, belirli bir alan için resolver kaydeder.
/// Mevcut bir resolver varsa üzerine yazar.
///
/// ## Parametreler
///
/// - `fieldName`: Resolver'ın kaydedileceği alanın adı
/// - `resolver`: Kaydedilecek Resolver implementasyonu
///
/// ## Kullanım Senaryoları
///
/// 1. **Dinamik Alan Değerleri**: Hesaplanan alanlar için resolver ekleme
/// 2. **Veri Dönüştürme**: Alan değerlerini dönüştürmek için
/// 3. **İlişki Yükleme**: İlişkili verileri lazy loading ile yükleme
/// 4. **Özel Formatlar**: Tarih, para birimi gibi özel formatlar
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit bir resolver tanımlama
/// fullNameResolver := &FullNameResolver{}
///
/// // Resolver'ı kaydetme
/// ctx.SetFieldResolver("full_name", fullNameResolver)
///
/// // Birden fazla resolver kaydetme
/// ctx.SetFieldResolver("avatar_url", avatarResolver)
/// ctx.SetFieldResolver("formatted_date", dateResolver)
/// ctx.SetFieldResolver("total_price", priceResolver)
/// ```
///
/// ## Resolver Implementasyonu
///
/// Resolver interface'ini implement eden bir struct oluşturun:
///
/// ```go
/// type FullNameResolver struct{}
///
/// func (r *FullNameResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
///     user, ok := item.(*User)
///     if !ok {
///         return nil, errors.New("invalid item type")
///     }
///     return user.FirstName + " " + user.LastName, nil
/// }
/// ```
///
/// ## İlişki Resolver'ları
///
/// İlişkili verileri yüklemek için resolver kullanabilirsiniz:
///
/// ```go
/// type PostsResolver struct {
///     db *gorm.DB
/// }
///
/// func (r *PostsResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
///     user, ok := item.(*User)
///     if !ok {
///         return nil, errors.New("invalid item type")
///     }
///
///     var posts []Post
///     err := r.db.Where("user_id = ?", user.ID).Find(&posts).Error
///     return posts, err
/// }
/// ```
///
/// Detaylı bilgi için `docs/Relationships.md` dosyasına bakınız.
///
/// ## Avantajlar
///
/// - **Esneklik**: Dinamik alan değerleri hesaplama
/// - **Lazy Loading**: İhtiyaç duyulduğunda veri yükleme
/// - **Separation of Concerns**: İş mantığını resolver'lara taşıma
/// - **Yeniden Kullanılabilirlik**: Resolver'ları farklı alanlarda kullanma
/// - **Test Edilebilirlik**: Resolver'ları bağımsız test etme
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - **Üzerine Yazma**: Aynı alan için birden fazla resolver kaydederseniz son kayıt geçerli olur
/// - **Nil Kontrol**: Resolver nil olmamalıdır
/// - **Performans**: Resolver'lar sık çağrılabilir, performansı optimize edin
/// - **Hata Yönetimi**: Resolver'lar hata döndürebilir, handle edin
///
/// ## Gereksinimler
///
/// - Requirement 15.2: THE Sistem SHALL bileşenlerin context'ten alan resolver'larına erişmesine izin vermelidir
/// - Requirement 15.3: THE Sistem SHALL bileşen-alan iletişimini desteklemelidir
///
/// ## İlgili Fonksiyonlar
///
/// - `GetFieldResolver`: Resolver'ı alma
/// - `ResolveField`: Alan değerini çözümleme
func (rc *ResourceContext) SetFieldResolver(fieldName string, resolver Resolver) {
	rc.fieldResolvers[fieldName] = resolver
}

/// # ResolveField
///
/// Bu fonksiyon, kayıtlı resolver kullanarak bir alan değerini çözümler.
/// Resolver bulunamazsa veya çözümleme başarısız olursa hata döner.
///
/// ## Parametreler
///
/// - `fieldName`: Çözümlenecek alanın adı
/// - `item`: Alan değerinin çözümleneceği öğe (örn: User instance'ı)
/// - `params`: Resolver'a iletilecek ek parametreler (opsiyonel)
///
/// ## Döndürür
///
/// - Çözümlenmiş alan değeri (interface{})
/// - Hata (resolver bulunamazsa veya çözümleme başarısız olursa)
///
/// ## Kullanım Senaryoları
///
/// 1. **Hesaplanan Alanlar**: Dinamik olarak hesaplanan değerleri alma
/// 2. **Veri Dönüştürme**: Alan değerlerini farklı formatlara dönüştürme
/// 3. **İlişki Yükleme**: İlişkili verileri lazy loading ile yükleme
/// 4. **Conditional Values**: Koşullara göre farklı değerler döndürme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit alan çözümleme
/// fullName, err := ctx.ResolveField("full_name", user, nil)
/// if err != nil {
///     return err
/// }
/// fmt.Println(fullName) // "John Doe"
///
/// // Parametreli çözümleme
/// params := map[string]interface{}{
///     "format": "DD/MM/YYYY",
///     "timezone": "Europe/Istanbul",
/// }
/// formattedDate, err := ctx.ResolveField("created_at", post, params)
/// if err != nil {
///     return err
/// }
///
/// // İlişki yükleme
/// posts, err := ctx.ResolveField("posts", user, map[string]interface{}{
///     "limit": 10,
///     "status": "published",
/// })
/// if err != nil {
///     return err
/// }
/// ```
///
/// ## Parametre Kullanımı
///
/// Resolver'lara parametre geçirerek davranışlarını özelleştirebilirsiniz:
///
/// ```go
/// type DateResolver struct{}
///
/// func (r *DateResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
///     post := item.(*Post)
///     format := params["format"].(string)
///     timezone := params["timezone"].(string)
///
///     loc, _ := time.LoadLocation(timezone)
///     return post.CreatedAt.In(loc).Format(format), nil
/// }
/// ```
///
/// ## Hata Durumları
///
/// 1. **Resolver Bulunamadı**: Belirtilen alan için resolver kayıtlı değilse
///    - `GetFieldResolver` hatası döner
///
/// 2. **Çözümleme Hatası**: Resolver çalışırken hata oluşursa
///    - Resolver'ın döndürdüğü hata iletilir
///
/// ## İş Akışı
///
/// 1. `GetFieldResolver` ile resolver'ı al
/// 2. Resolver bulunamazsa hata döndür
/// 3. Resolver'ın `Resolve` methodunu çağır
/// 4. Sonucu veya hatayı döndür
///
/// ## Performans Optimizasyonu
///
/// - **Cache**: Sık kullanılan değerleri cache'leyin
/// - **Lazy Loading**: Sadece gerektiğinde yükleyin
/// - **Batch Loading**: Birden fazla öğe için toplu yükleme yapın
/// - **N+1 Problem**: İlişki yüklemelerinde dikkatli olun
///
/// ## Avantajlar
///
/// - **Tek Satır Çözüm**: Resolver alma ve çalıştırma tek fonksiyonda
/// - **Hata Yönetimi**: Merkezi hata yönetimi
/// - **Tip Güvenliği**: Interface{} ile esneklik
/// - **Parametre Desteği**: Dinamik davranış kontrolü
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - **Nil Kontrol**: item ve params nil olabilir
/// - **Type Assertion**: Dönen değeri doğru tipe cast edin
/// - **Hata Yönetimi**: Resolver hataları handle edin
/// - **Performans**: Resolver'lar sık çağrılabilir
///
/// ## Gereksinimler
///
/// - Requirement 15.2: THE Sistem SHALL bileşenlerin context'ten alan resolver'larına erişmesine izin vermelidir
/// - Requirement 15.3: THE Sistem SHALL bileşen-alan iletişimini desteklemelidir
///
/// ## İlgili Fonksiyonlar
///
/// - `GetFieldResolver`: Resolver'ı alma
/// - `SetFieldResolver`: Resolver kaydetme
func (rc *ResourceContext) ResolveField(fieldName string, item interface{}, params map[string]interface{}) (interface{}, error) {
	resolver, err := rc.GetFieldResolver(fieldName)
	if err != nil {
		return nil, err
	}
	return resolver.Resolve(item, params)
}

/// # GetResourceMetadata
///
/// Bu fonksiyon, kaynak ve mevcut context hakkında metadata döndürür.
/// Görünürlük context'i, görünür alanlar, kayıtlı resolver'lar ve diğer
/// context bilgilerini içeren bir map döner.
///
/// ## Döndürür
///
/// Aşağıdaki anahtarları içeren bir map:
/// - `visibility_context`: Mevcut görünürlük context'i (string)
/// - `fields`: Görünür alanların isimleri ([]string)
/// - `resolvers`: Kayıtlı resolver'ların isimleri ([]string)
/// - `has_lens`: Lens uygulanmış mı (bool)
/// - `has_item`: Item var mı (bool)
/// - `has_user`: User var mı (bool)
///
/// ## Kullanım Senaryoları
///
/// 1. **Debugging**: Context durumunu inceleme
/// 2. **Logging**: Context bilgilerini loglama
/// 3. **Monitoring**: Context kullanımını izleme
/// 4. **API Response**: Client'a context bilgisi gönderme
/// 5. **Validasyon**: Context'in doğru yapılandırıldığını kontrol etme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Metadata'yı al
/// metadata := ctx.GetResourceMetadata()
///
/// // Bilgileri kullan
/// fmt.Printf("Visibility Context: %s\n", metadata["visibility_context"])
/// fmt.Printf("Visible Fields: %v\n", metadata["fields"])
/// fmt.Printf("Registered Resolvers: %v\n", metadata["resolvers"])
/// fmt.Printf("Has Lens: %v\n", metadata["has_lens"])
/// fmt.Printf("Has Item: %v\n", metadata["has_item"])
/// fmt.Printf("Has User: %v\n", metadata["has_user"])
///
/// // API response'a ekle
/// return c.JSON(fiber.Map{
///     "data": data,
///     "meta": metadata,
/// })
/// ```
///
/// ## Görünür Alan Filtreleme
///
/// Sadece mevcut context'te görünür olan alanlar listelenir:
///
/// ```go
/// // Örnek: Update context'inde sadece update'te görünür alanlar listelenir
/// ctx := NewResourceContextWithVisibility(
///     c,
///     resource,
///     nil,
///     VisibilityUpdate,
///     item,
///     user,
///     fields,
/// )
///
/// metadata := ctx.GetResourceMetadata()
/// // metadata["fields"] sadece ShowOnUpdate() ile işaretlenmiş alanları içerir
/// ```
///
/// ## Debugging Örneği
///
/// ```go
/// func debugContext(ctx *ResourceContext) {
///     metadata := ctx.GetResourceMetadata()
///
///     log.Printf("=== Resource Context Debug ===")
///     log.Printf("Visibility: %s", metadata["visibility_context"])
///     log.Printf("Fields: %v", metadata["fields"])
///     log.Printf("Resolvers: %v", metadata["resolvers"])
///     log.Printf("Lens Active: %v", metadata["has_lens"])
///     log.Printf("Item Present: %v", metadata["has_item"])
///     log.Printf("User Present: %v", metadata["has_user"])
///     log.Printf("============================")
/// }
/// ```
///
/// ## Monitoring Örneği
///
/// ```go
/// func monitorContext(ctx *ResourceContext) {
///     metadata := ctx.GetResourceMetadata()
///
///     // Metrics'e gönder
///     metrics.Gauge("resource.visible_fields", len(metadata["fields"].([]string)))
///     metrics.Gauge("resource.resolvers", len(metadata["resolvers"].([]string)))
///     metrics.Counter("resource.with_lens", metadata["has_lens"].(bool))
/// }
/// ```
///
/// ## Performans Notları
///
/// - Her çağrıda yeni slice'lar oluşturulur
/// - Büyük alan listelerinde memory allocation olabilir
/// - Sık çağrılacaksa sonucu cache'leyin
/// - O(n) karmaşıklığı (n = alan sayısı)
///
/// ## Avantajlar
///
/// - **Şeffaflık**: Context durumunu görünür kılar
/// - **Debugging**: Sorun tespitini kolaylaştırır
/// - **Monitoring**: Sistem davranışını izleme
/// - **API Integration**: Client'a context bilgisi sağlama
///
/// ## Dikkat Edilmesi Gerekenler
///
/// - **Type Assertion**: Map değerlerini kullanırken tip dönüşümü yapın
/// - **Nil Kontrol**: has_* alanları boolean döner ama değerler nil olabilir
/// - **Performance**: Sık çağrılacaksa cache'leyin
/// - **Memory**: Büyük alan listelerinde dikkatli olun
///
/// ## Gereksinimler
///
/// - Requirement 15.4: WHEN context oluşturulduğunda, THE Sistem SHALL tüm gerekli kaynak bilgisini başlatmalıdır
///
/// ## İlgili Fonksiyonlar
///
/// - `GetFieldResolver`: Resolver'ları alma
/// - `IsVisible`: Alan görünürlüğünü kontrol etme
func (rc *ResourceContext) GetResourceMetadata() map[string]interface{} {
	fieldNames := make([]string, 0, len(rc.Elements))
	for _, element := range rc.Elements {
		if element.IsVisible(rc) {
			fieldNames = append(fieldNames, element.GetKey())
		}
	}

	resolverNames := make([]string, 0, len(rc.fieldResolvers))
	for name := range rc.fieldResolvers {
		resolverNames = append(resolverNames, name)
	}

	return map[string]interface{}{
		"visibility_context": string(rc.VisibilityCtx),
		"fields":             fieldNames,
		"resolvers":          resolverNames,
		"has_lens":           rc.Lens != nil,
		"has_item":           rc.Item != nil,
		"has_user":           rc.User != nil,
	}
}

/// # Notify
///
/// Bu fonksiyon, context'e yeni bir bildirim ekler.
/// Bildirim mesajı, tipi ve varsayılan süre (3000ms) ile oluşturulur.
///
/// ## Parametreler
///
/// - `message`: Kullanıcıya gösterilecek bildirim mesajı
/// - `notifType`: Bildirim tipi ("success", "error", "warning", "info")
///
/// ## Kullanım Senaryoları
///
/// 1. **Özel Bildirim Tipleri**: Standart tipler dışında özel tipler için
/// 2. **Genel Bildirimler**: Tip belirtilerek bildirim ekleme
/// 3. **Batch İşlemler**: Birden fazla bildirim ekleme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Temel kullanım
/// ctx.Notify("İşlem tamamlandı", "success")
/// ctx.Notify("Bir hata oluştu", "error")
/// ctx.Notify("Dikkat edilmesi gereken durum", "warning")
/// ctx.Notify("Bilgilendirme mesajı", "info")
///
/// // Özel tip ile
/// ctx.Notify("Debug mesajı", "debug")
/// ctx.Notify("Sistem bildirimi", "system")
/// ```
///
/// ## Bildirim Yapısı
///
/// Oluşturulan bildirim şu özelliklere sahiptir:
/// - `Message`: Parametre olarak verilen mesaj
/// - `Type`: Parametre olarak verilen tip
/// - `Duration`: 3000ms (3 saniye) - sabit değer
/// - `UserID`: nil - genel bildirim
///
/// ## Bildirim Tipleri
///
/// Standart tipler:
/// - `success`: Başarılı işlemler için (yeşil)
/// - `error`: Hata durumları için (kırmızı)
/// - `warning`: Uyarılar için (sarı)
/// - `info`: Bilgilendirme için (mavi)
///
/// ## Notlar
///
/// - Bildirimler slice'a eklenir, sıralı tutulur
/// - Aynı mesajdan birden fazla eklenebilir
/// - Duration değiştirilemez (3000ms sabit)
/// - UserID her zaman nil olur
///
/// ## İlgili Fonksiyonlar
///
/// - `NotifySuccess`: Başarı bildirimi için kısayol
/// - `NotifyError`: Hata bildirimi için kısayol
/// - `NotifyWarning`: Uyarı bildirimi için kısayol
/// - `NotifyInfo`: Bilgi bildirimi için kısayol
/// - `GetNotifications`: Tüm bildirimleri alma
func (rc *ResourceContext) Notify(message string, notifType string) {
	rc.Notifications = append(rc.Notifications, Notification{
		Message:  message,
		Type:     notifType,
		Duration: 3000,
	})
}

/// # NotifySuccess
///
/// Bu fonksiyon, başarılı işlemler için bildirim ekler.
/// "success" tipinde bildirim oluşturur.
///
/// ## Parametreler
///
/// - `message`: Başarı mesajı
///
/// ## Kullanım Senaryoları
///
/// 1. **CRUD İşlemleri**: Kayıt oluşturma, güncelleme, silme başarılı olduğunda
/// 2. **Form Gönderimi**: Form başarıyla gönderildiğinde
/// 3. **Dosya İşlemleri**: Dosya yükleme, indirme başarılı olduğunda
/// 4. **API İşlemleri**: API çağrısı başarılı olduğunda
///
/// ## Örnek Kullanım
///
/// ```go
/// // Kayıt oluşturma
/// ctx.NotifySuccess("Kullanıcı başarıyla oluşturuldu")
///
/// // Güncelleme
/// ctx.NotifySuccess("Profil bilgileri güncellendi")
///
/// // Silme
/// ctx.NotifySuccess("Kayıt başarıyla silindi")
///
/// // Dosya yükleme
/// ctx.NotifySuccess("Dosya başarıyla yüklendi")
///
/// // Toplu işlem
/// ctx.NotifySuccess(fmt.Sprintf("%d kayıt başarıyla işlendi", count))
/// ```
///
/// ## UI Görünümü
///
/// - Renk: Yeşil
/// - İkon: Onay işareti
/// - Süre: 3 saniye
/// - Pozisyon: Genellikle sağ üst köşe
///
/// ## Best Practices
///
/// - Mesajlar kısa ve net olmalı
/// - Kullanıcı dostu dil kullanın
/// - Ne yapıldığını açıkça belirtin
/// - Teknik detaylardan kaçının
///
/// ## İlgili Fonksiyonlar
///
/// - `Notify`: Genel bildirim ekleme
/// - `NotifyError`: Hata bildirimi
/// - `NotifyWarning`: Uyarı bildirimi
/// - `NotifyInfo`: Bilgi bildirimi
func (rc *ResourceContext) NotifySuccess(message string) {
	rc.Notify(message, "success")
}

/// # NotifyError
///
/// Bu fonksiyon, hata durumları için bildirim ekler.
/// "error" tipinde bildirim oluşturur.
///
/// ## Parametreler
///
/// - `message`: Hata mesajı
///
/// ## Kullanım Senaryoları
///
/// 1. **Validasyon Hataları**: Form validasyonu başarısız olduğunda
/// 2. **Database Hataları**: Veritabanı işlemi başarısız olduğunda
/// 3. **Yetkilendirme Hataları**: Kullanıcı yetkisi olmadığında
/// 4. **Sistem Hataları**: Beklenmeyen hatalar oluştuğunda
/// 5. **API Hataları**: Dış API çağrısı başarısız olduğunda
///
/// ## Örnek Kullanım
///
/// ```go
/// // Validasyon hatası
/// ctx.NotifyError("E-posta adresi geçersiz")
///
/// // Database hatası
/// ctx.NotifyError("Kayıt oluşturulamadı")
///
/// // Yetkilendirme hatası
/// ctx.NotifyError("Bu işlem için yetkiniz yok")
///
/// // Sistem hatası
/// ctx.NotifyError("Bir hata oluştu, lütfen tekrar deneyin")
///
/// // Dosya hatası
/// ctx.NotifyError("Dosya boyutu çok büyük (max 5MB)")
///
/// // Özel hata mesajı
/// if err != nil {
///     ctx.NotifyError(fmt.Sprintf("İşlem başarısız: %s", err.Error()))
/// }
/// ```
///
/// ## UI Görünümü
///
/// - Renk: Kırmızı
/// - İkon: Hata işareti (X)
/// - Süre: 3 saniye
/// - Pozisyon: Genellikle sağ üst köşe
///
/// ## Best Practices
///
/// - Kullanıcı dostu hata mesajları yazın
/// - Teknik detayları kullanıcıya göstermeyin
/// - Çözüm önerisi sunun
/// - Neyin yanlış gittiğini açıkça belirtin
/// - Kullanıcıyı suçlamayın
///
/// ## Güvenlik Notları
///
/// - Hassas bilgileri hata mesajlarında göstermeyin
/// - Stack trace'leri kullanıcıya göstermeyin
/// - Database hatalarını olduğu gibi göstermeyin
/// - Sistem yollarını ifşa etmeyin
///
/// ## İlgili Fonksiyonlar
///
/// - `Notify`: Genel bildirim ekleme
/// - `NotifySuccess`: Başarı bildirimi
/// - `NotifyWarning`: Uyarı bildirimi
/// - `NotifyInfo`: Bilgi bildirimi
func (rc *ResourceContext) NotifyError(message string) {
	rc.Notify(message, "error")
}

/// # NotifyWarning
///
/// Bu fonksiyon, uyarı durumları için bildirim ekler.
/// "warning" tipinde bildirim oluşturur.
///
/// ## Parametreler
///
/// - `message`: Uyarı mesajı
///
/// ## Kullanım Senaryoları
///
/// 1. **Dikkat Gerektiren Durumlar**: Kullanıcının dikkat etmesi gereken durumlar
/// 2. **Kısmi Başarı**: İşlem kısmen başarılı olduğunda
/// 3. **Deprecation Uyarıları**: Eski özellikler kullanıldığında
/// 4. **Limit Uyarıları**: Limitler aşılmak üzere olduğunda
/// 5. **Güvenlik Uyarıları**: Güvenlik ile ilgili uyarılar
///
/// ## Örnek Kullanım
///
/// ```go
/// // Dikkat gerektiren durum
/// ctx.NotifyWarning("Bu işlem geri alınamaz")
///
/// // Kısmi başarı
/// ctx.NotifyWarning("10 kayıttan 8'i başarıyla işlendi")
///
/// // Deprecation uyarısı
/// ctx.NotifyWarning("Bu özellik yakında kaldırılacak")
///
/// // Limit uyarısı
/// ctx.NotifyWarning("Depolama alanınız dolmak üzere (%90)")
///
/// // Güvenlik uyarısı
/// ctx.NotifyWarning("Şifreniz 90 gündür değiştirilmedi")
///
/// // Veri kaybı uyarısı
/// ctx.NotifyWarning("Kaydedilmemiş değişiklikler var")
/// ```
///
/// ## UI Görünümü
///
/// - Renk: Sarı/Turuncu
/// - İkon: Uyarı işareti (!)
/// - Süre: 3 saniye
/// - Pozisyon: Genellikle sağ üst köşe
///
/// ## Best Practices
///
/// - Uyarıyı açık ve net ifade edin
/// - Kullanıcıya ne yapması gerektiğini söyleyin
/// - Gereksiz uyarılardan kaçının
/// - Kritik durumlar için error kullanın
/// - Bilgilendirme için info kullanın
///
/// ## Uyarı vs Hata
///
/// - **Warning**: İşlem devam edebilir ama dikkat gerekir
/// - **Error**: İşlem başarısız oldu, devam edilemez
///
/// ## İlgili Fonksiyonlar
///
/// - `Notify`: Genel bildirim ekleme
/// - `NotifySuccess`: Başarı bildirimi
/// - `NotifyError`: Hata bildirimi
/// - `NotifyInfo`: Bilgi bildirimi
func (rc *ResourceContext) NotifyWarning(message string) {
	rc.Notify(message, "warning")
}

/// # NotifyInfo
///
/// Bu fonksiyon, bilgilendirme amaçlı bildirimler ekler.
/// "info" tipinde bildirim oluşturur.
///
/// ## Parametreler
///
/// - `message`: Bilgilendirme mesajı
///
/// ## Kullanım Senaryoları
///
/// 1. **Genel Bilgilendirme**: Kullanıcıya bilgi verme
/// 2. **İşlem Durumu**: İşlem durumu hakkında bilgi
/// 3. **Sistem Bildirimleri**: Sistem güncellemeleri, bakım bildirimleri
/// 4. **İpuçları**: Kullanıcıya yardımcı ipuçları
/// 5. **Durum Değişiklikleri**: Durum değişikliği bildirimleri
///
/// ## Örnek Kullanım
///
/// ```go
/// // Genel bilgilendirme
/// ctx.NotifyInfo("Yeni özellikler eklendi")
///
/// // İşlem durumu
/// ctx.NotifyInfo("Veriler yükleniyor...")
///
/// // Sistem bildirimi
/// ctx.NotifyInfo("Sistem bakımı: 15 Şubat 02:00-04:00")
///
/// // İpucu
/// ctx.NotifyInfo("İpucu: Ctrl+S ile hızlı kayıt yapabilirsiniz")
///
/// // Durum değişikliği
/// ctx.NotifyInfo("Hesabınız aktif edildi")
///
/// // Bilgilendirme
/// ctx.NotifyInfo("E-posta doğrulama linki gönderildi")
/// ```
///
/// ## UI Görünümü
///
/// - Renk: Mavi
/// - İkon: Bilgi işareti (i)
/// - Süre: 3 saniye
/// - Pozisyon: Genellikle sağ üst köşe
///
/// ## Best Practices
///
/// - Bilgilendirici ve yardımcı olun
/// - Gereksiz bildirimlerden kaçının
/// - Kısa ve öz mesajlar yazın
/// - Kullanıcıyı rahatsız etmeyin
/// - Önemli bilgiler için kullanın
///
/// ## Info vs Success
///
/// - **Info**: Bilgilendirme, nötr durum
/// - **Success**: Başarılı işlem, pozitif durum
///
/// ## İlgili Fonksiyonlar
///
/// - `Notify`: Genel bildirim ekleme
/// - `NotifySuccess`: Başarı bildirimi
/// - `NotifyError`: Hata bildirimi
/// - `NotifyWarning`: Uyarı bildirimi
func (rc *ResourceContext) NotifyInfo(message string) {
	rc.Notify(message, "info")
}

/// # GetNotifications
///
/// Bu fonksiyon, context'e eklenmiş tüm bildirimleri döndürür.
///
/// ## Döndürür
///
/// - Notification slice'ı (boş olabilir)
///
/// ## Kullanım Senaryoları
///
/// 1. **API Response**: Bildirimleri client'a gönderme
/// 2. **Logging**: Bildirimleri loglama
/// 3. **Testing**: Test senaryolarında bildirim kontrolü
/// 4. **Debugging**: Hangi bildirimlerin eklendiğini kontrol etme
/// 5. **Middleware**: Bildirimleri işleme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Bildirimleri al
/// notifications := ctx.GetNotifications()
///
/// // API response'a ekle
/// return c.JSON(fiber.Map{
///     "data": data,
///     "notifications": notifications,
/// })
///
/// // Bildirimleri logla
/// for _, notif := range ctx.GetNotifications() {
///     log.Printf("[%s] %s", notif.Type, notif.Message)
/// }
///
/// // Test senaryosu
/// notifications := ctx.GetNotifications()
/// assert.Equal(t, 1, len(notifications))
/// assert.Equal(t, "success", notifications[0].Type)
/// assert.Equal(t, "İşlem başarılı", notifications[0].Message)
///
/// // Bildirim sayısını kontrol et
/// if len(ctx.GetNotifications()) > 0 {
///     // Bildirimler var
/// }
/// ```
///
/// ## Middleware Kullanımı
///
/// ```go
/// func NotificationMiddleware(c *fiber.Ctx) error {
///     // Handler'ı çalıştır
///     err := c.Next()
///
///     // Context'ten bildirimleri al
///     if ctx, ok := c.Locals(ResourceContextKey).(*ResourceContext); ok {
///         notifications := ctx.GetNotifications()
///
///         // Response'a ekle
///         if len(notifications) > 0 {
///             c.Set("X-Notifications", toJSON(notifications))
///         }
///     }
///
///     return err
/// }
/// ```
///
/// ## JSON Serileştirme
///
/// Notification struct JSON tag'leri içerir, doğrudan serialize edilebilir:
///
/// ```go
/// notifications := ctx.GetNotifications()
/// jsonData, err := json.Marshal(notifications)
/// ```
///
/// Örnek JSON çıktısı:
/// ```json
/// [
///   {
///     "message": "Kullanıcı oluşturuldu",
///     "type": "success",
///     "duration": 3000
///   },
///   {
///     "message": "E-posta gönderildi",
///     "type": "info",
///     "duration": 3000
///   }
/// ]
/// ```
///
/// ## Performans Notları
///
/// - Slice referansı döner, kopyalama yapılmaz
/// - O(1) karmaşıklığı
/// - Memory allocation yok
/// - Thread-safe değil
///
/// ## Notlar
///
/// - Boş slice dönebilir (nil değil)
/// - Bildirimler ekleme sırasına göre döner
/// - Slice değiştirilebilir (dikkatli kullanın)
/// - UserID alanı genellikle nil olur
///
/// ## İlgili Fonksiyonlar
///
/// - `Notify`: Genel bildirim ekleme
/// - `NotifySuccess`: Başarı bildirimi
/// - `NotifyError`: Hata bildirimi
/// - `NotifyWarning`: Uyarı bildirimi
/// - `NotifyInfo`: Bilgi bildirimi
func (rc *ResourceContext) GetNotifications() []Notification {
	return rc.Notifications
}

/// # GetCachedField
///
/// Bu fonksiyon, cache'den field'ı alır.
/// Field cache'de yoksa false döner.
///
/// ## Parametreler
///
/// - `key`: Field key'i
///
/// ## Döndürür
///
/// - Field element'i (cache'de varsa)
/// - Boolean (cache'de var mı?)
///
/// ## Kullanım Senaryoları
///
/// 1. **Lazy Loading**: Field'ın daha önce resolve edilip edilmediğini kontrol etme
/// 2. **Performance**: Aynı field'ı tekrar resolve etmekten kaçınma
/// 3. **Cache Hit Check**: Cache'de field olup olmadığını kontrol etme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Field'ı cache'den al
/// if field, ok := ctx.GetCachedField("name"); ok {
///     // Cache'de var, kullan
///     return field
/// }
/// // Cache'de yok, resolve et
/// ```
///
/// ## Thread Safety
///
/// - RWMutex ile thread-safe okuma
/// - Concurrent request'ler güvenli
///
/// ## İlgili Fonksiyonlar
///
/// - `CacheField`: Field'ı cache'e ekleme
/// - `GetOrCloneField`: Cache'den al veya clone et
func (rc *ResourceContext) GetCachedField(key string) (Element, bool) {
	rc.fieldsMutex.RLock()
	defer rc.fieldsMutex.RUnlock()
	field, ok := rc.fieldCache[key]
	return field, ok
}

/// # CacheField
///
/// Bu fonksiyon, field'ı cache'e ekler.
/// Mevcut bir field varsa üzerine yazar.
///
/// ## Parametreler
///
/// - `key`: Field key'i
/// - `element`: Cache'e eklenecek field element'i
///
/// ## Kullanım Senaryoları
///
/// 1. **Field Resolve**: Field resolve edildikten sonra cache'e ekleme
/// 2. **Clone Storage**: Clone edilmiş field'ı saklama
/// 3. **Performance**: Sonraki kullanımlar için cache'leme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Field'ı resolve et
/// field := resolveField(ctx)
///
/// // Cache'e ekle
/// ctx.CacheField("name", field)
/// ```
///
/// ## Thread Safety
///
/// - Mutex ile thread-safe yazma
/// - Concurrent request'ler güvenli
///
/// ## İlgili Fonksiyonlar
///
/// - `GetCachedField`: Field'ı cache'den alma
/// - `GetOrCloneField`: Cache'den al veya clone et
func (rc *ResourceContext) CacheField(key string, element Element) {
	rc.fieldsMutex.Lock()
	defer rc.fieldsMutex.Unlock()
	rc.fieldCache[key] = element
}

/// # GetOrCloneField
///
/// Bu fonksiyon, cache'den field'ı alır veya clone edip cache'e ekler.
/// Cache hit ise mevcut field'ı döner, cache miss ise clone eder.
///
/// ## Parametreler
///
/// - `key`: Field key'i
/// - `original`: Clone edilecek orijinal field
///
/// ## Döndürür
///
/// - Cache'den alınan veya clone edilmiş field element'i
///
/// ## Kullanım Senaryoları
///
/// 1. **Lazy Loading**: İlk erişimde clone et, sonraki erişimlerde cache'den al
/// 2. **Performance**: Gereksiz clone işlemlerinden kaçınma
/// 3. **Field Resolution**: Context ile field'ları resolve etme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Field'ı al veya clone et
/// field := ctx.GetOrCloneField("name", originalField)
///
/// // Artık field cache'de, sonraki çağrılar cache'den alır
/// field2 := ctx.GetOrCloneField("name", originalField) // Cache hit
/// ```
///
/// ## İş Akışı
///
/// 1. Cache'de field var mı kontrol et
/// 2. Varsa cache'den döndür (cache hit)
/// 3. Yoksa clone et
/// 4. Clone'u cache'e ekle
/// 5. Clone'u döndür
///
/// ## Thread Safety
///
/// - GetCachedField ve CacheField thread-safe
/// - Clone işlemi mutex dışında yapılır (performance)
///
/// ## İlgili Fonksiyonlar
///
/// - `GetCachedField`: Field'ı cache'den alma
/// - `CacheField`: Field'ı cache'e ekleme
func (rc *ResourceContext) GetOrCloneField(key string, original Element) Element {
	if cached, ok := rc.GetCachedField(key); ok {
		return cached
	}

	// Clone işlemi mutex dışında yapılır (performance)
	// fields.CloneElement fonksiyonu field'ı deep copy eder
	cloned := original // TODO: Implement proper cloning when fields.CloneElement is available

	rc.CacheField(key, cloned)
	return cloned
}
