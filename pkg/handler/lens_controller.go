// # Lens Controller Paketi
//
// Bu paket, kaynak (resource) bazlı lens işlemlerini yöneten HTTP handler fonksiyonlarını içerir.
// Lens'ler, Laravel Nova'dan esinlenerek geliştirilmiş, kaynaklar üzerinde özel filtreleme
// ve görünüm sağlayan güçlü bir özelliktir.
//
// ## Lens Nedir?
//
// Lens'ler, bir kaynağın verilerini özel bir perspektiften görüntülemenizi sağlayan
// özelleştirilmiş görünümlerdir. Örneğin:
// - "En Popüler Ürünler" lens'i
// - "Son 30 Günde Eklenen Kullanıcılar" lens'i
// - "Yüksek Öncelikli Görevler" lens'i
//
// ## Kullanım Senaryoları
//
// 1. **Özel Raporlama**: Belirli kriterlere göre filtrelenmiş veri görünümleri
// 2. **Dashboard Widgets**: Özel metrikler ve istatistikler
// 3. **Analitik Görünümler**: Karmaşık sorgulamalar ve agregasyonlar
// 4. **İş Mantığı Filtreleri**: Domain-specific veri görünümleri
//
// ## Mimari Yapı
//
// ```
// Client Request → Router → HandleLensIndex/HandleLens → Resource.Lenses()
//                                                       → Lens.Query()
//                                                       → Filtered Data
// ```
package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// # HandleLensIndex
//
// Bu fonksiyon, bir kaynak için mevcut tüm lens'leri listeler ve istemciye döndürür.
// Laravel Nova'nın `LensController@index` metoduna karşılık gelir.
//
// ## Amaç
//
// Bir kaynak için tanımlanmış tüm lens'lerin meta bilgilerini (isim ve slug) döndürerek,
// istemci tarafında lens seçimi yapılabilmesini sağlar. Bu, dinamik UI oluşturma için
// kritik öneme sahiptir.
//
// ## Parametreler
//
// - `h *FieldHandler`: Kaynak ve field işlemlerini yöneten handler. Resource bilgisini içerir.
//   - `h.Resource`: İşlem yapılacak kaynak nesnesi (nil kontrolü yapılır)
//   - `h.Resource.Lenses()`: Kaynağa ait lens listesini döndürür
//
// - `c *context.Context`: HTTP istek/yanıt context'i (Fiber context wrapper)
//   - İstek parametrelerine erişim
//   - Yanıt oluşturma ve durum kodu ayarlama
//
// ## Dönüş Değeri
//
// - `error`: İşlem başarılı ise nil, hata durumunda error nesnesi
//   - Resource bulunamadığında: 404 Not Found
//   - Başarılı durumda: 200 OK ile lens listesi
//
// ## Yanıt Formatı
//
// ```json
// {
//   "data": [
//     {
//       "name": "Most Popular Products",
//       "slug": "most-popular-products"
//     },
//     {
//       "name": "Recent Users",
//       "slug": "recent-users"
//     }
//   ]
// }
// ```
//
// ## Kullanım Örneği
//
// ```go
// // Router tanımlaması
// app.Get("/api/:resource/lenses", func(c *fiber.Ctx) error {
//     handler := NewFieldHandler(resource)
//     ctx := context.New(c)
//     return HandleLensIndex(handler, ctx)
// })
//
// // İstemci tarafı kullanım
// // GET /api/products/lenses
// // Response: {"data": [{"name": "Popular", "slug": "popular"}]}
// ```
//
// ## Hata Durumları
//
// 1. **Resource Bulunamadı (404)**
//    - Durum: `h.Resource == nil`
//    - Yanıt: `{"error": "Resource not found"}`
//    - Sebep: Geçersiz resource adı veya kayıt edilmemiş resource
//
// ## Önemli Notlar
//
// - ⚠️ **Nil Kontrolü**: Resource nil kontrolü mutlaka yapılmalıdır
// - 📝 **Lens Kayıt**: Lens'ler resource tanımında kayıt edilmelidir
// - 🔒 **Yetkilendirme**: Bu endpoint'e erişim kontrolü üst katmanda yapılmalıdır
// - 🚀 **Performans**: Lens listesi genellikle küçüktür, cache'leme gerekmez
//
// ## Avantajlar
//
// - ✅ Dinamik lens keşfi sağlar
// - ✅ Frontend'de otomatik UI oluşturma imkanı
// - ✅ Basit ve anlaşılır API
// - ✅ Laravel Nova ile uyumlu yapı
//
// ## Dikkat Edilmesi Gerekenler
//
// - Resource'un nil olma durumu kontrol edilmelidir
// - Lens'lerin Name() ve Slug() metodları implement edilmiş olmalıdır
// - Yanıt formatı frontend ile uyumlu olmalıdır
func HandleLensIndex(h *FieldHandler, c *context.Context) error {
	if h.Resource == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}

	lenses := h.Resource.Lenses()
	response := make([]map[string]interface{}, 0)

	for _, lens := range lenses {
		response = append(response, map[string]interface{}{
			"name": lens.Name(),
			"slug": lens.Slug(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

// # HandleLens
//
// Bu fonksiyon, belirli bir lens üzerinden filtrelenmiş kaynak verilerini listeler.
// Laravel Nova'nın `LensController@show` metoduna karşılık gelir.
//
// ## Amaç
//
// Seçilen lens'in tanımladığı özel sorgu ve filtreleme mantığını uygulayarak,
// kaynak verilerini özelleştirilmiş bir görünümde sunar. Lens'in query metodunu
// kullanarak veri setini filtreler ve standart index formatında döndürür.
//
// ## Çalışma Prensibi
//
// Bu fonksiyon, `NewLensHandler` tarafından önceden yapılandırılmış bir handler
// ile çalışır. Lens'in query mantığı handler oluşturulurken uygulanmıştır, bu
// nedenle doğrudan `HandleResourceIndex` fonksiyonunu kullanarak filtrelenmiş
// veri setini döndürür.
//
// ## İşlem Akışı
//
// ```
// 1. NewLensHandler() → Lens query'si uygulanır
// 2. HandleLens() çağrılır
// 3. HandleResourceIndex() → Filtrelenmiş veri döndürülür
// ```
//
// ## Parametreler
//
// - `h *FieldHandler`: Lens query'si ile önceden yapılandırılmış handler
//   - Lens'in Query() metodu zaten uygulanmış durumda
//   - Filtrelenmiş veri seti üzerinde çalışır
//   - Pagination, sorting gibi standart işlemler desteklenir
//
// - `c *context.Context`: HTTP istek/yanıt context'i
//   - Query parametreleri (page, per_page, sort, etc.)
//   - Filter parametreleri
//   - Yanıt oluşturma
//
// ## Dönüş Değeri
//
// - `error`: İşlem başarılı ise nil, hata durumunda error nesnesi
//   - HandleResourceIndex'in döndürdüğü tüm hatalar
//   - Lens query hatası (varsa)
//
// ## Yanıt Formatı
//
// ```json
// {
//   "data": [
//     {
//       "id": 1,
//       "name": "Product A",
//       "popularity_score": 95
//     }
//   ],
//   "meta": {
//     "current_page": 1,
//     "per_page": 15,
//     "total": 42
//   }
// }
// ```
//
// ## Kullanım Örneği
//
// ```go
// // Router tanımlaması
// app.Get("/api/:resource/lens/:lens", func(c *fiber.Ctx) error {
//     resourceName := c.Params("resource")
//     lensSlug := c.Params("lens")
//
//     // Lens handler oluştur (query otomatik uygulanır)
//     handler := NewLensHandler(resourceName, lensSlug)
//     ctx := context.New(c)
//
//     return HandleLens(handler, ctx)
// })
//
// // İstemci tarafı kullanım
// // GET /api/products/lens/most-popular?page=1&per_page=20
// // Response: Filtrelenmiş ve sayfalanmış ürün listesi
// ```
//
// ## Lens Query Örneği
//
// ```go
// type MostPopularLens struct {
//     base.Lens
// }
//
// func (l *MostPopularLens) Query(query interface{}) interface{} {
//     db := query.(*gorm.DB)
//     return db.Where("popularity_score > ?", 80).
//            Order("popularity_score DESC")
// }
// ```
//
// ## Önemli Notlar
//
// - 🔄 **Önceden Yapılandırma**: Handler, NewLensHandler ile oluşturulmalıdır
// - 🎯 **Lens Query**: Lens'in Query() metodu handler oluşturulurken uygulanır
// - 📊 **Standart Format**: Yanıt formatı normal index endpoint ile aynıdır
// - 🔍 **Ek Filtreler**: Lens query'sine ek olarak standart filtreler de uygulanabilir
// - ⚡ **Performans**: Lens query'leri optimize edilmiş olmalıdır (index kullanımı)
//
// ## Avantajlar
//
// - ✅ Karmaşık sorguları basit API'ye dönüştürür
// - ✅ Yeniden kullanılabilir veri görünümleri
// - ✅ Standart pagination ve sorting desteği
// - ✅ Mevcut index mantığını yeniden kullanır
// - ✅ Temiz ve maintainable kod yapısı
//
// ## Dezavantajlar
//
// - ⚠️ Lens query'si yanlış yazılırsa performans sorunları olabilir
// - ⚠️ Karmaşık agregasyonlar için ek optimizasyon gerekebilir
//
// ## Dikkat Edilmesi Gerekenler
//
// - Lens query'leri database index'leri ile uyumlu olmalıdır
// - N+1 query problemine dikkat edilmelidir
// - Büyük veri setlerinde pagination mutlaka kullanılmalıdır
// - Lens query'leri test edilmelidir
// - Cache stratejisi değerlendirilmelidir (sık kullanılan lens'ler için)
//
// ## Güvenlik Notları
//
// - 🔒 Lens erişim yetkileri kontrol edilmelidir
// - 🔒 SQL injection'a karşı parameterized query kullanılmalıdır
// - 🔒 Hassas veriler lens query'sinde filtrelenmelidir
func HandleLens(h *FieldHandler, c *context.Context) error {
	// Lens handler is already configured with filtered query via NewLensHandler
	// We can directly use the Index logic but with the lens's filtered dataset
	return HandleResourceIndex(h, c)
}
