package fields

import (
	"context"
	"fmt"
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

	// Get the related resource slug
	relatedResource := rs.field.GetRelatedResource()
	if relatedResource == "" {
		return []interface{}{}, nil
	}

	// Get database connection from context
	db, ok := ctx.Value("db").(interface {
		Where(query interface{}, args ...interface{}) interface{ Find(dest interface{}) error }
	})
	if !ok || db == nil {
		return []interface{}{}, fmt.Errorf("database connection not found in context")
	}

	// Build search query with OR conditions for each searchable column
	var results []map[string]interface{}
	query := db

	// Build WHERE clause with OR conditions
	whereClause := ""
	whereArgs := []interface{}{}

	for i, column := range columns {
		if i > 0 {
			whereClause += " OR "
		}
		whereClause += fmt.Sprintf("%s LIKE ?", column)
		whereArgs = append(whereArgs, "%"+term+"%")
	}

	// Execute query
	if err := query.Where(whereClause, whereArgs...).Find(&results); err != nil {
		return []interface{}{}, fmt.Errorf("search query failed: %w", err)
	}

	// Convert to []interface{}
	interfaceResults := make([]interface{}, len(results))
	for i, r := range results {
		interfaceResults[i] = r
	}

	return interfaceResults, nil
}

// GetSearchableColumns returns the searchable columns
func (rs *RelationshipSearchImpl) GetSearchableColumns() []string {
	// Try to get searchable columns from the field
	// All relationship fields should implement GetSearchableColumns()
	if searchableField, ok := rs.field.(interface{ GetSearchableColumns() []string }); ok {
		columns := searchableField.GetSearchableColumns()
		if len(columns) > 0 {
			return columns
		}
	}

	// Fallback: return common searchable columns
	// Most resources have at least "name" or "title" fields
	return []string{"name", "title", "email"}
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
