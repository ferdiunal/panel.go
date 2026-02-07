package fields

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
			Name: name,
			Key:  key,
			View: "has-many-field",
			Type: TYPE_RELATIONSHIP,
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
	// HasMany fields should not be shown in create/update forms
	// They should be managed in separate interfaces
	h.HideOnCreate()
	h.HideOnUpdate()
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

// GetRelatedResource returns the related resource slug
func (h *HasManyField) GetRelatedResource() string {
	return h.RelatedResourceSlug
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
