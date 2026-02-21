package handler

import (
	"strings"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceDetail, tek bir kaynağın detaylı görünümünü getiren HTTP handler fonksiyonudur.
//
// # Genel Bakış
//
// Bu fonksiyon, belirli bir ID'ye sahip kaynağın tüm detay alanlarını getirir, yetkilendirme
// kontrolü yapar ve alanları detay bağlamına göre filtreler. Sonuç olarak, kullanıcı arayüzünde
// gösterilmek üzere düzenlenmiş ve sıralanmış alan verilerini döndürür.
//
// # İşleyiş Akışı
//
// 1. **Kaynak Getirme**: URL parametresinden ID alınır ve Provider üzerinden kaynak getirilir
// 2. **Yetkilendirme**: Policy varsa, kullanıcının kaynağı görüntüleme yetkisi kontrol edilir
// 3. **Alan Filtreleme**: Alanlar detay bağlamına göre filtrelenir (HIDE_ON_DETAIL, ONLY_ON_LIST vb.)
// 4. **Değer Çözümleme**: Her alan için değerler çözümlenir ve formatlanır
// 5. **Sıralama**: Alanlar orijinal tanımlama sırasına göre düzenlenir
// 6. **Yanıt Dönme**: JSON formatında alanlar ve meta bilgiler döndürülür
//
// # Parametreler
//
// - `h *FieldHandler`: Alan işleyici yapısı. Şunları içerir:
//   - `Provider`: Veri sağlayıcı (Show metodunu çağırır)
//   - `Policy`: Yetkilendirme politikası (nil olabilir)
//   - `Elements`: Kaynak için tanımlanmış tüm alanlar
//   - `Title`: Kaynak başlığı
//
// - `c *context.Context`: Panel context yapısı. Şunları sağlar:
//   - `Params("id")`: URL'den kaynak ID'sini alır
//   - `Resource()`: Mevcut kaynak bilgisini döndürür
//   - `Ctx`: Fiber context'i
//
// # Dönüş Değeri
//
// - `error`: İşlem başarılıysa nil, aksi halde hata döner
//   - 404 Not Found: Kaynak bulunamadığında
//   - 403 Forbidden: Kullanıcı yetkisiz olduğunda
//   - 200 OK: Başarılı durumda JSON yanıt
//
// # Yanıt Formatı
//
// ```json
// {
//   "fields": [
//     {
//       "key": "name",
//       "label": "İsim",
//       "value": "Örnek Değer",
//       "type": "text",
//       ...
//     }
//   ],
//   "meta": {
//     "title": "Kaynak Başlığı"
//   }
// }
// ```
//
// # Kullanım Senaryoları
//
// ## Senaryo 1: Basit Kaynak Detayı
// ```go
// // Kullanıcı bir ürünün detaylarını görüntülemek istiyor
// // GET /api/products/123
// // Fonksiyon, ürünün tüm detay alanlarını getirir
// ```
//
// ## Senaryo 2: Yetkilendirme ile Detay
// ```go
// // Kullanıcı sadece kendi kayıtlarını görebilir
// // Policy.View() false dönerse 403 hatası alır
// ```
//
// ## Senaryo 3: Koşullu Alan Görünürlüğü
// ```go
// // Bazı alanlar sadece listeleme sayfasında gösterilir (ONLY_ON_LIST)
// // Bu alanlar detay sayfasında filtrelenir
// ```
//
// # Alan Filtreleme Kuralları
//
// Aşağıdaki bağlamlara sahip alanlar detay görünümünden **hariç tutulur**:
// - `HIDE_ON_DETAIL`: Detayda gizle
// - `ONLY_ON_LIST`: Sadece listede göster
// - `ONLY_ON_FORM`: Sadece formda göster
// - `HIDE_ON_UPDATE`: Güncellemede gizle
//
// # Önemli Notlar
//
// - **Sıralama Garantisi**: Alanlar, FieldHandler.Elements içindeki tanımlama sırasına göre döndürülür
// - **Görünürlük Kontrolü**: Her alan için `IsVisible()` metodu çağrılır
// - **Değer Çözümleme**: `resolveResourceFields()` metodu her alan için değerleri çözümler
// - **Policy Opsiyonel**: Policy nil ise yetkilendirme kontrolü atlanır
// - **ID Parametresi**: URL'den "id" parametresi zorunludur
//
// # Performans Notları
//
// - Provider.Show() tek bir veritabanı sorgusu yapar
// - Alan filtreleme bellekte yapılır (O(n) karmaşıklık)
// - Değer çözümleme her alan için ayrı ayrı yapılabilir (ilişkili alanlar için ek sorgular)
//
// # Hata Durumları
//
// 1. **Kaynak Bulunamadı**: Provider.Show() hata döndüğünde 404 yanıtı
// 2. **Yetkisiz Erişim**: Policy.View() false döndüğünde 403 yanıtı
// 3. **Geçersiz ID**: ID parametresi eksik veya geçersiz olduğunda Provider hatası
//
// # Örnek Kullanım
//
// ```go
// // Router tanımlaması
// app.Get("/api/:resource/:id", func(c *fiber.Ctx) error {
//     ctx := context.New(c)
//     handler := getFieldHandler(ctx.Resource())
//     return HandleResourceDetail(handler, ctx)
// })
// ```
//
// # İlgili Fonksiyonlar
//
// - `HandleResourceList`: Kaynak listesi için kullanılır
// - `HandleResourceCreate`: Yeni kaynak oluşturma için kullanılır
// - `HandleResourceUpdate`: Kaynak güncelleme için kullanılır
// - `resolveResourceFields`: Alan değerlerini çözümler (private metod)
//
// # Güvenlik Notları
//
// - Policy kontrolü mutlaka yapılmalıdır (nil kontrolü önemli)
// - ID parametresi SQL injection'a karşı Provider tarafından korunmalıdır
// - Hassas alanlar için IsVisible() metodunda ek kontroller yapılabilir
func HandleResourceDetail(h *FieldHandler, c *context.Context) error {
	resourceCtx := ensureResourceContext(c, h.Resource, nil, fields.ContextDetail)

	id := c.Params("id")
	item, err := h.Provider.Show(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.View(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Filter for Detail
	var detailElements []fields.Element
	elements := h.getElements(c)
	for _, element := range elements {
		if !element.IsVisible(resourceCtx) {
			continue
		}

		ctxStr := element.GetContext()

		// Skip if explicitly hidden on detail or restricted to other contexts
		if ctxStr != "" {
			contexts := strings.Fields(string(ctxStr))
			shouldSkip := false

			for _, ctx := range contexts {
				if ctx == string(fields.HIDE_ON_DETAIL) ||
					ctx == string(fields.ONLY_ON_LIST) ||
					ctx == string(fields.ONLY_ON_FORM) ||
					ctx == string(fields.HIDE_ON_UPDATE) {
					shouldSkip = true
					break
				}
			}

			if shouldSkip {
				continue
			}
		}

		detailElements = append(detailElements, element)
	}

	// Resolve fields with values
	resolvedMap, err := h.resolveResourceFields(c.Ctx, resourceCtx, item, detailElements)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Convert map to ordered slice based on h.Elements order (preserving filtered list order)
	var orderedFields []map[string]interface{}
	// Iterate over detailElements to preserve order
	for _, element := range detailElements {
		if val, ok := resolvedMap[element.GetKey()]; ok {
			if fieldMap, ok := val.(map[string]interface{}); ok {
				orderedFields = append(orderedFields, fieldMap)
			}
		}
	}

	dialogType, dialogSize := resolveDialogMeta(h)

	return c.JSON(fiber.Map{
		"fields": orderedFields,
		"meta": fiber.Map{
			"title":       h.Resource.TitleWithContext(c.Ctx),
			"dialog_type": dialogType,
			"dialog_size": dialogSize,
		},
	})
}
