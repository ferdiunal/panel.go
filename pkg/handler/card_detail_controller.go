package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

/// # HandleCardDetail
///
/// Bu fonksiyon, belirli bir index değerine sahip kart (card) detaylarını getirir.
/// Panel sisteminde dashboard veya özet görünümlerinde kullanılan kartların
/// detaylı bilgilerini almak için kullanılır.
///
/// ## Kullanım Senaryoları
///
/// - Dashboard'da bir kartın detaylı verilerini görüntüleme
/// - Kart metriklerinin gerçek zamanlı güncellenmesi
/// - Kart içeriğinin dinamik olarak yüklenmesi
/// - API endpoint'i üzerinden kart verilerine erişim
///
/// ## Parametreler
///
/// * `h` - `*FieldHandler`: Field handler instance'ı. Kart koleksiyonunu ve
///   veritabanı bağlantısını içerir.
///   - `h.Cards`: Sistemde tanımlı tüm kartların listesi
///   - `h.DB`: Veritabanı bağlantısı (GORM instance)
///
/// * `c` - `*context.Context`: HTTP request context'i. Fiber context'inin
///   genişletilmiş versiyonu.
///   - URL parametrelerini okumak için kullanılır
///   - HTTP response'ları göndermek için kullanılır
///
/// ## Dönüş Değeri
///
/// * `error`: İşlem başarılı ise `nil`, aksi halde hata döner
///   - 404: Kart bulunamadı (geçersiz index)
///   - 500: Kart verisi çözümlenirken hata oluştu
///   - 200: Başarılı, kart verisi JSON formatında döner
///
/// ## HTTP Response Formatı
///
/// ### Başarılı Response (200 OK)
/// ```json
/// {
///   "data": {
///     // Kartın çözümlenmiş verileri
///     // İçerik kart tipine göre değişir
///   }
/// }
/// ```
///
/// ### Hata Response (404 Not Found)
/// ```json
/// {
///   "error": "Card not found"
/// }
/// ```
///
/// ### Hata Response (500 Internal Server Error)
/// ```json
/// {
///   "error": "Hata mesajı detayı"
/// }
/// ```
///
/// ## Kullanım Örneği
///
/// ```go
/// // Router tanımlaması
/// app.Get("/api/cards/:index", func(c *fiber.Ctx) error {
///     ctx := context.New(c)
///     return HandleCardDetail(fieldHandler, ctx)
/// })
///
/// // HTTP Request
/// // GET /api/cards/0
/// // Response: İlk kartın detaylı verileri
///
/// // GET /api/cards/5
/// // Response: Altıncı kartın detaylı verileri
/// ```
///
/// ## İşleyiş Akışı
///
/// 1. **Index Parametresi Alma**: URL'den "index" parametresi integer olarak alınır
/// 2. **Validasyon**: Index değeri kontrol edilir
///    - Sayısal bir değer olmalı
///    - Negatif olmamalı
///    - Cards array uzunluğundan küçük olmalı
/// 3. **Kart Seçimi**: Geçerli index ile kart array'inden ilgili kart alınır
/// 4. **Veri Çözümleme**: Kartın `Resolve` metodu çağrılarak veriler hazırlanır
/// 5. **Response Dönme**: Çözümlenen veriler JSON formatında döndürülür
///
/// ## Önemli Notlar
///
/// ⚠️ **Index Sıfır Tabanlı**: Index değeri 0'dan başlar. İlk kart için index=0
///
/// ⚠️ **Thread Safety**: Bu fonksiyon concurrent request'lerde güvenli çalışır,
/// ancak `h.Cards` slice'ının runtime'da değiştirilmemesi gerekir.
///
/// ⚠️ **Performans**: Her request'te `Resolve` metodu çağrılır, bu nedenle
/// ağır hesaplamalar içeren kartlar için caching düşünülmelidir.
///
/// ## Güvenlik Notları
///
/// - Index parametresi otomatik olarak validate edilir
/// - Array bounds kontrolü yapılır, buffer overflow riski yoktur
/// - SQL injection koruması `Resolve` metodunun implementasyonuna bağlıdır
///
/// ## Hata Yönetimi
///
/// Fonksiyon üç farklı hata senaryosunu ele alır:
///
/// 1. **Parametre Hatası**: Index parametresi integer'a çevrilemezse 404 döner
/// 2. **Geçersiz Index**: Index negatif veya sınırların dışındaysa 404 döner
/// 3. **Resolve Hatası**: Veri çözümleme sırasında hata oluşursa 500 döner
///
/// ## Avantajlar
///
/// ✅ Basit ve anlaşılır API
/// ✅ Otomatik validasyon
/// ✅ Tutarlı hata yönetimi
/// ✅ Index tabanlı hızlı erişim
///
/// ## Dezavantajlar
///
/// ❌ Index değişebilir (kartlar eklenip çıkarılırsa)
/// ❌ Kart ID'si yerine index kullanımı daha az güvenli
/// ❌ Caching mekanizması yok
///
/// ## İlgili Fonksiyonlar
///
/// - `HandleCards`: Tüm kartları listeler
/// - `Card.Resolve`: Kart verilerini çözümler
///
/// ## Versiyon Notları
///
/// Bu fonksiyon panel.go'nun kart sistemi için temel endpoint'lerden biridir.
/// Gelecek versiyonlarda ID tabanlı erişim de eklenebilir.
func HandleCardDetail(h *FieldHandler, c *context.Context) error {
	index, err := c.ParamsInt("index")
	if err != nil || index < 0 || index >= len(h.Cards) {
		return c.Status(404).JSON(fiber.Map{"error": "Card not found"})
	}

	w := h.Cards[index]

	// Resolve data
	data, err := w.Resolve(c, h.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
