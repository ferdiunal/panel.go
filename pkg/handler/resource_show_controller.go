package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

/// # HandleResourceShow
///
/// Bu fonksiyon, belirli bir kaynağı (resource) ID'sine göre veritabanından getirir ve
/// kullanıcıya JSON formatında döndürür. Fonksiyon, yetkilendirme kontrolü yapar ve
/// kaynak üzerindeki kullanıcı izinlerini (view, update, delete) meta bilgisi olarak döner.
///
/// ## Kullanım Senaryoları
///
/// - **Detay Sayfası**: Bir kaynağın detay sayfasını görüntülemek için kullanılır
/// - **Düzenleme Formu**: Düzenleme formunu doldurmak için mevcut veriyi çekmek
/// - **Önizleme**: Kaynak önizlemesi göstermek için
/// - **API Entegrasyonu**: RESTful API üzerinden tek bir kaynağa erişim sağlamak
///
/// ## Parametreler
///
/// * `h` - `*FieldHandler`: Kaynak işleyici yapısı. Aşağıdaki bileşenleri içerir:
///   - `Provider`: Veritabanı işlemlerini yöneten provider
///   - `Policy`: Yetkilendirme kurallarını kontrol eden policy nesnesi
///   - `Resource`: Kaynak tanımı ve meta bilgileri
///   - `Elements`: Görüntülenecek alan (field) listesi
///
/// * `c` - `*context.Context`: Fiber context'i genişleten özel context yapısı. Şunları içerir:
///   - HTTP request/response bilgileri
///   - Route parametreleri (örn: `id`)
///   - Kullanıcı oturumu ve yetkilendirme bilgileri
///   - Middleware tarafından eklenen özel veriler
///
/// ## Dönüş Değeri
///
/// * `error`: İşlem başarılı ise `nil`, aksi halde hata döner
///   - `404 Not Found`: Kaynak bulunamadığında
///   - `403 Forbidden`: Kullanıcının görüntüleme yetkisi yoksa
///   - `200 OK`: Başarılı durumda JSON response
///
/// ## HTTP Response Yapısı
///
/// ```json
/// {
///   "data": {
///     "id": 1,
///     "name": "Örnek Kaynak",
///     "created_at": "2024-01-01T00:00:00Z",
///     // ... diğer alanlar
///   },
///   "meta": {
///     "title": "Kaynak Başlığı",
///     "policy": {
///       "view": true,
///       "update": true,
///       "delete": false
///     }
///   }
/// }
/// ```
///
/// ## İşlem Akışı
///
/// 1. **ID Çıkarma**: URL'den kaynak ID'si alınır (`/api/resource/:resource/:id`)
/// 2. **Veri Getirme**: Provider üzerinden kaynak veritabanından sorgulanır
/// 3. **Varlık Kontrolü**: Kaynak bulunamazsa 404 hatası döner
/// 4. **Yetkilendirme**: Policy.View() ile görüntüleme yetkisi kontrol edilir
/// 5. **Yetki Kontrolü**: Yetki yoksa 403 hatası döner
/// 6. **Element Belirleme**: Context'ten veya handler'dan görüntülenecek alanlar alınır
/// 7. **Field Resolution**: Alanlar resolve edilir ve formatlanır
/// 8. **Response Oluşturma**: Data ve meta bilgileri ile JSON response döner
///
/// ## Kullanım Örnekleri
///
/// ### Örnek 1: Temel Kullanım
///
/// ```go
/// // Route tanımı
/// app.Get("/api/resource/:resource/:id", func(c *fiber.Ctx) error {
///     ctx := context.FromFiber(c)
///     handler := getFieldHandler(ctx.Resource())
///     return HandleResourceShow(handler, ctx)
/// })
///
/// // HTTP Request
/// // GET /api/resource/users/123
///
/// // Response
/// // {
/// //   "data": { "id": 123, "name": "John Doe", "email": "john@example.com" },
/// //   "meta": { "title": "Users", "policy": { "view": true, "update": true, "delete": false } }
/// // }
/// ```
///
/// ### Örnek 2: Middleware ile Element Filtreleme
///
/// ```go
/// // Middleware ile sadece belirli alanları göster
/// app.Use("/api/resource/:resource/:id", func(c *fiber.Ctx) error {
///     ctx := context.FromFiber(c)
///     // Sadece public alanları göster
///     ctx.Elements = []fields.Element{
///         fields.NewText("name"),
///         fields.NewText("email"),
///     }
///     return c.Next()
/// })
/// ```
///
/// ### Örnek 3: Policy ile Yetkilendirme
///
/// ```go
/// type UserPolicy struct{}
///
/// func (p *UserPolicy) View(c *context.Context, item interface{}) bool {
///     user := item.(*User)
///     currentUser := c.User()
///     // Sadece kendi profilini veya admin görebilir
///     return currentUser.ID == user.ID || currentUser.IsAdmin
/// }
/// ```
///
/// ## Önemli Notlar
///
/// **⚠️ Yetkilendirme Önceliği**
/// - Policy nil ise tüm işlemler için yetki verilir (varsayılan: izin ver)
/// - Üretim ortamında mutlaka Policy tanımlanmalıdır
/// - Policy.View() false dönerse 403 Forbidden hatası verilir
///
/// **⚠️ Element Önceliği**
/// - Context'teki Elements varsa öncelikli olarak kullanılır
/// - Context'te yoksa Handler'daki Elements kullanılır
/// - Bu sayede middleware ile dinamik alan filtreleme yapılabilir
///
/// **⚠️ ID Formatı**
/// - ID string olarak alınır, Provider'ın sorgu tipine göre dönüştürülür
/// - UUID, integer veya string ID desteklenir
/// - Geçersiz ID formatı Provider tarafından hata döner
///
/// **⚠️ Performans**
/// - resolveResourceFields() fonksiyonu lazy loading yapabilir
/// - İlişkili veriler (relationships) gerektiğinde yüklenir
/// - Büyük veri setlerinde N+1 sorgu problemine dikkat edilmelidir
///
/// **⚠️ Hata Yönetimi**
/// - Provider.Show() hatası 404 olarak döner (kaynak bulunamadı)
/// - Policy.View() false dönerse 403 döner (yetkisiz erişim)
/// - Diğer hatalar (DB bağlantı vb.) 500 olarak döner
///
/// ## Avantajlar
///
/// - **Güvenli**: Policy tabanlı yetkilendirme
/// - **Esnek**: Context üzerinden dinamik alan filtreleme
/// - **Tutarlı**: Standart JSON response formatı
/// - **Bilgilendirici**: Meta bilgilerinde izin durumları
/// - **Genişletilebilir**: Middleware ile özelleştirilebilir
///
/// ## Dezavantajlar
///
/// - **N+1 Problem**: İlişkili veriler için dikkatli olunmalı
/// - **Bellek Kullanımı**: Büyük kaynaklarda tüm data memory'e yüklenir
/// - **Cache Yok**: Her istekte veritabanı sorgusu yapılır
///
/// ## İlgili Fonksiyonlar
///
/// - `HandleResourceIndex`: Kaynak listesini getirir
/// - `HandleResourceCreate`: Yeni kaynak oluşturur
/// - `HandleResourceUpdate`: Mevcut kaynağı günceller
/// - `HandleResourceDelete`: Kaynağı siler
/// - `resolveResourceFields`: Alanları resolve eder ve formatlar
///
/// ## Endpoint
///
/// ```
/// GET /api/resource/:resource/:id
/// ```
///
/// ## HTTP Status Kodları
///
/// - `200 OK`: Kaynak başarıyla getirildi
/// - `403 Forbidden`: Görüntüleme yetkisi yok
/// - `404 Not Found`: Kaynak bulunamadı
/// - `500 Internal Server Error`: Sunucu hatası
func HandleResourceShow(h *FieldHandler, c *context.Context) error {
	id := c.Params("id")

	// Determine elements to use (before Provider.Show to extract relationship fields)
	ctx := context.FromFiber(c.Ctx)

	// Set visibility context for proper field filtering
	if ctx != nil {
		ctx.VisibilityCtx = fields.ContextDetail
	}

	var elements []fields.Element
	if ctx != nil && len(ctx.Elements) > 0 {
		elements = ctx.Elements
	} else {
		elements = h.getElements(c)
	}

	// Extract relationship fields from elements and set to provider
	relationshipFields := []fields.RelationshipField{}
	for _, element := range elements {
		if relField, ok := fields.IsRelationshipField(element); ok {
			// relField nil olabilir (IsRelationshipField view'a göre true döndürebilir ama nil relField ile)
			if relField == nil {
				continue
			}

			relationshipFields = append(relationshipFields, relField)
		}
	}

	h.Provider.SetRelationshipFields(relationshipFields)

	// Fetch data
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.View(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return c.JSON(fiber.Map{
		"data": h.resolveResourceFields(c.Ctx, c.Resource(), item, elements),
		"meta": fiber.Map{
			"title": h.Resource.TitleWithContext(c.Ctx),
			"policy": fiber.Map{
				"view":   h.Policy == nil || h.Policy.View(c, item),
				"update": h.Policy == nil || h.Policy.Update(c, item),
				"delete": h.Policy == nil || h.Policy.Delete(c, item),
			},
		},
	})
}
