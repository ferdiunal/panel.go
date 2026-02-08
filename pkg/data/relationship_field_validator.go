package data

import (
	"reflect"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/iancoleman/strcase"
)

// ValidateRelationshipFields validates that all BelongsTo fields have corresponding relationship fields in the model.
//
// Bu fonksiyon, resource yüklenirken BelongsTo field'larının model struct'ında ilişki field'larına sahip olup olmadığını kontrol eder.
// Eksik ilişki field'ları için açıklayıcı hata mesajı döndürür.
//
// ## Parametreler
//   - model: Kontrol edilecek model (struct pointer)
//   - elements: Resource'un field'ları
//
// ## Dönüş Değeri
//   - error: Eksik ilişki field'ı varsa MissingRelationshipFieldError, yoksa nil
//
// ## Kullanım Örneği
//
//	type Post struct {
//	    ID       uint64
//	    AuthorID uint64
//	    // Author field'ı eksik - hata verecek
//	}
//
//	elements := []core.Element{
//	    fields.BelongsTo("Author", "author_id", "users"),
//	}
//
//	err := ValidateRelationshipFields(&Post{}, elements)
//	// err != nil - Author field'ı eksik
func ValidateRelationshipFields(model interface{}, elements []core.Element) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for _, element := range elements {
		// Sadece BelongsTo field'larını kontrol et
		if belongsTo, ok := element.(*fields.BelongsToField); ok {
			// İlişki field'ının beklenen adını al (BelongsTo'nun ilk parametresi)
			// Örnek: BelongsTo("Author", "author_id", ...) -> "Author"
			relationshipFieldName := belongsTo.GetName()

			// Model'de bu field var mı kontrol et
			_, found := modelType.FieldByName(relationshipFieldName)
			if !found {
				// Foreign key field'ının adını hesapla
				// Örnek: "author_id" -> "AuthorID"
				foreignKeyName := strcase.ToCamel(belongsTo.GetForeignKey())

				return NewMissingRelationshipFieldError(
					modelType.Name(),
					belongsTo.GetKey(),
					foreignKeyName,
					relationshipFieldName,
					belongsTo.GetRelatedResourceSlug(),
				)
			}
		}
	}

	return nil
}
