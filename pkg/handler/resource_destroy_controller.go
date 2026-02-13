package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

/// # HandleResourceDestroy
///
/// Bu fonksiyon, belirtilen ID'ye sahip bir kaynağı (resource) sistemden kalıcı olarak siler.
/// DELETE HTTP metodunu kullanarak `/api/resource/:resource/:id` endpoint'ine gelen istekleri işler.
///
/// ## Temel İşleyiş
///
/// Fonksiyon, kaynak silme işlemini güvenli ve kontrollü bir şekilde gerçekleştirir:
/// 1. URL'den kaynak ID'sini alır
/// 2. Kaynağın varlığını doğrular (Policy kontrolü için)
/// 3. Kullanıcının silme yetkisini kontrol eder (Policy)
/// 4. Kaynağı veritabanından siler
/// 5. Başarı bildirimi oluşturur
/// 6. Bildirimleri veritabanına kaydeder
/// 7. JSON yanıtı döndürür
///
/// ## Parametreler
///
/// * `h` - `*FieldHandler`: Kaynak işlemleri için gerekli servisleri içeren handler
///   - `Provider`: Veritabanı işlemlerini yöneten provider (Show, Delete metodları)
///   - `Policy`: Yetkilendirme kontrollerini yapan policy nesnesi (Delete metodu)
///   - `NotificationService`: Bildirim kaydetme servisini sağlayan servis
///
/// * `c` - `*context.Context`: Fiber context'ini genişleten özel context nesnesi
///   - HTTP istek/yanıt bilgilerini içerir
///   - Kaynak (Resource) bilgilerine erişim sağlar
///   - Bildirim (Notification) yönetimi için metodlar sunar
///
/// ## Dönüş Değeri
///
/// * `error`: İşlem başarılı ise `nil`, hata durumunda ilgili hata mesajı
///   - `nil`: Kaynak başarıyla silindi
///   - `fiber.Error`: HTTP hata durumları (404, 403, 500)
///
/// ## HTTP Durum Kodları
///
/// * `200 OK`: Kaynak başarıyla silindi
/// * `404 Not Found`: Belirtilen ID'ye sahip kaynak bulunamadı
/// * `403 Forbidden`: Kullanıcının silme yetkisi yok
/// * `500 Internal Server Error`: Silme işlemi sırasında beklenmeyen hata
///
/// ## Güvenlik Özellikleri
///
/// ### Policy Kontrolü
/// - Kaynak silinmeden önce kullanıcının yetkisi kontrol edilir
/// - Policy tanımlı değilse, tüm kullanıcılar silme yapabilir (dikkatli kullanın!)
/// - Policy.Delete() metodu false dönerse işlem reddedilir
///
/// ### Varlık Kontrolü
/// - Silme işleminden önce kaynağın var olduğu doğrulanır
/// - Olmayan bir kaynağı silmeye çalışmak 404 hatası döndürür
/// - Bu kontrol aynı zamanda policy kontrolü için de gereklidir
///
/// ## Bildirim Sistemi
///
/// ### Otomatik Bildirim
/// - Eğer kaynak için özel bildirim tanımlanmamışsa, otomatik başarı bildirimi eklenir
/// - Varsayılan mesaj: "Record deleted successfully"
///
/// ### Bildirim Kaydetme
/// - Bildirimler NotificationService aracılığıyla veritabanına kaydedilir
/// - Bildirim kaydetme hatası işlemi durdurmaz (silent fail)
/// - Hata durumunda log kaydı yapılabilir (şu an yorum satırında)
///
/// ### Yanıt Formatı
/// Bildirimler JSON yanıtında şu formatta döndürülür:
/// ```json
/// {
///   "message": "Deleted successfully",
///   "notifications": [
///     {
///       "message": "Record deleted successfully",
///       "type": "success",
///       "duration": 3000
///     }
///   ]
/// }
/// ```
///
/// ## Kullanım Senaryoları
///
/// ### Senaryo 1: Basit Kaynak Silme
/// ```go
/// // Kullanıcı bir blog yazısını silmek istiyor
/// // DELETE /api/resource/posts/123
/// // Sonuç: Post silinir, başarı bildirimi gösterilir
/// ```
///
/// ### Senaryo 2: Yetkisiz Silme Denemesi
/// ```go
/// // Kullanıcı başkasına ait bir kaydı silmeye çalışıyor
/// // Policy.Delete() false döndürür
/// // Sonuç: 403 Forbidden, "Unauthorized" mesajı
/// ```
///
/// ### Senaryo 3: Olmayan Kaynak Silme
/// ```go
/// // Kullanıcı silinmiş veya hiç var olmamış bir kaydı silmeye çalışıyor
/// // Provider.Show() hata döndürür
/// // Sonuç: 404 Not Found, "Not found" mesajı
/// ```
///
/// ### Senaryo 4: Özel Bildirimli Silme
/// ```go
/// // Kaynak için özel bildirim tanımlanmış
/// // c.Resource().NotifySuccess("Ürün başarıyla kaldırıldı")
/// // Sonuç: Özel bildirim gösterilir, varsayılan bildirim eklenmez
/// ```
///
/// ## Hata Yönetimi
///
/// ### Provider.Show() Hatası
/// - Kaynak bulunamadığında tetiklenir
/// - HTTP 404 döndürür
/// - Policy kontrolü yapılamaz, işlem sonlanır
///
/// ### Policy.Delete() Reddi
/// - Kullanıcı yetkisiz olduğunda tetiklenir
/// - HTTP 403 döndürür
/// - Kaynak silinmez
///
/// ### Provider.Delete() Hatası
/// - Veritabanı hatası veya constraint ihlali durumunda tetiklenir
/// - HTTP 500 döndürür
/// - Hata mesajı yanıtta döndürülür
///
/// ### NotificationService Hatası
/// - Bildirim kaydetme hatası işlemi durdurmaz
/// - Sessizce göz ardı edilir (silent fail)
/// - İsteğe bağlı olarak log kaydı yapılabilir
///
/// ## Önemli Notlar
///
/// ⚠️ **Policy Kontrolü**: Policy tanımlanmamışsa (`h.Policy == nil`), tüm kullanıcılar
/// silme işlemi yapabilir. Üretim ortamında mutlaka policy tanımlayın!
///
/// ⚠️ **Cascade Silme**: Provider'ın Delete metodu cascade silme yapıp yapmadığına bağlı
/// olarak ilişkili kayıtlar da silinebilir. GORM kullanıyorsanız, model tanımlarınızda
/// `OnDelete` davranışını kontrol edin.
///
/// ⚠️ **Soft Delete**: Provider soft delete kullanıyorsa, kayıt fiziksel olarak silinmez,
/// sadece `deleted_at` alanı güncellenir. Bu durumda kayıt hala veritabanında kalır.
///
/// ⚠️ **Transaction**: Bu fonksiyon transaction yönetimi yapmaz. Eğer silme işlemi
/// birden fazla tablo içeriyorsa, Provider'da transaction kullanmanız önerilir.
///
/// ⚠️ **Audit Log**: Silme işlemleri için audit log tutmak istiyorsanız, bunu
/// Provider.Delete() metodunda veya middleware'de implement etmelisiniz.
///
/// ## Performans Notları
///
/// - **İki Veritabanı Sorgusu**: Show() ve Delete() olmak üzere iki sorgu yapılır
/// - **Optimizasyon**: Eğer policy kontrolü gerekmiyorsa, Show() sorgusu atlanabilir
/// - **Bildirim Kaydetme**: Asenkron yapılabilir (şu an senkron)
///
/// ## Avantajlar
///
/// ✅ **Güvenli Silme**: Policy kontrolü ile yetkisiz silme engellenir
/// ✅ **Kullanıcı Bildirimi**: Otomatik bildirim sistemi ile kullanıcı bilgilendirilir
/// ✅ **Hata Yönetimi**: Tüm hata durumları uygun HTTP kodları ile yönetilir
/// ✅ **Esneklik**: Özel bildirimler tanımlanabilir
/// ✅ **Tutarlılık**: Standart JSON yanıt formatı
///
/// ## Dezavantajlar
///
/// ❌ **İki Sorgu**: Policy kontrolü için ekstra Show() sorgusu gerekir
/// ❌ **Silent Fail**: Bildirim kaydetme hatası sessizce göz ardı edilir
/// ❌ **Transaction Yok**: Atomik işlem garantisi yok
/// ❌ **Senkron Bildirim**: Bildirim kaydetme senkron yapılır, performans etkisi olabilir
///
/// ## Örnek Kullanım
///
/// ```go
/// // Router tanımı
/// app.Delete("/api/resource/:resource/:id", func(c *fiber.Ctx) error {
///     ctx := context.New(c)
///     return HandleResourceDestroy(fieldHandler, ctx)
/// })
///
/// // Policy tanımı
/// type PostPolicy struct{}
///
/// func (p *PostPolicy) Delete(c *context.Context, item interface{}) bool {
///     post := item.(*Post)
///     user := c.User()
///     // Sadece kendi postlarını silebilir
///     return post.UserID == user.ID || user.IsAdmin()
/// }
///
/// // Özel bildirim tanımı
/// resource.OnDeleting(func(c *context.Context, item interface{}) error {
///     c.Resource().NotifySuccess("Blog yazısı başarıyla silindi!")
///     return nil
/// })
/// ```
///
/// ## İlgili Fonksiyonlar
///
/// - `HandleResourceShow`: Tek bir kaynağı görüntüler
/// - `HandleResourceUpdate`: Kaynağı günceller
/// - `HandleResourceIndex`: Kaynak listesini getirir
/// - `Provider.Delete()`: Gerçek silme işlemini yapar
/// - `Policy.Delete()`: Silme yetkisini kontrol eder
///
func HandleResourceDestroy(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")

	// Fetch for Policy Check
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Delete(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.Provider.Delete(c, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Add default success notification if none exists
	if c.Resource() != nil {
		notifications := c.Resource().GetNotifications()
		if len(notifications) == 0 {
			c.Resource().NotifySuccess("Record deleted successfully")
		}
	}

	// Save notifications to database
	if c.Resource() != nil && h.NotificationService != nil {
		if err := h.NotificationService.SaveNotifications(c.Resource()); err != nil {
			// Log error but don't fail the request
			// fmt.Printf("Failed to save notifications: %v\n", err)
		}
	}

	// Get notifications for response
	var notificationsResponse []map[string]interface{}
	if c.Resource() != nil {
		for _, notif := range c.Resource().GetNotifications() {
			notificationsResponse = append(notificationsResponse, map[string]interface{}{
				"message":  notif.Message,
				"type":     notif.Type,
				"duration": notif.Duration,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message":       "Deleted successfully",
		"notifications": notificationsResponse,
	})
}
