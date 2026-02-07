package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

/// # HandleResourceCreate
///
/// Bu fonksiyon, kaynak oluşturma işlemi için gerekli form alanlarını döndürür.
/// GET isteklerini `/api/resource/:resource/create` endpoint'ine yönlendirir ve
/// oluşturma formunda gösterilecek alanları hazırlar.
///
/// ## Temel İşlevsellik
///
/// 1. **Yetkilendirme Kontrolü**: Policy varsa Create yetkisini kontrol eder
/// 2. **Alan Filtreleme**: Sadece oluşturma bağlamında görünür alanları seçer
/// 3. **Seçenek Çözümleme**: AutoOptions ve callback'leri çalıştırır
/// 4. **JSON Serileştirme**: Alanları frontend için uygun formata dönüştürür
///
/// ## Parametreler
///
/// * `h` - `*FieldHandler`: Kaynak için tanımlanmış alan işleyicisi
///   - `Policy`: Yetkilendirme politikası (opsiyonel)
///   - `Elements`: Kaynak için tanımlanmış tüm alanlar
///   - `ResolveFieldOptions`: Alan seçeneklerini çözümleyen metod
///
/// * `c` - `*context.Context`: Panel bağlamı
///   - HTTP istek/yanıt bilgilerini içerir
///   - Kaynak bilgilerine erişim sağlar
///   - Kullanıcı oturum bilgilerini taşır
///
/// ## Dönüş Değeri
///
/// * `error`: İşlem başarılıysa nil, aksi halde hata döner
///   - 403 Forbidden: Kullanıcının oluşturma yetkisi yoksa
///   - 200 OK: Başarılı durumda alan listesi ile JSON yanıtı
///
/// ## Alan Bağlam Filtreleme
///
/// Fonksiyon, aşağıdaki bağlamlara sahip alanları **hariç tutar**:
/// - `HIDE_ON_CREATE`: Oluşturma formunda gizli alanlar
/// - `ONLY_ON_LIST`: Sadece liste görünümünde gösterilen alanlar
/// - `ONLY_ON_DETAIL`: Sadece detay görünümünde gösterilen alanlar
/// - `ONLY_ON_UPDATE`: Sadece güncelleme formunda gösterilen alanlar
///
/// ## Kullanım Senaryoları
///
/// ### 1. Standart Kaynak Oluşturma
/// ```go
/// // Kullanıcı yeni bir kayıt oluşturmak için formu açar
/// // GET /api/resource/users/create
/// // Yanıt: { "fields": [...] }
/// ```
///
/// ### 2. İlişkisel Alan Seçenekleri
/// ```go
/// // BelongsTo veya HasMany alanları için seçenekler yüklenir
/// // AutoOptions callback'leri çalıştırılır
/// // Dropdown/select alanları için veri hazırlanır
/// ```
///
/// ### 3. Koşullu Alan Görünürlüğü
/// ```go
/// // IsVisible kontrolü ile dinamik alan gösterimi
/// // Kullanıcı rolüne veya kaynak durumuna göre alanlar filtrelenir
/// ```
///
/// ## Yanıt Formatı
///
/// ```json
/// {
///   "fields": [
///     {
///       "name": "title",
///       "type": "text",
///       "label": "Başlık",
///       "required": true,
///       "placeholder": "Başlık giriniz",
///       "rules": ["required", "min:3"]
///     },
///     {
///       "name": "category_id",
///       "type": "select",
///       "label": "Kategori",
///       "options": [
///         { "value": 1, "label": "Teknoloji" },
///         { "value": 2, "label": "Spor" }
///       ]
///     }
///   ]
/// }
/// ```
///
/// ## Güvenlik Özellikleri
///
/// * **Policy Kontrolü**: Create yetkisi olmayan kullanıcılar 403 hatası alır
/// * **Alan Görünürlüğü**: IsVisible ile hassas alanlar gizlenebilir
/// * **Bağlam Filtreleme**: Sadece oluşturma için uygun alanlar döndürülür
///
/// ## Performans Notları
///
/// * Alan sayısı arttıkça işlem süresi artar
/// * AutoOptions callback'leri veritabanı sorguları yapabilir
/// * Büyük seçenek listeleri için pagination düşünülmelidir
///
/// ## Önemli Uyarılar
///
/// ⚠️ **Policy Kontrolü**: Policy nil ise yetkilendirme atlanır, üretim ortamında
/// mutlaka Policy tanımlanmalıdır.
///
/// ⚠️ **AutoOptions**: Callback'ler senkron çalışır, yavaş sorgular response süresini
/// etkiler. Cache kullanımı önerilir.
///
/// ⚠️ **Görünürlük**: IsVisible false dönen alanlar frontend'e gönderilmez, bu
/// hassas bilgilerin korunması için kullanılabilir.
///
/// ## İlgili Fonksiyonlar
///
/// * `HandleResourceStore`: Form verilerini kaydeder
/// * `HandleResourceEdit`: Güncelleme formu için alanları döndürür
/// * `ResolveFieldOptions`: Alan seçeneklerini çözümler
///
/// ## Örnek Kullanım
///
/// ```go
/// // Router tanımı
/// app.Get("/api/resource/:resource/create", func(c *fiber.Ctx) error {
///     ctx := context.New(c)
///     handler := getFieldHandler(ctx.Resource())
///     return HandleResourceCreate(handler, ctx)
/// })
/// ```
///
/// ## Avantajlar
///
/// ✅ Dinamik form oluşturma
/// ✅ Yetkilendirme entegrasyonu
/// ✅ İlişkisel alan desteği
/// ✅ Koşullu görünürlük
/// ✅ Otomatik seçenek yükleme
///
/// ## Dezavantajlar
///
/// ❌ Çok sayıda alan için yavaş olabilir
/// ❌ AutoOptions callback'leri N+1 sorgu problemine yol açabilir
/// ❌ Büyük seçenek listeleri için bellek kullanımı artabilir
func HandleResourceCreate(h *FieldHandler, c *context.Context) error {
	// Adım 1: Yetkilendirme Kontrolü
	// Policy tanımlıysa, kullanıcının Create yetkisi olup olmadığını kontrol et
	// Yetki yoksa 403 Forbidden hatası döndür
	if h.Policy != nil && !h.Policy.Create(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Adım 2: Oluşturma Alanları Listesini Başlat
	// Frontend'e gönderilecek alanları tutacak boş bir slice oluştur
	// Her alan map[string]interface{} formatında JSON serileştirilecek
	createFields := make([]map[string]interface{}, 0)

	// Adım 3: Tüm Alanları Döngüyle İşle
	// Handler'da tanımlı tüm alanları tek tek kontrol et ve filtrele
	for _, element := range h.Elements {
		// Adım 3.1: Görünürlük Kontrolü
		// Alan, mevcut kaynak için görünür değilse atla
		// IsVisible, kullanıcı rolü, kaynak durumu gibi koşullara göre karar verir
		if !element.IsVisible(c.Resource()) {
			continue
		}

		// Adım 3.2: JSON Serileştirme
		// Alanı frontend için uygun JSON formatına dönüştür
		// name, type, label, rules, placeholder gibi özellikleri içerir
		serialized := element.JsonSerialize()

		// Adım 3.3: Seçenek Çözümleme
		// AutoOptions ve callback fonksiyonlarını çalıştır
		// BelongsTo, HasMany gibi ilişkisel alanlar için seçenekleri yükle
		// Örnek: Kategori dropdown'u için veritabanından kategorileri çek
		h.ResolveFieldOptions(element, serialized, nil)

		// Adım 3.4: Bağlam Kontrolü
		// Alanın hangi bağlamda gösterileceğini belirle
		// CREATE, UPDATE, LIST, DETAIL gibi bağlamlar var
		ctxStr := element.GetContext()

		// Adım 3.5: Bağlam Filtreleme
		// Oluşturma formunda gösterilmemesi gereken alanları filtrele
		// Sadece CREATE bağlamına uygun alanları listeye ekle
		if ctxStr != fields.HIDE_ON_CREATE &&      // Oluşturmada gizli değilse
			ctxStr != fields.ONLY_ON_LIST &&       // Sadece listede gösterilmiyorsa
			ctxStr != fields.ONLY_ON_DETAIL &&     // Sadece detayda gösterilmiyorsa
			ctxStr != fields.ONLY_ON_UPDATE {      // Sadece güncellemede gösterilmiyorsa
			// Alan oluşturma formu için uygun, listeye ekle
			createFields = append(createFields, serialized)
		}
	}

	// Adım 4: JSON Yanıtı Döndür
	// Filtrelenmiş ve hazırlanmış alanları frontend'e gönder
	// Frontend bu alanları kullanarak dinamik form oluşturacak
	return c.JSON(fiber.Map{
		"fields": createFields,
	})
}
