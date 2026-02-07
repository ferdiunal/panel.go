package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleFieldResolve, alan bazlı resolver fonksiyonlarını çağırmak için kullanılan HTTP handler fonksiyonudur.
//
// # Genel Bakış
//
// Bu fonksiyon, dinamik alan dönüşümlerini tetiklemek için frontend bileşenlerinin kullanabileceği
// bir API endpoint'i sağlar. Belirli bir kaynağın (resource) belirli bir alanı (field) için
// resolver fonksiyonunu çağırır ve sonucu döndürür.
//
// # Kullanım Senaryoları
//
// 1. **Dinamik Alan Dönüşümleri**: Bir alanın değerini runtime'da özel parametrelerle dönüştürme
// 2. **Hesaplanmış Alanlar**: Veritabanından gelen ham veriyi işleyerek hesaplanmış değerler üretme
// 3. **Koşullu Veri Gösterimi**: Kullanıcı rolüne veya diğer parametrelere göre farklı veri gösterme
// 4. **Veri Zenginleştirme**: Temel alan verisini ek bilgilerle zenginleştirme
// 5. **Format Dönüşümleri**: Tarihleri, sayıları veya metinleri farklı formatlara dönüştürme
//
// # HTTP Endpoint Detayları
//
// - **Route**: `/resources/:resource/:id/fields/:field/resolve`
// - **Method**: `POST`
// - **URL Parametreleri**:
//   - `:resource` - Kaynak adı (örn: "users", "products")
//   - `:id` - Kaynak ID'si (örn: "123", "abc-def")
//   - `:field` - Alan adı (örn: "email", "price")
//
// # Request Body
//
// Request body, resolver'a özgü parametreleri JSON formatında içermelidir:
//
// ```json
// {
//   "format": "currency",
//   "locale": "tr-TR",
//   "precision": 2
// }
// ```
//
// Body boş olabilir, bu durumda boş bir parametre map'i kullanılır.
//
// # Response Format
//
// Başarılı durumda:
// ```json
// {
//   "data": {
//     "key": "field_name",
//     "value": "resolved_value",
//     "label": "Field Label",
//     ...
//   }
// }
// ```
//
// Hata durumlarında:
// ```json
// {
//   "error": "Field not found"
// }
// ```
// veya
// ```json
// {
//   "error": "Resource not found"
// }
// ```
//
// # Parametreler
//
// - `h *FieldHandler`: Alan handler'ı, resolver'ların tanımlandığı ve yönetildiği yapı.
//   Bu yapı üzerinden alan listesine (Elements) ve veri sağlayıcısına (Provider) erişilir.
//
// - `c *context.Context`: Fiber context wrapper'ı. HTTP request/response işlemleri için kullanılır.
//   URL parametrelerine, request body'sine ve response yazma işlemlerine erişim sağlar.
//
// # Dönüş Değeri
//
// - `error`: İşlem başarılı ise nil, hata durumunda error döner.
//   Fiber framework'ü bu error'ı otomatik olarak HTTP response'a dönüştürür.
//
// # İşlem Akışı
//
// 1. URL'den field adı ve resource ID'si alınır
// 2. Field adına göre Elements listesinde arama yapılır
// 3. Field bulunamazsa 404 hatası döner
// 4. Provider üzerinden resource item'ı getirilir
// 5. Item bulunamazsa 404 hatası döner
// 6. Request body'den parametreler parse edilir (opsiyonel)
// 7. Field'ın Extract metodu çağrılarak item'dan değer çıkarılır
// 8. Field'ın JsonSerialize metodu çağrılarak veri serileştirilir
// 9. Serileştirilmiş veri JSON response olarak döndürülür
//
// # Kullanım Örnekleri
//
// ## Örnek 1: Basit Alan Çözümleme
//
// ```go
// // Route tanımı
// app.Post("/resources/:resource/:id/fields/:field/resolve",
//     func(c *fiber.Ctx) error {
//         ctx := context.New(c)
//         return HandleFieldResolve(fieldHandler, ctx)
//     })
//
// // Frontend'den çağrı
// fetch('/resources/users/123/fields/email/resolve', {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({})
// })
// ```
//
// ## Örnek 2: Parametreli Resolver Çağrısı
//
// ```go
// // Fiyat alanını farklı para birimlerinde gösterme
// fetch('/resources/products/456/fields/price/resolve', {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({
//         currency: 'USD',
//         includeVAT: true
//     })
// })
// ```
//
// # Avantajlar
//
// - **Esneklik**: Runtime'da dinamik veri dönüşümleri yapabilme
// - **Yeniden Kullanılabilirlik**: Aynı resolver'ı farklı parametrelerle kullanabilme
// - **Separation of Concerns**: İş mantığını API katmanından ayırma
// - **Frontend Kontrolü**: Frontend'in veri formatını kontrol edebilmesi
// - **Performans**: Sadece gerekli alanlar için resolver çağrılabilir
//
// # Dezavantajlar ve Dikkat Edilmesi Gerekenler
//
// - **Performans**: Her resolver çağrısı ayrı bir HTTP request gerektirir
// - **Güvenlik**: Resolver parametreleri dikkatli validate edilmelidir
// - **Karmaşıklık**: Çok fazla resolver kullanımı kod karmaşıklığını artırabilir
// - **Hata Yönetimi**: Resolver içindeki hatalar düzgün handle edilmelidir
//
// # Önemli Notlar
//
// ⚠️ **Güvenlik Uyarısı**: Bu endpoint, kullanıcının erişim yetkisi olan kaynaklara
// sınırlandırılmalıdır. Middleware'ler ile yetkilendirme kontrolü yapılmalıdır.
//
// ⚠️ **Performans Uyarısı**: Ağır işlemler yapan resolver'lar için caching mekanizması
// düşünülmelidir. Her çağrıda veritabanı sorgusu yapmak performans sorunlarına yol açabilir.
//
// 💡 **İpucu**: Resolver parametreleri için bir şema tanımlayarak, geçersiz parametrelerin
// erken aşamada yakalanmasını sağlayabilirsiniz.
//
// 💡 **İpucu**: Sık kullanılan resolver sonuçları için Redis gibi bir cache katmanı
// kullanmak performansı önemli ölçüde artırabilir.
//
// # Gereksinimler
//
// - **Requirement 16.1**: Sistem, alan resolver'larını API endpoint'leri aracılığıyla
//   erişilebilir hale getirmelidir
// - **Requirement 16.2**: Sistem, resolver'ların özel veri dönüşümleri gerçekleştirmesine
//   izin vermelidir
// - **Requirement 16.3**: Bir resolver çağrıldığında, sistem resolver-spesifik parametreleri
//   desteklemelidir
//
// # İlgili Tipler
//
// - `FieldHandler`: Alan yönetimi için kullanılan handler yapısı
// - `context.Context`: HTTP context wrapper'ı
// - `fiber.Map`: JSON response için kullanılan map tipi
//
// # Ayrıca Bakınız
//
// - `FieldHandler.Elements`: Tüm alan tanımlarının listesi
// - `FieldHandler.Provider`: Veri sağlayıcı interface'i
// - `Provider.Show()`: Tek bir kaynağı getiren metod
func HandleFieldResolve(h *FieldHandler, c *context.Context) error {
	// URL parametrelerinden alan adını al
	// Örnek: /resources/users/123/fields/email/resolve -> "email"
	fieldName := c.Params("field")

	// URL parametrelerinden kaynak ID'sini al
	// Örnek: /resources/users/123/fields/email/resolve -> "123"
	resourceID := c.Params("id")

	// ============================================================================
	// Adım 1: Alan Bulma (Field Lookup)
	// ============================================================================
	// Handler'da tanımlı tüm alanlar arasında istenen alanı bul.
	// Elements listesi, bu kaynak için tanımlanmış tüm alanları içerir.
	// Her alan bir GetKey() metoduna sahiptir ve bu metod alanın benzersiz adını döner.
	//
	// Not: Bu işlem O(n) karmaşıklığındadır. Çok sayıda alan varsa,
	// performans için bir map yapısı kullanılabilir.
	var targetField interface{}
	for _, element := range h.Elements {
		if element.GetKey() == fieldName {
			targetField = element
			break
		}
	}

	// Alan bulunamadıysa 404 hatası dön
	// Bu durum, geçersiz bir alan adı istendiğinde oluşur
	if targetField == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Field not found",
		})
	}

	// ============================================================================
	// Adım 2: Kaynak Verisi Getirme (Resource Retrieval)
	// ============================================================================
	// Provider'ın Show metodunu kullanarak belirtilen ID'ye sahip kaynağı getir.
	// Provider, veritabanı veya başka bir veri kaynağından veri çeker.
	//
	// Show metodu şunları yapabilir:
	// - Veritabanından tek bir kayıt çekme
	// - İlişkili verileri eager loading ile yükleme
	// - Yetkilendirme kontrolü yapma
	// - Cache'den veri okuma
	item, err := h.Provider.Show(c, resourceID)
	if err != nil {
		// Kaynak bulunamadıysa veya erişim yetkisi yoksa 404 hatası dön
		// Güvenlik nedeniyle "bulunamadı" ve "yetkisiz" durumları aynı hata ile döndürülür
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}

	// ============================================================================
	// Adım 3: Resolver Parametrelerini Parse Etme (Parameter Parsing)
	// ============================================================================
	// Request body'den resolver'a özgü parametreleri al.
	// Bu parametreler, resolver'ın davranışını özelleştirmek için kullanılır.
	//
	// Örnek parametreler:
	// - format: "currency", "date", "percentage"
	// - locale: "tr-TR", "en-US"
	// - precision: 2, 4
	// - includeVAT: true, false
	params := make(map[string]interface{})
	if err := c.Ctx.BodyParser(&params); err != nil {
		// Body parse edilemezse veya boşsa, boş bir parametre map'i kullan
		// Bu, parametresiz resolver çağrılarına izin verir
		params = make(map[string]interface{})
	}

	// ============================================================================
	// Adım 4: Alan Değerini Çıkarma (Field Value Extraction)
	// ============================================================================
	// Gelecekteki geliştirme: Burada resolver parametreleri kullanılarak
	// özel dönüşümler yapılabilir. Şu anki implementasyon temel değer
	// çıkarma işlemini gerçekleştirir.
	//
	// Örnek gelişmiş kullanım:
	// - Tarih alanları için timezone dönüşümü
	// - Para birimi alanları için kur çevrimi
	// - Metin alanları için dil çevirisi
	// - Resim alanları için boyut/format dönüşümü

	// Type assertion ile Extract metodunun varlığını kontrol et
	// Extract metodu, item'dan ilgili alan değerini çıkarır ve field'ın
	// internal state'ine kaydeder
	if field, ok := targetField.(interface{ Extract(interface{}) }); ok {
		field.Extract(item)
	}

	// ============================================================================
	// Adım 5: Alan Serileştirme (Field Serialization)
	// ============================================================================
	// Field'ı JSON formatına dönüştür.
	// JsonSerialize metodu, field'ın tüm özelliklerini (value, label, metadata vb.)
	// bir map olarak döner.
	//
	// Dönen map tipik olarak şunları içerir:
	// - key: Alan adı
	// - value: Alan değeri
	// - label: Görüntüleme etiketi
	// - type: Alan tipi
	// - metadata: Ek bilgiler
	var serialized map[string]interface{}
	if field, ok := targetField.(interface{ JsonSerialize() map[string]interface{} }); ok {
		serialized = field.JsonSerialize()
	}

	// ============================================================================
	// Adım 6: Response Dönme (Response Return)
	// ============================================================================
	// Serileştirilmiş alan verisini JSON response olarak dön.
	// Response formatı frontend tarafından kolayca işlenebilir.
	//
	// Başarılı response örneği:
	// {
	//   "data": {
	//     "key": "price",
	//     "value": 1234.56,
	//     "label": "Fiyat",
	//     "type": "number",
	//     "formatted": "1.234,56 TL"
	//   }
	// }
	return c.JSON(fiber.Map{
		"data": serialized,
	})
}
