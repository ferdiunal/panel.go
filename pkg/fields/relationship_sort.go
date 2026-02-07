package fields

import (
	"context"
)

// RelationshipSort, ilişkiler için sıralama işlevselliğini yönetir.
//
// Bu interface, ilişkili kayıtları sıralamak için metodlar sağlar.
// Tek veya çoklu sıralamalar uygulanabilir.
//
// # Kullanım Örneği
//
//	sort := fields.NewRelationshipSort(field)
//	results, err := sort.ApplySort(ctx, "created_at", "DESC")
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type RelationshipSort interface {
	// ApplySort, ilişki sorgusuna bir sıralama uygular
	ApplySort(ctx context.Context, column string, direction string) ([]interface{}, error)

	// ApplyMultipleSorts, birden fazla sıralama uygular
	ApplyMultipleSorts(ctx context.Context, sorts map[string]string) ([]interface{}, error)

	// RemoveSort, sıralamayı kaldırır ve varsayılan sıralama düzenini kullanır
	RemoveSort(ctx context.Context) ([]interface{}, error)
}

// RelationshipSortImpl, RelationshipSort interface'ini implement eder.
//
// Bu yapı, ilişki sıralama işlemlerini gerçekleştirir.
// Sıralamalar saklanır ve sorguya uygulanır.
type RelationshipSortImpl struct {
	field RelationshipField
	sorts map[string]string
}

// NewRelationshipSort, yeni bir relationship sort handler oluşturur.
//
// Bu fonksiyon, verilen field için sıralama handler'ı döndürür.
//
// # Parametreler
//
// - **field**: İlişki field'ı
//
// # Kullanım Örneği
//
//	sort := fields.NewRelationshipSort(field)
//
// Döndürür:
//   - Yapılandırılmış RelationshipSortImpl pointer'ı
func NewRelationshipSort(field RelationshipField) *RelationshipSortImpl {
	return &RelationshipSortImpl{
		field: field,
		sorts: make(map[string]string),
	}
}

// ApplySort applies a sort to the relationship query
func (rs *RelationshipSortImpl) ApplySort(ctx context.Context, column string, direction string) ([]interface{}, error) {
	if column == "" {
		return []interface{}{}, nil
	}

	// Validate direction
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}

	// Store the sort
	rs.sorts[column] = direction

	// In a real implementation, this would query the database with the sort
	// For now, return empty slice
	return []interface{}{}, nil
}

// ApplyMultipleSorts applies multiple sorts
func (rs *RelationshipSortImpl) ApplyMultipleSorts(ctx context.Context, sorts map[string]string) ([]interface{}, error) {
	if len(sorts) == 0 {
		return []interface{}{}, nil
	}

	// Store all sorts
	for column, direction := range sorts {
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}
		rs.sorts[column] = direction
	}

	// In a real implementation, this would query the database with all sorts applied in order
	// For now, return empty slice
	return []interface{}{}, nil
}

// RemoveSort removes a sort and uses default sort order
func (rs *RelationshipSortImpl) RemoveSort(ctx context.Context) ([]interface{}, error) {
	// Clear all sorts
	rs.sorts = make(map[string]string)

	// In a real implementation, this would load all related resources with default sort order
	// For now, return empty slice
	return []interface{}{}, nil
}

// GetSorts returns the current sorts
func (rs *RelationshipSortImpl) GetSorts() map[string]string {
	return rs.sorts
}
