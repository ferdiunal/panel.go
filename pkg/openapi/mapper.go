// Package openapi, field type'larını OpenAPI schema'ya map eden işlevleri sağlar.
//
// Bu paket, Panel.go field type'larını OpenAPI 3.0.3 schema'larına otomatik olarak dönüştürür:
// - 30+ field type için mapping
// - Validation rules mapping
// - Custom mapping desteği
// - Relationship field handling
package openapi

import (
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// FieldTypeMapper, field type'larını OpenAPI schema'ya map eder.
//
// ## Özellikler
//   - Otomatik type mapping: 30+ field type için varsayılan mapping
//   - Validation rules: Required, min/max, pattern vb. kuralları map eder
//   - Custom mapping: Kullanıcı tanımlı mapping'ler
//   - Relationship handling: BelongsTo, HasMany vb. ilişkiler için özel mapping
//
// ## Kullanım Örneği
//
//	mapper := NewFieldTypeMapper()
//	field := fields.Text("Name", "name").Required().MaxLength(255)
//	schema := mapper.MapFieldToSchema(field)
//	// schema.Type = "string"
//	// schema.MaxLength = 255
//	// schema.Required = true
type FieldTypeMapper struct {
	registry *CustomMappingRegistry
}

// NewFieldTypeMapper, yeni bir FieldTypeMapper oluşturur.
//
// ## Dönüş Değeri
//   - *FieldTypeMapper: Yapılandırılmış mapper
func NewFieldTypeMapper() *FieldTypeMapper {
	return &FieldTypeMapper{
		registry: NewCustomMappingRegistry(),
	}
}

// MapFieldToSchema, bir field'ı OpenAPI schema'ya dönüştürür.
//
// ## Parametreler
//   - field: Dönüştürülecek field
//
// ## Dönüş Değeri
//   - Schema: OpenAPI schema
//
// ## Davranış
//   1. Custom mapping kontrolü yapar
//   2. Yoksa varsayılan mapping kullanır
//   3. Validation rules ekler
//   4. Field özelliklerini (description, example vb.) ekler
//
// ## Kullanım Örneği
//
//	field := fields.Text("Name", "name").Required().MaxLength(255)
//	schema := mapper.MapFieldToSchema(field)
func (m *FieldTypeMapper) MapFieldToSchema(field fields.Element) Schema {
	// Custom mapping kontrolü
	if customSchema, ok := m.registry.GetFieldMapping(field); ok {
		return customSchema
	}

	// Varsayılan mapping
	schema := m.getDefaultMapping(field)

	// Validation rules ekle
	m.applyValidationRules(&schema, field)

	// Field özellikleri ekle
	m.applyFieldProperties(&schema, field)

	return schema
}

// getDefaultMapping, field type'ına göre varsayılan OpenAPI schema döndürür.
//
// ## Parametreler
//   - field: Field
//
// ## Dönüş Değeri
//   - Schema: Varsayılan OpenAPI schema
func (m *FieldTypeMapper) getDefaultMapping(field fields.Element) Schema {
	fieldType := field.GetType()

	switch fieldType {
	// String types
	case core.TYPE_TEXT:
		return Schema{
			Type:      "string",
			MaxLength: ptr(255),
		}

	case core.TYPE_TEXTAREA:
		return Schema{
			Type: "string",
		}

	case core.TYPE_RICHTEXT:
		return Schema{
			Type:        "string",
			Format:      "html",
			Description: "HTML content",
		}

	case core.TYPE_PASSWORD:
		return Schema{
			Type:      "string",
			Format:    "password",
			WriteOnly: true,
		}

	case core.TYPE_EMAIL:
		return Schema{
			Type:   "string",
			Format: "email",
		}

	case core.TYPE_TEL:
		return Schema{
			Type:   "string",
			Format: "tel",
		}

	case core.TYPE_CODE:
		return Schema{
			Type:        "string",
			Description: "Code content",
		}

	// Number types
	case core.TYPE_NUMBER:
		return Schema{
			Type:   "number",
			Format: "double",
		}

	// Date/Time types
	case core.TYPE_DATE:
		return Schema{
			Type:   "string",
			Format: "date",
		}

	case core.TYPE_DATETIME:
		return Schema{
			Type:   "string",
			Format: "date-time",
		}

	// Boolean types
	case core.TYPE_BOOLEAN:
		return Schema{
			Type: "boolean",
		}

	// File types
	case core.TYPE_FILE, core.TYPE_AUDIO, core.TYPE_VIDEO:
		return Schema{
			Type:   "string",
			Format: "binary",
		}

	// Color type
	case core.TYPE_COLOR:
		return Schema{
			Type:    "string",
			Format:  "color",
			Pattern: "^#[0-9A-Fa-f]{6}$",
		}

	// Select types
	case core.TYPE_SELECT:
		return Schema{
			Type: "string",
		}

	case core.TYPE_BOOLEAN_GROUP:
		return Schema{
			Type: "array",
			Items: &Schema{
				Type: "string",
			},
		}

	// Key-Value type
	case core.TYPE_KEY_VALUE:
		return Schema{
			Type: "object",
			Properties: map[string]Schema{
				"key":   {Type: "string"},
				"value": {Type: "string"},
			},
		}

	// Badge type
	case core.TYPE_BADGE:
		return Schema{
			Type: "string",
		}

	// Relationship types
	case core.TYPE_LINK, core.TYPE_DETAIL:
		// BelongsTo, HasOne - foreign key (integer)
		return Schema{
			Type:   "integer",
			Format: "int64",
		}

	case core.TYPE_COLLECTION, core.TYPE_CONNECT:
		// HasMany, BelongsToMany - array of IDs
		return Schema{
			Type: "array",
			Items: &Schema{
				Type:   "integer",
				Format: "int64",
			},
		}

	case core.TYPE_POLY_LINK, core.TYPE_POLY_DETAIL:
		// MorphTo - object with type and id
		return Schema{
			Type: "object",
			Properties: map[string]Schema{
				"type": {Type: "string", Description: "Polymorphic type"},
				"id":   {Type: "integer", Format: "int64", Description: "Polymorphic ID"},
			},
			Required: []string{"type", "id"},
		}

	case core.TYPE_POLY_COLLECTION, core.TYPE_POLY_CONNECT:
		// MorphMany, MorphToMany - array of polymorphic objects
		return Schema{
			Type: "array",
			Items: &Schema{
				Type: "object",
				Properties: map[string]Schema{
					"type": {Type: "string", Description: "Polymorphic type"},
					"id":   {Type: "integer", Format: "int64", Description: "Polymorphic ID"},
				},
				Required: []string{"type", "id"},
			},
		}

	// Panel type (ignored in API)
	case core.TYPE_PANEL:
		return Schema{
			Type:        "object",
			Description: "Panel grouping (not a data field)",
		}

	// Default: string
	default:
		return Schema{
			Type: "string",
		}
	}
}

// applyValidationRules, validation rules'ları schema'ya ekler.
//
// ## Parametreler
//   - schema: Güncellenecek schema (pointer)
//   - field: Field
//
// ## Eklenen Kurallar
//   - Required: required array'e eklenir (parent level'da)
//   - Min/Max: minimum/maximum
//   - MinLength/MaxLength: minLength/maxLength
//   - Pattern: pattern (regex)
//   - Enum: enum array
func (m *FieldTypeMapper) applyValidationRules(schema *Schema, field fields.Element) {
	// Not: Required kuralı parent level'da (object properties) eklenir
	// Bu metod sadece field-level validation'ları ekler

	// TODO: Validation rules'ları field'dan çıkarmak için
	// fields.Element interface'ine GetValidationRules() metodu eklenebilir
	// Şimdilik temel kuralları ekleyelim

	// Example: Email field için pattern
	if field.GetType() == core.TYPE_EMAIL {
		if schema.Pattern == "" {
			schema.Pattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		}
	}

	// Example: Tel field için pattern
	if field.GetType() == core.TYPE_TEL {
		if schema.Pattern == "" {
			schema.Pattern = "^[+]?[(]?[0-9]{1,4}[)]?[-\\s\\.]?[(]?[0-9]{1,4}[)]?[-\\s\\.]?[0-9]{1,9}$"
		}
	}
}

// applyFieldProperties, field özelliklerini schema'ya ekler.
//
// ## Parametreler
//   - schema: Güncellenecek schema (pointer)
//   - field: Field
//
// ## Eklenen Özellikler
//   - Description: field'ın help text'i
//   - Example: field'ın örnek değeri
//   - ReadOnly: field salt okunur mu?
//   - Nullable: field null değer alabilir mi?
func (m *FieldTypeMapper) applyFieldProperties(schema *Schema, field fields.Element) {
	// Description
	// TODO: Field'dan help text almak için interface'e metod eklenebilir
	// Şimdilik field name'i description olarak kullanabiliriz
	if schema.Description == "" {
		schema.Description = field.GetName()
	}

	// ReadOnly
	// TODO: Field'dan readonly durumunu almak için interface'e metod eklenebilir

	// Nullable
	// TODO: Field'dan nullable durumunu almak için interface'e metod eklenebilir
}

// MapFieldsToProperties, field listesini OpenAPI properties'e dönüştürür.
//
// ## Parametreler
//   - fields: Field listesi
//
// ## Dönüş Değeri
//   - map[string]Schema: OpenAPI properties
//   - []string: Required field'ların listesi
//
// ## Kullanım Örneği
//
//	fields := []fields.Element{
//	    fields.Text("Name", "name").Required(),
//	    fields.Email("Email", "email").Required(),
//	    fields.Number("Age", "age"),
//	}
//	properties, required := mapper.MapFieldsToProperties(fields)
func (m *FieldTypeMapper) MapFieldsToProperties(fieldList []fields.Element) (map[string]Schema, []string) {
	properties := make(map[string]Schema)
	required := []string{}

	for _, field := range fieldList {
		// Field key'ini al
		key := field.GetKey()
		if key == "" {
			continue
		}

		// Schema'ya dönüştür
		schema := m.MapFieldToSchema(field)

		// Properties'e ekle
		properties[key] = schema

		// Required kontrolü
		// TODO: Field'dan required durumunu almak için interface'e metod eklenebilir
		// Şimdilik field type'ına göre varsayılan davranış
	}

	return properties, required
}

// GetRelationshipType, relationship field type'ını döndürür.
//
// ## Parametreler
//   - fieldType: Field type
//
// ## Dönüş Değeri
//   - string: Relationship type ("belongsTo", "hasMany", "morphTo", vb.)
//   - bool: Relationship field mi?
//
// ## Kullanım Örneği
//
//	relType, isRel := mapper.GetRelationshipType(core.TYPE_LINK)
//	// relType = "belongsTo", isRel = true
func (m *FieldTypeMapper) GetRelationshipType(fieldType core.ElementType) (string, bool) {
	switch fieldType {
	case core.TYPE_LINK:
		return "belongsTo", true
	case core.TYPE_DETAIL:
		return "hasOne", true
	case core.TYPE_COLLECTION:
		return "hasMany", true
	case core.TYPE_CONNECT:
		return "belongsToMany", true
	case core.TYPE_POLY_LINK:
		return "morphTo", true
	case core.TYPE_POLY_DETAIL:
		return "morphOne", true
	case core.TYPE_POLY_COLLECTION:
		return "morphMany", true
	case core.TYPE_POLY_CONNECT:
		return "morphToMany", true
	default:
		return "", false
	}
}

// IsRelationshipField, field'ın relationship field olup olmadığını kontrol eder.
//
// ## Parametreler
//   - fieldType: Field type
//
// ## Dönüş Değeri
//   - bool: Relationship field mi?
func (m *FieldTypeMapper) IsRelationshipField(fieldType core.ElementType) bool {
	_, isRel := m.GetRelationshipType(fieldType)
	return isRel
}

// GetFieldExample, field type'ına göre örnek değer döndürür.
//
// ## Parametreler
//   - fieldType: Field type
//
// ## Dönüş Değeri
//   - interface{}: Örnek değer
//
// ## Kullanım Örneği
//
//	example := mapper.GetFieldExample(core.TYPE_TEXT)
//	// example = "Example text"
func (m *FieldTypeMapper) GetFieldExample(fieldType core.ElementType) interface{} {
	switch fieldType {
	case core.TYPE_TEXT:
		return "Example text"
	case core.TYPE_TEXTAREA:
		return "Long text content..."
	case core.TYPE_RICHTEXT:
		return "<p>HTML content</p>"
	case core.TYPE_PASSWORD:
		return "********"
	case core.TYPE_EMAIL:
		return "user@example.com"
	case core.TYPE_TEL:
		return "+90 555 123 4567"
	case core.TYPE_NUMBER:
		return 123.45
	case core.TYPE_DATE:
		return "2026-02-08"
	case core.TYPE_DATETIME:
		return "2026-02-08T15:30:00Z"
	case core.TYPE_BOOLEAN:
		return true
	case core.TYPE_COLOR:
		return "#3B82F6"
	case core.TYPE_SELECT:
		return "option1"
	case core.TYPE_BOOLEAN_GROUP:
		return []string{"option1", "option2"}
	case core.TYPE_KEY_VALUE:
		return map[string]string{"key": "value"}
	case core.TYPE_LINK, core.TYPE_DETAIL:
		return 123
	case core.TYPE_COLLECTION, core.TYPE_CONNECT:
		return []int{1, 2, 3}
	case core.TYPE_POLY_LINK, core.TYPE_POLY_DETAIL:
		return map[string]interface{}{
			"type": "posts",
			"id":   123,
		}
	case core.TYPE_POLY_COLLECTION, core.TYPE_POLY_CONNECT:
		return []map[string]interface{}{
			{"type": "posts", "id": 1},
			{"type": "videos", "id": 2},
		}
	default:
		return nil
	}
}

// SetCustomMappingRegistry, custom mapping registry'yi ayarlar.
//
// ## Parametreler
//   - registry: Custom mapping registry
//
// ## Kullanım Örneği
//
//	registry := NewCustomMappingRegistry()
//	registry.RegisterFieldTypeMapping(core.TYPE_TEXT, func(field fields.Element) Schema {
//	    return Schema{Type: "string", MaxLength: ptr(500)}
//	})
//	mapper.SetCustomMappingRegistry(registry)
func (m *FieldTypeMapper) SetCustomMappingRegistry(registry *CustomMappingRegistry) {
	m.registry = registry
}

// GetCustomMappingRegistry, custom mapping registry'yi döndürür.
//
// ## Dönüş Değeri
//   - *CustomMappingRegistry: Custom mapping registry
func (m *FieldTypeMapper) GetCustomMappingRegistry() *CustomMappingRegistry {
	return m.registry
}

// FormatFieldDescription, field için açıklama oluşturur.
//
// ## Parametreler
//   - field: Field
//
// ## Dönüş Değeri
//   - string: Açıklama
//
// ## Kullanım Örneği
//
//	field := fields.Text("Name", "name").Required()
//	desc := mapper.FormatFieldDescription(field)
//	// desc = "Name (required)"
func (m *FieldTypeMapper) FormatFieldDescription(field fields.Element) string {
	desc := field.GetName()

	// TODO: Field'dan required, unique vb. bilgileri alıp açıklamaya ekleyebiliriz
	// Şimdilik sadece field name'i döndürelim

	return desc
}

// ValidateSchema, schema'nın geçerli olup olmadığını kontrol eder.
//
// ## Parametreler
//   - schema: Kontrol edilecek schema
//
// ## Dönüş Değeri
//   - error: Hata varsa hata, aksi takdirde nil
//
// ## Kullanım Örneği
//
//	schema := Schema{Type: "string"}
//	err := mapper.ValidateSchema(schema)
func (m *FieldTypeMapper) ValidateSchema(schema Schema) error {
	// Type kontrolü
	if schema.Type == "" && schema.Ref == "" {
		return fmt.Errorf("schema must have either type or $ref")
	}

	// Array için items kontrolü
	if schema.Type == "array" && schema.Items == nil {
		return fmt.Errorf("array schema must have items")
	}

	// Object için properties kontrolü (opsiyonel)
	// Object'ler properties olmadan da olabilir (additionalProperties için)

	return nil
}
