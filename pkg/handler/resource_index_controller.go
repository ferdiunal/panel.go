package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/query"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceIndex
//
// Bu fonksiyon, kaynak (resource) listelerini sayfalama, sıralama ve filtreleme özellikleriyle birlikte
// yönetir ve istemciye döndürür. REST API'nin temel listeleme endpoint'i olarak çalışır ve
// `/api/resource/:resource` yoluna gelen GET isteklerini işler.
//
// # Temel İşlevsellik
//
// HandleResourceIndex, bir kaynağın tüm kayıtlarını listelemek için kullanılan ana controller fonksiyonudur.
// Bu fonksiyon, modern web uygulamalarında yaygın olarak kullanılan tablo görünümleri, veri grid'leri ve
// liste sayfaları için gerekli tüm veriyi hazırlar ve döndürür.
//
// # Parametreler
//
// - `h *FieldHandler`: Kaynak için tanımlanmış alan (field) yapılandırmalarını, policy kurallarını,
//   veri sağlayıcısını (provider) ve diğer kaynak özelliklerini içeren handler nesnesi.
//   Bu parametre üzerinden kaynağın tüm yapılandırmasına erişilir.
//
// - `c *context.Context`: HTTP isteği bağlamını (context) içeren özel context nesnesi. Fiber context'ini
//   wrap eder ve kaynak-spesifik bilgilere erişim sağlar. URL parametreleri, query string'ler,
//   authentication bilgileri bu nesne üzerinden alınır.
//
// # Dönüş Değeri
//
// - `error`: İşlem başarılı ise nil döner. Hata durumunda uygun HTTP status kodu ile birlikte
//   JSON formatında hata mesajı döndürülür ve error nesnesi return edilir.
//
// # Desteklenen Query Formatları
//
// Fonksiyon iki farklı query string formatını destekler:
//
// ## 1. Nested (İç İçe) Format (Önerilen)
//
// Modern ve daha organize bir yapı sunar. Kaynak adı prefix olarak kullanılır:
//
// ```
// GET /api/resource/users?users[search]=john&users[sort][id]=asc&users[filters][status][eq]=active
// ```
//
// Örnekler:
// - Arama: `users[search]=john`
// - Sıralama: `users[sort][created_at]=desc`
// - Filtreleme: `users[filters][role][eq]=admin`
// - Sayfalama: `users[page]=2&users[per_page]=25`
//
// ## 2. Legacy (Eski) Format
//
// Geriye dönük uyumluluk için desteklenir:
//
// ```
// GET /api/resource/users?search=john&sort_column=id&sort_direction=asc
// ```
//
// # İşlem Akışı
//
// 1. **Policy Kontrolü**: ViewAny policy kontrolü yapılır. Kullanıcının listeyi görme yetkisi yoksa
//    403 Forbidden hatası döner.
//
// 2. **Element Belirleme**: Gösterilecek alanlar (elements) belirlenir. Öncelik sırası:
//    - Context'ten gelen elements (varsa)
//    - Handler'dan gelen default elements
//    - Hiç element yoksa 500 Internal Server Error
//
// 3. **Query Parsing**: URL'den gelen parametreler parse edilir:
//    - Sayfa numarası (page)
//    - Sayfa başına kayıt sayısı (per_page)
//    - Sıralama kriterleri (sorts)
//    - Arama terimi (search)
//    - Filtreler (filters)
//
// 4. **Sıralama Varsayılanları**: Eğer sıralama belirtilmemişse:
//    - Önce Resource'tan tanımlı varsayılan sıralamalar kullanılır
//    - Hiç sıralama yoksa `created_at DESC` kullanılır
//
// 5. **Veri Çekme**: Provider üzerinden veriler çekilir. Provider, veritabanı sorgusunu
//    gerçekleştirir ve sonuçları döndürür.
//
// 6. **Resource Mapping**: Her kayıt için:
//    - Alanlar (fields) resolve edilir ve değerleri çıkarılır
//    - Kayıt-bazlı policy kontrolleri yapılır (view, update, delete)
//    - Policy sonuçları her kaydın içine eklenir
//
// 7. **Header Oluşturma**: Frontend tablo görünümü için header bilgileri oluşturulur.
//    Sadece liste görünümünde gösterilmesi gereken alanlar header'a eklenir.
//
// 8. **Response Dönme**: JSON formatında data ve meta bilgileri döndürülür.
//
// # Response Yapısı
//
// ```json
// {
//   "data": [
//     {
//       "id": 1,
//       "name": "John Doe",
//       "email": "john@example.com",
//       "policy": {
//         "view": true,
//         "update": true,
//         "delete": false
//       }
//     }
//   ],
//   "meta": {
//     "current_page": 1,
//     "per_page": 15,
//     "total": 100,
//     "dialog_type": "modal",
//     "title": "Users",
//     "headers": [
//       {
//         "key": "id",
//         "label": "ID",
//         "sortable": true
//       }
//     ],
//     "policy": {
//       "create": true,
//       "view_any": true,
//       "update": true,
//       "delete": true
//     }
//   }
// }
// ```
//
// # Kullanım Senaryoları
//
// ## Senaryo 1: Basit Listeleme
//
// Kullanıcı bir kaynağın tüm kayıtlarını görmek ister:
//
// ```
// GET /api/resource/users
// ```
//
// Varsayılan sayfalama ve sıralama ile tüm kullanıcılar listelenir.
//
// ## Senaryo 2: Arama ve Filtreleme
//
// Kullanıcı aktif admin kullanıcıları aramak ister:
//
// ```
// GET /api/resource/users?users[search]=john&users[filters][role][eq]=admin&users[filters][status][eq]=active
// ```
//
// ## Senaryo 3: Özel Sıralama ve Sayfalama
//
// Kullanıcı en yeni kayıtları görmek ister:
//
// ```
// GET /api/resource/users?users[sort][created_at]=desc&users[page]=1&users[per_page]=50
// ```
//
// ## Senaryo 4: Frontend Tablo Entegrasyonu
//
// Frontend bir veri tablosu render eder. Bu fonksiyon:
// - Tablo için gerekli header bilgilerini sağlar
// - Her satır için gösterilecek verileri hazırlar
// - Sıralama, filtreleme ve sayfalama için gerekli meta bilgileri döndürür
// - Her satır için kullanılabilir aksiyonları (view, edit, delete) belirtir
//
// # Policy Sistemi
//
// Fonksiyon iki seviyede policy kontrolü yapar:
//
// ## 1. Liste Seviyesi (ViewAny)
//
// Kullanıcının listeyi görme yetkisi kontrol edilir. Yetki yoksa 403 döner.
//
// ## 2. Kayıt Seviyesi (View, Update, Delete)
//
// Her kayıt için ayrı ayrı policy kontrolleri yapılır ve sonuçlar kayıtla birlikte döndürülür.
// Bu sayede frontend, her satır için hangi aksiyonların kullanılabilir olduğunu bilir.
//
// # Hata Durumları
//
// - **403 Forbidden**: Kullanıcının liste görme yetkisi yok
// - **500 Internal Server Error**:
//   - Kaynak için hiç field tanımlanmamış
//   - Provider'dan veri çekilirken hata oluştu
//
// # Performans Notları
//
// - Fonksiyon, her kayıt için field resolution yapar. Çok sayıda kayıt ve karmaşık field'lar
//   durumunda performans etkilenebilir.
// - Provider seviyesinde veritabanı sorgusu optimize edilmelidir (index'ler, eager loading vb.)
// - Header bilgileri her istekte yeniden oluşturulur. Cache mekanizması eklenebilir.
//
// # Önemli Uyarılar
//
// ⚠️ **Element Kontrolü**: Kaynak için mutlaka en az bir element tanımlanmalıdır. Aksi halde
// fonksiyon 500 hatası döndürür.
//
// ⚠️ **Policy Null Kontrolü**: Policy nil ise tüm işlemler izinli kabul edilir. Production
// ortamında mutlaka policy tanımlanmalıdır.
//
// ⚠️ **Varsayılan Sıralama**: Hiçbir sıralama belirtilmezse `created_at DESC` kullanılır.
// Bu alan tabloda yoksa hata oluşur. Resource tanımında uygun varsayılan sıralama belirtilmelidir.
//
// ⚠️ **Query Format**: Nested format kullanılması önerilir. Legacy format gelecekte
// kaldırılabilir.
//
// # Avantajlar
//
// ✅ **Esnek Query Sistemi**: İki farklı query formatı desteği ile geriye dönük uyumluluk
//
// ✅ **Granüler Policy Kontrolü**: Hem liste hem kayıt seviyesinde yetkilendirme
//
// ✅ **Frontend-Ready**: Frontend için gerekli tüm bilgiler (headers, policy, meta) hazır şekilde döner
//
// ✅ **Varsayılan Değerler**: Sıralama ve diğer parametreler için akıllı varsayılanlar
//
// ✅ **Context Desteği**: Context üzerinden element override edilebilir, farklı senaryolar için
// farklı field setleri kullanılabilir
//
// # Dezavantajlar
//
// ❌ **Performans**: Çok sayıda kayıt için field resolution maliyetli olabilir
//
// ❌ **N+1 Problem**: İlişkili veriler için eager loading yapılmazsa N+1 sorunu oluşabilir
//
// ❌ **Header Overhead**: Her istekte header bilgileri yeniden oluşturulur
//
// # İlgili Fonksiyonlar
//
// - `resolveResourceFields`: Kayıt için field değerlerini çıkarır
// - `Provider.Index`: Veritabanından veri çeker
// - `query.ParseResourceQuery`: Query string'i parse eder
//
// # Örnek Kullanım (Handler Tanımı)
//
// ```go
// userHandler := &handler.FieldHandler{
//     Resource: userResource,
//     Provider: gormProvider,
//     Elements: []fields.Element{
//         fields.NewID(),
//         fields.NewText("name").SetLabel("Name"),
//         fields.NewText("email").SetLabel("Email"),
//     },
//     Policy: userPolicy,
//     DialogType: "modal",
// }
//
// app.Get("/api/resource/users", func(c *fiber.Ctx) error {
//     ctx := context.New(c)
//     return handler.HandleResourceIndex(userHandler, ctx)
// })
// ```
func HandleResourceIndex(h *FieldHandler, c *context.Context) error {
	ctx := c.Resource()

	if h.Policy != nil && !h.Policy.ViewAny(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Determine elements to use: Context > Handler Defaults
	var elements []fields.Element
	if ctx != nil && len(ctx.Elements) > 0 {
		elements = ctx.Elements
	} else {
		elements = h.Elements
	}

	if len(elements) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No fields defined for this resource",
		})
	}

	// Parse Query Request using new parser
	resourceName := c.Params("resource")
	queryParams := query.ParseResourceQuery(c.Ctx, resourceName)

	// Convert query.Sort to data.Sort
	var sorts []data.Sort
	for _, s := range queryParams.Sorts {
		sorts = append(sorts, data.Sort{
			Column:    s.Column,
			Direction: s.Direction,
		})
	}

	// Apply defaults from Resource if no sorts provided
	if len(sorts) == 0 {
		if h.Resource != nil {
			for _, s := range h.Resource.GetSortable() {
				sorts = append(sorts, data.Sort{
					Column:    s.Column,
					Direction: s.Direction,
				})
			}
		}
		// Absolute fallback
		if len(sorts) == 0 {
			sorts = append(sorts, data.Sort{
				Column:    "created_at",
				Direction: "desc",
			})
		}
	}

	// Build QueryRequest
	req := data.QueryRequest{
		Page:    queryParams.Page,
		PerPage: queryParams.PerPage,
		Sorts:   sorts,
		Search:  queryParams.Search,
		Filters: queryParams.Filters,
	}

	// Fetch Data
	result, err := h.Provider.Index(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Map items to resources with fields extracted
	resources := make([]map[string]interface{}, 0)

	for _, item := range result.Items {
		res := h.resolveResourceFields(c.Ctx, c.Resource(), item, elements)
		// Inject per-item policy
		policy := map[string]bool{
			"view":   h.Policy == nil || h.Policy.View(c, item),
			"update": h.Policy == nil || h.Policy.Update(c, item),
			"delete": h.Policy == nil || h.Policy.Delete(c, item),
		}
		res["policy"] = policy
		resources = append(resources, res)
	}

	// Generate headers for frontend table order
	headers := make([]map[string]interface{}, 0)

	for _, element := range elements {
		if !element.IsVisible(c.Resource()) {
			continue
		}

		ctxStr := element.GetContext()
		serialized := element.JsonSerialize()

		// logic for headers (Index/List)
		if ctxStr != fields.HIDE_ON_LIST &&
			ctxStr != fields.ONLY_ON_CREATE &&
			ctxStr != fields.ONLY_ON_UPDATE &&
			ctxStr != fields.ONLY_ON_FORM &&
			ctxStr != fields.ONLY_ON_DETAIL {
			headers = append(headers, serialized)
		}
	}

	return c.JSON(fiber.Map{
		"data": resources,
		"meta": fiber.Map{
			"current_page": result.Page,
			"per_page":     result.PerPage,
			"total":        result.Total,
			"dialog_type":  h.DialogType,
			"title":        h.Resource.Title(),
			"headers":      headers,
			"policy": fiber.Map{
				"create":   h.Policy == nil || h.Policy.Create(c),
				"view_any": h.Policy == nil || h.Policy.ViewAny(c),
				"update":   h.Policy == nil || h.Policy.Update(c, nil),
				"delete":   h.Policy == nil || h.Policy.Delete(c, nil),
			},
		},
	})
}
