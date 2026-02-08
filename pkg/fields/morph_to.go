package fields

import (
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

// MorphTo, polimorfik ilişkiyi temsil eder (örn. Comment -> Commentable).
//
// MorphTo ilişkisi, bir kaydın farklı tiplerdeki kayıtlara ait olabileceğini belirtir.
// Bu, veritabanında morph_type ve morph_id sütunları ile temsil edilir.
//
// # Kullanım Senaryoları
//
// - **Comment -> Commentable**: Bir yorum hem Post'a hem de Video'ya ait olabilir
// - **Image -> Imageable**: Bir resim hem User'a hem de Product'a ait olabilir
// - **Tag -> Taggable**: Bir etiket hem Post'a hem de Video'ya ait olabilir
//
// # Özellikler
//
// - **Tip Eşlemeleri**: Veritabanı tip değerleri ile resource slug'ları arasında eşleme
// - **Görüntüleme Eşlemeleri**: Her tip için görüntüleme alanı özelleştirme
// - **Eager/Lazy Loading**: Yükleme stratejisi seçimi
// - **GORM Yapılandırması**: Polimorfik sütunlar özelleştirme
// - **Hover Card**: Index ve detail sayfalarında hover card desteği
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    Displays(map[string]string{
//	        "post":  "title",
//	        "video": "name",
//	    }).
//	    WithEagerLoad()
//
//	// Hover card ile
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    WithHoverCard(fields.NewHoverCardConfig().
//	        WithAvatar("thumbnail", "").
//	        WithGrid([]fields.HoverCardGridField{
//	            {Key: "title", Label: "Başlık", Type: "text"},
//	            {Key: "created_at", Label: "Tarih", Type: "date"},
//	        }, "2-column"))
//
// # Veritabanı Yapısı
//
// Polimorfik ilişki genellikle şu yapıya sahiptir:
//
//	CREATE TABLE comments (
//	    id INT PRIMARY KEY,
//	    commentable_type VARCHAR(255),  -- "post" veya "video"
//	    commentable_id INT,             -- İlgili kaydın ID'si
//	    content TEXT
//	);
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type MorphTo struct {
	Schema
	TypeMappings       map[string]string // Type => Resource slug mapping
	DisplayMappings    map[string]string // Type => Display field name
	QueryCallback      func(query interface{}) interface{}
	LoadingStrategy    LoadingStrategy
	GormRelationConfig *RelationshipGormConfig
	hoverCardConfig    *HoverCardConfig
}

// NewMorphTo, yeni bir MorphTo polimorfik ilişki alanı oluşturur.
//
// Bu fonksiyon, farklı tiplerdeki kayıtlara ait olabilen polimorfik ilişkiler için kullanılır.
// Veritabanında morph_type ve morph_id sütunları ile temsil edilir.
//
// # Parametreler
//
// - **name**: Alanın görünen adı (örn. "Commentable", "Imageable", "Taggable")
// - **key**: İlişki key'i (örn. "commentable", "imageable", "taggable")
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    Displays(map[string]string{
//	        "post":  "title",
//	        "video": "name",
//	    })
//
// # Varsayılan Değerler
//
// - **TypeMappings**: Boş map (Types() ile doldurulmalıdır)
// - **DisplayMappings**: Boş map (Displays() ile doldurulmalıdır)
// - **LoadingStrategy**: EAGER_LOADING (N+1 sorgu problemini önler)
// - **Polimorfik Sütunlar**: key + "_type" ve key + "_id" (örn. "commentable_type", "commentable_id")
//
// # Önemli Notlar
//
// - Types() metodu ile tip eşlemelerini tanımlamalısınız
// - Displays() metodu ile her tip için görüntüleme alanını belirtmelisiniz
// - Polimorfik sütunlar veritabanında mevcut olmalıdır
//
// Döndürür:
//   - Yapılandırılmış MorphTo pointer'ı
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func NewMorphTo(name, key string) *MorphTo {
	m := &MorphTo{
		Schema: Schema{
			Name: name,
			Key:  key,
			View: "morph-to-field",
			Type: TYPE_RELATIONSHIP,
			Props: map[string]interface{}{
				"types":    []map[string]string{},
				"displays": map[string]string{},
			},
		},
		TypeMappings:    make(map[string]string),
		DisplayMappings: make(map[string]string),
		LoadingStrategy: EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithPolymorphic(key+"_type", key+"_id"),
	}
	// MorphTo fields should not be shown in create/update forms
	// They should be managed in separate interfaces
	m.HideOnCreate()
	m.HideOnUpdate()
	return m
}

// Types, polimorfik ilişki için tip eşlemelerini ayarlar.
//
// Bu metod, veritabanındaki tip değerleri ile resource slug'ları arasında eşleme oluşturur.
// Her tip değeri, hangi resource'a karşılık geldiğini belirtir.
//
// # Parametreler
//
// - **types**: Tip değeri -> resource slug eşlemesi (örn. {"post": "posts", "video": "videos"})
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	        "article": "articles",
//	    })
//	// commentable_type = "post" ise posts resource'u kullanılır
//	// commentable_type = "video" ise videos resource'u kullanılır
//
// # Önemli Notlar
//
// - Tip değerleri veritabanında commentable_type sütununda saklanır
// - Resource slug'ları, ilgili resource'ların benzersiz tanımlayıcılarıdır
// - Tüm olası tip değerleri bu map'te tanımlanmalıdır
// - Frontend'de select dropdown olarak görüntülenir
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) Types(types map[string]string) *MorphTo {
	m.TypeMappings = types
	m.Props["types"] = m.formatTypesForFrontend(types)
	return m
}

// Displays, her tip için görüntüleme alanını ayarlar.
//
// Bu metod, her polimorfik tip için hangi alanın görüntüleneceğini belirler.
// Farklı tipler farklı görüntüleme alanlarına sahip olabilir.
//
// # Parametreler
//
// - **displays**: Tip değeri -> görüntüleme alanı eşlemesi (örn. {"post": "title", "video": "name"})
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    Displays(map[string]string{
//	        "post":  "title",    // Post'lar için title alanı görüntülenir
//	        "video": "name",     // Video'lar için name alanı görüntülenir
//	    })
//
// # Önemli Notlar
//
// - Her tip için görüntüleme alanı belirtilmelidir
// - Görüntüleme alanları, ilgili resource'larda mevcut olmalıdır
// - Frontend'de ilişkili kaydı gösterirken bu alanlar kullanılır
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) Displays(displays map[string]string) *MorphTo {
	m.DisplayMappings = displays
	m.Props["displays"] = displays
	return m
}

// formatTypesForFrontend, tip eşlemelerini frontend select seçeneklerine dönüştürür.
//
// Bu helper metod, backend tip eşlemelerini frontend'in anlayabileceği formata çevirir.
// Her tip için label, value ve slug bilgilerini içeren bir map oluşturur.
//
// # Parametreler
//
// - **types**: Tip değeri -> resource slug eşlemesi
//
// # Döndürür
//
// Frontend select dropdown için seçenek listesi:
//   - **label**: Görüntüleme etiketi (resource slug'ın capitalize edilmiş hali)
//   - **value**: Veritabanı tip değeri (örn. "post", "video")
//   - **slug**: Resource slug'ı (örn. "posts", "videos")
//
// # Örnek Çıktı
//
//	[
//	    {"label": "Posts", "value": "post", "slug": "posts"},
//	    {"label": "Videos", "value": "video", "slug": "videos"}
//	]
//
// # Önemli Notlar
//
// - Label, resource slug'ın ilk harfi büyük yapılarak oluşturulur
// - Value, veritabanında saklanacak tip değeridir
// - Slug, ilgili resource'un benzersiz tanımlayıcısıdır
func (m *MorphTo) formatTypesForFrontend(types map[string]string) []map[string]string {
	var options []map[string]string
	for dbType, resourceSlug := range types {
		label := resourceSlug
		if len(resourceSlug) > 0 {
			label = string(resourceSlug[0]-32) + resourceSlug[1:]
		}

		options = append(options, map[string]string{
			"label": label,
			"value": dbType,
			"slug":  resourceSlug,
		})
	}
	return options
}

// Query, ilişki sorgusunu özelleştirmek için callback ayarlar.
//
// Bu metod, polimorfik ilişki sorgularını özelleştirmek için kullanılır.
// Filtreleme, sıralama veya diğer sorgu modifikasyonları yapılabilir.
//
// # Parametreler
//
// - **fn**: Sorgu özelleştirme callback fonksiyonu
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Query(func(q interface{}) interface{} {
//	        // Sadece aktif kayıtları getir
//	        return q.(*gorm.DB).Where("status = ?", "active")
//	    })
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) Query(fn func(interface{}) interface{}) *MorphTo {
	m.QueryCallback = fn
	return m
}

// WithEagerLoad, yükleme stratejisini eager loading olarak ayarlar.
//
// Eager loading, ilişkili verileri önceden yükler ve N+1 sorgu problemini önler.
// Polimorfik ilişkiler için önerilir.
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    WithEagerLoad()
//	// İlişkili veriler önceden yüklenir
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) WithEagerLoad() *MorphTo {
	m.LoadingStrategy = EAGER_LOADING
	return m
}

// WithLazyLoad, yükleme stratejisini lazy loading olarak ayarlar.
//
// Lazy loading, ilişkili verileri ihtiyaç anında yükler.
// Bellek tasarrufu sağlar ancak N+1 sorgu problemine neden olabilir.
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    WithLazyLoad()
//	// İlişkili veriler ihtiyaç anında yüklenir
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) WithLazyLoad() *MorphTo {
	m.LoadingStrategy = LAZY_LOADING
	return m
}

// GetRelationshipType, ilişki türünü döndürür.
//
// MorphTo için her zaman "morphTo" döndürür.
//
// Döndürür:
//   - "morphTo" string değeri
func (m *MorphTo) GetRelationshipType() string {
	return "morphTo"
}

// GetRelatedResource, ilgili resource slug'ını döndürür.
//
// MorphTo için uygulanamaz çünkü birden fazla resource'a ait olabilir.
// Her zaman boş string döndürür.
//
// Döndürür:
//   - Boş string (MorphTo için uygulanamaz)
func (m *MorphTo) GetRelatedResource() string {
	return ""
}

// GetRelationshipName, ilişkinin adını döndürür.
//
// Döndürür:
//   - İlişkinin adı (örn. "Commentable", "Imageable")
func (m *MorphTo) GetRelationshipName() string {
	return m.Name
}

// ResolveRelationship, polimorfik ilişkiyi morph type'a göre çözümler.
//
// Bu metod, reflection kullanarak struct'tan morph_type ve morph_id değerlerini çıkarır.
// Farklı field adı varyasyonlarını (CamelCase, snake_case) destekler.
//
// # Parametreler
//
// - **item**: İlişkili verileri çözümlenecek kaynak
//
// # Döndürür
//
// Polimorfik ilişki bilgilerini içeren map:
//   - **type**: Morph type değeri (örn. "post", "video")
//   - **id**: Morph ID değeri (ilgili kaydın ID'si)
//   - **morphToType**: Morph type değeri (alias)
//   - **morphToId**: Morph ID değeri (alias)
//
// # Desteklenen Field Adları
//
// - Type için: "CommentableType", "Commentable_Type"
// - ID için: "CommentableID", "CommentableId", "Commentable_ID", "Commentable_Id"
//
// # Kullanım Örneği
//
//	type Comment struct {
//	    ID              int
//	    CommentableType string  // "post" veya "video"
//	    CommentableID   int     // İlgili kaydın ID'si
//	    Content         string
//	}
//
//	resolved, err := field.ResolveRelationship(comment)
//	// resolved = {"type": "post", "id": 123, "morphToType": "post", "morphToId": 123}
//
// Döndürür:
//   - İlişki bilgileri map'i veya nil
//   - Hata (çözümleme başarısız olursa)
func (m *MorphTo) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	val := reflect.ValueOf(item)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, nil
	}

	// Calculate field names
	// Key is usually camelCase or snake_case in JS, but we want the struct field name prefix.
	// e.g. "commentable" -> "Commentable"
	baseName := strcase.ToCamel(m.Key)

	typeFieldNames := []string{baseName + "Type", baseName + "_Type"}
	idFieldNames := []string{baseName + "ID", baseName + "Id", baseName + "_ID", baseName + "_Id"}

	var typeVal string
	var idVal interface{}
	var foundType, foundID bool

	// Find Type
	for _, name := range typeFieldNames {
		f := val.FieldByName(name)
		if f.IsValid() {
			typeVal = f.String()
			foundType = true
			break
		}
	}

	// Find ID
	for _, name := range idFieldNames {
		f := val.FieldByName(name)
		if f.IsValid() {
			idVal = f.Interface()
			foundID = true
			break
		}
	}

	if !foundType && !foundID {
		return nil, nil
	}

	return map[string]interface{}{
		"type":        typeVal,
		"id":          idVal,
		"morphToType": typeVal,
		"morphToId":   idVal,
	}, nil
}

// Extract, MorphTo'ya özgü veri çıkarma işlemini gerçekleştirir.
//
// Bu metod, Schema.Extract'ı override eder çünkü MorphTo struct'ta doğrudan bir field'a sahip değildir.
// Bunun yerine, morph_type ve morph_id field'larından polimorfik ilişkiyi çözümler.
//
// # Parametreler
//
// - **resource**: Veri çıkarılacak kaynak
//
// # Önemli Notlar
//
// - Schema.Extract çağrılmaz çünkü MorphTo doğrudan bir field'a sahip değildir
// - ResolveRelationship kullanılarak polimorfik ilişki çözümlenir
// - Çözümlenen veriler m.Data'ya atanır
//
// # Kullanım Örneği
//
//	field.Extract(comment)
//	// m.Data = {"type": "post", "id": 123, "morphToType": "post", "morphToId": 123}
func (m *MorphTo) Extract(resource interface{}) {
	// Don't call Schema.Extract because MorphTo doesn't have a direct field in the struct
	// Instead, directly resolve the polymorphic relationship from type and id fields
	resolved, _ := m.ResolveRelationship(resource)
	m.Data = resolved
}

// ValidateRelationship, ilişkiyi doğrular.
//
// Bu metod, polimorfik ilişkinin geçerli olup olmadığını kontrol eder.
// Gerçek implementasyonda, morph type'ın tip eşlemelerinde kayıtlı olup olmadığını kontrol eder.
//
// # Parametreler
//
// - **value**: Doğrulanacak değer
//
// # Döndürür
//
// - Hata (doğrulama başarısız olursa)
//
// # Önemli Notlar
//
// - Morph type, TypeMappings'te tanımlanmış olmalıdır
// - Gerçek implementasyonda veritabanı kısıtlamaları kontrol edilir
func (m *MorphTo) ValidateRelationship(value interface{}) error {
	// Validate that morph type is registered
	// In a real implementation, this would check that the type exists in the mapping
	return nil
}

// GetDisplayKey, görüntüleme key'ini döndürür.
//
// MorphTo için kullanılmaz çünkü her tip farklı görüntüleme alanına sahip olabilir.
// Her zaman boş string döndürür.
//
// Döndürür:
//   - Boş string (MorphTo için uygulanamaz)
func (m *MorphTo) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns, aranabilir sütunları döndürür.
//
// MorphTo için kullanılmaz çünkü polimorfik ilişkilerde doğrudan arama yapılmaz.
// Her zaman boş slice döndürür.
//
// Döndürür:
//   - Boş string slice (MorphTo için uygulanamaz)
func (m *MorphTo) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback, sorgu callback'ini döndürür.
//
// Sorgu özelleştirme callback'i tanımlanmışsa onu döndürür,
// aksi takdirde varsayılan (no-op) callback döndürür.
//
// Döndürür:
//   - Sorgu özelleştirme callback fonksiyonu
func (m *MorphTo) GetQueryCallback() func(interface{}) interface{} {
	if m.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return m.QueryCallback
}

// GetLoadingStrategy, yükleme stratejisini döndürür.
//
// Yükleme stratejisi tanımlanmışsa onu döndürür,
// aksi takdirde varsayılan EAGER_LOADING döndürür.
//
// Döndürür:
//   - EAGER_LOADING veya LAZY_LOADING
func (m *MorphTo) GetLoadingStrategy() LoadingStrategy {
	if m.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return m.LoadingStrategy
}

// GetTypes, tip eşlemelerini döndürür.
//
// Bu metod, polimorfik ilişki için tanımlanmış tüm tip eşlemelerini döndürür.
// Tip eşlemeleri, veritabanı tip değerleri ile resource slug'ları arasındaki ilişkiyi belirtir.
//
// Döndürür:
//   - Tip değeri -> resource slug eşlemesi (örn. {"post": "posts", "video": "videos"})
//   - Boş map (tip eşlemeleri tanımlanmamışsa)
func (m *MorphTo) GetTypes() map[string]string {
	if m.TypeMappings == nil {
		return make(map[string]string)
	}
	return m.TypeMappings
}

// GetResourceForType, verilen tip için resource slug'ını döndürür.
//
// Bu metod, morph type değerine karşılık gelen resource slug'ını bulur.
// Tip kayıtlı değilse hata döndürür.
//
// # Parametreler
//
// - **morphType**: Morph type değeri (örn. "post", "video")
//
// # Kullanım Örneği
//
//	resource, err := field.GetResourceForType("post")
//	// resource = "posts"
//
//	resource, err := field.GetResourceForType("unknown")
//	// err = "morph type 'unknown' is not registered"
//
// Döndürür:
//   - Resource slug'ı (tip kayıtlıysa)
//   - Hata (tip kayıtlı değilse)
func (m *MorphTo) GetResourceForType(morphType string) (string, error) {
	resource, ok := m.TypeMappings[morphType]
	if !ok {
		return "", fmt.Errorf("morph type '%s' is not registered", morphType)
	}
	return resource, nil
}

// Searchable, alanı aranabilir olarak işaretler.
//
// Bu metod, alanın global arama işlemlerine dahil edilmesini sağlar.
// Element interface'ini implement eder.
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Searchable()
//	// Bu alan global arama işlemlerine dahil edilir
//
// Döndürür:
//   - Element interface'i (method chaining için)
func (m *MorphTo) Searchable() Element {
	m.GlobalSearch = true
	return m
}

// IsRequired, alanın zorunlu olup olmadığını döndürür.
//
// Döndürür:
//   - true ise alan zorunludur
func (m *MorphTo) IsRequired() bool {
	return m.Schema.IsRequired
}

// WithHoverCard, hover card konfigürasyonunu ayarlar.
//
// Bu metod, index ve detail sayfalarında ilişkili kaydın hover card ile
// nasıl görüntüleneceğini belirler.
//
// # Parametreler
//
// - **config**: Hover card konfigürasyonu
//
// # Kullanım Örneği (Deprecated - Yeni API kullanın)
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    WithHoverCard(*fields.NewHoverCardConfig())
//
// # Yeni API (Önerilen)
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    HoverCard(&CommentableHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	        // Custom logic
//	        return &CommentableHoverCard{...}, nil
//	    })
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) WithHoverCard(config HoverCardConfig) *MorphTo {
	m.hoverCardConfig = &config
	m.WithProps("hover_card", config)
	return m
}

// HoverCard, hover card struct'ını ayarlar ve hover card'ı etkinleştirir.
//
// Bu metod, hover card için kullanılacak struct'ı belirler ve
// hover card özelliğini aktif eder.
//
// # Parametreler
//
// - **hoverStruct**: Hover card verisi için kullanılacak struct (örn. &CommentableHoverCard{})
//
// # Kullanım Örneği
//
//	type CommentableHoverCard struct {
//	    Thumbnail string `json:"thumbnail"`
//	    Title     string `json:"title"`
//	    Type      string `json:"type"`
//	}
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    HoverCard(&CommentableHoverCard{})
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) HoverCard(hoverStruct interface{}) *MorphTo {
	if m.hoverCardConfig == nil {
		m.hoverCardConfig = NewHoverCardConfig()
	}
	m.hoverCardConfig.SetStruct(hoverStruct)
	m.WithProps("hover_card_enabled", true)
	return m
}

// ResolveHoverCard, hover card verilerini çözmek için callback fonksiyonunu ayarlar.
//
// Bu metod, hover card açıldığında çağrılacak resolver fonksiyonunu belirler.
// Resolver, ilişkili kaydın hover card verilerini döndürür.
//
// # Parametreler
//
// - **resolver**: Hover card resolver callback fonksiyonu
//
// # Kullanım Örneği
//
//	field := fields.NewMorphTo("Commentable", "commentable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    HoverCard(&CommentableHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	        // MorphTo için tip bilgisini al
//	        morphType := record.(*Comment).CommentableType
//
//	        // Tip'e göre ilişkili kaydı al
//	        var data interface{}
//	        switch morphType {
//	        case "post":
//	            post := &Post{}
//	            if err := db.First(post, relatedID).Error; err != nil {
//	                return nil, err
//	            }
//	            data = post
//	        case "video":
//	            video := &Video{}
//	            if err := db.First(video, relatedID).Error; err != nil {
//	                return nil, err
//	            }
//	            data = video
//	        }
//
//	        // Hover card verisini döndür
//	        return &CommentableHoverCard{
//	            Thumbnail: data.Thumbnail,
//	            Title: data.Title,
//	            Type: morphType,
//	        }, nil
//	    })
//
// # API Endpoint
//
// Frontend, hover card açıldığında şu endpoint'e istek atar:
//
//	GET /api/resource/{resource}/resolver/{field_name}?id={related_id}&type={morph_type}
//	POST /api/resource/{resource}/resolver/{field_name} (body: {id: related_id, type: morph_type})
//
// Döndürür:
//   - MorphTo pointer'ı (method chaining için)
func (m *MorphTo) ResolveHoverCard(resolver HoverCardResolver) *MorphTo {
	if m.hoverCardConfig == nil {
		m.hoverCardConfig = NewHoverCardConfig()
	}
	m.hoverCardConfig.SetResolver(resolver)
	return m
}

// GetHoverCard, hover card konfigürasyonunu döndürür.
//
// Bu metod, hover card konfigürasyonunu alır.
//
// Döndürür:
//   - HoverCardConfig pointer'ı (nil olabilir)
func (m *MorphTo) GetHoverCard() *HoverCardConfig {
	return m.hoverCardConfig
}
