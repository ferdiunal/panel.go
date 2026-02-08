package handler

import (
	"github.com/ferdiunal/panel.go/pkg/openapi"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
)

// OpenAPIHandler, OpenAPI spesifikasyonu ve dokümantasyon endpoint'lerini yönetir.
//
// ## Özellikler
// - OpenAPI 3.0.3 spec JSON endpoint'i
// - Swagger UI arayüzü
// - ReDoc arayüzü
// - RapiDoc arayüzü
//
// ## Kullanım Örneği
//
//	handler := NewOpenAPIHandler(resources, config)
//	app.Get("/api/openapi.json", handler.GetSpec)
//	app.Get("/api/docs", handler.SwaggerUI)
type OpenAPIHandler struct {
	resources map[string]resource.Resource
	config    openapi.OpenAPIConfig
	spec      *openapi.OpenAPISpec
}

// NewOpenAPIHandler, yeni bir OpenAPIHandler oluşturur.
//
// ## Parametreler
// - resources: Resource map'i (slug -> Resource)
// - config: OpenAPI yapılandırması
//
// ## Döndürür
// - OpenAPIHandler pointer'ı
func NewOpenAPIHandler(resources map[string]resource.Resource, config openapi.OpenAPIConfig) *OpenAPIHandler {
	return &OpenAPIHandler{
		resources: resources,
		config:    config,
	}
}

// GetSpec, OpenAPI spesifikasyonunu JSON formatında döner.
//
// ## Endpoint
// GET /api/openapi.json
//
// ## Response
// - 200: OpenAPI spec JSON
// - 500: Internal server error
//
// ## Kullanım Örneği
//
//	app.Get("/api/openapi.json", handler.GetSpec)
func (h *OpenAPIHandler) GetSpec(c *fiber.Ctx) error {
	// Spec'i cache'den al veya oluştur
	if h.spec == nil {
		spec := openapi.NewOpenAPISpec(h.config)

		// Static endpoint'leri ekle
		staticGen := openapi.NewStaticSpecGenerator()
		staticPaths := staticGen.GeneratePaths()
		for path, pathItem := range staticPaths {
			spec.Paths[path] = pathItem
		}

		// Dynamic resource endpoint'lerini ekle
		dynamicGen := openapi.NewDynamicSpecGenerator()
		dynamicPaths := dynamicGen.GenerateResourcePaths(h.resources)
		for path, pathItem := range dynamicPaths {
			spec.Paths[path] = pathItem
		}

		// Resource schema'larını ekle
		dynamicSchemas := dynamicGen.GenerateResourceSchemas(h.resources)
		for name, schema := range dynamicSchemas {
			spec.Components.Schemas[name] = schema
		}

		h.spec = spec
	}

	return c.JSON(h.spec)
}

// SwaggerUI, Swagger UI arayüzünü gösterir.
//
// ## Endpoint
// GET /api/docs
//
// ## Response
// - 200: Swagger UI HTML
//
// ## Kullanım Örneği
//
//	app.Get("/api/docs", handler.SwaggerUI)
func (h *OpenAPIHandler) SwaggerUI(c *fiber.Ctx) error {
	// Spec URL'ini oluştur
	specURL := "/api/openapi.json"
	if h.config.BasePath != "" {
		specURL = h.config.BasePath + specURL
	}

	html := openapi.SwaggerUIHTML(specURL, h.config.Title)
	c.Type("html")
	return c.SendString(html)
}

// ReDocUI, ReDoc arayüzünü gösterir.
//
// ## Endpoint
// GET /api/docs/redoc
//
// ## Response
// - 200: ReDoc HTML
//
// ## Kullanım Örneği
//
//	app.Get("/api/docs/redoc", handler.ReDocUI)
func (h *OpenAPIHandler) ReDocUI(c *fiber.Ctx) error {
	// Spec URL'ini oluştur
	specURL := "/api/openapi.json"
	if h.config.BasePath != "" {
		specURL = h.config.BasePath + specURL
	}

	html := openapi.ReDocHTML(specURL, h.config.Title)
	c.Type("html")
	return c.SendString(html)
}

// RapidocUI, RapiDoc arayüzünü gösterir.
//
// ## Endpoint
// GET /api/docs/rapidoc
//
// ## Response
// - 200: RapiDoc HTML
//
// ## Kullanım Örneği
//
//	app.Get("/api/docs/rapidoc", handler.RapidocUI)
func (h *OpenAPIHandler) RapidocUI(c *fiber.Ctx) error {
	// Spec URL'ini oluştur
	specURL := "/api/openapi.json"
	if h.config.BasePath != "" {
		specURL = h.config.BasePath + specURL
	}

	html := openapi.RapidocHTML(specURL, h.config.Title)
	c.Type("html")
	return c.SendString(html)
}

// RefreshSpec, OpenAPI spec cache'ini temizler ve yeniden oluşturulmasını sağlar.
//
// ## Kullanım Örneği
//
//	handler.RefreshSpec()
func (h *OpenAPIHandler) RefreshSpec() {
	h.spec = nil
}
