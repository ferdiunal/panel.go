package fields

import (
	"context"
	"strings"
)

// RelationshipSearch, ilişkiler için arama işlevselliğini yönetir.
//
// Bu interface, ilişkili kayıtlarda arama yapmak için metodlar sağlar.
// Belirli sütunlarda veya tüm aranabilir sütunlarda arama yapılabilir.
//
// # Kullanım Örneği
//
//	search := fields.NewRelationshipSearch(field)
//	results, err := search.Search(ctx, "john")
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type RelationshipSearch interface {
	// Search, terime göre ilişkili kaynakları arar
	Search(ctx context.Context, term string) ([]interface{}, error)

	// SearchInColumns, belirli sütunlarda arama yapar
	SearchInColumns(ctx context.Context, term string, columns []string) ([]interface{}, error)

	// GetSearchableColumns, aranabilir sütunları döndürür
	GetSearchableColumns() []string
}

// RelationshipSearchImpl, RelationshipSearch interface'ini implement eder.
//
// Bu yapı, ilişki arama işlemlerini gerçekleştirir.
// Aranabilir sütunlar field'dan alınır.
type RelationshipSearchImpl struct {
	field RelationshipField
}

// NewRelationshipSearch, yeni bir relationship search handler oluşturur.
//
// Bu fonksiyon, verilen field için arama handler'ı döndürür.
//
// # Parametreler
//
// - **field**: İlişki field'ı
//
// # Kullanım Örneği
//
//	search := fields.NewRelationshipSearch(field)
//
// Döndürür:
//   - Yapılandırılmış RelationshipSearchImpl pointer'ı
func NewRelationshipSearch(field RelationshipField) *RelationshipSearchImpl {
	return &RelationshipSearchImpl{
		field: field,
	}
}

// Search searches for related resources by term
func (rs *RelationshipSearchImpl) Search(ctx context.Context, term string) ([]interface{}, error) {
	if term == "" {
		return []interface{}{}, nil
	}

	searchableColumns := rs.GetSearchableColumns()
	if len(searchableColumns) == 0 {
		return []interface{}{}, nil
	}

	return rs.SearchInColumns(ctx, term, searchableColumns)
}

// SearchInColumns searches in specific columns
func (rs *RelationshipSearchImpl) SearchInColumns(ctx context.Context, term string, columns []string) ([]interface{}, error) {
	if term == "" || len(columns) == 0 {
		return []interface{}{}, nil
	}

	// In a real implementation, this would query the database
	// For now, return empty slice
	return []interface{}{}, nil
}

// GetSearchableColumns returns the searchable columns
func (rs *RelationshipSearchImpl) GetSearchableColumns() []string {
	relationType := rs.field.GetRelationshipType()

	switch relationType {
	case "belongsTo":
		return rs.field.GetSearchableColumns()
	case "hasMany":
		return []string{}
	case "hasOne":
		return []string{}
	case "belongsToMany":
		return []string{}
	case "morphTo":
		return []string{}
	default:
		return []string{}
	}
}

// CaseInsensitiveSearch performs case-insensitive search
func (rs *RelationshipSearchImpl) CaseInsensitiveSearch(ctx context.Context, term string) ([]interface{}, error) {
	if term == "" {
		return []interface{}{}, nil
	}

	// Convert term to lowercase for case-insensitive search
	lowerTerm := strings.ToLower(term)

	// In a real implementation, this would query the database with LOWER() function
	// For now, return empty slice
	_ = lowerTerm
	return []interface{}{}, nil
}
