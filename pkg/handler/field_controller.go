// Package handler, panel.go framework'ünde HTTP isteklerini işleyen ve kaynak (resource)
// yönetimi için gerekli controller fonksiyonlarını sağlayan bir pakettir.
//
// Bu paket, özellikle field (alan) işlemleri için gerekli HTTP handler'ları içerir ve
// Fiber web framework'ü ile entegre çalışır.
//
// # Temel Özellikler
//
// - Field listesi oluşturma ve serileştirme
// - Resource context yönetimi
// - Otomatik değer çözümleme (value resolution)
// - JSON serileştirme desteği
//
// # Kullanım Senaryoları
//
// - Admin panel'de form alanlarının listelenmesi
// - Kaynak detay sayfalarında field bilgilerinin gösterilmesi
// - API endpoint'lerinde field metadata'sının sunulması
//
// # Örnek Kullanım
//
// ```go
// app := fiber.New()
// fieldHandler := &FieldHandler{}
//
//	app.Get("/api/resources/:resource/fields", func(c *fiber.Ctx) error {
//	    ctx := context.New(c)
//	    return HandleFieldList(fieldHandler, ctx)
//	})
//
// ```
package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleFieldList, bir kaynak (resource) için tanımlanmış tüm alanları (fields) listeleyen
// ve JSON formatında döndüren HTTP handler fonksiyonudur.
//
// Bu fonksiyon, resource context'inden field bilgilerini alır, görünürlük filtrelerini uygular
// ve eğer bir resource instance mevcutsa, her field için değer çözümlemesi yapar.
//
// # Parametreler
//
// - `h`: FieldHandler pointer'ı - Field işlemleri için gerekli handler instance'ı
// - `c`: Context pointer'ı - Panel.go custom context'i, Fiber context'ini wrap eder
//
// # Dönüş Değeri
//
// - `error`: İşlem başarılıysa nil, hata durumunda error döner
//
// # İşleyiş Akışı
//
// 1. Context'ten resource bilgisini alır
// 2. Resource bulunamazsa 500 Internal Server Error döner
// 3. Her field için:
//   - Eğer resource instance mevcutsa, field değerini extract eder
//   - Field'ı JSON formatına serialize eder
//
// 4. Tüm field'ları içeren JSON array döner
//
// # Kullanım Senaryoları
//
// ## Senaryo 1: Form Alanlarının Listelenmesi
// Admin panel'de yeni kayıt oluşturma veya düzenleme formunda hangi alanların
// gösterileceğini belirlemek için kullanılır.
//
// ## Senaryo 2: Detay Sayfası Metadata'sı
// Bir kaydın detay sayfasında hangi alanların nasıl gösterileceğine dair
// metadata bilgisi sağlar.
//
// ## Senaryo 3: API Documentation
// Frontend uygulamaların hangi field'ların mevcut olduğunu ve özelliklerini
// öğrenmesi için kullanılır.
//
// # Örnek Kullanım
//
// ```go
// // Route tanımlama
//
//	app.Get("/api/resources/:resource/fields", func(c *fiber.Ctx) error {
//	    ctx := context.New(c)
//	    fieldHandler := &FieldHandler{}
//	    return HandleFieldList(fieldHandler, ctx)
//	})
//
// // Örnek Response:
// // [
// //   {
// //     "name": "title",
// //     "type": "text",
// //     "label": "Başlık",
// //     "value": "Örnek Başlık",
// //     "required": true
// //   },
// //   {
// //     "name": "description",
// //     "type": "textarea",
// //     "label": "Açıklama",
// //     "value": "Örnek açıklama metni",
// //     "required": false
// //   }
// // ]
// ```
//
// # Hata Durumları
//
// ## Resource Context Bulunamadı
// Eğer context'ten resource bilgisi alınamazsa, 500 status code ile
// "Field context not found" hatası döner.
//
// ```json
//
//	{
//	  "error": "Field context not found"
//	}
//
// ```
//
// # Önemli Notlar
//
//   - **Thread Safety**: Bu fonksiyon her istek için yeni bir context ile çağrılır,
//     dolayısıyla thread-safe'dir.
//
//   - **Performance**: Field sayısı arttıkça response süresi artabilir. Büyük field
//     listelerinde pagination düşünülmelidir.
//
//   - **Value Extraction**: Resource instance mevcutsa, her field için Extract() metodu
//     çağrılır. Bu işlem veritabanı sorguları içerebilir (özellikle relationship field'lar için).
//
//   - **Visibility Filtering**: Mevcut implementasyonda tüm field'lar döndürülür.
//     İleride role-based veya context-based filtering eklenebilir.
//
// # Avantajlar
//
// - **Dinamik Form Oluşturma**: Frontend, field metadata'sına göre dinamik formlar oluşturabilir
// - **Type Safety**: Her field kendi JsonSerialize() metodunu implement eder
// - **Extensibility**: Yeni field tipleri kolayca eklenebilir
// - **Separation of Concerns**: Field logic'i handler'dan ayrıdır
//
// # Dezavantajlar
//
// - **N+1 Problem**: Relationship field'lar için Extract() çağrıları N+1 sorgu problemine yol açabilir
// - **Memory Usage**: Büyük field listelerinde tüm data memory'de tutulur
// - **No Pagination**: Şu an pagination desteği yok
//
// # İyileştirme Önerileri
//
// 1. **Lazy Loading**: Field değerlerini sadece gerektiğinde yükle
// 2. **Caching**: Field metadata'sını cache'le
// 3. **Pagination**: Büyük field listeleri için sayfalama ekle
// 4. **Filtering**: Query parameter'larına göre field filtreleme
// 5. **Eager Loading**: Relationship field'lar için eager loading kullan
//
// # İlgili Tipler
//
// - `FieldHandler`: Field işlemleri için handler struct'ı
// - `context.Context`: Panel.go custom context
// - `context.ResourceContext`: Resource ve field bilgilerini içeren context
//
// # Bağımlılıklar
//
// - `github.com/gofiber/fiber/v2`: HTTP framework
// - `github.com/ferdiunal/panel.go/pkg/context`: Custom context paketi
func HandleFieldList(h *FieldHandler, c *context.Context) error {
	// Resource context'ini al
	// Bu context, resource tanımı ve field listesini içerir
	ctx := c.Resource()
	if ctx == nil {
		// Resource context bulunamadıysa, bu ciddi bir hata durumudur
		// Muhtemelen middleware'de resource yüklenmemiştir
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Field context not found",
		})
	}

	// Response için boş bir slice oluştur
	// Capacity belirtilmemiş, çünkü field sayısı önceden bilinmiyor
	// Performans için ctx.Elements uzunluğu kullanılabilir: make([]map[string]interface{}, 0, len(ctx.Elements))
	response := make([]map[string]interface{}, 0)

	// Her field için işlem yap
	for _, element := range ctx.Elements {
		// Context bazlı görünürlük kontrolü - gizli field'ları backend'den gönderme
		if !element.IsVisible(ctx) {
			continue
		}

		// Eğer bir resource instance mevcutsa (örneğin edit/detail sayfası),
		// field'ın değerini resource'tan extract et
		// Bu işlem, field'ın Extract() metodunu çağırır ve değeri field struct'ına set eder
		if ctx.Resource != nil {
			element.Extract(ctx.Resource)
		}

		// Field'ı JSON formatına serialize et ve response'a ekle
		// JsonSerialize() metodu, field'ın tüm özelliklerini (name, type, label, value, vb.)
		// map[string]interface{} formatında döndürür
		response = append(response, element.JsonSerialize())
	}

	// Tüm field'ları içeren JSON array'i döndür
	return c.JSON(response)
}
