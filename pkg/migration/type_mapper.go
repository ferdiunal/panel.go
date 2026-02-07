package migration

import (
	"reflect"
	"time"

	"github.com/ferdiunal/panel.go/pkg/fields"
)

// GoType, Go dilindeki tip bilgisini temsil eder.
type GoType struct {
	Type        reflect.Type
	IsPointer   bool
	IsSlice     bool
	ElementType reflect.Type // Slice için element tipi
}

// TypeMapper, Field tiplerini Go ve SQL tiplerine eşler.
type TypeMapper struct{}

// NewTypeMapper, yeni bir TypeMapper oluşturur.
func NewTypeMapper() *TypeMapper {
	return &TypeMapper{}
}

// MapFieldTypeToGo, field tipini Go tipine dönüştürür.
func (tm *TypeMapper) MapFieldTypeToGo(fieldType fields.ElementType, nullable bool) GoType {
	var baseType reflect.Type

	switch fieldType {
	// Metin Tipleri
	case fields.TYPE_TEXT, fields.TYPE_PASSWORD, fields.TYPE_TEXTAREA, fields.TYPE_RICHTEXT:
		baseType = reflect.TypeOf("")
	case fields.TYPE_EMAIL:
		baseType = reflect.TypeOf("")
	case fields.TYPE_TEL:
		baseType = reflect.TypeOf("")

	// Sayısal Tipler
	case fields.TYPE_NUMBER:
		baseType = reflect.TypeOf(int64(0))

	// Boolean
	case fields.TYPE_BOOLEAN:
		baseType = reflect.TypeOf(false)

	// Tarih/Saat Tipleri
	case fields.TYPE_DATE, fields.TYPE_DATETIME:
		baseType = reflect.TypeOf(time.Time{})

	// Dosya Tipleri (URL string olarak saklanır)
	case fields.TYPE_FILE, fields.TYPE_VIDEO, fields.TYPE_AUDIO:
		baseType = reflect.TypeOf("")

	// Seçim Tipleri
	case fields.TYPE_SELECT:
		baseType = reflect.TypeOf("")

	// Key-Value (JSON olarak saklanır)
	case fields.TYPE_KEY_VALUE:
		baseType = reflect.TypeOf(map[string]interface{}{})

	// İlişki Tipleri
	case fields.TYPE_LINK: // BelongsTo -> Foreign Key
		baseType = reflect.TypeOf(uint(0))

	case fields.TYPE_DETAIL: // HasOne
		// Bu tip model struct'ında pointer olarak tanımlanır
		baseType = nil

	case fields.TYPE_COLLECTION: // HasMany
		// Bu tip model struct'ında slice olarak tanımlanır
		baseType = nil

	case fields.TYPE_CONNECT: // BelongsToMany
		// Bu tip model struct'ında slice olarak tanımlanır (many2many)
		baseType = nil

	// Polimorfik İlişkiler
	case fields.TYPE_POLY_LINK, fields.TYPE_POLY_DETAIL, fields.TYPE_POLY_COLLECTION, fields.TYPE_POLY_CONNECT:
		baseType = nil

	default:
		baseType = reflect.TypeOf("")
	}

	result := GoType{
		Type:      baseType,
		IsPointer: nullable && baseType != nil,
	}

	return result
}

// MapFieldTypeToSQL, field tipini SQL tipine dönüştürür.
func (tm *TypeMapper) MapFieldTypeToSQL(fieldType fields.ElementType, size int) string {
	switch fieldType {
	// Metin Tipleri
	case fields.TYPE_TEXT, fields.TYPE_EMAIL, fields.TYPE_TEL, fields.TYPE_PASSWORD:
		if size > 0 {
			return "varchar(" + itoa(size) + ")"
		}
		return "varchar(255)"

	// Uzun Metin Tipleri (TEXT column)
	case fields.TYPE_TEXTAREA, fields.TYPE_RICHTEXT:
		return "text"

	// Sayısal Tipler
	case fields.TYPE_NUMBER:
		return "bigint"

	// Boolean
	case fields.TYPE_BOOLEAN:
		return "boolean"

	// Tarih/Saat Tipleri
	case fields.TYPE_DATE:
		return "date"
	case fields.TYPE_DATETIME:
		return "datetime"

	// Dosya Tipleri (URL saklanır)
	case fields.TYPE_FILE, fields.TYPE_VIDEO, fields.TYPE_AUDIO:
		return "text"

	// Seçim Tipleri
	case fields.TYPE_SELECT:
		return "varchar(100)"

	// Key-Value
	case fields.TYPE_KEY_VALUE:
		return "json"

	// İlişki Tipleri (Foreign Key)
	case fields.TYPE_LINK:
		return "bigint unsigned"

	default:
		return "varchar(255)"
	}
}

// GetRelationshipType, field tipinin ilişki tipini döner.
func (tm *TypeMapper) GetRelationshipType(fieldType fields.ElementType) string {
	switch fieldType {
	case fields.TYPE_LINK:
		return "belongsTo"
	case fields.TYPE_DETAIL:
		return "hasOne"
	case fields.TYPE_COLLECTION:
		return "hasMany"
	case fields.TYPE_CONNECT:
		return "belongsToMany"
	case fields.TYPE_POLY_LINK:
		return "morphTo"
	case fields.TYPE_POLY_DETAIL:
		return "morphOne"
	case fields.TYPE_POLY_COLLECTION:
		return "morphMany"
	case fields.TYPE_POLY_CONNECT:
		return "morphToMany"
	default:
		return ""
	}
}

// IsRelationshipType, field tipinin ilişki tipi olup olmadığını kontrol eder.
func (tm *TypeMapper) IsRelationshipType(fieldType fields.ElementType) bool {
	switch fieldType {
	case fields.TYPE_LINK, fields.TYPE_DETAIL, fields.TYPE_COLLECTION, fields.TYPE_CONNECT,
		fields.TYPE_POLY_LINK, fields.TYPE_POLY_DETAIL, fields.TYPE_POLY_COLLECTION, fields.TYPE_POLY_CONNECT:
		return true
	default:
		return false
	}
}

// Yardımcı fonksiyon
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	negative := false
	if i < 0 {
		negative = true
		i = -i
	}
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
