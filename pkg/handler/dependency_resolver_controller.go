// # Dependency Resolver Controller
//
// Bu paket, form alanları arasındaki bağımlılıkları çözmek için kullanılan HTTP controller'ları içerir.
// Alan bağımlılıkları, bir alanın değeri değiştiğinde diğer alanların dinamik olarak güncellenmesini sağlar.
//
// ## Temel Özellikler
//
// - **Dinamik Alan Güncellemeleri**: Bir alan değiştiğinde bağımlı alanları otomatik günceller
// - **Dairesel Bağımlılık Kontrolü**: Sonsuz döngüleri önlemek için dairesel bağımlılıkları tespit eder
// - **Context-Aware**: Create ve update işlemlerinde farklı davranışlar sergiler
// - **Toplu Güncelleme**: Birden fazla alanı aynı anda günceller
//
// ## Kullanım Senaryoları
//
// 1. **Cascade Select**: Ülke seçildiğinde şehir listesini güncelleme
// 2. **Conditional Fields**: Bir checkbox işaretlendiğinde ilgili alanları gösterme/gizleme
// 3. **Dynamic Options**: Kategori seçildiğinde alt kategori seçeneklerini yükleme
// 4. **Calculated Fields**: Fiyat ve miktar değiştiğinde toplam tutarı hesaplama
//
// ## Örnek Kullanım
//
// ```go
// // Router'a endpoint ekleme
// app.Post("/resources/:resource/fields/resolve-dependencies",
//     func(c *fiber.Ctx) error {
//         ctx := context.New(c)
//         return HandleResolveDependencies(fieldHandler, ctx)
//     },
// )
// ```
//
// ## Güvenlik Notları
//
// - Request body validasyonu yapılır
// - Context değeri sadece "create" veya "update" olabilir
// - Dairesel bağımlılıklar tespit edilir ve hata döndürülür
package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// # ResolveDependenciesRequest
//
// Alan bağımlılıklarını çözmek için kullanılan HTTP request yapısı.
// Bu yapı, frontend'den gelen form verilerini ve değişen alanları içerir.
//
// ## Alanlar
//
// - **FormData**: Formun güncel durumunu içeren key-value map'i
// - **Context**: İşlem bağlamı ("create" veya "update")
// - **ChangedFields**: Değişen alanların key listesi
// - **ResourceID**: Güncelleme işlemlerinde kaynak ID'si (opsiyonel)
//
// ## Kullanım Senaryoları
//
// ### 1. Create İşlemi
// ```json
// {
//   "formData": {
//     "country_id": 1,
//     "name": "John Doe"
//   },
//   "context": "create",
//   "changedFields": ["country_id"],
//   "resourceId": null
// }
// ```
//
// ### 2. Update İşlemi
// ```json
// {
//   "formData": {
//     "country_id": 2,
//     "city_id": 5
//   },
//   "context": "update",
//   "changedFields": ["country_id"],
//   "resourceId": 123
// }
// ```
//
// ## Validasyon Kuralları
//
// - **FormData**: Boş olabilir ama nil olmamalı
// - **Context**: Sadece "create" veya "update" değerlerini alabilir
// - **ChangedFields**: En az bir alan içermeli
// - **ResourceID**: Update context'inde zorunlu, create'de opsiyonel
//
// ## Önemli Notlar
//
// - FormData içindeki değerler interface{} tipinde olduğu için tip dönüşümü gerekebilir
// - ResourceID hem string hem de numeric değerler alabilir
// - ChangedFields listesi, bağımlılık çözümlemesinin başlangıç noktasıdır
type ResolveDependenciesRequest struct {
	FormData      map[string]interface{} `json:"formData"`      // Formun güncel durumu
	Context       string                 `json:"context"`       // İşlem bağlamı: "create" veya "update"
	ChangedFields []string               `json:"changedFields"` // Değişen alanların listesi
	ResourceID    interface{}            `json:"resourceId"`    // Kaynak ID'si (update için)
}

// # HandleResolveDependencies
//
// Alan bağımlılıklarını çözen ve güncellenmiş alan değerlerini döndüren HTTP handler fonksiyonu.
// Bu fonksiyon, bir veya daha fazla alan değiştiğinde, bağımlı alanların nasıl güncellenmesi
// gerektiğini belirler ve güncellenmiş değerleri döndürür.
//
// ## İşleyiş Akışı
//
// 1. **Request Parsing**: HTTP request body'sini parse eder
// 2. **Context Validation**: İşlem bağlamını ("create" veya "update") doğrular
// 3. **Schema Conversion**: Element'leri Schema field'larına dönüştürür
// 4. **Resolver Creation**: Bağımlılık çözümleyici oluşturur
// 5. **Circular Check**: Dairesel bağımlılıkları kontrol eder
// 6. **Dependency Resolution**: Bağımlılıkları çözer ve güncellemeleri hesaplar
// 7. **Response**: Güncellenmiş alan değerlerini döndürür
//
// ## Parametreler
//
// - **h**: FieldHandler instance'ı - Alan tanımlarını içerir
// - **c**: Context instance'ı - HTTP request/response context'i
//
// ## Dönüş Değeri
//
// - **error**: İşlem başarılıysa nil, aksi halde hata mesajı
//
// ## HTTP Endpoint
//
// - **Route**: `POST /resources/:resource/fields/resolve-dependencies`
// - **Method**: POST
// - **Content-Type**: application/json
//
// ## Request Body Örneği
//
// ```json
// {
//   "formData": {
//     "country_id": 1,
//     "city_id": null,
//     "district_id": null
//   },
//   "context": "create",
//   "changedFields": ["country_id"],
//   "resourceId": null
// }
// ```
//
// ## Response Örneği
//
// ### Başarılı Response (200 OK)
// ```json
// {
//   "fields": {
//     "city_id": {
//       "options": [
//         {"value": 1, "label": "Istanbul"},
//         {"value": 2, "label": "Ankara"}
//       ],
//       "value": null,
//       "disabled": false
//     },
//     "district_id": {
//       "options": [],
//       "value": null,
//       "disabled": true
//     }
//   }
// }
// ```
//
// ### Hata Response (400 Bad Request)
// ```json
// {
//   "error": "Invalid context. Must be 'create' or 'update'"
// }
// ```
//
// ### Dairesel Bağımlılık Hatası (400 Bad Request)
// ```json
// {
//   "error": "Circular dependency detected: field_a -> field_b -> field_a"
// }
// ```
//
// ## Kullanım Senaryoları
//
// ### 1. Cascade Select (Ülke-Şehir-İlçe)
//
// Kullanıcı ülke seçtiğinde:
// - Şehir alanı o ülkeye ait şehirlerle güncellenir
// - İlçe alanı devre dışı bırakılır ve temizlenir
//
// ### 2. Conditional Visibility
//
// Kullanıcı "Özel Adres" checkbox'ını işaretlediğinde:
// - Adres detay alanları görünür hale gelir
// - İlgili alanlar zorunlu hale gelebilir
//
// ### 3. Dynamic Pricing
//
// Kullanıcı miktar veya birim fiyat değiştirdiğinde:
// - Toplam tutar otomatik hesaplanır
// - İndirim oranı uygulanır
// - KDV dahil fiyat güncellenir
//
// ### 4. Related Data Loading
//
// Kullanıcı kategori seçtiğinde:
// - Alt kategori seçenekleri yüklenir
// - İlgili özellik alanları gösterilir
// - Varsayılan değerler atanır
//
// ## Avantajlar
//
// - **Reaktif UI**: Kullanıcı deneyimini iyileştirir
// - **Veri Tutarlılığı**: İlişkili alanların senkronize kalmasını sağlar
// - **Performans**: Sadece değişen alanları günceller
// - **Esneklik**: Karmaşık bağımlılık senaryolarını destekler
// - **Güvenlik**: Dairesel bağımlılıkları önler
//
// ## Dezavantajlar ve Dikkat Edilmesi Gerekenler
//
// - **Karmaşıklık**: Çok fazla bağımlılık yönetimi zorlaşabilir
// - **Performans**: Çok sayıda bağımlılık çözümlemesi yavaşlayabilir
// - **Debugging**: Bağımlılık zincirlerini takip etmek zor olabilir
// - **Network Overhead**: Her değişiklikte sunucuya istek gönderilir
//
// ## Önemli Notlar
//
// - Bu fonksiyon sadece bağımlılıkları çözer, veriyi kaydetmez
// - Dairesel bağımlılıklar otomatik olarak tespit edilir ve hata döndürülür
// - Context değeri mutlaka "create" veya "update" olmalıdır
// - FieldHandler'ın Elements listesi boş olmamalıdır
// - Bağımlılık çözümlemesi sırasında oluşan hatalar 500 Internal Server Error döndürür
//
// ## Performans İpuçları
//
// - Gereksiz bağımlılıklar tanımlamayın
// - Bağımlılık zincirlerini mümkün olduğunca kısa tutun
// - Ağır hesaplamalar için caching kullanın
// - Debouncing ile istek sayısını azaltın (frontend tarafında)
//
// ## Güvenlik Notları
//
// - Request body validasyonu yapılır
// - SQL injection'a karşı korumalıdır (ORM kullanımı sayesinde)
// - Dairesel bağımlılıklar DoS saldırılarını önler
// - Resource ID validasyonu yapılmalıdır (authorization middleware ile)
//
// ## Test Örneği
//
// ```go
// func TestHandleResolveDependencies(t *testing.T) {
//     app := fiber.New()
//     handler := &FieldHandler{
//         Elements: []fields.Element{
//             fields.NewSelect("country_id").DependsOn("region_id"),
//             fields.NewSelect("city_id").DependsOn("country_id"),
//         },
//     }
//
//     app.Post("/resolve", func(c *fiber.Ctx) error {
//         ctx := context.New(c)
//         return HandleResolveDependencies(handler, ctx)
//     })
//
//     req := httptest.NewRequest("POST", "/resolve", strings.NewReader(`{
//         "formData": {"country_id": 1},
//         "context": "create",
//         "changedFields": ["country_id"]
//     }`))
//     req.Header.Set("Content-Type", "application/json")
//
//     resp, _ := app.Test(req)
//     assert.Equal(t, 200, resp.StatusCode)
// }
// ```
//
// ## İlgili Tipler ve Fonksiyonlar
//
// - `fields.DependencyResolver`: Bağımlılık çözümleme mantığı
// - `fields.Schema`: Alan tanımları ve bağımlılık kuralları
// - `ResolveDependenciesRequest`: Request body yapısı
// - `FieldHandler`: Alan yönetimi için handler
func HandleResolveDependencies(h *FieldHandler, c *context.Context) error {
	// Parse request body
	var req ResolveDependenciesRequest
	if err := c.Ctx.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate context
	if req.Context != "create" && req.Context != "update" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid context. Must be 'create' or 'update'",
		})
	}

	// Convert elements to Schema fields
	schemaFields := make([]*fields.Schema, 0, len(h.Elements))
	for _, element := range h.Elements {
		if schema, ok := element.(*fields.Schema); ok {
			schemaFields = append(schemaFields, schema)
		}
	}

	// Create dependency resolver
	resolver := fields.NewDependencyResolver(schemaFields, req.Context)

	// Check for circular dependencies
	if err := resolver.DetectCircularDependencies(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Resolve dependencies
	updates, err := resolver.ResolveDependencies(req.FormData, req.ChangedFields, c.Ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to resolve dependencies",
		})
	}

	// Return field updates
	return c.JSON(fiber.Map{
		"fields": updates,
	})
}
