package data

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

// MissingRelationshipFieldError, BelongsTo field'ı tanımlanmış ama model struct'ında ilişki field'ı eksik olduğunda döndürülen hata.
//
// Bu hata, kullanıcıya eksik ilişki field'ını ve nasıl düzelteceğini açıkça gösterir.
//
// ## Örnek Hata Mesajı
//
//	╔══════════════════════════════════════════════════════════════════╗
//	║ GORM Relationship Error: Missing Relationship Field             ║
//	╚══════════════════════════════════════════════════════════════════╝
//
//	Model: PriceList
//	Field: BelongsTo("CargoCompany", "cargo_company_id", "cargo_companies")
//
//	Problem: Struct is missing the relationship field
//
//	Fix: Add this field to your PriceList struct:
//
//	    CargoCompany *cargo_company.CargoCompany `gorm:"foreignKey:CargoCompanyID;references:ID"`
//
//	Docs: https://panel.go/docs/relationships#belongs-to
type MissingRelationshipFieldError struct {
	ModelName           string // Model adı (örn: "PriceList")
	BelongsToFieldKey   string // BelongsTo field key (örn: "cargo_company")
	ForeignKey          string // Foreign key field adı (örn: "CargoCompanyID")
	ExpectedFieldName   string // Beklenen ilişki field adı (örn: "CargoCompany")
	RelatedResourceSlug string // İlişkili resource slug (örn: "cargo_companies")
}

// NewMissingRelationshipFieldError, yeni bir MissingRelationshipFieldError oluşturur.
func NewMissingRelationshipFieldError(modelName, fieldKey, foreignKey, expectedFieldName, relatedSlug string) *MissingRelationshipFieldError {
	return &MissingRelationshipFieldError{
		ModelName:           modelName,
		BelongsToFieldKey:   fieldKey,
		ForeignKey:          foreignKey,
		ExpectedFieldName:   expectedFieldName,
		RelatedResourceSlug: relatedSlug,
	}
}

// Error, hata mesajını döndürür.
//
// Hata mesajı, kullanıcıya:
// 1. Hangi model'de sorun olduğunu
// 2. Hangi BelongsTo field'ının eksik olduğunu
// 3. Nasıl düzelteceğini (tam kod örneği ile)
// 4. Dokümantasyon linkini
//
// gösterir.
func (e *MissingRelationshipFieldError) Error() string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("╔══════════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║ GORM Relationship Error: Missing Relationship Field             ║\n")
	sb.WriteString("╚══════════════════════════════════════════════════════════════════╝\n")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Model: %s\n", e.ModelName))
	sb.WriteString(fmt.Sprintf("Field: BelongsTo(\"%s\", \"%s\", \"%s\")\n",
		e.ExpectedFieldName, e.BelongsToFieldKey, e.RelatedResourceSlug))
	sb.WriteString("\n")
	sb.WriteString("Problem: Struct is missing the relationship field\n")
	sb.WriteString("\n")
	sb.WriteString("Fix: Add this field to your struct:\n")
	sb.WriteString("\n")

	// İlişkili model tipini hesapla
	// Örnek: "cargo_companies" -> "CargoCompany"
	relatedType := strings.TrimSuffix(e.RelatedResourceSlug, "s") // Basit pluralization kaldırma
	relatedType = strcase.ToCamel(relatedType)

	// Package adını hesapla
	// Örnek: "cargo_companies" -> "cargo_company"
	packageName := strings.TrimSuffix(e.RelatedResourceSlug, "s")
	packageName = strings.ReplaceAll(packageName, "-", "_")

	sb.WriteString(fmt.Sprintf("    %s *%s.%s `gorm:\"foreignKey:%s;references:ID\"`\n",
		e.ExpectedFieldName, packageName, relatedType, e.ForeignKey))
	sb.WriteString("\n")
	sb.WriteString("Note: Make sure to import the related package:\n")
	sb.WriteString(fmt.Sprintf("    import \"your-project/resources/%s\"\n", packageName))
	sb.WriteString("\n")
	sb.WriteString("Docs: https://panel.go/docs/relationships#belongs-to\n")
	sb.WriteString("\n")

	return sb.String()
}
