package fields

import (
	"reflect"
)

// HasOneField, one-to-one ilişkiyi temsil eder (örn. User -> Profile).
//
// HasOne ilişkisi, bir kaydın tek bir ilişkili kayda sahip olduğunu belirtir.
// Bu, veritabanında ilişkili tabloda foreign key ile temsil edilir.
//
// # Kullanım Senaryoları
//
// - **User -> Profile**: Bir kullanıcının bir profili vardır
// - **Country -> Capital**: Bir ülkenin bir başkenti vardır
// - **Order -> Invoice**: Bir siparişin bir faturası vardır
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
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    WithEagerLoad()
//
//	// Resource instance ile (tip güvenli)
//	field := fields.HasOne("Profile", "profile", user.NewProfileResource()).
//	    ForeignKey("user_id").
//	    WithEagerLoad()
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type HasOneField struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	ForeignKeyColumn    string
	OwnerKeyColumn      string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
}

// HasOne, yeni bir HasOne ilişki alanı oluşturur.
//
// Bu fonksiyon, hem string slug hem de resource instance kabul eder.
// Resource instance kullanımı tip güvenliği sağlar ve refactoring'i kolaylaştırır.
//
// # Parametreler
//
// - **name**: Alanın görünen adı (örn. "Profile", "Profil")
// - **key**: İlişki key'i (örn. "profile")
// - **relatedResource**: İlgili resource (string slug veya resource instance)
//
// # String Slug Kullanımı
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//
// # Resource Instance Kullanımı (Önerilen)
//
//	field := fields.HasOne("Profile", "profile", user.NewProfileResource())
//
// **Avantajlar:**
// - ✅ Tip güvenliği (derleme zamanı kontrolü)
// - ✅ Refactoring desteği
// - ✅ IDE desteği (autocomplete, go-to-definition)
//
// # Varsayılan Değerler
//
// - **ForeignKeyColumn**: slug + "_id" (örn. "profiles_id")
// - **OwnerKeyColumn**: "id" (ana tablonun primary key'i)
// - **LoadingStrategy**: EAGER_LOADING (N+1 sorgu problemini önler)
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı
func HasOne(name, key string, relatedResource interface{}) *HasOneField {
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

	h := &HasOneField{
		Schema: Schema{
			Name:  name,
			Key:   key,
			View:  "has-one-field",
			Type:  TYPE_RELATIONSHIP,
			Props: make(map[string]interface{}),
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
	// Store relationship details in props for generic access (when Schema interface is used)
	h.WithProps("related_resource", slug)
	if resourceInstance != nil {
		h.WithProps("related_resource_instance", resourceInstance)
	}
	h.WithProps("foreign_key", h.ForeignKeyColumn)
	return h
}

// AutoOptions enables automatic options generation from the related table.
// displayField is the column name to use for the option label.
func (h *HasOneField) AutoOptions(displayField string) *HasOneField {
	h.Schema.AutoOptions(displayField)
	return h
}

// ForeignKey sets the foreign key column name
func (h *HasOneField) ForeignKey(key string) *HasOneField {
	h.ForeignKeyColumn = key
	h.WithProps("foreign_key", key)
	return h
}

// OwnerKey sets the owner key column name
func (h *HasOneField) OwnerKey(key string) *HasOneField {
	h.OwnerKeyColumn = key
	return h
}

// Query sets the query callback for customizing relationship query
func (h *HasOneField) Query(fn func(interface{}) interface{}) *HasOneField {
	h.QueryCallback = fn
	return h
}

// WithEagerLoad sets the loading strategy to eager loading
func (h *HasOneField) WithEagerLoad() *HasOneField {
	h.LoadingStrategy = EAGER_LOADING
	return h
}

// WithLazyLoad sets the loading strategy to lazy loading
func (h *HasOneField) WithLazyLoad() *HasOneField {
	h.LoadingStrategy = LAZY_LOADING
	return h
}

// Extract extracts the value from the resource.
// For HasOne, we want to extract the ID of the related resource if it's a struct.
func (h *HasOneField) Extract(resource interface{}) {
	h.Schema.Extract(resource)

	if h.Schema.Data != nil {
		v := reflect.ValueOf(h.Schema.Data)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				h.Schema.Data = nil
				return
			}
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			// Try to find ID or Id field
			idField := v.FieldByName("ID")
			if !idField.IsValid() {
				idField = v.FieldByName("Id")
			}

			if idField.IsValid() && idField.CanInterface() {
				h.Schema.Data = idField.Interface()
			}
		}
	}
}

// GetRelationshipType returns the relationship type
func (h *HasOneField) GetRelationshipType() string {
	return "hasOne"
}

// GetRelatedResource returns the related resource slug
func (h *HasOneField) GetRelatedResource() string {
	return h.RelatedResourceSlug
}

// GetRelationshipName returns the relationship name
func (h *HasOneField) GetRelationshipName() string {
	return h.Name
}

// ResolveRelationship resolves the relationship by loading single related resource
func (h *HasOneField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// In a real implementation, this would query the database
	// For now, return nil
	return nil, nil
}

// ValidateRelationship validates the relationship
func (h *HasOneField) ValidateRelationship(value interface{}) error {
	// Validate that at most one related resource exists
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey returns the display key (not used for HasOne)
func (h *HasOneField) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns returns the searchable columns (not used for HasOne)
func (h *HasOneField) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback returns the query callback
func (h *HasOneField) GetQueryCallback() func(interface{}) interface{} {
	if h.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return h.QueryCallback
}

// GetLoadingStrategy returns the loading strategy
func (h *HasOneField) GetLoadingStrategy() LoadingStrategy {
	if h.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return h.LoadingStrategy
}

// Searchable marks the element as searchable (implements Element interface)
func (h *HasOneField) Searchable() Element {
	h.GlobalSearch = true
	return h
}

// IsRequired returns whether the field is required
func (h *HasOneField) IsRequired() bool {
	return h.Schema.IsRequired
}

// GetTypes returns the type mappings (not used for HasOne)
func (h *HasOneField) GetTypes() map[string]string {
	return make(map[string]string)
}
