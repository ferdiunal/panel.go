package fields

import "context"

// HoverCardResolver, hover card verilerini çözmek için kullanılan callback fonksiyonu tipi.
//
// Bu fonksiyon, hover card açıldığında çağrılır ve ilişkili kaydın
// hover card verilerini döndürür.
//
// # Parametreler
//
// - **ctx**: Context
// - **record**: Ana kayıt (örn. Post, Comment, vb.)
// - **relatedID**: İlişkili kaydın ID'si
// - **field**: Field instance
//
// # Döndürür
//
// - **interface{}**: Hover card verisi (struct veya map)
// - **error**: Hata (varsa)
//
// # Kullanım Örneği
//
//	resolver := func(ctx context.Context, record interface{}, relatedID interface{}, field Field) (interface{}, error) {
//	    // İlişkili kaydı veritabanından al
//	    author := &Author{}
//	    db.First(author, relatedID)
//
//	    // Hover card verisini döndür
//	    return &AuthorHoverCard{
//	        Avatar: author.Avatar,
//	        Name:   author.Name,
//	        Email:  author.Email,
//	        Phone:  author.Phone,
//	    }, nil
//	}
type HoverCardResolver func(ctx context.Context, record interface{}, relatedID interface{}, field RelationshipField) (interface{}, error)

// HoverCardConfig, ilişki field'ları için hover card görüntüleme ayarlarını yönetir.
//
// Bu yapı, hasOne, belongsTo ve morphTo field'larının index ve detail sayfalarında
// hover card ile nasıl görüntüleneceğini kontrol eder.
//
// # Özellikler
//
// - **Enabled**: Hover card'ın aktif olup olmadığını belirler
// - **Struct**: Hover card verisi için kullanılacak struct (tip bilgisi)
// - **Resolver**: Hover card verilerini çözmek için callback fonksiyonu
// - **Width**: Hover card genişliği (örn. "sm", "md", "lg", "xl")
// - **OpenDelay**: Hover card açılma gecikmesi (ms)
// - **CloseDelay**: Hover card kapanma gecikmesi (ms)
//
// # Kullanım Örneği
//
//	// Hover card struct'ı tanımla
//	type AuthorHoverCard struct {
//	    Avatar string `json:"avatar"`
//	    Name   string `json:"name"`
//	    Email  string `json:"email"`
//	    Phone  string `json:"phone"`
//	}
//
//	// Field'a hover card ekle
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    DisplayUsing("name").
//	    HoverCard(&AuthorHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	        // İlişkili kaydı veritabanından al
//	        author := &Author{}
//	        db.First(author, relatedID)
//
//	        // Hover card verisini döndür
//	        return &AuthorHoverCard{
//	            Avatar: author.Avatar,
//	            Name:   author.Name,
//	            Email:  author.Email,
//	            Phone:  author.Phone,
//	        }, nil
//	    })
//
// # API Endpoint
//
// Frontend, hover card açıldığında şu endpoint'e istek atar:
//
//	GET /api/resource/{resource}/resolver/{field_name}?id={related_id}
//	POST /api/resource/{resource}/resolver/{field_name} (body: {id: related_id})
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type HoverCardConfig struct {
	// Enabled, hover card'ın aktif olup olmadığını belirler
	Enabled bool `json:"enabled"`

	// Struct, hover card verisi için kullanılacak struct (tip bilgisi)
	// Bu, JSON serialization için kullanılır
	Struct interface{} `json:"-"`

	// Resolver, hover card verilerini çözmek için callback fonksiyonu
	Resolver HoverCardResolver `json:"-"`

	// Width, hover card genişliği (örn. "sm", "md", "lg", "xl")
	Width string `json:"width,omitempty"`

	// OpenDelay, hover card açılma gecikmesi (ms)
	OpenDelay int `json:"open_delay,omitempty"`

	// CloseDelay, hover card kapanma gecikmesi (ms)
	CloseDelay int `json:"close_delay,omitempty"`
}

// NewHoverCardConfig, varsayılan hover card konfigürasyonu oluşturur.
//
// Bu fonksiyon, temel hover card ayarlarıyla bir konfigürasyon döndürür.
//
// # Varsayılan Değerler
//
// - **Enabled**: true
// - **Width**: "md"
// - **OpenDelay**: 200ms
// - **CloseDelay**: 300ms
//
// # Kullanım Örneği
//
//	config := fields.NewHoverCardConfig()
//
// Döndürür:
//   - Varsayılan HoverCardConfig pointer'ı
func NewHoverCardConfig() *HoverCardConfig {
	return &HoverCardConfig{
		Enabled:    true,
		Width:      "md",
		OpenDelay:  200,
		CloseDelay: 300,
	}
}

// SetStruct, hover card verisi için kullanılacak struct'ı ayarlar.
//
// Bu metod, hover card verisi için tip bilgisini belirler.
//
// # Parametreler
//
// - **s**: Hover card struct'ı (örn. &AuthorHoverCard{})
//
// # Kullanım Örneği
//
//	config := fields.NewHoverCardConfig()
//	config.SetStruct(&AuthorHoverCard{})
//
// Döndürür:
//   - HoverCardConfig pointer'ı (method chaining için)
func (h *HoverCardConfig) SetStruct(s interface{}) *HoverCardConfig {
	h.Struct = s
	return h
}

// SetResolver, hover card verilerini çözmek için callback fonksiyonunu ayarlar.
//
// Bu metod, hover card verilerini almak için kullanılacak resolver'ı belirler.
//
// # Parametreler
//
// - **resolver**: Hover card resolver callback fonksiyonu
//
// # Kullanım Örneği
//
//	config := fields.NewHoverCardConfig()
//	config.SetResolver(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	    // Custom logic
//	    return &AuthorHoverCard{...}, nil
//	})
//
// Döndürür:
//   - HoverCardConfig pointer'ı (method chaining için)
func (h *HoverCardConfig) SetResolver(resolver HoverCardResolver) *HoverCardConfig {
	h.Resolver = resolver
	return h
}

// SetWidth, hover card genişliğini ayarlar.
//
// Bu metod, hover card'ın genişliğini belirler.
//
// # Parametreler
//
// - **width**: Genişlik değeri ("sm", "md", "lg", "xl")
//
// # Kullanım Örneği
//
//	config := fields.NewHoverCardConfig()
//	config.SetWidth("lg")
//
// Döndürür:
//   - HoverCardConfig pointer'ı (method chaining için)
func (h *HoverCardConfig) SetWidth(width string) *HoverCardConfig {
	h.Width = width
	return h
}

// SetDelays, hover card açılma ve kapanma gecikmelerini ayarlar.
//
// Bu metod, hover card'ın açılma ve kapanma gecikmelerini belirler.
//
// # Parametreler
//
// - **openDelay**: Açılma gecikmesi (ms)
// - **closeDelay**: Kapanma gecikmesi (ms)
//
// # Kullanım Örneği
//
//	config := fields.NewHoverCardConfig()
//	config.SetDelays(300, 500)
//
// Döndürür:
//   - HoverCardConfig pointer'ı (method chaining için)
func (h *HoverCardConfig) SetDelays(openDelay, closeDelay int) *HoverCardConfig {
	h.OpenDelay = openDelay
	h.CloseDelay = closeDelay
	return h
}

// Disable, hover card'ı devre dışı bırakır.
//
// Bu metod, hover card gösterimini tamamen kapatır.
//
// # Kullanım Örneği
//
//	config := fields.NewHoverCardConfig()
//	config.Disable()
//
// Döndürür:
//   - HoverCardConfig pointer'ı (method chaining için)
func (h *HoverCardConfig) Disable() *HoverCardConfig {
	h.Enabled = false
	return h
}
