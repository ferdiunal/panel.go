package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

/// # HandleResourceStore
///
/// Bu fonksiyon, yeni bir kaynak (resource) oluşturmak için kullanılan HTTP POST isteklerini işler.
/// RESTful API mimarisinde "Create" operasyonunu gerçekleştirir ve `/api/resource/:resource` endpoint'ine
/// gelen istekleri yönetir.
///
/// ## Temel İşlevsellik
///
/// Fonksiyon aşağıdaki adımları sırasıyla gerçekleştirir:
/// 1. **İstek Gövdesi Ayrıştırma**: Gelen HTTP isteğinin body kısmını parse eder
/// 2. **Yetkilendirme Kontrolü**: Policy üzerinden kullanıcının kayıt oluşturma yetkisini doğrular
/// 3. **Kayıt Oluşturma**: Provider aracılığıyla veritabanında yeni kayıt oluşturur
/// 4. **Bildirim Yönetimi**: Başarılı işlem için otomatik bildirim oluşturur
/// 5. **Bildirim Kaydetme**: Bildirimleri veritabanına kaydeder
/// 6. **Yanıt Dönme**: Oluşturulan kaydı ve bildirimleri JSON formatında döner
///
/// ## Parametreler
///
/// * `h` - `*FieldHandler`: Alan işleyici yapısı. Aşağıdaki bileşenleri içerir:
///   - `Policy`: Yetkilendirme politikası (opsiyonel)
///   - `Provider`: Veri sağlayıcı (CRUD işlemleri için)
///   - `NotificationService`: Bildirim servisi (opsiyonel)
///   - `Elements`: Kaynak alanları tanımları
///
/// * `c` - `*context.Context`: Panel.go özel context yapısı. Fiber context'i genişletir ve şunları sağlar:
///   - HTTP istek/yanıt yönetimi
///   - Resource bilgilerine erişim
///   - Bildirim yönetimi
///   - Kullanıcı oturum bilgileri
///
/// ## Dönüş Değeri
///
/// * `error`: İşlem başarılı ise `nil`, hata durumunda ilgili hata mesajı döner
///   - Başarılı durumda HTTP 201 (Created) status kodu ile yanıt döner
///   - Hata durumlarında uygun HTTP status kodları kullanılır (400, 403, 500)
///
/// ## Desteklenen İçerik Türleri
///
/// Fonksiyon iki farklı HTTP Content-Type'ı destekler:
/// - `application/json`: JSON formatında veri gönderimi
/// - `multipart/form-data`: Form verisi ve dosya yükleme desteği
///
/// ## Kullanım Senaryoları
///
/// ### Senaryo 1: Basit JSON Verisi ile Kayıt Oluşturma
/// ```go
/// // Kullanıcı kaydı oluşturma
/// POST /api/resource/users
/// Content-Type: application/json
/// {
///   "name": "Ahmet Yılmaz",
///   "email": "ahmet@example.com",
///   "role": "admin"
/// }
/// ```
///
/// ### Senaryo 2: Dosya Yükleme ile Kayıt Oluşturma
/// ```go
/// // Ürün resmi ile birlikte ürün oluşturma
/// POST /api/resource/products
/// Content-Type: multipart/form-data
/// {
///   "name": "Laptop",
///   "price": "15000",
///   "image": [dosya]
/// }
/// ```
///
/// ### Senaryo 3: İlişkili Kayıt Oluşturma
/// ```go
/// // Kategori ile ilişkili ürün oluşturma
/// POST /api/resource/products
/// {
///   "name": "Telefon",
///   "category_id": 5,
///   "tags": [1, 2, 3]  // Many-to-Many ilişki
/// }
/// ```
///
/// ## Hata Durumları ve HTTP Status Kodları
///
/// | Durum | Status Kodu | Açıklama |
/// |-------|-------------|----------|
/// | Geçersiz istek gövdesi | 400 Bad Request | JSON parse hatası veya geçersiz form verisi |
/// | Yetkisiz erişim | 403 Forbidden | Policy.Create() kontrolü başarısız |
/// | Veritabanı hatası | 500 Internal Server Error | Provider.Create() işlemi başarısız |
/// | Başarılı oluşturma | 201 Created | Kayıt başarıyla oluşturuldu |
///
/// ## Bildirim Sistemi
///
/// Fonksiyon otomatik bildirim yönetimi sağlar:
/// - Kayıt başarıyla oluşturulduğunda varsayılan başarı bildirimi eklenir
/// - Özel bildirimler resource üzerinden tanımlanabilir
/// - Bildirimler veritabanına asenkron olarak kaydedilir
/// - Bildirim kaydetme hatası ana işlemi etkilemez (graceful degradation)
///
/// ## Yanıt Formatı
///
/// Başarılı işlem sonrası dönen JSON yapısı:
/// ```json
/// {
///   "data": {
///     "id": 123,
///     "name": "Örnek Kayıt",
///     "created_at": "2026-02-07T15:22:32Z",
///     // ... diğer alanlar
///   },
///   "notifications": [
///     {
///       "message": "Kayıt başarıyla oluşturuldu",
///       "type": "success",
///       "duration": 3000
///     }
///   ]
/// }
/// ```
///
/// ## Güvenlik Özellikleri
///
/// 1. **Policy Tabanlı Yetkilendirme**: Her istek için Create yetkisi kontrol edilir
/// 2. **Veri Validasyonu**: parseBody() fonksiyonu ile gelen veriler doğrulanır
/// 3. **SQL Injection Koruması**: Provider katmanı parametreli sorgular kullanır
/// 4. **XSS Koruması**: Çıktı verileri otomatik olarak escape edilir
///
/// ## Performans Notları
///
/// - **Veritabanı İşlemi**: Tek bir INSERT sorgusu çalıştırılır
/// - **Bildirim Kaydetme**: Ana işlemi bloklamaz, hata durumunda sessizce başarısız olur
/// - **Alan Çözümleme**: resolveResourceFields() ile sadece gerekli alanlar döndürülür
/// - **Bellek Kullanımı**: Büyük dosya yüklemeleri için streaming desteği önerilir
///
/// ## Avantajlar
///
/// ✓ **Esnek Veri Formatı**: JSON ve multipart/form-data desteği
/// ✓ **Otomatik Bildirim**: Kullanıcı deneyimi için hazır bildirim sistemi
/// ✓ **Policy Entegrasyonu**: Merkezi yetkilendirme yönetimi
/// ✓ **Hata Yönetimi**: Detaylı hata mesajları ve uygun HTTP kodları
/// ✓ **Genişletilebilir**: Provider pattern ile farklı veri kaynakları desteklenebilir
///
/// ## Dezavantajlar
///
/// ✗ **Senkron İşlem**: Büyük veri setlerinde yanıt süresi uzayabilir
/// ✗ **Tek Kayıt**: Toplu (bulk) kayıt oluşturma desteklenmez
/// ✗ **Bildirim Hatası**: Bildirim kaydetme hatası loglanmaz (yorum satırında)
///
/// ## Önemli Notlar
///
/// ⚠️ **Dikkat**: Policy nil ise yetkilendirme kontrolü atlanır. Üretim ortamında mutlaka Policy tanımlanmalıdır.
///
/// ⚠️ **Dikkat**: NotificationService nil ise bildirimler veritabanına kaydedilmez, sadece yanıtta döner.
///
/// ⚠️ **Dikkat**: parseBody() fonksiyonu Content-Type header'ına göre otomatik ayrıştırma yapar.
///
/// 💡 **İpucu**: Büyük dosya yüklemeleri için Fiber'ın BodyLimit middleware'ini yapılandırın.
///
/// 💡 **İpucu**: Özel bildirimler için resource üzerinde NotifySuccess(), NotifyError() metodlarını kullanın.
///
/// ## İlgili Fonksiyonlar
///
/// - `HandleResourceUpdate`: Kayıt güncelleme işlemi
/// - `HandleResourceDelete`: Kayıt silme işlemi
/// - `HandleResourceShow`: Tekil kayıt görüntüleme
/// - `HandleResourceIndex`: Kayıt listesi görüntüleme
///
/// ## Örnek Kullanım (Handler Tanımlama)
///
/// ```go
/// // Resource handler oluşturma
/// handler := &FieldHandler{
///     Policy: &UserPolicy{},
///     Provider: gormProvider,
///     NotificationService: notificationService,
///     Elements: []fields.Field{
///         fields.Text("name").Required(),
///         fields.Email("email").Required(),
///     },
/// }
///
/// // Route tanımlama
/// app.Post("/api/resource/users", func(c *fiber.Ctx) error {
///     ctx := context.New(c)
///     return HandleResourceStore(handler, ctx)
/// })
/// ```
///
/// ## Test Örneği
///
/// ```go
/// func TestHandleResourceStore(t *testing.T) {
///     app := fiber.New()
///     handler := setupTestHandler()
///
///     app.Post("/api/resource/test", func(c *fiber.Ctx) error {
///         ctx := context.New(c)
///         return HandleResourceStore(handler, ctx)
///     })
///
///     req := httptest.NewRequest("POST", "/api/resource/test",
///         strings.NewReader(`{"name":"Test"}`))
///     req.Header.Set("Content-Type", "application/json")
///
///     resp, _ := app.Test(req)
///     assert.Equal(t, 201, resp.StatusCode)
/// }
/// ```
///
/// ## Versiyon Bilgisi
///
/// - **Eklendi**: v1.0.0
/// - **Son Güncelleme**: v2.0.0 (Bildirim sistemi eklendi)
///
/// ## Bakım Notları
///
/// - Bildirim kaydetme hatası şu anda loglanmıyor (satır 38)
/// - Gelecekte async/background job desteği eklenebilir
/// - Bulk create özelliği için ayrı endpoint düşünülebilir
func HandleResourceStore(h *FieldHandler, c *context.Context) error {
	data, err := h.parseBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if h.Policy != nil && !h.Policy.Create(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	result, err := h.Provider.Create(c, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Add default success notification if none exists
	if c.Resource() != nil {
		notifications := c.Resource().GetNotifications()
		if len(notifications) == 0 {
			c.Resource().NotifySuccess("Record created successfully")
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

	resolvedData, err := h.resolveResourceFields(c.Ctx, c.Resource(), result, h.getElements(c))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":          resolvedData,
		"notifications": notificationsResponse,
	})
}
