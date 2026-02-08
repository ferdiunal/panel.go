package panel

import (
	"github.com/ferdiunal/panel.go/pkg/openapi"
)

// OpenAPI, OpenAPI custom mapping registry'sine erişim sağlar.
//
// Bu metod, kullanıcıların field type'larını OpenAPI schema'ya nasıl map edileceğini
// özelleştirmelerine olanak tanır.
//
// ## Kullanım Örneği
//
//	// Global field type mapping
//	panel.OpenAPI().MapFieldType(fields.TYPE_CUSTOM, func(element core.Element) *openapi.Schema {
//	    return &openapi.Schema{
//	        Type: "string",
//	        Format: "custom-format",
//	    }
//	})
//
//	// Specific field mapping
//	panel.OpenAPI().MapField("users", "avatar", func(element core.Element) *openapi.Schema {
//	    return &openapi.Schema{
//	        Type: "string",
//	        Format: "uri",
//	        Description: "User avatar URL",
//	    }
//	})
//
//	// Resource-level mapping
//	panel.OpenAPI().MapResource("users", func(element core.Element) *openapi.Schema {
//	    // Custom mapping logic for all fields in users resource
//	    return nil // Return nil to use default mapping
//	})
func (p *Panel) OpenAPI() *openapi.CustomMappingRegistry {
	return openapi.GetRegistry()
}

// RefreshOpenAPISpec, OpenAPI spec cache'ini temizler ve yeniden oluşturulmasını sağlar.
//
// Bu metod, custom mapping'ler eklendikten sonra spec'in yeniden oluşturulması için kullanılır.
//
// ## Kullanım Örneği
//
//	// Custom mapping ekle
//	panel.OpenAPI().MapFieldType(fields.TYPE_CUSTOM, customMapper)
//
//	// Spec'i yenile
//	panel.RefreshOpenAPISpec()
func (p *Panel) RefreshOpenAPISpec() {
	if p.openAPIHandler != nil {
		p.openAPIHandler.RefreshSpec()
	}
}
