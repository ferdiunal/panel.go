/// # Context Paketi - Fiber HTTP İstekleri için Tür-Güvenli Bağlam Yönetimi
///
/// Bu paket, Fiber web framework'ü ile çalışırken tür-güvenli bir bağlam (context) sağlar.
/// Kullanıcı, oturum ve kaynak bilgilerine güvenli bir şekilde erişim sağlar.
///
/// ## Temel Özellikler
/// - Fiber.Ctx'i wrap ederek tür-güvenli erişim
/// - Kimlik doğrulanmış kullanıcı bilgisine erişim
/// - Oturum yönetimi
/// - Rol ve izin kontrolü
/// - SSE (Server-Sent Events) streaming desteği
///
/// ## Kullanım Senaryoları
/// 1. HTTP isteklerinde kullanıcı bilgisine erişim
/// 2. Rol tabanlı erişim kontrolü (RBAC)
/// 3. İzin kontrolü
/// 4. Oturum yönetimi
/// 5. Gerçek zamanlı veri akışı (SSE)
package context

import (
	stdcontext "context"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/user"

	"github.com/gofiber/fiber/v2"
)

/// ## FromFiber Fonksiyonu
///
/// Fiber.Ctx'ten ResourceContext'i alır ve döndürür.
///
/// ### Açıklama
/// Bu fonksiyon, Fiber HTTP bağlamından ResourceContext'i çıkartır.
/// ResourceContext, uygulamanın kaynak yönetimi için gerekli bilgileri içerir.
///
/// ### Parametreler
/// - `c *fiber.Ctx`: Fiber HTTP bağlamı
///
/// ### Dönüş Değeri
/// - `*core.ResourceContext`: Kaynak bağlamı (nil olabilir)
///
/// ### Kullanım Örneği
/// ```go
/// func MyHandler(c *fiber.Ctx) error {
///     resourceCtx := FromFiber(c)
///     if resourceCtx == nil {
///         return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
///             "error": "ResourceContext bulunamadı",
///         })
///     }
///     // ResourceContext'i kullan
///     return nil
/// }
/// ```
///
/// ### Önemli Notlar
/// - Eğer ResourceContext Fiber.Locals'ta yoksa nil döndürür
/// - Middleware tarafından önceden ayarlanmış olması gerekir
/// - Type assertion kullanarak güvenli bir şekilde erişim sağlar
func FromFiber(c *fiber.Ctx) *core.ResourceContext {
	val := c.Locals(core.ResourceContextKey)
	if val == nil {
		return nil
	}
	return val.(*core.ResourceContext)
}

/// ## Context Struct'ı
///
/// Fiber.Ctx'i wrap ederek tür-güvenli erişim sağlayan yapı.
///
/// ### Açıklama
/// Context struct'ı, Fiber'ın fiber.Ctx'ini embed ederek, ek metotlar
/// aracılığıyla tür-güvenli erişim sağlar. Bu sayede, Locals'tan veri
/// çekerken type assertion hatalarından kaçınılır.
///
/// ### Avantajları
/// - Tür-güvenli erişim (compile-time kontrol)
/// - Nil pointer dereference'den korunma
/// - Daha temiz ve okunabilir kod
/// - IDE'de otomatik tamamlama desteği
///
/// ### Dezavantajları
/// - Ek bir wrapper katmanı (minimal performans etkisi)
/// - Fiber'ın tüm metotlarına erişim için embed kullanılması gerekir
///
/// ### Kullanım Örneği
/// ```go
/// func MyHandler(c *Context) error {
///     user := c.User()
///     if user == nil {
///         return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
///             "error": "Kullanıcı kimlik doğrulaması gerekli",
///         })
///     }
///     return c.JSON(fiber.Map{
///         "message": "Hoş geldiniz, " + user.Name,
///     })
/// }
/// ```
type Context struct {
	*fiber.Ctx
}

/// ## Handler Tipi
///
/// Özel handler tipi, tür-güvenli Context kullanır.
///
/// ### Açıklama
/// Handler, standart Fiber handler'ı yerine, tür-güvenli Context
/// parametresi alan bir fonksiyon tipidir.
///
/// ### Parametreler
/// - `*Context`: Tür-güvenli HTTP bağlamı
///
/// ### Dönüş Değeri
/// - `error`: İşlem sırasında oluşan hata (nil başarılı demektir)
///
/// ### Kullanım Örneği
/// ```go
/// var getUserHandler Handler = func(c *Context) error {
///     user := c.User()
///     if user == nil {
///         return c.Status(fiber.StatusUnauthorized).SendString("Yetkisiz")
///     }
///     return c.JSON(user)
/// }
/// ```
///
/// ### Avantajları
/// - Tür-güvenli handler tanımı
/// - Daha az boilerplate kod
/// - Middleware zincirinde tutarlılık
type Handler func(*Context) error

/// ## User Metodu
///
/// Kimlik doğrulanmış kullanıcı bilgisini alır.
///
/// ### Açıklama
/// Bu metot, HTTP bağlamının Locals'ından "user" anahtarı altında
/// depolanan kullanıcı bilgisini alır. Type assertion kullanarak
/// güvenli bir şekilde *user.User türüne dönüştürür.
///
/// ### Dönüş Değeri
/// - `*user.User`: Kullanıcı nesnesi (nil olabilir)
///
/// ### Kullanım Örneği
/// ```go
/// func ProfileHandler(c *Context) error {
///     user := c.User()
///     if user == nil {
///         return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
///             "error": "Lütfen giriş yapınız",
///         })
///     }
///     return c.JSON(fiber.Map{
///         "id": user.ID,
///         "email": user.Email,
///         "name": user.Name,
///     })
/// }
/// ```
///
/// ### Önemli Notlar
/// - Middleware tarafından önceden ayarlanmış olması gerekir
/// - Nil döndürülmesi, kullanıcının kimlik doğrulanmadığı anlamına gelir
/// - Type assertion başarısız olursa nil döndürür
///
/// ### Avantajları
/// - Nil-safe erişim
/// - Panic'ten korunma
/// - Temiz hata yönetimi
func (c *Context) User() *user.User {
	if u, ok := c.Locals("user").(*user.User); ok {
		return u
	}
	return nil
}

/// ## Session Metodu
///
/// Aktif oturum bilgisini alır.
///
/// ### Açıklama
/// Bu metot, HTTP bağlamının Locals'ından "session" anahtarı altında
/// depolanan oturum bilgisini alır. Type assertion kullanarak
/// güvenli bir şekilde *session.Session türüne dönüştürür.
///
/// ### Dönüş Değeri
/// - `*session.Session`: Oturum nesnesi (nil olabilir)
///
/// ### Kullanım Örneği
/// ```go
/// func SessionInfoHandler(c *Context) error {
///     sess := c.Session()
///     if sess == nil {
///         return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
///             "error": "Oturum bulunamadı",
///         })
///     }
///     return c.JSON(fiber.Map{
///         "session_id": sess.ID,
///         "created_at": sess.CreatedAt,
///         "expires_at": sess.ExpiresAt,
///     })
/// }
/// ```
///
/// ### Önemli Notlar
/// - Oturum middleware tarafından önceden ayarlanmış olması gerekir
/// - Nil döndürülmesi, oturum olmadığı anlamına gelir
/// - Type assertion başarısız olursa nil döndürür
///
/// ### Avantajları
/// - Nil-safe erişim
/// - Panic'ten korunma
/// - Oturum yönetiminde tutarlılık
func (c *Context) Session() *session.Session {
	if s, ok := c.Locals("session").(*session.Session); ok {
		return s
	}
	return nil
}

/// ## Resource Metodu
///
/// ResourceContext'i alır.
///
/// ### Açıklama
/// Bu metot, FromFiber fonksiyonunu kullanarak ResourceContext'i alır.
/// Kaynak yönetimi ve uygulamanın genel bağlamı için kullanılır.
///
/// ### Dönüş Değeri
/// - `*core.ResourceContext`: Kaynak bağlamı (nil olabilir)
///
/// ### Kullanım Örneği
/// ```go
/// func ResourceHandler(c *Context) error {
///     resource := c.Resource()
///     if resource == nil {
///         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
///             "error": "Kaynak bağlamı bulunamadı",
///         })
///     }
///     // Kaynak bağlamını kullan
///     return c.JSON(fiber.Map{
///         "resource": resource,
///     })
/// }
/// ```
///
/// ### Önemli Notlar
/// - FromFiber fonksiyonunun wrapper'ı olarak çalışır
/// - Nil döndürülmesi, kaynak bağlamının ayarlanmadığı anlamına gelir
func (c *Context) Resource() *core.ResourceContext {
	return FromFiber(c.Ctx)
}

/// ## Wrap Fonksiyonu
///
/// Özel Handler'ı standart Fiber Handler'a dönüştürür.
///
/// ### Açıklama
/// Bu fonksiyon, tür-güvenli Handler'ı Fiber'ın fiber.Handler türüne
/// dönüştürür. Bu sayede, özel handler'lar Fiber'ın middleware zincirinde
/// kullanılabilir.
///
/// ### Parametreler
/// - `h Handler`: Tür-güvenli handler fonksiyonu
///
/// ### Dönüş Değeri
/// - `fiber.Handler`: Fiber uyumlu handler fonksiyonu
///
/// ### Kullanım Örneği
/// ```go
/// app := fiber.New()
///
/// // Özel handler tanımı
/// getUserHandler := func(c *Context) error {
///     user := c.User()
///     if user == nil {
///         return c.Status(fiber.StatusUnauthorized).SendString("Yetkisiz")
///     }
///     return c.JSON(user)
/// }
///
/// // Wrap kullanarak route'a ekle
/// app.Get("/user", Wrap(getUserHandler))
/// ```
///
/// ### Avantajları
/// - Tür-güvenli handler'ları Fiber'da kullanabilme
/// - Middleware zincirinde tutarlılık
/// - Daha az boilerplate kod
///
/// ### Dezavantajları
/// - Ek bir wrapper katmanı (minimal performans etkisi)
func Wrap(h Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return h(&Context{Ctx: c})
	}
}

/// ## HasRole Metodu
///
/// Kimlik doğrulanmış kullanıcının belirli bir role'ü olup olmadığını kontrol eder.
///
/// ### Açıklama
/// Bu metot, kullanıcının belirtilen role'ü olup olmadığını kontrol eder.
/// Admin role'ü her zaman tüm rolleri içerir (admin her şeyi yapabilir).
///
/// ### Parametreler
/// - `role string`: Kontrol edilecek role adı (örn: "editor", "viewer")
///
/// ### Dönüş Değeri
/// - `bool`: Kullanıcının role'ü varsa true, yoksa false
///
/// ### Kullanım Örneği
/// ```go
/// func EditPostHandler(c *Context) error {
///     if !c.HasRole("editor") {
///         return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
///             "error": "Bu işlem için 'editor' role'ü gereklidir",
///         })
///     }
///     // Post'u düzenle
///     return c.JSON(fiber.Map{
///         "message": "Post başarıyla düzenlendi",
///     })
/// }
/// ```
///
/// ### Önemli Notlar
/// - Kullanıcı nil ise false döndürür
/// - Admin role'ü tüm rolleri içerir
/// - Role adı case-sensitive'dir
///
/// ### Avantajları
/// - Basit rol kontrolü
/// - Admin override mekanizması
/// - Temiz ve okunabilir kod
///
/// ### Dezavantajları
/// - Sadece tam rol eşleşmesi kontrol eder
/// - Hiyerarşik rol yapısını desteklemez
/// - Dinamik izin kontrolü için yetersiz
func (c *Context) HasRole(role string) bool {
	u := c.User()
	if u == nil {
		return false
	}
	return u.Role == role || u.Role == "admin"
}

/// ## HasPermission Metodu
///
/// Kimlik doğrulanmış kullanıcının belirli bir işlem için izni olup olmadığını kontrol eder.
///
/// ### Açıklama
/// Bu metot, kullanıcının belirtilen işlem (action) için izni olup olmadığını kontrol eder.
/// Şu anda, admin role'ü tüm işlemlere izin verir ve diğer kullanıcılar da izin alır.
/// (TODO: Gerçek izin mantığı entegre edilmesi gerekir)
///
/// ### Parametreler
/// - `action string`: Kontrol edilecek işlem adı (örn: "create", "delete", "publish")
///
/// ### Dönüş Değeri
/// - `bool`: Kullanıcının işlem için izni varsa true, yoksa false
///
/// ### Kullanım Örneği
/// ```go
/// func DeletePostHandler(c *Context) error {
///     if !c.HasPermission("delete_post") {
///         return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
///             "error": "Bu işlem için yetkiniz yok",
///         })
///     }
///     // Post'u sil
///     return c.JSON(fiber.Map{
///         "message": "Post başarıyla silindi",
///     })
/// }
/// ```
///
/// ### Önemli Notlar
/// - Kullanıcı nil ise false döndürür
/// - Admin role'ü tüm işlemlere izin verir
/// - TODO: Gerçek izin mantığı entegre edilmesi gerekir
/// - Şu anda tüm kimlik doğrulanmış kullanıcılar izin alır
///
/// ### Avantajları
/// - Basit izin kontrolü
/// - Admin override mekanizması
/// - Gelecekte genişletilebilir
///
/// ### Dezavantajları
/// - Şu anda gerçek izin kontrolü yapmaz
/// - Dinamik izin sistemi eksik
/// - Veritabanı sorgusu yapılmaz
///
/// ### Gelecek İyileştirmeler
/// - Veritabanından izin bilgisini oku
/// - Rol tabanlı izin sistemi
/// - Kaynak tabanlı izin sistemi (RBAC)
/// - İzin önbelleği
func (c *Context) HasPermission(action string) bool {
	u := c.User()
	if u == nil {
		return false
	}
	if u.Role == "admin" {
		return true
	}
	// TODO: Gerçek izin mantığını entegre et
	return true
}

/// ## Context Metodu
///
/// Fiber.Ctx'ten temel context.Context'i alır.
///
/// ### Açıklama
/// Bu metot, Fiber'ın fiber.Ctx'inden standart Go context.Context'i alır.
/// Bu, Go'nun context mekanizmasını (timeout, cancellation, vb.) kullanmak için gereklidir.
///
/// ### Dönüş Değeri
/// - `context.Context`: Standart Go context nesnesi
///
/// ### Kullanım Örneği
/// ```go
/// func DatabaseQueryHandler(c *Context) error {
///     ctx := c.Context()
///
///     // Context ile timeout ayarla
///     ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
///     defer cancel()
///
///     // Veritabanı sorgusu yap
///     var user *user.User
///     if err := db.WithContext(ctx).First(&user).Error; err != nil {
///         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
///             "error": "Veritabanı hatası",
///         })
///     }
///     return c.JSON(user)
/// }
/// ```
///
/// ### Önemli Notlar
/// - Fiber'ın context'i request lifecycle'ı ile bağlıdır
/// - Request bittiğinde context cancel edilir
/// - Goroutine'lerde kullanmak için dikkatli olun
///
/// ### Avantajları
/// - Go'nun context mekanizmasını kullanabilme
/// - Timeout ve cancellation desteği
/// - Goroutine'lerde kontrol
///
/// ### Dezavantajları
/// - Request bittiğinde context cancel edilir
/// - Goroutine'lerde context'i kopyalamanız gerekir
func (c *Context) Context() stdcontext.Context {
	return c.Ctx.Context()
}

/// ## Flush Metodu
///
/// Response buffer'ı flush eder (SSE streaming için).
///
/// ### Açıklama
/// Bu metot, HTTP response buffer'ını flush eder. Bu, Server-Sent Events (SSE)
/// gibi streaming protokolleri için gereklidir. Flush yapıldığında, buffer'daki
/// tüm veriler istemciye gönderilir.
///
/// ### Dönüş Değeri
/// - `error`: İşlem sırasında oluşan hata (nil başarılı demektir)
///
/// ### Kullanım Örneği
/// ```go
/// func StreamHandler(c *Context) error {
///     c.Set("Content-Type", "text/event-stream")
///     c.Set("Cache-Control", "no-cache")
///     c.Set("Connection", "keep-alive")
///
///     for i := 0; i < 10; i++ {
///         c.WriteString(fmt.Sprintf("data: Mesaj %d\n\n", i))
///         if err := c.Flush(); err != nil {
///             return err
///         }
///         time.Sleep(1 * time.Second)
///     }
///     return nil
/// }
/// ```
///
/// ### Önemli Notlar
/// - Sadece streaming protokolleri için kullanılır
/// - Response writer Flush() metodunu desteklemiyorsa sessizce başarısız olur
/// - SSE, WebSocket gibi protokollerde gereklidir
///
/// ### Avantajları
/// - Gerçek zamanlı veri akışı
/// - Düşük latency
/// - Tarayıcı uyumluluğu
///
/// ### Dezavantajları
/// - Sadece streaming protokolleri için uygun
/// - Bağlantı açık kalması gerekir
/// - Kaynakları daha fazla tüketir
///
/// ### Kullanım Senaryoları
/// 1. Server-Sent Events (SSE)
/// 2. Gerçek zamanlı bildirimler
/// 3. Canlı veri akışı
/// 4. İlerleme güncellemeleri
func (c *Context) Flush() error {
	// Type assertion: writer'ın Flush() metodunu destekleyip desteklemediğini kontrol et
	if flusher, ok := c.Ctx.Context().Response.BodyWriter().(interface{ Flush() }); ok {
		flusher.Flush()
	}
	return nil
}
