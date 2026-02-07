// Package fields, veritabanı ilişkilerini alan sisteminde temsil eden yapıları sağlar.
//
// Bu dosya, ilişki alanlarının (BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo)
// temel interface'lerini ve tiplerini tanımlar.
//
// # İlişki Türleri
//
// - **BelongsTo**: Ters one-to-one veya one-to-many ilişki (bir Post bir Author'a aittir)
// - **HasMany**: One-to-many ilişki (bir Author birden fazla Post'a sahiptir)
// - **HasOne**: One-to-one ilişki (bir User bir Profile'a sahiptir)
// - **BelongsToMany**: Many-to-many ilişki (bir User birden fazla Role'e sahiptir)
// - **MorphTo**: Polimorfik ilişki (bir Comment farklı tiplere ait olabilir)
//
// # Yükleme Stratejileri
//
// - **Eager Loading**: N+1 sorgu problemini önler, ilişkili verileri önceden yükler
// - **Lazy Loading**: İhtiyaç anında yükler, bellek tasarrufu sağlar
//
// # Kullanım Örneği
//
//	// BelongsTo ilişkisi
//	field := fields.NewBelongsTo("Author", "author_id", "authors").
//	    DisplayUsing("name").
//	    WithSearchableColumns("name", "email").
//	    WithEagerLoad()
//
//	// HasMany ilişkisi
//	field := fields.NewHasMany("Posts", "posts", "posts").
//	    ForeignKey("author_id").
//	    WithEagerLoad()
//
// Daha fazla bilgi için docs/Relationships.md ve .docs/RESOURCE_BASED_RELATIONSHIPS.md dosyalarına bakın.
package fields

import (
	"context"
	"fmt"
)

// LoadingStrategy, ilişkilerin nasıl yükleneceğini tanımlar.
//
// İki ana strateji vardır:
// - EAGER_LOADING: İlişkili verileri önceden yükler (N+1 sorgu problemini önler)
// - LAZY_LOADING: İlişkili verileri ihtiyaç anında yükler (bellek tasarrufu)
type LoadingStrategy string

const (
	// EAGER_LOADING, ilişkili verileri önceden yükler.
	// N+1 sorgu problemini önlemek için önerilir.
	EAGER_LOADING LoadingStrategy = "eager"

	// LAZY_LOADING, ilişkili verileri ihtiyaç anında yükler.
	// Bellek tasarrufu sağlar ancak N+1 sorgu problemine neden olabilir.
	LAZY_LOADING  LoadingStrategy = "lazy"
)

// RelationshipField, alan sisteminde bir veritabanı ilişkisini temsil eder.
//
// Bu interface, tüm ilişki türleri (BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo)
// için ortak metodları tanımlar.
//
// # Temel Özellikler
//
// - **Tip Bilgisi**: İlişki türünü ve ilgili resource'u belirtir
// - **Çözümleme**: İlişkili verileri yükler ve çözümler
// - **Sorgu Özelleştirme**: İlişki sorgularını özelleştirme callback'leri
// - **Yükleme Stratejisi**: Eager veya lazy loading seçimi
// - **Doğrulama**: İlişki verilerini doğrulama
// - **Görüntüleme**: İlişkili verilerin nasıl gösterileceğini kontrol eder
//
// # Kullanım
//
// RelationshipField interface'i doğrudan kullanılmaz, bunun yerine
// BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo gibi somut tipler kullanılır.
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type RelationshipField interface {
	Element

	// İlişki Türü Metodları

	// GetRelationshipType, ilişki türünü döndürür.
	// Döndürür: "belongsTo", "hasMany", "hasOne", "belongsToMany", "morphTo"
	GetRelationshipType() string

	// GetRelatedResource, ilgili resource'un slug'ını döndürür.
	// Döndürür: İlgili resource'un benzersiz tanımlayıcısı
	GetRelatedResource() string

	// GetRelationshipName, ilişkinin adını döndürür.
	// Döndürür: İlişkinin adı (örn. "author", "posts", "roles")
	GetRelationshipName() string

	// İlişki Çözümleme

	// ResolveRelationship, verilen item için ilişkili verileri çözümler.
	// Parametreler:
	//   - item: İlişkili verileri çözümlenecek kaynak
	// Döndürür:
	//   - İlişkili veriler (tek kayıt veya kayıt listesi)
	//   - Hata (çözümleme başarısız olursa)
	ResolveRelationship(item interface{}) (interface{}, error)

	// Sorgu Özelleştirme

	// GetQueryCallback, sorgu özelleştirme callback'ini döndürür.
	// Döndürür: Sorguyu özelleştiren callback fonksiyonu
	GetQueryCallback() func(interface{}) interface{}

	// Yükleme Stratejisi

	// GetLoadingStrategy, yükleme stratejisini döndürür.
	// Döndürür: EAGER_LOADING veya LAZY_LOADING
	GetLoadingStrategy() LoadingStrategy

	// İlişki Doğrulama

	// ValidateRelationship, ilişki değerini doğrular.
	// Parametreler:
	//   - value: Doğrulanacak değer
	// Döndürür:
	//   - Hata (doğrulama başarısız olursa)
	ValidateRelationship(value interface{}) error

	// İlişki Görüntüleme

	// GetDisplayKey, BelongsTo için görüntülenecek key'i döndürür.
	// Döndürür: Görüntüleme için kullanılacak alan adı (örn. "name", "title")
	GetDisplayKey() string

	// GetSearchableColumns, BelongsTo için aranabilir sütunları döndürür.
	// Döndürür: Aranabilir sütun adlarının listesi
	GetSearchableColumns() []string

	// Zorunluluk Kontrolü

	// IsRequired, ilişkinin zorunlu olup olmadığını döndürür.
	// Döndürür: true ise ilişki zorunludur
	IsRequired() bool

	// MorphTo için Tipler

	// GetTypes, MorphTo için kullanılabilir tipleri döndürür.
	// Döndürür: Tip adı -> resource slug eşlemesi
	GetTypes() map[string]string
}

// RelationshipError, ilişki işlemleri sırasında oluşan bir hatayı temsil eder.
//
// Bu hata tipi, ilişki işlemlerinde oluşan hataları daha iyi anlamak ve
// debug etmek için ek bağlam bilgisi sağlar.
//
// # Özellikler
//
// - **FieldName**: Hatanın oluştuğu alan adı
// - **RelationshipType**: İlişki türü (belongsTo, hasMany, vb.)
// - **Message**: Hata mesajı
// - **Context**: Ek bağlam bilgisi (map formatında)
//
// # Kullanım Örneği
//
//	err := &RelationshipError{
//	    FieldName: "author",
//	    RelationshipType: "belongsTo",
//	    Message: "related resource not found",
//	    Context: map[string]interface{}{
//	        "author_id": 123,
//	        "resource": "authors",
//	    },
//	}
type RelationshipError struct {
	FieldName        string
	RelationshipType string
	Message          string
	Context          map[string]interface{}
}

// Error, error interface'ini implement eder.
// Formatlanmış hata mesajı döndürür.
//
// Döndürür:
//   - Formatlanmış hata mesajı (alan adı, ilişki türü ve mesaj içerir)
func (e *RelationshipError) Error() string {
	return fmt.Sprintf("relationship error in field '%s' (%s): %s", e.FieldName, e.RelationshipType, e.Message)
}

// RelationshipLoader handles loading relationships with different strategies
type RelationshipLoader interface {
	// Load related data using eager loading strategy
	EagerLoad(ctx context.Context, items []interface{}, field RelationshipField) error

	// Load related data using lazy loading strategy
	LazyLoad(ctx context.Context, item interface{}, field RelationshipField) (interface{}, error)

	// Load with constraints applied
	LoadWithConstraints(ctx context.Context, item interface{}, field RelationshipField, constraints map[string]interface{}) (interface{}, error)
}

// RelationshipValidator handles validation of relationships
type RelationshipValidator interface {
	// Validate that related resource exists
	ValidateExists(ctx context.Context, value interface{}, field RelationshipField) error

	// Validate foreign key references
	ValidateForeignKey(ctx context.Context, value interface{}, field RelationshipField) error

	// Validate pivot table entries
	ValidatePivot(ctx context.Context, value interface{}, field RelationshipField) error

	// Validate morph type is registered
	ValidateMorphType(ctx context.Context, value interface{}, field RelationshipField) error
}

// RelationshipQuery represents a query builder for relationships
type RelationshipQuery interface {
	// Apply WHERE clause
	Where(column string, operator string, value interface{}) RelationshipQuery

	// Apply WHERE IN clause
	WhereIn(column string, values []interface{}) RelationshipQuery

	// Apply ORDER BY clause
	OrderBy(column string, direction string) RelationshipQuery

	// Apply LIMIT clause
	Limit(limit int) RelationshipQuery

	// Apply OFFSET clause
	Offset(offset int) RelationshipQuery

	// Get count of results
	Count(ctx context.Context) (int64, error)

	// Check if results exist
	Exists(ctx context.Context) (bool, error)

	// Execute query and get results
	Get(ctx context.Context) ([]interface{}, error)

	// Execute query and get first result
	First(ctx context.Context) (interface{}, error)
}

// IsRelationshipField checks if an element is a relationship field
func IsRelationshipField(e Element) (RelationshipField, bool) {
	rf, ok := e.(RelationshipField)
	return rf, ok
}
