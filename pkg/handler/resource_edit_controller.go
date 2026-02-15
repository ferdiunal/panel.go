package handler

import (
	"strings"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

/// # HandleResourceEdit
///
/// Bu fonksiyon, belirli bir kaynağın düzenleme formunu hazırlar ve mevcut değerleriyle birlikte
/// form alanlarını döndürür. Kaynak düzenleme işlemlerinde kullanıcı arayüzüne gerekli form
/// yapısını ve mevcut veri değerlerini sağlar.
///
/// ## Temel İşlevsellik
///
/// 1. **Kaynak Yükleme**: Belirtilen ID'ye sahip kaynağı veritabanından yükler
/// 2. **Yetkilendirme Kontrolü**: Policy üzerinden güncelleme yetkisi kontrolü yapar
/// 3. **Alan Filtreleme**: Sadece güncelleme bağlamında görünür alanları filtreler
/// 4. **Değer Çözümleme**: Her alan için mevcut değerleri çözümler ve hazırlar
/// 5. **Sıralı Yanıt**: Alanları tanımlı sıraya göre düzenleyerek döndürür
///
/// ## Parametreler
///
/// * `h` - `*FieldHandler`: Alan işleyici yapısı. Şunları içerir:
///   - `Provider`: Veri sağlayıcı (veritabanı işlemleri için)
///   - `Policy`: Yetkilendirme politikası (opsiyonel)
///   - `Elements`: Form alanları listesi
///   - `Title`: Kaynak başlığı
///
/// * `c` - `*context.Context`: İstek bağlamı. Şunları sağlar:
///   - URL parametrelerine erişim (kaynak ID'si)
///   - HTTP yanıt yönetimi
///   - Kaynak meta verilerine erişim
///
/// ## Dönüş Değeri
///
/// * `error`: İşlem başarılı ise `nil`, aksi halde hata döner
///   - `404 Not Found`: Kaynak bulunamadığında
///   - `403 Forbidden`: Yetkilendirme başarısız olduğunda
///   - `200 OK`: Başarılı durumda JSON yanıt
///
/// ## JSON Yanıt Yapısı
///
/// ```json
/// {
///   "fields": [
///     {
///       "key": "name",
///       "label": "İsim",
///       "type": "text",
///       "value": "Mevcut Değer",
///       "rules": ["required"],
///       ...
///     }
///   ],
///   "meta": {
///     "title": "Kaynak Başlığı"
///   }
/// }
/// ```
///
/// ## Kullanım Senaryoları
///
/// ### Senaryo 1: Standart Kaynak Düzenleme
/// ```go
/// // Router tanımı
/// app.Get("/api/resource/:resource/:id/edit", func(c *fiber.Ctx) error {
///     ctx := context.New(c)
///     return HandleResourceEdit(fieldHandler, ctx)
/// })
/// ```
///
/// ### Senaryo 2: Yetkilendirme ile Düzenleme
/// ```go
/// // Policy tanımlı handler
/// handler := &FieldHandler{
///     Provider: gormProvider,
///     Policy: &UserPolicy{},
///     Elements: []fields.Element{
///         fields.NewText("name").SetLabel("İsim"),
///         fields.NewEmail("email").SetLabel("E-posta"),
///     },
///     Title: "Kullanıcı Düzenle",
/// }
/// ```
///
/// ### Senaryo 3: Özel Alan Bağlamları
/// ```go
/// // Sadece güncelleme formunda görünen alanlar
/// elements := []fields.Element{
///     fields.NewText("name"),                    // Tüm bağlamlarda görünür
///     fields.NewText("slug").HideOnUpdate(),     // Güncelleme formunda gizli
///     fields.NewText("created_at").OnlyOnDetail(), // Sadece detay sayfasında
/// }
/// ```
///
/// ## Alan Filtreleme Mantığı
///
/// Fonksiyon, aşağıdaki bağlamlarda alanları **hariç tutar**:
/// - `HIDE_ON_UPDATE`: Güncelleme formunda gizli
/// - `ONLY_ON_LIST`: Sadece liste görünümünde
/// - `ONLY_ON_DETAIL`: Sadece detay görünümünde
/// - `ONLY_ON_CREATE`: Sadece oluşturma formunda
///
/// Diğer tüm bağlamlar (varsayılan, `SHOW_ON_UPDATE`, vb.) dahil edilir.
///
/// ## Güvenlik Özellikleri
///
/// 1. **Kaynak Doğrulama**: ID'ye göre kaynağın varlığı kontrol edilir
/// 2. **Yetkilendirme**: Policy.Update() metodu ile yetki kontrolü
/// 3. **Alan Görünürlüğü**: IsVisible() ile kaynak bazlı görünürlük kontrolü
/// 4. **Bağlam Filtreleme**: Sadece uygun bağlamdaki alanlar döndürülür
///
/// ## Performans Notları
///
/// * **Veritabanı Sorgusu**: Tek bir Show() çağrısı ile kaynak yüklenir
/// * **Alan Çözümleme**: resolveResourceFields() ile toplu değer çözümleme
/// * **Bellek Kullanımı**: Büyük alan listeleri için filtreleme öncesi yapılır
/// * **Sıralama**: Orijinal element sırası korunur (O(n) karmaşıklık)
///
/// ## Hata Yönetimi
///
/// ### 404 Not Found
/// ```json
/// {
///   "error": "Not found"
/// }
/// ```
/// **Sebep**: Belirtilen ID'ye sahip kaynak bulunamadı
///
/// ### 403 Forbidden
/// ```json
/// {
///   "error": "Unauthorized"
/// }
/// ```
/// **Sebep**: Kullanıcının kaynağı güncelleme yetkisi yok
///
/// ## Önemli Notlar
///
/// * ⚠️ **ID Parametresi**: URL'den alınan ID string formatındadır
/// * ⚠️ **Policy Kontrolü**: Policy nil ise yetkilendirme atlanır
/// * ⚠️ **Alan Sırası**: Yanıttaki alan sırası Elements listesi ile aynıdır
/// * ⚠️ **Değer Çözümleme**: resolveResourceFields() ilişkili verileri de yükler
/// * ⚠️ **Bağlam Hassasiyeti**: Alan bağlamları büyük/küçük harf duyarlıdır
///
/// ## Avantajlar
///
/// * ✅ **Otomatik Filtreleme**: Bağlam bazlı alan filtreleme otomatiktir
/// * ✅ **Güvenli**: Yetkilendirme ve doğrulama katmanları mevcuttur
/// * ✅ **Esnek**: Policy ve Element yapılandırması ile özelleştirilebilir
/// * ✅ **Performanslı**: Tek sorgu ile kaynak ve değerler yüklenir
/// * ✅ **Tutarlı**: Sıralı yanıt yapısı UI tutarlılığını sağlar
///
/// ## Dezavantajlar
///
/// * ❌ **Bellek**: Büyük alan listeleri için filtreleme öncesi tüm liste bellekte tutulur
/// * ❌ **N+1 Sorunu**: İlişkili alanlar için ek sorgular gerekebilir
/// * ❌ **Esneklik**: Özel filtreleme mantığı için kod değişikliği gerekir
///
/// ## İlgili Fonksiyonlar
///
/// * `HandleResourceCreate`: Yeni kaynak oluşturma formu
/// * `HandleResourceUpdate`: Kaynak güncelleme işlemi
/// * `HandleResourceShow`: Kaynak detay görünümü
/// * `resolveResourceFields`: Alan değerlerini çözümleme
///
/// ## Endpoint Bilgisi
///
/// * **Method**: GET
/// * **Path**: `/api/resource/:resource/:id/edit`
/// * **Parametreler**:
///   - `:resource` - Kaynak adı (örn: "users", "posts")
///   - `:id` - Kaynak ID'si (örn: "123", "uuid-string")
///
/// ## Örnek HTTP İsteği
///
/// ```http
/// GET /api/resource/users/123/edit HTTP/1.1
/// Host: example.com
/// Authorization: Bearer <token>
/// ```
///
/// ## Örnek HTTP Yanıtı
///
/// ```http
/// HTTP/1.1 200 OK
/// Content-Type: application/json
///
/// {
///   "fields": [
///     {
///       "key": "name",
///       "label": "İsim",
///       "type": "text",
///       "value": "Ahmet Yılmaz",
///       "rules": ["required", "min:3"]
///     },
///     {
///       "key": "email",
///       "label": "E-posta",
///       "type": "email",
///       "value": "ahmet@example.com",
///       "rules": ["required", "email"]
///     }
///   ],
///   "meta": {
///     "title": "Kullanıcı Düzenle"
///   }
/// }
/// ```
func HandleResourceEdit(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")

	// Set visibility context for proper field filtering
	ctx := c.Resource()
	if ctx != nil {
		ctx.VisibilityCtx = fields.ContextUpdate
	}

	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Update(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Filter for Update
	var updateElements []fields.Element
	elements := h.getElements(c)

	// Nested relationship context'inde parent kaynağa geri dönen ilişki alanlarını gizle.
	if viaResource := strings.TrimSpace(c.Query("viaResource")); viaResource != "" {
		filtered := make([]fields.Element, 0, len(elements))
		for _, element := range elements {
			if shouldSkipViaBackReferenceField(element, viaResource) {
				continue
			}
			filtered = append(filtered, element)
		}
		elements = filtered
	}

	for _, element := range elements {
		if !element.IsVisible(c.Resource()) {
			continue
		}

		// Context string'i space-separated olabilir (örn: "hide_on_create hide_on_update")
		// Bu yüzden string'i parse edip her bir context'i kontrol etmeliyiz
		ctxStr := element.GetContext()
		if ctxStr != "" {
			contexts := strings.Fields(string(ctxStr))
			shouldSkip := false

			for _, ctx := range contexts {
				// Güncelleme formunda gösterilmemesi gereken alanları filtrele
				if ctx == string(fields.HIDE_ON_UPDATE) ||
					ctx == string(fields.ONLY_ON_LIST) ||
					ctx == string(fields.ONLY_ON_DETAIL) ||
					ctx == string(fields.ONLY_ON_CREATE) {
					shouldSkip = true
					break
				}
			}

			if shouldSkip {
				continue
			}
		}

		updateElements = append(updateElements, element)
	}

	// Resolve fields with values
	resolvedMap := h.resolveResourceFields(c.Ctx, c.Resource(), item, updateElements)

	// Convert map to ordered slice based on h.Elements order
	var orderedFields []map[string]interface{}
	for _, element := range updateElements {
		if val, ok := resolvedMap[element.GetKey()]; ok {
			// Cast val to map
			if fieldMap, ok := val.(map[string]interface{}); ok {
				orderedFields = append(orderedFields, fieldMap)
			}
		}
	}

	return c.JSON(fiber.Map{
		"fields": orderedFields,
		"meta": fiber.Map{
			"title": h.Resource.TitleWithContext(c.Ctx),
		},
	})
}
