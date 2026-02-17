package fields

import (
	"fmt"
	"reflect"
	"strings"
)

// BelongsToField, ters one-to-one ilişkiyi temsil eder (örn. Post -> Author).
//
// BelongsTo ilişkisi, bir kaydın başka bir kayda ait olduğunu belirtir.
// Bu, veritabanında foreign key ile temsil edilir.
//
// # Kullanım Senaryoları
//
// - **Post -> Author**: Bir yazı bir yazara aittir
// - **Comment -> User**: Bir yorum bir kullanıcıya aittir
// - **Order -> Customer**: Bir sipariş bir müşteriye aittir
//
// # Özellikler
//
// - **Tip Güvenliği**: Resource instance veya string slug kullanılabilir
// - **Otomatik Seçenekler**: AutoOptions ile veritabanından otomatik seçenek yükleme
// - **Arama Desteği**: İlişkili kayıtlarda arama yapabilme
// - **Eager/Lazy Loading**: Yükleme stratejisi seçimi
// - **GORM Yapılandırması**: Foreign key ve references özelleştirme
// - **Hover Card**: Index ve detail sayfalarında hover card desteği
//
// # Kullanım Örneği
//
//	// String slug ile
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    DisplayUsing("name").
//	    WithSearchableColumns("name", "email").
//	    AutoOptions("name").
//	    WithEagerLoad()
//
//	// Resource instance ile (tip güvenli)
//	field := fields.BelongsTo("Author", "author_id", blog.NewAuthorResource()).
//	    DisplayUsing("name").
//	    WithSearchableColumns("name", "email").
//	    AutoOptions("name").
//	    WithEagerLoad()
//
//	// Hover card ile
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    DisplayUsing("name").
//	    WithHoverCard(fields.NewHoverCardConfig().
//	        WithAvatar("avatar", "").
//	        WithGrid([]fields.HoverCardGridField{
//	            {Key: "email", Label: "Email", Type: "email", Icon: "mail"},
//	            {Key: "role", Label: "Rol", Type: "badge"},
//	        }, "2-column"))
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type BelongsToField struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	DisplayKey          string
	SearchableColumns   []string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
	hoverCardConfig     *HoverCardConfig
}

// BelongsTo, yeni bir BelongsTo ilişki alanı oluşturur.
//
// Bu fonksiyon, hem string slug hem de resource instance kabul eder.
// Resource instance kullanımı tip güvenliği sağlar ve refactoring'i kolaylaştırır.
//
// # Parametreler
//
// - **name**: Alanın görünen adı (örn. "Author", "Yazar")
// - **key**: Foreign key sütun adı (örn. "author_id")
// - **relatedResource**: İlgili resource (string slug veya resource instance)
//
// # String Slug Kullanımı
//
//	field := fields.BelongsTo("Author", "author_id", "authors")
//
// **Avantajlar:**
// - Basit ve hızlı
// - Circular import sorunu yok
//
// **Dezavantajlar:**
// - Tip güvenliği yok
// - Refactoring zor
// - IDE desteği sınırlı
//
// # Resource Instance Kullanımı (Önerilen)
//
//	field := fields.BelongsTo("Author", "author_id", blog.NewAuthorResource())
//
// **Avantajlar:**
// - ✅ Tip güvenliği (derleme zamanı kontrolü)
// - ✅ Refactoring desteği (resource adı değişirse otomatik güncellenir)
// - ✅ IDE desteği (autocomplete, go-to-definition)
// - ✅ Tablo adı otomatik alınır (resource.Slug())
//
// **Dezavantajlar:**
// - Circular import'a dikkat edilmeli
//
// # Varsayılan Değerler
//
// - **DisplayKey**: "name" (görüntüleme için kullanılacak alan)
// - **SearchableColumns**: ["name"] (aranabilir sütunlar)
// - **LoadingStrategy**: EAGER_LOADING (N+1 sorgu problemini önler)
// - **Foreign Key**: key parametresi (örn. "author_id")
// - **References**: "id" (ilgili tablonun primary key'i)
//
// Döndürür:
//   - Yapılandırılmış BelongsToField pointer'ı
//
// Daha fazla bilgi için docs/Relationships.md ve .docs/RESOURCE_BASED_RELATIONSHIPS.md dosyalarına bakın.
func BelongsTo(name, key string, relatedResource interface{}) *BelongsToField {
	// Resource interface'inden slug'ı al
	type resourceSlugger interface {
		Slug() string
	}

	var slug string
	var resourceInstance interface{}

	// Check if relatedResource is a string or a resource instance
	if slugStr, ok := relatedResource.(string); ok {
		// String slug provided
		slug = slugStr
	} else if res, ok := relatedResource.(resourceSlugger); ok {
		// Resource instance provided
		slug = res.Slug()
		resourceInstance = relatedResource
	} else {
		// Fallback: empty slug
		slug = ""
	}

	b := &BelongsToField{
		Schema: Schema{
			LabelText: name,
			Name:      name,
			Key:       key,
			View:      "belongs-to-field",
			Type:      TYPE_RELATIONSHIP,
			Props:     make(map[string]interface{}),
		},
		RelatedResourceSlug: slug,
		RelatedResource:     resourceInstance,
		DisplayKey:          "name",
		SearchableColumns:   []string{"name"},
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithForeignKey(key). // Don't add "_id" suffix - key should already include it
			WithReferences("id"),
	}
	b.WithProps("related_resource", slug)
	if resourceInstance != nil {
		b.WithProps("related_resource_instance", resourceInstance)
	}
	return b
}

// AutoOptions, ilgili tablodan otomatik seçenek oluşturmayı etkinleştirir.
//
// Bu metod, backend'in veritabanından otomatik olarak tüm kayıtları çekip
// form elemanları (Combobox/Select) için seçenekler oluşturmasını sağlar.
// Manuel olarak Options callback'i tanımlamaya gerek kalmaz.
//
// # Parametreler
//
// - **displayField**: Seçenek etiketi için kullanılacak sütun adı (örn. "name", "title", "email")
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    AutoOptions("name")
//	// Backend otomatik olarak authors tablosundan tüm kayıtları çeker
//	// ve "name" sütununu etiket olarak kullanır
//
// # Önemli Notlar
//
// - AutoOptions kullanıldığında, backend otomatik olarak ilgili tablodan tüm kayıtları çeker
// - Büyük tablolar için performans sorunu olabilir, bu durumda Query() ile filtreleme yapılmalıdır
// - displayField, ilgili tabloda mevcut bir sütun olmalıdır
//
// Döndürür:
//   - BelongsToField pointer'ı (method chaining için)
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func (b *BelongsToField) AutoOptions(displayField string) *BelongsToField {
	b.Schema.AutoOptions(displayField)
	return b
}

// DisplayUsing, ilişkili resource'u göstermek için kullanılacak key'i ayarlar.
//
// Bu metod, ilişkili kaydın hangi field'ının görüntüleneceğini belirler.
// Varsayılan olarak "name" field'ı kullanılır.
//
// # Parametreler
//
// - **key**: Görüntüleme için kullanılacak field adı (örn. "name", "title", "email", "username")
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    DisplayUsing("email")
//	// Author'un email'i görüntülenir
//
// # Yaygın Kullanım Senaryoları
//
// - **name**: Genel amaçlı görüntüleme (varsayılan)
// - **title**: Başlık alanları için
// - **email**: E-posta adresi görüntüleme
// - **username**: Kullanıcı adı görüntüleme
// - **full_name**: Tam ad görüntüleme
//
// Döndürür:
//   - BelongsToField pointer'ı (method chaining için)
func (b *BelongsToField) DisplayUsing(key string) *BelongsToField {
	b.DisplayKey = key
	return b
}

// WithSearchableColumns, BelongsTo için aranabilir sütunları ayarlar.
//
// Bu metod, ilişkili kayıtlarda arama yapılabilecek sütunları belirler.
// Bu sütunlar, combobox'ta arama yaparken kullanılır.
//
// # Parametreler
//
// - **columns**: Aranabilir sütun adlarının listesi (örn. "name", "email", "username")
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    WithSearchableColumns("name", "email", "username")
//	// Kullanıcı combobox'ta arama yaptığında name, email ve username sütunlarında arama yapılır
//
// # Önemli Notlar
//
// - Aranabilir sütunlar, ilgili tabloda mevcut olmalıdır
// - Çok fazla sütun eklemek performans sorunlarına neden olabilir
// - Genellikle 2-4 sütun yeterlidir
//
// Döndürür:
//   - BelongsToField pointer'ı (method chaining için)
func (b *BelongsToField) WithSearchableColumns(columns ...string) *BelongsToField {
	b.SearchableColumns = columns
	return b
}

// Searchable, alanı aranabilir olarak işaretler (Element interface'ini implement eder).
//
// Bu metod, alanın global arama işlemlerine dahil edilmesini sağlar.
// Global arama yapıldığında, bu alan da arama sonuçlarına dahil edilir.
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    Searchable()
//	// Bu alan global arama işlemlerine dahil edilir
//
// # Önemli Notlar
//
// - Global arama, tüm aranabilir alanlarda arama yapar
// - WithSearchableColumns() ile birlikte kullanılmalıdır
// - Performans için dikkatli kullanılmalıdır
//
// Döndürür:
//   - Element interface'i (method chaining için)
func (b *BelongsToField) Searchable() Element {
	b.GlobalSearch = true
	return b
}

// Query sets the query callback for customizing relationship query
func (b *BelongsToField) Query(fn func(interface{}) interface{}) *BelongsToField {
	b.QueryCallback = fn
	return b
}

// WithEagerLoad sets the loading strategy to eager loading
func (b *BelongsToField) WithEagerLoad() *BelongsToField {
	b.LoadingStrategy = EAGER_LOADING
	return b
}

// WithLazyLoad sets the loading strategy to lazy loading
func (b *BelongsToField) WithLazyLoad() *BelongsToField {
	b.LoadingStrategy = LAZY_LOADING
	return b
}

// GetRelationshipType returns the relationship type
func (b *BelongsToField) GetRelationshipType() string {
	return "belongsTo"
}

// GetRelatedResource returns the related resource slug
func (b *BelongsToField) GetRelatedResource() string {
	return b.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (b *BelongsToField) GetRelationshipName() string {
	return b.Name
}

// ResolveRelationship resolves the relationship value using reflection
func (b *BelongsToField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		if b.Schema.IsRequired {
			return nil, &RelationshipError{
				FieldName:        b.Name,
				RelationshipType: "belongsTo",
				Message:          "Related resource is required",
				Context: map[string]interface{}{
					"related_resource": b.RelatedResourceSlug,
					"display_key":      b.DisplayKey,
				},
			}
		}
		return nil, nil
	}

	// Extract the relationship value using reflection
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Try to get the field by key name
	fieldVal := v.FieldByName(b.Key)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("field %s not found in %T", b.Key, item)
	}

	return fieldVal.Interface(), nil
}

// ValidateRelationship validates the relationship
func (b *BelongsToField) ValidateRelationship(value interface{}) error {
	if value == nil {
		if b.Schema.IsRequired {
			return &RelationshipError{
				FieldName:        b.Name,
				RelationshipType: "belongsTo",
				Message:          "Related resource is required",
				Context: map[string]interface{}{
					"related_resource": b.RelatedResourceSlug,
				},
			}
		}
		return nil
	}

	// Additional validation logic can be added here
	return nil
}

// GetDisplayKey returns the display key
func (b *BelongsToField) GetDisplayKey() string {
	if b.DisplayKey == "" {
		return "name"
	}
	return b.DisplayKey
}

// GetSearchableColumns returns the searchable columns
func (b *BelongsToField) GetSearchableColumns() []string {
	if b.SearchableColumns == nil {
		return []string{}
	}
	return b.SearchableColumns
}

// GetQueryCallback returns the query callback
func (b *BelongsToField) GetQueryCallback() func(interface{}) interface{} {
	if b.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return b.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (b *BelongsToField) GetLoadingStrategy() LoadingStrategy {
	if b.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return b.LoadingStrategy
}

// IsRequired returns whether the field is required
func (b *BelongsToField) IsRequired() bool {
	return b.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for BelongsTo)
func (b *BelongsToField) GetTypes() map[string]string {
	return make(map[string]string)
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
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    WithHoverCard(*fields.NewHoverCardConfig())
//
// # Yeni API (Önerilen)
//
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    HoverCard(&AuthorHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	        // Custom logic
//	        return &AuthorHoverCard{...}, nil
//	    })
//
// Döndürür:
//   - BelongsToField pointer'ı (method chaining için)
func (b *BelongsToField) WithHoverCard(config HoverCardConfig) *BelongsToField {
	b.hoverCardConfig = &config
	b.WithProps("hover_card", config)
	return b
}

// HoverCard, hover card struct'ını ayarlar ve hover card'ı etkinleştirir.
//
// Bu metod, hover card için kullanılacak struct'ı belirler ve
// hover card özelliğini aktif eder.
//
// # Parametreler
//
// - **hoverStruct**: Hover card verisi için kullanılacak struct (örn. &AuthorHoverCard{})
//
// # Kullanım Örneği
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
//	    HoverCard(&AuthorHoverCard{})
//
// Döndürür:
//   - BelongsToField pointer'ı (method chaining için)
func (b *BelongsToField) HoverCard(hoverStruct interface{}) *BelongsToField {
	if b.hoverCardConfig == nil {
		b.hoverCardConfig = NewHoverCardConfig()
	}
	b.hoverCardConfig.SetStruct(hoverStruct)
	b.WithProps("hover_card_enabled", true)
	return b
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
//	field := fields.BelongsTo("Author", "author_id", "authors").
//	    DisplayUsing("name").
//	    HoverCard(&AuthorHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	        // İlişkili kaydı veritabanından al
//	        author := &Author{}
//	        if err := db.First(author, relatedID).Error; err != nil {
//	            return nil, err
//	        }
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
// Döndürür:
//   - BelongsToField pointer'ı (method chaining için)
func (b *BelongsToField) ResolveHoverCard(resolver HoverCardResolver) *BelongsToField {
	if b.hoverCardConfig == nil {
		b.hoverCardConfig = NewHoverCardConfig()
	}
	b.hoverCardConfig.SetResolver(resolver)
	return b
}

// GetHoverCard, hover card konfigürasyonunu döndürür.
//
// Bu metod, hover card konfigürasyonunu alır.
//
// Döndürür:
//   - HoverCardConfig pointer'ı (nil olabilir)
func (b *BelongsToField) GetHoverCard() *HoverCardConfig {
	return b.hoverCardConfig
}

// GetForeignKey, foreign key sütun adını döndürür.
//
// Bu metod, BelongsTo ilişkisinde kullanılan foreign key sütununun adını döndürür.
// GormRelationConfig'den foreign key bilgisini alır.
//
// # Dönüş Değeri
//
// - Foreign key sütun adı (örn. "author_id", "user_id")
// - GormRelationConfig nil ise boş string döner
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors")
//	foreignKey := field.GetForeignKey() // "author_id"
//
// Döndürür:
//   - Foreign key sütun adı
func (b *BelongsToField) GetForeignKey() string {
	if b.GormRelationConfig != nil {
		return b.GormRelationConfig.ForeignKey
	}
	return ""
}

// GetRelatedResourceSlug, ilişkili resource'un slug'ını döndürür.
//
// Bu metod, BelongsTo ilişkisinde kullanılan ilişkili resource'un slug'ını döndürür.
//
// # Dönüş Değeri
//
// - İlişkili resource slug'ı (örn. "authors", "users", "categories")
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors")
//	slug := field.GetRelatedResourceSlug() // "authors"
//
// Döndürür:
//   - İlişkili resource slug'ı
func (b *BelongsToField) GetRelatedResourceSlug() string {
	return b.RelatedResourceSlug
}

// GetRelatedTableName, ilişkili tablo adını döndürür.
//
// Bu metod, BelongsTo ilişkisinde kullanılan ilişkili tablonun adını döndürür.
// Raw SQL sorguları için kullanılır.
//
// # Dönüş Değeri
//
// - İlişkili tablo adı (örn. "authors", "users", "categories")
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors")
//	tableName := field.GetRelatedTableName() // "authors"
//
// Döndürür:
//   - İlişkili tablo adı
func (b *BelongsToField) GetRelatedTableName() string {
	return b.RelatedResourceSlug
}

// GetOwnerKeyColumn, owner key sütun adını döndürür.
//
// Bu metod, BelongsTo ilişkisinde kullanılan owner key (references) sütununun adını döndürür.
// Owner key, ilişkili tablodaki referans sütunudur (genellikle primary key).
//
// # Dönüş Değeri
//
// - Owner key sütun adı (örn. "id", "uuid")
// - GormRelationConfig nil ise varsayılan olarak "id" döner
//
// # Kullanım Örneği
//
//	field := fields.BelongsTo("Author", "author_id", "authors")
//	ownerKey := field.GetOwnerKeyColumn() // "id"
//
// Döndürür:
//   - Owner key sütun adı
func (b *BelongsToField) GetOwnerKeyColumn() string {
	if b.GormRelationConfig != nil {
		return b.GormRelationConfig.References
	}
	return "id"
}

// Extract, BelongsTo ilişkisini minimal formatta serialize eder.
//
// Bu metod, ilişkili kaydı {"id": ..., "title": ...} formatında döndürür.
// Record title pattern'ini takip eder.
//
// # Çalışma Mantığı
//
// 1. Schema.Extract ile foreign key değerini alır
// 2. İlişki eager loaded ise ve RelatedResource varsa:
//   - İlişkili kaydın ID'sini alır
//   - RelatedResource.RecordTitle ile başlığı alır
//   - {"id": foreignKey, "title": recordTitle} formatında döndürür
//
// 3. İlişki yüklenmemişse veya nil ise null döndürür
//
// # Önemli Notlar
//
// - İlişki eager loaded olmalı (WithEagerLoad veya Preload kullanılmalı)
// - RelatedResource nil ise sadece foreign key değeri döner
// - DisplayUsing callback'i varsa öncelik verilir
//
// # Kullanım Örneği
//
//	// Post -> Author ilişkisi
//	post := &Post{
//	    ID: 1,
//	    AuthorID: 5,
//	    Author: &User{ID: 5, Name: "John Doe"},
//	}
//	field.Extract(&post)
//	// field.Data = {"id": 5, "title": "John Doe"}
//
//	// İlişki yüklenmemiş
//	post := &Post{ID: 1, AuthorID: 5, Author: nil}
//	field.Extract(&post)
//	// field.Data = nil
//
// Parametreler:
//   - record: Kayıt instance'ı (pointer olmalı)
func (b *BelongsToField) Extract(record any) {
	// Schema.Extract ile foreign key değerini al
	b.Schema.Extract(record)

	// Display callback varsa Schema.Extract sonucu korunur
	if b.DisplayCallback != nil {
		return
	}

	// RelatedResource yoksa sadece foreign key değerini döndür
	if b.RelatedResource == nil {
		return
	}

	// Record'un reflection değerini al
	v := reflect.ValueOf(record)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Struct değilse çık
	if v.Kind() != reflect.Struct {
		return
	}

	// İlişki field'ını bul (Schema.Key değeri)
	var relationField reflect.Value
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		// Schema.Key'daki foreign key field adını kullan (örn: "author_id" -> "Author")
		// GORM convention: "author_id" -> "Author" field'ı
		relationName := strings.TrimSuffix(b.Schema.Key, "_id")
		relationName = strings.TrimSuffix(relationName, "Id")
		relationName = strings.TrimSuffix(relationName, "ID")

		if strings.EqualFold(field.Name, relationName) {
			relationField = v.Field(i)
			break
		}
	}

	// İlişki field'ı bulunamadıysa veya nil ise çık
	if !relationField.IsValid() || relationField.IsNil() {
		b.Data = nil
		return
	}

	// İlişkili kaydı al
	relatedRecord := relationField.Interface()

	// Foreign key değerini al
	foreignKeyValue := b.Schema.Data

	// RelatedResource'dan RecordTitle metodunu çağır (type assertion ile)
	// RelatedResource interface{} tipinde olduğu için type assertion gerekli
	type ResourceWithTitle interface {
		RecordTitle(any) string
	}

	res, ok := b.RelatedResource.(ResourceWithTitle)
	if !ok {
		// RelatedResource RecordTitle metoduna sahip değilse sadece foreign key döndür
		return
	}

	// İlişkili kaydın title'ını al (gerekirse struct field fallback kullan)
	recordTitle := resolveRelationshipRecordTitle(res, relatedRecord, foreignKeyValue)

	// Minimal format: {"id": foreignKey, "title": recordTitle}
	b.Data = map[string]any{
		"id":    foreignKeyValue,
		"title": recordTitle,
	}
}
