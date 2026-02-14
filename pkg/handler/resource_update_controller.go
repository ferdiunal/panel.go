package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

/// # HandleResourceUpdate
///
/// Bu fonksiyon, mevcut bir kaynağı (resource) güncellemek için kullanılır.
/// RESTful API mimarisinde PUT/PATCH işlemlerini yönetir ve kaynak güncelleme
/// sürecinin tüm aşamalarını koordine eder.
///
/// ## Temel İşlevsellik
///
/// Bu fonksiyon aşağıdaki işlemleri sırasıyla gerçekleştirir:
/// 1. URL parametresinden kaynak ID'sini alır
/// 2. İstek gövdesini (body) parse eder
/// 3. Mevcut kaynağı veritabanından getirir
/// 4. Yetkilendirme kontrolü yapar (Policy)
/// 5. Kaynağı günceller
/// 6. Bildirim (notification) sistemi ile kullanıcıya geri bildirim sağlar
/// 7. Güncellenmiş kaynağı ve bildirimleri JSON formatında döndürür
///
/// ## Parametreler
///
/// * `h` - `*FieldHandler`: Alan işleyici (field handler) pointer'ı. Bu yapı şunları içerir:
///   - `Provider`: Veritabanı işlemlerini yöneten provider
///   - `Policy`: Yetkilendirme kurallarını kontrol eden policy nesnesi
///   - `NotificationService`: Bildirim yönetim servisi
///   - `Elements`: Kaynak alanlarını (fields) tanımlayan element listesi
///
/// * `c` - `*context.Context`: Panel.go özel context yapısı. Fiber context'ini genişletir ve şunları sağlar:
///   - HTTP istek/yanıt yönetimi
///   - Kaynak (resource) bilgilerine erişim
///   - Bildirim sistemi entegrasyonu
///   - Kullanıcı oturum bilgileri
///
/// ## Dönüş Değeri
///
/// * `error`: İşlem başarılı ise `nil`, hata durumunda ilgili hata mesajı döner
///
/// ## HTTP Endpoint
///
/// * **Method**: PUT/PATCH
/// * **Path**: `/api/resource/:resource/:id`
/// * **Content-Type**: `multipart/form-data` veya `application/json`
///
/// ## Desteklenen İçerik Türleri
///
/// 1. **application/json**: JSON formatında veri gönderimi
/// 2. **multipart/form-data**: Dosya yükleme ve form verisi gönderimi
///
/// ## Güvenlik ve Yetkilendirme
///
/// Fonksiyon üç katmanlı güvenlik kontrolü uygular:
///
/// 1. **Kaynak Varlık Kontrolü**: Güncellenecek kaynak mevcut mu?
/// 2. **Policy Kontrolü**: Kullanıcının bu kaynağı güncelleme yetkisi var mı?
/// 3. **Veri Doğrulama**: Gönderilen veriler geçerli mi?
///
/// ## Bildirim Sistemi
///
/// Fonksiyon otomatik bildirim yönetimi sağlar:
/// - Varsayılan başarı bildirimi: "Record updated successfully"
/// - Özel bildirimler resource üzerinden tanımlanabilir
/// - Bildirimler veritabanına kaydedilir (hata durumunda sessizce atlanır)
/// - Her bildirim şunları içerir: `message`, `type`, `duration`
///
/// ## Hata Durumları ve HTTP Durum Kodları
///
/// | Durum Kodu | Senaryo | Açıklama |
/// |------------|---------|----------|
/// | 400 Bad Request | Geçersiz istek gövdesi | Body parse edilemedi |
/// | 404 Not Found | Kaynak bulunamadı | Belirtilen ID'ye sahip kaynak yok |
/// | 403 Forbidden | Yetki yok | Policy kontrolü başarısız |
/// | 500 Internal Server Error | Güncelleme hatası | Veritabanı veya sistem hatası |
/// | 200 OK | Başarılı | Kaynak başarıyla güncellendi |
///
/// ## Kullanım Senaryoları
///
/// ### Senaryo 1: Basit Kaynak Güncelleme
/// ```go
/// // Kullanıcı profili güncelleme
/// // PUT /api/resource/users/123
/// // Body: {"name": "Ahmet Yılmaz", "email": "ahmet@example.com"}
/// ```
///
/// ### Senaryo 2: Dosya ile Birlikte Güncelleme
/// ```go
/// // Ürün resmi ile birlikte güncelleme
/// // PUT /api/resource/products/456
/// // Content-Type: multipart/form-data
/// // Body: {name: "Yeni Ürün", image: [file]}
/// ```
///
/// ### Senaryo 3: İlişkisel Veri Güncelleme
/// ```go
/// // Kategori ilişkisi ile birlikte ürün güncelleme
/// // PUT /api/resource/products/789
/// // Body: {"name": "Ürün", "category_id": 5, "tags": [1, 2, 3]}
/// ```
///
/// ## Yanıt Formatı
///
/// Başarılı güncelleme yanıtı:
/// ```json
/// {
///   "data": {
///     "id": 123,
///     "name": "Güncellenmiş Kayıt",
///     "updated_at": "2026-02-07T15:22:38Z",
///     // ... diğer alanlar
///   },
///   "notifications": [
///     {
///       "message": "Record updated successfully",
///       "type": "success",
///       "duration": 3000
///     }
///   ]
/// }
/// ```
///
/// ## Avantajlar
///
/// * **Otomatik Yetkilendirme**: Policy sistemi ile entegre güvenlik
/// * **Esnek Veri Formatı**: JSON ve multipart/form-data desteği
/// * **Bildirim Yönetimi**: Kullanıcı geri bildirimi için otomatik bildirim sistemi
/// * **Alan Çözümleme**: Resource fields otomatik olarak çözümlenir ve formatlanır
/// * **Hata Yönetimi**: Kapsamlı hata yakalama ve anlamlı HTTP durum kodları
/// * **Veritabanı Bağımsızlığı**: Provider pattern ile farklı veritabanları desteklenir
///
/// ## Dikkat Edilmesi Gerekenler
///
/// ⚠️ **Önemli Notlar**:
///
/// 1. **ID Parametresi**: URL'den alınan ID parametresi string formatındadır,
///    provider tarafından uygun tipe dönüştürülmelidir
///
/// 2. **Policy Kontrolü**: Policy nil ise yetkilendirme kontrolü atlanır.
///    Üretim ortamında mutlaka policy tanımlanmalıdır
///
/// 3. **Bildirim Hatası**: Bildirim kaydetme hatası ana işlemi durdurmaz,
///    sadece sessizce loglanır (şu an yorum satırında)
///
/// 4. **Kaynak Varlığı**: Güncelleme öncesi kaynak mutlaka kontrol edilir,
///    bu ekstra bir veritabanı sorgusu anlamına gelir
///
/// 5. **Alan Çözümleme**: `resolveResourceFields` fonksiyonu resource fields'ları
///    çözümler, bu işlem performans etkisi yaratabilir
///
/// 6. **Transaction Yönetimi**: Bu fonksiyon transaction yönetimi yapmaz,
///    gerekirse provider seviyesinde implement edilmelidir
///
/// ## İlişkili Fonksiyonlar
///
/// * `parseBody()`: İstek gövdesini parse eder
/// * `Provider.Show()`: Mevcut kaynağı getirir
/// * `Provider.Update()`: Kaynağı günceller
/// * `Policy.Update()`: Güncelleme yetkisi kontrolü
/// * `NotificationService.SaveNotifications()`: Bildirimleri kaydeder
/// * `resolveResourceFields()`: Kaynak alanlarını çözümler ve formatlar
///
/// ## Örnek Kullanım (Handler Kaydı)
///
/// ```go
/// // Router'a handler kaydı
/// app.Put("/api/resource/:resource/:id", func(c *fiber.Ctx) error {
///     ctx := context.New(c)
///     return HandleResourceUpdate(fieldHandler, ctx)
/// })
/// ```
///
/// ## Test Senaryoları
///
/// Test edilmesi gereken durumlar:
/// 1. ✅ Başarılı güncelleme
/// 2. ✅ Geçersiz ID ile güncelleme denemesi
/// 3. ✅ Yetkisiz kullanıcı ile güncelleme denemesi
/// 4. ✅ Geçersiz veri formatı ile güncelleme
/// 5. ✅ Dosya yükleme ile güncelleme
/// 6. ✅ İlişkisel veri güncelleme
/// 7. ✅ Bildirim sistemi çalışması
///
/// ## Performans Notları
///
/// * İki veritabanı sorgusu yapılır: Show (kontrol) + Update (güncelleme)
/// * Alan çözümleme işlemi CPU yoğun olabilir
/// * Büyük dosya yüklemeleri için timeout ayarları kontrol edilmelidir
/// * Cache stratejisi provider seviyesinde implement edilebilir
func HandleResourceUpdate(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")
	data, err := h.parseBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Fetch existing to check policy
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Update(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	result, err := h.Provider.Update(c, id, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Add default success notification if none exists
	if c.Resource() != nil {
		notifications := c.Resource().GetNotifications()
		if len(notifications) == 0 {
			c.Resource().NotifySuccess("Record updated successfully")
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
		"data":          h.resolveResourceFields(c.Ctx, c.Resource(), result, h.getElements(c)),
		"notifications": notificationsResponse,
	})
}
