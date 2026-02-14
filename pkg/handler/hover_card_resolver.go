package handler

import (
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// HandleHoverCardResolve, hover card verilerini çözmek için kullanılan HTTP handler fonksiyonudur.
//
// # Genel Bakış
//
// Bu fonksiyon, ilişki field'larının (BelongsTo, HasOne, MorphTo) hover card verilerini
// çözmek için kullanılır. Frontend, hover card açıldığında bu endpoint'e istek atar ve
// field'ın ResolveHoverCard callback'i çağrılarak hover card verisi döndürülür.
//
// # Kullanım Senaryoları
//
// 1. **BelongsTo Hover Card**: Bir yazının yazarının hover card verilerini gösterme
// 2. **HasOne Hover Card**: Bir kullanıcının profilinin hover card verilerini gösterme
// 3. **MorphTo Hover Card**: Polimorfik ilişkinin hover card verilerini gösterme
//
// # HTTP Endpoint Detayları
//
// - **Route**: `/api/resource/:resource/resolver/:field`
// - **Methods**: `GET`, `POST`, `PATCH`, `DELETE`
// - **URL Parametreleri**:
//   - `:resource` - Kaynak adı (örn: "posts", "users")
//   - `:field` - Field adı (örn: "author_id", "profile")
//
// # Request Parameters (Query veya Body)
//
// - **id** (required): İlişkili kaydın ID'si
// - **type** (optional): MorphTo için morph type (örn: "post", "video")
// - **record_id** (optional): Ana kaydın ID'si (context için)
//
// # Request Örnekleri
//
// GET isteği:
// ```
// GET /api/resource/posts/resolver/author_id?id=5
// GET /api/resource/comments/resolver/commentable?id=10&type=post
// ```
//
// POST isteği:
// ```json
// POST /api/resource/posts/resolver/author_id
//
//	{
//	  "id": 5,
//	  "record_id": 123
//	}
//
// ```
//
// # Response Format
//
// Başarılı durumda:
// ```json
//
//	{
//	  "data": {
//	    "avatar": "https://example.com/avatar.jpg",
//	    "name": "John Doe",
//	    "email": "john@example.com",
//	    "phone": "+1 234 567 8900"
//	  }
//	}
//
// ```
//
// Hata durumlarında:
// ```json
//
//	{
//	  "error": "Field not found"
//	}
//
// ```
// veya
// ```json
//
//	{
//	  "error": "Hover card not configured for this field"
//	}
//
// ```
//
// # Parametreler
//
//   - `h *FieldHandler`: Field handler'ı, field'ların tanımlandığı ve yönetildiği yapı.
//     Bu yapı üzerinden field listesine (Elements) erişilir.
//
//   - `c *context.Context`: Fiber context wrapper'ı. HTTP request/response işlemleri için kullanılır.
//     URL parametrelerine, query parametrelerine, request body'sine ve response yazma işlemlerine erişim sağlar.
//
// # Dönüş Değeri
//
//   - `error`: İşlem başarılı ise nil, hata durumunda error döner.
//     Fiber framework'ü bu error'ı otomatik olarak HTTP response'a dönüştürür.
//
// # İşlem Akışı
//
// 1. URL'den field adı alınır
// 2. Query veya body'den parametreler alınır (id, type, record_id)
// 3. Field adına göre Elements listesinde arama yapılır
// 4. Field bulunamazsa 404 hatası döner
// 5. Field'ın hover card konfigürasyonu kontrol edilir
// 6. Hover card konfigürasyonu yoksa veya resolver yoksa 400 hatası döner
// 7. Resolver callback'i çağrılır ve hover card verisi alınır
// 8. Hover card verisi JSON response olarak döndürülür
//
// # Kullanım Örneği
//
// Backend'de field tanımı:
//
//	type AuthorHoverCard struct {
//	    Avatar string `json:"avatar"`
//	    Name   string `json:"name"`
//	    Email  string `json:"email"`
//	    Phone  string `json:"phone"`
//	}
//
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    DisplayUsing("name").
//	    HoverCard(&AuthorHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
//	        author := &Author{}
//	        if err := db.First(author, relatedID).Error; err != nil {
//	            return nil, err
//	        }
//	        return &AuthorHoverCard{
//	            Avatar: author.Avatar,
//	            Name:   author.Name,
//	            Email:  author.Email,
//	            Phone:  author.Phone,
//	        }, nil
//	    })
//
// Frontend'de kullanım:
//
//	const response = await axios.get('/api/resource/posts/resolver/author_id?id=5')
//	// response.data = { avatar: "...", name: "...", email: "...", phone: "..." }
//
// # Güvenlik Notları
//
// - Resolver callback'leri içinde authorization kontrolü yapılmalıdır
// - Hassas veriler döndürülmeden önce filtrelenmelidir
// - Rate limiting uygulanmalıdır
//
// # Performans Notları
//
// - Resolver callback'leri cache'lenebilir
// - N+1 sorgu problemine dikkat edilmelidir
// - Gereksiz veri döndürülmemelidir
func HandleHoverCardResolve(h *FieldHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// URL'den field adını al
		fieldName := c.Params("field")
		if fieldName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Field name is required",
			})
		}

		// Query veya body'den parametreleri al
		var params struct {
			ID       interface{} `json:"id" query:"id"`
			Type     string      `json:"type" query:"type"`           // MorphTo için
			RecordID interface{} `json:"record_id" query:"record_id"` // Ana kayıt ID'si (context için)
		}

		// Query parametrelerini parse et
		if err := c.QueryParser(&params); err != nil {
			// Query'de yoksa body'den parse et
			if err := c.BodyParser(&params); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid request parameters",
				})
			}
		}

		// ID parametresi zorunlu
		if params.ID == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID parameter is required",
			})
		}

		// Field'ı bul
		var targetField fields.Element
		for _, field := range h.getElements(&context.Context{Ctx: c}) {
			if field.GetKey() == fieldName {
				targetField = field
				break
			}
		}

		if targetField == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Field not found",
			})
		}

		// Field'ın RelationshipField olup olmadığını kontrol et
		relationshipField, ok := targetField.(fields.RelationshipField)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Field is not a relationship field",
			})
		}

		// Hover card konfigürasyonunu al
		var hoverCardConfig *fields.HoverCardConfig

		// Field tipine göre hover card konfigürasyonunu al
		switch field := relationshipField.(type) {
		case *fields.BelongsToField:
			hoverCardConfig = field.GetHoverCard()
		case *fields.HasOneField:
			hoverCardConfig = field.GetHoverCard()
		case *fields.MorphTo:
			hoverCardConfig = field.GetHoverCard()
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Hover card not supported for this field type",
			})
		}

		// Hover card konfigürasyonu kontrolü
		if hoverCardConfig == nil || !hoverCardConfig.Enabled {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Hover card not configured for this field",
			})
		}

		// Resolver kontrolü
		if hoverCardConfig.Resolver == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Hover card resolver not configured",
			})
		}

		// Ana kaydı al (eğer record_id varsa)
		var record interface{}
		if params.RecordID != nil && h.Provider != nil {
			// RecordID'yi string'e çevir
			recordIDStr := fmt.Sprintf("%v", params.RecordID)
			var err error
			record, err = h.Provider.Show(&context.Context{Ctx: c}, recordIDStr)
			if err != nil {
				// Record bulunamadı ama devam edebiliriz (resolver'da nil olabilir)
				record = nil
			}
		}

		// Resolver'ı çağır
		hoverCardData, err := hoverCardConfig.Resolver(c.Context(), record, params.ID, relationshipField)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Hover card verisini döndür
		return c.JSON(fiber.Map{
			"data": hoverCardData,
		})
	}
}
