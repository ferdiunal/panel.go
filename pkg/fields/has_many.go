package fields

import (
	"encoding/json"
	"reflect"
)

// HasManyField, one-to-many ilişkiyi temsil eder (örn. Author -> Posts).
//
// HasMany ilişkisi, bir kaydın birden fazla ilişkili kayda sahip olduğunu belirtir.
// Bu, veritabanında ilişkili tabloda foreign key ile temsil edilir.
//
// # Kullanım Senaryoları
//
// - **Author -> Posts**: Bir yazar birden fazla yazıya sahiptir
// - **Category -> Products**: Bir kategori birden fazla ürüne sahiptir
// - **User -> Comments**: Bir kullanıcı birden fazla yoruma sahiptir
//
// # Özellikler
//
// - **Tip Güvenliği**: Resource instance veya string slug kullanılabilir
// - **Foreign Key Özelleştirme**: İlişkili tablodaki foreign key sütunu özelleştirilebilir
// - **Owner Key Özelleştirme**: Ana tablodaki referans sütunu özelleştirilebilir
// - **Eager/Lazy Loading**: Yükleme stratejisi seçimi
// - **GORM Yapılandırması**: Foreign key ve references özelleştirme
//
// # Kullanım Örneği
//
//	// String slug ile
//	field := fields.HasMany("Posts", "posts", "posts").
//	    ForeignKey("author_id").
//	    WithEagerLoad()
//
//	// Resource instance ile (tip güvenli)
//	field := fields.HasMany("Posts", "posts", blog.NewPostResource()).
//	    ForeignKey("author_id").
//	    WithEagerLoad()
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type HasManyField struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	ForeignKeyColumn    string
	OwnerKeyColumn      string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
	FullDataMode        bool // true: raw data + title, false (default): minimal format (id + title)
}

// HasMany, yeni bir HasMany ilişki alanı oluşturur.
//
// Bu fonksiyon, hem string slug hem de resource instance kabul eder.
// Resource instance kullanımı tip güvenliği sağlar ve refactoring'i kolaylaştırır.
//
// # Parametreler
//
// - **name**: Alanın görünen adı (örn. "Posts", "Yazılar")
// - **key**: İlişki key'i (örn. "posts")
// - **relatedResource**: İlgili resource (string slug veya resource instance)
//
// # String Slug Kullanımı
//
//	field := fields.HasMany("Posts", "posts", "posts")
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
//	field := fields.HasMany("Posts", "posts", blog.NewPostResource())
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
// - **ForeignKeyColumn**: slug + "_id" (örn. "posts_id")
// - **OwnerKeyColumn**: "id" (ana tablonun primary key'i)
// - **LoadingStrategy**: EAGER_LOADING (N+1 sorgu problemini önler)
//
// Döndürür:
//   - Yapılandırılmış HasManyField pointer'ı
//
// Daha fazla bilgi için docs/Relationships.md ve .docs/RESOURCE_BASED_RELATIONSHIPS.md dosyalarına bakın.
func HasMany(name, key string, relatedResource interface{}) *HasManyField {
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

	h := &HasManyField{
		Schema: Schema{
			LabelText: name,
			Name:      name,
			Key:       key,
			View:      "has-many-field",
			Type:      TYPE_RELATIONSHIP,
		},
		RelatedResourceSlug: slug,
		RelatedResource:     resourceInstance,
		ForeignKeyColumn:    slug + "_id",
		OwnerKeyColumn:      "id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithForeignKey(slug + "_id").
			WithReferences("id"),
	}
	h.WithProps("related_resource", slug)
	if resourceInstance != nil {
		h.WithProps("related_resource_instance", resourceInstance)
	}

	// HasMany relationship requires the parent record to exist first
	h.HideOnCreate()

	return h
}

// ForeignKey, foreign key sütun adını ayarlar.
//
// Bu metod, ilişkili tablodaki foreign key sütununun adını özelleştirir.
// Varsayılan olarak slug + "_id" kullanılır.
//
// # Parametreler
//
// - **key**: Foreign key sütun adı (örn. "author_id", "user_id", "category_id")
//
// # Kullanım Örneği
//
//	field := fields.HasMany("Posts", "posts", "posts").
//	    ForeignKey("author_id")
//	// Posts tablosunda "author_id" sütunu foreign key olarak kullanılır
//
// # Önemli Notlar
//
// - Foreign key sütunu, ilişkili tabloda mevcut olmalıdır
// - Genellikle "parent_table_id" formatında olur (örn. "author_id", "category_id")
// - GORM tarafından otomatik olarak ilişki kurulumu için kullanılır
//
// Döndürür:
//   - HasManyField pointer'ı (method chaining için)
func (h *HasManyField) ForeignKey(key string) *HasManyField {
	h.ForeignKeyColumn = key
	return h
}

// OwnerKey, owner key sütun adını ayarlar.
//
// Bu metod, ana tablodaki referans sütununun adını özelleştirir.
// Varsayılan olarak "id" (primary key) kullanılır.
//
// # Parametreler
//
// - **key**: Owner key sütun adı (örn. "id", "uuid", "author_id")
//
// # Kullanım Örneği
//
//	field := fields.HasMany("Posts", "posts", "posts").
//	    ForeignKey("author_id").
//	    OwnerKey("id")
//	// Authors tablosundaki "id" sütunu ile Posts tablosundaki "author_id" sütunu eşleştirilir
//
// # Önemli Notlar
//
// - Owner key sütunu, ana tabloda mevcut olmalıdır
// - Genellikle primary key kullanılır ("id")
// - Özel durumlar için farklı sütunlar kullanılabilir (örn. "uuid")
//
// Döndürür:
//   - HasManyField pointer'ı (method chaining için)
func (h *HasManyField) OwnerKey(key string) *HasManyField {
	h.OwnerKeyColumn = key
	return h
}

// Query sets the query callback for customizing relationship query
func (h *HasManyField) Query(fn func(interface{}) interface{}) *HasManyField {
	h.QueryCallback = fn
	return h
}

// WithEagerLoad sets the loading strategy to eager loading
func (h *HasManyField) WithEagerLoad() *HasManyField {
	h.LoadingStrategy = EAGER_LOADING
	return h
}

// WithLazyLoad sets the loading strategy to lazy loading
func (h *HasManyField) WithLazyLoad() *HasManyField {
	h.LoadingStrategy = LAZY_LOADING
	return h
}

// GetRelationshipType returns the relationship type
func (h *HasManyField) GetRelationshipType() string {
	return "hasMany"
}

// GetRelationshipName returns the relationship name
func (h *HasManyField) GetRelationshipName() string {
	return h.Name
}

// ResolveRelationship resolves the relationship by loading all related resources
func (h *HasManyField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}

	// In a real implementation, this would query the database
	// For now, return empty slice
	return []interface{}{}, nil
}

// ValidateRelationship validates the relationship
func (h *HasManyField) ValidateRelationship(value interface{}) error {
	// Validate that foreign key references are valid
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for HasMany)
func (h *HasManyField) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for HasMany)
func (h *HasManyField) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (h *HasManyField) GetQueryCallback() func(interface{}) interface{} {
	if h.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return h.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (h *HasManyField) GetLoadingStrategy() LoadingStrategy {
	if h.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return h.LoadingStrategy
}

// Label sets the field label while preserving HasMany concrete type in fluent chains.
func (h *HasManyField) Label(label string) Element {
	h.Schema.LabelText = label
	return h
}

// Placeholder sets the field placeholder while preserving HasMany concrete type in fluent chains.
func (h *HasManyField) Placeholder(placeholder string) Element {
	h.Schema.PlaceholderText = placeholder
	return h
}

// Searchable marks the element as searchable (implements Element interface)
func (h *HasManyField) Searchable() Element {
	h.GlobalSearch = true
	return h
}

// Count returns the count of related resources
func (h *HasManyField) Count() int64 {
	// In a real implementation, this would execute a COUNT query
	return 0
}

// IsRequired returns whether the field is required
func (h *HasManyField) IsRequired() bool {
	return h.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for HasMany)
func (h *HasManyField) GetTypes() map[string]string {
	return make(map[string]string)
}

// AutoOptions, ilişkili tablodan otomatik options oluşturmayı etkinleştirir.
//
// Bu metod, HasMany ilişkisinde ilişkili kayıtların otomatik olarak yüklenmesini
// ve frontend'de multi-select combobox'ta gösterilmesini sağlar.
//
// # Parametreler
//
// - **displayField**: Option label'ı için kullanılacak sütun adı (örn. "title", "name")
//
// # Kullanım Örneği
//
//	// Author resource'unda Posts ilişkisi
//	field := fields.HasMany("Posts", "posts", "posts").
//	    AutoOptions("title").  // Post'ların "title" sütunu label olarak kullanılır
//	    ForeignKey("author_id").
//	    WithEagerLoad()
//
// # Backend Response Formatı
//
// AutoOptions etkinleştirildiğinde, backend response'unda options otomatik olarak eklenir:
//
//	{
//	  "posts": {
//	    "key": "posts",
//	    "type": "relationship",
//	    "view": "has-many-field",
//	    "data": [1, 2, 3],
//	    "props": {
//	      "related_resource": "posts",
//	      "options": {
//	        "1": "First Post",
//	        "2": "Second Post",
//	        "3": "Third Post"
//	      }
//	    }
//	  }
//	}
//
// # Frontend Rendering
//
// Frontend'de HasManyField componenti otomatik olarak:
// - Pre-loaded options'ları multi-select combobox'ta gösterir
// - Search fonksiyonu ile filtreleme yapar
// - Seçili değerleri chips olarak gösterir
// - Kullanıcı yeni kayıtlar seçebilir veya mevcut seçimleri kaldırabilir
//
// # Önemli Notlar
//
// - displayField sütunu, ilişkili tabloda mevcut olmalıdır
// - Büyük veri setleri için (10,000+ kayıt) performans sorunları olabilir
// - Best practice: AutoOptions sadece küçük-orta veri setleri için kullanın
// - Null değerler için otomatik fallback kontrolü yapılır
//
// Döndürür:
//   - HasManyField pointer'ı (method chaining için)
func (h *HasManyField) AutoOptions(displayField string) *HasManyField {
	h.Schema.AutoOptions(displayField)
	return h
}

// GetRelatedTableName, ilişkili tablo adını döndürür.
//
// Bu metod, HasMany ilişkisinde kullanılan ilişkili tablonun adını döndürür.
// Raw SQL sorguları için kullanılır.
//
// # Dönüş Değeri
//
// - İlişkili tablo adı (örn. "posts", "comments", "orders")
//
// # Kullanım Örneği
//
//	field := fields.HasMany("Posts", "posts", "posts")
//	tableName := field.GetRelatedTableName() // "posts"
//
// Döndürür:
//   - İlişkili tablo adı
func (h *HasManyField) GetRelatedTableName() string {
	return h.RelatedResourceSlug
}

// GetForeignKeyColumn, foreign key sütun adını döndürür.
//
// Bu metod, HasMany ilişkisinde kullanılan foreign key sütununun adını döndürür.
// Foreign key, ilişkili tablodaki referans sütunudur.
//
// # Dönüş Değeri
//
// - Foreign key sütun adı (örn. "author_id", "user_id", "category_id")
//
// # Kullanım Örneği
//
//	field := fields.HasMany("Posts", "posts", "posts").ForeignKey("author_id")
//	foreignKey := field.GetForeignKeyColumn() // "author_id"
//
// Döndürür:
//   - Foreign key sütun adı
func (h *HasManyField) GetForeignKeyColumn() string {
	return h.ForeignKeyColumn
}

// GetOwnerKeyColumn, owner key sütun adını döndürür.
//
// Bu metod, HasMany ilişkisinde kullanılan owner key sütununun adını döndürür.
// Owner key, ana tablodaki referans sütunudur (genellikle primary key).
//
// # Dönüş Değeri
//
// - Owner key sütun adı (örn. "id", "uuid")
//
// # Kullanım Örneği
//
//	field := fields.HasMany("Posts", "posts", "posts").OwnerKey("id")
//	ownerKey := field.GetOwnerKeyColumn() // "id"
//
// Döndürür:
//   - Owner key sütun adı
func (h *HasManyField) GetOwnerKeyColumn() string {
	return h.OwnerKeyColumn
}

// WithFullData, ilişkili kayıtların tam verilerini (raw data + title) döndürmesini sağlar.
//
// Bu metod, Extract metodunun davranışını değiştirir:
// - Varsayılan (FullDataMode=false): Minimal format (sadece id + title)
// - WithFullData() çağrıldığında (FullDataMode=true): Full format (tüm raw data + title)
//
// # Kullanım Senaryoları
//
// - Frontend'de ilişkili kayıtların detaylı bilgilerini göstermek
// - Relationship field'larında tüm veriyi kullanmak
// - Custom rendering için ek field'lara ihtiyaç duymak
//
// # Kullanım Örneği
//
//	// Minimal format (varsayılan): {"id": 1, "title": "Adres 1"}
//	fields.HasMany("Addresses", "addresses", "addresses")
//
//	// Full format: {"ID": 1, "Name": "Adres 1", "City": "İstanbul", ..., "title": "Adres 1"}
//	fields.HasMany("Addresses", "addresses", "addresses").
//	    WithFullData()
//
// Döndürür:
//   - HasManyField pointer'ı (method chaining için)
func (h *HasManyField) WithFullData() *HasManyField {
	h.FullDataMode = true
	return h
}

// SetRelatedResource, ilişkili resource instance'ını set eder.
//
// Bu metod, RelatedResource'u runtime'da set etmek için kullanılır.
// Genellikle field_handler.go'da Extract çağrılmadan önce kullanılır.
//
// # Kullanım Senaryoları
//
// - Circular dependency önlemek için string slug kullanıldığında
// - Runtime'da resource registry'den resource instance'ı alındığında
// - Field handler'da Extract öncesi RelatedResource'u set etmek
//
// # Parametreler
//
// - res: Resource instance'ı (resource.Resource interface)
//
// # Kullanım Örneği
//
//	// field_handler.go'da Extract öncesi
//	if relField, ok := element.(*fields.HasManyField); ok {
//	    if relField.RelatedResource == nil {
//	        relatedResource := resource.Get(relField.RelatedResourceSlug)
//	        relField.SetRelatedResource(relatedResource)
//	    }
//	}
//	element.Extract(item)
//
// Döndürür:
//   - HasManyField pointer'ı (method chaining için)
func (h *HasManyField) SetRelatedResource(res interface{}) *HasManyField {
	h.RelatedResource = res
	return h
}

// GetRelatedResourceSlug, ilişkili resource'un slug'ını döndürür.
//
// Bu metod, RelatedResourceSlug'a erişim sağlar.
// Genellikle field_handler.go'da resource registry'den resource instance'ı almak için kullanılır.
//
// # Dönüş Değeri
//
// - string: Resource slug'ı (örn. "users", "posts", "addresses")
//
// # Kullanım Örneği
//
//	// field_handler.go'da
//	if relField, ok := element.(*fields.HasManyField); ok {
//	    slug := relField.GetRelatedResourceSlug()
//	    relatedResource := resource.Get(slug)
//	    relField.SetRelatedResource(relatedResource)
//	}
//
// Döndürür:
//   - Resource slug'ı
func (h *HasManyField) GetRelatedResourceSlug() string {
	return h.RelatedResourceSlug
}

// GetRelatedResource, ilişkili resource instance'ını döndürür.
//
// Bu metod, RelatedResource'a erişim sağlar.
//
// # Dönüş Değeri
//
// - interface{}: Resource instance'ı (resource.Resource interface)
// - nil: RelatedResource set edilmemişse
//
// # Kullanım Örneği
//
//	// Extract metodunda
//	if h.GetRelatedResource() == nil {
//	    // RelatedResource set edilmemiş, minimal format kullan
//	    return
//	}
//
// Döndürür:
//   - Resource instance'ı veya nil
func (h *HasManyField) GetRelatedResource() interface{} {
	return h.RelatedResource
}

// Extract, HasMany ilişkisi için özel veri çıkarma metodudur.
//
// Bu metod, Schema'nın genel Extract metodunu override eder ve HasMany ilişkisi için
// özel işlem yapar. İlişkili kayıtları JSON tag'lere göre serialize eder.
//
// # İşleyiş
//
//  1. Schema'nın Extract metodunu çağırır (mevcut davranışı korumak için)
//  2. Eğer Data nil ise, boş array olarak ayarlar
//  3. Eğer Data bir slice ise:
//     a. Her kayıt için JSON marshaling yapar (JSON tag'ler otomatik kullanılır)
//     b. JSON'dan map[string]interface{}'e unmarshal eder
//     c. Bu map'leri data array'ine ekler
//
// # JSON Tag Kullanımı
//
// Model'de tanımlı JSON tag'ler otomatik olarak kullanılır:
//
//	type Address struct {
//	    ID             uint   `json:"id"`
//	    OrganizationID uint   `json:"organization_id"`
//	    Name           string `json:"name"`
//	}
//
// Sonuç: {"id": 1, "organization_id": 16, "name": "Adres"}
//
// # Parametreler
//
// - resource: Veri çıkarılacak kaynak (struct veya map)
//
// # Kullanım Örneği
//
//	// Organization model'inde Addresses field'ı var
//	type Organization struct {
//	    ID        uint
//	    Name      string
//	    Addresses []Address
//	}
//
//	// HasMany field Extract çağrıldığında
//	field := fields.HasMany("Addresses", "addresses", "addresses")
//	field.Extract(&organization)
//	// field.Data artık [{"id": 1, "organization_id": 16, ...}] gibi JSON tag'lere göre serialize edilmiş kayıtlar içerir
func (h *HasManyField) Extract(resource interface{}) {
	// Schema.Extract ile ilişki verilerini al
	h.Schema.Extract(resource)

	// Data nil ise boş array olarak ayarla
	if h.Schema.Data == nil {
		h.Schema.Data = []interface{}{}
		return
	}

	// RelatedResource yoksa mevcut veriyi kullan
	if h.RelatedResource == nil {
		// Data bir slice değilse boş array döndür
		v := reflect.ValueOf(h.Schema.Data)
		if v.Kind() != reflect.Slice {
			h.Schema.Data = []interface{}{}
		}
		return
	}

	// Data'yı slice olarak işle
	v := reflect.ValueOf(h.Schema.Data)
	if v.Kind() != reflect.Slice {
		h.Schema.Data = []interface{}{}
		return
	}

	// Her kayıt için format oluştur
	// - FullDataMode=false (varsayılan): Minimal format (id + title)
	// - FullDataMode=true: Full format (raw data + title)
	serializedRecords := make([]interface{}, 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)

		// Pointer ise dereference et
		if elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				continue
			}
			elem = elem.Elem()
		}

		// Struct değilse atla
		if elem.Kind() != reflect.Struct {
			continue
		}

		record := elem.Interface()

		// ID field'ını bul
		var idValue interface{}
		for j := 0; j < elem.NumField(); j++ {
			field := elem.Type().Field(j)
			if field.Name == "ID" || field.Name == "Id" {
				idValue = elem.Field(j).Interface()
				break
			}
		}

		// ID bulunamadıysa atla
		if idValue == nil {
			continue
		}

		// RelatedResource'dan RecordTitle metodunu çağır (type assertion ile)
		// RelatedResource interface{} tipinde olduğu için type assertion gerekli
		type ResourceWithTitle interface {
			RecordTitle(any) string
		}

		res, ok := h.RelatedResource.(ResourceWithTitle)
		if !ok {
			// RelatedResource RecordTitle metoduna sahip değilse atla
			continue
		}

		// RecordTitle ile başlığı al (gerekirse struct field fallback kullan)
		recordTitle := resolveRelationshipRecordTitle(res, record, idValue)

		// Format seçimi: Minimal veya Full
		if h.FullDataMode {
			// Full format: raw data + title
			// JSON marshal/unmarshal ile raw data'yı map'e çevir
			jsonData, err := json.Marshal(record)
			if err != nil {
				// JSON marshal hatası, minimal format kullan
				serializedRecords = append(serializedRecords, map[string]interface{}{
					"id":    idValue,
					"title": recordTitle,
				})
				continue
			}

			var recordMap map[string]interface{}
			err = json.Unmarshal(jsonData, &recordMap)
			if err != nil {
				// JSON unmarshal hatası, minimal format kullan
				serializedRecords = append(serializedRecords, map[string]interface{}{
					"id":    idValue,
					"title": recordTitle,
				})
				continue
			}

			// Title ekle
			recordMap["title"] = recordTitle
			serializedRecords = append(serializedRecords, recordMap)
		} else {
			// Minimal format: {"id": ..., "title": ...}
			serializedRecords = append(serializedRecords, map[string]interface{}{
				"id":    idValue,
				"title": recordTitle,
			})
		}
	}

	h.Schema.Data = serializedRecords
}

// Compile-time check: HasManyField implements RelationshipField interface
var _ RelationshipField = (*HasManyField)(nil)
